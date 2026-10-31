package workspace

import (
	"context"
	"encoding/json"
	"go-gateway/internal/datalink/connector"
	"go-gateway/internal/datalink/decoding"
	"go-gateway/internal/datalink/measurement"
	"go-gateway/internal/datalink/schema"
)

type sourceRevisionSource struct {
	DeviceID, DeviceName, Protocol, ConnectionConfig                                                          string
	PointID, PointDeviceID, PointName, PointAddress, PointFunction, PointDataType, PointDataFormat, PointMode string
	PointEnabled                                                                                              bool
	TagID, TagKey, TagDisplayName, TagDataType, TagUnit, TagStatus, TagLabels                                 string
	pointFormat                                                                                               schema.Point
	profile                                                                                                   connector.ReadProfile
}

func (s sourceRevisionSource) records() (schema.Device, schema.Point, schema.Tag) {
	p := s.pointFormat
	p.ID = s.PointID
	p.DeviceID = s.PointDeviceID
	p.Name = s.PointName
	p.Address = s.PointAddress
	p.Function = s.PointFunction
	p.DataType = schema.DataType(s.PointDataType)
	p.DataFormat = s.PointDataFormat
	p.Mode = schema.PointMode(s.PointMode)
	p.Enabled = s.PointEnabled
	d := schema.Device{ID: s.DeviceID, Name: s.DeviceName, Protocol: schema.ProtocolType(s.Protocol), ConnectionConfig: s.ConnectionConfig}
	t := schema.Tag{ID: s.TagID, Key: s.TagKey, DisplayName: s.TagDisplayName, DataType: schema.DataType(s.TagDataType), Unit: s.TagUnit, Status: schema.TagStatus(s.TagStatus), Labels: s.TagLabels}
	return d, p, t
}
func (s sourceRevisionSource) revision() (string, error) {
	d, p, t := s.records()
	return measurement.SourceRevision(d, p, t, s.profile)
}

// Load from the same SQL transaction as source validation. Do not silently infer
// inheritance from a stale projected DataFormat or a reused address.
func loadSourceFormatProfile(ctx context.Context, runner writeGroupSQLRunner, query func(string) string, s *sourceRevisionSource) error {
	if s.PointID == "" || s.DeviceID == "" {
		return nil
	}
	p := &s.pointFormat
	var address, function, dtype, dataFormat, name, mode, protocol, connectionConfig string
	var enabled bool
	err := runner.QueryRowContext(ctx, query(`SELECT COALESCE(p.format_mode,''), COALESCE(p.format_override,''), p.format_version, p.revision_id,p.address,p.function,p.data_type,COALESCE(p.data_format,''),p.name,p.mode,p.enabled,d.protocol,d.connection_config FROM points p JOIN devices d ON d.id=p.device_id WHERE p.id=$1 AND p.device_id=$2`), s.PointID, s.DeviceID).Scan(&p.FormatMode, &p.FormatOverride, &p.FormatVersion, &p.RevisionID, &address, &function, &dtype, &dataFormat, &name, &mode, &enabled, &protocol, &connectionConfig)
	if err != nil {
		return err
	}
	if address != s.PointAddress || function != s.PointFunction || dtype != s.PointDataType || dataFormat != s.PointDataFormat || name != s.PointName || mode != s.PointMode || enabled != s.PointEnabled || protocol != s.Protocol || connectionConfig != s.ConnectionConfig {
		return connector.ErrReadProfile
	}
	if p.FormatMode == "" || p.FormatMode == schema.PointFormatLegacy {
		return nil
	}
	d, point, _ := s.records()
	defaults := decoding.FormatContext{Protocol: d.Protocol, ConnectionRevision: connector.ReadConnectionRevision(d.ConnectionConfig)}
	var config struct {
		DataFormat string `json:"data_format"`
	}
	if json.Unmarshal([]byte(d.ConnectionConfig), &config) != nil {
		return connector.ErrReadProfile
	}
	defaults.ConnectionFormat = config.DataFormat
	rows, err := runner.QueryContext(ctx, query(`SELECT l.rule_id,l.address,COALESCE(r.device_id,''),COALESCE(r.revision_id,''),COALESCE(r.data_format,'') FROM source_rule_links l LEFT JOIN source_rules r ON r.id=l.rule_id WHERE l.point_id=$1`), s.PointID)
	if err != nil {
		return err
	}
	defer rows.Close()
	count := 0
	for rows.Next() {
		count++
		var address, deviceID string
		if err := rows.Scan(&defaults.RuleID, &address, &deviceID, &defaults.RuleRevision, &defaults.RuleFormat); err != nil {
			return err
		}
		if count > 1 || address != point.Address || deviceID != d.ID || defaults.RuleRevision == "" {
			return connector.ErrReadProfile
		}
		defaults.RuleOwned = true
	}
	if err := rows.Err(); err != nil {
		return err
	}
	s.profile, err = connector.ResolveReadProfile(point, defaults)
	return err
}
