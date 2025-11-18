package postgres

import (
	"context"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/ssokov/pr-reviewer-service/internal/model/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func cleanupPRs(t *testing.T, pool *pgxpool.Pool) {
	ctx := context.Background()
	_, err := pool.Exec(ctx, "TRUNCATE TABLE pr_system.pull_requests CASCADE")
	require.NoError(t, err)
	_, err = pool.Exec(ctx, "TRUNCATE TABLE pr_system.users CASCADE")
	require.NoError(t, err)
	_, err = pool.Exec(ctx, "TRUNCATE TABLE pr_system.teams CASCADE")
	require.NoError(t, err)
}

func TestPRRepo_Create(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	pool := setupTestDB(t)
	prRepo := NewPRRepository(pool)
	userRepo := NewUserRepository(pool)
	teamRepo := NewTeamRepository(pool)
	cleanupPRs(t, pool)

	ctx := context.Background()

	team := &domain.Team{TeamName: "test-team"}
	createdTeam, err := teamRepo.Create(ctx, team)
	require.NoError(t, err)

	author := &domain.User{
		UserID:   "author1",
		Username: "Author",
		TeamID:   createdTeam.ID,
		IsActive: true,
	}
	createdAuthor, err := userRepo.Create(ctx, author)
	require.NoError(t, err)

	reviewer := &domain.User{
		UserID:   "reviewer1",
		Username: "Reviewer",
		TeamID:   createdTeam.ID,
		IsActive: true,
	}
	_, err = userRepo.Create(ctx, reviewer)
	require.NoError(t, err)

	t.Run("create PR with reviewers", func(t *testing.T) {
		pr := &domain.PullRequest{
			PullRequestID:     "pr-001",
			PullRequestName:   "Test PR",
			AuthorID:          createdAuthor.UserID,
			Status:            domain.PRStatusOpen,
			AssignedReviewers: []string{"reviewer1"},
		}

		createdPR, err := prRepo.Create(ctx, pr)
		require.NoError(t, err)
		assert.NotZero(t, createdPR.ID)
		assert.Equal(t, "pr-001", createdPR.PullRequestID)
		assert.Equal(t, domain.PRStatusOpen, createdPR.Status)
		assert.Len(t, createdPR.AssignedReviewers, 1)
	})
}

func TestPRRepo_GetByPRID(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	pool := setupTestDB(t)
	prRepo := NewPRRepository(pool)
	userRepo := NewUserRepository(pool)
	teamRepo := NewTeamRepository(pool)
	cleanupPRs(t, pool)

	ctx := context.Background()

	team := &domain.Team{TeamName: "test-team"}
	createdTeam, err := teamRepo.Create(ctx, team)
	require.NoError(t, err)

	author := &domain.User{
		UserID:   "author2",
		Username: "Author2",
		TeamID:   createdTeam.ID,
		IsActive: true,
	}
	createdAuthor, err := userRepo.Create(ctx, author)
	require.NoError(t, err)

	pr := &domain.PullRequest{
		PullRequestID:   "pr-002",
		PullRequestName: "Get Test PR",
		AuthorID:        createdAuthor.UserID,
		Status:          domain.PRStatusOpen,
	}
	_, err = prRepo.Create(ctx, pr)
	require.NoError(t, err)

	t.Run("get existing PR", func(t *testing.T) {
		foundPR, err := prRepo.GetByPRID(ctx, "pr-002")
		require.NoError(t, err)
		assert.NotNil(t, foundPR)
		assert.Equal(t, "pr-002", foundPR.PullRequestID)
		assert.Equal(t, "Get Test PR", foundPR.PullRequestName)
	})

	t.Run("get non-existing PR", func(t *testing.T) {
		foundPR, err := prRepo.GetByPRID(ctx, "pr-999")
		require.NoError(t, err)
		assert.Nil(t, foundPR)
	})
}

func TestPRRepo_Update(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	pool := setupTestDB(t)
	prRepo := NewPRRepository(pool)
	userRepo := NewUserRepository(pool)
	teamRepo := NewTeamRepository(pool)
	cleanupPRs(t, pool)

	ctx := context.Background()

	team := &domain.Team{TeamName: "test-team"}
	createdTeam, err := teamRepo.Create(ctx, team)
	require.NoError(t, err)

	author := &domain.User{
		UserID:   "author3",
		Username: "Author3",
		TeamID:   createdTeam.ID,
		IsActive: true,
	}
	createdAuthor, err := userRepo.Create(ctx, author)
	require.NoError(t, err)

	pr := &domain.PullRequest{
		PullRequestID:   "pr-003",
		PullRequestName: "Update Test PR",
		AuthorID:        createdAuthor.UserID,
		Status:          domain.PRStatusOpen,
	}
	createdPR, err := prRepo.Create(ctx, pr)
	require.NoError(t, err)

	t.Run("update PR status to merged", func(t *testing.T) {
		now := time.Now()
		createdPR.Status = domain.PRStatusMerged
		createdPR.MergedAt = &now

		updatedPR, err := prRepo.Update(ctx, createdPR)
		require.NoError(t, err)
		assert.Equal(t, domain.PRStatusMerged, updatedPR.Status)
		assert.NotNil(t, updatedPR.MergedAt)
	})
}

func TestPRRepo_GetByReviewerID(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	pool := setupTestDB(t)
	prRepo := NewPRRepository(pool)
	userRepo := NewUserRepository(pool)
	teamRepo := NewTeamRepository(pool)
	cleanupPRs(t, pool)

	ctx := context.Background()

	team := &domain.Team{TeamName: "test-team"}
	createdTeam, err := teamRepo.Create(ctx, team)
	require.NoError(t, err)

	author := &domain.User{
		UserID:   "author4",
		Username: "Author4",
		TeamID:   createdTeam.ID,
		IsActive: true,
	}
	createdAuthor, err := userRepo.Create(ctx, author)
	require.NoError(t, err)

	reviewer := &domain.User{
		UserID:   "reviewer2",
		Username: "Reviewer2",
		TeamID:   createdTeam.ID,
		IsActive: true,
	}
	_, err = userRepo.Create(ctx, reviewer)
	require.NoError(t, err)

	pr := &domain.PullRequest{
		PullRequestID:     "pr-004",
		PullRequestName:   "Reviewer Test PR",
		AuthorID:          createdAuthor.UserID,
		Status:            domain.PRStatusOpen,
		AssignedReviewers: []string{"reviewer2"},
	}
	_, err = prRepo.Create(ctx, pr)
	require.NoError(t, err)

	t.Run("get PRs by reviewer", func(t *testing.T) {
		prs, err := prRepo.GetByReviewerID(ctx, "reviewer2")
		require.NoError(t, err)
		assert.Len(t, prs, 1)
		assert.Equal(t, "pr-004", prs[0].PullRequestID)
	})

	t.Run("get PRs for non-reviewer", func(t *testing.T) {
		prs, err := prRepo.GetByReviewerID(ctx, "nobody")
		require.NoError(t, err)
		assert.Empty(t, prs)
	})
}
