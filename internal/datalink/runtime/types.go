package runtime

import "go-gateway/internal/datalink/schema"

// Snapshot describes the static runtime resources required by datalink pipeline.
type Snapshot struct {
	WorkspaceID   string
	Devices       []*schema.Device
	Points        []*schema.Point
	PollingGroups []*schema.PollingGroup
	Mappings      []*schema.Mapping
	Tags          []*schema.Tag
}

func (s Snapshot) workspaceID(fallback string) string {
	if s.WorkspaceID != "" {
		return s.WorkspaceID
	}
	return fallback
}

// Metrics exposes ingestion counters used by load tests and health checks.
type Metrics struct {
	TotalReads        int64
	TotalWrites       int64
	ErrorCount        int64
	MappingErrorCount int64
	WriteErrorCount   int64
}
