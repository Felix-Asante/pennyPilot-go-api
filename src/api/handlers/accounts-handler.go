package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/felix-Asante/pennyPilot-go-api/src/api/repositories"
	accountsServices "github.com/felix-Asante/pennyPilot-go-api/src/api/services/accountsService"
	customErrors "github.com/felix-Asante/pennyPilot-go-api/src/utils/errors"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
	"github.com/google/uuid"
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
type addBalanceRequest struct {
	Amount float64 `json:"amount" validate:"required,min=1"`
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

	newAccount, statusCode, err := accountServices.Create(account)
	jsonResponse, _ := json.Marshal(newAccount)

	if err != nil {
		customErrors.RespondWithError(w, statusCode, customErrors.StatusCodes[statusCode], err.Error(), nil)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}

func (h *accountsRoutesHandler) get(w http.ResponseWriter, r *http.Request) {

	accountId := chi.URLParam(r, "accountId")

	if err := uuid.Validate(accountId); err != nil {
		customErrors.RespondWithError(w, http.StatusBadRequest, customErrors.BadRequest, customErrors.InvalidAccountIDError, nil)
		return
	}

	_, claims, error := jwtauth.FromContext(r.Context())

	if error != nil {
		customErrors.RespondWithError(w, http.StatusInternalServerError, customErrors.InternalServerError, error.Error(), nil)
		return
	}

	accountsServices := newAccountServices(h.db)
	account, statusCode, err := accountsServices.Find(accountId, claims["id"].(string))

	if err != nil {
		customErrors.RespondWithError(w, statusCode, customErrors.StatusCodes[statusCode], err.Error(), nil)
		return
	}

	jsonResponse, _ := json.Marshal(account)
	w.WriteHeader(statusCode)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)

}

func (h *accountsRoutesHandler) updateBalance(w http.ResponseWriter, r *http.Request) {

	accountId := chi.URLParam(r, "accountId")

	if err := uuid.Validate(accountId); err != nil {
		customErrors.RespondWithError(w, http.StatusBadRequest, customErrors.BadRequest, customErrors.InvalidAccountIDError, nil)
		return
	}
	_, claims, error := jwtauth.FromContext(r.Context())

	if error != nil {
		customErrors.RespondWithError(w, http.StatusInternalServerError, customErrors.InternalServerError, error.Error(), nil)
		return
	}

	var request addBalanceRequest

	if err := customErrors.DecodeAndValidate(r, &request); err != nil {
		customErrors.RespondWithError(w, http.StatusBadRequest, customErrors.BadRequest, err.Error(), nil)
		return
	}

	accountsServices := newAccountServices(h.db)
	updateAmount := request.Amount
	userId := claims["id"].(string)

	_, statusCode, err := accountsServices.UpdateBalance(accountId, updateAmount, userId)
	if err != nil {
		customErrors.RespondWithError(w, statusCode, customErrors.StatusCodes[statusCode], err.Error(), nil)
		return
	}

	jsonResponse, _ := json.Marshal(map[string]interface{}{"success": true})
	w.WriteHeader(statusCode)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)

}

func (h *accountsRoutesHandler) update(w http.ResponseWriter, r *http.Request) {

}

func (h *accountsRoutesHandler) delete(w http.ResponseWriter, r *http.Request) {
	accountId := chi.URLParam(r, "accountId")

	if err := uuid.Validate(accountId); err != nil {
		customErrors.RespondWithError(w, http.StatusBadRequest, customErrors.BadRequest, customErrors.InvalidAccountIDError, nil)
		return
	}
	_, claims, error := jwtauth.FromContext(r.Context())

	if error != nil {
		customErrors.RespondWithError(w, http.StatusInternalServerError, customErrors.InternalServerError, error.Error(), nil)
		return
	}

	accountsServices := newAccountServices(h.db)
	userId := claims["id"].(string)
	statusCode, error := accountsServices.Remove(accountId, userId)

	if error != nil {
		customErrors.RespondWithError(w, statusCode, customErrors.StatusCodes[statusCode], error.Error(), nil)
		return
	}

	jsonResponse, _ := json.Marshal(map[string]interface{}{"success": true})
	w.WriteHeader(statusCode)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)

}

func (h *accountsRoutesHandler) transfer(w http.ResponseWriter, r *http.Request) {

}
