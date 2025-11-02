package utility

import (
	repository "gapi/internal/db"
)

type UserSettingsResponse struct {
	JiraID             string
	ProjectID          string
	IssueQuery         string
	ContentTemplate    string
	TicketItemTemplate string
	MailRecipient      string
	MailSubject        string
	MailAuthor         string
}

func MapUserSettingsToResponse(settings *repository.GetUserWithSettingsRow) *UserSettingsResponse {
	s := &UserSettingsResponse{
		JiraID:             settings.JiraID,
		ProjectID:          settings.ProjectID.String,
		IssueQuery:         settings.IssueQuery.String,
		ContentTemplate:    settings.ContentTemplate.String,
		TicketItemTemplate: settings.TicketItemTemplate.String,
		MailRecipient:      settings.MailRecipient.String,
		MailSubject:        settings.MailSubject.String,
		MailAuthor:         settings.MailAuthor.String,
	}

	return s
}
