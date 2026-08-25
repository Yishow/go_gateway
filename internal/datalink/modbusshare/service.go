package modbusshare

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"

	"go-gateway/internal/datalink/schema"
	datalinksettings "go-gateway/internal/datalink/settings"
	"go-gateway/internal/datalink/tag"
	"go-gateway/internal/virtual/memory"
	virtualmodbus "go-gateway/internal/virtual/server/modbus"
)

// modbusMaxRegs is the maximum capacity of the in-process holding-register
// bank. The bank is byte-addressed and a holding register is two bytes, so the
// 64 KiB bank exposes 32 Ki registers (not 64 Ki registers).
const modbusMaxRegs = 32768

const settingsKey = datalinksettings.KeyModbusShare

// TagMirrorMapping defines one tag-to-register mirror binding.
type TagMirrorMapping struct {
	WorkspaceID        string          `json:"workspace_id,omitempty"`
	SourceRuleID       string          `json:"source_rule_id,omitempty"`
	SourceRuleRevision string          `json:"source_rule_revision,omitempty"`
	TagID              string          `json:"tag_id"`
	MappingID          string          `json:"mapping_id,omitempty"`
	Register           uint16          `json:"register"`
	ShareStartRegister uint32          `json:"share_start_register"`
	ZeroBasedRegister  uint16          `json:"zero_based_register"`
	SpanRegisters      int             `json:"span_registers"`
	StrideRegisters    int             `json:"stride_registers"`
	CapacityRegisters  int             `json:"capacity_registers"`
	DataType           schema.DataType `json:"data_type"`
	TagKey             string          `json:"tag_key,omitempty"`
	DisplayName        string          `json:"display_name,omitempty"`
	UpdatedAt          time.Time       `json:"updated_at"`
}

// Service mirrors tag values into virtual Modbus memory and serves them over Modbus TCP.
type Service struct {
	lifecycleMu            sync.Mutex
	lifecycleEpoch         uint64
	mu                     sync.RWMutex
	tagSvc                 *tag.Service
	bank                   *memory.MemoryBank
	server                 *virtualmodbus.Server
	settingsRepo           datalinksettings.Repository
	mappings               map[string]TagMirrorMapping
	settings               Settings
	state                  string
	authoritativeSettings  bool
	hydration              HydrationState
	hydrationErr           error
	lastFailure            *Diagnostic
	ownership              OwnershipChecker
	desiredOwnership       DesiredOwnershipChecker
	restoreMemoryHook      func([]byte) error
	projectionSwapHook     func(map[string]TagMirrorMapping) error
	workspaceRevisionStore WorkspaceRevisionStore
}

// SetWorkspaceRevisionStore wires the durable workspace CAS boundary used by
// membership invalidation and reconciler operations.
func (s *Service) SetWorkspaceRevisionStore(store WorkspaceRevisionStore) {
	s.mu.Lock()
	s.workspaceRevisionStore = store
	s.mu.Unlock()
}

func (s *Service) lifecycleSnapshot() (settings Settings, authoritative bool, epoch uint64) {
	s.lifecycleMu.Lock()
	defer s.lifecycleMu.Unlock()
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.settings, s.authoritativeSettings, s.lifecycleEpoch
}

// SetProjectionSwapHook injects a runtime projection failure for integration
// tests and leaves production projection behavior unchanged.
func (s *Service) SetProjectionSwapHook(hook func(map[string]TagMirrorMapping) error) {
	s.mu.Lock()
	s.projectionSwapHook = hook
	s.mu.Unlock()
}

func (s *Service) projectionSwapError(staged map[string]TagMirrorMapping) error {
	s.mu.RLock()
	hook := s.projectionSwapHook
	s.mu.RUnlock()
	if hook == nil {
		return nil
	}
	return hook(staged)
}

// SetSettingsRepository connects the service to the durable system settings seam.
// The repository is intentionally injected so process-local snapshots cannot be
// mistaken for the persisted Share configuration.
func (s *Service) SetSettingsRepository(repo datalinksettings.Repository) error {
	s.mu.Lock()
	s.settingsRepo = repo
	s.mu.Unlock()
	return nil
}

// NewService creates a new Modbus share service.
func NewService(tagSvc *tag.Service, memorySizeBytes int) *Service {
	if memorySizeBytes <= 0 {
		memorySizeBytes = modbusMaxRegs * 2
	}

	bank := memory.NewMemoryBank(memorySizeBytes)
	defaults := DefaultSettings()
	if capacity := bank.Size() / 2; defaults.CapacityRegisters > capacity {
		defaults.CapacityRegisters = capacity
	}
	return &Service{
		tagSvc:    tagSvc,
		bank:      bank,
		server:    virtualmodbus.NewServer(bank),
		mappings:  make(map[string]TagMirrorMapping),
		settings:  defaults,
		state:     shareStateDisabled,
		hydration: HydrationState{State: HydrationStatePending, Readiness: false},
	}
}

func (s *Service) bankCapacityRegisters() int {
	return s.bank.Size() / 2
}

func (s *Service) effectiveCapacityRegisters() int {
	capRegs := s.bankCapacityRegisters()
	s.mu.RLock()
	configured := s.settings.CapacityRegisters
	s.mu.RUnlock()
	if configured > 0 && configured < capRegs {
		return configured
	}
	return capRegs
}

// Settings returns the latest settings snapshot used by lifecycle operations.
func (s *Service) Settings() Settings {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.settings
}

// SettingsAuthoritative reports whether lifecycle settings are backed by durable storage.
func (s *Service) SettingsAuthoritative() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.authoritativeSettings
}

