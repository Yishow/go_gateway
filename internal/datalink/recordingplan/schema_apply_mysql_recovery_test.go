package recordingplan

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"path/filepath"
	"slices"
	"testing"
	"time"

	_ "modernc.org/sqlite"

	"github.com/stretchr/testify/require"
)

type mysqlApplyProgressTestDetail struct {
	Version     int                          `json:"version"`
	TokenDigest string                       `json:"token_digest"`
	ScopeKey    string                       `json:"scope_key"`
	Steps       []mysqlApplyProgressTestStep `json:"steps"`
}

type mysqlApplyProgressTestStep struct {
	Kind        string `json:"kind"`
	Name        string `json:"name"`
	SQLHash     string `json:"sql_hash"`
	Status      string `json:"status"`
	Fingerprint string `json:"fingerprint"`
}

type mysqlDDLApplyFixture struct {
	layout        GroupSchemaLayout
	token         *SchemaPreviewToken
	statements    []string
	groupExact    TargetTableInspection
	receiptExact  TargetTableInspection
	states        map[string]TargetTableInspection
	calls         [][]string
	failStatement string
	failErr       error
	failureUsed   bool
	loseAck       bool
}

func newMySQLDDLApplyFixture(t *testing.T, repo Repository) *mysqlDDLApplyFixture {
	t.Helper()
	layout, _, fullGroup, statements := mysqlCatalogFixture(t)
	fullGroup.Status = inspectionExists
	fullGroup.IndexesInspected = true
	fullGroup.MySQL.Database = "recordings"
	receipt := mysqlReceiptStateForApply(fullGroup.MySQL)

	missingGroup := TargetTableInspection{
		Status: inspectionMissing,
		MySQL:  cloneMySQLTableEvidence(fullGroup.MySQL),
	}
	missingReceipt := TargetTableInspection{
		Status: inspectionMissing,
		MySQL:  cloneMySQLTableEvidence(fullGroup.MySQL),
	}
	fixture := &mysqlDDLApplyFixture{
		layout:       layout,
		statements:   statements,
		groupExact:   fullGroup,
		receiptExact: receipt,
		states: map[string]TargetTableInspection{
			layout.TableName:          missingGroup,
			managedEffectReceiptTable: missingReceipt,
		},
	}
	fixture.token = mysqlApplyPreviewToken(t, repo, layout, statements)
	return fixture
}

func mysqlApplyPreviewToken(t *testing.T, repo Repository, layout GroupSchemaLayout, statements []string) *SchemaPreviewToken {
	t.Helper()
	token := mysqlTestToken(layout)
	token.TablePrefix = ""
	token.Database, token.Schema = "recordings", "recordings"
	token.Statements = append([]string(nil), statements...)
	token.Tables = []SchemaPreviewTable{
		{Name: layout.TableName, Action: SchemaTableActionCreate, Columns: groupColumnNames(layout.Columns)},
		{Name: managedEffectReceiptTable, Action: SchemaTableActionCreate,
			Columns: []string{"effect_key", "payload_digest", "committed_at", managedEffectReceiptMarker}},
	}
	now := time.Now().UTC()
	token.CreatedAt, token.ExpiresAt = now, now.Add(time.Hour)
	var err error
	token.Digest, err = previewDigest(&token)
	require.NoError(t, err)
	require.NoError(t, repo.SavePreviewToken(t.Context(), &token))
	return &token
}

func (f *mysqlDDLApplyFixture) target() SchemaApplyTarget {
	return SchemaApplyTarget{Inspect: f.inspect, Execute: f.execute}
}

func (f *mysqlDDLApplyFixture) inspect(_ context.Context, table string) (TargetTableInspection, error) {
	state, ok := f.states[table]
	if !ok {
		return TargetTableInspection{}, fmt.Errorf("unexpected catalog table %q", table)
	}
	return state, nil
}

