package service_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/raul/glovebox-go-migration/go-service/internal/domain"
	"github.com/raul/glovebox-go-migration/go-service/internal/service"
)

type fakeNoteRepo struct {
	notes map[uuid.UUID]domain.Note
}

func newFakeNoteRepo() *fakeNoteRepo {
	return &fakeNoteRepo{notes: make(map[uuid.UUID]domain.Note)}
}

func (r *fakeNoteRepo) Create(ctx context.Context, in domain.CreateNoteInput) (domain.Note, error) {
	id := uuid.New()
	n := domain.Note{
		ID:             id,
		PolicyholderID: in.PolicyholderID,
		Author:         in.Author,
		Body:           in.Body,
		CreatedAt:      time.Now(),
	}
	r.notes[id] = n
	return n, nil
}

func (r *fakeNoteRepo) Get(ctx context.Context, id uuid.UUID) (domain.Note, error) {
	n, ok := r.notes[id]
	if !ok {
		return domain.Note{}, errors.New("not found")
	}
	return n, nil
}

func (r *fakeNoteRepo) ListByPolicyholder(ctx context.Context, policyholderID uuid.UUID) ([]domain.Note, error) {
	var out []domain.Note
	for _, n := range r.notes {
		if n.PolicyholderID == policyholderID {
			out = append(out, n)
		}
	}
	return out, nil
}

func (r *fakeNoteRepo) Delete(ctx context.Context, id uuid.UUID) error {
	delete(r.notes, id)
	return nil
}

func TestNoteService_Create(t *testing.T) {
	phID := uuid.New()

	tests := []struct {
		name    string
		input   domain.CreateNoteInput
		wantErr bool
	}{
		{
			name: "valid note",
			input: domain.CreateNoteInput{
				PolicyholderID: phID,
				Author:         "Agent Smith",
				Body:           "Client called to update policy address.",
			},
			wantErr: false,
		},
		{
			name: "missing policyholder ID",
			input: domain.CreateNoteInput{
				Author: "Agent Smith",
				Body:   "Some note",
			},
			wantErr: true,
		},
		{
			name: "empty body",
			input: domain.CreateNoteInput{
				PolicyholderID: phID,
				Author:         "Agent Smith",
				Body:           "   ",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := newFakeNoteRepo()
			svc := service.NewNoteService(repo)

			created, err := svc.Create(context.Background(), tt.input)
			if (err != nil) != tt.wantErr {
				t.Fatalf("Create() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr {
				if created.ID == uuid.Nil {
					t.Errorf("expected non-zero ID")
				}
				if created.Body != tt.input.Body {
					t.Errorf("got body %q, want %q", created.Body, tt.input.Body)
				}
			}
		})
	}
}
