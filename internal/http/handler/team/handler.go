package team

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/ssokov/pr-reviewer-service/internal/http/mapper"
	"github.com/ssokov/pr-reviewer-service/internal/http/response"
	"github.com/ssokov/pr-reviewer-service/internal/model/dto"
	"github.com/ssokov/pr-reviewer-service/internal/service"
	"github.com/vmkteam/embedlog"
)

type TeamHandler struct {
	teamService service.TeamService
	logger      embedlog.Logger
}

func NewHandler(service service.TeamService, logger embedlog.Logger) *TeamHandler {
	return &TeamHandler{
		teamService: service,
		logger:      logger,
	}
}

// AddTeam godoc
// @Summary Add a new team
// @Description Create a new team with members
// @Tags team
// @Accept json
// @Produce json
// @Param request body dto.AddTeamRequest true "Team data"
// @Success 201 {object} dto.AddTeamResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 409 {object} dto.ErrorResponse "Team already exists"
// @Failure 500 {object} dto.ErrorResponse
// @Router /team/add [post]
func (t *TeamHandler) AddTeam(c echo.Context) error {
	var req dto.AddTeamRequest
	if err := c.Bind(&req); err != nil {
		t.logger.Errorf("failed to bind request: %v", err)
		return response.Error(c, http.StatusBadRequest, "INVALID_INPUT", "invalid request body")
	}

	domainTeam := mapper.AddTeamRequestToDomain(req)

	ctx := c.Request().Context()
	result, err := t.teamService.AddTeam(ctx, domainTeam)
	if err != nil {
		t.logger.Errorf("failed to add team: %v", err)
		return response.HandleError(c, err)
	}

	return c.JSON(http.StatusCreated, dto.AddTeamResponse{
		Team: mapper.TeamToResponse(result),
	})
}

// GetTeam godoc
// @Summary Get team by name
// @Description Get team information by team name
// @Tags team
// @Accept json
// @Produce json
// @Param team_name query string true "Team name"
// @Success 200 {object} dto.TeamResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse "Team not found"
// @Failure 500 {object} dto.ErrorResponse
// @Router /team/get [get]
func (t *TeamHandler) GetTeam(c echo.Context) error {
	teamName := c.QueryParam("team_name")
	if teamName == "" {
		t.logger.Errorf("failed to get team: team_name is required")
		return response.Error(c, http.StatusBadRequest, "INVALID_INPUT", "team_name is required")
	}

	ctx := c.Request().Context()
	team, err := t.teamService.GetTeam(ctx, teamName)
	if err != nil {
		t.logger.Errorf("failed to get team: %v", err)
		return response.HandleError(c, err)
	}

	return c.JSON(http.StatusOK, mapper.TeamToResponse(team))
}
