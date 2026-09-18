import { describe, expect, it } from 'vitest';
import { parseRecordingTestWriteResult } from '../../../src/utils/recordingPlanJson';

const fullUnknownResult = {
  status: 'unknown',
  record_id: 'rec-9',
  table: 'gw_record_samples',
  observed_at: '2026-09-18T09:00:00Z',
  delivered_at: '2026-09-18T09:00:01Z',
  operation_id: 'op-7',
  message: 'connection lost after dispatch',
};

describe('parseRecordingTestWriteResult', () => {
  it('accepts the unknown status and keeps the remaining fields', () => {
    expect(parseRecordingTestWriteResult(fullUnknownResult)).toEqual({
      status: 'unknown',
      record_id: 'rec-9',
      table: 'gw_record_samples',
      observed_at: '2026-09-18T09:00:00Z',
      delivered_at: '2026-09-18T09:00:01Z',
      operation_id: 'op-7',
      message: 'connection lost after dispatch',
    });
  });

  it('still rejects an unknown status whose required evidence is missing', () => {
    expect(
      parseRecordingTestWriteResult({ status: 'unknown', record_id: '' }),
    ).toBeNull();
  });
});
