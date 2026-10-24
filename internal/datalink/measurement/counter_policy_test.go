package measurement

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func counterDefinition(policy *CounterPolicy) *MeasurementDefinition {
	return &MeasurementDefinition{
		ID: "meas-energy", WorkspaceID: "ws", DeviceID: "dev", PointID: "pt", Name: "Energy", Quantity: "energy", Unit: "kWh",
		SemanticKind: SemanticKindCounter, NumericEncoding: "uint16", CounterPolicy: policy,
	}
}

// TestCounterPolicyExactRolloverIsValidatedWhereItIsSaved pins the contract of
// the exact rollover evidence: positive exact decimal text, a monotonic counter
// only, never a float, and a rate only together with the modulus it protects.
func TestCounterPolicyExactRolloverIsValidatedWhereItIsSaved(t *testing.T) {
	for name, policy := range map[string]*CounterPolicy{
		"modulus alone":                 {RolloverModulus: "65536"},
		"modulus and rate":              {RolloverModulus: "65536", MaxRatePerSecond: "100"},
		"decimal rate":                  {RolloverModulus: "65536", MaxRatePerSecond: "0.125"},
		"a modulus above 2^64":          {RolloverModulus: "18446744073709551616", MaxRatePerSecond: "1000000000000000"},
		"scale is preserved":            {RolloverModulus: "65536.000", MaxRatePerSecond: "100.0"},
		"legacy float fields untouched": {MaxRollOver: ptrFloat(65536), ResetThreshold: ptrFloat(1)},
		"legacy and exact side by side": {MaxRollOver: ptrFloat(65536), RolloverModulus: "65536", MaxRatePerSecond: "100"},
	} {
		t.Run("accepts "+name, func(t *testing.T) {
			require.NoError(t, counterDefinition(policy).Validate())
		})
	}

	for name, policy := range map[string]*CounterPolicy{
		"a zero modulus":             {RolloverModulus: "0"},
		"a zero decimal modulus":     {RolloverModulus: "0.000"},
		"a negative modulus":         {RolloverModulus: "-65536"},
		"a negative rate":            {RolloverModulus: "65536", MaxRatePerSecond: "-1"},
		"a zero rate":                {RolloverModulus: "65536", MaxRatePerSecond: "0"},
		"an exponent":                {RolloverModulus: "6.5536e4"},
		"hex":                        {RolloverModulus: "0x10000"},
		"NaN":                        {RolloverModulus: "NaN"},
		"infinity":                   {RolloverModulus: "Infinity"},
		"surrounding spaces":         {RolloverModulus: " 65536"},
		"trailing spaces":            {RolloverModulus: "65536 "},
		"a thousands separator":      {RolloverModulus: "65,536"},
		"a rate without the modulus": {MaxRatePerSecond: "100"},
		"an over-long digit string":  {RolloverModulus: strings.Repeat("9", 65)},
		"an invalid rate string":     {RolloverModulus: "65536", MaxRatePerSecond: "fast"},
		"signed negatives allowed":   {RolloverModulus: "65536", AllowNegative: true},
	} {
		t.Run("rejects "+name, func(t *testing.T) {
			require.Error(t, counterDefinition(policy).Validate())
		})
	}

	policy := &CounterPolicy{RolloverModulus: "65536", MaxRatePerSecond: "100"}
	for _, kind := range []SemanticKind{SemanticKindGauge, SemanticKindRate, SemanticKindSignedCounter, SemanticKindDelta, SemanticKindState, SemanticKindEvent, SemanticKindText} {
		def := counterDefinition(policy)
		def.SemanticKind = kind
		require.Error(t, def.Validate(), "rollover belongs to a monotonic counter, not %s", kind)
	}
	for _, encoding := range []string{"float64", "Float32", "float"} {
		def := counterDefinition(policy)
		def.NumericEncoding = encoding
		require.Error(t, def.Validate(), "a float source cannot prove a wrap: %s", encoding)
	}
	def := counterDefinition(policy)
	def.NumericEncoding = ""
	require.NoError(t, def.Validate(), "an unspecified encoding is not a float encoding")
}

func ptrFloat(v float64) *float64 { return &v }

