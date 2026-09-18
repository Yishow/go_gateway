package handlers

import (
	"testing"

	"go-gateway/internal/datalink/dbtarget"

	"github.com/stretchr/testify/require"
)

func workspaceInspection(status dbtarget.TableInspectionStatus, table string) *dbtarget.TableInspection {
	return &dbtarget.TableInspection{Status: status, Table: table}
}

// EnsureWorkspaceSchema readiness is a pure decision over the dry-run
// statement count and corroborating table inspections, so every acceptance
// state is pinned here and the wiring stays covered by the fixture tests.
func TestEvaluateWorkspaceSchemaReadiness(t *testing.T) {
	t.Run("empty statements with every inspected table present are ready", func(t *testing.T) {
		require.NoError(t, evaluateWorkspaceSchemaReadiness(0, []*dbtarget.TableInspection{
			workspaceInspection(dbtarget.TableInspectionExists, "sensor_values"),
		}))
	})

	t.Run("a missing table requires preparation even without statements", func(t *testing.T) {
		err := evaluateWorkspaceSchemaReadiness(0, []*dbtarget.TableInspection{
			workspaceInspection(dbtarget.TableInspectionExists, "sensor_values"),
			workspaceInspection(dbtarget.TableInspectionMissing, "gw_record_events"),
		})
		require.ErrorIs(t, err, errWorkspaceSchemaPreparationRequired)
	})

	t.Run("remaining dry-run statements require preparation", func(t *testing.T) {
		err := evaluateWorkspaceSchemaReadiness(2, []*dbtarget.TableInspection{
			workspaceInspection(dbtarget.TableInspectionExists, "sensor_values"),
		})
		require.ErrorIs(t, err, errWorkspaceSchemaPreparationRequired)
	})

	t.Run("an unconfirmable inspection is a safe error instead of false readiness", func(t *testing.T) {
		for _, status := range []dbtarget.TableInspectionStatus{dbtarget.TableInspectionForbidden, dbtarget.TableInspectionFailed} {
			err := evaluateWorkspaceSchemaReadiness(0, []*dbtarget.TableInspection{
				workspaceInspection(status, "sensor_values"),
			})
			require.Error(t, err, "status %s must not be ready", status)
			require.NotErrorIs(t, err, errWorkspaceSchemaPreparationRequired, "status %s must not request preparation", status)
			require.Contains(t, err.Error(), "sensor_values")
		}
	})

	t.Run("an empty inspection list falls back to the statement count", func(t *testing.T) {
		require.NoError(t, evaluateWorkspaceSchemaReadiness(0, nil))
		require.ErrorIs(t, evaluateWorkspaceSchemaReadiness(1, nil), errWorkspaceSchemaPreparationRequired)
	})
}
