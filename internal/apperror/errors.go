package apperror

import (
	"errors"
	"fmt"
)

type AppError struct {
	Code    ErrorCode
	Message string
	Err     error
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("[%s] %s: %v", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

func (e *AppError) Unwrap() error {
	return e.Err
}

func New(code ErrorCode, message string) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
	}
}

func Wrap(code ErrorCode, message string, err error) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

func Is(err error, code ErrorCode) bool {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr.Code == code
	}
	return false
}

func NewTeamExistsError(teamName string) *AppError {
	return New(ErrCodeTeamExists, fmt.Sprintf("team '%s' already exists", teamName))
}

func NewTeamNotFoundError(teamName string) *AppError {
	return New(ErrCodeTeamNotFound, fmt.Sprintf("team '%s' not found", teamName))
}

func NewUserNotFoundError(userID string) *AppError {
	return New(ErrCodeUserNotFound, fmt.Sprintf("user '%s' not found", userID))
}

func NewPRExistsError(prID string) *AppError {
	return New(ErrCodePRExists, fmt.Sprintf("pull request '%s' already exists", prID))
}

func NewPRNotFoundError(prID string) *AppError {
	return New(ErrCodePRNotFound, fmt.Sprintf("pull request '%s' not found", prID))
}

func NewPRMergedError(prID string) *AppError {
	return New(ErrCodePRMerged, fmt.Sprintf("cannot modify merged pull request '%s'", prID))
}

func NewNotAssignedError(userID, prID string) *AppError {
	return New(ErrCodeNotAssigned, fmt.Sprintf("user '%s' is not assigned to PR '%s'", userID, prID))
}

func NewNoCandidateError(teamName string) *AppError {
	return New(ErrCodeNoCandidate, fmt.Sprintf("no active candidate available in team '%s'", teamName))
}

func NewNotFoundError(resource string) *AppError {
	return New(ErrCodeNotFound, fmt.Sprintf("%s not found", resource))
}

func NewInvalidInputError(message string) *AppError {
	return New(ErrCodeInvalidInput, message)
}

func NewInternalError(message string, err error) *AppError {
	return Wrap(ErrCodeInternalError, message, err)
}
