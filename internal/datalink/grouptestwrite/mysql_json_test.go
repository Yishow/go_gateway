package grouptestwrite

import (
	"strings"
	"testing"
	"time"

	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/measurement"
	"go-gateway/internal/datalink/schema"

	"github.com/stretchr/testify/require"
)

func TestMySQLVerifyAcceptsCanonicalJSONWithoutLosingIntegerPrecision(t *testing.T) {
	for _, test := range []struct {
		name, stored, outcome string
	}{
		{"canonical-spacing", `{"n": 9007199254740993, "sample_id": "gw-test-op1-0"}`, WriteVerified},
		{"different-large-integer", `{"n": 9007199254740992, "sample_id": "gw-test-op1-0"}`, WriteUnverified},
	} {
		t.Run(test.name, func(t *testing.T) {
			db := openMySQLReadbackDB(t)
			recordID := strings.Repeat("a", 64)
			state := mysqlReadbackState(recordID)
			_, err := db.ExecContext(t.Context(), "INSERT INTO gw_long_test VALUES (?,?,?,?,?,?,?)",
				recordID, state.WorkspaceID, state.GroupID, state.GroupRevision, test.stored, "gateway", "12000.400000000000000000000000000000")
			require.NoError(t, err)
			_, err = db.ExecContext(t.Context(), "INSERT INTO gw_effect_receipts VALUES (?,?,?)", "effect-1", "digest-1", "2026-01-01T00:00:00.000000Z")
			require.NoError(t, err)
			value, err := measurement.NewDecimal("12000.4")
			require.NoError(t, err)
			p := mysqlReadbackPlan(t, value, nil)
			row := dbtarget.EncodedRow{RecordID: recordID, EffectKey: "effect-1", Cells: []dbtarget.EncodedCell{
				{Column: "record_id", Value: recordID}, {Column: "value_decimal", Value: "12000.4"},
				{Column: mysqlReadbackOwner, Value: "gateway"},
				{Column: "provenance", Value: `{"sample_id":"gw-test-op1-0","n":9007199254740993}`},
			}}
			svc := New(Dependencies{}, Config{StepTimeout: time.Second})
			result := svc.verifyMySQL(t.Context(), &dbtarget.OpenedDestination{DB: db, Kind: schema.DatabaseConnectorKindMySQL}, p, state, []dbtarget.EncodedRow{row})
			require.Equal(t, test.outcome, result.outcome, result.reason)
		})
	}
}