func (f *mysqlDDLApplyFixture) execute(_ context.Context, statements []string) (SchemaExecution, error) {
	f.calls = append(f.calls, append([]string(nil), statements...))
	if len(statements) != 1 {
		return SchemaExecution{}, fmt.Errorf("MySQL execution must receive one statement, got %d", len(statements))
	}
	statement := statements[0]
	if f.failStatement == statement && !f.failureUsed {
		f.failureUsed = true
		return SchemaExecution{}, f.failErr
	}
	switch statement {
	case f.statements[0]:
		group := cloneTargetTableInspection(f.groupExact)
		group.Indexes = nil
		f.states[f.layout.TableName] = group
	case f.statements[len(f.statements)-1]:
		f.states[managedEffectReceiptTable] = cloneTargetTableInspection(f.receiptExact)
	default:
		if len(f.statements) < 2 {
			return SchemaExecution{}, fmt.Errorf("unexpected MySQL DDL statement %q", statement)
		}
		indexStatements := f.statements[1 : len(f.statements)-1]
		indexPosition := slices.Index(indexStatements, statement)
		if indexPosition < 0 || indexPosition >= len(f.layout.Indexes) {
			return SchemaExecution{}, fmt.Errorf("unexpected MySQL DDL statement %q", statement)
		}
		index := f.layout.Indexes[indexPosition]
		group := f.states[f.layout.TableName]
		group.Indexes = append(group.Indexes, TargetIndexInspection{Name: index.Name, Columns: append([]string(nil), index.Columns...)})
		f.states[f.layout.TableName] = group
	}
	if f.loseAck {
		return SchemaExecution{}, errors.New("connection lost after DDL")
	}
	return SchemaExecution{Committed: 1}, nil
}

func mysqlReceiptStateForApply(binding *TargetMySQLTableEvidence) TargetTableInspection {
	state := TargetTableInspection{
		Status: inspectionExists, IndexesInspected: true,
		ColumnTypes: map[string]string{
			"effect_key": "VARCHAR(191)", "payload_digest": "CHAR(64)",
			"committed_at": "DATETIME(6)", managedEffectReceiptMarker: "VARCHAR(64)",
		},
		ColumnNullable:   map[string]bool{"effect_key": false, "payload_digest": false, "committed_at": false, managedEffectReceiptMarker: false},
		ColumnPrimaryKey: map[string]bool{"effect_key": true, "payload_digest": false, "committed_at": false, managedEffectReceiptMarker: false},
		Columns:          []string{"effect_key", "payload_digest", "committed_at", managedEffectReceiptMarker},
		MySQL:            cloneMySQLTableEvidence(binding),
	}
	state.MySQL.PrimaryKey = []string{"effect_key"}
	state.MySQL.ColumnCollations["effect_key"] = mysqlTestCollation
	state.MySQL.ColumnCollations["payload_digest"] = mysqlTestCollation
	state.MySQL.ColumnCollations[managedEffectReceiptMarker] = mysqlTestCollation
	ownerDefault := "gateway"
	state.MySQL.ColumnDefaults[managedEffectReceiptMarker] = &ownerDefault
	return state
}

func mysqlApplyTestDetail(t *testing.T, detail string) mysqlApplyProgressTestDetail {
	t.Helper()
	var got mysqlApplyProgressTestDetail
	require.NoError(t, json.Unmarshal([]byte(detail), &got))
	return got
}

func mysqlApplySQLHash(statement string) string {
	sum := sha256.Sum256([]byte(statement))
	return fmt.Sprintf("%x", sum[:])
}

func TestMySQLApplyRequiresVerifiedFamilyBaseline(t *testing.T) {
	repo := NewMemoryRepository()
	fixture := newMySQLDDLApplyFixture(t, repo)
	layout := fixture.layout
	options := *layout.MySQL
	options.ServerVersion = "8.3.12"
	layout.MySQL = &options
	unverified := *fixture.token
	unverified.Token, unverified.OperationID = "tok-mysql-unsupported", "op-mysql-unsupported"
	unverified.GroupLayout = &layout
	var err error
	unverified.Digest, err = previewDigest(&unverified)
	require.NoError(t, err)
	for table, state := range fixture.states {
		state.MySQL.ServerVersion = "8.3.12"
		fixture.states[table] = state
	}
	require.NoError(t, repo.SavePreviewToken(t.Context(), &unverified))
	service := NewService(repo)

	_, _, err = service.ApplySchemaPreview(t.Context(), &unverified, fixture.target())

	require.ErrorIs(t, err, ErrSchemaAdapterUnverified)
	require.Empty(t, fixture.calls)
	_, err = repo.GetSchemaOperation(t.Context(), unverified.WorkspaceID, unverified.OperationID)
	require.ErrorIs(t, err, ErrSchemaOperationNotFound)
}

