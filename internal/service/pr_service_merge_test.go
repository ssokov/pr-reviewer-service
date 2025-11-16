package service

import (
	"context"
	"testing"
	"time"

	"github.com/ssokov/pr-reviewer-service/internal/apperror"
	"github.com/ssokov/pr-reviewer-service/internal/model/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/vmkteam/embedlog"
)

func TestPRService_MergePR(t *testing.T) {
	ctx := context.Background()
	logger := embedlog.NewLogger(false, false)

	t.Run("success - first merge", func(t *testing.T) {
		mockPRRepo := new(MockPRRepository)
		mockUserRepo := new(MockUserRepository)
		mockTeamRepo := new(MockTeamRepository)
		service := NewPRService(mockPRRepo, mockUserRepo, mockTeamRepo, logger)

		existingPR := &domain.PullRequest{
			ID:            1,
			PullRequestID: "pr-1",
			Status:        domain.PRStatusOpen,
		}

		mockPRRepo.On("GetByPRID", ctx, "pr-1").Return(existingPR, nil)
		mockPRRepo.On("Update", ctx, mock.Anything).Return(&domain.PullRequest{
			ID:            1,
			PullRequestID: "pr-1",
			Status:        domain.PRStatusMerged,
			MergedAt:      &time.Time{},
		}, nil)

		result, err := service.MergePR(ctx, "pr-1")
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, domain.PRStatusMerged, result.Status)
	})

	t.Run("success - already merged (idempotent)", func(t *testing.T) {
		mockPRRepo := new(MockPRRepository)
		mockUserRepo := new(MockUserRepository)
		mockTeamRepo := new(MockTeamRepository)
		service := NewPRService(mockPRRepo, mockUserRepo, mockTeamRepo, logger)

		now := time.Now()
		existingPR := &domain.PullRequest{
			ID:            1,
			PullRequestID: "pr-1",
			Status:        domain.PRStatusMerged,
			MergedAt:      &now,
		}

		mockPRRepo.On("GetByPRID", ctx, "pr-1").Return(existingPR, nil)

		result, err := service.MergePR(ctx, "pr-1")
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, domain.PRStatusMerged, result.Status)
		assert.NotNil(t, result.MergedAt)
	})

	t.Run("error - PR not found", func(t *testing.T) {
		mockPRRepo := new(MockPRRepository)
		mockUserRepo := new(MockUserRepository)
		mockTeamRepo := new(MockTeamRepository)
		service := NewPRService(mockPRRepo, mockUserRepo, mockTeamRepo, logger)

		mockPRRepo.On("GetByPRID", ctx, "pr-unknown").Return((*domain.PullRequest)(nil), nil)

		result, err := service.MergePR(ctx, "pr-unknown")
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.True(t, apperror.Is(err, apperror.ErrCodePRNotFound))
	})

	t.Run("error - pr_id required", func(t *testing.T) {
		mockPRRepo := new(MockPRRepository)
		mockUserRepo := new(MockUserRepository)
		mockTeamRepo := new(MockTeamRepository)
		service := NewPRService(mockPRRepo, mockUserRepo, mockTeamRepo, logger)

		result, err := service.MergePR(ctx, "")
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.True(t, apperror.Is(err, apperror.ErrCodeInvalidInput))
	})
}

