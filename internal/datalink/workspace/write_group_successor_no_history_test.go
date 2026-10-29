package workspace

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/require"
	"go-gateway/internal/datalink/recordingplan"
	"testing"
	"time"
)

func TestSuccessorNoHistoryRejectsForeignCutoverNamingLineage(t *testing.T) {
	h, id, stage, _ := successorSchemaFixture(t)
	_, err := h.db.Exec(`INSERT INTO write_group_cutovers(id,workspace_id,kind,source_group_id,boundary,review_digest,detail,created_at,intent) VALUES('foreign-cut','foreign-workspace','forward',?,?, 'review',?,CURRENT_TIMESTAMP,'{}')`, id, time.Now().Format(time.RFC3339Nano), fmt.Sprintf(`{"ended":[{"group_id":%q,"revision":"R"}],"started":[{"group_id":"foreign","revision":"R2"}]}`, id))
	require.NoError(t, err)
	tx, err := h.db.BeginTx(t.Context(), &sql.TxOptions{ReadOnly: true})
	require.NoError(t, err)
	defer tx.Rollback()
	err = h.service.repo.requireSuccessorNoHistoryInTx(t.Context(), tx, stage.Envelope.Request.WorkspaceID, stage.Envelope.Request.LineageID, id)
	require.Error(t, err, "a foreign workspace label cannot erase a contradictory cutover naming this lineage")
}

func TestSuccessorNoHistoryRejectsForeignStartNamingLineage(t *testing.T) {
	h, id, stage, _ := successorSchemaFixture(t)
	request, digest, err := normalizeRecordingStart(RecordingStartRequest{WorkspaceID: "foreign-workspace", RequestID: "foreign-start", DeviceIDs: []string{"D1"}, Groups: []RecordingStartGroupIntent{{GroupID: id, ExpectedGroupRevision: "R", ExpectedConnectorRevision: "C"}}})
	require.NoError(t, err)
	raw, err := json.Marshal(initialRecordingStart(request, digest))
	require.NoError(t, err)
	ledger := recordingplan.NewService(recordingplan.NewSQLRepository(h.db))
	_, _, err = ledger.ClaimRecordingStart(t.Context(), "foreign-workspace", "foreign-start", digest, string(raw))
	require.NoError(t, err)
	tx, err := h.db.BeginTx(t.Context(), nil)
	require.NoError(t, err)
	defer tx.Rollback()
	err = h.service.repo.requireSuccessorNoHistoryInTx(t.Context(), tx, stage.Envelope.Request.WorkspaceID, stage.Envelope.Request.LineageID, id)
	require.Error(t, err, "foreign workspace claim naming this lineage contradicts no-history proof")
}
