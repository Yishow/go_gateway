package point

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go-gateway/internal/datalink/modbusshare"
)

func TestFormatKnownInvalidationReleasesOnlyOwnedDirty(t *testing.T) {
	for _, foreign := range []bool{false, true} {
		t.Run(map[bool]string{false: "owned", true: "foreign-same-revision"}[foreign], func(t *testing.T) {
			db, svc, req := formatFixture(t)
			revisions := modbusshare.NewSQLWorkspaceRevisionStore(db)
			if foreign {
				svc.SetFormatInvalidator(&recordingFormatInvalidator{run: func() error { return revisions.MarkDirty(t.Context(), req.WorkspaceID) }})
			}
			result, err := svc.MutateFormat(t.Context(), req)
			if foreign {
				require.ErrorIs(t, err, ErrFormatUnavailable)
				require.Nil(t, result)
			} else {
				require.NoError(t, err)
				require.Equal(t, "pending", result.RuntimeApplyStatus)
			}
			_, dirty, err := revisions.GetRevision(t.Context(), req.WorkspaceID)
			require.NoError(t, err)
			require.Equal(t, foreign, dirty)
			receipt, err := svc.InspectFormatMutation(t.Context(), req.FormatScope)
			require.NoError(t, err)
			if foreign {
				require.Equal(t, "pending", receipt.InvalidationStatus)
			} else {
				require.Equal(t, "invalidated", receipt.InvalidationStatus)
			}
		})
	}
}

func TestFormatDirtyOwnershipCannotBeForgedOrOutliveFailedFinalize(t *testing.T) {
	for _, kind := range []string{"missing", "future-version", "foreign-point", "foreign-receipt", "finalize-fault", "callback-fault"} {
		t.Run(kind, func(t *testing.T) {
			db, svc, req := formatFixture(t)
			svc.SetFormatInvalidator(&recordingFormatInvalidator{run: func() error {
				var statement string
				switch kind {
				case "missing":
					statement = `UPDATE system_settings SET value=json_remove(value,'$.point_format_pending_v1') WHERE key='modbus_share_workspace_revision:ws'`
				case "future-version":
					statement = `UPDATE system_settings SET value=json_set(value,'$.point_format_pending_v1.version',2) WHERE key='modbus_share_workspace_revision:ws'`
				case "foreign-point":
					statement = `UPDATE system_settings SET value=json_set(value,'$.point_format_pending_v1.point_id','foreign') WHERE key='modbus_share_workspace_revision:ws'`
				case "foreign-receipt":
					statement = `UPDATE system_settings SET value=json_set(value,'$.point_format_pending_v1.receipt_key','foreign') WHERE key='modbus_share_workspace_revision:ws'`
				case "finalize-fault":
					statement = `CREATE TRIGGER stop_finalize BEFORE UPDATE ON system_settings WHEN NEW.key LIKE 'point_format_receipt:%' BEGIN SELECT RAISE(ABORT,'private failure'); END`
				case "callback-fault":
					return ErrFormatUnavailable
				}
				_, err := db.Exec(statement)
				return err
			}})
			out, err := svc.MutateFormat(t.Context(), req)
			require.ErrorIs(t, err, ErrFormatUnavailable)
			require.Nil(t, out)
			_, dirty, err := modbusshare.NewSQLWorkspaceRevisionStore(db).GetRevision(t.Context(), req.WorkspaceID)
			require.NoError(t, err)
			require.True(t, dirty)
			receipt, err := svc.InspectFormatMutation(t.Context(), req.FormatScope)
			require.NoError(t, err)
			require.Equal(t, "pending", receipt.InvalidationStatus)
		})
	}
}
