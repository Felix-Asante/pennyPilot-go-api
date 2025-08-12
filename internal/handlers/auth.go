package handlers

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/Felix-Asante/pennyPilot-go-api/internal/dto"
	"github.com/Felix-Asante/pennyPilot-go-api/internal/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func (h *Handler) login(w http.ResponseWriter, r *http.Request) {
	var loginDto dto.LoginDto

	if err := readJSON(w, r, &loginDto); err != nil {
		h.badRequestResponse(w, r, err)
		return
	}

	if err := Validate.Struct(loginDto); err != nil {
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
		"user": models.SerializeUser(user),
		"access_token": tokenString,
	}

	ctx := context.WithValue(r.Context(), "claims", claims)

	defer ctx.Done()
	
	writeJSON(w, http.StatusOK, response)
}

func (h *Handler) getMe(w http.ResponseWriter, r *http.Request) {

	user, err := h.getCurrentUser(r)
	if err != nil {
		h.unauthorizedErrorResponse(w, r, err)
		return
	}
	writeJSON(w, http.StatusOK, models.SerializeUser(user))
}