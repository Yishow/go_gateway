import * as React from 'react';
import { Icon } from '../components';
import type { StudioV2WorkspaceAuditEntry } from '../../../../types/studioV2WorkspaceAudit';

interface WorkspaceAuditHistoryPanelProps {
  entries?: StudioV2WorkspaceAuditEntry[];
  unavailable?: boolean;
}

export const WorkspaceAuditHistoryPanel: React.FC<WorkspaceAuditHistoryPanelProps> = ({
  entries = [],
  unavailable = false,
}) => {
  const recentEntries = entries.slice(0, 3);

  return (
    <div
      className="rounded-xl border border-slate-700/60 bg-slate-900/40 p-3"
      data-testid="summary-rail-audit-history"
    >
      <div className="mb-2 flex items-center gap-2">
        <Icon name="info" className="h-3.5 w-3.5 text-slate-400" />
        <span className="text-[11px] font-semibold uppercase tracking-wider text-slate-400">
          最近事件
        </span>
      </div>

      {unavailable ? (
        <div className="text-[11px] text-amber-200/80">Audit history unavailable</div>
      ) : recentEntries.length === 0 ? (
        <div className="text-[11px] italic text-slate-500">尚無事件紀錄</div>
      ) : (
        <div className="space-y-2">
          {recentEntries.map((entry) => (
            <div key={entry.id} className="rounded-lg border border-slate-800/80 bg-slate-950/40 p-2">
              <div className="flex items-center justify-between gap-2">
                <span className="truncate font-mono text-[10px] text-slate-200">
                  {entry.event_type}
                </span>
                <span className="font-mono text-[10px] text-cyan-200">{entry.result}</span>
              </div>
              <div className="mt-1 truncate font-mono text-[10px] text-slate-500">
                {entry.scope || 'workspace'}
              </div>
              <div className="mt-1 font-mono text-[10px] text-slate-500">
                {entry.occurred_at}
              </div>
            </div>
          ))}
        </div>
      )}
    </div>
  );
};
