// Package runtimeview defines an UNMOUNTED, read-only projection, not recording authority.
// No API redaction, ownership authorization, reader or runtime activation is implemented here.
// JSON tags are the exhaustive safe-key allowlist; enum tags enumerate accepted wire values.
// Strings are bounded safe identifiers/labels, never raw connector DTOs, SQL or errors.
// Pointer fields serialize as null; missing nullable fields mean null, never inferred evidence.
// All other keys are required. Unknown keys/enums and malformed values fail unavailable.
// Scope/revision disagreement fails foreign-scope/stale-reference; no partial mixed result.
package runtimeview

import (
	"encoding/json"
	"errors"
	"strings"
)

// Stage observes a fact, never permission to write/retry/activate. Runtime commit is not readback.
type Stage string

const (
	StageUnavailable      Stage = "unavailable"
	StageAcquisition      Stage = "acquisition"
	StageAccepted         Stage = "accepted"
	StageCollecting       Stage = "collecting"
	StageQueued           Stage = "queued"
	StageRetrying         Stage = "retrying"
	StageBlocked          Stage = "blocked"
	StageQuarantined      Stage = "quarantined"
	StageUnknown          Stage = "unknown"
	StageSkipped          Stage = "skipped"
	StageSQLCommitted     Stage = "sql-committed"
	StageReadbackVerified Stage = "readback-verified"
	StageCleanupCompleted Stage = "cleanup-completed"
	StageTransportEvent   Stage = "transport-event"
)

// Action is a read-only next step. None confers a mutation capability or bypasses CAS.
type Action string

const (
	ActionNone                     Action = "none"
	ActionRefresh                  Action = "refresh"
	ActionInspectSettings          Action = "inspect-settings"
	ActionInspectDelivery          Action = "inspect-delivery"
	ActionOpenWorkspaceList        Action = "open-workspace-list"
	ActionInspectImmutableEvidence Action = "inspect-immutable-evidence"
)

// Reason projects owner-defined codes only; unknown owner codes must become unavailable.
type Reason string

const (
	ReasonNone                 Reason = "none"
	ReasonUnavailable          Reason = "unavailable"
	ReasonStale                Reason = "stale-reference"
	ReasonForeign              Reason = "foreign-scope"
	ReasonDeleted              Reason = "deleted"
	ReasonUnsupported          Reason = "unsupported"
	ReasonMissingProof         Reason = "effective-proof-missing"
	ReasonCounterMissing       Reason = "counter-replay-proof-missing"
	ReasonCounterConflict      Reason = "counter-replay-proof-conflict"
	ReasonCounterUnsupported   Reason = "counter-replay-proof-unsupported"
	ReasonOwnershipUnavailable Reason = "recording-ownership-unavailable"
	ReasonOwnershipConflict    Reason = "recording-ownership-conflict"
	ReasonHandoffConflict      Reason = "recording-handoff-proof-conflict"
	ReasonFinalizedOverlap     Reason = "recording-cutover-finalized-overlap"
)

// ContractError contains only safe reason/action, never the rejected input.
type ContractError struct {
	Reason Reason
	Action Action
}

func (e *ContractError) Error() string { return string(e.Reason) }

var errUnavailable = &ContractError{ReasonUnavailable, ActionRefresh}
var errStale = &ContractError{ReasonStale, ActionRefresh}
var errForeign = &ContractError{ReasonForeign, ActionOpenWorkspaceList}

// Scope is supplied by the caller's selected workspace/device read, not by names or the payload.
// Matching IDs is structural checking only; a future reader/resolver must verify ownership.
type Scope struct{ WorkspaceID, WorkspaceRevision, DeviceID string }

// Binding binds each independently sourced projection to ONE observation revision.
// EntityRevision is source/group revision, not WorkspaceRevision (draft B may coexist with A).
type Binding struct {
	WorkspaceID       string `json:"workspace_id"`
	WorkspaceRevision string `json:"workspace_revision"`
	EntityRevision    string `json:"entity_revision"`
}

// EntityLocator owns setting selection, not mutation. GroupRevision is a historical hint only.
// Steps 1/2/3 require device/rule/mapping respectively; step 4 requires group.
// Applicable IDs are nullable; unrelated IDs are refused. Resolver must check tombstones/ownership.
type EntityLocator struct {
	Step          int     `json:"step" min:"1" max:"4"`
	WorkspaceID   string  `json:"workspace_id"`
	DeviceID      *string `json:"device_id"`
	RuleID        *string `json:"rule_id"`
	MappingID     *string `json:"mapping_id"`
	GroupID       *string `json:"group_id"`
	GroupRevision *string `json:"group_revision"`
}

