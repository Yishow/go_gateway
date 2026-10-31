package tag

import (
	"context"
	"database/sql"
	"go-gateway/internal/datalink/schema"
)

// ReadByIDTx reads a detached tag using the caller's existing transaction and
// the repository's canonical decoder. It never opens another connection.
func ReadByIDTx(ctx context.Context, tx *sql.Tx, id string) (*schema.Tag, error) {
	if tx == nil {
		return nil, ErrTagNotFound
	}
	reader := &SQLRepository{}
	value, err := reader.scanTag(tx.QueryRowContext(ctx, `SELECT id,key,display_name,data_type,unit,description,status,labels,created_at,updated_at FROM tags WHERE id=?`, id))
	if err != nil {
		return nil, err
	}
	if err = tx.QueryRowContext(ctx, `SELECT key_lower FROM tags WHERE id=?`, id).Scan(&value.KeyLower); err != nil {
		return nil, err
	}
	return value, nil
}
