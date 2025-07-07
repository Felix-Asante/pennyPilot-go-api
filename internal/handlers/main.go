package handlers

import (
	"log/slog"

	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

type Handler struct {
	DB     *gorm.DB
	Logger *slog.Logger
	Router *chi.Mux
}

func NewHandler(config *Handler) *Handler {
	return &Handler{
		DB:     config.DB,
		Logger: config.Logger,
		Router: config.Router,
	}
}

func (h *Handler) CreateRoutes() {
	h.Router.Route("/api/v1", func(r chi.Router) {

		// user routes
		r.Route("/users", func(r chi.Router) {
			r.Post("/", h.createUser)
			r.Get("/me", h.getMe)
		})
	})
}
