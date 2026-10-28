package recordingplan

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"
	"strings"
)

const (
	// SuccessorReserveAction is reserved for atomic local successor reservation.
	SuccessorReserveAction = "successor_reserve"
	// SuccessorFirstStartAction is reserved for atomic local first activation.
	SuccessorFirstStartAction = "successor_first_start"
	successorDetailVersion    = "successor-operation/v1"
)

type successorStoredScope struct {
	WorkspaceID        string `json:"workspace_id"`
	RequestID          string `json:"request_id"`
	LineageID          string `json:"lineage_id"`
	PredecessorGroupID string `json:"predecessor_group_id"`
	SuccessorGroupID   string `json:"successor_group_id"`
}

func storedSuccessorScope(s SuccessorOperationScope) successorStoredScope {
	return successorStoredScope{s.WorkspaceID, s.RequestID, s.LineageID, s.PredecessorGroupID, s.SuccessorGroupID}
}
func (s successorStoredScope) public() SuccessorOperationScope {
	return SuccessorOperationScope{s.WorkspaceID, s.RequestID, s.LineageID, s.PredecessorGroupID, s.SuccessorGroupID}
}

type successorOperationDetail struct {
	Version             string               `json:"version"`
	Action              string               `json:"action"`
	Scope               successorStoredScope `json:"scope"`
	ReserveOperationID  string               `json:"reserve_operation_id,omitempty"`
	IntentDigest        string               `json:"intent_digest"`
	Intent              json.RawMessage      `json:"intent"`
	Result              json.RawMessage      `json:"result"`
	SchemaSourceBinding json.RawMessage      `json:"schema_source_binding,omitempty"`
	SourceEvidence      json.RawMessage      `json:"source_evidence,omitempty"`
}

func successorOperationID(domain, workspaceID, requestID string) string {
	encoded, _ := json.Marshal([]string{domain, workspaceID, requestID})
	digest := sha256.Sum256(encoded)
	return strings.TrimSuffix(domain, "/v1") + "-" + hex.EncodeToString(digest[:])
}

func (d successorOperationDetail) id() string {
	if d.Action == SuccessorReserveAction {
		return SuccessorReserveOperationID(d.Scope.WorkspaceID, d.Scope.RequestID)
	}
	return SuccessorFirstStartOperationID(d.Scope.WorkspaceID, d.ReserveOperationID)
}
func (d successorOperationDetail) scopeKey() string {
	encoded, _ := json.Marshal(d.Scope)
	return successorOperationID(d.Action+"/v1", d.Scope.WorkspaceID, string(encoded))
}

func (d *successorOperationDetail) normalize() error {
	if d.Version != successorDetailVersion || !successorIdentity(d.Scope.WorkspaceID) ||
		!successorIdentity(d.Scope.RequestID) || len(d.Scope.RequestID) > 128 || !successorIdentity(d.Scope.LineageID) ||
		!successorIdentity(d.Scope.PredecessorGroupID) || !successorIdentity(d.Scope.SuccessorGroupID) ||
		d.Scope.PredecessorGroupID == d.Scope.SuccessorGroupID || len(d.IntentDigest) != 64 {
		return ErrSchemaOperationResult
	}
	if _, err := hex.DecodeString(d.IntentDigest); err != nil {
		return ErrSchemaOperationResult
	}
	var err error
	if d.Intent, err = successorJSONObject(d.Intent); err != nil {
		return err
	}
	if d.Result, err = successorJSONObject(d.Result); err != nil {
		return err
	}
	switch d.Action {
	case SuccessorReserveAction:
		if d.ReserveOperationID != "" {
			return ErrSchemaOperationMismatch
		}
		if d.SchemaSourceBinding, err = successorJSONObject(d.SchemaSourceBinding); err != nil {
			return err
		}
		if d.SourceEvidence, err = successorJSONObject(d.SourceEvidence); err != nil {
			return err
		}
	case SuccessorFirstStartAction:
		if d.ReserveOperationID != SuccessorReserveOperationID(d.Scope.WorkspaceID, d.Scope.RequestID) || len(d.SchemaSourceBinding) != 0 || len(d.SourceEvidence) != 0 {
			return ErrSchemaOperationMismatch
		}
	default:
		return ErrSchemaOperationMismatch
	}
	return nil
}

func successorIdentity(v string) bool {
	return v != "" && v == strings.TrimSpace(v) && !strings.ContainsRune(v, 0)
}

// Canonicalize complete JSON, retaining numeric precision and rejecting duplicate
// keys rather than allowing an ambiguous intent to collapse into another one.
func successorJSONObject(raw json.RawMessage) (json.RawMessage, error) {
	dec := json.NewDecoder(bytes.NewReader(raw))
	dec.UseNumber()
	value, err := successorJSONValue(dec)
	if err != nil {
		return nil, ErrSchemaOperationResult
	}
	if _, ok := value.(map[string]any); !ok {
		return nil, ErrSchemaOperationResult
	}
	if _, err = dec.Token(); err != io.EOF {
		return nil, ErrSchemaOperationResult
	}
	result, err := json.Marshal(value)
	if err != nil {
		return nil, ErrSchemaOperationResult
	}
	return result, nil
}

func successorJSONValue(dec *json.Decoder) (any, error) {
	token, err := dec.Token()
	if err != nil {
		return nil, err
	}
	delimiter, ok := token.(json.Delim)
	if !ok {
		return token, nil
	}
	switch delimiter {
	case '{':
		object := map[string]any{}
		for dec.More() {
			keyToken, err := dec.Token()
			if err != nil {
				return nil, err
			}
			key, ok := keyToken.(string)
			if !ok {
				return nil, ErrSchemaOperationResult
			}
			if _, exists := object[key]; exists {
				return nil, ErrSchemaOperationResult
			}
			object[key], err = successorJSONValue(dec)
			if err != nil {
				return nil, err
			}
		}
		end, err := dec.Token()
		if err != nil || end != json.Delim('}') {
			return nil, ErrSchemaOperationResult
		}
		return object, nil
	case '[':
		array := []any{}
		for dec.More() {
			value, err := successorJSONValue(dec)
			if err != nil {
				return nil, err
			}
			array = append(array, value)
		}
		end, err := dec.Token()
		if err != nil || end != json.Delim(']') {
			return nil, ErrSchemaOperationResult
		}
		return array, nil
	default:
		return nil, ErrSchemaOperationResult
	}
}
