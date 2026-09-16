package repo

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/raul/glovebox-go-migration/go-service/internal/domain"
	"github.com/raul/glovebox-go-migration/go-service/internal/repo/sqlcgen"
)

type PolicyPG struct {
	q *sqlcgen.Queries
}

func NewPolicyPG(pool *pgxpool.Pool) *PolicyPG {
	return &PolicyPG{q: sqlcgen.New(pool)}
}

func (r *PolicyPG) Create(ctx context.Context, in domain.CreatePolicyInput) (domain.Policy, error) {
	row, err := r.q.CreatePolicy(ctx, sqlcgen.CreatePolicyParams{
		PolicyholderID: pgUUID(in.PolicyholderID),
		Carrier:        in.Carrier,
		PolicyNumber:   in.PolicyNumber,
		LineOfBusiness: in.LineOfBusiness,
		EffectiveDate:  pgtype.Date{Time: in.EffectiveDate, Valid: true},
		ExpirationDate: pgtype.Date{Time: in.ExpirationDate, Valid: true},
		PremiumCents:   in.PremiumCents,
		Status:         string(domain.PolicyStatusActive),
	})
	if err != nil {
		return domain.Policy{}, fmt.Errorf("create policy: %w", err)
	}
	return toDomainPolicy(row), nil
}

func (r *PolicyPG) Get(ctx context.Context, id uuid.UUID) (domain.Policy, error) {
	row, err := r.q.GetPolicy(ctx, pgUUID(id))
	if err != nil {
		return domain.Policy{}, fmt.Errorf("get policy: %w", err)
	}
	return toDomainPolicy(row), nil
}

func (r *PolicyPG) ListByPolicyholder(ctx context.Context, policyholderID uuid.UUID) ([]domain.Policy, error) {
	rows, err := r.q.ListPoliciesByPolicyholder(ctx, pgUUID(policyholderID))
	if err != nil {
		return nil, fmt.Errorf("list policies by policyholder: %w", err)
	}
	out := make([]domain.Policy, len(rows))
	for i, row := range rows {
		out[i] = toDomainPolicy(row)
	}
	return out, nil
}

func (r *PolicyPG) ListExpiringWithin(ctx context.Context, days int32) ([]domain.Policy, error) {
	rows, err := r.q.ListExpiringPolicies(ctx, days)
	if err != nil {
		return nil, fmt.Errorf("list expiring policies: %w", err)
	}
	out := make([]domain.Policy, len(rows))
	for i, row := range rows {
		out[i] = toDomainPolicy(row)
	}
	return out, nil
}

func (r *PolicyPG) UpdateStatus(ctx context.Context, id uuid.UUID, status domain.PolicyStatus) (domain.Policy, error) {
	row, err := r.q.UpdatePolicyStatus(ctx, sqlcgen.UpdatePolicyStatusParams{
		ID:     pgUUID(id),
		Status: string(status),
	})
	if err != nil {
		return domain.Policy{}, fmt.Errorf("update policy status: %w", err)
	}
	return toDomainPolicy(row), nil
}

func (r *PolicyPG) Delete(ctx context.Context, id uuid.UUID) error {
	if err := r.q.DeletePolicy(ctx, pgUUID(id)); err != nil {
		return fmt.Errorf("delete policy: %w", err)
	}
	return nil
}

func toDomainPolicy(row sqlcgen.Policy) domain.Policy {
	return domain.Policy{
		ID:             uuidFromPG(row.ID),
		PolicyholderID: uuidFromPG(row.PolicyholderID),
		Carrier:        row.Carrier,
		PolicyNumber:   row.PolicyNumber,
		LineOfBusiness: row.LineOfBusiness,
		EffectiveDate:  row.EffectiveDate.Time,
		ExpirationDate: row.ExpirationDate.Time,
		PremiumCents:   row.PremiumCents,
		Status:         domain.PolicyStatus(row.Status),
		CreatedAt:      row.CreatedAt.Time,
		UpdatedAt:      row.UpdatedAt.Time,
	}
}
