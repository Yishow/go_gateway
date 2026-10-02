package workspace

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"slices"
	"strings"

	"go-gateway/internal/datalink/measurement"
	"go-gateway/internal/datalink/recordingplan"
	"go-gateway/internal/datalink/schema"
)

const (
	recordingPlanSourceStatusResolved = "resolved"
	recordingPlanSourceStatusBlocked  = "blocked"
)

type recordingPlanMigrationSourceRead struct {
	measurementID string
	definition    *measurement.MeasurementDefinition
	source        *writeGroupMigrationSource
	snapshot      WriteGroupRecordingPlanMigrationSource
	issues        []WriteGroupMigrationFinding
}

func readRecordingPlanMigrationSource(
	ctx context.Context,
	tx *sql.Tx,
	repo *SQLWriteGroupRepository,
	measurementRepo *measurement.SQLRepository,
	record *Record,
	measurementID string,
) (recordingPlanMigrationSourceRead, error) {
	result := recordingPlanMigrationSourceRead{
		measurementID: measurementID,
		snapshot: WriteGroupRecordingPlanMigrationSource{
			MeasurementID: measurementID,
			Status:        recordingPlanSourceStatusBlocked,
		},
		issues: make([]WriteGroupMigrationFinding, 0, 4),
	}
	definition, err := measurementRepo.GetByIDInTx(ctx, tx, measurementID)
	if errors.Is(err, sql.ErrNoRows) {
		result.issues = append(result.issues, WriteGroupMigrationFinding{
			Code:    "recording-plan-source-missing",
			Message: "the recording plan references a measurement that no longer exists",
		})
		return result, nil
	}
	if err != nil {
		return result, fmt.Errorf("read recording plan measurement %s: %w", measurementID, err)
	}
	if definition.WorkspaceID != record.ID || !slices.Contains(record.OrderedDeviceIDs, definition.DeviceID) {
		return result, writeGroupNotFound("recording plan migration source")
	}
	result.definition = definition
	result.snapshot.DeviceID = definition.DeviceID
	result.snapshot.PointID = definition.PointID
	if definition.TagID != nil {
		result.snapshot.TagID = strings.TrimSpace(*definition.TagID)
	}
	result.snapshot.DefinitionRevision = definition.DefinitionRevision
	result.snapshot.SourceBindingRevision = definition.SourceBindingRevision
	result.snapshot.SeriesEpoch = definition.SeriesEpoch
	if strings.TrimSpace(definition.DefinitionRevision) == "" ||
		strings.TrimSpace(definition.SourceBindingRevision) == "" ||
		strings.TrimSpace(definition.SeriesEpoch) == "" {
		result.issues = append(result.issues, WriteGroupMigrationFinding{
			Code:    "recording-plan-source-revision-missing",
			Message: "the measurement source does not have a complete persisted revision identity",
		})
	}

	if strings.TrimSpace(result.snapshot.PointID) == "" || result.snapshot.TagID == "" {
		result.snapshot.SourceRevision, err = recordingPlanSourceRevision(definition, nil, "")
		if err != nil {
			return result, err
		}
		result.issues = append(result.issues, WriteGroupMigrationFinding{
			Code:    "recording-plan-source-unresolved",
			Message: "the measurement does not identify a complete point and tag source",
		})
		return result, nil
	}
	allSources, err := readMigrationPreviewSources(ctx, tx, repo, result.snapshot.TagID)
	if err != nil {
		return result, err
	}
	source, sourceIssues, err := selectRowGroupMigrationSource(record, allSources, result.snapshot.PointID)
	if err != nil {
		return result, err
	}
	result.issues = append(result.issues, sourceIssues...)
	if source == nil {
		result.snapshot.SourceRevision, err = recordingPlanSourceRevision(definition, nil, "", allSources)
		if err != nil {
			return result, err
		}
		return result, nil
	}
	result.source = source
	if source.source.DeviceID != result.snapshot.DeviceID ||
		source.source.PointID != result.snapshot.PointID ||
		source.source.TagID != result.snapshot.TagID {
		result.issues = append(result.issues, WriteGroupMigrationFinding{
			Code:    "recording-plan-source-mismatch",
			Message: "the measurement source chain does not match its persisted device, point and tag identity",
		})
	}
	if source.deviceStatus != string(schema.DeviceStatusActive) {
		result.issues = append(result.issues, WriteGroupMigrationFinding{
			Code:    "recording-plan-source-device-not-active",
			Message: "the measurement source device is not active",
		})
	}
	mappingRevision, err := source.mapping.revision()
	if err != nil {
		return result, fmt.Errorf("derive recording plan mapping revision: %w", err)
	}
	result.snapshot.MappingRevision = mappingRevision
	sourceRevision, err := recordingPlanSourceRevision(definition, source, mappingRevision)
	if err != nil {
		return result, err
	}
	result.snapshot.SourceRevision = sourceRevision
	result.snapshot.Status = recordingPlanSourceStatusResolved
	if len(result.issues) > 0 {
		result.snapshot.Status = recordingPlanSourceStatusBlocked
	}
	return result, nil
}

