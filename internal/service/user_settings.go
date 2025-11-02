package service

import (
	"context"
	"database/sql"
	repository "gapi/internal/db"
)

type UserSettingsService struct {
	db      *sql.DB
	queries *repository.Queries
}

func NewUserSettingsService(
	db *sql.DB,
	queries *repository.Queries,
) *UserSettingsService {
	return &UserSettingsService{
		db:      db,
		queries: queries,
	}
}

func (s *UserSettingsService) GetSettingsByUserJiraID(
	ctx context.Context,
	jiraId string,
) (*repository.GetUserWithSettingsRow, error) {
	settings, err := s.queries.GetUserWithSettings(ctx, jiraId)

	if err != nil {
		return nil, err
	}

	return &settings, nil
}

func (s *UserSettingsService) UpdateUserSettings(
	ctx context.Context,
	jiraId string,
	params repository.UpsertSettingsParams,
) error {
	user, err := s.queries.GetUserByJiraID(ctx, jiraId)

	if err != nil {
		return err
	}

	params.UserID = user.ID

	err = s.queries.UpsertSettings(ctx, params)

	return err
}
