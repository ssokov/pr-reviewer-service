package max_superuser

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/ssokov/pr-reviewer-service/internal/model/db"
	"github.com/ssokov/pr-reviewer-service/internal/model/domain"
	"github.com/ssokov/pr-reviewer-service/internal/repository"
	"github.com/ssokov/pr-reviewer-service/internal/repository/postgres/mappers"
)

type prRepo struct {
	db *pgxpool.Pool
}

func NewPRRepository(dbPool *pgxpool.Pool) repository.PRRepository {
	return &prRepo{
		db: dbPool,
	}
}

func (r *prRepo) Create(ctx context.Context, pr *domain.PullRequest) (*domain.PullRequest, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx) //nolint:errcheck

	var authorInternalID int64
	err = tx.QueryRow(ctx, `SELECT id FROM pr_system.users WHERE user_id = $1`, pr.AuthorID).Scan(&authorInternalID)
	if err != nil {
		return nil, err
	}

	var statusID int
	err = tx.QueryRow(ctx, `SELECT id FROM pr_system.statuses WHERE name = $1`, pr.Status).Scan(&statusID)
	if err != nil {
		return nil, err
	}

	query := `
		INSERT INTO pr_system.pull_requests (pull_request_id, pull_request_name, author_id, status_id, merged_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, pull_request_id, pull_request_name, created_at, merged_at
	`

	var dbPR db.PullRequest
	err = tx.QueryRow(ctx, query, pr.PullRequestID, pr.PullRequestName, authorInternalID, statusID, pr.MergedAt).Scan(
		&dbPR.ID,
		&dbPR.PullRequestID,
		&dbPR.PullRequestName,
		&dbPR.CreatedAt,
		&dbPR.MergedAt,
	)
	if err != nil {
		return nil, err
	}

	dbPR.AuthorID = authorInternalID
	dbPR.StatusID = statusID

	if len(pr.AssignedReviewers) > 0 {
		for _, reviewerUserID := range pr.AssignedReviewers {
			var reviewerInternalID int64
			err = tx.QueryRow(ctx, `SELECT id FROM pr_system.users WHERE user_id = $1`, reviewerUserID).Scan(&reviewerInternalID)
			if err != nil {
				return nil, err
			}

			_, err = tx.Exec(ctx, `
				INSERT INTO pr_system.pr_reviewers (pr_id, reviewer_id)
				VALUES ($1, $2)
			`, dbPR.ID, reviewerInternalID)
			if err != nil {
				return nil, err
			}
		}
	}

	if err = tx.Commit(ctx); err != nil {
		return nil, err
	}

	return mappers.PRDBToDomain(&dbPR, pr.AuthorID, pr.Status, pr.AssignedReviewers), nil
}

