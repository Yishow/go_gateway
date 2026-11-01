import { useId, useLayoutEffect, useRef, useState, type ReactNode } from 'react';
import { useQuery } from '@tanstack/react-query';
import { useTranslation } from 'react-i18next';
import type { StudioLocatorEntryRead } from '@/types/studioV2LocatorEntry';
import { dbTargetAPI } from '@/services/datalink';
import { safeRuntimeLabel } from '@/utils/runtimeOperatorSafety';
import { useRecordingDraftGuard } from '../../../recording-editor/RecordingDraftNavigation';
import { useRecordingChangeController } from './lifecycle/useRecordingChangeController';
import { RecordingChangeDialog } from './lifecycle/RecordingChangeDialog';
import type { SavedRecordingActionHost } from './lifecycle/SavedRecordingActions';

/** Scoped presentation around the original successor owner, never workspace hydration. */
export function LocatorRecordingLifecycleHost({ entry, scopeLease, readonly, canUseScope, onReload, children }: {
  entry: StudioLocatorEntryRead; scopeLease: object; readonly: boolean; canUseScope: () => boolean;
  onReload: () => void; children: (host: SavedRecordingActionHost, locked: boolean, isOriginalScope: () => boolean) => ReactNode;
}) {
  const { t } = useTranslation('recording-lifecycle');
  const region = useRef<HTMLDivElement>(null);
  const target = entry.selection.state === 'selected' ? entry.selection.target : undefined;
  const group = entry.entry.write_groups?.groups.find(value => value.id === target?.groupId);
  const exact = Boolean(entry.entry.editor_hydrated && target?.step === 4 && target.workspaceId === entry.workspace.id &&
    entry.locator.workspace_id === target.workspaceId && entry.locator.group_id === target.groupId && group && group.status !== 'deleted' &&
    group.workspace_id === target.workspaceId && group.revision === target.groupRevision && entry.entry.group_heads[group.id] === group.revision &&
    entry.entry.workspace_revision === target.workspaceRevision && entry.entry.write_groups?.workspace_revision === target.workspaceRevision);
  const live = useRef({ entry, scopeLease, readonly }); live.current = { entry, scopeLease, readonly };
  const owner = scopeLease;
  // Reuse the WGS lease. A retained editor cannot borrow a new A → B → A route lease.
  const admittedScope = useRef(scopeLease);
  const scopeAdmitted = admittedScope.current === scopeLease;
  const isOriginalScope = () => admittedScope.current === live.current.scopeLease;
  const mounted = useRef(false);
  const scopeReader = useRef(canUseScope); scopeReader.current = canUseScope;
  const current = () => mounted.current && live.current.entry === entry && live.current.scopeLease === owner && !readonly && exact && scopeAdmitted && scopeReader.current();
  const changed = useRef<object | undefined>(undefined);
  const foreignRecoveryRef = useRef(false);
  const view = useRecordingChangeController({ workspaceId: target?.workspaceId ?? '', readonly: readonly || !exact, canUseScope: () => current() && !foreignRecoveryRef.current,
    onChanged: () => { if (current()) changed.current = owner; } });
  const { controller, state } = view;
  const recoveryGroup = state.unreadableScope?.groupId ?? state.recovery?.envelope.predecessor_group_id;
  const foreignRecovery = Boolean(recoveryGroup && recoveryGroup !== target?.groupId);
  // The original recovery parser binds scheduled receipts to their submitted
  // confirmation. A retained terminal history is not a pending foreign lease.
  const verifiedTerminal = !state.unreadableScope && !state.recovery?.unreadable && !view.unresolved &&
    state.recovery?.phase === 'scheduled' && state.recovery.outcome === 'recorded' && Boolean(state.recovery.result);
  const foreignBlocked = foreignRecovery && !verifiedTerminal;
  foreignRecoveryRef.current = foreignBlocked;
  const recoveryAllowed = !foreignRecovery && exact && scopeAdmitted;
  const operationLock = useRef(false);
  const [destinationDemand, setDestinationDemand] = useState<object>();
  const instance = useId();
  const destinations = useQuery({
    queryKey: ['locator-lifecycle-destinations', instance, target?.workspaceId, target?.groupId, target?.workspaceRevision, target?.groupRevision],
    // Only safe display metadata is cached. Credentials and raw diagnostics never enter this surface.
    queryFn: async () => (await dbTargetAPI.listConnectors({ enabled: true })).flatMap(value => value.enabled && value.id && value.identity_revision
      ? [{ id: value.id, revision: value.identity_revision, name: safeRuntimeLabel(value.name) || value.id }] : []),
    enabled: destinationDemand === owner && state.isOpen && recoveryAllowed && current(),
    retry: false, refetchOnWindowFocus: false, staleTime: 0, gcTime: 0,
  });
  const destinationBlocked = destinationDemand === owner && (destinations.isPending || destinations.isFetching || destinations.isError);
  useLayoutEffect(() => { mounted.current = true; return () => { mounted.current = false; }; }, []);
  const previousOwner = useRef({ owner, entry });
  useLayoutEffect(() => {
    if (previousOwner.current.owner !== owner || previousOwner.current.entry !== entry) { previousOwner.current = { owner, entry }; controller.dismiss(); }
    // A different group's persisted operation is a denial, never a reason to query its scope here.
    controller.readonly = !current() || foreignBlocked || destinationBlocked;
  });
  useLayoutEffect(() => {
    if (changed.current !== owner || state.isOpen || state.busy || view.unresolved || !current()) return;
    changed.current = undefined; onReload();
  });
  useRecordingDraftGuard({ version: JSON.stringify([target?.workspaceId, target?.groupId, target?.workspaceRevision, target?.groupRevision, state.recovery?.envelope.request.operation_id, state.phase]),
    dirty: false, pending: recoveryAllowed && (state.busy || view.unresolved) });
  const host: SavedRecordingActionHost = {
    tryAcquire: () => {
      if (!current() || foreignBlocked || operationLock.current || controller.state.isOpen || controller.state.busy || controller.unresolved) return false;
      operationLock.current = true; return true;
    },
    // Ordinary stop/remove already own their busy/read-only lease. Do not invalidate that lease mid-request.
    onBusyChange: active => { operationLock.current = active; },
    onAction: intent => {
      if (!current() || foreignBlocked || operationLock.current || intent.workspaceId !== target?.workspaceId || intent.groupId !== target.groupId ||
          intent.workspaceRevision !== target.workspaceRevision || intent.groupRevision !== target.groupRevision) return;
      if (intent.action === 'edit') { region.current?.querySelector<HTMLInputElement>('[data-testid="group-name"]')?.focus(); return; }
      if (intent.action !== 'transition' && intent.action !== 'recreate') return;
      operationLock.current = true; setDestinationDemand(owner);
      void view.requestAction(intent).finally(() => { operationLock.current = false; });
    },
  };
  controller.readonly = !current() || foreignBlocked || destinationBlocked;
  const locked = readonly || !exact || foreignBlocked || state.isOpen || state.busy || view.unresolved;
  return <div ref={region} data-testid="locator-recording-lifecycle-host">
    {children(host, locked, isOriginalScope)}
    {(foreignBlocked || !scopeAdmitted) && <p role="alert">{t('reason.lifecycle_unavailable')} {t('next.review_lifecycle')}</p>}
    {recoveryAllowed && <>
      {destinationDemand === owner && destinations.isError && <div role="alert"><p>{t('reason.connector_unavailable')}</p>
        <button type="button" className="gw-secondary" onClick={() => { void destinations.refetch(); }}>{t('refresh')}</button></div>}
      <RecordingChangeDialog controller={view} destinations={!destinationBlocked && destinationDemand === owner ? destinations.data : []}
        fallbackFocus={() => region.current?.querySelector<HTMLElement>('[data-testid="group-lifecycle-return"]') ?? null} />
    </>}
  </div>;
}
