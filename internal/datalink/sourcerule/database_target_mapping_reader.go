package sourcerule

import (
	"context"
	"sync"

	"go-gateway/internal/datalink/schema"
)

var databaseTargetMappingReaders sync.Map

type DatabaseTargetMappingReader interface {
	List(ctx context.Context) ([]*schema.DatabaseTargetMapping, error)
}

type DatabaseTargetMappingListFunc func(ctx context.Context) ([]*schema.DatabaseTargetMapping, error)

func (fn DatabaseTargetMappingListFunc) List(ctx context.Context) ([]*schema.DatabaseTargetMapping, error) {
	return fn(ctx)
}

func (s *Service) SetDatabaseTargetMappingReader(reader DatabaseTargetMappingReader) {
	if reader == nil {
		databaseTargetMappingReaders.Delete(s)
		return
	}
	databaseTargetMappingReaders.Store(s, reader)
}

func (s *Service) databaseTargetMappingReader() DatabaseTargetMappingReader {
	reader, ok := databaseTargetMappingReaders.Load(s)
	if !ok {
		return nil
	}
	mappingReader, _ := reader.(DatabaseTargetMappingReader)
	return mappingReader
}
