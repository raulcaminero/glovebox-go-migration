package http

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/raul/glovebox-go-migration/go-service/internal/domain"
	"github.com/raul/glovebox-go-migration/go-service/internal/service"
)

type PolicyHandler struct {
	svc *service.PolicyService
}

func NewPolicyHandler(svc *service.PolicyService) *PolicyHandler {
	return &PolicyHandler{svc: svc}
}

func (h *PolicyHandler) Routes(r chi.Router) {
	r.Post("/policies", h.create)
	r.Get("/policies/{id}", h.get)
	r.Get("/policyholders/{id}/policies", h.listByPolicyholder)
	r.Get("/policies/expiring", h.expiring)
	r.Patch("/policies/{id}/status", h.updateStatus)
	r.Delete("/policies/{id}", h.delete)
}

type createPolicyRequest struct {
	PolicyholderID string `json:"policyholder_id"`
	Carrier        string `json:"carrier"`
	PolicyNumber   string `json:"policy_number"`
	LineOfBusiness string `json:"line_of_business"`
	EffectiveDate  string `json:"effective_date"` // YYYY-MM-DD
	ExpirationDate string `json:"expiration_date"`
	PremiumCents   int64  `json:"premium_cents"`
}

func (h *PolicyHandler) create(w http.ResponseWriter, r *http.Request) {
	var req createPolicyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}
	phID, err := uuid.Parse(req.PolicyholderID)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid policyholder_id")
		return
	}
	eff, err := time.Parse("2006-01-02", req.EffectiveDate)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid effective_date, want YYYY-MM-DD")
		return
	}
	exp, err := time.Parse("2006-01-02", req.ExpirationDate)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid expiration_date, want YYYY-MM-DD")
		return
	}

	policy, err := h.svc.Create(r.Context(), domain.CreatePolicyInput{
		PolicyholderID: phID,
		Carrier:        req.Carrier,
		PolicyNumber:   req.PolicyNumber,
		LineOfBusiness: req.LineOfBusiness,
		EffectiveDate:  eff,
		ExpirationDate: exp,
		PremiumCents:   req.PremiumCents,
	})
	if err != nil {
		handleServiceErr(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, policy)
}

func (h *PolicyHandler) get(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}
	policy, err := h.svc.Get(r.Context(), id)
	if err != nil {
		handleServiceErr(w, err)
		return
	}
	writeJSON(w, http.StatusOK, policy)
}

func (h *PolicyHandler) listByPolicyholder(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}
	policies, err := h.svc.ListByPolicyholder(r.Context(), id)
	if err != nil {
		handleServiceErr(w, err)
		return
	}
	writeJSON(w, http.StatusOK, policies)
}

// expiring answers "which policies renew in the next N days" — the same
// question surfaced as an MCP tool for AI clients (docs/AI_WORKFLOW.md).
func (h *PolicyHandler) expiring(w http.ResponseWriter, r *http.Request) {
	days, err := strconv.Atoi(r.URL.Query().Get("within_days"))
	if err != nil || days <= 0 {
		days = 30
	}
	policies, err := h.svc.ExpiringWithin(r.Context(), int32(days))
	if err != nil {
		handleServiceErr(w, err)
		return
	}
	writeJSON(w, http.StatusOK, policies)
}

type updateStatusRequest struct {
	Status string `json:"status"`
}

func (h *PolicyHandler) updateStatus(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}
	var req updateStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}
	policy, err := h.svc.UpdateStatus(r.Context(), id, domain.PolicyStatus(req.Status))
	if err != nil {
		handleServiceErr(w, err)
		return
	}
	writeJSON(w, http.StatusOK, policy)
}

func (h *PolicyHandler) delete(w http.ResponseWriter, r *http.Request) {
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
