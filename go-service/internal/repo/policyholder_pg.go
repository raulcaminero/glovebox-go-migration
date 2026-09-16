// Package repo contains the Postgres adapters that implement the domain
// repo interfaces. sqlc generates the low-level query code into
// internal/repo/sqlcgen; this file adapts that generated code to the
// domain.PolicyholderRepo port so the service layer stays DB-agnostic.
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

type PolicyholderPG struct {
	q *sqlcgen.Queries
}

func NewPolicyholderPG(pool *pgxpool.Pool) *PolicyholderPG {
	return &PolicyholderPG{q: sqlcgen.New(pool)}
}

func (r *PolicyholderPG) Create(ctx context.Context, in domain.CreatePolicyholderInput) (domain.Policyholder, error) {
	row, err := r.q.CreatePolicyholder(ctx, sqlcgen.CreatePolicyholderParams{
		FullName: in.FullName,
		Email:    in.Email,
		Phone:    pgtype.Text{String: in.Phone, Valid: in.Phone != ""},
	})
	if err != nil {
		return domain.Policyholder{}, fmt.Errorf("create policyholder: %w", err)
	}
	return toDomainPolicyholder(row), nil
}

func (r *PolicyholderPG) Get(ctx context.Context, id uuid.UUID) (domain.Policyholder, error) {
	row, err := r.q.GetPolicyholder(ctx, pgUUID(id))
	if err != nil {
		return domain.Policyholder{}, fmt.Errorf("get policyholder: %w", err)
	}
	return toDomainPolicyholder(row), nil
}

func (r *PolicyholderPG) List(ctx context.Context, p domain.ListParams) ([]domain.Policyholder, error) {
	var rows []sqlcgen.Policyholder
	var err error

	if p.Query != "" {
		rows, err = r.q.SearchPolicyholders(ctx, sqlcgen.SearchPolicyholdersParams{
			Query:  p.Query,
			Limit:  p.Limit,
			Offset: p.Offset,
		})
	} else {
		rows, err = r.q.ListPolicyholders(ctx, sqlcgen.ListPolicyholdersParams{
			Limit:  p.Limit,
			Offset: p.Offset,
		})
	}
	if err != nil {
		return nil, fmt.Errorf("list policyholders: %w", err)
	}

	out := make([]domain.Policyholder, len(rows))
	for i, row := range rows {
		out[i] = toDomainPolicyholder(row)
	}
	return out, nil
}

func (r *PolicyholderPG) Update(ctx context.Context, in domain.UpdatePolicyholderInput) (domain.Policyholder, error) {
	row, err := r.q.UpdatePolicyholder(ctx, sqlcgen.UpdatePolicyholderParams{
		ID:       pgUUID(in.ID),
		FullName: in.FullName,
		Email:    in.Email,
		Phone:    pgtype.Text{String: in.Phone, Valid: in.Phone != ""},
	})
	if err != nil {
		return domain.Policyholder{}, fmt.Errorf("update policyholder: %w", err)
	}
	return toDomainPolicyholder(row), nil
}

func (r *PolicyholderPG) Delete(ctx context.Context, id uuid.UUID) error {
	if err := r.q.DeletePolicyholder(ctx, pgUUID(id)); err != nil {
		return fmt.Errorf("delete policyholder: %w", err)
	}
	return nil
}

func toDomainPolicyholder(row sqlcgen.Policyholder) domain.Policyholder {
	return domain.Policyholder{
		ID:        uuidFromPG(row.ID),
		FullName:  row.FullName,
		Email:     row.Email,
		Phone:     row.Phone.String,
		CreatedAt: row.CreatedAt.Time,
		UpdatedAt: row.UpdatedAt.Time,
	}
}
