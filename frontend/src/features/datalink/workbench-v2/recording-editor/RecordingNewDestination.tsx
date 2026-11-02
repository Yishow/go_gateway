import { useEffect, useRef, useState } from 'react';
import { useQuery, useQueryClient } from '@tanstack/react-query';
import { useTranslation } from 'react-i18next';
import { dbTargetAPI } from '../../../../services/datalink';
import { dbTargetKeys } from '../../../../hooks/datalink/keys';
import type { DatabaseConnector } from '../../../../types/datalink';
import type { DbConnector } from '../state/types';
import { safeRuntimeLabel } from '../../../../utils/runtimeOperatorSafety';

/** Connection identity only; no credentials, endpoint, inherited table or operation evidence enters a draft. */
function destination(value: DatabaseConnector, workspaceId: string): DbConnector | undefined {
  if (!value.enabled || !value.id?.trim() || !value.identity_revision?.trim() ||
      !['sqlite', 'postgres', 'mysql', 'sqlserver'].includes(value.kind)) return undefined;
  const config = value.connection_config ?? {};
  const text = (...keys: string[]): string | undefined => {
    for (const key of keys) { const item = config[key]; if (typeof item === 'string' && item.trim()) return item; }
    return undefined;
  };
  const database = value.kind === 'sqlite' ? text('database', 'path') : text('database', 'dbname');
  if (!database) return undefined;
  return { connector_id: value.id, identity_revision: value.identity_revision, workspace_id: workspaceId,
    kind: value.kind, name: safeRuntimeLabel(value.name), database, schema: text('schema', 'table_schema') ?? '', table: '',
    host: '', port: 0, username: '', write_mode: 'insert', write_interval_seconds: 5, timestamp_column: 'ts',
    status: value.status, persisted: true, save_state: 'saved' };
}
type Owner = { workspaceId: string; enabled: boolean };
type Choice = { owner: Owner; connector: DbConnector };
export interface NewRecordingDestination {
  connector: DbConnector;
  options: readonly DbConnector[];
  initialSetup: boolean;
  status?: 'loading' | 'error' | 'empty' | 'stale';
  current: () => boolean;
  choose: (id: string) => DbConnector | undefined;
  refresh: () => void;
}
/** Local selection is not authority: every use rechecks the existing connector query and original owner. */
export function useRecordingNewDestination(enabled: boolean, workspaceId: string, fallback: DbConnector): NewRecordingDestination | undefined {
  const client = useQueryClient();
  const query = useQuery({ queryKey: dbTargetKeys.connectors(), queryFn: () => dbTargetAPI.listConnectors(),
    enabled, retry: false, refetchOnWindowFocus: false, notifyOnChangeProps: 'all' });
  const owner = useRef<Owner>({ workspaceId, enabled });
  if (owner.current.workspaceId !== workspaceId || owner.current.enabled !== enabled) owner.current = { workspaceId, enabled };
  const renderOwner = owner.current;
  const mounted = useRef(false);
  useEffect(() => { mounted.current = true; return () => { mounted.current = false; }; }, []);
  const [choice, setChoice] = useState<Choice>();
  const liveChoice = useRef(choice); liveChoice.current = choice;
  const initialSetupOwner = useRef<Owner | undefined>(undefined);
  const read = () => {
    const state = client.getQueryState<DatabaseConnector[]>(dbTargetKeys.connectors());
    return state?.status === 'success' && state.fetchStatus === 'idle' && !state.isInvalidated && Array.isArray(state.data)
      ? state.data.flatMap(value => { const next = destination(value, workspaceId); return next ? [next] : []; }) : undefined;
  };
  const options = read();
  if (enabled && options && client.getQueryState<DatabaseConnector[]>(dbTargetKeys.connectors())?.data?.length === 0 &&
      !fallback.connector_id && !fallback.persisted) initialSetupOwner.current = renderOwner;
  const initialSetup = initialSetupOwner.current === renderOwner;
  const activeChoice = choice?.owner === renderOwner ? choice.connector : undefined;
  const expected = activeChoice ?? fallback;
  const matches = (value: DbConnector) => value.connector_id === expected.connector_id && value.identity_revision === expected.identity_revision;
  const observed = options?.find(matches);
  const connector = initialSetup ? fallback : observed ?? { ...expected, password: undefined, persisted: false, save_state: 'draft-invalid' as const };
  const current = () => enabled && mounted.current && owner.current === renderOwner && liveChoice.current === choice &&
    (initialSetup || Boolean(read()?.some(matches)));
  const status = initialSetup ? undefined : query.isError ? 'error' : query.isFetching || query.isPending ? 'loading'
    : !options ? 'stale' : !observed ? (activeChoice ? 'stale' : options.length ? undefined : 'empty') : undefined;
  return enabled ? { connector, options: options ?? [], initialSetup, status, current,
    choose: id => {
      if (!mounted.current || owner.current !== renderOwner || initialSetup) return undefined;
      const offered = options?.find(value => value.connector_id === id);
      const exact = offered && read()?.find(value => value.connector_id === id && value.identity_revision === offered.identity_revision);
      if (!exact) return undefined;
      const next = { owner: renderOwner, connector: exact };
      liveChoice.current = next; setChoice(next); return exact;
    }, refresh: () => { if (owner.current === renderOwner && mounted.current) void query.refetch(); },
  } : undefined;
}
export function RecordingNewDestinationChooser({ value, selectedId, disabled, onSelect }: {
  value: NewRecordingDestination; selectedId: string; disabled: boolean; onSelect: (id: string) => void;
}) {
  const { t } = useTranslation('recording-editor');
  return <div className="space-y-2" data-testid="recording-new-destination">
    <label className="block text-sm">{t('newDestination.label')}
      <select className="mt-1 w-full rounded border border-[var(--gw-control)] bg-[var(--gw-surface)] px-2 py-2"
        value={value.status === 'stale' ? '' : selectedId} disabled={disabled || value.status === 'loading' || value.status === 'error'} onChange={event => onSelect(event.target.value)}>
        <option value="">{t('newDestination.choose')}</option>
        {value.options.map(option => <option key={option.connector_id} value={option.connector_id}>{option.name || option.kind}</option>)}
      </select>
    </label>
    <p className="text-xs text-[var(--gw-muted)]">{t('newDestination.scope')}</p>
    {value.status && <p data-testid="new-destination-status" role={value.status === 'loading' ? 'status' : 'alert'}>{t(`newDestination.${value.status}`)}</p>}
    {value.status && value.status !== 'loading' && <button type="button" className="text-xs underline" disabled={disabled} onClick={value.refresh}>{t('newDestination.refresh')}</button>}
  </div>;
}
/** Existing mismatch presentation, kept separate so the editor retains its original draft/save owner. */
export function RecordingDestinationMismatch({ locatorUnknown, saved, readonly, onUseCurrent }: {
  locatorUnknown: boolean; saved: boolean; readonly: boolean; onUseCurrent: () => void;
}) {
  const { t } = useTranslation('workbench-v2');
  return <div role="alert" className="space-y-1 text-xs text-[var(--gw-warning-fg)]" data-testid="group-destination-mismatch">
    <p>{locatorUnknown ? t('runtime-dashboard:locatorEntry.destinationUnknown') : t('step4.group.editor.destination_mismatch')}</p>
    {saved && !readonly && <button type="button" className="underline" data-testid="group-use-current-destination" onClick={onUseCurrent}>
      {t('step4.group.editor.use_current_destination')}
    </button>}
  </div>;
}
