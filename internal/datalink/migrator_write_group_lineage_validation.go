package datalink

import (
	"context"
	"fmt"
)

func rejectSQLiteLineageRows(ctx context.Context, reader sqliteLineageReader, reason, query string) error {
	var invalid bool
	if err := reader.QueryRowContext(ctx, `SELECT EXISTS(`+query+`)`).Scan(&invalid); err != nil {
		return fmt.Errorf("validate write-group lineage %s: %w", reason, err)
	}
	if invalid {
		return fmt.Errorf("write-group lineage metadata is corrupt: %s", reason)
	}
	return nil
}

// Key shape and workspace checks also run before first initialization, while
// neither key table has a lineage reference yet. No existing payload is changed.
func validateSQLiteLineageKeyScopes(ctx context.Context, reader sqliteLineageReader) error {
	checks := []struct{ reason, query string }{
		{"canonical key scope", `SELECT 1 FROM (
			SELECT workspace_id,group_id FROM write_group_basic_keys UNION ALL
			SELECT workspace_id,group_id FROM write_group_profile_keys
		) k LEFT JOIN write_groups g ON g.id=k.group_id
		WHERE g.id IS NULL OR TRIM(k.workspace_id)='' OR TRIM(k.group_id)=''
			OR k.workspace_id<>g.workspace_id`},
		{"basic key shape", `SELECT 1 FROM write_group_basic_keys
			WHERE TRIM(device_id)='' OR canonical_role<>'basic'`},
		{"profile key shape", `SELECT 1 FROM write_group_profile_keys
			WHERE profile_kind NOT IN ('instantaneous','cumulative') OR NOT (
			(canonical_role='circuit' AND TRIM(circuit_id)<>'' AND selection_scope='') OR
			(canonical_role='unassigned' AND circuit_id='' AND TRIM(selection_scope)<>'' AND profile_kind='instantaneous'))`},
	}
	for _, check := range checks {
		if err := rejectSQLiteLineageRows(ctx, reader, check.reason, check.query); err != nil {
			return err
		}
	}
	return nil
}

func validateSQLiteWriteGroupLineages(ctx context.Context, reader sqliteLineageReader) error {
	checks := []struct{ reason, query string }{
		{"canonical genesis reference", `SELECT 1 FROM (
			SELECT workspace_id,group_id,lineage_id FROM write_group_basic_keys UNION ALL
			SELECT workspace_id,group_id,lineage_id FROM write_group_profile_keys
		) k LEFT JOIN write_group_generations g ON g.group_id=k.group_id
		LEFT JOIN write_group_lineages l ON l.lineage_id=k.lineage_id
		WHERE k.lineage_id IS NULL OR TRIM(k.lineage_id)='' OR g.group_id IS NULL
			OR l.lineage_id IS NULL OR g.lineage_id<>k.lineage_id OR g.generation<>1
			OR g.predecessor_group_id IS NOT NULL OR l.workspace_id<>k.workspace_id`},
		{"generation membership or predecessor", `SELECT 1 FROM write_group_generations g
		LEFT JOIN write_group_lineages l ON l.lineage_id=g.lineage_id
		LEFT JOIN write_groups w ON w.id=g.group_id
		LEFT JOIN write_group_generations p ON p.group_id=g.predecessor_group_id
		WHERE l.lineage_id IS NULL OR w.id IS NULL OR l.workspace_id<>w.workspace_id
			OR typeof(g.generation)<>'integer' OR g.generation<1
			OR (g.generation=1 AND g.predecessor_group_id IS NOT NULL)
			OR (g.generation>1 AND (p.group_id IS NULL OR p.lineage_id<>g.lineage_id OR p.generation<>g.generation-1))`},
		{"head membership or revision", `SELECT 1 FROM write_group_lineages l
		LEFT JOIN write_group_generations c ON c.group_id=l.current_group_id
		LEFT JOIN write_group_generations r ON r.group_id=l.reserved_group_id
		WHERE TRIM(l.lineage_id)='' OR TRIM(l.workspace_id)=''
			OR typeof(l.revision)<>'integer' OR l.revision<1
			OR c.group_id IS NULL OR c.lineage_id<>l.lineage_id
			OR (l.reserved_group_id IS NOT NULL AND (r.group_id IS NULL OR r.lineage_id<>l.lineage_id OR l.reserved_group_id=l.current_group_id))`},
		{"reservation tuple", `SELECT 1 FROM write_group_lineages
		WHERE (reserved_group_id IS NULL AND (reservation_operation_id IS NOT NULL OR reservation_intent_digest IS NOT NULL))
			OR (reserved_group_id IS NOT NULL AND (TRIM(reserved_group_id)='' OR reservation_operation_id IS NULL
				OR TRIM(reservation_operation_id)='' OR reservation_intent_digest IS NULL OR TRIM(reservation_intent_digest)=''))`},
		{"missing canonical genesis", `SELECT 1 FROM write_group_lineages l
		WHERE NOT EXISTS (SELECT 1 FROM write_group_generations g
			WHERE g.lineage_id=l.lineage_id AND g.generation=1 AND (
				EXISTS(SELECT 1 FROM write_group_basic_keys k WHERE k.group_id=g.group_id AND k.lineage_id=l.lineage_id) OR
				EXISTS(SELECT 1 FROM write_group_profile_keys k WHERE k.group_id=g.group_id AND k.lineage_id=l.lineage_id)))`},
	}
	for _, check := range checks {
		if err := rejectSQLiteLineageRows(ctx, reader, check.reason, check.query); err != nil {
			return err
		}
	}
	return nil
}
