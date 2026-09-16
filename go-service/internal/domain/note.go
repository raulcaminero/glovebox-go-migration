package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Note struct {
	ID             uuid.UUID `json:"id"`
	PolicyholderID uuid.UUID `json:"policyholder_id"`
	Author         string    `json:"author"`
	Body           string    `json:"body"`
	CreatedAt      time.Time `json:"created_at"`
}

type CreateNoteInput struct {
	PolicyholderID uuid.UUID
	Author         string
	Body           string
}

// NoteRepo is the port for note persistence.
type NoteRepo interface {
	Create(ctx context.Context, in CreateNoteInput) (Note, error)
	Get(ctx context.Context, id uuid.UUID) (Note, error)
	ListByPolicyholder(ctx context.Context, policyholderID uuid.UUID) ([]Note, error)
	Delete(ctx context.Context, id uuid.UUID) error
}
