package service

import (
	"context"

	"github.com/ssokov/pr-reviewer-service/internal/model/domain"
)

type UserService interface {
	SetIsActive(ctx context.Context, userID string, isActive bool) (*domain.User, error)
	GetReview(ctx context.Context, userID string) ([]domain.PullRequest, error)
}

type PRService interface {
	CreatePR(ctx context.Context, authorID string, pr *domain.PullRequest) (*domain.PullRequest, error)
	MergePR(ctx context.Context, prID string) (*domain.PullRequest, error)
	ReassignReviewer(ctx context.Context, prID string, oldUserID string) (*domain.PullRequest, string, error)
}

type TeamService interface {
	AddTeam(ctx context.Context, team *domain.Team) (*domain.Team, error)
	GetTeam(ctx context.Context, teamName string) (*domain.Team, error)
	DeactivateTeam(ctx context.Context, teamName string) ([]domain.User, []domain.PullRequest, error)
}
