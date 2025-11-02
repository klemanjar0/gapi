package handler

import (
	"encoding/json"
	"net/http"

	builder "gapi/internal/response_builder"
	"gapi/internal/service"
	"github.com/gorilla/mux"
)

type UserHandler struct {
	userService *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

func (h *UserHandler) GetUserByJiraId(w http.ResponseWriter, r *http.Request) {
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

	user, err := h.userService.GetUserByJiraID(r.Context(), jiraId)

	if err != nil {
		rb.Error(
			err.Error(),
			http.StatusInternalServerError,
		)
		return
	}

	if user == nil {
		rb.Error(
			"User not found",
			http.StatusNotFound,
		)
		return
	}

	rb.Success(user)
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	rb := builder.NewResponseBuilder(w)
	var req struct {
		Username string `json:"username"`
		JiraID   string `json:"jira_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.JiraID == "" || req.Username == "" {
		rb.Error("Invalid request body", http.StatusBadRequest)
		return
	}

	existingUser, _ := h.userService.GetUserByJiraID(r.Context(), req.JiraID)

	if existingUser != nil {
		rb.Error("User with this Jira ID already exists", http.StatusConflict)
		return
	}

	user, err := h.userService.CreateUser(r.Context(), req.JiraID, req.Username)

	if err != nil {
		rb.Error(err.Error(), http.StatusInternalServerError)
		return
	}

	rb.Created(user)
}

func (h *UserHandler) RefreshLogin(w http.ResponseWriter, r *http.Request) {
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

	existingUser, err := h.userService.GetUserByJiraID(r.Context(), jiraId)

	if existingUser == nil {
		rb.Error("User not found", http.StatusNotFound)
		return
	}

	if err != nil {
		rb.Error(err.Error(), http.StatusInternalServerError)
		return
	}

	if err := h.userService.RefreshLogin(r.Context(), jiraId); err != nil {
		rb.Error(err.Error(), http.StatusInternalServerError)
		return
	}

	rb.Success(map[string]string{"message": "Connection successful"})
}
