package rest

import (
	"testing"
	"time"

	"github.com/ssokov/pr-reviewer-service/internal/pr"
)

func TestCreatePRRequestToDomain(t *testing.T) {
	req := CreatePRRequest{
		PullRequestID:   "pr-1",
		PullRequestName: "Test PR",
		AuthorID:        "u1",
	}

	got := CreatePRRequestToDomain(req)
	if got.PullRequestID != "pr-1" || got.PullRequestName != "Test PR" || got.AuthorID != "u1" || got.Status != pr.PRStatusOpen {
		t.Fatalf("unexpected result: %#v", got)
	}
}

func TestTeamMappings(t *testing.T) {
	req := AddTeamRequest{
		TeamName: "backend",
		Members: []TeamMember{
			{UserID: "u1", Username: "Alice", IsActive: true},
			{UserID: "u2", Username: "Bob", IsActive: false},
		},
	}

	domain := AddTeamRequestToDomain(req)
	if len(domain.Members) != 2 || domain.Members[0].UserID != "u1" {
		t.Fatalf("unexpected domain mapping: %#v", domain)
	}

	resp := TeamToResponse(domain)
	if resp.TeamName != "backend" || len(resp.Members) != 2 || resp.Members[1].UserID != "u2" {
		t.Fatalf("unexpected response mapping: %#v", resp)
	}
}

func TestPullRequestsToShort(t *testing.T) {
	now := time.Now()
	in := []pr.PullRequest{
		{PullRequestID: "pr-1", PullRequestName: "one", AuthorID: "u1", Status: pr.PRStatusOpen, CreatedAt: now},
		{PullRequestID: "pr-2", PullRequestName: "two", AuthorID: "u2", Status: pr.PRStatusMerged, CreatedAt: now},
	}

	got := PullRequestsToShort(in)
	if len(got) != 2 {
		t.Fatalf("unexpected len: %d", len(got))
	}
	if got[0].Status != "OPEN" || got[1].Status != "MERGED" {
		t.Fatalf("unexpected status mapping: %#v", got)
	}
}
