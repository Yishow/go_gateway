package grouppipeline

import (
	"crypto/sha256"
	"encoding/hex"
	"strings"

	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/groupdelivery"
	"go-gateway/internal/datalink/recordingplan"
	"go-gateway/internal/datalink/workspace"
)

func mixedRecordingGroup(group *workspace.WriteGroup) bool {
	return group != nil && group.RecordingProfile != nil && group.RecordingProfile.Kind == workspace.WriteGroupProfileMixed
}

// validMixedRuntimeStorageProof checks the frozen storage projection, not
// counter replay meaning or live schema readiness. Finalized rows deliberately
// retain delivery even when replay-only counter bindings cannot be loaded.
func validMixedRuntimeStorageProof(group *workspace.WriteGroup, layout *workspace.WriteGroupRuntimeLayout) bool {
	if !mixedRecordingGroup(group) || layout == nil || group.WorkspaceID == "" || group.ID == "" || group.AppliedRevision == "" || group.AppliedRevision != group.Revision || group.RowPolicy.IntervalSeconds != 0 {
		return false
	}
	if layout.Dialect != dbtarget.SQLDialectSQLite && layout.Dialect != dbtarget.SQLDialectPostgres {
		return false
	}
	if (group.StorageMode != "" && group.StorageMode != workspace.WriteGroupStorageModePlain) || layout.Hypertable != nil || layout.MySQL != nil || layout.SQLServer != nil {
		return false
	}
	// Confirmation and inspection digests have different inputs and purposes.
	// Preserve their distinct presence/format contracts; equality proves nothing.
	if strings.TrimSpace(group.Destination.SchemaRevision) == "" || !frozenDigestPresent(group.Destination.SchemaDigest) || !frozenDigestPresent(layout.SchemaDigest) {
		return false
	}
	expected, err := workspace.MixedManagedSchemaLayout(group, layout.TagTypes, string(layout.Dialect))
	if err != nil || expected.TableName != group.Destination.TableName || len(layout.Columns) == 0 {
		return false
	}
	actual := recordingplan.TargetTableInspection{Status: "exists", ColumnTypes: map[string]string{}, ColumnNullable: map[string]bool{}, ColumnPrimaryKey: map[string]bool{}}
	for _, c := range layout.Columns {
		actual.Columns = append(actual.Columns, c.Name)
		actual.ColumnTypes[c.Name] = c.DataType
		actual.ColumnNullable[c.Name] = c.Nullable
		actual.ColumnPrimaryKey[c.Name] = c.PrimaryKey
	}
	if err = recordingplan.ValidateManagedGroupColumns(*expected, actual); err != nil {
		return false
	}
	return !strings.EqualFold(strings.TrimSpace(group.WritePolicy.DedupeCapability), groupdelivery.DedupeReceipt) || layout.ReceiptTableReady
}
func frozenDigestPresent(raw string) bool {
	value := strings.TrimSpace(raw)
	if len(value) != sha256.Size*2 {
		return false
	}
	_, err := hex.DecodeString(value)
	return err == nil
}