// RuleEvidenceView is the safe read envelope. Observation time/freshness come from the reader,
// not transport health. Missing timestamp cannot be fresh; stale preserves an explicit snapshot.
type RuleEvidenceView struct {
	ContractVersion   int            `json:"contract_version" min:"1" max:"1"`
	WorkspaceID       string         `json:"workspace_id"`
	WorkspaceRevision string         `json:"workspace_revision"`
	DeviceID          string         `json:"device_id"`
	ObservedAt        *string        `json:"observed_at" format:"time"`
	Freshness         string         `json:"freshness" enum:"fresh,stale,unavailable"`
	Source            *SourceView    `json:"source"`
	Groups            []GroupView    `json:"groups"`
	Locator           *EntityLocator `json:"locator"`
}

// SourceView saved facts come ONLY from SourceRule/Mapping/Point/PollingGroup.
// Effective requires immutable-source-facts, matching revision and full facts; projection version
// alone is insufficient. Missing proof clears effective and labels unavailable, never copies saved.
type SourceView struct {
	RuleID            string       `json:"rule_id"`
	PointID           string       `json:"point_id"`
	TagID             *string      `json:"tag_id"`
	MappingID         *string      `json:"mapping_id"`
	SavedRevision     *string      `json:"saved_revision"`
	EffectiveRevision *string      `json:"effective_revision"`
	State             string       `json:"state" enum:"saved-only,effective,unavailable"`
	Proof             string       `json:"proof" enum:"unavailable,immutable-source-facts"`
	Saved             *SourceFacts `json:"saved"`
	Effective         *SourceFacts `json:"effective"`
}

// SourceFacts: binding/source revision = SourceRule.RevisionID; mapping_revision = pinned mapping
// signature/revision, not current draft. enabled/address/read_type = SourceRule + linked Point;
// protocol_function/data_format = Point/frozen decode; result_type/transforms = validated Mapping.
// Unknown facts stay null; unsupported transform params are never copied or stringified.
type SourceFacts struct {
	Binding          Binding          `json:"binding"`
	MappingRevision  *string          `json:"mapping_revision"`
	Enabled          bool             `json:"enabled"`
	Address          *string          `json:"address"`
	ProtocolFunction *string          `json:"protocol_function"`
	ReadType         *string          `json:"read_type" enum:"bool,int16,uint16,int32,uint32,int64,uint64,float32,float64,string"`
	ResultType       *string          `json:"result_type" enum:"bool,int16,uint16,int32,uint32,int64,uint64,float32,float64,string"`
	Polling          *PollingEvidence `json:"polling"`
	Transforms       []TransformStep  `json:"transforms"`
}

// PollingEvidence is Point.PollingGroupID -> PollingGroup.IntervalMs, never write interval.
type PollingEvidence struct {
	GroupID    string `json:"group_id"`
	IntervalMS int64  `json:"interval_ms" min:"1"`
	Source     string `json:"source" enum:"polling-group"`
}

// TransformStep is a closed union: byte-order/order, scale/multiplier/offset, cast/target_type,
// or unsupported with no parameters. Nullable irrelevant fields MUST be null.
type TransformStep struct {
	Kind       string   `json:"kind" enum:"byte-order,scale,cast,unsupported"`
	Order      *string  `json:"order" enum:"ABCD,BADC,CDAB,DCBA"`
	Multiplier *float64 `json:"multiplier" signed:"true"`
	Offset     *float64 `json:"offset" signed:"true"`
	TargetType *string  `json:"target_type" enum:"bool,int16,uint16,int32,uint32,int64,uint64,float32,float64,string"`
}

// GroupView draft = WriteGroup; effective = ResolveAppliedAt + frozen RuntimeLayout.
// AppliedRevision is latest requested (possibly scheduled), NOT effective at AsOf.
// AsOf must match the server observation; absent new evidence safely degrades legacy payloads.
// AppliedEffectiveAt is the requested snapshot boundary, null unless independently proven;
// Effective.EffectiveAt belongs to the as-of effective snapshot, never copied from the request.
type GroupView struct {
	GroupID            string            `json:"group_id"`
	SavedRevision      *string           `json:"saved_revision"`
	EffectiveRevision  *string           `json:"effective_revision"`
	AsOf               *string           `json:"as_of" format:"time"`
	AppliedEffectiveAt *string           `json:"applied_effective_at" format:"time"`
	AppliedRevision    *string           `json:"applied_revision"`
	State              string            `json:"state" enum:"saved-only,effective,unavailable"`
	Proof              string            `json:"proof" enum:"unavailable,applied-snapshot-and-layout"`
	Saved              *GroupFacts       `json:"saved"`
	Effective          *GroupFacts       `json:"effective"`
	Delivery           *DeliveryEvidence `json:"delivery"`
}

