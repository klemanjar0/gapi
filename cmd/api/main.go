package main

import (
	"database/sql"
	repository "gapi/internal/db"
	"gapi/internal/utility"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
}

func main() {
	dbConnString := os.Getenv(utility.DBConnEnv)
	if dbConnString == "" {
		log.Fatal("Database connection string env variable is missing. Check .env file.")
	}

	db, err := sql.Open("postgres", dbConnString)

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
