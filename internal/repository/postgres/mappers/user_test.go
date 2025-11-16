package mappers

import (
	"testing"
	"time"

	"github.com/ssokov/pr-reviewer-service/internal/model/db"
	"github.com/ssokov/pr-reviewer-service/internal/model/domain"
	"github.com/stretchr/testify/assert"
)

func TestUserDBToDomain(t *testing.T) {
	now := time.Now()
	dbUser := &db.User{
		ID:        1,
		UserID:    "user-123",
		Username:  "Alice",
		TeamID:    10,
		IsActive:  true,
		CreatedAt: now,
	}

	teamName := "Backend Team"

	result := UserDBToDomain(dbUser, teamName)

	assert.Equal(t, int64(1), result.ID)
	assert.Equal(t, "user-123", result.UserID)
	assert.Equal(t, "Alice", result.Username)
	assert.Equal(t, int64(10), result.TeamID)
	assert.Equal(t, "Backend Team", result.TeamName)
	assert.True(t, result.IsActive)
	assert.Equal(t, now, result.CreatedAt)
}

func TestUserDBToDomain_Inactive(t *testing.T) {
	now := time.Now()
	dbUser := &db.User{
		ID:        2,
		UserID:    "user-456",
		Username:  "Bob",
		TeamID:    20,
		IsActive:  false,
		CreatedAt: now,
	}

	result := UserDBToDomain(dbUser, "Frontend Team")

	assert.Equal(t, int64(2), result.ID)
	assert.Equal(t, "user-456", result.UserID)
	assert.Equal(t, "Bob", result.Username)
	assert.False(t, result.IsActive)
	assert.Equal(t, "Frontend Team", result.TeamName)
}

func TestUserDBToDomain_EmptyTeamName(t *testing.T) {
	now := time.Now()
	dbUser := &db.User{
		ID:        3,
		UserID:    "user-789",
		Username:  "Charlie",
		TeamID:    0,
		IsActive:  true,
		CreatedAt: now,
	}

	result := UserDBToDomain(dbUser, "")

	assert.Equal(t, "", result.TeamName)
	assert.Equal(t, int64(0), result.TeamID)
}

func TestUserDomainToDB(t *testing.T) {
	now := time.Now()
	domainUser := &domain.User{
		ID:        1,
		UserID:    "user-123",
		Username:  "Alice",
		TeamID:    10,
		TeamName:  "Backend Team",
		IsActive:  true,
		CreatedAt: now,
	}

	result := UserDomainToDB(domainUser)

	assert.Equal(t, int64(1), result.ID)
	assert.Equal(t, "user-123", result.UserID)
	assert.Equal(t, "Alice", result.Username)
	assert.Equal(t, int64(10), result.TeamID)
	assert.True(t, result.IsActive)
	assert.Equal(t, now, result.CreatedAt)
}

func TestUserDomainToDB_Inactive(t *testing.T) {
	now := time.Now()
	domainUser := &domain.User{
		ID:        2,
		UserID:    "user-456",
		Username:  "Bob",
		TeamID:    20,
		IsActive:  false,
		CreatedAt: now,
	}

	result := UserDomainToDB(domainUser)

	assert.Equal(t, int64(2), result.ID)
	assert.False(t, result.IsActive)
}
