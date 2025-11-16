package mapper

import (
	"github.com/ssokov/pr-reviewer-service/internal/model/domain"
	"github.com/ssokov/pr-reviewer-service/internal/model/dto"
)

func AddTeamRequestToDomain(req dto.AddTeamRequest) *domain.Team {
	members := make([]domain.User, len(req.Members))
	for i, m := range req.Members {
		members[i] = domain.User{
			UserID:   m.UserID,
			Username: m.Username,
			IsActive: m.IsActive,
		}
	}
	return &domain.Team{
		TeamName: req.TeamName,
		Members:  members,
	}
}

func TeamToResponse(team *domain.Team) dto.TeamResponse {
	members := make([]dto.TeamMember, len(team.Members))
	for i, m := range team.Members {
		members[i] = dto.TeamMember{
			UserID:   m.UserID,
			Username: m.Username,
			IsActive: m.IsActive,
		}
	}
	return dto.TeamResponse{
		TeamName: team.TeamName,
		Members:  members,
	}
}
