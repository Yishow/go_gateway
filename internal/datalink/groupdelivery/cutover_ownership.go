package groupdelivery

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"

	"go-gateway/internal/datalink/snapshot"
)

// OwnershipSource names a source without conflating independent profile kinds.
type OwnershipSource struct {
	WorkspaceID string `json:"workspace_id"`
	DeviceID    string `json:"device_id"`
	PointID     string `json:"point_id"`
	TagID       string `json:"tag_id"`
}

// CutoverOwnership separates affected revisions from their eligible owners.
// An empty result proves no relevant cutover; it is not a missing-table fallback.
type CutoverOwnership struct {
	Affected     map[GroupKey]bool
	Owners       map[GroupKey]bool
	predecessors map[GroupKey]map[GroupKey]bool
	bounds       map[GroupKey]ownershipBound
	roots        map[string]bool
}

type ownershipQuery interface {
	QueryContext(context.Context, string, ...any) (*sql.Rows, error)
	QueryRowContext(context.Context, string, ...any) *sql.Row
}
type ownershipRef struct {
	GroupID  string `json:"group_id"`
	Revision string `json:"revision"`
}
type ownershipDetail struct {
	Fence   string            `json:"ownership_fence_revision"`
	Ended   []ownershipRef    `json:"ended"`
	Started []ownershipRef    `json:"started"`
	Moved   []OwnershipSource `json:"moved"`
}
type ownershipDestination struct {
	ConnectorID       string `json:"connector_id"`
	ConnectorRevision string `json:"connector_revision"`
	Database          string `json:"database"`
	Schema            string `json:"table_schema"`
	Table             string `json:"table_name"`
}
type ownershipVersion struct {
	ID              string               `json:"id"`
	WorkspaceID     string               `json:"workspace_id"`
	Revision        string               `json:"revision"`
	AppliedRevision string               `json:"applied_revision"`
	Members         []OwnershipSource    `json:"members"`
	Destination     ownershipDestination `json:"destination"`
	Profile         struct {
		Kind string `json:"kind"`
	} `json:"recording_profile"`
}
type ownershipScope struct {
	Key                           GroupKey
	Direction, Device, Point, Tag string
}
type ownershipBound struct {
	Start, End    time.Time
	matchesSource bool
}

// ResolveCutoverOwnership is a consistent routing hint. AppendSample repeats
// the same validation in its write transaction; this read never authorizes ACK.
func (s *Store) ResolveCutoverOwnership(ctx context.Context, source OwnershipSource, observedAt time.Time) (CutoverOwnership, error) {
	if source.WorkspaceID == "" || source.DeviceID == "" || source.PointID == "" || source.TagID == "" || observedAt.IsZero() {
		return CutoverOwnership{}, ErrOwnershipUnavailable
	}
	tx, err := s.db.BeginTx(ctx, &sql.TxOptions{ReadOnly: true})
	if err != nil {
		return CutoverOwnership{}, ErrOwnershipUnavailable
	}
	defer tx.Rollback()
	result, err := resolveCutoverOwnership(ctx, tx, source, observedAt, GroupKey{})
	if err != nil {
		return CutoverOwnership{}, err
	}
	if err = tx.Commit(); err != nil {
		return CutoverOwnership{}, ErrOwnershipUnavailable
	}
	return result, nil
}

func validateCutoverAcceptance(ctx context.Context, tx *sql.Tx, key GroupKey, sample snapshot.Sample) error {
	var parts []string
	source := OwnershipSource{WorkspaceID: key.WorkspaceID}
	if json.Unmarshal([]byte(sample.MemberKey), &parts) == nil && len(parts) == 3 {
		source.DeviceID, source.PointID, source.TagID = parts[0], parts[1], parts[2]
	}
	ownership, err := resolveCutoverOwnership(ctx, tx, source, sample.ObservedAt, key)
	if err != nil {
		return err
	}
	if ownership.Affected[key] && !ownership.Owners[key] {
		return ErrOwnershipConflict
	}
	return rejectAcceptedPredecessor(ctx, tx, ownership, key, source, sample)
}

