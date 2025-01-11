package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/felix-Asante/pennyPilot-go-api/src/api/repositories"
	"github.com/felix-Asante/pennyPilot-go-api/src/api/services/authServices"
	"github.com/felix-Asante/pennyPilot-go-api/src/pkgs/jwt"
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

	user, err := authServices.Register(request)

	if err != nil {
		customErrors.RespondWithError(w, http.StatusBadRequest, customErrors.BadRequest, err.Error(), nil)
		return
	}

	claims := map[string]interface{}{"id": user.ID, "email": user.Email}
	jwtService := jwt.NewJWTService()
	_, token, err := jwtService.Encode(claims)

	if err != nil {
		customErrors.RespondWithError(w, http.StatusBadRequest, customErrors.InternalServerError, err.Error(), nil)
		return
	}

	w.WriteHeader(http.StatusCreated)
	jsonResponse := map[string]interface{}{"user": user, "token": token}
	json.NewEncoder(w).Encode(jsonResponse)
}

func (handler *authRoutesHandler) loginHandler(w http.ResponseWriter, r *http.Request) {
	var request authServices.LoginRequest

	authServices := newAuthServices(handler.db)

	if err := customErrors.DecodeAndValidate(r, &request); err != nil {
		customErrors.RespondWithError(w, http.StatusBadRequest, customErrors.BadRequest, err.Error(), nil)
		return
	}

	user, err := authServices.Login(request.Email, request.Password)

	if err != nil {
		customErrors.RespondWithError(w, http.StatusBadRequest, customErrors.BadRequest, err.Error(), nil)
		return
	}

	claims := map[string]interface{}{"id": user.ID, "email": user.Email}
	jwtService := jwt.NewJWTService()
	_, token, err := jwtService.Encode(claims)

	if err != nil {
		customErrors.RespondWithError(w, http.StatusBadRequest, customErrors.InternalServerError, err.Error(), nil)
		return
	}

	jsonResponse := map[string]interface{}{"user": user, "token": token}
	json.NewEncoder(w).Encode(jsonResponse)
}

func (handler *authRoutesHandler) resetPasswordHandler(w http.ResponseWriter, r *http.Request) {

}

func (h *authRoutesHandler) requestResetPasswordCode(w http.ResponseWriter, r *http.Request) {
	var request authServices.ForgetPasswordRequest

	authServices := newAuthServices(h.db)

	if err := customErrors.DecodeAndValidate(r, &request); err != nil {
		customErrors.RespondWithError(w, http.StatusBadRequest, customErrors.BadRequest, err.Error(), nil)
		return
	}

	token, err := authServices.ResetPasswordRequest(request.Email)

	if err != nil {
		customErrors.RespondWithError(w, http.StatusBadRequest, customErrors.BadRequest, err.Error(), nil)
		return
	}

	jsonResponse := map[string]interface{}{"code": token}
	json.NewEncoder(w).Encode(jsonResponse)
}
