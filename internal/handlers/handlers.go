package handlers

import (
	"TodoList/internal/models"
	"TodoList/internal/storage"

	"github.com/gofiber/fiber/v3"
	"go.uber.org/zap"
)

type Handler struct {
	logger  *zap.Logger
	storage storage.TodoStorage
}

func NewHandler(logger *zap.Logger, storage storage.TodoStorage) *Handler {
	return &Handler{
		logger:  logger,
		storage: storage,
	}
}

func (h *Handler) CreateTask(c fiber.Ctx) error {
	var task models.Task

	if err := c.Bind().Body(&task); err != nil {
		h.logger.Error("CreateTask: body parse error", zap.Error(err))
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if task.Status == "" {
		task.Status = "new"
	}

	if err := h.storage.CreateTask(task); err != nil {
		h.logger.Error("CreateTask: storage error", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create task",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(task)
}

func (h *Handler) GetAllTasks(c fiber.Ctx) error {

	tasks, err := h.storage.GetAllTasks()
	if err != nil {
		h.logger.Error("GetAllTasks: storage error", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get all tasks",
		})
	}

	return c.Status(fiber.StatusOK).JSON(tasks)
}
