package handlers

import (
	"TodoList/internal/models"
	"TodoList/internal/storage"
	"errors"
	"strconv"

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

func (h *Handler) GetTask(c fiber.Ctx) error {
	idParam := c.Params("id")
	if idParam == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Task ID is required",
		})
	}

	id, err := strconv.Atoi(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid task ID",
			"id":    idParam,
		})
	}

	task, err := h.storage.GetTask(id)
	if err != nil {

		if errors.Is(err, storage.ErrTaskNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Task not found",
				"id":    id,
			})
		}

		h.logger.Error("GetTask: storage error", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to select task",
		})
	}

	return c.Status(fiber.StatusOK).JSON(task)
}

func (h *Handler) DeleteTask(c fiber.Ctx) error {

	idParam := c.Params("id")
	if idParam == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Task ID is required",
		})
	}

	id, err := strconv.Atoi(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid task ID",
			"id":    idParam,
		})
	}

	if err := h.storage.DeleteTask(id); err != nil {

		if errors.Is(err, storage.ErrTaskNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Task not found",
				"id":    id,
			})
		}

		h.logger.Error("DeleteTask: storage error", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete task",
		})
	}

	return c.SendStatus(fiber.StatusNoContent)
}

func (h *Handler) UpdateTask(c fiber.Ctx) error {
	idParam := c.Params("id")
	if idParam == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Task ID is required",
		})
	}

	id, err := strconv.Atoi(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid task ID",
			"id":    idParam,
		})
	}

	var task models.Task
	if err := c.Bind().Body(&task); err != nil {
		h.logger.Error("UpdateTask: body parse error", zap.Error(err))
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	task.ID = id
	if task.Status == "" {
		task.Status = "new"
	}

	if err := h.storage.UpdateTask(task); err != nil {
		if errors.Is(err, storage.ErrTaskNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Task not found",
			})
		}

		h.logger.Error("UpdateTask: storage error", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update task",
		})
	}

	return c.Status(fiber.StatusOK).JSON(task)
}
