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
	_, claims, error := jwtauth.FromContext(r.Context())

	if error != nil {
		customErrors.RespondWithError(w, http.StatusInternalServerError, customErrors.InternalServerError, error.Error(), nil)
		return
	}
	userId := claims["id"].(string)
	id := r.URL.Query().Get("id")

	if id !=userId {
		customErrors.RespondWithError(w, http.StatusUnauthorized, customErrors.ForbiddenError, "Unauthorized", nil)
		return
	}

	foRepo := repositories.NewFinancialObligationsRepository(fo.db)
	fob,error := foRepo.FindByIDAndUserID(id,userId)
	if error!= nil {
		customErrors.RespondWithError(w, http.StatusInternalServerError, customErrors.InternalServerError,"Internal server error", nil)
		return
	}
	jsonResponse, _ := json.Marshal(fob)
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}


func (fo *financialObligationRoutesHandler) delete(w http.ResponseWriter, r *http.Request) {
	_, claims, error := jwtauth.FromContext(r.Context())
	if error!= nil {
		customErrors.RespondWithError(w, http.StatusInternalServerError, customErrors.InternalServerError, error.Error(), nil)
		return
	}
	userId := claims["id"].(string)
	id := r.URL.Query().Get("id")
	if id!=userId {
		customErrors.RespondWithError(w, http.StatusUnauthorized, customErrors.ForbiddenError, "Unauthorized", nil)
		return
	}

	foRepo := repositories.NewFinancialObligationsRepository(fo.db)
	fob,error :=foRepo.FindByIDAndUserID(id,userId)
	if error!= nil {
		customErrors.RespondWithError(w, http.StatusInternalServerError, customErrors.InternalServerError, error.Error(), nil)
		return
	}

	if fob==nil {
		customErrors.RespondWithError(w, http.StatusNotFound, customErrors.NotFoundError, "Financial obligation not found", nil)
		return
	}

	if status,error := foRepo.Remove(id); error!= nil || !status {
		customErrors.RespondWithError(w, http.StatusInternalServerError, customErrors.InternalServerError, customErrors.InternalServerError, nil)
		return
	}

	jsonResponse, _ := json.Marshal(map[string]interface{}{"message": "Financial obligation deleted successfully"})
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
	w.Write(jsonResponse)
}