package workspace

import (
	"testing"

	"go-gateway/internal/datalink/schema"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

// seedConfirmedCounter seeds a measurement definition whose saved evidence
// confirms counter semantics through the existing ConfirmCounter rules.
func seedConfirmedCounter(t *testing.T, h *circuitProfileHarness, source circuitSource) string {
	t.Helper()
	id := uuid.NewString()
	_, err := h.db.ExecContext(t.Context(), `INSERT INTO measurement_definitions (
		id, workspace_id, device_id, point_id, tag_id, equipment_id,
		definition_revision, source_binding_revision, series_epoch,
		name, quantity, unit, semantic_kind, numeric_encoding, counter_policy
	) VALUES (?, ?, ?, ?, ?, ?, 'definition-1', 'source-1', 'epoch-1',
		'Energy', 'energy', 'kWh', 'counter', 'uint32', '{"allow_negative":false}')`,
		id, h.fixture.workspaceID, source.deviceID, source.pointID, source.tagID, source.deviceID)
	require.NoError(t, err)
	return id
}

func seedBrokenCounterPolicy(t *testing.T, h *circuitProfileHarness, source circuitSource) string {
	t.Helper()
	id := uuid.NewString()
	_, err := h.db.ExecContext(t.Context(), `INSERT INTO measurement_definitions (
		id, workspace_id, device_id, point_id, tag_id, equipment_id,
		definition_revision, source_binding_revision, series_epoch,
		name, quantity, unit, semantic_kind, numeric_encoding, counter_policy
	) VALUES (?, ?, ?, ?, ?, ?, 'definition-1', 'source-1', 'epoch-1',
		'Energy', 'energy', 'kWh', 'counter', 'uint32', 'not-json')`,
		id, h.fixture.workspaceID, source.deviceID, source.pointID, source.tagID, source.deviceID)
	require.NoError(t, err)
	return id
}

func counterRawDraft(t *testing.T, h *circuitProfileHarness, circuit string, source circuitSource, measurementID *string) *WriteGroup {
	draft := h.draft(WriteGroupProfileInstantaneous, circuit, "Counter", source)
	draft.RecordingProfile = meanInstantaneousProfile()
	draft.Members[0] = memberWithJSON(t, draft.Members[0], `{"value_role":"counter_raw"}`)
	draft.Members[0].MeasurementID = measurementID
	return draft
}

// TestRecordingCounterBindingValidation pins the confirmed counter binding of
// a counter_raw member: an empty or blank measurement id is not a binding, a
// bool or text source is never a counter, and a broken saved definition is an
// internal availability failure instead of an operator validation refusal.
func TestRecordingCounterBindingValidation(t *testing.T) {
	t.Run("an empty or blank measurement id cannot save counter_raw", func(t *testing.T) {
		for name, id := range map[string]*string{"empty": strPtr(""), "blank": strPtr("   ")} {
			h := newCircuitProfileHarness(t, ":memory:")
			_, err := h.ensure(counterRawDraft(t, h, "circuit-blank", h.primary(), id))
			require.Error(t, err, "%s measurement id is not a confirmed binding", name)
			require.ErrorIs(t, err, ErrWriteGroupValidation, name)
			require.Zero(t, h.countRows(`SELECT COUNT(*) FROM write_groups`), name)
		}
	})

	t.Run("counter_raw on a bool tag is rejected even with a confirmed definition", func(t *testing.T) {
		h := newCircuitProfileHarness(t, ":memory:")
		boolSource := seedTypedSource(t, h, "state-bit", schema.DataTypeBool)
		measurementID := seedConfirmedCounter(t, h, boolSource)
		_, err := h.ensure(counterRawDraft(t, h, "circuit-bool-counter", boolSource, &measurementID))
		require.Error(t, err, "a bool source cannot feed a counter regardless of the saved definition")
		require.ErrorIs(t, err, ErrWriteGroupValidation)
	})

	t.Run("a broken saved definition is an availability failure", func(t *testing.T) {
		h := newCircuitProfileHarness(t, ":memory:")
		measurementID := seedBrokenCounterPolicy(t, h, h.primary())
		_, err := h.ensure(counterRawDraft(t, h, "circuit-broken", h.primary(), &measurementID))
		require.Error(t, err)
		require.NotErrorIs(t, err, ErrWriteGroupValidation,
			"an unreadable saved definition is an internal failure, not an operator mistake")
	})
}

// TestRecordingValueRoleRequired pins that an explicit value role is part of
// each new instantaneous shape: a mean member without one is refused, a role
// without a recording profile is refused, and legacy members stay omitted.
func TestRecordingValueRoleRequired(t *testing.T) {
	t.Run("a mean member without an explicit role is rejected", func(t *testing.T) {
		h := newCircuitProfileHarness(t, ":memory:")
		draft := h.draft(WriteGroupProfileInstantaneous, "circuit-norole", "NoRole", h.primary())
		draft.RecordingProfile = meanInstantaneousProfile()
		_, err := h.ensure(draft)
		require.Error(t, err, "a mean member must declare its role explicitly")
		require.ErrorIs(t, err, ErrWriteGroupValidation)
		require.Zero(t, h.countRows(`SELECT COUNT(*) FROM write_groups`))
	})

	t.Run("a role without a recording profile is rejected", func(t *testing.T) {
		h := newCircuitProfileHarness(t, ":memory:")
		draft := h.draft(WriteGroupProfileInstantaneous, "circuit-orphan", "Orphan", h.primary())
		draft.RecordingProfile = nil
		draft.Members = []WriteGroupMember{{DeviceID: h.fixture.deviceID, PointID: h.fixture.pointID, TagID: h.fixture.tagID, TargetColumn: "value"}}
		draft.Members[0] = memberWithJSON(t, draft.Members[0], `{"value_role":"gauge"}`)
		_, err := h.service.Create(t.Context(), h.mutation(t.Context(), draft))
		require.Error(t, err, "roles record instantaneous profiles only")
		require.ErrorIs(t, err, ErrWriteGroupValidation)
	})
}

// newNewWideHarness pins the fixture connector to its main SQLite scope: the
// preparation entry inspects it, so the saved group carries it and the scope
// call reaches the layout decision instead of a schema-name refusal.
func newNewWideHarness(t *testing.T) *circuitProfileHarness {
	t.Helper()
	h := newCircuitProfileHarness(t, ":memory:")
	_, err := h.db.ExecContext(t.Context(), `UPDATE database_connectors SET connection_config = ? WHERE id = ?`,
		`{"database":"persisted-db","schema":"main"}`, h.fixture.connectorID)
	require.NoError(t, err)
	return h
}

// TestRecordingNewWideProfileLifecycle pins the new wide profile layout end to
// end on the real SQL path: it saves, reloads and replays to the same single
// group, and its schema preparation follows the frozen result types: a gauge is a
// nullable float, a confirmed cumulative counter keeps its raw integer, and an
// unconfirmed one is refused instead of falling back to the legacy wide DDL.
func TestRecordingNewWideProfileLifecycle(t *testing.T) {
	newWideDraft := func(h *circuitProfileHarness) *WriteGroup {
		draft := h.draft(WriteGroupProfileInstantaneous, "circuit-newwide", "NewWide", h.primary())
		draft.LayoutVersion = WriteGroupLayoutManagedWideProfileV1
		draft.Destination.TableName = ""
		draft.Members[0] = memberWithJSON(t, draft.Members[0], `{"value_role":"gauge"}`)
		return draft
	}

	t.Run("a new wide profile group saves, reloads and replays to one group", func(t *testing.T) {
		h := newNewWideHarness(t)
		created := h.mustEnsure(newWideDraft(h))
		require.Equal(t, WriteGroupLayoutManagedWideProfileV1, created.EffectiveLayoutVersion())
		wire := memberWire(t, created.Members[0])
		require.Equal(t, "float32", wire["source_value_type"], "the new wide layout is a recording shape")

		reloaded, err := h.service.Get(t.Context(), created.ID)
		require.NoError(t, err)
		require.Equal(t, WriteGroupLayoutManagedWideProfileV1, reloaded.Group.EffectiveLayoutVersion())
		require.Equal(t, wire, memberWire(t, reloaded.Group.Members[0]))

		again := h.mustEnsure(newWideDraft(h))
		require.Equal(t, created.ID, again.ID, "a resend resolves to the saved group")
		require.Equal(t, 1, h.countRows(`SELECT COUNT(*) FROM write_groups`))
	})

	t.Run("a new wide cumulative profile keeps raw types without a role", func(t *testing.T) {
		h := newNewWideHarness(t)
		intSource := seedTypedSource(t, h, "energy64", schema.DataTypeInt64)
		measurementID := seedConfirmedCounter(t, h, intSource)
		draft := h.draft(WriteGroupProfileCumulative, "circuit-wide-cum", "WideCum", intSource)
		draft.LayoutVersion = WriteGroupLayoutManagedWideProfileV1
		draft.Destination.TableName = ""
		draft.Members[0].MeasurementID = &measurementID

		created := h.mustEnsure(draft)
		require.Equal(t, WriteGroupLayoutManagedWideProfileV1, created.EffectiveLayoutVersion())
		wire := memberWire(t, created.Members[0])
		require.NotContains(t, wire, "value_role", "a cumulative member carries no role")
		require.Equal(t, "int64", wire["source_value_type"], "the raw tag type is verified and saved")
		require.Equal(t, "int64", wire["result_value_type"], "a cumulative profile keeps the raw mapping")

		reloaded, err := h.service.Get(t.Context(), created.ID)
		require.NoError(t, err)
		require.Equal(t, wire, memberWire(t, reloaded.Group.Members[0]))

		again := h.mustEnsure(draft)
		require.Equal(t, created.ID, again.ID)
		require.Equal(t, 1, h.countRows(`SELECT COUNT(*) FROM write_groups`))

		gated := finalizedLong(reloaded.Group)
		require.False(t, gated.Ready, "the new wide profile stays unwritable without its runtime proof")
		require.Contains(t, issueCodesOf(gated), "recording-reduction-unavailable")
	})

	t.Run("a verified shape is prepared from its result types and an unserved one is refused", func(t *testing.T) {
		h := newNewWideHarness(t)
		scopeOf := func(group *WriteGroup) (columnTypes map[string]string, err error) {
			current, err := h.workspace.GetOrCreate(t.Context())
			require.NoError(t, err)
			scope, err := h.service.ManagedSchemaScope(t.Context(), group.ID, WriteGroupMutation{
				WorkspaceID: h.fixture.workspaceID, ExpectedWorkspaceRevision: current.DatabaseSetupRevision,
				ExpectedGroupRevision: group.Revision, ExpectedConnectorRevision: "connector-revision-1",
			})
			if err != nil {
				return nil, err
			}
			columnTypes = map[string]string{}
			for _, column := range scope.GroupLayout.Columns {
				if column.Nullable {
					columnTypes[column.Name] = column.SQLType
				}
			}
			return columnTypes, nil
		}

		gauge := h.mustEnsure(newWideDraft(h))
		types, err := scopeOf(gauge)
		require.NoError(t, err, "an instantaneous gauge is prepared from its frozen result type")
		require.Equal(t, map[string]string{gauge.Members[0].TargetColumn: "REAL"}, types,
			"the member column is typed by the frozen float64 result and may hold NULL")

		intSource := seedTypedSource(t, h, "energy64", schema.DataTypeInt64)
		measurementID := seedConfirmedCounter(t, h, intSource)
		cumulative := h.draft(WriteGroupProfileCumulative, "circuit-wide-cum", "WideCum", intSource)
		cumulative.LayoutVersion = WriteGroupLayoutManagedWideProfileV1
		cumulative.Destination.TableName = ""
		cumulative.Members[0].MeasurementID = &measurementID
		saved := h.mustEnsure(cumulative)
		types, err = scopeOf(saved)
		require.NoError(t, err, "a confirmed cumulative wide profile is prepared from its raw result type")
		require.Equal(t, map[string]string{saved.Members[0].TargetColumn: "INTEGER"}, types,
			"the member column keeps the raw int64 reading and may hold NULL; it is never a decimal or the legacy DDL")

		unconfirmed := h.draft(WriteGroupProfileCumulative, "circuit-wide-unconfirmed", "WideCumNone", intSource)
		unconfirmed.LayoutVersion = WriteGroupLayoutManagedWideProfileV1
		unconfirmed.Destination.TableName = ""
		_, err = h.service.Create(t.Context(), h.mutation(t.Context(), unconfirmed))
		require.Error(t, err, "a counter without its saved measurement is never saved, so it is never prepared")
		require.ErrorIs(t, err, ErrWriteGroupValidation)
	})
}

// TestRecordingCounterAndStatePositive pins the confirmed positive shapes: a
// confirmed counter keeps its raw result type beside an independent cumulative
// owner, and bool and text states save in a mean profile with latest results.
func TestRecordingCounterAndStatePositive(t *testing.T) {
	t.Run("a confirmed counter keeps its raw result beside a cumulative owner", func(t *testing.T) {
		h := newCircuitProfileHarness(t, ":memory:")
		intSource := seedTypedSource(t, h, "energy16", schema.DataTypeInt16)
		measurementID := seedConfirmedCounter(t, h, intSource)

		instantaneous := counterRawDraft(t, h, "circuit-energy", intSource, &measurementID)
		created := h.mustEnsure(instantaneous)
		wire := memberWire(t, created.Members[0])
		require.Equal(t, "int16", wire["source_value_type"])
		require.Equal(t, "int64", wire["result_value_type"],
			"counter_raw keeps the raw value: an integer tag maps to int64, never decimal")

		cumulative := h.draft(WriteGroupProfileCumulative, "circuit-energy", "Counter", intSource)
		cumulative.Members[0].MeasurementID = &measurementID
		cumulativeCreated := h.mustEnsure(cumulative)
		require.NotEqual(t, created.ID, cumulativeCreated.ID,
			"the two kinds are different canonical owners of the same source")
		require.Equal(t, 2, h.countRows(`SELECT COUNT(*) FROM write_groups`))

		require.Equal(t, created.ID, h.mustEnsure(instantaneous).ID)
		require.Equal(t, cumulativeCreated.ID, h.mustEnsure(cumulative).ID)
		require.Equal(t, 2, h.countRows(`SELECT COUNT(*) FROM write_groups`))
	})

	t.Run("bool and text states save in a mean profile", func(t *testing.T) {
		h := newCircuitProfileHarness(t, ":memory:")
		boolSource := seedTypedSource(t, h, "pump-state", schema.DataTypeBool)
		textSource := seedTypedSource(t, h, "mode-text", schema.DataTypeString)
		draft := h.draft(WriteGroupProfileInstantaneous, "circuit-states", "States", boolSource, textSource)
		draft.RecordingProfile = meanInstantaneousProfile()
		for i := range draft.Members {
			draft.Members[i] = memberWithJSON(t, draft.Members[i], `{"value_role":"state"}`)
		}
		created := h.mustEnsure(draft)
		require.Len(t, created.Members, 2)
		require.Equal(t, "bool", memberWire(t, created.Members[0])["source_value_type"])
		require.Equal(t, "bool", memberWire(t, created.Members[0])["result_value_type"], "a bool state keeps latest")
		require.Equal(t, "string", memberWire(t, created.Members[1])["source_value_type"])
		require.Equal(t, "text", memberWire(t, created.Members[1])["result_value_type"], "a text state keeps latest")
	})
}

func strPtr(value string) *string { return &value }

// TestRecordingMetadataRequiresProfile pins that new recording metadata never
// attaches to a group without a recording profile: the rejection is typed and
// nothing is persisted, instead of a nil-profile panic or silent adoption.
func TestRecordingMetadataRequiresProfile(t *testing.T) {
	h := newCircuitProfileHarness(t, ":memory:")
	// A complete legacy custom group without a recording profile; the only
	// difference is the member's new source type, so the refusal can only come
	// from the metadata shape guard itself.
	draft := newBasicWriteGroup(h.fixture)
	draft.Members[0] = memberWithJSON(t, draft.Members[0], `{"source_value_type":"float32"}`)

	_, err := h.service.Create(t.Context(), h.mutation(t.Context(), draft))
	require.Error(t, err, "recording metadata requires a recording profile")
	require.ErrorIs(t, err, ErrWriteGroupValidation)
	require.Zero(t, h.countRows(`SELECT COUNT(*) FROM write_groups`))
}
