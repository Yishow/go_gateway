package workspace

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"testing"

	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/recordingplan"
	"go-gateway/internal/datalink/schema"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

type fakeMySQLProbeInspector struct {
	capability *dbtarget.MySQLRecordingCapability
	err        error
}

func (f *fakeMySQLProbeInspector) InspectTableAtRevision(ctx context.Context, connectorID, expectedRevision, schemaName, tableName string) (*dbtarget.TableInspection, error) {
	return &dbtarget.TableInspection{Status: "missing"}, nil
}

func (f *fakeMySQLProbeInspector) ProbeMySQLRecordingAtRevision(ctx context.Context, connectorID, expectedRevision string) (*dbtarget.MySQLRecordingCapability, error) {
	if f.err != nil {
		return nil, f.err
	}
	return f.capability, nil
}

func seedMySQLConnector(t *testing.T, db *sql.DB, configJSON string) string {
	t.Helper()
	id := uuid.NewString()
	name := "MySQL Target " + id[:8]
	_, err := db.ExecContext(t.Context(), `
		INSERT INTO database_connectors (id, name, kind, connection_config, identity_revision, status, enabled)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`, id, name, "mysql", configJSON, "connector-revision-mysql", schema.DatabaseConnectorStatusReady, true)
	require.NoError(t, err)
	return id
}

func ensureMySQLGroup(t *testing.T, h *circuitProfileHarness, group *WriteGroup) (*WriteGroupSaveResult, error) {
	t.Helper()
	current, err := h.workspace.GetOrCreate(t.Context())
	require.NoError(t, err)
	mutation := WriteGroupMutation{
		WorkspaceID:               h.fixture.workspaceID,
		ExpectedWorkspaceRevision: current.DatabaseSetupRevision,
		ExpectedConnectorRevision: group.Destination.ConnectorRevision,
		Group:                     group,
	}
	return h.service.EnsureCircuitProfile(t.Context(), mutation)
}

func mustEnsureMySQLGroup(t *testing.T, h *circuitProfileHarness, group *WriteGroup) *WriteGroup {
	t.Helper()
	result, err := ensureMySQLGroup(t, h, group)
	require.NoError(t, err)
	require.NotNil(t, result)
	require.NotNil(t, result.Group)
	return result.Group
}

func TestMySQLWriteGroupCreation_ScopeResolutionAndDraftStorage(t *testing.T) {
	h := newCircuitProfileHarness(t, ":memory:")
	connectorID := seedMySQLConnector(t, h.db, `{"dsn":"app:secret@tcp(127.0.0.1:3306)/telemetry_db?tls=true"}`)
	power := h.primary()

	group := h.longDraft(WriteGroupProfileInstantaneous, "mysql-main", power)
	group.Destination.ConnectorID = connectorID
	group.Destination.ConnectorRevision = "connector-revision-mysql"
	group.Destination.Database = ""
	group.Destination.TableSchema = ""

	saved := mustEnsureMySQLGroup(t, h, group)
	require.Equal(t, WriteGroupLayoutManagedLongV1, saved.EffectiveLayoutVersion())
	require.Equal(t, WriteGroupStorageModePlain, saved.StorageMode)
	require.Equal(t, "telemetry_db", saved.Destination.Database)
	require.Equal(t, "telemetry_db", saved.Destination.TableSchema)

	scope := longStorageScopeOf(saved)
	require.Equal(t, scope.tableName(), saved.Destination.TableName)

	var tableName, ownerColumn, storedDigest, confirmedLayout string
	var generation int
	err := h.db.QueryRowContext(t.Context(), `
		SELECT table_name, owner_column, layout_digest, confirmed_layout, generation
		FROM write_group_storage_scopes
		WHERE workspace_id = ? AND connector_id = ? AND database_name = ? AND table_schema = ?
			AND storage_mode = ? AND layout_version = ?
	`, scope.WorkspaceID, scope.ConnectorID, scope.Database, scope.TableSchema, string(scope.Mode), string(scope.Layout),
	).Scan(&tableName, &ownerColumn, &storedDigest, &confirmedLayout, &generation)
	require.NoError(t, err)
	require.Equal(t, scope.tableName(), tableName)
	require.Equal(t, scope.ownerColumn(), ownerColumn)
	require.Equal(t, 1, generation)
	require.Empty(t, confirmedLayout, "draft storage scope must not have confirmed_layout")
	expectedSeed := recordingplan.UnconfirmedMySQLStorageDigest(StorageScopeOf(saved))
	require.Equal(t, expectedSeed, storedDigest, "unconfirmed draft digest must match versioned scope seed")
	require.Equal(t, scope.digest("mysql-storage-unconfirmed-v1"), expectedSeed, "seed must match scope digest")
}

