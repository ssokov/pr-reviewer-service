package service

import (
	"context"

	"github.com/ssokov/pr-reviewer-service/internal/model/domain"
	"github.com/stretchr/testify/mock"
)

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(ctx context.Context, user *domain.User) (*domain.User, error) {
	args := m.Called(ctx, user)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserRepository) Update(ctx context.Context, user *domain.User) (*domain.User, error) {
	args := m.Called(ctx, user)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserRepository) GetByUserID(ctx context.Context, userID string) (*domain.User, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserRepository) GetByTeamID(ctx context.Context, teamID int64) ([]domain.User, error) {
	args := m.Called(ctx, teamID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.User), args.Error(1)
}

func (m *MockUserRepository) SetIsActive(ctx context.Context, userID string, isActive bool) (*domain.User, error) {
	args := m.Called(ctx, userID, isActive)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserRepository) GetByReviewerID(ctx context.Context, userID string) ([]domain.PullRequest, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.PullRequest), args.Error(1)
}

func (m *MockUserRepository) DeactivateByTeamID(ctx context.Context, teamID int64) ([]domain.User, error) {
	args := m.Called(ctx, teamID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.User), args.Error(1)
}

type MockPRRepository struct {
	mock.Mock
}

func (m *MockPRRepository) Create(ctx context.Context, pr *domain.PullRequest) (*domain.PullRequest, error) {
	args := m.Called(ctx, pr)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.PullRequest), args.Error(1)
}

func (m *MockPRRepository) Update(ctx context.Context, pr *domain.PullRequest) (*domain.PullRequest, error) {
	args := m.Called(ctx, pr)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.PullRequest), args.Error(1)
}

func (m *MockPRRepository) GetByPRID(ctx context.Context, prID string) (*domain.PullRequest, error) {
	args := m.Called(ctx, prID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.PullRequest), args.Error(1)
}

func (m *MockPRRepository) GetByReviewerID(ctx context.Context, reviewerID string) ([]domain.PullRequest, error) {
	args := m.Called(ctx, reviewerID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.PullRequest), args.Error(1)
}

func (m *MockPRRepository) GetOpenPRsByUserIDs(ctx context.Context, userIDs []string) ([]domain.PullRequest, error) {
	args := m.Called(ctx, userIDs)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.PullRequest), args.Error(1)
}

type MockTeamRepository struct {
	mock.Mock
}

func (m *MockTeamRepository) Create(ctx context.Context, team *domain.Team) (*domain.Team, error) {
	args := m.Called(ctx, team)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Team), args.Error(1)
}

func (m *MockTeamRepository) GetByName(ctx context.Context, teamName string) (*domain.Team, error) {
	args := m.Called(ctx, teamName)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Team), args.Error(1)
}

func (m *MockTeamRepository) ExistsByName(ctx context.Context, teamName string) (bool, error) {
	args := m.Called(ctx, teamName)
	return args.Bool(0), args.Error(1)
}

type MockStatsRepository struct {
	mock.Mock
}

func (m *MockStatsRepository) GetTotalPRs(ctx context.Context) (int, error) {
	args := m.Called(ctx)
	return args.Int(0), args.Error(1)
}

func (m *MockStatsRepository) GetTotalUsers(ctx context.Context) (int, error) {
	args := m.Called(ctx)
	return args.Int(0), args.Error(1)
}

func (m *MockStatsRepository) GetActiveUsers(ctx context.Context) (int, error) {
	args := m.Called(ctx)
	return args.Int(0), args.Error(1)
}

func (m *MockStatsRepository) GetPRsByStatus(ctx context.Context) (map[string]int, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[string]int), args.Error(1)
}

func (m *MockStatsRepository) GetTopReviewers(ctx context.Context, limit int) ([]domain.ReviewerStats, error) {
	args := m.Called(ctx, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.ReviewerStats), args.Error(1)
}
