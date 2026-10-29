package workspace

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"time"
)

// Intent normalization is deliberately independent of live connectors/source
// rows: a frozen retry is a readback, not authorization for another mutation.
func normalizeRecreateReservationRequest(req writeGroupRecreateReservationRequest) (writeGroupRecreateReservationRequest, string, error) {
	fields := []*string{&req.WorkspaceID, &req.ExpectedWorkspaceRevision, &req.LineageID, &req.PredecessorGroupID, &req.ExpectedPredecessorRevision, &req.ExpectedPredecessorConnectorRevision, &req.ExpectedTargetConnectorRevision, &req.OperationID}
	for _, field := range fields {
		*field = strings.TrimSpace(*field)
		if *field == "" {
			return req, "", fmt.Errorf("incomplete successor reservation: %w", ErrWriteGroupValidation)
		}
	}
	if req.ExpectedLineageRevision < 1 || req.Group == nil {
		return req, "", fmt.Errorf("incomplete successor reservation: %w", ErrWriteGroupValidation)
	}
	group := cloneWriteGroup(req.Group)
	if strings.TrimSpace(group.WorkspaceID) != req.WorkspaceID {
		return req, "", recreateReservationConflict("scope_changed")
	}
	if provided := strings.TrimSpace(group.Destination.ConnectorRevision); provided != "" && provided != req.ExpectedTargetConnectorRevision {
		return req, "", recreateReservationConflict("target_connector_changed")
	}
	group.WorkspaceID = req.WorkspaceID
	group.ID, group.Revision, group.AppliedRevision, group.BasicManagedDeviceID = "", "", "", ""
	group.Status = WriteGroupStatusDraft
	group.CreatedAt, group.UpdatedAt = time.Time{}, time.Time{}
	group.Migration = WriteGroupMigration{}
	group.Destination.SchemaRevision, group.Destination.SchemaDigest = "", ""
	group.Destination.ConnectorRevision = req.ExpectedTargetConnectorRevision
	group.Name = strings.TrimSpace(group.Name)
	group.Destination.ConnectorID = strings.TrimSpace(group.Destination.ConnectorID)
	group.Destination.Database = strings.TrimSpace(group.Destination.Database)
	group.Destination.TableSchema = strings.TrimSpace(group.Destination.TableSchema)
	group.Destination.TableName = strings.TrimSpace(group.Destination.TableName)
	if (group.isWideProfile() || group.RecordingProfile == nil && group.EffectiveLayoutVersion() == WriteGroupLayoutManagedWideV1) && group.Destination.TableName == managedGroupTableName(req.PredecessorGroupID) {
		group.Destination.TableName = ""
	}
	group.RowPolicy.GroupKeyColumns = normalizeStringSet(group.RowPolicy.GroupKeyColumns)
	group.RowPolicy.UniqueKeyColumns = normalizeStringSet(group.RowPolicy.UniqueKeyColumns)
	for i := range group.Members {
		m := &group.Members[i]
		if strings.TrimSpace(m.EntityKey) == "" {
			m.EntityKey = ""
		}
		for _, field := range []*string{&m.DeviceID, &m.PointID, &m.TagID, &m.SourceRevision, &m.MappingRevision, &m.CircuitID, &m.CircuitName, &m.TargetColumn} {
			*field = strings.TrimSpace(*field)
		}
	}
	sort.SliceStable(group.Members, func(i, j int) bool {
		a, b := group.Members[i], group.Members[j]
		return a.DeviceID+"\x00"+a.PointID+"\x00"+a.TagID < b.DeviceID+"\x00"+b.PointID+"\x00"+b.TagID
	})
	group.RecordingProfile = canonicalDigestProfile(group.RecordingProfile)
	req.Group = group
	encoded, err := json.Marshal(req)
	if err != nil {
		return req, "", fmt.Errorf("encode successor intent: %w", ErrWriteGroupValidation)
	}
	digest := sha256.Sum256(encoded)
	return req, hex.EncodeToString(digest[:]), nil
}
