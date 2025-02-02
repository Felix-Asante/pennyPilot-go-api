package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/felix-Asante/pennyPilot-go-api/src/api/repositories"
	"github.com/felix-Asante/pennyPilot-go-api/src/api/services/financialObligationService"
	customErrors "github.com/felix-Asante/pennyPilot-go-api/src/utils/errors"
	"github.com/go-chi/jwtauth/v5"
	"gorm.io/gorm"
)

type financialObligationRoutesHandler struct {
	db *gorm.DB
}

func (fo *financialObligationRoutesHandler) new(w http.ResponseWriter, r *http.Request) {
	var request repositories.CreateFinancialObligationDto
	_, claims, error := jwtauth.FromContext(r.Context())

	if error != nil {
		customErrors.RespondWithError(w, http.StatusInternalServerError, customErrors.InternalServerError, error.Error(), nil)
		return
	}

	if err := customErrors.DecodeAndValidate(r, &request); err != nil {
		customErrors.RespondWithError(w, http.StatusBadRequest, customErrors.BadRequest, err.Error(), nil)
		return
	}

	userId := claims["id"].(string)

	fos := financialObligationService.NewFinancialObligationService(fo.db)
	newObligation, err := fos.Create(userId, request)

	if err != nil {
		customErrors.RespondWithError(w, http.StatusInternalServerError, customErrors.InternalServerError, err.Error(), nil)
		return
	}

	jsonResponse, _ := json.Marshal(newObligation)
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)

}
func (fo *financialObligationRoutesHandler) get(w http.ResponseWriter, r *http.Request) {
}
