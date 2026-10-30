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
  /** Denial-only incoming navigation fence; never a read scope or write permission. */
  locatorGuard?: boolean;
  group: WriteGroup;
  workspaceId: string;
  workspaceRevision: string;
  /** The exact saved revision is ready and the editor has no unsaved edits. */
  canApply: boolean;
  applyBlockedReason?: string;
  readonly: boolean;
  onChanged: (group: WriteGroup, workspaceRevision?: string) => void;
  onConflict: () => void;
  onReviewScopedStart?: (group: WriteGroup) => void;
  lifecycleHost?: SavedRecordingActionHost;
  /** Basic keeps its existing reviewed start; this bar only manages saved recordings. */
  hideApply?: boolean;
}

/** One action surface in Basic and Advanced, with server authority refreshed before ordinary mutations. */
export const GroupLifecycleBar: React.FC<GroupLifecycleBarProps> = ({
  group, workspaceId, workspaceRevision, canApply, applyBlockedReason, readonly, onChanged, onConflict,
  onReviewScopedStart, lifecycleHost, hideApply = false, locatorScope, locatorGuard,
}) => {
  const { t } = useTranslation('workbench-v2');
  const { t: lifecycleT } = useTranslation('recording-lifecycle');
  const apply = useWriteGroupLifecycleMutation('apply');
  const disable = useWriteGroupLifecycleMutation('disable');
  const remove = useWriteGroupLifecycleMutation('remove');
  const lifecycle = useWriteGroupLifecycleQuery(workspaceId, group.id);
  const [deleteIntent, setDeleteIntent] = useState<object | null>(null);
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
  const confirmingDelete = deleteIntent === scopeOwner.current;
  const returnContext = useRef<HTMLHeadingElement>(null);
  // This transient denial lease cannot be revived by an A -> B -> A navigation.
  const locatorLeaseKey = locatorScope || locatorGuard ? JSON.stringify([scope, locatorScope?.workspaceId, locatorScope?.workspaceRevision, readonly]) : null;
  const locatorLease = useRef({ key: locatorLeaseKey, blocked: readonly });
  if (locatorLease.current.key !== locatorLeaseKey) locatorLease.current = { key: locatorLeaseKey, blocked: readonly };
  const renderLocatorLease = locatorLease.current;
  useEffect(() => {
    mountedRef.current = true;
    return () => { mountedRef.current = false; };
  }, []);
  const intentScope = { workspaceId, groupId: group.id, workspaceRevision, groupRevision: group.revision };
  const snapshot = !lifecycle.isError && lifecycleMatches(lifecycle.data, intentScope) ? lifecycle.data : undefined;
  const deleted = group.status === 'deleted';
  const busy = checking || apply.isPending || disable.isPending || remove.isPending;
  const reviewStart = onReviewScopedStart ?? lifecycleHost?.onReviewScopedStart;
  const routed = reviewStart !== undefined && managedProfileDeviceId(group) !== undefined;

  const run = async (action: 'apply' | 'stop' | 'remove', confirm = false) => {
    if (confirm && deleteIntent !== scopeOwner.current) return;
    if (readonly || busyRef.current || (action === 'apply' && (!canApply || deleted)) || locatorLease.current !== renderLocatorLease || renderLocatorLease.key !== null &&
      (!mountedRef.current || locatorLease.current.blocked)) return;
    if (lifecycleHost?.tryAcquire && !lifecycleHost.tryAcquire()) return;
    busyRef.current = true;
    lifecycleHost?.onBusyChange?.(true);
    const requestScope = scopeRef.current;
    const requestOwner = scopeOwner.current;
    const requestLease = renderLocatorLease;
    setError(null); setRefreshBlocked(false); setChecking(true);
    try {
      if (action !== 'apply') {
        const refreshed = await lifecycle.refetch();
        const fresh = refreshed.isError || refreshed.error ? undefined : refreshed.data;
        if (!mountedRef.current || scopeRef.current !== requestScope || scopeOwner.current !== requestOwner ||
          locatorLease.current !== requestLease || requestLease.key !== null && locatorLease.current.blocked) return;
        if (!lifecycleMatches(fresh, intentScope) || !fresh.actions.some((entry) => entry.action === action && entry.allowed)) {
          setRefreshBlocked(true); onConflict(); return;
        }
        if (action === 'remove' && !confirm) { setDeleteIntent(requestOwner); return; }
      }
      const mutation = action === 'apply' ? apply : action === 'stop' ? disable : remove;
      const result = await mutation.mutateAsync({ id: group.id, request: {
        workspace_id: workspaceId, expected_workspace_revision: workspaceRevision,
        expected_group_revision: group.revision, expected_connector_revision: group.destination.connector_revision,
      } });
      if (mountedRef.current && scopeRef.current === requestScope && scopeOwner.current === requestOwner) {
        setDeleteIntent(null); onChanged(result.group, result.workspace_revision);
      }
    } catch (cause) {
      if (mountedRef.current && scopeRef.current === requestScope && scopeOwner.current === requestOwner) {
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

  return <div className="space-y-3" data-testid="group-lifecycle" data-focus-scope={scope}>
    <h3 ref={returnContext} tabIndex={-1} data-testid="group-lifecycle-return" className="sr-only focus:not-sr-only">{group.name}</h3>
    {!hideApply && <>
      <button type="button" onClick={() => (routed ? reviewStart?.(group) : void run('apply'))}
        disabled={readonly || busy || deleted || !canApply} className="rounded bg-emerald-600 px-3 py-1 text-xs text-white disabled:opacity-40" data-testid="group-apply">
        {t(routed ? 'step4.group.lifecycle.review_start' : 'step4.group.lifecycle.apply')}
      </button>
      {!canApply && !deleted && applyBlockedReason && <p className="text-xs text-amber-300" data-testid="group-apply-blocked">{applyBlockedReason}</p>}
      <p className="text-xs text-slate-500" data-testid="group-apply-note">{t(routed ? 'step4.group.lifecycle.review_start_note' : 'step4.group.lifecycle.apply_note')}</p>
    </>}
    <SavedRecordingActions lifecycle={snapshot} scope={intentScope} disabled={readonly || busy || confirmingDelete} onAction={handleAction} canReview={Boolean(lifecycleHost)} />
    {(!snapshot || refreshBlocked) && <button type="button" onClick={() => { void lifecycle.refetch(); onConflict(); }} className="text-xs underline" data-testid="saved-recording-refresh">{lifecycleT('refresh')}</button>}
    <p className="text-xs text-slate-400">{lifecycleT('stop.keeps')}</p>
    <SavedRecordingRemoveDialog open={confirmingDelete} busy={busy || readonly} history={snapshot?.history}
      onClose={() => setDeleteIntent(null)} onConfirm={() => void run('remove', true)} returnFocus={returnFocus.current} objectLabel={group.name} focusOwner={scopeOwner.current} fallbackFocus={returnContext.current} />
    {refreshBlocked && <p role="alert" className="text-xs text-amber-300">{lifecycleT('reason.stale_scope')}</p>}
    {error !== null && <p role="alert" className="text-xs text-red-300" data-testid="group-lifecycle-error">{getSafeErrorMessage(error, t).message}</p>}
  </div>;
};
