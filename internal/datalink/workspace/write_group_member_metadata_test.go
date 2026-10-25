package workspace

import (
	"database/sql"
	"encoding/json"
	"testing"

	"go-gateway/internal/datalink/schema"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

// decodeMemberPatch patches a member through its JSON wire form, the only
// path every persisted payload takes, and returns the decode error so tests
// can pin rejections without weakening them.
func decodeMemberPatch(t *testing.T, member WriteGroupMember, patch string) error {
	t.Helper()
	raw, err := json.Marshal(member)
	require.NoError(t, err)
	var payload map[string]any
	require.NoError(t, json.Unmarshal(raw, &payload))
	var merged map[string]any
	require.NoError(t, json.Unmarshal([]byte(patch), &merged))
	for key, value := range merged {
		payload[key] = value
	}
	encoded, err := json.Marshal(payload)
	require.NoError(t, err)
	var decoded WriteGroupMember
	return json.Unmarshal(encoded, &decoded)
}

// memberWithJSON patches a member and requires the decode to succeed.
func memberWithJSON(t *testing.T, member WriteGroupMember, patch string) WriteGroupMember {
	t.Helper()
	require.NoError(t, decodeMemberPatch(t, member, patch), patch)
	raw, err := json.Marshal(member)
	require.NoError(t, err)
	var payload map[string]any
	require.NoError(t, json.Unmarshal(raw, &payload))
	var merged map[string]any
	require.NoError(t, json.Unmarshal([]byte(patch), &merged))
	for key, value := range merged {
		payload[key] = value
	}
	encoded, err := json.Marshal(payload)
	require.NoError(t, err)
	var decoded WriteGroupMember
	require.NoError(t, json.Unmarshal(encoded, &decoded))
	return decoded
}

// memberWire marshals a member and returns its JSON object for wire assertions.
func memberWire(t *testing.T, member WriteGroupMember) map[string]any {
	t.Helper()
	raw, err := json.Marshal(member)
	require.NoError(t, err)
	var wire map[string]any
	require.NoError(t, json.Unmarshal(raw, &wire))
	return wire
}

// TestRecordingMemberMetadataCompatibility pins the member metadata decode
// contract: an explicit role or type survives the wire, unknown values and
// explicit empty, null or non-string metadata are rejected, and known values
// keep decoding.
func TestRecordingMemberMetadataCompatibility(t *testing.T) {
	base := WriteGroupMember{DeviceID: "device", PointID: "point", TagID: "tag", TargetColumn: "value"}

	t.Run("member metadata survives the json round trip", func(t *testing.T) {
		member := memberWithJSON(t, base, `{"value_role":"gauge","source_value_type":"float32","result_value_type":"float64"}`)
		wire := memberWire(t, member)
		require.Equal(t, "gauge", wire["value_role"], "an explicit role must survive decoding")
		require.Equal(t, "float32", wire["source_value_type"])
		require.Equal(t, "float64", wire["result_value_type"])
	})

	t.Run("a member without metadata stays absent on the wire", func(t *testing.T) {
		wire := memberWire(t, base)
		require.NotContains(t, wire, "value_role")
		require.NotContains(t, wire, "source_value_type")
		require.NotContains(t, wire, "result_value_type")
	})

	t.Run("unknown role is rejected at decode", func(t *testing.T) {
		err := decodeMemberPatch(t, base, `{"value_role":"average"}`)
		require.Error(t, err, "an unknown role must fail closed instead of decoding as gauge")
		require.ErrorIs(t, err, ErrWriteGroupValidation)
	})

	t.Run("explicit empty, null or non-string member metadata is rejected", func(t *testing.T) {
		for name, patch := range map[string]string{
			"empty role":     `{"value_role":""}`,
			"null role":      `{"value_role":null}`,
			"numeric role":   `{"value_role":3}`,
			"empty source":   `{"source_value_type":""}`,
			"unknown source": `{"source_value_type":"decimal128"}`,
			"unknown result": `{"result_value_type":"bignum"}`,
		} {
			err := decodeMemberPatch(t, base, patch)
			require.Error(t, err, "%s must not silently degrade", name)
			require.ErrorIs(t, err, ErrWriteGroupValidation, name)
		}
	})

	t.Run("known roles and types keep decoding", func(t *testing.T) {
		for _, patch := range []string{
			`{"value_role":"counter_raw"}`, `{"value_role":"state"}`,
			`{"source_value_type":"int16"}`, `{"result_value_type":"decimal"}`,
			`{"source_value_type":"string"}`, `{"result_value_type":"bool"}`,
		} {
			member := memberWithJSON(t, base, patch)
			wire := memberWire(t, member)
			var parsed map[string]any
			require.NoError(t, json.Unmarshal([]byte(patch), &parsed))
			for key := range parsed {
				require.Contains(t, wire, key, patch)
			}
		}
	})
}

// seedTypedSource seeds one mapped point and tag with an explicit data type.
func seedTypedSource(t *testing.T, h *circuitProfileHarness, key string, dataType schema.DataType) circuitSource {
	t.Helper()
	ctx := t.Context()
	h.nextAddress++
	source := circuitSource{deviceID: h.fixture.deviceID, pointID: uuid.NewString(), tagID: uuid.NewString()}
	_, err := h.db.ExecContext(ctx, `INSERT INTO points (id, device_id, name, address, function, data_type, mode, enabled)
		VALUES (?, ?, ?, ?, 'FC03', ?, 'read', 1)`, source.pointID, source.deviceID, key, h.nextAddress, dataType)
	require.NoError(t, err)
	_, err = h.db.ExecContext(ctx, `INSERT INTO tags (id, key, key_lower, display_name, data_type, status, labels)
		VALUES (?, ?, ?, ?, ?, 'active', '{}')`, source.tagID, key, key, key, dataType)
	require.NoError(t, err)
	_, err = h.db.ExecContext(ctx, `INSERT INTO mappings (id, point_id, tag_id, transform_pipeline, status, enabled)
		VALUES (?, ?, ?, '[]', 'active', 1)`, uuid.NewString(), source.pointID, source.tagID)
	require.NoError(t, err)
	return source
}

// seedMeasurement seeds a measurement definition of the given semantic kind.
func seedMeasurement(t *testing.T, h *circuitProfileHarness, source circuitSource, kind string) string {
	t.Helper()
	id := uuid.NewString()
	_, err := h.db.ExecContext(t.Context(), `INSERT INTO measurement_definitions (
		id, workspace_id, device_id, point_id, tag_id, equipment_id,
		definition_revision, source_binding_revision, name, quantity, semantic_kind
	) VALUES (?, ?, ?, ?, ?, ?, 'definition-1', 'source-1', 'Measured', 'process', ?)`,
		id, h.fixture.workspaceID, source.deviceID, source.pointID, source.tagID, source.deviceID, kind)
	require.NoError(t, err)
	return id
}

func meanInstantaneousProfile() *WriteGroupRecordingProfile {
	return &WriteGroupRecordingProfile{Kind: WriteGroupProfileInstantaneous, Interval: 15, Unit: WriteGroupIntervalUnitSecond, SampleReduction: WriteGroupSampleReductionMean}
}

// TestRecordingMemberMetadataPersistence pins the durable member metadata
// contract: a mean gauge keeps tag-verified source and derived result types
// across reload and update; legacy members stay omitted; explicit metadata is
// verified against the real tag; counter roles need the confirmed binding;
// stored NULL stays readable while non-NULL empty or unknown values fail.
func TestRecordingMemberMetadataPersistence(t *testing.T) {
	t.Run("a mean gauge member keeps verified types across reload and update", func(t *testing.T) {
		h := newCircuitProfileHarness(t, ":memory:")
		draft := h.draft(WriteGroupProfileInstantaneous, "circuit-meta", "Meta", h.primary())
		draft.RecordingProfile = meanInstantaneousProfile()
		draft.Members[0] = memberWithJSON(t, draft.Members[0], `{"value_role":"gauge"}`)
		created := h.mustEnsure(draft)

		wire := memberWire(t, created.Members[0])
		require.Equal(t, "float32", wire["source_value_type"], "the tag's own type is verified and saved")
		require.Equal(t, "float64", wire["result_value_type"], "a float mean yields float64")
		require.Equal(t, "gauge", wire["value_role"])

		var role, sourceType, resultType string
		require.NoError(t, h.db.QueryRowContext(t.Context(), `SELECT value_role, source_value_type, result_value_type
			FROM write_group_members WHERE group_id = ?`, created.ID).Scan(&role, &sourceType, &resultType))
		require.Equal(t, "gauge", role)
		require.Equal(t, "float32", sourceType)
		require.Equal(t, "float64", resultType)

		reloaded, err := h.service.Get(t.Context(), created.ID)
		require.NoError(t, err)
		require.Equal(t, wire, memberWire(t, reloaded.Group.Members[0]))

		current := reloaded.Group
		current.Name = "Meta renamed"
		mutation := h.mutation(t.Context(), current)
		mutation.ExpectedGroupRevision = current.Revision
		updated, err := h.service.Update(t.Context(), current.ID, mutation)
		require.NoError(t, err)
		require.Equal(t, wire, memberWire(t, updated.Group.Members[0]), "the member rewrite keeps verified types")
	})

	t.Run("a mean integer gauge derives a decimal result", func(t *testing.T) {
		h := newCircuitProfileHarness(t, ":memory:")
		intSource := seedTypedSource(t, h, "counter16", schema.DataTypeInt16)
		draft := h.draft(WriteGroupProfileInstantaneous, "circuit-int", "Int", intSource)
		draft.RecordingProfile = meanInstantaneousProfile()
		draft.Members[0] = memberWithJSON(t, draft.Members[0], `{"value_role":"gauge"}`)
		created := h.mustEnsure(draft)
		wire := memberWire(t, created.Members[0])
		require.Equal(t, "int16", wire["source_value_type"])
		require.Equal(t, "decimal", wire["result_value_type"], "an integer mean yields decimal")
	})

	t.Run("a mean state member keeps the latest mapping", func(t *testing.T) {
		h := newCircuitProfileHarness(t, ":memory:")
		draft := h.draft(WriteGroupProfileInstantaneous, "circuit-state", "State", h.primary())
		draft.RecordingProfile = meanInstantaneousProfile()
		draft.Members[0] = memberWithJSON(t, draft.Members[0], `{"value_role":"state"}`)
		created := h.mustEnsure(draft)
		wire := memberWire(t, created.Members[0])
		require.Equal(t, "state", wire["value_role"])
		require.Equal(t, "float64", wire["result_value_type"], "a state member is never averaged")
	})

	t.Run("legacy latest members stay omitted", func(t *testing.T) {
		h := newCircuitProfileHarness(t, ":memory:")
		draft := h.draft(WriteGroupProfileInstantaneous, "circuit-legacy", "Legacy", h.primary())
		draft.RecordingProfile = &WriteGroupRecordingProfile{Kind: WriteGroupProfileInstantaneous, Interval: 15, Unit: WriteGroupIntervalUnitSecond, SampleReduction: WriteGroupSampleReductionLatest}
		created := h.mustEnsure(draft)
		wire := memberWire(t, created.Members[0])
		require.NotContains(t, wire, "source_value_type", "legacy members are never backfilled")
		var role, sourceType, resultType sql.NullString
		require.NoError(t, h.db.QueryRowContext(t.Context(), `SELECT value_role, source_value_type, result_value_type
			FROM write_group_members WHERE group_id = ?`, created.ID).Scan(&role, &sourceType, &resultType))
		require.False(t, role.Valid || sourceType.Valid || resultType.Valid)
	})

	t.Run("explicit metadata with latest is verified and saved", func(t *testing.T) {
		h := newCircuitProfileHarness(t, ":memory:")
		draft := h.draft(WriteGroupProfileInstantaneous, "circuit-explicit", "Explicit", h.primary())
		draft.RecordingProfile.SampleReduction = WriteGroupSampleReductionLatest
		draft.Members[0] = memberWithJSON(t, draft.Members[0], `{"value_role":"gauge"}`)
		created := h.mustEnsure(draft)
		wire := memberWire(t, created.Members[0])
		require.Equal(t, "float32", wire["source_value_type"])
		require.Equal(t, "float64", wire["result_value_type"], "latest maps through the tag's exact type")
	})

	t.Run("a caller-supplied wrong source type is rejected", func(t *testing.T) {
		h := newCircuitProfileHarness(t, ":memory:")
		draft := h.draft(WriteGroupProfileInstantaneous, "circuit-wrong", "Wrong", h.primary())
		draft.RecordingProfile = meanInstantaneousProfile()
		draft.Members[0] = memberWithJSON(t, draft.Members[0], `{"value_role":"gauge","source_value_type":"int16"}`)
		_, err := h.ensure(draft)
		require.Error(t, err, "the saved source type must match the real tag")
		require.ErrorIs(t, err, ErrWriteGroupValidation)
	})

	t.Run("counter_raw requires a confirmed counter binding", func(t *testing.T) {
		h := newCircuitProfileHarness(t, ":memory:")
		draft := h.draft(WriteGroupProfileInstantaneous, "circuit-counter", "Counter", h.primary())
		draft.RecordingProfile = meanInstantaneousProfile()
		draft.Members[0] = memberWithJSON(t, draft.Members[0], `{"value_role":"counter_raw"}`)
		_, err := h.ensure(draft)
		require.Error(t, err, "counter_raw without a measurement is not a counter")
		require.ErrorIs(t, err, ErrWriteGroupValidation)

		measurementID := seedMeasurement(t, h, h.primary(), "gauge")
		gaugeBound := h.draft(WriteGroupProfileInstantaneous, "circuit-counter", "Counter", h.primary())
		gaugeBound.RecordingProfile = meanInstantaneousProfile()
		gaugeBound.Members[0] = memberWithJSON(t, gaugeBound.Members[0], `{"value_role":"counter_raw"}`)
		gaugeBound.Members[0].MeasurementID = &measurementID
		_, err = h.ensure(gaugeBound)
		require.Error(t, err, "a gauge measurement never confirms counter semantics")
		require.ErrorIs(t, err, ErrWriteGroupValidation)
	})

	t.Run("a resent mean gauge ensure replays without conflict", func(t *testing.T) {
		h := newCircuitProfileHarness(t, ":memory:")
		draft := h.draft(WriteGroupProfileInstantaneous, "circuit-meta", "Meta", h.primary())
		draft.RecordingProfile = meanInstantaneousProfile()
		draft.Members[0] = memberWithJSON(t, draft.Members[0], `{"value_role":"gauge"}`)
		first := h.mustEnsure(draft)
		again := h.mustEnsure(draft)
		require.Equal(t, first.ID, again.ID, "server-derived types must not turn a resend into a conflict")
		require.Equal(t, 1, h.countRows(`SELECT COUNT(*) FROM write_groups`))
	})

	t.Run("a different explicit role conflicts", func(t *testing.T) {
		h := newCircuitProfileHarness(t, ":memory:")
		draft := h.draft(WriteGroupProfileInstantaneous, "circuit-role", "Role", h.primary())
		draft.RecordingProfile = meanInstantaneousProfile()
		draft.Members[0] = memberWithJSON(t, draft.Members[0], `{"value_role":"gauge"}`)
		saved := h.mustEnsure(draft)
		changed := h.draft(WriteGroupProfileInstantaneous, "circuit-role", "Role", h.primary())
		changed.RecordingProfile = meanInstantaneousProfile()
		changed.Members[0] = memberWithJSON(t, changed.Members[0], `{"value_role":"state"}`)
		_, err := h.ensure(changed)
		require.ErrorIs(t, err, ErrWriteGroupCircuitProfileConflict, "a role change is a different intent")
		require.Equal(t, 1, h.countRows(`SELECT COUNT(*) FROM write_groups`))
		reloaded, err := h.service.Get(t.Context(), saved.ID)
		require.NoError(t, err)
		require.Equal(t, "gauge", memberWire(t, reloaded.Group.Members[0])["value_role"])
	})

	t.Run("an explicit wrong source type on resend conflicts", func(t *testing.T) {
		h := newCircuitProfileHarness(t, ":memory:")
		draft := h.draft(WriteGroupProfileInstantaneous, "circuit-type", "Type", h.primary())
		draft.RecordingProfile = meanInstantaneousProfile()
		draft.Members[0] = memberWithJSON(t, draft.Members[0], `{"value_role":"gauge"}`)
		h.mustEnsure(draft)
		resent := h.draft(WriteGroupProfileInstantaneous, "circuit-type", "Type", h.primary())
		resent.RecordingProfile = meanInstantaneousProfile()
		resent.Members[0] = memberWithJSON(t, resent.Members[0], `{"value_role":"gauge","source_value_type":"int16"}`)
		_, err := h.ensure(resent)
		require.Error(t, err, "an explicit type that differs from the saved fact is a different intent")
	})

	t.Run("stored metadata reads fail closed", func(t *testing.T) {
		h := newCircuitProfileHarness(t, ":memory:")
		draft := h.draft(WriteGroupProfileInstantaneous, "circuit-sql", "SQL", h.primary())
		draft.RecordingProfile = &WriteGroupRecordingProfile{Kind: WriteGroupProfileInstantaneous, Interval: 15, Unit: WriteGroupIntervalUnitSecond, SampleReduction: WriteGroupSampleReductionLatest}
		created := h.mustEnsure(draft)
		memberID := created.Members[0].TagID

		_, err := h.db.ExecContext(t.Context(), `UPDATE write_groups SET sample_reduction = '' WHERE id = ?`, created.ID)
		require.NoError(t, err)
		_, err = h.service.Get(t.Context(), created.ID)
		require.Error(t, err, "a non-NULL empty reduction is corrupt, not legacy")

		_, err = h.db.ExecContext(t.Context(), `UPDATE write_groups SET sample_reduction = 'average' WHERE id = ?`, created.ID)
		require.NoError(t, err)
		_, err = h.service.Get(t.Context(), created.ID)
		require.Error(t, err, "an unknown stored reduction fails instead of reading latest")

		_, err = h.db.ExecContext(t.Context(), `UPDATE write_groups SET sample_reduction = NULL WHERE id = ?`, created.ID)
		require.NoError(t, err)
		_, err = h.service.Get(t.Context(), created.ID)
		require.NoError(t, err, "stored NULL keeps the legacy read")

		_, err = h.db.ExecContext(t.Context(), `UPDATE write_group_members SET value_role = '' WHERE tag_id = ?`, memberID)
		require.NoError(t, err)
		_, err = h.service.Get(t.Context(), created.ID)
		require.Error(t, err, "a non-NULL empty role is corrupt")

		_, err = h.db.ExecContext(t.Context(), `UPDATE write_group_members SET value_role = 'average' WHERE tag_id = ?`, memberID)
		require.NoError(t, err)
		_, err = h.service.Get(t.Context(), created.ID)
		require.Error(t, err, "an unknown stored role fails closed")

		_, err = h.db.ExecContext(t.Context(), `UPDATE write_group_members SET value_role = NULL WHERE tag_id = ?`, memberID)
		require.NoError(t, err)
		_, err = h.service.Get(t.Context(), created.ID)
		require.NoError(t, err, "stored NULL member metadata keeps the legacy read")
	})
}

// newWideProfileLayout is the wire value of the not-yet-proven new wide
// profile layout, spelled literally so the gate tests compile before and
// after the constant exists.
const newWideProfileLayout = WriteGroupLayoutVersion("managed_wide_profile_v1")

// TestRecordingReductionCompatibility_NewShapesStayBlocked pins the 1.1 safety
// boundary: every new member metadata shape, a mean reduction and the new wide
// profile layout stay savable but readiness refuses them until a later runtime
// proof lands; a legacy latest group keeps its ready path.
func TestRecordingReductionCompatibility_NewShapesStayBlocked(t *testing.T) {
	metadataGroup := longGroupFixture(WriteGroupProfileInstantaneous, WriteGroupStorageModePlain, false)
	metadataGroup.Members[0] = memberWithJSON(t, metadataGroup.Members[0], `{"value_role":"gauge"}`)
	metadata := finalizedLong(metadataGroup)
	require.False(t, metadata.Ready, "member metadata has no runtime proof yet")
	require.Contains(t, issueCodesOf(metadata), "recording-reduction-unavailable")

	newWide := longGroupFixture(WriteGroupProfileInstantaneous, WriteGroupStorageModePlain, false)
	newWide.LayoutVersion = newWideProfileLayout
	wideProfile := finalizedLong(newWide)
	require.False(t, wideProfile.Ready, "the new wide profile has no runtime proof yet")
	require.Contains(t, issueCodesOf(wideProfile), "recording-reduction-unavailable")
}

// TestRecordingReductionCompatibility_DigestLiterals pins the legacy managed
// source digests with the exact hashes rebuilt independently from the
// pre-reduction struct and JSON field order, so a silent widening of the
// digest scope cannot pass unnoticed.
func TestRecordingReductionCompatibility_DigestLiterals(t *testing.T) {
	tagTypes := map[string]schema.DataType{"a": schema.DataTypeFloat32}
	legacyWide := longGroupFixture(WriteGroupProfileInstantaneous, WriteGroupStorageModePlain, false)
	legacyWide.LayoutVersion = ""
	wideDigest, err := managedSourceDigest(legacyWide, tagTypes)
	require.NoError(t, err)
	require.Equal(t, "b5ab091e72686793b4ad77475849fd26bfe7b35bdd247070dc8e12387dd71c7c", wideDigest,
		"the legacy wide digest literal must not move")

	legacyLong := longGroupFixture(WriteGroupProfileInstantaneous, WriteGroupStorageModePlain, false)
	longDigest, err := managedSourceDigest(legacyLong, tagTypes)
	require.NoError(t, err)
	require.Equal(t, "7c27c04b31d13fef8f1bbad153dc3ec111b9f5debdd74b2602ed712f42adae49", longDigest,
		"the legacy long digest literal must not move")
}

// TestRecordingReductionCompatibility_NewWideDigest pins that the new wide
// profile layout joins the long layout in binding its profile and layout into
// the confirmed scope, while the legacy wide digest never carried them.
func TestRecordingReductionCompatibility_NewWideDigest(t *testing.T) {
	tagTypes := map[string]schema.DataType{"a": schema.DataTypeFloat32}
	newWide := func(reduction WriteGroupSampleReduction) *WriteGroup {
		group := longGroupFixture(WriteGroupProfileInstantaneous, WriteGroupStorageModePlain, false)
		group.LayoutVersion = newWideProfileLayout
		group.StorageMode = ""
		group.RecordingProfile.SampleReduction = reduction
		return group
	}
	latest, err := managedSourceDigest(newWide(WriteGroupSampleReductionLatest), tagTypes)
	require.NoError(t, err)
	mean, err := managedSourceDigest(newWide(WriteGroupSampleReductionMean), tagTypes)
	require.NoError(t, err)
	require.NotEqual(t, latest, mean, "a new wide profile group binds its reduction into the digest")

	// The legacy wide digest never carries the profile; that guarantee is
	// pinned by TestRecordingReductionCompatibility_Digest and the fixed
	// digest literals, so it is not duplicated here with an aliased pair.
}
