package pr

import (
	"context"
	"time"

	"github.com/ssokov/pr-reviewer-service/internal/db"
	"github.com/vmkteam/embedlog"
)

type PrService struct {
	repo   *db.PrReviewerServiceRepo
	logger embedlog.Logger
}

func NewPrService(repo *db.PrReviewerServiceRepo, logger embedlog.Logger) *PrService {
	return &PrService{
		repo:   repo,
		logger: logger,
	}
}

func (s *PrService) CreatePR(ctx context.Context, authorID string, pullRequest *PullRequest) (*PullRequest, error) {
	if pullRequest.PullRequestID == "" {
		return nil, NewInvalidInputError("pull_request_id is required")
	}
	if pullRequest.PullRequestName == "" {
		return nil, NewInvalidInputError("pull_request_name is required")
	}
	if authorID == "" {
		return nil, NewInvalidInputError("author_id is required")
	}

	s.logger.Print(ctx, "creating PR", "pr_id", pullRequest.PullRequestID, "author_id", authorID)

	authorDB, err := s.repo.GetByUserID(ctx, authorID)
	if err != nil {
		s.logger.Errorf("failed to get author: %v", err)
		return nil, NewInternalError("failed to get author", err)
	}
	if authorDB == nil {
		s.logger.Print(ctx, "author not found", "author_id", authorID)
		return nil, NewUserNotFoundError(authorID)
	}
	author := NewUser(authorDB)

	if !author.IsActive {
		s.logger.Print(ctx, "author is not active", "author_id", authorID)
		return nil, NewInvalidInputError("author is not active")
	}

	reviewers, err := s.autoAssignReviewers(ctx, author)
	if err != nil {
		s.logger.Errorf("failed to assign reviewers: %v", err)
		return nil, err
	}

	pullRequest.AssignedReviewers = reviewers
	pullRequest.Status = PRStatusOpen

	createdDB, err := s.repo.CreatePR(ctx, NewDBPullRequest(pullRequest))
	if err != nil {
		s.logger.Errorf("failed to create PR: %v", err)
		return nil, NewInternalError("failed to create PR", err)
	}
	createdPR := NewPullRequest(createdDB)

	s.logger.Print(ctx, "PR created successfully", "pr_id", createdPR.PullRequestID, "reviewers_count", len(reviewers))
	return createdPR, nil
}

func (s *PrService) MergePR(ctx context.Context, prID string) (*PullRequest, error) {
	if prID == "" {
		return nil, NewInvalidInputError("pull_request_id is required")
	}

	s.logger.Print(ctx, "merging PR", "pr_id", prID)

	pullRequestDB, err := s.repo.GetByPRID(ctx, prID)
	if err != nil {
		s.logger.Errorf("failed to get PR: %v", err)
		return nil, NewInternalError("failed to get PR", err)
	}
	if pullRequestDB == nil {
		s.logger.Print(ctx, "PR not found", "pr_id", prID)
		return nil, NewPRNotFoundError(prID)
	}
	pullRequest := NewPullRequest(pullRequestDB)

	if pullRequest.Status == PRStatusMerged {
		s.logger.Print(ctx, "PR already merged, returning current state", "pr_id", prID)
		return pullRequest, nil
	}

	now := time.Now()
	pullRequest.Status = PRStatusMerged
	pullRequest.MergedAt = &now

	updatedDB, err := s.repo.UpdatePR(ctx, NewDBPullRequest(pullRequest))
	if err != nil {
		s.logger.Errorf("failed to merge PR: %v", err)
		return nil, NewInternalError("failed to merge PR", err)
	}
	updatedPR := NewPullRequest(updatedDB)

	s.logger.Print(ctx, "PR merged successfully", "pr_id", prID)
	return updatedPR, nil
}

