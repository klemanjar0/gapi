package handler

import (
	"encoding/json"
	repository "gapi/internal/db"
	builder "gapi/internal/response_builder"
	"gapi/internal/service"
	utility "gapi/internal/utility"
	"github.com/gorilla/mux"
	"net/http"
)

type UserSettingsHandler struct {
	userSettingsService *service.UserSettingsService
}

func NewUserSettingsHandler(userSettingsService *service.UserSettingsService) *UserSettingsHandler {
	return &UserSettingsHandler{userSettingsService: userSettingsService}
}

func (h *UserSettingsHandler) GetUserSettings(w http.ResponseWriter, r *http.Request) {
	rb := builder.NewResponseBuilder(w)
	vars := mux.Vars(r)
	jiraId := vars["jira_id"]

	if jiraId == "" {
		rb.Error(
			"No Jira ID provided (jira_id var missing)",
			http.StatusBadRequest,
		)
		return
	}

	settings, err := h.userSettingsService.GetSettingsByUserJiraID(r.Context(), jiraId)

	if err != nil {
		rb.Error(
			err.Error(),
			http.StatusInternalServerError,
		)
		return
	}

	rb.Success(utility.MapUserSettingsToResponse(settings))
}

func (h *UserSettingsHandler) UpdateUserSettings(w http.ResponseWriter, r *http.Request) {
	rb := builder.NewResponseBuilder(w)
	vars := mux.Vars(r)
	jiraId := vars["jira_id"]
	if jiraId == "" {
		rb.Error(
			"No Jira ID provided (jira_id var missing)",
			http.StatusBadRequest,
		)
		return
	}
	var req struct {
		ProjectID          string `json:"project_id"`
		IssueQuery         string `json:"issue_query"`
		ContentTemplate    string `json:"content_template"`
		TicketItemTemplate string `json:"ticket_item_template"`
		MailRecipient      string `json:"mail_recipient"`
		MailSubject        string `json:"mail_subject"`
		MailAuthor         string `json:"mail_author"`
	}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		rb.Error(
			"Invalid request body",
			http.StatusBadRequest,
		)
		return
	}
	params := repository.UpsertSettingsParams{
		ProjectID:          req.ProjectID,
		IssueQuery:         req.IssueQuery,
		ContentTemplate:    req.ContentTemplate,
		TicketItemTemplate: req.TicketItemTemplate,
		MailRecipient:      req.MailRecipient,
		MailSubject:        req.MailSubject,
		MailAuthor:         req.MailAuthor,
	}
	err = h.userSettingsService.UpdateUserSettings(r.Context(), jiraId, params)
	if err != nil {
		rb.Error(
			err.Error(),
			http.StatusInternalServerError,
		)
		return
	}

	rb.Success(map[string]string{"status": "success"})
}
