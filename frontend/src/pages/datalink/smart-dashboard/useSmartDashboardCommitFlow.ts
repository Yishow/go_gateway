import { useCallback, useEffect, useMemo, useState } from 'react';
import type { TFunction } from 'i18next';
import type { PlannedAllocation } from '../../../components/datalink/MemoryGrid';
import type { FlowSegment } from '../../../features/flow/stateMachine';
import {
  canExecuteCommit,
  executeCommitLifecycle,
  retryFailedLifecycle,
  rollbackCommitLifecycle,
  summarizeCommitImpact,
  type CommitQueueViewStatus,
  type QueueBaseStatus,
} from '../../../features/datalink/commitLifecycle';
import { buildCommitAuditPayload, type CommitAuditPayload } from '../../../features/datalink/commitAudit';
import { estimatePollingLoadDelta } from '../../../features/datalink/pollingLoadEstimate';
import { MOTION_TOKENS, buildMotionReadabilityGate } from '../../../features/datalink/motionGuidance';
import type { Point, PollingGroup } from '../../../types/datalink';

const DEFAULT_COMMIT_CHUNK_SIZE = 8;

function isCommitStatusMapEqual(
  left: Record<string, 'success' | 'failed'>,
  right: Record<string, 'success' | 'failed'>,
): boolean {
  const leftKeys = Object.keys(left);
  const rightKeys = Object.keys(right);
  if (leftKeys.length !== rightKeys.length) return false;
  return leftKeys.every((key) => left[key] === right[key]);
}

interface UseSmartDashboardCommitFlowInput {
  allPoints: Point[];
  linkedAddresses: string[];
  plannedAllocations: PlannedAllocation[];
  pendingGlobalTagEdit: boolean;
  pollingGroups: PollingGroup[];
  canActivate: boolean;
  mappingEnabled: boolean;
  markActive: () => void;
  markError: (segment: FlowSegment, message: string) => void;
  setGuideStage: (stage: 'idle' | 'grid' | 'commit') => void;
  scheduleGuideStageReset: () => void;
  t: TFunction;
  chunkSize?: number;
}

interface CommitQueueBaseItem {
  id: string;
  label: string;
  type: string;
  addresses: string[];
  status: QueueBaseStatus;
  order: number;
}

interface CommitQueueItem extends CommitQueueBaseItem {
  viewStatus: CommitQueueViewStatus;
}

interface CommitQueueSummary {
  total: number;
  pending: number;
  linked: number;
  conflict: number;
  committed: number;
  failed: number;
}

