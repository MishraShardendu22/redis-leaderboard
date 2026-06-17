package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	redisclient "github.com/redis/go-redis/v9"

	"redis-leaderboard/internal/handler"
	appredis "redis-leaderboard/internal/redis"
)

func main() {
	_ = godotenv.Load()

	redisAddr := getEnv("REDIS_ADDR", "localhost:6379")
	serverPort := getEnv("SERVER_PORT", "3000")

	client := redisclient.NewClient(&redisclient.Options{Addr: redisAddr})
	service := appredis.NewService(client)

	if err := service.EnsureSeedData(); err != nil {
		log.Fatalf("seed leaderboard: %v", err)
	}

	app := fiber.New(fiber.Config{
		AppName: "redis-leaderboard",
	})

	app.Static("/static", "./static")

	h := handler.New(service)
	h.RegisterRoutes(app)

	log.Printf("listening on :%s with redis %s", serverPort, redisAddr)
	if err := app.Listen(":" + serverPort); err != nil {
		log.Fatalf("listen: %v", err)
	}
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
