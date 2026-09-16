-- name: CreatePolicy :one
INSERT INTO policies (
    policyholder_id, carrier, policy_number, line_of_business,
    effective_date, expiration_date, premium_cents, status
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING *;

-- name: GetPolicy :one
SELECT * FROM policies WHERE id = $1;

-- name: ListPoliciesByPolicyholder :many
SELECT * FROM policies
WHERE policyholder_id = $1
ORDER BY effective_date DESC;

-- name: ListExpiringPolicies :many
SELECT * FROM policies
WHERE expiration_date BETWEEN now()::date AND (now()::date + sqlc.arg(within_days)::int)
  AND status = 'active'
ORDER BY expiration_date ASC;

-- name: UpdatePolicyStatus :one
UPDATE policies
SET status = $2, updated_at = now()
WHERE id = $1
RETURNING *;

-- name: DeletePolicy :exec
DELETE FROM policies WHERE id = $1;
