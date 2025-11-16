package dto

import "time"

type CreatePRRequest struct {
	PullRequestID   string `json:"pull_request_id" validate:"required"`
	PullRequestName string `json:"pull_request_name" validate:"required"`
	AuthorID        string `json:"author_id" validate:"required"`
}

type MergePRRequest struct {
	PullRequestID string `json:"pull_request_id" validate:"required"`
}

type ReassignRequest struct {
	PullRequestID string `json:"pull_request_id" validate:"required"`
	OldUserID     string `json:"old_user_id" validate:"required"`
}

type PullRequestResponse struct {
	PullRequestID     string     `json:"pull_request_id"`
	PullRequestName   string     `json:"pull_request_name"`
	AuthorID          string     `json:"author_id"`
	Status            string     `json:"status"`
	AssignedReviewers []string   `json:"assigned_reviewers"`
	CreatedAt         *time.Time `json:"createdAt,omitempty"`
	MergedAt          *time.Time `json:"mergedAt,omitempty"`
}

type CreatePRResponse struct {
	PR PullRequestResponse `json:"pr"`
}

type MergePRResponse struct {
	PR PullRequestResponse `json:"pr"`
}

type ReassignResponse struct {
	PR         PullRequestResponse `json:"pr"`
	ReplacedBy string              `json:"replaced_by"`
}
