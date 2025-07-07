package main

import (
	"log/slog"
	"os"

	"github.com/Felix-Asante/pennyPilot-go-api/cmd/api"
	"github.com/Felix-Asante/pennyPilot-go-api/pkg/db"
	"github.com/Felix-Asante/pennyPilot-go-api/pkg/env"
	"github.com/go-chi/chi/v5"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))
	slog.SetDefault(logger)

	dbConfig := &db.DbConfig{
		DbUser:     env.GetEnv("DB_USER"),
		DbHost:     env.GetEnv("DB_HOST"),
		DbPassword: env.GetEnv("DB_PWD"),
		DbName:     env.GetEnv("DB_NAME"),
		DbPort:     env.GetEnv("DB_PORT"),
	}

	db, err := db.NewPgDB(*dbConfig).Init()
	if err != nil {
		panic(err)
	}

	apiConfig := &api.Server{
		Router: chi.NewRouter(),
		DB:     db,
		Logger: logger,
		Port:   env.GetEnv("PORT"),
	}

	server := api.Init(apiConfig)
	server.Run()
}