func TestMySQLWriteGroupCreation_Rejections(t *testing.T) {
	h := newCircuitProfileHarness(t, ":memory:")
	connectorID := seedMySQLConnector(t, h.db, `{"dsn":"app:secret@tcp(127.0.0.1:3306)/telemetry_db?tls=true"}`)
	emptyDBConnID := seedMySQLConnector(t, h.db, `{"host":"127.0.0.1","user":"app","password":"sec"}`)
	power := h.primary()

	t.Run("empty database in connector config", func(t *testing.T) {
		group := h.longDraft(WriteGroupProfileInstantaneous, "empty-db", power)
		group.Destination.ConnectorID = emptyDBConnID
		group.Destination.ConnectorRevision = "connector-revision-mysql"
		_, err := ensureMySQLGroup(t, h, group)
		require.Error(t, err)
	})

	t.Run("conflicting destination database rejected", func(t *testing.T) {
		group := h.longDraft(WriteGroupProfileInstantaneous, "conflict-db", power)
		group.Destination.ConnectorID = connectorID
		group.Destination.ConnectorRevision = "connector-revision-mysql"
		group.Destination.Database = "other_database"
		_, err := ensureMySQLGroup(t, h, group)
		require.Error(t, err)
		require.ErrorIs(t, err, ErrWriteGroupValidation)
	})

	t.Run("conflicting destination schema rejected", func(t *testing.T) {
		group := h.longDraft(WriteGroupProfileInstantaneous, "conflict-schema", power)
		group.Destination.ConnectorID = connectorID
		group.Destination.ConnectorRevision = "connector-revision-mysql"
		group.Destination.TableSchema = "other_schema"
		_, err := ensureMySQLGroup(t, h, group)
		require.Error(t, err)
		require.ErrorIs(t, err, ErrWriteGroupValidation)
	})

	t.Run("non-plain storage mode rejected", func(t *testing.T) {
		group := h.longDraft(WriteGroupProfileInstantaneous, "timescale-mysql", power)
		group.Destination.ConnectorID = connectorID
		group.Destination.ConnectorRevision = "connector-revision-mysql"
		group.StorageMode = WriteGroupStorageModeTimescale
		_, err := ensureMySQLGroup(t, h, group)
		require.Error(t, err)
		require.ErrorIs(t, err, ErrWriteGroupValidation)
	})

	t.Run("wide layout with mysql rejected", func(t *testing.T) {
		group := h.draft(WriteGroupProfileInstantaneous, "wide-mysql", "Wide MySQL", power)
		group.LayoutVersion = WriteGroupLayoutManagedWideProfileV1
		group.Destination.ConnectorID = connectorID
		group.Destination.ConnectorRevision = "connector-revision-mysql"
		_, err := ensureMySQLGroup(t, h, group)
		require.Error(t, err)
		require.ErrorIs(t, err, ErrWriteGroupValidation)
	})
}

