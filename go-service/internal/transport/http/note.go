package http

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/raul/glovebox-go-migration/go-service/internal/domain"
	"github.com/raul/glovebox-go-migration/go-service/internal/service"
)

type NoteHandler struct {
	svc *service.NoteService
}

func NewNoteHandler(svc *service.NoteService) *NoteHandler {
	return &NoteHandler{svc: svc}
}

func (h *NoteHandler) Routes(r chi.Router) {
	r.Post("/policyholders/{id}/notes", h.Create)
	r.Get("/policyholders/{id}/notes", h.ListByPolicyholder)
	r.Delete("/notes/{id}", h.Delete)
}

type createNoteRequest struct {
	Author string `json:"author"`
	Body   string `json:"body"`
}

func (h *NoteHandler) Create(w http.ResponseWriter, r *http.Request) {
	phID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid policyholder id"})
		return
	}

	var req createNoteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid JSON body"})
		return
	}

	created, err := h.svc.Create(r.Context(), domain.CreateNoteInput{
		PolicyholderID: phID,
		Author:         req.Author,
		Body:           req.Body,
	})
	if err != nil {
		if errors.Is(err, service.ErrInvalidInput) {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
			return
		}
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal server error"})
		return
	}

	writeJSON(w, http.StatusCreated, created)
}

func (h *NoteHandler) ListByPolicyholder(w http.ResponseWriter, r *http.Request) {
	phID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid policyholder id"})
		return
	}

	notes, err := h.svc.ListByPolicyholder(r.Context(), phID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal server error"})
		return
	}

	if notes == nil {
		notes = []domain.Note{}
	}

	writeJSON(w, http.StatusOK, notes)
}

func (h *NoteHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid note id"})
		return
	}

	if err := h.svc.Delete(r.Context(), id); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal server error"})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
