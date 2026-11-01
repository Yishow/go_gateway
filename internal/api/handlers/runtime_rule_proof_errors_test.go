package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"go-gateway/internal/datalink/grouppipeline"
	"go-gateway/internal/datalink/runtimeview"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/workspace"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func TestRuntimeSavedSnapshotRejectsUnencodableProof(t *testing.T) {
	for _, value := range []float64{math.NaN(), math.Inf(1), math.Inf(-1)} {
		t.Run(fmtProofValue(value), func(t *testing.T) {
			h, record := runtimeEvidenceScopeFixture(t)
			rules := []*schema.SourceRule{{ID: "rule-1", DeviceID: "device-1", ScaleMultiplier: &value}}
			h.sourceRuleReader = runtimeSetupRuleReaderStub{rules: rules}
			before, err := json.Marshal(rules)
			require.Error(t, err)
			require.Nil(t, before)
			require.Error(t, h.checkRuntimeSavedSnapshot(t.Context(), record, before), "two failed encodings must not establish snapshot equality")
		})
	}
}

func fmtProofValue(value float64) string {
	if math.IsNaN(value) {
		return "nan"
	}
	if math.IsInf(value, 1) {
		return "positive-infinity"
	}
	return "negative-infinity"
}

func TestRuntimeSetupRejectsUnencodableInitialSnapshot(t *testing.T) {
	h, record := runtimeEvidenceScopeFixture(t)
	value := math.NaN()
	h.sourceRuleReader = runtimeSetupRuleReaderStub{rules: []*schema.SourceRule{{ID: "rule-1", DeviceID: "device-1", ScaleOffset: &value}}}
	_, err := h.buildWorkspaceSetupContext(t.Context(), record)
	require.Error(t, err)
	rec := runtimeEvidenceScopeRequest(t, h, "")
	require.Equal(t, http.StatusServiceUnavailable, rec.Code)
	require.NotContains(t, rec.Body.String(), "NaN")
	require.NotContains(t, rec.Body.String(), "scale_offset")
}

type runtimeAppliedErrorReader struct{ err error }

func (r runtimeAppliedErrorReader) ResolveAppliedConfigurationAt(context.Context, string, time.Time) (*workspace.WriteGroupAppliedSnapshot, error) {
	return nil, r.err
}

func TestRuntimeAppliedProviderErrorIsNotMissingProof(t *testing.T) {
	for _, tc := range []struct {
		name   string
		err    error
		status int
	}{
		{"provider-error", errors.New("postgres://private-secret"), http.StatusServiceUnavailable},
		{"canceled-read", context.Canceled, http.StatusServiceUnavailable},
		{"deadline-read", context.DeadlineExceeded, http.StatusServiceUnavailable},
		{"stale", grouppipeline.ErrAppliedConfigurationStale, http.StatusConflict},
		{"missing-group", workspace.ErrWriteGroupNotFound, http.StatusConflict},
		{"missing-proof", nil, 0},
		{"known-unavailable", grouppipeline.ErrAppliedConfigurationUnavailable, 0},
		{"wrapped-unavailable", fmt.Errorf("context: %w", grouppipeline.ErrAppliedConfigurationUnavailable), 0},
	} {
		t.Run(tc.name, func(t *testing.T) {
			h := &RuntimeHandler{ruleApplied: runtimeAppliedErrorReader{err: tc.err}}
			snapshot, err := h.runtimeRuleSnapshotAt(t.Context(), &workspace.WriteGroup{ID: "group-1", AppliedRevision: "revision-1"}, runtimeview.Scope{}, time.Now())
			require.Nil(t, snapshot)
			if tc.status == 0 {
				require.NoError(t, err, "a genuinely absent proof retains the unavailable projection")
				return
			}
			require.Error(t, err)
			rec := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(rec)
			renderRuntimeEvidenceError(c, err)
			require.Equal(t, tc.status, rec.Code)
			require.NotContains(t, rec.Body.String(), "private-secret")
		})
	}
}

func TestRuntimeAppliedProviderFailureHTTPIsSanitized(t *testing.T) {
	for _, failAt := range []int{1, 2} {
		t.Run(fmtProofRead(failAt), func(t *testing.T) {
			h, record := runtimeEvidenceScopeFixture(t)
			saved, _, _, _, at := runtimeDeliveryProjectionFixture()
			saved.WorkspaceID = record.ID
			saved.Members[0].DeviceID = "device-1"
			reader := &runtimeAppliedFailingRead{failAt: failAt, err: errors.New("postgres://private-secret")}
			h.WithRuleEvidenceServices(runtimeDeliveryGroupsStub{&workspace.WriteGroupListResult{WorkspaceID: record.ID, WorkspaceRevision: record.DatabaseSetupRevision, Groups: []*workspace.WriteGroup{saved}}}, reader).WithRuleEvidenceClock(func() time.Time { return at })
			query := url.Values{"workspace_id": {record.ID}, "workspace_revision": {record.DatabaseSetupRevision}, "device_id": {"device-1"}}
			rec := runtimeEvidenceScopeRequest(t, h, "?"+query.Encode())
			require.Equal(t, failAt, reader.calls)
			require.Equal(t, http.StatusServiceUnavailable, rec.Code, rec.Body.String())
			require.NotContains(t, rec.Body.String(), "private-secret")
			require.NotContains(t, rec.Body.String(), `"rule_evidence"`)
		})
	}
}

