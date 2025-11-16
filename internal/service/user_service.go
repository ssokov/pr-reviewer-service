package service

import (
	"context"

	"github.com/ssokov/pr-reviewer-service/internal/apperror"
	"github.com/ssokov/pr-reviewer-service/internal/model/domain"
	"github.com/ssokov/pr-reviewer-service/internal/repository"
	"github.com/vmkteam/embedlog"
)

type userService struct {
	userRepo repository.UserRepository
	teamRepo repository.TeamRepository
	logger   embedlog.Logger
}

func NewUserService(userRepo repository.UserRepository, teamRepo repository.TeamRepository, logger embedlog.Logger) UserService {
	return &userService{
		userRepo: userRepo,
		teamRepo: teamRepo,
		logger:   logger,
	}
}

func (s *userService) SetIsActive(ctx context.Context, userID string, isActive bool) (*domain.User, error) {
	if userID == "" {
		return nil, apperror.NewInvalidInputError("user_id is required")
	}

	s.logger.Print(ctx, "setting user active status", "user_id", userID, "is_active", isActive)

	user, err := s.userRepo.SetIsActive(ctx, userID, isActive)
	if err != nil {
		s.logger.Errorf("failed to set user active status: %v", err)
		return nil, apperror.NewInternalError("failed to set user active status", err)
	}

	if user == nil {
		s.logger.Print(ctx, "user not found", "user_id", userID)
		return nil, apperror.NewUserNotFoundError(userID)
	}

	s.logger.Print(ctx, "user active status updated", "user_id", userID, "is_active", isActive)
	return user, nil
}

func (s *userService) GetReview(ctx context.Context, userID string) ([]domain.PullRequest, error) {
	if userID == "" {
		return nil, apperror.NewInvalidInputError("user_id is required")
	}

	s.logger.Print(ctx, "getting reviews for user", "user_id", userID)

	user, err := s.userRepo.GetByUserID(ctx, userID)
	if err != nil {
		s.logger.Errorf("failed to get user: %v", err)
		return nil, apperror.NewInternalError("failed to get user", err)
	}

	if user == nil {
		s.logger.Print(ctx, "user not found", "user_id", userID)
		return nil, apperror.NewUserNotFoundError(userID)
	}

	pullRequests, err := s.userRepo.GetByReviewerID(ctx, userID)
	if err != nil {
		s.logger.Errorf("failed to get reviews: %v", err)
		return nil, apperror.NewInternalError("failed to get reviews", err)
	}

	s.logger.Print(ctx, "reviews retrieved", "user_id", userID, "count", len(pullRequests))
	return pullRequests, nil
}
