package counterproof

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"strings"
)

// counterProofWire fixes the v1 field order. Its digest is omitted when hashing;
// neither JSON indentation nor incoming object order affects proof identity.
type counterProofWire CounterBindingProof

func (p CounterBindingProof) canonicalDigest() (string, error) {
	wire := counterProofWire(p)
	wire.Digest = ""
	data, err := json.Marshal(wire)
	if err != nil {
		return "", p.conflict()
	}
	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:]), nil
}

// Validate verifies both the supported contract and its complete digest.
// It never reads a measurement repository or fills gaps from current settings.
func (p CounterBindingProof) Validate() error {
	if p.Version != CounterBindingProofVersion {
		return &CounterBindingProofError{Reason: ErrCounterBindingProofUnsupported, Version: p.Version}
	}
	if err := p.validateFacts(); err != nil {
		return err
	}
	digest, err := p.canonicalDigest()
	if err != nil || p.Digest != digest {
		return p.conflict()
	}
	return nil
}

// MarshalJSON refuses modified facts instead of silently creating a new proof.
func (p CounterBindingProof) MarshalJSON() ([]byte, error) {
	if err := p.Validate(); err != nil {
		return nil, err
	}
	return json.Marshal(counterProofWire(p))
}

// UnmarshalJSON accepts only a complete supported proof. A failed decode leaves
// the receiver untouched, and exposes only a safe typed reason, not raw facts.
func (p *CounterBindingProof) UnmarshalJSON(data []byte) error {
	var wire counterProofWire
	conflict := CounterBindingProof{}.conflict()
	if !UniqueJSONObject(data, false) {
		return conflict
	}
	var header struct {
		Version json.RawMessage `json:"version"`
	}
	if err := json.Unmarshal(data, &header); err != nil {
		return conflict
	}
	var version int
	if len(header.Version) != 0 {
		if bytes.Equal(bytes.TrimSpace(header.Version), []byte("null")) || json.Unmarshal(header.Version, &version) != nil {
			return conflict
		}
	}
	if version != CounterBindingProofVersion {
		return &CounterBindingProofError{Reason: ErrCounterBindingProofUnsupported, Version: version}
	}
	if !UniqueJSONObject(data, true) {
		return conflict
	}
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&wire); err != nil {
		return conflict
	}
	proof := CounterBindingProof(wire)
	if err := proof.Validate(); err != nil {
		return err
	}
	*p = proof
	return nil
}

// MarshalJSON verifies that lookup identity and sealed member identity agree.
func (proofs CounterBindingProofs) MarshalJSON() ([]byte, error) {
	for key, proof := range proofs {
		if key != proof.Member.Key() {
			return nil, proof.conflict()
		}
	}
	return json.Marshal(map[string]CounterBindingProof(proofs))
}

// UnmarshalJSON validates every map entry before replacing the receiver.
// An absent or null optional map remains compatible with legacy snapshots.
func (proofs *CounterBindingProofs) UnmarshalJSON(data []byte) error {
	if bytes.Equal(bytes.TrimSpace(data), []byte("null")) {
		*proofs = nil
		return nil
	}
	if !UniqueJSONObject(data, false) {
		return CounterBindingProof{}.conflict()
	}
	var decoded map[string]CounterBindingProof
	if err := json.Unmarshal(data, &decoded); err != nil {
		return err
	}
	for key, proof := range decoded {
		if key != proof.Member.Key() {
			return proof.conflict()
		}
	}
	*proofs = decoded
	return nil
}

// UniqueJSONObject rejects duplicate-key last-value-wins interpretations.
// It is shared by proof and referencing journal codecs to avoid divergent parsers.
// V1 contains lower-case object fields and scalars, never arrays or null facts.
// Future versions are checked for duplicate keys without assuming a v1 shape.
func UniqueJSONObject(data []byte, v1 bool) bool {
	if !json.Valid(data) {
		return false
	}
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.UseNumber()
	var object func() bool
	var value func() bool
	value = func() bool {
		token, err := decoder.Token()
		if err != nil || (v1 && token == nil) {
			return false
		}
		switch token {
		case json.Delim('{'):
			return object()
		case json.Delim('['):
			if v1 {
				return false
			}
			for decoder.More() {
				if !value() {
					return false
				}
			}
			end, err := decoder.Token()
			return err == nil && end == json.Delim(']')
		}
		return true
	}
	object = func() bool {
		seen := make(map[string]bool)
		for decoder.More() {
			key, err := decoder.Token()
			if err != nil {
				return false
			}
			name, ok := key.(string)
			if !ok || seen[name] || (v1 && strings.IndexFunc(name, func(r rune) bool {
				// encoding/json folds Unicode aliases such as long s to ASCII.
				// V1 field names are ASCII; opaque string values are not restricted.
				return r != '_' && (r < 'a' || r > 'z')
			}) >= 0) {
				return false
			}
			seen[name] = true
			if !value() {
				return false
			}
		}
		end, err := decoder.Token()
		return err == nil && end == json.Delim('}')
	}
	start, err := decoder.Token()
	return err == nil && start == json.Delim('{') && object()
}
