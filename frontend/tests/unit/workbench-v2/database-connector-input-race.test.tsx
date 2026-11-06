import * as React from 'react';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { act, fireEvent, render, screen, waitFor } from '@testing-library/react';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import { ConnectorSection } from '@/features/datalink/workbench-v2/steps/step4/ConnectorSection';
import { INITIAL_STATE, workbenchV2Reducer, type WorkbenchV2Action } from '@/features/datalink/workbench-v2/state/useWorkbenchV2State';
import { hydrateStudioV2DatabaseConnector } from '@/features/datalink/workbench-v2/state/studioV2DatabaseAutosave';
import { useStudioV2DatabaseAutosave } from '@/pages/datalink/workbench-v2/useStudioV2DatabaseAutosave';
import { studioV2WorkspaceKeys } from '@/hooks/datalink/keys';
import { recordingHomeSourceContext } from '@/hooks/datalink/recordingHomeSourceCache';
import { studioV2WorkspaceDatabaseAPI } from '@/services/studioV2WorkspaceDatabase';
import { savedGroupState } from '../../fixtures/writeGroupState';
import type { StudioV2RuntimeAppliedRecord } from '@/types/studioV2RuntimeApply';
import type { StudioV2WorkspaceDatabaseConfigRecord } from '@/types/datalink';

vi.mock('react-i18next', () => ({ useTranslation: () => ({ t: (key: string) => key }) }));
vi.mock('@/services/studioV2WorkspaceDatabase', () => ({ studioV2WorkspaceDatabaseAPI: {
  getConfig: vi.fn(), listTargets: vi.fn(), updateConfig: vi.fn(), upsertTarget: vi.fn(),
} }));
const record = (patch: Partial<StudioV2WorkspaceDatabaseConfigRecord> = {}): StudioV2RuntimeAppliedRecord<StudioV2WorkspaceDatabaseConfigRecord> => ({
  runtime_apply_status: 'not_running',
  id: 'db-1', workspace_id: 'ws-1', identity_revision: 'identity-1', setup_revision: 'setup-1',
  kind: 'sqlite', name: 'SQLite', host: '', port: 0, database: 'gateway.db', username: '', schema: '',
  table: 'sensor_readings', write_mode: 'insert', write_interval_seconds: 5, timestamp_column: 'ts',
  status: 'ready', save_state: 'saved', created_at: '2026-10-31T08:00:00Z', updated_at: '2026-10-31T08:00:00Z', ...patch,
});
function deferred<T>() {
  let resolve!: (value: T) => void;
  let reject!: (error: unknown) => void;
  const promise = new Promise<T>((yes, no) => { resolve = yes; reject = no; });
  return { promise, resolve, reject };
}
function Consumer({ fresh = false }: { fresh?: boolean }) {
  const [state, dispatch] = React.useReducer(workbenchV2Reducer, {
    ...savedGroupState(), db: { ...INITIAL_STATE.db, connector: { ...hydrateStudioV2DatabaseConnector(record()), ...(fresh ? {
      connector_id: undefined, identity_revision: undefined, setup_revision: undefined, workspace_id: undefined, persisted: false, save_state: 'idle' as const,
    } : {}) } },
  });
  const stateRef = React.useRef(state);
  const autosave = useStudioV2DatabaseAutosave({ dispatch }, stateRef, true, undefined, 'ws-1');
  React.useEffect(() => { stateRef.current = state; }, [state]);
  const perform = (action: WorkbenchV2Action) => {
    stateRef.current = workbenchV2Reducer(stateRef.current, action);
    dispatch(action);
    autosave.afterDatabaseAction(action, stateRef.current);
  };
  return <>
    <ConnectorSection connector={state.db.connector}
      onUpdateConnector={patch => perform({ type: 'updateDbConnector', patch })}
      onKindChange={kind => perform({ type: 'updateDbConnector', patch: { kind } })} />
    <button onClick={() => perform({ type: 'upsertDbTarget', pointId: 'rule-1-p-0', target: { tag_id: 'tag-0', column_name: 'value', enabled: true } })}>save target</button>
    <output data-testid="save-state">{state.db.connector.save_state}</output>
  </>;
}
function renderConsumer(fresh = false) {
  // Production's global retry remains on: connector mutation must explicitly opt out.
  const client = new QueryClient({ defaultOptions: { queries: { retry: false }, mutations: { retry: 1, retryDelay: 0 } } });
  render(<QueryClientProvider client={client}><Consumer fresh={fresh} /></QueryClientProvider>);
  return client;
}

