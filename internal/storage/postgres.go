package storage

import (
	"context"
	"errors"
	"fmt"
	"time"

	"TodoList/internal/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

var (
	ErrTaskNotFound = errors.New("task not found")
)

type PostgresRepository struct {
	pool   *pgxpool.Pool
	logger *zap.Logger
}

func NewPostgresRepository(pool *pgxpool.Pool, logger *zap.Logger) *PostgresRepository {
	return &PostgresRepository{
		pool:   pool,
		logger: logger.With(zap.String("component", "postgres_repository"))}
}

func NewPostgresDB(ctx context.Context, connString string) (*pgxpool.Pool, error) {
	config, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, fmt.Errorf("failed to parse conn string: %v", err)
	}

	config.MaxConns = 10
	config.MinConns = 2
	config.HealthCheckPeriod = 1 * time.Minute

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("unable to create connection pool: %v", err)
	}

	return pool, nil
}

func (r *PostgresRepository) CreateTask(task models.Task) error {
	r.logger.Debug("CreateTask started",
		zap.String("Title", task.Title),
		zap.Int("ID", task.ID),
	)

	query := `
	INSERT INTO tasks
	(title, description, status)
	VALUES
	(@title, @description, @status)
	RETURNING id, created_at, updated_at`

	args := pgx.NamedArgs{
		"title":       task.Title,
		"description": task.Description,
		"status":      task.Status,
	}

	allowedStatuses := map[string]bool{"new": true, "in_progress": true, "done": true}
	if !allowedStatuses[task.Status] {
		r.logger.Error("task.Status wrong",
			zap.String("query", query),
			zap.Any("arguments", args),
		)
		return fmt.Errorf("invalid status: %s", task.Status)
	}

	row := r.pool.QueryRow(context.TODO(), query, args)
	if err := row.Scan(&task.ID, &task.CreatedAt, &task.UpdatedAt); err != nil {
		r.logger.Error("Failed to create task",
			zap.Error(err),
			zap.String("query", query),
			zap.Any("arguments", args),
		)
		return fmt.Errorf("create person: %w", err)
	}

	r.logger.Info("Task created successfully",
		zap.Int("id", task.ID),
		zap.Time("created_at", task.CreatedAt),
		zap.Time("updated_at", task.UpdatedAt),
	)
	return nil
}

func (r *PostgresRepository) GetAllTasks() ([]models.Task, error) {
	r.logger.Debug("GetAllTasks started")

	query := `SELECT * FROM tasks ORDER BY created_at DESC`

	rows, err := r.pool.Query(context.TODO(), query)
	if err != nil {
		r.logger.Error("GetAllTasks failed", zap.Error(err))
		return nil, fmt.Errorf("get tasks: %w", err)
	}
	defer rows.Close()

	var tasks []models.Task
	for rows.Next() {
		var task models.Task
		err := rows.Scan(
			&task.ID,
			&task.Title,
			&task.Description,
			&task.Status,
			&task.CreatedAt,
			&task.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("scan task: %w", err)
		}
		tasks = append(tasks, task)
	}

	r.logger.Info("Geted all tasks successfully",
		zap.Int("number of pieces", len(tasks)),
	)
	return tasks, nil
}

func (r *PostgresRepository) UpdateTask(task models.Task) error {
	r.logger.Debug("UpdateTask started",
		zap.String("Title", task.Title),
		zap.Int("ID", task.ID),
	)

	query := `
        UPDATE tasks 
        SET 
            title = $1,
            description = $2,
            status = $3,
            updated_at = NOW()
        WHERE id = $4
    `

	result, err := r.pool.Exec(
		context.TODO(),
		query,
		task.Title,
		task.Description,
		task.Status,
		task.ID,
	)

	if err != nil {
		r.logger.Error("UpdateTask failed", zap.Error(err))
		return fmt.Errorf("update task: %w", err)
	}

	if result.RowsAffected() == 0 {
		return ErrTaskNotFound
	}

	r.logger.Info("Task updated successfully",
		zap.Int("id", task.ID),
		zap.Time("created_at", task.CreatedAt),
		zap.Time("updated_at", task.UpdatedAt),
	)
	return nil
}

func (r *PostgresRepository) DeleteTask(id int) error {
	r.logger.Debug("DeleteTask started",
		zap.Int("ID", id),
	)

	query := `DELETE FROM tasks WHERE id = $1`

	result, err := r.pool.Exec(context.Background(), query, id)
	if err != nil {
		if result.RowsAffected() == 0 {
			r.logger.Debug("UpdateTask: no rows affected", zap.Int("id", id))
			return ErrTaskNotFound
		}
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("task with ID %d not found", id)
	}

	r.logger.Info("Task deleted successfully",
		zap.Int("id", id),
	)
	return nil
}

func (r *PostgresRepository) GetTask(id int) (models.Task, error) {
	r.logger.Debug("GetTask started",
		zap.Int("ID", id),
	)

	var task models.Task
	query := `
        SELECT id, title, description, status, created_at, updated_at
        FROM tasks
        WHERE id = $1
    `
	row := r.pool.QueryRow(context.Background(), query, id)
	if err := row.Scan(
		&task.ID,
		&task.Title,
		&task.Description,
		&task.Status,
		&task.CreatedAt,
		&task.UpdatedAt,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			r.logger.Debug("GetTask: no rows", zap.Int("id", id))
			return models.Task{}, ErrTaskNotFound
		}
		r.logger.Error("GetTask failed", zap.Error(err))
		return models.Task{}, fmt.Errorf("select task: %w", err)
	}

	r.logger.Info("Task geted successfully",
		zap.Int("id", task.ID),
		zap.Time("created_at", task.CreatedAt),
		zap.Time("updated_at", task.UpdatedAt),
	)
	return task, nil
}
