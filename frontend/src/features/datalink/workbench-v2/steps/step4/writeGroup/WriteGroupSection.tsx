import * as React from 'react';
import { useMemo, useState } from 'react';
import { useTranslation } from 'react-i18next';
import { useWriteGroupDeliveryQuery, useWriteGroupsQuery } from '../../../../../../hooks/datalink/useStudioV2WriteGroups';
import type { WriteGroup } from '../../../../../../types/studioV2WriteGroup';
import type { WorkbenchV2State } from '../../../state/types';
import { buildGroupCandidates } from '../../../state/writeGroup/candidates';
import type { WorkspaceReadinessStepNumber } from '../../../components/WorkspaceReadinessPanel';
import { GroupEditor } from './GroupEditor';

export interface WriteGroupSectionProps {
  state: WorkbenchV2State;
  workspaceId?: string;
  readonly: boolean;
  onNavigateStep?: (step: WorkspaceReadinessStepNumber) => void;
}

type Selection = { kind: 'none' } | { kind: 'new' } | { kind: 'edit'; id: string; saved?: WriteGroup };

type Prerequisite = 'device' | 'tags' | 'destination' | 'table' | null;

function GroupListItem({ group, current, onOpen }: { group: WriteGroup; current: boolean; onOpen: () => void }) {
  const { t } = useTranslation('workbench-v2');
  const delivery = useWriteGroupDeliveryQuery(group.id, true);
  return (
    <li className="flex flex-wrap items-center justify-between gap-2 rounded border border-slate-800 p-3 text-xs text-slate-300" data-testid={`group-list-item-${group.id}`}>
      <div>
        <div className="text-sm font-medium text-slate-100">{group.name}</div>
        <div className="text-slate-500">
          {group.destination.table_name} · {t(`step4.group.status.${group.status}`)} · {group.members.length} {t('step4.group.list.members')}
          {!current && <span className="ml-2 text-amber-300" data-testid={`group-list-mismatch-${group.id}`}>{t('step4.group.list.other_destination')}</span>}
        </div>
        {delivery.data && (
          <div className="text-slate-500" data-testid={`group-list-delivery-${group.id}`}>
            {t('step4.group.list.delivery', { committed: delivery.data.stages.sql_committed, waiting: delivery.data.stages.queued + delivery.data.stages.retrying })}
          </div>
        )}
      </div>
      <button type="button" onClick={onOpen} className="rounded border border-slate-600 px-3 py-1 text-slate-100" data-testid={`group-open-${group.id}`}>
        {group.status === 'deleted' ? t('step4.group.list.view') : t('step4.group.list.edit')}
      </button>
    </li>
  );
}

/**
 * The only database output editor of Step 4. It reads the saved groups of the
 * workspace and says what is missing before anything can be created, keeping a
 * failed load, an empty list and groups for another destination apart.
 */
export const WriteGroupSection: React.FC<WriteGroupSectionProps> = ({ state, workspaceId, readonly, onNavigateStep }) => {
  const { t } = useTranslation('workbench-v2');
  const [selection, setSelection] = useState<Selection>({ kind: 'none' });
  const groupsQuery = useWriteGroupsQuery(Boolean(workspaceId));
  const { candidates, excluded } = useMemo(() => buildGroupCandidates(state), [state]);
  const connector = state.db.connector;
  const destinationSaved = connector.persisted === true && connector.save_state === 'saved' &&
    Boolean(connector.connector_id && connector.identity_revision);

  const prerequisite: Prerequisite = state.devices.length === 0 ? 'device'
    : candidates.length === 0 ? 'tags'
      : !destinationSaved ? 'destination'
        : !connector.table.trim() ? 'table' : null;

  const groups = groupsQuery.data?.groups ?? [];
  const live = groups.filter((group) => group.status !== 'deleted');
  const removed = groups.filter((group) => group.status === 'deleted');
  // Right after a save the refetched list may not carry the group yet; the returned group bridges that gap.
  const selected = selection.kind === 'edit' ? groups.find((group) => group.id === selection.id) ?? selection.saved : undefined;
  const matches = (group: WriteGroup) => group.destination.connector_id === connector.connector_id &&
    group.destination.connector_revision === connector.identity_revision;

  if (groupsQuery.isError) {
    return (
      <section role="alert" className="space-y-2 rounded-xl border border-amber-700 p-5 text-xs text-amber-200" data-testid="group-section-error">
        <p>{t('step4.group.section.load_failed')}</p>
        <button type="button" onClick={() => void groupsQuery.refetch()} className="underline" data-testid="group-section-retry">{t('step4.group.retry')}</button>
      </section>
    );
  }
  if (groupsQuery.isLoading) {
    return <section className="p-5 text-xs text-slate-400" role="status" data-testid="group-section-loading">{t('step4.group.section.loading')}</section>;
  }

  const editing = selection.kind === 'new' || (selection.kind === 'edit' && selected);
  return (
    <section className="space-y-4 rounded-xl border border-slate-800 bg-slate-900/20 p-5" aria-labelledby="group-section-title" data-testid="group-section">
      <div>
        <h3 id="group-section-title" className="text-sm font-semibold text-slate-100">{t('step4.group.section.title')}</h3>
        <p className="mt-1 text-xs text-slate-400">{t('step4.group.section.subtitle')}</p>
      </div>

      {prerequisite && (
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
              {live.map((group) => <GroupListItem key={group.id} group={group} current={matches(group)} onOpen={() => setSelection({ kind: 'edit', id: group.id })} />)}
            </ul>
          )}
          {live.length > 0 && !live.some(matches) && destinationSaved && (
            <p className="text-xs text-amber-300" role="status" data-testid="group-none-match">{t('step4.group.list.none_match')}</p>
          )}
          <button
            type="button" disabled={readonly || prerequisite !== null || !groupsQuery.data}
            onClick={() => setSelection({ kind: 'new' })}
            className="rounded bg-blue-600 px-4 py-1.5 text-xs font-semibold text-white disabled:opacity-40" data-testid="group-create"
          >
            {t('step4.group.list.create')}
          </button>
          {removed.length > 0 && (
            <details data-testid="group-removed">
              <summary className="cursor-pointer text-xs text-slate-400">{t('step4.group.list.removed', { count: removed.length })}</summary>
              <ul className="mt-2 space-y-2">
                {removed.map((group) => <GroupListItem key={group.id} group={group} current={matches(group)} onOpen={() => setSelection({ kind: 'edit', id: group.id })} />)}
              </ul>
            </details>
          )}
        </>
      )}

      {editing && groupsQuery.data && (
        <GroupEditor
          key={selection.kind === 'edit' ? selection.id : 'new'}
          group={selected} workspaceId={groupsQuery.data.workspace_id} workspaceRevision={groupsQuery.data.workspace_revision}
          connector={connector} candidates={candidates} excluded={excluded} readonly={readonly || selected?.status === 'deleted'}
          onClose={() => setSelection({ kind: 'none' })}
          onSaved={(group) => setSelection({ kind: 'edit', id: group.id, saved: group })}
          onReload={() => void groupsQuery.refetch()}
        />
      )}
    </section>
  );
};
