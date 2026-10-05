import * as React from 'react';
import { useEffect, useMemo, useRef, useState } from 'react';
import { useTranslation } from 'react-i18next';
import {
  useCreateWriteGroupMutation,
  useUpdateWriteGroupMutation,
  useWriteGroupReadinessQuery,
} from '../../../../../../hooks/datalink/useStudioV2WriteGroups';
import type { WriteGroup, WriteGroupStorageStrategy } from '../../../../../../types/studioV2WriteGroup';
import { normalizeTypedEnvelope } from '../../../../../../utils/safeJson';
import { getSafeErrorMessage } from '../../../../../../utils/typedErrors';
import type { DbConnector } from '../../../state/types';
import type { ExcludedCandidate, GroupCandidate } from '../../../state/writeGroup/candidates';
import { suggestColumns } from '../../../state/writeGroup/columns';
import {
  draftIsDirty, draftToRequestGroup, groupToDraft, newDraft, validateDraft,
  type EditorDraft, type IncompletePolicy,
} from '../../../state/writeGroup/draft';
import { useStep4TargetColumns } from '../useStep4TargetColumns';
import { GroupColumnProposal } from './GroupColumnProposal';
import { GroupDeliveryStrip } from './GroupDeliveryStrip';
import { GroupLifecycleBar } from './GroupLifecycleBar';
import { GroupMemberTable } from './GroupMemberTable';
import { GroupReadinessPanel } from './GroupReadinessPanel';
import { GroupSchemaPanel } from './GroupSchemaPanel';
import { GroupTestWritePanel } from './GroupTestWritePanel';

export interface GroupEditorProps {
  /** Undefined creates a new group. */
  group: WriteGroup | undefined;
  workspaceId: string;
  workspaceRevision: string;
  connector: DbConnector;
  candidates: GroupCandidate[];
  excluded: ExcludedCandidate[];
  readonly: boolean;
  onClose: () => void;
  onSaved: (group: WriteGroup) => void;
  onReload: () => void;
  /**
   * Reads the current workspace revision and group just before a save. Other
   * setup saves (devices, rules, Tags) move the workspace revision without
   * touching this group, so a stale copy would be refused for no real conflict.
   */
  fetchLatest: () => Promise<{ workspace_revision: string; groups: WriteGroup[] } | undefined>;
}

const FIELD = 'mt-1 w-full rounded border border-slate-700 bg-slate-950 px-2 py-1 text-sm text-slate-100';

function destinationSaved(connector: DbConnector): boolean {
  return connector.persisted === true && connector.save_state === 'saved' &&
    Boolean(connector.connector_id && connector.identity_revision);
}

/**
 * One saved group: its name, destination table, row rules and members. The
 * same editor serves managed and custom storage; the backend stays the
 * authority for readiness, so Apply follows the server's answer for the saved
 * revision and never the local checks alone.
 */
