package service

import (
	"context"
	"time"

	"github.com/ssokov/pr-reviewer-service/internal/apperror"
	"github.com/ssokov/pr-reviewer-service/internal/model/domain"
	"github.com/ssokov/pr-reviewer-service/internal/repository"
	"github.com/vmkteam/embedlog"
)

type prService struct {
	prRepo   repository.PRRepository
	userRepo repository.UserRepository
	teamRepo repository.TeamRepository
	logger   embedlog.Logger
}

func NewPRService(prRepo repository.PRRepository, userRepo repository.UserRepository, teamRepo repository.TeamRepository, logger embedlog.Logger) PRService {
	return &prService{
		prRepo:   prRepo,
		userRepo: userRepo,
		teamRepo: teamRepo,
		logger:   logger,
	}
}

func (s *prService) CreatePR(ctx context.Context, authorID string, pr *domain.PullRequest) (*domain.PullRequest, error) {
	if pr.PullRequestID == "" {
		return nil, apperror.NewInvalidInputError("pull_request_id is required")
	}
	if pr.PullRequestName == "" {
		return nil, apperror.NewInvalidInputError("pull_request_name is required")
	}
	if authorID == "" {
		return nil, apperror.NewInvalidInputError("author_id is required")
	}

	s.logger.Print(ctx, "creating PR", "pr_id", pr.PullRequestID, "author_id", authorID)

	author, err := s.userRepo.GetByUserID(ctx, authorID)
	if err != nil {
		s.logger.Errorf("failed to get author: %v", err)
		return nil, apperror.NewInternalError("failed to get author", err)
	}
	if author == nil {
		s.logger.Print(ctx, "author not found", "author_id", authorID)
		return nil, apperror.NewUserNotFoundError(authorID)
	}

	if !author.IsActive {
		s.logger.Print(ctx, "author is not active", "author_id", authorID)
		return nil, apperror.NewInvalidInputError("author is not active")
	}

	reviewers, err := s.autoAssignReviewers(ctx, author)
	if err != nil {
		s.logger.Errorf("failed to assign reviewers: %v", err)
		return nil, err
	}

	pr.AssignedReviewers = reviewers
	pr.Status = domain.PRStatusOpen

	createdPR, err := s.prRepo.Create(ctx, pr)
	if err != nil {
		s.logger.Errorf("failed to create PR: %v", err)
		return nil, apperror.NewInternalError("failed to create PR", err)
	}

	s.logger.Print(ctx, "PR created successfully", "pr_id", createdPR.PullRequestID, "reviewers_count", len(reviewers))
	return createdPR, nil
}

func (s *prService) MergePR(ctx context.Context, prID string) (*domain.PullRequest, error) {
	if prID == "" {
		return nil, apperror.NewInvalidInputError("pull_request_id is required")
	}

	s.logger.Print(ctx, "merging PR", "pr_id", prID)

	pr, err := s.prRepo.GetByPRID(ctx, prID)
	if err != nil {
		s.logger.Errorf("failed to get PR: %v", err)
		return nil, apperror.NewInternalError("failed to get PR", err)
	}
	if pr == nil {
		s.logger.Print(ctx, "PR not found", "pr_id", prID)
		return nil, apperror.NewPRNotFoundError(prID)
	}

	if pr.Status == domain.PRStatusMerged {
		s.logger.Print(ctx, "PR already merged, returning current state", "pr_id", prID)
		return pr, nil
	}

	now := time.Now()
	pr.Status = domain.PRStatusMerged
	pr.MergedAt = &now

	updatedPR, err := s.prRepo.Update(ctx, pr)
	if err != nil {
		s.logger.Errorf("failed to merge PR: %v", err)
		return nil, apperror.NewInternalError("failed to merge PR", err)
	}

	s.logger.Print(ctx, "PR merged successfully", "pr_id", prID)
	return updatedPR, nil
}

