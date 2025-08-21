package api

import (
	"log/slog"
	"net/http"

	"github.com/Felix-Asante/pennyPilot-go-api/internal/handlers"
	"github.com/Felix-Asante/pennyPilot-go-api/internal/models"
	"github.com/Felix-Asante/pennyPilot-go-api/internal/notifications"
	"github.com/Felix-Asante/pennyPilot-go-api/internal/utils"
	"github.com/Felix-Asante/pennyPilot-go-api/pkg/env"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/jwtauth/v5"
	"gorm.io/gorm"
)

type Server struct {
	DB            *gorm.DB
	Logger        *slog.Logger
	Router        *chi.Mux
	Port          string
	JWTAuth       *jwtauth.JWTAuth
	Notifications *notifications.NotificationService
}

func Init(apiConfig *Server) *Server {
	return &Server{
		DB:            apiConfig.DB,
		Logger:        apiConfig.Logger,
		Router:        apiConfig.Router,
		Port:          apiConfig.Port,
		JWTAuth:       jwtauth.New("HS256", []byte(env.GetEnv("JWT_SECRET")), nil),
		Notifications: apiConfig.Notifications,
	}
}

func (s *Server) Run() {
	runMigrations(s.DB)
	setUpMiddleware(s.Router)
	utils.InitializeValidator()

	handler := handlers.NewHandler(&handlers.Handler{
		DB:            s.DB,
		Logger:        s.Logger,
		Router:        s.Router,
		Models:        models.NewModels(s.DB),
		JWTAuth:       s.JWTAuth,
		Notifications: s.Notifications,
	})

	handler.CreateRoutes()

	s.Logger.Info("Server started on port " + s.Port)
	http.ListenAndServe(":"+s.Port, s.Router)
}

func runMigrations(db *gorm.DB) {
	models := []interface{}{
		&models.User{},
		&models.Code{},
		&models.Income{},
	}

	if err := db.AutoMigrate(models...); err != nil {
		panic(err)
	}
}

func setUpMiddleware(r *chi.Mux) {
	r.Use(r.Middlewares()...)
	r.Use(middleware.RealIP)
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
