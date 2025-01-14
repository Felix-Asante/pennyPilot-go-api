package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/felix-Asante/pennyPilot-go-api/src/api/repositories"
	accountsServices "github.com/felix-Asante/pennyPilot-go-api/src/api/services/accountsService"
	customErrors "github.com/felix-Asante/pennyPilot-go-api/src/utils/errors"
	"github.com/go-chi/jwtauth/v5"
	"gorm.io/gorm"
)

type accountsRoutesHandler struct {
	db *gorm.DB
}

type createAccountRequest struct {
	Name            string  `json:"name" validate:"required,min=2"`
	TargetBalance   float64 `json:"target_balance" validate:"required,min=0.00"`
	AllocationPoint float64 `json:"allocation_point" validate:"required,min=1,max=100"`
}

func newAccountServices(db *gorm.DB) *accountsServices.AccountsServices {
	newAccountRepository := repositories.NewAccountsRepository(db)

	accountServices := accountsServices.NewAccountServices(newAccountRepository)
	return accountServices
}

func (h *accountsRoutesHandler) new(w http.ResponseWriter, r *http.Request) {

	var request createAccountRequest
	_, claims, error := jwtauth.FromContext(r.Context())

	if error != nil {
		customErrors.RespondWithError(w, http.StatusInternalServerError, customErrors.InternalServerError, error.Error(), nil)
		return
	}

	if err := customErrors.DecodeAndValidate(r, &request); err != nil {
		customErrors.RespondWithError(w, http.StatusBadRequest, customErrors.BadRequest, err.Error(), nil)
		return
	}

	accountServices := newAccountServices(h.db)

	account := repositories.CreateAccountDto{
		UserID:          claims["id"].(string),
		Name:            request.Name,
		TargetBalance:   request.TargetBalance,
		AllocationPoint: request.AllocationPoint,
	}

	newAccount, err := accountServices.Create(account)
	jsonResponse, _ := json.Marshal(newAccount)

	if err != nil {
		customErrors.RespondWithError(w, http.StatusInternalServerError, customErrors.InternalServerError, err.Error(), nil)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}

func (h *accountsRoutesHandler) get(w http.ResponseWriter, r *http.Request) {

}

func (h *accountsRoutesHandler) update(w http.ResponseWriter, r *http.Request) {

}

func (h *accountsRoutesHandler) delete(w http.ResponseWriter, r *http.Request) {

}

func (h *accountsRoutesHandler) transfer(w http.ResponseWriter, r *http.Request) {

}
