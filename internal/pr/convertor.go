package pr

import "github.com/ssokov/pr-reviewer-service/internal/db"

func NewPullRequest(in *db.PullRequest) *PullRequest {
	return &PullRequest{
		ID:                in.ID,
		PullRequestID:     in.PullRequestID,
		PullRequestName:   in.PullRequestName,
		AuthorID:          in.AuthorUserID,
		Status:            PRStatus(in.Status),
		AssignedReviewers: append([]string(nil), in.AssignedReviewers...),
		CreatedAt:         in.CreatedAt,
		MergedAt:          in.MergedAt,
	}
}

func NewReviewerStats(in db.ReviewerStats) ReviewerStats {
	return ReviewerStats{
		UserID:         in.UserID,
		Username:       in.Username,
		AssignedCount:  in.AssignedCount,
		CompletedCount: in.CompletedCount,
		ActiveCount:    in.ActiveCount,
	}
}

func NewTeam(in *db.Team) *Team {
	return &Team{
		ID:        in.ID,
		TeamName:  in.TeamName,
		Members:   NewUsers(in.Members),
		CreatedAt: in.CreatedAt,
	}
}

func NewUser(in *db.User) *User {
	return &User{
		ID:        in.ID,
		UserID:    in.UserID,
		Username:  in.Username,
		TeamID:    in.TeamID,
		TeamName:  in.TeamName,
		IsActive:  in.IsActive,
		CreatedAt: in.CreatedAt,
	}
}

func NewDBPullRequest(in *PullRequest) *db.PullRequest {
	return &db.PullRequest{
		ID:                in.ID,
		PullRequestID:     in.PullRequestID,
		PullRequestName:   in.PullRequestName,
		AuthorUserID:      in.AuthorID,
		Status:            db.PRStatus(in.Status),
		AssignedReviewers: append([]string(nil), in.AssignedReviewers...),
		CreatedAt:         in.CreatedAt,
		MergedAt:          in.MergedAt,
	}
}

func NewDBUser(in User) db.User {
	return db.User{
		ID:        in.ID,
		UserID:    in.UserID,
		Username:  in.Username,
		TeamID:    in.TeamID,
		TeamName:  in.TeamName,
		IsActive:  in.IsActive,
		CreatedAt: in.CreatedAt,
	}
}

func NewDBTeam(in *Team) *db.Team {
	return &db.Team{
		ID:        in.ID,
		TeamName:  in.TeamName,
		Members:   Map(in.Members, NewDBUser),
		CreatedAt: in.CreatedAt,
	}
}

// MapP converts slice of type T to slice of type M using pointer converter.
func MapP[T, M any](a []T, f func(*T) *M) []M {
	n := make([]M, len(a))
	for i := range a {
		n[i] = *f(&a[i])
	}
	return n
}

// Map converts slice of type T to slice of type M using value converter.
func Map[T, M any](a []T, f func(T) M) []M {
	n := make([]M, len(a))
	for i := range a {
		n[i] = f(a[i])
	}
	return n
}
