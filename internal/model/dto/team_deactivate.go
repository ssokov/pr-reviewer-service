package dto

type DeactivateTeamRequest struct {
	TeamName string `json:"team_name" validate:"required"`
}

type DeactivateTeamResponse struct {
	DeactivatedUsers int                   `json:"deactivated_users"`
	ReassignedPRs    int                   `json:"reassigned_prs"`
	Users            []DeactivatedUserInfo `json:"users"`
}

type DeactivatedUserInfo struct {
	UserID       string `json:"user_id"`
	Username     string `json:"username"`
	OpenPRsCount int    `json:"open_prs_count"`
}
