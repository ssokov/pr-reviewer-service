package postgres

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

type teamRepo struct {
	db *pgxpool.Pool
}

func NewTeamRepository(dbPool *pgxpool.Pool) repository.TeamRepository {
	return &teamRepo{
		db: dbPool,
	}
}

func (r *teamRepo) Create(ctx context.Context, team *domain.Team) (*domain.Team, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx) //nolint:errcheck

	query := `
		INSERT INTO pr_system.teams (name)
		VALUES ($1)
		RETURNING id, name, created_at
	`

	var dbTeam db.Team
	err = tx.QueryRow(ctx, query, team.TeamName).Scan(
		&dbTeam.ID,
		&dbTeam.TeamName,
		&dbTeam.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	if len(team.Members) > 0 {
		for _, member := range team.Members {
			updateUserQuery := `
				UPDATE pr_system.users
				SET team_id = $1
				WHERE user_id = $2
			`
			_, err = tx.Exec(ctx, updateUserQuery, dbTeam.ID, member.UserID)
			if err != nil {
				return nil, err
			}
		}

		membersQuery := `
			SELECT id, user_id, username, is_active, team_id, created_at
			FROM pr_system.users
			WHERE team_id = $1
		`
		rows, err := tx.Query(ctx, membersQuery, dbTeam.ID)
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		var members []domain.User
		for rows.Next() {
			var user domain.User
			var teamID *int64
			err = rows.Scan(
				&user.ID,
				&user.UserID,
				&user.Username,
				&user.IsActive,
				&teamID,
				&user.CreatedAt,
			)
			if err != nil {
				return nil, err
			}
			if teamID != nil {
				user.TeamID = *teamID
			}
			members = append(members, user)
		}

		if err = tx.Commit(ctx); err != nil {
			return nil, err
		}

		return mappers.TeamDBToDomain(&dbTeam, members), nil
	}

	if err = tx.Commit(ctx); err != nil {
		return nil, err
	}

	return mappers.TeamDBToDomain(&dbTeam, nil), nil
}

func (r *teamRepo) GetByName(ctx context.Context, teamName string) (*domain.Team, error) {
	query := `
		SELECT id, name, created_at
		FROM pr_system.teams
		WHERE name = $1
	`

	var dbTeam db.Team
	err := r.db.QueryRow(ctx, query, teamName).Scan(
		&dbTeam.ID,
		&dbTeam.TeamName,
		&dbTeam.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	membersQuery := `
		SELECT id, user_id, username, is_active, team_id, created_at
		FROM pr_system.users
		WHERE team_id = $1
	`
	rows, err := r.db.Query(ctx, membersQuery, dbTeam.ID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var members []domain.User
	for rows.Next() {
		var user domain.User
		var teamID *int64
		err = rows.Scan(
			&user.ID,
			&user.UserID,
			&user.Username,
			&user.IsActive,
			&teamID,
			&user.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		if teamID != nil {
			user.TeamID = *teamID
		}
		members = append(members, user)
	}

	return mappers.TeamDBToDomain(&dbTeam, members), nil
}

func (r *teamRepo) ExistsByName(ctx context.Context, teamName string) (bool, error) {
	query := `
		SELECT EXISTS(
			SELECT 1
			FROM pr_system.teams
			WHERE name = $1
		)
	`

	var exists bool
	err := r.db.QueryRow(ctx, query, teamName).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}
