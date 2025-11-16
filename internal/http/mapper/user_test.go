package mapper

import (
	"testing"

	"github.com/ssokov/pr-reviewer-service/internal/model/domain"
	"github.com/stretchr/testify/assert"
)

func TestUserToResponse(t *testing.T) {
	user := &domain.User{
		UserID:   "u1",
		Username: "Alice",
		TeamName: "backend",
		IsActive: true,
	}

	result := UserToResponse(user)

	assert.Equal(t, "u1", result.UserID)
	assert.Equal(t, "Alice", result.Username)
	assert.Equal(t, "backend", result.TeamName)
	assert.True(t, result.IsActive)
}
