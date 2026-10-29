package groupdelivery

import (
	"context"
	"database/sql"
	"encoding/json"
	"sort"
	"strings"
)

// LifecycleOutboxCounts labels SQL rows separately from collecting samples.
// Every state is a subset of Total; Outstanding excludes committed/skipped rows.
type LifecycleOutboxCounts struct {
	Outstanding  int `json:"outstanding"`
	Total        int `json:"total"`
	Queued       int `json:"queued"`
	Retrying     int `json:"retrying"`
	Blocked      int `json:"blocked"`
	Quarantined  int `json:"quarantined"`
	Unknown      int `json:"unknown"`
	SQLCommitted int `json:"sql_committed"`
	Skipped      int `json:"skipped"`
}

// LifecycleHandoff describes verified audit edges, not additional SQL rows.
// Fanout edges must not be summed into a unique-observation count.
type LifecycleHandoff struct {
	Origin  GroupKey `json:"origin"`
	Source  GroupKey `json:"source"`
	Target  GroupKey `json:"target"`
	State   string   `json:"state"`
	Samples int      `json:"samples"`
}

// LifecyclePending has no connector configuration, raw error or sample values.
type LifecyclePending struct {
	Samples  int                   `json:"samples"`
	Outbox   LifecycleOutboxCounts `json:"outbox"`
	Handoffs []LifecycleHandoff    `json:"handoffs"`
}

// ReadLifecyclePending reads one immutable scope using the caller's snapshot.
func ReadLifecyclePending(ctx context.Context, tx *sql.Tx, key GroupKey, destination FrozenDestination) (LifecyclePending, error) {
	if tx == nil || !key.Valid() || strings.TrimSpace(destination.Scope) == "" || destination.ConnectorID == "" || destination.ConnectorRevision == "" || destination.TableName == "" {
		return LifecyclePending{}, ErrOwnershipUnavailable
	}
	if _, err := LifecyclePendingRevisions(ctx, tx, key.WorkspaceID, key.GroupID); err != nil {
		return LifecyclePending{}, err
	}
	if err := validateHandoffScope(ctx, tx, key.GroupID, key.GroupRevision); err != nil {
		return LifecyclePending{}, err
	}
	var result LifecyclePending
	if err := tx.QueryRowContext(ctx, `SELECT COUNT(*) FROM wg_delivery_samples s WHERE workspace_id=? AND group_id=? AND group_revision=? AND consumed=0 AND NOT EXISTS (SELECT 1 FROM wg_delivery_sample_handoffs h WHERE h.source_seq=s.seq)`, key.WorkspaceID, key.GroupID, key.GroupRevision).Scan(&result.Samples); err != nil {
		return LifecyclePending{}, &handoffStorageError{cause: err}
	}
	counts, err := lifecycleOutbox(ctx, tx, key, destination)
	if err != nil {
		return LifecyclePending{}, err
	}
	result.Outbox = counts
	receipts, err := handoffReceiptsForKey(ctx, tx, key)
	if err != nil {
		return LifecyclePending{}, err
	}
	for _, r := range receipts {
		// Validate each edge before projecting its allowlisted identities. This also
		// retains compact finalized audit evidence after samples are reclaimed.
		if err := verifyHandoffReceipt(ctx, tx, r, map[string]bool{}); err != nil {
			return LifecyclePending{}, err
		}
		edge := LifecycleHandoff{Origin: r.Origin, Source: r.Source, Target: r.Target, State: r.State}
		found := false
		for i := range result.Handoffs {
			old := result.Handoffs[i]
			old.Samples = 0
			if old == edge {
				result.Handoffs[i].Samples++
				found = true
				break
			}
		}
		if !found {
			edge.Samples = 1
			result.Handoffs = append(result.Handoffs, edge)
		}
	}
	sort.Slice(result.Handoffs, func(i, j int) bool {
		a, b := result.Handoffs[i], result.Handoffs[j]
		ak := []string{a.Origin.WorkspaceID, a.Origin.GroupID, a.Origin.GroupRevision, a.Source.GroupID, a.Source.GroupRevision, a.Target.GroupID, a.Target.GroupRevision, a.State}
		bk := []string{b.Origin.WorkspaceID, b.Origin.GroupID, b.Origin.GroupRevision, b.Source.GroupID, b.Source.GroupRevision, b.Target.GroupID, b.Target.GroupRevision, b.State}
		for n := range ak {
			if ak[n] != bk[n] {
				return ak[n] < bk[n]
			}
		}
		return false
	})
	return result, nil
}

