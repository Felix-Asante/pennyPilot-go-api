package handlers

import (
	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

type Handlers struct {
	router *chi.Router
}

type Handler interface {
	SetupRoutes()
}

func (h *Handlers) SetupRoutes(db *gorm.DB) {
	authHandler := NewAuthHandler(h, db)
	usersHandler := NewUsersHandler(h, db)
	accountsHandler := NewAccountHandler(h, db)
	incomesHandler := NewIncomeHandler(h, db)
	goalsHandler := NewGoalsHandler(h, db)
	expenseHandler := NewExpenseHandler(h, db)
	handlers := []Handler{authHandler, usersHandler, accountsHandler, incomesHandler, goalsHandler, expenseHandler}
	createAllRoutes(handlers)
}

func NewHandlers(router *chi.Router) *Handlers {

	return &Handlers{router}
}

func createAllRoutes(handlers []Handler) {

	for _, handler := range handlers {
		handler.SetupRoutes()
	}
}
