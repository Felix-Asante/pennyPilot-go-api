package handlers

import (
	"github.com/felix-Asante/pennyPilot-go-api/src/pkgs/jwt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
	"gorm.io/gorm"
)

type UsersHandler struct {
	handlers *Handlers
	db       *gorm.DB
}

func NewUsersHandler(h *Handlers, db *gorm.DB) *UsersHandler {
	return &UsersHandler{handlers: h, db: db}
}

func (h *UsersHandler) SetupRoutes() {
	router := *h.handlers.router

	usersRoutesHandler := usersRoutesHandler{db: h.db}

	router.Route("/users", func(route chi.Router) {

		route.Use(jwtauth.Verifier(jwt.InitAuthToken()))
		route.Use(jwtauth.Authenticator(jwt.InitAuthToken()))
		route.Get("/{userId}/accounts", usersRoutesHandler.getAccounts)
	})
}
