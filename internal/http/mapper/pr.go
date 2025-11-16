package mapper

import (
	"github.com/ssokov/pr-reviewer-service/internal/model/domain"
	"github.com/ssokov/pr-reviewer-service/internal/model/dto"
)

func CreatePRRequestToDomain(req dto.CreatePRRequest) *domain.PullRequest {
	return &domain.PullRequest{
		PullRequestID:   req.PullRequestID,
		PullRequestName: req.PullRequestName,
		AuthorID:        req.AuthorID,
		Status:          domain.PRStatusOpen,
	}
}

func PullRequestToResponse(pr *domain.PullRequest) dto.PullRequestResponse {
	return dto.PullRequestResponse{
		PullRequestID:     pr.PullRequestID,
		PullRequestName:   pr.PullRequestName,
		AuthorID:          pr.AuthorID,
		Status:            string(pr.Status),
		AssignedReviewers: pr.AssignedReviewers,
		CreatedAt:         &pr.CreatedAt,
		MergedAt:          pr.MergedAt,
	}
}
