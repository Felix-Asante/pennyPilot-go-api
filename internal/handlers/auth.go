package handlers

import (
	"errors"
	"net/http"
	"time"

	"github.com/Felix-Asante/pennyPilot-go-api/internal/dto"
	"golang.org/x/crypto/bcrypt"
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

	if err != nil {
		h.internalServerError(w, r, err)
		return
	}
	if user == nil {
		h.notFoundResponse(w, r, errors.New("user not found"))
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(loginDto.Password)); err != nil {
		h.unauthorizedErrorResponse(w, r, err)
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
		"user": user,
		"accessToken": tokenString,
	}
	writeJSON(w, http.StatusOK, response)
}

func (h *Handler) getMe(w http.ResponseWriter, r *http.Request) {

	user, err := h.getCurrentUser(r)
	if err != nil {
		h.internalServerError(w, r, err)
		return
	}
	writeJSON(w, http.StatusOK, user)
}