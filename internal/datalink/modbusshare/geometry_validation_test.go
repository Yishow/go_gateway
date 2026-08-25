package modbusshare

import (
	"context"
	"errors"
	"testing"

	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/tag"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRegisterGeometry_HumanAndZeroBasedConversion(t *testing.T) {
	// 40001 -> 0
	zero, err := HumanToZeroBased(40001)
	require.NoError(t, err)
	assert.Equal(t, uint16(0), zero)
	assert.Equal(t, uint32(40001), ZeroBasedToHuman(0))

	// 40010 -> 9
	zero, err = HumanToZeroBased(40010)
	require.NoError(t, err)
	assert.Equal(t, uint16(9), zero)
	assert.Equal(t, uint32(40010), ZeroBasedToHuman(9))

	// Invalid < 40001
	_, err = HumanToZeroBased(40000)
	assert.Error(t, err)

	// Overflow > 40001 + 65535
	_, err = HumanToZeroBased(40001 + 65536)
	assert.Error(t, err)
}

func TestRegisterGeometry_DatatypeSpansAndBytes(t *testing.T) {
	cases := []struct {
		dt        schema.DataType
		wantSpan  int
		wantBytes int
	}{
		{schema.DataTypeBool, 1, 2},
		{schema.DataTypeInt16, 1, 2},
		{schema.DataTypeUint16, 1, 2},
		{schema.DataTypeInt32, 2, 4},
		{schema.DataTypeUint32, 2, 4},
		{schema.DataTypeFloat32, 2, 4},
		{schema.DataTypeInt64, 4, 8},
		{schema.DataTypeUint64, 4, 8},
		{schema.DataTypeFloat64, 4, 8},
	}

	for _, tc := range cases {
		span := DataTypeSpan(tc.dt)
		bytesCount := DataTypeBytes(tc.dt)
		assert.Equal(t, tc.wantSpan, span, "span for %s", tc.dt)
		assert.Equal(t, tc.wantBytes, bytesCount, "bytes for %s", tc.dt)
	}
}

func TestRegisterGeometry_StrideValidation(t *testing.T) {
	// Stride < Span should fail validation
	mappings := []DesiredMapping{
		{
			TagID:              "tag-float",
			DataType:           schema.DataTypeFloat32, // span = 2
			ZeroBasedRegister:  0,
			ShareStartRegister: 40001,
			SpanRegisters:      2,
			StrideRegisters:    1, // INVALID: stride 1 < span 2
		},
	}

	err := ValidateDesiredMappings(mappings, 100)
	require.Error(t, err)
	var shareErr *Error
	require.ErrorAs(t, err, &shareErr)
	assert.Equal(t, ErrCodeInvalidGeometry, shareErr.Code)
}

func TestRegisterGeometry_CapacityValidation(t *testing.T) {
	// Capacity is 8 registers [0, 8)
	// int64 at start 6 occupies [6, 10), which exceeds capacity 8
	mappings := []DesiredMapping{
		{
			TagID:              "tag-int64",
			DataType:           schema.DataTypeInt64, // span = 4
			ZeroBasedRegister:  6,
			ShareStartRegister: 40007,
			SpanRegisters:      4,
			StrideRegisters:    4,
		},
	}

	err := ValidateDesiredMappings(mappings, 8)
	require.Error(t, err)
	var shareErr *Error
	require.ErrorAs(t, err, &shareErr)
	assert.Equal(t, ErrCodeCapacityExceeded, shareErr.Code)
}

func TestRegisterGeometry_RangeCollisionDetection(t *testing.T) {
	// Mapping 1: int32 at 0, span=2 -> occupies [0, 2)
	// Mapping 2: float32 at 1, span=2 -> occupies [1, 3) -> COLLISION with mapping 1!
	mappings := []DesiredMapping{
		{
			TagID:              "tag-1",
			DataType:           schema.DataTypeInt32,
			ZeroBasedRegister:  0,
			ShareStartRegister: 40001,
			SpanRegisters:      2,
			StrideRegisters:    2,
		},
		{
			TagID:              "tag-2",
			DataType:           schema.DataTypeFloat32,
			ZeroBasedRegister:  1,
			ShareStartRegister: 40002,
			SpanRegisters:      2,
			StrideRegisters:    2,
		},
	}

	err := ValidateDesiredMappings(mappings, 100)
	require.Error(t, err)
	var shareErr *Error
	require.ErrorAs(t, err, &shareErr)
	assert.Equal(t, ErrCodeRangeCollision, shareErr.Code)
}

