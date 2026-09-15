package mcprotocol

import "testing"

func TestClient_RandomRead_RejectsInvalidInputsWithoutSending(t *testing.T) {
	items256 := make([]RandomReadItem, 256)
	for i := range items256 {
		items256[i] = RandomReadItem{Device: "D", Addr: i}
	}

	tests := []struct {
		name  string
		items []RandomReadItem
	}{
		{name: "empty", items: nil},
		{name: "count above wire limit", items: items256},
		{name: "negative address", items: []RandomReadItem{{Device: "D", Addr: -1}}},
		{name: "address above 24 bit limit", items: []RandomReadItem{{Device: "D", Addr: 0x1000000}}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			transport := NewMockMCTransport()
			client := NewClientWithTransport(transport)
			before := len(transport.sentData)

			if _, err := client.RandomRead(tt.items); err == nil {
				t.Fatal("RandomRead should reject invalid input")
			}
			if got := len(transport.sentData); got != before {
				t.Fatalf("invalid input sent %d requests, want %d", got-before, 0)
			}
		})
	}
}
