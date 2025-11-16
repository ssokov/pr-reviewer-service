package pr

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/ssokov/pr-reviewer-service/internal/http/mapper"
	"github.com/ssokov/pr-reviewer-service/internal/http/response"
	"github.com/ssokov/pr-reviewer-service/internal/model/dto"
	"github.com/ssokov/pr-reviewer-service/internal/service"
	"github.com/vmkteam/embedlog"
)

type PRHandler struct {
	prService service.PRService
	logger    embedlog.Logger
}

func NewHandler(service service.PRService, logger embedlog.Logger) *PRHandler {
	return &PRHandler{
		prService: service,
		logger:    logger,
	}
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
	var req dto.CreatePRRequest
	if err := c.Bind(&req); err != nil {
		p.logger.Errorf("failed to bind request: %v", err)
		return response.Error(c, http.StatusBadRequest, "INVALID_INPUT", "invalid request body")
	}

	domainPR := mapper.CreatePRRequestToDomain(req)

	ctx := c.Request().Context()
	createdPR, err := p.prService.CreatePR(ctx, req.AuthorID, domainPR)
	if err != nil {
		p.logger.Errorf("failed to create PR: %v", err)
		return response.HandleError(c, err)
	}

	return c.JSON(http.StatusCreated, dto.CreatePRResponse{
		PR: mapper.PullRequestToResponse(createdPR),
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
// @Failure 409 {object} dto.ErrorResponse "PR already merged"
// @Failure 500 {object} dto.ErrorResponse
// @Router /pullRequest/merge [post]
func (p *PRHandler) MergePR(c echo.Context) error {
	var req dto.MergePRRequest
	if err := c.Bind(&req); err != nil {
		p.logger.Errorf("failed to bind request: %v", err)
		return response.Error(c, http.StatusBadRequest, "INVALID_INPUT", "invalid request body")
	}

	ctx := c.Request().Context()
	mergedPR, err := p.prService.MergePR(ctx, req.PullRequestID)
	if err != nil {
		p.logger.Errorf("failed to merge PR: %v", err)
		return response.HandleError(c, err)
	}

	return c.JSON(http.StatusOK, dto.MergePRResponse{
		PR: mapper.PullRequestToResponse(mergedPR),
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
	var req dto.ReassignRequest
	if err := c.Bind(&req); err != nil {
		p.logger.Errorf("failed to bind request: %v", err)
		return response.Error(c, http.StatusBadRequest, "INVALID_INPUT", "invalid request body")
	}

	ctx := c.Request().Context()
	pr, newReviewerID, err := p.prService.ReassignReviewer(ctx, req.PullRequestID, req.OldUserID)
	if err != nil {
		p.logger.Errorf("failed to reassign reviewer: %v", err)
		return response.HandleError(c, err)
	}

	return c.JSON(http.StatusOK, dto.ReassignResponse{
		PR:         mapper.PullRequestToResponse(pr),
		ReplacedBy: newReviewerID,
	})
}
