package repository

import (
	"context"
	"errors"
	"fmt"
	"user-service/internal/domain"

	"gorm.io/gorm"
)

type UserRepository interface {
	Create(ctx context.Context, user *domain.User) (int, error)
	FindByUsername(ctx context.Context, username string) (*domain.User, error)
	GetUserById(ctx context.Context, id int) (*domain.User, error)
	UpdateBalance(ctx context.Context, userID int, newBalance int) error
	GetTopUsersByBalance(ctx context.Context, limit int) ([]domain.User, error)
	AddReferrer(ctx context.Context, userID, referrerID int) error

	SaveRefreshToken(ctx context.Context, userID int, token string) error
	FindByRefreshToken(ctx context.Context, token string) (*domain.User, error)
	RevokeRefreshToken(ctx context.Context, userID int) error
}

type PostgresUserRepository struct {
	db *gorm.DB
}

func NewPostgresUserRepository(db *gorm.DB) *PostgresUserRepository {
	return &PostgresUserRepository{
		db: db,
	}
}

func (r *PostgresUserRepository) Create(ctx context.Context, user *domain.User) (int, error) {
	result := r.db.WithContext(ctx).Create(user)
	if result.Error != nil {
		return 0, fmt.Errorf("failed to create new user: %w", result.Error)
	}
	return int(user.ID), nil
}

func (r *PostgresUserRepository) FindByUsername(ctx context.Context, username string) (*domain.User, error) {
	var user domain.User

	result := r.db.WithContext(ctx).Where("username = ?", username).First(&user)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	if result.Error != nil {
		return nil, fmt.Errorf("error to get user by username: %w", result.Error)
	}

	return &user, nil
}

func (r *PostgresUserRepository) GetUserById(ctx context.Context, id int) (*domain.User, error) {
	var user domain.User

	result := r.db.WithContext(ctx).First(&user, id)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, errors.New("user is not found")
	}

	if result.Error != nil {
		return nil, fmt.Errorf("error to get user by id: %w", result.Error)
	}

	return &user, nil
}

func (r *PostgresUserRepository) SaveRefreshToken(ctx context.Context, userID int, token string) error {
	result := r.db.WithContext(ctx).Model(&domain.User{}).
		Where("id = ?", userID).
		Update("refresh_token", token)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("user is not found")
	}

	return nil
}

func (r *PostgresUserRepository) FindByRefreshToken(ctx context.Context, token string) (*domain.User, error) {
	var user domain.User

	result := r.db.WithContext(ctx).Where("refresh_token = ?", token).First(&user)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, errors.New("user by token not found, may be token is invalid")
	}

	if result.Error != nil {
		return nil, fmt.Errorf("error to get user by id: %w", result.Error)
	}

	return &user, nil
}

func (r *PostgresUserRepository) RevokeRefreshToken(ctx context.Context, userID int) error {
	result := r.db.WithContext(ctx).Model(&domain.User{}).
		Where("id = ?", userID).
		Update("refresh_token", nil)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (r *PostgresUserRepository) UpdateBalance(ctx context.Context, userID int, newBalance int) error {
	result := r.db.WithContext(ctx).Model(&domain.User{}).
		Where("id = ?", userID).
		Update("balance", newBalance)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("user is not found")
	}

	return nil
}

func (r *PostgresUserRepository) GetTopUsersByBalance(ctx context.Context, limit int) ([]domain.User, error) {
	var users []domain.User

	result := r.db.WithContext(ctx).
		Order("balance DESC").
		Limit(limit).
		Find(&users)

	if result.Error != nil {
		return nil, result.Error
	}

	return users, nil
}

func (r *PostgresUserRepository) AddReferrer(ctx context.Context, userID, referrerID int) error {
	var user domain.User
	if err := r.db.WithContext(ctx).First(&user, userID).Error; err != nil {
		return err
	}

	if user.ReferrerID != nil {
		return errors.New("the user already has a referrer")
	}

	if userID == referrerID {
		return errors.New("you cannot be a referal for yourself")
	}

	var referrer domain.User
	if err := r.db.WithContext(ctx).First(&referrer, referrerID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("referrer not found")
		}
		return err
	}

	result := r.db.WithContext(ctx).Model(&domain.User{}).
		Where("id = ?", userID).
		Update("referrer_id", referrerID)

	return result.Error
}
