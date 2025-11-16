package repository

import (
	"context"

	"github.com/ssokov/pr-reviewer-service/internal/model/domain"
)

type UserRepository interface {
	Create(ctx context.Context, user *domain.User) (*domain.User, error)
	Update(ctx context.Context, user *domain.User) (*domain.User, error)
	GetByUserID(ctx context.Context, userID string) (*domain.User, error)
	GetByTeamID(ctx context.Context, teamID int64) ([]domain.User, error)
	SetIsActive(ctx context.Context, userID string, isActive bool) (*domain.User, error)
	GetByReviewerID(ctx context.Context, userID string) ([]domain.PullRequest, error)
	DeactivateByTeamID(ctx context.Context, teamID int64) ([]domain.User, error)
}

type TeamRepository interface {
	Create(ctx context.Context, team *domain.Team) (*domain.Team, error)
	GetByName(ctx context.Context, teamName string) (*domain.Team, error)
	ExistsByName(ctx context.Context, teamName string) (bool, error)
}

type PRRepository interface {
	Create(ctx context.Context, pr *domain.PullRequest) (*domain.PullRequest, error)
	Update(ctx context.Context, pr *domain.PullRequest) (*domain.PullRequest, error)
	GetByPRID(ctx context.Context, prID string) (*domain.PullRequest, error)
	GetByReviewerID(ctx context.Context, reviewerID string) ([]domain.PullRequest, error)
	GetOpenPRsByUserIDs(ctx context.Context, userIDs []string) ([]domain.PullRequest, error)
}

type StatsRepository interface {
	GetTotalPRs(ctx context.Context) (int, error)
	GetTotalUsers(ctx context.Context) (int, error)
	GetActiveUsers(ctx context.Context) (int, error)
	GetPRsByStatus(ctx context.Context) (map[string]int, error)
	GetTopReviewers(ctx context.Context, limit int) ([]domain.ReviewerStats, error)
}
