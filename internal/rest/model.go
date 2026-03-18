package rest

import (
	"errors"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/ssokov/pr-reviewer-service/internal/pr"
)

//go:generate colgen -imports=github.com/ssokov/pr-reviewer-service/internal/pr
//colgen:TeamMember:Map(pr.User)
//colgen:PullRequestShort:Map(pr.PullRequest)

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

type UserStatsItem struct {
	UserID         string `json:"user_id"`
	Username       string `json:"username"`
	AssignedCount  int    `json:"assigned_count"`
	CompletedCount int    `json:"completed_count"`
	ActiveCount    int    `json:"active_count"`
}

type PRStatsItem struct {
	Status string `json:"status"`
	Count  int    `json:"count"`
}

type StatsResponse struct {
	TotalPRs     int             `json:"total_prs"`
	TotalUsers   int             `json:"total_users"`
	ActiveUsers  int             `json:"active_users"`
	PRsByStatus  []PRStatsItem   `json:"prs_by_status"`
	TopReviewers []UserStatsItem `json:"top_reviewers"`
}

type TeamMember struct {
	UserID   string `json:"user_id" validate:"required"`
	Username string `json:"username" validate:"required"`
	IsActive bool   `json:"is_active"`
}

type AddTeamRequest struct {
	TeamName string       `json:"team_name" validate:"required"`
	Members  []TeamMember `json:"members" validate:"required,min=1"`
}

type TeamResponse struct {
	TeamName string       `json:"team_name"`
	Members  []TeamMember `json:"members"`
}

type AddTeamResponse struct {
	Team TeamResponse `json:"team"`
}

type SetIsActiveRequest struct {
	UserID   string `json:"user_id" validate:"required"`
	IsActive bool   `json:"is_active"`
}

type UserResponse struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	TeamName string `json:"team_name"`
	IsActive bool   `json:"is_active"`
}

type SetIsActiveResponse struct {
	User UserResponse `json:"user"`
}

type PullRequestShort struct {
	PullRequestID   string `json:"pull_request_id"`
	PullRequestName string `json:"pull_request_name"`
	AuthorID        string `json:"author_id"`
	Status          string `json:"status"`
}

type GetReviewResponse struct {
	UserID       string             `json:"user_id"`
	PullRequests []PullRequestShort `json:"pull_requests"`
}

type ErrorDetail struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type ErrorResponse struct {
	Error ErrorDetail `json:"error"`
}

type ErrorCode string

const (
	ErrCodeTeamExists   ErrorCode = "TEAM_EXISTS"
	ErrCodeTeamNotFound ErrorCode = "TEAM_NOT_FOUND"

	ErrCodeUserNotFound ErrorCode = "USER_NOT_FOUND"

	ErrCodePRExists    ErrorCode = "PR_EXISTS"
	ErrCodePRNotFound  ErrorCode = "PR_NOT_FOUND"
	ErrCodePRMerged    ErrorCode = "PR_MERGED"
	ErrCodeNotAssigned ErrorCode = "NOT_ASSIGNED"
	ErrCodeNoCandidate ErrorCode = "NO_CANDIDATE"

	ErrCodeNotFound      ErrorCode = "NOT_FOUND"
	ErrCodeInvalidInput  ErrorCode = "INVALID_INPUT"
	ErrCodeInternalError ErrorCode = "INTERNAL_ERROR"
)

var errorStatusMap = map[ErrorCode]int{
	ErrCodeTeamExists:   http.StatusBadRequest,
	ErrCodeInvalidInput: http.StatusBadRequest,
	ErrCodePRExists:     http.StatusConflict,
	ErrCodePRMerged:     http.StatusConflict,
	ErrCodeNotAssigned:  http.StatusConflict,
	ErrCodeNoCandidate:  http.StatusConflict,
	ErrCodeTeamNotFound: http.StatusNotFound,
	ErrCodeUserNotFound: http.StatusNotFound,
	ErrCodePRNotFound:   http.StatusNotFound,
	ErrCodeNotFound:     http.StatusNotFound,
}

func HandleError(c echo.Context, err error) error {
	var appErr *pr.AppError
	if !errors.As(err, &appErr) {
		return Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "internal server error")
	}

	statusCode := getHTTPStatus(ErrorCode(appErr.Code))
	return Error(c, statusCode, string(appErr.Code), appErr.Message)
}

func getHTTPStatus(code ErrorCode) int {
	if status, ok := errorStatusMap[code]; ok {
		return status
	}
	return http.StatusInternalServerError
}

func Error(c echo.Context, status int, code, message string) error {
	return c.JSON(status, NewErrorResponse(code, message))
}

func NewErrorResponse(code, message string) ErrorResponse {
	return ErrorResponse{
		Error: ErrorDetail{
			Code:    code,
			Message: message,
		},
	}
}
