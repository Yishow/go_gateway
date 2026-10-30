package workspace

import (
	"context"
	"encoding/json"
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

type observedCutover struct {
	Available                                                                                                    bool
	WorkspaceID, WorkspaceRevision, GroupID, GroupRevision, ConnectorID, ConnectorRevision, OperationID, Outcome string
}
type cutoverObservationFixture struct {
	f            counterApplyFixture
	req          CutoverRequest
	target       *WriteGroup
	observations []observedCutover
	scopeCalls   int
	beforeScope  func(int)
	onObserve    func(observedCutover)
}

func newCutoverObservationFixture(t *testing.T, info, warn bool) *cutoverObservationFixture {
	t.Helper()
	f, target, req := counterCutoverFixture(t)
	h := &cutoverObservationFixture{f: f, req: req, target: target}
	h.install(t, info, warn)
	return h
}

// The absent optional constructor seam leaves the real service unchanged for the
// initial behavioral RED: ApplyCutover still commits, but no event is observed.
func (h *cutoverObservationFixture) install(t *testing.T, info, warn bool) {
	t.Helper()
	ctor := reflect.ValueOf(NewWriteGroupService)
	if !ctor.Type().IsVariadic() {
		return
	}
	typ := ctor.Type().In(2).Elem()
	config := reflect.New(typ).Elem()
	config.FieldByName("InfoEnabled").SetBool(info)
	config.FieldByName("WarnEnabled").SetBool(warn)
	config.FieldByName("RecordLimit").SetInt(1)
	scope := config.FieldByName("ScopeRevision")
	scope.Set(reflect.MakeFunc(scope.Type(), func(args []reflect.Value) []reflect.Value {
		h.scopeCalls++
		if h.beforeScope != nil {
			h.beforeScope(h.scopeCalls)
		}
		ctx := args[0].Interface().(context.Context)
		record, err := h.f.service.workspaceSvc.ReadSnapshot(ctx)
		valid := err == nil && record != nil && record.ID == args[1].String() && record.ID == h.f.group.WorkspaceID && args[2].String() == h.f.group.ID && args[3].String() == h.f.group.AppliedRevision && args[4].String() == h.f.group.Destination.ConnectorID && args[5].String() == h.f.group.Destination.ConnectorRevision
		revision := ""
		if valid {
			revision = record.DatabaseSetupRevision
		}
		return []reflect.Value{reflect.ValueOf(revision), reflect.ValueOf(valid)}
	}))
	observe := config.FieldByName("Observe")
	observe.Set(reflect.MakeFunc(observe.Type(), func(args []reflect.Value) []reflect.Value {
		raw, err := json.Marshal(args[1].Interface())
		require.NoError(t, err)
		var v observedCutover
		require.NoError(t, json.Unmarshal(raw, &v))
		h.observations = append(h.observations, v)
		if h.onObserve != nil {
			h.onObserve(v)
		}
		return nil
	}))
	previous := h.f.service
	result := ctor.Call([]reflect.Value{reflect.ValueOf(previous.workspaceSvc), reflect.ValueOf(previous.repo), config})
	next := result[0].Interface().(*WriteGroupService)
	next.now, next.tableInspector, next.managedTableInspector = previous.now, previous.tableInspector, previous.managedTableInspector
	next.backlogGuard, next.writerOwnershipActivation = previous.backlogGuard, previous.writerOwnershipActivation
	h.f.service = next
}
func (h *cutoverObservationFixture) apply(t *testing.T) (*WriteGroupCutover, error) {
	t.Helper()
	return h.f.service.ApplyCutover(t.Context(), h.req)
}
func (h *cutoverObservationFixture) requireUnavailable(t *testing.T) {
	t.Helper()
	require.Equal(t, []observedCutover{{Outcome: "unavailable"}}, h.observations, "unproved operation must not retain partial identity")
}
