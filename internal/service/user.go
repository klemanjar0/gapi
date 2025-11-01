package service

import (
	"context"
	"fmt"
	repository "gapi/internal/db"
)

type UserService struct {
	queries *repository.Queries
}

func NewUserService(queries *repository.Queries) *UserService {
	return &UserService{
		queries: queries,
	}
}

func (s *UserService) CreateUser(ctx context.Context, jiraID string, username string) (*repository.User, error) {
	user, err := s.queries.CreateUser(ctx, repository.CreateUserParams{
		JiraID:   jiraID,
		Username: username,
	})

	if err != nil {
		return nil, err
	}

	settingsError := s.queries.UpsertSettings(ctx, repository.UpsertSettingsParams{
		UserID:             user.ID,
		ProjectID:          "",
		IssueQuery:         "",
		ContentTemplate:    "",
		TicketItemTemplate: "",
		MailRecipient:      "",
		MailSubject:        "",
		MailAuthor:         "",
	}) // possible error ignored, should be replaced with transaction later

	if settingsError != nil {
		fmt.Println("Failed to create default settings for user:", settingsError)
	}

	return &user, err
}

func (s *UserService) GetUserByJiraID(ctx context.Context, id string) (*repository.User, error) {
	user, err := s.queries.GetUserByJiraID(ctx, id)

	if err != nil {
		return nil, err
	}

	return &user, nil
}
