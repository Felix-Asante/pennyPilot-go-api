package handlers

import (
	"errors"
	"net/http"

	"github.com/Felix-Asante/pennyPilot-go-api/internal/dto"
	"github.com/Felix-Asante/pennyPilot-go-api/internal/models"
	"github.com/Felix-Asante/pennyPilot-go-api/internal/utils"
	"github.com/go-chi/jwtauth/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func (h *Handler) createUser(w http.ResponseWriter, r *http.Request) {
	var createUserDto dto.CreateUserDto

	if err := utils.ReadAndValidateJSON(w, r, &createUserDto); err != nil {
		h.badRequestResponse(w, r, err)
		return
	}

	user, err := h.Models.Users.GetUserByEmail(createUserDto.Email)

	if err != nil && err != gorm.ErrRecordNotFound {
		h.internalServerError(w, r, err)
		return
	}
	if user != nil {
		h.conflictResponse(w, r, errors.New("user with this email already exists"))
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(createUserDto.Password), bcrypt.DefaultCost)

	if err != nil {
		h.internalServerError(w, r, err)
		return
	}

	createUserDto.Password = string(hashedPassword)

	user, err = h.Models.Users.Create(&createUserDto)
	if err != nil {
		h.internalServerError(w, r, err)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, models.SerializeUser(user))
}

func (h *Handler) getCurrentUser(r *http.Request) (*models.User, error) {
	_, claims, err := jwtauth.FromContext(r.Context())
	if err != nil {
		h.Logger.Error("Missing or invalid token", "method", r.Method, "path", r.URL.Path, "error", err.Error())
		return nil, err
	}

	userId := claims["email"].(string)
	user, err := h.Models.Users.GetUserByEmail(userId)
	if err != nil && err != gorm.ErrRecordNotFound {
		h.Logger.Error("internal error", "method", r.Method, "path", r.URL.Path, "error", err.Error())
		return nil, err
	}

	if user == nil || err == gorm.ErrRecordNotFound {
		h.Logger.Error("user not found", "method", r.Method, "path", r.URL.Path)
		return nil, errors.New("user not found")
	}
	return user, nil
}
