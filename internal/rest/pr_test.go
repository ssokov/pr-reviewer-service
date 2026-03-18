package rest

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vmkteam/embedlog"
)

func TestCreatePR_InvalidJSON(t *testing.T) {
	handler := NewHandler(nil, embedlog.NewLogger(false, false))
	e := echo.New()

	req := httptest.NewRequest(http.MethodPost, "/pullRequest/create", bytes.NewBufferString("{"))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	err := handler.CreatePR(e.NewContext(req, rec))
	require.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestSetIsActive_InvalidJSON(t *testing.T) {
	handler := NewHandler(nil, embedlog.NewLogger(false, false))
	e := echo.New()

	req := httptest.NewRequest(http.MethodPost, "/user/setIsActive", bytes.NewBufferString("{"))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	err := handler.SetIsActive(e.NewContext(req, rec))
	require.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestGetReview_MissingParam(t *testing.T) {
	handler := NewHandler(nil, embedlog.NewLogger(false, false))
	e := echo.New()

	req := httptest.NewRequest(http.MethodGet, "/user/getReview", nil)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)

	err := handler.GetReview(ctx)
	require.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestGetTeam_MissingParam(t *testing.T) {
	handler := NewHandler(nil, embedlog.NewLogger(false, false))
	e := echo.New()

	req := httptest.NewRequest(http.MethodGet, "/team/get", nil)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)

	err := handler.GetTeam(ctx)
	require.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}
