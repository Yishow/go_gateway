package handlers

import (
	"context"
	"database/sql"
	"strings"

	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/workspace"
)

// prepareWorkspaceConnector validates and probes the connector change before
// the local setup transaction opens; nothing is written here.
func (h *StudioV2WorkspaceDatabaseHandler) prepareWorkspaceConnector(ctx context.Context, record *workspace.Record, req studioV2WorkspaceDatabaseConfigRequest) (*dbtarget.PreparedConnector, error) {
	if req.ConnectorID != "" {
		saved, err := h.connectorSvc.ResolveSavedTarget(ctx, req.ConnectorID, req.ExpectedConnectorRevision)
		if err != nil {
			return nil, err
		}
		if saved.ID != record.DatabaseConnectorID {
			// Settings exposes the saved gateway connector pool. Selection binds
			// that reference; masked browser fields never replace its credentials.
			return dbtarget.PrepareExisting(saved), nil
		}
	}
	connectionConfig := dbtarget.ConnectionConfig{
		"host":             strings.TrimSpace(req.Host),
		"port":             req.Port,
		"user":             strings.TrimSpace(req.Username),
		"database":         strings.TrimSpace(req.Database),
		"schema":           strings.TrimSpace(req.Schema),
		"table":            strings.TrimSpace(req.Table),
		"write_mode":       string(req.WriteMode),
		"timestamp_column": strings.TrimSpace(req.TimestampColumn),
	}
	if req.Kind == schema.DatabaseConnectorKindSQLite {
		connectionConfig["dsn"] = strings.TrimSpace(req.Database)
	}
	// 僅在使用者實際輸入密碼時才寫入；更新時留空則交由
	// preserveSensitiveConnectionConfigValues 保留既有密碼。
	if workspaceDatabasePasswordRequired(req.Kind) && !req.ClearPassword && req.Password != "" {
		connectionConfig["password"] = req.Password
	}
	clearPassword := req.ClearPassword || !workspaceDatabasePasswordRequired(req.Kind)

	if strings.TrimSpace(record.DatabaseConnectorID) == "" {
		return h.connectorSvc.PrepareCreate(ctx, dbtarget.CreateConnectorRequest{
			Name:                        strings.TrimSpace(req.Name),
			Kind:                        req.Kind,
			ConnectionConfig:            connectionConfig,
			DefaultWriteIntervalSeconds: workspaceOptionalInt(req.WriteIntervalSeconds),
		})
	}

	return h.connectorSvc.PrepareUpdate(ctx, record.DatabaseConnectorID, dbtarget.UpdateConnectorRequest{
		ExpectedIdentityRevision:    optionalConnectorRevision(req.ExpectedConnectorRevision),
		Name:                        workspaceOptionalString(strings.TrimSpace(req.Name)),
		Kind:                        &req.Kind,
		ConnectionConfig:            &connectionConfig,
		ClearPassword:               boolPtr(clearPassword),
		DefaultWriteIntervalSeconds: workspaceOptionalInt(req.WriteIntervalSeconds),
	})
}

// commitWorkspaceConnectorSetup persists the prepared connector, its workspace
// binding and optional row groups in one local transaction guarded by the
// expected setup revision. boundConnectorID is the binding the change was
// prepared against; a different binding inside the transaction is stale.
func (h *StudioV2WorkspaceDatabaseHandler) commitWorkspaceConnectorSetup(ctx context.Context, boundConnectorID string, req studioV2WorkspaceDatabaseConfigRequest, prepared *dbtarget.PreparedConnector) (*workspace.Record, error) {
	saved := prepared.Connector
	return h.workspaceSvc.UpdateDatabaseSetup(ctx, req.ExpectedSetupRevision, func(ctx context.Context, tx *sql.Tx, current *workspace.Record) error {
		if strings.TrimSpace(current.DatabaseConnectorID) != strings.TrimSpace(boundConnectorID) {
			return workspace.ErrSetupRevisionConflict
		}
		if err := h.connectorSvc.PersistInTx(ctx, tx, prepared); err != nil {
			return err
		}
		current.DatabaseConnectorID = saved.ID
		if req.RowGroups == nil {
			return nil
		}
		return workspace.ReplaceDatabaseRowGroups(current, saved.ID,
			connectorConfigString(saved, "schema"), connectorConfigString(saved, "table"), req.RowGroups)
	})
}

// commitWorkspaceTargetSetup persists a prepared target mapping and its
// workspace row-group reference together; a failed reference leaves the prior
// mapping untouched.
func (h *StudioV2WorkspaceDatabaseHandler) commitWorkspaceTargetSetup(ctx context.Context, connectorID, pointID, rowGroupID, expectedRevision string, prepared *dbtarget.PreparedTargetMapping) (*workspace.Record, error) {
	return h.workspaceSvc.UpdateDatabaseSetup(ctx, expectedRevision, func(ctx context.Context, tx *sql.Tx, current *workspace.Record) error {
		if strings.TrimSpace(current.DatabaseConnectorID) != connectorID {
			return workspace.ErrSetupRevisionConflict
		}
		if err := h.mappingSvc.PersistInTx(ctx, tx, prepared); err != nil {
			return err
		}
		return workspace.SetDatabaseTargetReference(current, pointID, rowGroupID)
	})
}