export const GroupEditor: React.FC<GroupEditorProps> = ({
  group: propGroup, workspaceId, workspaceRevision, connector, candidates, excluded, readonly, onClose, onSaved, onReload, fetchLatest,
}) => {
  const { t } = useTranslation('workbench-v2');
  const saved = destinationSaved(connector);
  // Keep the just-returned canonical group visible while the list query is
  // refetching. This prevents stale list data from making server-assigned
  // managed table/columns look dirty again.
  const [editorGroup, setEditorGroup] = useState<WriteGroup | undefined>(propGroup);
  const previousPropGroup = useRef(propGroup);
  const group = editorGroup;
  const [draft, setDraft] = useState<EditorDraft>(() => group ? groupToDraft(group) : newDraft({
    connector_id: connector.connector_id ?? '', connector_revision: connector.identity_revision ?? '',
    table_schema: connector.schema ?? '', table_name: connector.table ?? '',
  }));
  const [saveError, setSaveError] = useState<unknown>(null);
  const [conflict, setConflict] = useState(false);
  const createMutation = useCreateWriteGroupMutation();
  const updateMutation = useUpdateWriteGroupMutation();
  const savingRef = useRef(false);
  const customDedupeRef = useRef(draft.dedupe_capability);
  const mountedRef = useRef(false);
  useEffect(() => {
    mountedRef.current = true;
    return () => { mountedRef.current = false; };
  }, []);

  const dirty = group ? draftIsDirty(draft, group) : true;
  // Adopt a new server revision only while there is nothing unsaved to lose.
  const revision = group?.revision;
  const dirtyRef = useRef(dirty);
  dirtyRef.current = dirty;
  useEffect(() => {
    if (previousPropGroup.current !== propGroup && !dirtyRef.current) {
      previousPropGroup.current = propGroup;
      setEditorGroup(propGroup);
    }
  }, [propGroup]);
  useEffect(() => {
    if (group && !dirtyRef.current) {
      const next = groupToDraft(group);
      customDedupeRef.current = next.dedupe_capability;
      setDraft(next);
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [revision]);

  const destinationMatches = draft.connector_id === connector.connector_id &&
    draft.connector_revision === connector.identity_revision;
  const scopedConnector = useMemo<DbConnector>(() => ({
    ...connector, table: draft.table_name, schema: draft.table_schema,
  }), [connector, draft.table_name, draft.table_schema]);
  const tableSaved = group
    ? draft.table_name === group.destination.table_name && draft.table_schema === group.destination.table_schema
    : draft.table_name === connector.table && draft.table_schema === connector.schema;
  const metadata = useStep4TargetColumns(destinationMatches && tableSaved && (draft.storage_strategy === 'custom' || Boolean(group)) ? scopedConnector : { ...scopedConnector, persisted: false },
    group ? { groupId: group.id, groupRevision: group.revision } : undefined);
  // A new managed group has no target table to inspect. Present an empty
  // server-owned column view so the editor never asks for manual assignments.
  const editorMetadata = draft.storage_strategy === 'managed' && !group
    ? { ...metadata, status: 'exists' as const, columns: [] }
    : metadata;
  const dialect = connector.kind === 'postgres' ? 'postgres' : connector.kind === 'sqlite' ? 'sqlite' : 'portable';
  const confirmed = useMemo(
    () => Object.fromEntries(draft.members.filter((m) => m.target_column).map((m) => [m.key, m.target_column])),
    [draft.members],
  );
  const suggestions = useMemo(
    () => (editorMetadata.status === 'exists' ? suggestColumns(candidates, editorMetadata.columns, confirmed, dialect) : { assignments: {}, unmatched: candidates.map((c) => c.key) }),
    [candidates, editorMetadata.columns, editorMetadata.status, confirmed, dialect],
  );
  const issues = validateDraft(draft, { destinationSaved: saved && destinationMatches, storageStrategy: draft.storage_strategy });
  const readiness = useWriteGroupReadinessQuery(group?.id, Boolean(group) && !dirty);
  const set = (patch: Partial<EditorDraft>) => setDraft((current) => ({ ...current, ...patch }));
  const changeStorageStrategy = (storage_strategy: WriteGroupStorageStrategy) => setDraft((current) => {
    if (storage_strategy === current.storage_strategy) return current;
    if (storage_strategy === 'managed') {
      customDedupeRef.current = current.dedupe_capability;
      return { ...current, storage_strategy, dedupe_capability: 'receipt' };
    }
    return { ...current, storage_strategy, dedupe_capability: customDedupeRef.current };
  });

  const handleSave = async () => {
    if (savingRef.current || issues.length > 0 || readonly) return;
    savingRef.current = true;
    setSaveError(null);
    try {
      const body = draftToRequestGroup(draft, workspaceId);
      const latest = await fetchLatest();
      const freshRevision = latest?.workspace_revision ?? workspaceRevision;
      if (group) {
        const current = latest?.groups.find((entry) => entry.id === group.id);
        if (current && current.revision !== group.revision) {
          // This group itself changed elsewhere: that is a real conflict.
          setConflict(true);
          return;
        }
      }
      const result = group
        ? await updateMutation.mutateAsync({ id: group.id, request: {
          workspace_id: workspaceId, expected_workspace_revision: freshRevision, expected_group_revision: group.revision,
          expected_connector_revision: draft.connector_revision, group: body,
        } })
        : await createMutation.mutateAsync({
          workspace_id: workspaceId, expected_workspace_revision: freshRevision,
          expected_connector_revision: draft.connector_revision, group: body,
      });
      if (mountedRef.current) {
        // Adopt every server-assigned managed column and binding before the
        // parent refetches. Server proof stays outside the editable draft.
        setDraft(groupToDraft(result.group));
        setEditorGroup(result.group);
        setConflict(false);
        onSaved(result.group);
      }
    } catch (cause) {
      setSaveError(cause);
      if (normalizeTypedEnvelope(cause).code === 'revision_mismatch') setConflict(true);
    } finally {
      savingRef.current = false;
    }
  };

  const saving = createMutation.isPending || updateMutation.isPending;
  const canApply = Boolean(group) && !dirty && issues.length === 0 && readiness.data?.ready === true &&
    readiness.data.group_revision === group?.revision;
  const applyBlocked = !group ? t('step4.group.apply_blocked_unsaved')
    : dirty ? t('step4.group.apply_blocked_dirty')
      : readiness.isLoading ? t('step4.group.apply_blocked_loading')
        : readiness.isError ? t('step4.group.apply_blocked_unknown')
          : readiness.data && !readiness.data.ready ? t('step4.group.apply_blocked_not_ready') : undefined;
  const testBlocked = !group ? t('step4.group.test_write.blocked_unsaved')
    : dirty ? t('step4.group.test_write.blocked_dirty')
      : !destinationMatches ? t('step4.group.test_write.blocked_destination') : undefined;

  return (
    <div className="space-y-5 rounded-xl border border-slate-800 p-5" data-testid="group-editor">
      <div className="flex items-center justify-between">
        <h3 className="text-sm font-semibold text-slate-100">{group ? group.name : t('step4.group.editor.new_title')}</h3>
        <button type="button" onClick={onClose} className="text-xs text-slate-400 underline" data-testid="group-editor-close">{t('step4.group.editor.back')}</button>
      </div>

      {group && (
        <dl className="grid grid-cols-2 gap-x-4 gap-y-1 text-xs text-slate-400 sm:grid-cols-4" data-testid="group-revisions">
          <dt>{t('step4.group.editor.status')}</dt><dd className="text-slate-200">{group.status}</dd>
          <dt>{t('step4.group.editor.draft_revision')}</dt><dd className="font-mono text-slate-200">{group.revision}</dd>
          <dt>{t('step4.group.editor.applied_revision')}</dt><dd className="font-mono text-slate-200" data-testid="group-applied-revision">{group.applied_revision || t('step4.group.editor.not_applied')}</dd>
          <dt>{t('step4.group.editor.connector_revision')}</dt><dd className="font-mono text-slate-200">{group.destination.connector_revision}</dd>
        </dl>
      )}

      {!destinationMatches && (
        <div role="alert" className="space-y-1 text-xs text-amber-300" data-testid="group-destination-mismatch">
          <p>{t('step4.group.editor.destination_mismatch')}</p>
          {saved && !readonly && (
            <button
              type="button" className="underline" data-testid="group-use-current-destination"
              onClick={() => set({ connector_id: connector.connector_id ?? '', connector_revision: connector.identity_revision ?? '', table_schema: connector.schema ?? '', table_name: connector.table ?? '' })}
            >
              {t('step4.group.editor.use_current_destination')}
            </button>
          )}
        </div>
      )}

      <div className="grid gap-3 sm:grid-cols-2">
        <label className="text-xs text-slate-400">
          {t('step4.group.editor.name')}
          <input type="text" value={draft.name} disabled={readonly} onChange={(event) => set({ name: event.target.value })} className={FIELD} data-testid="group-name" />
        </label>
        <fieldset className="text-xs text-slate-400">
          <legend>{t('step4.group.editor.storage')}</legend>
          {(['custom', 'managed'] as WriteGroupStorageStrategy[]).map((strategy) => (
            <label key={strategy} className="mr-4 inline-flex items-center gap-1 text-slate-200">
              <input type="radio" name="group-storage" value={strategy} checked={draft.storage_strategy === strategy} disabled={readonly} onChange={() => changeStorageStrategy(strategy)} data-testid={`group-storage-${strategy}`} />
              {t(`step4.group.editor.storage_${strategy}`)}
            </label>
          ))}
          <p className="mt-1 text-slate-500" data-testid="group-storage-note">{t(`step4.group.editor.storage_${draft.storage_strategy}_note`)}</p>
        </fieldset>
        <label className="text-xs text-slate-400">
          {t('step4.group.editor.table')}
          <input type="text" value={draft.table_name} disabled={readonly} onChange={(event) => set({ table_name: event.target.value })} className={FIELD} data-testid="group-table" />
          {draft.storage_strategy === 'managed' && <span className="mt-1 block text-[11px] text-slate-500">{t('step4.group.editor.table_managed_note')}</span>}
        </label>
        <label className="text-xs text-slate-400">
          {t('step4.group.editor.schema')}
          <input type="text" value={draft.table_schema} disabled={readonly} onChange={(event) => set({ table_schema: event.target.value })} className={FIELD} data-testid="group-schema" />
        </label>
      </div>

      <details className="rounded border border-slate-800 p-3" data-testid="group-advanced">
        <summary className="cursor-pointer text-xs font-medium text-slate-300">{t('step4.group.editor.advanced')}</summary>
        <div className="mt-3 grid gap-3 sm:grid-cols-2">
          <label className="text-xs text-slate-400">{t('step4.group.editor.interval')}
            <input type="number" min={1} value={draft.interval_seconds} disabled={readonly} onChange={(event) => set({ interval_seconds: Number(event.target.value) })} className={FIELD} data-testid="group-interval" />
          </label>
          <label className="text-xs text-slate-400">{t('step4.group.editor.lateness')}
            <input type="number" min={0} value={draft.allowed_lateness_seconds} disabled={readonly} onChange={(event) => set({ allowed_lateness_seconds: Number(event.target.value) })} className={FIELD} data-testid="group-lateness" />
          </label>
          <label className="text-xs text-slate-400">{t('step4.group.editor.incomplete')}
            <select value={draft.incomplete_policy} disabled={readonly} onChange={(event) => set({ incomplete_policy: event.target.value as IncompletePolicy })} className={FIELD} data-testid="group-incomplete">
              <option value="skip_row">{t('step4.group.editor.incomplete_skip')}</option>
              <option value="partial">{t('step4.group.editor.incomplete_partial')}</option>
            </select>
          </label>
          <label className="text-xs text-slate-400">{t('step4.group.editor.dedupe')}
            <select value={draft.storage_strategy === 'managed' ? 'receipt' : draft.dedupe_capability} disabled={readonly || draft.storage_strategy === 'managed'} onChange={(event) => { customDedupeRef.current = event.target.value; set({ dedupe_capability: event.target.value }); }} className={FIELD} data-testid="group-dedupe">
              {draft.storage_strategy !== 'managed' && <option value="">{t('step4.group.editor.dedupe_none')}</option>}
              <option value="receipt">{t('step4.group.editor.dedupe_receipt')}</option>
            </select>
          </label>
          <label className="text-xs text-slate-400">{t('step4.group.editor.entity_column')}
            <input type="text" value={draft.entity_key_column} disabled={readonly} onChange={(event) => set({ entity_key_column: event.target.value })} className={FIELD} data-testid="group-entity-column" />
          </label>
          <label className="text-xs text-slate-400">{t('step4.group.editor.provenance_column')}
            <input type="text" value={draft.provenance_column} disabled={readonly} onChange={(event) => set({ provenance_column: event.target.value })} className={FIELD} data-testid="group-provenance-column" />
          </label>
        </div>
      </details>

      <GroupMemberTable
        dialect={dialect}
        candidates={candidates} excluded={excluded} members={draft.members} columns={editorMetadata.columns} metadataStatus={editorMetadata.status}
        suggestions={suggestions} issues={issues} showEntityKey={draft.entity_key_column.trim() !== ''} disabled={readonly}
        onChange={(members) => set({ members })}
      />
      <GroupColumnProposal
        dialect={dialect}
        managed={draft.storage_strategy === 'managed'}
        existingColumns={editorMetadata.columns.map((column) => column.name)}
        unmatched={candidates.filter((candidate) => {
          const member = draft.members.find((entry) => entry.key === candidate.key);
          if (!member) return false;
          const column = metadata.columns.find((entry) => entry.name === member.target_column);
          return editorMetadata.status === 'exists' && (!member.target_column || !column);
        })}
      />
      {editorMetadata.status === 'failed' && (
        <button type="button" onClick={metadata.refetch} className="text-xs underline" data-testid="group-metadata-retry">{t('step4.group.retry')}</button>
      )}

      {draft.storage_strategy === 'managed' && (
        <GroupSchemaPanel
          group={group}
          workspaceId={workspaceId}
          workspaceRevision={workspaceRevision}
          destinationMatches={destinationMatches}
          dirty={dirty}
          readonly={readonly}
          onApplied={() => {
            metadata.refetch();
            void readiness.refetch();
            onReload();
          }}
        />
      )}

      {issues.length > 0 && (
        <ul role="alert" className="list-disc pl-5 text-xs text-red-300" data-testid="group-issues">
          {issues.map((issue, index) => <li key={`${issue.code}-${issue.member_key ?? issue.column ?? index}`}>{t(`step4.group.issue.${issue.code.replace(/-/g, '_')}`, { column: issue.column ?? '' })}</li>)}
        </ul>
      )}
      {conflict && (
        <div role="alert" className="text-xs text-amber-300" data-testid="group-conflict">
          {t('step4.group.editor.conflict')}
          <button type="button" onClick={() => { setConflict(false); onReload(); }} className="ml-2 underline" data-testid="group-conflict-reload">{t('step4.group.editor.reload')}</button>
        </div>
      )}
      {saveError !== null && !conflict && (
        <p role="alert" className="text-xs text-red-300" data-testid="group-save-error">{getSafeErrorMessage(saveError, t).message}</p>
      )}

      <div className="flex items-center gap-3">
        <button
          type="button" onClick={() => void handleSave()} disabled={readonly || saving || issues.length > 0 || (Boolean(group) && !dirty)}
          className="rounded bg-blue-600 px-4 py-1.5 text-xs font-semibold text-white disabled:opacity-40" data-testid="group-save"
        >
          {group ? t('step4.group.editor.save') : t('step4.group.editor.create')}
        </button>
        {group && !dirty && <span className="text-xs text-emerald-300" data-testid="group-saved-note">{t('step4.group.editor.saved')}</span>}
        {dirty && group && <span className="text-xs text-amber-300" data-testid="group-dirty-note">{t('step4.group.editor.unsaved')}</span>}
      </div>

      {group && (
        <>
          <section aria-label={t('step4.group.readiness.title')} className="space-y-2">
            <h4 className="text-sm font-semibold text-slate-200">{t('step4.group.readiness.title')}</h4>
            <GroupReadinessPanel readiness={readiness.data} loading={readiness.isLoading && !dirty} failed={readiness.isError} dirty={dirty} onRetry={() => void readiness.refetch()} />
          </section>
          <GroupDeliveryStrip group={group} readonly={readonly} />
          <GroupLifecycleBar
            group={group} workspaceId={workspaceId} workspaceRevision={workspaceRevision}
            canApply={canApply} applyBlockedReason={applyBlocked} readonly={readonly}
            onChanged={onSaved} onConflict={() => setConflict(true)}
          />
          <GroupTestWritePanel
            groupId={group.id} scopeKey={`${group.id}|${group.revision}|${workspaceRevision}`}
            enabled={!readonly && !dirty && destinationMatches && group.status !== 'deleted'} disabledReason={readonly ? t('step4.group.test_write.blocked_readonly') : testBlocked}
          />
        </>
      )}
    </div>
  );
};