// Source discovery uses indexed fence candidates PLUS exact canonical JSON
// fallback. The latter is essential when every fence row of a later cutover is
// missing: an earlier valid fence must not disguise that absence. Only matched
// roots are materialized; unrelated snapshot corruption is not interpreted.
const ownershipRootsSQL = `SELECT c.id,c.boundary,c.detail
 FROM write_group_cutovers c WHERE c.workspace_id=? AND (
 EXISTS(SELECT 1 FROM write_group_ownership_fences f WHERE f.fence_revision=c.id AND f.workspace_id=c.workspace_id AND
  ((f.device_id=? AND f.point_id=? AND f.tag_id=? AND f.device_id<>'') OR (f.group_id=? AND f.group_revision=? AND f.group_id<>'')))
 OR EXISTS(SELECT 1 FROM json_each(CASE WHEN json_valid(c.detail) THEN c.detail ELSE '{}' END,'$.moved') m
  WHERE json_extract(CASE WHEN json_valid(m.value) THEN m.value ELSE '{}' END,'$.device_id')<>'' AND json_extract(CASE WHEN json_valid(m.value) THEN m.value ELSE '{}' END,'$.device_id')=? AND json_extract(CASE WHEN json_valid(m.value) THEN m.value ELSE '{}' END,'$.point_id')=? AND json_extract(CASE WHEN json_valid(m.value) THEN m.value ELSE '{}' END,'$.tag_id')=?)
 OR (c.source_group_id=? AND c.source_group_id<>'')
 OR EXISTS(SELECT 1 FROM json_each(CASE WHEN json_valid(c.detail) THEN c.detail ELSE '{}' END,'$.ended') r
  WHERE json_extract(CASE WHEN json_valid(r.value) THEN r.value ELSE '{}' END,'$.group_id')<>'' AND json_extract(CASE WHEN json_valid(r.value) THEN r.value ELSE '{}' END,'$.group_id')=? AND json_extract(CASE WHEN json_valid(r.value) THEN r.value ELSE '{}' END,'$.revision')=?)
 OR EXISTS(SELECT 1 FROM json_each(CASE WHEN json_valid(c.detail) THEN c.detail ELSE '{}' END,'$.started') r
  WHERE json_extract(CASE WHEN json_valid(r.value) THEN r.value ELSE '{}' END,'$.group_id')<>'' AND json_extract(CASE WHEN json_valid(r.value) THEN r.value ELSE '{}' END,'$.group_id')=? AND json_extract(CASE WHEN json_valid(r.value) THEN r.value ELSE '{}' END,'$.revision')=?)
 OR EXISTS(SELECT 1 FROM write_group_versions v,
  json_each(CASE WHEN json_valid(v.payload) THEN v.payload ELSE '{}' END,'$.members') m
  WHERE v.group_id=c.source_group_id AND json_extract(CASE WHEN json_valid(m.value) THEN m.value ELSE '{}' END,'$.device_id')<>'' AND json_extract(CASE WHEN json_valid(m.value) THEN m.value ELSE '{}' END,'$.device_id')=?
   AND json_extract(CASE WHEN json_valid(m.value) THEN m.value ELSE '{}' END,'$.point_id')=? AND json_extract(CASE WHEN json_valid(m.value) THEN m.value ELSE '{}' END,'$.tag_id')=? AND
   (NOT json_valid(c.detail) OR json_type(CASE WHEN json_valid(c.detail) THEN c.detail ELSE '{}' END,'$.ended') IS NOT 'array' OR json_type(CASE WHEN json_valid(c.detail) THEN c.detail ELSE '{}' END,'$.started') IS NOT 'array' OR EXISTS(SELECT 1 FROM json_each(CASE WHEN json_valid(c.detail) THEN c.detail ELSE '{}' END,'$.ended') r
     WHERE json_extract(CASE WHEN json_valid(r.value) THEN r.value ELSE '{}' END,'$.group_id')=v.group_id AND json_extract(CASE WHEN json_valid(r.value) THEN r.value ELSE '{}' END,'$.revision')=v.group_revision)
    OR EXISTS(SELECT 1 FROM json_each(CASE WHEN json_valid(c.detail) THEN c.detail ELSE '{}' END,'$.started') r
     WHERE json_extract(CASE WHEN json_valid(r.value) THEN r.value ELSE '{}' END,'$.group_id')=v.group_id AND json_extract(CASE WHEN json_valid(r.value) THEN r.value ELSE '{}' END,'$.revision')=v.group_revision)))
 )`