// RecordingProfile is the exact saved/frozen profile; Unit measures interval, NOT physical units.
type RecordingProfile struct {
	Kind     string `json:"kind" enum:"instantaneous,cumulative"`
	Interval int64  `json:"interval" min:"1" max:"600"`
	Unit     string `json:"unit" enum:"second,minute"`
}

// GroupFacts fields come from the corresponding saved/applied WriteGroup, never a browser join.
// effective_at = snapshot.EffectiveAt (null for draft); interval normalized to seconds by authority;
// reduction = profile latest/mean/cumulative_delta; policies are not guessed for legacy/custom.
// EndsAt is that snapshot Group.EndsAt; ended/disabled/tombstoned retained facts grant no intake.
type GroupFacts struct {
	MixedRecording         *MixedRecordingFacts `json:"mixed_recording,omitempty"`
	RecordingProfile       *RecordingProfile    `json:"recording_profile"`
	AllowedLatenessSeconds *int64               `json:"allowed_lateness_seconds"`

	Binding          Binding     `json:"binding"`
	EndsAt           *string     `json:"ends_at" format:"time"`
	EffectiveAt      *string     `json:"effective_at" format:"time"`
	Enabled          bool        `json:"enabled"`
	Destination      Destination `json:"destination"`
	Members          []Member    `json:"members"`
	Layout           string      `json:"layout" enum:"managed_wide_v1,managed_long_v1,managed_wide_profile_v1,unavailable"`
	StorageMode      string      `json:"storage_mode" enum:"plain,timescale,native_monthly,unavailable"`
	IntervalSeconds  *int64      `json:"interval_seconds" min:"1"`
	Reduction        string      `json:"reduction" enum:"latest,mean,cumulative_delta,member_calculation,unavailable"`
	WriteMode        string      `json:"write_mode" enum:"append,latest,unavailable"`
	IncompletePolicy string      `json:"incomplete_policy" enum:"skip_row,partial,unavailable"`
}

// Destination is an explicit safe projection of frozen connector identity + group table.
// Alias is validated display-only text; schema/table/column are identifiers. Database is omitted
// because SQLite database can be a private path. No endpoint/user/DSN/error/SQL/params key exists.
type Destination struct {
	ConnectorID       string  `json:"connector_id"`
	ConnectorRevision string  `json:"connector_revision"`
	Alias             *string `json:"alias" format:"label"`
	Kind              string  `json:"kind" enum:"sqlite,postgres,mysql,sqlserver,unavailable"`
	Schema            *string `json:"schema"`
	Table             string  `json:"table"`
}

// Member is the pinned WriteGroupMember source/column identity, without values or arbitrary metadata.
// A group is selected when either saved or applied membership contains the selected device.
type Member struct {
	Calculation *MemberCalculationFacts `json:"calculation,omitempty"`
	// SourceValueType comes from this persisted member/verified frozen source, never current Tag.
	SourceValueType *string `json:"source_value_type" enum:"bool,int16,uint16,int32,uint32,int64,uint64,float32,float64,string"`

	// CircuitID is persisted WriteGroupMember.CircuitID from THIS snapshot, never a column/name.
	// Missing legacy evidence stays null; adapters normalize absent source identity, not infer it.
	CircuitID *string `json:"circuit_id"`
	DeviceID  string  `json:"device_id"`
	PointID   string  `json:"point_id"`
	TagID     string  `json:"tag_id"`
	Column    string  `json:"column"`
}

// DeliveryEvidence is ONLY DeliveryForRevision/DeliveryView projected in the same workspace read.
// Current is revision-specific; backlog retains historical destinations. Oldest pending age is
// DeliveryView.OldestPendingSeconds across ALL group revisions. TestWrite is the separate
// operation ledger (null without evidence), never used to populate runtime receipts or counts.
type DeliveryEvidence struct {
	CounterReplayQuarantines *[]CounterReplayQuarantine `json:"counter_replay_quarantines"`
	OldestPendingScope       string                     `json:"oldest_pending_scope" enum:"group-all-revisions"`
	OldestPendingSeconds     *float64                   `json:"oldest_pending_seconds"`
	Current                  *CurrentDelivery           `json:"current"`
	Backlog                  []Backlog                  `json:"backlog"`
	TestWrite                *TestWriteEvidence         `json:"test_write"`
}

