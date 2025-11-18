package max_superuser

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/ssokov/pr-reviewer-service/internal/model/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func cleanupTeams(t *testing.T, pool *pgxpool.Pool) {
	ctx := context.Background()
	_, err := pool.Exec(ctx, "TRUNCATE TABLE pr_system.teams CASCADE")
	require.NoError(t, err)
}

func TestTeamRepo_Create(t *testing.T) {
	pool := setupTestDB(t)
	repo := NewTeamRepository(pool)
	cleanupTeams(t, pool)

	ctx := context.Background()

	t.Run("create team without members", func(t *testing.T) {
		team := &domain.Team{
			TeamName: "Backend Team",
			Members:  []domain.User{},
		}

		createdTeam, err := repo.Create(ctx, team)
		require.NoError(t, err)
		assert.NotZero(t, createdTeam.ID)
		assert.Equal(t, "Backend Team", createdTeam.TeamName)
		assert.NotZero(t, createdTeam.CreatedAt)
	})

	t.Run("create team with duplicate name should fail", func(t *testing.T) {
		team := &domain.Team{
			TeamName: "Backend Team",
			Members:  []domain.User{},
		}

		_, err := repo.Create(ctx, team)
		assert.Error(t, err)
	})
}

func TestTeamRepo_GetByName(t *testing.T) {
	pool := setupTestDB(t)
	repo := NewTeamRepository(pool)
	cleanupTeams(t, pool)

	ctx := context.Background()

	t.Run("get existing team", func(t *testing.T) {
		team := &domain.Team{
			TeamName: "Frontend Team",
			Members:  []domain.User{},
		}
		_, err := repo.Create(ctx, team)
		require.NoError(t, err)

		foundTeam, err := repo.GetByName(ctx, "Frontend Team")
		require.NoError(t, err)
		assert.NotNil(t, foundTeam)
		assert.Equal(t, "Frontend Team", foundTeam.TeamName)
	})

	t.Run("get non-existing team", func(t *testing.T) {
		foundTeam, err := repo.GetByName(ctx, "NonExistent Team")
		require.NoError(t, err)
		assert.Nil(t, foundTeam)
	})
}

func TestTeamRepo_ExistsByName(t *testing.T) {
	pool := setupTestDB(t)
	repo := NewTeamRepository(pool)
	cleanupTeams(t, pool)

	ctx := context.Background()

	t.Run("check existing team", func(t *testing.T) {
		team := &domain.Team{
			TeamName: "DevOps Team",
			Members:  []domain.User{},
		}
		_, err := repo.Create(ctx, team)
		require.NoError(t, err)

		exists, err := repo.ExistsByName(ctx, "DevOps Team")
		require.NoError(t, err)
		assert.True(t, exists)
	})

	t.Run("check non-existing team", func(t *testing.T) {
		exists, err := repo.ExistsByName(ctx, "Ghost Team")
		require.NoError(t, err)
		assert.False(t, exists)
	})
}