func resolveCutoverOwnership(ctx context.Context, q ownershipQuery, source OwnershipSource, observedAt time.Time, key GroupKey) (CutoverOwnership, error) {
	result := CutoverOwnership{Affected: map[GroupKey]bool{}, Owners: map[GroupKey]bool{}, predecessors: map[GroupKey]map[GroupKey]bool{}}
	// Even an empty candidate set must prove both authority tables are readable.
	var count int
	if err := q.QueryRowContext(ctx, `SELECT COUNT(*) FROM write_group_revision_ends WHERE group_id=? AND group_revision=?`, key.GroupID, key.GroupRevision).Scan(&count); err != nil {
		return result, ErrOwnershipUnavailable
	}
	discovery := source
	if key.Valid() {
		discovery.DeviceID, discovery.PointID, discovery.TagID = "", "", ""
	}
	rows, err := q.QueryContext(ctx, ownershipRootsSQL, source.WorkspaceID,
		discovery.DeviceID, discovery.PointID, discovery.TagID, key.GroupID, key.GroupRevision,
		discovery.DeviceID, discovery.PointID, discovery.TagID, key.GroupID, key.GroupID, key.GroupRevision, key.GroupID, key.GroupRevision,
		discovery.DeviceID, discovery.PointID, discovery.TagID)
	if err != nil {
		return result, ErrOwnershipUnavailable
	}
	defer rows.Close()
	bounds := map[GroupKey]ownershipBound{}
	roots := map[string]bool{}
	for rows.Next() {
		var id, boundaryText, raw string
		if err := rows.Scan(&id, &boundaryText, &raw); err != nil {
			return result, ErrOwnershipUnavailable
		}
		var detail ownershipDetail
		boundary, err := time.Parse(time.RFC3339Nano, boundaryText)
		if err != nil || json.Unmarshal([]byte(raw), &detail) != nil || id == "" || len(detail.Ended) == 0 || len(detail.Started) == 0 {
			return result, ErrOwnershipUnavailable
		}
		if key.Valid() && !detailReferencesKey(detail, key) {
			var n int
			if q.QueryRowContext(ctx, `SELECT COUNT(*) FROM write_group_ownership_fences WHERE fence_revision=? AND group_id=? AND group_revision=?`, id, key.GroupID, key.GroupRevision).Scan(&n) != nil || n != 0 {
				return result, ErrOwnershipUnavailable
			}
			continue
		}
		if detail.Fence != id {
			return result, ErrOwnershipUnavailable
		}
		if err := validateOwnershipRoot(ctx, q, source, id, boundary, detail, key, bounds, result.Affected, result.predecessors); err != nil {
			return result, err
		}
		roots[id] = true
	}
	if rows.Err() != nil {
		return result, ErrOwnershipUnavailable
	}
	if err := rows.Close(); err != nil {
		return result, ErrOwnershipUnavailable
	}
	if err := validateOwnershipEnds(ctx, q, source, key, roots); err != nil {
		return result, err
	}
	result.bounds, result.roots = bounds, roots
	for owner, bound := range bounds {
		if bound.matchesSource && !observedAt.Before(bound.Start) && (bound.End.IsZero() || observedAt.Before(bound.End)) {
			result.Owners[owner] = true
		}
	}
	return result, nil
}