// CounterReplayQuarantine preserves backend accepted-input evidence, NEVER finalized output.
// Null collection means unavailable, [] means observed empty. Counts must be JS-safe exact integers.
// Action is fixed backend repair guidance, not an executable action or new EvidenceAction.
type CounterReplayQuarantine struct {
	WorkspaceID    string `json:"workspace_id"`
	GroupID        string `json:"group_id"`
	GroupRevision  string `json:"group_revision"`
	PendingSamples int64  `json:"pending_samples"`
	Reason         Reason `json:"reason"`
	Action         string `json:"action" enum:"preserve-inputs-and-restore-original-revision-evidence-or-contact-maintenance"`
}

// IntakeEvidence requires Pipeline.Status GroupID/Revision matched to the enclosing group/current
// binding. Aggregate DeliveryView.Intake can prefer another revision and MUST NOT populate this.
// Missing revision proof means null. Missing/healthy reason is null; raw text is never copied.
type IntakeEvidence struct {
	GroupRevision string  `json:"group_revision"`
	State         string  `json:"state" enum:"active,retiring,blocked,not_running,unavailable"`
	Reason        *Reason `json:"reason"`
}

// CurrentDelivery totals are group-level, shared when members span devices; never per-Tag success.
// Timing is scheduler/outbox/attempt evidence ONLY; no interval-derived countdown is permitted.
type CurrentDelivery struct {
	Intake            *IntakeEvidence  `json:"intake"`
	Binding           Binding          `json:"binding"`
	TotalsScope       string           `json:"totals_scope" enum:"group"`
	Shared            bool             `json:"shared"`
	Stages            StageCounts      `json:"stages"`
	LastCommit        *Receipt         `json:"last_commit"`
	Sources           []SourceEvidence `json:"sources"`
	SourceState       string           `json:"source_state" enum:"available,unavailable"`
	NextDueAt         *string          `json:"next_due_at" format:"time"`
	NextRetryAt       *string          `json:"next_retry_at" format:"time"`
	AttemptDurationMS *float64         `json:"attempt_duration_ms"`
	Quota             *QuotaEvidence   `json:"quota"`
	Storage           *StorageEvidence `json:"storage"`
}

// StageCounts mirrors durable RevisionStages, not collector health or transport ACKs.
type StageCounts struct {
	Collecting   int64 `json:"collecting"`
	Queued       int64 `json:"queued"`
	Retrying     int64 `json:"retrying"`
	Blocked      int64 `json:"blocked"`
	Quarantined  int64 `json:"quarantined"`
	Unknown      int64 `json:"unknown"`
	SQLCommitted int64 `json:"sql_committed"`
	Skipped      int64 `json:"skipped"`
}

// Receipt projects matching LastCommittedEffect/SourceRecord receipt identity, runtime only.
// No readback/cleanup success may be derived from CommittedAt.
type Receipt struct {
	GroupRevision     string `json:"group_revision"`
	ConnectorRevision string `json:"connector_revision"`
	RecordID          string `json:"record_id"`
	EffectKey         string `json:"effect_key"`
	PayloadDigest     string `json:"payload_digest"`
	CommittedAt       string `json:"committed_at" format:"time"`
	Provenance        string `json:"provenance" enum:"runtime"`
}

// SourceEvidence projects source windows/coverage/delta-known without values. Null receipt never
// means zero commits; unavailable legacy/custom source evidence remains unavailable, not [] success.
type SourceEvidence struct {
	DeviceID      string        `json:"device_id"`
	PointID       string        `json:"point_id"`
	TagID         string        `json:"tag_id"`
	Latest        *SourceResult `json:"latest"`
	LastCommitted *SourceResult `json:"last_committed"`
}

// SourceResult keeps latest local windows separate from the last receipt-confirmed window.
type SourceResult struct {
	Calculation *CalculationResult `json:"calculation,omitempty"`
	State       string             `json:"state" enum:"available,missing,unknown,partial,not_due,unavailable"`
	WindowStart *string            `json:"window_start" format:"time"`
	WindowEnd   *string            `json:"window_end" format:"time"`
	DeltaKnown  *bool              `json:"delta_known"`
	Coverage    *string            `json:"coverage" enum:"complete,partial,unknown"`
	Receipt     *Receipt           `json:"receipt"`
}

