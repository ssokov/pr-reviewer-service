package apperror

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
