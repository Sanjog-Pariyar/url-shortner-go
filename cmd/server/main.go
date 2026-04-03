package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sanjog-pariyar/url-shorten-api/internal/config"
	"github.com/sanjog-pariyar/url-shorten-api/internal/database"
	"github.com/sanjog-pariyar/url-shorten-api/internal/handlers"
	"github.com/sanjog-pariyar/url-shorten-api/internal/middleware"
	"github.com/sanjog-pariyar/url-shorten-api/internal/repository"
	"github.com/sanjog-pariyar/url-shorten-api/internal/services"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Load()

	db, err := database.Connect(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	urlRepo := repository.NewURLRepository(db)
	userRepo := repository.NewUserRepository(db)

	urlService := services.NewURLService(urlRepo, cfg.BaseURL)
	authService := services.NewAuthService(userRepo, cfg.JWTSecret)

	urlHandler := handlers.NewURLHandler(urlService)
	authHandler := handlers.NewAuthHandler(authService)

	rateLimiter := middleware.NewRateLimiter(cfg.RateLimit.RequestsPerMinute)

	router := gin.Default()

	router.Use(middleware.Logger())
	router.Use(middleware.RateLimit(rateLimiter))

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	auth := router.Group("/auth")
	{
		auth.POST("/register", authHandler.Register)
		auth.POST("/login", authHandler.Login)
	}

	router.POST("/shorten", middleware.OptionalAuth(authService), urlHandler.CreateShortURL)

	router.GET("/:shortCode", urlHandler.Redirect)

	authRequired := router.Group("/")
	authRequired.Use(middleware.JWTAuth(authService))
	{
		authRequired.GET("/stats/:shortCode", urlHandler.GetStats)
		authRequired.DELETE("/:shortCode", urlHandler.DeleteURL)
	}

	srv := &http.Server{
		Addr:    ":" + cfg.ServerPort,
		Handler: router,
	}

	go func() {
		log.Printf("Server starting on port %s", cfg.ServerPort)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	fmt.Println("Server exited")
}
