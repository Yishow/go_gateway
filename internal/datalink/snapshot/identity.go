package snapshot

import (
	"fmt"
	"strings"
	"time"

	"go-gateway/internal/datalink/measurement"
)

const (
	recordIDVersion     = "write-group-record-v1"
	effectKeyVersion    = "write-group-effect-v1"
	destinationVersion  = "write-group-destination-v1"
	fixedGroupScopeMark = "group-scope"
	entityScopeMark     = "entity"
)

func missingIdentity(field string) error {
	return fmt.Errorf("%w: %s is required for row identity", ErrInvalidConfig, field)
}

// RecordID derives the stable identity of one row from workspace, group,
// group revision, entity key and bucket start. An empty entity key means the
// fixed group scope and can never collide with a literal entity key. The
// timestamp alone is deliberately insufficient: distinct groups and entities
// in the same bucket always get distinct IDs.
func RecordID(workspaceID, groupID, groupRevision, entityKey string, bucketStart time.Time) (string, error) {
	switch {
	case strings.TrimSpace(workspaceID) == "":
		return "", missingIdentity("workspace")
	case strings.TrimSpace(groupID) == "":
		return "", missingIdentity("group")
	case strings.TrimSpace(groupRevision) == "":
		return "", missingIdentity("group revision")
	case bucketStart.IsZero():
		return "", missingIdentity("bucket start")
	}
	scope := []string{fixedGroupScopeMark}
	if entityKey != "" {
		scope = []string{entityScopeMark, entityKey}
	}
	parts := append([]string{
		recordIDVersion, workspaceID, groupID, groupRevision,
		bucketStart.UTC().Format(time.RFC3339Nano),
	}, scope...)
	return measurement.HashStringTuple(parts...)
}

// DestinationScope derives the frozen destination namespace used for effect
// keys from the saved connector identity and the destination table.
func DestinationScope(connectorID, connectorRevision, database, tableSchema, table string) (string, error) {
	if strings.TrimSpace(connectorID) == "" {
		return "", missingIdentity("connector")
	}
	if strings.TrimSpace(connectorRevision) == "" {
		return "", missingIdentity("connector revision")
	}
	if strings.TrimSpace(table) == "" {
		return "", missingIdentity("table")
	}
	return measurement.HashStringTuple(destinationVersion, connectorID, connectorRevision, database, tableSchema, table)
}

// EffectKey scopes a record to a destination so the same record resubmitted
// to the same destination reuses one key and unrelated targets never share it.
func EffectKey(destinationScope, recordID string) (string, error) {
	if strings.TrimSpace(destinationScope) == "" {
		return "", missingIdentity("destination scope")
	}
	if strings.TrimSpace(recordID) == "" {
		return "", missingIdentity("record")
	}
	return measurement.HashStringTuple(effectKeyVersion, destinationScope, recordID)
}
