import { readStudioLocatorEditorEntry } from '../../../../../../services/studioV2LocatorEntry';
import { studioV2LocatorHydration } from '../../../../../../pages/datalink/workbench-v2/studioV2LocatorHydration';
import type { StudioLocatorEntryRead } from '../../../../../../types/studioV2LocatorEntry';
import * as React from 'react';
import { useEffect, useMemo, useRef, useState } from 'react';
import { useTranslation } from 'react-i18next';
import { useWriteGroupsQuery } from '../../../../../../hooks/datalink/useStudioV2WriteGroups';
import type { WriteGroup } from '../../../../../../types/studioV2WriteGroup';
import type { WorkbenchV2State } from '../../../state/types';
import { buildGroupCandidates } from '../../../state/writeGroup/candidates';
import type { WorkspaceReadinessStepNumber } from '../../../components/WorkspaceReadinessPanel';
import { CutoverPanel } from './CutoverPanel';
import { GroupEditor } from './GroupEditor';
import { SavedRecordingCard } from '../SavedRecordingCard';
import type { SavedRecordingActionHost } from './lifecycle/SavedRecordingActions';

export interface WriteGroupSectionProps {
  locatorEntry?: StudioLocatorEntryRead;
  locatorOnly?: boolean;
  locatorSuspended?: boolean;
  onEditorMountedChange?: (mounted: boolean) => void;
  state: WorkbenchV2State;
  workspaceId?: string;
  readonly: boolean;
  lifecycleHost?: SavedRecordingActionHost;
  onNavigateStep?: (step: WorkspaceReadinessStepNumber) => void;
  openGroupRequest?: { id: string };
  /** The host's reviewed scoped start: a managed profile group's effective action goes there instead of applying on its own. */
  onReviewScopedStart?: (group: WriteGroup) => void;
}

type Selection = { kind: 'none' } | { kind: 'new' } | { kind: 'edit'; id: string; saved?: WriteGroup };

type Prerequisite = 'device' | 'tags' | 'destination' | null;

/**
 * The only database output editor of Step 4. It reads the saved groups of the
 * workspace and says what is missing before anything can be created, keeping a
 * failed load, an empty list and groups for another destination apart.
 */
