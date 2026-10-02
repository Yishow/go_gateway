/**
 * Explicit, operator-confirmed test write of one saved write group. The write
 * result and the cleanup of the test row are independent facts: a verified
 * write can still leave its row behind, and neither is ever inferred.
 */
export type WriteGroupTestWriteOutcome = 'written_verified' | 'written_unverified' | 'failed' | 'unknown';
export type WriteGroupTestCleanupStatus = 'not_attempted' | 'cleaned' | 'failed' | 'unknown';
export type WriteGroupTestWriteStatus = 'pending' | 'running' | 'succeeded' | 'partial' | 'failed' | 'unknown';

export interface WriteGroupTestWriteValue {
  column: string;
  type: string;
  value: string;
}

export interface WriteGroupTestWritePreview {
  token: string;
  operation_id: string;
  action: 'test_write';
  group_id: string;
  group_revision: string;
  expires_at: string;
  target: { connector_id: string; dialect: string; database: string; schema: string; table: string };
  /** The column whose operation-owned value identifies the only row cleanup may remove. */
  owner_column: string;
  owner_value: string;
  dedupe: string;
  values: WriteGroupTestWriteValue[];
  cleanup: string;
}

export interface WriteGroupTestWriteConfirmation {
  token: string;
  operation_id: string;
}

export interface WriteGroupTestWriteOperation {
  operation_id: string;
  action: 'test_write';
  status: WriteGroupTestWriteStatus;
  /** Present once the operation ended; absent while it is still running. */
  write_outcome?: WriteGroupTestWriteOutcome;
  cleanup_status?: WriteGroupTestCleanupStatus;
  /** Safe reason codes only, never values or connection details. */
  reason?: string;
  cleanup_reason?: string;
  created_at: string;
  updated_at: string;
  completed_at?: string;
}
