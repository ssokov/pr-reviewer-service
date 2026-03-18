// nolint
//
//lint:file-ignore U1000 ignore unused code, it's generated
package db

import (
	"time"

	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
)

const condition = "?.? = ?"

// base filters
type applier func(query *orm.Query) (*orm.Query, error)

type search struct {
	appliers []applier
}

func (s *search) apply(query *orm.Query) {
	for _, applier := range s.appliers {
		query.Apply(applier)
	}
}

func (s *search) where(query *orm.Query, table, field string, value interface{}) {

	query.Where(condition, pg.Ident(table), pg.Ident(field), value)

}

func (s *search) WithApply(a applier) {
	if s.appliers == nil {
		s.appliers = []applier{}
	}
	s.appliers = append(s.appliers, a)
}

func (s *search) With(condition string, params ...interface{}) {
	s.WithApply(func(query *orm.Query) (*orm.Query, error) {
		return query.Where(condition, params...), nil
	})
}

// Searcher is interface for every generated filter
type Searcher interface {
	Apply(query *orm.Query) *orm.Query
	Q() applier

	With(condition string, params ...interface{})
	WithApply(a applier)
}

type PrSystemPrReviewerSearch struct {
	search

	ID         *int64
	PrID       *int64
	ReviewerID *int64
	AssignedAt *time.Time
}

func (s *PrSystemPrReviewerSearch) Apply(query *orm.Query) *orm.Query {
	if s.ID != nil {
		s.where(query, Tables.PrSystemPrReviewer.Alias, Columns.PrSystemPrReviewer.ID, s.ID)
	}
	if s.PrID != nil {
		s.where(query, Tables.PrSystemPrReviewer.Alias, Columns.PrSystemPrReviewer.PrID, s.PrID)
	}
	if s.ReviewerID != nil {
		s.where(query, Tables.PrSystemPrReviewer.Alias, Columns.PrSystemPrReviewer.ReviewerID, s.ReviewerID)
	}
	if s.AssignedAt != nil {
		s.where(query, Tables.PrSystemPrReviewer.Alias, Columns.PrSystemPrReviewer.AssignedAt, s.AssignedAt)
	}

	s.apply(query)

	return query
}

func (s *PrSystemPrReviewerSearch) Q() applier {
	return func(query *orm.Query) (*orm.Query, error) {
		return s.Apply(query), nil
	}
}

type PrSystemPullRequestSearch struct {
	search

	ID              *int64
	PullRequestID   *string
	PullRequestName *string
	AuthorID        *int64
	StatusID        *int
	MergedAt        *time.Time
	CreatedAt       *time.Time
}

func (s *PrSystemPullRequestSearch) Apply(query *orm.Query) *orm.Query {
	if s.ID != nil {
		s.where(query, Tables.PrSystemPullRequest.Alias, Columns.PrSystemPullRequest.ID, s.ID)
	}
	if s.PullRequestID != nil {
		s.where(query, Tables.PrSystemPullRequest.Alias, Columns.PrSystemPullRequest.PullRequestID, s.PullRequestID)
	}
	if s.PullRequestName != nil {
		s.where(query, Tables.PrSystemPullRequest.Alias, Columns.PrSystemPullRequest.PullRequestName, s.PullRequestName)
	}
	if s.AuthorID != nil {
		s.where(query, Tables.PrSystemPullRequest.Alias, Columns.PrSystemPullRequest.AuthorID, s.AuthorID)
	}
	if s.StatusID != nil {
		s.where(query, Tables.PrSystemPullRequest.Alias, Columns.PrSystemPullRequest.StatusID, s.StatusID)
	}
	if s.MergedAt != nil {
		s.where(query, Tables.PrSystemPullRequest.Alias, Columns.PrSystemPullRequest.MergedAt, s.MergedAt)
	}
	if s.CreatedAt != nil {
		s.where(query, Tables.PrSystemPullRequest.Alias, Columns.PrSystemPullRequest.CreatedAt, s.CreatedAt)
	}

	s.apply(query)

	return query
}

func (s *PrSystemPullRequestSearch) Q() applier {
	return func(query *orm.Query) (*orm.Query, error) {
		return s.Apply(query), nil
	}
}

type PrSystemStatusSearch struct {
	search

	ID        *int
	Name      *string
	CreatedAt *time.Time
}

func (s *PrSystemStatusSearch) Apply(query *orm.Query) *orm.Query {
	if s.ID != nil {
		s.where(query, Tables.PrSystemStatus.Alias, Columns.PrSystemStatus.ID, s.ID)
	}
	if s.Name != nil {
		s.where(query, Tables.PrSystemStatus.Alias, Columns.PrSystemStatus.Name, s.Name)
	}
	if s.CreatedAt != nil {
		s.where(query, Tables.PrSystemStatus.Alias, Columns.PrSystemStatus.CreatedAt, s.CreatedAt)
	}

	s.apply(query)

	return query
}

func (s *PrSystemStatusSearch) Q() applier {
	return func(query *orm.Query) (*orm.Query, error) {
		return s.Apply(query), nil
	}
}

type PrSystemTeamSearch struct {
	search

	ID        *int64
	Name      *string
	CreatedAt *time.Time
}

func (s *PrSystemTeamSearch) Apply(query *orm.Query) *orm.Query {
	if s.ID != nil {
		s.where(query, Tables.PrSystemTeam.Alias, Columns.PrSystemTeam.ID, s.ID)
	}
	if s.Name != nil {
		s.where(query, Tables.PrSystemTeam.Alias, Columns.PrSystemTeam.Name, s.Name)
	}
	if s.CreatedAt != nil {
		s.where(query, Tables.PrSystemTeam.Alias, Columns.PrSystemTeam.CreatedAt, s.CreatedAt)
	}

	s.apply(query)

	return query
}

func (s *PrSystemTeamSearch) Q() applier {
	return func(query *orm.Query) (*orm.Query, error) {
		return s.Apply(query), nil
	}
}

type PrSystemUserSearch struct {
	search

	ID        *int64
	UserID    *string
	Username  *string
	IsActive  *bool
	TeamID    *int64
	CreatedAt *time.Time
}

func (s *PrSystemUserSearch) Apply(query *orm.Query) *orm.Query {
	if s.ID != nil {
		s.where(query, Tables.PrSystemUser.Alias, Columns.PrSystemUser.ID, s.ID)
	}
	if s.UserID != nil {
		s.where(query, Tables.PrSystemUser.Alias, Columns.PrSystemUser.UserID, s.UserID)
	}
	if s.Username != nil {
		s.where(query, Tables.PrSystemUser.Alias, Columns.PrSystemUser.Username, s.Username)
	}
	if s.IsActive != nil {
		s.where(query, Tables.PrSystemUser.Alias, Columns.PrSystemUser.IsActive, s.IsActive)
	}
	if s.TeamID != nil {
		s.where(query, Tables.PrSystemUser.Alias, Columns.PrSystemUser.TeamID, s.TeamID)
	}
	if s.CreatedAt != nil {
		s.where(query, Tables.PrSystemUser.Alias, Columns.PrSystemUser.CreatedAt, s.CreatedAt)
	}

	s.apply(query)

	return query
}

func (s *PrSystemUserSearch) Q() applier {
	return func(query *orm.Query) (*orm.Query, error) {
		return s.Apply(query), nil
	}
}
