package handlers

import (
	"github.com/felix-Asante/pennyPilot-go-api/src/pkgs/jwt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
	"gorm.io/gorm"
)

type AccountsHandler struct {
	handlers *Handlers
	db       *gorm.DB
}

func NewAccountHandler(h *Handlers, db *gorm.DB) *AccountsHandler {
	return &AccountsHandler{handlers: h, db: db}
}

func (h *AccountsHandler) SetupRoutes() {
	router := *h.handlers.router

	accountRoutesHandler := accountsRoutesHandler{db: h.db}

	router.Route("/accounts", func(route chi.Router) {

		route.Use(jwtauth.Verifier(jwt.InitAuthToken()))
		route.Use(jwtauth.Authenticator(jwt.InitAuthToken()))

		route.Post("/", accountRoutesHandler.new)
		route.Route("/{accountId}", func(r chi.Router) {
			r.Get("/", accountRoutesHandler.get)
			r.Get("/transactions", accountRoutesHandler.getTransactions)
			r.Get("/goals", accountRoutesHandler.getGoals)
			r.Put("/", accountRoutesHandler.update)
			r.Delete("/", accountRoutesHandler.delete)
			r.Put("/add-balance", accountRoutesHandler.updateBalance)
			r.Put("/allocate", accountRoutesHandler.allocate)
			r.Get("/expenses", accountRoutesHandler.getExpenses)
			r.Get("/expenses/total", accountRoutesHandler.getTotalExpenses)
		})

		route.Post("/transfer", accountRoutesHandler.transfer)

	})

}
