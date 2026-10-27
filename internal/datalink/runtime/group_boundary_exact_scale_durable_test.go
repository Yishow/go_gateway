package runtime

import (
	"encoding/json"
	"path/filepath"
	"testing"
	"time"

	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/groupdelivery"
	"go-gateway/internal/datalink/measurement"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/workspace"

	"github.com/stretchr/testify/require"
)

func scaledCounterFixture(t *testing.T, cumulative bool) *longFixture {
	t.Helper()
	f := newCounterFixture(t)
	member := f.members[0]
	member.SourceValueType, member.ResultValueType = schema.DataTypeInt64, measurement.ExactDecimal
	member.ExactRecordingSource = &workspace.ExactRecordingSource{
		MeasurementID: *member.MeasurementID, DefinitionRevision: "def-1", SourceBindingRevision: "bind-1",
		SeriesEpoch: "epoch-1", OriginalDataType: schema.DataTypeInt64, Scale: "0.1", Unit: "kWh",
		SemanticKind: measurement.SemanticKindCounter,
	}
	if !cumulative {
		member.ValueRole = workspace.WriteGroupValueRoleCounterRaw
		f.config.Group.RecordingProfile = &workspace.WriteGroupRecordingProfile{
			Kind: workspace.WriteGroupProfileInstantaneous, Interval: 15, Unit: workspace.WriteGroupIntervalUnitSecond,
		}
		f.config.Group.RowPolicy.IntervalSeconds = 15
	}
	f.members = []workspace.WriteGroupMember{member}
	f.config.Group.Members = f.members
	f.config.Counters = nil // frozen proof is sufficient; editable definitions are not read
	return f
}

func TestExactRecordingScaleBoundary(t *testing.T) {
	for _, cumulative := range []bool{false, true} {
		name, interval := "instantaneous", 15*time.Second
		if cumulative {
			name, interval = "cumulative", counterWindow
		}
		t.Run(name, func(t *testing.T) {
			f := scaledCounterFixture(t, cumulative)
			db := openDurableDB(t, filepath.Join(t.TempDir(), "scaled-journal.db"))
			f.config.Ledger = longLedger(db)
			boundary := f.build(t)
			f.accept(t, boundary, f.envelope(0, "scaled-first", at(interval-time.Second), int64(120000)))
			var payload string
			require.NoError(t, db.QueryRowContext(t.Context(), `SELECT payload FROM wg_delivery_samples`).Scan(&payload))
			var sample struct {
				Value measurement.ExactValue `json:"value"`
			}
			require.NoError(t, json.Unmarshal([]byte(payload), &sample))
			require.Equal(t, measurement.ExactDecimal, sample.Value.Type(), "scaling precedes the journal ACK")
			digits, ok := sample.Value.DecimalDigits()
			require.True(t, ok)
			require.Equal(t, "12000.0", digits)

			// Restore the accepted decimal journal into a new boundary; it must not
			// pass through the integer source conversion a second time.
			recovered := f.build(t)
			require.NoError(t, recovered.Tick(t.Context(), at(interval)))
			first := outboxRows(t, db)["tag-int"]
			require.Equal(t, "12000.0", first["value_decimal"])
			require.Nil(t, first["delta_value"])
			f.clock = at(2 * interval)
			f.accept(t, recovered, f.envelope(0, "scaled-second", at(2*interval-time.Second), int64(120004)))
			require.NoError(t, recovered.Tick(t.Context(), at(2*interval)))
			require.NoError(t, db.QueryRowContext(t.Context(),
				`SELECT payload FROM wg_delivery_outbox ORDER BY bucket_start DESC LIMIT 1`).Scan(&payload))
			row, err := groupdelivery.DecodeRowPayload([]byte(payload))
			require.NoError(t, err)
			second := make(map[string]any, len(row.Cells))
			for _, cell := range row.Cells {
				second[cell.Column] = cell.Value
			}
			require.Equal(t, "decimal", second["value_type"])
			require.Equal(t, "12000.4", second["value_decimal"])
			if cumulative {
				require.Equal(t, "0.4", second["delta_value"])
				require.Equal(t, "ok", second["status"])
			} else {
				require.Nil(t, second["delta_value"])
			}
		})
	}
}

func TestExactRecordingScaleRejectsUnrepresentableBeforeJournal(t *testing.T) {
	f := scaledCounterFixture(t, false)
	f.members[0].ExactRecordingSource.Scale = "0.0000000000000000000000000000001"
	db := openDurableDB(t, filepath.Join(t.TempDir(), "scaled-range.db"))
	f.config.Ledger = longLedger(db)
	boundary := f.build(t)
	// Only the declared destination value range is substituted. This is an
	// intake test, without a claim of a real MySQL catalog or delivery.
	_, declared, _ := mysqlIntakeFixture(t, "DECIMAL(65,30)")
	boundary.layout = declared.layout
	envelope := f.envelope(0, "same-sample", at(time.Second), int64(1))
	require.ErrorIs(t, boundary.AcceptSample(t.Context(), envelope), dbtarget.ErrExactValueUnrepresentable)
	require.Zero(t, countJournal(t, db))
	require.Zero(t, baselineCount(t, db))
	envelope.Value = int64(0)
	require.NoError(t, boundary.AcceptSample(t.Context(), envelope))
	require.Equal(t, 1, countJournal(t, db), "range refusal cannot consume the sample identity")
}

func TestExactRecordingScaleActivationRejectsBadProof(t *testing.T) {
	for name, mutate := range map[string]func(*longFixture){
		"missing proof":     func(f *longFixture) { f.config.Group.Members[0].ExactRecordingSource = nil },
		"invalid scale":     func(f *longFixture) { f.config.Group.Members[0].ExactRecordingSource.Scale = "NaN" },
		"wrong result type": func(f *longFixture) { f.config.Group.Members[0].ResultValueType = measurement.ExactInt64 },
		"wrong original type": func(f *longFixture) {
			f.config.Group.Members[0].ExactRecordingSource.OriginalDataType = schema.DataTypeUint64
		},
		"missing measurement":     func(f *longFixture) { f.config.Group.Members[0].MeasurementID = nil },
		"missing source revision": func(f *longFixture) { f.config.Group.Members[0].SourceRevision = "" },
		"wide layout":             func(f *longFixture) { f.config.Group.LayoutVersion = workspace.WriteGroupLayoutManagedWideV1 },
	} {
		t.Run(name, func(t *testing.T) {
			f := scaledCounterFixture(t, true)
			mutate(f)
			_, err := NewGroupBoundary(f.config)
			require.Error(t, err, "invalid frozen evidence cannot fall back to unscaled intake")
		})
	}
}
