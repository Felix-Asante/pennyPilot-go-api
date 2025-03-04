package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/felix-Asante/pennyPilot-go-api/src/api/repositories"
	goalsservice "github.com/felix-Asante/pennyPilot-go-api/src/api/services/goalsService"
	"github.com/felix-Asante/pennyPilot-go-api/src/api/services/transactionsService"
	"github.com/felix-Asante/pennyPilot-go-api/src/api/services/usersServices"
	"github.com/felix-Asante/pennyPilot-go-api/src/utils/dates"
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
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	query := r.URL.Query().Get("query")
	sort := r.URL.Query().Get("sort")

	if page == 0 {
		page = 1
	}

	if pageSize == 0 {
		pageSize = 10
	}

	if query == "" {
		query = ""
	}

	if sort == "" {
		sort = ""
	}

	accountService := newAccountServices(u.db)
	queries := repositories.AccountQueries{
		Page:  page,
		Limit: pageSize,
		Query: query,
		Sort:  sort,
	}

	accounts, err := accountService.FindUserAccounts(userId, queries)

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

func (u *usersRoutesHandler) getExpenses(w http.ResponseWriter, r *http.Request) {
	_, claims, error := jwtauth.FromContext(r.Context())
	if error != nil {
		customErrors.RespondWithError(w, http.StatusInternalServerError, customErrors.InternalServerError, error.Error(), nil)
		return
	}

	id := chi.URLParam(r, "userId")
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	year, _ := strconv.Atoi(chi.URLParam(r, "year"))
	startDate, _ := dates.ParseDate(r.URL.Query().Get("start_date"))
	endDate, _ := dates.ParseDate(r.URL.Query().Get("end_date"))

	if year == 0 {
		year = time.Now().Year()
	}
	if !startDate.IsZero() && endDate.IsZero() {
		customErrors.RespondWithError(w, http.StatusBadRequest, customErrors.BadRequest, "end_date is required if start_date is provided", nil)
		return
	}
	if startDate.IsZero() && !endDate.IsZero() {
		customErrors.RespondWithError(w, http.StatusBadRequest, customErrors.BadRequest, "start_date is required if end_date is provided", nil)
		return
	}

	if err := uuid.Validate(id); err != nil {
		customErrors.RespondWithError(w, http.StatusBadRequest, customErrors.BadRequest, "invalid user id", nil)
		return
	}

	userId := claims["id"].(string)
	if id != userId {
		customErrors.RespondWithError(w, http.StatusForbidden, customErrors.ForbiddenError, "Access to this user denied", nil)
		return
	}

	expenseRepo := repositories.NewExpenseRepository(u.db)
	paginationOptions := &repositories.PaginationOptions{
		Page:  page,
		Limit: pageSize,
	}
	paginationOptions.SetDefaultValues()

	expenses, err := expenseRepo.Get(repositories.GetExpenseDto{
		User:       userId,
		Year:       year,
		StartDate:  &startDate,
		EndDate:    &endDate,
		Pagination: paginationOptions,
	})

	if err != nil {
		customErrors.RespondWithError(w, http.StatusInternalServerError, customErrors.InternalServerError, error.Error(), nil)
		return
	}

	jsonResponse, _ := json.Marshal(expenses)
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}

func (u *usersRoutesHandler) getTotalExpenses(w http.ResponseWriter, r *http.Request) {
	_, claims, error := jwtauth.FromContext(r.Context())
	if error != nil {
		customErrors.RespondWithError(w, http.StatusInternalServerError, customErrors.InternalServerError, error.Error(), nil)
		return
	}

	id := chi.URLParam(r, "userId")

	userId := claims["id"].(string)
	if err := uuid.Validate(id); err != nil {
		customErrors.RespondWithError(w, http.StatusBadRequest, customErrors.BadRequest, "invalid user id", nil)
		return
	}

	if id != userId {
		customErrors.RespondWithError(w, http.StatusForbidden, customErrors.ForbiddenError, "Access to this user denied", nil)
		return
	}

	expenseRepo := repositories.NewExpenseRepository(u.db)
	total, err := expenseRepo.GetUserTotal(userId)
	if err != nil {
		customErrors.RespondWithError(w, http.StatusInternalServerError, customErrors.InternalServerError, error.Error(), nil)
		return
	}

	jsonResponse, _ := json.Marshal(map[string]interface{}{"total_expense": total})
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}

