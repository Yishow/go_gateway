package dbtarget

import (
	"context"
	"fmt"

	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/tag"
)

type refreshTagList interface {
	List(context.Context, tag.ListFilter) ([]*schema.Tag, error)
}

type refreshTagResult struct {
	entity *schema.Tag
	err    error
}

type refreshTagReader struct {
	source  ConnectorTagReader
	entries map[string]refreshTagResult
	batch   bool
}

func newRefreshTagReader(ctx context.Context, reader ConnectorTagReader) (ConnectorTagReader, error) {
	if reader == nil {
		return nil, nil
	}
	cached := &refreshTagReader{source: reader, entries: make(map[string]refreshTagResult)}
	if lister, ok := reader.(refreshTagList); ok {
		entities, err := lister.List(ctx, tag.ListFilter{})
		if err != nil {
			return nil, fmt.Errorf("批次讀取標籤失敗: %w", err)
		}
		cached.batch = true
		for _, entity := range entities {
			if entity != nil {
				cached.entries[entity.ID] = refreshTagResult{entity: entity}
			}
		}
	}
	return cached, nil
}

func (r *refreshTagReader) GetByID(ctx context.Context, id string) (*schema.Tag, error) {
	if result, ok := r.entries[id]; ok {
		return result.entity, result.err
	}
	if r.batch {
		return nil, tag.ErrTagNotFound
	}
	entity, err := r.source.GetByID(ctx, id)
	r.entries[id] = refreshTagResult{entity: entity, err: err}
	return entity, err
}
