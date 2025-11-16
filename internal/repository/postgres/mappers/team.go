package mappers

import (
	"github.com/ssokov/pr-reviewer-service/internal/model/db"
	"github.com/ssokov/pr-reviewer-service/internal/model/domain"
)

func TeamDBToDomain(dbTeam *db.Team, members []domain.User) *domain.Team {
	return &domain.Team{
		ID:        dbTeam.ID,
		TeamName:  dbTeam.TeamName,
		Members:   members,
		CreatedAt: dbTeam.CreatedAt,
	}
}

func TeamDomainToDB(domainTeam *domain.Team) *db.Team {
	return &db.Team{
		ID:        domainTeam.ID,
		TeamName:  domainTeam.TeamName,
		CreatedAt: domainTeam.CreatedAt,
	}
}
