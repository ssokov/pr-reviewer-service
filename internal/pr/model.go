package pr

import (
	"time"
)

//go:generate colgen -imports=github.com/ssokov/pr-reviewer-service/internal/db
//colgen:PullRequest:MapP(db)
//colgen:User:MapP(db)
//colgen:Team:MapP(db)
//colgen:ReviewerStats:Map(db)

type PRStatus string

const (
	PRStatusOpen   PRStatus = "OPEN"
	PRStatusMerged PRStatus = "MERGED"
)

type PullRequest struct {
	ID                int64
	PullRequestID     string
	PullRequestName   string
	AuthorID          string
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

type UserStatsItem struct {
	UserID         string
	Username       string
	AssignedCount  int
	CompletedCount int
	ActiveCount    int
}

type PRStatsItem struct {
	Status string
	Count  int
}

type StatsResponse struct {
	TotalPRs     int
	TotalUsers   int
	ActiveUsers  int
	PRsByStatus  []PRStatsItem
	TopReviewers []UserStatsItem
}

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
