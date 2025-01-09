package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/felix-Asante/pennyPilot-go-api/src/api/repositories"
	"github.com/felix-Asante/pennyPilot-go-api/src/api/services/authServices"
	customErrors "github.com/felix-Asante/pennyPilot-go-api/src/utils/errors"
	"gorm.io/gorm"
)

type authRoutesHandler struct {
	db *gorm.DB
}

func newAuthServices(db *gorm.DB) *authServices.AuthServices {
	newUserRepository := repositories.NewUsersRepository(db)

	authServices := authServices.NewAuthServices(newUserRepository)
	return authServices
}

func (handler *authRoutesHandler) signupHandler(w http.ResponseWriter, r *http.Request) {

	var request repositories.CreateUserRequest

	authServices := newAuthServices(handler.db)

	if err := customErrors.DecodeAndValidate(r, &request); err != nil {
		customErrors.RespondWithError(w, http.StatusBadRequest, customErrors.BadRequest, err.Error(), nil)
		return
	}

	_, err := authServices.Register(request)

	if err != nil {
		customErrors.RespondWithError(w, http.StatusBadRequest, customErrors.BadRequest, err.Error(), nil)
		return
	}
	w.WriteHeader(http.StatusCreated)
	jsonResponse := map[string]bool{"success": true}
	json.NewEncoder(w).Encode(jsonResponse)
}

func (h *authRoutesHandler) loginHandler(w http.ResponseWriter, r *http.Request) {
	authServices := newAuthServices(h.db)
	authServices.Login("email", "password")
}