func readOwnershipVersion(ctx context.Context, q ownershipQuery, workspace string, ref ownershipRef) (ownershipVersion, time.Time, error) {
	var v ownershipVersion
	var effective time.Time
	var payload string
	err := q.QueryRowContext(ctx, `SELECT effective_at,payload FROM write_group_versions WHERE group_id=? AND group_revision=?`, ref.GroupID, ref.Revision).Scan(&effective, &payload)
	if err != nil || json.Unmarshal([]byte(payload), &v) != nil || v.ID != ref.GroupID || v.Revision != ref.Revision || v.AppliedRevision != ref.Revision || v.WorkspaceID != workspace || len(v.Members) == 0 {
		return v, effective, ErrOwnershipUnavailable
	}
	return v, effective.UTC(), nil
}

func validateOwnershipRoot(ctx context.Context, q ownershipQuery, source OwnershipSource, id string, boundary time.Time, detail ownershipDetail, requested GroupKey, bounds map[GroupKey]ownershipBound, affected map[GroupKey]bool, predecessors map[GroupKey]map[GroupKey]bool) error {
	expected := map[ownershipScope]ownershipVersion{}
	members := map[string]map[[3]string]bool{"ended": {}, "started": {}}
	matched := map[string][]GroupKey{}
	for _, side := range []struct {
		direction string
		refs      []ownershipRef
	}{{"ended", detail.Ended}, {"started", detail.Started}} {
		seen := map[ownershipRef]bool{}
		for _, ref := range side.refs {
			if ref.GroupID == "" || ref.Revision == "" || seen[ref] {
				return ErrOwnershipUnavailable
			}
			seen[ref] = true
			version, effective, err := readOwnershipVersion(ctx, q, source.WorkspaceID, ref)
			if err != nil {
				return err
			}
			if side.direction == "started" && !effective.Equal(boundary) {
				return ErrOwnershipUnavailable
			}
			if side.direction == "ended" {
				var endText, cutoverID, workspace string
				err := q.QueryRowContext(ctx, `SELECT workspace_id,ends_at,cutover_id FROM write_group_revision_ends WHERE group_id=? AND group_revision=?`, ref.GroupID, ref.Revision).Scan(&workspace, &endText, &cutoverID)
				end, parseErr := time.Parse(time.RFC3339Nano, endText)
				if err != nil || parseErr != nil || workspace != source.WorkspaceID || cutoverID != id || !end.Equal(boundary) {
					return ErrOwnershipUnavailable
				}
			}
			groupKey := GroupKey{WorkspaceID: source.WorkspaceID, GroupID: ref.GroupID, GroupRevision: ref.Revision}
			for _, member := range version.Members {
				scope := ownershipScope{Key: groupKey, Direction: side.direction, Device: member.DeviceID, Point: member.PointID, Tag: member.TagID}
				if scope.Device == "" || scope.Point == "" || scope.Tag == "" {
					return ErrOwnershipUnavailable
				}
				if _, duplicate := expected[scope]; duplicate {
					return ErrOwnershipUnavailable
				}
				expected[scope] = version
				members[side.direction][[3]string{scope.Device, scope.Point, scope.Tag}] = true
				if groupKey == requested {
					affected[groupKey] = true
				}
				matchesSource := member.DeviceID == source.DeviceID && member.PointID == source.PointID && member.TagID == source.TagID
				if groupKey == requested || matchesSource {
					if matchesSource {
						matched[side.direction] = append(matched[side.direction], groupKey)
					}
					affected[groupKey] = true
					b := bounds[groupKey]
					b.matchesSource = b.matchesSource || matchesSource
					if !b.Start.IsZero() && !b.Start.Equal(effective) {
						return ErrOwnershipUnavailable
					}
					b.Start = effective
					if side.direction == "ended" {
						if !b.End.IsZero() && !b.End.Equal(boundary) {
							return ErrOwnershipUnavailable
						}
						b.End = boundary
					}
					bounds[groupKey] = b
				}
			}
		}
	}
	if len(detail.Moved) == 0 || len(members["ended"]) != len(members["started"]) {
		return ErrOwnershipUnavailable
	}
	for source := range members["ended"] {
		if !members["started"][source] {
			return ErrOwnershipUnavailable
		}
	}
	moved := map[[3]string]bool{}
	for _, source := range detail.Moved {
		ref := [3]string{source.DeviceID, source.PointID, source.TagID}
		if !members["ended"][ref] || !members["started"][ref] || moved[ref] {
			return ErrOwnershipUnavailable
		}
		moved[ref] = true
	}
	rows, err := q.QueryContext(ctx, `SELECT workspace_id,direction,group_id,group_revision,device_id,point_id,tag_id,profile_kind,connector_id,connector_revision,database_name,table_schema,table_name FROM write_group_ownership_fences WHERE fence_revision=?`, id)
	if err != nil {
		return ErrOwnershipUnavailable
	}
	defer rows.Close()
	for rows.Next() {
		var scope ownershipScope
		var kind string
		var destination ownershipDestination
		if err := rows.Scan(&scope.Key.WorkspaceID, &scope.Direction, &scope.Key.GroupID, &scope.Key.GroupRevision, &scope.Device, &scope.Point, &scope.Tag, &kind, &destination.ConnectorID, &destination.ConnectorRevision, &destination.Database, &destination.Schema, &destination.Table); err != nil {
			return ErrOwnershipUnavailable
		}
		version, ok := expected[scope]
		if !ok || version.Profile.Kind != kind || version.Destination != destination {
			return ErrOwnershipUnavailable
		}
		delete(expected, scope)
	}
	if rows.Err() != nil || len(expected) != 0 {
		return ErrOwnershipUnavailable
	}
	for _, successor := range matched["started"] {
		if predecessors[successor] == nil {
			predecessors[successor] = map[GroupKey]bool{}
		}
		for _, prior := range matched["ended"] {
			predecessors[successor][prior] = true
		}
	}
	return nil
}

