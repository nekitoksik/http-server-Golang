package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"user-service/internal/config"
	"user-service/internal/db"
	handler "user-service/internal/handlers"
	"user-service/internal/middleware"
	"user-service/internal/repository"
	"user-service/internal/services"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "user-service/docs"
)

// @title           User Service API
// @version         1.0
// @description     API для тестового задания
// @termsOfService  http://swagger.io/terms/

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Config load error: %v", err)
	}

	if err := db.RunMigrations(cfg.Database.URL, cfg.Database.MigrationsPath); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	dbConn, err := db.NewPostgresDB(&cfg.Database)
	if err != nil {
		log.Fatalf("DB connection error: %v", err)
	}

	userRepo := repository.NewPostgresUserRepository(dbConn)
	taskRepo := repository.NewTaskRepository(dbConn)

	jwtServices := services.NewJWTService(cfg.JWT.Secret, int(cfg.JWT.AccessTokenDuration), int(cfg.JWT.RefreshTokenDuration))
	authServices := services.NewAuthService(userRepo, jwtServices)
	userService := services.NewUserService(userRepo, taskRepo)

	authHandler := handler.NewAuthHandler(authServices)
	userHandler := handler.NewUserHandler(userService)
	authMw := middleware.NewAuthMiddleware(jwtServices)

	router := gin.Default()

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// ---- PUBLIC ROUTERS ----
	router.POST("/register", authHandler.Register)
	router.POST("/login", authHandler.Login)

	api := router.Group("/api")
	api.Use(authMw.JWT())
	{
		api.POST("/logout", authHandler.Logout)
		api.GET("/users/:id/status", userHandler.GetStatus)
		api.GET("/users/leaderboard", userHandler.GetLeaderBoard)
		api.POST("/users/:id/task/complete", userHandler.CompleteTask)
		api.POST("/users/:id/referrer", userHandler.AddReferrer)
	}

	srv := &http.Server{
		Addr:         cfg.Server.Address,
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	log.Println("HTTP server initialized")

	go func() {
		log.Printf("Starting HTTP server on %s", cfg.Server.Address)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("HTTP server error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	<-quit
	log.Println("Shutting down gracefully...")
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()

	log.Println("⏸️  Stopping HTTP server...")
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("⚠️  HTTP server shutdown error: %v", err)
	}

	defer func() {
		sqlDB, _ := dbConn.DB()
		if sqlDB != nil {
			sqlDB.Close()
		}
	}()
}
