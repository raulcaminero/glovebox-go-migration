// Package domain holds core business types. It has zero dependencies on
// transport (HTTP), persistence (pgx/sqlc), or any framework. This is what
// lets us swap the DB driver or the web framework without touching business
// rules — the classic "ports and adapters" boundary.
package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Policyholder struct {
	ID        uuid.UUID `json:"id"`
	FullName  string    `json:"full_name"`
	Email     string    `json:"email"`
	Phone     string    `json:"phone,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreatePolicyholderInput struct {
	FullName string
	Email    string
	Phone    string
}

type UpdatePolicyholderInput struct {
	ID       uuid.UUID
	FullName string
	Email    string
	Phone    string
}

type ListParams struct {
	Query  string // optional search term over name/email
	Limit  int32
	Offset int32
}

// PolicyholderRepo is the port the service layer depends on. Any storage
// backend (Postgres via pgx/sqlc today, something else tomorrow) implements
// this interface. The service layer never imports pgx or sqlc directly.
type PolicyholderRepo interface {
	Create(ctx context.Context, in CreatePolicyholderInput) (Policyholder, error)
	Get(ctx context.Context, id uuid.UUID) (Policyholder, error)
	List(ctx context.Context, p ListParams) ([]Policyholder, error)
	Update(ctx context.Context, in UpdatePolicyholderInput) (Policyholder, error)
	Delete(ctx context.Context, id uuid.UUID) error
}
