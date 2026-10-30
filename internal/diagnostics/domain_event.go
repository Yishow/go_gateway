package diagnostics

import (
	"bytes"
	"encoding/json"
	"io"
	"reflect"
	"regexp"
	"strconv"
	"time"

	"go-gateway/internal/datalink/runtimeview"
)

// DomainEventCode identifies a fixed observation, never a free-text error.
type DomainEventCode string

// DomainEvidenceKind separates evidence authorities; a connection or process event
// is not a delivery reference and is intentionally outside this domain contract.
type DomainEvidenceKind string

// DomainEventRef is an UNMOUNTED reference contract, not a delivery truth store.
// Binding.EntityRevision is the observed original source/mapping/group revision, NOT a
// latest draft. Locator is runtimeview's read-only historical hint, never permission.
// Group events do not attribute a shared result to individual points/tags.
type DomainEventRef struct {
	ContractVersion   int                       `json:"contract_version"`
	Code              DomainEventCode           `json:"code"`
	Binding           runtimeview.Binding       `json:"binding"`
	Locator           runtimeview.EntityLocator `json:"locator"`
	PointID           *string                   `json:"point_id"`
	TagID             *string                   `json:"tag_id"`
	ConnectorID       *string                   `json:"connector_id"`
	ConnectorRevision *string                   `json:"connector_revision"`
	Stage             runtimeview.Stage         `json:"stage"`
	Reason            runtimeview.Reason        `json:"reason"`
	Action            runtimeview.Action        `json:"action"`
	Evidence          DomainEvidenceRef         `json:"evidence"`
}

// DomainEvidenceRef is a closed union of references to existing authorities.
// ReferenceID is an opaque operation/acquisition/effect reference, never an error.
// Receipt needs actual durable identity, not a fabricated "committed" boolean.
// Intake requires exact Pipeline.Status revision. Quota keeps global/group scope.
// PendingSamples counts accepted input samples, NOT outbox rows/attention entries.
// No current/group pending age metric is exposed: its authority spans all revisions.
type DomainEvidenceRef struct {
	Kind           DomainEvidenceKind          `json:"kind"`
	ReferenceID    string                      `json:"reference_id"`
	Receipt        *runtimeview.Receipt        `json:"receipt"`
	Intake         *runtimeview.IntakeEvidence `json:"intake"`
	Quota          *runtimeview.QuotaEvidence  `json:"quota"`
	PendingSamples *string                     `json:"pending_samples"`
}

// DomainEventRule is a fixed code's structural evidence requirement.
type DomainEventRule struct {
	Code        DomainEventCode
	Stage       runtimeview.Stage
	Kind        DomainEvidenceKind
	LocatorStep int
	Success     bool
}

// DomainProducerContract records a reviewable integration plan, not active wiring.
type DomainProducerContract struct {
	Category                                                      string
	Events                                                        []DomainEventRule
	Authority, EmitBoundary, Identity, Evidence, UnknownOrFailure string
}

