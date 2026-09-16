package repo

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/raul/glovebox-go-migration/go-service/internal/domain"
	"github.com/raul/glovebox-go-migration/go-service/internal/repo/sqlcgen"
)

type NotePG struct {
	q *sqlcgen.Queries
}

func NewNotePG(pool *pgxpool.Pool) *NotePG {
	return &NotePG{q: sqlcgen.New(pool)}
}

func (r *NotePG) Create(ctx context.Context, in domain.CreateNoteInput) (domain.Note, error) {
	row, err := r.q.CreateNote(ctx, sqlcgen.CreateNoteParams{
		PolicyholderID: pgUUID(in.PolicyholderID),
		Author:         in.Author,
		Body:           in.Body,
	})
	if err != nil {
		return domain.Note{}, fmt.Errorf("creating note in DB: %w", err)
	}
	return noteFromSQLC(row), nil
}

func (r *NotePG) Get(ctx context.Context, id uuid.UUID) (domain.Note, error) {
	row, err := r.q.GetNoteByID(ctx, pgUUID(id))
	if err != nil {
		return domain.Note{}, fmt.Errorf("fetching note %s: %w", id, err)
	}
	return noteFromSQLC(row), nil
}

func (r *NotePG) ListByPolicyholder(ctx context.Context, policyholderID uuid.UUID) ([]domain.Note, error) {
	rows, err := r.q.ListNotesByPolicyholder(ctx, pgUUID(policyholderID))
	if err != nil {
		return nil, fmt.Errorf("listing notes for policyholder %s: %w", policyholderID, err)
	}
	out := make([]domain.Note, len(rows))
	for i, row := range rows {
		out[i] = noteFromSQLC(row)
	}
	return out, nil
}

func (r *NotePG) Delete(ctx context.Context, id uuid.UUID) error {
	if err := r.q.DeleteNote(ctx, pgUUID(id)); err != nil {
		return fmt.Errorf("deleting note %s: %w", id, err)
	}
	return nil
}

func noteFromSQLC(n sqlcgen.Note) domain.Note {
	return domain.Note{
		ID:             uuidFromPG(n.ID),
		PolicyholderID: uuidFromPG(n.PolicyholderID),
		Author:         n.Author,
		Body:           n.Body,
		CreatedAt:      n.CreatedAt.Time,
	}
}
