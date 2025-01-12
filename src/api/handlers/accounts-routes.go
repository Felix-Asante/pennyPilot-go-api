package handlers

import (
	"github.com/go-chi/chi/v5"
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

		route.Post("/", accountRoutesHandler.new)
		route.Get("/{accountsId}", accountRoutesHandler.get)
		route.Put("/{accountsId}", accountRoutesHandler.update)
		route.Delete("/{accountsId}", accountRoutesHandler.delete)
		route.Post("/transfer", accountRoutesHandler.transfer)

	})

}
