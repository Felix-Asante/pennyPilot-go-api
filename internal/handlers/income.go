package handlers

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"sync"

	"github.com/Felix-Asante/pennyPilot-go-api/internal/dto"
	"github.com/Felix-Asante/pennyPilot-go-api/internal/models"
	"github.com/Felix-Asante/pennyPilot-go-api/internal/utils"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (h *Handler) createIncome(w http.ResponseWriter, r *http.Request) {
	var incomeDto dto.CreateIncomeDto

	if err := utils.ReadAndValidateJSON(w, r, &incomeDto); err != nil {
		h.badRequestResponse(w, r, err)
		return
	}

	userId, err := getUserIdFromContext(r.Context())
	if err != nil {
		h.unauthorizedErrorResponse(w, r, err)
		return
	}

	income := &models.Income{
		ID:           uuid.New(),
		UserID:       userId,
		Amount:       incomeDto.Amount,
		Category:     incomeDto.Category,
		DateRecieved: incomeDto.DateRecieved,
		Type:         incomeDto.Type,
		Frequency:    incomeDto.Frequency,
	}

	h.DB.Transaction(func(tx *gorm.DB) error {
		userIncomeBalance, balanceError := h.Models.Income.GetIncomeBalanceByUserId(r.Context(), userId, tx)

		if balanceError != nil && balanceError != gorm.ErrRecordNotFound {
			h.internalServerError(w, r, balanceError)
			return balanceError
		}

		if userIncomeBalance == nil || balanceError == gorm.ErrRecordNotFound {
			userIncomeBalance = &models.IncomeBalance{
				UserID:      userId,
				TotalIncome: incomeDto.Amount,
				Allocated:   0,
				Unallocated: incomeDto.Amount,
			}
		} else {
			userIncomeBalance.TotalIncome += incomeDto.Amount
			userIncomeBalance.Unallocated += incomeDto.Amount
		}

		savedIncome, err := h.Models.Income.Create(income, tx)

		if err != nil {
			h.internalServerError(w, r, err)
			return err
		}

		if err := h.Models.Income.SaveIncomeBalance(r.Context(), userIncomeBalance, tx); err != nil {
			return err
		}

		utils.WriteJSON(w, http.StatusCreated, models.SerializeIncome(savedIncome))
		return nil
	})
}

func (h *Handler) getUserIncome(w http.ResponseWriter, r *http.Request) {
	_, claims, err := jwtauth.FromContext(r.Context())
	if err != nil {
		h.unauthorizedErrorResponse(w, r, err)
		return
	}

	userID := claims["user_id"].(string)

	if userID == "" {
		h.unauthorizedErrorResponse(w, r, errors.New("unauthorized"))
		return
	}

	income, err := h.Models.Income.GetAllByUserID(userID, nil)
	if err != nil {
		h.internalServerError(w, r, err)
		return
	}

	serializedIncome := make([]models.IncomeSerializer, len(income))
	for i, income := range income {
		serializedIncome[i] = *models.SerializeIncome(income)
	}

	utils.WriteJSON(w, http.StatusOK, serializedIncome)
}

func (h *Handler) updateIncome(w http.ResponseWriter, r *http.Request) {
	var incomeDto dto.UpdateIncomeDto

	if err := utils.ReadAndValidateJSON(w, r, &incomeDto); err != nil {
		h.badRequestResponse(w, r, err)
		return
	}

	incomeID := chi.URLParam(r, "id")

	income, err := h.Models.Income.GetByID(incomeID, nil)
	if err != nil && err != gorm.ErrRecordNotFound {
		h.internalServerError(w, r, err)
		return
	}

	if income == nil || err == gorm.ErrRecordNotFound {
		h.notFoundResponse(w, r, errors.New("income not found"))
		return
	}

	_, claims, err := jwtauth.FromContext(r.Context())
	if err != nil {
		h.unauthorizedErrorResponse(w, r, err)
		return
	}

	if claims["user_id"].(string) != income.UserID {
		h.forbiddenResponse(w, r)
		return
	}

	if incomeDto.Amount != nil {
		income.Amount = *incomeDto.Amount
	}
	if incomeDto.Category != nil {
		income.Category = incomeDto.Category
	}
	if incomeDto.DateRecieved != nil {
		income.DateRecieved = *incomeDto.DateRecieved
	}
	if incomeDto.Type != nil {
		income.Type = *incomeDto.Type
	}
	if incomeDto.Frequency != nil {
		income.Frequency = *incomeDto.Frequency
	}

	if err := h.Models.Income.Save(income, nil); err != nil {
		h.internalServerError(w, r, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, models.SerializeIncome(income))
}

func (h *Handler) transferIncomeToAccount(w http.ResponseWriter, r *http.Request) {
	var transferDto dto.TransferIncome
	userID, err := getUserIdFromContext(r.Context())
	if err != nil {
		h.unauthorizedErrorResponse(w, r, err)
		return
	}

	if err := utils.ReadAndValidateJSON(w, r, &transferDto); err != nil {
		h.badRequestResponse(w, r, err)
		return
	}

	if transferDto.Amount <= 0 {
		h.badRequestResponse(w, r, errors.New("amount must be greater than 0"))
		return
	}

	if len(transferDto.Accounts) == 0 {
		h.badRequestResponse(w, r, errors.New("accounts is required"))
		return
	}

	totalTransferAmount := transferDto.Amount * float64(len(transferDto.Accounts))
	incomeBalance, err := h.Models.Income.GetIncomeBalanceByUserId(r.Context(), userID, nil)

	if err != nil {
		h.internalServerError(w, r, err)
		return
	}

	if incomeBalance == nil {
		h.notFoundResponse(w, r, errors.New("Insufficient balance"))
		return
	}

	if incomeBalance.Unallocated < totalTransferAmount {
		h.badRequestResponse(w, r, errors.New("Insufficient balance"))
		return
	}

	wg := sync.WaitGroup{}
	var mu sync.Mutex
	var errs []error

	for _, accountID := range transferDto.Accounts {
		wg.Add(1)
		go func(id string) {
			defer wg.Done()
			if err := h.moveIncomeToAccountWorker(r.Context(), userID, id, transferDto.Amount); err != nil {
				mu.Lock()
				errs = append(errs, err)
				mu.Unlock()
			}
		}(accountID)
	}

	wg.Wait()

	if len(errs) > 0 {
		h.internalServerError(w, r, fmt.Errorf("errors occurred: %v", errs))
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
	})
}

func (h *Handler) moveIncomeToAccountWorker(ctx context.Context, userId, accountId string, amount float64) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:

		return h.DB.Transaction(func(tx *gorm.DB) error {
			incomeBalance, err := h.Models.Income.GetIncomeBalanceByUserId(ctx, userId, tx)
			if err != nil {
				return err
			}

			if incomeBalance == nil {
				return errors.New("Insufficient balance")
			}

			account, err := h.Models.Account.GetByIDAndUserID(ctx, accountId, userId, tx)
			if err != nil {
				return err
			}

			if account == nil {
				return errors.New(fmt.Sprintf("account(%s) not found", accountId))
			}

			account.Balance += amount

			if err := h.Models.Account.Save(ctx, account, tx); err != nil {
				return err
			}

			incomeBalance.Unallocated -= amount
			incomeBalance.Allocated += amount

			if err := h.Models.Income.SaveIncomeBalance(ctx, incomeBalance, tx); err != nil {
				return err
			}

			return nil
		})
	}
}
