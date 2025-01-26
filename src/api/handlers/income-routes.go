package handlers

import (
	"github.com/felix-Asante/pennyPilot-go-api/src/pkgs/jwt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
	"gorm.io/gorm"
)

type IncomeHandler struct {
	handlers *Handlers
	db       *gorm.DB
}

func NewIncomeHandler(h *Handlers, db *gorm.DB) *IncomeHandler {
	return &IncomeHandler{handlers: h, db: db}
}

func (h *IncomeHandler) SetupRoutes() {
	router := *h.handlers.router

	incomeRoutesHandler := incomeRoutesHandler{db: h.db}

	router.Route("/incomes", func(route chi.Router) {
		route.Use(jwtauth.Verifier(jwt.InitAuthToken()))
		route.Use(jwtauth.Authenticator(jwt.InitAuthToken()))

		route.Post("/", incomeRoutesHandler.create)
		route.Put("/{incomeId}", incomeRoutesHandler.update)
		route.Get("/{incomeId}", incomeRoutesHandler.get)
		route.Put("/allocate", incomeRoutesHandler.allocate)

	})

}