func TestMySQLPartialDDLRecovery(t *testing.T) {
	path := filepath.Join(t.TempDir(), "mysql-apply-ledger.db")
	ledger := openLedgerDB(t, path)
	applyRecordingDDL(t, ledger)
	repo := NewSQLRepository(ledger)
	fixture := newMySQLDDLApplyFixture(t, repo)
	fixture.failStatement = fixture.statements[1]
	fixture.failErr = ErrTargetPermissionDenied
	service := NewService(repo)
	target := fixture.target()

	first, outcome, err := service.ApplySchemaPreview(t.Context(), fixture.token, target)
	require.NoError(t, err)
	require.Equal(t, ClaimAcquired, outcome)
	require.Equal(t, fixture.token.OperationID, first.OperationID)
	require.Equal(t, SchemaOperationPartial, first.Status)
	require.Len(t, fixture.calls, 2)
	for _, call := range fixture.calls {
		require.Len(t, call, 1, "each MySQL DDL call must contain exactly one statement")
	}
	if ManagedSchemaExecutionVerified(dialectMySQL) {
		t.Fatal("the specialized group path must not enable generic MySQL schema execution")
	}
	require.False(t, service.GetConnectorCapability(dialectMySQL).SupportsManagedSchema)

	stored, err := repo.GetSchemaOperation(t.Context(), fixture.token.WorkspaceID, fixture.token.OperationID)
	require.NoError(t, err)
	firstOwner := stored.Owner
	require.NotEmpty(t, firstOwner)
	manifest := mysqlApplyTestDetail(t, stored.Detail)
	require.Equal(t, 1, manifest.Version)
	require.Equal(t, fixture.token.Digest, manifest.TokenDigest)
	require.Equal(t, fixture.token.ScopeKey(), manifest.ScopeKey)
	require.Len(t, manifest.Steps, 3)
	require.Equal(t, []string{"verified", "failed", "pending"}, []string{
		manifest.Steps[0].Status, manifest.Steps[1].Status, manifest.Steps[2].Status,
	})
	for index, statement := range fixture.statements {
		require.Equal(t, mysqlApplySQLHash(statement), manifest.Steps[index].SQLHash)
	}
	groupEvidence, err := inspectMySQLApplyStep(t.Context(), fixture.token, fixture.layout,
		mysqlStepReference(mysqlStepKindGroupTable, fixture.layout.TableName, fixture.statements[0]), target.Inspect)
	require.NoError(t, err)
	require.Equal(t, mysqlStepStateExact, groupEvidence.State)
	indexEvidence, err := inspectMySQLApplyStep(t.Context(), fixture.token, fixture.layout,
		mysqlStepReference(mysqlStepKindIndex, fixture.layout.Indexes[0].Name, fixture.statements[1]), target.Inspect)
	require.NoError(t, err)
	require.Equal(t, mysqlStepStateMissing, indexEvidence.State)
	receiptEvidence, err := inspectMySQLApplyStep(t.Context(), fixture.token, fixture.layout,
		mysqlStepReference(mysqlStepKindReceiptTable, managedEffectReceiptTable, fixture.statements[2]), target.Inspect)
	require.NoError(t, err)
	require.Equal(t, mysqlStepStateMissing, receiptEvidence.State)

	// A changed target scope cannot reclaim the saved partial operation.
	drifted := *fixture.token
	drifted.Database, drifted.Schema = "other", "other"
	_, _, err = service.ApplySchemaPreview(t.Context(), &drifted, target)
	require.ErrorIs(t, err, ErrPreviewTokenStale)
	require.Len(t, fixture.calls, 2)

	require.NoError(t, ledger.Close())
	reopened := openLedgerDB(t, path)
	restartedRepo := NewSQLRepository(reopened)
	restored, err := restartedRepo.GetSchemaOperation(t.Context(), fixture.token.WorkspaceID, fixture.token.OperationID)
	require.NoError(t, err)
	require.Equal(t, SchemaOperationPartial, restored.Status)
	require.Equal(t, stored.Detail, restored.Detail)
	restarted := NewService(restartedRepo)

	fixture.failStatement = ""
	fixture.failErr = nil
	resumed, outcome, err := restarted.ApplySchemaPreview(t.Context(), fixture.token, fixture.target())
	require.NoError(t, err)
	require.Equal(t, ClaimAcquired, outcome)
	require.Equal(t, fixture.token.OperationID, resumed.OperationID)
	require.Equal(t, SchemaOperationSucceeded, resumed.Status)
	require.Len(t, fixture.calls, 4)
	require.Equal(t, fixture.statements[0], fixture.calls[0][0], "the completed table must not be rebuilt")
	require.Equal(t, fixture.statements[1], fixture.calls[1][0])
	require.Equal(t, fixture.statements[1], fixture.calls[2][0], "only the missing index is retried")
	require.Equal(t, fixture.statements[2], fixture.calls[3][0])

	final, err := restartedRepo.GetSchemaOperation(t.Context(), fixture.token.WorkspaceID, fixture.token.OperationID)
	require.NoError(t, err)
	require.NotEmpty(t, final.Owner)
	require.NotEqual(t, firstOwner, final.Owner, "reclaim must fence the previous owner")
	finalManifest := mysqlApplyTestDetail(t, final.Detail)
	for _, step := range finalManifest.Steps {
		require.Equal(t, "verified", step.Status)
	}
}

