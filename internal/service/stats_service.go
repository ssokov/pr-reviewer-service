package service

import (
	"context"

	"github.com/ssokov/pr-reviewer-service/internal/model/dto"
	"github.com/ssokov/pr-reviewer-service/internal/repository"
	"github.com/vmkteam/embedlog"
)

type StatsService interface {
	GetStats(ctx context.Context) (*dto.StatsResponse, error)
}

type statsService struct {
	statsRepo repository.StatsRepository
	logger    embedlog.Logger
}

func NewStatsService(statsRepo repository.StatsRepository, logger embedlog.Logger) StatsService {
	return &statsService{
		statsRepo: statsRepo,
		logger:    logger,
	}
}

func (s *statsService) GetStats(ctx context.Context) (*dto.StatsResponse, error) {
	s.logger.Print(ctx, "getting statistics")

	totalPRs, err := s.statsRepo.GetTotalPRs(ctx)
	if err != nil {
		s.logger.Print(ctx, "failed to get total PRs", "error", err)
		return nil, err
	}

	totalUsers, err := s.statsRepo.GetTotalUsers(ctx)
	if err != nil {
		s.logger.Print(ctx, "failed to get total users", "error", err)
		return nil, err
	}

	activeUsers, err := s.statsRepo.GetActiveUsers(ctx)
	if err != nil {
		s.logger.Print(ctx, "failed to get active users", "error", err)
		return nil, err
	}

	prsByStatus, err := s.statsRepo.GetPRsByStatus(ctx)
	if err != nil {
		s.logger.Print(ctx, "failed to get PRs by status", "error", err)
		return nil, err
	}

	topReviewers, err := s.statsRepo.GetTopReviewers(ctx, 10)
	if err != nil {
		s.logger.Print(ctx, "failed to get top reviewers", "error", err)
		return nil, err
	}

	var prsByStatusDTO []dto.PRStatsItem
	for status, count := range prsByStatus {
		prsByStatusDTO = append(prsByStatusDTO, dto.PRStatsItem{
			Status: status,
			Count:  count,
		})
	}

	var topReviewersDTO []dto.UserStatsItem
	for _, reviewer := range topReviewers {
		topReviewersDTO = append(topReviewersDTO, dto.UserStatsItem{
			UserID:         reviewer.UserID,
			Username:       reviewer.Username,
			AssignedCount:  reviewer.AssignedCount,
			CompletedCount: reviewer.CompletedCount,
			ActiveCount:    reviewer.ActiveCount,
		})
	}

	return &dto.StatsResponse{
		TotalPRs:     totalPRs,
		TotalUsers:   totalUsers,
		ActiveUsers:  activeUsers,
		PRsByStatus:  prsByStatusDTO,
		TopReviewers: topReviewersDTO,
	}, nil
}
