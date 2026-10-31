/** Where one write group's accepted rows stand. Only `sql_committed` is destination evidence. */
export type WriteGroupIntakeState = 'active' | 'retiring' | 'blocked' | 'not_running';

export interface WriteGroupDeliveryStages {
  /** Accepted samples whose bucket has not closed yet. */
  collecting: number;
  /** Closed rows durably waiting for an attempt, including one being sent. */
  queued: number;
  retrying: number;
  blocked: number;
  quarantined: number;
  unknown: number;
  /** Rows the destination confirmed. Buffered or acknowledged data never counts here. */
  sql_committed: number;
  skipped: number;
}

export interface WriteGroupDeliveryBacklog {
  group_revision: string;
  connector_id: string;
  /** The destination revision the rows were accepted for; never retargeted. */
  connector_revision: string;
  table_schema: string;
  table_name: string;
  pending: number;
  error_codes: string[];
}

export type WriteGroupQuotaState = 'unconfigured' | 'ok' | 'warning' | 'hard_limit';

export interface WriteGroupDeliveryQuota {
  configured: boolean;
  state: WriteGroupQuotaState;
  scope: 'global' | 'group' | '';
  used_bytes: number;
  max_bytes: number;
  intake_refused: boolean;
  loss_risk_notice?: string;
}

export type WriteGroupStorageEventKind = 'confirmed' | 'recreated_by_confirmation' | 'reconstructed' | 'repair_blocked';

/** One value-free diagnostic about a managed long table; it never carries SQL, DSNs or driver text. */
export interface WriteGroupStorageEvent {
  generation: number;
  group_id?: string;
  kind: WriteGroupStorageEventKind;
  code?: string;
  reason?: string;
  operation_id?: string;
  occurrences: number;
  occurred_at: string;
}

/**
 * Which generation of the shared long table a group writes to, and whether a
 * rebuilt table means earlier history may be gone. A rebuilt structure is never
 * evidence that history was restored.
 */
export interface WriteGroupStorage {
  generation: number;
  history_may_be_lost: boolean;
  reconstructions: number;
  last_reconstructed_at?: string;
  /** Present only while an automatic restore ended in a result nobody can vouch for. */
  repair_state?: 'blocked';
  repair_reason?: string;
  /** Absent when no decision was ever made; never read as on. */
  recreate_missing?: boolean;
  events: WriteGroupStorageEvent[];
}

export interface WriteGroupDelivery {
  attention?: DeliveryAttentionItem[];
  group_id: string;
  intake: { state: WriteGroupIntakeState; reason?: string };
  stages: WriteGroupDeliveryStages;
  /** Latest destination-confirmed commit; null when nothing was ever confirmed. */
  last_sql_committed_at: string | null;
  /** Receipt evidence for the current applied revision; never a test-write or readback claim. */
  last_sql_committed_effect?: WriteGroupCommittedEffect | null;
  /** Accepted samples and rows belonging to the current applied revision only. */
  revision_stages?: { group_revision: string; stages: WriteGroupDeliveryStages } | null;
  oldest_pending_seconds: number;
  no_data_buckets: number;
  skipped_buckets: number;
  recent_bucket_issues?: WriteGroupDeliveryBucketIssue[];
  backlog: WriteGroupDeliveryBacklog[];
  quota: WriteGroupDeliveryQuota;
  /** Only for groups on managed long storage. */
  storage?: WriteGroupStorage;
  /** Per applied source result of managed long groups; null or absent for every other layout. */
  source_evidence?: WriteGroupSourceEvidence | null;
}

/** What one long record stored for its source. Zero and false are values; only NULL is absent. */
export type WriteGroupSourceStatus = 'ok' | 'partial' | 'baseline' | 'missing' | 'bad' | 'stale' | 'invalid' | 'reset_or_uncertain' | 'gap';
/** Local delivery stage of a record. Only `sql_committed` with a matching receipt is destination evidence. */
export type WriteGroupSourceDelivery = 'queued' | 'retrying' | 'blocked' | 'quarantined' | 'unknown' | 'sql_committed' | 'operator_skipped' | 'window_skipped';
export type WriteGroupSourceUnavailableCause = 'payload_invalid' | 'digest_mismatch' | 'identity_mismatch' | 'status_unknown' | 'window_invalid' | 'calculation_projection_unsupported';

/**
 * Value-free evidence of a mean window, verified against what the applied layout froze for the source: the method, the
 * actual counts and the coverage. Coverage is never complete (no stored proof can back it) and never a percentage.
 */
export interface WriteGroupSourceAggregation {
  method: 'arithmetic-mean';
  version: string;
  basis: 'exact' | 'float64';
  /** The Tag's own type and the type the mean is stored as are separate facts. */
  source_type: string;
  result_type: string;
  valid_count: number;
  rejected_count: number;
  coverage: 'partial' | 'unknown' | 'missing';
  /** A stable typed code, or `reason_unavailable` when the stored reason is not a known code. */
  reason?: string;
}