func (s *PrService) ReassignReviewer(ctx context.Context, prID string, oldUserID string) (*PullRequest, string, error) {
	if prID == "" {
		return nil, "", NewInvalidInputError("pull_request_id is required")
	}
	if oldUserID == "" {
		return nil, "", NewInvalidInputError("old_user_id is required")
	}

	s.logger.Print(ctx, "reassigning reviewer", "pr_id", prID, "old_user_id", oldUserID)

	pullRequestDB, err := s.repo.GetByPRID(ctx, prID)
	if err != nil {
		s.logger.Errorf("failed to get PR: %v", err)
		return nil, "", NewInternalError("failed to get PR", err)
	}
	if pullRequestDB == nil {
		s.logger.Print(ctx, "PR not found", "pr_id", prID)
		return nil, "", NewPRNotFoundError(prID)
	}
	pullRequest := NewPullRequest(pullRequestDB)

	if pullRequest.Status == PRStatusMerged {
		s.logger.Print(ctx, "cannot reassign on merged PR", "pr_id", prID)
		return nil, "", NewPRMergedError(prID)
	}

	isAssigned := false
	for _, reviewerID := range pullRequest.AssignedReviewers {
		if reviewerID == oldUserID {
			isAssigned = true
			break
		}
	}
	if !isAssigned {
		s.logger.Print(ctx, "user not assigned to PR", "pr_id", prID, "user_id", oldUserID)
		return nil, "", NewNotAssignedError(oldUserID, prID)
	}

	oldUserDB, err := s.repo.GetByUserID(ctx, oldUserID)
	if err != nil {
		s.logger.Errorf("failed to get old user: %v", err)
		return nil, "", NewInternalError("failed to get old user", err)
	}
	if oldUserDB == nil {
		s.logger.Print(ctx, "old user not found", "user_id", oldUserID)
		return nil, "", NewUserNotFoundError(oldUserID)
	}
	oldUser := NewUser(oldUserDB)

	newReviewers, err := s.autoAssignReviewers(ctx, oldUser)
	if err != nil {
		s.logger.Errorf("failed to assign new reviewer: %v", err)
		return nil, "", err
	}

	if len(newReviewers) == 0 {
		s.logger.Print(ctx, "no available reviewers", "team_id", oldUser.TeamID)
		return nil, "", NewInvalidInputError("no available reviewers in team")
	}

	newReviewerID := newReviewers[0]

	updatedReviewers := make([]string, 0, len(pullRequest.AssignedReviewers))
	for _, reviewer := range pullRequest.AssignedReviewers {
		if reviewer != oldUserID {
			updatedReviewers = append(updatedReviewers, reviewer)
		}
	}
	updatedReviewers = append(updatedReviewers, newReviewerID)
	pullRequest.AssignedReviewers = updatedReviewers

	updatedDB, err := s.repo.UpdatePR(ctx, NewDBPullRequest(pullRequest))
	if err != nil {
		s.logger.Errorf("failed to update PR: %v", err)
		return nil, "", NewInternalError("failed to update PR", err)
	}
	updatedPR := NewPullRequest(updatedDB)

	s.logger.Print(ctx, "reviewer reassigned", "pr_id", prID, "old_user_id", oldUserID, "new_user_id", newReviewerID)
	return updatedPR, newReviewerID, nil
}

func (s *PrService) autoAssignReviewers(ctx context.Context, user *User) ([]string, error) {
	if user.TeamID == 0 {
		return nil, NewInvalidInputError("user has no team")
	}

	teamMembersDB, err := s.repo.GetByTeamID(ctx, user.TeamID)
	if err != nil {
		return nil, NewInternalError("failed to get team members", err)
	}
	teamMembers := NewUsers(teamMembersDB)

	var activeReviewers []string
	for _, member := range teamMembers {
		if member.UserID != user.UserID && member.IsActive {
			activeReviewers = append(activeReviewers, member.UserID)
		}
	}

	if len(activeReviewers) == 0 {
		return nil, NewInvalidInputError("no active reviewers in team")
	}

	return activeReviewers, nil
}

func (s *PrService) GetStats(ctx context.Context) (*StatsResponse, error) {
	s.logger.Print(ctx, "getting statistics")

	totalPRs, err := s.repo.GetTotalPRs(ctx)
	if err != nil {
		s.logger.Print(ctx, "failed to get total PRs", "error", err)
		return nil, err
	}

	totalUsers, err := s.repo.GetTotalUsers(ctx)
	if err != nil {
		s.logger.Print(ctx, "failed to get total users", "error", err)
		return nil, err
	}

	activeUsers, err := s.repo.GetActiveUsers(ctx)
	if err != nil {
		s.logger.Print(ctx, "failed to get active users", "error", err)
		return nil, err
	}

	prsByStatus, err := s.repo.GetPRsByStatus(ctx)
	if err != nil {
		s.logger.Print(ctx, "failed to get PRs by status", "error", err)
		return nil, err
	}

	topReviewersDB, err := s.repo.GetTopReviewers(ctx, 10)
	if err != nil {
		s.logger.Print(ctx, "failed to get top reviewers", "error", err)
		return nil, err
	}
	topReviewers := NewReviewerStatsList(topReviewersDB)

	var prsByStatusDTO []PRStatsItem
	for status, count := range prsByStatus {
		prsByStatusDTO = append(prsByStatusDTO, PRStatsItem{
			Status: status,
			Count:  count,
		})
	}

	var topReviewersDTO []UserStatsItem
	for _, reviewer := range topReviewers {
		topReviewersDTO = append(topReviewersDTO, UserStatsItem(reviewer))
	}

	return &StatsResponse{
		TotalPRs:     totalPRs,
		TotalUsers:   totalUsers,
		ActiveUsers:  activeUsers,
		PRsByStatus:  prsByStatusDTO,
		TopReviewers: topReviewersDTO,
	}, nil
}