export const WriteGroupSection: React.FC<WriteGroupSectionProps> = ({ state, workspaceId, readonly, onNavigateStep, openGroupRequest, onReviewScopedStart, lifecycleHost, locatorEntry, locatorOnly, locatorSuspended, onEditorMountedChange }) => {
  const { t } = useTranslation('workbench-v2');
  const { t: runtime } = useTranslation('runtime-dashboard');
  const region = useRef<HTMLElement>(null);
  const [selection, setSelection] = useState<Selection>({ kind: 'none' });
  const groupsQuery = useWriteGroupsQuery(Boolean(workspaceId) && !locatorOnly);
  const [editorEntry, setEditorEntry] = useState<StudioLocatorEntryRead>();
  const [openError, setOpenError] = useState<'unavailable' | 'stale-reference' | 'foreign-scope'>();
  const [opening, setOpening] = useState(false);
  const openingRequest = useRef<AbortController | undefined>(undefined);
  useEffect(() => {
    openingRequest.current?.abort(); openingRequest.current = undefined; setOpening(false);
    return () => { openingRequest.current?.abort(); openingRequest.current = undefined; };
  }, [locatorEntry, locatorSuspended]);
  const observed = locatorEntry?.entry.write_groups ?? groupsQuery.data;
  const editorState = useMemo(() => editorEntry ? studioV2LocatorHydration(editorEntry) : null, [editorEntry]);
  const editorObserved = editorEntry?.entry.write_groups ?? observed;
  const open = async (group: WriteGroup) => {
    if (!locatorEntry) { setSelection({ kind: 'edit', id: group.id, saved: group }); return; }
    const target = locatorEntry.selection.state === 'selected' ? locatorEntry.selection.target : undefined;
    if (!target?.groupRevision || target.groupId !== group.id || target.workspaceId !== workspaceId ||
      locatorSuspended || selection.kind !== 'none' || openingRequest.current) return;
    const request = new AbortController(); openingRequest.current = request; setOpening(true); setOpenError(undefined);
    try {
      const result = await readStudioLocatorEditorEntry(locatorEntry.locator,
        { workspaceRevision: target.workspaceRevision, groupRevision: target.groupRevision }, request.signal);
      if (request.signal.aborted || openingRequest.current !== request) return;
      if (!result.ok) { setOpenError(result.reason); return; }
      const next = result.value;
      const selected = next.selection.state === 'selected' ? next.selection.target : undefined;
      const saved = next.entry.write_groups?.groups.find(value => value.id === group.id && value.status !== 'deleted');
      if (!next.entry.editor_hydrated || !selected || selected.workspaceId !== target.workspaceId ||
        selected.workspaceRevision !== target.workspaceRevision || selected.groupId !== target.groupId ||
        selected.groupRevision !== target.groupRevision || !saved || !studioV2LocatorHydration(next)) {
        setOpenError('stale-reference'); return;
      }
      setEditorEntry(next); setSelection({ kind: 'edit', id: group.id, saved });
    } catch { if (!request.signal.aborted) setOpenError('unavailable'); }
    finally { if (openingRequest.current === request) { openingRequest.current = undefined; setOpening(false); } }
  };
  const reloadEditor = async () => {
    if (!editorEntry || openingRequest.current || locatorSuspended || selection.kind !== 'edit') return;
    const original = editorEntry.selection.state === 'selected' ? editorEntry.selection.target : undefined;
    if (!original?.groupId) return;
    const request = new AbortController(); openingRequest.current = request; setOpening(true); setOpenError(undefined);
    try {
      const latest = await groupsQuery.refetch();
      if (request.signal.aborted || openingRequest.current !== request) return;
      if (latest.isError || !latest.data) { setOpenError('unavailable'); return; }
      if (latest.data.workspace_id !== original.workspaceId) { setOpenError('foreign-scope'); return; }
      const exact = latest.data.groups.find(group => group.id === original.groupId && group.status !== 'deleted');
      if (!exact) { setOpenError('stale-reference'); return; }
      const result = await readStudioLocatorEditorEntry(editorEntry.locator,
        { workspaceRevision: latest.data.workspace_revision, groupRevision: exact.revision }, request.signal);
      if (request.signal.aborted || openingRequest.current !== request) return;
      if (!result.ok) { setOpenError(result.reason); return; }
      const target = result.value.selection.state === 'selected' ? result.value.selection.target : undefined;
      const group = result.value.entry.write_groups?.groups.find(value => value.id === original.groupId && value.status !== 'deleted');
      if (!target || target.workspaceId !== original.workspaceId || target.groupId !== original.groupId ||
        target.workspaceRevision !== latest.data.workspace_revision || target.groupRevision !== exact.revision ||
        !result.value.entry.editor_hydrated || !group || !studioV2LocatorHydration(result.value)) {
        setOpenError('stale-reference'); return;
      }
      setEditorEntry(result.value); setSelection({ kind: 'edit', id: group.id, saved: group });
    } catch { if (!request.signal.aborted) setOpenError('unavailable'); }
    finally { if (openingRequest.current === request) { openingRequest.current = undefined; setOpening(false); } }
  };
  const located = locatorEntry?.selection.state === 'selected' ? locatorEntry.selection.target : undefined;
  useEffect(() => { if (located) region.current?.querySelector<HTMLElement>('[data-locator-selected="true"]')?.focus(); }, [located]);
  const consumedOpen = useRef<typeof openGroupRequest>(undefined);
  useEffect(() => {
    if (openGroupRequest !== consumedOpen.current && openGroupRequest && groupsQuery.data?.groups.some(group => group.id === openGroupRequest.id)) {
      consumedOpen.current = openGroupRequest;
      setSelection({ kind: 'edit', id: openGroupRequest.id });
    }
  }, [openGroupRequest, groupsQuery.data]);
  const { candidates, excluded } = useMemo(() => buildGroupCandidates(editorState ?? state), [editorState, state]);
  const connector = (editorState ?? state).db.connector;
  const destinationSaved = connector.persisted === true && connector.save_state === 'saved' &&
    Boolean(connector.connector_id && connector.identity_revision);

  const prerequisite: Prerequisite = state.devices.length === 0 ? 'device'
    : candidates.length === 0 ? 'tags'
      : !destinationSaved ? 'destination' : null;

  const groups = observed?.groups ?? [];
  const live = groups.filter((group) => group.status !== 'deleted' && (!located || group.id === located.groupId));
  const removed = groups.filter((group) => group.status === 'deleted' && !locatorEntry);
  // Right after a save the refetched list may not carry the group yet; the returned group bridges that gap.
  const selected = selection.kind === 'edit' ? editorEntry ? selection.saved ?? editorObserved?.groups.find(group => group.id === selection.id)
    : groups.find(group => group.id === selection.id) ?? selection.saved : undefined;
  const management = { workspaceId, workspaceRevision: observed?.workspace_revision ?? '', lifecycleHost,
    readonly: readonly || !workspaceId || observed?.workspace_id !== workspaceId || !observed?.workspace_revision || !locatorOnly && groupsQuery.isFetching,
    onReload: () => { void groupsQuery.refetch(); } };
  const matches = (group: WriteGroup) => !destinationSaved ? null : group.destination.connector_id === connector.connector_id &&
    group.destination.connector_revision === connector.identity_revision;

  const editing = selection.kind === 'new' || Boolean(selection.kind === 'edit' && selected);
  const retention = useRef(onEditorMountedChange); retention.current = onEditorMountedChange;
  useEffect(() => { retention.current?.(editing); return () => retention.current?.(false); }, [editing]);
  if (!locatorOnly && groupsQuery.isError) {
    return (
      <section role="alert" className="space-y-2 rounded-xl border border-amber-700 p-5 text-xs text-amber-200" data-testid="group-section-error">
        <p>{t('step4.group.section.load_failed')}</p>
        <button type="button" onClick={() => void groupsQuery.refetch()} className="underline" data-testid="group-section-retry">{t('step4.group.retry')}</button>
      </section>
    );
  }
  if (!locatorOnly && groupsQuery.isLoading) {
    return <section className="p-5 text-xs text-slate-400" role="status" data-testid="group-section-loading">{t('step4.group.section.loading')}</section>;
  }

  return (
    <section ref={region} className="space-y-4 rounded-xl border border-slate-800 bg-slate-900/20 p-5" aria-labelledby="group-section-title" data-testid="group-section">
      <div>
        <h3 id="group-section-title" className="text-sm font-semibold text-slate-100">{t('step4.group.section.title')}</h3>
        <p className="mt-1 text-xs text-slate-400">{t('step4.group.section.subtitle')}</p>
      </div>

      {opening && <p role="status">{runtime('locatorEntry.loading')}</p>}
      {openError && <p role="alert" data-testid="studio-locator-open-error">{runtime(`locatorEntry.${openError === 'stale-reference' ? 'stale' : openError === 'foreign-scope' ? 'foreign' : 'unavailable'}`)}</p>}
      {locatorOnly && <p role="status">{runtime('locatorEntry.unrequested')}</p>}
      {located && selection.kind !== 'none' && (selection.kind !== 'edit' || selection.id !== located.groupId) && <p role="alert">{runtime('locatorEntry.draftBlocked')}</p>}
      {locatorEntry?.selection.state === 'selected' && locatorEntry.selection.history && <p>{runtime('locatorEntry.historicalHint', { revision: locatorEntry.selection.history.revisionHint })}</p>}
      {located?.groupId && located.groupRevision && located.workspaceId === workspaceId &&
        <p>{runtime('locatorEntry.currentRevision')}: <span>{located.groupRevision}</span></p>}
      {!locatorOnly && prerequisite && (
        <div className="space-y-1 text-xs text-amber-300" role="status" data-testid={`group-prerequisite-${prerequisite}`}>
          <p>{t(`step4.group.prerequisite.${prerequisite}`)}</p>
          {prerequisite === 'device' && onNavigateStep && (
            <button type="button" onClick={() => onNavigateStep(1)} className="underline" data-testid="group-prerequisite-go">{t('step4.group.prerequisite.go_step1')}</button>
          )}
          {prerequisite === 'tags' && (
            <>
              {excluded.length > 0 && <p data-testid="group-prerequisite-excluded">{t('step4.group.prerequisite.excluded', { count: excluded.length })}</p>}
              {onNavigateStep && <button type="button" onClick={() => onNavigateStep(3)} className="underline" data-testid="group-prerequisite-go">{t('step4.group.prerequisite.go_step3')}</button>}
            </>
          )}
        </div>
      )}

      {!editing && (
        <>
          {live.length === 0 && removed.length === 0 ? (
            <p className="text-xs text-slate-400" data-testid="group-empty">{t('step4.group.list.empty')}</p>
          ) : (
            <ul className="space-y-2" data-testid="group-list">
              {live.map((group) => <SavedRecordingCard key={group.id} group={group} highlighted={located?.groupId === group.id} current={matches(group)} onOpen={() => void open(group)}
                management={{ ...management, readonly: management.readonly || group.workspace_id !== workspaceId }} />)}
            </ul>
          )}
          {live.length > 0 && !live.some(matches) && destinationSaved && (
            <p className="text-xs text-amber-300" role="status" data-testid="group-none-match">{t('step4.group.list.none_match')}</p>
          )}
          {!locatorOnly && groupsQuery.data && observed && live.length > 0 && (
            <CutoverPanel groups={live} workspaceId={observed.workspace_id} workspaceRevision={observed.workspace_revision}
              connectorRevision={connector.identity_revision ?? ''} readonly={readonly || !destinationSaved} onOpenGroup={(id) => setSelection({ kind: 'edit', id })}
              candidates={candidates} />
          )}
          <button
            type="button" disabled={readonly || prerequisite !== null || !observed || Boolean(locatorOnly)}
            onClick={() => setSelection({ kind: 'new' })}
            className="rounded bg-blue-600 px-4 py-1.5 text-xs font-semibold text-white disabled:opacity-40" data-testid="group-create"
          >
            {t('step4.group.list.create')}
          </button>
          {removed.length > 0 && (
            <details data-testid="group-removed">
              <summary className="cursor-pointer text-xs text-slate-400">{t('step4.group.list.removed', { count: removed.length })}</summary>
              <ul className="mt-2 space-y-2">
                {removed.map((group) => <SavedRecordingCard key={group.id} group={group} highlighted={located?.groupId === group.id} current={matches(group)} onOpen={() => void open(group)}
                management={{ ...management, readonly: management.readonly || group.workspace_id !== workspaceId }} />)}
              </ul>
            </details>
          )}
        </>
      )}

      {editing && editorObserved && (
        <GroupEditor
          key={selection.kind === 'edit' ? selection.id : 'new'}
          group={selected} workspaceId={editorObserved.workspace_id} workspaceRevision={editorObserved.workspace_revision}
          locatorGuard={Boolean(locatorEntry || locatorSuspended)}
          locatorScope={editorEntry ? { workspaceId: editorEntry.workspace.id, workspaceRevision: editorEntry.entry.workspace_revision } : undefined}
          connector={connector} candidates={candidates} excluded={excluded} readonly={readonly || Boolean(locatorSuspended)} editorReadonly={selected?.status === 'deleted'}
          onClose={() => { openingRequest.current?.abort(); openingRequest.current = undefined; setOpening(false); setSelection({ kind: 'none' }); setEditorEntry(undefined); }}
          onSaved={(group) => { setSelection({ kind: 'edit', id: group.id, saved: group }); if (editorEntry) void reloadEditor(); }}
          onReload={() => { if (editorEntry) void reloadEditor(); else void groupsQuery.refetch(); }}
          fetchLatest={async () => {
            const latest = await groupsQuery.refetch();
            if (editorEntry && (latest.isError || !latest.data || latest.data.workspace_id !== editorObserved.workspace_id ||
              selected && !latest.data.groups.some(group => group.id === selected.id && group.status !== 'deleted')))
              throw new Error('Recording scope is unavailable.');
            return latest.data;
          }}
          onReviewScopedStart={onReviewScopedStart} lifecycleHost={lifecycleHost}
        />
      )}
    </section>
  );
};
