package handlers

import (
	"net/http"

	"github.com/felix-Asante/pennyPilot-go-api/src/api/repositories"
	"github.com/felix-Asante/pennyPilot-go-api/src/api/services/usersServices"
	"gorm.io/gorm"
)

type usersRoutesHandler struct {
	db *gorm.DB
}

func newUserServices(db *gorm.DB) *usersServices.UsersServices {
	newUserRepository := repositories.NewUsersRepository(db)

	usersServices := usersServices.NewUsersServices(newUserRepository)
	return usersServices
}

func (h *usersRoutesHandler) getUser(w http.ResponseWriter, r *http.Request) {
	usersServices := newUserServices(h.db)
	usersServices.CreateNewUser("id")
}
