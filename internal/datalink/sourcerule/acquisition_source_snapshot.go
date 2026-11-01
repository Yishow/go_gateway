package sourcerule

import (
	"context"
	"database/sql"
	"errors"

	"go-gateway/internal/datalink/schema"
)

// ErrAcquisitionSourceSnapshot does not expose SQL, endpoint or stored values.
var ErrAcquisitionSourceSnapshot = errors.New("acquisition source snapshot unavailable or changed")

// ReadAcquisitionSourceRowsTx enumerates the complete device source set in the
// caller's transaction. It never filters disabled or standalone points, opens a
// second connection, probes a device, or consults recording/output selectors.
func ReadAcquisitionSourceRowsTx(ctx context.Context, tx *sql.Tx, deviceID string) (AcquisitionSourceRows, error) {
	var out AcquisitionSourceRows
	if tx == nil || deviceID == "" {
		return out, ErrAcquisitionSourceSnapshot
	}
	if tx.QueryRowContext(ctx, `SELECT id,protocol,status,connection_config FROM devices WHERE id=?`, deviceID).Scan(&out.Device.ID, &out.Device.Protocol, &out.Device.Status, &out.Device.ConnectionConfig) != nil {
		return out, ErrAcquisitionSourceSnapshot
	}
	var err error
	out.Rules, err = readShareList(ctx, tx, ListFilter{DeviceID: &deviceID})
	if err != nil {
		return out, ErrAcquisitionSourceSnapshot
	}
	for _, rule := range out.Rules {
		// The legacy shared scanner does not load layout fields. Do not silently
		// treat a persisted mixed rule as homogeneous source evidence.
		if tx.QueryRowContext(ctx, `SELECT COALESCE(layout_mode,''),COALESCE(mixed_items,'') FROM source_rules WHERE id=? AND device_id=?`, rule.ID, deviceID).Scan(&rule.LayoutMode, &rule.MixedItems) != nil {
			return out, ErrAcquisitionSourceSnapshot
		}
		links, err := readShareListLinks(ctx, tx, rule.ID)
		if err != nil {
			return out, ErrAcquisitionSourceSnapshot
		}
		out.Links = append(out.Links, links...)
	}
	out.Points, err = readAcquisitionDevicePoints(ctx, tx, deviceID)
	if err != nil || validateAcquisitionOwnership(ctx, tx, out) != nil {
		return out, ErrAcquisitionSourceSnapshot
	}
	return out, nil
}

func readAcquisitionDevicePoints(ctx context.Context, tx *sql.Tx, deviceID string) ([]*schema.Point, error) {
	rows, err := tx.QueryContext(ctx, `SELECT id,device_id,address,COALESCE(function,''),data_type,COALESCE(data_format,''),COALESCE(format_mode,''),COALESCE(format_override,''),format_version,revision_id,enabled,mode,polling_group_id FROM points WHERE device_id=? ORDER BY id`, deviceID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	points := make([]*schema.Point, 0)
	for rows.Next() {
		var point schema.Point
		var group sql.NullString
		if err := rows.Scan(&point.ID, &point.DeviceID, &point.Address, &point.Function, &point.DataType, &point.DataFormat, &point.FormatMode, &point.FormatOverride, &point.FormatVersion, &point.RevisionID, &point.Enabled, &point.Mode, &group); err != nil {
			return nil, err
		}
		if group.Valid {
			point.PollingGroupID = &group.String
		}
		points = append(points, &point)
	}
	return points, rows.Err()
}

func validateAcquisitionOwnership(ctx context.Context, tx *sql.Tx, input AcquisitionSourceRows) error {
	points := make(map[string]*schema.Point, len(input.Points))
	owners := make(map[string]*schema.SourceRuleLink, len(input.Links))
	for _, point := range input.Points {
		points[point.ID] = point
	}
	for _, link := range input.Links {
		point := points[link.PointID]
		if point == nil || point.DeviceID != input.Device.ID || point.Address != link.Address || owners[point.ID] != nil {
			return ErrAcquisitionSourceSnapshot
		}
		owners[point.ID] = link
		if _, err := readShareRecoveryPointRows(ctx, tx, link.PointID, input.Device.ID); err != nil {
			return ErrAcquisitionSourceSnapshot
		}
	}
	for _, point := range input.Points {
		// Include foreign-owner insertions, not merely links of this device's rules.
		var count int
		if tx.QueryRowContext(ctx, `SELECT count(*) FROM source_rule_links WHERE point_id=?`, point.ID).Scan(&count) != nil {
			return ErrAcquisitionSourceSnapshot
		}
		if (owners[point.ID] == nil && count != 0) || (owners[point.ID] != nil && count != 1) {
			return ErrAcquisitionSourceSnapshot
		}
	}
	return nil
}

// Shared with the published Share recovery witness; preserve its exact columns,
// row representation and ordering rather than creating a second point reader.
func readShareRecoveryPointRows(ctx context.Context, tx *sql.Tx, pointID, deviceID string) ([]any, error) {
	var pointDevice string
	if tx.QueryRowContext(ctx, `SELECT device_id FROM points WHERE id=?`, pointID).Scan(&pointDevice) != nil || pointDevice != deviceID {
		return nil, ErrShareRecoverySnapshot
	}
	return shareRecoveryRows(ctx, tx, `SELECT id,device_id,address,function,data_type,data_format,format_mode,format_override,format_version,revision_id,enabled,mode,polling_group_id FROM points WHERE id=?`, pointID)
}
