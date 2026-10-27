package measurement

import (
	"database/sql"
	"os"
	"testing"

	"go-gateway/internal/datalink/schema/migrations"

	"github.com/stretchr/testify/require"
)

func scaleCounterDefinition(scale string) *MeasurementDefinition {
	return &MeasurementDefinition{
		ID: "meas-scale", WorkspaceID: "ws", DeviceID: "dev", PointID: "pt", Name: "Energy",
		Quantity: "energy", Unit: "kWh", SemanticKind: SemanticKindCounter, NumericEncoding: "uint16",
		DefinitionRevision: "rev-1", SourceBindingRevision: "rev-1", SeriesEpoch: "epoch-1",
		ExactRecordingScale: scale,
	}
}

func TestExactRecordingScaleDefinition(t *testing.T) {
	for _, scale := range []string{"0.1", "0.001", "1", "12000.4"} {
		require.NoError(t, scaleCounterDefinition(scale).Validate(), scale)
	}
	require.NoError(t, scaleCounterDefinition("").Validate(), "empty scale preserves old behavior")

	signed := scaleCounterDefinition("0.01")
	signed.SemanticKind = SemanticKindSignedCounter
	signed.CounterPolicy = &CounterPolicy{AllowNegative: true}
	require.NoError(t, signed.Validate())

	for name, mutate := range map[string]func(*MeasurementDefinition){
		"zero":           func(d *MeasurementDefinition) { d.ExactRecordingScale = "0" },
		"negative":       func(d *MeasurementDefinition) { d.ExactRecordingScale = "-0.1" },
		"exponent":       func(d *MeasurementDefinition) { d.ExactRecordingScale = "1e-2" },
		"nan":            func(d *MeasurementDefinition) { d.ExactRecordingScale = "NaN" },
		"infinity":       func(d *MeasurementDefinition) { d.ExactRecordingScale = "Infinity" },
		"spaces":         func(d *MeasurementDefinition) { d.ExactRecordingScale = " 0.1" },
		"hex":            func(d *MeasurementDefinition) { d.ExactRecordingScale = "0x10" },
		"gauge":          func(d *MeasurementDefinition) { d.SemanticKind = SemanticKindGauge },
		"float encoding": func(d *MeasurementDefinition) { d.NumericEncoding = "float64" },
		"empty encoding": func(d *MeasurementDefinition) { d.NumericEncoding = "" },
		"signed no policy": func(d *MeasurementDefinition) {
			d.SemanticKind = SemanticKindSignedCounter
			d.CounterPolicy = nil
		},
	} {
		t.Run(name, func(t *testing.T) {
			def := scaleCounterDefinition("0.1")
			mutate(def)
			require.ErrorIs(t, def.Validate(), ErrExactRecordingScaleInvalid)
		})
	}

	require.Equal(t, CounterConfirmed, ConfirmCounter(*scaleCounterDefinition("0.1")))
	invalid := scaleCounterDefinition("0")
	require.Equal(t, CounterUnconfirmed, ConfirmCounter(*invalid))
	gauge := scaleCounterDefinition("0.1")
	gauge.SemanticKind = SemanticKindGauge
	require.Equal(t, CounterNotCounter, ConfirmCounter(*gauge))
}

func TestExactRecordingScaleEpoch(t *testing.T) {
	for name, scales := range map[string][2]string{
		"changed":          {"0.1", "0.2"},
		"removed":          {"0.1", ""},
		"added from empty": {"", "0.1"},
	} {
		oldDef := *scaleCounterDefinition(scales[0])
		newDef := oldDef
		newDef.ExactRecordingScale = scales[1]
		changed, epoch, reason := CheckEpochTransition(oldDef, newDef)
		require.True(t, changed, name)
		require.Equal(t, "epoch-2", epoch, name)
		require.Contains(t, reason, "exact recording scale", name)
	}
	oldDef := *scaleCounterDefinition("0.1")
	same := oldDef
	same.ExactRecordingScale = "0.10"
	changed, epoch, reason := CheckEpochTransition(oldDef, same)
	require.False(t, changed, "equal numeric scale does not split the epoch")
	require.Equal(t, "epoch-1", epoch)
	require.Empty(t, reason)

	repo := NewMemoryRepository()
	service := NewService(repo)
	def := scaleCounterDefinition("0.1")
	require.NoError(t, service.CreateMeasurement(t.Context(), def))
	def.ExactRecordingScale = "0.2"
	changed, err := service.UpdateMeasurement(t.Context(), def)
	require.NoError(t, err)
	require.True(t, changed)
	stored, err := service.GetMeasurement(t.Context(), def.ID)
	require.NoError(t, err)
	require.Equal(t, "epoch-2", stored.SeriesEpoch)
	require.Equal(t, "rev-2", stored.DefinitionRevision)

	stored.ExactRecordingScale = "0.20"
	changed, err = service.UpdateMeasurement(t.Context(), stored)
	require.NoError(t, err)
	require.False(t, changed)
	again, err := service.GetMeasurement(t.Context(), def.ID)
	require.NoError(t, err)
	require.Equal(t, "epoch-2", again.SeriesEpoch)
	require.Equal(t, "rev-3", again.DefinitionRevision)
}

