package modbusshare

import (
	"context"
	"encoding/binary"
	"fmt"
	"math"
	"net"
	"strings"
	"sync"
	"time"

	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/tag"
	"go-gateway/internal/virtual/memory"
	virtualmodbus "go-gateway/internal/virtual/server/modbus"
)

// TagMirrorMapping defines one tag-to-register mirror binding.
type TagMirrorMapping struct {
	TagID     string          `json:"tag_id"`
	Register  uint16          `json:"register"`
	DataType  schema.DataType `json:"data_type"`
	UpdatedAt time.Time       `json:"updated_at"`
}

// Status describes runtime state of the local Modbus share server.
type Status struct {
	Enabled      bool   `json:"enabled"`
	Port         int    `json:"port"`
	Address      string `json:"address"`
	BindState    string `json:"bind_state"`
	MappingCount int    `json:"mapping_count"`
}

// Service mirrors tag values into virtual Modbus memory and serves them over Modbus TCP.
type Service struct {
	mu       sync.RWMutex
	tagSvc   *tag.Service
	bank     *memory.MemoryBank
	server   *virtualmodbus.Server
	mappings map[string]TagMirrorMapping
}

// NewService creates a new Modbus share service.
func NewService(tagSvc *tag.Service, memorySizeBytes int) *Service {
	if memorySizeBytes <= 0 {
		memorySizeBytes = 65536 // 32K registers
	}

	bank := memory.NewMemoryBank(memorySizeBytes)
	return &Service{
		tagSvc:   tagSvc,
		bank:     bank,
		server:   virtualmodbus.NewServer(bank),
		mappings: make(map[string]TagMirrorMapping),
	}
}

// Start starts local Modbus TCP server on specified port.
func (s *Service) Start(port int) error {
	if port <= 0 {
		return fmt.Errorf("invalid port: %d", port)
	}
	if err := preflightPortAvailable(port); err != nil {
		return err
	}
	if err := s.server.Start(port); err != nil {
		if isAddressInUseError(err) {
			return fmt.Errorf("port %d is already in use, please release the port or change service port: %w", port, err)
		}
		return fmt.Errorf("start modbus share server failed on port %d: %w", port, err)
	}
	return nil
}

// Stop stops the local Modbus server.
func (s *Service) Stop() error {
	return s.server.Stop()
}

// Status returns runtime status.
func (s *Service) Status() Status {
	s.mu.RLock()
	defer s.mu.RUnlock()

	port := s.server.Port()
	enabled := port > 0
	addr := ""
	if enabled {
		addr = s.server.Address()
	}
	bindState := "fail"
	if enabled && addr != "" {
		bindState = "pass"
	}

	return Status{
		Enabled:      enabled,
		Port:         port,
		Address:      addr,
		BindState:    bindState,
		MappingCount: len(s.mappings),
	}
}

// UpsertMapping creates or updates tag-to-register mapping.
func (s *Service) UpsertMapping(ctx context.Context, tagID string, register uint16) (*TagMirrorMapping, error) {
	if tagID == "" {
		return nil, fmt.Errorf("tag_id is required")
	}

	t, err := s.tagSvc.GetByID(ctx, tagID)
	if err != nil {
		return nil, fmt.Errorf("tag not found: %w", err)
	}

	m := TagMirrorMapping{
		TagID:     tagID,
		Register:  register,
		DataType:  t.DataType,
		UpdatedAt: time.Now().UTC(),
	}

	s.mu.Lock()
	s.mappings[tagID] = m
	s.mu.Unlock()

	return &m, nil
}

// RemoveMapping removes an existing mapping.
func (s *Service) RemoveMapping(tagID string) {
	s.mu.Lock()
	delete(s.mappings, tagID)
	s.mu.Unlock()
}

// ListMappings lists all mappings.
func (s *Service) ListMappings() []TagMirrorMapping {
	s.mu.RLock()
	defer s.mu.RUnlock()

	out := make([]TagMirrorMapping, 0, len(s.mappings))
	for _, m := range s.mappings {
		out = append(out, m)
	}
	return out
}

// HasMapping checks whether a tag mirror mapping exists.
func (s *Service) HasMapping(tagID string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	_, ok := s.mappings[tagID]
	return ok
}

// WriteTagValue writes tag value to mapped Modbus registers.
func (s *Service) WriteTagValue(ctx context.Context, tagID string, value interface{}) error {
	s.mu.RLock()
	mapping, ok := s.mappings[tagID]
	s.mu.RUnlock()
	if !ok {
		return fmt.Errorf("mapping not found for tag_id: %s", tagID)
	}

	t, err := s.tagSvc.GetByID(ctx, tagID)
	if err != nil {
		return fmt.Errorf("tag not found: %w", err)
	}

	words, err := encodeToWords(t.DataType, value)
	if err != nil {
		return err
	}

	for i, word := range words {
		offset := int(mapping.Register+uint16(i)) * 2
		if err := s.bank.WriteWord(offset, word); err != nil {
			return fmt.Errorf("write register %d failed: %w", int(mapping.Register)+i, err)
		}
	}

	return nil
}

