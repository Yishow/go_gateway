import { savedRecordingStartDevices } from '../../../recording-editor/savedRecordingStartScope';
import type { StudioV2LocatorReadScope } from '../../../../../../services/studioV2WorkspaceDatabase';
import * as React from 'react';
import { useEffect, useRef, useState } from 'react';
import { useTranslation } from 'react-i18next';
import { useWriteGroupLifecycleMutation } from '../../../../../../hooks/datalink/useStudioV2WriteGroups';
import { useWriteGroupLifecycleQuery } from '@/hooks/datalink/useStudioV2WriteGroupLifecycle';
import type { WriteGroup } from '../../../../../../types/studioV2WriteGroup';
import type { SavedRecordingActionIntent } from '@/types/studioV2WriteGroupLifecycle';
import { normalizeTypedEnvelope } from '../../../../../../utils/safeJson';
import { getSafeErrorMessage } from '../../../../../../utils/typedErrors';
import { managedProfileDeviceId } from '../../../state/circuitProfileBasic';
import { lifecycleMatches, SavedRecordingActions, type SavedRecordingActionHost } from './lifecycle/SavedRecordingActions';
import { SavedRecordingRemoveDialog } from './lifecycle/SavedRecordingRemoveDialog';

export interface GroupLifecycleBarProps {
  locatorScope?: StudioV2LocatorReadScope;
  compact?: boolean;
  /** Denial-only incoming navigation fence; never a read scope or write permission. */
  locatorGuard?: boolean;
  group: WriteGroup;
  workspaceId: string;
  workspaceRevision: string;
  /** The exact saved revision is ready and the editor has no unsaved edits. */
  canApply: boolean;
  applyBlockedReason?: string;
  sharedReadinessReason?: string;
  readonly: boolean;
  onChanged: (group: WriteGroup, workspaceRevision?: string) => void;
  onConflict: () => void;
  onReviewScopedStart?: (group: WriteGroup) => void;
  lifecycleHost?: SavedRecordingActionHost;
  /** Basic keeps its existing reviewed start; this bar only manages saved recordings. */
  hideApply?: boolean;
  /** Simple presentation asks for explicit scope/impact consent before Stop. */
  confirmStop?: boolean;
}