func fmtProofRead(read int) string {
	if read == 1 {
		return "initial"
	}
	return "recheck"
}

type runtimeAppliedFailingRead struct {
	calls, failAt int
	err           error
}

func (r *runtimeAppliedFailingRead) ResolveAppliedConfigurationAt(context.Context, string, time.Time) (*workspace.WriteGroupAppliedSnapshot, error) {
	r.calls++
	if r.calls == r.failAt {
		return nil, r.err
	}
	return nil, nil
}

func TestRuntimeAppliedStaleRecheckHTTPRemainsConflict(t *testing.T) {
	h, record := runtimeEvidenceScopeFixture(t)
	saved, _, _, _, at := runtimeDeliveryProjectionFixture()
	saved.WorkspaceID = record.ID
	saved.Members[0].DeviceID = "device-1"
	reader := &runtimeAppliedFailingRead{failAt: 2, err: grouppipeline.ErrAppliedConfigurationStale}
	h.WithRuleEvidenceServices(runtimeDeliveryGroupsStub{&workspace.WriteGroupListResult{WorkspaceID: record.ID, WorkspaceRevision: record.DatabaseSetupRevision, Groups: []*workspace.WriteGroup{saved}}}, reader).WithRuleEvidenceClock(func() time.Time { return at })
	query := url.Values{"workspace_id": {record.ID}, "workspace_revision": {record.DatabaseSetupRevision}, "device_id": {"device-1"}}
	rec := runtimeEvidenceScopeRequest(t, h, "?"+query.Encode())
	require.Equal(t, 2, reader.calls)
	require.Equal(t, http.StatusConflict, rec.Code, rec.Body.String())
}

func TestRuntimeExpectedUnavailableProofHTTPRetainsHistoricalProjection(t *testing.T) {
	for _, failAt := range []int{1, 2} {
		t.Run(fmtProofRead(failAt), func(t *testing.T) {
			h, record := runtimeEvidenceScopeFixture(t)
			saved, _, _, _, at := runtimeDeliveryProjectionFixture()
			saved.WorkspaceID = record.ID
			saved.Members[0].DeviceID = "device-1"
			reader := &runtimeAppliedFailingRead{failAt: failAt, err: fmt.Errorf("context: %w", grouppipeline.ErrAppliedConfigurationUnavailable)}
			h.WithRuleEvidenceServices(runtimeDeliveryGroupsStub{&workspace.WriteGroupListResult{WorkspaceID: record.ID, WorkspaceRevision: record.DatabaseSetupRevision, Groups: []*workspace.WriteGroup{saved}}}, reader).WithRuleEvidenceClock(func() time.Time { return at })
			query := url.Values{"workspace_id": {record.ID}, "workspace_revision": {record.DatabaseSetupRevision}, "device_id": {"device-1"}}
			rec := runtimeEvidenceScopeRequest(t, h, "?"+query.Encode())
			require.Equal(t, 2, reader.calls)
			require.Equal(t, http.StatusOK, rec.Code, rec.Body.String())
			require.Contains(t, rec.Body.String(), `"proof":"unavailable"`)
			require.Contains(t, rec.Body.String(), `"effective_revision":null`)
			require.NotContains(t, rec.Body.String(), "context:")
		})
	}
}

func TestRuntimeWrappedStaleSnapshotHTTPPreservesConflict(t *testing.T) {
	h, _ := runtimeEvidenceScopeFixture(t)
	h.sourceRuleReader = runtimeWrappedStaleRules{}
	rec := runtimeEvidenceScopeRequest(t, h, "")
	require.Equal(t, http.StatusConflict, rec.Code, rec.Body.String())
	require.Contains(t, rec.Body.String(), "stale-reference")
	require.NotContains(t, rec.Body.String(), "private-provider")
}

type runtimeWrappedStaleRules struct {
	runtimeWorkspaceSourceRuleReader
}

func (runtimeWrappedStaleRules) ListByDeviceIDs(context.Context, []string) ([]*schema.SourceRule, error) {
	return nil, fmt.Errorf("private-provider: %w", errRuntimeStaleReference)
}