func encodeToWords(dataType schema.DataType, value interface{}) ([]uint16, error) {
	switch dataType {
	case schema.DataTypeBool:
		b, err := toBool(value)
		if err != nil {
			return nil, err
		}
		if b {
			return []uint16{1}, nil
		}
		return []uint16{0}, nil

	case schema.DataTypeInt16:
		v, err := toInt64(value)
		if err != nil {
			return nil, err
		}
		return []uint16{uint16(int16(v))}, nil

	case schema.DataTypeUint16:
		v, err := toUint64(value)
		if err != nil {
			return nil, err
		}
		return []uint16{uint16(v)}, nil

	case schema.DataTypeInt32:
		v, err := toInt64(value)
		if err != nil {
			return nil, err
		}
		raw := uint32(int32(v))
		return []uint16{uint16(raw >> 16), uint16(raw)}, nil

	case schema.DataTypeUint32:
		v, err := toUint64(value)
		if err != nil {
			return nil, err
		}
		raw := uint32(v)
		return []uint16{uint16(raw >> 16), uint16(raw)}, nil

	case schema.DataTypeFloat32:
		f, err := toFloat64(value)
		if err != nil {
			return nil, err
		}
		raw := math.Float32bits(float32(f))
		return []uint16{uint16(raw >> 16), uint16(raw)}, nil

	case schema.DataTypeInt64:
		v, err := toInt64(value)
		if err != nil {
			return nil, err
		}
		raw := uint64(v)
		return []uint16{
			uint16(raw >> 48),
			uint16(raw >> 32),
			uint16(raw >> 16),
			uint16(raw),
		}, nil

	case schema.DataTypeUint64:
		raw, err := toUint64(value)
		if err != nil {
			return nil, err
		}
		return []uint16{
			uint16(raw >> 48),
			uint16(raw >> 32),
			uint16(raw >> 16),
			uint16(raw),
		}, nil

	case schema.DataTypeFloat64:
		f, err := toFloat64(value)
		if err != nil {
			return nil, err
		}
		raw := math.Float64bits(f)
		return []uint16{
			uint16(raw >> 48),
			uint16(raw >> 32),
			uint16(raw >> 16),
			uint16(raw),
		}, nil

	default:
		return nil, fmt.Errorf("unsupported data_type for Modbus mirror: %s", dataType)
	}
}

func toInt64(v interface{}) (int64, error) {
	switch n := v.(type) {
	case int:
		return int64(n), nil
	case int8:
		return int64(n), nil
	case int16:
		return int64(n), nil
	case int32:
		return int64(n), nil
	case int64:
		return n, nil
	case uint:
		return int64(n), nil
	case uint8:
		return int64(n), nil
	case uint16:
		return int64(n), nil
	case uint32:
		return int64(n), nil
	case uint64:
		return int64(n), nil
	case float32:
		return int64(n), nil
	case float64:
		return int64(n), nil
	default:
		return 0, fmt.Errorf("value is not numeric")
	}
}

func toUint64(v interface{}) (uint64, error) {
	switch n := v.(type) {
	case int:
		return uint64(n), nil
	case int8:
		return uint64(n), nil
	case int16:
		return uint64(n), nil
	case int32:
		return uint64(n), nil
	case int64:
		return uint64(n), nil
	case uint:
		return uint64(n), nil
	case uint8:
		return uint64(n), nil
	case uint16:
		return uint64(n), nil
	case uint32:
		return uint64(n), nil
	case uint64:
		return n, nil
	case float32:
		return uint64(n), nil
	case float64:
		return uint64(n), nil
	default:
		return 0, fmt.Errorf("value is not numeric")
	}
}

func toFloat64(v interface{}) (float64, error) {
	switch n := v.(type) {
	case int:
		return float64(n), nil
	case int8:
		return float64(n), nil
	case int16:
		return float64(n), nil
	case int32:
		return float64(n), nil
	case int64:
		return float64(n), nil
	case uint:
		return float64(n), nil
	case uint8:
		return float64(n), nil
	case uint16:
		return float64(n), nil
	case uint32:
		return float64(n), nil
	case uint64:
		return float64(n), nil
	case float32:
		return float64(n), nil
	case float64:
		return n, nil
	default:
		return 0, fmt.Errorf("value is not numeric")
	}
}

func toBool(v interface{}) (bool, error) {
	switch b := v.(type) {
	case bool:
		return b, nil
	case int:
		return b != 0, nil
	case int64:
		return b != 0, nil
	case float64:
		return b != 0, nil
	default:
		return false, fmt.Errorf("value is not boolean")
	}
}

// ReadHoldingWords is a helper for tests and diagnostics.
func (s *Service) ReadHoldingWords(startRegister uint16, quantity uint16) ([]uint16, error) {
	if quantity == 0 {
		return []uint16{}, nil
	}
	buf := make([]byte, quantity*2)
	for i := uint16(0); i < quantity; i++ {
		offset := int(startRegister+i) * 2
		w, err := s.bank.ReadWord(offset)
		if err != nil {
			return nil, err
		}
		binary.BigEndian.PutUint16(buf[i*2:], w)
	}
	out := make([]uint16, quantity)
	for i := range out {
		out[i] = binary.BigEndian.Uint16(buf[i*2:])
	}
	return out, nil
}

func preflightPortAvailable(port int) error {
	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		if isAddressInUseError(err) {
			return fmt.Errorf("port %d is already in use, please stop the conflicting process or choose another port", port)
		}
		return fmt.Errorf("cannot bind port %d: %w", port, err)
	}
	return ln.Close()
}

func isAddressInUseError(err error) bool {
	if err == nil {
		return false
	}
	msg := strings.ToLower(err.Error())
	return strings.Contains(msg, "address already in use") ||
		strings.Contains(msg, "only one usage of each socket address")
}
