package workspace

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
)

var ErrNotFound = errors.New("studio v2 workspace not found")

type Kind string

const (
	WorkspaceKindSingle Kind = "single"
)

type Status string

const (
	WorkspaceStatusEmpty Status = "empty"
	WorkspaceStatusReady Status = "ready"
)

type Record struct {
	ID                  string `json:"id"`
	Kind                Kind   `json:"kind"`
	Status              Status `json:"status"`
	DatabaseConnectorID string `json:"database_connector_id,omitempty"`
	// DatabaseSetupRevision advances with every Step 4 database setup change so
	// stale browser saves are rejected before they mutate anything.
	DatabaseSetupRevision string              `json:"database_setup_revision,omitempty"`
	DatabaseRowGroups     []DatabaseRowGroup  `json:"database_row_groups,omitempty"`
	DatabaseTargetRefs    []DatabaseTargetRef `json:"database_target_refs,omitempty"`
	OrderedDeviceIDs      []string            `json:"ordered_device_ids"`
	CreatedAt             time.Time           `json:"created_at"`
	UpdatedAt             time.Time           `json:"updated_at"`
}

type Repository interface {
	Get(ctx context.Context) (*Record, error)
	Save(ctx context.Context, record *Record) error
}

// MembershipScopeInvalidator advances the formal Share scope when device
// membership changes. It is intentionally narrower than the Share service.
type MembershipScopeInvalidator interface {
	InvalidateWorkspaceMembership(context.Context, string) error
}

// Service 提供 Workspace 的讀寫操作。
// 所有 read-modify-write 方法持有 mu 期間執行完整週期，
// 確保單程序內的 goroutine 不會互相覆寫。
type Service struct {
	mu                   sync.Mutex
	repo                 Repository
	now                  func() time.Time
	newID                func() string
	readinessDevices     readinessDeviceService
	readinessRules       readinessSourceRuleService
	readinessConnectors  readinessConnectorService
	readinessMappings    readinessDBTargetService
	readinessTags        readinessTagService
	readinessLinkRecords readinessMappingRecordService
	projectionDevices    runtimeProjectionDeviceService
	projectionRules      runtimeProjectionSourceRuleService
	projectionPoints     runtimeProjectionPointService
	projectionMappings   runtimeProjectionMappingService
	projectionTags       runtimeProjectionTagService
	projectionConnectors runtimeProjectionConnectorService
	projectionTargets    runtimeProjectionTargetService
	projectionGroups     runtimeProjectionPollingGroupService
	shareScope           MembershipScopeInvalidator
}

// SetMembershipScopeInvalidator wires the optional Local Modbus Share scope
// barrier without coupling workspace persistence to Share lifecycle internals.
func (s *Service) SetMembershipScopeInvalidator(invalidator MembershipScopeInvalidator) {
	s.mu.Lock()
	s.shareScope = invalidator
	s.mu.Unlock()
}

func NewService(repo Repository) *Service {
	return &Service{
		repo: repo,
		now: func() time.Time {
			return time.Now().UTC()
		},
		newID: uuid.NewString,
	}
}

// GetOrCreate 返回工作區的當前快照（mutex 保護，並行安全）。
// 若需要 read-modify-write（如新增/移除設備），請使用 AttachDevice、DetachDevice
// 等 Service 方法，而非自行修改快照後另行儲存（那會繞過互斥保護造成 lost update）。
func (s *Service) GetOrCreate(ctx context.Context) (*Record, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.getOrCreate(ctx)
}

// ShareOwnershipSnapshot exposes only the persisted workspace scope needed by
// the Local Modbus ownership seam, avoiding coupling callers to workspace
// mutation details.
func (s *Service) ShareOwnershipSnapshot(ctx context.Context) (workspaceID string, deviceIDs []string, err error) {
	record, err := s.GetOrCreate(ctx)
	if err != nil {
		return "", nil, err
	}
	return record.ID, append([]string(nil), record.OrderedDeviceIDs...), nil
}

// getOrCreate 為內部使用：必須在 s.mu 已鎖定的情況下呼叫。
func (s *Service) getOrCreate(ctx context.Context) (*Record, error) {
	record, err := s.repo.Get(ctx)
	if err == nil {
		return cloneRecord(record), nil
	}
	if !errors.Is(err, ErrNotFound) {
		return nil, fmt.Errorf("read studio v2 workspace: %w", err)
	}

	record = s.newRecord()
	if err := s.repo.Save(ctx, record); err != nil {
		return nil, fmt.Errorf("create studio v2 workspace: %w", err)
	}

	return cloneRecord(record), nil
}

func cloneRecord(record *Record) *Record {
	if record == nil {
		return nil
	}

	cloned := *record
	cloned.OrderedDeviceIDs = append([]string{}, record.OrderedDeviceIDs...)
	cloned.DatabaseRowGroups = cloneDatabaseRowGroups(record.DatabaseRowGroups)
	cloned.DatabaseTargetRefs = cloneDatabaseTargetRefs(record.DatabaseTargetRefs)

	return &cloned
}
