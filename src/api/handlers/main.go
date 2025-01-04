package handlers

import "github.com/go-chi/chi/v5"

type Handlers struct {
	router *chi.Router
}

type Handler interface {
	SetupRoutes()
}

func (h *Handlers) SetupRoutes() {
	authHandler := NewAuthHandler(h)
	usersHandler := NewUsersHandler(h)
	handlers := []Handler{authHandler, usersHandler}
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