func TestMySQLManagedSchemaScope_BackendFactsBinding(t *testing.T) {
	h := newCircuitProfileHarness(t, ":memory:")
	connectorID := seedMySQLConnector(t, h.db, `{"dsn":"app:secret@tcp(127.0.0.1:3306)/metrics_db?tls=true"}`)
	power := h.primary()

	group := h.longDraft(WriteGroupProfileInstantaneous, "mysql-facts", power)
	group.Destination.ConnectorID = connectorID
	group.Destination.ConnectorRevision = "connector-revision-mysql"
	saved := mustEnsureMySQLGroup(t, h, group)

	inspector := &fakeMySQLProbeInspector{
		capability: &dbtarget.MySQLRecordingCapability{
			ConnectorID:         connectorID,
			ConnectorRevision:   "connector-revision-mysql",
			Database:            "metrics_db",
			ServerFamily:        "mysql",
			ServerVersion:       "8.4.13",
			CapabilityRevision:  "cap-mysql-rev1",
			Connected:           true,
			MetadataStatus:      "verified",
			TransactionalEngine: true,
		},
	}
	h.service.WithManagedTableInspector(inspector)

	workspaceRec, err := h.workspace.GetOrCreate(t.Context())
	require.NoError(t, err)

	scope, err := h.service.ManagedSchemaScope(t.Context(), saved.ID, WriteGroupMutation{
		WorkspaceID:               workspaceRec.ID,
		ExpectedWorkspaceRevision: workspaceRec.DatabaseSetupRevision,
		ExpectedGroupRevision:     saved.Revision,
		ExpectedConnectorRevision: saved.Destination.ConnectorRevision,
		Group:                     saved,
	})
	require.NoError(t, err)
	require.Equal(t, "mysql", scope.Dialect)
	require.Equal(t, "metrics_db", scope.Database)
	require.Equal(t, "metrics_db", scope.Schema)
	require.NotNil(t, scope.GroupLayout)
	require.NotNil(t, scope.GroupLayout.MySQL)

	mySQLOpts := scope.GroupLayout.MySQL
	require.Equal(t, "mysql", mySQLOpts.ServerFamily)
	require.Equal(t, "8.4.13", mySQLOpts.ServerVersion)
	require.Equal(t, "cap-mysql-rev1", mySQLOpts.CapabilityRevision)
	require.Equal(t, "InnoDB", mySQLOpts.Engine)
	require.Equal(t, "utf8mb4_0900_bin", mySQLOpts.Collation)

	expectedTable, expectedOwner, ok := LongStorageIdentity(saved)
	require.True(t, ok)
	require.Equal(t, expectedTable, scope.GroupLayout.TableName)
	require.Equal(t, expectedOwner, scope.GroupLayout.OwnerColumn)

	// MariaDB verified probe identity
	mariaInspector := &fakeMySQLProbeInspector{
		capability: &dbtarget.MySQLRecordingCapability{
			ConnectorID:         connectorID,
			ConnectorRevision:   "connector-revision-mysql",
			Database:            "metrics_db",
			ServerFamily:        "mariadb",
			ServerVersion:       "11.4.13-MariaDB-ubu2404",
			CapabilityRevision:  "cap-maria-rev1",
			Connected:           true,
			MetadataStatus:      "verified",
			TransactionalEngine: true,
		},
	}
	h.service.WithManagedTableInspector(mariaInspector)

	mariaScope, err := h.service.ManagedSchemaScope(t.Context(), saved.ID, WriteGroupMutation{
		WorkspaceID:               workspaceRec.ID,
		ExpectedWorkspaceRevision: workspaceRec.DatabaseSetupRevision,
		ExpectedGroupRevision:     saved.Revision,
		ExpectedConnectorRevision: saved.Destination.ConnectorRevision,
		Group:                     saved,
	})
	require.NoError(t, err)
	require.Equal(t, "mariadb", mariaScope.GroupLayout.MySQL.ServerFamily)
	require.Equal(t, "utf8mb4_nopad_bin", mariaScope.GroupLayout.MySQL.Collation)
}

