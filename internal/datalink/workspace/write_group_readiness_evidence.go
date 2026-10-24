package workspace

import (
	"cmp"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"slices"
	"strings"

	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/recordingplan"
)

type writeGroupSchemaDigest struct {
	GroupID     string                    `json:"group_id"`
	ConnectorID string                    `json:"connector_id"`
	Schema      string                    `json:"schema"`
	Table       string                    `json:"table"`
	Columns     []dbtarget.ColumnInfo     `json:"columns"`
	Bindings    []writeGroupSchemaBinding `json:"bindings"`
	// Indexes joins the digest of a long table so that a changed index
	// invalidates the confirmation like a changed column does.
	Indexes []dbtarget.IndexInfo `json:"indexes,omitempty"`
	// Hypertable joins the digest of a timescale table so that a changed time
	// dimension or interval invalidates the confirmation the same way.
	Hypertable *dbtarget.HypertableInfo `json:"hypertable,omitempty"`
}

type writeGroupSchemaBinding struct {
	DeviceID     string `json:"device_id"`
	PointID      string `json:"point_id"`
	TagID        string `json:"tag_id"`
	TargetColumn string `json:"target_column"`
}

func recordingPlanTableInspection(inspection *dbtarget.TableInspection) recordingplan.TargetTableInspection {
	if inspection == nil {
		return recordingplan.TargetTableInspection{}
	}
	state := recordingplan.TargetTableInspection{
		Status:           string(inspection.Status),
		Columns:          make([]string, 0, len(inspection.Columns)),
		ColumnTypes:      make(map[string]string, len(inspection.Columns)),
		ColumnNullable:   make(map[string]bool, len(inspection.Columns)),
		ColumnPrimaryKey: make(map[string]bool, len(inspection.Columns)),
	}
	for _, column := range inspection.Columns {
		state.Columns = append(state.Columns, column.Name)
		state.ColumnTypes[column.Name] = column.DataType
		state.ColumnNullable[column.Name] = column.Nullable
		state.ColumnPrimaryKey[column.Name] = column.PrimaryKey
	}
	if inspection.HypertableInspected {
		state.Timescale = &recordingplan.TargetTimescaleEvidence{}
		if hypertable := inspection.Hypertable; hypertable != nil {
			state.Timescale.Hypertable = &recordingplan.TargetHypertable{
				TimeColumn: hypertable.TimeColumn, ChunkInterval: hypertable.ChunkInterval, Dimensions: hypertable.Dimensions,
				CompressionEnabled: hypertable.CompressionEnabled, Jobs: hypertable.Jobs,
			}
		}
	}
	state.IndexesInspected = inspection.IndexesInspected
	for _, index := range inspection.Indexes {
		state.Indexes = append(state.Indexes, recordingplan.TargetIndexInspection{Name: index.Name, Columns: index.Columns, Unique: index.Unique})
	}
	return state
}

func digestWriteGroupInspection(group *WriteGroup, inspection *dbtarget.TableInspection) (string, error) {
	digest := writeGroupSchemaDigest{
		GroupID: group.ID, ConnectorID: group.Destination.ConnectorID,
		Schema: group.Destination.TableSchema, Table: group.Destination.TableName,
		Columns:    slices.Clone(inspection.Columns),
		Bindings:   make([]writeGroupSchemaBinding, 0, len(group.Members)),
		Indexes:    slices.Clone(inspection.Indexes),
		Hypertable: inspection.Hypertable,
	}
	slices.SortFunc(digest.Columns, func(a, b dbtarget.ColumnInfo) int {
		left, right := strings.ToLower(a.Name), strings.ToLower(b.Name)
		if c := cmp.Compare(left, right); c != 0 {
			return c
		}
		return cmp.Compare(a.Name, b.Name)
	})
	for _, member := range group.Members {
		digest.Bindings = append(digest.Bindings, writeGroupSchemaBinding{
			DeviceID: member.DeviceID, PointID: member.PointID,
			TagID: member.TagID, TargetColumn: member.TargetColumn,
		})
	}
	payload, err := json.Marshal(digest)
	if err != nil {
		return "", err
	}
	sum := sha256.Sum256(payload)
	return hex.EncodeToString(sum[:]), nil
}
