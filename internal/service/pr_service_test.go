package service

import (
	"context"
	"testing"

	"github.com/ssokov/pr-reviewer-service/internal/apperror"
	"github.com/ssokov/pr-reviewer-service/internal/model/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/vmkteam/embedlog"
)

func TestPRService_CreatePR(t *testing.T) {
	ctx := context.Background()
	logger := embedlog.NewLogger(false, false)

	t.Run("success", func(t *testing.T) {
		mockPRRepo := new(MockPRRepository)
		mockUserRepo := new(MockUserRepository)
		mockTeamRepo := new(MockTeamRepository)
		service := NewPRService(mockPRRepo, mockUserRepo, mockTeamRepo, logger)

		pr := &domain.PullRequest{
			PullRequestID:   "pr123",
			PullRequestName: "Feature A",
			AuthorID:        "user1",
		}

		author := &domain.User{
			UserID:   "user1",
			IsActive: true,
			TeamID:   1,
		}

		teamMembers := []domain.User{
			{UserID: "user1", IsActive: true},
			{UserID: "user2", IsActive: true},
		}

		mockUserRepo.On("GetByUserID", ctx, "user1").Return(author, nil)
		mockUserRepo.On("GetByTeamID", ctx, int64(1)).Return(teamMembers, nil)
		mockPRRepo.On("Create", ctx, mock.Anything).Return(&domain.PullRequest{
			ID:              1,
			PullRequestID:   "pr123",
			PullRequestName: "Feature A",
			Status:          domain.PRStatusOpen,
		}, nil)

		result, err := service.CreatePR(ctx, "user1", pr)
		assert.NoError(t, err)
		assert.NotNil(t, result)
	})

	t.Run("error - author not active", func(t *testing.T) {
		mockPRRepo := new(MockPRRepository)
		mockUserRepo := new(MockUserRepository)
		mockTeamRepo := new(MockTeamRepository)
		service := NewPRService(mockPRRepo, mockUserRepo, mockTeamRepo, logger)

		pr := &domain.PullRequest{
			PullRequestID:   "pr123",
			PullRequestName: "Feature A",
		}

		author := &domain.User{
			UserID:   "user1",
			IsActive: false,
		}

		mockUserRepo.On("GetByUserID", ctx, "user1").Return(author, nil)

		result, err := service.CreatePR(ctx, "user1", pr)
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.True(t, apperror.Is(err, apperror.ErrCodeInvalidInput))
	})

	t.Run("error - pr_id required", func(t *testing.T) {
		mockPRRepo := new(MockPRRepository)
		mockUserRepo := new(MockUserRepository)
		mockTeamRepo := new(MockTeamRepository)
		service := NewPRService(mockPRRepo, mockUserRepo, mockTeamRepo, logger)

		pr := &domain.PullRequest{
			PullRequestID:   "",
			PullRequestName: "Feature A",
		}

		result, err := service.CreatePR(ctx, "user1", pr)
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.True(t, apperror.Is(err, apperror.ErrCodeInvalidInput))
	})

	t.Run("error - pr_name required", func(t *testing.T) {
		mockPRRepo := new(MockPRRepository)
		mockUserRepo := new(MockUserRepository)
		mockTeamRepo := new(MockTeamRepository)
		service := NewPRService(mockPRRepo, mockUserRepo, mockTeamRepo, logger)

		pr := &domain.PullRequest{
			PullRequestID:   "pr123",
			PullRequestName: "",
		}

		result, err := service.CreatePR(ctx, "user1", pr)
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.True(t, apperror.Is(err, apperror.ErrCodeInvalidInput))
	})

	t.Run("error - author_id required", func(t *testing.T) {
		mockPRRepo := new(MockPRRepository)
		mockUserRepo := new(MockUserRepository)
		mockTeamRepo := new(MockTeamRepository)
		service := NewPRService(mockPRRepo, mockUserRepo, mockTeamRepo, logger)

		pr := &domain.PullRequest{
			PullRequestID:   "pr123",
			PullRequestName: "Feature A",
		}

		result, err := service.CreatePR(ctx, "", pr)
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.True(t, apperror.Is(err, apperror.ErrCodeInvalidInput))
	})

	t.Run("error - author not found", func(t *testing.T) {
		mockPRRepo := new(MockPRRepository)
		mockUserRepo := new(MockUserRepository)
		mockTeamRepo := new(MockTeamRepository)
		service := NewPRService(mockPRRepo, mockUserRepo, mockTeamRepo, logger)

		pr := &domain.PullRequest{
			PullRequestID:   "pr123",
			PullRequestName: "Feature A",
		}

		mockUserRepo.On("GetByUserID", ctx, "user1").Return(nil, nil)

		result, err := service.CreatePR(ctx, "user1", pr)
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.True(t, apperror.Is(err, apperror.ErrCodeUserNotFound))
	})

	t.Run("error - author has no team", func(t *testing.T) {
		mockPRRepo := new(MockPRRepository)
		mockUserRepo := new(MockUserRepository)
		mockTeamRepo := new(MockTeamRepository)
		service := NewPRService(mockPRRepo, mockUserRepo, mockTeamRepo, logger)

		pr := &domain.PullRequest{
			PullRequestID:   "pr123",
			PullRequestName: "Feature A",
		}

		author := &domain.User{
			UserID:   "user1",
			IsActive: true,
			TeamID:   0,
		}

		mockUserRepo.On("GetByUserID", ctx, "user1").Return(author, nil)

		result, err := service.CreatePR(ctx, "user1", pr)
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.True(t, apperror.Is(err, apperror.ErrCodeInvalidInput))
	})

	t.Run("error - no active reviewers in team", func(t *testing.T) {
		mockPRRepo := new(MockPRRepository)
		mockUserRepo := new(MockUserRepository)
		mockTeamRepo := new(MockTeamRepository)
		service := NewPRService(mockPRRepo, mockUserRepo, mockTeamRepo, logger)

		pr := &domain.PullRequest{
			PullRequestID:   "pr123",
			PullRequestName: "Feature A",
		}

		author := &domain.User{
			UserID:   "user1",
			IsActive: true,
			TeamID:   1,
		}

		teamMembers := []domain.User{
			{UserID: "user1", IsActive: true},
		}

		mockUserRepo.On("GetByUserID", ctx, "user1").Return(author, nil)
		mockUserRepo.On("GetByTeamID", ctx, int64(1)).Return(teamMembers, nil)

		result, err := service.CreatePR(ctx, "user1", pr)
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.True(t, apperror.Is(err, apperror.ErrCodeInvalidInput))
	})
}
