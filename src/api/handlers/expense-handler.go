package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/felix-Asante/pennyPilot-go-api/src/api/repositories"
	customErrors "github.com/felix-Asante/pennyPilot-go-api/src/utils/errors"
	"github.com/go-chi/jwtauth/v5"
	"gorm.io/gorm"
)

type expenseRoutesHandler struct {
	db *gorm.DB
}

type CreateExpenseCategoryRequest struct {
	Name string `json:"name" validate:"required,min=2"`
	Icon string `json:"icon" validate:"required"`
}



func (h *expenseRoutesHandler) new(w http.ResponseWriter, r *http.Request) {
	var request repositories.CreateExpenseDto
	_, claims, error := jwtauth.FromContext(r.Context())

	if error != nil {
		customErrors.RespondWithError(w, http.StatusInternalServerError, customErrors.InternalServerError, error.Error(), nil)
		return
	}

	if err := customErrors.DecodeAndValidate(r, &request); err != nil {
		customErrors.RespondWithError(w, http.StatusBadRequest, customErrors.BadRequest, err.Error(), nil)
		return
	}

	tx := h.db.Begin()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	userId := claims["id"].(string)
	userRepo := repositories.NewUsersRepository(tx)
	accountRepo := repositories.NewAccountsRepository(tx)
	user, err := userRepo.FindUserById(userId)

	if err != nil {
		tx.Rollback()
		customErrors.RespondWithError(w, http.StatusInternalServerError, customErrors.InternalServerError, err.Error(), nil)
		return
	}

	account := &repositories.Accounts{}

	if request.Account != "" {

		account, err = creditAccount(request.Account, userId, request.Amount, tx)
		if err != nil {
			tx.Rollback()
			customErrors.RespondWithError(w, http.StatusInternalServerError, customErrors.InternalServerError, err.Error(), nil)
			return
		}
	} else {

		if user.TotalIncome < request.Amount {
			tx.Rollback()
			customErrors.RespondWithError(w, http.StatusBadRequest, customErrors.BadRequest, "insufficient funds", nil)
			return
		}
		user.TotalIncome -= request.Amount
	}

	expenseRepo := repositories.NewExpensesRepository(tx)

	if _, expenseError := expenseRepo.Create(request); expenseError != nil {
		tx.Rollback()
		customErrors.RespondWithError(w, http.StatusInternalServerError, customErrors.InternalServerError, expenseError.Error(), nil)
		return
	}

	if _, error := userRepo.Save(user); error != nil {
		tx.Rollback()
		customErrors.RespondWithError(w, http.StatusInternalServerError, customErrors.InternalServerError, error.Error(), nil)
		return
	}
	if _, error := accountRepo.Save(account); error != nil {
		tx.Rollback()
		customErrors.RespondWithError(w, http.StatusInternalServerError, customErrors.InternalServerError, error.Error(), nil)
		return
	}

	if err = tx.Commit().Error; err != nil {
		tx.Rollback()
		customErrors.RespondWithError(w, http.StatusInternalServerError, customErrors.InternalServerError, err.Error(), nil)
		return
	}

	jsonResponse, _ := json.Marshal(map[string]interface{}{"success": true})
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)

}

func (h *expenseRoutesHandler) newExpenseCategory(w http.ResponseWriter, r *http.Request) {
	var request CreateExpenseCategoryRequest
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
	expenseCategoryRepo := repositories.NewExpenseCategoryRepository(h.db)
	newExpenseCategory := repositories.CreateExpenseCategoryDto{
		Name: request.Name,
		User: userId,
		Icon: request.Icon,
	}
	category, err := expenseCategoryRepo.Create(newExpenseCategory)

	if err != nil {
		customErrors.RespondWithError(w, http.StatusInternalServerError, customErrors.InternalServerError, err.Error(), nil)
		return
	}

	newCategory := repositories.ExpenseCategoryResponse{CreatedAt:category.CreatedAt,UpdatedAt: category.UpdatedAt,Name: category.Name,Icon: category.Icon,ID: int(category.ID) }

	jsonResponse, _ := json.Marshal(newCategory)
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)

}

func creditAccount(accountId string, userId string, amount float64, tx *gorm.DB) (*repositories.Accounts, error) {
	accountRepo := repositories.NewAccountsRepository(tx)
	account, err := accountRepo.FindByIDAndUserID(accountId, userId)
	if err != nil {
		return nil, errors.New("failed to fetch account")
	}
	if account.Name == "" {
		return nil, errors.New("account not found")
	}

	if account.CurrentBalance < amount {
		return nil, errors.New("insufficient funds")
	}
	account.CurrentBalance -= amount
	return account, nil
}
