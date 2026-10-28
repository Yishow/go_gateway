package workspace

import (
	"encoding/json"
	"github.com/stretchr/testify/require"
	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/recordingplan"
	"math"
	"testing"
)

func TestSuccessorReserveDecoderRejectsPrivateProofFamilyContradiction(t *testing.T) {
	h, id, req, _ := successorMySQLNeverAppliedFixture(t, "mysql")
	v, err := h.service.SuccessorReview(t.Context(), id, req)
	require.NoError(t, err)
	_, err = h.service.ReserveSuccessor(t.Context(), id, WriteGroupSuccessorReservationRequest{Envelope: v.Envelope, ReviewDigest: v.ReviewDigest})
	require.NoError(t, err)
	tx, err := h.db.BeginTx(t.Context(), nil)
	require.NoError(t, err)
	defer tx.Rollback()
	root, err := recordingplan.ReadSuccessorReserveInTx(t.Context(), tx, req.WorkspaceID, req.OperationID)
	require.NoError(t, err)
	var binding map[string]any
	require.NoError(t, json.Unmarshal(root.SchemaSourceBinding, &binding))
	binding["source_family"] = "mariadb"
	root.SchemaSourceBinding, err = json.Marshal(binding)
	require.NoError(t, err)
	_, _, err = decodeSuccessorReserveRecord(root)
	require.Error(t, err, "schema source family must equal the immutable dual-probe facts, not merely be valid JSON")
}

func TestSuccessorReserveDecoderRejectsMySQLProofOnNonMySQLRoot(t *testing.T) {
	h, id, req := successorReviewFixture(t)
	v, err := h.service.SuccessorReview(t.Context(), id, req)
	require.NoError(t, err)
	_, err = h.service.ReserveSuccessor(t.Context(), id, WriteGroupSuccessorReservationRequest{Envelope: v.Envelope, ReviewDigest: v.ReviewDigest})
	require.NoError(t, err)
	tx, err := h.db.BeginTx(t.Context(), nil)
	require.NoError(t, err)
	defer tx.Rollback()
	root, err := recordingplan.ReadSuccessorReserveInTx(t.Context(), tx, req.WorkspaceID, req.OperationID)
	require.NoError(t, err)
	var proof successorReserveSourceEvidence
	require.NoError(t, json.Unmarshal(root.SourceEvidence, &proof))
	var reserved WriteGroupSuccessorPreparation
	require.NoError(t, json.Unmarshal(root.Result, &reserved))
	capability := func(g *WriteGroup) *dbtarget.MySQLRecordingCapability {
		return &dbtarget.MySQLRecordingCapability{ConnectorID: g.Destination.ConnectorID, ConnectorRevision: g.Destination.ConnectorRevision, Database: g.Destination.Database, ServerFamily: "mysql", ServerVersion: "8.4.13", CapabilityRevision: "forged-cap", Connected: true, MetadataStatus: "verified", TransactionalEngine: true, SupportsManagedLong: true}
	}
	proof.SourceMySQL = capability(proof.Source)
	proof.TargetMySQL = capability(reserved.Group)
	root.SourceEvidence, err = json.Marshal(proof)
	require.NoError(t, err)
	_, _, err = decodeSuccessorReserveRecord(root)
	require.Error(t, err, "mysql authority cannot be attached to another engine's reserve")
}

func TestSuccessorReserveDecoderRejectsOverflowedGenerationAndLineageRevision(t *testing.T) {
	h, id, req := successorReviewFixture(t)
	v, err := h.service.SuccessorReview(t.Context(), id, req)
	require.NoError(t, err)
	_, err = h.service.ReserveSuccessor(t.Context(), id, WriteGroupSuccessorReservationRequest{Envelope: v.Envelope, ReviewDigest: v.ReviewDigest})
	require.NoError(t, err)
	tx, err := h.db.BeginTx(t.Context(), nil)
	require.NoError(t, err)
	defer tx.Rollback()
	root, err := recordingplan.ReadSuccessorReserveInTx(t.Context(), tx, req.WorkspaceID, req.OperationID)
	require.NoError(t, err)
	var envelope WriteGroupSuccessorEnvelope
	require.NoError(t, json.Unmarshal(root.Intent, &envelope))
	envelope.Generation = math.MaxInt64
	envelope.Request.ExpectedLineageRevision = math.MaxInt64
	envelope, _, digest, err := normalizeSuccessorEnvelope(envelope)
	require.NoError(t, err)
	root.IntentDigest = digest
	root.Intent, err = json.Marshal(envelope)
	require.NoError(t, err)
	var result WriteGroupSuccessorPreparation
	require.NoError(t, json.Unmarshal(root.Result, &result))
	result.Generation = math.MinInt64
	result.LineageRevision = math.MinInt64
	result.IntentDigest = digest
	root.Result, err = json.Marshal(result)
	require.NoError(t, err)
	var binding successorReserveSchemaBinding
	require.NoError(t, json.Unmarshal(root.SchemaSourceBinding, &binding))
	binding.IntentDigest = digest
	binding.SourceDigest, err = binding.sourceDigest()
	require.NoError(t, err)
	root.SchemaSourceBinding, err = json.Marshal(binding)
	require.NoError(t, err)
	_, _, err = decodeSuccessorReserveRecord(root)
	require.Error(t, err, "corrupt max-int origin cannot become a negative successful generation by integer wrap")
}

func TestSuccessorFirstStartDecoderRejectsOverflowedLineageRevision(t *testing.T) {
	h, id, req := firstStartAtomicFixture(t)
	_, err := h.service.ConfirmSuccessorFirstStart(t.Context(), id, req)
	require.NoError(t, err)
	e := req.Envelope.Request.Stage.Envelope
	tx, err := h.db.BeginTx(t.Context(), nil)
	require.NoError(t, err)
	defer tx.Rollback()
	stored, err := recordingplan.ReadSuccessorFirstStartInTx(t.Context(), tx, e.Request.WorkspaceID, recordingplan.SuccessorReserveOperationID(e.Request.WorkspaceID, e.Request.OperationID))
	require.NoError(t, err)
	req.Envelope.Request.Stage.Current.LineageRevision = math.MaxInt64
	req.ReviewDigest, err = successorActivationDigest(req.Envelope)
	require.NoError(t, err)
	_, stored.Intent, err = successorFirstStartIntent(id, req)
	require.NoError(t, err)
	stored.IntentDigest = req.ReviewDigest
	var result WriteGroupSuccessorFirstStartResult
	require.NoError(t, json.Unmarshal(stored.Result, &result))
	result.LineageRevision = math.MinInt64
	stored.Result, err = json.Marshal(result)
	require.NoError(t, err)
	_, _, err = decodeSuccessorFirstStartRecord(stored)
	require.Error(t, err, "an overflowed lineage CAS cannot be a successful first-start result")
}