export function useSmartDashboardCommitFlow({
  allPoints,
  linkedAddresses,
  plannedAllocations,
  pendingGlobalTagEdit,
  pollingGroups,
  canActivate,
  mappingEnabled,
  markActive,
  markError,
  setGuideStage,
  scheduleGuideStageReset,
  t,
  chunkSize = DEFAULT_COMMIT_CHUNK_SIZE,
}: UseSmartDashboardCommitFlowInput) {
  const [commitActionMessage, setCommitActionMessage] = useState('');
  const [commitQueueRunStatus, setCommitQueueRunStatus] = useState<Record<string, 'success' | 'failed'>>({});
  const [lastCommitSnapshot, setLastCommitSnapshot] = useState<Record<string, 'success' | 'failed'> | null>(null);
  const [isCommitRunning, setIsCommitRunning] = useState(false);
  const [commitChunkResults, setCommitChunkResults] = useState<
    Array<{ chunk: number; totalChunks: number; success: number; failed: number; status: 'success' | 'failed' }>
  >([]);
  const [commitAuditPayload, setCommitAuditPayload] = useState<CommitAuditPayload | null>(null);
  const [failedChunkRetryQueue, setFailedChunkRetryQueue] = useState<number[]>([]);

  const baseCommitQueueItems = useMemo<CommitQueueBaseItem[]>(() => {
    const usedSet = new Set(allPoints.map((point) => point.address));
    const linkedSet = new Set(linkedAddresses);
    return plannedAllocations.map((allocation, index) => {
      const conflictCount = allocation.addresses.filter((address) => usedSet.has(address)).length;
      const linkedCount = allocation.addresses.filter((address) => linkedSet.has(address)).length;
      const status: QueueBaseStatus = conflictCount > 0 ? 'conflict' : linkedCount > 0 ? 'linked' : 'pending';
      return {
        id: allocation.id,
        label: allocation.label,
        type: allocation.dataType,
        addresses: allocation.addresses,
        status,
        order: index + 1,
      };
    });
  }, [allPoints, linkedAddresses, plannedAllocations]);

  const commitQueueItems = useMemo<CommitQueueItem[]>(() => {
    return baseCommitQueueItems.map((item) => {
      const runStatus = commitQueueRunStatus[item.id];
      const viewStatus: CommitQueueViewStatus =
        runStatus === 'success' ? 'committed' : runStatus === 'failed' ? 'failed' : item.status;
      return {
        ...item,
        viewStatus,
      };
    });
  }, [baseCommitQueueItems, commitQueueRunStatus]);

  useEffect(() => {
    setCommitQueueRunStatus((prev) => {
      const next: Record<string, 'success' | 'failed'> = {};
      baseCommitQueueItems.forEach((item) => {
        if (prev[item.id]) next[item.id] = prev[item.id];
      });
      return isCommitStatusMapEqual(prev, next) ? prev : next;
    });
  }, [baseCommitQueueItems]);

  const commitQueueSummary = useMemo<CommitQueueSummary>(() => {
    return commitQueueItems.reduce(
      (acc, item) => {
        acc.total += 1;
        if (item.viewStatus === 'conflict') acc.conflict += 1;
        if (item.viewStatus === 'linked') acc.linked += 1;
        if (item.viewStatus === 'pending') acc.pending += 1;
        if (item.viewStatus === 'failed') acc.failed += 1;
        if (item.viewStatus === 'committed') acc.committed += 1;
        return acc;
      },
      { total: 0, pending: 0, linked: 0, conflict: 0, failed: 0, committed: 0 },
    );
  }, [commitQueueItems]);

  const commitImpactSummary = useMemo(
    () => summarizeCommitImpact(commitQueueItems, pendingGlobalTagEdit),
    [commitQueueItems, pendingGlobalTagEdit],
  );
  const preCommitLoadEstimate = useMemo(
    () => estimatePollingLoadDelta(allPoints, pollingGroups, commitImpactSummary.newPoints),
    [allPoints, commitImpactSummary.newPoints, pollingGroups],
  );
  const motionQAGate = useMemo(
    () =>
      buildMotionReadabilityGate({
        stageHandoffMs: MOTION_TOKENS.stageHandoffMs,
        commitFeedbackMs: MOTION_TOKENS.commitFeedbackMs,
        hasReducedMotionFallback: true,
        intentOnlyAnimations: true,
      }),
    [],
  );

  const handleCommitFlow = useCallback(() => {
    const commitGate = canExecuteCommit({
      canActivate,
      mappingEnabled,
    });
    if (!commitGate.ok) {
      if (commitGate.reason === 'not_validated') {
        markError('sink', t('smartDashboard.flowErrors.notValidated'));
        setCommitActionMessage('Commit 失敗：請先完成 Validate。');
        return;
      }
      markError('sink', t('smartDashboard.flowErrors.mappingDisabled'));
      setCommitActionMessage('Commit 失敗：Mapping 尚未啟用。');
      return;
    }

    setIsCommitRunning(true);
    setLastCommitSnapshot(commitQueueRunStatus);
    setCommitChunkResults([]);
    setFailedChunkRetryQueue([]);

    const chunks = Array.from({ length: Math.ceil(baseCommitQueueItems.length / chunkSize) }, (_, index) =>
      baseCommitQueueItems.slice(index * chunkSize, (index + 1) * chunkSize),
    );
    let rollingStatus = { ...commitQueueRunStatus };
    const nextChunkResults: Array<{
      chunk: number;
      totalChunks: number;
      success: number;
      failed: number;
      status: 'success' | 'failed';
    }> = [];
    let successCount = 0;
    let failedCount = 0;

    chunks.forEach((chunkItems, index) => {
      const chunkResult = executeCommitLifecycle(chunkItems, rollingStatus);
      rollingStatus = chunkResult.nextStatus;
      successCount += chunkResult.successCount;
      failedCount += chunkResult.failedCount;
      nextChunkResults.push({
        chunk: index + 1,
        totalChunks: chunks.length,
        success: chunkResult.successCount,
        failed: chunkResult.failedCount,
        status: chunkResult.failedCount > 0 ? 'failed' : 'success',
      });
    });

    setCommitQueueRunStatus(rollingStatus);
    setCommitChunkResults(nextChunkResults);
    setFailedChunkRetryQueue(
      nextChunkResults.filter((chunkResult) => chunkResult.status === 'failed').map((chunkResult) => chunkResult.chunk - 1),
    );

    const nextQueueItems = baseCommitQueueItems.map((item) => {
      const runStatus = rollingStatus[item.id];
      const viewStatus: CommitQueueViewStatus =
        runStatus === 'success' ? 'committed' : runStatus === 'failed' ? 'failed' : item.status;
      return {
        ...item,
        viewStatus,
      };
    });
    const nextImpact = summarizeCommitImpact(nextQueueItems, pendingGlobalTagEdit);
    setCommitAuditPayload(
      buildCommitAuditPayload({
        queueItems: nextQueueItems,
        chunkResults: nextChunkResults,
        impact: nextImpact,
      }),
    );

    setIsCommitRunning(false);

    if (failedCount > 0) {
      markError('sink', `Commit 部分失敗：${failedCount} 筆失敗，請執行 Retry 或 Rollback。`);
      setCommitActionMessage(`Commit 部分成功：成功 ${successCount}、失敗 ${failedCount}。`);
      return;
    }

    markActive();
    setGuideStage('commit');
    scheduleGuideStageReset();
    setCommitActionMessage(`Commit 成功：${successCount} 筆已提交並啟用流程。`);
  }, [
    baseCommitQueueItems,
    canActivate,
    chunkSize,
    commitQueueRunStatus,
    mappingEnabled,
    markActive,
    markError,
    pendingGlobalTagEdit,
    scheduleGuideStageReset,
    setGuideStage,
    t,
  ]);

  const handleRetryFailedCommits = useCallback(() => {
    if (failedChunkRetryQueue.length === 0) {
      setCommitActionMessage('沒有可重試的失敗項目。');
      return;
    }

    const chunks = Array.from({ length: Math.ceil(baseCommitQueueItems.length / chunkSize) }, (_, index) =>
      baseCommitQueueItems.slice(index * chunkSize, (index + 1) * chunkSize),
    );

    let rollingStatus = { ...commitQueueRunStatus };
    let recovered = 0;
    let remainingFailed = 0;
    const nextFailedChunkQueue: number[] = [];

    failedChunkRetryQueue.forEach((chunkIndex) => {
      const chunkItems = chunks[chunkIndex] || [];
      const chunkRetry = retryFailedLifecycle(chunkItems, rollingStatus);
      rollingStatus = chunkRetry.nextStatus;
      recovered += chunkRetry.recovered;
      remainingFailed += chunkRetry.remainingFailed;
      if (chunkRetry.remainingFailed > 0) nextFailedChunkQueue.push(chunkIndex);
    });

    setCommitQueueRunStatus(rollingStatus);
    setFailedChunkRetryQueue(nextFailedChunkQueue);

    if (remainingFailed > 0) {
      setCommitActionMessage(`Retry 完成：恢復 ${recovered} 筆，仍有 ${remainingFailed} 筆衝突。`);
      return;
    }

    setCommitActionMessage(`Retry 成功：已恢復 ${recovered} 筆失敗項目。`);
  }, [baseCommitQueueItems, chunkSize, commitQueueRunStatus, failedChunkRetryQueue]);

  const handleRollbackCommitRun = useCallback(() => {
    const rolledBack = rollbackCommitLifecycle(lastCommitSnapshot);
    if (!rolledBack) {
      setCommitActionMessage('目前沒有可回滾的提交快照。');
      return;
    }
    setCommitQueueRunStatus(rolledBack);
    setFailedChunkRetryQueue([]);
    setLastCommitSnapshot(null);
    setCommitActionMessage('已回滾到上次 Commit 前的佇列狀態。');
  }, [lastCommitSnapshot]);

  return {
    commitQueueItems,
    commitQueueSummary,
    commitImpactSummary,
    preCommitLoadEstimate,
    motionQAGate,
    commitActionMessage,
    setCommitActionMessage,
    isCommitRunning,
    commitChunkResults,
    commitAuditPayload,
    hasFailedChunk: failedChunkRetryQueue.length > 0,
    canRollback: Boolean(lastCommitSnapshot),
    handleCommitFlow,
    handleRetryFailedCommits,
    handleRollbackCommitRun,
  };
}
