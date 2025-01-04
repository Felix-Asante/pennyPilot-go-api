package handlers

import "net/http"

type UsersHandler struct {
	handlers *Handlers
}

func NewUsersHandler(h *Handlers) *UsersHandler {
	return &UsersHandler{handlers: h}
}

func (h *UsersHandler) SetupRoutes() {
	router := *h.handlers.router

	router.Get("/users", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello Users!"))
	})
}
