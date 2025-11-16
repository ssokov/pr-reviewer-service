package service

import (
	"context"

	"github.com/ssokov/pr-reviewer-service/internal/apperror"
	"github.com/ssokov/pr-reviewer-service/internal/model/domain"
	"github.com/ssokov/pr-reviewer-service/internal/repository"
	"github.com/vmkteam/embedlog"
)

type teamService struct {
	teamRepo repository.TeamRepository
	userRepo repository.UserRepository
	prRepo   repository.PRRepository
	logger   embedlog.Logger
}

func NewTeamService(teamRepo repository.TeamRepository, userRepo repository.UserRepository, prRepo repository.PRRepository, logger embedlog.Logger) TeamService {
	return &teamService{
		teamRepo: teamRepo,
		userRepo: userRepo,
		prRepo:   prRepo,
		logger:   logger,
	}
}

func (s *teamService) AddTeam(ctx context.Context, team *domain.Team) (*domain.Team, error) {
	if team.TeamName == "" {
		return nil, apperror.NewInvalidInputError("team_name is required")
	}

	if len(team.Members) == 0 {
		return nil, apperror.NewInvalidInputError("team must have at least one member")
	}

	s.logger.Print(ctx, "creating team", "team_name", team.TeamName, "members_count", len(team.Members))

	exists, err := s.teamRepo.ExistsByName(ctx, team.TeamName)
	if err != nil {
		s.logger.Errorf("failed to check team existence: %v", err)
		return nil, apperror.NewInternalError("failed to check team existence", err)
	}
	if exists {
		s.logger.Print(ctx, "team already exists", "team_name", team.TeamName)
		return nil, apperror.NewTeamExistsError(team.TeamName)
	}

	createdTeam, err := s.teamRepo.Create(ctx, team)
	if err != nil {
		s.logger.Errorf("failed to create team: %v", err)
		return nil, apperror.NewInternalError("failed to create team", err)
	}

	s.logger.Print(ctx, "team created successfully", "team_name", createdTeam.TeamName, "team_id", createdTeam.ID)
	return createdTeam, nil
}

func (s *teamService) GetTeam(ctx context.Context, teamName string) (*domain.Team, error) {
	s.logger.Print(ctx, "getting team", "team_name", teamName)

	team, err := s.teamRepo.GetByName(ctx, teamName)
	if err != nil {
		s.logger.Errorf("failed to get team from repository: %v", err)
		return nil, apperror.NewInternalError("failed to get team", err)
	}
	if team == nil {
		s.logger.Print(ctx, "team not found", "team_name", teamName)
		return nil, apperror.NewTeamNotFoundError(teamName)
	}

	s.logger.Print(ctx, "team found", "team_name", team.TeamName, "members_count", len(team.Members))
	return team, nil
}
