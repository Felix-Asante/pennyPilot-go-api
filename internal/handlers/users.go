package handlers

import (
	"errors"
	"net/http"

	"github.com/Felix-Asante/pennyPilot-go-api/internal/dto"
)

func (h *Handler) createUser(w http.ResponseWriter, r *http.Request) {
	var createUserDto dto.CreateUserDto

	if err := readJSON(w, r, &createUserDto); err != nil {
		h.badRequestResponse(w, r, err)
		return
	}

	if err := Validate.Struct(createUserDto); err != nil {
		h.badRequestResponse(w, r, err)
		return
	}

	user, err := h.Models.Users.GetUserByEmail(createUserDto.Email)

	if err != nil {
		h.internalServerError(w, r, err)
		return
	}
	if user != nil {
		h.conflictResponse(w, r, errors.New("user with this email already exists"))
		return
	}

	user, err = h.Models.Users.Create(&createUserDto)
	if err != nil {
		h.internalServerError(w, r, err)
		return
	}

	writeJSON(w, http.StatusCreated, user)
}

func (h *Handler) getMe(w http.ResponseWriter, r *http.Request) {

	w.Write([]byte("get me"))
}
