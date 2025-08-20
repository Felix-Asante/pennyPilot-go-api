package api

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/Felix-Asante/pennyPilot-go-api/internal/handlers"
	"github.com/Felix-Asante/pennyPilot-go-api/internal/models"
	"github.com/Felix-Asante/pennyPilot-go-api/pkg/env"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/httprate"
	"github.com/go-chi/jwtauth/v5"
	"gorm.io/gorm"
)

type Server struct {
	DB      *gorm.DB
	Logger  *slog.Logger
	Router  *chi.Mux
	Port    string
	JWTAuth *jwtauth.JWTAuth
}

func Init(apiConfig *Server) *Server {
	return &Server{
		DB:      apiConfig.DB,
		Logger:  apiConfig.Logger,
		Router:  apiConfig.Router,
		Port:    apiConfig.Port,
		JWTAuth: jwtauth.New("HS256", []byte(env.GetEnv("JWT_SECRET")), nil),
	}
}

func (s *Server) Run() {
	runMigrations(s.DB)
	setUpMiddleware(s.Router)

	handler := handlers.NewHandler(&handlers.Handler{
		DB:      s.DB,
		Logger:  s.Logger,
		Router:  s.Router,
		Models:  models.NewModels(s.DB),
		JWTAuth: s.JWTAuth,
	})

	handler.CreateRoutes()

	s.Logger.Info("Server started on port " + s.Port)
	http.ListenAndServe(":"+s.Port, s.Router)
}

func runMigrations(db *gorm.DB) {
	models := []interface{}{
		&models.User{},
	}

	if err := db.AutoMigrate(models...); err != nil {
		panic(err)
	}
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
