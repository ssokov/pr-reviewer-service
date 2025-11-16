package mappers

import (
	"github.com/ssokov/pr-reviewer-service/internal/model/db"
	"github.com/ssokov/pr-reviewer-service/internal/model/domain"
)

func UserDBToDomain(dbUser *db.User, teamName string) *domain.User {
	return &domain.User{
		ID:        dbUser.ID,
		UserID:    dbUser.UserID,
		Username:  dbUser.Username,
		TeamID:    dbUser.TeamID,
		TeamName:  teamName,
		IsActive:  dbUser.IsActive,
		CreatedAt: dbUser.CreatedAt,
	}
}

func UserDomainToDB(domainUser *domain.User) *db.User {
	return &db.User{
		ID:        domainUser.ID,
		UserID:    domainUser.UserID,
		Username:  domainUser.Username,
		TeamID:    domainUser.TeamID,
		IsActive:  domainUser.IsActive,
		CreatedAt: domainUser.CreatedAt,
	}
}
