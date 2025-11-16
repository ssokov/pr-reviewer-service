package service

import (
	"context"
	"testing"

	"github.com/ssokov/pr-reviewer-service/internal/apperror"
	"github.com/ssokov/pr-reviewer-service/internal/model/domain"
	"github.com/stretchr/testify/assert"
	"github.com/vmkteam/embedlog"
)

func TestUserService_SetIsActive(t *testing.T) {
	ctx := context.Background()
	logger := embedlog.NewLogger(false, false)

	t.Run("success - set user inactive", func(t *testing.T) {
		mockUserRepo := new(MockUserRepository)
		mockTeamRepo := new(MockTeamRepository)
		service := NewUserService(mockUserRepo, mockTeamRepo, logger)

		expectedUser := &domain.User{
			ID:       1,
			UserID:   "user123",
			Username: "Alice",
			IsActive: false,
		}

		mockUserRepo.On("SetIsActive", ctx, "user123", false).Return(expectedUser, nil)

		result, err := service.SetIsActive(ctx, "user123", false)
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.False(t, result.IsActive)

		mockUserRepo.AssertExpectations(t)
	})

	t.Run("error - user_id required", func(t *testing.T) {
		mockUserRepo := new(MockUserRepository)
		mockTeamRepo := new(MockTeamRepository)
		service := NewUserService(mockUserRepo, mockTeamRepo, logger)

		result, err := service.SetIsActive(ctx, "", false)
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.True(t, apperror.Is(err, apperror.ErrCodeInvalidInput))
	})

	t.Run("error - user not found", func(t *testing.T) {
		mockUserRepo := new(MockUserRepository)
		mockTeamRepo := new(MockTeamRepository)
		service := NewUserService(mockUserRepo, mockTeamRepo, logger)

		mockUserRepo.On("SetIsActive", ctx, "ghost", false).Return(nil, nil)

		result, err := service.SetIsActive(ctx, "ghost", false)
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.True(t, apperror.Is(err, apperror.ErrCodeUserNotFound))

		mockUserRepo.AssertExpectations(t)
	})
}

func TestUserService_GetReview(t *testing.T) {
	ctx := context.Background()
	logger := embedlog.NewLogger(false, false)

	t.Run("success - get reviews", func(t *testing.T) {
		mockUserRepo := new(MockUserRepository)
		mockTeamRepo := new(MockTeamRepository)
		service := NewUserService(mockUserRepo, mockTeamRepo, logger)

		user := &domain.User{
			ID:       1,
			UserID:   "user123",
			Username: "Alice",
			IsActive: true,
		}

		prs := []domain.PullRequest{
			{
				PullRequestID:   "pr1",
				PullRequestName: "Feature A",
				AuthorID:        "user456",
				Status:          domain.PRStatusOpen,
			},
		}

		mockUserRepo.On("GetByUserID", ctx, "user123").Return(user, nil)
		mockUserRepo.On("GetByReviewerID", ctx, "user123").Return(prs, nil)

		result, err := service.GetReview(ctx, "user123")
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result, 1)
		assert.Equal(t, "pr1", result[0].PullRequestID)

		mockUserRepo.AssertExpectations(t)
	})

	t.Run("error - user_id required", func(t *testing.T) {
		mockUserRepo := new(MockUserRepository)
		mockTeamRepo := new(MockTeamRepository)
		service := NewUserService(mockUserRepo, mockTeamRepo, logger)

		result, err := service.GetReview(ctx, "")
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.True(t, apperror.Is(err, apperror.ErrCodeInvalidInput))
	})

	t.Run("error - user not found", func(t *testing.T) {
		mockUserRepo := new(MockUserRepository)
		mockTeamRepo := new(MockTeamRepository)
		service := NewUserService(mockUserRepo, mockTeamRepo, logger)

		mockUserRepo.On("GetByUserID", ctx, "ghost").Return(nil, nil)

		result, err := service.GetReview(ctx, "ghost")
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.True(t, apperror.Is(err, apperror.ErrCodeUserNotFound))

		mockUserRepo.AssertExpectations(t)
	})
}
