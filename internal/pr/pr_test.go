package pr

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vmkteam/embedlog"
)

func TestPRService_CreatePR_Validation(t *testing.T) {
	svc := NewPrService(nil, embedlog.NewLogger(false, false))
	ctx := context.Background()

	_, err := svc.CreatePR(ctx, "u1", &PullRequest{PullRequestName: "name"})
	require.Error(t, err)
	assert.True(t, Is(err, ErrCodeInvalidInput))

	_, err = svc.CreatePR(ctx, "u1", &PullRequest{PullRequestID: "pr-1"})
	require.Error(t, err)
	assert.True(t, Is(err, ErrCodeInvalidInput))

	_, err = svc.CreatePR(ctx, "", &PullRequest{PullRequestID: "pr-1", PullRequestName: "name"})
	require.Error(t, err)
	assert.True(t, Is(err, ErrCodeInvalidInput))
}

func TestPRService_MergePR_Validation(t *testing.T) {
	svc := NewPrService(nil, embedlog.NewLogger(false, false))
	_, err := svc.MergePR(context.Background(), "")
	require.Error(t, err)
	assert.True(t, Is(err, ErrCodeInvalidInput))
}

func TestPRService_ReassignReviewer_Validation(t *testing.T) {
	svc := NewPrService(nil, embedlog.NewLogger(false, false))
	ctx := context.Background()

	_, _, err := svc.ReassignReviewer(ctx, "", "u1")
	require.Error(t, err)
	assert.True(t, Is(err, ErrCodeInvalidInput))

	_, _, err = svc.ReassignReviewer(ctx, "pr-1", "")
	require.Error(t, err)
	assert.True(t, Is(err, ErrCodeInvalidInput))
}

func TestPRService_SetIsActive_Validation(t *testing.T) {
	svc := NewPrService(nil, embedlog.NewLogger(false, false))
	_, err := svc.SetIsActive(context.Background(), "", true)
	require.Error(t, err)
	assert.True(t, Is(err, ErrCodeInvalidInput))
}

func TestPRService_GetReview_Validation(t *testing.T) {
	svc := NewPrService(nil, embedlog.NewLogger(false, false))
	_, err := svc.GetReview(context.Background(), "")
	require.Error(t, err)
	assert.True(t, Is(err, ErrCodeInvalidInput))
}

func TestPRService_AddTeam_Validation(t *testing.T) {
	svc := NewPrService(nil, embedlog.NewLogger(false, false))
	ctx := context.Background()

	_, err := svc.AddTeam(ctx, &Team{Members: []User{{UserID: "u1"}}})
	require.Error(t, err)
	assert.True(t, Is(err, ErrCodeInvalidInput))

	_, err = svc.AddTeam(ctx, &Team{TeamName: "backend"})
	require.Error(t, err)
	assert.True(t, Is(err, ErrCodeInvalidInput))
}
