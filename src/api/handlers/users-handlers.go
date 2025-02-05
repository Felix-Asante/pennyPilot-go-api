package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	goalsservice "github.com/felix-Asante/pennyPilot-go-api/src/api/services/goalsService"
	"github.com/felix-Asante/pennyPilot-go-api/src/api/services/transactionsService"
	"github.com/felix-Asante/pennyPilot-go-api/src/api/services/usersServices"
	customErrors "github.com/felix-Asante/pennyPilot-go-api/src/utils/errors"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type usersRoutesHandler struct {
	db *gorm.DB
}

func newUserServices(db *gorm.DB) *usersServices.UsersServices {
	usersServices := usersServices.NewUsersServices(db)
	return usersServices
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
func (u *usersRoutesHandler) getGoals(w http.ResponseWriter, r *http.Request) {
	_, claims, error := jwtauth.FromContext(r.Context())
	if error != nil {
		customErrors.RespondWithError(w, http.StatusInternalServerError, customErrors.InternalServerError, error.Error(), nil)
		return
	}

	id := chi.URLParam(r, "userId")
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("limit"))

	if page == 0 {
		page = 1
	}

	if pageSize == 0 {
		pageSize = 10
	}

	if err := uuid.Validate(id); err != nil {
		customErrors.RespondWithError(w, http.StatusBadRequest, customErrors.BadRequest, "invalid account id", nil)
		return
	}

	userId := claims["id"].(string)
	if id != userId {
		customErrors.RespondWithError(w, http.StatusForbidden, customErrors.ForbiddenError, "Access to this account denied", nil)
		return
	}
	goalService := goalsservice.NewGoalsService(u.db)
	goals, err := goalService.FindUserGoals(id, page, pageSize)

	if err != nil {
		customErrors.RespondWithError(w, http.StatusInternalServerError, customErrors.InternalServerError, error.Error(), nil)
		return
	}

	jsonResponse, _ := json.Marshal(goals)
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}

func (u *usersRoutesHandler) getTransactions(w http.ResponseWriter, r *http.Request) {
	_, claims, error := jwtauth.FromContext(r.Context())
	if error != nil {
		customErrors.RespondWithError(w, http.StatusInternalServerError, customErrors.InternalServerError, error.Error(), nil)
		return
	}

	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("limit"))

	if page == 0 {
		page = 1
	}

	if pageSize == 0 {
		pageSize = 10
	}

	userId := claims["id"].(string)

	transactionService := transactionsService.NewTransactionsService(u.db)
	transactions, err := transactionService.FindAllByUserId(userId, page, pageSize)

	if err != nil {
		customErrors.RespondWithError(w, http.StatusInternalServerError, customErrors.InternalServerError, error.Error(), nil)
		return
	}

	jsonResponse, _ := json.Marshal(transactions)
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}

func (u *usersRoutesHandler) getMe(w http.ResponseWriter, r *http.Request) {
	_, claims, error := jwtauth.FromContext(r.Context())
	if error != nil {
		customErrors.RespondWithError(w, http.StatusInternalServerError, customErrors.InternalServerError, error.Error(), nil)
		return
	}

	userId := claims["id"].(string)
	usersService := newUserServices(u.db)
	user, err := usersService.Me(userId)

	if err != nil {
		customErrors.RespondWithError(w, http.StatusInternalServerError, customErrors.InternalServerError, error.Error(), nil)
		return
	}

	jsonResponse, _ := json.Marshal(user)
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}
