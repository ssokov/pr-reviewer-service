package db

import (
	"context"
	"errors"
	"fmt"
	"sort"

	"github.com/go-pg/pg/v10"
)

type PrReviewerServiceRepo struct {
	db *pg.DB
}

func NewPrReviewerServiceRepo(db *pg.DB) *PrReviewerServiceRepo {
	return &PrReviewerServiceRepo{db: db}
}

func (r *PrReviewerServiceRepo) CreatePR(ctx context.Context, pr *PullRequest) (*PullRequest, error) {
	var created *PullRequest

	err := r.db.RunInTransaction(ctx, func(tx *pg.Tx) error {
		authorInternalID, err := resolveUserInternalID(ctx, tx, pr.AuthorUserID)
		if err != nil {
			return err
		}
		statusID, err := resolveStatusID(ctx, tx, pr.Status)
		if err != nil {
			return err
		}

		model := &PrSystemPullRequest{
			PullRequestID:   pr.PullRequestID,
			PullRequestName: pr.PullRequestName,
			AuthorID:        authorInternalID,
			StatusID:        statusID,
			MergedAt:        pr.MergedAt,
		}
		if validateErrs, ok := model.Validate(); !ok {
			return fmt.Errorf("CreatePR: validate model: %v", validateErrs)
		}

		if _, err = tx.ModelContext(ctx, model).Returning("*").Insert(); err != nil {
			return err
		}

		if err = replacePRReviewers(ctx, tx, model.ID, pr.AssignedReviewers); err != nil {
			return err
		}

		result := toPullRequest(model, pr.AuthorUserID, pr.Status, pr.AssignedReviewers)
		created = &result
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("CreatePR: %w", err)
	}

	return created, nil
}

func (r *PrReviewerServiceRepo) GetByPRID(ctx context.Context, prID string) (*PullRequest, error) {
	model := new(PrSystemPullRequest)
	search := PrSystemPullRequestSearch{PullRequestID: &prID}

	q := r.db.ModelContext(ctx, model).
		Relation(Columns.PrSystemPullRequest.Author).
		Relation(Columns.PrSystemPullRequest.Status)
	search.Apply(q)
	if err := q.Select(); err != nil {
		if errors.Is(err, pg.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("GetByPRID: %w", err)
	}

	reviewers, err := getReviewerIDsByPR(ctx, r.db, model.ID)
	if err != nil {
		return nil, fmt.Errorf("GetByPRID reviewers: %w", err)
	}

	authorUserID := ""
	if model.Author != nil {
		authorUserID = model.Author.UserID
	}

	status := StatusIDToName(model.StatusID)
	if model.Status != nil && model.Status.Name != "" {
		status = PRStatus(model.Status.Name)
	}

	result := toPullRequest(model, authorUserID, status, reviewers)
	return &result, nil
}

func (r *PrReviewerServiceRepo) UpdatePR(ctx context.Context, pr *PullRequest) (*PullRequest, error) {
	var updated *PullRequest
	found := false

	err := r.db.RunInTransaction(ctx, func(tx *pg.Tx) error {
		statusID, err := resolveStatusID(ctx, tx, pr.Status)
		if err != nil {
			return err
		}

		model := new(PrSystemPullRequest)
		search := PrSystemPullRequestSearch{PullRequestID: &pr.PullRequestID}

		q := tx.ModelContext(ctx, model).
			Set("? = ?", pg.Ident(Columns.PrSystemPullRequest.PullRequestName), pr.PullRequestName).
			Set("? = ?", pg.Ident(Columns.PrSystemPullRequest.StatusID), statusID).
			Set("? = ?", pg.Ident(Columns.PrSystemPullRequest.MergedAt), pr.MergedAt).
			Returning("*")
		search.Apply(q)

		res, err := q.Update()
		if err != nil {
			if errors.Is(err, pg.ErrNoRows) {
				return nil
			}
			return err
		}
		if res.RowsAffected() == 0 {
			return nil
		}
		found = true

		if err = replacePRReviewers(ctx, tx, model.ID, pr.AssignedReviewers); err != nil {
			return err
		}

		result := toPullRequest(model, pr.AuthorUserID, pr.Status, pr.AssignedReviewers)
		updated = &result
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("UpdatePR: %w", err)
	}
	if !found {
		return nil, nil
	}

	return updated, nil
}

func (r *PrReviewerServiceRepo) GetByReviewerID(ctx context.Context, userID string) ([]PullRequest, error) {
	reviewer := new(PrSystemUser)
	reviewerSearch := PrSystemUserSearch{UserID: &userID}
	qUser := r.db.ModelContext(ctx, reviewer).Column(Columns.PrSystemUser.ID)
	reviewerSearch.Apply(qUser)
	if err := qUser.Select(); err != nil {
		if errors.Is(err, pg.ErrNoRows) {
			return []PullRequest{}, nil
		}
		return nil, fmt.Errorf("GetByReviewerID reviewer lookup: %w", err)
	}

	search := PrSystemPrReviewerSearch{ReviewerID: &reviewer.ID}
	var rels []PrSystemPrReviewer
	q := r.db.ModelContext(ctx, &rels).
		Relation(Columns.PrSystemPrReviewer.Pr).
		Relation(Columns.PrSystemPrReviewer.Pr + "." + Columns.PrSystemPullRequest.Author).
		Relation(Columns.PrSystemPrReviewer.Pr + "." + Columns.PrSystemPullRequest.Status)
	search.Apply(q)
	if err := q.Select(); err != nil {
		if errors.Is(err, pg.ErrNoRows) {
			return []PullRequest{}, nil
		}
		return nil, fmt.Errorf("GetByReviewerID: %w", err)
	}

	result := make([]PullRequest, 0, len(rels))
	for _, rel := range rels {
		if rel.Pr == nil {
			continue
		}
		authorUserID := ""
		if rel.Pr.Author != nil {
			authorUserID = rel.Pr.Author.UserID
		}
		status := StatusIDToName(rel.Pr.StatusID)
		if rel.Pr.Status != nil && rel.Pr.Status.Name != "" {
			status = PRStatus(rel.Pr.Status.Name)
		}
		prItem := toPullRequest(rel.Pr, authorUserID, status, nil)
		result = append(result, prItem)
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].CreatedAt.After(result[j].CreatedAt)
	})

	return result, nil
}

func getReviewerIDsByPR(ctx context.Context, db *pg.DB, prID int64) ([]string, error) {
	var rels []PrSystemPrReviewer
	search := PrSystemPrReviewerSearch{PrID: &prID}

	q := db.ModelContext(ctx, &rels).Relation(Columns.PrSystemPrReviewer.Reviewer)
	search.Apply(q)
	if err := q.Select(); err != nil {
		if errors.Is(err, pg.ErrNoRows) {
			return []string{}, nil
		}
		return nil, err
	}

	reviewers := make([]string, 0, len(rels))
	for _, rel := range rels {
		if rel.Reviewer != nil {
			reviewers = append(reviewers, rel.Reviewer.UserID)
		}
	}
	return reviewers, nil
}

func (r *PrReviewerServiceRepo) GetTotalPRs(ctx context.Context) (int, error) {
	return r.db.ModelContext(ctx, (*PrSystemPullRequest)(nil)).Count()
}

func (r *PrReviewerServiceRepo) GetTotalUsers(ctx context.Context) (int, error) {
	return r.db.ModelContext(ctx, (*PrSystemUser)(nil)).Count()
}

func (r *PrReviewerServiceRepo) GetActiveUsers(ctx context.Context) (int, error) {
	return r.db.ModelContext(ctx, (*PrSystemUser)(nil)).
		Where("? = ?", pg.Ident(Columns.PrSystemUser.IsActive), true).
		Count()
}

func (r *PrReviewerServiceRepo) GetPRsByStatus(ctx context.Context) (map[string]int, error) {
	var statuses []PrSystemStatus
	if err := r.db.ModelContext(ctx, &statuses).Select(); err != nil {
		return nil, err
	}

	result := make(map[string]int, len(statuses))
	for _, st := range statuses {
		count, err := r.db.ModelContext(ctx, (*PrSystemPullRequest)(nil)).
			Where("? = ?", pg.Ident(Columns.PrSystemPullRequest.StatusID), st.ID).
			Count()
		if err != nil {
			return nil, err
		}
		result[st.Name] = count
	}

	return result, nil
}

func (r *PrReviewerServiceRepo) GetTopReviewers(ctx context.Context, limit int) ([]ReviewerStats, error) {
	users, err := r.loadActiveUsers(ctx)
	if err != nil {
		return nil, err
	}

	assignments, err := r.loadReviewerAssignments(ctx)
	if err != nil {
		return nil, err
	}

	byReviewer := initReviewerStats(users)
	aggregateReviewerAssignments(byReviewer, assignments)

	result := flattenReviewerStats(byReviewer)
	sortReviewerStats(result)
	return applyReviewerLimit(result, limit), nil
}

func (r *PrReviewerServiceRepo) loadActiveUsers(ctx context.Context) ([]PrSystemUser, error) {
	var users []PrSystemUser
	if err := r.db.ModelContext(ctx, &users).
		Where("? = ?", pg.Ident(Columns.PrSystemUser.IsActive), true).
		Select(); err != nil {
		return nil, err
	}
	return users, nil
}

func (r *PrReviewerServiceRepo) loadReviewerAssignments(ctx context.Context) ([]PrSystemPrReviewer, error) {
	var assignments []PrSystemPrReviewer
	if err := r.db.ModelContext(ctx, &assignments).
		Relation(Columns.PrSystemPrReviewer.Pr).
		Where("? IS NOT NULL", pg.Ident(Columns.PrSystemPrReviewer.ReviewerID)).
		Select(); err != nil {
		return nil, err
	}
	return assignments, nil
}

func initReviewerStats(users []PrSystemUser) map[int64]ReviewerStats {
	byReviewer := make(map[int64]ReviewerStats, len(users))
	for _, u := range users {
		byReviewer[u.ID] = ReviewerStats{UserID: u.UserID, Username: u.Username}
	}
	return byReviewer
}

func aggregateReviewerAssignments(byReviewer map[int64]ReviewerStats, assignments []PrSystemPrReviewer) {
	distinctPRs := make(map[int64]map[int64]struct{})
	for _, a := range assignments {
		stat, ok := byReviewer[a.ReviewerID]
		if !ok || a.Pr == nil {
			continue
		}
		if distinctPRs[a.ReviewerID] == nil {
			distinctPRs[a.ReviewerID] = map[int64]struct{}{}
		}
		if _, exists := distinctPRs[a.ReviewerID][a.PrID]; exists {
			continue
		}
		distinctPRs[a.ReviewerID][a.PrID] = struct{}{}

		stat.AssignedCount++
		switch StatusIDToName(a.Pr.StatusID) {
		case PRStatusMerged:
			stat.CompletedCount++
		case PRStatusOpen:
			stat.ActiveCount++
		}
		byReviewer[a.ReviewerID] = stat
	}
}

func flattenReviewerStats(byReviewer map[int64]ReviewerStats) []ReviewerStats {
	result := make([]ReviewerStats, 0, len(byReviewer))
	for _, stat := range byReviewer {
		if stat.AssignedCount == 0 {
			continue
		}
		result = append(result, stat)
	}
	return result
}

func sortReviewerStats(result []ReviewerStats) {
	sort.Slice(result, func(i, j int) bool {
		if result[i].AssignedCount == result[j].AssignedCount {
			return result[i].UserID < result[j].UserID
		}
		return result[i].AssignedCount > result[j].AssignedCount
	})
}

func applyReviewerLimit(result []ReviewerStats, limit int) []ReviewerStats {
	if limit > 0 && len(result) > limit {
		return result[:limit]
	}
	return result
}

func (r *PrReviewerServiceRepo) CreateTeam(ctx context.Context, team *Team) (*Team, error) {
	teamModel := &PrSystemTeam{Name: team.TeamName}
	if validateErrs, ok := teamModel.Validate(); !ok {
		return nil, fmt.Errorf("CreateTeam: validate team: %v", validateErrs)
	}

	err := r.db.RunInTransaction(ctx, func(tx *pg.Tx) error {
		if _, err := tx.ModelContext(ctx, teamModel).Returning("*").Insert(); err != nil {
			return err
		}

		for _, m := range team.Members {
			search := PrSystemUserSearch{UserID: &m.UserID}
			q := tx.ModelContext(ctx, (*PrSystemUser)(nil)).
				Set("? = ?", pg.Ident(Columns.PrSystemUser.TeamID), teamModel.ID)
			search.Apply(q)
			res, err := q.Update()
			if err != nil {
				return err
			}
			if res.RowsAffected() == 0 {
				return fmt.Errorf("CreateTeam: user not found: %s", m.UserID)
			}
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("CreateTeam: %w", err)
	}

	members, err := fetchTeamMembers(ctx, r.db, teamModel.ID)
	if err != nil {
		return nil, fmt.Errorf("CreateTeam fetch members: %w", err)
	}

	out := toTeam(teamModel, members)
	return &out, nil
}

func (r *PrReviewerServiceRepo) GetByName(ctx context.Context, teamName string) (*Team, error) {
	teamModel := new(PrSystemTeam)
	search := PrSystemTeamSearch{Name: &teamName}

	q := r.db.ModelContext(ctx, teamModel)
	search.Apply(q)
	if err := q.Select(); err != nil {
		if errors.Is(err, pg.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("GetByName: %w", err)
	}

	members, err := fetchTeamMembers(ctx, r.db, teamModel.ID)
	if err != nil {
		return nil, fmt.Errorf("GetByName fetch members: %w", err)
	}

	out := toTeam(teamModel, members)
	return &out, nil
}

func (r *PrReviewerServiceRepo) ExistsByName(ctx context.Context, teamName string) (bool, error) {
	count, err := r.db.ModelContext(ctx, (*PrSystemTeam)(nil)).
		Where("? = ?", pg.Ident(Columns.PrSystemTeam.Name), teamName).
		Count()
	if err != nil {
		return false, fmt.Errorf("ExistsByName: %w", err)
	}
	return count > 0, nil
}

func fetchTeamMembers(ctx context.Context, db *pg.DB, teamID int64) ([]PrSystemUser, error) {
	var members []PrSystemUser
	search := PrSystemUserSearch{TeamID: &teamID}

	query := db.ModelContext(ctx, &members)
	search.Apply(query)
	if err := query.Select(); err != nil {
		if errors.Is(err, pg.ErrNoRows) {
			return []PrSystemUser{}, nil
		}
		return nil, err
	}

	return members, nil
}

func (r *PrReviewerServiceRepo) GetByUserID(ctx context.Context, userID string) (*User, error) {
	model := new(PrSystemUser)
	search := PrSystemUserSearch{UserID: &userID}

	q := r.db.ModelContext(ctx, model).Relation(Columns.PrSystemUser.Team)
	search.Apply(q)
	if err := q.Select(); err != nil {
		if errors.Is(err, pg.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("GetByUserID: %w", err)
	}

	out := toUser(model)
	return &out, nil
}

func (r *PrReviewerServiceRepo) GetByTeamID(ctx context.Context, teamID int64) ([]User, error) {
	var models []PrSystemUser
	search := PrSystemUserSearch{TeamID: &teamID}

	q := r.db.ModelContext(ctx, &models).Relation(Columns.PrSystemUser.Team)
	search.Apply(q)
	if err := q.Select(); err != nil {
		if errors.Is(err, pg.ErrNoRows) {
			return []User{}, nil
		}
		return nil, fmt.Errorf("GetByTeamID: %w", err)
	}

	return toUsers(models), nil
}

func (r *PrReviewerServiceRepo) SetIsActive(ctx context.Context, userID string, isActive bool) (*User, error) {
	model := new(PrSystemUser)
	search := PrSystemUserSearch{UserID: &userID}

	q := r.db.ModelContext(ctx, model).
		Set("? = ?", pg.Ident(Columns.PrSystemUser.IsActive), isActive).
		Returning("*")
	search.Apply(q)

	res, err := q.Update()
	if err != nil {
		if errors.Is(err, pg.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("SetIsActive: %w", err)
	}
	if res.RowsAffected() == 0 {
		return nil, nil
	}

	out := toUser(model)
	return &out, nil
}

func resolveUserInternalID(ctx context.Context, tx *pg.Tx, userID string) (int64, error) {
	m := new(PrSystemUser)
	search := PrSystemUserSearch{UserID: &userID}

	q := tx.ModelContext(ctx, m).Column(Columns.PrSystemUser.ID)
	search.Apply(q)
	if err := q.Select(); err != nil {
		return 0, err
	}
	return m.ID, nil
}

func resolveStatusID(ctx context.Context, tx *pg.Tx, status PRStatus) (int, error) {
	name := string(status)
	m := new(PrSystemStatus)
	search := PrSystemStatusSearch{Name: &name}

	q := tx.ModelContext(ctx, m).Column(Columns.PrSystemStatus.ID)
	search.Apply(q)
	if err := q.Select(); err != nil {
		return 0, err
	}
	return m.ID, nil
}

func replacePRReviewers(ctx context.Context, tx *pg.Tx, prID int64, reviewerUserIDs []string) error {
	search := PrSystemPrReviewerSearch{PrID: &prID}
	deleteQuery := tx.ModelContext(ctx, (*PrSystemPrReviewer)(nil))
	search.Apply(deleteQuery)
	if _, err := deleteQuery.Delete(); err != nil {
		return err
	}

	for _, reviewerUserID := range reviewerUserIDs {
		reviewerInternalID, err := resolveUserInternalID(ctx, tx, reviewerUserID)
		if err != nil {
			return err
		}

		row := &PrSystemPrReviewer{
			PrID:       prID,
			ReviewerID: reviewerInternalID,
		}
		if _, err = tx.ModelContext(ctx, row).Insert(); err != nil {
			return err
		}
	}
	return nil
}
