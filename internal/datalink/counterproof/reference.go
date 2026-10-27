package counterproof

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"errors"
	"strings"
	"unicode/utf8"
)

// ErrCounterBindingProofMissing marks accepted counter meaning without evidence.
var ErrCounterBindingProofMissing = errors.New("counter-replay-proof-missing")

// CounterBindingReference identifies a sealed proof in an immutable snapshot.
// It contains only comparable values so copying a sample cannot share mutation.
type CounterBindingReference struct {
	Version int                          `json:"version"`
	Member  CounterBindingMemberIdentity `json:"member"`
	Digest  string                       `json:"digest"`
}

// IsZero distinguishes the original journal format from a proof reference.
func (r CounterBindingReference) IsZero() bool { return r == (CounterBindingReference{}) }

// Validate rejects absent, unsupported or incomplete references before use.
func (r CounterBindingReference) Validate() error {
	if r.IsZero() {
		return &CounterBindingProofError{Reason: ErrCounterBindingProofMissing}
	}
	if r.Version != CounterBindingProofVersion {
		return &CounterBindingProofError{Reason: ErrCounterBindingProofUnsupported, Version: r.Version}
	}
	for _, value := range []string{r.Member.WorkspaceID, r.Member.GroupID, r.Member.AppliedRevision, r.Member.DeviceID, r.Member.PointID, r.Member.TagID} {
		if strings.TrimSpace(value) == "" || !utf8.ValidString(value) {
			return CounterBindingProof{}.conflict()
		}
	}
	if decoded, err := hex.DecodeString(r.Digest); err != nil || len(decoded) != 32 || strings.ToLower(r.Digest) != r.Digest {
		return CounterBindingProof{}.conflict()
	}
	return nil
}

// UnmarshalJSON refuses ambiguous fields rather than interpreting an alias.
func (r *CounterBindingReference) UnmarshalJSON(data []byte) error {
	if !UniqueJSONObject(data, false) {
		return CounterBindingProof{}.conflict()
	}
	var header struct {
		Version int `json:"version"`
	}
	if json.Unmarshal(data, &header) != nil {
		return CounterBindingProof{}.conflict()
	}
	if header.Version != CounterBindingProofVersion {
		return &CounterBindingProofError{Reason: ErrCounterBindingProofUnsupported, Version: header.Version}
	}
	if !UniqueJSONObject(data, true) {
		return CounterBindingProof{}.conflict()
	}
	type wireReference CounterBindingReference
	var decoded wireReference
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	if decoder.Decode(&decoded) != nil {
		return CounterBindingProof{}.conflict()
	}
	ref := CounterBindingReference(decoded)
	if err := ref.Validate(); err != nil {
		return err
	}
	*r = ref
	return nil
}

// Reference verifies the proof before exposing its journal identity.
func (p CounterBindingProof) Reference() (CounterBindingReference, error) {
	if err := p.Validate(); err != nil {
		return CounterBindingReference{}, err
	}
	return CounterBindingReference{Version: p.Version, Member: p.Member, Digest: p.Digest}, nil
}

// SampleMemberKey is the existing snapshot selector's device/point/tag key.
func (m CounterBindingMemberIdentity) SampleMemberKey() string {
	key, _ := json.Marshal([3]string{m.DeviceID, m.PointID, m.TagID})
	return string(key)
}
