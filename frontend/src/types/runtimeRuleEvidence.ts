/** UNMOUNTED safe read contract, not a recording/authorization authority.
 * JSON keys exactly match Go runtimeview/contract.go. Every non-nullable key is required;
 * omitted nullable keys normalize to null. Unknown keys/enums fail unavailable.
 * Cross-workspace/device fails foreign-scope; mixed observation/entity revisions fail
 * stale-reference. Missing effective proof drops effective, never substitutes saved facts.
 * No endpoint, credential, database path, raw SQL/error or arbitrary transform params exist.
 */
export type EvidenceStage = 'unavailable' | 'acquisition' | 'accepted' | 'collecting' |
  'queued' | 'retrying' | 'blocked' | 'quarantined' | 'unknown' | 'skipped' |
  'sql-committed' | 'readback-verified' | 'cleanup-completed' | 'transport-event'
export type EvidenceAction = 'none' | 'refresh' | 'inspect-settings' | 'inspect-delivery' |
  'open-workspace-list' | 'inspect-immutable-evidence'
export type EvidenceReason = 'none' | 'unavailable' | 'stale-reference' | 'foreign-scope' |
  'deleted' | 'unsupported' | 'effective-proof-missing' | 'counter-replay-proof-missing' |
  'counter-replay-proof-conflict' | 'counter-replay-proof-unsupported' |
  'recording-ownership-unavailable' | 'recording-ownership-conflict' |
  'recording-handoff-proof-conflict' | 'recording-cutover-finalized-overlap'
export type EvidenceResult<T> = { ok: true; value: T } |
  { ok: false; reason: 'unavailable' | 'stale-reference' | 'foreign-scope'; action: EvidenceAction }
/** Caller-selected scope; ID equality is NOT ownership authorization. */
export interface EvidenceScope { workspace_id: string; workspace_revision: string; device_id: string }
/** workspace_revision is the read observation; entity_revision is independently saved/applied. */
export interface RevisionBinding { workspace_id: string; workspace_revision: string; entity_revision: string }
/** Historical group_revision only selects; never replaces mutation CAS. Steps require exact IDs. */
export interface EntityLocator {
  step: 1 | 2 | 3 | 4
  workspace_id: string
  device_id: string | null
  rule_id: string | null
  mapping_id: string | null
  group_id: string | null
  group_revision: string | null
}
/** Reader observation freshness, never inferred from SSE/transport health. */
export interface RuleEvidenceView extends EvidenceScope {
  contract_version: 1
  observed_at: string | null
  freshness: 'fresh' | 'stale' | 'unavailable'
  source: SourceEvidenceView | null
  groups: RecordingGroupView[]
  locator: EntityLocator | null
}
export type EvidenceState = 'saved-only' | 'effective' | 'unavailable'
export type EvidenceDataType = 'bool' | 'int16' | 'uint16' | 'int32' | 'uint32' |
  'int64' | 'uint64' | 'float32' | 'float64' | 'string'
/** saved: SourceRule/Point/Mapping; effective: proven immutable source facts, not projection version. */
export interface SourceEvidenceView {
  rule_id: string; point_id: string; tag_id: string | null; mapping_id: string | null
  saved_revision: string | null; effective_revision: string | null
  state: EvidenceState; proof: 'unavailable' | 'immutable-source-facts'
  saved: SourceFacts | null; effective: SourceFacts | null
}
/** SourceRule.RevisionID and pinned mapping signature; enabled/address/function/read type from
 * SourceRule + Point. Result type and ordered transforms from validated Mapping only. */
