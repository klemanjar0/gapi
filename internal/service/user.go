package service

import (
	"context"
	"database/sql"
	repository "gapi/internal/db"
)

type UserService struct {
	db      *sql.DB
	queries *repository.Queries
}

func NewUserService(db *sql.DB, queries *repository.Queries) *UserService {
	return &UserService{
		db:      db,
		queries: queries,
	}
}

func (s *UserService) CreateUser(ctx context.Context, jiraID string, username string) (*repository.User, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	qtx := s.queries.WithTx(tx)

	user, err := qtx.CreateUser(ctx, repository.CreateUserParams{
		JiraID:   jiraID,
		Username: username,
	})

	if err != nil {
		return nil, err
	}

	settingsError := qtx.UpsertSettings(ctx, repository.UpsertSettingsParams{
		UserID:             user.ID,
		ProjectID:          "",
		IssueQuery:         "",
		ContentTemplate:    "",
		TicketItemTemplate: "",
		MailRecipient:      "",
		MailSubject:        "",
		MailAuthor:         "",
	})

	if settingsError != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
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

func (s *UserService) RefreshLogin(ctx context.Context, id string) error {
	return s.queries.UpdateLastLogin(ctx, id)
}