func TestRegisterGeometry_ValidateDesiredMappingsSuccess(t *testing.T) {
	// Mapping 1: int32 at 0, span=2 -> [0, 2)
	// Mapping 2: int64 at 2, span=4 -> [2, 6)
	// Mapping 3: int16 at 6, span=1 -> [6, 7)
	mappings := []DesiredMapping{
		{
			TagID:              "tag-1",
			DataType:           schema.DataTypeInt32,
			ZeroBasedRegister:  0,
			ShareStartRegister: 40001,
			SpanRegisters:      2,
			StrideRegisters:    2,
		},
		{
			TagID:              "tag-2",
			DataType:           schema.DataTypeInt64,
			ZeroBasedRegister:  2,
			ShareStartRegister: 40003,
			SpanRegisters:      4,
			StrideRegisters:    4,
		},
		{
			TagID:              "tag-3",
			DataType:           schema.DataTypeInt16,
			ZeroBasedRegister:  6,
			ShareStartRegister: 40007,
			SpanRegisters:      1,
			StrideRegisters:    1,
		},
	}

	err := ValidateDesiredMappings(mappings, 100)
	require.NoError(t, err)
}

func TestRegisterGeometry_RejectsDuplicateTagIdentityBeforeMutation(t *testing.T) {
	mappings := []DesiredMapping{
		{TagID: "tag-duplicate", DataType: schema.DataTypeInt16, ZeroBasedRegister: 0, ShareStartRegister: 40001, SpanRegisters: 1, StrideRegisters: 1},
		{TagID: "tag-duplicate", DataType: schema.DataTypeInt16, ZeroBasedRegister: 2, ShareStartRegister: 40003, SpanRegisters: 1, StrideRegisters: 1},
	}
	err := ValidateDesiredMappings(mappings, 16)
	var shareErr *Error
	if !errors.As(err, &shareErr) {
		t.Fatalf("expected typed duplicate identity error, got %v", err)
	}
	if shareErr.Code != ErrCodeInvalidGeometry {
		t.Fatalf("expected invalid geometry, got %s", shareErr.Code)
	}
}

func TestRegisterGeometry_InvalidSetDoesNotMutateService(t *testing.T) {
	tagRepo := tag.NewMemoryRepository()
	tagSvc := tag.NewService(tagRepo)
	svc := NewService(tagSvc, 65536)

	// Valid initial mapping
	ctx := context.Background()
	t1, err := tagSvc.Create(ctx, tag.CreateTagRequest{Key: "t1", DisplayName: "Tag 1", DataType: schema.DataTypeInt16})
	require.NoError(t, err)
	_, err = svc.UpsertMapping(ctx, t1.ID, 0)
	require.NoError(t, err)
	require.True(t, svc.HasMapping(t1.ID))

	// Attempting invalid desired batch with collision
	collidingMappings := []DesiredMapping{
		{
			TagID:             t1.ID,
			DataType:          schema.DataTypeInt32,
			ZeroBasedRegister: 0,
			SpanRegisters:     2,
			StrideRegisters:   2,
		},
		{
			TagID:             "tag-bad",
			DataType:          schema.DataTypeInt32,
			ZeroBasedRegister: 1, // collision
			SpanRegisters:     2,
			StrideRegisters:   2,
		},
	}

	err = ValidateDesiredMappings(collidingMappings, 100)
	require.Error(t, err)

	// Service state should be unchanged
	assert.True(t, svc.HasMapping(t1.ID))
	assert.False(t, svc.HasMapping("tag-bad"))
}
