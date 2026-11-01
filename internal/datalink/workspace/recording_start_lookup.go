package workspace

import (
	"context"
	"errors"

	"go-gateway/internal/datalink/recordingplan"
)

// Lookup recovers a lost acknowledgement using the original complete request.
// It reads only the existing workspace and ledger, never claims or resumes an
// operation, and deliberately does not validate today's source/configuration.
func (s *RecordingStartService) Lookup(ctx context.Context, request RecordingStartRequest) (*RecordingStartOperation, error) {
	if s == nil || s.workspace == nil || s.ledger == nil {
		return nil, ErrWriteGroupServiceUnavailable
	}
	request, digest, err := normalizeRecordingStart(request)
	if err != nil {
		return nil, err
	}
	record, err := s.workspace.repo.Get(ctx)
	if errors.Is(err, ErrNotFound) {
		return nil, recordingplan.ErrSchemaOperationNotFound
	}
	if err != nil {
		return nil, err
	}
	if record.ID != request.WorkspaceID {
		return nil, ErrRecordingStartInvalid
	}
	operation, err := s.ledger.FindRecordingStart(ctx, record.ID, request.RequestID, digest)
	if err != nil {
		return nil, err
	}
	progress, err := decodeRecordingStart(operation)
	if err != nil {
		return nil, err
	}
	return recordingStartView(operation, progress), nil
}
