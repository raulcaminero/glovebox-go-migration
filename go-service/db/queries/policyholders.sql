-- name: CreatePolicyholder :one
INSERT INTO policyholders (full_name, email, phone)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetPolicyholder :one
SELECT * FROM policyholders WHERE id = $1;

-- name: ListPolicyholders :many
SELECT * FROM policyholders
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: SearchPolicyholders :many
SELECT * FROM policyholders
WHERE full_name ILIKE '%' || sqlc.arg(query)::text || '%'
   OR email ILIKE '%' || sqlc.arg(query)::text || '%'
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: UpdatePolicyholder :one
UPDATE policyholders
SET full_name = $2, email = $3, phone = $4, updated_at = now()
WHERE id = $1
RETURNING *;

-- name: DeletePolicyholder :exec
DELETE FROM policyholders WHERE id = $1;
