package main

import (
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"t-challenge/internal"
	"t-challenge/internal/durable"
	"t-challenge/internal/model"
)

func init() {
	// setup logger
	durable.SetupLogger()

	// load .env file
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("Error loading .env file")
	}

	if err := durable.ConnectDB(os.Getenv("DB_DSN")); err != nil {
		log.Fatal("Error connecting to database")
	}

	// migrate database
	if err := durable.Connection().AutoMigrate(
		&model.User{},
		&model.Account{},
		&model.Balance{},
		&model.Transaction{},
	); err != nil {
		log.Fatal(err)
	}
}

func main() {
	server := internal.Server{}

	mux := http.NewServeMux()
	server.SetupRoutes(mux)

	middlewareMux := server.SetupMiddleware(mux)
	server.StartServer(middlewareMux, os.Getenv("SERVER_PORT"))
}
