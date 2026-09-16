// Package http contains thin HTTP handlers: decode request, call service,
// encode response. No business logic lives here — that's the point of the
// service layer boundary.
package http

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/raul/glovebox-go-migration/go-service/internal/domain"
	"github.com/raul/glovebox-go-migration/go-service/internal/service"
)

type PolicyholderHandler struct {
	svc *service.PolicyholderService
}

func NewPolicyholderHandler(svc *service.PolicyholderService) *PolicyholderHandler {
	return &PolicyholderHandler{svc: svc}
}

func (h *PolicyholderHandler) Routes(r chi.Router) {
	r.Post("/policyholders", h.create)
	r.Get("/policyholders", h.list)
	r.Get("/policyholders/{id}", h.get)
	r.Put("/policyholders/{id}", h.update)
	r.Delete("/policyholders/{id}", h.delete)
}

type createPolicyholderRequest struct {
	FullName string `json:"full_name"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
}

func (h *PolicyholderHandler) create(w http.ResponseWriter, r *http.Request) {
	var req createPolicyholderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}

	ph, err := h.svc.Create(r.Context(), domain.CreatePolicyholderInput{
		FullName: req.FullName,
		Email:    req.Email,
		Phone:    req.Phone,
	})
	if err != nil {
		handleServiceErr(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, ph)
}

func (h *PolicyholderHandler) get(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}
	ph, err := h.svc.Get(r.Context(), id)
	if err != nil {
		handleServiceErr(w, err)
		return
	}
	writeJSON(w, http.StatusOK, ph)
}

func (h *PolicyholderHandler) list(w http.ResponseWriter, r *http.Request) {
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))

	items, err := h.svc.List(r.Context(), domain.ListParams{
		Query:  r.URL.Query().Get("q"),
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		handleServiceErr(w, err)
		return
	}
	writeJSON(w, http.StatusOK, items)
}

type updatePolicyholderRequest struct {
	FullName string `json:"full_name"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
}

func (h *PolicyholderHandler) update(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}
	var req updatePolicyholderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}
	ph, err := h.svc.Update(r.Context(), domain.UpdatePolicyholderInput{
		ID: id, FullName: req.FullName, Email: req.Email, Phone: req.Phone,
	})
	if err != nil {
		handleServiceErr(w, err)
		return
	}
	writeJSON(w, http.StatusOK, ph)
}

func (h *PolicyholderHandler) delete(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}
	if err := h.svc.Delete(r.Context(), id); err != nil {
		handleServiceErr(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func handleServiceErr(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, service.ErrInvalidInput):
		writeError(w, http.StatusBadRequest, err.Error())
	case errors.Is(err, service.ErrNotFound):
		writeError(w, http.StatusNotFound, err.Error())
	default:
		writeError(w, http.StatusInternalServerError, "internal error")
	}
}

func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}
