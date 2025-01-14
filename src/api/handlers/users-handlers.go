package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/felix-Asante/pennyPilot-go-api/src/api/repositories"
	"github.com/felix-Asante/pennyPilot-go-api/src/api/services/usersServices"
	customErrors "github.com/felix-Asante/pennyPilot-go-api/src/utils/errors"
	"github.com/go-chi/jwtauth/v5"
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

func (u *usersRoutesHandler) getAccounts(w http.ResponseWriter, r *http.Request) {
	_, claims, error := jwtauth.FromContext(r.Context())
	if error != nil {
		customErrors.RespondWithError(w, http.StatusInternalServerError, customErrors.InternalServerError, error.Error(), nil)
		return
	}

	userId := claims["id"].(string)
	accountService := newAccountServices(u.db)
	accounts, err := accountService.FindUserAccounts(userId)

	if err != nil {
		customErrors.RespondWithError(w, http.StatusInternalServerError, customErrors.InternalServerError, error.Error(), nil)
		return
	}

	jsonResponse, _ := json.Marshal(accounts)
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}
