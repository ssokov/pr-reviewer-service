package service

import (
	"context"
	"errors"
	"testing"

	"github.com/ssokov/pr-reviewer-service/internal/model/domain"
	"github.com/stretchr/testify/assert"
	"github.com/vmkteam/embedlog"
)

func TestStatsService_GetStats(t *testing.T) {
	ctx := context.Background()
	logger := embedlog.NewLogger(false, false)

	t.Run("success - full stats", func(t *testing.T) {
		mockStatsRepo := new(MockStatsRepository)
		service := NewStatsService(mockStatsRepo, logger)

		prsByStatus := map[string]int{
			"open":   5,
			"merged": 10,
		}

		topReviewers := []domain.ReviewerStats{
			{UserID: "u1", Username: "Alice", AssignedCount: 10, CompletedCount: 8, ActiveCount: 2},
			{UserID: "u2", Username: "Bob", AssignedCount: 7, CompletedCount: 5, ActiveCount: 2},
		}

		mockStatsRepo.On("GetTotalPRs", ctx).Return(15, nil)
		mockStatsRepo.On("GetTotalUsers", ctx).Return(20, nil)
		mockStatsRepo.On("GetActiveUsers", ctx).Return(18, nil)
		mockStatsRepo.On("GetPRsByStatus", ctx).Return(prsByStatus, nil)
		mockStatsRepo.On("GetTopReviewers", ctx, 10).Return(topReviewers, nil)

		result, err := service.GetStats(ctx)
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, 15, result.TotalPRs)
		assert.Equal(t, 20, result.TotalUsers)
		assert.Equal(t, 18, result.ActiveUsers)
		assert.Len(t, result.PRsByStatus, 2)
		assert.Len(t, result.TopReviewers, 2)

		mockStatsRepo.AssertExpectations(t)
	})

	t.Run("error - failed to get total PRs", func(t *testing.T) {
		mockStatsRepo := new(MockStatsRepository)
		service := NewStatsService(mockStatsRepo, logger)

		mockStatsRepo.On("GetTotalPRs", ctx).Return(0, errors.New("db error"))

		result, err := service.GetStats(ctx)
		assert.Error(t, err)
		assert.Nil(t, result)

		mockStatsRepo.AssertExpectations(t)
	})

	t.Run("error - failed to get total users", func(t *testing.T) {
		mockStatsRepo := new(MockStatsRepository)
		service := NewStatsService(mockStatsRepo, logger)

		mockStatsRepo.On("GetTotalPRs", ctx).Return(15, nil)
		mockStatsRepo.On("GetTotalUsers", ctx).Return(0, errors.New("db error"))

		result, err := service.GetStats(ctx)
		assert.Error(t, err)
		assert.Nil(t, result)

		mockStatsRepo.AssertExpectations(t)
	})

	t.Run("error - failed to get active users", func(t *testing.T) {
		mockStatsRepo := new(MockStatsRepository)
		service := NewStatsService(mockStatsRepo, logger)

		mockStatsRepo.On("GetTotalPRs", ctx).Return(15, nil)
		mockStatsRepo.On("GetTotalUsers", ctx).Return(20, nil)
		mockStatsRepo.On("GetActiveUsers", ctx).Return(0, errors.New("db error"))

		result, err := service.GetStats(ctx)
		assert.Error(t, err)
		assert.Nil(t, result)

		mockStatsRepo.AssertExpectations(t)
	})

	t.Run("error - failed to get PRs by status", func(t *testing.T) {
		mockStatsRepo := new(MockStatsRepository)
		service := NewStatsService(mockStatsRepo, logger)

		mockStatsRepo.On("GetTotalPRs", ctx).Return(15, nil)
		mockStatsRepo.On("GetTotalUsers", ctx).Return(20, nil)
		mockStatsRepo.On("GetActiveUsers", ctx).Return(18, nil)
		mockStatsRepo.On("GetPRsByStatus", ctx).Return(nil, errors.New("db error"))

		result, err := service.GetStats(ctx)
		assert.Error(t, err)
		assert.Nil(t, result)

		mockStatsRepo.AssertExpectations(t)
	})

	t.Run("error - failed to get top reviewers", func(t *testing.T) {
		mockStatsRepo := new(MockStatsRepository)
		service := NewStatsService(mockStatsRepo, logger)

		prsByStatus := map[string]int{"open": 5}

		mockStatsRepo.On("GetTotalPRs", ctx).Return(15, nil)
		mockStatsRepo.On("GetTotalUsers", ctx).Return(20, nil)
		mockStatsRepo.On("GetActiveUsers", ctx).Return(18, nil)
		mockStatsRepo.On("GetPRsByStatus", ctx).Return(prsByStatus, nil)
		mockStatsRepo.On("GetTopReviewers", ctx, 10).Return(nil, errors.New("db error"))

		result, err := service.GetStats(ctx)
		assert.Error(t, err)
		assert.Nil(t, result)

		mockStatsRepo.AssertExpectations(t)
	})

	t.Run("success - empty data", func(t *testing.T) {
		mockStatsRepo := new(MockStatsRepository)
		service := NewStatsService(mockStatsRepo, logger)

		prsByStatus := map[string]int{}
		topReviewers := []domain.ReviewerStats{}

		mockStatsRepo.On("GetTotalPRs", ctx).Return(0, nil)
		mockStatsRepo.On("GetTotalUsers", ctx).Return(0, nil)
		mockStatsRepo.On("GetActiveUsers", ctx).Return(0, nil)
		mockStatsRepo.On("GetPRsByStatus", ctx).Return(prsByStatus, nil)
		mockStatsRepo.On("GetTopReviewers", ctx, 10).Return(topReviewers, nil)

		result, err := service.GetStats(ctx)
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, 0, result.TotalPRs)
		assert.Equal(t, 0, result.TotalUsers)
		assert.Equal(t, 0, result.ActiveUsers)

		mockStatsRepo.AssertExpectations(t)
	})
}
