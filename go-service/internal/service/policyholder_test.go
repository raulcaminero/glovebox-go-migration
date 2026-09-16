package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"

	"github.com/raul/glovebox-go-migration/go-service/internal/domain"
	"github.com/raul/glovebox-go-migration/go-service/internal/service"
)

// fakePolicyholderRepo is an in-memory stand-in for domain.PolicyholderRepo.
// Because the service layer depends on the interface (not on pgx/sqlc), we
// can test all business rules here with zero database, zero containers,
// milliseconds per test.
type fakePolicyholderRepo struct {
	created    domain.CreatePolicyholderInput
	lastParams domain.ListParams
}

func (f *fakePolicyholderRepo) Create(_ context.Context, in domain.CreatePolicyholderInput) (domain.Policyholder, error) {
	f.created = in
	return domain.Policyholder{ID: uuid.New(), FullName: in.FullName, Email: in.Email, Phone: in.Phone}, nil
}
func (f *fakePolicyholderRepo) Get(_ context.Context, id uuid.UUID) (domain.Policyholder, error) {
	return domain.Policyholder{ID: id}, nil
}
func (f *fakePolicyholderRepo) List(_ context.Context, p domain.ListParams) ([]domain.Policyholder, error) {
	f.lastParams = p
	return nil, nil
}
func (f *fakePolicyholderRepo) Update(_ context.Context, in domain.UpdatePolicyholderInput) (domain.Policyholder, error) {
	return domain.Policyholder{ID: in.ID, FullName: in.FullName, Email: in.Email}, nil
}
func (f *fakePolicyholderRepo) Delete(_ context.Context, _ uuid.UUID) error { return nil }

func TestPolicyholderService_Create(t *testing.T) {
	tests := []struct {
		name    string
		input   domain.CreatePolicyholderInput
		wantErr bool
	}{
		{
			name:  "valid input succeeds and normalizes email",
			input: domain.CreatePolicyholderInput{FullName: "  Jane Doe ", Email: " JANE@Example.com "},
		},
		{
			name:    "empty name is rejected",
			input:   domain.CreatePolicyholderInput{FullName: "   ", Email: "jane@example.com"},
			wantErr: true,
		},
		{
			name:    "email without @ is rejected",
			input:   domain.CreatePolicyholderInput{FullName: "Jane Doe", Email: "not-an-email"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &fakePolicyholderRepo{}
			svc := service.NewPolicyholderService(repo)

			_, err := svc.Create(context.Background(), tt.input)

			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				if !errors.Is(err, service.ErrInvalidInput) {
					t.Fatalf("expected ErrInvalidInput, got %v", err)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if repo.created.Email != "jane@example.com" {
				t.Errorf("expected email to be trimmed+lowercased, got %q", repo.created.Email)
			}
			if repo.created.FullName != "Jane Doe" {
				t.Errorf("expected full_name to be trimmed, got %q", repo.created.FullName)
			}
		})
	}
}

func TestPolicyholderService_List_DefaultsLimit(t *testing.T) {
	repo := &fakePolicyholderRepo{}
	svc := service.NewPolicyholderService(repo)

	// Limit of 0 and out-of-range limits should both be clamped by the
	// service — the repo/DB should never see an unbounded or absurd LIMIT.
	for _, limit := range []int32{0, -5, 500} {
		if _, err := svc.List(context.Background(), domain.ListParams{Limit: limit}); err != nil {
			t.Fatalf("unexpected error for limit=%d: %v", limit, err)
		}
		if repo.lastParams.Limit != 25 {
			t.Errorf("input limit=%d: expected clamped limit 25, got %d", limit, repo.lastParams.Limit)
		}
	}
}
