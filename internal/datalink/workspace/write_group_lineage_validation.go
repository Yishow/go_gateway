package workspace

import (
	"context"
	"database/sql"
	"strings"
)

type lineageGeneration struct {
	number      int64
	predecessor sql.NullString
}

func (r *SQLWriteGroupRepository) validateLineageGenerations(ctx context.Context, runner writeGroupSQLRunner, metadata *writeGroupLineageMetadata) error {
	rows, err := runner.QueryContext(ctx, r.query(`SELECT g.group_id,g.generation,g.predecessor_group_id,w.workspace_id
		FROM write_group_generations g LEFT JOIN write_groups w ON w.id=g.group_id WHERE g.lineage_id=$1`), metadata.LineageID)
	if err != nil {
		return lineageCorrupt("read lineage generations: %v", err)
	}
	defer rows.Close()
	members := map[string]lineageGeneration{}
	for rows.Next() {
		var id string
		var workspace sql.NullString
		var member lineageGeneration
		if err := rows.Scan(&id, &member.number, &member.predecessor, &workspace); err != nil {
			return lineageCorrupt("decode lineage generation: %v", err)
		}
		if strings.TrimSpace(id) == "" || !workspace.Valid || workspace.String != metadata.WorkspaceID || member.number < 1 {
			return lineageCorrupt("generation group is missing or outside workspace")
		}
		members[id] = member
		if member.number == 1 {
			if metadata.GenesisGroupID != "" {
				return lineageCorrupt("multiple genesis generations")
			}
			metadata.GenesisGroupID = id
		}
	}
	if err := rows.Err(); err != nil {
		return lineageCorrupt("iterate lineage generations: %v", err)
	}
	if err := rows.Close(); err != nil {
		return lineageCorrupt("close lineage generations: %v", err)
	}
	if metadata.GenesisGroupID == "" {
		return lineageCorrupt("lineage has no genesis")
	}
	for id, member := range members {
		if member.number == 1 {
			if member.predecessor.Valid {
				return lineageCorrupt("genesis has a predecessor")
			}
			continue
		}
		previous, found := members[member.predecessor.String]
		if !member.predecessor.Valid || !found || id == member.predecessor.String || previous.number != member.number-1 {
			return lineageCorrupt("generation predecessor is not its prior lineage member")
		}
	}
	if _, found := members[metadata.CurrentGroupID]; !found {
		return lineageCorrupt("current head is not a lineage member")
	}
	reserved := metadata.ReservedGroupID
	if reserved.Valid != metadata.ReservationOperationID.Valid || reserved.Valid != metadata.ReservationIntentDigest.Valid {
		return lineageCorrupt("incomplete lineage reservation")
	}
	if reserved.Valid {
		_, found := members[reserved.String]
		if !found || reserved.String == metadata.CurrentGroupID || strings.TrimSpace(metadata.ReservationOperationID.String) == "" ||
			strings.TrimSpace(metadata.ReservationIntentDigest.String) == "" {
			return lineageCorrupt("invalid reserved lineage member or intent")
		}
	}
	return nil
}

func (r *SQLWriteGroupRepository) validateLineageAliases(ctx context.Context, runner writeGroupSQLRunner, metadata *writeGroupLineageMetadata) error {
	rows, err := runner.QueryContext(ctx, r.query(`SELECT workspace_id,group_id,lineage_id FROM write_group_basic_keys
		WHERE lineage_id=$1 OR group_id IN (SELECT group_id FROM write_group_generations WHERE lineage_id=$2)
		UNION ALL SELECT workspace_id,group_id,lineage_id FROM write_group_profile_keys
		WHERE lineage_id=$3 OR group_id IN (SELECT group_id FROM write_group_generations WHERE lineage_id=$4)`),
		metadata.LineageID, metadata.LineageID, metadata.LineageID, metadata.LineageID)
	if err != nil {
		return lineageCorrupt("read canonical aliases: %v", err)
	}
	defer rows.Close()
	count := 0
	for rows.Next() {
		var workspace, anchor string
		var lineage sql.NullString
		if err := rows.Scan(&workspace, &anchor, &lineage); err != nil {
			return lineageCorrupt("decode canonical alias: %v", err)
		}
		if workspace != metadata.WorkspaceID || anchor != metadata.GenesisGroupID || !lineage.Valid || lineage.String != metadata.LineageID {
			return lineageCorrupt("canonical aliases disagree on workspace, lineage or genesis")
		}
		count++
	}
	if err := rows.Err(); err != nil {
		return lineageCorrupt("iterate canonical aliases: %v", err)
	}
	if count == 0 {
		return lineageCorrupt("lineage has no canonical alias")
	}
	return nil
}