/** One action surface in Basic and Advanced, with server authority refreshed before ordinary mutations. */
export const GroupLifecycleBar: React.FC<GroupLifecycleBarProps> = ({
  group, workspaceId, workspaceRevision, canApply, applyBlockedReason, sharedReadinessReason, readonly, onChanged, onConflict,
  onReviewScopedStart, lifecycleHost, hideApply = false, confirmStop = false, locatorScope, locatorGuard, compact = false,
}) => {
  const { t } = useTranslation('workbench-v2');
  const { t: recordingT } = useTranslation('recording-editor');
  const { t: lifecycleT } = useTranslation('recording-lifecycle');
  const apply = useWriteGroupLifecycleMutation('apply');
  const disable = useWriteGroupLifecycleMutation('disable');
  const remove = useWriteGroupLifecycleMutation('remove');
  const lifecycle = useWriteGroupLifecycleQuery(workspaceId, group.id);
  const [confirmation, setConfirmation] = useState<{ action: 'stop' | 'remove'; owner: object; lease: object } | null>(null);
  const [error, setError] = useState<unknown>(null);
  const [refreshBlocked, setRefreshBlocked] = useState(false);
  const [checking, setChecking] = useState(false);
  const busyRef = useRef(false);
  const mountedRef = useRef(false);
  const returnFocus = useRef<HTMLElement | null>(null);
  const scope = JSON.stringify([group.id, group.revision, workspaceId, workspaceRevision]);
  const scopeRef = useRef(scope);
  scopeRef.current = scope;
  const scopeOwner = useRef({ scope });
  if (scopeOwner.current.scope !== scope) scopeOwner.current = { scope };
  const returnContext = useRef<HTMLHeadingElement>(null);
  // This transient denial lease cannot be revived by an A -> B -> A navigation.
  const locatorLeaseKey = locatorScope || locatorGuard ? JSON.stringify([scope, locatorScope?.workspaceId, locatorScope?.workspaceRevision, readonly]) : null;
  const locatorLease = useRef({ key: locatorLeaseKey, blocked: readonly });
  if (locatorLease.current.key !== locatorLeaseKey) locatorLease.current = { key: locatorLeaseKey, blocked: readonly };
  const renderLocatorLease = locatorLease.current;
  const confirmingAction = confirmation?.owner === scopeOwner.current && confirmation.lease === renderLocatorLease ? confirmation.action : undefined;
  useEffect(() => {
    mountedRef.current = true;
    return () => { mountedRef.current = false; };
  }, []);
  const intentScope = { workspaceId, groupId: group.id, workspaceRevision, groupRevision: group.revision };
  const snapshot = !lifecycle.isError && lifecycleMatches(lifecycle.data, intentScope) ? lifecycle.data : undefined;
  const deleted = group.status === 'deleted';
  const busy = checking || apply.isPending || disable.isPending || remove.isPending;
  const reviewStart = onReviewScopedStart ?? lifecycleHost?.onReviewScopedStart;
  const routed = reviewStart !== undefined && (managedProfileDeviceId(group) !== undefined || savedRecordingStartDevices(group).length > 0);

  const run = async (action: 'apply' | 'stop' | 'remove', confirm = false) => {
    if (confirm && (confirmation?.action !== action || confirmation.owner !== scopeOwner.current || confirmation.lease !== renderLocatorLease)) return;
    if (readonly || busyRef.current || (action === 'apply' && (!canApply || deleted)) || locatorLease.current !== renderLocatorLease || renderLocatorLease.key !== null &&
      (!mountedRef.current || locatorLease.current.blocked)) return;
    if (lifecycleHost?.tryAcquire && !lifecycleHost.tryAcquire()) return;
    busyRef.current = true;
    lifecycleHost?.onBusyChange?.(true);
    const requestScope = scopeRef.current;
    const requestOwner = scopeOwner.current;
    const requestLease = renderLocatorLease;
    let commandStarted = false;
    const dismissPreCommandConfirmation = () => setConfirmation((current) =>
      current?.action === action && current.owner === requestOwner && current.lease === requestLease ? null : current);
    setError(null); setRefreshBlocked(false); setChecking(true);
    try {
      if (action !== 'apply') {
        const refreshed = await lifecycle.refetch();
        const fresh = refreshed.isError || refreshed.error ? undefined : refreshed.data;
        if (!mountedRef.current || scopeRef.current !== requestScope || scopeOwner.current !== requestOwner ||
          locatorLease.current !== requestLease || requestLease.key !== null && locatorLease.current.blocked) return;
        if (!lifecycleMatches(fresh, intentScope) || !fresh.actions.some((entry) => entry.action === action && entry.allowed)) {
          dismissPreCommandConfirmation(); setRefreshBlocked(true); onConflict(); return;
        }
        if (!confirm && (action === 'remove' || action === 'stop' && confirmStop)) { setConfirmation({ action, owner: requestOwner, lease: requestLease }); return; }
      }
      const mutation = action === 'apply' ? apply : action === 'stop' ? disable : remove;
      commandStarted = true;
      const result = await mutation.mutateAsync({ id: group.id, request: {
        workspace_id: workspaceId, expected_workspace_revision: workspaceRevision,
        expected_group_revision: group.revision, expected_connector_revision: group.destination.connector_revision,
      } });
      if (mountedRef.current && scopeRef.current === requestScope && scopeOwner.current === requestOwner) {
        setConfirmation(null); onChanged(result.group, result.workspace_revision);
      }
    } catch (cause) {
      if (mountedRef.current && scopeRef.current === requestScope && scopeOwner.current === requestOwner) {
        if (!commandStarted) {
          if (locatorLease.current !== requestLease || requestLease.key !== null && locatorLease.current.blocked) return;
          dismissPreCommandConfirmation(); setRefreshBlocked(true);
        }
        setError(cause);
        if (normalizeTypedEnvelope(cause).code === 'revision_mismatch') onConflict();
      }
    } finally {
      busyRef.current = false;
      lifecycleHost?.onBusyChange?.(false);
      if (mountedRef.current) setChecking(false);
    }
  };
  const handleAction = (intent: SavedRecordingActionIntent) => {
    if (intent.action === 'stop' || intent.action === 'remove') {
      returnFocus.current = document.activeElement instanceof HTMLElement ? document.activeElement : null;
      void run(intent.action);
    } else if (lifecycleHost) lifecycleHost.onAction(intent);
    else if (intent.action === 'edit') document.querySelector<HTMLInputElement>('[data-testid="group-name"]')?.focus();
  };

  const management = <SavedRecordingActions compact={compact} lifecycle={snapshot} scope={intentScope} disabled={readonly || busy || Boolean(confirmingAction)} onAction={handleAction} canReview={Boolean(lifecycleHost)} />;
  return <div className="space-y-3" data-testid="group-lifecycle" data-focus-scope={scope}>
    <h3 ref={returnContext} tabIndex={-1} data-testid="group-lifecycle-return" className="sr-only focus:not-sr-only">{group.name}</h3>
    {!hideApply && <>
      <button type="button" onClick={() => (routed ? reviewStart?.(group) : void run('apply'))}
        aria-describedby={!canApply ? sharedReadinessReason : undefined} disabled={readonly || busy || deleted || !canApply} className="rounded border border-[var(--gw-success-border)] bg-[var(--gw-success-bg)] px-3 py-1 text-xs text-[var(--gw-success-fg)] disabled:text-[var(--gw-muted)]" data-testid="group-apply">
        {routed ? recordingT('start.reviewSaved') : t('step4.group.lifecycle.apply')}
      </button>
      {!sharedReadinessReason && !canApply && !deleted && applyBlockedReason && <p className="text-xs text-[var(--gw-warning-fg)]" data-testid="group-apply-blocked">{applyBlockedReason}</p>}
      <p className="text-xs text-[var(--gw-muted)]" data-testid="group-apply-note">{routed ? recordingT('start.reviewSavedNote') : t('step4.group.lifecycle.apply_note')}</p>
    </>}
    {compact ? <details className="recording-secondary-management" data-testid="recording-secondary-management"><summary>{recordingT('workflow.manage')}</summary>{management}<p className="text-xs text-[var(--gw-muted)]">{lifecycleT('stop.keeps')}</p></details> : management}
    {compact && !snapshot && <p role="status" className="text-sm text-[var(--gw-warning-fg)]">{lifecycleT('reason.lifecycle_unavailable')}</p>}
    {(!snapshot || refreshBlocked) && <button type="button" onClick={() => { void lifecycle.refetch(); onConflict(); }} className="text-xs underline" data-testid="saved-recording-refresh">{lifecycleT('refresh')}</button>}
    {!compact && <p className="text-xs text-[var(--gw-muted)]">{lifecycleT('stop.keeps')}</p>}
    <SavedRecordingRemoveDialog action={confirmation?.action} open={Boolean(confirmingAction)} busy={busy || readonly} history={snapshot?.history}
      onClose={() => setConfirmation(null)} onConfirm={() => { if (confirmingAction) void run(confirmingAction, true); }} returnFocus={returnFocus.current} objectLabel={group.name} focusOwner={scopeOwner.current} fallbackFocus={returnContext.current} />
    {refreshBlocked && <p role="alert" className="text-xs text-[var(--gw-warning-fg)]">{lifecycleT('reason.stale_scope')}</p>}
    {error !== null && <p role="alert" className="text-xs text-[var(--gw-danger-fg)]" data-testid="group-lifecycle-error">{getSafeErrorMessage(error, t).message}</p>}
  </div>;
};
