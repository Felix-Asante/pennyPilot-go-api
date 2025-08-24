package handlers

import (
	"context"
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
			// auth
			r.Post("/auth/forgot-password", h.forgotPassword)
			r.Post("/auth/reset-password", h.resetPassword)
			r.Get("/auth/me", h.getMe)

			// income
			r.Post("/income", h.createIncome)
			r.Put("/income/{id}", h.updateIncome)

			// users
			r.Get("/user/income", h.getUserIncome)
			r.Get("/user/accounts", h.getAccounts)
			r.Get("/user/envelopes", h.getUserOwnedEnvelopes)

			// accounts
			r.Post("/account", h.createAccount)
			r.Get("/account/{id}", h.getAccount)
			r.Put("/account/{id}", h.updateAccount)
			r.Delete("/account/{id}", h.deleteAccount)
			r.Get("/account/{id}/envelopes", h.getEnvelopes)

			// envelopes
			r.Post("/envelope", h.createEnvelope)
			r.Get("/envelope/{id}", h.getEnvelope)
			r.Put("/envelope/{id}", h.updateEnvelope)
			r.Delete("/envelope/{id}", h.deleteEnvelope)
		})

		// public routes
		r.Route("/auth", func(r chi.Router) {
			r.Post("/register", h.createUser)
			r.Post("/login", h.login)
		})
	})
}

func getUserIdFromContext(ctx context.Context) (string, error) {
	_, claims, err := jwtauth.FromContext(ctx)
	if err != nil {
		return "", err
	}
	return claims["user_id"].(string), nil
}
