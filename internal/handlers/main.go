package handlers

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/Felix-Asante/pennyPilot-go-api/internal/models"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/httprate"
	"github.com/go-chi/jwtauth/v5"
	"gorm.io/gorm"
)

type Handler struct {
	DB      *gorm.DB
	Logger  *slog.Logger
	Router  *chi.Mux
	Models  *models.Models
	JWTAuth *jwtauth.JWTAuth
}

func NewHandler(config *Handler) *Handler {
	return &Handler{
		DB:      config.DB,
		Logger:  config.Logger,
		Router:  config.Router,
		Models:  config.Models,
		JWTAuth: config.JWTAuth,
	}
}

func (h *Handler) CreateRoutes() {
	initValidator()

	h.Router.Route("/api/v1", func(r chi.Router) {

		// protected routes
		r.Group(func(r chi.Router) {
			r.Use(jwtauth.Verifier(h.JWTAuth))
			r.Use(httprate.Limit(
				5,
				time.Minute,
				httprate.WithLimitHandler(func(w http.ResponseWriter, r *http.Request) {
					h.rateLimitExceededResponse(w, r, "1 minute")
				}),
			))
			r.Get("/auth/me", h.getMe)
		})

		// public routes
		r.Route("/auth", func(r chi.Router) {
			r.Post("/register", h.createUser)
			r.Post("/login", h.login)
		})
	})
}