func (s *PrService) AddTeam(ctx context.Context, team *Team) (*Team, error) {
	if team.TeamName == "" {
		return nil, NewInvalidInputError("team_name is required")
	}

	if len(team.Members) == 0 {
		return nil, NewInvalidInputError("team must have at least one member")
	}

	s.logger.Print(ctx, "creating team", "team_name", team.TeamName, "members_count", len(team.Members))

	exists, err := s.repo.ExistsByName(ctx, team.TeamName)
	if err != nil {
		s.logger.Errorf("failed to check team existence: %v", err)
		return nil, NewInternalError("failed to check team existence", err)
	}
	if exists {
		s.logger.Print(ctx, "team already exists", "team_name", team.TeamName)
		return nil, NewTeamExistsError(team.TeamName)
	}

	createdDB, err := s.repo.CreateTeam(ctx, NewDBTeam(team))
	if err != nil {
		s.logger.Errorf("failed to create team: %v", err)
		return nil, NewInternalError("failed to create team", err)
	}
	createdTeam := NewTeam(createdDB)

	s.logger.Print(ctx, "team created successfully", "team_name", createdTeam.TeamName, "team_id", createdTeam.ID)
	return createdTeam, nil
}

func (s *PrService) GetTeam(ctx context.Context, teamName string) (*Team, error) {
	s.logger.Print(ctx, "getting team", "team_name", teamName)

	teamDB, err := s.repo.GetByName(ctx, teamName)
	if err != nil {
		s.logger.Errorf("failed to get team from repository: %v", err)
		return nil, NewInternalError("failed to get team", err)
	}
	if teamDB == nil {
		s.logger.Print(ctx, "team not found", "team_name", teamName)
		return nil, NewTeamNotFoundError(teamName)
	}
	team := NewTeam(teamDB)

	s.logger.Print(ctx, "team found", "team_name", team.TeamName, "members_count", len(team.Members))
	return team, nil
}

func (s *PrService) SetIsActive(ctx context.Context, userID string, isActive bool) (*User, error) {
	if userID == "" {
		return nil, NewInvalidInputError("user_id is required")
	}

	s.logger.Print(ctx, "setting user active status", "user_id", userID, "is_active", isActive)

	userDB, err := s.repo.SetIsActive(ctx, userID, isActive)
	if err != nil {
		s.logger.Errorf("failed to set user active status: %v", err)
		return nil, NewInternalError("failed to set user active status", err)
	}

	if userDB == nil {
		s.logger.Print(ctx, "user not found", "user_id", userID)
		return nil, NewUserNotFoundError(userID)
	}
	user := NewUser(userDB)

	s.logger.Print(ctx, "user active status updated", "user_id", userID, "is_active", isActive)
	return user, nil
}

func (s *PrService) GetReview(ctx context.Context, userID string) ([]PullRequest, error) {
	if userID == "" {
		return nil, NewInvalidInputError("user_id is required")
	}

	s.logger.Print(ctx, "getting reviews for user", "user_id", userID)

	userDB, err := s.repo.GetByUserID(ctx, userID)
	if err != nil {
		s.logger.Errorf("failed to get user: %v", err)
		return nil, NewInternalError("failed to get user", err)
	}

	if userDB == nil {
		s.logger.Print(ctx, "user not found", "user_id", userID)
		return nil, NewUserNotFoundError(userID)
	}

	pullRequestsDB, err := s.repo.GetByReviewerID(ctx, userID)
	if err != nil {
		s.logger.Errorf("failed to get reviews: %v", err)
		return nil, NewInternalError("failed to get reviews", err)
	}
	pullRequests := NewPullRequests(pullRequestsDB)

	s.logger.Print(ctx, "reviews retrieved", "user_id", userID, "count", len(pullRequests))
	return pullRequests, nil
}
