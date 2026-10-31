package modbusshare

import (
	"context"
	"encoding/json"
)

type formatRecoveryWitnessKey struct{}
type formatRecoveryWitness struct{ workspaceID, raw string }

// WithFormatRecoveryWitness carries server-captured durable evidence only. It
// must never be populated from a browser's desired mappings or restore flag.
func WithFormatRecoveryWitness(ctx context.Context, workspaceID, raw string) context.Context {
	return context.WithValue(ctx, formatRecoveryWitnessKey{}, formatRecoveryWitness{workspaceID, raw})
}

// CheckFormatRecoveryWitness is read-only. Final acceptance still needs its own
// transaction; this check does not turn a stale snapshot into new authority.
func (s *SQLWorkspaceRevisionStore) CheckFormatRecoveryWitness(ctx context.Context, workspaceID, witness string) error {
	var raw string
	if s == nil || s.db == nil {
		return NewError(ErrCodeProjectionRequired, "recovery storage unavailable", true)
	}
	if err := s.db.QueryRowContext(ctx, `SELECT value FROM system_settings WHERE key=?`, workspaceRevisionKey(workspaceID)).Scan(&raw); err != nil {
		return newStorageError("read recovery witness", err)
	}
	return validateFormatRecoveryWitness(workspaceID, witness, raw)
}

func validateFormatRecoveryWitness(workspaceID, witness, raw string) error {
	var state struct {
		Revision string `json:"revision"`
		Dirty    *bool  `json:"dirty"`
	}
	if workspaceID == "" || witness == "" || json.Unmarshal([]byte(raw), &state) != nil || state.Revision == "" || state.Dirty == nil || *state.Dirty {
		return NewError(ErrCodeDirtyUnknown, "format recovery requires explicit clean authority", true)
	}
	if raw != witness {
		return NewError(ErrCodeRevisionConflict, "format recovery authority changed", true)
	}
	return nil
}
