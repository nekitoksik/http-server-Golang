package services

import (
	"context"
	"fmt"
	"user-service/internal/domain"
	"user-service/internal/dto"
	"user-service/internal/repository"
)

type UserService struct {
	userRepo repository.UserRepository
	taskRepo repository.TaskRepository
}

func NewUserService(
	userRepo repository.UserRepository,
	taskRepo repository.TaskRepository,
) *UserService {
	return &UserService{
		userRepo: userRepo,
		taskRepo: taskRepo,
	}
}

func (s *UserService) GetUserInfoByID(ctx context.Context, userID int) (*dto.UserStatusResponse, error) {
	user, err := s.userRepo.GetUserById(ctx, userID)
	if err != nil {
		return nil, err
	}

	response := dto.ToUserStatusResponse(user)
	return &response, nil
}

func (s *UserService) GetLeaderBoard(ctx context.Context, limit int) (*[]dto.LeaderboardUserDTO, error) {
	users, err := s.userRepo.GetTopUsersByBalance(ctx, limit)
	if err != nil {
		return nil, err
	}

	userDTOs := make([]dto.LeaderboardUserDTO, len(users))
	for i, user := range users {
		userDTOs[i] = dto.ToLeaderboardUserDTO(&user, i+1)
	}

	return &userDTOs, nil
}

func (s *UserService) CompleteTask(ctx context.Context, userID, taskID int) error {
	task, err := s.taskRepo.GetTaskByID(ctx, taskID)
	if err != nil {
		return fmt.Errorf("task is not found: %w", err)
	}

	userTask := &domain.UserTask{
		UserID: userID,
		TaskID: taskID,
	}

	if err := s.taskRepo.CompleteTask(ctx, userTask); err != nil {
		return err
	}

	user, err := s.userRepo.GetUserById(ctx, userID)
	if err != nil {
		return err
	}

	newBalance := user.Balance + task.Points
	if err := s.userRepo.UpdateBalance(ctx, userID, newBalance); err != nil {
		return err
	}

	return nil
}

func (s *UserService) AddReferrer(ctx context.Context, userID, referrerID int) error {
	if err := s.userRepo.AddReferrer(ctx, userID, referrerID); err != nil {
		return err
	}

	const referalBonus = 100

	referrer, err := s.userRepo.GetUserById(ctx, referrerID)
	if err != nil {
		return err
	}

	newBalance := referrer.Balance + referalBonus
	if err := s.userRepo.UpdateBalance(ctx, referrerID, newBalance); err != nil {
		return err
	}

	return nil
}

func (s *UserService) GetUserCompletedTasks(ctx context.Context, userID int) ([]domain.Task, error) {
	return s.taskRepo.GetUserCompletedTasks(ctx, userID)
}
