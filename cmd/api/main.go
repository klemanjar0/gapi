package main

import (
	//"context"
	"database/sql"
	"gapi/internal/handler"
	"gapi/internal/service"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"log"
	"net/http"

	repository "gapi/internal/db"
	"gapi/internal/utility"
)

type Server struct {
	db      *sql.DB
	queries *repository.Queries
	router  *mux.Router
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
}

func NewServer(db *sql.DB, queries *repository.Queries) *Server {
	server := &Server{
		db:      db,
		queries: queries,
		router:  mux.NewRouter(),
	}

	server.setupRoutes()

	return server
}

const connString = "postgres://postgres:postgres@localhost:5432/gapi?sslmode=disable"

func main() {
	db, err := sql.Open("postgres", connString)

	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	bootstrap(db)

	queries := repository.New(db)
	server := NewServer(db, queries)

	utility.LogInfo("Server starting on port 8080")
	if err := http.ListenAndServe(":8080", server.router); err != nil {
		utility.LogError("Failed to start server: " + err.Error())
		log.Fatal(err)
	}
}
