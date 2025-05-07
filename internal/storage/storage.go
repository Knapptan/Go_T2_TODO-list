package storage

import (
	"TodoList/internal/models"
)

// CRUD interface
type TodoStorage interface {
	CreateTask(task models.Task) error
	GetTask() (models.Task, error)
	UpdateTask(task models.Task) error
	DeleteTask(id int) error
}
