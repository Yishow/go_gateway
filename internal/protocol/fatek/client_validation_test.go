package fatek

import "testing"

func TestClient_ReadRejectsNonPositiveCountWithoutSending(t *testing.T) {
	tests := []struct {
		name string
		call func(*FatekClient) error
	}{
		{
			name: "status zero",
			call: func(client *FatekClient) error {
				_, err := client.ReadStatus("X", 0, 0)
				return err
			},
		},
		{
			name: "status negative",
			call: func(client *FatekClient) error {
				_, err := client.ReadStatus("X", 0, -1)
				return err
			},
		},
		{
			name: "register zero",
			call: func(client *FatekClient) error {
				_, err := client.ReadRegisters("D", 0, 0)
				return err
			},
		},
		{
			name: "register negative",
			call: func(client *FatekClient) error {
				_, err := client.ReadRegisters("D", 0, -1)
				return err
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			transport := NewMockFatekTransport()
			client := NewClient(transport, 1)
			before := len(transport.sentData)

			if err := tt.call(client); err == nil {
				t.Fatal("read should reject non-positive count")
			}
			if got := len(transport.sentData); got != before {
				t.Fatalf("invalid count sent %d requests, want %d", got-before, 0)
			}
		})
	}
}
