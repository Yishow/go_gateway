package sourcerule

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"slices"

	"go-gateway/internal/datalink/schema"
)

// AcquisitionSourceRows is detached source input, not an installation grant.
// Only the workspace SQL snapshot/acceptance boundary can grant completeness.
type AcquisitionSourceRows struct {
	Device schema.Device
	Rules  []*schema.SourceRule
	Links  []*schema.SourceRuleLink
	Points []*schema.Point
}

type acquisitionRuleIdentity struct {
	ID, DeviceID, StartAddress string
	Count                      int
	DataType                   schema.DataType
	Enabled                    bool
	Skipped, DataFormat        string
	Layout, Mixed              string
}
type acquisitionPointIdentity struct {
	ID, DeviceID, Address, Function string
	DataType                        schema.DataType
	DataFormat                      string
	FormatMode                      schema.PointFormatMode
	FormatOverride                  string
	FormatVersion                   int
	Mode                            schema.PointMode
	PollingGroupID                  *string
	Enabled                         bool
	Revision                        string `json:",omitempty"`
}
type acquisitionLinkIdentity struct {
	RuleID, Address, PointID string
	ID                       string `json:",omitempty"`
}

// AcquisitionSourceDigest deliberately excludes recording/Share fields, tags,
// mapping IDs, names and timestamps. Opaque rule revisions can advance for an
// output-only edit; actual source fields, rather than that token, are witnessed.
// Witness=true additionally retains point revisions and link installation IDs.
func AcquisitionSourceDigest(input AcquisitionSourceRows, witness bool) (string, error) {
	rules := make([]acquisitionRuleIdentity, 0, len(input.Rules))
	links := make([]acquisitionLinkIdentity, 0, len(input.Links))
	points := make([]acquisitionPointIdentity, 0, len(input.Points))
	for _, rule := range input.Rules {
		if rule == nil {
			return "", ErrAcquisitionSourceSnapshot
		}
		layout := rule.LayoutMode
		if layout == "" {
			layout = "homogeneous"
		}
		rules = append(rules, acquisitionRuleIdentity{rule.ID, rule.DeviceID, rule.StartAddress, rule.Count, rule.DataType, rule.Enabled, rule.SkippedAddresses, rule.DataFormat, layout, rule.MixedItems})
	}
	for _, link := range input.Links {
		if link == nil {
			return "", ErrAcquisitionSourceSnapshot
		}
		value := acquisitionLinkIdentity{RuleID: link.RuleID, Address: link.Address, PointID: link.PointID}
		if witness {
			value.ID = link.ID
		}
		links = append(links, value)
	}
	for _, point := range input.Points {
		if point == nil {
			return "", ErrAcquisitionSourceSnapshot
		}
		value := acquisitionPointIdentity{ID: point.ID, DeviceID: point.DeviceID, Address: point.Address, Function: point.Function, DataType: point.DataType, DataFormat: point.DataFormat, FormatMode: point.FormatMode, FormatOverride: point.FormatOverride, FormatVersion: point.FormatVersion, Mode: point.Mode, PollingGroupID: point.PollingGroupID, Enabled: point.Enabled}
		if witness {
			value.Revision = point.RevisionID
		}
		points = append(points, value)
	}
	if sortAcquisitionIdentity(rules) != nil || sortAcquisitionIdentity(links) != nil || sortAcquisitionIdentity(points) != nil {
		return "", ErrAcquisitionSourceSnapshot
	}
	device := struct {
		ID, Config string
		Protocol   schema.ProtocolType
		Status     schema.DeviceStatus
	}{input.Device.ID, input.Device.ConnectionConfig, input.Device.Protocol, input.Device.Status}
	raw, err := json.Marshal(struct {
		Version string
		Witness bool
		Device  any
		Rules   []acquisitionRuleIdentity
		Links   []acquisitionLinkIdentity
		Points  []acquisitionPointIdentity
	}{"acquisition-source-v1", witness, device, rules, links, points})
	if err != nil {
		return "", ErrAcquisitionSourceSnapshot
	}
	digest := sha256.Sum256(raw)
	return hex.EncodeToString(digest[:]), nil
}

func sortAcquisitionIdentity[T any](values []T) error {
	type encodedValue struct {
		value T
		key   string
	}
	encoded := make([]encodedValue, len(values))
	for i, value := range values {
		raw, err := json.Marshal(value)
		if err != nil {
			return err
		}
		encoded[i] = encodedValue{value: value, key: string(raw)}
	}
	slices.SortFunc(encoded, func(a, b encodedValue) int {
		if a.key < b.key {
			return -1
		}
		if a.key > b.key {
			return 1
		}
		return 0
	})
	for i, value := range encoded {
		values[i] = value.value
	}
	return nil
}
