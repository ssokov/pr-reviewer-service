package pr

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vmkteam/embedlog"
)

func TestPRService_CreatePR_Validation(t *testing.T) {
	t.Run("given missing pull_request_id when CreatePR then invalid_input is returned", func(t *testing.T) {
		// Arrange
		svc := NewPrService(nil, embedlog.NewLogger(false, false))
		ctx := context.Background()

		// Act
		_, err := svc.CreatePR(ctx, "u1", &PullRequest{PullRequestName: "name"})

		// Assert
		require.Error(t, err)
		assert.True(t, Is(err, ErrCodeInvalidInput))
	})

	t.Run("given missing pull_request_name when CreatePR then invalid_input is returned", func(t *testing.T) {
		// Arrange
		svc := NewPrService(nil, embedlog.NewLogger(false, false))
		ctx := context.Background()

		// Act
		_, err := svc.CreatePR(ctx, "u1", &PullRequest{PullRequestID: "pr-1"})

		// Assert
		require.Error(t, err)
		assert.True(t, Is(err, ErrCodeInvalidInput))
	})

	t.Run("given missing author_id when CreatePR then invalid_input is returned", func(t *testing.T) {
		// Arrange
		svc := NewPrService(nil, embedlog.NewLogger(false, false))
		ctx := context.Background()

		// Act
		_, err := svc.CreatePR(ctx, "", &PullRequest{PullRequestID: "pr-1", PullRequestName: "name"})

		// Assert
		require.Error(t, err)
		assert.True(t, Is(err, ErrCodeInvalidInput))
	})
}

func TestPRService_MergePR_Validation(t *testing.T) {
	t.Run("given empty pull_request_id when MergePR then invalid_input is returned", func(t *testing.T) {
		// Arrange
		svc := NewPrService(nil, embedlog.NewLogger(false, false))
		ctx := context.Background()

		// Act
		_, err := svc.MergePR(ctx, "")

		// Assert
		require.Error(t, err)
		assert.True(t, Is(err, ErrCodeInvalidInput))
	})
}

func TestPRService_ReassignReviewer_Validation(t *testing.T) {
	t.Run("given empty pull_request_id when ReassignReviewer then invalid_input is returned", func(t *testing.T) {
		// Arrange
		svc := NewPrService(nil, embedlog.NewLogger(false, false))
		ctx := context.Background()

		// Act
		_, _, err := svc.ReassignReviewer(ctx, "", "u1")

		// Assert
		require.Error(t, err)
		assert.True(t, Is(err, ErrCodeInvalidInput))
	})

	t.Run("given empty old_user_id when ReassignReviewer then invalid_input is returned", func(t *testing.T) {
		// Arrange
		svc := NewPrService(nil, embedlog.NewLogger(false, false))
		ctx := context.Background()

		// Act
		_, _, err := svc.ReassignReviewer(ctx, "pr-1", "")

		// Assert
		require.Error(t, err)
		assert.True(t, Is(err, ErrCodeInvalidInput))
	})
}

func TestPRService_SetIsActive_Validation(t *testing.T) {
	t.Run("given empty user_id when SetIsActive then invalid_input is returned", func(t *testing.T) {
		// Arrange
		svc := NewPrService(nil, embedlog.NewLogger(false, false))
		ctx := context.Background()

		// Act
		_, err := svc.SetIsActive(ctx, "", true)

		// Assert
		require.Error(t, err)
		assert.True(t, Is(err, ErrCodeInvalidInput))
	})
}

func TestPRService_GetReview_Validation(t *testing.T) {
	t.Run("given empty user_id when GetReview then invalid_input is returned", func(t *testing.T) {
		// Arrange
		svc := NewPrService(nil, embedlog.NewLogger(false, false))
		ctx := context.Background()

		// Act
		_, err := svc.GetReview(ctx, "")

		// Assert
		require.Error(t, err)
		assert.True(t, Is(err, ErrCodeInvalidInput))
	})
}

func TestPRService_AddTeam_Validation(t *testing.T) {
	t.Run("given missing team_name when AddTeam then invalid_input is returned", func(t *testing.T) {
		// Arrange
		svc := NewPrService(nil, embedlog.NewLogger(false, false))
		ctx := context.Background()

		// Act
		_, err := svc.AddTeam(ctx, &Team{Members: []User{{UserID: "u1"}}})

		// Assert
		require.Error(t, err)
		assert.True(t, Is(err, ErrCodeInvalidInput))
	})

	t.Run("given empty members when AddTeam then invalid_input is returned", func(t *testing.T) {
		// Arrange
		svc := NewPrService(nil, embedlog.NewLogger(false, false))
		ctx := context.Background()

		// Act
		_, err := svc.AddTeam(ctx, &Team{TeamName: "backend"})

		// Assert
		require.Error(t, err)
		assert.True(t, Is(err, ErrCodeInvalidInput))
	})
}
