package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/AKA333/URLShortner/internal/config"
	"github.com/AKA333/URLShortner/internal/handlers"
	"github.com/AKA333/URLShortner/internal/repository"
	"github.com/AKA333/URLShortner/internal/service"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found")
	}

	cfg := config.LoadConfig()

	postgresRepo, err := repository.NewPostgresRepository(
		fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName),
	)

	if err != nil {
		log.Fatalf("Failed to connect to PostgresSQL %v", err)
	}
	if postgresRepo != nil {
		defer postgresRepo.Close()
	}

	redisRepo, err := repository.NewRedisRepository(cfg.RedisURL)
	if err != nil {
		log.Printf("Failed to connect to Redis: %v", err)
		redisRepo = nil
	}
	if redisRepo != nil {
		defer redisRepo.Close()
	}

	urlService := service.NewURLService(postgresRepo, redisRepo, cfg.BaseURL)

	urlHandler := handlers.NewURLHandler(urlService)

	app := fiber.New(fiber.Config{
		AppName:      "URL Shortener",
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	})

	app.Use(recover.New())
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,OPTIONS",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))

	api := app.Group("/api/v1")

	api.Get("/health", urlHandler.HealthCheck)

	api.Post("/shorten", urlHandler.ShortenURL)
	api.Get("/:shortCode", urlHandler.RedirectURL)
	api.Get("/stats/:shortCode", urlHandler.GetURLStats)

	go func() {
		if err := app.Listen(":" + cfg.ServerPort); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	log.Printf("Server started on port %s", cfg.ServerPort)
	log.Printf("Base URL: %s", cfg.BaseURL)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	if err := app.Shutdown(); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited gracefully")

}
