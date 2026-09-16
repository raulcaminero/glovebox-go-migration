package repo

import (
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

// pgUUID and uuidFromPG bridge google/uuid (used in the domain layer) and
// pgtype.UUID (used by the generated sqlc code). Keeping the conversion in
// one place means the domain package never imports pgx types.
func pgUUID(id uuid.UUID) pgtype.UUID {
	return pgtype.UUID{Bytes: id, Valid: true}
}

func uuidFromPG(id pgtype.UUID) uuid.UUID {
	return uuid.UUID(id.Bytes)
}
