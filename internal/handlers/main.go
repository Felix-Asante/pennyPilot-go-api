package handlers

import (
	"log/slog"

	"github.com/Felix-Asante/pennyPilot-go-api/internal/models"
	"github.com/Felix-Asante/pennyPilot-go-api/internal/notifications"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
	"gorm.io/gorm"
)

type Handler struct {
	DB            *gorm.DB
	Logger        *slog.Logger
	Router        *chi.Mux
	Models        *models.Models
	JWTAuth       *jwtauth.JWTAuth
	Notifications *notifications.NotificationService
}

func NewHandler(config *Handler) *Handler {
	return &Handler{
		DB:            config.DB,
		Logger:        config.Logger,
		Router:        config.Router,
		Models:        config.Models,
		JWTAuth:       config.JWTAuth,
		Notifications: config.Notifications,
	}
}

func (h *Handler) CreateRoutes() {

	h.Router.Route("/api/v1", func(r chi.Router) {

		// protected routes
		r.Group(func(r chi.Router) {
			r.Use(jwtauth.Verifier(h.JWTAuth))
			r.Post("/auth/forgot-password", h.forgotPassword)
			r.Post("/auth/reset-password", h.resetPassword)
			r.Get("/auth/me", h.getMe)
		})

		// public routes
		r.Route("/auth", func(r chi.Router) {
			r.Post("/register", h.createUser)
			r.Post("/login", h.login)
		})
	})
}
