package main

import (
	"context"

	"TodoList/internal/config"
	"TodoList/internal/handlers"
	"TodoList/internal/storage"

	"github.com/gofiber/fiber/v3"
	"go.uber.org/zap"
)

// @title Todo List API
// @version 1.0
// @description REST API для управления задачами (TODO-лист)

// @host localhost:8080
// @BasePath /api/v1

func main() {
	app := fiber.New()

	logger, _ := zap.NewProduction()
	defer logger.Sync()

	cfg, err := config.Load()
	if err != nil {
		logger.Fatal("Config load error", zap.Error(err))
	}

	pool, err := storage.NewPostgresDB(context.Background(), cfg.DatabaseURL())
	if err != nil {
		logger.Fatal("DB connection error", zap.Error(err))
	}
	defer pool.Close()

	repo := storage.NewPostgresRepository(pool, logger)

	h := handlers.NewHandler(logger, repo)

	app.Get("/tasks", h.GetAllTasks)
	app.Post("/tasks", h.CreateTask)
	app.Get("/tasks/:id", h.GetTask)
	app.Put("/tasks/:id", h.UpdateTask)
	app.Delete("/tasks/:id", h.DeleteTask)

	port := ":" + cfg.Port
	logger.Sugar().Fatalln(app.Listen(port))
}
