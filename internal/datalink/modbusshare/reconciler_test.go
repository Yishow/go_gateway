package modbusshare

import (
	"context"
	"fmt"
	"sync"
	"testing"

	"go-gateway/internal/datalink/tag"
)

type mockWorkspaceRevisionStore struct {
	mu               sync.Mutex
	revisions        map[string]string
	dirty            map[string]bool
	desired          map[string][]DesiredMapping
	commitCalls      int
	failCommitAt     int
	updateErr        error
	markDirtyErr     error
	beforeGetDesired func()
}

func newMockWorkspaceRevisionStore() *mockWorkspaceRevisionStore {
	return &mockWorkspaceRevisionStore{
		revisions: make(map[string]string),
		dirty:     make(map[string]bool),
		desired:   make(map[string][]DesiredMapping),
	}
}

func (m *mockWorkspaceRevisionStore) GetDesiredMappings(_ context.Context, wsID string) ([]DesiredMapping, error) {
	if m.beforeGetDesired != nil {
		m.beforeGetDesired()
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	return append([]DesiredMapping(nil), m.desired[wsID]...), nil
}

func (m *mockWorkspaceRevisionStore) CommitDesiredMappings(ctx context.Context, wsID, expectedRev, newRev string, mappings []DesiredMapping) error {
	m.mu.Lock()
	m.commitCalls++
	shouldFail := m.failCommitAt > 0 && m.commitCalls >= m.failCommitAt
	injectedErr := m.updateErr
	m.mu.Unlock()
	if shouldFail && injectedErr != nil {
		return injectedErr
	}
	if err := m.UpdateRevision(ctx, wsID, expectedRev, newRev); err != nil {
		return err
	}
	m.mu.Lock()
	m.desired[wsID] = append([]DesiredMapping(nil), mappings...)
	m.mu.Unlock()
	return nil
}

func (m *mockWorkspaceRevisionStore) GetRevision(ctx context.Context, wsID string) (revision string, isDirty bool, err error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	rev, ok := m.revisions[wsID]
	if !ok {
		return "rev-1", m.dirty[wsID], nil
	}
	return rev, m.dirty[wsID], nil
}

func (m *mockWorkspaceRevisionStore) UpdateRevision(ctx context.Context, wsID, expectedRev, newRev string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.updateErr != nil && (m.failCommitAt == 0 || m.commitCalls >= m.failCommitAt) {
		return m.updateErr
	}
	current, ok := m.revisions[wsID]
	if !ok {
		current = "rev-1"
	}
	if current != expectedRev {
		return &Error{
			Code:      ErrCodeRevisionConflict,
			Message:   fmt.Sprintf("stale workspace revision: expected %s, current %s", expectedRev, current),
			Retryable: true,
		}
	}
	m.revisions[wsID] = newRev
	m.dirty[wsID] = false
	return nil
}

func (m *mockWorkspaceRevisionStore) MarkDirty(ctx context.Context, wsID string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.markDirtyErr != nil {
		return m.markDirtyErr
	}
	m.dirty[wsID] = true
	return nil
}

func setupReconcilerFixture(t *testing.T) (reconciler *Reconciler, service *Service, tagService *tag.Service, revStore *mockWorkspaceRevisionStore) {
	t.Helper()
	tagRepo := tag.NewMemoryRepository()
	tagSvc := tag.NewService(tagRepo)
	svc := NewService(tagSvc, 65536)
	svc.SetHydrationState(HydrationState{State: HydrationStateReady, WorkspaceID: "ws-1", Readiness: true})
	revStore = newMockWorkspaceRevisionStore()
	reconciler = NewReconciler(svc, revStore).WithOwnershipValidator(func(context.Context, DesiredMapping) error {
		return nil
	})
	return reconciler, svc, tagSvc, revStore
}
