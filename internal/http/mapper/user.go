package mapper

import (
	"github.com/ssokov/pr-reviewer-service/internal/model/domain"
	"github.com/ssokov/pr-reviewer-service/internal/model/dto"
)

func UserToResponse(user *domain.User) dto.UserResponse {
	return dto.UserResponse{
		UserID:   user.UserID,
		Username: user.Username,
		TeamName: user.TeamName,
		IsActive: user.IsActive,
	}
}

func PullRequestsToShort(prs []domain.PullRequest) []dto.PullRequestShort {
	result := make([]dto.PullRequestShort, 0, len(prs))
	for _, pr := range prs {
		result = append(result, dto.PullRequestShort{
			PullRequestID:   pr.PullRequestID,
			PullRequestName: pr.PullRequestName,
			AuthorID:        pr.AuthorID,
			Status:          string(pr.Status),
		})
	}
	return result
}
