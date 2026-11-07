package workspace

import (
	"context"
	"log"

	"go-gateway/internal/datalink/audit"
)

type writeGroupAuditRecorder interface {
	Record(context.Context, audit.RecordEvent) error
}

// WithAudit records successful permanent deletions without affecting their commit.
func (s *WriteGroupService) WithAudit(svc *audit.Service) *WriteGroupService {
	if s != nil && svc != nil {
		s.audit = svc
	}
	return s
}

func (s *WriteGroupService) recordPermanentDeleteAudit(ctx context.Context, result *WriteGroupSaveResult) {
	if result == nil || result.Group == nil {
		return
	}
	if s.audit == nil {
		log.Printf("record write group deletion audit unavailable")
		return
	}
	if err := s.audit.Record(ctx, audit.RecordEvent{WorkspaceID: result.Group.WorkspaceID, EventType: audit.EventTypeWriteGroupDeleted,
		Result: audit.ResultSuccess, Scope: "workspace", ReferenceID: result.GroupID,
		Details: map[string]any{"name": result.Group.Name, "lineage_id": result.LineageID, "deleted_generations": result.DeletedGenerations}}); err != nil {
		log.Printf("record write group deletion audit failed")
	}
}
