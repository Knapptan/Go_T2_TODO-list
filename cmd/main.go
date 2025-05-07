package main

import (
	"context"
	"fmt"

	"TodoList/internal/config"
	"TodoList/internal/models"
	"TodoList/internal/storage"

	"go.uber.org/zap"
)

// @title Todo List API
// @version 1.0
// @description REST API для управления задачами (TODO-лист)

// @contact.name API Support

// @host localhost:8080
// @BasePath /api/v1

func main() {
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

	repo.CreateTask(models.Task{
		Title:       "Example Title1",
		Description: "Example Description1",
		Status:      "new",
	})

	fmt.Println(repo.GetAllTasks())

}
