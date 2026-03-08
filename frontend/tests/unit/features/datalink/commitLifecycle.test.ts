import { describe, expect, it } from 'vitest';
import {
  canExecuteCommit,
  executeCommitLifecycle,
  retryFailedLifecycle,
  rollbackCommitLifecycle,
  summarizeCommitImpact,
  type CommitQueueLifecycleItem,
} from '@/features/datalink/commitLifecycle';

const queueItems: CommitQueueLifecycleItem[] = [
  { id: 'q1', status: 'pending' },
  { id: 'q2', status: 'linked' },
  { id: 'q3', status: 'conflict' },
];

describe('commitLifecycle', () => {
  it('validates commit gates before execution', () => {
    expect(
      canExecuteCommit({ canActivate: false, mappingEnabled: true })
    ).toEqual({ ok: false, reason: 'not_validated' });
    expect(
      canExecuteCommit({ canActivate: true, mappingEnabled: false })
    ).toEqual({ ok: false, reason: 'mapping_disabled' });
    expect(
      canExecuteCommit({ canActivate: true, mappingEnabled: true })
    ).toEqual({ ok: true, reason: 'ok' });
  });

  it('executes commit and returns partial failure results', () => {
    const result = executeCommitLifecycle(queueItems, {});
    expect(result.successCount).toBe(2);
    expect(result.failedCount).toBe(1);
    expect(result.nextStatus).toEqual({
      q1: 'success',
      q2: 'success',
      q3: 'failed',
    });
  });

  it('retries failed items and only recovers non-conflict entries', () => {
    const retryInputItems: CommitQueueLifecycleItem[] = [
      { id: 'q3', status: 'conflict' },
      { id: 'q4', status: 'pending' },
    ];
    const current = { q3: 'failed', q4: 'failed' } as const;
    const result = retryFailedLifecycle(retryInputItems, current);

    expect(result.recovered).toBe(1);
    expect(result.remainingFailed).toBe(1);
    expect(result.nextStatus).toEqual({
      q3: 'failed',
      q4: 'success',
    });
  });

  it('rolls back queue run status to previous snapshot', () => {
    const snapshot = { q1: 'success', q3: 'failed' } as const;
    expect(rollbackCommitLifecycle(snapshot)).toEqual(snapshot);
    expect(rollbackCommitLifecycle(null)).toBeNull();
  });

  it('summarizes commit impact for queue and pending global tag edit', () => {
    const impact = summarizeCommitImpact(
      [
        { viewStatus: 'pending' },
        { viewStatus: 'linked' },
        { viewStatus: 'conflict' },
        { viewStatus: 'failed' },
      ],
      true
    );

    expect(impact).toEqual({
      newPoints: 2,
      globalTagUpdates: 1,
      conflicts: 2,
    });
  });
});