func TestMySQLManagedSchemaScope_ClientCannotSelfClaimFamily(t *testing.T) {
	h := newCircuitProfileHarness(t, ":memory:")
	connectorID := seedMySQLConnector(t, h.db, `{"dsn":"app:secret@tcp(127.0.0.1:3306)/metrics_db?tls=true"}`)
	power := h.primary()

	group := h.longDraft(WriteGroupProfileInstantaneous, "mysql-unverified", power)
	group.Destination.ConnectorID = connectorID
	group.Destination.ConnectorRevision = "connector-revision-mysql"
	saved := mustEnsureMySQLGroup(t, h, group)

	workspaceRec, err := h.workspace.GetOrCreate(t.Context())
	require.NoError(t, err)

	mutation := WriteGroupMutation{
		WorkspaceID:               workspaceRec.ID,
		ExpectedWorkspaceRevision: workspaceRec.DatabaseSetupRevision,
		ExpectedGroupRevision:     saved.Revision,
		ExpectedConnectorRevision: saved.Destination.ConnectorRevision,
		Group:                     saved,
	}

	for name, mutate := range map[string]func(*dbtarget.MySQLRecordingCapability){
		"unsupported version 8.0":       func(c *dbtarget.MySQLRecordingCapability) { c.ServerVersion = "8.0.35" },
		"unsupported version 8.40.1":    func(c *dbtarget.MySQLRecordingCapability) { c.ServerVersion = "8.40.1" },
		"unverified mysql patch 8.4.4":  func(c *dbtarget.MySQLRecordingCapability) { c.ServerVersion = "8.4.4" },
		"unverified mysql patch 8.4.14": func(c *dbtarget.MySQLRecordingCapability) { c.ServerVersion = "8.4.14" },
		"unsupported version 11.40.1":   func(c *dbtarget.MySQLRecordingCapability) { c.ServerFamily = "mariadb"; c.ServerVersion = "11.40.1" },
		"unverified mariadb distro": func(c *dbtarget.MySQLRecordingCapability) {
			c.ServerFamily = "mariadb"
			c.ServerVersion = "11.4.13-MariaDB-otherdistro"
		},
		"unsupported family tidb": func(c *dbtarget.MySQLRecordingCapability) { c.ServerFamily = "tidb" },
		"disconnected":            func(c *dbtarget.MySQLRecordingCapability) { c.Connected = false },
		"unverified metadata":     func(c *dbtarget.MySQLRecordingCapability) { c.MetadataStatus = "forbidden" },
		"non transactional":       func(c *dbtarget.MySQLRecordingCapability) { c.TransactionalEngine = false },
		"mismatched database":     func(c *dbtarget.MySQLRecordingCapability) { c.Database = "wrong_db" },
		"mismatched connector":    func(c *dbtarget.MySQLRecordingCapability) { c.ConnectorID = "wrong_conn" },
		"mismatched revision":     func(c *dbtarget.MySQLRecordingCapability) { c.ConnectorRevision = "wrong_rev" },
	} {
		t.Run(name, func(t *testing.T) {
			targetCap := &dbtarget.MySQLRecordingCapability{
				ConnectorID:         connectorID,
				ConnectorRevision:   "connector-revision-mysql",
				Database:            "metrics_db",
				ServerFamily:        "mysql",
				ServerVersion:       "8.4.13",
				CapabilityRevision:  "cap-rev1",
				Connected:           true,
				MetadataStatus:      "verified",
				TransactionalEngine: true,
			}
			mutate(targetCap)
			inspector := &fakeMySQLProbeInspector{capability: targetCap}
			h.service.WithManagedTableInspector(inspector)

			_, err := h.service.ManagedSchemaScope(t.Context(), saved.ID, mutation)
			require.Error(t, err)
			require.ErrorIs(t, err, ErrWriteGroupValidation)
		})
	}
}

func TestMySQLStorageScope_ConfirmedLayoutProofCannotBeOverwritten(t *testing.T) {
	h := newCircuitProfileHarness(t, ":memory:")
	connectorID := seedMySQLConnector(t, h.db, `{"dsn":"app:secret@tcp(127.0.0.1:3306)/metrics_db?tls=true"}`)
	power := h.primary()

	group := h.longDraft(WriteGroupProfileInstantaneous, "mysql-confirmed", power)
	group.Destination.ConnectorID = connectorID
	group.Destination.ConnectorRevision = "connector-revision-mysql"
	saved := mustEnsureMySQLGroup(t, h, group)
	scope := longStorageScopeOf(saved)

	layout, err := recordingplan.ManagedMySQLLongGroupLayout(scope.tableName(), scope.ownerColumn(), recordingplan.MySQLGroupSchemaOptions{
		ServerFamily: "mysql", ServerVersion: "8.4.13", CapabilityRevision: "cap-1", Engine: "InnoDB", Collation: "utf8mb4_0900_bin",
	})
	require.NoError(t, err)

	encoded, err := json.Marshal(layout)
	require.NoError(t, err)
	sum := sha256.Sum256(encoded)
	validDigest := hex.EncodeToString(sum[:])

	// Simulate main agent recording storage confirmation
	_, err = h.db.ExecContext(t.Context(), `
		UPDATE write_group_storage_scopes
		SET layout_digest = ?, confirmed_layout = ?, generation = 1
		WHERE workspace_id = ? AND connector_id = ? AND database_name = ? AND table_schema = ?
			AND storage_mode = ? AND layout_version = ?
	`, validDigest, string(encoded), scope.WorkspaceID, scope.ConnectorID, scope.Database, scope.TableSchema, string(scope.Mode), string(scope.Layout))
	require.NoError(t, err)

	// Re-saving with valid confirmed layout succeeds
	current, err := h.service.Get(t.Context(), saved.ID)
	require.NoError(t, err)
	_, err = h.service.Update(t.Context(), saved.ID, WriteGroupMutation{
		WorkspaceID: saved.WorkspaceID, ExpectedWorkspaceRevision: current.WorkspaceRevision,
		ExpectedGroupRevision: saved.Revision, ExpectedConnectorRevision: saved.Destination.ConnectorRevision, Group: saved,
	})
	require.NoError(t, err)

	// Digest remains the confirmed digest, not reverted to unconfirmed seed
	var currentDigest string
	err = h.db.QueryRowContext(t.Context(), `
		SELECT layout_digest FROM write_group_storage_scopes
		WHERE workspace_id = ? AND connector_id = ? AND database_name = ? AND table_schema = ?
			AND storage_mode = ? AND layout_version = ?
	`, scope.WorkspaceID, scope.ConnectorID, scope.Database, scope.TableSchema, string(scope.Mode), string(scope.Layout)).Scan(&currentDigest)
	require.NoError(t, err)
	require.Equal(t, validDigest, currentDigest)
}

