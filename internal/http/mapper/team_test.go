package mapper

import (
	"testing"
	"time"

	"github.com/ssokov/pr-reviewer-service/internal/model/domain"
	"github.com/ssokov/pr-reviewer-service/internal/model/dto"
	"github.com/stretchr/testify/assert"
)

func TestAddTeamRequestToDomain(t *testing.T) {
	req := dto.AddTeamRequest{
		TeamName: "backend",
		Members: []dto.TeamMember{
			{UserID: "u1", Username: "Alice", IsActive: true},
			{UserID: "u2", Username: "Bob", IsActive: false},
		},
	}

	result := AddTeamRequestToDomain(req)

	assert.Equal(t, "backend", result.TeamName)
	assert.Len(t, result.Members, 2)
	assert.Equal(t, "u1", result.Members[0].UserID)
	assert.Equal(t, "Alice", result.Members[0].Username)
	assert.True(t, result.Members[0].IsActive)
}

func TestTeamToResponse(t *testing.T) {
	team := &domain.Team{
		ID:        1,
		TeamName:  "backend",
		CreatedAt: time.Now(),
		Members: []domain.User{
			{UserID: "u1", Username: "Alice", IsActive: true},
		},
	}

	result := TeamToResponse(team)

	assert.Equal(t, "backend", result.TeamName)
	assert.Len(t, result.Members, 1)
	assert.Equal(t, "u1", result.Members[0].UserID)
}
