import { describe, expect, it } from 'vitest';
import { parseRecordingSchemaOperation } from '@/utils/recordingPlanJson';

const operation = {
  operation_id: 'op-1',
  action: 'schema_apply',
  status: 'succeeded',
  executed_statements: 9,
  verified_digest: 'a'.repeat(64),
  created_at: '2026-09-16T00:50:00Z',
  updated_at: '2026-09-16T00:51:00Z',
  completed_at: '2026-09-16T00:51:00Z',
};

describe('parseRecordingSchemaOperation', () => {
  it('keeps the recorded outcome of a finished operation', () => {
    expect(parseRecordingSchemaOperation(operation)).toEqual({
      operation_id: 'op-1',
      action: 'schema_apply',
      status: 'succeeded',
      executed_statements: 9,
      verified_digest: 'a'.repeat(64),
      created_at: '2026-09-16T00:50:00Z',
      updated_at: '2026-09-16T00:51:00Z',
      completed_at: '2026-09-16T00:51:00Z',
    });
  });

  it('keeps partial and unknown outcomes with their reason and next step', () => {
    // An unverified outcome carries no digest at all, so the key is absent.
    const { verified_digest: _digest, ...unverified } = operation;
    for (const status of ['partial', 'unknown'] as const) {
      const parsed = parseRecordingSchemaOperation({
        ...unverified, status, executed_statements: 3,
        reason: 'statement_failed_partially_applied', next_action: 'preview again to create the missing tables',
      });
      expect(parsed?.status).toBe(status);
      expect(parsed?.executed_statements).toBe(3);
      expect(parsed?.reason).toBe('statement_failed_partially_applied');
      expect(parsed?.next_action).toBe('preview again to create the missing tables');
      expect(parsed?.verified_digest).toBeUndefined();
    }
  });

  it('keeps a running operation that has not finished yet', () => {
    const parsed = parseRecordingSchemaOperation({
      operation_id: 'op-2', status: 'running', executed_statements: 0,
      created_at: '2026-09-16T00:50:00Z', updated_at: '2026-09-16T00:50:00Z',
    });
    expect(parsed?.status).toBe('running');
    expect(parsed?.completed_at).toBeUndefined();
  });

  it('rejects an unknown status instead of guessing', () => {
    expect(parseRecordingSchemaOperation({ ...operation, status: 'done' })).toBeNull();
    expect(parseRecordingSchemaOperation({ ...operation, status: '' })).toBeNull();
  });

  it('rejects a result without its operation identity or counts', () => {
    expect(parseRecordingSchemaOperation({ ...operation, operation_id: '' })).toBeNull();
    expect(parseRecordingSchemaOperation({ ...operation, executed_statements: '9' })).toBeNull();
    expect(parseRecordingSchemaOperation({ ...operation, executed_statements: -1 })).toBeNull();
    expect(parseRecordingSchemaOperation({ ...operation, created_at: '' })).toBeNull();
    expect(parseRecordingSchemaOperation(null)).toBeNull();
  });
});
