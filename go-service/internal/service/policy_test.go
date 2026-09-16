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

type fakePolicyRepo struct {
	lastExpiringDays int32
}

func (f *fakePolicyRepo) Create(_ context.Context, in domain.CreatePolicyInput) (domain.Policy, error) {
	return domain.Policy{ID: uuid.New(), PolicyholderID: in.PolicyholderID, Carrier: in.Carrier,
		PolicyNumber: in.PolicyNumber, EffectiveDate: in.EffectiveDate, ExpirationDate: in.ExpirationDate,
		PremiumCents: in.PremiumCents, Status: domain.PolicyStatusActive}, nil
}
func (f *fakePolicyRepo) Get(_ context.Context, id uuid.UUID) (domain.Policy, error) {
	return domain.Policy{ID: id}, nil
}
func (f *fakePolicyRepo) ListByPolicyholder(_ context.Context, _ uuid.UUID) ([]domain.Policy, error) {
	return nil, nil
}
func (f *fakePolicyRepo) ListExpiringWithin(_ context.Context, days int32) ([]domain.Policy, error) {
	f.lastExpiringDays = days
	return nil, nil
}
func (f *fakePolicyRepo) UpdateStatus(_ context.Context, id uuid.UUID, status domain.PolicyStatus) (domain.Policy, error) {
	return domain.Policy{ID: id, Status: status}, nil
}
func (f *fakePolicyRepo) Delete(_ context.Context, _ uuid.UUID) error { return nil }

func TestPolicyService_Create(t *testing.T) {
	base := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name    string
		input   domain.CreatePolicyInput
		wantErr bool
	}{
		{
			name: "valid policy succeeds",
			input: domain.CreatePolicyInput{
				Carrier: "Progressive", PolicyNumber: "PG-123",
				EffectiveDate: base, ExpirationDate: base.AddDate(1, 0, 0), PremiumCents: 120000,
			},
		},
		{
			name: "missing carrier is rejected",
			input: domain.CreatePolicyInput{
				PolicyNumber: "PG-123", EffectiveDate: base, ExpirationDate: base.AddDate(1, 0, 0),
			},
			wantErr: true,
		},
		{
			name: "expiration before effective is rejected",
			input: domain.CreatePolicyInput{
				Carrier: "Progressive", PolicyNumber: "PG-123",
				EffectiveDate: base, ExpirationDate: base.AddDate(0, 0, -1),
			},
			wantErr: true,
		},
		{
			name: "negative premium is rejected",
			input: domain.CreatePolicyInput{
				Carrier: "Progressive", PolicyNumber: "PG-123",
				EffectiveDate: base, ExpirationDate: base.AddDate(1, 0, 0), PremiumCents: -1,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := service.NewPolicyService(&fakePolicyRepo{})
			_, err := svc.Create(context.Background(), tt.input)

			if tt.wantErr && !errors.Is(err, service.ErrInvalidInput) {
				t.Fatalf("expected ErrInvalidInput, got %v", err)
			}
			if !tt.wantErr && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}

func TestPolicyService_ExpiringWithin_ClampsUpperBound(t *testing.T) {
	repo := &fakePolicyRepo{}
	svc := service.NewPolicyService(repo)

	if _, err := svc.ExpiringWithin(context.Background(), 10000); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if repo.lastExpiringDays != 365 {
		t.Errorf("expected days clamped to 365, got %d", repo.lastExpiringDays)
	}
}

func TestPolicyService_ExpiringWithin_RejectsNonPositive(t *testing.T) {
	svc := service.NewPolicyService(&fakePolicyRepo{})
	if _, err := svc.ExpiringWithin(context.Background(), 0); !errors.Is(err, service.ErrInvalidInput) {
		t.Fatalf("expected ErrInvalidInput, got %v", err)
	}
}

func TestPolicyService_UpdateStatus_RejectsUnknownStatus(t *testing.T) {
	svc := service.NewPolicyService(&fakePolicyRepo{})
	_, err := svc.UpdateStatus(context.Background(), uuid.New(), domain.PolicyStatus("not_a_status"))
	if !errors.Is(err, service.ErrInvalidInput) {
		t.Fatalf("expected ErrInvalidInput, got %v", err)
	}
}