func TestMySQLDDLResponseLost(t *testing.T) {
	repo := NewMemoryRepository()
	fixture := newMySQLDDLApplyFixture(t, repo)
	fixture.loseAck = true
	service := NewService(repo)

	operation, outcome, err := service.ApplySchemaPreview(t.Context(), fixture.token, fixture.target())

	require.NoError(t, err)
	require.Equal(t, ClaimAcquired, outcome)
	require.Equal(t, fixture.token.OperationID, operation.OperationID)
	require.Equal(t, SchemaOperationSucceeded, operation.Status)
	require.Len(t, fixture.calls, 3)
	for index, call := range fixture.calls {
		require.Len(t, call, 1)
		require.Equal(t, fixture.statements[index], call[0])
	}
	stored, err := repo.GetSchemaOperation(t.Context(), fixture.token.WorkspaceID, fixture.token.OperationID)
	require.NoError(t, err)
	manifest := mysqlApplyTestDetail(t, stored.Detail)
	for _, step := range manifest.Steps {
		require.Equal(t, mysqlSchemaStepVerified, step.Status)
		require.NotEmpty(t, step.Fingerprint)
	}
}

type mysqlApplyFinishOnceRepository struct {
	*SQLRepository
	fail bool
}

func (r *mysqlApplyFinishOnceRepository) FinishSchemaOperation(ctx context.Context, operationID, owner string, result SchemaOperationResult) (*SchemaOperation, error) {
	if r.fail {
		r.fail = false
		return nil, errors.New("local operation acknowledgement lost")
	}
	return r.SQLRepository.FinishSchemaOperation(ctx, operationID, owner, result)
}

func TestMySQLStaleRunningOperationReconcilesSamePreview(t *testing.T) {
	path := filepath.Join(t.TempDir(), "mysql-stale-ledger.db")
	ledger := openMySQLApplyLedger(t, path)
	applyRecordingDDL(t, ledger)
	baseRepo := NewSQLRepository(ledger)
	repo := &mysqlApplyFinishOnceRepository{SQLRepository: baseRepo, fail: true}
	fixture := newMySQLDDLApplyFixture(t, repo)
	service := NewService(repo)

	first, outcome, err := service.ApplySchemaPreview(t.Context(), fixture.token, fixture.target())
	require.ErrorIs(t, err, ErrSchemaOperationUnacknowledged)
	require.Equal(t, ClaimAcquired, outcome)
	require.Equal(t, fixture.token.OperationID, first.OperationID)
	require.Len(t, fixture.calls, 3)
	beforeRestart, err := baseRepo.GetSchemaOperation(t.Context(), fixture.token.WorkspaceID, fixture.token.OperationID)
	require.NoError(t, err)
	oldOwner := beforeRestart.Owner
	require.NotEmpty(t, oldOwner)
	for _, step := range mysqlApplyTestDetail(t, beforeRestart.Detail).Steps {
		require.Equal(t, mysqlSchemaStepVerified, step.Status)
	}

	require.NoError(t, ledger.Close())
	reopened := openMySQLApplyLedger(t, path)
	_, err = reopened.ExecContext(t.Context(), `UPDATE managed_schema_operations SET updated_at = ? WHERE operation_id = ?`,
		time.Now().UTC().Add(-2*SchemaOperationLease), fixture.token.OperationID)
	require.NoError(t, err)
	restartedRepo := NewSQLRepository(reopened)
	restarted := NewService(restartedRepo)

	resolved, outcome, err := restarted.ApplySchemaPreview(t.Context(), fixture.token, fixture.target())

	require.NoError(t, err)
	require.Equal(t, ClaimAcquired, outcome)
	require.Equal(t, fixture.token.OperationID, resolved.OperationID)
	require.Equal(t, SchemaOperationSucceeded, resolved.Status)
	require.Len(t, fixture.calls, 3, "stale reconciliation must never re-execute confirmed DDL")
	afterRestart, err := restartedRepo.GetSchemaOperation(t.Context(), fixture.token.WorkspaceID, fixture.token.OperationID)
	require.NoError(t, err)
	require.NotEqual(t, oldOwner, afterRestart.Owner)
	require.Equal(t, beforeRestart.Detail, afterRestart.Detail)
}

func openMySQLApplyLedger(t *testing.T, path string) *sql.DB {
	t.Helper()
	db, err := sql.Open("sqlite", "file:"+path+"?_pragma=busy_timeout(10000)&_pragma=journal_mode(WAL)")
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })
	return db
}
