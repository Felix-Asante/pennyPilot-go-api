package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/felix-Asante/pennyPilot-go-api/src/api/repositories"
	goalsservice "github.com/felix-Asante/pennyPilot-go-api/src/api/services/goalsService"
	customErrors "github.com/felix-Asante/pennyPilot-go-api/src/utils/errors"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type goalsRoutesHandler struct {
	db *gorm.DB
}

func (g *goalsRoutesHandler) create(w http.ResponseWriter, r *http.Request) {
	var request repositories.CreateGoalDto
	_, claims, error := jwtauth.FromContext(r.Context())

	if error != nil {
		customErrors.RespondWithError(w, http.StatusInternalServerError, customErrors.InternalServerError, error.Error(), nil)
		return
	}

	if err := customErrors.DecodeAndValidate(r, &request); err != nil {
		customErrors.RespondWithError(w, http.StatusBadRequest, customErrors.BadRequest, err.Error(), nil)
		return
	}

	if err := customErrors.ValidateUUIDs([]string{request.Account}); err != nil {
		customErrors.RespondWithError(w, http.StatusBadRequest, customErrors.BadRequest, err.Error(), nil)
		return
	}

	goalsService := goalsservice.NewGoalsService(g.db)
	userId := claims["id"].(string)

	newGoal, statusCode, err := goalsService.Create(userId, request)

	if err != nil {
		customErrors.RespondWithError(w, statusCode, customErrors.StatusCodes[statusCode], err.Error(), nil)
		return
	}

	jsonResponse, _ := json.Marshal(newGoal)
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}

func (g *goalsRoutesHandler) update(w http.ResponseWriter, r *http.Request) {
	var request repositories.UpdateGoalDto
	_, claims, error := jwtauth.FromContext(r.Context())
	goalId := chi.URLParam(r, "goalId")

	if err := uuid.Validate(goalId); err != nil {
		customErrors.RespondWithError(w, http.StatusBadRequest, customErrors.BadRequest, "invalid goal id", nil)
		return
	}

	if error != nil {
		customErrors.RespondWithError(w, http.StatusInternalServerError, customErrors.InternalServerError, error.Error(), nil)
		return
	}

	if err := customErrors.DecodeAndValidate(r, &request); err != nil {
		customErrors.RespondWithError(w, http.StatusBadRequest, customErrors.BadRequest, err.Error(), nil)
		return
	}

	// if err := customErrors.ValidateUUIDs([]string{request.Account}); err != nil {
	// 	customErrors.RespondWithError(w, http.StatusBadRequest, customErrors.BadRequest, err.Error(), nil)
	// 	return
	// }

	goalsService := goalsservice.NewGoalsService(g.db)
	userId := claims["id"].(string)

	newGoal, statusCode, err := goalsService.Update(goalId, userId, request)

	if err != nil {
		customErrors.RespondWithError(w, statusCode, customErrors.StatusCodes[statusCode], err.Error(), nil)
		return
	}

	jsonResponse, _ := json.Marshal(newGoal)
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}

func (g *goalsRoutesHandler) delete(w http.ResponseWriter, r *http.Request) {
	_, claims, error := jwtauth.FromContext(r.Context())
	goalId := chi.URLParam(r, "goalId")

	if err := uuid.Validate(goalId); err != nil {
		customErrors.RespondWithError(w, http.StatusBadRequest, customErrors.BadRequest, "invalid goal id", nil)
		return
	}

	if error != nil {
		customErrors.RespondWithError(w, http.StatusInternalServerError, customErrors.InternalServerError, error.Error(), nil)
		return
	}

	goalsService := goalsservice.NewGoalsService(g.db)
	userId := claims["id"].(string)

	statusCode, err := goalsService.Delete(goalId, userId)

	if err != nil {
		customErrors.RespondWithError(w, statusCode, customErrors.StatusCodes[statusCode], err.Error(), nil)
		return
	}

	jsonResponse, _ := json.Marshal(map[string]bool{"success": true})
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}

func (g *goalsRoutesHandler) get(w http.ResponseWriter, r *http.Request) {
	_, claims, error := jwtauth.FromContext(r.Context())
	goalId := chi.URLParam(r, "goalId")

	if err := uuid.Validate(goalId); err != nil {
		customErrors.RespondWithError(w, http.StatusBadRequest, customErrors.BadRequest, "invalid goal id", nil)
		return
	}

	if error != nil {
		customErrors.RespondWithError(w, http.StatusInternalServerError, customErrors.InternalServerError, error.Error(), nil)
		return
	}

	goalsService := goalsservice.NewGoalsService(g.db)
	userId := claims["id"].(string)

	goal, statusCode, err := goalsService.Get(goalId, userId)

	if err != nil {
		customErrors.RespondWithError(w, statusCode, customErrors.StatusCodes[statusCode], err.Error(), nil)
		return
	}

	jsonResponse, _ := json.Marshal(goal)
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}
