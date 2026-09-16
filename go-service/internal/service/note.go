package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"

	"github.com/raul/glovebox-go-migration/go-service/internal/domain"
)

type NoteService struct {
	repo domain.NoteRepo
}

func NewNoteService(repo domain.NoteRepo) *NoteService {
	return &NoteService{repo: repo}
}

func (s *NoteService) Create(ctx context.Context, in domain.CreateNoteInput) (domain.Note, error) {
	in.Author = strings.TrimSpace(in.Author)
	in.Body = strings.TrimSpace(in.Body)

	if in.PolicyholderID == uuid.Nil {
		return domain.Note{}, fmt.Errorf("%w: policyholder_id is required", ErrInvalidInput)
	}
	if in.Author == "" {
		return domain.Note{}, fmt.Errorf("%w: author cannot be empty", ErrInvalidInput)
	}
	if in.Body == "" {
		return domain.Note{}, fmt.Errorf("%w: note body cannot be empty", ErrInvalidInput)
	}

	return s.repo.Create(ctx, in)
}

func (s *NoteService) Get(ctx context.Context, id uuid.UUID) (domain.Note, error) {
	if id == uuid.Nil {
		return domain.Note{}, fmt.Errorf("%w: id is required", ErrInvalidInput)
	}
	return s.repo.Get(ctx, id)
}

func (s *NoteService) ListByPolicyholder(ctx context.Context, policyholderID uuid.UUID) ([]domain.Note, error) {
	if policyholderID == uuid.Nil {
		return nil, fmt.Errorf("%w: policyholder_id is required", ErrInvalidInput)
	}
	return s.repo.ListByPolicyholder(ctx, policyholderID)
}

func (s *NoteService) Delete(ctx context.Context, id uuid.UUID) error {
	if id == uuid.Nil {
		return fmt.Errorf("%w: id is required", ErrInvalidInput)
	}
	return s.repo.Delete(ctx, id)
}
