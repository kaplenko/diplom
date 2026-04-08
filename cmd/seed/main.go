// Temporary database seed script.
// Run from the project root:
//
//	go run ./cmd/seed
package main

import (
	"log"
	"os"
	"strings"

	"golang.org/x/crypto/bcrypt"

	"github.com/kaplenko/diplom/internal/config"
	"github.com/kaplenko/diplom/internal/repository"
)

const (
	seedFile     = "cmd/seed/data/seed.sql"
	seedPassword = "password123"
)

func main() {
	log.Println("Starting database seed...")

	// config.Load() requires JWT_SECRET; set a dummy value when it is
	// missing so the seed script can run without a fully configured .env.
	if os.Getenv("JWT_SECRET") == "" {
		os.Setenv("JWT_SECRET", "seed-placeholder")
	}

	// Initialise configuration exactly as in the main application.
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Open database connection (reuses repository.NewPostgresDB).
	db, err := repository.NewPostgresDB(cfg.DB)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()
	log.Println("Connected to database")

	// Generate a bcrypt hash so the seed users can actually log in.
	hash, err := bcrypt.GenerateFromPassword([]byte(seedPassword), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("Failed to generate password hash: %v", err)
	}

	// Read the SQL seed file (path relative to project root).
	raw, err := os.ReadFile(seedFile)
	if err != nil {
		log.Fatalf("Failed to read seed file %s: %v", seedFile, err)
	}

	// Replace the password placeholder with the real hash.
	query := strings.ReplaceAll(string(raw), "{{PASSWORD_HASH}}", string(hash))

	// Execute the entire SQL script in one call.
	if _, err := db.Exec(query); err != nil {
		log.Fatalf("Failed to execute seed SQL: %v", err)
	}

	log.Println("Database seeded successfully!")
	log.Printf("Seed users (password: %s):", seedPassword)
	log.Println("  admin@example.com   — role: admin")
	log.Println("  student@example.com — role: student")
}
