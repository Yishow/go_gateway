package adapters

import (
	"encoding/binary"
	"fmt"
	"sync"
	"testing"

	"go-gateway/internal/datalink/connector"
	"go-gateway/internal/datalink/mapping"
	"go-gateway/internal/datalink/point"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/protocol/modbus"

	"github.com/stretchr/testify/require"
)

type plcModbusAddressParseCase struct {
	name     string
	address  string
	function string
	wantAddr uint16
	wantFunc string
}

func TestPLCModbusAddress_ParseAliasesAndTraditional(t *testing.T) {
	tests := []plcModbusAddressParseCase{
		{name: "five digit coil prefix wins", address: "00001", function: "04", wantAddr: 0, wantFunc: "01"},
		{name: "five digit discrete input prefix wins", address: "10001", function: "03", wantAddr: 0, wantFunc: "02"},
		{name: "five digit input register prefix wins", address: "30001", function: "01", wantAddr: 0, wantFunc: "04"},
		{name: "five digit holding register prefix wins", address: "40001", function: "04", wantAddr: 0, wantFunc: "03"},
		{name: "holding register alias", address: "HR100", function: "01", wantAddr: 100, wantFunc: "03"},
		{name: "input register alias", address: "IR100", function: "03", wantAddr: 100, wantFunc: "04"},
		{name: "coil alias", address: "C100", function: "04", wantAddr: 100, wantFunc: "01"},
		{name: "coil long alias", address: "CO100", function: "04", wantAddr: 100, wantFunc: "01"},
		{name: "discrete input alias", address: "DI100", function: "01", wantAddr: 100, wantFunc: "02"},
		{name: "explicit raw register function", address: "100", function: "03", wantAddr: 100, wantFunc: "03"},
		{name: "explicit raw input function", address: "100", function: "04", wantAddr: 100, wantFunc: "04"},
		{name: "six digit traditional-looking value stays raw", address: "040001", function: "03", wantAddr: 40001, wantFunc: "03"},
		{name: "legacy no-function one-based coil", address: "1", wantAddr: 0, wantFunc: "01"},
		{name: "legacy no-function short numeric coil", address: "100", wantAddr: 99, wantFunc: "01"},
		{name: "no-function zero raw register", address: "0", wantAddr: 0, wantFunc: "03"},
		{name: "no-function raw register at default boundary", address: "10000", wantAddr: 10000, wantFunc: "03"},
		{name: "no-function raw register at maximum", address: "65535", wantAddr: 65535, wantFunc: "03"},
		{name: "raw max address", address: "65535", function: "03", wantAddr: 65535, wantFunc: "03"},
		{name: "raw zero address", address: "0", function: "04", wantAddr: 0, wantFunc: "04"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			address, function, err := parseModbusAddress(tt.address, tt.function)
			require.NoError(t, err)
			require.Equal(t, tt.wantAddr, address)
			require.Equal(t, tt.wantFunc, function)
		})
	}
}

func TestPLCModbusAddress_NoFunctionOverflowFailsBeforeWire(t *testing.T) {
	transport := &plcModbusAddressRecordingTransport{}
	protocol := newPLCAdapterModbusConnector("tcp", transport)

	result, err := protocol.Read(t.Context(), connector.ReadRequest{
		Address:  "65536",
		DataType: schema.DataTypeUint16,
	})
	require.Error(t, err)
	require.Equal(t, schema.QualityBad, result.Quality)
	require.Nil(t, result.Value)
	require.Empty(t, transport.requests)
}

type plcModbusAddressWireRequest struct {
	function byte
	address  uint16
}

type plcModbusAddressRecordingTransport struct {
	mu            sync.Mutex
	transactionID uint16
	requests      []plcModbusAddressWireRequest
}

func (t *plcModbusAddressRecordingTransport) Connect() error { return nil }

func (t *plcModbusAddressRecordingTransport) Close() error { return nil }

func (t *plcModbusAddressRecordingTransport) IsConnected() bool { return true }

func (t *plcModbusAddressRecordingTransport) GetNextTransactionID() uint16 {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.transactionID++
	return t.transactionID
}

func (t *plcModbusAddressRecordingTransport) SendReceive(data []byte) ([]byte, error) {
	t.mu.Lock()
	defer t.mu.Unlock()

	transactionID, unitID, function, pdu, err := modbus.ParseTCPFrame(data)
	if err != nil {
		return nil, err
	}
	if len(pdu) < 4 {
		return nil, fmt.Errorf("modbus request too short: %d", len(pdu))
	}
	t.requests = append(t.requests, plcModbusAddressWireRequest{
		function: function,
		address:  binary.BigEndian.Uint16(pdu[:2]),
	})

	var responseData []byte
	switch function {
	case modbus.FuncReadCoils, modbus.FuncReadDiscreteInputs:
		responseData = []byte{1, 1}
	case modbus.FuncReadHoldingRegisters, modbus.FuncReadInputRegisters:
		responseData = []byte{2, 0x12, 0x34}
	case modbus.FuncWriteSingleCoil, modbus.FuncWriteSingleRegister:
		responseData = append([]byte(nil), pdu...)
	default:
		return nil, fmt.Errorf("unsupported Modbus function 0x%02X", function)
	}

	return modbus.BuildTCPFrame(transactionID, unitID, function, responseData), nil
}

