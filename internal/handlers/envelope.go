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

			allocationStrategy := utils.AllocationStrategy(*createEnvelopeDto.AllocationStrategy)
			allocationValue := *createEnvelopeDto.AllocationValue

			if allocationStrategy == utils.AllocationStrategyPercentage && allocationValue > 100 {
				h.badRequestResponse(w, r, errors.New("allocation value must be less than or equal to 100"))
				return errors.New("allocation value must be less than or equal to 100")
			}

			allocationRule.ID = uuid.New()
			allocationRule.Strategy = allocationStrategy
			allocationRule.Value = allocationValue
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
	envelopeID := chi.URLParam(r, "id")
	userID, err := getUserIdFromContext(r.Context())
	if err != nil {
		h.unauthorizedErrorResponse(w, r, err)
		return
	}

	if uuid.Validate(envelopeID) != nil || uuid.Validate(userID) != nil {
		h.badRequestResponse(w, r, errors.New("invalid envelope ID or user ID"))
		return
	}

	envelope, err := h.Models.Envelope.GetByIDAndUserID(r.Context(), envelopeID, userID, nil)
	if err != nil && err != gorm.ErrRecordNotFound {
		h.internalServerError(w, r, err)
		return
	}

	if envelope == nil || err == gorm.ErrRecordNotFound {
		h.notFoundResponse(w, r, errors.New("envelope not found"))
		return
	}

	utils.WriteJSON(w, http.StatusOK, models.SerializeEnvelope(envelope))
}

func (h *Handler) getUserOwnedEnvelopes(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserIdFromContext(r.Context())
	if err != nil {
		h.unauthorizedErrorResponse(w, r, err)
		return
	}

	envelopes, err := h.Models.Envelope.GetAllByUserID(r.Context(), userID, nil)
	if err != nil {
		h.internalServerError(w, r, err)
		return
	}

	// improve with goroutines and channels
	serializedEnvelopes := make([]*models.EnvelopeSerializer, len(envelopes))
	for i, envelope := range envelopes {
		serializedEnvelopes[i] = models.SerializeEnvelope(envelope)
	}

	utils.WriteJSON(w, http.StatusOK, serializedEnvelopes)
}

func (h *Handler) updateEnvelope(w http.ResponseWriter, r *http.Request) {
	envelopeID := chi.URLParam(r, "id")
	userID, err := getUserIdFromContext(r.Context())

	if err != nil {
		h.unauthorizedErrorResponse(w, r, err)
		return
	}

	if uuid.Validate(envelopeID) != nil || uuid.Validate(userID) != nil {
		h.badRequestResponse(w, r, errors.New("invalid envelope ID or user ID"))
		return
	}

	var updateEnvelopeDto dto.UpdateEnvelopeDto

	if err := utils.ReadAndValidateJSON(w, r, &updateEnvelopeDto); err != nil {
		h.badRequestResponse(w, r, err)
		return
	}

	h.DB.Transaction(func(tx *gorm.DB) error {
		envelope, err := h.Models.Envelope.GetByIDAndUserID(r.Context(), envelopeID, userID, nil)
		if err != nil && err != gorm.ErrRecordNotFound {
			h.internalServerError(w, r, err)
			return err
		}

		if envelope == nil || err == gorm.ErrRecordNotFound {
			h.notFoundResponse(w, r, errors.New("envelope not found"))
			return errors.New("envelope not found")
		}

		if updateEnvelopeDto.Name != nil && len(*updateEnvelopeDto.Name) > 0 {
			envelope.Name = *updateEnvelopeDto.Name
		}

		if updateEnvelopeDto.TargetAmount != nil {
			envelope.TargetAmount = *updateEnvelopeDto.TargetAmount
		}

		if updateEnvelopeDto.TargetedDate != nil {
			envelope.TargetedDate = updateEnvelopeDto.TargetedDate
		}

		if updateEnvelopeDto.IsActive != nil {
			envelope.IsActive = *updateEnvelopeDto.IsActive
		}

		oldAllocationRule := envelope.AllocationRule

		if updateEnvelopeDto.AutoAllocate != nil {
			isAutoAllocate := *updateEnvelopeDto.AutoAllocate

			envelope.AutoAllocate = isAutoAllocate

			if oldAllocationRule == nil && isAutoAllocate {
				oldAllocationRule = &models.AllocationRule{
					ID:       uuid.New(),
					Active:   true,
					TargetID: envelope.ID,
				}
			}

			if isAutoAllocate && (updateEnvelopeDto.AllocationStrategy == nil || updateEnvelopeDto.AllocationValue == nil) {
				h.badRequestResponse(w, r, errors.New("allocation strategy and value are required"))
				return errors.New("allocation strategy and value are required")
			}

			if updateEnvelopeDto.AllocationStrategy != nil && len(*updateEnvelopeDto.AllocationStrategy) > 0 {
				oldAllocationRule.Strategy = utils.AllocationStrategy(*updateEnvelopeDto.AllocationStrategy)
			}

			if updateEnvelopeDto.AllocationValue != nil {
				if oldAllocationRule.Strategy == utils.AllocationStrategyPercentage && *updateEnvelopeDto.AllocationValue > 100 {
					h.badRequestResponse(w, r, errors.New("allocation value must be less than or equal to 100"))
					return errors.New("allocation value must be less than or equal to 100")
				}
				oldAllocationRule.Value = *updateEnvelopeDto.AllocationValue
			}

			oldAllocationRule.Active = isAutoAllocate

			envelope.AllocationRule = oldAllocationRule
		}

		if err := h.Models.Envelope.Save(r.Context(), envelope, nil); err != nil {
			h.internalServerError(w, r, err)
			return err
		}

		if oldAllocationRule != nil {
			if err := h.Models.AllocationRule.Save(r.Context(), oldAllocationRule, nil); err != nil {
				h.internalServerError(w, r, err)
				return err
			}
		}

		utils.WriteJSON(w, http.StatusOK, models.SerializeEnvelope(envelope))
		return nil
	})

}

func (h *Handler) deleteEnvelope(w http.ResponseWriter, r *http.Request) {

	// decide what to with unallocated balance
}
