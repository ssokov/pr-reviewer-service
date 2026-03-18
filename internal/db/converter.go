package db

func StatusIDToName(statusID int) PRStatus {
	switch statusID {
	case 1:
		return PRStatusOpen
	case 2:
		return PRStatusMerged
	default:
		return PRStatusOpen
	}
}

func toUser(in *PrSystemUser) User {
	out := User{
		ID:       in.ID,
		UserID:   in.UserID,
		Username: in.Username,
		IsActive: in.IsActive,
	}
	if in.TeamID != nil {
		out.TeamID = *in.TeamID
	}
	if in.Team != nil {
		out.TeamName = in.Team.Name
	}
	if in.CreatedAt != nil {
		out.CreatedAt = *in.CreatedAt
	}
	return out
}

func toUsers(in []PrSystemUser) []User {
	out := make([]User, 0, len(in))
	for i := range in {
		out = append(out, toUser(&in[i]))
	}
	return out
}

func toTeam(in *PrSystemTeam, members []PrSystemUser) Team {
	out := Team{
		ID:       in.ID,
		TeamName: in.Name,
		Members:  toUsers(members),
	}
	if in.CreatedAt != nil {
		out.CreatedAt = *in.CreatedAt
	}
	return out
}

func toPullRequest(in *PrSystemPullRequest, authorUserID string, status PRStatus, reviewers []string) PullRequest {
	out := PullRequest{
		ID:                in.ID,
		PullRequestID:     in.PullRequestID,
		PullRequestName:   in.PullRequestName,
		AuthorID:          in.AuthorID,
		AuthorUserID:      authorUserID,
		Status:            status,
		AssignedReviewers: append([]string(nil), reviewers...),
		MergedAt:          in.MergedAt,
	}
	if in.CreatedAt != nil {
		out.CreatedAt = *in.CreatedAt
	}
	return out
}
