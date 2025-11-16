package domain

type ReviewerStats struct {
	UserID         string
	Username       string
	AssignedCount  int
	CompletedCount int
	ActiveCount    int
}
