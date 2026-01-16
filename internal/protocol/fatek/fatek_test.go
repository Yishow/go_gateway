package fatek

import (
	"fmt"
	"testing"
)

// Mock Transport for Benchmark
type MockTransport struct{}

func (m *MockTransport) Connect() error { return nil }
func (m *MockTransport) Close() error   { return nil }
func (m *MockTransport) SendReceive(data []byte) ([]byte, error) {
	// Simulate simple response: STX + Station + Cmd + 0 (Success) + LRC + ETX
	// Min length 9 bytes.
	return []byte{0x02, '0', '1', '4', '4', '0', 'F', 'F', 0x03}, nil
}

func TestBuildFrameToBuffer(t *testing.T) {
	buf := GetBuffer()
	defer PutBuffer(buf)

	BuildFrameToBuffer(buf, 1, "44", "TEST")

	if buf.Len() == 0 {
		t.Error("Buffer should not be empty")
	}
	// Basic check: STX(02) + "01" + "44" + "TEST" + LRC + ETX(03)
	// Len = 1 + 2 + 2 + 4 + 2 + 1 = 12
	if buf.Len() != 12 {
		t.Errorf("Expected length 12, got %d", buf.Len())
	}
}

// Benchmark the legacy BuildFrame (which allocates new slice)
func BenchmarkBuildFrame_Alloc(b *testing.B) {
	// Re-implement legacy logic locally to compare
	legacyBuild := func(station int, cmd string, body string) []byte {
		stationStr := fmt.Sprintf("%02X", station)
		content := stationStr + cmd + body
		lrcData := make([]byte, 0, 1+len(content))
		lrcData = append(lrcData, STX)
		lrcData = append(lrcData, []byte(content)...)
		lrc := CalculateLRC(lrcData)
		frame := make([]byte, 0, len(lrcData)+3)
		frame = append(frame, lrcData...)
		frame = append(frame, []byte(lrc)...)
		frame = append(frame, ETX)
		return frame
	}

	body := "0AF00000" // Example body
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = legacyBuild(1, "44", body)
	}
}

// Benchmark the new BuildFrameToBuffer (which uses sync.Pool)
func BenchmarkBuildFrame_Pool(b *testing.B) {
	body := "0AF00000"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		buf := GetBuffer()
		BuildFrameToBuffer(buf, 1, "44", body)
		_ = buf.Bytes() // Simulate usage
		PutBuffer(buf)
	}
}

// Benchmark Client Execute (Full Cycle)
func BenchmarkClient_Execute(b *testing.B) {
	transport := &MockTransport{}
	client := NewClient(transport, 1)

	// Pre-calculate args
	cmd := "44"
	body := "0A00000"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// This uses the internal pool logic
		_, _ = client.execute(cmd, body)
	}
}
