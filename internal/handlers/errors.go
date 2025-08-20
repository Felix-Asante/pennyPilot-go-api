package handlers

import (
	"net/http"

	customErrors "github.com/Felix-Asante/pennyPilot-go-api/internal/errors"
	"github.com/Felix-Asante/pennyPilot-go-api/internal/utils"
)

func (h *Handler) internalServerError(w http.ResponseWriter, r *http.Request, err error) {
	h.Logger.Error("internal error", "method", r.Method, "path", r.URL.Path, "error", err.Error(), "ip", r.RemoteAddr)

	utils.WriteJSONError(w, http.StatusInternalServerError, map[string]string{"message": "the server encountered a problem"})
}

func (h *Handler) forbiddenResponse(w http.ResponseWriter, r *http.Request) {
	h.Logger.Warn("forbidden", "method", r.Method, "path", r.URL.Path)

	utils.WriteJSONError(w, http.StatusForbidden, map[string]string{"message": "forbidden", "ip": r.RemoteAddr})
}

func (h *Handler) badRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	h.Logger.Warn("bad request", "method", r.Method, "path", r.URL.Path, "error", err.Error(), "ip", r.RemoteAddr)

	if mapErr, ok := err.(*customErrors.MapError); ok {
		utils.WriteJSONError(w, http.StatusBadRequest, mapErr.Errors)
		return
	}

	utils.WriteJSONError(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
}

func (h *Handler) conflictResponse(w http.ResponseWriter, r *http.Request, err error) {
	h.Logger.Error("conflict response", "method", r.Method, "path", r.URL.Path, "error", err.Error(), "ip", r.RemoteAddr)

	utils.WriteJSONError(w, http.StatusConflict, map[string]string{"message": err.Error()})
}

func (h *Handler) notFoundResponse(w http.ResponseWriter, r *http.Request, err error) {
	h.Logger.Warn("not found error", "method", r.Method, "path", r.URL.Path, "error", err.Error(), "ip", r.RemoteAddr)

	utils.WriteJSONError(w, http.StatusNotFound, map[string]string{"message": err.Error()})
}

func (h *Handler) unauthorizedErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	h.Logger.Warn("unauthorized error", "method", r.Method, "path", r.URL.Path, "error", err.Error(), "ip", r.RemoteAddr)

	utils.WriteJSONError(w, http.StatusUnauthorized, map[string]string{"message": err.Error()})
}

func (h *Handler) unauthorizedBasicErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	h.Logger.Warn("unauthorized basic error", "method", r.Method, "path", r.URL.Path, "error", err.Error(), "ip", r.RemoteAddr)

	w.Header().Set("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)

	utils.WriteJSONError(w, http.StatusUnauthorized, map[string]string{"message": "unauthorized"})
}

func (h *Handler) rateLimitExceededResponse(w http.ResponseWriter, r *http.Request, retryAfter string) {
	h.Logger.Warn("rate limit exceeded", "method", r.Method, "path", r.URL.Path, "ip", r.RemoteAddr)

	w.Header().Set("Retry-After", retryAfter)

	utils.WriteJSONError(w, http.StatusTooManyRequests, map[string]string{"message": "rate limit exceeded, retry after: " + retryAfter})
}