// TestCounterPolicyJSONNeverTakesANumberForTheExactEvidence keeps the exact
// policy out of float64: a JSON number is refused rather than coerced, and the
// legacy float fields keep their meaning.
func TestCounterPolicyJSONNeverTakesANumberForTheExactEvidence(t *testing.T) {
	for _, body := range []string{
		`{"allow_negative":false,"rollover_modulus":65536}`,
		`{"allow_negative":false,"rollover_modulus":"65536","max_rate_per_second":100}`,
		`{"allow_negative":false,"rollover_modulus":9007199254740993}`,
	} {
		var policy CounterPolicy
		require.Error(t, json.Unmarshal([]byte(body), &policy), body)
	}

	var exact CounterPolicy
	require.NoError(t, json.Unmarshal([]byte(`{"allow_negative":false,"rollover_modulus":"9007199254740993","max_rate_per_second":"0.125"}`), &exact))
	require.Equal(t, "9007199254740993", exact.RolloverModulus, "digits above 2^53 are kept as text")
	require.Equal(t, "0.125", exact.MaxRatePerSecond)

	var legacy CounterPolicy
	require.NoError(t, json.Unmarshal([]byte(`{"allow_negative":false,"max_roll_over":65536,"reset_threshold":1}`), &legacy))
	require.NotNil(t, legacy.MaxRollOver)
	require.Equal(t, 65536.0, *legacy.MaxRollOver)
	require.Empty(t, legacy.RolloverModulus, "a saved float MaxRollOver is never converted into exact evidence")
	require.Empty(t, legacy.MaxRatePerSecond)
	encoded, err := json.Marshal(legacy)
	require.NoError(t, err)
	require.NotContains(t, string(encoded), "rollover_modulus", "and it round-trips as it was saved")
}

// TestCounterPolicyRevisionStartsANewEpoch pins that changed rollover evidence
// breaks comparability, while the same number written with another scale does
// not.
func TestCounterPolicyRevisionStartsANewEpoch(t *testing.T) {
	base := func() *MeasurementDefinition {
		return counterDefinition(&CounterPolicy{RolloverModulus: "65536", MaxRatePerSecond: "100"})
	}
	changed := func(mutate func(*CounterPolicy)) bool {
		next := base()
		mutate(next.CounterPolicy)
		changed, epoch, reason := CheckEpochTransition(*base(), *next)
		if changed {
			require.Equal(t, "epoch-2", epoch)
			require.NotEmpty(t, reason)
		}
		return changed
	}
	require.True(t, changed(func(p *CounterPolicy) { p.RolloverModulus = "32768" }), "a different modulus")
	require.True(t, changed(func(p *CounterPolicy) { p.MaxRatePerSecond = "50" }), "a different rate bound")
	require.True(t, changed(func(p *CounterPolicy) { p.MaxRatePerSecond = "" }), "a removed rate bound")
	require.True(t, changed(func(p *CounterPolicy) { p.RolloverModulus, p.MaxRatePerSecond = "", "" }), "a removed policy")
	require.False(t, changed(func(p *CounterPolicy) { p.RolloverModulus = "65536.00" }), "the same modulus with another scale is the same evidence")
	require.False(t, changed(func(p *CounterPolicy) { p.MaxRatePerSecond = "100.000" }), "the same rate with another scale is the same evidence")
	require.False(t, changed(func(p *CounterPolicy) {}), "nothing changed")

	plain := counterDefinition(&CounterPolicy{})
	withPolicy := counterDefinition(&CounterPolicy{RolloverModulus: "65536"})
	moved, _, _ := CheckEpochTransition(*plain, *withPolicy)
	require.True(t, moved, "gaining rollover evidence changes what a decrease means")
}

