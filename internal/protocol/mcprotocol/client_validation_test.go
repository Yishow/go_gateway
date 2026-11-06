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
		{name: "unaligned bit device M5", items: []RandomReadItem{{Device: "M", Addr: 5}}},
		{name: "bit device as double word", items: []RandomReadItem{{Device: "M", Addr: 0, DoubleWord: true}}},
		{name: "long timer current value LTN", items: []RandomReadItem{{Device: "LTN", Addr: 0}}},
		{name: "long timer contact LTS", items: []RandomReadItem{{Device: "LTS", Addr: 0}}},
		{name: "long timer coil LTC", items: []RandomReadItem{{Device: "LTC", Addr: 0}}},
		{name: "long retentive timer LSTS", items: []RandomReadItem{{Device: "LSTS", Addr: 0}}},
		{name: "long retentive timer LSTC", items: []RandomReadItem{{Device: "LSTC", Addr: 0}}},
		{name: "long retentive timer LSTN", items: []RandomReadItem{{Device: "LSTN", Addr: 0}}},
		{name: "long counter contact LCS", items: []RandomReadItem{{Device: "LCS", Addr: 0}}},
		{name: "long counter coil LCC", items: []RandomReadItem{{Device: "LCC", Addr: 0}}},
		{name: "long counter current value LCN", items: []RandomReadItem{{Device: "LCN", Addr: 0}}},
		{name: "string data type", items: []RandomReadItem{{Device: "D", Addr: 0, DataType: "string"}}},
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