export interface SourceFacts {
  binding: RevisionBinding; mapping_revision: string | null; enabled: boolean
  address: string | null; protocol_function: string | null
  read_type: EvidenceDataType | null; result_type: EvidenceDataType | null
  polling: PollingEvidence | null; transforms: SafeTransformStep[]
}
/** Point.PollingGroupID -> PollingGroup.IntervalMs, never the recording interval. */
export interface PollingEvidence { group_id: string; interval_ms: number; source: 'polling-group' }
/** Exactly one typed step; irrelevant parameters must be null, unsupported has no params. */
export interface SafeTransformStep {
  kind: 'byte-order' | 'scale' | 'cast' | 'unsupported'
  order: 'ABCD' | 'BADC' | 'CDAB' | 'DCBA' | null
  multiplier: number | null; offset: number | null; target_type: EvidenceDataType | null
}
/** WriteGroup saved; ResolveAppliedAt + verified frozen RuntimeLayout effective. */
export interface RecordingGroupView {
  /** applied_revision is latest requested, possibly scheduled; never current authority. */
  /** Requested boundary is separate from effective.effective_at, null unless proven. */
  effective_revision: string | null; as_of: string | null; applied_effective_at: string | null
  group_id: string; saved_revision: string | null; applied_revision: string | null
  state: EvidenceState; proof: 'unavailable' | 'applied-snapshot-and-layout'
  saved: RecordingGroupFacts | null; effective: RecordingGroupFacts | null
  delivery: DeliveryEvidence | null
}
/** Exact snapshot profile; unit is interval unit, never the physical measurement unit. */
export interface RecordingProfile { kind: 'instantaneous' | 'cumulative'; interval: number; unit: 'second' | 'minute' }
/** Values from that exact WriteGroup snapshot; effective_at from applied snapshot; nullable
 * interval remains unknown. No legacy layout/policy defaults are synthesized by this parser. */
export interface RecordingGroupFacts {
  recording_profile: RecordingProfile | null
  /** Exact snapshot RowPolicy, not an invented scheduler delay. */
  allowed_lateness_seconds: number | null
  /** Historical ended/disabled snapshots never grant intake permission. */
  ends_at: string | null
  binding: RevisionBinding; effective_at: string | null; enabled: boolean
  destination: SafeDestination; members: RecordingMember[]
  layout: 'managed_wide_v1' | 'managed_long_v1' | 'managed_wide_profile_v1' | 'unavailable'
  storage_mode: 'plain' | 'timescale' | 'native_monthly' | 'unavailable'
  interval_seconds: number | null; reduction: 'latest' | 'mean' | 'cumulative_delta' | 'unavailable'
  write_mode: 'append' | 'latest' | 'unavailable'; incomplete_policy: 'skip_row' | 'partial' | 'unavailable'
}
/** Frozen connector ID/revision, validated display alias/kind, schema/table; no database path. */
export interface SafeDestination {
  connector_id: string; connector_revision: string; alias: string | null
  kind: 'sqlite' | 'postgres' | 'mysql' | 'sqlserver' | 'unavailable'
  schema: string | null; table: string
}
/** WriteGroupMember identity and target_column, never arbitrary metadata/values.
 * Group relevance is membership in either saved or applied, without merging the two. */
/** circuit_id is persisted in this exact WriteGroupMember snapshot; missing legacy evidence
 * is null, never inferred from a column, display name or newer saved membership. */
/** source_value_type is persisted/verified frozen evidence, never a current Tag lookup. */
export interface RecordingMember { source_value_type: EvidenceDataType | null; circuit_id: string | null; device_id: string; point_id: string; tag_id: string; column: string }
/** DeliveryForRevision is the ONLY runtime delivery authority. Operation ledger stays separate. */
export interface DeliveryEvidence {
  counter_replay_quarantines: CounterReplayQuarantine[] | null
  /** DeliveryView.OldestPendingSeconds includes every group revision, not only current. */
  oldest_pending_scope: 'group-all-revisions'; oldest_pending_seconds: number | null
  current: CurrentDelivery | null; backlog: HistoricalBacklog[]; test_write: TestWriteEvidence | null
}
/** Requires Pipeline.Status GroupID/Revision match. Aggregate DeliveryView.Intake may prefer
 * another revision and cannot populate this. Missing revision proof means null, no invented active. */
