package handlers

import (
	"net/http"

	"github.com/felix-Asante/pennyPilot-go-api/src/api/repositories"
	"github.com/felix-Asante/pennyPilot-go-api/src/api/services/authServices"
	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

type AuthHandler struct {
	handlers *Handlers
	db       *gorm.DB
}

func NewAuthHandler(h *Handlers, db *gorm.DB) *AuthHandler {
	return &AuthHandler{handlers: h, db: db}
}

func (h *AuthHandler) SetupRoutes() {
	router := *h.handlers.router
	newUserRepository := repositories.NewAuthRepository(h.db)

	authServices := authServices.NewAuthServices(newUserRepository)

	router.Route("/auth", func(route chi.Router) {

		route.Get("/login", func(w http.ResponseWriter, r *http.Request) {
			authServices.Login("email", "password")
		})
	})

}
