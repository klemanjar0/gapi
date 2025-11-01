package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

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
	vars := mux.Vars(r)
	jiraId := vars["jira_id"]

	if jiraId == "" {
		http.Error(w, "invalid jiraId", http.StatusBadRequest)
		return
	}

	user, err := h.userService.GetUserByJiraID(r.Context(), jiraId)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Username string `json:"username"`
		JiraID   string `json:"jira_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	fmt.Println("Request")
	fmt.Println(req)

	existingUser, _ := h.userService.GetUserByJiraID(r.Context(), req.JiraID)
	fmt.Println("Existing user:", existingUser)
	if existingUser != nil {
		http.Error(w, "user with this Jira ID already exists", http.StatusConflict)
		return
	}

	user, err := h.userService.CreateUser(r.Context(), req.JiraID, req.Username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}
