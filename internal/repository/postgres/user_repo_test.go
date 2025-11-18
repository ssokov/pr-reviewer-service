package postgres

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/ssokov/pr-reviewer-service/internal/model/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func cleanupUsers(t *testing.T, pool *pgxpool.Pool) {
	ctx := context.Background()
	_, err := pool.Exec(ctx, "TRUNCATE TABLE pr_system.users CASCADE")
	require.NoError(t, err)
}

func TestUserRepo_Create(t *testing.T) {
	pool := setupTestDB(t)
	repo := NewUserRepository(pool)
	cleanupUsers(t, pool)

	ctx := context.Background()

	t.Run("create user successfully", func(t *testing.T) {
		user := &domain.User{
			UserID:   "user123",
			Username: "John Doe",
			IsActive: true,
		}

		createdUser, err := repo.Create(ctx, user)
		require.NoError(t, err)
		assert.NotZero(t, createdUser.ID)
		assert.Equal(t, "user123", createdUser.UserID)
		assert.Equal(t, "John Doe", createdUser.Username)
		assert.True(t, createdUser.IsActive)
	})

	t.Run("create user with duplicate user_id should fail", func(t *testing.T) {
		user := &domain.User{
			UserID:   "user123",
			Username: "Jane Doe",
			IsActive: true,
		}

		_, err := repo.Create(ctx, user)
		assert.Error(t, err)
	})
}

func TestUserRepo_GetByUserID(t *testing.T) {
	pool := setupTestDB(t)
	repo := NewUserRepository(pool)
	cleanupUsers(t, pool)

	ctx := context.Background()

	t.Run("get existing user", func(t *testing.T) {
		user := &domain.User{
			UserID:   "user456",
			Username: "Alice",
			IsActive: true,
		}
		_, err := repo.Create(ctx, user)
		require.NoError(t, err)

		foundUser, err := repo.GetByUserID(ctx, "user456")
		require.NoError(t, err)
		assert.NotNil(t, foundUser)
		assert.Equal(t, "user456", foundUser.UserID)
		assert.Equal(t, "Alice", foundUser.Username)
	})

	t.Run("get non-existing user", func(t *testing.T) {
		foundUser, err := repo.GetByUserID(ctx, "ghost")
		require.NoError(t, err)
		assert.Nil(t, foundUser)
	})
}

func TestUserRepo_SetIsActive(t *testing.T) {
	pool := setupTestDB(t)
	repo := NewUserRepository(pool)
	cleanupUsers(t, pool)

	ctx := context.Background()

	t.Run("set user inactive", func(t *testing.T) {
		user := &domain.User{
			UserID:   "user789",
			Username: "Bob",
			IsActive: true,
		}
		_, err := repo.Create(ctx, user)
		require.NoError(t, err)

		updatedUser, err := repo.SetIsActive(ctx, "user789", false)
		require.NoError(t, err)
		assert.NotNil(t, updatedUser)
		assert.False(t, updatedUser.IsActive)
	})

	t.Run("set non-existing user inactive", func(t *testing.T) {
		updatedUser, err := repo.SetIsActive(ctx, "ghost", false)
		require.NoError(t, err)
		assert.Nil(t, updatedUser)
	})
}
