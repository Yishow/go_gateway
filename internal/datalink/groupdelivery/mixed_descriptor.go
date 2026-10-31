package groupdelivery

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"go-gateway/internal/datalink/counterproof"
	"go-gateway/internal/datalink/snapshot"
)

// mixedDescriptor reads only scheduling projections from the immutable runtime
// descriptor. Counter meaning continues to use ResolveSnapshotEvidence.
type mixedDescriptor struct {
	WorkspaceID     string `json:"workspace_id"`
	GroupID         string `json:"group_id"`
	AppliedRevision string `json:"applied_revision"`
	Group           *struct {
		ID              string `json:"id"`
		WorkspaceID     string `json:"workspace_id"`
		Revision        string `json:"revision"`
		AppliedRevision string `json:"applied_revision"`
		LayoutVersion   string `json:"layout_version"`
		Profile         *struct {
			Kind  string `json:"kind"`
			Shape string `json:"result_shape"`
		} `json:"recording_profile"`
		Members []struct {
			DeviceID    string `json:"device_id"`
			PointID     string `json:"point_id"`
			TagID       string `json:"tag_id"`
			Calculation string `json:"calculation"`
			Cadence     *struct {
				Interval  int64  `json:"interval"`
				Unit      string `json:"unit"`
				Alignment string `json:"alignment"`
			} `json:"cadence"`
		} `json:"members"`
	} `json:"group"`
}

func readMixedDescriptor(ctx context.Context, q ownershipQuery, key GroupKey) (*mixedDescriptor, string, error) {
	var payload string
	err := q.QueryRowContext(ctx, `SELECT payload FROM wg_runtime_versions WHERE workspace_id=? AND group_id=? AND group_revision=?`, key.WorkspaceID, key.GroupID, key.GroupRevision).Scan(&payload)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, "", nil
	}
	if err != nil {
		return nil, "", err
	}
	var d mixedDescriptor
	if json.Unmarshal([]byte(payload), &d) != nil {
		return nil, "", ErrInvalidClosure
	}
	if d.Group == nil || d.Group.Profile == nil || d.Group.Profile.Kind != "mixed" {
		return nil, "", nil
	}
	if !counterproof.UniqueJSONObject([]byte(payload), false) || d.WorkspaceID != key.WorkspaceID || d.GroupID != key.GroupID || d.AppliedRevision != key.GroupRevision || d.Group.WorkspaceID != key.WorkspaceID || d.Group.ID != key.GroupID || d.Group.Revision != key.GroupRevision || d.Group.AppliedRevision != key.GroupRevision || d.Group.Profile.Shape != "member-calculation-v1" {
		return nil, "", ErrInvalidClosure
	}
	if _, err := counterproof.ResolveSnapshotEvidence([]byte(payload), counterproof.SnapshotScope{WorkspaceID: key.WorkspaceID, GroupID: key.GroupID, AppliedRevision: key.GroupRevision}); err != nil {
		return nil, "", err
	}
	return &d, handoffHash(payload), nil
}
func (d *mixedDescriptor) intervals() (map[string]time.Duration, error) {
	result := map[string]time.Duration{}
	for _, m := range d.Group.Members {
		if m.DeviceID == "" || m.PointID == "" || m.TagID == "" || m.Cadence == nil || m.Cadence.Alignment != "utc_epoch" {
			return nil, ErrInvalidClosure
		}
		key := (counterproof.CounterBindingMemberIdentity{DeviceID: m.DeviceID, PointID: m.PointID, TagID: m.TagID}).SampleMemberKey()
		interval, err := snapshot.CheckedCadenceDuration(m.Cadence.Interval, m.Cadence.Unit)
		if err != nil || result[key] != 0 {
			return nil, ErrInvalidClosure
		}
		result[key] = interval
	}
	if len(result) == 0 {
		return nil, ErrInvalidClosure
	}
	return result, nil
}
func validateMixedCheckpoint(key GroupKey, cp snapshot.MixedCheckpoint, intervals map[string]time.Duration) error {
	if cp.Version != 1 || cp.Progress.WorkspaceID != key.WorkspaceID || cp.Progress.GroupID != key.GroupID || cp.Progress.GroupRevision != key.GroupRevision || cp.StartsAt.IsZero() || len(cp.Members) != len(intervals) {
		return ErrInvalidClosure
	}
	// Reuse the assembler's scheduling projection; this instance never receives
	// input or computes results. It introduces no counter or reducer authority.
	return snapshot.ValidateMixedCheckpointCadences(cp, intervals)
}
func mixedInvalid(reason string) error { return fmt.Errorf("%w: mixed %s", ErrInvalidClosure, reason) }
