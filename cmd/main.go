package main

import (
	"Anastasia/songs/internal/api"
	"Anastasia/songs/internal/repository"
	"Anastasia/songs/internal/services"
	"log"
	"net/http"
	"os"

	_ "Anastasia/songs/docs"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

//	@title			Swagger Music API
//	@version		1.0
//	@description	API for online songs library.

//	@host	localhost:8080

func main() {

	db, err := repository.NewStorage(repository.Config{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   os.Getenv("DB_NAME"),
		SSLMode:  os.Getenv("DB_SSLMODE"),
	})
	if err != nil {
		logrus.WithError(err).Fatalf("Failed to connect to DB")
	}

	defer db.Close()

	repo := repository.NewRepo(db)
	srv := services.NewService(repo)

	api := api.New(srv)

	logrus.Info("Service is running...")
	err = http.ListenAndServe(os.Getenv("PORT"), api.Router())
	if err != nil {
		logrus.WithError(err).Fatalf("Failed starting the service")
	}
}

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	logLevel := os.Getenv("LOG_LEVEL")
	level, err := logrus.ParseLevel(logLevel)
	if err != nil {
		log.Fatalf("Invalid log level: %v", err)
	}
	logrus.SetLevel(level)
	logrus.SetFormatter(&logrus.JSONFormatter{})
}
