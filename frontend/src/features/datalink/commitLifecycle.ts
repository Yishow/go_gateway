export type QueueBaseStatus = 'pending' | 'linked' | 'conflict';

export interface CommitQueueLifecycleItem {
  id: string;
  status: QueueBaseStatus;
}

export type CommitRunStatusMap = Record<string, 'success' | 'failed'>;
export type CommitQueueViewStatus = QueueBaseStatus | 'committed' | 'failed';

export interface CommitImpactSummaryItem {
  viewStatus: CommitQueueViewStatus;
}

export function summarizeCommitImpact(
  items: CommitImpactSummaryItem[],
  pendingGlobalTagEdit: boolean
): {
  newPoints: number;
  globalTagUpdates: number;
  conflicts: number;
} {
  const newPoints = items.filter((item) => item.viewStatus !== 'conflict' && item.viewStatus !== 'failed').length;
  const conflicts = items.filter((item) => item.viewStatus === 'conflict' || item.viewStatus === 'failed').length;
  return {
    newPoints,
    globalTagUpdates: pendingGlobalTagEdit ? 1 : 0,
    conflicts,
  };
}

export function canExecuteCommit(input: { canActivate: boolean; mappingEnabled: boolean }): {
  ok: boolean;
  reason: 'ok' | 'not_validated' | 'mapping_disabled';
} {
  if (!input.canActivate) return { ok: false, reason: 'not_validated' };
  if (!input.mappingEnabled) return { ok: false, reason: 'mapping_disabled' };
  return { ok: true, reason: 'ok' };
}

export function executeCommitLifecycle(
  items: CommitQueueLifecycleItem[],
  previousStatus: CommitRunStatusMap
): {
  nextStatus: CommitRunStatusMap;
  successCount: number;
  failedCount: number;
} {
  const nextStatus: CommitRunStatusMap = { ...previousStatus };
  let successCount = 0;
  let failedCount = 0;

  items.forEach((item) => {
    if (item.status === 'conflict') {
      nextStatus[item.id] = 'failed';
      failedCount += 1;
      return;
    }
    nextStatus[item.id] = 'success';
    successCount += 1;
  });

  return { nextStatus, successCount, failedCount };
}

export function retryFailedLifecycle(
  items: CommitQueueLifecycleItem[],
  currentStatus: CommitRunStatusMap
): {
  nextStatus: CommitRunStatusMap;
  recovered: number;
  remainingFailed: number;
} {
  const nextStatus: CommitRunStatusMap = { ...currentStatus };
  let recovered = 0;
  let remainingFailed = 0;

  items.forEach((item) => {
    if (currentStatus[item.id] !== 'failed') return;
    if (item.status === 'conflict') {
      nextStatus[item.id] = 'failed';
      remainingFailed += 1;
      return;
    }
    nextStatus[item.id] = 'success';
    recovered += 1;
  });

  return { nextStatus, recovered, remainingFailed };
}

export function rollbackCommitLifecycle(snapshot: CommitRunStatusMap | null): CommitRunStatusMap | null {
  if (!snapshot) return null;
  return { ...snapshot };
}
