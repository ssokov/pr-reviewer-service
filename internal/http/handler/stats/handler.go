package stats

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/ssokov/pr-reviewer-service/internal/http/response"
	"github.com/ssokov/pr-reviewer-service/internal/service"
	"github.com/vmkteam/embedlog"
)

type Handler struct {
	statsService service.StatsService
	logger       embedlog.Logger
}

func NewHandler(statsService service.StatsService, logger embedlog.Logger) *Handler {
	return &Handler{
		statsService: statsService,
		logger:       logger,
	}
}

// GetStats godoc
// @Summary Get system statistics
// @Summary Get statistics
// @Description Get system statistics including PR counts, user counts, and top reviewers
// @Tags stats
// @Produce json
// @Success 200 {object} dto.StatsResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /stats [get]
func (h *Handler) GetStats(c echo.Context) error {
	ctx := c.Request().Context()

	stats, err := h.statsService.GetStats(ctx)
	if err != nil {
		h.logger.Print(ctx, "failed to get stats", "error", err)
		return response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to get statistics")
	}

	return c.JSON(http.StatusOK, stats)
}