func recordingPlanSourceRevision(
	definition *measurement.MeasurementDefinition,
	source *writeGroupMigrationSource,
	mappingRevision string,
	candidates ...[]writeGroupMigrationSource,
) (string, error) {
	measurementPayload, err := json.Marshal(definition)
	if err != nil {
		return "", fmt.Errorf("encode recording plan measurement revision: %w", err)
	}
	sourceRevision := ""
	sourceID := ""
	mappingParts := make([]string, 0)
	if source != nil {
		sourceRevision, err = source.source.revision()
		if err != nil {
			return "", fmt.Errorf("derive recording plan source revision: %w", err)
		}
		sourceID = source.mapping.ID
		mappingParts = append(mappingParts, source.mapping.ID, mappingRevision, sourceRevision)
	} else if len(candidates) > 0 {
		for _, candidate := range candidates[0] {
			candidateMappingRevision, revisionErr := candidate.mapping.revision()
			if revisionErr != nil {
				return "", fmt.Errorf("derive unresolved recording plan mapping revision: %w", revisionErr)
			}
			candidateSourceRevision, revisionErr := candidate.source.revision()
			if revisionErr != nil {
				return "", fmt.Errorf("derive unresolved recording plan source revision: %w", revisionErr)
			}
			mappingParts = append(mappingParts, candidate.mapping.ID, candidateMappingRevision, candidateSourceRevision)
		}
	}
	return hashWriteGroupRevision(
		"recording-plan-source",
		string(measurementPayload),
		sourceID,
		mappingRevision,
		sourceRevision,
		strings.Join(mappingParts, "\x00"),
	)
}

func recordingPlanSourceSnapshot(read recordingPlanMigrationSourceRead) WriteGroupRecordingPlanMigrationSource {
	snapshot := read.snapshot
	if read.definition == nil {
		return snapshot
	}
	if read.source == nil {
		return snapshot
	}
	if read.source.source.DeviceID != "" {
		snapshot.DeviceID = read.source.source.DeviceID
	}
	if read.source.source.PointID != "" {
		snapshot.PointID = read.source.source.PointID
	}
	if read.source.source.TagID != "" {
		snapshot.TagID = read.source.source.TagID
	}
	if len(read.issues) > 0 {
		snapshot.Status = recordingPlanSourceStatusBlocked
	}
	return snapshot
}

func recordingPlanSourceIDs(plan *recordingplan.RecordingPlan) []string {
	seen := make(map[string]struct{}, len(plan.Members)+len(plan.Streams))
	for _, member := range plan.Members {
		if id := strings.TrimSpace(member.MeasurementID); id != "" {
			seen[id] = struct{}{}
		}
	}
	for _, stream := range plan.Streams {
		if id := strings.TrimSpace(stream.MeasurementID); id != "" {
			seen[id] = struct{}{}
		}
	}
	ids := make([]string, 0, len(seen))
	for id := range seen {
		ids = append(ids, id)
	}
	slices.Sort(ids)
	return ids
}
