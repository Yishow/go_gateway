package handlers

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
	"go-gateway/internal/datalink/device"
	"go-gateway/internal/datalink/runtimeview"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/workspace"
)

func TestStudioLocatorStrictQueryShapes(t *testing.T) {
	for _, q := range []string{"entry=locator&workspace_id=W&step=1&device_id=D", "entry=locator&workspace_id=W&step=2&device_id=D&rule_id=R", "entry=locator&workspace_id=W&step=3&device_id=D&mapping_id=M", "entry=locator&workspace_id=W&step=4&group_id=G", "entry=locator&workspace_id=W&step=4&group_id=G&device_id=D&group_revision=old"} {
		t.Run(q, func(t *testing.T) {
			got, err := parseStudioLocatorQuery(q)
			require.NoError(t, err)
			require.NotNil(t, got)
		})
	}
	for _, q := range []string{"workspace_id=W", "entry=locator", "entry=locator&workspace_id=W&step=1&device_id=D&rule_id=R", "entry=locator&workspace_id=W&step=2&device_id=D", "entry=locator&workspace_id=W&step=3&device_id=D&mapping_id=M&rule_id=R", "entry=locator&workspace_id=W&step=4&group_id=G&point_id=P", "entry=locator&workspace_id=W&step=4&group_id=G&readiness=1", "entry=locator&workspace_id=W&step=01&device_id=D", "entry=locator&workspace_id=W&step=1&device_id=D&device_id=D", "entry=locator&workspace_id=W&step=1&device_id=%ZZ", "entry=locator;workspace_id=W", "entry=locator&workspace_id=W&step=1&device_id=postgres%3A%2F%2Fsecret", "entry=locator&workspace_id=W&step=1&device_id=D&unknown="} {
		t.Run(q, func(t *testing.T) { _, err := parseStudioLocatorQuery(q); require.Error(t, err) })
	}
	got, err := parseStudioLocatorQuery("")
	require.NoError(t, err)
	require.Nil(t, got)
}

func TestStudioLocatorReadErrorTyped(t *testing.T) {
	for _, err := range []error{sql.ErrNoRows, workspace.ErrNotFound, workspace.ErrWriteGroupNotFound} {
		require.Equal(t, runtimeEvidenceError(runtimeview.ReasonStale), studioLocatorReadError(err))
	}
	require.Equal(t, runtimeEvidenceError(runtimeview.ReasonUnavailable), studioLocatorReadError(errors.New("password=secret")))
}

func TestStudioLocatorTypedNilReaderFailsClosed(t *testing.T) {
	var devices *device.Service
	h := NewStudioV2WorkspaceHandler(nil).WithLocatorEntryServices(StudioV2LocatorReaders{Devices: devices})
	_, err := h.readStudioLocator(context.Background(), &workspace.Record{OrderedDeviceIDs: []string{"D"}}, runtimeview.EntityLocator{Step: 1, DeviceID: runtimeEvidenceString("D")})
	require.Equal(t, runtimeEvidenceError(runtimeview.ReasonUnavailable), err)
}

func TestStudioLocatorComparisonExcludesOnlyVolatileTelemetry(t *testing.T) {
	a := &studioLocatorObservation{Device: &schema.Device{ID: "D", ConnectionConfig: "saved"}, Point: &schema.Point{ID: "P", Address: "40001"}}
	b, _ := cloneStudioLocatorValue(a)
	b.Device.LastTestError = "raw failure"
	b.Point.LastError = "raw failure"
	require.True(t, sameStudioLocatorObservation(a, b))
	b.Point.Address = "40002"
	require.False(t, sameStudioLocatorObservation(a, b))
	b.Point.Address = "40001"
	b.Device.ConnectionConfig = "changed"
	require.False(t, sameStudioLocatorObservation(a, b))
}
