package recordingplan

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSuccessorLedgerReadRejectsCorruptScopeActionAndState(t *testing.T) {
	cases := map[string]string{
		"action":             `UPDATE managed_schema_operations SET action='schema_apply'`,
		"workspace":          `UPDATE managed_schema_operations SET workspace_id='W9'`,
		"token":              `UPDATE managed_schema_operations SET token='foreign-token'`,
		"scope":              `UPDATE managed_schema_operations SET scope_key='foreign-scope'`,
		"detail-predecessor": `UPDATE managed_schema_operations SET detail=replace(detail,'"predecessor_group_id":"G1"','"predecessor_group_id":"G9"')`,
		"detail-successor":   `UPDATE managed_schema_operations SET detail=replace(detail,'"successor_group_id":"G2"','"successor_group_id":"G9"')`,
		"digest":             `UPDATE managed_schema_operations SET payload_digest='bad'`,
		"unknown":            `UPDATE managed_schema_operations SET status='unknown'`,
		"owner":              `UPDATE managed_schema_operations SET owner='foreign-owner'`,
		"verified-proof":     `UPDATE managed_schema_operations SET verified_digest='pretend-schema-proof'`,
		"executed":           `UPDATE managed_schema_operations SET executed_statements=1`,
		"incomplete":         `UPDATE managed_schema_operations SET completed_at=NULL`,
		"malformed":          `UPDATE managed_schema_operations SET detail='{'`,
		"unknown-version":    `UPDATE managed_schema_operations SET detail=replace(detail,'successor-operation/v1','successor-operation/v9')`,
		"extra-field":        `UPDATE managed_schema_operations SET detail=substr(detail,1,length(detail)-1)||',"extra":true}'`,
	}
	for name, query := range cases {
		t.Run(name, func(t *testing.T) {
			db := recordingRepositoryDB(t)
			tx := beginSuccessorTx(t, db)
			_, err := RecordSuccessorReserveInTx(t.Context(), tx, successorReserveFixture())
			require.NoError(t, err)
			_, err = tx.Exec(query)
			require.NoError(t, err)
			got, err := ReadSuccessorReserveInTx(t.Context(), tx, "W1", "O2")
			require.Error(t, err)
			require.Nil(t, got)
		})
	}
}

func TestSuccessorLedgerCannotBeAcceptedAsSchemaProof(t *testing.T) {
	db := recordingRepositoryDB(t)
	repo := NewSQLRepository(db)
	svc := NewService(repo)
	token := preparedPreview(t, svc, previewScope())
	op, _, err := svc.ClaimSchemaApply(t.Context(), token)
	require.NoError(t, err)
	digest := strings.Repeat("c", 64)
	_, err = repo.FinishSchemaOperation(t.Context(), op.OperationID, op.Owner, SchemaOperationResult{Status: SchemaOperationSucceeded, VerifiedDigest: digest})
	require.NoError(t, err)
	tx := beginSuccessorTx(t, db)
	for _, action := range []string{SuccessorReserveAction, SuccessorFirstStartAction} {
		_, err = tx.Exec(`UPDATE managed_schema_operations SET action=? WHERE operation_id=?`, action, op.OperationID)
		require.NoError(t, err)
		_, err = SucceededSchemaTokenInTx(t.Context(), tx, token.WorkspaceID, op.OperationID, digest)
		require.ErrorIs(t, err, ErrSchemaOperationResult, "even a matching schema token must not upgrade a successor action")
	}
}

func TestSuccessorLedgerInvalidInputAndForeignRead(t *testing.T) {
	db := recordingRepositoryDB(t)
	tx := beginSuccessorTx(t, db)
	_, err := RecordSuccessorReserveInTx(t.Context(), nil, successorReserveFixture())
	require.ErrorIs(t, err, ErrSchemaOperationResult)
	_, err = ReadSuccessorReserveInTx(t.Context(), nil, "W1", "O2")
	require.ErrorIs(t, err, ErrSchemaOperationResult)
	_, err = RecordSuccessorReserveInTx(t.Context(), tx, nil)
	require.ErrorIs(t, err, ErrSchemaOperationResult)
	for _, change := range []func(*SuccessorReserveRecord){
		func(r *SuccessorReserveRecord) { r.Scope.RequestID = "" },
		func(r *SuccessorReserveRecord) { r.Scope.RequestID = strings.Repeat("a", 129) },
		func(r *SuccessorReserveRecord) { r.Scope.WorkspaceID = " W1" },
		func(r *SuccessorReserveRecord) { r.Scope.SuccessorGroupID = r.Scope.PredecessorGroupID },
		func(r *SuccessorReserveRecord) { r.IntentDigest = strings.Repeat("z", 64) },
		func(r *SuccessorReserveRecord) { r.SourceEvidence = json.RawMessage(`null`) },
	} {
		record := successorReserveFixture()
		change(record)
		_, err = RecordSuccessorReserveInTx(t.Context(), tx, record)
		require.ErrorIs(t, err, ErrSchemaOperationResult)
	}
	_, err = RecordSuccessorReserveInTx(t.Context(), tx, successorReserveFixture())
	require.NoError(t, err)
	got, err := ReadSuccessorReserveInTx(t.Context(), tx, "W9", "O2")
	require.NoError(t, err)
	require.Nil(t, got)
}
