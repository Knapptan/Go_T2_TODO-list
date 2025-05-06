package models

import "time"

type Task struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description,omitempty"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Config struct {
	Port     string
	DBConfig DBConfig
}

type DBConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DbName   string
}