beforeEach(() => {
  vi.resetAllMocks();
  vi.mocked(studioV2WorkspaceDatabaseAPI.getConfig).mockResolvedValue(record());
  vi.mocked(studioV2WorkspaceDatabaseAPI.listTargets).mockResolvedValue([]);
});

describe('connector input through real query, mutation, reducer and form consumer', () => {
  it.each([
    ['step4.field_name', 'name', 'L', 'Linux QA', ' SQLite', 'Linux QA SQLite'],
    ['step4.field_db_sqlite', 'database', '/', '/workspace/scratch/b2a0a232be39', '/linux-qa/target.sqlite', '/workspace/scratch/b2a0a232be39/linux-qa/target.sqlite'],
  ])('keeps later characters in %s across mutation refetch and ACK', async (label, field, firstText, middleText, suffix, finalText) => {
    const first = deferred<ReturnType<typeof record>>();
    const second = deferred<ReturnType<typeof record>>();
    const third = deferred<ReturnType<typeof record>>();
    vi.mocked(studioV2WorkspaceDatabaseAPI.updateConfig)
      .mockReturnValueOnce(first.promise).mockReturnValueOnce(second.promise).mockReturnValueOnce(third.promise);
    const client = renderConsumer();
    await waitFor(() => expect(studioV2WorkspaceDatabaseAPI.getConfig).toHaveBeenCalledTimes(1));
    await act(async () => {});
    const input = screen.getByRole('textbox', { name: label }) as HTMLInputElement;
    fireEvent.change(input, { target: { value: firstText } });
    await waitFor(() => expect(studioV2WorkspaceDatabaseAPI.updateConfig).toHaveBeenCalledTimes(1));
    fireEvent.change(input, { target: { value: middleText } });
    const firstSaved = record({ [field]: firstText, setup_revision: 'setup-2', updated_at: '2026-10-31T08:00:01Z' });
    vi.mocked(studioV2WorkspaceDatabaseAPI.getConfig).mockResolvedValue(firstSaved);
    await act(async () => { first.resolve(firstSaved); });
    await waitFor(() => expect(studioV2WorkspaceDatabaseAPI.updateConfig).toHaveBeenCalledTimes(2));
    // The queued write coalesces the first ACK's refresh. Explicitly inject an
    // older GET while that second write is pending to preserve the race coverage.
    await act(async () => { await client.invalidateQueries({ queryKey: studioV2WorkspaceKeys.databaseConfig() }); });
    await waitFor(() => expect(studioV2WorkspaceDatabaseAPI.getConfig).toHaveBeenCalledTimes(2));
    expect(input).toHaveValue(middleText);
    expect(screen.getByTestId('save-state')).toHaveTextContent('saving');
    // Browser typing appends to the actual controlled DOM value, not the desired string.
    for (const character of suffix) fireEvent.change(input, { target: { value: input.value + character } });
    const secondSaved = record({ [field]: middleText, setup_revision: 'setup-3', updated_at: '2026-10-31T08:00:02Z' });
    vi.mocked(studioV2WorkspaceDatabaseAPI.getConfig).mockResolvedValue(secondSaved);
    await act(async () => { second.resolve(secondSaved); });
    await waitFor(() => expect(studioV2WorkspaceDatabaseAPI.updateConfig).toHaveBeenCalledTimes(3));
    expect(input).toHaveValue(finalText);
    expect(studioV2WorkspaceDatabaseAPI.updateConfig).toHaveBeenLastCalledWith(expect.objectContaining({
      [field]: finalText, expected_setup_revision: 'setup-3', connector_id: 'db-1',
    }), { baseURL: recordingHomeSourceContext().apiOrigin });
    const finalSaved = record({ [field]: finalText, setup_revision: 'setup-4', updated_at: '2026-10-31T08:00:03Z' });
    vi.mocked(studioV2WorkspaceDatabaseAPI.getConfig).mockResolvedValue(finalSaved);
    await act(async () => { third.resolve(finalSaved); });
    await waitFor(() => expect(screen.getByTestId('save-state')).toHaveTextContent('saved'));
    expect(input).toHaveValue(finalText);
    expect(studioV2WorkspaceDatabaseAPI.updateConfig).toHaveBeenCalledTimes(3);
  });

  it('does not resend a connector mutation on an unknown transport outcome', async () => {
    vi.mocked(studioV2WorkspaceDatabaseAPI.updateConfig).mockRejectedValue(new Error('Network Error'));
    renderConsumer();
    await act(async () => {});
    fireEvent.change(screen.getByRole('textbox', { name: 'step4.field_name' }), { target: { value: 'Linux QA SQLite' } });
    await waitFor(() => expect(screen.getByTestId('save-state')).toHaveTextContent('save-error'));
    expect(studioV2WorkspaceDatabaseAPI.updateConfig).toHaveBeenCalledTimes(1);
  });
  it('holds queued and later intent after unknown outcome without resending or adopting canonical text', async () => {
    const first = deferred<ReturnType<typeof record>>();
    vi.mocked(studioV2WorkspaceDatabaseAPI.updateConfig).mockReturnValue(first.promise);
    const client = renderConsumer();
    await act(async () => {});
    const input = screen.getByRole('textbox', { name: 'step4.field_name' });
    fireEvent.change(input, { target: { value: 'L' } });
    await waitFor(() => expect(studioV2WorkspaceDatabaseAPI.updateConfig).toHaveBeenCalledTimes(1));
    fireEvent.change(input, { target: { value: 'Linux QA' } });
    await act(async () => { first.reject(new Error('Network Error')); });
    await waitFor(() => expect(screen.getByTestId('save-state')).toHaveTextContent('save-error'));
    vi.mocked(studioV2WorkspaceDatabaseAPI.getConfig).mockResolvedValue(record({ name: 'L', setup_revision: 'setup-2', updated_at: '2026-10-31T08:00:01Z' }));
    await act(async () => { await client.invalidateQueries(); });
    expect(input).toHaveValue('Linux QA');
    fireEvent.change(input, { target: { value: 'Linux QA SQLite' } });
    await act(async () => {});
    expect(input).toHaveValue('Linux QA SQLite');
    expect(screen.getByTestId('save-state')).toHaveTextContent('save-error');
    expect(studioV2WorkspaceDatabaseAPI.updateConfig).toHaveBeenCalledTimes(1);
  });

  it('keeps unknown ownership after a first create assigns the workspace to an already queued update', async () => {
    const first = deferred<ReturnType<typeof record>>();
    const second = deferred<ReturnType<typeof record>>();
    vi.mocked(studioV2WorkspaceDatabaseAPI.getConfig).mockResolvedValue(null);
    vi.mocked(studioV2WorkspaceDatabaseAPI.updateConfig).mockReturnValueOnce(first.promise).mockReturnValueOnce(second.promise);
    renderConsumer(true);
    await act(async () => {});
    const input = screen.getByRole('textbox', { name: 'step4.field_name' });
    fireEvent.change(input, { target: { value: 'L' } });
    await waitFor(() => expect(studioV2WorkspaceDatabaseAPI.updateConfig).toHaveBeenCalledTimes(1));
    fireEvent.change(input, { target: { value: 'Linux QA' } });
    const saved = record({ name: 'L' });
    vi.mocked(studioV2WorkspaceDatabaseAPI.getConfig).mockResolvedValue(saved);
    await act(async () => { first.resolve(saved); });
    await waitFor(() => expect(studioV2WorkspaceDatabaseAPI.updateConfig).toHaveBeenCalledTimes(2));
    expect(studioV2WorkspaceDatabaseAPI.updateConfig).toHaveBeenLastCalledWith(expect.objectContaining({ connector_id: 'db-1', name: 'Linux QA' }), { baseURL: recordingHomeSourceContext().apiOrigin });
    await act(async () => { second.reject(new Error('Network Error')); });
    await waitFor(() => expect(screen.getByTestId('save-state')).toHaveTextContent('save-error'));
    fireEvent.change(input, { target: { value: 'Linux QA SQLite' } });
    await act(async () => {});
    expect(input).toHaveValue('Linux QA SQLite');
    expect(studioV2WorkspaceDatabaseAPI.updateConfig).toHaveBeenCalledTimes(2);
  });

  it.each([{ id: 'wrong-connector' }, { workspace_id: 'wrong-workspace' }])('rejects an ACK bound to a different identity: %j', async mismatch => {
    const first = deferred<ReturnType<typeof record>>();
    vi.mocked(studioV2WorkspaceDatabaseAPI.updateConfig).mockReturnValue(first.promise);
    renderConsumer();
    await act(async () => {});
    const input = screen.getByRole('textbox', { name: 'step4.field_name' });
    fireEvent.change(input, { target: { value: 'L' } });
    await waitFor(() => expect(studioV2WorkspaceDatabaseAPI.updateConfig).toHaveBeenCalledTimes(1));
    fireEvent.change(input, { target: { value: 'Linux QA SQLite' } });
    const wrong = record({ ...mismatch, name: 'L', updated_at: '2026-10-31T08:00:01Z' });
    vi.mocked(studioV2WorkspaceDatabaseAPI.getConfig).mockResolvedValue(wrong);
    await act(async () => { first.resolve(wrong); });
    await waitFor(() => expect(screen.getByTestId('save-state')).toHaveTextContent('save-error'));
    expect(input).toHaveValue('Linux QA SQLite');
    expect(studioV2WorkspaceDatabaseAPI.updateConfig).toHaveBeenCalledTimes(1);
  });

  it('protects a connector edit queued behind a target write from a late config GET', async () => {
    const target = deferred<Awaited<ReturnType<typeof studioV2WorkspaceDatabaseAPI.upsertTarget>>>();
    vi.mocked(studioV2WorkspaceDatabaseAPI.upsertTarget).mockReturnValue(target.promise);
    const client = renderConsumer();
    await act(async () => {});
    fireEvent.click(screen.getByRole('button', { name: 'save target' }));
    await waitFor(() => expect(studioV2WorkspaceDatabaseAPI.upsertTarget).toHaveBeenCalledTimes(1));
    const input = screen.getByRole('textbox', { name: 'step4.field_name' });
    fireEvent.change(input, { target: { value: 'Linux QA SQLite' } });
    expect(studioV2WorkspaceDatabaseAPI.updateConfig).not.toHaveBeenCalled();
    vi.mocked(studioV2WorkspaceDatabaseAPI.getConfig).mockResolvedValue(record({ updated_at: '2026-10-31T08:00:01Z' }));
    await act(async () => { await client.invalidateQueries(); });
    await waitFor(() => expect(studioV2WorkspaceDatabaseAPI.getConfig).toHaveBeenCalledTimes(2));
    expect(input).toHaveValue('Linux QA SQLite');
    expect(screen.getByTestId('save-state')).toHaveTextContent('saving');
  });

});
