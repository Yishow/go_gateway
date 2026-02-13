import type { CommitImpactSummaryItem, CommitQueueViewStatus } from './commitLifecycle';

export interface CommitAuditChunkResult {
  chunk: number;
  totalChunks: number;
  success: number;
  failed: number;
  status: 'success' | 'failed';
}

export interface CommitAuditQueueItem extends CommitImpactSummaryItem {
  id: string;
  label: string;
  type: string;
  addresses: string[];
}

export interface CommitAuditPayload {
  action: 'commit';
  createdAt: string;
  traceLink: '#commit-audit-trace';
  summary: {
    total: number;
    newPoints: number;
    globalTagUpdates: number;
    conflicts: number;
    committed: number;
    failed: number;
  };
  queue: Array<{
    id: string;
    label: string;
    type: string;
    startAddress: string;
    endAddress: string;
    viewStatus: CommitQueueViewStatus;
  }>;
  chunks: CommitAuditChunkResult[];
}

export function buildCommitAuditPayload(input: {
  queueItems: CommitAuditQueueItem[];
  chunkResults: CommitAuditChunkResult[];
  impact: {
    newPoints: number;
    globalTagUpdates: number;
    conflicts: number;
  };
}): CommitAuditPayload {
  const committed = input.queueItems.filter((item) => item.viewStatus === 'committed').length;
  const failed = input.queueItems.filter((item) => item.viewStatus === 'failed').length;
  return {
    action: 'commit',
    createdAt: new Date().toISOString(),
    traceLink: '#commit-audit-trace',
    summary: {
      total: input.queueItems.length,
      newPoints: input.impact.newPoints,
      globalTagUpdates: input.impact.globalTagUpdates,
      conflicts: input.impact.conflicts,
      committed,
      failed,
    },
    queue: input.queueItems.map((item) => ({
      id: item.id,
      label: item.label,
      type: item.type,
      startAddress: item.addresses[0] || '',
      endAddress: item.addresses[item.addresses.length - 1] || '',
      viewStatus: item.viewStatus,
    })),
    chunks: input.chunkResults,
  };
}
