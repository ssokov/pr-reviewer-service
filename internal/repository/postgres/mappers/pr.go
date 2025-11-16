package mappers

import (
	"github.com/ssokov/pr-reviewer-service/internal/model/db"
	"github.com/ssokov/pr-reviewer-service/internal/model/domain"
)

func PRDBToDomain(dbPR *db.PullRequest, authorUserID string, status domain.PRStatus, reviewers []string) *domain.PullRequest {
	return &domain.PullRequest{
		ID:                dbPR.ID,
		PullRequestID:     dbPR.PullRequestID,
		PullRequestName:   dbPR.PullRequestName,
		AuthorID:          authorUserID,
		Status:            status,
		AssignedReviewers: reviewers,
		CreatedAt:         dbPR.CreatedAt,
		MergedAt:          dbPR.MergedAt,
	}
}

func PRDomainToDB(domainPR *domain.PullRequest, authorInternalID int64, statusID int) *db.PullRequest {
	return &db.PullRequest{
		ID:              domainPR.ID,
		PullRequestID:   domainPR.PullRequestID,
		PullRequestName: domainPR.PullRequestName,
		AuthorID:        authorInternalID,
		StatusID:        statusID,
		CreatedAt:       domainPR.CreatedAt,
		MergedAt:        domainPR.MergedAt,
	}
}

func StatusNameToID(status domain.PRStatus) int {
	switch status {
	case domain.PRStatusOpen:
		return 1
	case domain.PRStatusMerged:
		return 2
	default:
		return 1
	}
}

func StatusIDToName(statusID int) domain.PRStatus {
	switch statusID {
	case 1:
		return domain.PRStatusOpen
	case 2:
		return domain.PRStatusMerged
	default:
		return domain.PRStatusOpen
	}
}
