package user

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/ssokov/pr-reviewer-service/internal/http/mapper"
	"github.com/ssokov/pr-reviewer-service/internal/http/response"
	"github.com/ssokov/pr-reviewer-service/internal/model/dto"
	"github.com/ssokov/pr-reviewer-service/internal/service"
	"github.com/vmkteam/embedlog"
)

type UserHandler struct {
	userService service.UserService
	logger      embedlog.Logger
}

func NewHandler(service service.UserService, logger embedlog.Logger) *UserHandler {
	return &UserHandler{
		userService: service,
		logger:      logger,
	}
}

// SetIsActive godoc
// @Summary Set user active status
// @Description Set whether a user is active or inactive
// @Tags user
// @Accept json
// @Produce json
// @Param request body dto.SetIsActiveRequest true "User active status"
// @Success 200 {object} dto.SetIsActiveResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse "User not found"
// @Failure 500 {object} dto.ErrorResponse
// @Router /user/setIsActive [post]
func (h *UserHandler) SetIsActive(c echo.Context) error {
	var req dto.SetIsActiveRequest
	if err := c.Bind(&req); err != nil {
		h.logger.Errorf("failed to bind request: %v", err)
		return response.Error(c, http.StatusBadRequest, "INVALID_INPUT", "invalid request body")
	}

	ctx := c.Request().Context()
	user, err := h.userService.SetIsActive(ctx, req.UserID, req.IsActive)
	if err != nil {
		h.logger.Errorf("failed to set user active status: %v", err)
		return response.HandleError(c, err)
	}

	return c.JSON(http.StatusOK, dto.SetIsActiveResponse{
		User: mapper.UserToResponse(user),
	})
}

// GetReview godoc
// @Summary Get user's pull requests for review
// @Description Get all pull requests assigned to a user for review
// @Tags user
// @Accept json
// @Produce json
// @Param user_id query string true "User ID"
// @Success 200 {object} dto.GetReviewResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse "User not found"
// @Failure 500 {object} dto.ErrorResponse
// @Router /user/getReview [get]
func (h *UserHandler) GetReview(c echo.Context) error {
	userID := c.QueryParam("user_id")
	if userID == "" {
		h.logger.Errorf("failed to get review: user_id is required")
		return response.Error(c, http.StatusBadRequest, "INVALID_INPUT", "user_id is required")
	}

	ctx := c.Request().Context()
	pullRequests, err := h.userService.GetReview(ctx, userID)
	if err != nil {
		h.logger.Errorf("failed to get review: %v", err)
		return response.HandleError(c, err)
	}

	prShorts := mapper.PullRequestsToShort(pullRequests)

	return c.JSON(http.StatusOK, dto.GetReviewResponse{
		UserID:       userID,
		PullRequests: prShorts,
	})
}