func TestPLCModbusAddress_ReadExplicitFunctionThroughAdapters(t *testing.T) {
	tests := []struct {
		name         string
		address      string
		function     string
		wantFunction byte
		wantAddress  uint16
	}{
		{name: "raw holding register", address: "100", function: "03", wantFunction: modbus.FuncReadHoldingRegisters, wantAddress: 100},
		{name: "raw input register at zero", address: "0", function: "04", wantFunction: modbus.FuncReadInputRegisters, wantAddress: 0},
		{name: "raw input register at maximum", address: "65535", function: "04", wantFunction: modbus.FuncReadInputRegisters, wantAddress: 65535},
	}

	for _, protocolName := range []string{"tcp", "rtu", "udp"} {
		for _, tt := range tests {
			t.Run(fmt.Sprintf("%s/%s", protocolName, tt.name), func(t *testing.T) {
				transport := &plcModbusAddressRecordingTransport{}
				protocol := newPLCAdapterModbusConnector(protocolName, transport)

				result, err := protocol.Read(t.Context(), connector.ReadRequest{
					Address:  tt.address,
					Function: tt.function,
					DataType: schema.DataTypeUint16,
				})
				require.NoError(t, err)
				require.Equal(t, schema.QualityGood, result.Quality)
				require.Equal(t, uint16(0x1234), result.Value)
				require.Len(t, transport.requests, 1)
				require.Equal(t, tt.wantFunction, transport.requests[0].function)
				require.Equal(t, tt.wantAddress, transport.requests[0].address)
			})
		}
	}
}

func TestPLCModbusAddress_WriteExplicitFunctionThroughAdapters(t *testing.T) {
	tests := []struct {
		name         string
		function     string
		value        any
		wantFunction byte
	}{
		{name: "single register", function: "06", value: uint16(0x1234), wantFunction: modbus.FuncWriteSingleRegister},
		{name: "single coil", function: "05", value: true, wantFunction: modbus.FuncWriteSingleCoil},
	}

	for _, protocolName := range []string{"tcp", "rtu", "udp"} {
		for _, tt := range tests {
			t.Run(fmt.Sprintf("%s/%s", protocolName, tt.name), func(t *testing.T) {
				transport := &plcModbusAddressRecordingTransport{}
				protocol := newPLCAdapterModbusConnector(protocolName, transport)

				err := protocol.Write(t.Context(), connector.WriteRequest{
					Address:  "100",
					Function: tt.function,
					Value:    tt.value,
				})
				require.NoError(t, err)
				require.Len(t, transport.requests, 1)
				require.Equal(t, tt.wantFunction, transport.requests[0].function)
				require.Equal(t, uint16(100), transport.requests[0].address)
			})
		}
	}
}

func TestPLCModbusAddress_RejectsInvalidRawAddresses(t *testing.T) {
	for _, address := range []string{"-1", "bad", "65536", "4294967296"} {
		t.Run(address, func(t *testing.T) {
			transport := &plcModbusAddressRecordingTransport{}
			protocol := newPLCAdapterModbusConnector("tcp", transport)

			result, err := protocol.Read(t.Context(), connector.ReadRequest{
				Address:  address,
				Function: "03",
				DataType: schema.DataTypeUint16,
			})
			require.Error(t, err)
			require.Equal(t, schema.QualityBad, result.Quality)
			require.Empty(t, transport.requests)
		})
	}
}

func TestPLCModbusAddress_PersistedShortNumericPointPreservesIdentity(t *testing.T) {
	pointID := "point-raw-100"
	deviceID := "device-raw-100"
	pollingGroupID := "group-raw-100"

	pointRepo := point.NewMemoryRepository()
	persistedPoint := &schema.Point{
		ID:             pointID,
		DeviceID:       deviceID,
		Name:           "raw holding register",
		Address:        "100",
		Function:       "03",
		DataType:       schema.DataTypeUint16,
		DataFormat:     "ABCD",
		Mode:           schema.PointModeReadWrite,
		PollingGroupID: &pollingGroupID,
		Enabled:        true,
	}
	require.NoError(t, pointRepo.Create(t.Context(), persistedPoint))
	pointBeforeRead, err := pointRepo.GetByID(t.Context(), pointID)
	require.NoError(t, err)

	mappingRepo := mapping.NewMemoryRepository()
	persistedMapping := &schema.Mapping{
		ID:                "mapping-raw-100",
		PointID:           pointID,
		TagID:             "tag-raw-100",
		TransformPipeline: "[]",
		Status:            schema.MappingStatusActive,
		Enabled:           true,
	}
	require.NoError(t, mappingRepo.Create(t.Context(), persistedMapping))
	mappingsBeforeRead, err := mappingRepo.GetByPointID(t.Context(), pointID)
	require.NoError(t, err)
	require.Len(t, mappingsBeforeRead, 1)

	transport := &plcModbusAddressRecordingTransport{}
	protocol := newPLCAdapterModbusConnector("tcp", transport)
	result, err := protocol.Read(t.Context(), connector.ReadRequest{
		Address:    pointBeforeRead.Address,
		Function:   pointBeforeRead.Function,
		DataType:   pointBeforeRead.DataType,
		DataFormat: pointBeforeRead.DataFormat,
	})
	require.NoError(t, err)
	require.Equal(t, schema.QualityGood, result.Quality)
	require.Equal(t, uint16(0x1234), result.Value)
	require.Len(t, transport.requests, 1)
	require.Equal(t, byte(modbus.FuncReadHoldingRegisters), transport.requests[0].function)
	require.Equal(t, uint16(100), transport.requests[0].address)

	pointAfterRead, err := pointRepo.GetByID(t.Context(), pointID)
	require.NoError(t, err)
	require.Equal(t, pointBeforeRead, pointAfterRead)
	mappingsAfterRead, err := mappingRepo.GetByPointID(t.Context(), pointID)
	require.NoError(t, err)
	require.Equal(t, mappingsBeforeRead, mappingsAfterRead)
}
