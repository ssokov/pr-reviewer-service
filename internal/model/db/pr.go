package db

import "time"

type PullRequest struct {
	ID              int64
	PullRequestID   string
	PullRequestName string
	AuthorID        int64
	StatusID        int
	CreatedAt       time.Time
	MergedAt        *time.Time
}

type PRReviewer struct {
	PRID       int64
	ReviewerID int64
}