// TestCounterPolicyRoundTripsThroughTheSavedDefinition saves a definition with
// exact evidence beyond 2^53 and a decimal rate through the service and the SQL
// repository, reads it back, and revises it: the digits survive untouched and a
// changed policy moves the series to a new epoch and revision.
func TestCounterPolicyRoundTripsThroughTheSavedDefinition(t *testing.T) {
	db := measurementRepositoryDB(t)
	service := NewService(NewSQLRepository(db))
	policy := &CounterPolicy{MaxRollOver: ptrFloat(65536), RolloverModulus: "9007199254740993", MaxRatePerSecond: "0.125"}
	def := counterDefinition(policy)
	require.NoError(t, service.CreateMeasurement(t.Context(), def))

	saved, err := service.GetMeasurement(t.Context(), def.ID)
	require.NoError(t, err)
	require.NotNil(t, saved.CounterPolicy)
	require.Equal(t, "9007199254740993", saved.CounterPolicy.RolloverModulus, "no float64 detour")
	require.Equal(t, "0.125", saved.CounterPolicy.MaxRatePerSecond)
	require.Equal(t, 65536.0, *saved.CounterPolicy.MaxRollOver, "the legacy float field is kept as saved")
	require.Equal(t, "epoch-1", saved.SeriesEpoch)
	var stored string
	require.NoError(t, db.QueryRowContext(t.Context(), `SELECT counter_policy FROM measurement_definitions WHERE id = ?`, def.ID).Scan(&stored))
	require.Contains(t, stored, `"rollover_modulus":"9007199254740993"`, "stored as text in the existing counter_policy column")

	same := counterDefinition(&CounterPolicy{RolloverModulus: "9007199254740993.00", MaxRatePerSecond: "0.1250", MaxRollOver: ptrFloat(65536)})
	same.ID, same.WorkspaceID = def.ID, def.WorkspaceID
	moved, err := service.UpdateMeasurement(t.Context(), same)
	require.NoError(t, err)
	require.False(t, moved, "the same numbers written with another scale keep the series")

	revised := counterDefinition(&CounterPolicy{RolloverModulus: "65536", MaxRatePerSecond: "100"})
	revised.ID, revised.WorkspaceID = def.ID, def.WorkspaceID
	moved, err = service.UpdateMeasurement(t.Context(), revised)
	require.NoError(t, err)
	require.True(t, moved, "changed evidence is a new series epoch")
	after, err := service.GetMeasurement(t.Context(), def.ID)
	require.NoError(t, err)
	require.Equal(t, "epoch-2", after.SeriesEpoch)
	require.Equal(t, "rev-3", after.DefinitionRevision)
	require.Equal(t, "65536", after.CounterPolicy.RolloverModulus)

	invalid := counterDefinition(&CounterPolicy{RolloverModulus: "0"})
	invalid.ID, invalid.WorkspaceID = def.ID, def.WorkspaceID
	_, err = service.UpdateMeasurement(t.Context(), invalid)
	require.ErrorIs(t, err, ErrCounterPolicyInvalid, "an invalid policy is refused where it is saved")
	unchanged, err := service.GetMeasurement(t.Context(), def.ID)
	require.NoError(t, err)
	require.Equal(t, "65536", unchanged.CounterPolicy.RolloverModulus)
}

func TestConfirmCounterIsSingleSourceOfTruthForCounterReadinessAndRuntime(t *testing.T) {
	confirmed := func() MeasurementDefinition {
		def := *counterDefinition(&CounterPolicy{RolloverModulus: "65536", MaxRatePerSecond: "100"})
		def.DefinitionRevision = "def-1"
		def.SourceBindingRevision = "bind-1"
		def.SeriesEpoch = "epoch-1"
		return def
	}

	require.Equal(t, CounterConfirmed, ConfirmCounter(confirmed()))

	signed := confirmed()
	signed.SemanticKind = SemanticKindSignedCounter
	signed.CounterPolicy = &CounterPolicy{AllowNegative: true}
	require.Equal(t, CounterConfirmed, ConfirmCounter(signed))

	for _, kind := range []SemanticKind{SemanticKindGauge, SemanticKindRate, SemanticKindDelta, SemanticKindState, SemanticKindEvent, SemanticKindText} {
		def := confirmed()
		def.SemanticKind = kind
		require.Equal(t, CounterNotCounter, ConfirmCounter(def), "kind=%s", kind)
	}

	unconfirmedCases := map[string]func(*MeasurementDefinition){
		"monotonic counter allowing negative": func(d *MeasurementDefinition) { d.CounterPolicy = &CounterPolicy{AllowNegative: true} },
		"signed counter disallowing negative": func(d *MeasurementDefinition) {
			d.SemanticKind = SemanticKindSignedCounter
			d.CounterPolicy = &CounterPolicy{AllowNegative: false}
		},
		"signed counter without policy":   func(d *MeasurementDefinition) { d.SemanticKind = SemanticKindSignedCounter; d.CounterPolicy = nil },
		"missing id":                      func(d *MeasurementDefinition) { d.ID = " " },
		"missing definition revision":     func(d *MeasurementDefinition) { d.DefinitionRevision = "" },
		"missing source binding revision": func(d *MeasurementDefinition) { d.SourceBindingRevision = "" },
		"missing series epoch":            func(d *MeasurementDefinition) { d.SeriesEpoch = "" },
		"missing unit":                    func(d *MeasurementDefinition) { d.Unit = "" },
		"invalid rollover policy":         func(d *MeasurementDefinition) { d.CounterPolicy = &CounterPolicy{RolloverModulus: "0"} },
		"float encoding with rollover":    func(d *MeasurementDefinition) { d.NumericEncoding = "float64" },
	}
	for name, mutate := range unconfirmedCases {
		t.Run(name, func(t *testing.T) {
			def := confirmed()
			mutate(&def)
			require.Equal(t, CounterUnconfirmed, ConfirmCounter(def))
		})
	}
}
