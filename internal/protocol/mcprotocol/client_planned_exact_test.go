package mcprotocol

import (
	"bytes"
	"testing"
)

type plannedExactTransport struct {
	response, request []byte
	sends             int
}

func (*plannedExactTransport) Connect() error { return nil }
func (*plannedExactTransport) Close() error   { return nil }
func (tr *plannedExactTransport) SendReceive(req []byte) ([]byte, error) {
	tr.sends++
	tr.request = append([]byte(nil), req...)
	return tr.response, nil
}

func TestPlannedExactRandomRejectsMalformedPayload(t *testing.T) {
	for _, response := range [][]byte{nil, {0}, {0, 0}, {0, 0, 1}, {0, 0, 1, 0, 2}, {0, 0, 1, 0, 2, 0}, {0x59, 0xc0}} {
		tr := &plannedExactTransport{response: response}
		got, err := NewClientWithTransport(tr).RandomReadExact([]RandomReadItem{{Device: "D", Addr: 0}})
		if err == nil || got != nil || tr.sends != 1 {
			t.Fatalf("response %x accepted: %v %v sends %d", response, got, err, tr.sends)
		}
	}
}
func TestPlannedExactBitsRejectsMalformedPayload(t *testing.T) {
	for _, response := range [][]byte{nil, {0}, {0, 0}, {0, 0, 0x10, 0}, {0, 0, 0x20}, {0, 0, 0x12}, {0x59, 0xc0}} {
		tr := &plannedExactTransport{response: response}
		got, err := NewClientWithTransport(tr).BatchReadBitExact("M", 0, 2)
		if err == nil || got != nil || tr.sends != 1 {
			t.Fatalf("response %x accepted: %v %v sends %d", response, got, err, tr.sends)
		}
	}
}

// Binary 0403 layout: Mitsubishi SH080008AB printed p101 (PDF p103),
// https://www.mitsubishielectric.com/dl/fa/document/manual/plc/sh080008/sh080008ab.pdf
// Each item is LE24 device number THEN device code, not the ASCII layout.
func TestRandomReadBinary0403OfficialGolden(t *testing.T) {
	for _, exact := range []bool{false, true} {
		tr := &plannedExactTransport{response: []byte{0, 0, 11, 0, 22, 0, 33, 0}}
		items := []RandomReadItem{{Device: "D", Addr: 0}, {Device: "D", Addr: 10}, {Device: "D", Addr: 20}}
		client := NewClientWithTransport(tr)
		var err error
		if exact {
			_, err = client.RandomReadExact(items)
		} else {
			_, err = client.RandomRead(items)
		}
		if err != nil {
			t.Fatal(err)
		}
		want := []byte{0x50, 0, 0, 0xff, 0xff, 3, 0, 20, 0, 0x10, 0, 3, 4, 0, 0, 3, 0, 0, 0, 0, 0xa8, 10, 0, 0, 0xa8, 20, 0, 0, 0xa8}
		if !bytes.Equal(tr.request, want) {
			t.Fatalf("exact=%v got %x want official binary %x", exact, tr.request, want)
		}
	}
}
func TestRandomReadBinary0403AddressBytesAndWordCode(t *testing.T) {
	for _, exact := range []bool{false, true} {
		tr := &plannedExactTransport{response: []byte{0, 0, 1, 0, 2, 0}}
		client := NewClientWithTransport(tr)
		items := []RandomReadItem{{Device: "D", Addr: 0x123456}, {Device: "W", Addr: 0xabcdef}}
		var err error
		if exact {
			_, err = client.RandomReadExact(items)
		} else {
			_, err = client.RandomRead(items)
		}
		if err != nil {
			t.Fatal(err)
		}
		want := []byte{2, 0, 0x56, 0x34, 0x12, 0xa8, 0xef, 0xcd, 0xab, 0xb4}
		if !bytes.Equal(tr.request[15:], want) {
			t.Fatalf("exact=%v got %x want %x", exact, tr.request[15:], want)
		}
	}
}

