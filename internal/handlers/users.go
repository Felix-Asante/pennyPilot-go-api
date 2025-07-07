package handlers

import "net/http"

func (h *Handler) createUser(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("create user"))
}

func (h *Handler) getMe(w http.ResponseWriter, r *http.Request) {

	w.Write([]byte("get me"))
}
