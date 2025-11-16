package pr

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

type MockPRService struct {
	mock.Mock
}

func (m *MockPRService) CreatePR(ctx context.Context, authorID string, pr *domain.PullRequest) (*domain.PullRequest, error) {
	args := m.Called(ctx, authorID, pr)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.PullRequest), args.Error(1)
}

func (m *MockPRService) MergePR(ctx context.Context, prID string) (*domain.PullRequest, error) {
	args := m.Called(ctx, prID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.PullRequest), args.Error(1)
}

func (m *MockPRService) ReassignReviewer(ctx context.Context, prID, oldUserID string) (*domain.PullRequest, string, error) {
	args := m.Called(ctx, prID, oldUserID)
	if args.Get(0) == nil {
		return nil, "", args.Error(2)
	}
	return args.Get(0).(*domain.PullRequest), args.String(1), args.Error(2)
}

func TestCreatePR_Success(t *testing.T) {
	e := echo.New()
	mockService := new(MockPRService)
	logger := embedlog.NewLogger(false, false)
	handler := NewHandler(mockService, logger)

	reqBody := dto.CreatePRRequest{
		PullRequestID:   "pr-1",
		PullRequestName: "Test PR",
		AuthorID:        "u1",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/pullRequest/create", bytes.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	expectedPR := &domain.PullRequest{
		PullRequestID:     "pr-1",
		PullRequestName:   "Test PR",
		AuthorID:          "u1",
		Status:            domain.PRStatusOpen,
		AssignedReviewers: []string{"u2"},
	}
	mockService.On("CreatePR", mock.Anything, "u1", mock.Anything).Return(expectedPR, nil)

	err := handler.CreatePR(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, rec.Code)
	mockService.AssertExpectations(t)
}

func TestMergePR_Success(t *testing.T) {
	e := echo.New()
	mockService := new(MockPRService)
	logger := embedlog.NewLogger(false, false)
	handler := NewHandler(mockService, logger)

	reqBody := dto.MergePRRequest{PullRequestID: "pr-1"}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/pullRequest/merge", bytes.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	expectedPR := &domain.PullRequest{
		PullRequestID: "pr-1",
		Status:        domain.PRStatusMerged,
	}
	mockService.On("MergePR", mock.Anything, "pr-1").Return(expectedPR, nil)

	err := handler.MergePR(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestReassignReviewer_Success(t *testing.T) {
	e := echo.New()
	mockService := new(MockPRService)
	logger := embedlog.NewLogger(false, false)
	handler := NewHandler(mockService, logger)

	reqBody := dto.ReassignRequest{
		PullRequestID: "pr-1",
		OldUserID:     "u2",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/pullRequest/reassign", bytes.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	expectedPR := &domain.PullRequest{
		PullRequestID:     "pr-1",
		AssignedReviewers: []string{"u3"},
	}
	mockService.On("ReassignReviewer", mock.Anything, "pr-1", "u2").Return(expectedPR, "u3", nil)

	err := handler.ReassignReviewer(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestReassignReviewer_PRMerged(t *testing.T) {
	e := echo.New()
	mockService := new(MockPRService)
	logger := embedlog.NewLogger(false, false)
	handler := NewHandler(mockService, logger)

	reqBody := dto.ReassignRequest{PullRequestID: "pr-1", OldUserID: "u2"}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/pullRequest/reassign", bytes.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	mockService.On("ReassignReviewer", mock.Anything, "pr-1", "u2").Return(nil, "", apperror.NewPRMergedError("pr-1"))

	err := handler.ReassignReviewer(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusConflict, rec.Code)
}
