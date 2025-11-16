package dto

type UserStatsItem struct {
	UserID         string `json:"user_id"`
	Username       string `json:"username"`
	AssignedCount  int    `json:"assigned_count"`
	CompletedCount int    `json:"completed_count"`
	ActiveCount    int    `json:"active_count"`
}

type PRStatsItem struct {
	Status string `json:"status"`
	Count  int    `json:"count"`
}

type StatsResponse struct {
	TotalPRs     int             `json:"total_prs"`
	TotalUsers   int             `json:"total_users"`
	ActiveUsers  int             `json:"active_users"`
	PRsByStatus  []PRStatsItem   `json:"prs_by_status"`
	TopReviewers []UserStatsItem `json:"top_reviewers"`
}
