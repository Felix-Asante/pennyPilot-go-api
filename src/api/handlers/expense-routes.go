package handlers

import (
	"github.com/felix-Asante/pennyPilot-go-api/src/pkgs/jwt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
	"gorm.io/gorm"
)

type ExpenseRoutesHandler struct {
	handlers *Handlers
	db       *gorm.DB
}

func NewExpenseHandler(h *Handlers, db *gorm.DB) *ExpenseRoutesHandler {
	return &ExpenseRoutesHandler{handlers: h, db: db}
}

func (h *ExpenseRoutesHandler) SetupRoutes() {
	router := *h.handlers.router

	expenseRoutesHandler := expenseRoutesHandler{db: h.db}
	router.Route("/expenses", func(route chi.Router) {

		route.Use(jwtauth.Verifier(jwt.InitAuthToken()))
		route.Use(jwtauth.Authenticator(jwt.InitAuthToken()))

		route.Post("/", expenseRoutesHandler.new)
		route.Post("/category", expenseRoutesHandler.newExpenseCategory)

	})

}