func TestPlannedExactBitOddPaddingIsValidated(t *testing.T) {
	for _, b := range []byte{0x01, 0x11, 0x12} {
		tr := &plannedExactTransport{response: []byte{0, 0, b}}
		got, err := NewClientWithTransport(tr).BatchReadBitExact("M", 0, 1)
		if err == nil || got != nil {
			t.Fatalf("nonzero odd padding accepted: %x", b)
		}
	}
	for _, b := range []byte{0x00, 0x10} {
		tr := &plannedExactTransport{response: []byte{0, 0, b}}
		got, err := NewClientWithTransport(tr).BatchReadBitExact("M", 0, 1)
		if err != nil || len(got) != 1 || got[0] != (b == 0x10) {
			t.Fatalf("valid bit %x: %v %v", b, got, err)
		}
	}
}

func TestPlannedExactBitSerializerGolden(t *testing.T) {
	for _, exact := range []bool{false, true} {
		tr := &plannedExactTransport{response: []byte{0, 0, 0x10}}
		client := NewClientWithTransport(tr)
		var err error
		if exact {
			_, err = client.BatchReadBitExact("X", 0x10, 2)
		} else {
			_, err = client.BatchReadBit("X", 0x10, 2)
		}
		if err != nil {
			t.Fatal(err)
		}
		want := []byte{0x50, 0, 0, 0xff, 0xff, 3, 0, 12, 0, 0x10, 0, 1, 4, 1, 0, 0x10, 0, 0, 0x9c, 2, 0}
		if !bytes.Equal(tr.request, want) {
			t.Fatalf("exact=%v got %x want %x", exact, tr.request, want)
		}
	}
}

func TestRandomReadBinary0403WordDoubleWordAndAlignedBitSH080008(t *testing.T) {
	// D39 (word) = 0x1234, M0..M15 (word) = 0x0020 (M5=1), D100 (dword) = 0x12345678
	resp := []byte{
		0x00, 0x00, // EndCode = 0
		0x34, 0x12, // word[0]: D39 = 0x1234
		0x20, 0x00, // word[1]: M0  = 0x0020 (bit 5 set)
		0x78, 0x56, 0x34, 0x12, // dword[0]: D100 = 0x12345678
	}
	items := []RandomReadItem{
		{Device: "D", Addr: 39},
		{Device: "D", Addr: 100, DoubleWord: true},
		{Device: "M", Addr: 0},
	}
	wantReq := []byte{
		0x50, 0x00, 0x00, 0xff, 0xff, 0x03, 0x00, 20, 0x00, 0x10, 0x00,
		0x03, 0x04, 0x00, 0x00,
		0x02, 0x01, // 2 word points (D39, M0), 1 dword point (D100)
		39, 0x00, 0x00, 0xa8, // D39 in word section
		0x00, 0x00, 0x00, 0x90, // M0 in word section
		100, 0x00, 0x00, 0xa8, // D100 in dword section
	}
	for _, exact := range []bool{false, true} {
		tr := &plannedExactTransport{response: resp}
		client := NewClientWithTransport(tr)
		var (
			got []int
			err error
		)
		if exact {
			got, err = client.RandomReadExact(items)
		} else {
			got, err = client.RandomRead(items)
		}
		if err != nil {
			t.Fatalf("exact=%v RandomRead failed: %v", exact, err)
		}
		if !bytes.Equal(tr.request, wantReq) {
			t.Fatalf("exact=%v got request %x want SH-080008 %x", exact, tr.request, wantReq)
		}
		if len(got) != 3 || got[0] != 0x1234 || got[1] != 0x12345678 || got[2] != 0x0020 {
			t.Fatalf("exact=%v got values %#v, want [0x1234, 0x12345678, 0x0020]", exact, got)
		}
		m5, err := ExtractBitFromWord(got[2], 5)
		if err != nil || !m5 {
			t.Fatalf("exact=%v expected M5=true from M0 word 0x%04x, got %v err=%v", exact, got[2], m5, err)
		}
		m4, err := ExtractBitFromWord(got[2], 4)
		if err != nil || m4 {
			t.Fatalf("exact=%v expected M4=false from M0 word 0x%04x, got %v err=%v", exact, got[2], m4, err)
		}
	}
}