func (r *prRepo) GetByPRID(ctx context.Context, prID string) (*domain.PullRequest, error) {
	query := `
		SELECT 
			pr.id,
			pr.pull_request_id,
			pr.pull_request_name,
			u.user_id as author_user_id,
			s.name as status,
			pr.created_at,
			pr.merged_at
		FROM pr_system.pull_requests pr
		INNER JOIN pr_system.users u ON pr.author_id = u.id
		INNER JOIN pr_system.statuses s ON pr.status_id = s.id
		WHERE pr.pull_request_id = $1
	`

	var dbPR db.PullRequest
	var authorUserID string
	var statusStr string
	err := r.db.QueryRow(ctx, query, prID).Scan(
		&dbPR.ID,
		&dbPR.PullRequestID,
		&dbPR.PullRequestName,
		&authorUserID,
		&statusStr,
		&dbPR.CreatedAt,
		&dbPR.MergedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	reviewersQuery := `
		SELECT u.user_id
		FROM pr_system.pr_reviewers rev
		INNER JOIN pr_system.users u ON rev.reviewer_id = u.id
		WHERE rev.pr_id = $1
	`
	rows, err := r.db.Query(ctx, reviewersQuery, dbPR.ID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reviewers []string
	for rows.Next() {
		var userID string
		if err := rows.Scan(&userID); err != nil {
			return nil, err
		}
		reviewers = append(reviewers, userID)
	}

	return mappers.PRDBToDomain(&dbPR, authorUserID, domain.PRStatus(statusStr), reviewers), nil
}

func (r *prRepo) Update(ctx context.Context, pr *domain.PullRequest) (*domain.PullRequest, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx) //nolint:errcheck

	var statusID int
	err = tx.QueryRow(ctx, `SELECT id FROM pr_system.statuses WHERE name = $1`, pr.Status).Scan(&statusID)
	if err != nil {
		return nil, err
	}

	query := `
		UPDATE pr_system.pull_requests
		SET pull_request_name = $1, status_id = $2, merged_at = $3
		WHERE pull_request_id = $4
		RETURNING id, pull_request_id, pull_request_name, created_at, merged_at
	`

	var dbPR db.PullRequest
	err = tx.QueryRow(ctx, query, pr.PullRequestName, statusID, pr.MergedAt, pr.PullRequestID).Scan(
		&dbPR.ID,
		&dbPR.PullRequestID,
		&dbPR.PullRequestName,
		&dbPR.CreatedAt,
		&dbPR.MergedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	dbPR.StatusID = statusID

	_, err = tx.Exec(ctx, `DELETE FROM pr_system.pr_reviewers WHERE pr_id = $1`, dbPR.ID)
	if err != nil {
		return nil, err
	}

	if len(pr.AssignedReviewers) > 0 {
		for _, reviewerUserID := range pr.AssignedReviewers {
			var reviewerInternalID int64
			err = tx.QueryRow(ctx, `SELECT id FROM pr_system.users WHERE user_id = $1`, reviewerUserID).Scan(&reviewerInternalID)
			if err != nil {
				return nil, err
			}

			_, err = tx.Exec(ctx, `
				INSERT INTO pr_system.pr_reviewers (pr_id, reviewer_id)
				VALUES ($1, $2)
			`, dbPR.ID, reviewerInternalID)
			if err != nil {
				return nil, err
			}
		}
	}

	if err = tx.Commit(ctx); err != nil {
		return nil, err
	}

	return mappers.PRDBToDomain(&dbPR, pr.AuthorID, pr.Status, pr.AssignedReviewers), nil
}

func (r *prRepo) GetByReviewerID(ctx context.Context, userID string) ([]domain.PullRequest, error) {
	query := `
		SELECT 
			pr.id,
			pr.pull_request_id,
			pr.pull_request_name,
			u.user_id as author_user_id,
			s.name as status,
			pr.created_at,
			pr.merged_at
		FROM pr_system.pull_requests pr
		INNER JOIN pr_system.pr_reviewers rev ON pr.id = rev.pr_id
		INNER JOIN pr_system.users u ON pr.author_id = u.id
		INNER JOIN pr_system.users reviewer ON rev.reviewer_id = reviewer.id
		INNER JOIN pr_system.statuses s ON pr.status_id = s.id
		WHERE reviewer.user_id = $1
		ORDER BY pr.created_at DESC
	`

	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var pullRequests []domain.PullRequest
	for rows.Next() {
		var dbPR db.PullRequest
		var authorUserID string
		var statusStr string
		err = rows.Scan(
			&dbPR.ID,
			&dbPR.PullRequestID,
			&dbPR.PullRequestName,
			&authorUserID,
			&statusStr,
			&dbPR.CreatedAt,
			&dbPR.MergedAt,
		)
		if err != nil {
			return nil, err
		}
		pullRequests = append(pullRequests, *mappers.PRDBToDomain(&dbPR, authorUserID, domain.PRStatus(statusStr), nil))
	}

	return pullRequests, nil
}

func (r *prRepo) GetOpenPRsByUserIDs(ctx context.Context, userIDs []string) ([]domain.PullRequest, error) {
	if len(userIDs) == 0 {
		return []domain.PullRequest{}, nil
	}

	query := `
		SELECT DISTINCT
			pr.id,
			pr.pull_request_id,
			pr.pull_request_name,
			u.user_id as author_user_id,
			s.name as status,
			pr.created_at,
			pr.merged_at
		FROM pr_system.pull_requests pr
		INNER JOIN pr_system.pr_reviewers rev ON pr.id = rev.pr_id
		INNER JOIN pr_system.users u ON pr.author_id = u.id
		INNER JOIN pr_system.users reviewer ON rev.reviewer_id = reviewer.id
		INNER JOIN pr_system.statuses s ON pr.status_id = s.id
		WHERE reviewer.user_id = ANY($1) AND s.name = 'OPEN'
	`

	rows, err := r.db.Query(ctx, query, userIDs)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var pullRequests []domain.PullRequest
	for rows.Next() {
		var dbPR db.PullRequest
		var authorUserID string
		var statusStr string
		if err := rows.Scan(
			&dbPR.ID,
			&dbPR.PullRequestID,
			&dbPR.PullRequestName,
			&authorUserID,
			&statusStr,
			&dbPR.CreatedAt,
			&dbPR.MergedAt,
		); err != nil {
			return nil, err
		}
		pullRequests = append(pullRequests, *mappers.PRDBToDomain(&dbPR, authorUserID, domain.PRStatus(statusStr), nil))
	}

	return pullRequests, rows.Err()
}
