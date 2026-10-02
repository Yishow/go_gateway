package workspace

import (
	"fmt"
	"slices"
	"strings"

	"go-gateway/internal/datalink/recordingplan"
)

func recordingPlanLayoutIssues(
	plan *recordingplan.RecordingPlan,
	connector recordingPlanMigrationConnector,
) []WriteGroupMigrationFinding {
	issues := make([]WriteGroupMigrationFinding, 0, 10)
	addIssue := func(code, message string) {
		issues = append(issues, WriteGroupMigrationFinding{Code: code, Message: message})
	}
	if !connector.enabled {
		addIssue("recording-plan-destination-disabled", "the selected saved connector is disabled")
	}
	if connector.kind != writeGroupConnectorKindSQLite && connector.kind != writeGroupConnectorKindPostgres {
		addIssue("recording-plan-destination-kind-unsupported", "the selected saved connector kind is outside the local migration adapter")
	}
	switch len(plan.Destinations) {
	case 0:
		addIssue("recording-plan-destination-missing", "the recording plan has no saved destination")
	case 1:
		destination := plan.Destinations[0]
		if strings.TrimSpace(destination.ConnectorID) == "" {
			addIssue("recording-plan-destination-missing", "the recording plan destination has no connector identity")
		} else if strings.TrimSpace(destination.ConnectorID) != connector.id {
			addIssue("recording-plan-destination-connector-mismatch", "the recording plan destination is outside the selected saved connector")
		}
		if revision := strings.TrimSpace(destination.ConnectorRevision); revision != "" && revision != connector.identityRevision {
			addIssue("recording-plan-destination-revision-stale", "the recording plan destination revision differs from the selected saved connector")
		}
	default:
		addIssue("recording-plan-multiple-destinations", "the recording plan references more than one destination")
	}
	// RecordingPlan has no canonical column or row-identity binding. Keep the
	// original plan and explain the repair path instead of fabricating either.
	addIssue("recording-plan-column-binding-missing", "the recording plan has no canonical destination column binding")
	addIssue("recording-plan-row-identity-missing", "the recording plan has no canonical row identity binding")
	if len(plan.Streams) == 0 {
		addIssue("recording-plan-stream-missing", "the recording plan has no persisted stream")
	}
	destinationIDs := make([]string, 0, len(plan.Destinations))
	for _, destination := range plan.Destinations {
		if id := strings.TrimSpace(destination.DestinationID); id != "" {
			destinationIDs = append(destinationIDs, id)
		}
	}
	for _, stream := range plan.Streams {
		switch stream.Mode {
		case recordingplan.StreamModeRawHistory:
			if stream.RawPolicy != recordingplan.RawPolicyEverySample {
				addIssue("recording-plan-raw-policy-unsupported", fmt.Sprintf("stream %s does not use the persisted every-sample raw policy", stream.StreamID))
			} else {
				addIssue("recording-plan-every-sample-not-equivalent", fmt.Sprintf("stream %s every-sample intent cannot be changed into a canonical snapshot implicitly", stream.StreamID))
			}
		default:
			addIssue("recording-plan-stream-mode-unsupported", fmt.Sprintf("stream %s uses unsupported mode %q", stream.StreamID, stream.Mode))
		}
		if len(stream.DestinationIDs) > 1 {
			addIssue("recording-plan-multiple-targets", fmt.Sprintf("stream %s references more than one destination", stream.StreamID))
		}
		for _, destinationID := range stream.DestinationIDs {
			if !slices.Contains(destinationIDs, strings.TrimSpace(destinationID)) {
				addIssue("recording-plan-destination-unresolved", fmt.Sprintf("stream %s references an unknown destination", stream.StreamID))
			}
		}
	}
	return issues
}
