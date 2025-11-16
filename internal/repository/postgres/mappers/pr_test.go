package mappers

import (
	"testing"
	"time"

	"github.com/ssokov/pr-reviewer-service/internal/model/db"
	"github.com/ssokov/pr-reviewer-service/internal/model/domain"
	"github.com/stretchr/testify/assert"
)

func TestPRDBToDomain(t *testing.T) {
	now := time.Now()
	dbPR := &db.PullRequest{
		ID:              1,
		PullRequestID:   "pr-123",
		PullRequestName: "Feature X",
		AuthorID:        100,
		StatusID:        1,
		CreatedAt:       now,
		MergedAt:        nil,
	}

	authorUserID := "user-456"
	status := domain.PRStatusOpen
	reviewers := []string{"reviewer1", "reviewer2"}

	result := PRDBToDomain(dbPR, authorUserID, status, reviewers)

	assert.Equal(t, int64(1), result.ID)
	assert.Equal(t, "pr-123", result.PullRequestID)
	assert.Equal(t, "Feature X", result.PullRequestName)
	assert.Equal(t, "user-456", result.AuthorID)
	assert.Equal(t, domain.PRStatusOpen, result.Status)
	assert.Equal(t, []string{"reviewer1", "reviewer2"}, result.AssignedReviewers)
	assert.Equal(t, now, result.CreatedAt)
	assert.Nil(t, result.MergedAt)
}

func TestPRDBToDomain_WithMergedAt(t *testing.T) {
	now := time.Now()
	mergedAt := now.Add(1 * time.Hour)
	dbPR := &db.PullRequest{
		ID:              2,
		PullRequestID:   "pr-456",
		PullRequestName: "Feature Y",
		AuthorID:        200,
		StatusID:        2,
		CreatedAt:       now,
		MergedAt:        &mergedAt,
	}

	result := PRDBToDomain(dbPR, "user-789", domain.PRStatusMerged, []string{})

	assert.Equal(t, int64(2), result.ID)
	assert.Equal(t, domain.PRStatusMerged, result.Status)
	assert.NotNil(t, result.MergedAt)
	assert.Equal(t, mergedAt, *result.MergedAt)
}

func TestPRDomainToDB(t *testing.T) {
	now := time.Now()
	domainPR := &domain.PullRequest{
		ID:                1,
		PullRequestID:     "pr-123",
		PullRequestName:   "Feature X",
		AuthorID:          "user-456",
		Status:            domain.PRStatusOpen,
		AssignedReviewers: []string{"r1", "r2"},
		CreatedAt:         now,
		MergedAt:          nil,
	}

	authorInternalID := int64(100)
	statusID := 1

	result := PRDomainToDB(domainPR, authorInternalID, statusID)

	assert.Equal(t, int64(1), result.ID)
	assert.Equal(t, "pr-123", result.PullRequestID)
	assert.Equal(t, "Feature X", result.PullRequestName)
	assert.Equal(t, int64(100), result.AuthorID)
	assert.Equal(t, 1, result.StatusID)
	assert.Equal(t, now, result.CreatedAt)
	assert.Nil(t, result.MergedAt)
}

func TestPRDomainToDB_WithMergedAt(t *testing.T) {
	now := time.Now()
	mergedAt := now.Add(1 * time.Hour)
	domainPR := &domain.PullRequest{
		ID:              2,
		PullRequestID:   "pr-456",
		PullRequestName: "Feature Y",
		CreatedAt:       now,
		MergedAt:        &mergedAt,
	}

	result := PRDomainToDB(domainPR, 200, 2)

	assert.NotNil(t, result.MergedAt)
	assert.Equal(t, mergedAt, *result.MergedAt)
}

func TestStatusNameToID(t *testing.T) {
	tests := []struct {
		name     string
		status   domain.PRStatus
		expected int
	}{
		{
			name:     "open status",
			status:   domain.PRStatusOpen,
			expected: 1,
		},
		{
			name:     "merged status",
			status:   domain.PRStatusMerged,
			expected: 2,
		},
		{
			name:     "unknown status defaults to open",
			status:   "unknown",
			expected: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := StatusNameToID(tt.status)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestStatusIDToName(t *testing.T) {
	tests := []struct {
		name     string
		statusID int
		expected domain.PRStatus
	}{
		{
			name:     "status ID 1 is open",
			statusID: 1,
			expected: domain.PRStatusOpen,
		},
		{
			name:     "status ID 2 is merged",
			statusID: 2,
			expected: domain.PRStatusMerged,
		},
		{
			name:     "unknown status ID defaults to open",
			statusID: 999,
			expected: domain.PRStatusOpen,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := StatusIDToName(tt.statusID)
			assert.Equal(t, tt.expected, result)
		})
	}
}
