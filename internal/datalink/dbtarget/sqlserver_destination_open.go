package dbtarget

import (
	"context"
	"fmt"
	"slices"

	"go-gateway/internal/datalink/schema"
)

// openSQLServerDeliveryDestination shares the normal connection manager and
// reads the actual saved scope before judging the operation-specific release
// gate. Merely registering the driver never enables delivery.
func openSQLServerDeliveryDestination(ctx context.Context, connector *schema.DatabaseConnector, config ConnectionConfig) (*OpenedDestination, error) {
	if connector == nil || connector.Kind != schema.DatabaseConnectorKindSQLServer {
		return nil, fmt.Errorf("%w: SQL Server destination is unavailable", ErrDestinationBlocked)
	}
	settings, err := sqlServerConnectionSettingsFor(config)
	if err != nil {
		return nil, fmt.Errorf("%w: SQL Server connection settings are unverified", ErrDestinationBlocked)
	}
	session, err := openSQLServerManagedSessionFunc(connector.Kind, config)
	if err != nil {
		if session != nil {
			_ = session.Close()
		}
		reason := sqlServerConnectionFailureReason(err)
		if reason == sqlServerTLSVerificationFailed || reason == inspectionReasonAuthenticationFailed {
			return nil, fmt.Errorf("%w: SQL Server transport or authentication was rejected", ErrDestinationBlocked)
		}
		return nil, fmt.Errorf("%w: SQL Server connection is unavailable", ErrDestinationUnreachable)
	}
	if session == nil || session.DB() == nil {
		if session != nil {
			_ = session.Close()
		}
		return nil, fmt.Errorf("%w: SQL Server connection is unavailable", ErrDestinationUnreachable)
	}
	metadata := inspectSQLServerRecordingMetadata(ctx, session.DB(), settings.Database, settings.Schema)
	capability := projectSQLServerRecordingCapability(connector, metadata, true)
	if capability.MetadataStatus != sqlServerRecordingMetadataVerified {
		_ = session.Close()
		if capability.MetadataStatus == sqlServerRecordingMetadataForbidden || capability.MetadataStatus == sqlServerRecordingMetadataMissing ||
			slices.Contains(capability.Reasons, "database_scope_mismatch") || slices.Contains(capability.Reasons, "schema_scope_mismatch") {
			return nil, fmt.Errorf("%w: SQL Server destination metadata or scope was rejected", ErrDestinationBlocked)
		}
		return nil, fmt.Errorf("%w: SQL Server metadata is unavailable", ErrDestinationUnreachable)
	}
	if !capability.SupportsManagedLong {
		_ = session.Close()
		return nil, fmt.Errorf("%w: SQL Server recording adapter is unverified", ErrDestinationBlocked)
	}
	return &OpenedDestination{DB: session.DB(), Kind: connector.Kind, SQLServer: capability, Close: session.Close}, nil
}
