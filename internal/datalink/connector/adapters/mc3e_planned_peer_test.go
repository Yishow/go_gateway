package adapters

import (
	"encoding/binary"
	"fmt"
	"testing"

	"go-gateway/internal/datalink/schema"

	"github.com/stretchr/testify/require"
)

type mcPlannedWordAddress struct {
	code    byte
	address int
}

// The peer decodes binary 0403 independently of the production serializer:
// SH080008AB printed p101: LE24 device number, then one-byte device code.
func mcPlannedRandomAddresses(frame []byte) ([]mcPlannedWordAddress, error) {
	if len(frame) < 17 || frame[0] != 0x50 || frame[1] != 0 || int(binary.LittleEndian.Uint16(frame[7:])) != len(frame)-9 || binary.LittleEndian.Uint16(frame[11:]) != 0x0403 || binary.LittleEndian.Uint16(frame[13:]) != 0 {
		return nil, fmt.Errorf("invalid binary random request framing")
	}
	wordCount := int(frame[15])
	dwordCount := int(frame[16])
	totalEntries := wordCount + dwordCount
	if totalEntries == 0 || len(frame) != 17+totalEntries*4 {
		return nil, fmt.Errorf("invalid word/doubleword counts")
	}
	addresses := make([]mcPlannedWordAddress, 0, wordCount+dwordCount*2)
	for i := range totalEntries {
		entry := frame[17+4*i : 21+4*i]
		// This peer fixture exposes D, W and M memory only. Unknown codes cannot
		// be silently read as zeros or accepted because response values were preset.
		if entry[3] != 0xa8 && entry[3] != 0xb4 && entry[3] != 0x90 {
			return nil, fmt.Errorf("unknown fixture word device code %02x", entry[3])
		}
		addr := int(entry[0]) | int(entry[1])<<8 | int(entry[2])<<16
		addresses = append(addresses, mcPlannedWordAddress{entry[3], addr})
		if i >= wordCount {
			addresses = append(addresses, mcPlannedWordAddress{entry[3], addr + 1})
		}
	}
	return addresses, nil
}
func TestMCPlannedPeerRejectsASCIIOrderInBinaryEntries(t *testing.T) {
	wrong := []byte{0x50, 0, 0, 0xff, 0xff, 3, 0, 12, 0, 0x10, 0, 3, 4, 0, 0, 1, 0, 0xa8, 0, 0, 0}
	_, err := mcPlannedRandomAddresses(wrong)
	require.Error(t, err)
	right := []byte{0x50, 0, 0, 0xff, 0xff, 3, 0, 12, 0, 0x10, 0, 3, 4, 0, 0, 1, 0, 0, 0, 0, 0xa8}
	got, err := mcPlannedRandomAddresses(right)
	require.NoError(t, err)
	require.Equal(t, []mcPlannedWordAddress{{0xa8, 0}}, got)
}
func TestMCPlannedRandomReadsAddressedPeerMemory(t *testing.T) {
	req := mcPlannedRequest(t, mcPlannedPoint("wide", "D0", schema.DataTypeUint32, "ABCD"), mcPlannedPoint("swapped", "D0", schema.DataTypeUint32, "CDAB"), mcPlannedPoint("far", "D10", schema.DataTypeInt16, ""))
	tr := &mcPlannedTransport{wordMemory: map[mcPlannedWordAddress]uint16{{0xa8, 0}: 0x0102, {0xa8, 1}: 0x0304, {0xa8, 10}: 0x1234}}
	got, err := mcPlannedConnector(tr).ReadPlanned(t.Context(), req)
	require.NoError(t, err)
	assertMCPlannedResults(t, req, got, map[string]any{"wide": uint32(0x01020304), "swapped": uint32(0x03040102), "far": int16(0x1234)})
	require.Len(t, tr.frames, 1)
	require.Equal(t, []byte{1, 1, 10, 0, 0, 0xa8, 0, 0, 0, 0xa8}, tr.frames[0][15:])
}

func TestMCPlannedRandomMissingPeerWordDoesNotInventGoodZero(t *testing.T) {
	req := mcPlannedRequest(t, mcPlannedPoint("near", "D0", schema.DataTypeInt16, ""), mcPlannedPoint("missing", "D10", schema.DataTypeInt16, ""))
	tr := &mcPlannedTransport{wordMemory: map[mcPlannedWordAddress]uint16{{0xa8, 0}: 11}}
	got, err := mcPlannedConnector(tr).ReadPlanned(t.Context(), req)
	require.Error(t, err)
	require.Len(t, tr.frames, 1)
	require.Len(t, got, 2)
	for _, point := range got {
		require.Equal(t, schema.QualityBad, point.Quality)
		require.Nil(t, point.Value)
	}
}

func TestMCPlannedRandomMixedWordDoubleWordAndAlignedBit(t *testing.T) {
	req := mcPlannedRequest(t,
		mcPlannedPoint("d39", "D39", schema.DataTypeInt16, ""),
		mcPlannedPoint("d100", "D100", schema.DataTypeInt32, "CDAB"),
		mcPlannedPoint("m5", "M5", schema.DataTypeBool, ""),
	)
	require.Len(t, req.Plan.Requests(), 1)
	require.Equal(t, "random", req.Plan.Requests()[0].Kind)
	tr := &mcPlannedTransport{wordMemory: map[mcPlannedWordAddress]uint16{
		{0xa8, 39}:  0x0100, // D39 = 256
		{0x90, 0}:   0x0020, // M0..M15: bit 5 (M5) = 1
		{0xa8, 100}: 0x5678, // D100 low word
		{0xa8, 101}: 0x1234, // D100 high word -> CDAB int32 = 0x12345678
	}}
	got, err := mcPlannedConnector(tr).ReadPlanned(t.Context(), req)
	require.NoError(t, err)
	assertMCPlannedResults(t, req, got, map[string]any{
		"d39":  int16(256),
		"d100": int32(0x12345678),
		"m5":   true,
	})
	require.Len(t, tr.frames, 1)
	require.Equal(t, []byte{
		2, 1,
		39, 0, 0, 0xa8,
		0, 0, 0, 0x90,
		100, 0, 0, 0xa8,
	}, tr.frames[0][15:])
}
