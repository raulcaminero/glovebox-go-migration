package service

import (
	"context"
	"errors"

	"github.com/google/uuid"

	"github.com/raul/glovebox-go-migration/go-service/internal/domain"
)

type PolicyService struct {
	repo domain.PolicyRepo
}

func NewPolicyService(repo domain.PolicyRepo) *PolicyService {
	return &PolicyService{repo: repo}
}

func (s *PolicyService) Create(ctx context.Context, in domain.CreatePolicyInput) (domain.Policy, error) {
	if in.Carrier == "" || in.PolicyNumber == "" {
		return domain.Policy{}, wrapInvalid("carrier and policy_number are required")
	}
	if !in.ExpirationDate.After(in.EffectiveDate) {
		return domain.Policy{}, wrapInvalid("expiration_date must be after effective_date")
	}
	if in.PremiumCents < 0 {
		return domain.Policy{}, wrapInvalid("premium_cents cannot be negative")
	}
	return s.repo.Create(ctx, in)
}

func (s *PolicyService) Get(ctx context.Context, id uuid.UUID) (domain.Policy, error) {
	return s.repo.Get(ctx, id)
}

func (s *PolicyService) ListByPolicyholder(ctx context.Context, policyholderID uuid.UUID) ([]domain.Policy, error) {
	return s.repo.ListByPolicyholder(ctx, policyholderID)
}

// ExpiringWithin powers the "which clients renew soon" use case — the same
// question a GloveBox agent would ask, and the same query the MCP tool
// (see docs/AI_WORKFLOW.md) exposes to an LLM client.
func (s *PolicyService) ExpiringWithin(ctx context.Context, days int32) ([]domain.Policy, error) {
	if days <= 0 {
		return nil, wrapInvalid("days must be positive")
	}
	if days > 365 {
		days = 365
	}
	return s.repo.ListExpiringWithin(ctx, days)
}

func (s *PolicyService) UpdateStatus(ctx context.Context, id uuid.UUID, status domain.PolicyStatus) (domain.Policy, error) {
	switch status {
	case domain.PolicyStatusActive, domain.PolicyStatusExpired, domain.PolicyStatusCancelled:
	default:
		return domain.Policy{}, wrapInvalid("invalid status")
	}
	return s.repo.UpdateStatus(ctx, id, status)
}

func (s *PolicyService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}

var ErrNotFound = errors.New("not found")
