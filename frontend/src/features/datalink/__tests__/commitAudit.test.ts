import { describe, expect, it, vi } from 'vitest';
import { buildCommitAuditPayload } from '../commitAudit';

describe('commitAudit', () => {
  it('builds audit payload with summary and trace link', () => {
    vi.useFakeTimers();
    vi.setSystemTime(new Date('2026-02-13T10:30:00.000Z'));

    const payload = buildCommitAuditPayload({
      queueItems: [
        {
          id: 'q1',
          label: 'SRC_001',
          type: 'int16',
          addresses: ['40001'],
          viewStatus: 'committed',
        },
        {
          id: 'q2',
          label: 'SRC_002',
          type: 'float32',
          addresses: ['40002', '40003'],
          viewStatus: 'failed',
        },
      ],
      chunkResults: [{ chunk: 1, totalChunks: 1, success: 1, failed: 1, status: 'failed' }],
      impact: {
        newPoints: 1,
        globalTagUpdates: 1,
        conflicts: 1,
      },
    });

    expect(payload.createdAt).toBe('2026-02-13T10:30:00.000Z');
    expect(payload.traceLink).toBe('#commit-audit-trace');
    expect(payload.summary).toEqual({
      total: 2,
      newPoints: 1,
      globalTagUpdates: 1,
      conflicts: 1,
      committed: 1,
      failed: 1,
    });
    expect(payload.queue[1]).toEqual({
      id: 'q2',
      label: 'SRC_002',
      type: 'float32',
      startAddress: '40002',
      endAddress: '40003',
      viewStatus: 'failed',
    });

    vi.useRealTimers();
  });
});
