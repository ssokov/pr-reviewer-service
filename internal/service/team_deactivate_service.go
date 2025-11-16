package service

import (
	"context"

	"github.com/ssokov/pr-reviewer-service/internal/model/domain"
)

func (s *teamService) DeactivateTeam(ctx context.Context, teamName string) ([]domain.User, []domain.PullRequest, error) {
	s.logger.Print(ctx, "deactivating team", "team_name", teamName)

	team, err := s.teamRepo.GetByName(ctx, teamName)
	if err != nil {
		s.logger.Print(ctx, "failed to get team", "error", err)
		return nil, nil, err
	}

	deactivatedUsers, err := s.userRepo.DeactivateByTeamID(ctx, team.ID)
	if err != nil {
		s.logger.Print(ctx, "failed to deactivate users", "error", err)
		return nil, nil, err
	}

	if len(deactivatedUsers) == 0 {
		s.logger.Print(ctx, "no users to deactivate")
		return []domain.User{}, []domain.PullRequest{}, nil
	}

	userIDs := make([]string, len(deactivatedUsers))
	for i, user := range deactivatedUsers {
		userIDs[i] = user.UserID
	}

	openPRs, err := s.prRepo.GetOpenPRsByUserIDs(ctx, userIDs)
	if err != nil {
		s.logger.Print(ctx, "failed to get open PRs", "error", err)
		return nil, nil, err
	}

	s.logger.Print(ctx, "team deactivated", "users_count", len(deactivatedUsers), "open_prs_count", len(openPRs))

	return deactivatedUsers, openPRs, nil
}
