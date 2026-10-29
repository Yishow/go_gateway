package workspace

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

type lifecyclePrivateCause struct{}

func (*lifecyclePrivateCause) Error() string { return "private endpoint credential sample payload" }

func TestLifecycleRemoveClassifiesGuardAndConnectorFailures(t *testing.T) {
	for _, stage := range []string{"guard-missing", "guard-failed", "connector-stale", "connector-missing", "connector-read-failed"} {
		for _, stopped := range []bool{false, true} {
			t.Run(fmt.Sprintf("%s/stopped=%t", stage, stopped), func(t *testing.T) {
				db, service, group := lifecycleReadFixture(t, true)
				cause := &lifecyclePrivateCause{}
				service.WithBacklogOwnershipGuard(writeGroupBacklogGuardFunc(func(context.Context, *sql.Tx, *WriteGroup) error { return nil }))
				if stopped {
					saved, err := service.Disable(t.Context(), group.ID, lifecycleReadMutation(t, service, group))
					require.NoError(t, err)
					group = saved.Group
				}
				switch stage {
				case "guard-missing":
					service.WithBacklogOwnershipGuard(nil)
				case "guard-failed":
					service.WithBacklogOwnershipGuard(writeGroupBacklogGuardFunc(func(context.Context, *sql.Tx, *WriteGroup) error { return cause }))
				case "connector-stale":
					_, err := db.Exec(`UPDATE database_connectors SET identity_revision='C2' WHERE id=?`, group.Destination.ConnectorID)
					require.NoError(t, err)
				case "connector-missing":
					_, err := db.Exec(`DELETE FROM database_connectors WHERE id=?`, group.Destination.ConnectorID)
					require.NoError(t, err)
				case "connector-read-failed":
					_, err := db.Exec(`ALTER TABLE database_connectors RENAME COLUMN identity_revision TO unavailable_revision`)
					require.NoError(t, err)
				}
				before := lifecycleReadState(t, db)
				_, err := service.Delete(t.Context(), group.ID, lifecycleReadMutation(t, service, group))
				require.Error(t, err)
				require.Equal(t, before, lifecycleReadState(t, db), "refused removal never changes the group or durable accepted data")
				switch stage {
				case "guard-missing", "guard-failed":
					require.ErrorIs(t, err, ErrWriteGroupLifecycleBlocked)
				case "connector-stale":
					require.ErrorIs(t, err, ErrWriteGroupConnectorRevisionConflict)
					require.ErrorIs(t, err, ErrSetupRevisionConflict)
					require.ErrorIs(t, err, ErrWriteGroupRevisionConflict)
				case "connector-missing":
					require.ErrorIs(t, err, ErrWriteGroupConnectorNotFound)
					require.ErrorIs(t, err, ErrNotFound)
				}
				if stage == "guard-failed" {
					var actual *lifecyclePrivateCause
					require.ErrorAs(t, err, &actual)
					require.Same(t, cause, actual)
				}
				var classified *WriteGroupLifecycleActionError
				require.ErrorAs(t, err, &classified, "removal returns structured repair metadata instead of only a generic cause")
				require.Equal(t, "remove", classified.Action)
				reason, next := "ownership_unavailable", "inspect_ownership"
				if stage == "connector-stale" {
					reason, next = "connector_identity_stale", "review_connector"
				}
				if stage == "connector-missing" || stage == "connector-read-failed" {
					reason, next = "connector_unavailable", "review_connector"
				}
				if !stopped {
					next = "stop"
				}
				require.Equal(t, reason, classified.ReasonCode)
				require.Equal(t, next, classified.NextAction)
				require.NotContains(t, err.Error(), "private")
				require.NotContains(t, err.Error(), "unavailable_revision")
			})
		}
	}
}

func TestLifecycleRemoveKeepsCASAndForeignScopeErrors(t *testing.T) {
	for _, field := range []string{"workspace", "group", "connector", "scope"} {
		t.Run(field, func(t *testing.T) {
			db, service, group := lifecycleReadFixture(t, true)
			guardCalls := 0
			service.WithBacklogOwnershipGuard(writeGroupBacklogGuardFunc(func(context.Context, *sql.Tx, *WriteGroup) error { guardCalls++; return nil }))
			mutation := lifecycleReadMutation(t, service, group)
			switch field {
			case "workspace":
				mutation.ExpectedWorkspaceRevision = "stale"
			case "group":
				mutation.ExpectedGroupRevision = "stale"
			case "connector":
				mutation.ExpectedConnectorRevision = "stale"
			case "scope":
				mutation.WorkspaceID = "foreign"
			}
			before := lifecycleReadState(t, db)
			_, err := service.Delete(t.Context(), group.ID, mutation)
			if field == "scope" {
				require.ErrorIs(t, err, ErrWriteGroupNotFound)
			} else {
				require.ErrorIs(t, err, ErrWriteGroupRevisionConflict)
			}
			var actionError *WriteGroupLifecycleActionError
			require.False(t, errors.As(err, &actionError), "stale request revisions still require reload, not a stop fallback")
			require.Equal(t, before, lifecycleReadState(t, db))
			if field == "connector" {
				require.Equal(t, 1, guardCalls, "existing ownership-before-saved-connector check order is preserved")
			} else {
				require.Zero(t, guardCalls)
			}
		})
	}
}

var _ error = (*lifecyclePrivateCause)(nil)
