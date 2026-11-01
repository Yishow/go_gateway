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
      className="rounded-xl border border-[var(--gw-border)] bg-[var(--gw-surface)] p-3"
      data-testid="summary-rail-audit-history"
    >
      <div className="mb-2 flex items-center gap-2">
        <Icon name="info" className="h-3.5 w-3.5 text-[var(--gw-muted)]" />
        <span className="text-xs font-semibold text-[var(--gw-muted)]">
          最近事件
        </span>
      </div>

      {unavailable ? (
        <div className="text-xs text-[var(--gw-warning-fg)]">Audit history unavailable</div>
      ) : recentEntries.length === 0 ? (
        <div className="text-xs italic text-[var(--gw-muted)]">尚無事件紀錄</div>
      ) : (
        <div className="space-y-2">
          {recentEntries.map((entry) => (
            <div key={entry.id} className="rounded-lg border border-[var(--gw-border)] bg-[var(--gw-surface)] p-2">
              <div className="flex items-center justify-between gap-2">
                <span className="truncate font-mono text-xs text-[var(--gw-text)]">
                  {entry.event_type}
                </span>
                <span className="font-mono text-xs text-[var(--gw-accent)]">{entry.result}</span>
              </div>
              <div className="mt-1 truncate font-mono text-xs text-[var(--gw-muted)]">
                {entry.scope || 'workspace'}
              </div>
              <div className="mt-1 font-mono text-xs text-[var(--gw-muted)]">
                {entry.occurred_at}
              </div>
            </div>
          ))}
        </div>
      )}
    </div>
  );
};