// SetHydrationState updates the backend bootstrap gate.
func (s *Service) SetHydrationState(state HydrationState) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if (state.State != HydrationStateReady || !state.Readiness) && s.state == shareStateRunning {
		if err := s.server.Stop(); err != nil {
			s.state = HydrationStateFailed
		} else {
			s.state = shareStateStopped
		}
	}
	s.hydration = state
}

// InvalidateWorkspaceMembership advances the durable Share scope barrier when
// workspace device membership changes. Until bootstrap hydrates the new
// membership, candidate/apply requests fail closed as hydration_required.
func (s *Service) InvalidateWorkspaceMembership(ctx context.Context, workspaceID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	workspaceRevisionStore := s.workspaceRevisionStore
	state := s.hydration
	state.WorkspaceID = workspaceID
	state.WorkspaceRevision = uuid.NewString()
	// Membership is persisted before this boundary is invoked. Once the
	// durable revision CAS succeeds, the new workspace scope is authoritative
	// and can be exposed as the rehydrated ready barrier immediately.
	state.State = HydrationStateReady
	state.Readiness = true
	state.ReadinessToken = uuid.NewString()
	if workspaceRevisionStore != nil {
		current, _, err := workspaceRevisionStore.GetRevision(ctx, workspaceID)
		if err != nil {
			return err
		}
		if err := workspaceRevisionStore.UpdateRevision(ctx, workspaceID, current, state.WorkspaceRevision); err != nil {
			return err
		}
	}
	s.hydration = state
	return nil
}

// CheckHydration implements the handler Share gate.
func (s *Service) CheckHydration(_ context.Context) (HydrationState, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.hydration, s.hydrationErr
}

// SetHydrationError installs a deterministic bootstrap-read failure seam.
// Production bootstrap leaves this unset; tests and recovery probes can use it
// to verify that consumers fail closed when the status read itself fails.
func (s *Service) SetHydrationError(err error) {
	s.mu.Lock()
	s.hydrationErr = err
	s.mu.Unlock()
}

// GetSettings implements the handler Share gate.
func (s *Service) GetSettings(_ context.Context) (Settings, error) { return s.Settings(), nil }

// SetOwnershipChecker injects the durable workspace relationship query.
func (s *Service) SetOwnershipChecker(checker OwnershipChecker) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.ownership = checker
}

// SetDesiredMappingOwnershipChecker installs the identity-aware ownership
// proof used when projecting mappings back to an API workspace response.
func (s *Service) SetDesiredMappingOwnershipChecker(checker DesiredOwnershipChecker) {
	s.mu.Lock()
	s.desiredOwnership = checker
	s.mu.Unlock()
}

// ValidateOwnership implements the handler Share gate.
func (s *Service) ValidateOwnership(ctx context.Context, workspaceID, tagID string) bool {
	s.mu.RLock()
	checker := s.ownership
	s.mu.RUnlock()
	return checker != nil && checker(ctx, workspaceID, tagID)
}

// ValidateDirectMutation rejects the legacy tag/register mutation path. The
// persisted candidate/apply/reconcile seam is the only formal projection API.
func (s *Service) ValidateDirectMutation(_ context.Context, _, _ string) error {
	return NewError(ErrCodeProjectionRequired, "local Modbus mappings must be applied through the workspace candidate projection", false)
}

// UpsertMapping creates or updates tag-to-register mapping.
func (s *Service) UpsertMapping(ctx context.Context, tagID string, register uint16) (*TagMirrorMapping, error) {
	if tagID == "" {
		return nil, fmt.Errorf("tag_id is required")
	}

	t, err := s.tagSvc.GetByID(ctx, tagID)
	if err != nil {
		return nil, fmt.Errorf("tag not found: %w", err)
	}
	if !t.DataType.IsValid() || t.DataType == schema.DataTypeString {
		return nil, &Error{Code: ErrCodeInvalidGeometry, Message: fmt.Sprintf("tag %s has unsupported datatype %q", tagID, t.DataType), Retryable: false}
	}

	span := DataTypeSpan(t.DataType)
	stride := span
	capRegs := s.effectiveCapacityRegisters()

	if int(register)+span > capRegs {
		return nil, &Error{
			Code:      ErrCodeCapacityExceeded,
			Message:   fmt.Sprintf("tag %s register span [%d, %d) exceeds capacity %d", tagID, register, int(register)+span, capRegs),
			Retryable: false,
		}
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	for otherID, existing := range s.mappings {
		if otherID == tagID {
			continue
		}
		otherSpan := existing.SpanRegisters
		if otherSpan <= 0 {
			otherSpan = DataTypeSpan(existing.DataType)
		}
		if int(register) < int(existing.Register)+otherSpan && int(existing.Register) < int(register)+span {
			return nil, &Error{
				Code: ErrCodeRangeCollision,
				Message: fmt.Sprintf("range collision between tag %s [%d, %d) and tag %s [%d, %d)",
					tagID, register, int(register)+span,
					otherID, existing.Register, int(existing.Register)+otherSpan,
				),
				Retryable: false,
			}
		}
	}

	m := TagMirrorMapping{
		TagID:              tagID,
		Register:           register,
		ShareStartRegister: ZeroBasedToHuman(register),
		ZeroBasedRegister:  register,
		SpanRegisters:      span,
		StrideRegisters:    stride,
		CapacityRegisters:  capRegs,
		DataType:           t.DataType,
		UpdatedAt:          time.Now().UTC(),
	}

	s.mappings[tagID] = m
	return &m, nil
}

// RemoveMapping removes an existing mapping.
func (s *Service) RemoveMapping(tagID string) {
	s.mu.Lock()
	delete(s.mappings, tagID)
	s.mu.Unlock()
}
