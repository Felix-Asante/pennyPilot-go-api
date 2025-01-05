package handlers

import (
	"net/http"

	"github.com/felix-Asante/pennyPilot-go-api/src/api/repositories"
	"github.com/felix-Asante/pennyPilot-go-api/src/api/services/usersServices"
	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

type UsersHandler struct {
	handlers *Handlers
	db       *gorm.DB
}

func NewUsersHandler(h *Handlers, db *gorm.DB) *UsersHandler {
	return &UsersHandler{handlers: h, db: db}
}

func (h *UsersHandler) SetupRoutes() {
	router := *h.handlers.router
	newUserRepository := repositories.NewUsersRepository(h.db)

	usersServices := usersServices.NewUsersServices(newUserRepository)

	router.Route("/users", func(route chi.Router) {

		route.Get("/{usersId}", func(w http.ResponseWriter, r *http.Request) {
			userId := chi.URLParam(r, "usersId")
			usersServices.CreateNewUser(userId)
		})
	})
}
