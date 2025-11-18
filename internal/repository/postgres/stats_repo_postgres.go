package max_superuser

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/ssokov/pr-reviewer-service/internal/model/domain"
	"github.com/ssokov/pr-reviewer-service/internal/repository"
)

type statsRepo struct {
	db *pgxpool.Pool
}

func NewStatsRepository(dbPool *pgxpool.Pool) repository.StatsRepository {
	return &statsRepo{
		db: dbPool,
	}
}

func (r *statsRepo) GetTotalPRs(ctx context.Context) (int, error) {
	var count int
	err := r.db.QueryRow(ctx, `SELECT COUNT(*) FROM pr_system.pull_requests`).Scan(&count)
	return count, err
}

func (r *statsRepo) GetTotalUsers(ctx context.Context) (int, error) {
	var count int
	err := r.db.QueryRow(ctx, `SELECT COUNT(*) FROM pr_system.users`).Scan(&count)
	return count, err
}

func (r *statsRepo) GetActiveUsers(ctx context.Context) (int, error) {
	var count int
	err := r.db.QueryRow(ctx, `SELECT COUNT(*) FROM pr_system.users WHERE is_active = true`).Scan(&count)
	return count, err
}

func (r *statsRepo) GetPRsByStatus(ctx context.Context) (map[string]int, error) {
	query := `
		SELECT s.name, COUNT(pr.id)
		FROM pr_system.statuses s
		LEFT JOIN pr_system.pull_requests pr ON pr.status_id = s.id
		GROUP BY s.name
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make(map[string]int)
	for rows.Next() {
		var status string
		var count int
		if err := rows.Scan(&status, &count); err != nil {
			return nil, err
		}
		result[status] = count
	}

	return result, rows.Err()
}

func (r *statsRepo) GetTopReviewers(ctx context.Context, limit int) ([]domain.ReviewerStats, error) {
	query := `
		SELECT 
			u.user_id,
			u.username,
			COUNT(DISTINCT pr.id) as assigned_count,
			COUNT(DISTINCT CASE WHEN pr.status_id = 2 THEN pr.id END) as completed_count,
			COUNT(DISTINCT CASE WHEN pr.status_id = 1 THEN pr.id END) as active_count
		FROM pr_system.users u
		LEFT JOIN pr_system.pr_reviewers prr ON prr.reviewer_id = u.id
		LEFT JOIN pr_system.pull_requests pr ON pr.id = prr.pr_id
		WHERE u.is_active = true
		GROUP BY u.id, u.user_id, u.username
		HAVING COUNT(DISTINCT pr.id) > 0
		ORDER BY assigned_count DESC
		LIMIT $1
	`

	rows, err := r.db.Query(ctx, query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []domain.ReviewerStats

	for rows.Next() {
		var item domain.ReviewerStats
		if err := rows.Scan(&item.UserID, &item.Username, &item.AssignedCount, &item.CompletedCount, &item.ActiveCount); err != nil {
			return nil, err
		}
		result = append(result, item)
	}

	return result, rows.Err()
}