func (s *prService) ReassignReviewer(ctx context.Context, prID string, oldUserID string) (*domain.PullRequest, string, error) {
	if prID == "" {
		return nil, "", apperror.NewInvalidInputError("pull_request_id is required")
	}
	if oldUserID == "" {
		return nil, "", apperror.NewInvalidInputError("old_user_id is required")
	}

	s.logger.Print(ctx, "reassigning reviewer", "pr_id", prID, "old_user_id", oldUserID)

	pr, err := s.prRepo.GetByPRID(ctx, prID)
	if err != nil {
		s.logger.Errorf("failed to get PR: %v", err)
		return nil, "", apperror.NewInternalError("failed to get PR", err)
	}
	if pr == nil {
		s.logger.Print(ctx, "PR not found", "pr_id", prID)
		return nil, "", apperror.NewPRNotFoundError(prID)
	}

	if pr.Status == domain.PRStatusMerged {
		s.logger.Print(ctx, "cannot reassign on merged PR", "pr_id", prID)
		return nil, "", apperror.NewPRMergedError(prID)
	}

	isAssigned := false
	for _, reviewerID := range pr.AssignedReviewers {
		if reviewerID == oldUserID {
			isAssigned = true
			break
		}
	}
	if !isAssigned {
		s.logger.Print(ctx, "user not assigned to PR", "pr_id", prID, "user_id", oldUserID)
		return nil, "", apperror.NewNotAssignedError(prID, oldUserID)
	}

	oldUser, err := s.userRepo.GetByUserID(ctx, oldUserID)
	if err != nil {
		s.logger.Errorf("failed to get old user: %v", err)
		return nil, "", apperror.NewInternalError("failed to get old user", err)
	}
	if oldUser == nil {
		s.logger.Print(ctx, "old user not found", "user_id", oldUserID)
		return nil, "", apperror.NewUserNotFoundError(oldUserID)
	}

	newReviewers, err := s.autoAssignReviewers(ctx, oldUser)
	if err != nil {
		s.logger.Errorf("failed to assign new reviewer: %v", err)
		return nil, "", err
	}

	if len(newReviewers) == 0 {
		s.logger.Print(ctx, "no available reviewers", "team_id", oldUser.TeamID)
		return nil, "", apperror.NewInvalidInputError("no available reviewers in team")
	}

	newReviewerID := newReviewers[0]

	updatedReviewers := make([]string, 0, len(pr.AssignedReviewers))
	for _, reviewer := range pr.AssignedReviewers {
		if reviewer != oldUserID {
			updatedReviewers = append(updatedReviewers, reviewer)
		}
	}
	updatedReviewers = append(updatedReviewers, newReviewerID)
	pr.AssignedReviewers = updatedReviewers

	updatedPR, err := s.prRepo.Update(ctx, pr)
	if err != nil {
		s.logger.Errorf("failed to update PR: %v", err)
		return nil, "", apperror.NewInternalError("failed to update PR", err)
	}

	s.logger.Print(ctx, "reviewer reassigned", "pr_id", prID, "old_user_id", oldUserID, "new_user_id", newReviewerID)
	return updatedPR, newReviewerID, nil
}

func (s *prService) autoAssignReviewers(ctx context.Context, user *domain.User) ([]string, error) {
	if user.TeamID == 0 {
		return nil, apperror.NewInvalidInputError("user has no team")
	}

	teamMembers, err := s.userRepo.GetByTeamID(ctx, user.TeamID)
	if err != nil {
		return nil, apperror.NewInternalError("failed to get team members", err)
	}

	var activeReviewers []string
	for _, member := range teamMembers {
		if member.UserID != user.UserID && member.IsActive {
			activeReviewers = append(activeReviewers, member.UserID)
		}
	}

	if len(activeReviewers) == 0 {
		return nil, apperror.NewInvalidInputError("no active reviewers in team")
	}

	return activeReviewers, nil
}
