package service

import (
	"context"
	"errors"
	"testing"

	"github.com/ssokov/pr-reviewer-service/internal/apperror"
	"github.com/ssokov/pr-reviewer-service/internal/model/domain"
	"github.com/stretchr/testify/assert"
	"github.com/vmkteam/embedlog"
)

func TestTeamService_AddTeam(t *testing.T) {
	ctx := context.Background()
	logger := embedlog.NewLogger(false, false)

	t.Run("success - create team", func(t *testing.T) {
		mockTeamRepo := new(MockTeamRepository)
		mockUserRepo := new(MockUserRepository)
		mockPRRepo := new(MockPRRepository)
		service := NewTeamService(mockTeamRepo, mockUserRepo, mockPRRepo, logger)

		team := &domain.Team{
			TeamName: "Backend Team",
			Members: []domain.User{
				{UserID: "user1", Username: "Alice"},
			},
		}

		mockTeamRepo.On("ExistsByName", ctx, "Backend Team").Return(false, nil)
		mockTeamRepo.On("Create", ctx, team).Return(&domain.Team{
			ID:       1,
			TeamName: "Backend Team",
			Members:  team.Members,
		}, nil)

		result, err := service.AddTeam(ctx, team)
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, int64(1), result.ID)
		assert.Equal(t, "Backend Team", result.TeamName)

		mockTeamRepo.AssertExpectations(t)
	})

	t.Run("error - team name required", func(t *testing.T) {
		mockTeamRepo := new(MockTeamRepository)
		mockUserRepo := new(MockUserRepository)
		mockPRRepo := new(MockPRRepository)
		service := NewTeamService(mockTeamRepo, mockUserRepo, mockPRRepo, logger)

		team := &domain.Team{
			TeamName: "",
			Members:  []domain.User{{UserID: "user1"}},
		}

		result, err := service.AddTeam(ctx, team)
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.True(t, apperror.Is(err, apperror.ErrCodeInvalidInput))
	})

	t.Run("error - team must have members", func(t *testing.T) {
		mockTeamRepo := new(MockTeamRepository)
		mockUserRepo := new(MockUserRepository)
		mockPRRepo := new(MockPRRepository)
		service := NewTeamService(mockTeamRepo, mockUserRepo, mockPRRepo, logger)

		team := &domain.Team{
			TeamName: "Backend Team",
			Members:  []domain.User{},
		}

		result, err := service.AddTeam(ctx, team)
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.True(t, apperror.Is(err, apperror.ErrCodeInvalidInput))
	})

	t.Run("error - team already exists", func(t *testing.T) {
		mockTeamRepo := new(MockTeamRepository)
		mockUserRepo := new(MockUserRepository)
		mockPRRepo := new(MockPRRepository)
		service := NewTeamService(mockTeamRepo, mockUserRepo, mockPRRepo, logger)

		team := &domain.Team{
			TeamName: "Backend Team",
			Members:  []domain.User{{UserID: "user1"}},
		}

		mockTeamRepo.On("ExistsByName", ctx, "Backend Team").Return(true, nil)

		result, err := service.AddTeam(ctx, team)
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.True(t, apperror.Is(err, apperror.ErrCodeTeamExists))

		mockTeamRepo.AssertExpectations(t)
	})

	t.Run("error - repository error", func(t *testing.T) {
		mockTeamRepo := new(MockTeamRepository)
		mockUserRepo := new(MockUserRepository)
		mockPRRepo := new(MockPRRepository)
		service := NewTeamService(mockTeamRepo, mockUserRepo, mockPRRepo, logger)

		team := &domain.Team{
			TeamName: "Backend Team",
			Members:  []domain.User{{UserID: "user1"}},
		}

		mockTeamRepo.On("ExistsByName", ctx, "Backend Team").Return(false, errors.New("db error"))

		result, err := service.AddTeam(ctx, team)
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.True(t, apperror.Is(err, apperror.ErrCodeInternalError))

		mockTeamRepo.AssertExpectations(t)
	})
}

func TestTeamService_GetTeam(t *testing.T) {
	ctx := context.Background()
	logger := embedlog.NewLogger(false, false)

	t.Run("success - get team", func(t *testing.T) {
		mockTeamRepo := new(MockTeamRepository)
		mockUserRepo := new(MockUserRepository)
		mockPRRepo := new(MockPRRepository)
		service := NewTeamService(mockTeamRepo, mockUserRepo, mockPRRepo, logger)

		expectedTeam := &domain.Team{
			ID:       1,
			TeamName: "Backend Team",
			Members:  []domain.User{{UserID: "user1", Username: "Alice"}},
		}

		mockTeamRepo.On("GetByName", ctx, "Backend Team").Return(expectedTeam, nil)

		result, err := service.GetTeam(ctx, "Backend Team")
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "Backend Team", result.TeamName)

		mockTeamRepo.AssertExpectations(t)
	})

	t.Run("error - team not found", func(t *testing.T) {
		mockTeamRepo := new(MockTeamRepository)
		mockUserRepo := new(MockUserRepository)
		mockPRRepo := new(MockPRRepository)
		service := NewTeamService(mockTeamRepo, mockUserRepo, mockPRRepo, logger)

		mockTeamRepo.On("GetByName", ctx, "Ghost Team").Return(nil, nil)

		result, err := service.GetTeam(ctx, "Ghost Team")
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.True(t, apperror.Is(err, apperror.ErrCodeTeamNotFound))

		mockTeamRepo.AssertExpectations(t)
	})
}
