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

export interface WriteGroupDelivery {
  group_id: string;
  intake: { state: WriteGroupIntakeState; reason?: string };
  stages: WriteGroupDeliveryStages;
  /** Latest destination-confirmed commit; null when nothing was ever confirmed. */
  last_sql_committed_at: string | null;
  oldest_pending_seconds: number;
  no_data_buckets: number;
  skipped_buckets: number;
  recent_bucket_issues?: WriteGroupDeliveryBucketIssue[];
  backlog: WriteGroupDeliveryBacklog[];
  quota: WriteGroupDeliveryQuota;
}

export type WriteGroupDeliveryBucketIssueKind = 'skipped' | 'no_data';
export type WriteGroupDeliveryBucketCause = 'missing' | 'bad' | 'stale' | 'invalid' | 'no_data' | 'unavailable';

export interface WriteGroupDeliveryBucketIssue {
  group_revision: string;
  bucket_start: string;
  kind: WriteGroupDeliveryBucketIssueKind;
  causes: WriteGroupDeliveryBucketCause[];
}