func lifecycleOutbox(ctx context.Context, tx *sql.Tx, key GroupKey, destination FrozenDestination) (LifecycleOutboxCounts, error) {
	rows, err := tx.QueryContext(ctx, `SELECT destination_scope,connector_id,connector_revision,database_name,table_schema,table_name,dedupe_capability,record_key_column,state,COUNT(*) FROM wg_delivery_outbox WHERE workspace_id=? AND group_id=? AND group_revision=? GROUP BY destination_scope,connector_id,connector_revision,database_name,table_schema,table_name,dedupe_capability,record_key_column,state`, key.WorkspaceID, key.GroupID, key.GroupRevision)
	if err != nil {
		return LifecycleOutboxCounts{}, &handoffStorageError{cause: err}
	}
	defer rows.Close()
	var counts LifecycleOutboxCounts
	for rows.Next() {
		var actual FrozenDestination
		var state string
		var n int
		if err := rows.Scan(&actual.Scope, &actual.ConnectorID, &actual.ConnectorRevision, &actual.Database, &actual.TableSchema, &actual.TableName, &actual.DedupeCapability, &actual.RecordKeyColumn, &state, &n); err != nil {
			return LifecycleOutboxCounts{}, &handoffStorageError{cause: err}
		}
		if actual != destination {
			return LifecycleOutboxCounts{}, ErrOwnershipUnavailable
		}
		switch state {
		case StatePending, StateSending:
			counts.Queued += n
		case StateRetrying:
			counts.Retrying += n
		case StateBlocked:
			counts.Blocked += n
		case StateQuarantined:
			counts.Quarantined += n
		case StateUnknown:
			counts.Unknown += n
		case StateCommitted:
			counts.SQLCommitted += n
		case StateSkipped:
			counts.Skipped += n
		default:
			return LifecycleOutboxCounts{}, ErrOwnershipUnavailable
		}
		counts.Total += n
		if state != StateCommitted && state != StateSkipped {
			counts.Outstanding += n
		}
	}
	if err := rows.Err(); err != nil {
		return LifecycleOutboxCounts{}, &handoffStorageError{cause: err}
	}
	return counts, nil
}

// LifecyclePendingRevisions discovers durable evidence, including orphan revisions.
// It does not infer that an absent immutable definition means an empty backlog.
func LifecyclePendingRevisions(ctx context.Context, tx *sql.Tx, workspaceID, groupID string) ([]string, error) {
	if tx == nil || strings.TrimSpace(workspaceID) == "" || strings.TrimSpace(groupID) == "" {
		return nil, ErrOwnershipUnavailable
	}
	rows, err := tx.QueryContext(ctx, `SELECT workspace_id,group_revision FROM wg_delivery_samples WHERE group_id=? UNION SELECT workspace_id,group_revision FROM wg_delivery_outbox WHERE group_id=? UNION SELECT workspace_id,source_revision FROM wg_delivery_sample_handoffs WHERE source_group_id=? UNION SELECT workspace_id,target_revision FROM wg_delivery_sample_handoffs WHERE target_group_id=?`, groupID, groupID, groupID, groupID)
	if err != nil {
		return nil, &handoffStorageError{cause: err}
	}
	defer rows.Close()
	var revisions []string
	for rows.Next() {
		var workspace, revision string
		if err := rows.Scan(&workspace, &revision); err != nil {
			return nil, &handoffStorageError{cause: err}
		}
		if workspace != workspaceID || strings.TrimSpace(revision) == "" {
			return nil, ErrOwnershipUnavailable
		}
		revisions = append(revisions, revision)
	}
	if err := rows.Err(); err != nil {
		return nil, &handoffStorageError{cause: err}
	}
	sort.Strings(revisions)
	return revisions, nil
}

// MarshalJSON projects only public identity fields without changing GroupKey's
// internal serialization, which participates in existing immutable audit hashes.
func (h LifecycleHandoff) MarshalJSON() ([]byte, error) {
	type identity struct {
		WorkspaceID   string `json:"workspace_id"`
		GroupID       string `json:"group_id"`
		GroupRevision string `json:"group_revision"`
	}
	key := func(k GroupKey) identity { return identity{k.WorkspaceID, k.GroupID, k.GroupRevision} }
	return json.Marshal(struct {
		Origin  identity `json:"origin"`
		Source  identity `json:"source"`
		Target  identity `json:"target"`
		State   string   `json:"state"`
		Samples int      `json:"samples"`
	}{key(h.Origin), key(h.Source), key(h.Target), h.State, h.Samples})
}
