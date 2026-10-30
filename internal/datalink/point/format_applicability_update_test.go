package point

import (
	"github.com/stretchr/testify/require"
	"go-gateway/internal/datalink/schema"
	"testing"
)

func TestOrdinaryUpdateCannotMakeOverrideInapplicable(t *testing.T) {
	repo := NewMemoryRepository()
	svc := NewService(repo, nil)
	p, err := svc.Create(t.Context(), CreatePointRequest{DeviceID: "d", Name: "p", Address: "40001", DataType: schema.DataTypeInt32})
	require.NoError(t, err)
	p.FormatMode = schema.PointFormatOverride
	p.FormatOverride = "ABCD"
	p.DataFormat = "ABCD"
	require.NoError(t, repo.Update(t.Context(), p))
	for _, dtype := range []schema.DataType{schema.DataTypeBool, schema.DataTypeInt16, schema.DataTypeString} {
		_, err = svc.Update(t.Context(), p.ID, UpdatePointRequest{DataType: &dtype})
		require.ErrorIs(t, err, ErrFormatInvalid)
		stored, e := svc.GetByID(t.Context(), p.ID)
		require.NoError(t, e)
		require.Equal(t, p.RevisionID, stored.RevisionID)
		require.Equal(t, p.DataType, stored.DataType)
	}
}

func TestOrdinarySQLUpdateValidatesActualOverrideContext(t *testing.T) {
	for _, dtype := range []schema.DataType{schema.DataTypeBool, schema.DataTypeInt16, schema.DataTypeString, schema.DataTypeFloat32} {
		t.Run(string(dtype), func(t *testing.T) {
			_, svc, req := formatFixture(t)
			original, err := svc.MutateFormat(t.Context(), req)
			require.NoError(t, err)
			got, err := svc.Update(t.Context(), req.PointID, UpdatePointRequest{DataType: &dtype})
			if dtype == schema.DataTypeFloat32 {
				require.NoError(t, err)
				require.Equal(t, dtype, got.DataType)
				require.Equal(t, original.Point.FormatOverride, got.FormatOverride)
			} else {
				require.ErrorIs(t, err, ErrFormatInvalid)
				stored, e := svc.GetByID(t.Context(), req.PointID)
				require.NoError(t, e)
				require.Equal(t, original.Point.RevisionID, stored.RevisionID)
			}
		})
	}
}
func TestOrdinarySQLNameEditKeepsLegacyAndOverrideIntent(t *testing.T) {
	for _, override := range []bool{false, true} {
		_, svc, req := formatFixture(t)
		if override {
			_, err := svc.MutateFormat(t.Context(), req)
			require.NoError(t, err)
		}
		before, err := svc.GetByID(t.Context(), req.PointID)
		require.NoError(t, err)
		name := "renamed"
		got, err := svc.Update(t.Context(), req.PointID, UpdatePointRequest{Name: &name})
		require.NoError(t, err)
		require.Equal(t, before.FormatMode, got.FormatMode)
		require.Equal(t, before.FormatOverride, got.FormatOverride)
		require.Equal(t, before.DataFormat, got.DataFormat)
	}
}

func TestOrdinarySQLOverrideUsesActualModbusAddressFunction(t *testing.T) {
	for _, tc := range []struct {
		name, address, function string
		valid                   bool
	}{
		{"raw-register", "100", "03", true}, {"raw-bit", "50000", "01", false}, {"alias-wins", "HR100", "01", true}, {"traditional-wins", "40001", "01", true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			_, svc, req := formatFixture(t)
			original, err := svc.MutateFormat(t.Context(), req)
			require.NoError(t, err)
			got, err := svc.Update(t.Context(), req.PointID, UpdatePointRequest{Address: &tc.address, Function: &tc.function})
			if tc.valid {
				require.NoError(t, err)
				require.Equal(t, tc.address, got.Address)
			} else {
				require.ErrorIs(t, err, ErrFormatInvalid)
				got, err = svc.GetByID(t.Context(), req.PointID)
				require.NoError(t, err)
				require.Equal(t, original.Point.RevisionID, got.RevisionID)
			}
		})
	}
}

func TestGuardedFormatMutationAndReceiptUseActualFunction(t *testing.T) {
	for _, function := range []string{"01", "03", "99"} {
		t.Run(function, func(t *testing.T) {
			_, svc, req := formatFixture(t)
			address := "50000"
			req.Patch = &UpdatePointRequest{Address: &address, Function: &function}
			result, err := svc.MutateFormat(t.Context(), req)
			if function != "03" {
				require.ErrorIs(t, err, ErrFormatInvalid)
				return
			}
			require.NoError(t, err)
			require.True(t, result.Resolution.Applicable)
			replay, err := svc.MutateFormat(t.Context(), req)
			require.NoError(t, err)
			require.Equal(t, result, replay)
			inspected, err := svc.InspectFormatMutation(t.Context(), req.FormatScope)
			require.NoError(t, err)
			require.Equal(t, result, inspected)
		})
	}
}
