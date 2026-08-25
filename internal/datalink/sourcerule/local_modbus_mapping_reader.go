package sourcerule

import (
	"context"
	"sync"
	"time"

	"go-gateway/internal/datalink/schema"
)

var localModbusMappingReaders sync.Map

// LocalModbusMappingRecord is the persisted Local Modbus mapping metadata used
// to rebuild source-rule candidate projections.
type LocalModbusMappingRecord struct {
	MappingID string
	TagID     string
	Register  uint16
	DataType  schema.DataType
	UpdatedAt time.Time
}

// LocalModbusMappingReader lists persisted Local Modbus mappings.
type LocalModbusMappingReader interface {
	List(ctx context.Context) ([]LocalModbusMappingRecord, error)
}

// LocalModbusMappingListFunc adapts a function to LocalModbusMappingReader.
type LocalModbusMappingListFunc func(ctx context.Context) ([]LocalModbusMappingRecord, error)

// List invokes the adapted Local Modbus mapping listing function.
func (fn LocalModbusMappingListFunc) List(ctx context.Context) ([]LocalModbusMappingRecord, error) {
	return fn(ctx)
}

// LocalModbusMappingWriter persists Local Modbus mapping metadata.
type LocalModbusMappingWriter interface {
	Upsert(ctx context.Context, mappingRecord LocalModbusMappingRecord) error
}

// LocalModbusMappingUpsertFunc adapts a function to LocalModbusMappingWriter.
type LocalModbusMappingUpsertFunc func(ctx context.Context, mappingRecord LocalModbusMappingRecord) error

// Upsert invokes the adapted Local Modbus mapping persistence function.
func (fn LocalModbusMappingUpsertFunc) Upsert(ctx context.Context, mappingRecord LocalModbusMappingRecord) error {
	return fn(ctx, mappingRecord)
}

// SetLocalModbusMappingReader configures the persisted Local Modbus mapping
// reader used by candidate snapshot generation.
func (s *Service) SetLocalModbusMappingReader(reader LocalModbusMappingReader) {
	if reader == nil {
		localModbusMappingReaders.Delete(s)
		return
	}
	localModbusMappingReaders.Store(s, reader)
}

func (s *Service) localModbusMappingReader() LocalModbusMappingReader {
	reader, ok := localModbusMappingReaders.Load(s)
	if !ok {
		return nil
	}
	mappingReader, ok := reader.(LocalModbusMappingReader)
	if !ok {
		return nil
	}
	return mappingReader
}
