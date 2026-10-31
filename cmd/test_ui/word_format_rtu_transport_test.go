package main

import (
	"encoding/binary"
	"fmt"
	"path/filepath"
	"reflect"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go-gateway/internal/datalink/connector"
	"go-gateway/internal/datalink/connector/adapters"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/protocol/modbus"
)

// Simulated RTU byte transport, not a decoded Protocol. The actual ModbusClient
// builds/checks RTU address/function/CRC and the production adapter decodes words.
type wordRTUPeer struct{ plc *wordWirePLC }

func (*wordRTUPeer) Connect() error            { return nil }
func (*wordRTUPeer) Close() error              { return nil }
func (*wordRTUPeer) IsConnected() bool         { return true }
func (*wordRTUPeer) GetOriginalTransport() any { return &modbus.RTUTransport{} }
func (p *wordRTUPeer) SendReceive(raw []byte) ([]byte, error) {
	unit, function, payload, err := modbus.ParseRTUFrame(raw)
	if err != nil {
		return nil, err
	}
	if len(payload) != 4 || function != 3 && function != 4 {
		return nil, fmt.Errorf("unexpected simulated RTU request")
	}
	words, bad := p.plc.read(int(binary.BigEndian.Uint16(payload[:2])), int(binary.BigEndian.Uint16(payload[2:])))
	response := []byte{byte(len(words) * 2)}
	for _, word := range words {
		response = append(response, byte(word>>8), byte(word))
	}
	if bad {
		response = response[:len(response)-1]
	}
	return modbus.BuildRTUFrame(unit, function, response), nil
}
func wordScopedRTUTransport(plc *wordWirePLC, port string) func(schema.ConnectionConfigModbusRTU) modbus.Transport {
	return func(config schema.ConnectionConfigModbusRTU) modbus.Transport {
		if config.SerialPort == port {
			return &wordRTUPeer{plc}
		}
		return modbus.NewRTUTransport(config.SerialPort, config.BaudRate, config.DataBits, config.StopBits, config.Parity, time.Duration(config.Timeout)*time.Second)
	}
}
func registerWordRTUFixture(t *testing.T, plc *wordWirePLC, port string) {
	t.Helper()
	original, ok := connector.GetFactory(schema.ProtocolModbusRTU)
	require.True(t, ok)
	t.Cleanup(func() { connector.Register(schema.ProtocolModbusRTU, original) })
	connector.Register(schema.ProtocolModbusRTU, func() connector.Protocol {
		return adapters.NewModbusRTUConnectorWithTransportFactory(wordScopedRTUTransport(plc, port))
	})
}
func wordRTUFixture(t *testing.T) wordProductionFixture {
	t.Helper()
	plc := &wordWirePLC{words: map[int]uint16{}}
	port := filepath.Join(t.TempDir(), "simulated-raw-rtu")
	registerWordRTUFixture(t, plc, port)
	config := fmt.Sprintf(`{"serial_port":%q,"baud_rate":9600,"slave_id":1,"timeout":1}`, port)
	return newWordTransportFixture(t, schema.ProtocolModbusRTU, plc, config)
}
func TestProductionWordFormatRTUToIngestorMatrix(t *testing.T) {
	exerciseWordTransportMatrix(t, wordRTUFixture(t), schema.ProtocolModbusRTU)
}
func TestProductionWordFormatRTUZeroShortAndLegacy(t *testing.T) {
	exerciseWordTransportBoundaries(t, wordRTUFixture(t))
}

func TestWordRTUFixtureScopesTransportAndRestoresRegistry(t *testing.T) {
	original, ok := connector.GetFactory(schema.ProtocolModbusRTU)
	require.True(t, ok)
	t.Run("scoped", func(t *testing.T) {
		plc := &wordWirePLC{words: map[int]uint16{}}
		port := filepath.Join(t.TempDir(), "raw-only-this-fixture")
		registerWordRTUFixture(t, plc, port)
		factory := wordScopedRTUTransport(plc, port)
		require.IsType(t, &wordRTUPeer{}, factory(schema.ConnectionConfigModbusRTU{SerialPort: port}))
		var wg sync.WaitGroup
		results := make(chan modbus.Transport, 8)
		for i := 0; i < 8; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				results <- factory(schema.ConnectionConfigModbusRTU{SerialPort: "foreign-unopened", BaudRate: 19200})
			}()
		}
		wg.Wait()
		close(results)
		for transport := range results {
			native, ok := transport.(*modbus.RTUTransport)
			require.True(t, ok, "foreign config must retain real transport; never return fixture bytes")
			require.Equal(t, "foreign-unopened", native.Port)
			require.Equal(t, 19200, native.BaudRate)
		}
		require.Zero(t, plc.calls, "constructing fallback transports performs no serial I/O")
	})
	restored, ok := connector.GetFactory(schema.ProtocolModbusRTU)
	require.True(t, ok)
	require.Equal(t, reflect.ValueOf(original).Pointer(), reflect.ValueOf(restored).Pointer())
}
