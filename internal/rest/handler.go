package rest

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/ssokov/pr-reviewer-service/internal/pr"
	"github.com/vmkteam/embedlog"
)

type PRHandler struct {
	prService *pr.PrService
	logger    embedlog.Logger
}

func NewHandler(prService *pr.PrService, logger embedlog.Logger) *PRHandler {
	return &PRHandler{
		prService: prService,
		logger:    logger,
	}
}

func (h *PRHandler) bindRequest(c echo.Context, req any) error {
	if err := c.Bind(req); err != nil {
		h.logger.Errorf("failed to bind request: %v", err)
		return Error(c, http.StatusBadRequest, "INVALID_INPUT", "invalid request body")
	}
	return nil
}

// CreatePR godoc
// @Summary Create a new pull request
// @Description Create a new pull request and automatically assign reviewers
// @Tags pullRequest
// @Accept json
// @Produce json
// @Param request body dto.CreatePRRequest true "Pull request data"
// @Success 201 {object} dto.CreatePRResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse "Author or team not found"
// @Failure 409 {object} dto.ErrorResponse "PR already exists"
// @Failure 500 {object} dto.ErrorResponse
// @Router /pullRequest/create [post]
func (p *PRHandler) CreatePR(c echo.Context) error {
	var req CreatePRRequest
	if err := p.bindRequest(c, &req); err != nil {
		return err
	}

	domainPR := CreatePRRequestToDomain(req)

	ctx := c.Request().Context()
	createdPR, err := p.prService.CreatePR(ctx, req.AuthorID, domainPR)
	if err != nil {
		p.logger.Errorf("failed to create PR: %v", err)
		return HandleError(c, err)
	}

	return c.JSON(http.StatusCreated, CreatePRResponse{
		PR: PullRequestToResponse(createdPR),
	})
}

// MergePR godoc
// @Summary Merge a pull request
// @Description Merge an existing pull request
// @Tags pullRequest
// @Accept json
// @Produce json
// @Param request body dto.MergePRRequest true "Pull request ID"
// @Success 200 {object} dto.MergePRResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse "PR not found"
// @Failure 500 {object} dto.ErrorResponse
// @Router /pullRequest/merge [post]
func (p *PRHandler) MergePR(c echo.Context) error {
	var req MergePRRequest
	if err := p.bindRequest(c, &req); err != nil {
		return err
	}

	ctx := c.Request().Context()
	mergedPR, err := p.prService.MergePR(ctx, req.PullRequestID)
	if err != nil {
		p.logger.Errorf("failed to merge PR: %v", err)
		return HandleError(c, err)
	}

	return c.JSON(http.StatusOK, MergePRResponse{
		PR: PullRequestToResponse(mergedPR),
	})
}

// ReassignReviewer godoc
// @Summary Reassign a reviewer
// @Description Replace a reviewer with another active team member
// @Tags pullRequest
// @Accept json
// @Produce json
// @Param request body dto.ReassignRequest true "Reassign data"
// @Success 200 {object} dto.ReassignResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse "PR or user not found"
// @Failure 409 {object} dto.ErrorResponse "User not assigned to PR"
// @Failure 500 {object} dto.ErrorResponse
// @Router /pullRequest/reassign [post]
func (p *PRHandler) ReassignReviewer(c echo.Context) error {
	var req ReassignRequest
	if err := p.bindRequest(c, &req); err != nil {
		return err
	}

	ctx := c.Request().Context()
	pr, newReviewerID, err := p.prService.ReassignReviewer(ctx, req.PullRequestID, req.OldUserID)
	if err != nil {
		p.logger.Errorf("failed to reassign reviewer: %v", err)
		return HandleError(c, err)
	}

	return c.JSON(http.StatusOK, ReassignResponse{
		PR:         PullRequestToResponse(pr),
		ReplacedBy: newReviewerID,
	})
}

// GetStats godoc
// @Summary Get statistics
// @Description Get system statistics including PR counts, user counts, and top reviewers
// @Tags stats
// @Produce json
// @Success 200 {object} dto.StatsResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /stats [get]
func (h *PRHandler) GetStats(c echo.Context) error {
	ctx := c.Request().Context()

	stats, err := h.prService.GetStats(ctx)
	if err != nil {
		h.logger.Print(ctx, "failed to get stats", "error", err)
		return Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to get statistics")
	}

	return c.JSON(http.StatusOK, stats)
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
func (t *PRHandler) AddTeam(c echo.Context) error {
	var req AddTeamRequest
	if err := t.bindRequest(c, &req); err != nil {
		return err
	}

	domainTeam := AddTeamRequestToDomain(req)

	ctx := c.Request().Context()
	result, err := t.prService.AddTeam(ctx, domainTeam)
	if err != nil {
		t.logger.Errorf("failed to add team: %v", err)
		return HandleError(c, err)
	}

	return c.JSON(http.StatusCreated, AddTeamResponse{
		Team: TeamToResponse(result),
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
func (t *PRHandler) GetTeam(c echo.Context) error {
	teamName := c.QueryParam("team_name")
	if teamName == "" {
		t.logger.Errorf("failed to get team: team_name is required")
		return Error(c, http.StatusBadRequest, "INVALID_INPUT", "team_name is required")
	}

	ctx := c.Request().Context()
	team, err := t.prService.GetTeam(ctx, teamName)
	if err != nil {
		t.logger.Errorf("failed to get team: %v", err)
		return HandleError(c, err)
	}

	return c.JSON(http.StatusOK, TeamToResponse(team))
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
func (h *PRHandler) SetIsActive(c echo.Context) error {
	var req SetIsActiveRequest
	if err := h.bindRequest(c, &req); err != nil {
		return err
	}

	ctx := c.Request().Context()
	user, err := h.prService.SetIsActive(ctx, req.UserID, req.IsActive)
	if err != nil {
		h.logger.Errorf("failed to set user active status: %v", err)
		return HandleError(c, err)
	}

	return c.JSON(http.StatusOK, SetIsActiveResponse{
		User: UserToResponse(user),
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
func (h *PRHandler) GetReview(c echo.Context) error {
	userID := c.QueryParam("user_id")
	if userID == "" {
		h.logger.Errorf("failed to get review: user_id is required")
		return Error(c, http.StatusBadRequest, "INVALID_INPUT", "user_id is required")
	}

	ctx := c.Request().Context()
	pullRequests, err := h.prService.GetReview(ctx, userID)
	if err != nil {
		h.logger.Errorf("failed to get review: %v", err)
		return HandleError(c, err)
	}

	prShorts := PullRequestsToShort(pullRequests)

	return c.JSON(http.StatusOK, GetReviewResponse{
		UserID:       userID,
		PullRequests: prShorts,
	})
}
