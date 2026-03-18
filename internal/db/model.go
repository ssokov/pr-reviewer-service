package db

import "time"

type PRStatus string

const (
	PRStatusOpen   PRStatus = "OPEN"
	PRStatusMerged PRStatus = "MERGED"
)

type Team struct {
	ID        int64
	TeamName  string
	Members   []User
	CreatedAt time.Time
}

type User struct {
	ID        int64
	UserID    string
	Username  string
	TeamID    int64
	TeamName  string
	IsActive  bool
	CreatedAt time.Time
}

type PullRequest struct {
	ID                int64
	PullRequestID     string
	PullRequestName   string
	AuthorID          int64
	AuthorUserID      string
	Status            PRStatus
	AssignedReviewers []string
	CreatedAt         time.Time
	MergedAt          *time.Time
}

type ReviewerStats struct {
	UserID         string
	Username       string
	AssignedCount  int
	CompletedCount int
	ActiveCount    int
}
