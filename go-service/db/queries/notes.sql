-- name: CreateNote :one
INSERT INTO notes (
    policyholder_id,
    author,
    body
) VALUES (
    $1, $2, $3
) RETURNING *;

-- name: GetNoteByID :one
SELECT * FROM notes
WHERE id = $1;

-- name: ListNotesByPolicyholder :many
SELECT * FROM notes
WHERE policyholder_id = $1
ORDER BY created_at DESC;

-- name: DeleteNote :exec
DELETE FROM notes
WHERE id = $1;
