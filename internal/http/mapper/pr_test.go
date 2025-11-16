package mapper

import (
	"testing"
	"time"

	"github.com/ssokov/pr-reviewer-service/internal/model/domain"
	"github.com/ssokov/pr-reviewer-service/internal/model/dto"
	"github.com/stretchr/testify/assert"
)

func TestCreatePRRequestToDomain(t *testing.T) {
	req := dto.CreatePRRequest{
		PullRequestID:   "pr-1",
		PullRequestName: "Test PR",
		AuthorID:        "u1",
	}

	result := CreatePRRequestToDomain(req)

	assert.Equal(t, "pr-1", result.PullRequestID)
	assert.Equal(t, "Test PR", result.PullRequestName)
	assert.Equal(t, "u1", result.AuthorID)
	assert.Equal(t, domain.PRStatusOpen, result.Status)
}

func TestPullRequestToResponse(t *testing.T) {
	now := time.Now()
	pr := &domain.PullRequest{
		PullRequestID:     "pr-1",
		PullRequestName:   "Test PR",
		AuthorID:          "u1",
		Status:            domain.PRStatusOpen,
		AssignedReviewers: []string{"u2", "u3"},
		CreatedAt:         now,
		MergedAt:          nil,
	}

	result := PullRequestToResponse(pr)

	assert.Equal(t, "pr-1", result.PullRequestID)
	assert.Equal(t, "Test PR", result.PullRequestName)
	assert.Equal(t, "u1", result.AuthorID)
	assert.Equal(t, "OPEN", result.Status)
	assert.Len(t, result.AssignedReviewers, 2)
	assert.NotNil(t, result.CreatedAt)
	assert.Nil(t, result.MergedAt)
}

func TestPullRequestsToShort(t *testing.T) {
	prs := []domain.PullRequest{
		{
			PullRequestID:   "pr-1",
			PullRequestName: "Test 1",
			AuthorID:        "u1",
			Status:          domain.PRStatusOpen,
		},
		{
			PullRequestID:   "pr-2",
			PullRequestName: "Test 2",
			AuthorID:        "u2",
			Status:          domain.PRStatusMerged,
		},
	}

	result := PullRequestsToShort(prs)

	assert.Len(t, result, 2)
	assert.Equal(t, "pr-1", result[0].PullRequestID)
	assert.Equal(t, "OPEN", result[0].Status)
	assert.Equal(t, "pr-2", result[1].PullRequestID)
	assert.Equal(t, "MERGED", result[1].Status)
}
