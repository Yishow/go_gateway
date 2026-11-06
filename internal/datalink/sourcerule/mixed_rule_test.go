package sourcerule

import (
	"encoding/json"
	"errors"
	"testing"
	"time"

	"go-gateway/internal/datalink/measurement"
	"go-gateway/internal/datalink/schema"

	"github.com/stretchr/testify/require"
)

func mixedItem(id, address string, dataType schema.DataType) measurement.MixedItem {
	kind := measurement.SemanticKindGauge
	if dataType == schema.DataTypeBool {
		kind = measurement.SemanticKindState
	}
	return measurement.MixedItem{ItemID: id, Address: address, DataType: dataType, SemanticKind: kind}
}

func requireMixedFieldError(t *testing.T, err error, field string, itemIDs ...string) *FieldError {
	t.Helper()
	require.ErrorIs(t, err, ErrValidation)
	var fieldErr *FieldError
	require.True(t, errors.As(err, &fieldErr), "expected FieldError, got %v", err)
	require.Equal(t, field, fieldErr.Field)
	require.Equal(t, itemIDs, fieldErr.ItemIDs)
	return fieldErr
}

func TestPlanMixedRuleItemsExpandsEachItemAcrossAreas(t *testing.T) {
	items := []measurement.MixedItem{
		mixedItem("d39", " d39 ", schema.DataTypeInt16),
		mixedItem("d100", "D100", schema.DataTypeInt32),
		mixedItem("m5", "M5", schema.DataTypeBool),
	}
	plan, err := planMixedRuleItems(items, schema.ProtocolMC3E)
	require.NoError(t, err)
	require.Equal(t, []plannedRulePoint{
		{ItemID: "d39", Address: "D39", DataType: schema.DataTypeInt16},
		{ItemID: "d100", Address: "D100", DataType: schema.DataTypeInt32},
		{ItemID: "m5", Address: "M5", DataType: schema.DataTypeBool},
	}, plan.points)
	var stored []measurement.MixedItem
	require.NoError(t, json.Unmarshal([]byte(plan.encoded), &stored))
	require.Equal(t, "D39", stored[0].Address, "stored items keep the normalized address")
}

func TestPlanMixedRuleItemsRejectsOverlapNamingBothItems(t *testing.T) {
	_, err := planMixedRuleItems([]measurement.MixedItem{
		mixedItem("d3", "D3", schema.DataTypeUint32),
		mixedItem("d4", "D4", schema.DataTypeInt16),
	}, schema.ProtocolMC3E)
	fieldErr := requireMixedFieldError(t, err, "mixed_items", "d3", "d4")
	require.Contains(t, fieldErr.Error(), "D3")
	require.Contains(t, fieldErr.Error(), "D4")
}

func TestPlanMixedRuleItemsRejectsInvalidItems(t *testing.T) {
	scale := 2.0
	scaled := mixedItem("d1", "D1", schema.DataTypeInt16)
	scaled.ScaleMultiplier = &scale
	noKind := mixedItem("d2", "D2", schema.DataTypeInt16)
	noKind.SemanticKind = ""
	cases := []struct {
		name  string
		items []measurement.MixedItem
		field string
		ids   []string
	}{
		{"empty", nil, "mixed_items", nil},
		{"duplicate id", []measurement.MixedItem{mixedItem("a", "D1", schema.DataTypeInt16), mixedItem("a", "D9", schema.DataTypeInt16)}, "mixed_items[1].item_id", []string{"a"}},
		{"bad address", []measurement.MixedItem{mixedItem("bad", "QQ9", schema.DataTypeInt16)}, "mixed_items[0].address", []string{"bad"}},
		{"missing address", []measurement.MixedItem{mixedItem("blank", " ", schema.DataTypeInt16)}, "mixed_items[0].address", []string{"blank"}},
		{"bad data type", []measurement.MixedItem{mixedItem("dt", "D1", "decimal")}, "mixed_items[0]", []string{"dt"}},
		{"missing semantic kind", []measurement.MixedItem{noKind}, "mixed_items[0]", []string{"d2"}},
		{"unsupported per-item scale", []measurement.MixedItem{scaled}, "mixed_items[0]", []string{"d1"}},
		{"same address twice", []measurement.MixedItem{mixedItem("x", "D7", schema.DataTypeInt16), mixedItem("y", "d7", schema.DataTypeInt16)}, "mixed_items", []string{"x", "y"}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := planMixedRuleItems(tc.items, schema.ProtocolMC3E)
			requireMixedFieldError(t, err, tc.field, tc.ids...)
		})
	}
}

