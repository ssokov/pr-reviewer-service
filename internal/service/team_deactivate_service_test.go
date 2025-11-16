package service

import (
	"context"
	"errors"
	"testing"

	"github.com/ssokov/pr-reviewer-service/internal/model/domain"
	"github.com/stretchr/testify/assert"
	"github.com/vmkteam/embedlog"
)

func TestTeamService_DeactivateTeam(t *testing.T) {
	ctx := context.Background()
	logger := embedlog.NewLogger(false, false)

	t.Run("success - deactivate team with users and PRs", func(t *testing.T) {
		mockTeamRepo := new(MockTeamRepository)
		mockUserRepo := new(MockUserRepository)
		mockPRRepo := new(MockPRRepository)
		service := NewTeamService(mockTeamRepo, mockUserRepo, mockPRRepo, logger)

		team := &domain.Team{
			ID:       1,
			TeamName: "Backend Team",
		}

		deactivatedUsers := []domain.User{
			{UserID: "u1", Username: "Alice", IsActive: false, TeamID: 1},
			{UserID: "u2", Username: "Bob", IsActive: false, TeamID: 1},
		}

		openPRs := []domain.PullRequest{
			{PullRequestID: "pr1", Status: domain.PRStatusOpen, AssignedReviewers: []string{"u1"}},
			{PullRequestID: "pr2", Status: domain.PRStatusOpen, AssignedReviewers: []string{"u2"}},
		}

		mockTeamRepo.On("GetByName", ctx, "Backend Team").Return(team, nil)
		mockUserRepo.On("DeactivateByTeamID", ctx, int64(1)).Return(deactivatedUsers, nil)
		mockPRRepo.On("GetOpenPRsByUserIDs", ctx, []string{"u1", "u2"}).Return(openPRs, nil)

		users, prs, err := service.DeactivateTeam(ctx, "Backend Team")
		assert.NoError(t, err)
		assert.NotNil(t, users)
		assert.NotNil(t, prs)
		assert.Len(t, users, 2)
		assert.Len(t, prs, 2)

		mockTeamRepo.AssertExpectations(t)
		mockUserRepo.AssertExpectations(t)
		mockPRRepo.AssertExpectations(t)
	})

	t.Run("success - no users to deactivate", func(t *testing.T) {
		mockTeamRepo := new(MockTeamRepository)
		mockUserRepo := new(MockUserRepository)
		mockPRRepo := new(MockPRRepository)
		service := NewTeamService(mockTeamRepo, mockUserRepo, mockPRRepo, logger)

		team := &domain.Team{
			ID:       1,
			TeamName: "Backend Team",
		}

		deactivatedUsers := []domain.User{}

		mockTeamRepo.On("GetByName", ctx, "Backend Team").Return(team, nil)
		mockUserRepo.On("DeactivateByTeamID", ctx, int64(1)).Return(deactivatedUsers, nil)

		users, prs, err := service.DeactivateTeam(ctx, "Backend Team")
		assert.NoError(t, err)
		assert.NotNil(t, users)
		assert.NotNil(t, prs)
		assert.Len(t, users, 0)
		assert.Len(t, prs, 0)

		mockTeamRepo.AssertExpectations(t)
		mockUserRepo.AssertExpectations(t)
	})

	t.Run("error - team not found", func(t *testing.T) {
		mockTeamRepo := new(MockTeamRepository)
		mockUserRepo := new(MockUserRepository)
		mockPRRepo := new(MockPRRepository)
		service := NewTeamService(mockTeamRepo, mockUserRepo, mockPRRepo, logger)

		mockTeamRepo.On("GetByName", ctx, "Ghost Team").Return(nil, errors.New("not found"))

		users, prs, err := service.DeactivateTeam(ctx, "Ghost Team")
		assert.Error(t, err)
		assert.Nil(t, users)
		assert.Nil(t, prs)

		mockTeamRepo.AssertExpectations(t)
	})

	t.Run("error - failed to deactivate users", func(t *testing.T) {
		mockTeamRepo := new(MockTeamRepository)
		mockUserRepo := new(MockUserRepository)
		mockPRRepo := new(MockPRRepository)
		service := NewTeamService(mockTeamRepo, mockUserRepo, mockPRRepo, logger)

		team := &domain.Team{
			ID:       1,
			TeamName: "Backend Team",
		}

		mockTeamRepo.On("GetByName", ctx, "Backend Team").Return(team, nil)
		mockUserRepo.On("DeactivateByTeamID", ctx, int64(1)).Return(nil, errors.New("db error"))

		users, prs, err := service.DeactivateTeam(ctx, "Backend Team")
		assert.Error(t, err)
		assert.Nil(t, users)
		assert.Nil(t, prs)

		mockTeamRepo.AssertExpectations(t)
		mockUserRepo.AssertExpectations(t)
	})

	t.Run("error - failed to get open PRs", func(t *testing.T) {
		mockTeamRepo := new(MockTeamRepository)
		mockUserRepo := new(MockUserRepository)
		mockPRRepo := new(MockPRRepository)
		service := NewTeamService(mockTeamRepo, mockUserRepo, mockPRRepo, logger)

		team := &domain.Team{
			ID:       1,
			TeamName: "Backend Team",
		}

		deactivatedUsers := []domain.User{
			{UserID: "u1", Username: "Alice", IsActive: false, TeamID: 1},
		}

		mockTeamRepo.On("GetByName", ctx, "Backend Team").Return(team, nil)
		mockUserRepo.On("DeactivateByTeamID", ctx, int64(1)).Return(deactivatedUsers, nil)
		mockPRRepo.On("GetOpenPRsByUserIDs", ctx, []string{"u1"}).Return(nil, errors.New("db error"))

		users, prs, err := service.DeactivateTeam(ctx, "Backend Team")
		assert.Error(t, err)
		assert.Nil(t, users)
		assert.Nil(t, prs)

		mockTeamRepo.AssertExpectations(t)
		mockUserRepo.AssertExpectations(t)
		mockPRRepo.AssertExpectations(t)
	})
}
