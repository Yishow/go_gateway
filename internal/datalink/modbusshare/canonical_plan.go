package modbusshare

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"sort"
)

// CanonicalSharePlan is the server-owned candidate snapshot used by bootstrap,
// status, and reconcile. Its signature covers the workspace/settings revisions
// and every persisted mapping geometry, so a browser cursor cannot become a
// projection authority.
type CanonicalSharePlan struct {
	WorkspaceID       string           `json:"workspace_id"`
	WorkspaceRevision string           `json:"workspace_revision"`
	SettingsRevision  string           `json:"settings_revision"`
	DesiredMappings   []DesiredMapping `json:"desired_mappings"`
	Signature         string           `json:"signature"`
}

// NewCanonicalSharePlan creates a deterministic server-owned plan snapshot.
func NewCanonicalSharePlan(workspaceID, workspaceRevision, settingsRevision string, desired []DesiredMapping) CanonicalSharePlan {
	canonical := append([]DesiredMapping(nil), desired...)
	sort.SliceStable(canonical, func(i, j int) bool {
		if canonical[i].TagID != canonical[j].TagID {
			return canonical[i].TagID < canonical[j].TagID
		}
		return canonical[i].SourceRuleID < canonical[j].SourceRuleID
	})
	plan := CanonicalSharePlan{
		WorkspaceID:       workspaceID,
		WorkspaceRevision: workspaceRevision,
		SettingsRevision:  settingsRevision,
		DesiredMappings:   canonical,
	}
	plan.Signature = plan.signature()
	return plan
}

func (p CanonicalSharePlan) signature() string {
	p.Signature = ""
	payload, err := json.Marshal(p)
	if err != nil {
		return ""
	}
	digest := sha256.Sum256(payload)
	return hex.EncodeToString(digest[:])
}

// Matches reports whether a supplied signature identifies this exact plan.
func (p CanonicalSharePlan) Matches(signature string) bool {
	return signature != "" && signature == p.signature()
}