func (u *usersRoutesHandler) getExpenseCategories(w http.ResponseWriter, r *http.Request) {
	_, claims, error := jwtauth.FromContext(r.Context())
	if error!= nil {
		customErrors.RespondWithError(w, http.StatusInternalServerError, customErrors.InternalServerError, error.Error(), nil)
		return
	}
	id := chi.URLParam(r, "userId")
	userId := claims["id"].(string)
	if err := uuid.Validate(id); err!= nil {
		customErrors.RespondWithError(w, http.StatusBadRequest, customErrors.BadRequest, "invalid user id", nil)
		return
	}
	if id!= userId {
		customErrors.RespondWithError(w, http.StatusForbidden, customErrors.ForbiddenError, "Access to this ressource denied", nil)
		return
	}
	expenseRepo := repositories.NewExpenseCategoryRepository(u.db)
	categories, err := expenseRepo.FindAllByUserID(id)
	if err!= nil {
		customErrors.RespondWithError(w, http.StatusInternalServerError, customErrors.InternalServerError, error.Error(), nil)
		return
	}
	jsonResponse, _ := json.Marshal(categories)
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}

func (fo *usersRoutesHandler) getFinancialObligations(w http.ResponseWriter, r *http.Request) {
	_, claims, error := jwtauth.FromContext(r.Context())
	if error!= nil {
		customErrors.RespondWithError(w, http.StatusInternalServerError, customErrors.InternalServerError, error.Error(), nil)
		return
	}
	userId := claims["id"].(string)
	page,_ := strconv.Atoi(r.URL.Query().Get("page"))
	limit,_ := strconv.Atoi(r.URL.Query().Get("limit"))

	foRepo := repositories.NewFinancialObligationsRepository(fo.db)

	fosArgs := repositories.GetUserFiancialObligationDto{
		User: userId,
		Pagination: &repositories.PaginationOptions{
			Page: page,
			Limit: limit,
			Sort: "desc",
			Query: "",
		},
	}
	fosArgs.Pagination.SetDefaultValues()
	
	fobs,error := foRepo.FindAllUser(fosArgs)
	if error!= nil {
		customErrors.RespondWithError(w, http.StatusInternalServerError, customErrors.InternalServerError,"Internal server error", nil)
		return
	}
	jsonResponse, _ := json.Marshal(fobs)
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}

func (u *usersRoutesHandler) totalSavings(w http.ResponseWriter, r *http.Request) {
	_, claims, error := jwtauth.FromContext(r.Context())
	if error!= nil {
		customErrors.RespondWithError(w, http.StatusInternalServerError, customErrors.InternalServerError, error.Error(), nil)
		return
	}
	id := chi.URLParam(r, "userId")
	userId := claims["id"].(string)
	if err := uuid.Validate(id); err!= nil {
		customErrors.RespondWithError(w, http.StatusBadRequest, customErrors.BadRequest, "invalid user id", nil)
		return
	}
	if id!= userId {
		customErrors.RespondWithError(w, http.StatusForbidden, customErrors.ForbiddenError, "Access to this ressource denied", nil)
		return
	}

	accountRepo := repositories.NewAccountsRepository(u.db)
	totalSavings, error := accountRepo.FindTotalSavings(id)
	if error!= nil {
		customErrors.RespondWithError(w, http.StatusInternalServerError, customErrors.InternalServerError, error.Error(), nil)
		return
	}

	jsonResponse := map[string]float64{"total_savings": totalSavings}
	jsonData, _ := json.Marshal(jsonResponse)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}