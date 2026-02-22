interface SmartDashboardCommitPanelProps {
  commitQueueSummary: {
    total: number;
    pending: number;
    linked: number;
    conflict: number;
    committed: number;
    failed: number;
  };
  commitImpactSummary: {
    newPoints: number;
    globalTagUpdates: number;
    conflicts: number;
  };
  preCommitLoadEstimate: {
    deltaReadsPerSec: number;
    baselineReadsPerSec: number;
    projectedReadsPerSec: number;
    assumedIntervalMs: number;
  };
  motionQAGate: {
    pass: boolean;
    checklist: Array<{ pass: boolean }>;
  };
  commitQueueItems: Array<{
    id: string;
    order: number;
    label: string;
    viewStatus: 'pending' | 'linked' | 'conflict' | 'committed' | 'failed';
    type: string;
    addresses: string[];
  }>;
  onValidateFlow: () => void;
  canValidate: boolean;
  validating: boolean;
  onCommitFlow: () => void;
  canCommit: boolean;
  isCommitRunning: boolean;
  onRetryFailed: () => void;
  hasFailedChunk: boolean;
  onRollback: () => void;
  canRollback: boolean;
  segmentFeedback: Array<{ id: string; label: string; ok: boolean; message: string }>;
  commitActionMessage: string;
  commitAuditPayload: {
    traceId: string;
    traceLink: string;
    createdAt: string;
    summary: { newPoints: number; globalTagUpdates: number; conflicts: number };
  } | null;
  commitChunkResults: Array<{ chunk: number; totalChunks: number; success: number; failed: number; status: 'success' | 'failed' }>;
  hasError: boolean;
  onRecoverFlow: () => void;
  t: import('i18next').TFunction;
}

