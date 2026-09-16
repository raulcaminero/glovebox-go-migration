// Package service holds business logic. Handlers stay thin (parse request,
// call service, write response); services own validation and orchestration
// across repos. This separation is what makes the domain testable without
// spinning up an HTTP server.
package service

import (
	"context"
	"errors"
	"strings"

	"github.com/google/uuid"

	"github.com/raul/glovebox-go-migration/go-service/internal/domain"
)

var ErrInvalidInput = errors.New("invalid input")

type PolicyholderService struct {
	repo domain.PolicyholderRepo
}

func NewPolicyholderService(repo domain.PolicyholderRepo) *PolicyholderService {
	return &PolicyholderService{repo: repo}
}

func (s *PolicyholderService) Create(ctx context.Context, in domain.CreatePolicyholderInput) (domain.Policyholder, error) {
	in.FullName = strings.TrimSpace(in.FullName)
	in.Email = strings.TrimSpace(strings.ToLower(in.Email))

	if in.FullName == "" {
		return domain.Policyholder{}, wrapInvalid("full_name is required")
	}
	if !strings.Contains(in.Email, "@") {
		return domain.Policyholder{}, wrapInvalid("email is invalid")
	}
	return s.repo.Create(ctx, in)
}

func (s *PolicyholderService) Get(ctx context.Context, id uuid.UUID) (domain.Policyholder, error) {
	return s.repo.Get(ctx, id)
}

func (s *PolicyholderService) List(ctx context.Context, p domain.ListParams) ([]domain.Policyholder, error) {
	if p.Limit <= 0 || p.Limit > 100 {
		p.Limit = 25 // sane default + hard ceiling to prevent unbounded scans
	}
	return s.repo.List(ctx, p)
}

func (s *PolicyholderService) Update(ctx context.Context, in domain.UpdatePolicyholderInput) (domain.Policyholder, error) {
	in.FullName = strings.TrimSpace(in.FullName)
	if in.FullName == "" {
		return domain.Policyholder{}, wrapInvalid("full_name is required")
	}
	return s.repo.Update(ctx, in)
}

func (s *PolicyholderService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}

func wrapInvalid(msg string) error {
	return errors.Join(ErrInvalidInput, errors.New(msg))
}
