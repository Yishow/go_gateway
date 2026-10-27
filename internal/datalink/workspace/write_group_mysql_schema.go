package workspace

import (
	"context"
	"fmt"
	"strings"

	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/recordingplan"
)

// ManagedWriteGroupMySQLProbeInspector reads live MySQL/MariaDB recording capability
// facts at a specific connector identity revision.
type ManagedWriteGroupMySQLProbeInspector interface {
	ProbeMySQLRecordingAtRevision(ctx context.Context, connectorID, expectedRevision string) (*dbtarget.MySQLRecordingCapability, error)
}

// managedMySQLSchemaLayout builds a verified MySQL long layout strictly from
// the live server capability probe facts and fixed family collation mappings,
// never accepting client-supplied family or engine assertions.
func (s *WriteGroupService) managedMySQLSchemaLayout(ctx context.Context, group *WriteGroup) (*recordingplan.GroupSchemaLayout, error) {
	layout, _, err := s.managedMySQLSchemaLayoutWithCapability(ctx, group)
	return layout, err
}

func (s *WriteGroupService) managedMySQLSchemaLayoutWithCapability(ctx context.Context, group *WriteGroup) (*recordingplan.GroupSchemaLayout, *dbtarget.MySQLRecordingCapability, error) {
	if s == nil || s.managedTableInspector == nil {
		return nil, nil, fmt.Errorf("%w: managed inspector is unavailable", ErrWriteGroupValidation)
	}
	prober, ok := s.managedTableInspector.(ManagedWriteGroupMySQLProbeInspector)
	if !ok {
		return nil, nil, fmt.Errorf("%w: mysql capability probe is not supported by inspector", ErrWriteGroupValidation)
	}

	capability, err := prober.ProbeMySQLRecordingAtRevision(ctx, group.Destination.ConnectorID, group.Destination.ConnectorRevision)
	if err != nil {
		return nil, nil, fmt.Errorf("%w: probe mysql capability failed: %w", ErrWriteGroupValidation, err)
	}
	if capability == nil {
		return nil, nil, fmt.Errorf("%w: mysql capability probe returned nil", ErrWriteGroupValidation)
	}

	if !capability.Connected {
		return nil, nil, fmt.Errorf("%w: mysql destination is not connected", ErrWriteGroupValidation)
	}
	if capability.MetadataStatus != "verified" {
		return nil, nil, fmt.Errorf("%w: mysql destination metadata is not verified (status=%s)", ErrWriteGroupValidation, capability.MetadataStatus)
	}
	if !capability.TransactionalEngine {
		return nil, nil, fmt.Errorf("%w: mysql destination requires a transactional engine", ErrWriteGroupValidation)
	}
	if capability.ConnectorID != group.Destination.ConnectorID {
		return nil, nil, fmt.Errorf("%w: mysql probe connector mismatch", ErrWriteGroupValidation)
	}
	if capability.ConnectorRevision != group.Destination.ConnectorRevision {
		return nil, nil, fmt.Errorf("%w: mysql probe connector revision mismatch", ErrWriteGroupValidation)
	}
	if capability.Database == "" || capability.Database != group.Destination.Database {
		return nil, nil, fmt.Errorf("%w: mysql probe database mismatch: %s != %s", ErrWriteGroupValidation, capability.Database, group.Destination.Database)
	}

	if !dbtarget.MySQLRecordingVersionVerified(capability.ServerFamily, capability.ServerVersion) {
		return nil, nil, fmt.Errorf("%w: unsupported mysql server version: %s", ErrWriteGroupValidation, capability.ServerVersion)
	}
	family := strings.ToLower(strings.TrimSpace(capability.ServerFamily))
	version := capability.ServerVersion
	var collation string
	switch family {
	case writeGroupConnectorKindMySQL:
		collation = "utf8mb4_0900_bin"
	case "mariadb":
		collation = "utf8mb4_nopad_bin"
	default:
		return nil, nil, fmt.Errorf("%w: unsupported mysql server family: %s", ErrWriteGroupValidation, family)
	}

	tableName, ownerColumn, ok := LongStorageIdentity(group)
	if !ok {
		return nil, nil, fmt.Errorf("%w: group is not a long group", ErrWriteGroupValidation)
	}

	options := recordingplan.MySQLGroupSchemaOptions{
		ServerFamily:       family,
		ServerVersion:      version,
		CapabilityRevision: capability.CapabilityRevision,
		Engine:             "InnoDB",
		Collation:          collation,
		Charset:            "utf8mb4",
	}
	if group.StorageMode == WriteGroupStorageModeNativeMonthly {
		options.NativeMonthly, err = s.mysqlNativeMonthlyLayout(ctx, group)
		if err != nil {
			return nil, nil, err
		}
	}

	layout, err := recordingplan.ManagedMySQLLongGroupLayout(tableName, ownerColumn, options)
	if err != nil {
		return nil, nil, fmt.Errorf("%w: build mysql long layout: %w", ErrWriteGroupValidation, err)
	}
	return &layout, capability, nil
}
