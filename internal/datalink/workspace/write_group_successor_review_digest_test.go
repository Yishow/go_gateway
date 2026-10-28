package workspace

import (
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
	"time"
)

func TestSuccessorEnvelopeBindsEveryScopeAndMeaning(t *testing.T) {
	h, id, req := successorReviewFixture(t)
	view, err := h.service.SuccessorReview(t.Context(), id, req)
	require.NoError(t, err)
	base, review, reservation, err := normalizeSuccessorEnvelope(view.Envelope)
	require.NoError(t, err)
	require.NotEqual(t, review, reservation, "review and reservation have different domains")
	mutations := map[string]func(*WriteGroupSuccessorEnvelope){
		"workspace": func(e *WriteGroupSuccessorEnvelope) {
			e.Request.WorkspaceID = "foreign"
			e.Request.Group.WorkspaceID = "foreign"
		},
		"workspace-revision":        func(e *WriteGroupSuccessorEnvelope) { e.Request.ExpectedWorkspaceRevision = "W2" },
		"lineage":                   func(e *WriteGroupSuccessorEnvelope) { e.Request.LineageID = "L2" },
		"lineage-revision":          func(e *WriteGroupSuccessorEnvelope) { e.Request.ExpectedLineageRevision++ },
		"predecessor":               func(e *WriteGroupSuccessorEnvelope) { e.PredecessorGroupID = "G2" },
		"predecessor-revision":      func(e *WriteGroupSuccessorEnvelope) { e.Request.ExpectedPredecessorRevision = "R2" },
		"source-connector-revision": func(e *WriteGroupSuccessorEnvelope) { e.Request.ExpectedPredecessorConnectorRevision = "C2" },
		"target-connector-revision": func(e *WriteGroupSuccessorEnvelope) {
			e.Request.ExpectedTargetConnectorRevision = "C2"
			e.Request.Group.Destination.ConnectorRevision = "C2"
		},
		"generation":          func(e *WriteGroupSuccessorEnvelope) { e.Generation++ },
		"source-applied":      func(e *WriteGroupSuccessorEnvelope) { e.SourceAppliedRevision = "A2" },
		"source-definition":   func(e *WriteGroupSuccessorEnvelope) { e.SourceDefinitionDigest = strings.Repeat("d", 64) },
		"operation":           func(e *WriteGroupSuccessorEnvelope) { e.Request.OperationID = "O3" },
		"boundary-nanosecond": func(e *WriteGroupSuccessorEnvelope) { e.Request.Boundary = e.Request.Boundary.Add(time.Nanosecond) },
		"target-name":         func(e *WriteGroupSuccessorEnvelope) { e.Request.Group.Name = "different" },
		"source-revision":     func(e *WriteGroupSuccessorEnvelope) { e.Request.Group.Members[0].SourceRevision = "S2" },
	}
	for name, change := range mutations {
		t.Run(name, func(t *testing.T) {
			e := base
			e.Request.Group = cloneWriteGroup(base.Request.Group)
			change(&e)
			_, rd, sd, err := normalizeSuccessorEnvelope(e)
			require.NoError(t, err)
			require.NotEqual(t, review, rd)
			require.NotEqual(t, reservation, sd)
		})
	}
	e := base
	e.Request.Group = cloneWriteGroup(base.Request.Group)
	e.Request.Boundary = e.Request.Boundary.In(time.FixedZone("same instant", 8*3600))
	e.Request.Group.ID = "discard-server-id"
	e.Request.Group.Revision = "discard-server-revision"
	_, rd, sd, err := normalizeSuccessorEnvelope(e)
	require.NoError(t, err)
	require.Equal(t, review, rd)
	require.Equal(t, reservation, sd)
}

func TestSuccessorEnvelopeRejectsMalformedAuthority(t *testing.T) {
	h, id, req := successorReviewFixture(t)
	view, err := h.service.SuccessorReview(t.Context(), id, req)
	require.NoError(t, err)
	for _, change := range []func(*WriteGroupSuccessorEnvelope){
		func(e *WriteGroupSuccessorEnvelope) { e.Version = "future" }, func(e *WriteGroupSuccessorEnvelope) { e.Generation = 0 },
		func(e *WriteGroupSuccessorEnvelope) { e.SourceDefinitionDigest = strings.Repeat("z", 64) }, func(e *WriteGroupSuccessorEnvelope) { e.Request.Boundary = time.Time{} },
		func(e *WriteGroupSuccessorEnvelope) { e.PredecessorGroupID = "" }, func(e *WriteGroupSuccessorEnvelope) { e.Request.Group = nil },
	} {
		e := view.Envelope
		change(&e)
		_, _, _, err := normalizeSuccessorEnvelope(e)
		require.Error(t, err)
	}
}
