package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"

	repository "gapi/internal/db"
)

type Server struct {
	queries *repository.Queries
	router  *mux.Router
}

func (s *Server) setupRoutes() {
	// Define your routes here
}

func NewServer(db *sql.DB, queries *repository.Queries) *Server {
	server := &Server{
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

	queries := repository.New(db)
	server := NewServer(db, queries)

	http.ListenAndServe(":8080", server.router)

	//ctx := context.Background()

	// user, err := queries.CreateUser(ctx, server.CreateUserParams{
	// 	JiraID:   "abc-123",
	// 	Username: "klym",
	// })
	// if err != nil {
	// 	log.Fatal(err)
	// }

	//log.Printf("new user: %+v\n", user)
}
