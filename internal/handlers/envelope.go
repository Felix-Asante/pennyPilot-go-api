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

func (h *Handler) createEnvelope(w http.ResponseWriter, r *http.Request) {
	var createEnvelopeDto dto.CreateEnvelopeDto

	if err := utils.ReadAndValidateJSON(w, r, &createEnvelopeDto); err != nil {
		h.badRequestResponse(w, r, err)
		return
	}

	userID, err := getUserIdFromContext(r.Context())
	if err != nil {
		h.unauthorizedErrorResponse(w, r, err)
		return
	}

	h.DB.Transaction(func(tx *gorm.DB) error {
		account, err := h.Models.Account.GetByIDAndUserID(r.Context(), createEnvelopeDto.Account, userID, tx)
		if err != nil && err != gorm.ErrRecordNotFound {
			h.internalServerError(w, r, err)
			return err
		}

		if account == nil || err == gorm.ErrRecordNotFound {
			h.notFoundResponse(w, r, errors.New("account not found"))
			return errors.New("account not found")
		}

		oldEnvelope, err := h.Models.Envelope.GetByNameAndAccountID(r.Context(), createEnvelopeDto.Name, account.ID, tx)
		if err != nil && err != gorm.ErrRecordNotFound {
			h.internalServerError(w, r, err)
			return err
		}

		if oldEnvelope != nil {
			h.badRequestResponse(w, r, errors.New("envelope with this name already exists"))
			return errors.New("envelope with this name already exists")
		}

		envelope := &models.Envelope{
			ID:           uuid.New(),
			Name:         createEnvelopeDto.Name,
			AccountID:    account.ID,
			TargetAmount: *createEnvelopeDto.TargetAmount,
			AutoAllocate: *createEnvelopeDto.AutoAllocate,
			TargetedDate: createEnvelopeDto.TargetedDate,
		}

		allocationRule := &models.AllocationRule{}

		if createEnvelopeDto.AutoAllocate != nil && *createEnvelopeDto.AutoAllocate {
			if createEnvelopeDto.AllocationStrategy == nil || createEnvelopeDto.AllocationValue == nil {
				h.badRequestResponse(w, r, errors.New("allocation strategy and value are required"))
				return errors.New("allocation strategy and value are required")
			}
			allocationRule.Strategy = utils.AllocationStrategy(*createEnvelopeDto.AllocationStrategy)
			allocationRule.Value = *createEnvelopeDto.AllocationValue
			allocationRule.TargetID = envelope.ID
		}

		if err := h.Models.Envelope.Create(r.Context(), envelope, tx); err != nil {
			h.internalServerError(w, r, err)
			return err
		}

		if createEnvelopeDto.AutoAllocate != nil && *createEnvelopeDto.AutoAllocate {
			if err := h.Models.AllocationRule.Create(r.Context(), allocationRule, tx); err != nil {
				h.internalServerError(w, r, err)
				return err
			}
		}

		utils.WriteJSON(w, http.StatusCreated, models.SerializeEnvelope(envelope))
		return nil
	})

}

func (h *Handler) getEnvelopes(w http.ResponseWriter, r *http.Request) {
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

	envelopes, err := h.Models.Envelope.GetAllByAccountID(r.Context(), account.ID, nil)
	if err != nil {
		h.internalServerError(w, r, err)
		return
	}

	serializedEnvelopes := make([]*models.EnvelopeSerializer, len(envelopes))
	for i, envelope := range envelopes {
		serializedEnvelopes[i] = models.SerializeEnvelope(envelope)
	}

	utils.WriteJSON(w, http.StatusOK, serializedEnvelopes)

}

func (h *Handler) getEnvelope(w http.ResponseWriter, r *http.Request) {

}

func (h *Handler) updateEnvelope(w http.ResponseWriter, r *http.Request) {

}

func (h *Handler) deleteEnvelope(w http.ResponseWriter, r *http.Request) {

	// decide what to with unallocated balance
}
