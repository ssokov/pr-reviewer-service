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

type userRepo struct {
	db *pgxpool.Pool
}

func NewUserRepository(dbPool *pgxpool.Pool) repository.UserRepository {
	return &userRepo{
		db: dbPool,
	}
}

func (r *userRepo) Create(ctx context.Context, user *domain.User) (*domain.User, error) {
	query := `
		INSERT INTO pr_system.users (user_id, username, is_active, team_id)
		VALUES ($1, $2, $3, $4)
		RETURNING id, user_id, username, is_active, team_id, created_at
	`

	var dbUser db.User
	var teamID *int64
	err := r.db.QueryRow(ctx, query, user.UserID, user.Username, user.IsActive, nullInt64(user.TeamID)).Scan(
		&dbUser.ID,
		&dbUser.UserID,
		&dbUser.Username,
		&dbUser.IsActive,
		&teamID,
		&dbUser.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	if teamID != nil {
		dbUser.TeamID = *teamID
	}

	return mappers.UserDBToDomain(&dbUser, ""), nil
}

func (r *userRepo) Update(ctx context.Context, user *domain.User) (*domain.User, error) {
	query := `
		UPDATE pr_system.users
		SET username = $1, is_active = $2, team_id = $3
		WHERE user_id = $4
		RETURNING id, user_id, username, is_active, team_id, created_at
	`

	var dbUser db.User
	var teamID *int64
	err := r.db.QueryRow(ctx, query, user.Username, user.IsActive, nullInt64(user.TeamID), user.UserID).Scan(
		&dbUser.ID,
		&dbUser.UserID,
		&dbUser.Username,
		&dbUser.IsActive,
		&teamID,
		&dbUser.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	if teamID != nil {
		dbUser.TeamID = *teamID
	}

	return mappers.UserDBToDomain(&dbUser, ""), nil
}

func (r *userRepo) GetByUserID(ctx context.Context, userID string) (*domain.User, error) {
	query := `
		SELECT u.id, u.user_id, u.username, u.is_active, u.team_id, u.created_at, t.name as team_name
		FROM pr_system.users u
		LEFT JOIN pr_system.teams t ON u.team_id = t.id
		WHERE u.user_id = $1
	`

	var dbUser db.User
	var teamID *int64
	var teamName *string
	err := r.db.QueryRow(ctx, query, userID).Scan(
		&dbUser.ID,
		&dbUser.UserID,
		&dbUser.Username,
		&dbUser.IsActive,
		&teamID,
		&dbUser.CreatedAt,
		&teamName,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	if teamID != nil {
		dbUser.TeamID = *teamID
	}
	teamNameStr := ""
	if teamName != nil {
		teamNameStr = *teamName
	}

	return mappers.UserDBToDomain(&dbUser, teamNameStr), nil
}

func (r *userRepo) GetByTeamID(ctx context.Context, teamID int64) ([]domain.User, error) {
	query := `
		SELECT u.id, u.user_id, u.username, u.is_active, u.team_id, u.created_at, t.name as team_name
		FROM pr_system.users u
		LEFT JOIN pr_system.teams t ON u.team_id = t.id
		WHERE u.team_id = $1
	`

	rows, err := r.db.Query(ctx, query, teamID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []domain.User
	for rows.Next() {
		var dbUser db.User
		var teamIDPtr *int64
		var teamName *string
		err = rows.Scan(
			&dbUser.ID,
			&dbUser.UserID,
			&dbUser.Username,
			&dbUser.IsActive,
			&teamIDPtr,
			&dbUser.CreatedAt,
			&teamName,
		)
		if err != nil {
			return nil, err
		}

		if teamIDPtr != nil {
			dbUser.TeamID = *teamIDPtr
		}
		teamNameStr := ""
		if teamName != nil {
			teamNameStr = *teamName
		}
		users = append(users, *mappers.UserDBToDomain(&dbUser, teamNameStr))
	}

	return users, nil
}

func (r *userRepo) SetIsActive(ctx context.Context, userID string, isActive bool) (*domain.User, error) {
	query := `
		UPDATE pr_system.users
		SET is_active = $1
		WHERE user_id = $2
		RETURNING id, user_id, username, is_active, team_id, created_at
	`

	var dbUser db.User
	var teamID *int64
	err := r.db.QueryRow(ctx, query, isActive, userID).Scan(
		&dbUser.ID,
		&dbUser.UserID,
		&dbUser.Username,
		&dbUser.IsActive,
		&teamID,
		&dbUser.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	if teamID != nil {
		dbUser.TeamID = *teamID
	}

	return mappers.UserDBToDomain(&dbUser, ""), nil
}

func (r *userRepo) GetByReviewerID(ctx context.Context, userID string) ([]domain.PullRequest, error) {
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

func (r *userRepo) DeactivateByTeamID(ctx context.Context, teamID int64) ([]domain.User, error) {
	query := `
		UPDATE pr_system.users u
		SET is_active = false
		FROM pr_system.teams t
		WHERE u.team_id = $1 AND u.is_active = true AND u.team_id = t.id
		RETURNING u.id, u.user_id, u.username, u.is_active, u.team_id, u.created_at, t.team_name
	`

	rows, err := r.db.Query(ctx, query, teamID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []domain.User
	for rows.Next() {
		var dbUser db.User
		var teamName string
		if err := rows.Scan(&dbUser.ID, &dbUser.UserID, &dbUser.Username, &dbUser.IsActive, &dbUser.TeamID, &dbUser.CreatedAt, &teamName); err != nil {
			return nil, err
		}
		users = append(users, *mappers.UserDBToDomain(&dbUser, teamName))
	}

	return users, rows.Err()
}

func nullInt64(val int64) *int64 {
	if val == 0 {
		return nil
	}
	return &val
}
