package mappers

import (
	"testing"
	"time"

	"github.com/ssokov/pr-reviewer-service/internal/model/db"
	"github.com/ssokov/pr-reviewer-service/internal/model/domain"
	"github.com/stretchr/testify/assert"
)

func TestTeamDBToDomain(t *testing.T) {
	now := time.Now()
	dbTeam := &db.Team{
		ID:        1,
		TeamName:  "Backend Team",
		CreatedAt: now,
	}

	members := []domain.User{
		{UserID: "user1", Username: "Alice", IsActive: true},
		{UserID: "user2", Username: "Bob", IsActive: true},
	}

	result := TeamDBToDomain(dbTeam, members)

	assert.Equal(t, int64(1), result.ID)
	assert.Equal(t, "Backend Team", result.TeamName)
	assert.Len(t, result.Members, 2)
	assert.Equal(t, "user1", result.Members[0].UserID)
	assert.Equal(t, "Alice", result.Members[0].Username)
	assert.Equal(t, "user2", result.Members[1].UserID)
	assert.Equal(t, "Bob", result.Members[1].Username)
	assert.Equal(t, now, result.CreatedAt)
}

func TestTeamDBToDomain_EmptyMembers(t *testing.T) {
	now := time.Now()
	dbTeam := &db.Team{
		ID:        2,
		TeamName:  "Empty Team",
		CreatedAt: now,
	}

	result := TeamDBToDomain(dbTeam, []domain.User{})

	assert.Equal(t, int64(2), result.ID)
	assert.Equal(t, "Empty Team", result.TeamName)
	assert.Empty(t, result.Members)
}

func TestTeamDomainToDB(t *testing.T) {
	now := time.Now()
	domainTeam := &domain.Team{
		ID:       1,
		TeamName: "Backend Team",
		Members: []domain.User{
			{UserID: "user1", Username: "Alice"},
			{UserID: "user2", Username: "Bob"},
		},
		CreatedAt: now,
	}

	result := TeamDomainToDB(domainTeam)

	assert.Equal(t, int64(1), result.ID)
	assert.Equal(t, "Backend Team", result.TeamName)
	assert.Equal(t, now, result.CreatedAt)
}

func TestTeamDomainToDB_EmptyMembers(t *testing.T) {
	now := time.Now()
	domainTeam := &domain.Team{
		ID:        2,
		TeamName:  "Empty Team",
		Members:   []domain.User{},
		CreatedAt: now,
	}

	result := TeamDomainToDB(domainTeam)

	assert.Equal(t, int64(2), result.ID)
	assert.Equal(t, "Empty Team", result.TeamName)
}
