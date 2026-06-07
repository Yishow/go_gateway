package audit

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

const (
	defaultListLimit = 20
	maxListLimit     = 100
)

// Repository persists workspace audit history entries.
type Repository interface {
	Create(ctx context.Context, entry *Entry) error
	List(ctx context.Context, filter ListFilter) ([]Entry, error)
}

// Service records and queries workspace audit history.
type Service struct {
	repo  Repository
	now   func() time.Time
	newID func() string
}

// NewService creates a workspace audit history service.
func NewService(repo Repository) *Service {
	return &Service{
		repo:  repo,
		now:   func() time.Time { return time.Now().UTC() },
		newID: uuid.NewString,
	}
}

// Record persists one workspace audit history event.
func (s *Service) Record(ctx context.Context, event RecordEvent) error {
	if s == nil || s.repo == nil {
		return fmt.Errorf("audit repository unavailable")
	}

	workspaceID := strings.TrimSpace(event.WorkspaceID)
	if workspaceID == "" {
		return fmt.Errorf("workspace_id is required")
	}
	if event.EventType == "" {
		return fmt.Errorf("event_type is required")
	}
	if event.Result == "" {
		return fmt.Errorf("result is required")
	}

	occurredAt := event.OccurredAt.UTC()
	if occurredAt.IsZero() {
		occurredAt = s.now()
	}
	createdAt := s.now()

	entry := &Entry{
		ID:          s.newID(),
		WorkspaceID: workspaceID,
		EventType:   event.EventType,
		Result:      event.Result,
		Scope:       strings.TrimSpace(event.Scope),
		ReferenceID: strings.TrimSpace(event.ReferenceID),
		Details:     encodeDetails(event.Details),
		OccurredAt:  occurredAt,
		CreatedAt:   createdAt,
	}
	if entry.Scope == "" {
		entry.Scope = "workspace"
	}

	return s.repo.Create(ctx, entry)
}

// List returns recent audit history entries for one workspace.
func (s *Service) List(ctx context.Context, filter ListFilter) ([]Entry, error) {
	if s == nil || s.repo == nil {
		return nil, fmt.Errorf("audit repository unavailable")
	}

	filter.WorkspaceID = strings.TrimSpace(filter.WorkspaceID)
	if filter.WorkspaceID == "" {
		return nil, fmt.Errorf("workspace_id is required")
	}
	filter.Limit = normalizeLimit(filter.Limit)

	return s.repo.List(ctx, filter)
}

func normalizeLimit(limit int) int {
	if limit <= 0 {
		return defaultListLimit
	}
	if limit > maxListLimit {
		return maxListLimit
	}
	return limit
}

func encodeDetails(details any) string {
	if details == nil {
		return "{}"
	}
	encoded, err := json.Marshal(details)
	if err != nil {
		return "{}"
	}
	return string(encoded)
}
