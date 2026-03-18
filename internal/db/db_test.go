package db

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/go-pg/pg/v10"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestDB(t *testing.T) *pg.DB {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	ctx := context.Background()
	connString := os.Getenv("TEST_DATABASE_URL")
	if connString == "" {
		connString = "postgres://db:db@localhost:5433/pr_system?sslmode=disable"
	}

	opt, err := pg.ParseURL(connString)
	if err != nil {
		t.Skipf("Skipping integration test: parse DSN failed: %v", err)
		return nil
	}

	db := pg.Connect(opt)
	if err = db.Ping(ctx); err != nil {
		_ = db.Close()
		t.Skipf("Skipping integration test: cannot ping database: %v", err)
		return nil
	}

	t.Cleanup(func() {
		_ = db.Close()
	})

	return db
}

func cleanupDB(t *testing.T, db *pg.DB) {
	t.Helper()
	ctx := context.Background()
	_, err := db.ExecContext(ctx, "TRUNCATE TABLE pr_system.pull_requests, pr_system.users, pr_system.teams RESTART IDENTITY CASCADE")
	require.NoError(t, err)
}

func ensureStatuses(t *testing.T, db *pg.DB) {
	t.Helper()
	ctx := context.Background()
	_, err := db.ExecContext(ctx, `
		INSERT INTO pr_system.statuses (name) VALUES ('OPEN')
		ON CONFLICT (name) DO NOTHING
	`)
	require.NoError(t, err)
	_, err = db.ExecContext(ctx, `
		INSERT INTO pr_system.statuses (name) VALUES ('MERGED')
		ON CONFLICT (name) DO NOTHING
	`)
	require.NoError(t, err)
}

func insertTeam(t *testing.T, db *pg.DB, name string) int64 {
	t.Helper()
	ctx := context.Background()
	var teamID int64
	_, err := db.QueryOneContext(ctx, pg.Scan(&teamID), `
		INSERT INTO pr_system.teams (name)
		VALUES (?)
		RETURNING id
	`, name)
	require.NoError(t, err)
	return teamID
}

func insertUser(t *testing.T, db *pg.DB, userID, username string, teamID *int64, isActive bool) int64 {
	t.Helper()
	ctx := context.Background()
	var id int64
	_, err := db.QueryOneContext(ctx, pg.Scan(&id), `
		INSERT INTO pr_system.users (user_id, username, is_active, team_id)
		VALUES (?, ?, ?, ?)
		RETURNING id
	`, userID, username, isActive, teamID)
	require.NoError(t, err)
	return id
}

func TestPRRepo_CreateGetUpdateAndReviewerLookup(t *testing.T) {
	db := setupTestDB(t)
	repo := NewPrReviewerServiceRepo(db)
	cleanupDB(t, db)
	ensureStatuses(t, db)

	teamID := insertTeam(t, db, "backend")
	insertUser(t, db, "author-1", "Author", &teamID, true)
	insertUser(t, db, "reviewer-1", "Reviewer 1", &teamID, true)

	ctx := context.Background()
	created, err := repo.CreatePR(ctx, &PullRequest{
		PullRequestID:     "pr-1",
		PullRequestName:   "Test PR",
		AuthorUserID:      "author-1",
		Status:            PRStatusOpen,
		AssignedReviewers: []string{"reviewer-1"},
	})
	require.NoError(t, err)
	require.NotNil(t, created)
	assert.Equal(t, "pr-1", created.PullRequestID)
	assert.Equal(t, PRStatusOpen, created.Status)
	assert.Equal(t, []string{"reviewer-1"}, created.AssignedReviewers)

	got, err := repo.GetByPRID(ctx, "pr-1")
	require.NoError(t, err)
	require.NotNil(t, got)
	assert.Equal(t, "Test PR", got.PullRequestName)
	assert.Equal(t, PRStatusOpen, got.Status)
	assert.Equal(t, []string{"reviewer-1"}, got.AssignedReviewers)

	now := time.Now().UTC()
	updated, err := repo.UpdatePR(ctx, &PullRequest{
		PullRequestID:     "pr-1",
		PullRequestName:   "Test PR merged",
		AuthorUserID:      "author-1",
		Status:            PRStatusMerged,
		AssignedReviewers: []string{"reviewer-1"},
		MergedAt:          &now,
	})
	require.NoError(t, err)
	require.NotNil(t, updated)
	assert.Equal(t, PRStatusMerged, updated.Status)
	require.NotNil(t, updated.MergedAt)

	byReviewer, err := repo.GetByReviewerID(ctx, "reviewer-1")
	require.NoError(t, err)
	require.Len(t, byReviewer, 1)
	assert.Equal(t, "pr-1", byReviewer[0].PullRequestID)
	assert.Equal(t, PRStatusMerged, byReviewer[0].Status)
}

func TestTeamRepo_CreateAndGetByName(t *testing.T) {
	db := setupTestDB(t)
	repo := NewPrReviewerServiceRepo(db)
	cleanupDB(t, db)
	ensureStatuses(t, db)

	insertUser(t, db, "u1", "Alice", nil, true)
	insertUser(t, db, "u2", "Bob", nil, true)

	ctx := context.Background()
	created, err := repo.CreateTeam(ctx, &Team{
		TeamName: "platform",
		Members: []User{
			{UserID: "u1"},
			{UserID: "u2"},
		},
	})
	require.NoError(t, err)
	require.NotNil(t, created)
	assert.Equal(t, "platform", created.TeamName)
	require.Len(t, created.Members, 2)

	got, err := repo.GetByName(ctx, "platform")
	require.NoError(t, err)
	require.NotNil(t, got)
	assert.Equal(t, "platform", got.TeamName)
	require.Len(t, got.Members, 2)

	exists, err := repo.ExistsByName(ctx, "platform")
	require.NoError(t, err)
	assert.True(t, exists)
}

func TestUserRepo_GetByUserIDGetByTeamIDSetIsActive(t *testing.T) {
	db := setupTestDB(t)
	repo := NewPrReviewerServiceRepo(db)
	cleanupDB(t, db)
	ensureStatuses(t, db)

	teamID := insertTeam(t, db, "data")
	insertUser(t, db, "user-1", "Alice", &teamID, true)
	insertUser(t, db, "user-2", "Bob", &teamID, true)

	ctx := context.Background()
	user, err := repo.GetByUserID(ctx, "user-1")
	require.NoError(t, err)
	require.NotNil(t, user)
	assert.Equal(t, "Alice", user.Username)
	assert.Equal(t, "data", user.TeamName)
	assert.True(t, user.IsActive)

	users, err := repo.GetByTeamID(ctx, teamID)
	require.NoError(t, err)
	require.Len(t, users, 2)

	updated, err := repo.SetIsActive(ctx, "user-1", false)
	require.NoError(t, err)
	require.NotNil(t, updated)
	assert.False(t, updated.IsActive)
}
