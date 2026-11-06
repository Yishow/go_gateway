package workspace

import (
	"context"
	"errors"
)

// RecordingHomeSourceInput contains only source fields used by the home projection.
// Readers must return complete lists, not pages, without repairing persisted state.
type RecordingHomeSourceInput struct {
	Workspace RecordingHomeSourceWorkspace `json:"workspace"`
	Devices   []RecordingHomeSourceDevice  `json:"devices"`
	Rules     []RecordingHomeSourceRule    `json:"rules"`
	Mappings  []RecordingHomeSourceMapping `json:"mappings"`
	Points    []RecordingHomeSourcePoint   `json:"points"`
	Links     []RecordingHomeSourceLink    `json:"links"`
}

// RecordingHomeSourceLink retains canonical link identity even without a mapping.
// It is comparison evidence only, never part of the presentation response/digest.
type RecordingHomeSourceLink struct {
	ID        string  `json:"id"`
	RuleID    string  `json:"rule_id"`
	PointID   string  `json:"point_id"`
	Address   string  `json:"address"`
	MappingID *string `json:"mapping_id"`
	TagID     *string `json:"tag_id"`
}

// RecordingHomeSourceWorkspace identifies the canonical workspace and ordered membership.
type RecordingHomeSourceWorkspace struct {
	ID               string   `json:"id"`
	OrderedDeviceIDs []string `json:"ordered_device_ids"`
}

// RecordingHomeSourceDevice excludes connection configuration and credentials.
type RecordingHomeSourceDevice struct {
	ID                 string  `json:"id"`
	Name               string  `json:"name"`
	Protocol           *string `json:"protocol"`
	AvailabilityStatus *string `json:"availability_status"`
}

// RecordingHomeSourceRule retains revision evidence even when candidates are unchanged.
type RecordingHomeSourceRule struct {
	ID               string   `json:"id"`
	WorkspaceID      string   `json:"workspace_id"`
	DeviceID         string   `json:"device_id"`
	RevisionID       string   `json:"revision_id"`
	Enabled          bool     `json:"enabled"`
	SkippedAddresses []string `json:"skipped_addresses"`
}

// RecordingHomeSourceMapping identifies a persisted mapping and its display fields.
type RecordingHomeSourceMapping struct {
	ID          string `json:"id"`
	WorkspaceID string `json:"workspace_id"`
	PointID     string `json:"point_id"`
	RuleID      string `json:"rule_id"`
	DeviceID    string `json:"device_id"`
	Address     string `json:"address"`
	TagID       string `json:"tag_id"`
	TagKey      string `json:"tag_key"`
	DisplayName string `json:"display_name"`
	TargetType  string `json:"target_type"`
	Enabled     bool   `json:"enabled"`
}

// RecordingHomeSourcePoint deliberately excludes live values and protocol configuration.
type RecordingHomeSourcePoint struct {
	ID       string `json:"id"`
	DeviceID string `json:"device_id"`
	Address  string `json:"address"`
	Enabled  bool   `json:"enabled"`
}

// RecordingHomeSourceReader reads one complete canonical projection without mutations.
type RecordingHomeSourceReader func(context.Context, string) (RecordingHomeSourceInput, error)

// RecordingHomeDeviceSummary is safe device presentation data.
type RecordingHomeDeviceSummary struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	Protocol     string `json:"protocol"`
	Availability string `json:"availability"`
}

// RecordingHomeCandidate identifies one saved Tag eligible for a group.
type RecordingHomeCandidate struct {
	Key        string `json:"key"`
	DeviceID   string `json:"device_id"`
	PointID    string `json:"point_id"`
	TagID      string `json:"tag_id"`
	TagKey     string `json:"tag_key"`
	Label      string `json:"label"`
	DeviceName string `json:"device_name"`
	Address    string `json:"address"`
	TargetType string `json:"target_type"`
}

// RecordingHomeExcluded describes a saved source lacking a Tag.
type RecordingHomeExcluded struct {
	Label  string `json:"label"`
	Reason string `json:"reason"`
}

// RecordingHomeSourceSummary preserves the legacy candidate ordering and exclusions.
type RecordingHomeSourceSummary struct {
	Devices    []RecordingHomeDeviceSummary `json:"devices"`
	Candidates []RecordingHomeCandidate     `json:"candidates"`
	Excluded   []RecordingHomeExcluded      `json:"excluded"`
}

// RecordingHomeSourcesObservation is presentation evidence, never a readiness token.
// Unknown observations contain no partial summary or digest.
type RecordingHomeSourcesObservation struct {
	WorkspaceID  string `json:"workspace_id"`
	State        string `json:"state"`
	SourceDigest string `json:"source_digest,omitempty"`
	ReasonCode   string `json:"reason_code,omitempty"`
	*RecordingHomeSourceSummary
}

// ErrRecordingHomeSourcesUnavailable indicates a failed complete source read.
var ErrRecordingHomeSourcesUnavailable = errors.New("recording home sources unavailable")

// ErrRecordingHomeInvalidScope indicates a missing workspace scope.
var ErrRecordingHomeInvalidScope = errors.New("recording home workspace scope required")

// ErrRecordingHomeWorkspaceNotFound can be returned by a reader for an unknown workspace.
var ErrRecordingHomeWorkspaceNotFound = errors.New("recording home workspace not found")

// ErrRecordingHomeSourcesInconsistent reports dangling or conflicting canonical links.
// Readers use it instead of repairing records or returning partial projections.
var ErrRecordingHomeSourcesInconsistent = errors.New("recording home sources inconsistent")