func validateOwnershipEnds(ctx context.Context, q ownershipQuery, source OwnershipSource, key GroupKey, roots map[string]bool) error {
	rows, err := q.QueryContext(ctx, `SELECT e.workspace_id,e.group_id,e.group_revision,e.ends_at,e.cutover_id,c.detail
 FROM write_group_revision_ends e
 LEFT JOIN write_group_cutovers c ON c.id=e.cutover_id
 LEFT JOIN write_group_versions v ON v.group_id=e.group_id AND v.group_revision=e.group_revision
 WHERE (e.group_id=? AND e.group_revision=?) OR (e.workspace_id=? AND EXISTS(
 SELECT 1 FROM json_each(CASE WHEN json_valid(v.payload) THEN v.payload ELSE '{}' END,'$.members') m
 WHERE json_extract(CASE WHEN json_valid(m.value) THEN m.value ELSE '{}' END,'$.device_id')<>'' AND json_extract(CASE WHEN json_valid(m.value) THEN m.value ELSE '{}' END,'$.device_id')=? AND json_extract(CASE WHEN json_valid(m.value) THEN m.value ELSE '{}' END,'$.point_id')=? AND json_extract(CASE WHEN json_valid(m.value) THEN m.value ELSE '{}' END,'$.tag_id')=?))`, key.GroupID, key.GroupRevision, source.WorkspaceID, source.DeviceID, source.PointID, source.TagID)
	if err != nil {
		return ErrOwnershipUnavailable
	}
	defer rows.Close()
	for rows.Next() {
		var ended GroupKey
		var endText, id string
		var raw sql.NullString
		if rows.Scan(&ended.WorkspaceID, &ended.GroupID, &ended.GroupRevision, &endText, &id, &raw) != nil {
			return ErrOwnershipUnavailable
		}
		if ended == key || !key.Valid() {
			if !roots[id] {
				return ErrOwnershipUnavailable
			}
			continue
		}
		if ended.GroupID == key.GroupID && ended.GroupRevision == key.GroupRevision {
			return ErrOwnershipUnavailable
		}
		var detail ownershipDetail
		if raw.Valid && json.Unmarshal([]byte(raw.String), &detail) == nil && len(detail.Ended) > 0 && len(detail.Started) > 0 {
			if detailReferencesKey(detail, key) && !roots[id] {
				return ErrOwnershipUnavailable
			}
			continue // the intact immutable intent proves this distinct key is unrelated
		}
		// Without the root, matching start/end instants are ambiguity, not ownership.
		// A bare key cannot prove independence in a relevant orphan-source scope.
		version, effective, err := readOwnershipVersion(ctx, q, source.WorkspaceID, ownershipRef{GroupID: key.GroupID, Revision: key.GroupRevision})
		end, parseErr := time.Parse(time.RFC3339Nano, endText)
		if err != nil || parseErr != nil || effective.Equal(end) {
			return ErrOwnershipUnavailable
		}
		memberFound := false
		for _, member := range version.Members {
			if member.DeviceID == source.DeviceID && member.PointID == source.PointID && member.TagID == source.TagID {
				memberFound = true
			}
		}
		if !memberFound {
			return ErrOwnershipUnavailable
		}
	}
	if rows.Err() != nil {
		return ErrOwnershipUnavailable
	}
	var orphan bool
	err = q.QueryRowContext(ctx, `SELECT EXISTS(SELECT 1 FROM write_group_ownership_fences f LEFT JOIN write_group_cutovers c ON c.id=f.fence_revision
 WHERE c.id IS NULL AND f.workspace_id=? AND ((f.group_id=? AND f.group_revision=?) OR (f.device_id=? AND f.point_id=? AND f.tag_id=?)))`, source.WorkspaceID, key.GroupID, key.GroupRevision, source.DeviceID, source.PointID, source.TagID).Scan(&orphan)
	if err != nil || orphan {
		return ErrOwnershipUnavailable
	}
	return nil
}