func TestValidateRuleLayoutRejectsUnsupportedCombinations(t *testing.T) {
	items := []measurement.MixedItem{mixedItem("d1", "D1", schema.DataTypeInt16)}
	start := 10
	scale := 2.0
	target := schema.DataTypeFloat32
	cases := []struct {
		name  string
		req   CreateRuleRequest
		field string
	}{
		{"unknown layout", CreateRuleRequest{LayoutMode: "scattered", MixedItems: items}, "layout_mode"},
		{"items without mixed", CreateRuleRequest{LayoutMode: "homogeneous", MixedItems: items}, "mixed_items"},
		{"share", CreateRuleRequest{LayoutMode: "mixed", MixedItems: items, ShareEnabled: true, ShareStartRegister: &start}, "share_enabled"},
		{"skipped", CreateRuleRequest{LayoutMode: "mixed", MixedItems: items, SkippedAddresses: []string{"D1"}}, "skipped_addresses"},
		{"rule scale", CreateRuleRequest{LayoutMode: "mixed", MixedItems: items, ScaleMultiplier: &scale}, "scale_multiplier"},
		{"rule target", CreateRuleRequest{LayoutMode: "mixed", MixedItems: items, TargetDataType: &target}, "target_data_type"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := validateMixedRuleFields(normalizeLayoutMode(tc.req.LayoutMode), len(tc.req.MixedItems) > 0, tc.req.ShareEnabled, tc.req.SkippedAddresses, tc.req.TargetDataType, tc.req.ScaleMultiplier, tc.req.ScaleOffset)
			requireMixedFieldError(t, err, tc.field)
		})
	}
}

func TestSQLRepositoryPersistsMixedLayout(t *testing.T) {
	db := setupSQLRepoDB(t)
	defer db.Close()
	repo := NewSQLRepository(db)
	ctx := t.Context()
	now := time.Now().UTC().Truncate(time.Second)
	items := `[{"item_id":"d39","name":"","address":"D39","register_offset":0,"data_type":"int16","semantic_kind":"gauge"}]`
	mixed := &schema.SourceRule{ID: "mixed", DeviceID: "device-1", StartAddress: "D39", Count: 1, DataType: schema.DataTypeInt16, NamingPrefix: "MIX", Enabled: true, Origin: "manual", SkippedAddresses: "[]", LayoutMode: "mixed", MixedItems: items, RevisionID: "r1", CreatedAt: now, UpdatedAt: now}
	plain := &schema.SourceRule{ID: "plain", DeviceID: "device-1", StartAddress: "D0", Count: 2, DataType: schema.DataTypeInt16, NamingPrefix: "P", Enabled: true, Origin: "manual", SkippedAddresses: "[]", RevisionID: "r1", CreatedAt: now, UpdatedAt: now.Add(time.Second)}
	require.NoError(t, repo.Create(ctx, mixed))
	require.NoError(t, repo.Create(ctx, plain))

	got, err := repo.GetByID(ctx, "mixed")
	require.NoError(t, err)
	require.Equal(t, "mixed", got.LayoutMode)
	require.JSONEq(t, items, got.MixedItems)
	gotPlain, err := repo.GetByID(ctx, "plain")
	require.NoError(t, err)
	require.Empty(t, gotPlain.LayoutMode, "homogeneous rules keep their existing in-memory shape")
	require.Empty(t, gotPlain.MixedItems)
	var stored string
	require.NoError(t, db.QueryRowContext(ctx, `SELECT layout_mode FROM source_rules WHERE id='plain'`).Scan(&stored))
	require.Equal(t, "homogeneous", stored)

	listed, err := repo.List(ctx, ListFilter{})
	require.NoError(t, err)
	require.Len(t, listed, 2)
	require.Equal(t, "mixed", listed[0].LayoutMode)
	require.JSONEq(t, items, listed[0].MixedItems)

	next := *got
	next.MixedItems = `[{"item_id":"d40","name":"","address":"D40","register_offset":0,"data_type":"int16","semantic_kind":"gauge"}]`
	next.RevisionID = "r2"
	require.NoError(t, repo.UpdateRevision(ctx, &next, got))
	reread, err := repo.GetByID(ctx, "mixed")
	require.NoError(t, err)
	require.JSONEq(t, next.MixedItems, reread.MixedItems)

	stale := *got
	require.ErrorIs(t, repo.DeleteRevision(ctx, &stale), ErrRuleRevisionConflict, "a stale mixed layout must not match the CAS predicate")
	require.NoError(t, repo.DeleteRevision(ctx, reread))
}