export default function SmartDashboardCommitPanel({
  commitQueueSummary,
  commitImpactSummary,
  preCommitLoadEstimate,
  motionQAGate,
  commitQueueItems,
  onValidateFlow,
  canValidate,
  validating,
  onCommitFlow,
  canCommit,
  isCommitRunning,
  onRetryFailed,
  hasFailedChunk,
  onRollback,
  canRollback,
  segmentFeedback,
  commitActionMessage,
  commitAuditPayload,
  commitChunkResults,
  hasError,
  onRecoverFlow,
  t,
}: SmartDashboardCommitPanelProps) {
  return (
    <>
      <div className="space-y-2 border-t border-white/5 p-4">
        <div className="flex items-center justify-between">
          <p className="text-xs font-semibold tracking-wide text-slate-200">Commit Queue</p>
          <span className="text-[10px] text-slate-400">Total {commitQueueSummary.total}</span>
        </div>
        <div className="flex flex-wrap gap-1.5 text-[10px]">
          <span className="rounded border border-slate-700 bg-slate-800/60 px-2 py-1 text-slate-200">
            Pending <strong>{commitQueueSummary.pending}</strong>
          </span>
          <span className="rounded border border-indigo-500/30 bg-indigo-500/10 px-2 py-1 text-indigo-100">
            Linked <strong>{commitQueueSummary.linked}</strong>
          </span>
          <span className="rounded border border-rose-500/30 bg-rose-500/10 px-2 py-1 text-rose-100">
            Conflict <strong>{commitQueueSummary.conflict}</strong>
          </span>
          <span className="rounded border border-emerald-500/30 bg-emerald-500/10 px-2 py-1 text-emerald-100">
            Committed <strong>{commitQueueSummary.committed}</strong>
          </span>
          <span className="rounded border border-amber-500/30 bg-amber-500/10 px-2 py-1 text-amber-100">
            Failed <strong>{commitQueueSummary.failed}</strong>
          </span>
        </div>
        <div className="space-y-1 rounded border border-cyan-500/30 bg-cyan-500/10 px-2 py-2 text-[10px] text-cyan-100">
          <p className="font-semibold tracking-wide">Commit Impact</p>
          <div className="grid grid-cols-3 gap-2">
            <div>New Points {commitImpactSummary.newPoints}</div>
            <div>Global Tag Updates {commitImpactSummary.globalTagUpdates}</div>
            <div>Conflicts {commitImpactSummary.conflicts}</div>
          </div>
          <p className="text-[10px] text-cyan-50/90">
            Polling Load Δ +{preCommitLoadEstimate.deltaReadsPerSec}/s ({preCommitLoadEstimate.baselineReadsPerSec}
            /s → {preCommitLoadEstimate.projectedReadsPerSec}/s, 假設週期 {preCommitLoadEstimate.assumedIntervalMs}ms)
          </p>
          <p className={`text-[10px] ${motionQAGate.pass ? 'text-emerald-200' : 'text-rose-200'}`}>
            Motion QA Gate: {motionQAGate.pass ? 'PASS' : 'FAIL'} ({motionQAGate.checklist.filter((item) => item.pass).length}/
            {motionQAGate.checklist.length})
          </p>
        </div>
        <div className="max-h-40 space-y-1.5 overflow-y-auto pr-1">
          {commitQueueItems.length === 0 ? (
            <p className="rounded border border-slate-700 bg-slate-900/60 px-2 py-2 text-[11px] text-slate-400">
              尚無待提交規劃
            </p>
          ) : (
            commitQueueItems.map((item) => (
              <div
                key={item.id}
                className={`rounded border px-2 py-1.5 text-[11px] ${
                  item.viewStatus === 'conflict' || item.viewStatus === 'failed'
                    ? 'border-rose-500/30 bg-rose-500/10 text-rose-100'
                    : item.viewStatus === 'linked'
                      ? 'border-indigo-500/30 bg-indigo-500/10 text-indigo-100'
                      : item.viewStatus === 'committed'
                        ? 'border-emerald-500/30 bg-emerald-500/10 text-emerald-100'
                        : 'border-slate-700 bg-slate-900/60 text-slate-200'
                }`}
              >
                <div className="flex items-center justify-between">
                  <span className="font-mono">#{item.order} {item.label}</span>
                  <span className="text-[10px] uppercase">{item.viewStatus}</span>
                </div>
                <div className="mt-0.5 text-[10px] opacity-80">
                  {item.type} · {item.addresses[0]}..{item.addresses[item.addresses.length - 1]}
                </div>
              </div>
            ))
          )}
        </div>
      </div>
      <div className="space-y-2 border-t border-white/5 p-4">
        <button
          type="button"
          onClick={onValidateFlow}
          disabled={!canValidate || validating}
          aria-keyshortcuts="Control+Enter"
          className="w-full rounded-lg border border-blue-500/40 bg-blue-500/20 px-3 py-2 text-xs font-medium text-blue-100 hover:bg-blue-500/30 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-blue-500 focus-visible:ring-offset-2 focus-visible:ring-offset-slate-900 disabled:cursor-not-allowed disabled:opacity-50"
        >
          {validating ? t('smartDashboard.validating') : t('smartDashboard.validateFlow')}
        </button>
        <button
          type="button"
          onClick={onCommitFlow}
          disabled={!canCommit || isCommitRunning}
          aria-keyshortcuts="Control+Shift+Enter"
          className="w-full rounded-lg border border-emerald-500/40 bg-emerald-500/20 px-3 py-2 text-xs font-medium text-emerald-100 hover:bg-emerald-500/30 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-emerald-500 focus-visible:ring-offset-2 focus-visible:ring-offset-slate-900 disabled:cursor-not-allowed disabled:opacity-50"
        >
          {isCommitRunning ? 'Commit 執行中...' : 'Commit 到 DB'}
        </button>
        <div className="grid grid-cols-2 gap-2">
          <button
            type="button"
            onClick={onRetryFailed}
            disabled={!hasFailedChunk}
            className="rounded-lg border border-amber-500/40 bg-amber-500/20 px-3 py-2 text-xs font-medium text-amber-100 hover:bg-amber-500/30 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-amber-500 disabled:cursor-not-allowed disabled:opacity-50"
          >
            Retry 失敗 Chunk
          </button>
          <button
            type="button"
            onClick={onRollback}
            disabled={!canRollback}
            className="rounded-lg border border-slate-600 bg-slate-700/50 px-3 py-2 text-xs font-medium text-slate-100 hover:bg-slate-700 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-slate-400 disabled:cursor-not-allowed disabled:opacity-50"
          >
            Rollback
          </button>
        </div>
        <div className="space-y-1 rounded-lg border border-white/10 bg-slate-900/60 p-2">
          {segmentFeedback.map((segment) => (
            <div key={segment.id} className="flex items-center justify-between text-[11px]">
              <span className="text-slate-300">{segment.label}</span>
              <span className={segment.ok ? 'text-emerald-300' : 'text-amber-300'}>
                {segment.message}
              </span>
            </div>
          ))}
        </div>
        {commitActionMessage && (
          <p className="rounded border border-slate-700 bg-slate-900/70 px-2 py-1.5 text-[11px] text-slate-300">
            {commitActionMessage}
          </p>
        )}
        {commitChunkResults.length > 0 && (
          <div className="space-y-1 rounded-lg border border-white/10 bg-slate-900/60 p-2">
            <p className="text-[11px] font-semibold text-slate-200">Chunk 結果</p>
            {commitChunkResults.map((chunkResult) => (
              <div key={chunkResult.chunk} className="flex items-center justify-between text-[10px]">
                <span className="text-slate-300">
                  Chunk {chunkResult.chunk}/{chunkResult.totalChunks}
                </span>
                <span className={chunkResult.status === 'success' ? 'text-emerald-300' : 'text-amber-300'}>
                  success {chunkResult.success} / failed {chunkResult.failed}
                </span>
              </div>
            ))}
          </div>
        )}
        {commitAuditPayload && (
          <div id="commit-audit-trace" className="space-y-1.5 rounded-lg border border-cyan-500/20 bg-slate-900/70 p-2">
            <div className="flex items-center justify-between">
              <p className="text-[11px] font-semibold text-cyan-200">Audit Trace</p>
              <a
                href={commitAuditPayload.traceLink}
                className="text-[10px] text-cyan-300 underline decoration-cyan-400/40 underline-offset-2 hover:text-cyan-200"
              >
                {commitAuditPayload.traceId}
              </a>
            </div>
            <p className="text-[10px] text-slate-400">{commitAuditPayload.createdAt}</p>
            <div className="flex flex-wrap gap-1.5 text-[10px]">
              <span className="rounded bg-emerald-500/10 px-1.5 py-0.5 text-emerald-200">
                new {commitAuditPayload.summary.newPoints}
              </span>
              <span className="rounded bg-blue-500/10 px-1.5 py-0.5 text-blue-200">
                tag {commitAuditPayload.summary.globalTagUpdates}
              </span>
              <span className="rounded bg-rose-500/10 px-1.5 py-0.5 text-rose-200">
                conflict {commitAuditPayload.summary.conflicts}
              </span>
            </div>
          </div>
        )}
        {hasError && (
          <button
            type="button"
            onClick={onRecoverFlow}
            aria-keyshortcuts="Alt+R"
            className="w-full rounded-lg border border-amber-500/40 bg-amber-500/20 px-3 py-2 text-xs font-medium text-amber-100 hover:bg-amber-500/30 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-amber-500 focus-visible:ring-offset-2 focus-visible:ring-offset-slate-900"
          >
            {t('smartDashboard.recoverFlow')}
          </button>
        )}
      </div>
    </>
  );
}
