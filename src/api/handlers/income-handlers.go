package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/felix-Asante/pennyPilot-go-api/src/api/repositories"
	incomeservices "github.com/felix-Asante/pennyPilot-go-api/src/api/services/incomeServices"
	customErrors "github.com/felix-Asante/pennyPilot-go-api/src/utils/errors"
	"github.com/go-chi/jwtauth/v5"
	"gorm.io/gorm"
)

type incomeRoutesHandler struct {
	db *gorm.DB
}

type createIncomeRequest struct {
	Amount       float64    `json:"amount" validate:"required,min=1"`
	DateReceived *time.Time `json:"date_received" validate:"required"`
	IncomeType   string     `json:"type" validate:"required"`
	Frequency    string     `json:"frequency" validate:"required"`
}

func newIncomeServices(db *gorm.DB) *incomeservices.IncomeServices {
	accountServices := incomeservices.NewIncomeServices(db)
	return accountServices
}

func (h *incomeRoutesHandler) create(w http.ResponseWriter, r *http.Request) {
	var request createIncomeRequest
	_, claims, error := jwtauth.FromContext(r.Context())

	if error != nil {
		fmt.Print(error)
		customErrors.RespondWithError(w, http.StatusInternalServerError, customErrors.InternalServerError, error.Error(), nil)
		return
	}

	if err := customErrors.DecodeAndValidate(r, &request); err != nil {
		customErrors.RespondWithError(w, http.StatusBadRequest, customErrors.BadRequest, err.Error(), nil)
		return
	}

	incomeService := newIncomeServices(h.db)

	income := repositories.CreateIncomeDto{
		User:         claims["id"].(string),
		Amount:       request.Amount,
		Type:         request.IncomeType,
		Frequency:    request.Frequency,
		DateReceived: request.DateReceived,
	}

	newIncome, status, error := incomeService.Create(income)

	if error != nil {
		customErrors.RespondWithError(w, status, customErrors.InternalServerError, error.Error(), nil)
		return
	}

	jsonResponse, _ := json.Marshal(newIncome)
	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}
