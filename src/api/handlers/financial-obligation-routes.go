package handlers

import (
	"github.com/felix-Asante/pennyPilot-go-api/src/pkgs/jwt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
	"gorm.io/gorm"
)

type FinancialObligationHandler struct {
	handlers *Handlers
	db       *gorm.DB
}

func NewFinancialObligationHandler(h *Handlers, db *gorm.DB) *FinancialObligationHandler {
	return &FinancialObligationHandler{handlers: h, db: db}
}

func (h *FinancialObligationHandler) SetupRoutes() {
	router := *h.handlers.router

	routesHandler := financialObligationRoutesHandler{db: h.db}

	router.Route("/financial-obligations", func(route chi.Router) {

		route.Use(jwtauth.Verifier(jwt.InitAuthToken()))
		route.Use(jwtauth.Authenticator(jwt.InitAuthToken()))

		route.Post("/", routesHandler.new)
		route.Route("/{obligationId}", func(r chi.Router) {
			r.Get("/", routesHandler.get)
			r.Delete("/", routesHandler.delete)
		})

	})

}
