package workspace

import (
	"encoding/json"
	"fmt"
	"maps"
)

// FrozenCopy detaches the existing runtime input before provider callbacks can
// mutate caller-owned records or slices. It preserves only this input's opaque
// source evidence; copying hand-built/JSON input cannot manufacture authority.
// This is a memory snapshot, not a fresh database read or completeness check.
func (p *RuntimeProjection) FrozenCopy() (*RuntimeProjection, error) {
	if p == nil {
		return nil, fmt.Errorf("runtime projection unavailable")
	}
	raw, err := json.Marshal(p)
	if err != nil {
		return nil, fmt.Errorf("freeze runtime projection: %w", err)
	}
	var frozen RuntimeProjection
	if err := json.Unmarshal(raw, &frozen); err != nil {
		return nil, fmt.Errorf("freeze runtime projection: %w", err)
	}
	frozen.acquisitionGraphs = maps.Clone(p.acquisitionGraphs)
	return &frozen, nil
}