export interface IntakeEvidence {
  group_revision: string
  state: 'active' | 'retiring' | 'blocked' | 'not_running' | 'unavailable'
  reason: EvidenceReason | null
}
/** Group totals (shared if multiple devices), not per-Tag successes. Sources remain independent.
 * Nullable timing: actual scheduler/outbox/attempt evidence only; no invented countdown. */
export interface CurrentDelivery {
  intake: IntakeEvidence | null
  binding: RevisionBinding; totals_scope: 'group'; shared: boolean; stages: DeliveryStageCounts
  last_commit: RuntimeReceipt | null; sources: PerSourceEvidence[]; source_state: 'available' | 'unavailable'
  next_due_at: string | null
  next_retry_at: string | null; attempt_duration_ms: number | null
  quota: QuotaEvidence | null; storage: StorageEvidence | null
}
/** Durable revision stage counts, never transport ACK/collector success. */
export interface DeliveryStageCounts {
  collecting: number; queued: number; retrying: number; blocked: number
  quarantined: number; unknown: number; sql_committed: number; skipped: number
}
/** Matching LastCommittedEffect / SourceRecord, not test-write/readback/cleanup. */
export interface RuntimeReceipt {
  group_revision: string; connector_revision: string; record_id: string
  effect_key: string; payload_digest: string; committed_at: string; provenance: 'runtime'
}
/** Separate latest local window from the last receipt-confirmed window. */
export interface PerSourceEvidence {
  device_id: string; point_id: string; tag_id: string
  latest: SourceResult | null; last_committed: SourceResult | null
}
/** SourceRecord verified window/coverage/delta-known, with no raw values. Unknown is not zero. */
export interface SourceResult {
  state: 'available' | 'missing' | 'unknown' | 'partial' | 'unavailable'
  window_start: string | null; window_end: string | null; delta_known: boolean | null
  coverage: 'complete' | 'partial' | 'unknown' | null; receipt: RuntimeReceipt | null
}
/** DeliveryView.Backlog retains its own frozen revision and destination. */
export interface HistoricalBacklog {
  binding: RevisionBinding; destination: SafeDestination; pending: number; reasons: EvidenceReason[]
}
/** Independent operation ledger provenance; unavailable until real test-write has evidence.
 * Cleanup can succeed after failed readback; neither timestamp implies the other. */
export interface TestWriteEvidence {
  operation_id: string; provenance: 'test-write'; written_at: string | null
  readback_at: string | null; cleanup_at: string | null
}
/** DeliveryView.Quota scope is global OR group (not revision-specific), no raw diagnostics. */
export interface QuotaEvidence {
  scope: 'global' | 'group'
  state: 'unconfigured' | 'ok' | 'warning' | 'hard_limit' | 'unavailable'
  used_bytes: number; max_bytes: number; intake_refused: boolean
}
/** Null unless adapter proves frozen destination/revision scope against storage lineage.
 * Editable-group DeliveryView.Storage alone is insufficient. No repair authority/raw payload. */
export interface StorageEvidence {
  generation: number; history_may_be_lost: boolean; reconstructions: number
  last_reconstructed_at: string | null; repair_state: 'unavailable' | 'none' | 'blocked'
  events: StorageEvidenceEvent[]
}
export interface StorageEvidenceEvent {
  kind: 'confirmed' | 'recreated_by_confirmation' | 'reconstructed' | 'repair_blocked' | 'unavailable'
  at: string | null
}

/** Exact backend accepted-input scope; not SQL rows, backlog or permission to repair. Null list
 * means unavailable, [] means observed empty; unsafe integer counts are rejected, not rounded. */
export interface CounterReplayQuarantine {
  workspace_id: string; group_id: string; group_revision: string; pending_samples: number
  reason: 'counter-replay-proof-missing' | 'counter-replay-proof-conflict' | 'counter-replay-proof-unsupported'
  action: 'preserve-inputs-and-restore-original-revision-evidence-or-contact-maintenance'
}
