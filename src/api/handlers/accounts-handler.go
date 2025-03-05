package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/felix-Asante/pennyPilot-go-api/src/api/repositories"
	accountsServices "github.com/felix-Asante/pennyPilot-go-api/src/api/services/accountsService"
	"github.com/felix-Asante/pennyPilot-go-api/src/utils/dates"
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

type allocateGoalsRequest struct {
	Goals []string `json:"goals" validate:"required"`
}

func newAccountServices(db *gorm.DB) *accountsServices.AccountsServices {
	accountServices := accountsServices.NewAccountServices(db)
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
	userId := claims["id"].(string)

	account := repositories.CreateAccountDto{
		UserID:          userId,
		Name:            request.Name,
		TargetBalance:   request.TargetBalance,
		AllocationPoint: request.AllocationPoint,
	}

	newAccount, statusCode, err := accountServices.Create(account)

	if err != nil {
		customErrors.RespondWithError(w, statusCode, customErrors.StatusCodes[statusCode], err.Error(), nil)
		return
	}

	jsonResponse, _ := json.Marshal(newAccount)
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

func (h *accountsRoutesHandler) allocate(w http.ResponseWriter, r *http.Request) {
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

	var request allocateGoalsRequest
	if err := customErrors.DecodeAndValidate(r, &request); err != nil {
		customErrors.RespondWithError(w, http.StatusBadRequest, customErrors.BadRequest, err.Error(), nil)
		return
	}

	if err := customErrors.ValidateUUIDs(request.Goals); err != nil {
		customErrors.RespondWithError(w, http.StatusBadRequest, customErrors.BadRequest, err.Error(), nil)
		return
	}

	userId := claims["id"].(string)

	accountService := newAccountServices(h.db)

	if status, error := accountService.AllocateToGoals(accountId, userId, request.Goals); error != nil {
		customErrors.RespondWithError(w, status, customErrors.InternalServerError, error.Error(), nil)
		return
	}

	jsonResponse, _ := json.Marshal(map[string]interface{}{"success": true})
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}

func (h *accountsRoutesHandler) getTransactions(w http.ResponseWriter, r *http.Request) {
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
	userId := claims["id"].(string)
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	query := r.URL.Query().Get("query")
	transactionType := r.URL.Query().Get("type")
	startDate, _ := dates.ParseDate(r.URL.Query().Get("start_date"))
	endDate, _ := dates.ParseDate(r.URL.Query().Get("end_date"))

	if page == 0 {
		page = 1
	}

	if pageSize == 0 {
		pageSize = 10
	}

	if query == "" {
		query = ""
	}

	if transactionType == "" {
		transactionType = ""
	}

	if !startDate.IsZero() && endDate.IsZero() {
		customErrors.RespondWithError(w, http.StatusBadRequest, customErrors.BadRequest, "end_date is required if start_date is provided", nil)
		return
	}
	if startDate.IsZero() && !endDate.IsZero() {
		customErrors.RespondWithError(w, http.StatusBadRequest, customErrors.BadRequest, "start_date is required if end_date is provided", nil)
		return
	}

	// if (startDate != "" && !dates.IsValidDate(startDate)) || (endDate != "" && !dates.IsValidDate(endDate)) {
	// 	customErrors.RespondWithError(w, http.StatusBadRequest, customErrors.BadRequest, "start_date and end_date must be valid dates", nil)
	// 	return
	// }

	accountService := newAccountServices(h.db)

	transactionArgs := repositories.GetAccountTransactions{
		AccountId: accountId,
		UserId:    userId,
		Page:      page,
		PageSize:  pageSize,
		Type:      repositories.TransactionType(transactionType),
		Query:     query,
		StartDate: &startDate,
		EndDate:   &endDate,
	}

	if transactions, status, error := accountService.GetTransactions(transactionArgs); error != nil {
		customErrors.RespondWithError(w, status, customErrors.InternalServerError, error.Error(), nil)
		return
	} else {
		jsonResponse, _ := json.Marshal(transactions)
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonResponse)
	}
}

