package repository

import (
	"database/sql"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	_ "github.com/lib/pq"

	"github.com/kaplenko/diplom/internal/config"
)

func NewPostgresDB(cfg config.DBConfig) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.DSN())
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(10)

	return db, nil
}

// RunMigrations reads SQL files from the migrations directory and executes them in order.
func RunMigrations(db *sql.DB, migrationsDir string) error {
	var files []string

	err := filepath.WalkDir(migrationsDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && strings.HasSuffix(d.Name(), ".sql") {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("failed to read migrations directory: %w", err)
	}

	sort.Strings(files)

	for _, f := range files {
		content, err := os.ReadFile(f)
		if err != nil {
			return fmt.Errorf("failed to read migration %s: %w", f, err)
		}

		if _, err := db.Exec(string(content)); err != nil {
			return fmt.Errorf("failed to execute migration %s: %w", f, err)
		}

		log.Printf("Applied migration: %s", filepath.Base(f))
	}

	return nil
}