func TestExactRecordingScaleRepositoryRoundTrip(t *testing.T) {
	db := measurementRepositoryDB(t)
	repo := NewSQLRepository(db)
	def := scaleCounterDefinition("0.001")
	require.NoError(t, repo.Create(t.Context(), def))

	got, err := repo.GetByID(t.Context(), def.ID)
	require.NoError(t, err)
	require.Equal(t, "0.001", got.ExactRecordingScale)

	list, err := repo.ListByWorkspace(t.Context(), def.WorkspaceID)
	require.NoError(t, err)
	require.Len(t, list, 1)
	require.Equal(t, "0.001", list[0].ExactRecordingScale)

	require.NoError(t, db.QueryRowContext(t.Context(), `SELECT COUNT(*) FROM measurement_definitions WHERE exact_recording_scale='0.001'`).Scan(new(int)))

	def.ExactRecordingScale = "0.10"
	require.NoError(t, repo.Update(t.Context(), def))
	updated, err := repo.GetByID(t.Context(), def.ID)
	require.NoError(t, err)
	require.Equal(t, "0.10", updated.ExactRecordingScale, "the exact scale text is preserved")

	byDevice, err := repo.ListByDevice(t.Context(), def.WorkspaceID, def.DeviceID)
	require.NoError(t, err)
	require.Len(t, byDevice, 1)
	require.Equal(t, "0.10", byDevice[0].ExactRecordingScale)

	tx, err := db.BeginTx(t.Context(), nil)
	require.NoError(t, err)
	t.Cleanup(func() { _ = tx.Rollback() })
	inTx, err := repo.GetByIDInTx(t.Context(), tx, def.ID)
	require.NoError(t, err)
	require.Equal(t, "0.10", inTx.ExactRecordingScale)
}

func TestExactRecordingScaleMigrationKeepsLegacyRows(t *testing.T) {
	db, err := sql.Open("sqlite", t.TempDir()+"/legacy.db")
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })
	legacy, err := os.ReadFile("../schema/migrations/018_measurement_semantics_sqlite.up.sql")
	require.NoError(t, err)
	_, err = db.ExecContext(t.Context(), string(legacy))
	require.NoError(t, err)
	_, err = db.ExecContext(t.Context(), `INSERT INTO measurement_definitions
		(id,workspace_id,device_id,point_id,equipment_id,definition_revision,source_binding_revision,series_epoch,name,quantity,semantic_kind)
		VALUES ('legacy','ws','dev','pt','dev','rev-1','rev-1','epoch-1','Legacy','energy','counter')`)
	require.NoError(t, err)
	migration, err := os.ReadFile("../schema/migrations/039_measurement_exact_recording_scale_sqlite.up.sql")
	require.NoError(t, err)
	_, err = db.ExecContext(t.Context(), string(migration))
	require.NoError(t, err)
	got, err := NewSQLRepository(db).GetByID(t.Context(), "legacy")
	require.NoError(t, err)
	require.Empty(t, got.ExactRecordingScale)
}

func TestScaleExactInteger(t *testing.T) {
	cases := []struct {
		name  string
		value ExactValue
		scale string
		want  string
	}{
		{"120000 at 0.1", NewInt64(120000), "0.1", "12000.0"},
		{"120004 at 0.1", NewInt64(120004), "0.1", "12000.4"},
		{"negative signed", NewInt64(-120004), "0.1", "-12000.4"},
		{"three fractional digits", NewInt64(120000), "0.001", "120.000"},
		{"uint64 max at 0.1", NewUint64(^uint64(0)), "0.1", "1844674407370955161.5"},
	}
	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			got, err := ScaleExactInteger(testCase.value, testCase.scale)
			require.NoError(t, err)
			require.Equal(t, ExactDecimal, got.Type())
			digits, ok := got.DecimalDigits()
			require.True(t, ok)
			require.Equal(t, testCase.want, digits)
		})
	}

	for name, testCase := range map[string]struct {
		value ExactValue
		scale string
	}{
		"decimal value":  {mustDecimal(t, "120000"), "0.1"},
		"float value":    {mustFloat(t, 120000), "0.1"},
		"text value":     {NewText("120000"), "0.1"},
		"zero scale":     {NewInt64(120000), "0"},
		"negative scale": {NewInt64(120000), "-0.1"},
		"exponent scale": {NewInt64(120000), "1e-1"},
		"invalid scale":  {NewInt64(120000), "abc"},
	} {
		t.Run("rejects "+name, func(t *testing.T) {
			_, err := ScaleExactInteger(testCase.value, testCase.scale)
			require.Error(t, err)
		})
	}
}

func TestExactRecordingScaleMigrationCompanions(t *testing.T) {
	for _, name := range []string{
		"039_measurement_exact_recording_scale.up.sql",
		"039_measurement_exact_recording_scale.down.sql",
		"039_measurement_exact_recording_scale_sqlite.up.sql",
		"039_measurement_exact_recording_scale_sqlite.down.sql",
	} {
		content, err := migrations.FS.ReadFile(name)
		require.NoError(t, err, name)
		require.NotEmpty(t, content, name)
	}
	db, err := sql.Open("sqlite", t.TempDir()+"/down.db")
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })
	base, err := migrations.FS.ReadFile("018_measurement_semantics_sqlite.up.sql")
	require.NoError(t, err)
	_, err = db.ExecContext(t.Context(), string(base))
	require.NoError(t, err)
	up, err := migrations.FS.ReadFile("039_measurement_exact_recording_scale_sqlite.up.sql")
	require.NoError(t, err)
	_, err = db.ExecContext(t.Context(), string(up))
	require.NoError(t, err)
	down, err := migrations.FS.ReadFile("039_measurement_exact_recording_scale_sqlite.down.sql")
	require.NoError(t, err)
	_, err = db.ExecContext(t.Context(), string(down))
	require.NoError(t, err)
	var count int
	require.NoError(t, db.QueryRowContext(t.Context(), `SELECT COUNT(*) FROM pragma_table_info('measurement_definitions') WHERE name='exact_recording_scale'`).Scan(&count))
	require.Zero(t, count)
}
