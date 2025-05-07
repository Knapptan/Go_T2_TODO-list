package storage

import (
	"TodoList/internal/models"
)

// CRUD interface
type TodoStorage interface {
	CreateTask(task models.Task) error
	GetAllTasks() ([]models.Task, error)
	GetTask(id int) (models.Task, error)
	UpdateTask(task models.Task) error
	DeleteTask(id int) error
}
