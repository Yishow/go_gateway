package dbtarget

import "context"

// ProbeSQLServerRecordingAtRevision exposes the saved connector's existing
// read-only probe to workspace schema preparation without a second connection
// manager or a separate capability decision.
func (i *ManagedTableInspector) ProbeSQLServerRecordingAtRevision(ctx context.Context, connectorID, expectedRevision string) (*SQLServerRecordingCapability, error) {
	if i == nil || i.service == nil || i.service.repo == nil {
		return nil, managedDestinationInvalid("saved connector service is unavailable")
	}
	return i.service.ProbeSQLServerRecordingAtRevision(ctx, connectorID, expectedRevision)
}
