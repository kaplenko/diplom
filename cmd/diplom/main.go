package main

import (
	"log"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "github.com/kaplenko/diplom/docs"
	"github.com/kaplenko/diplom/internal/config"
	"github.com/kaplenko/diplom/internal/handler"
	"github.com/kaplenko/diplom/internal/middleware"
	"github.com/kaplenko/diplom/internal/repository"
	"github.com/kaplenko/diplom/internal/service"
)

// @title           Diplom API
// @version         1.0
// @description     Online programming learning platform API
// @host            localhost:8080
// @BasePath        /api/v1
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	db, err := repository.NewPostgresDB(cfg.DB)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Repositories
	userRepo := repository.NewUserRepository(db)
	courseRepo := repository.NewCourseRepository(db)
	lessonRepo := repository.NewLessonRepository(db)
	taskRepo := repository.NewTaskRepository(db)
	submissionRepo := repository.NewSubmissionRepository(db)
	progressRepo := repository.NewProgressRepository(db)

	// Services
	authSvc := service.NewAuthService(userRepo, cfg.JWT)
	userSvc := service.NewUserService(userRepo)
	courseSvc := service.NewCourseService(courseRepo)
	lessonSvc := service.NewLessonService(lessonRepo, courseRepo)
	taskSvc := service.NewTaskService(taskRepo, lessonRepo)
	submissionSvc := service.NewSubmissionService(submissionRepo, taskRepo)
	progressSvc := service.NewProgressService(progressRepo)

	// Handlers
	authH := handler.NewAuthHandler(authSvc)
	userH := handler.NewUserHandler(userSvc)
	courseH := handler.NewCourseHandler(courseSvc)
	lessonH := handler.NewLessonHandler(lessonSvc)
	taskH := handler.NewTaskHandler(taskSvc)
	submissionH := handler.NewSubmissionHandler(submissionSvc)
	progressH := handler.NewProgressHandler(progressSvc)

	router := gin.Default()
	router.Use(middleware.RequestLogger())

	// Swagger UI
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	v1 := router.Group("/api/v1")

	// Public routes
	auth := v1.Group("/auth")
	{
		auth.POST("/register", authH.Register)
		auth.POST("/login", authH.Login)
		auth.POST("/refresh", authH.Refresh)
	}

	// Protected routes
	protected := v1.Group("")
	protected.Use(middleware.AuthRequired(authSvc))
	{
		// Users
		protected.GET("/users/me", userH.GetMe)
		protected.PUT("/users/me", userH.UpdateMe)

		// Courses (read)
		protected.GET("/courses", courseH.List)
		protected.GET("/courses/:course_id", courseH.GetByID)

		// Lessons (read)
		protected.GET("/courses/:course_id/lessons", lessonH.ListByCourse)
		protected.GET("/lessons/:lesson_id", lessonH.GetByID)

		// Tasks (read)
		protected.GET("/lessons/:lesson_id/tasks", taskH.ListByLesson)
		protected.GET("/tasks/:task_id", taskH.GetByID)

		// Submissions
		protected.POST("/tasks/:task_id/submissions", submissionH.Create)
		protected.GET("/tasks/:task_id/submissions", submissionH.ListByTask)
		protected.GET("/submissions/:submission_id", submissionH.GetByID)

		// Progress
		protected.GET("/courses/:course_id/progress", progressH.GetCourseProgress)
		protected.GET("/progress", progressH.GetAllProgress)
	}

	// Admin routes
	admin := protected.Group("")
	admin.Use(middleware.AdminOnly())
	{
		// Users (admin)
		admin.GET("/users", userH.ListUsers)
		admin.DELETE("/users/:user_id", userH.DeleteUser)

		// Courses (admin)
		admin.POST("/courses", courseH.Create)
		admin.PUT("/courses/:course_id", courseH.Update)
		admin.DELETE("/courses/:course_id", courseH.Delete)

		// Lessons (admin)
		admin.POST("/courses/:course_id/lessons", lessonH.Create)
		admin.PUT("/lessons/:lesson_id", lessonH.Update)
		admin.DELETE("/lessons/:lesson_id", lessonH.Delete)

		// Tasks (admin)
		admin.POST("/lessons/:lesson_id/tasks", taskH.Create)
		admin.PUT("/tasks/:task_id", taskH.Update)
		admin.DELETE("/tasks/:task_id", taskH.Delete)

		// Submissions (admin)
		admin.GET("/submissions", submissionH.ListAll)
	}

	log.Printf("Server starting on port %s", cfg.ServerPort)
	if err := router.Run(":" + cfg.ServerPort); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
