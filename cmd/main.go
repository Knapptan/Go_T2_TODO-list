package main

import (
	"TodoList/internal/config"
	"fmt"

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

	fmt.Println(cfg)

}