// DomainProducerMatrix returns independent copies of the integration plan. Paths
// are existing authorities to RECHECK after backend-owner handoff, not emitters now.
// Every boundary is after a confirmed result; capture failure must never change it.
// Restart/read-model replay must NOT emit an old receipt as a new commitment.
func DomainProducerMatrix() []DomainProducerContract {
	return []DomainProducerContract{
		{"source", []DomainEventRule{
			{"domain.source.applied", "accepted", "source-runtime-transition", 2, true},
			{"domain.source.rejected", "blocked", "source-runtime-transition", 2, false}},
			"internal/datalink/runtime/service_source_rule_reconcile.go: ReconcileSourceRule",
			"After confirmed aligned/rejected runtime reconciliation; never after draft save alone.",
			"workspace/device/rule and original source revision", "opaque reconcile operation; immutable effective source facts for applied",
			"Deferred/not-running/missing immutable revision proof is unavailable; never label a draft effective."},
		{"collector", []DomainEventRule{
			{"domain.collector.failed", "acquisition", "acquisition-result", 1, false},
			{"domain.collector.recovered", "acquisition", "acquisition-result", 1, true}},
			"internal/datalink/collector/scheduler_polling.go: pollPointWithAcquisition/collectedValueWithReason",
			"After final typed acquisition result; recovery only after same-scope failure and success.",
			"workspace/device/point and frozen source revision", "acquisition ID and classified result, never sample value",
			"Missing workspace/revision association is unavailable; connection health is not acquisition or SQL proof."},
		{"mapping", []DomainEventRule{
			{"domain.mapping.failed", "blocked", "mapping-result", 3, false},
			{"domain.mapping.recovered", "acquisition", "mapping-result", 3, true}},
			"internal/datalink/runtime/ingestor.go: mapping.ExecutePipeline result",
			"After transform result against the confirmed runtime binding, not mapping editor changes.",
			"workspace/device/mapping/point/tag and pinned mapping revision", "opaque acquisition/operation reference",
			"Missing binding/revision is unavailable; never expose input/output values or raw transform errors."},
		{"intake", []DomainEventRule{
			{"domain.group.applied", "collecting", "intake-status", 4, true},
			{"domain.group.paused", "blocked", "intake-status", 4, true},
			{"domain.group.retiring", "collecting", "intake-status", 4, true}},
			"internal/datalink/workspace/write_group_lifecycle.go; internal/datalink/grouppipeline/reconcile.go and pipeline.go: Status",
			"After applied/paused/retiring transition and matching Pipeline.Status GroupID/Revision observation.",
			"workspace/group/applied revision and frozen connector/revision", "lifecycle operation and exact-revision intake status",
			"Aggregate DeliveryView.Intake may select another revision: no exact status means unavailable, not current intake."},
		{"delivery", []DomainEventRule{
			{"domain.delivery.retrying", "retrying", "durable-transition", 4, false},
			{"domain.delivery.blocked", "blocked", "durable-transition", 4, false},
			{"domain.delivery.quarantined", "quarantined", "durable-transition", 4, false},
			{"domain.delivery.unknown", "unknown", "durable-transition", 4, false}},
			"internal/datalink/groupdelivery/sender.go: settle/retry/blockClaimed; outbox.go: MarkRetry/MarkUnknown/MarkQuarantined/MarkBlocked",
			"Only after the corresponding durable outbox transition succeeds.",
			"workspace/group/original revision and frozen connector/revision", "outbox effect reference; classified owner reason",
			"Failed/ambiguous persistence stays unavailable/unknown; never use HTTP 200, buffer ACK, or nil OnError as commitment."},
		{"commit", []DomainEventRule{
			{"domain.delivery.first-confirmed", "sql-committed", "runtime-receipt", 4, true},
			{"domain.delivery.recovered", "sql-committed", "runtime-receipt", 4, true},
			{"domain.delivery.summary", "sql-committed", "runtime-receipt", 4, true}},
			"internal/datalink/groupdelivery/outbox.go: CompleteDelivery; committed_effect.go: LastCommittedEffect",
			"After destination confirmation AND durable matching local receipt/outbox transaction; recovery matches original fault scope.",
			"workspace/group/original revision and frozen connector/revision", "actual record/effect/digest/time runtime receipt; never a receipt boolean",
			"Receipt mismatch/local receipt failure is unavailable; no readback/cleanup inference; Already/restart is not first-confirmed."},
		{"counter-replay", []DomainEventRule{
			{"domain.counter-replay.quarantined", "quarantined", "counter-replay-inputs", 4, false}},
			"internal/datalink/grouppipeline/counter_quarantine.go: setCounterQuarantine/CounterReplayQuarantines; groupdelivery/pending_samples.go",
			"After original-scope proof quarantine and exact PendingSampleCount read; recheck owner handoff.",
			"workspace/group/original revision and frozen connector/revision", "original proof/operation reference; pending_samples counts INPUTS, not outbox/attention length",
			"Count read failure is unavailable, not zero; preserve immutable evidence; never relabel epochs or resend finalized output."},
		{"recording-refusal", []DomainEventRule{{"domain.recording.refused", "blocked", "recording-refusal", 4, false}},
			"internal/datalink/grouppipeline/owner_observation.go: confirmed original-scope owner refusal",
			"After a typed native refusal with frozen owner evidence and exact pending INPUT count; never an outbox transition.",
			"workspace/original group revision and frozen connector/revision", "ephemeral observation UUID and exact pending input count, not durable effect or receipt",
			"Missing or changed original evidence is unavailable; no repair, reattribution or finalized-output resend."},
		{"quota", []DomainEventRule{
			{"domain.quota.refused", "blocked", "quota-refusal", 4, false},
			{"domain.quota.recovered", "accepted", "quota-status", 4, true}},
			"internal/datalink/groupdelivery/quota.go: checkQuota/QuotaStatus; ledger.go: AppendSample/CommitClosure",
			"After typed refusal, or exact retained-episode probe and native below-hard measurement; recovery is capacity only, not resumed intake or SQL.",
			"workspace/affected group/revision and frozen connector/revision", "refused operation reference and global/group quota evidence",
			"Unknown usage is unavailable, never zero; global usage is NOT group/current-revision usage; no pending-age metric here."},
		{"repair", []DomainEventRule{
			{"domain.repair.completed", "accepted", "operation-result", 4, true},
			{"domain.repair.blocked", "blocked", "operation-result", 4, false}},
			"internal/datalink/recordingplan/storage_restore_exec.go: runStorageRepair; repository_storage_repair.go: finishStorageRepair",
			"After operation ledger/lineage settlement, with frozen group destination mapped to that storage scope.",
			"workspace/group/original revision and frozen connector/revision", "opaque schema operation ID, not SQL or ledger payload",
			"Missing frozen scope/uncertain settlement is unavailable; never infer SQL delivery or permission to repair."},
		{"cutover", []DomainEventRule{
			{"domain.cutover.applied", "accepted", "operation-result", 4, true},
			{"domain.cutover.rejected", "blocked", "operation-result", 4, false}},
			"internal/datalink/workspace/write_group_cutover_apply.go: ApplyCutover/applyCutoverInTx",
			"After committed cutover or confirmed refusal, one original source-group reference per event, not a draft successor.",
			"workspace/source group/original revision and frozen connector/revision", "opaque cutover operation ID and owner-defined reason",
			"Ownership/handoff/finalized-overlap faults retain immutable evidence action; no retry/skip/DDL mutation authority."},
	}
}

