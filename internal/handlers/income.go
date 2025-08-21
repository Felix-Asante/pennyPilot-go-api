package handlers

import (
	"errors"
	"net/http"

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

	user, err := h.getCurrentUser(r)
	if err != nil {
		h.unauthorizedErrorResponse(w, r, err)
		return
	}

	income := &models.Income{
		ID:           uuid.New(),
		UserID:       user.ID,
		Amount:       incomeDto.Amount,
		Category:     incomeDto.Category,
		DateRecieved: incomeDto.DateRecieved,
		Type:         incomeDto.Type,
		Frequency:    incomeDto.Frequency,
	}

	savedIncome, err := h.Models.Income.Create(income, nil)

	if err != nil {
		h.internalServerError(w, r, err)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, models.SerializeIncome(savedIncome))
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
