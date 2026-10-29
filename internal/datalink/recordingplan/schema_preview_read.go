package recordingplan

import (
	"context"
	"crypto/subtle"
	"database/sql"
	"errors"
)

// StoredSchemaPreviewTokenInTx authenticates historical stored preview evidence.
// It is not schema or activation authorization, and deliberately ignores expiry.
// Callers must independently verify root, source and scope bindings, and must
// never expose this internal token directly in an operator response.
func StoredSchemaPreviewTokenInTx(ctx context.Context, tx *sql.Tx, workspaceID, tokenID string) (*SchemaPreviewToken, error) {
	if tx == nil || !successorIdentity(workspaceID) || !successorIdentity(tokenID) {
		return nil, ErrSchemaOperationResult
	}
	token, err := readSchemaPreviewToken(ctx, tx, tokenID)
	if errors.Is(err, ErrPreviewTokenNotFound) {
		return nil, ErrPreviewTokenNotFound
	}
	if err != nil {
		return nil, err
	}
	if token.WorkspaceID != workspaceID || token.Token != tokenID {
		return nil, ErrSchemaOperationMismatch
	}
	if isLegacyPreviewToken(token) {
		return nil, ErrPreviewTokenLegacy
	}
	digest, err := previewDigest(token)
	if err != nil {
		return nil, err
	}
	if subtle.ConstantTimeCompare([]byte(digest), []byte(token.Digest)) != 1 {
		return nil, ErrPreviewTokenStale
	}
	return token, nil
}
