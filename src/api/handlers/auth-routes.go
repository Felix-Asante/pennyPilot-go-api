package handlers

import (
	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

type AuthHandler struct {
	handlers *Handlers
	db       *gorm.DB
}

func NewAuthHandler(h *Handlers, db *gorm.DB) *AuthHandler {
	return &AuthHandler{handlers: h, db: db}
}

func (h *AuthHandler) SetupRoutes() {
	router := *h.handlers.router

	authRoutesHandler := authRoutesHandler{db: h.db}

	router.Route("/auth", func(route chi.Router) {

		route.Post("/login", authRoutesHandler.loginHandler)
		route.Post("/signup", authRoutesHandler.signupHandler)
	})

}
