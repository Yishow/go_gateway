package sourcerule

import (
	"context"
	"sync"
	"time"

	"go-gateway/internal/datalink/schema"
)

var localModbusMappingReaders sync.Map

type LocalModbusMappingRecord struct {
	TagID     string
	Register  uint16
	DataType  schema.DataType
	UpdatedAt time.Time
}

type LocalModbusMappingReader interface {
	List(ctx context.Context) ([]LocalModbusMappingRecord, error)
}

type LocalModbusMappingListFunc func(ctx context.Context) ([]LocalModbusMappingRecord, error)

func (fn LocalModbusMappingListFunc) List(ctx context.Context) ([]LocalModbusMappingRecord, error) {
	return fn(ctx)
}

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