// Backlog identity is historical DeliveryView.Backlog, not today's editable group destination.
type Backlog struct {
	Binding     Binding     `json:"binding"`
	Destination Destination `json:"destination"`
	Pending     int64       `json:"pending"`
	Reasons     []Reason    `json:"reasons"`
}

// TestWriteEvidence uses independent operation ledger timestamps; all absent until proven.
// Successful owned cleanup can coexist with absent/failed readback; neither implies the other.
type TestWriteEvidence struct {
	OperationID string  `json:"operation_id"`
	Provenance  string  `json:"provenance" enum:"test-write"`
	WrittenAt   *string `json:"written_at" format:"time"`
	ReadbackAt  *string `json:"readback_at" format:"time"`
	CleanupAt   *string `json:"cleanup_at" format:"time"`
}

// QuotaEvidence preserves DeliveryView.Quota scope: global usage is not group usage.
// Scope is not revision-specific; absent quota is null, never a fabricated capacity or raw notice.
type QuotaEvidence struct {
	Scope         string `json:"scope" enum:"global,group"`
	State         string `json:"state" enum:"unconfigured,ok,warning,hard_limit,unavailable"`
	UsedBytes     int64  `json:"used_bytes"`
	MaxBytes      int64  `json:"max_bytes"`
	IntakeRefused bool   `json:"intake_refused"`
}

// StorageEvidence is nullable unless an adapter proves frozen destination/revision scope against
// storage lineage. DeliveryView.Storage alone uses editable group scope and is NOT sufficient.
// Events are typed summaries, no raw ledger payload.
type StorageEvidence struct {
	Generation          int64          `json:"generation"`
	HistoryMayBeLost    bool           `json:"history_may_be_lost"`
	Reconstructions     int64          `json:"reconstructions"`
	LastReconstructedAt *string        `json:"last_reconstructed_at" format:"time"`
	RepairState         string         `json:"repair_state" enum:"unavailable,none,blocked"`
	Events              []StorageEvent `json:"events"`
}

// StorageEvent is a safe storage reconstruction observation, not permission to repair.
type StorageEvent struct {
	Kind string  `json:"kind" enum:"confirmed,recreated_by_confirmation,reconstructed,repair_blocked,unavailable"`
	At   *string `json:"at" format:"time"`
}

// IsContractError checks a safe error without exposing rejected payload details.
func IsContractError(err error, reason Reason) bool {
	var e *ContractError
	return errors.As(err, &e) && e.Reason == reason
}

// MarshalRuleEvidence serializes only a validated projection, preserving nullable evidence.
func MarshalRuleEvidence(v RuleEvidenceView, scope Scope) ([]byte, error) {
	raw, err := json.Marshal(v)
	if err != nil {
		return nil, errUnavailable
	}
	checked, err := ParseRuleEvidence(raw, scope)
	if err != nil {
		return nil, err
	}
	return json.Marshal(checked)
}

const reasonValues = "none,unavailable,stale-reference,foreign-scope,deleted,unsupported,effective-proof-missing,counter-replay-proof-missing,counter-replay-proof-conflict,counter-replay-proof-unsupported,recording-ownership-unavailable,recording-ownership-conflict,recording-handoff-proof-conflict,recording-cutover-finalized-overlap"
const stageValues = "unavailable,acquisition,accepted,collecting,queued,retrying,blocked,quarantined,unknown,skipped,sql-committed,readback-verified,cleanup-completed,transport-event"

// ParseStage degrades future stages safely; it cannot assert any evidence by itself.
func ParseStage(value string) Stage {
	if member(value, stageValues) {
		return Stage(value)
	}
	return StageUnavailable
}

// ParseReason refuses raw errors and future reason codes without reflecting their contents.
func ParseReason(value string) Reason {
	if member(value, reasonValues) {
		return Reason(value)
	}
	return ReasonUnavailable
}

// ActionForReason is the single shared, read-only reason/action mapping for diagnostics.
func ActionForReason(reason Reason) Action {
	switch reason {
	case ReasonNone:
		return ActionNone
	case ReasonForeign, ReasonDeleted:
		return ActionOpenWorkspaceList
	case ReasonUnsupported:
		return ActionInspectSettings
	case ReasonCounterMissing, ReasonCounterConflict, ReasonCounterUnsupported, ReasonOwnershipUnavailable, ReasonOwnershipConflict, ReasonHandoffConflict, ReasonFinalizedOverlap:
		return ActionInspectImmutableEvidence
	default:
		return ActionRefresh
	}
}
func member(value, values string) bool {
	for _, v := range strings.Split(values, ",") {
		if value == v {
			return true
		}
	}
	return false
}
