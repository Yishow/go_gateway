package dbtarget

import (
	"context"
	"sync"

	"go-gateway/internal/datalink/schema"
)

type ConnectorTagReader interface {
	GetByID(ctx context.Context, id string) (*schema.Tag, error)
}

var connectorTagReaders sync.Map

func (s *ConnectorService) SetTagReader(reader ConnectorTagReader) {
	if reader == nil {
		connectorTagReaders.Delete(s)
		return
	}
	connectorTagReaders.Store(s, reader)
}

func (s *ConnectorService) connectorTagReader() ConnectorTagReader {
	reader, ok := connectorTagReaders.Load(s)
	if !ok {
		return nil
	}
	tagReader, _ := reader.(ConnectorTagReader)
	return tagReader
}
