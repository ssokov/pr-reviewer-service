package rest

import (
	"testing"
	"time"

	"github.com/ssokov/pr-reviewer-service/internal/pr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreatePRRequestToDomain(t *testing.T) {
	t.Run("given create request when mapped to domain then required fields are preserved", func(t *testing.T) {
		// Arrange
		req := CreatePRRequest{
			PullRequestID:   "pr-1",
			PullRequestName: "Test PR",
			AuthorID:        "u1",
		}

		// Act
		got := CreatePRRequestToDomain(req)

		// Assert
		assert.Equal(t, "pr-1", got.PullRequestID)
		assert.Equal(t, "Test PR", got.PullRequestName)
		assert.Equal(t, "u1", got.AuthorID)
		assert.Equal(t, pr.PRStatusOpen, got.Status)
	})
}

func TestTeamMappings(t *testing.T) {
	t.Run("given add-team request when mapped to domain and response then members are preserved", func(t *testing.T) {
		// Arrange
		req := AddTeamRequest{
			TeamName: "backend",
			Members: []TeamMember{
				{UserID: "u1", Username: "Alice", IsActive: true},
				{UserID: "u2", Username: "Bob", IsActive: false},
			},
		}

		// Act
		domain := AddTeamRequestToDomain(req)
		resp := TeamToResponse(domain)

		// Assert
		require.Len(t, domain.Members, 2)
		assert.Equal(t, "u1", domain.Members[0].UserID)

		assert.Equal(t, "backend", resp.TeamName)
		require.Len(t, resp.Members, 2)
		assert.Equal(t, "u2", resp.Members[1].UserID)
	})
}

func TestPullRequestsToShort(t *testing.T) {
	t.Run("given domain pull requests when mapped to short response then statuses are converted to strings", func(t *testing.T) {
		// Arrange
		now := time.Now()
		in := []pr.PullRequest{
			{PullRequestID: "pr-1", PullRequestName: "one", AuthorID: "u1", Status: pr.PRStatusOpen, CreatedAt: now},
			{PullRequestID: "pr-2", PullRequestName: "two", AuthorID: "u2", Status: pr.PRStatusMerged, CreatedAt: now},
		}

		// Act
		got := PullRequestsToShort(in)

		// Assert
		require.Len(t, got, 2)
		assert.Equal(t, "OPEN", got[0].Status)
		assert.Equal(t, "MERGED", got[1].Status)
	})
}
