package repository

import (
	"context"
	"errors"
	"user-service/internal/domain"

	"gorm.io/gorm"
)

type TaskRepository interface {
	CreateTask(ctx context.Context, task *domain.Task) error
	GetTaskByID(ctx context.Context, id int) (*domain.Task, error)
	CompleteTask(ctx context.Context, userTask *domain.UserTask) error
	GetUserCompletedTasks(ctx context.Context, userID int) ([]domain.Task, error)
}

type PostgresTaskRepository struct {
	db *gorm.DB
}

func NewTaskRepository(db *gorm.DB) *PostgresTaskRepository {
	return &PostgresTaskRepository{
		db: db,
	}
}

func (r *PostgresTaskRepository) CreateTask(ctx context.Context, task *domain.Task) error {
	result := r.db.WithContext(ctx).Create(task)
	return result.Error
}

func (r *PostgresTaskRepository) GetTaskByID(ctx context.Context, id int) (*domain.Task, error) {
	var task domain.Task

	result := r.db.WithContext(ctx).First(&task, id)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, errors.New("task not found")
	}

	return &task, result.Error
}

func (r *PostgresTaskRepository) CompleteTask(ctx context.Context, userTask *domain.UserTask) error {
	result := r.db.WithContext(ctx).Create(userTask)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
			return errors.New("task")
		}
		return result.Error
	}

	return nil
}

func (r *PostgresTaskRepository) GetUserCompletedTasks(ctx context.Context, userID int) ([]domain.Task, error) {
	var tasks []domain.Task

	result := r.db.WithContext(ctx).
		Joins("JOIN user_tasks WHERE user_tasks.task_id = tasks.id").
		Where("user_tasks.user_id = ?", userID).
		Find(&tasks)

	return tasks, result.Error
}