func TestMySQLStorageScope_TamperedColumnsInConfirmedLayoutRejected(t *testing.T) {
	h := newCircuitProfileHarness(t, ":memory:")
	connectorID := seedMySQLConnector(t, h.db, `{"dsn":"app:secret@tcp(127.0.0.1:3306)/metrics_db?tls=true"}`)
	power := h.primary()

	group := h.longDraft(WriteGroupProfileInstantaneous, "mysql-tampered", power)
	group.Destination.ConnectorID = connectorID
	group.Destination.ConnectorRevision = "connector-revision-mysql"
	saved := mustEnsureMySQLGroup(t, h, group)
	scope := longStorageScopeOf(saved)

	layout, err := recordingplan.ManagedMySQLLongGroupLayout(scope.tableName(), scope.ownerColumn(), recordingplan.MySQLGroupSchemaOptions{
		ServerFamily: "mysql", ServerVersion: "8.4.13", CapabilityRevision: "cap-1", Engine: "InnoDB", Collation: "utf8mb4_0900_bin",
	})
	require.NoError(t, err)

	// Tamper with columns: only keep 2 columns instead of complete long layout
	tampered := layout
	tampered.Columns = tampered.Columns[:2]

	encoded, err := json.Marshal(tampered)
	require.NoError(t, err)
	sum := sha256.Sum256(encoded)
	matchingDigest := hex.EncodeToString(sum[:])

	// Store tampered confirmed layout with matching hash in write_group_storage_scopes
	_, err = h.db.ExecContext(t.Context(), `
		UPDATE write_group_storage_scopes
		SET layout_digest = ?, confirmed_layout = ?, generation = 1
		WHERE workspace_id = ? AND connector_id = ? AND database_name = ? AND table_schema = ?
			AND storage_mode = ? AND layout_version = ?
	`, matchingDigest, string(encoded), scope.WorkspaceID, scope.ConnectorID, scope.Database, scope.TableSchema, string(scope.Mode), string(scope.Layout))
	require.NoError(t, err)

	// Use an actual update: EnsureCircuitProfile's no-op replay intentionally
	// returns an existing group without opening a persistence transaction.
	// The update must reject malformed proof even with a matching stored hash.
	current, err := h.service.Get(t.Context(), saved.ID)
	require.NoError(t, err)
	_, err = h.service.Update(t.Context(), saved.ID, WriteGroupMutation{
		WorkspaceID: saved.WorkspaceID, ExpectedWorkspaceRevision: current.WorkspaceRevision,
		ExpectedGroupRevision: saved.Revision, ExpectedConnectorRevision: saved.Destination.ConnectorRevision, Group: saved,
	})
	require.Error(t, err)
	require.ErrorIs(t, err, ErrWriteGroupValidation)
}

func TestMySQLWriteGroup_RuntimeAndReadinessRemainClosed(t *testing.T) {
	h := newCircuitProfileHarness(t, ":memory:")
	connectorID := seedMySQLConnector(t, h.db, `{"dsn":"app:secret@tcp(127.0.0.1:3306)/metrics_db?tls=true"}`)
	power := h.primary()

	group := h.longDraft(WriteGroupProfileInstantaneous, "mysql-not-ready", power)
	group.Destination.ConnectorID = connectorID
	group.Destination.ConnectorRevision = "connector-revision-mysql"
	saved := mustEnsureMySQLGroup(t, h, group)

	readiness, err := h.service.Readiness(t.Context(), saved.ID)
	require.NoError(t, err)
	require.False(t, readiness.Ready, "MySQL write group must not be ready for production runtime writing")
	require.NotEmpty(t, readiness.Issues, "must have readiness issues blocking writing")

	for _, issue := range readiness.Issues {
		t.Logf("issue code=%s message=%s", issue.Code, issue.Message)
	}
}
