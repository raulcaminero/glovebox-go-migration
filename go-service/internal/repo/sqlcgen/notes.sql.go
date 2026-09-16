package sqlcgen

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const createNote = `-- name: CreateNote :one
INSERT INTO notes (
    policyholder_id,
    author,
    body
) VALUES (
    $1, $2, $3
) RETURNING id, policyholder_id, author, body, created_at
`

type CreateNoteParams struct {
	PolicyholderID pgtype.UUID
	Author         string
	Body           string
}

func (q *Queries) CreateNote(ctx context.Context, arg CreateNoteParams) (Note, error) {
	row := q.pool.QueryRow(ctx, createNote, arg.PolicyholderID, arg.Author, arg.Body)
	var i Note
	err := row.Scan(
		&i.ID,
		&i.PolicyholderID,
		&i.Author,
		&i.Body,
		&i.CreatedAt,
	)
	return i, err
}

const getNoteByID = `-- name: GetNoteByID :one
SELECT id, policyholder_id, author, body, created_at FROM notes
WHERE id = $1
`

func (q *Queries) GetNoteByID(ctx context.Context, id pgtype.UUID) (Note, error) {
	row := q.pool.QueryRow(ctx, getNoteByID, id)
	var i Note
	err := row.Scan(
		&i.ID,
		&i.PolicyholderID,
		&i.Author,
		&i.Body,
		&i.CreatedAt,
	)
	return i, err
}

const listNotesByPolicyholder = `-- name: ListNotesByPolicyholder :many
SELECT id, policyholder_id, author, body, created_at FROM notes
WHERE policyholder_id = $1
ORDER BY created_at DESC
`

func (q *Queries) ListNotesByPolicyholder(ctx context.Context, policyholderID pgtype.UUID) ([]Note, error) {
	rows, err := q.pool.Query(ctx, listNotesByPolicyholder, policyholderID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Note
	for rows.Next() {
		var i Note
		if err := rows.Scan(
			&i.ID,
			&i.PolicyholderID,
			&i.Author,
			&i.Body,
			&i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const deleteNote = `-- name: DeleteNote :exec
DELETE FROM notes
WHERE id = $1
`

func (q *Queries) DeleteNote(ctx context.Context, id pgtype.UUID) error {
	_, err := q.pool.Exec(ctx, deleteNote, id)
	return err
}
