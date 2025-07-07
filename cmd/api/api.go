package api

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/Felix-Asante/pennyPilot-go-api/internal/models"
	"github.com/Felix-Asante/pennyPilot-go-api/pkg/db"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/httprate"
)

type Server struct {
	DbConfig *db.DbConfig
	Logger   *slog.Logger
	Router   *chi.Mux
	Port     string
}

func Init(apiConfig *Server) *Server {
	return &Server{
		DbConfig: apiConfig.DbConfig,
		Logger:   apiConfig.Logger,
		Router:   apiConfig.Router,
		Port:     apiConfig.Port,
	}
}

func (s *Server) Run() {
	setUpDbWithAutoMigrate(s.DbConfig)
	setUpMiddleware(s.Router)

	s.Router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World"))
	})

	s.Logger.Info("Server started on port " + s.Port)
	http.ListenAndServe(":"+s.Port, s.Router)
}

func setUpDbWithAutoMigrate(dbConfig *db.DbConfig) {
	db, err := db.NewPgDB(*dbConfig).Init()
	if err != nil {
		panic(err)
	}

	models := []interface{}{
		&models.User{},
	}

	db.AutoMigrate(models...)
}

func setUpMiddleware(r *chi.Mux) {
	r.Use(r.Middlewares()...)
	r.Use(httprate.LimitByIP(10, 1*time.Minute))
	r.Use(middleware.AllowContentType("application/json", "multipart/form-data", "text/xml"))
	r.Use(middleware.CleanPath)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"http://localhost:3000"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "PATCH"},
	}))
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.StripSlashes)
}
