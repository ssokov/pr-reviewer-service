package user

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/ssokov/pr-reviewer-service/internal/apperror"
	"github.com/ssokov/pr-reviewer-service/internal/model/domain"
	"github.com/ssokov/pr-reviewer-service/internal/model/dto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/vmkteam/embedlog"
)

type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) SetIsActive(ctx context.Context, userID string, isActive bool) (*domain.User, error) {
	args := m.Called(ctx, userID, isActive)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserService) GetReview(ctx context.Context, userID string) ([]domain.PullRequest, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.PullRequest), args.Error(1)
}

func TestSetIsActive_Success(t *testing.T) {
	e := echo.New()
	mockService := new(MockUserService)
	logger := embedlog.NewLogger(false, false)
	handler := NewHandler(mockService, logger)

	reqBody := dto.SetIsActiveRequest{UserID: "u1", IsActive: false}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/users/setIsActive", bytes.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	expectedUser := &domain.User{UserID: "u1", Username: "Alice", IsActive: false}
	mockService.On("SetIsActive", mock.Anything, "u1", false).Return(expectedUser, nil)

	err := handler.SetIsActive(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	mockService.AssertExpectations(t)
}

func TestSetIsActive_UserNotFound(t *testing.T) {
	e := echo.New()
	mockService := new(MockUserService)
	logger := embedlog.NewLogger(false, false)
	handler := NewHandler(mockService, logger)

	reqBody := dto.SetIsActiveRequest{UserID: "unknown", IsActive: false}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/users/setIsActive", bytes.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	mockService.On("SetIsActive", mock.Anything, "unknown", false).Return(nil, apperror.NewUserNotFoundError("unknown"))

	err := handler.SetIsActive(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestGetReview_Success(t *testing.T) {
	e := echo.New()
	mockService := new(MockUserService)
	logger := embedlog.NewLogger(false, false)
	handler := NewHandler(mockService, logger)

	req := httptest.NewRequest(http.MethodGet, "/users/getReview?user_id=u1", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	expectedPRs := []domain.PullRequest{
		{PullRequestID: "pr-1", PullRequestName: "Test", AuthorID: "u2", Status: domain.PRStatusOpen},
	}
	mockService.On("GetReview", mock.Anything, "u1").Return(expectedPRs, nil)

	err := handler.GetReview(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	mockService.AssertExpectations(t)
}

func TestGetReview_MissingParam(t *testing.T) {
	e := echo.New()
	mockService := new(MockUserService)
	logger := embedlog.NewLogger(false, false)
	handler := NewHandler(mockService, logger)

	req := httptest.NewRequest(http.MethodGet, "/users/getReview", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := handler.GetReview(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}
