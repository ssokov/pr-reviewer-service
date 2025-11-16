package team

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/ssokov/pr-reviewer-service/internal/http/response"
	"github.com/ssokov/pr-reviewer-service/internal/model/dto"
)

// DeactivateTeam godoc
// @Summary Deactivate team members
// @Description Deactivate all active members of a team and return their open PRs for reassignment

// @Tags team
// @Accept json
// @Produce json
// @Param request body dto.DeactivateTeamRequest true "Team deactivation request"
// @Success 200 {object} dto.DeactivateTeamResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse "Team not found"
// @Failure 500 {object} dto.ErrorResponse
// @Router /team/deactivate [post]
func (t *TeamHandler) DeactivateTeam(c echo.Context) error {
	ctx := c.Request().Context()

	var req dto.DeactivateTeamRequest
	if err := c.Bind(&req); err != nil {
		t.logger.Errorf("failed to bind request: %v", err)
		return response.Error(c, http.StatusBadRequest, "INVALID_INPUT", "invalid request body")
	}

	if req.TeamName == "" {
		return response.Error(c, http.StatusBadRequest, "INVALID_INPUT", "team_name is required")
	}

	deactivatedUsers, openPRs, err := t.teamService.DeactivateTeam(ctx, req.TeamName)
	if err != nil {
		t.logger.Errorf("failed to deactivate team: %v", err)
		return response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to deactivate team")
	}

	usersInfo := make([]dto.DeactivatedUserInfo, len(deactivatedUsers))
	userPRCount := make(map[string]int)

	for _, pr := range openPRs {
		if pr.AssignedReviewers != nil {
			for _, reviewerID := range pr.AssignedReviewers {
				userPRCount[reviewerID]++
			}
		}
	}

	for i, user := range deactivatedUsers {
		usersInfo[i] = dto.DeactivatedUserInfo{
			UserID:       user.UserID,
			Username:     user.Username,
			OpenPRsCount: userPRCount[user.UserID],
		}
	}

	resp := dto.DeactivateTeamResponse{
		DeactivatedUsers: len(deactivatedUsers),
		ReassignedPRs:    len(openPRs),
		Users:            usersInfo,
	}

	return c.JSON(http.StatusOK, resp)
}
