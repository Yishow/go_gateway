import { beforeEach, describe, expect, it, vi } from 'vitest';
import { act, renderHook, waitFor } from '@testing-library/react';
import type * as React from 'react';
import { useStudioV2DatabaseAutosave } from '../../../src/pages/datalink/workbench-v2/useStudioV2DatabaseAutosave';
import { INITIAL_STATE, workbenchV2Reducer, type WorkbenchV2Action } from '../../../src/features/datalink/workbench-v2/state/useWorkbenchV2State';
import type { Mapping, WorkbenchV2State } from '../../../src/features/datalink/workbench-v2/state/types';

const mocks = vi.hoisted(() => ({ updateConfig: vi.fn(), upsertTarget: vi.fn() }));
vi.mock('react-i18next', () => ({ useTranslation: () => ({ t: (key: string) => key }) }));
vi.mock('../../../src/hooks/datalink/useStudioV2WorkspaceDatabase', () => ({
  useStudioV2DatabaseConfigQuery: () => ({ isSuccess: false, data: undefined }),
  useStudioV2DatabaseTargetsQuery: () => ({ isSuccess: false, data: undefined }),
  useUpdateStudioV2DatabaseConfigMutation: () => ({ mutateAsync: mocks.updateConfig }),
  useUpsertStudioV2DatabaseTargetMutation: () => ({ mutateAsync: mocks.upsertTarget }),
}));

interface Deferred<T> { promise: Promise<T>; resolve: (value: T) => void }
function deferred<T>(): Deferred<T> {
  let resolve!: (value: T) => void;
  const promise = new Promise<T>((next) => { resolve = next; });
  return { promise, resolve };
}

const timestamp = '2026-09-16T00:00:00Z';
const configRecord = (setupRevision: string) => ({
  id: 'db-1', identity_revision: 'identity-1', setup_revision: setupRevision, workspace_id: 'ws-1',
  kind: 'sqlite', name: 'Line A', host: '', port: 0, database: '/tmp/line-a.db', username: '',
  schema: 'main', table: 'sensor_values', write_mode: 'insert', write_interval_seconds: 5,
  timestamp_column: 'ts', status: 'ready', row_groups: [], created_at: timestamp, updated_at: timestamp,
});
const targetRecord = (pointId: string, setupRevision: string) => ({
  id: `row-${pointId}`, workspace_id: 'ws-1', point_id: pointId, tag_id: `tag-${pointId}`,
  column_name: `column_${pointId}`, enabled: true, setup_revision: setupRevision, created_at: timestamp, updated_at: timestamp,
});
const mapping = (pointId: string) => ({
  point_id: pointId, tag_key: `line.${pointId}`, tag_id: `tag-${pointId}`, display_name: pointId, unit: '',
  target_type: 'float64', scale: 1, offset: 0, enabled: true, persisted: true,
}) as unknown as Mapping;

function savedSetupState(): WorkbenchV2State {
  return {
    ...INITIAL_STATE,
    mappings: { p1: mapping('p1'), p2: mapping('p2') },
    db: {
      ...INITIAL_STATE.db,
      connector: {
        ...INITIAL_STATE.db.connector, kind: 'sqlite', name: 'Line A', database: '/tmp/line-a.db',
        table: 'sensor_values', write_interval_seconds: 5, connector_id: 'db-1', identity_revision: 'identity-1',
        setup_revision: 'setup-1', workspace_id: 'ws-1', persisted: true, save_state: 'saved',
      },
      targets: {
        p1: { tag_id: 'tag-p1', column_name: 'column_p1', enabled: true },
        p2: { tag_id: 'tag-p2', column_name: 'column_p2', enabled: true },
      },
    },
  };
}

