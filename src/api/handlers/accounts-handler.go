package handlers

import (
	"net/http"

	"gorm.io/gorm"
)

type accountsRoutesHandler struct {
	db *gorm.DB
}

func (h *accountsRoutesHandler) new(w http.ResponseWriter, r *http.Request) {

}

func (h *accountsRoutesHandler) get(w http.ResponseWriter, r *http.Request) {

}

func (h *accountsRoutesHandler) update(w http.ResponseWriter, r *http.Request) {

}

func (h *accountsRoutesHandler) delete(w http.ResponseWriter, r *http.Request) {

}

func (h *accountsRoutesHandler) transfer(w http.ResponseWriter, r *http.Request) {

}