func (h *accountsRoutesHandler) getGoals(w http.ResponseWriter, r *http.Request) {
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
	userId := claims["id"].(string)

	accountService := newAccountServices(h.db)

	if goals, status, error := accountService.GetGoals(accountId, userId); error != nil {
		customErrors.RespondWithError(w, status, customErrors.InternalServerError, error.Error(), nil)
		return
	} else {
		jsonResponse, _ := json.Marshal(goals)
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonResponse)
	}
}

func (h *accountsRoutesHandler) getExpenses(w http.ResponseWriter, r *http.Request) {
	accountId := chi.URLParam(r, "accountId")
	year, _ := strconv.Atoi(chi.URLParam(r, "year"))
	startDate, _ := dates.ParseDate(r.URL.Query().Get("start_date"))
	endDate, _ := dates.ParseDate(r.URL.Query().Get("end_date"))

	if year == 0 {
		year = time.Now().Year()
	}

	if err := uuid.Validate(accountId); err != nil {
		customErrors.RespondWithError(w, http.StatusBadRequest, customErrors.BadRequest, customErrors.InvalidAccountIDError, nil)
		return
	}

	if !startDate.IsZero() && endDate.IsZero() {
		customErrors.RespondWithError(w, http.StatusBadRequest, customErrors.BadRequest, "end_date is required if start_date is provided", nil)
		return
	}
	if startDate.IsZero() && !endDate.IsZero() {
		customErrors.RespondWithError(w, http.StatusBadRequest, customErrors.BadRequest, "start_date is required if end_date is provided", nil)
		return
	}

	expenseRepo := repositories.NewExpenseRepository(h.db)
	expenseArgs := repositories.GetAccountExpensesDto{
		Account:   accountId,
		StartDate: &startDate,
		EndDate:   &endDate,
		Year:      year,
	}

	expenses, error := expenseRepo.GetByAccount(expenseArgs)

	if error != nil {
		customErrors.RespondWithError(w, http.StatusInternalServerError, customErrors.InternalServerError, error.Error(), nil)
		return
	}


	jsonResponse, _ := json.Marshal(expenses)
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}

func (h *accountsRoutesHandler) getTotalExpenses(w http.ResponseWriter, r *http.Request) {
	accountId := chi.URLParam(r, "accountId")
	year, _ := strconv.Atoi(chi.URLParam(r, "year"))
	startDate, _ := dates.ParseDate(r.URL.Query().Get("start_date"))
	endDate, _ := dates.ParseDate(r.URL.Query().Get("end_date"))
	if year == 0 {
		year = time.Now().Year()
	}
	if err := uuid.Validate(accountId); err!= nil {
		customErrors.RespondWithError(w, http.StatusBadRequest, customErrors.BadRequest, customErrors.InvalidAccountIDError, nil)
		return
	}
	if!startDate.IsZero() && endDate.IsZero() {
		customErrors.RespondWithError(w, http.StatusBadRequest, customErrors.BadRequest, "end_date is required if start_date is provided", nil)
		return
	}
	if startDate.IsZero() &&!endDate.IsZero() {
		customErrors.RespondWithError(w, http.StatusBadRequest, customErrors.BadRequest, "start_date is required if end_date is provided", nil)
		return
	}
	expenseRepo := repositories.NewExpenseRepository(h.db)
	expenseArgs := repositories.GetAccountExpensesDto{
		Account:   accountId,
		StartDate: &startDate,
		EndDate:   &endDate,
		Year:      year,
	}
	totalExpenses, error := expenseRepo.GetTotalByAccount(expenseArgs)
	if error!= nil {
		customErrors.RespondWithError(w, http.StatusInternalServerError, customErrors.InternalServerError, error.Error(), nil)
		return
	}
	jsonResponse, _ := json.Marshal(map[string]interface{}{"total": totalExpenses})
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}

func (h *accountsRoutesHandler) transfer(w http.ResponseWriter, r *http.Request) {

}
