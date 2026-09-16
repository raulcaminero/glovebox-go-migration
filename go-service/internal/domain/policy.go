package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type PolicyStatus string

const (
	PolicyStatusActive    PolicyStatus = "active"
	PolicyStatusExpired   PolicyStatus = "expired"
	PolicyStatusCancelled PolicyStatus = "cancelled"
)

type Policy struct {
	ID             uuid.UUID    `json:"id"`
	PolicyholderID uuid.UUID    `json:"policyholder_id"`
	Carrier        string       `json:"carrier"`
	PolicyNumber   string       `json:"policy_number"`
	LineOfBusiness string       `json:"line_of_business"`
	EffectiveDate  time.Time    `json:"effective_date"`
	ExpirationDate time.Time    `json:"expiration_date"`
	PremiumCents   int64        `json:"premium_cents"`
	Status         PolicyStatus `json:"status"`
	CreatedAt      time.Time    `json:"created_at"`
	UpdatedAt      time.Time    `json:"updated_at"`
}

type CreatePolicyInput struct {
	PolicyholderID uuid.UUID
	Carrier        string
	PolicyNumber   string
	LineOfBusiness string
	EffectiveDate  time.Time
	ExpirationDate time.Time
	PremiumCents   int64
}

// PolicyRepo mirrors PolicyholderRepo's pattern: an interface the service
// layer depends on, implemented by the Postgres adapter in internal/repo.
type PolicyRepo interface {
	Create(ctx context.Context, in CreatePolicyInput) (Policy, error)
	Get(ctx context.Context, id uuid.UUID) (Policy, error)
	ListByPolicyholder(ctx context.Context, policyholderID uuid.UUID) ([]Policy, error)
	ListExpiringWithin(ctx context.Context, days int32) ([]Policy, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status PolicyStatus) (Policy, error)
	Delete(ctx context.Context, id uuid.UUID) error
}
