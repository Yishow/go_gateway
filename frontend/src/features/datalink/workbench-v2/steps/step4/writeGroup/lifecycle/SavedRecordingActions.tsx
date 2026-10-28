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
export function SavedRecordingActions({ lifecycle, scope, disabled, onAction, canReview = true }: {
  lifecycle?: WriteGroupLifecycle;
  scope: Omit<SavedRecordingActionIntent, 'action'>;
  disabled: boolean;
  onAction: (intent: SavedRecordingActionIntent) => void;
  canReview?: boolean;
}) {
  const { t } = useTranslation('recording-lifecycle');
  const instanceId = useId();
  const matching = lifecycleMatches(lifecycle, scope);
  return <div className="space-y-2" data-testid="saved-recording-actions">
    <div className="flex flex-wrap gap-2">
      {ACTIONS.map((action) => {
        const availability = matching ? lifecycle.actions.find((item) => item.action === action) : undefined;
        const blocked = !availability?.allowed || disabled || ((!canReview) && (action === 'recreate' || action === 'transition'));
        return <button key={action} type="button" disabled={blocked} aria-describedby={!availability?.allowed ? `saved-action-${instanceId}-${action}` : undefined}
          data-testid={TEST_IDS[action]} onClick={() => onAction({ ...scope, action })}
          className="rounded border border-slate-600 px-3 py-1 text-xs text-slate-200 disabled:opacity-40">
          {t(`action.${action}`)}{availability?.requires_review && <span className="ml-1 text-[10px]">· {t('review_required')}</span>}
        </button>;
      })}
    </div>
    {ACTIONS.map((action) => {
      const availability = matching ? lifecycle.actions.find((item) => item.action === action) : undefined;
      if (availability?.allowed) return null;
      const reason = availability?.reason_code ?? 'lifecycle_unavailable';
      const next = availability?.next_action ?? 'reload';
      return <p key={action} id={`saved-action-${instanceId}-${action}`} data-testid={`saved-recording-reason-${action}`} className="text-xs text-amber-300">
        {t(`action.${action}`)}: {t(`reason.${REASONS.has(reason) ? reason : 'lifecycle_unavailable'}`)}{' '}
        {t(`next.${NEXT.has(next) ? next : 'review_lifecycle'}`)}
      </p>;
    })}
  </div>;
}
