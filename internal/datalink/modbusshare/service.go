package modbusshare

import (
	"context"
	"fmt"
	"sync"
	"time"

	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/tag"
	"go-gateway/internal/virtual/memory"
	virtualmodbus "go-gateway/internal/virtual/server/modbus"
)

const modbusMaxRegs = 65536

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
