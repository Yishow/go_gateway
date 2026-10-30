package workspace

import (
	"encoding/json"
	"fmt"
)

func isMixedWriteGroup(group *WriteGroup) bool {
	return group != nil && group.RecordingProfile != nil && group.RecordingProfile.Kind == WriteGroupProfileMixed
}

// marshalWriteGroupRowPolicy is shared by public and persisted JSON. A mixed
// record has no group window; legacy byte representation remains unchanged.
func marshalWriteGroupRowPolicy(group *WriteGroup) ([]byte, error) {
	if !isMixedWriteGroup(group) {
		return json.Marshal(group.RowPolicy)
	}
	type policyAlias WriteGroupRowPolicy
	return json.Marshal(struct {
		*policyAlias
		IntervalSeconds *int `json:"interval_seconds,omitempty"`
	}{policyAlias: (*policyAlias)(&group.RowPolicy)})
}

func (group WriteGroup) MarshalJSON() ([]byte, error) {
	type groupAlias WriteGroup
	if !isMixedWriteGroup(&group) {
		return json.Marshal(groupAlias(group))
	}
	policy, err := marshalWriteGroupRowPolicy(&group)
	if err != nil {
		return nil, err
	}
	return json.Marshal(struct {
		*groupAlias
		RowPolicy json.RawMessage `json:"row_policy"`
	}{groupAlias: (*groupAlias)(&group), RowPolicy: policy})
}

// UnmarshalJSON distinguishes an absent mixed group interval from an explicit
// legacy field (including zero/null), while retaining legacy decoding exactly.
func (group *WriteGroup) UnmarshalJSON(data []byte) error {
	type groupAlias WriteGroup
	var decoded groupAlias
	if err := json.Unmarshal(data, &decoded); err != nil {
		return err
	}
	if isMixedWriteGroup((*WriteGroup)(&decoded)) {
		var raw struct {
			RowPolicy map[string]json.RawMessage `json:"row_policy"`
		}
		if err := json.Unmarshal(data, &raw); err != nil {
			return err
		}
		if _, present := raw.RowPolicy["interval_seconds"]; present {
			return fmt.Errorf("%w: mixed row policy cannot contain interval_seconds", ErrWriteGroupValidation)
		}
	}
	*group = WriteGroup(decoded)
	return nil
}
