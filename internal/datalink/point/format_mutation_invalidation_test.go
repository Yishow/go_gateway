package point

import (
	"context"
	"errors"
	"github.com/stretchr/testify/require"
	"go-gateway/internal/datalink/schema"
	"testing"
)

type recordingFormatInvalidator struct {
	calls int
	run   func() error
	seen  FormatInvalidation
}

func (i *recordingFormatInvalidator) InvalidatePointFormat(_ context.Context, change FormatInvalidation) error {
	i.calls++
	i.seen = change
	if i.run != nil {
		return i.run()
	}
	return nil
}

func TestFormatInvalidationIsOnceAndPendingCanBeInspected(t *testing.T) {
	for _, kind := range []string{"success", "callback-failure", "finalize-failure", "newer-revision"} {
		t.Run(kind, func(t *testing.T) {
			db, svc, req := formatFixture(t)
			spy := &recordingFormatInvalidator{}
			switch kind {
			case "callback-failure":
				spy.run = func() error { return errors.New("secret://callback") }
			case "finalize-failure":
				spy.run = func() error {
					_, err := db.Exec(`CREATE TRIGGER reject_receipt_finalize BEFORE UPDATE ON system_settings WHEN NEW.key LIKE 'point_format_receipt:%' BEGIN SELECT RAISE(ABORT,'secret://finalize'); END;`)
					return err
				}
			case "newer-revision":
				spy.run = func() error { _, err := db.Exec(`UPDATE points SET revision_id='newer'`); return err }
			}
			svc.SetFormatInvalidator(spy)
			first, err := svc.MutateFormat(t.Context(), req)
			if kind == "success" {
				require.NoError(t, err)
				require.Equal(t, "invalidated", first.InvalidationStatus)
			} else {
				require.ErrorIs(t, err, ErrFormatUnavailable)
				require.Nil(t, first)
				require.NotContains(t, err.Error(), "secret")
			}
			require.Equal(t, 1, spy.calls)
			status, err := svc.InspectFormatMutation(t.Context(), req.FormatScope)
			require.NoError(t, err)
			require.Equal(t, req.PointID, status.Point.ID)
			if kind == "success" {
				require.Equal(t, "invalidated", status.InvalidationStatus)
			} else {
				require.Equal(t, "pending", status.InvalidationStatus)
			}
			replay, err := svc.MutateFormat(t.Context(), req)
			if kind == "success" {
				require.NoError(t, err)
				require.Equal(t, first, replay)
			} else {
				require.ErrorIs(t, err, ErrFormatUnavailable)
				require.Nil(t, replay)
			}
			require.Equal(t, 1, spy.calls, "pending retries must not replay a destructive invalidation")
			foreign := req.FormatScope
			foreign.DeviceID = "foreign"
			_, err = svc.InspectFormatMutation(t.Context(), foreign)
			require.ErrorIs(t, err, ErrFormatScope)
		})
	}
}

func TestFormatInvalidationDoesNotRunOnRejectedOrFailedTransaction(t *testing.T) {
	db, svc, req := formatFixture(t)
	spy := &recordingFormatInvalidator{}
	svc.SetFormatInvalidator(spy)
	bad := req
	bad.ExpectedWorkspaceRevision = "stale"
	_, err := svc.MutateFormat(t.Context(), bad)
	require.ErrorIs(t, err, ErrFormatConflict)
	require.Zero(t, spy.calls)
	_, err = db.Exec(`CREATE TRIGGER reject_point BEFORE UPDATE ON points BEGIN SELECT RAISE(ABORT,'failure'); END;`)
	require.NoError(t, err)
	_, err = svc.MutateFormat(t.Context(), req)
	require.ErrorIs(t, err, ErrFormatUnavailable)
	require.Zero(t, spy.calls)
}

func TestFormatReceiptUnknownVersionAndResolutionAreUnavailable(t *testing.T) {
	for _, path := range []string{"$.version", "$.result.resolution.point_revision"} {
		t.Run(path, func(t *testing.T) {
			db, svc, req := formatFixture(t)
			_, err := svc.MutateFormat(t.Context(), req)
			require.NoError(t, err)
			_, err = db.Exec(`UPDATE system_settings SET value=json_set(value,?,?) WHERE key LIKE 'point_format_receipt:%'`, path, "unknown")
			require.NoError(t, err)
			_, err = svc.MutateFormat(t.Context(), req)
			require.ErrorIs(t, err, ErrFormatUnavailable)
			_, err = svc.InspectFormatMutation(t.Context(), req.FormatScope)
			require.ErrorIs(t, err, ErrFormatUnavailable)
		})
	}
}

