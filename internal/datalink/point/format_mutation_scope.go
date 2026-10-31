package point

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"slices"
	"time"

	"go-gateway/internal/datalink/connector"
	"go-gateway/internal/datalink/decoding"
	"go-gateway/internal/datalink/schema"
)

type formatSnapshot struct {
	connectionConfigRaw string
	point               *schema.Point
	ruleRevision        string
	defaults            decoding.FormatContext
}

const pointFormatSelect = `SELECT id,device_id,name,description,address,function,data_type,data_format,COALESCE(format_mode,''),COALESCE(format_override,''),format_version,revision_id,mode,polling_group_id,last_read_at,last_value,last_error,enabled,created_at,updated_at FROM points WHERE id=?`

func (r *SQLRepository) formatSnapshot(ctx context.Context, tx *sql.Tx, req FormatMutationRequest) (formatSnapshot, error) {
	var workspaceRaw string
	if err := tx.QueryRowContext(ctx, `SELECT value FROM system_settings WHERE key='studio_v2_workspace'`).Scan(&workspaceRaw); err != nil {
		return formatSnapshot{}, ErrFormatUnavailable
	}
	var workspace struct {
		ID        string   `json:"id"`
		DeviceIDs []string `json:"ordered_device_ids"`
	}
	if json.Unmarshal([]byte(workspaceRaw), &workspace) != nil {
		return formatSnapshot{}, ErrFormatUnavailable
	}
	if workspace.ID != req.WorkspaceID || !slices.Contains(workspace.DeviceIDs, req.DeviceID) {
		return formatSnapshot{}, ErrFormatScope
	}
	p, err := r.scanPoint(tx.QueryRowContext(ctx, pointFormatSelect, req.PointID))
	if errors.Is(err, ErrPointNotFound) {
		return formatSnapshot{}, ErrFormatScope
	}
	if err != nil {
		return formatSnapshot{}, ErrFormatUnavailable
	}
	if p.DeviceID != req.DeviceID {
		return formatSnapshot{}, ErrFormatScope
	}
	var protocol schema.ProtocolType
	var configRaw string
	if err = tx.QueryRowContext(ctx, `SELECT protocol,connection_config FROM devices WHERE id=?`, p.DeviceID).Scan(&protocol, &configRaw); err != nil {
		return formatSnapshot{}, ErrFormatUnavailable
	}
	result := formatSnapshot{point: p, connectionConfigRaw: configRaw, defaults: decoding.FormatContext{Protocol: protocol, ConnectionRevision: connector.ReadConnectionRevision(configRaw)}}

	var ruleID, linkAddress string
	err = tx.QueryRowContext(ctx, `SELECT rule_id,address FROM source_rule_links WHERE point_id=?`, p.ID).Scan(&ruleID, &linkAddress)
	if errors.Is(err, sql.ErrNoRows) {
		if req.RuleID != "" {
			return formatSnapshot{}, ErrFormatScope
		}
		return result, nil
	}
	if err != nil {
		return formatSnapshot{}, ErrFormatUnavailable
	}
	if ruleID != req.RuleID || linkAddress != p.Address {
		return formatSnapshot{}, ErrFormatScope
	}
	var deviceID string
	var format sql.NullString
	if err = tx.QueryRowContext(ctx, `SELECT device_id,revision_id,data_format FROM source_rules WHERE id=?`, ruleID).Scan(&deviceID, &result.ruleRevision, &format); err != nil {
		return formatSnapshot{}, ErrFormatUnavailable
	}
	if deviceID != p.DeviceID {
		return formatSnapshot{}, ErrFormatScope
	}
	result.defaults.RuleOwned = true
	result.defaults.RuleFormat = format.String
	result.defaults.RuleRevision = result.ruleRevision
	return result, nil
}

type formatWorkspaceRevision struct {
	Revision string `json:"revision"`
	Dirty    bool   `json:"dirty"`
}

func loadFormatWorkspaceRevision(ctx context.Context, tx *sql.Tx, id string) (formatWorkspaceRevision, string, error) {
	var raw string
	err := tx.QueryRowContext(ctx, `SELECT value FROM system_settings WHERE key=?`, "modbus_share_workspace_revision:"+id).Scan(&raw)
	if errors.Is(err, sql.ErrNoRows) {
		return formatWorkspaceRevision{Revision: "rev-1"}, "", nil
	}
	if err != nil {
		return formatWorkspaceRevision{}, "", ErrFormatUnavailable
	}
	var record formatWorkspaceRevision
	if json.Unmarshal([]byte(raw), &record) != nil || record.Revision == "" {
		return formatWorkspaceRevision{}, "", ErrFormatUnavailable
	}
	return record, raw, nil
}

func advanceFormatWorkspace(ctx context.Context, tx *sql.Tx, id, previous, next string, owner formatPendingDirty) error {
	var record map[string]json.RawMessage
	if previous != "" {
		if json.Unmarshal([]byte(previous), &record) != nil {
			return ErrFormatUnavailable
		}
	} else {
		record = make(map[string]json.RawMessage)
	}
	record["revision"], _ = json.Marshal(next)
	record["dirty"] = json.RawMessage("true")
	record[formatPendingDirtyKey], _ = json.Marshal(owner)
	payload, _ := json.Marshal(record)
	key := "modbus_share_workspace_revision:" + id
	var result sql.Result
	var err error
	if previous == "" {
		result, err = tx.ExecContext(ctx, `INSERT INTO system_settings(key,value,updated_at) VALUES(?,?,?)`, key, string(payload), time.Now().UTC())
	} else {
		result, err = tx.ExecContext(ctx, `UPDATE system_settings SET value=?,updated_at=? WHERE key=? AND value=?`, string(payload), time.Now().UTC(), key, previous)
	}
	if err != nil {
		return ErrFormatUnavailable
	}
	n, err := result.RowsAffected()
	if err != nil {
		return ErrFormatUnavailable
	}
	if n != 1 {
		return ErrFormatConflict
	}
	return nil
}