// Already accepted observations need the later transfer protocol. Reconstructing
// a target proof is not permission to append the same observation to a new key.
// Independent sibling profiles are not predecessors, and a new observation at a
// different time may legitimately reuse an ID under a different revision.
func rejectAcceptedPredecessor(ctx context.Context, q ownershipQuery, ownership CutoverOwnership, key GroupKey, source OwnershipSource, sample snapshot.Sample) error {
	queue := []GroupKey{key}
	seen := map[GroupKey]bool{key: true}
	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]
		if current != key {
			var err error
			ownership, err = resolveCutoverOwnership(ctx, q, source, sample.ObservedAt, current)
			if err != nil {
				return err
			}
		}
		for prior := range ownership.predecessors[current] {
			if seen[prior] {
				continue
			}
			seen[prior] = true
			queue = append(queue, prior)
			rows, err := q.QueryContext(ctx, `SELECT member_key,observed_at FROM wg_delivery_samples WHERE workspace_id=? AND group_id=? AND group_revision=? AND COALESCE(NULLIF(origin_sample_id,''),sample_id)=?`, prior.WorkspaceID, prior.GroupID, prior.GroupRevision, sample.SampleID)
			if err != nil {
				return ErrOwnershipUnavailable
			}
			conflict := false
			for rows.Next() {
				var member, observed string
				if rows.Scan(&member, &observed) != nil {
					rows.Close()
					return ErrOwnershipUnavailable
				}
				at, err := time.Parse(time.RFC3339Nano, observed)
				if err != nil {
					rows.Close()
					return ErrOwnershipUnavailable
				}
				if member == sample.MemberKey && at.Equal(sample.ObservedAt) {
					conflict = true
				}
			}
			if rows.Err() != nil {
				rows.Close()
				return ErrOwnershipUnavailable
			}
			if rows.Close() != nil {
				return ErrOwnershipUnavailable
			}
			if conflict {
				return ErrOwnershipConflict
			}
		}
	}
	return nil
}

func detailReferencesKey(detail ownershipDetail, key GroupKey) bool {
	for _, refs := range [][]ownershipRef{detail.Ended, detail.Started} {
		for _, ref := range refs {
			if ref.GroupID == key.GroupID && ref.Revision == key.GroupRevision {
				return true
			}
		}
	}
	return false
}
