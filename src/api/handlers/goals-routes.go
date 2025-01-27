package handlers

import (
	"github.com/felix-Asante/pennyPilot-go-api/src/pkgs/jwt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
	"gorm.io/gorm"
)

type GoalsHandler struct {
	handlers *Handlers
	db       *gorm.DB
}

func NewGoalsHandler(h *Handlers, db *gorm.DB) *GoalsHandler {
	return &GoalsHandler{handlers: h, db: db}
}

func (h *GoalsHandler) SetupRoutes() {
	router := *h.handlers.router

	goalsRoutesHandler := goalsRoutesHandler{db: h.db}

	router.Route("/goals", func(route chi.Router) {
		route.Use(jwtauth.Verifier(jwt.InitAuthToken()))
		route.Use(jwtauth.Authenticator(jwt.InitAuthToken()))

		route.Post("/", goalsRoutesHandler.create)

	})

}
