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
	t.Run("given malformed json when CreatePR then bad_request is returned", func(t *testing.T) {
		// Arrange
		handler := NewHandler(nil, embedlog.NewLogger(false, false))
		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/pullRequest/create", bytes.NewBufferString("{"))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)

		// Act
		err := handler.CreatePR(ctx)

		// Assert
		require.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})
}

func TestSetIsActive_InvalidJSON(t *testing.T) {
	t.Run("given malformed json when SetIsActive then bad_request is returned", func(t *testing.T) {
		// Arrange
		handler := NewHandler(nil, embedlog.NewLogger(false, false))
		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/user/setIsActive", bytes.NewBufferString("{"))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)

		// Act
		err := handler.SetIsActive(ctx)

		// Assert
		require.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})
}

func TestGetReview_MissingParam(t *testing.T) {
	t.Run("given missing user_id when GetReview then bad_request is returned", func(t *testing.T) {
		// Arrange
		handler := NewHandler(nil, embedlog.NewLogger(false, false))
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/user/getReview", nil)
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)

		// Act
		err := handler.GetReview(ctx)

		// Assert
		require.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})
}

func TestGetTeam_MissingParam(t *testing.T) {
	t.Run("given missing team_name when GetTeam then bad_request is returned", func(t *testing.T) {
		// Arrange
		handler := NewHandler(nil, embedlog.NewLogger(false, false))
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/team/get", nil)
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)

		// Act
		err := handler.GetTeam(ctx)

		// Assert
		require.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})
}