func TestFormatReceiptCorruptResolvedSemanticsAreUnavailable(t *testing.T) {
	for _, change := range []struct {
		path  string
		value any
	}{
		{"$.result.resolution.resolved", false}, {"$.result.resolution.effective_format", "UNKNOWN"},
		{"$.result.resolution.format_mode", "legacy"}, {"$.result.resolution.configured_format", "BADC"},
		{"$.result.resolution.effective_source", "native"}, {"$.result.point.format_version", 999},
	} {
		t.Run(change.path, func(t *testing.T) {
			db, svc, req := formatFixture(t)
			_, err := svc.MutateFormat(t.Context(), req)
			require.NoError(t, err)
			_, err = db.Exec(`UPDATE system_settings SET value=json_set(value,?,?) WHERE key LIKE 'point_format_receipt:%'`, change.path, change.value)
			require.NoError(t, err)
			_, err = svc.MutateFormat(t.Context(), req)
			require.ErrorIs(t, err, ErrFormatUnavailable)
			_, err = svc.InspectFormatMutation(t.Context(), req.FormatScope)
			require.ErrorIs(t, err, ErrFormatUnavailable)
		})
	}
}

func TestFormatReceiptReplaysFrozenDefaultsNotCurrentConnection(t *testing.T) {
	db, svc, req := formatFixture(t)
	req.FormatMode = schema.PointFormatInherit
	req.ConfiguredFormat = ""
	first, err := svc.MutateFormat(t.Context(), req)
	require.NoError(t, err)
	_, err = db.Exec(`UPDATE devices SET connection_config='{"data_format":"DCBA"}'`)
	require.NoError(t, err)
	replay, err := svc.MutateFormat(t.Context(), req)
	require.NoError(t, err)
	require.Equal(t, first, replay)
	status, err := svc.InspectFormatMutation(t.Context(), req.FormatScope)
	require.NoError(t, err)
	require.Equal(t, first, status)
}

func TestFormatReceiptReplayDoesNotParseChangedConnectionConfig(t *testing.T) {
	db, svc, req := formatFixture(t)
	first, err := svc.MutateFormat(t.Context(), req)
	require.NoError(t, err)
	_, err = db.Exec(`UPDATE devices SET connection_config='invalid-json'`)
	require.NoError(t, err)
	got, err := svc.MutateFormat(t.Context(), req)
	require.NoError(t, err)
	require.Equal(t, first, got)
	_, err = svc.InspectFormatMutation(t.Context(), req.FormatScope)
	require.NoError(t, err)
}

func TestFormatReceiptInheritedNativeAndSingleWordPoliciesRemainValid(t *testing.T) {
	for _, tc := range []struct {
		protocol schema.ProtocolType
		address  string
		dtype    schema.DataType
		format   string
	}{
		{schema.ProtocolFatekFBs, "D0", schema.DataTypeUint32, ""},
		{schema.ProtocolFatekFBs, "DR0", schema.DataTypeUint64, ""},
		{schema.ProtocolFatekFBs, "D0", schema.DataTypeUint16, "DCBA"},
		{schema.ProtocolMC3E, "D0", schema.DataTypeUint16, "DCBA"},
	} {
		t.Run(string(tc.protocol)+string(tc.dtype), func(t *testing.T) {
			db, svc, req := formatFixture(t)
			_, err := db.Exec(`UPDATE devices SET protocol=?,connection_config=json_object('data_format',?)`, tc.protocol, tc.format)
			require.NoError(t, err)
			_, err = db.Exec(`UPDATE points SET address=?,data_type=?`, tc.address, tc.dtype)
			require.NoError(t, err)
			req.FormatMode = schema.PointFormatInherit
			req.ConfiguredFormat = ""
			first, err := svc.MutateFormat(t.Context(), req)
			require.NoError(t, err)
			replay, err := svc.MutateFormat(t.Context(), req)
			require.NoError(t, err)
			require.Equal(t, first, replay)
			if tc.protocol == schema.ProtocolFatekFBs {
				require.Empty(t, replay.Resolution.EffectiveFormat)
				require.True(t, replay.Resolution.NativeCompatibility)
			}
		})
	}
}
