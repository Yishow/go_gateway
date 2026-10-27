package counterproof

import "encoding/json"

// SnapshotStorageProjection removes only replay counter bindings from an
// immutable descriptor for finalized delivery. It verifies the same duplicate
// and structural-name rules as replay; callers must still validate the original
// scope, layout and destination. The returned copy is not replay evidence and
// must never replace the persisted descriptor.
func SnapshotStorageProjection(data []byte) ([]byte, error) {
	if !UniqueJSONObject(data, false) || !snapshotFieldNames(data) {
		return nil, &CounterBindingProofError{Reason: ErrCounterBindingProofConflict}
	}
	var fields map[string]json.RawMessage
	if err := json.Unmarshal(data, &fields); err != nil {
		return nil, &CounterBindingProofError{Reason: ErrCounterBindingProofConflict}
	}
	delete(fields, "counter_bindings")
	return json.Marshal(fields)
}
