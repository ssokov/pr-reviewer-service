package rest

import (
	"github.com/ssokov/pr-reviewer-service/internal/pr"
)

func CreatePRRequestToDomain(req CreatePRRequest) *pr.PullRequest {
	return &pr.PullRequest{
		PullRequestID:   req.PullRequestID,
		PullRequestName: req.PullRequestName,
		AuthorID:        req.AuthorID,
		Status:          pr.PRStatusOpen,
	}
}

func PullRequestToResponse(pr *pr.PullRequest) PullRequestResponse {
	return PullRequestResponse{
		PullRequestID:     pr.PullRequestID,
		PullRequestName:   pr.PullRequestName,
		AuthorID:          pr.AuthorID,
		Status:            string(pr.Status),
		AssignedReviewers: pr.AssignedReviewers,
		CreatedAt:         &pr.CreatedAt,
		MergedAt:          pr.MergedAt,
	}
}

func AddTeamRequestToDomain(req AddTeamRequest) *pr.Team {
	return &pr.Team{
		TeamName: req.TeamName,
		Members: Map(req.Members, func(member TeamMember) pr.User {
			return pr.User{
				UserID:   member.UserID,
				Username: member.Username,
				IsActive: member.IsActive,
			}
		}),
	}
}

func TeamToResponse(team *pr.Team) TeamResponse {
	return TeamResponse{
		TeamName: team.TeamName,
		Members:  NewTeamMembers(team.Members),
	}
}

func UserToResponse(user *pr.User) UserResponse {
	return UserResponse{
		UserID:   user.UserID,
		Username: user.Username,
		TeamName: user.TeamName,
		IsActive: user.IsActive,
	}
}

func PullRequestsToShort(prs []pr.PullRequest) []PullRequestShort {
	return NewPullRequestShorts(prs)
}

func NewTeamMember(in pr.User) TeamMember {
	return TeamMember{
		UserID:   in.UserID,
		Username: in.Username,
		IsActive: in.IsActive,
	}
}

func NewPullRequestShort(in pr.PullRequest) PullRequestShort {
	return PullRequestShort{
		PullRequestID:   in.PullRequestID,
		PullRequestName: in.PullRequestName,
		AuthorID:        in.AuthorID,
		Status:          string(in.Status),
	}
}

// Map converts slice of type T to slice of type M using value converter.
func Map[T, M any](a []T, f func(T) M) []M {
	n := make([]M, len(a))
	for i := range a {
		n[i] = f(a[i])
	}
	return n
}