/** A value-free, verified result read from the frozen outbox payload of one source and window. */
export interface WriteGroupLegacySourceRecord {
  calculation?: never;
  available: true;
  group_revision: string;
  /** The destination revision the record was frozen to; never retargeted. */
  connector_revision: string;
  record_id: string;
  effect_key: string;
  payload_digest: string;
  profile_kind: 'instantaneous' | 'cumulative';
  window_start: string;
  window_end: string;
  status: WriteGroupSourceStatus;
  /** A stable typed code, or `reason_unavailable` when the stored reason is not a known code. */
  reason?: string;
  raw_present: boolean;
  /** Present on the window of a mean source only. */
  aggregation?: WriteGroupSourceAggregation;
  /** Cumulative records only: whether a trustworthy difference was stored. Unknown is not zero. */
  delta_known?: boolean;
  delivery: WriteGroupSourceDelivery;
  /** Present only when the destination receipt matched this payload. */
  committed_at?: string;
}


/** Value-free calculated result; raw presence is provenance, never source health. */
export interface WriteGroupSourceCalculation {
  result_shape: 'member-calculation-v1'; method: 'latest' | 'mean' | 'delta'; interval_seconds: number; alignment: 'utc_epoch';
  due_at: string; next_due?: string; outcome_digest: string; outcome_kind: 'row' | 'skipped' | 'no_data';
  status: 'ok' | 'missing' | 'bad' | 'stale' | 'invalid' | 'not_due'; reason?: string;
  result_present: boolean; row_emitted: boolean; raw_present: boolean; source_type: string; result_type: string;
  aligned_start?: string; aligned_end?: string; actual_start?: string; actual_end?: string; rounded?: boolean;
}
export interface WriteGroupMixedSourceRecord extends Omit<WriteGroupLegacySourceRecord, 'profile_kind' | 'window_start' | 'window_end' | 'status' | 'calculation'> {
  profile_kind: 'mixed';
  /** null means not due; the Go zero-time wire sentinel is never a data window. */
  window_start: string | null; window_end: string | null;
  status: WriteGroupSourceStatus | 'not_due';
  calculation: WriteGroupSourceCalculation;
}
export type WriteGroupSourceRecord = WriteGroupLegacySourceRecord | WriteGroupMixedSourceRecord;

/** A record whose frozen evidence could not be verified; it never counts as written. */
export interface WriteGroupSourceRecordUnavailable {
  available: false;
  cause: WriteGroupSourceUnavailableCause;
  delivery: WriteGroupSourceDelivery;
}

export interface WriteGroupSourceEvidenceItem {
  circuit_id?: string;
  device_id: string;
  point_id: string;
  tag_id: string;
  latest: WriteGroupSourceRecord | WriteGroupSourceRecordUnavailable | null;
  last_sql_committed: WriteGroupSourceRecord | WriteGroupSourceRecordUnavailable | null;
}

export type WriteGroupSourceEvidenceUnavailableReason = 'not_long_layout' | 'applied_snapshot_unavailable' | 'revision_mismatch';

export interface WriteGroupSourceEvidence {
  state: 'available' | 'unavailable';
  unavailable_reason?: WriteGroupSourceEvidenceUnavailableReason;
  /** The applied revision every record below belongs to. */
  group_revision: string;
  /** The destination the applied revision froze; every record was written there, whatever the draft says now. */
  connector_id?: string;
  connector_revision?: string;
  sources: WriteGroupSourceEvidenceItem[];
}

export interface WriteGroupCommittedEffect {
  group_revision: string;
  connector_revision: string;
  record_id: string;
  effect_key: string;
  payload_digest: string;
  committed_at: string;
}

/** `truncated` is a row that was written for a window a cutover cut short. */
export type WriteGroupDeliveryBucketIssueKind = 'skipped' | 'no_data' | 'truncated';
export type WriteGroupDeliveryBucketCause = 'missing' | 'bad' | 'stale' | 'invalid' | 'no_data' | 'unavailable' | 'boundary';

export interface WriteGroupDeliveryBucketIssue {
  group_revision: string;
  bucket_start: string;
  kind: WriteGroupDeliveryBucketIssueKind;
  causes: WriteGroupDeliveryBucketCause[];
}

export interface DeliveryAttentionItem {
  effect_key: string; state: 'blocked' | 'quarantined' | 'unknown'; state_revision: number; payload_digest: string;
  group_revision: string; connector_id: string; connector_revision: string;
  table_schema: string; table_name: string; error_code: string;
}
export interface DeliveryDecision {
  decision_id: string; effect_key: string; expected_state: string; expected_state_revision: number; payload_digest: string;
  resolution: 'retry' | 'skip'; reason: string; confirm_skip: boolean;
}
export interface DeliveryDecisionResult { decision_id: string; effect_key: string; state: string; duplicate: boolean }
