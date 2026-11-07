import { useId } from 'react';
import { useTranslation } from 'react-i18next';
import type { WriteGroup } from '@/types/studioV2WriteGroup';
import type { SavedRecordingActionIntent, WriteGroupLifecycle, WriteGroupLifecycleActionName } from '@/types/studioV2WriteGroupLifecycle';

export interface SavedRecordingActionHost {
  onAction: (intent: SavedRecordingActionIntent) => void;
  tryAcquire?: () => boolean;
  onBusyChange?: (active: boolean) => void;
  onReviewScopedStart?: (group: WriteGroup) => void;
}

const ACTIONS: WriteGroupLifecycleActionName[] = ['edit', 'stop', 'remove', 'recreate', 'transition'];
const TEST_IDS: Record<WriteGroupLifecycleActionName, string> = {
  edit: 'saved-recording-edit', stop: 'group-disable', remove: 'group-delete', recreate: 'saved-recording-recreate', transition: 'saved-recording-transition',
};
const REASONS = new Set(['successor_reserved', 'recording_removed', 'already_stopped', 'connector_identity_stale', 'connector_unavailable', 'ownership_unavailable', 'successor_unavailable', 'source_basis_unavailable', 'successor_source_unavailable', 'lifecycle_unavailable', 'source_not_canonical', 'source_not_current', 'successor_requires_applied', 'successor_requires_removed', 'canonical_required', 'historical_generation', 'active_basic_unsupported', 'source_definition_unavailable', 'recording_not_removed', 'source_draft_not_active']);
const NEXT = new Set(['review_lifecycle', 'review_connector', 'inspect_ownership', 'stop', 'reload', 'review_successor', 'recreate']);

export function lifecycleMatches(data: WriteGroupLifecycle | undefined, scope: Omit<SavedRecordingActionIntent, 'action'>): data is WriteGroupLifecycle {
  return Boolean(data && data.workspace_id === scope.workspaceId && data.group_id === scope.groupId &&
    data.workspace_revision === scope.workspaceRevision && data.group_revision === scope.groupRevision);
}

/** Both saved surfaces render the same server action facts. No status-derived capability guesses. */
export function SavedRecordingActions({ lifecycle, scope, disabled, onAction, canReview = true, compact = false }: {
  lifecycle?: WriteGroupLifecycle;
  scope: Omit<SavedRecordingActionIntent, 'action'>;
  disabled: boolean;
  onAction: (intent: SavedRecordingActionIntent) => void;
  canReview?: boolean;
  compact?: boolean;
}) {
  const { t } = useTranslation('recording-lifecycle');
  const { t: workbench } = useTranslation('workbench-v2');
  const { t: deletion } = useTranslation('recording-home');
  const actionName = (action: WriteGroupLifecycleActionName) => action === 'remove' ? deletion('deleteDialogs.recording.confirm') : t(`action.${action}`);
  const instanceId = useId();
  const ReasonTag = compact ? 'div' : 'p';
  const matching = lifecycleMatches(lifecycle, scope);
  // Presentation groups only identical server reason/next pairs; permissions stay per action.
  const reasons = ACTIONS.flatMap(action => {
    const availability = matching ? lifecycle.actions.find(item => item.action === action) : undefined;
    return availability?.allowed ? [] : [{ action, reason: availability?.reason_code ?? 'lifecycle_unavailable', next: availability?.next_action ?? 'reload' }];
  });
  const firstReason = (action: WriteGroupLifecycleActionName) => {
    const item = reasons.find(row => row.action === action);
    return compact && item ? reasons.find(row => row.reason === item.reason && row.next === item.next)!.action : action;
  };
  return <div className="space-y-2" data-testid="saved-recording-actions">
    <div className="flex flex-wrap gap-2">
      {ACTIONS.map((action) => {
        const availability = matching ? lifecycle.actions.find((item) => item.action === action) : undefined;
        const blocked = !availability?.allowed || disabled || ((!canReview) && (action === 'recreate' || action === 'transition'));
        return <button key={action} type="button" disabled={blocked} aria-describedby={!availability?.allowed ? `saved-action-${instanceId}-${firstReason(action)}` : undefined}
          data-testid={TEST_IDS[action]} onClick={() => onAction({ ...scope, action })}
          className={compact ? 'rounded border border-[var(--gw-control)] bg-[var(--gw-surface)] px-3 py-2 text-sm text-[var(--gw-text)] disabled:bg-[var(--gw-ground)] disabled:text-[var(--gw-muted)] disabled:opacity-100' : 'rounded border border-slate-600 px-3 py-1 text-xs text-slate-200 disabled:opacity-40'}>
          {actionName(action)}{availability?.requires_review && <span className="ml-1 text-[10px]">· {t('review_required')}</span>}
        </button>;
      })}
    </div>
    {ACTIONS.map((action) => {
      const availability = matching ? lifecycle.actions.find((item) => item.action === action) : undefined;
      if (availability?.allowed) return null;
      const reason = availability?.reason_code ?? 'lifecycle_unavailable';
      const next = availability?.next_action ?? 'reload';
      if (compact && firstReason(action) !== action) return null;
      const names = compact ? reasons.filter(row => row.reason === reason && row.next === next).map(row => actionName(row.action)).join(' / ') : actionName(action);
      return <ReasonTag key={action} id={`saved-action-${instanceId}-${action}`} data-testid={`saved-recording-reason-${action}`} className={compact ? 'text-sm text-[var(--gw-warning-fg)]' : 'text-xs text-amber-300'}>
        {names}: {t(`reason.${REASONS.has(reason) ? reason : 'lifecycle_unavailable'}`)}{' '}
        {t(`next.${NEXT.has(next) ? next : 'review_lifecycle'}`)}
        {compact && (!REASONS.has(reason) || !NEXT.has(next)) && <details className="text-[var(--gw-muted)]"><summary>{workbench('step4.basic.diagnostics')}</summary><p>{reason} / {next}</p></details>}
      </ReasonTag>;
    })}
  </div>;
}
