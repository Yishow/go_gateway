package dbtarget

import (
	"context"
	"fmt"
	"strings"

	"go-gateway/internal/datalink/schema"
)

// ProbeTimescaleAtRevision runs the read-only capability probe on the confirmed
// saved connector revision, through the same destination guard as table
// inspection. A SQLite destination never has the capability: it reports the
// dialect reason instead of opening anything, so the selection is refused and
// not rewritten.
func (i *ManagedTableInspector) ProbeTimescaleAtRevision(ctx context.Context, connectorID, expectedRevision, schemaName string) (*TimescaleCapability, error) {
	if i == nil || i.service == nil || i.service.repo == nil {
		return nil, managedDestinationInvalid("saved connector service is unavailable")
	}
	connector, err := i.service.ResolveSavedTarget(ctx, strings.TrimSpace(connectorID), strings.TrimSpace(expectedRevision))
	if err != nil {
		return nil, err
	}
	if connector.Kind != schema.DatabaseConnectorKindPostgres {
		return &TimescaleCapability{Reasons: []string{TimescaleReasonDialect}}, nil
	}
	if _, err := i.guardManagedConnector(ctx, connector, schemaName); err != nil {
		return nil, err
	}
	config, err := parseConnectionConfig(connector.ConnectionConfig)
	if err != nil {
		return nil, managedDestinationInvalid("saved connector configuration is invalid")
	}
	manager, err := openExternalDBManagerFunc(connector.Kind, config)
	if err != nil {
		return nil, fmt.Errorf("open PostgreSQL destination for the timescale probe: %w", err)
	}
	defer manager.Close()
	capability, err := ProbeTimescale(ctx, manager.DB(), defaultString(strings.TrimSpace(schemaName), "public"))
	if err != nil {
		return nil, err
	}
	return &capability, nil
}

// InspectHypertableAtRevision reads the hypertable catalog evidence of one table
// at the confirmed connector revision. It returns nil for a table that is not a
// hypertable, and never creates or converts anything.
func (i *ManagedTableInspector) InspectHypertableAtRevision(ctx context.Context, connectorID, expectedRevision, schemaName, tableName string) (*HypertableInfo, error) {
	if i == nil || i.service == nil || i.service.repo == nil {
		return nil, managedDestinationInvalid("saved connector service is unavailable")
	}
	connector, err := i.service.ResolveSavedTarget(ctx, strings.TrimSpace(connectorID), strings.TrimSpace(expectedRevision))
	if err != nil {
		return nil, err
	}
	if connector.Kind != schema.DatabaseConnectorKindPostgres {
		return nil, nil
	}
	if _, err := i.guardManagedConnector(ctx, connector, schemaName); err != nil {
		return nil, err
	}
	config, err := parseConnectionConfig(connector.ConnectionConfig)
	if err != nil {
		return nil, managedDestinationInvalid("saved connector configuration is invalid")
	}
	manager, err := openExternalDBManagerFunc(connector.Kind, config)
	if err != nil {
		return nil, fmt.Errorf("open PostgreSQL destination for hypertable inspection: %w", err)
	}
	defer manager.Close()
	return InspectHypertable(ctx, manager.DB(), defaultString(strings.TrimSpace(schemaName), "public"), strings.TrimSpace(tableName))
}
