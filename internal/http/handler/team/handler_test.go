package team

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

type MockTeamService struct {
	mock.Mock
}

func (m *MockTeamService) AddTeam(ctx context.Context, team *domain.Team) (*domain.Team, error) {
	args := m.Called(ctx, team)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Team), args.Error(1)
}

func (m *MockTeamService) GetTeam(ctx context.Context, teamName string) (*domain.Team, error) {
	args := m.Called(ctx, teamName)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Team), args.Error(1)
}

func (m *MockTeamService) DeactivateTeam(ctx context.Context, teamName string) ([]domain.User, []domain.PullRequest, error) {
	args := m.Called(ctx, teamName)
	if args.Get(0) == nil {
		return nil, nil, args.Error(2)
	}
	users := args.Get(0).([]domain.User)
	prs := args.Get(1).([]domain.PullRequest)
	return users, prs, args.Error(2)
}

func TestAddTeam_Success(t *testing.T) {
	e := echo.New()
	mockService := new(MockTeamService)
	logger := embedlog.NewLogger(false, false)
	handler := NewHandler(mockService, logger)

	reqBody := dto.AddTeamRequest{
		TeamName: "backend",
		Members: []dto.TeamMember{
			{UserID: "u1", Username: "Alice", IsActive: true},
		},
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/team/add", bytes.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	expectedTeam := &domain.Team{
		ID:       1,
		TeamName: "backend",
		Members:  []domain.User{{UserID: "u1", Username: "Alice", IsActive: true}},
	}
	mockService.On("AddTeam", mock.Anything, mock.Anything).Return(expectedTeam, nil)

	err := handler.AddTeam(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, rec.Code)
	mockService.AssertExpectations(t)
}

func TestAddTeam_TeamExists(t *testing.T) {
	e := echo.New()
	mockService := new(MockTeamService)
	logger := embedlog.NewLogger(false, false)
	handler := NewHandler(mockService, logger)

	reqBody := dto.AddTeamRequest{TeamName: "backend"}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/team/add", bytes.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	mockService.On("AddTeam", mock.Anything, mock.Anything).Return(nil, apperror.NewTeamExistsError("backend"))

	err := handler.AddTeam(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestGetTeam_Success(t *testing.T) {
	e := echo.New()
	mockService := new(MockTeamService)
	logger := embedlog.NewLogger(false, false)
	handler := NewHandler(mockService, logger)

	req := httptest.NewRequest(http.MethodGet, "/team/get?team_name=backend", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	expectedTeam := &domain.Team{
		ID:       1,
		TeamName: "backend",
		Members:  []domain.User{{UserID: "u1", Username: "Alice"}},
	}
	mockService.On("GetTeam", mock.Anything, "backend").Return(expectedTeam, nil)

	err := handler.GetTeam(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	mockService.AssertExpectations(t)
}

func TestGetTeam_NotFound(t *testing.T) {
	e := echo.New()
	mockService := new(MockTeamService)
	logger := embedlog.NewLogger(false, false)
	handler := NewHandler(mockService, logger)

	req := httptest.NewRequest(http.MethodGet, "/team/get?team_name=unknown", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	mockService.On("GetTeam", mock.Anything, "unknown").Return(nil, apperror.NewTeamNotFoundError("unknown"))

	err := handler.GetTeam(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestGetTeam_MissingParam(t *testing.T) {
	e := echo.New()
	mockService := new(MockTeamService)
	logger := embedlog.NewLogger(false, false)
	handler := NewHandler(mockService, logger)

	req := httptest.NewRequest(http.MethodGet, "/team/get", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := handler.GetTeam(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}
