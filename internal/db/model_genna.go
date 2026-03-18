// nolint
//
//lint:file-ignore U1000 ignore unused code, it's generated
package db

import (
	"time"
)

type ColumnsPrSystemPrReviewer struct {
	ID, PrID, ReviewerID, AssignedAt string
	Pr, Reviewer                     string
}

type ColumnsPrSystemPullRequest struct {
	ID, PullRequestID, PullRequestName, AuthorID, StatusID, MergedAt, CreatedAt string
	Author, Status                                                              string
}

type ColumnsPrSystemStatus struct {
	ID, Name, CreatedAt string
}

type ColumnsPrSystemTeam struct {
	ID, Name, CreatedAt string
}

type ColumnsPrSystemUser struct {
	ID, UserID, Username, IsActive, TeamID, CreatedAt string
	Team                                              string
}

type ColumnsSt struct {
	PrSystemPrReviewer  ColumnsPrSystemPrReviewer
	PrSystemPullRequest ColumnsPrSystemPullRequest
	PrSystemStatus      ColumnsPrSystemStatus
	PrSystemTeam        ColumnsPrSystemTeam
	PrSystemUser        ColumnsPrSystemUser
}

var Columns = ColumnsSt{
	PrSystemPrReviewer: ColumnsPrSystemPrReviewer{
		ID:         "id",
		PrID:       "pr_id",
		ReviewerID: "reviewer_id",
		AssignedAt: "assigned_at",

		Pr:       "Pr",
		Reviewer: "Reviewer",
	},
	PrSystemPullRequest: ColumnsPrSystemPullRequest{
		ID:              "id",
		PullRequestID:   "pull_request_id",
		PullRequestName: "pull_request_name",
		AuthorID:        "author_id",
		StatusID:        "status_id",
		MergedAt:        "merged_at",
		CreatedAt:       "created_at",

		Author: "Author",
		Status: "Status",
	},
	PrSystemStatus: ColumnsPrSystemStatus{
		ID:        "id",
		Name:      "name",
		CreatedAt: "created_at",
	},
	PrSystemTeam: ColumnsPrSystemTeam{
		ID:        "id",
		Name:      "name",
		CreatedAt: "created_at",
	},
	PrSystemUser: ColumnsPrSystemUser{
		ID:        "id",
		UserID:    "user_id",
		Username:  "username",
		IsActive:  "is_active",
		TeamID:    "team_id",
		CreatedAt: "created_at",

		Team: "Team",
	},
}

type TablePrSystemPrReviewer struct {
	Name, Alias string
}

type TablePrSystemPullRequest struct {
	Name, Alias string
}

type TablePrSystemStatus struct {
	Name, Alias string
}

type TablePrSystemTeam struct {
	Name, Alias string
}

type TablePrSystemUser struct {
	Name, Alias string
}

type TablesSt struct {
	PrSystemPrReviewer  TablePrSystemPrReviewer
	PrSystemPullRequest TablePrSystemPullRequest
	PrSystemStatus      TablePrSystemStatus
	PrSystemTeam        TablePrSystemTeam
	PrSystemUser        TablePrSystemUser
}

var Tables = TablesSt{
	PrSystemPrReviewer: TablePrSystemPrReviewer{
		Name:  "pr_system.pr_reviewers",
		Alias: "t",
	},
	PrSystemPullRequest: TablePrSystemPullRequest{
		Name:  "pr_system.pull_requests",
		Alias: "t",
	},
	PrSystemStatus: TablePrSystemStatus{
		Name:  "pr_system.statuses",
		Alias: "t",
	},
	PrSystemTeam: TablePrSystemTeam{
		Name:  "pr_system.teams",
		Alias: "t",
	},
	PrSystemUser: TablePrSystemUser{
		Name:  "pr_system.users",
		Alias: "t",
	},
}

type PrSystemPrReviewer struct {
	tableName struct{} `pg:"pr_system.pr_reviewers,alias:t,discard_unknown_columns"`

	ID         int64      `pg:"id,pk"`
	PrID       int64      `pg:"pr_id,use_zero"`
	ReviewerID int64      `pg:"reviewer_id,use_zero"`
	AssignedAt *time.Time `pg:"assigned_at"`

	Pr       *PrSystemPullRequest `pg:"fk:pr_id,rel:has-one"`
	Reviewer *PrSystemUser        `pg:"fk:reviewer_id,rel:has-one"`
}

type PrSystemPullRequest struct {
	tableName struct{} `pg:"pr_system.pull_requests,alias:t,discard_unknown_columns"`

	ID              int64      `pg:"id,pk"`
	PullRequestID   string     `pg:"pull_request_id,use_zero"`
	PullRequestName string     `pg:"pull_request_name,use_zero"`
	AuthorID        int64      `pg:"author_id,use_zero"`
	StatusID        int        `pg:"status_id,use_zero"`
	MergedAt        *time.Time `pg:"merged_at"`
	CreatedAt       *time.Time `pg:"created_at"`

	Author *PrSystemUser   `pg:"fk:author_id,rel:has-one"`
	Status *PrSystemStatus `pg:"fk:status_id,rel:has-one"`
}

type PrSystemStatus struct {
	tableName struct{} `pg:"pr_system.statuses,alias:t,discard_unknown_columns"`

	ID        int        `pg:"id,pk"`
	Name      string     `pg:"name,use_zero"`
	CreatedAt *time.Time `pg:"created_at"`
}

type PrSystemTeam struct {
	tableName struct{} `pg:"pr_system.teams,alias:t,discard_unknown_columns"`

	ID        int64      `pg:"id,pk"`
	Name      string     `pg:"name,use_zero"`
	CreatedAt *time.Time `pg:"created_at"`
}

type PrSystemUser struct {
	tableName struct{} `pg:"pr_system.users,alias:t,discard_unknown_columns"`

	ID        int64      `pg:"id,pk"`
	UserID    string     `pg:"user_id,use_zero"`
	Username  string     `pg:"username,use_zero"`
	IsActive  bool       `pg:"is_active,use_zero"`
	TeamID    *int64     `pg:"team_id"`
	CreatedAt *time.Time `pg:"created_at"`

	Team *PrSystemTeam `pg:"fk:team_id,rel:has-one"`
}
