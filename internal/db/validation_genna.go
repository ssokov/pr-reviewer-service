// nolint
//
//lint:file-ignore U1000 ignore unused code, it's generated
package db

import (
	"unicode/utf8"
)

const (
	ErrEmptyValue = "empty"
	ErrMaxLength  = "len"
	ErrWrongValue = "value"
)

func (m PrSystemPrReviewer) Validate() (errors map[string]string, valid bool) {
	errors = map[string]string{}

	if m.PrID == 0 {
		errors[Columns.PrSystemPrReviewer.PrID] = ErrEmptyValue
	}

	if m.ReviewerID == 0 {
		errors[Columns.PrSystemPrReviewer.ReviewerID] = ErrEmptyValue
	}

	return errors, len(errors) == 0
}

func (m PrSystemPullRequest) Validate() (errors map[string]string, valid bool) {
	errors = map[string]string{}

	if utf8.RuneCountInString(m.PullRequestID) > 255 {
		errors[Columns.PrSystemPullRequest.PullRequestID] = ErrMaxLength
	}

	if utf8.RuneCountInString(m.PullRequestName) > 255 {
		errors[Columns.PrSystemPullRequest.PullRequestName] = ErrMaxLength
	}

	if m.AuthorID == 0 {
		errors[Columns.PrSystemPullRequest.AuthorID] = ErrEmptyValue
	}

	if m.StatusID == 0 {
		errors[Columns.PrSystemPullRequest.StatusID] = ErrEmptyValue
	}

	return errors, len(errors) == 0
}

func (m PrSystemStatus) Validate() (errors map[string]string, valid bool) {
	errors = map[string]string{}

	if utf8.RuneCountInString(m.Name) > 50 {
		errors[Columns.PrSystemStatus.Name] = ErrMaxLength
	}

	return errors, len(errors) == 0
}

func (m PrSystemTeam) Validate() (errors map[string]string, valid bool) {
	errors = map[string]string{}

	if utf8.RuneCountInString(m.Name) > 255 {
		errors[Columns.PrSystemTeam.Name] = ErrMaxLength
	}

	return errors, len(errors) == 0
}

func (m PrSystemUser) Validate() (errors map[string]string, valid bool) {
	errors = map[string]string{}

	if utf8.RuneCountInString(m.UserID) > 255 {
		errors[Columns.PrSystemUser.UserID] = ErrMaxLength
	}

	if utf8.RuneCountInString(m.Username) > 255 {
		errors[Columns.PrSystemUser.Username] = ErrMaxLength
	}

	if m.TeamID != nil && *m.TeamID == 0 {
		errors[Columns.PrSystemUser.TeamID] = ErrEmptyValue
	}

	return errors, len(errors) == 0
}
