package handlers

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/Felix-Asante/pennyPilot-go-api/internal/dto"
	customErrors "github.com/Felix-Asante/pennyPilot-go-api/internal/errors"
	"github.com/Felix-Asante/pennyPilot-go-api/internal/models"
	"github.com/Felix-Asante/pennyPilot-go-api/internal/notifications"
	"github.com/Felix-Asante/pennyPilot-go-api/internal/utils"
	"github.com/go-chi/httprate"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func (h *Handler) login(w http.ResponseWriter, r *http.Request) {
	var loginDto dto.LoginDto

	if err := utils.ReadAndValidateJSON(w, r, &loginDto); err != nil {
		h.badRequestResponse(w, r, err)
		return
	}

	user, err := h.Models.Users.GetUserByEmail(loginDto.Email)

	if err != nil && err != gorm.ErrRecordNotFound {
		h.internalServerError(w, r, err)
		return
	}

	if user == nil {
		h.notFoundResponse(w, r, errors.New("no user exists with this email"))
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(loginDto.Password)); err != nil {
		h.unauthorizedErrorResponse(w, r, errors.New("invalid password"))
		return
	}

	claims := map[string]interface{}{
		"user_id": user.ID,
		"email":   user.Email,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
		"iat":     time.Now().Unix(),
	}

	_, tokenString, err := h.JWTAuth.Encode(claims)
	if err != nil {
		h.internalServerError(w, r, err)
		return
	}

	response := map[string]interface{}{
		"user":         models.SerializeUser(user),
		"access_token": tokenString,
	}

	ctx := context.WithValue(r.Context(), "claims", claims)

	defer ctx.Done()

	utils.WriteJSON(w, http.StatusOK, response)
}

func (h *Handler) getMe(w http.ResponseWriter, r *http.Request) {

	user, err := h.getCurrentUser(r)
	if err != nil {
		h.unauthorizedErrorResponse(w, r, err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, models.SerializeUser(user))
}

func (h *Handler) forgotPassword(w http.ResponseWriter, r *http.Request) {
	// TODO: FIX RATE LIMITER
	forgotPasswordRateLimit := httprate.NewRateLimiter(1, 30*time.Minute)
	var forgotPasswordDto dto.ForgotPasswordDto

	if err := utils.ReadAndValidateJSON(w, r, &forgotPasswordDto); err != nil {
		h.badRequestResponse(w, r, err)
		return
	}

	if forgotPasswordRateLimit.RespondOnLimit(w, r, forgotPasswordDto.Email) {
		h.rateLimitExceededResponse(w, r, "30 minutes")
		return
	}

	user, err := h.Models.Users.GetUserByEmail(forgotPasswordDto.Email)

	if err != nil && err != gorm.ErrRecordNotFound {
		h.internalServerError(w, r, err)
		return
	}

	if user == nil {
		utils.WriteJSON(w, http.StatusOK, map[string]string{"message": customErrors.RESET_PASSWORD_LINK_SENT})
		return
	}

	// generate token
	token, err := utils.GenerateRandomTokens(32)
	if err != nil {
		h.internalServerError(w, r, err)
		return
	}
	// hash token with bcrypt
	hashedToken, err := utils.HashString(token)
	if err != nil {
		h.internalServerError(w, r, err)
		return
	}

	// store token in db
	tokenExpiresAt := time.Now().Add(30 * time.Minute)

	code := &models.Code{
		Code:      hashedToken,
		UserID:    user.ID,
		Type:      utils.CodeTypeForgotPassword,
		ExpiresAt: &tokenExpiresAt,
		Used:      false,
	}
	if _, err = h.Models.Code.Create(code); err != nil {
		h.internalServerError(w, r, err)
		return
	}

	// send email with reset password link
	resetPaswordLink := fmt.Sprintf("%s/reset-password?token=%s&email=%s", utils.GetFrontendUrl(), token, user.Email)
	_, notificationError := h.Notifications.Mailer.Send(notifications.ForgotPasswordMessageTemplate, []string{user.Email}, "Reset Your Password - Penny Pilot", map[string]interface{}{
		"username": user.FullName,
		"link":     resetPaswordLink,
	})
	if notificationError != nil {
		h.internalServerError(w, r, notificationError)
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]string{"message": customErrors.RESET_PASSWORD_LINK_SENT})
}