func TestPRService_ReassignReviewer(t *testing.T) {
	ctx := context.Background()
	logger := embedlog.NewLogger(false, false)

	t.Run("success", func(t *testing.T) {
		mockPRRepo := new(MockPRRepository)
		mockUserRepo := new(MockUserRepository)
		mockTeamRepo := new(MockTeamRepository)
		service := NewPRService(mockPRRepo, mockUserRepo, mockTeamRepo, logger)

		existingPR := &domain.PullRequest{
			ID:                1,
			PullRequestID:     "pr-1",
			Status:            domain.PRStatusOpen,
			AssignedReviewers: []string{"u2", "u3"},
		}

		oldUser := &domain.User{UserID: "u2", TeamID: 1, IsActive: true}
		teamMembers := []domain.User{
			{UserID: "u2", IsActive: true},
			{UserID: "u3", IsActive: true},
			{UserID: "u4", IsActive: true},
		}

		mockPRRepo.On("GetByPRID", ctx, "pr-1").Return(existingPR, nil)
		mockUserRepo.On("GetByUserID", ctx, "u2").Return(oldUser, nil)
		mockUserRepo.On("GetByTeamID", ctx, int64(1)).Return(teamMembers, nil)
		mockPRRepo.On("Update", ctx, mock.Anything).Return(&domain.PullRequest{
			ID:                1,
			PullRequestID:     "pr-1",
			Status:            domain.PRStatusOpen,
			AssignedReviewers: []string{"u3", "u4"},
		}, nil)

		result, newReviewer, err := service.ReassignReviewer(ctx, "pr-1", "u2")
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.NotEmpty(t, newReviewer)
	})

	t.Run("error - PR merged", func(t *testing.T) {
		mockPRRepo := new(MockPRRepository)
		mockUserRepo := new(MockUserRepository)
		mockTeamRepo := new(MockTeamRepository)
		service := NewPRService(mockPRRepo, mockUserRepo, mockTeamRepo, logger)

		now := time.Now()
		existingPR := &domain.PullRequest{
			ID:            1,
			PullRequestID: "pr-1",
			Status:        domain.PRStatusMerged,
			MergedAt:      &now,
		}

		mockPRRepo.On("GetByPRID", ctx, "pr-1").Return(existingPR, nil)

		result, newReviewer, err := service.ReassignReviewer(ctx, "pr-1", "u2")
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Empty(t, newReviewer)
	})

	t.Run("error - user not assigned", func(t *testing.T) {
		mockPRRepo := new(MockPRRepository)
		mockUserRepo := new(MockUserRepository)
		mockTeamRepo := new(MockTeamRepository)
		service := NewPRService(mockPRRepo, mockUserRepo, mockTeamRepo, logger)

		existingPR := &domain.PullRequest{
			ID:                1,
			PullRequestID:     "pr-1",
			Status:            domain.PRStatusOpen,
			AssignedReviewers: []string{"u3"},
		}

		mockPRRepo.On("GetByPRID", ctx, "pr-1").Return(existingPR, nil)

		result, newReviewer, err := service.ReassignReviewer(ctx, "pr-1", "u2")
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Empty(t, newReviewer)
		assert.True(t, apperror.Is(err, apperror.ErrCodeNotAssigned))
	})

	t.Run("error - pr_id required", func(t *testing.T) {
		mockPRRepo := new(MockPRRepository)
		mockUserRepo := new(MockUserRepository)
		mockTeamRepo := new(MockTeamRepository)
		service := NewPRService(mockPRRepo, mockUserRepo, mockTeamRepo, logger)

		result, newReviewer, err := service.ReassignReviewer(ctx, "", "u2")
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Empty(t, newReviewer)
		assert.True(t, apperror.Is(err, apperror.ErrCodeInvalidInput))
	})

	t.Run("error - old_user_id required", func(t *testing.T) {
		mockPRRepo := new(MockPRRepository)
		mockUserRepo := new(MockUserRepository)
		mockTeamRepo := new(MockTeamRepository)
		service := NewPRService(mockPRRepo, mockUserRepo, mockTeamRepo, logger)

		result, newReviewer, err := service.ReassignReviewer(ctx, "pr-1", "")
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Empty(t, newReviewer)
		assert.True(t, apperror.Is(err, apperror.ErrCodeInvalidInput))
	})

	t.Run("error - PR not found", func(t *testing.T) {
		mockPRRepo := new(MockPRRepository)
		mockUserRepo := new(MockUserRepository)
		mockTeamRepo := new(MockTeamRepository)
		service := NewPRService(mockPRRepo, mockUserRepo, mockTeamRepo, logger)

		mockPRRepo.On("GetByPRID", ctx, "pr-unknown").Return(nil, nil)

		result, newReviewer, err := service.ReassignReviewer(ctx, "pr-unknown", "u2")
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Empty(t, newReviewer)
		assert.True(t, apperror.Is(err, apperror.ErrCodePRNotFound))
	})

	t.Run("error - old user not found", func(t *testing.T) {
		mockPRRepo := new(MockPRRepository)
		mockUserRepo := new(MockUserRepository)
		mockTeamRepo := new(MockTeamRepository)
		service := NewPRService(mockPRRepo, mockUserRepo, mockTeamRepo, logger)

		existingPR := &domain.PullRequest{
			ID:                1,
			PullRequestID:     "pr-1",
			Status:            domain.PRStatusOpen,
			AssignedReviewers: []string{"u2"},
		}

		mockPRRepo.On("GetByPRID", ctx, "pr-1").Return(existingPR, nil)
		mockUserRepo.On("GetByUserID", ctx, "u2").Return(nil, nil)

		result, newReviewer, err := service.ReassignReviewer(ctx, "pr-1", "u2")
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Empty(t, newReviewer)
		assert.True(t, apperror.Is(err, apperror.ErrCodeUserNotFound))
	})

	t.Run("error - no available reviewers in team", func(t *testing.T) {
		mockPRRepo := new(MockPRRepository)
		mockUserRepo := new(MockUserRepository)
		mockTeamRepo := new(MockTeamRepository)
		service := NewPRService(mockPRRepo, mockUserRepo, mockTeamRepo, logger)

		existingPR := &domain.PullRequest{
			ID:                1,
			PullRequestID:     "pr-1",
			Status:            domain.PRStatusOpen,
			AssignedReviewers: []string{"u2"},
		}

		oldUser := &domain.User{UserID: "u2", TeamID: 1, IsActive: true}
		teamMembers := []domain.User{
			{UserID: "u2", IsActive: true},
		}

		mockPRRepo.On("GetByPRID", ctx, "pr-1").Return(existingPR, nil)
		mockUserRepo.On("GetByUserID", ctx, "u2").Return(oldUser, nil)
		mockUserRepo.On("GetByTeamID", ctx, int64(1)).Return(teamMembers, nil)

		result, newReviewer, err := service.ReassignReviewer(ctx, "pr-1", "u2")
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Empty(t, newReviewer)
		assert.True(t, apperror.Is(err, apperror.ErrCodeInvalidInput))
	})

}
