package main

import (
	"database/sql"
	repository "gapi/internal/db"
	"gapi/internal/handler"
	"gapi/internal/service"
	"gapi/internal/utility"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"

	"github.com/gorilla/mux"
)

type Server struct {
	db      *sql.DB
	queries *repository.Queries
	router  *mux.Router
	proxy   *httputil.ReverseProxy
}

func NewServer(db *sql.DB, queries *repository.Queries) *Server {
	server := &Server{
		db:      db,
		queries: queries,
		router:  mux.NewRouter(),
	}

	server.setupProxy()
	server.setupRoutes()

	return server
}

func (s *Server) setupProxy() {
	jiraHost := os.Getenv(utility.JiraHostEnv)
	if jiraHost == "" {
		log.Println("Jira Host env variable is missing, proxy disabled. Check .env file.")
		return
	}

	targetURL, err := url.Parse(jiraHost)
	if err != nil {
		log.Printf("Failed to parse Jira URL: %v", err)
		return
	}

	s.proxy = httputil.NewSingleHostReverseProxy(targetURL)

	originalDirector := s.proxy.Director
	s.proxy.Director = func(req *http.Request) {
		originalDirector(req)
		req.Host = targetURL.Host
		req.URL.Path = strings.TrimPrefix(req.URL.Path, "/api")
		log.Printf("Proxying request to: %s%s", req.URL.Host, req.URL.Path)
	}
}

func (s *Server) handleProxy(w http.ResponseWriter, r *http.Request) {
	if s.proxy == nil {
		http.Error(w, "Proxy is not configured", http.StatusServiceUnavailable)
		return
	}

	// CORS headers
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	s.proxy.ServeHTTP(w, r)
}

func (s *Server) setupRoutes() {
	userService := service.NewUserService(s.db, s.queries)
	userHandler := handler.NewUserHandler(userService)

	userSettingsService := service.NewUserSettingsService(s.db, s.queries)
	userSettingsHandler := handler.NewUserSettingsHandler(userSettingsService)

	s.router.HandleFunc("/users/{jira_id}", userHandler.GetUserByJiraId).Methods("GET")
	s.router.HandleFunc("/users", userHandler.CreateUser).Methods("POST")
	s.router.HandleFunc("/users/{jira_id}/check", userHandler.RefreshLogin).Methods("PATCH")

	s.router.HandleFunc("/users/{jira_id}/settings", userSettingsHandler.GetUserSettings).Methods("GET")
	s.router.HandleFunc("/users/{jira_id}/settings", userSettingsHandler.UpdateUserSettings).Methods("PUT")

	// Setup proxy for /api/ routes, important to be after other route definitions
	s.router.PathPrefix("/api/").HandlerFunc(s.handleProxy)
}