describe('AtomicSetupSave ordered browser saves', () => {
  beforeEach(() => {
    mocks.updateConfig.mockReset();
    mocks.upsertTarget.mockReset();
  });

  it('sends one database save at a time with the revision from the previous reply', async () => {
    const config = deferred<ReturnType<typeof configRecord>>();
    const first = deferred<ReturnType<typeof targetRecord>>();
    const second = deferred<ReturnType<typeof targetRecord>>();
    mocks.updateConfig.mockReturnValueOnce(config.promise);
    mocks.upsertTarget.mockReturnValueOnce(first.promise).mockReturnValueOnce(second.promise);
    const stateRef = { current: savedSetupState() } as React.MutableRefObject<WorkbenchV2State>;
    const { result } = renderHook(() => useStudioV2DatabaseAutosave({ dispatch: vi.fn() }, stateRef, true));

    act(() => {
      result.current.afterDatabaseAction({ type: 'updateDbConnector', patch: {} }, stateRef.current);
      result.current.afterDatabaseAction({ type: 'setAllDbTargetsEnabled', enabled: true }, stateRef.current);
    });
    expect(mocks.updateConfig).toHaveBeenCalledWith(expect.objectContaining({ expected_setup_revision: 'setup-1' }));
    expect(mocks.upsertTarget).not.toHaveBeenCalled();

    await act(async () => { config.resolve(configRecord('setup-2')); });
    await waitFor(() => expect(mocks.upsertTarget).toHaveBeenCalled());
    expect(mocks.upsertTarget).toHaveBeenCalledTimes(1);
    expect(mocks.upsertTarget.mock.calls[0][0]).toEqual({
      pointId: 'p1', request: expect.objectContaining({ expected_setup_revision: 'setup-2' }),
    });

    await act(async () => { first.resolve(targetRecord('p1', 'setup-3')); });
    await waitFor(() => expect(mocks.upsertTarget).toHaveBeenCalledTimes(2));
    expect(mocks.upsertTarget.mock.calls[1][0].request.expected_setup_revision).toBe('setup-3');

    await act(async () => { second.resolve(targetRecord('p2', 'setup-4')); });
    await waitFor(() => expect(stateRef.current.db.connector.setup_revision).toBe('setup-4'));
    expect(stateRef.current.db.targets.p2.save_state).toBe('saved');
  });

  it('keeps a stale revision rejection as a save error without retrying', async () => {
    const conflict = Object.assign(new Error('conflict'), {
      response: { status: 409, data: { success: false, error: { code: 'revision_mismatch', request_id: 'req-stale' } } },
    });
    mocks.upsertTarget.mockRejectedValueOnce(conflict);
    const stateRef = { current: savedSetupState() } as React.MutableRefObject<WorkbenchV2State>;
    const { result } = renderHook(() => useStudioV2DatabaseAutosave({ dispatch: vi.fn() }, stateRef, true));

    act(() => {
      result.current.afterDatabaseAction({ type: 'updateDbTarget', pointId: 'p1', patch: { column_name: 'column_new' } }, stateRef.current);
    });

    await waitFor(() => expect(stateRef.current.db.targets.p1.save_state).toBe('save-error'));
    expect(mocks.upsertTarget).toHaveBeenCalledTimes(1);
    expect(stateRef.current.db.connector.setup_revision).toBe('setup-1');
  });

  it('keeps a newer same-row edit when the earlier reply arrives and resends it', async () => {
    const first = deferred<ReturnType<typeof targetRecord>>();
    const second = deferred<ReturnType<typeof targetRecord>>();
    mocks.upsertTarget.mockReturnValueOnce(first.promise).mockReturnValueOnce(second.promise);
    const { stateRef, perform } = renderAutosave(savedSetupState());

    perform({ type: 'updateDbTarget', pointId: 'p1', patch: { column_name: 'column_a' } });
    expect(mocks.upsertTarget.mock.calls[0][0].request.column_name).toBe('column_a');
    perform({ type: 'updateDbTarget', pointId: 'p1', patch: { column_name: 'column_b' } });

    await act(async () => { first.resolve({ ...targetRecord('p1', 'setup-2'), column_name: 'column_a' }); });
    expect(stateRef.current.db.targets.p1.column_name).toBe('column_b');
    expect(stateRef.current.db.targets.p1.save_state).toBe('saving');
    await waitFor(() => expect(mocks.upsertTarget).toHaveBeenCalledTimes(2));
    expect(mocks.upsertTarget.mock.calls[1][0].request).toEqual(expect.objectContaining({
      column_name: 'column_b', expected_setup_revision: 'setup-2',
    }));

    await act(async () => { second.resolve({ ...targetRecord('p1', 'setup-3'), column_name: 'column_b' }); });
    await waitFor(() => expect(stateRef.current.db.targets.p1.save_state).toBe('saved'));
    expect(stateRef.current.db.targets.p1.column_name).toBe('column_b');
  });

  it('sends a queued connector edit with the latest identity revision and keeps the newer form value', async () => {
    const first = deferred<ReturnType<typeof configRecord>>();
    const second = deferred<ReturnType<typeof configRecord>>();
    mocks.updateConfig.mockReturnValueOnce(first.promise).mockReturnValueOnce(second.promise);
    const { stateRef, perform } = renderAutosave(savedSetupState());

    perform({ type: 'updateDbConnector', patch: { database: '/tmp/line-b.db' } });
    perform({ type: 'updateDbConnector', patch: { database: '/tmp/line-c.db' } });
    await act(async () => {
      first.resolve({ ...configRecord('setup-2'), identity_revision: 'identity-2', database: '/tmp/line-b.db' });
    });

    expect(stateRef.current.db.connector.database).toBe('/tmp/line-c.db');
    await waitFor(() => expect(mocks.updateConfig).toHaveBeenCalledTimes(2));
    expect(mocks.updateConfig.mock.calls[1][0]).toEqual(expect.objectContaining({
      database: '/tmp/line-c.db', expected_connector_revision: 'identity-2', expected_setup_revision: 'setup-2',
    }));

    await act(async () => {
      second.resolve({ ...configRecord('setup-3'), identity_revision: 'identity-3', database: '/tmp/line-c.db' });
    });
    await waitFor(() => expect(stateRef.current.db.connector.save_state).toBe('saved'));
    expect(stateRef.current.db.connector.identity_revision).toBe('identity-3');
  });

  it('keeps a saved connection chosen while an earlier save is in flight', async () => {
    const first = deferred<ReturnType<typeof configRecord>>();
    const second = deferred<ReturnType<typeof configRecord>>();
    mocks.updateConfig.mockReturnValueOnce(first.promise).mockReturnValueOnce(second.promise).mockResolvedValue(configRecord('setup-4'));
    const { stateRef, perform } = renderAutosave(savedSetupState());

    perform({ type: 'updateDbConnector', patch: { database: '/tmp/line-b.db' } });
    perform({ type: 'updateDbConnector', patch: {
      connector_id: 'db-2', identity_revision: 'identity-9', name: 'Pool B', database: '/tmp/pool-b.db',
    } });
    await act(async () => {
      first.resolve({ ...configRecord('setup-2'), identity_revision: 'identity-2', database: '/tmp/line-b.db' });
    });

    expect(stateRef.current.db.connector).toEqual(expect.objectContaining({
      connector_id: 'db-2', identity_revision: 'identity-9', setup_revision: 'setup-2', database: '/tmp/pool-b.db',
    }));
    await waitFor(() => expect(mocks.updateConfig).toHaveBeenCalledTimes(2));
    perform({ type: 'updateDbConnector', patch: { write_interval_seconds: 10 } });
    await act(async () => {
      second.resolve({ ...configRecord('setup-3'), id: 'db-2', identity_revision: 'identity-9', name: 'Pool B', database: '/tmp/pool-b.db' });
    });

    await waitFor(() => expect(mocks.updateConfig).toHaveBeenCalledTimes(3));
    expect(mocks.updateConfig.mock.calls[2][0]).toEqual(expect.objectContaining({
      connector_id: 'db-2', expected_connector_revision: 'identity-9', expected_setup_revision: 'setup-3', write_interval_seconds: 10,
    }));
  });

  it('marks a target edit as pending while a failed connector save holds target saves', async () => {
    mocks.updateConfig.mockRejectedValueOnce(new Error('connector failed'));
    const state = savedSetupState();
    state.db.targets.p1 = { ...state.db.targets.p1, save_state: 'saved' };
    const { stateRef, perform } = renderAutosave(state);

    perform({ type: 'updateDbConnector', patch: { database: '/tmp/line-b.db' } });
    await waitFor(() => expect(stateRef.current.db.connector.save_state).toBe('save-error'));
    perform({ type: 'updateDbTarget', pointId: 'p1', patch: { column_name: 'column_new' } });

    expect(stateRef.current.db.targets.p1.save_state).toBe('saving');
    expect(mocks.upsertTarget).not.toHaveBeenCalled();
  });
});

function renderAutosave(state: WorkbenchV2State) {
  const stateRef = { current: state } as React.MutableRefObject<WorkbenchV2State>;
  const { result } = renderHook(() => useStudioV2DatabaseAutosave({ dispatch: vi.fn() }, stateRef, true));
  // Mirrors the page: the reducer applies the edit, then autosave sees the next state.
  const perform = (action: WorkbenchV2Action) => act(() => {
    stateRef.current = workbenchV2Reducer(stateRef.current, action);
    result.current.afterDatabaseAction(action, stateRef.current);
  });
  return { stateRef, perform };
}
