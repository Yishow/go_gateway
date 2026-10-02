package workspace

import (
	"slices"
	"strings"

	"go-gateway/internal/datalink/schema"
)

// selectRowGroupMigrationSource keeps the source-rule check independent from
// the single-mapping adapter. A matched applied signature is provenance; only
// pending, mismatched or otherwise unsafe rule state blocks this batch.
func selectRowGroupMigrationSource(
	record *Record,
	sources []writeGroupMigrationSource,
	pointID string,
) (*writeGroupMigrationSource, []WriteGroupMigrationFinding, error) {
	pointSources := make([]writeGroupMigrationSource, 0)
	for _, source := range sources {
		if source.source.PointID == pointID || source.mapping.PointID == pointID {
			pointSources = append(pointSources, source)
		}
	}
	for _, source := range pointSources {
		if !slices.Contains(record.OrderedDeviceIDs, source.source.DeviceID) {
			return nil, nil, writeGroupNotFound("preview row-group migration source")
		}
	}
	enabledCount := 0
	for _, source := range sources {
		if source.mapping.Enabled {
			enabledCount++
		}
	}
	active := make([]writeGroupMigrationSource, 0, len(pointSources))
	for _, source := range pointSources {
		if source.mapping.Enabled && source.mappingState == string(schema.MappingStatusActive) {
			active = append(active, source)
		}
	}
	if enabledCount > 1 {
		return nil, []WriteGroupMigrationFinding{{Code: "multiple-enabled-source-mappings", Message: "more than one enabled source mapping reaches this row-group Tag"}}, nil
	}
	if len(active) == 0 {
		return nil, []WriteGroupMigrationFinding{{Code: "source-mapping-missing-or-disabled", Message: "no enabled active source mapping reaches this row-group point"}}, nil
	}
	source := active[0]
	if !slices.Contains(record.OrderedDeviceIDs, source.source.DeviceID) {
		return nil, nil, writeGroupNotFound("preview row-group migration source")
	}
	issues := make([]WriteGroupMigrationFinding, 0, 4)
	if source.deviceStatus == string(schema.DeviceStatusDisabled) {
		issues = append(issues, WriteGroupMigrationFinding{Code: "source-device-disabled", Message: "the source device is disabled"})
	}
	if !source.source.PointEnabled {
		issues = append(issues, WriteGroupMigrationFinding{Code: "source-point-disabled", Message: "the source point is disabled"})
	}
	if source.tagStatus != string(schema.TagStatusActive) {
		issues = append(issues, WriteGroupMigrationFinding{Code: "source-tag-missing-or-disabled", Message: "the source tag is not active"})
	}
	if pipeline := strings.TrimSpace(source.mapping.TransformPipeline); pipeline != "" && pipeline != "[]" {
		issues = append(issues, WriteGroupMigrationFinding{Code: "advanced-transform-pipeline", Message: "the mapping has a non-empty transform pipeline"})
	}
	if reason := strings.TrimSpace(source.mapping.BlockingReason); reason != "" {
		issues = append(issues, WriteGroupMigrationFinding{Code: "source-rule-blocking-reason", Message: "the source mapping has a persisted blocking reason"})
	}
	proposed := strings.TrimSpace(source.mapping.ProposedSignature)
	applied := strings.TrimSpace(source.mapping.LastAppliedSignature)
	switch {
	case proposed != "" && applied == "":
		issues = append(issues, WriteGroupMigrationFinding{Code: "source-rule-pending-signature", Message: "the source mapping has a proposed signature that has not been applied"})
	case proposed != "" && applied != "" && proposed != applied:
		issues = append(issues, WriteGroupMigrationFinding{Code: "source-rule-mismatched-signature", Message: "the source mapping proposed signature differs from its applied signature"})
	case proposed == "" && applied != "":
		issues = append(issues, WriteGroupMigrationFinding{Code: "source-rule-mismatched-signature", Message: "the source mapping has an applied signature without the matching proposed signature"})
	}
	if strings.TrimSpace(source.mapping.RuleCandidateID) != "" && proposed == "" && applied == "" {
		issues = append(issues, WriteGroupMigrationFinding{Code: "source-rule-metadata-incomplete", Message: "the source rule candidate has no verifiable applied signature"})
	}
	return &source, issues, nil
}
