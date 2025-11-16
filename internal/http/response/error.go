package response

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/ssokov/pr-reviewer-service/internal/apperror"
	"github.com/ssokov/pr-reviewer-service/internal/model/dto"
)

var errorStatusMap = map[apperror.ErrorCode]int{
	apperror.ErrCodeTeamExists:   http.StatusBadRequest,
	apperror.ErrCodeInvalidInput: http.StatusBadRequest,
	apperror.ErrCodePRExists:     http.StatusConflict,
	apperror.ErrCodePRMerged:     http.StatusConflict,
	apperror.ErrCodeNotAssigned:  http.StatusConflict,
	apperror.ErrCodeNoCandidate:  http.StatusConflict,
	apperror.ErrCodeTeamNotFound: http.StatusNotFound,
	apperror.ErrCodeUserNotFound: http.StatusNotFound,
	apperror.ErrCodePRNotFound:   http.StatusNotFound,
	apperror.ErrCodeNotFound:     http.StatusNotFound,
}

func HandleError(c echo.Context, err error) error {
	var appErr *apperror.AppError
	if !errors.As(err, &appErr) {
		return Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "internal server error")
	}

	statusCode := getHTTPStatus(appErr.Code)
	return Error(c, statusCode, string(appErr.Code), appErr.Message)
}

func getHTTPStatus(code apperror.ErrorCode) int {
	if status, ok := errorStatusMap[code]; ok {
		return status
	}
	return http.StatusInternalServerError
}

func Error(c echo.Context, status int, code, message string) error {
	return c.JSON(status, NewErrorResponse(code, message))
}

func NewErrorResponse(code, message string) dto.ErrorResponse {
	return dto.ErrorResponse{
		Error: dto.ErrorDetail{
			Code:    code,
			Message: message,
		},
	}
}
