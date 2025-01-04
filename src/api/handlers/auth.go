package handlers

import "net/http"

type AuthHandler struct {
	handlers *Handlers
}

func NewAuthHandler(h *Handlers) *AuthHandler {
	return &AuthHandler{handlers: h}
}

func (h *AuthHandler) SetupRoutes() {
	router := *h.handlers.router

	router.Get("/auth", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello Auth!"))
	})
}
