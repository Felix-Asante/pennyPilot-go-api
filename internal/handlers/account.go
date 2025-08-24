package handlers

import (
	"errors"
	"net/http"

	"github.com/Felix-Asante/pennyPilot-go-api/internal/dto"
	"github.com/Felix-Asante/pennyPilot-go-api/internal/models"
	"github.com/Felix-Asante/pennyPilot-go-api/internal/utils"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (h *Handler) createAccount(w http.ResponseWriter, r *http.Request) {
	var createAccountData dto.CreateAccountDto

	if err := utils.ReadAndValidateJSON(w, r, &createAccountData); err != nil {
		h.badRequestResponse(w, r, err)
		return
	}

	userID, err := getUserIdFromContext(r.Context())
	if err != nil {
		h.unauthorizedErrorResponse(w, r, err)
		return
	}

	oldAccount, err := h.Models.Account.GetByNameAndUserID(r.Context(), createAccountData.Name, userID, nil)

	if err != nil && err != gorm.ErrRecordNotFound {
		h.internalServerError(w, r, err)
		return
	}

	if oldAccount != nil {
		h.badRequestResponse(w, r, errors.New("account with this name already exists"))
		return
	}

	account := &models.Account{
		ID:       uuid.New(),
		Name:     createAccountData.Name,
		Currency: createAccountData.Currency,
		UserID:   userID,
	}

	if err := h.Models.Account.Create(r.Context(), account, nil); err != nil {
		h.internalServerError(w, r, err)
		return
	}

	// trigger automatic envelope allocation

	utils.WriteJSON(w, http.StatusCreated, models.SerializeAccount(account))
}

func (h *Handler) getAccounts(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserIdFromContext(r.Context())
	if err != nil {
		h.unauthorizedErrorResponse(w, r, err)
		return
	}

	accounts, err := h.Models.Account.GetAllByUserID(r.Context(), userID, nil)
	if err != nil {
		h.internalServerError(w, r, err)
		return
	}

	serializedAccounts := make([]*models.AccountSerializer, len(accounts))
	for i, account := range accounts {
		serializedAccounts[i] = models.SerializeAccount(account)
	}

	utils.WriteJSON(w, http.StatusOK, serializedAccounts)

}

func (h *Handler) getAccount(w http.ResponseWriter, r *http.Request) {
	accountID := chi.URLParam(r, "id")
	userID, err := getUserIdFromContext(r.Context())
	if err != nil {
		h.unauthorizedErrorResponse(w, r, err)
		return
	}

	account, err := h.Models.Account.GetByIDAndUserID(r.Context(), accountID, userID, nil)
	if err != nil && err != gorm.ErrRecordNotFound {
		h.internalServerError(w, r, err)
		return
	}

	if account == nil || err == gorm.ErrRecordNotFound {
		h.notFoundResponse(w, r, errors.New("account not found"))
		return
	}

	utils.WriteJSON(w, http.StatusOK, models.SerializeAccount(account))

}

func (h *Handler) updateAccount(w http.ResponseWriter, r *http.Request) {
	accountID := chi.URLParam(r, "id")
	userID, err := getUserIdFromContext(r.Context())
	if err != nil {
		h.unauthorizedErrorResponse(w, r, err)
		return
	}

	account, err := h.Models.Account.GetByIDAndUserID(r.Context(), accountID, userID, nil)
	if err != nil && err != gorm.ErrRecordNotFound {
		h.internalServerError(w, r, err)
		return
	}

	if account == nil || err == gorm.ErrRecordNotFound {
		h.notFoundResponse(w, r, errors.New("account not found"))
		return
	}

	var updateAccountDto dto.UpdateAccountDto

	if err := utils.ReadAndValidateJSON(w, r, &updateAccountDto); err != nil {
		h.badRequestResponse(w, r, err)
		return
	}

	if updateAccountDto.Name != nil && len(*updateAccountDto.Name) >= 3 {
		account.Name = *updateAccountDto.Name
	}
	if updateAccountDto.Currency != nil && len(*updateAccountDto.Currency) >= 3 {
		account.Currency = *updateAccountDto.Currency
	}
	if updateAccountDto.IsActive != nil {
		account.IsActive = *updateAccountDto.IsActive
	}

	if err := h.Models.Account.Save(r.Context(), account, nil); err != nil {
		h.internalServerError(w, r, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, models.SerializeAccount(account))

}

func (h *Handler) deleteAccount(w http.ResponseWriter, r *http.Request) {
	accountID := chi.URLParam(r, "id")
	userID, err := getUserIdFromContext(r.Context())
	if err != nil {
		h.unauthorizedErrorResponse(w, r, err)
		return
	}

	account, err := h.Models.Account.GetByIDAndUserID(r.Context(), accountID, userID, nil)
	if err != nil && err != gorm.ErrRecordNotFound {
		h.internalServerError(w, r, err)
		return
	}

	if account == nil || err == gorm.ErrRecordNotFound {
		h.notFoundResponse(w, r, errors.New("account not found"))
		return
	}

	if err := h.Models.Account.Delete(r.Context(), account, nil); err != nil {
		h.internalServerError(w, r, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]bool{"success": true})
}