var domainID = regexp.MustCompile(`^[A-Za-z0-9_-][A-Za-z0-9_.-]{0,127}$`)
var domainUTCTime = regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}(\.\d{1,9})?Z$`)

func domainUnavailable() (DomainEventRef, error) {
	return DomainEventRef{ContractVersion: 1, Code: "domain.unavailable", Stage: runtimeview.StageUnavailable,
			Reason: runtimeview.ReasonUnavailable, Action: runtimeview.ActionForReason(runtimeview.ReasonUnavailable),
			Evidence: DomainEvidenceRef{Kind: "unavailable"}}, &runtimeview.ContractError{
			Reason: runtimeview.ReasonUnavailable, Action: runtimeview.ActionForReason(runtimeview.ReasonUnavailable)}
}

// ParseDomainEventRef is pure fail-closed structural validation, NOT authorization,
// proof verification, a projector or emitter. A trusted future adapter must verify
// ownership, original revision, actual receipt origin and confirmed transition.
// The caller supplies the observation scope, not today's editable draft revision.
// Step 4 group-only references need no invented device; Steps 1-3 and device-bearing
// locators require an actual matching device scope. Ownership remains adapter work.
// On any missing/foreign/inconsistent evidence, no rejected identity is retained.
func ParseDomainEventRef(data []byte, scope runtimeview.Scope) (DomainEventRef, error) {
	if len(data) > MaxEventBytes || !domainJSONShape(data, reflect.TypeFor[DomainEventRef]()) {
		return domainUnavailable()
	}
	var v DomainEventRef
	if json.Unmarshal(data, &v) != nil || !domainID.MatchString(scope.WorkspaceID) || !domainID.MatchString(scope.WorkspaceRevision) {
		return domainUnavailable()
	}
	needsDevice := v.Locator.Step != 4 || v.Locator.DeviceID != nil
	if (needsDevice || scope.DeviceID != "") && !domainID.MatchString(scope.DeviceID) {
		return domainUnavailable()
	}
	if v.ContractVersion != 1 || v.Binding.WorkspaceID != scope.WorkspaceID || v.Binding.WorkspaceRevision != scope.WorkspaceRevision || !domainID.MatchString(v.Binding.EntityRevision) || runtimeview.ValidateEntityLocator(v.Locator, scope) != nil {
		return domainUnavailable()
	}
	if runtimeview.ParseStage(string(v.Stage)) != v.Stage || runtimeview.ParseReason(string(v.Reason)) != v.Reason || v.Action != runtimeview.ActionForReason(v.Reason) || !domainID.MatchString(v.Evidence.ReferenceID) {
		return domainUnavailable()
	}
	for _, id := range []*string{v.PointID, v.TagID, v.ConnectorID, v.ConnectorRevision} {
		if id != nil && !domainID.MatchString(*id) {
			return domainUnavailable()
		}
	}
	for _, producer := range DomainProducerMatrix() {
		for _, rule := range producer.Events {
			if rule.Code == v.Code && rule.Stage == v.Stage && rule.Kind == v.Evidence.Kind && rule.LocatorStep == v.Locator.Step && rule.Success == (v.Reason == runtimeview.ReasonNone) && validDomainIdentity(v) && validDomainEvidence(v) {
				return v, nil
			}
		}
	}
	return domainUnavailable()
}

func validDomainIdentity(v DomainEventRef) bool {
	if v.Locator.Step == 4 {
		return v.Locator.GroupRevision != nil && *v.Locator.GroupRevision == v.Binding.EntityRevision && v.ConnectorID != nil && v.ConnectorRevision != nil && v.PointID == nil && v.TagID == nil
	}
	if v.ConnectorID != nil || v.ConnectorRevision != nil {
		return false
	}
	switch v.Locator.Step {
	case 1:
		return v.PointID != nil && v.TagID == nil
	case 2:
		return v.PointID == nil && v.TagID == nil
	case 3:
		return v.PointID != nil && v.TagID != nil
	default:
		return false
	}
}

func validDomainEvidence(v DomainEventRef) bool {
	e := v.Evidence
	if !compatibleDomainReason(v) {
		return false
	}
	if (e.Receipt != nil) != (e.Kind == "runtime-receipt") || (e.Intake != nil) != (e.Kind == "intake-status") || (e.Quota != nil) != (e.Kind == "quota-refusal" || e.Kind == "quota-status") || (e.PendingSamples != nil) != (e.Kind == "counter-replay-inputs" || e.Kind == "recording-refusal") {
		return false
	}
	if r := e.Receipt; r != nil {
		if r.GroupRevision != v.Binding.EntityRevision || r.ConnectorRevision != *v.ConnectorRevision || r.EffectKey != e.ReferenceID || r.Provenance != "runtime" {
			return false
		}
		for _, id := range []string{r.RecordID, r.EffectKey, r.PayloadDigest} {
			if !domainID.MatchString(id) {
				return false
			}
		}
		if _, err := time.Parse(time.RFC3339Nano, r.CommittedAt); err != nil || !domainUTCTime.MatchString(r.CommittedAt) {
			return false
		}
	}
	if intake := e.Intake; intake != nil {
		state := map[DomainEventCode]string{"domain.group.applied": "active", "domain.group.paused": "not_running", "domain.group.retiring": "retiring"}[v.Code]
		if intake.GroupRevision != v.Binding.EntityRevision || intake.State != state || intake.Reason != nil {
			return false
		}
	}
	if q := e.Quota; q != nil {
		if (q.Scope != "global" && q.Scope != "group") || q.UsedBytes < 0 || q.MaxBytes < 0 || q.UsedBytes > 9007199254740991 || q.MaxBytes > 9007199254740991 {
			return false
		}
		if v.Code == "domain.quota.recovered" {
			if (q.State != "ok" && q.State != "warning") || q.IntakeRefused || q.MaxBytes == 0 || q.UsedBytes >= q.MaxBytes {
				return false
			}
		} else if q.State != "hard_limit" || !q.IntakeRefused {
			return false
		}
	}
	if e.PendingSamples != nil {
		n, err := strconv.ParseUint(*e.PendingSamples, 10, 64)
		if err != nil || strconv.FormatUint(n, 10) != *e.PendingSamples {
			return false
		}
	}
	return true
}

// This only constrains which producer can reference an existing owner reason;
// runtimeview still owns all reason values/actions. Task 2.4 preserves ownership
// faults on group delivery as well as cutover results. Their actual durable origin
// must be checked at backend handoff, never inferred from the event's code alone.
func compatibleDomainReason(v DomainEventRef) bool {
	switch v.Reason {
	case runtimeview.ReasonNone, runtimeview.ReasonUnavailable, runtimeview.ReasonStale,
		runtimeview.ReasonForeign, runtimeview.ReasonDeleted, runtimeview.ReasonUnsupported, runtimeview.ReasonMissingProof:
		return v.Evidence.Kind != "counter-replay-inputs" && v.Evidence.Kind != "recording-refusal"
	case runtimeview.ReasonCounterMissing, runtimeview.ReasonCounterConflict, runtimeview.ReasonCounterUnsupported:
		return v.Code == "domain.counter-replay.quarantined" && v.Evidence.Kind == "counter-replay-inputs"
	case runtimeview.ReasonOwnershipUnavailable, runtimeview.ReasonOwnershipConflict, runtimeview.ReasonHandoffConflict, runtimeview.ReasonFinalizedOverlap:
		return v.Locator.Step == 4 && (v.Evidence.Kind == "durable-transition" || (v.Code == "domain.recording.refused" && v.Evidence.Kind == "recording-refusal") || (v.Code == "domain.cutover.rejected" && v.Evidence.Kind == "operation-result"))
	default:
		return false
	}
}

// Only declared exact-case JSON keys are admitted, including inside shared DTOs.
// Duplicate keys, null required fields, trailing input and arbitrary nested payloads fail.
// This validates the event shape, not a parallel runtime reason/action registry.
func domainJSONShape(raw []byte, typ reflect.Type) bool {
	null := bytes.Equal(bytes.TrimSpace(raw), []byte("null"))
	if typ.Kind() == reflect.Pointer {
		return len(raw) == 0 || null || domainJSONShape(raw, typ.Elem())
	}
	if null || len(raw) == 0 {
		return false
	}
	if typ.Kind() != reflect.Struct {
		return true
	} // json.Unmarshal enforces scalar types.
	d := json.NewDecoder(bytes.NewReader(raw))
	tok, err := d.Token()
	if err != nil || tok != json.Delim('{') {
		return false
	}
	fields := make(map[string]reflect.Type, typ.NumField())
	for i := 0; i < typ.NumField(); i++ {
		f := typ.Field(i)
		fields[f.Tag.Get("json")] = f.Type
	}
	seen := map[string]bool{}
	for d.More() {
		tok, err = d.Token()
		key, ok := tok.(string)
		ft, exists := fields[key]
		if err != nil || !ok || !exists || seen[key] {
			return false
		}
		seen[key] = true
		var value json.RawMessage
		if d.Decode(&value) != nil || !domainJSONShape(value, ft) {
			return false
		}
	}
	if _, err = d.Token(); err != nil {
		return false
	}
	if _, err = d.Token(); err != io.EOF {
		return false
	}
	for key, ft := range fields {
		if !seen[key] && ft.Kind() != reflect.Pointer {
			return false
		}
	}
	return true
}
