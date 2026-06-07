package dbtarget

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/tag"
)

type liveTargetProjection struct {
	Mappings        []*schema.DatabaseTargetMapping
	TagsByMappingID map[string]*schema.Tag
	Excluded        []liveTargetProjectionExclusion
}

type liveTargetProjectionExclusion struct {
	MappingID string
	TagID     string
	Code      string
}

func listLiveTargetProjection(
	ctx context.Context,
	repo TargetMappingRepository,
	filter TargetMappingListFilter,
	tagReader ConnectorTagReader,
) (*liveTargetProjection, error) {
	enabled := true
	filter.Enabled = &enabled
	mappings, err := repo.List(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("列出資料庫目標映射失敗: %w", err)
	}
	return filterLiveTargetProjection(ctx, mappings, tagReader)
}

func filterLiveTargetProjection(
	ctx context.Context,
	mappings []*schema.DatabaseTargetMapping,
	tagReader ConnectorTagReader,
) (*liveTargetProjection, error) {
	projection := &liveTargetProjection{
		Mappings:        make([]*schema.DatabaseTargetMapping, 0, len(mappings)),
		TagsByMappingID: make(map[string]*schema.Tag, len(mappings)),
		Excluded:        make([]liveTargetProjectionExclusion, 0),
	}

	for _, mappingRecord := range mappings {
		if mappingRecord == nil || !mappingRecord.Enabled {
			continue
		}
		if strings.TrimSpace(mappingRecord.TagID) == "" {
			projection.Excluded = append(projection.Excluded, liveTargetProjectionExclusion{
				MappingID: mappingRecord.ID,
				TagID:     mappingRecord.TagID,
				Code:      "tag_missing",
			})
			continue
		}
		if tagReader == nil {
			projection.Mappings = append(projection.Mappings, mappingRecord)
			continue
		}

		tagEntity, err := tagReader.GetByID(ctx, mappingRecord.TagID)
		if err != nil {
			if errors.Is(err, tag.ErrTagNotFound) {
				projection.Excluded = append(projection.Excluded, liveTargetProjectionExclusion{
					MappingID: mappingRecord.ID,
					TagID:     mappingRecord.TagID,
					Code:      "tag_missing",
				})
				continue
			}
			return nil, fmt.Errorf("取得標籤失敗: %w", err)
		}

		projection.Mappings = append(projection.Mappings, mappingRecord)
		projection.TagsByMappingID[mappingRecord.ID] = tagEntity
	}

	return projection, nil
}
