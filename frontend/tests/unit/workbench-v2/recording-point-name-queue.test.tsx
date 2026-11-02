import { act, renderHook, waitFor } from '@testing-library/react';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { useRef, type PropsWithChildren } from 'react';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import { useRecordingPointNameSubmission } from '../../../src/pages/datalink/workbench-v2/useRecordingPointNameSubmission';
import { useStudioV2MappingAutosave } from '../../../src/pages/datalink/workbench-v2/useStudioV2MappingAutosave';
import { INITIAL_STATE, workbenchV2Reducer, type WorkbenchV2Action } from '../../../src/features/datalink/workbench-v2/state/useWorkbenchV2State';
import type { WorkbenchV2State } from '../../../src/features/datalink/workbench-v2/state/types';
import type { StudioV2MappingSaveResult } from '../../../src/services/studioV2Mappings';
import { locatorRead, nativeDevice, nativeMapping, nativeRule } from './studio-locator-fixtures';
const api = vi.hoisted(() => ({ workspace: vi.fn(), list: vi.fn(), read: vi.fn(), candidates: vi.fn(), update: vi.fn(), create: vi.fn(), remove: vi.fn(), refetch: vi.fn() }));
vi.mock('../../../src/services/studioV2Workspace', () => ({ studioV2WorkspaceAPI: { get: api.workspace } }));
vi.mock('../../../src/services/studioV2Mappings', () => ({ studioV2MappingsAPI: { list: api.list } }));
vi.mock('../../../src/services/studioV2LocatorEntry', () => ({ readStudioLocatorEntry: api.read }));
vi.mock('../../../src/services/sourceRuleCandidates', () => ({ sourceRuleCandidateAPI: { get: api.candidates } }));
vi.mock('../../../src/hooks/datalink/useStudioV2Mappings', () => ({ useStudioV2MappingsQuery: () => ({ isSuccess: false, data: [], refetch: api.refetch }), useCreateStudioV2MappingMutation: () => ({ mutateAsync: api.create }), useUpdateStudioV2MappingMutation: () => ({ mutateAsync: api.update }), useDeleteStudioV2MappingMutation: () => ({ mutateAsync: api.remove }) }));
const identity = { device_id: 'D', point_id: 'P-native', tag_id: 'T' };
let record = { ...nativeMapping(), address: 'D0', target_type: 'int16' as const, scale: 1, unit: '' };
function read() {
  const value = locatorRead(3);
  value.entry.devices = [{ ...nativeDevice(), protocol: 'mc_3e', connection_config: '{"host":"owned.example","port":5000,"network_no":0,"pc_no":255,"io_no":1023,"station_no":0}' }];
  value.entry.source_rules = [{ ...nativeRule(), start_address: 'D0', count: 1 }];
  value.entry.mappings = [record];
  value.workspace.modbus_share = { hydration_state: 'ready', readiness: true, workspace_revision: 'W2' };
  return { ok: true, value };
}
function setup() {
  const client = new QueryClient({ defaultOptions: { queries: { retry: false } } });
  const stateRef: { current: WorkbenchV2State } = { current: { ...structuredClone(INITIAL_STATE), devices: [], rules: [], points: [], mappings: {} } };
  const dispatch = vi.fn();
  const wrapper = ({ children }: PropsWithChildren) => <QueryClientProvider client={client}>{children}</QueryClientProvider>;
  const hook = renderHook(({ token, generallyAllowed }) => {
    const queue = useRef<(before: WorkbenchV2State, action: WorkbenchV2Action, after: WorkbenchV2State) => void>(() => undefined);
    const canMutate = useRef<(id: string) => boolean>(() => false);
    const canQueue = useRef<(id: string) => boolean>(() => false);
    const apply = (action: WorkbenchV2Action) => { stateRef.current = workbenchV2Reducer(stateRef.current, action); dispatch(action); };
    const names = useRecordingPointNameSubmission({ workspaceId: 'W', enabled: true, scopeToken: token, stateRef, apply, canWriteRule: () => true, canQueueRule: id => canQueue.current(id), queue: (...args) => queue.current(...args), canMutateRule: id => canMutate.current(id) });
    canMutate.current = () => generallyAllowed;
    const original = useStudioV2MappingAutosave({ dispatch }, stateRef, true, canMutate.current, names.observation, names.allowPreparedMapping);
    canQueue.current = original.canSubmitPreparedMapping;
    queue.current = original.afterMappingAction;
    return { ...names, original };
  }, { initialProps: { token: {}, generallyAllowed: false }, wrapper });
  return { ...hook, stateRef, dispatch };
}
beforeEach(() => {
  vi.resetAllMocks(); record = { ...nativeMapping(), address: 'D0', target_type: 'int16', scale: 1, unit: '' };
  api.workspace.mockImplementation(async () => read().value.workspace);
  api.list.mockImplementation(async () => [record]); api.read.mockImplementation(async () => read());
  api.candidates.mockResolvedValue({ workspace_id: 'W', workspace_revision: 'W2', source_rule_id: 'R', revision_id: 'R2', tags: { status: 'ready', candidates: [{ id: 'C', identity: { source_rule_id: 'R' }, address: 'D0', point_id: 'P-native', mapping_id: 'M', tag_id: 'T', tag_key: 'temperature', display_name: 'Mapping M', data_type: 'int16', status: 'draft' }] } });
  api.refetch.mockImplementation(async () => ({ data: [record], isError: false }));
  api.update.mockImplementation(async ({ request }) => { record = { ...record, ...request }; return { ...record, runtime_apply_status: 'applied' }; });
});
describe('shared point naming through the original mapping queue', () => {
  it('loads a missing MC source for a saved locator and confirms the same mapping through actual queue ACK and readback', async () => {
    const h = setup(), database = h.stateRef.current.db, settings = h.stateRef.current.settings;
    expect(h.stateRef.current.devices).toEqual([]); expect(h.stateRef.current.rules).toEqual([]);
    const ready = await h.result.current.port.prepare(identity); expect(ready?.status).toBe('ready');
    expect(h.stateRef.current.rules[0].mc_persisted_points?.status).toBe('compatible');
    let outcome; await act(async () => { outcome = await h.result.current.port.submit(ready!.session.id, 'Motor temperature'); });
    expect(outcome).toMatchObject({ status: 'confirmed', currentName: 'Motor temperature' });
    expect(api.update).toHaveBeenCalledExactlyOnceWith({ mappingId: 'M', request: expect.objectContaining({ rule_id: 'R', address: 'D0', tag_key: 'temperature', display_name: 'Motor temperature' }) });
    expect(api.create).not.toHaveBeenCalled(); expect(api.remove).not.toHaveBeenCalled();
    expect(h.stateRef.current.db).toBe(database); expect(h.stateRef.current.settings).toBe(settings);
    expect(h.stateRef.current.mappings['P-native']).toMatchObject({ mapping_id: 'M', persisted_point_id: 'P-native', tag_id: 'T', save_state: 'saved', display_name: 'Motor temperature' });
  });
  it('cancel releases only the local name buffer with zero mutation', async () => {
    const h = setup(), ready = await h.result.current.port.prepare(identity);
    h.result.current.port.release(ready!.session.id);
    expect(await h.result.current.port.submit(ready!.session.id, 'Not sent')).toBeUndefined();
    expect(api.update).not.toHaveBeenCalled(); expect(api.create).not.toHaveBeenCalled(); expect(api.remove).not.toHaveBeenCalled();
  });
  it('one-flights even the read-before-submit interval, and a failed read can be tried again', async () => {
    const h = setup(), ready = await h.result.current.port.prepare(identity);
    let resolve!: (value: typeof record[]) => void;
    api.list.mockImplementationOnce(() => new Promise<typeof record[]>(done => { resolve = done; }));
    const first = h.result.current.port.submit(ready!.session.id, 'One');
    expect(await h.result.current.port.submit(ready!.session.id, 'Two')).toBeUndefined();
    expect(api.update).not.toHaveBeenCalled();
    await act(async () => { resolve([record]); await first; });
    expect(api.update).toHaveBeenCalledTimes(1);
    const next = await h.result.current.port.prepare(identity);
    api.list.mockRejectedValueOnce(new Error('read unavailable'));
    await expect(h.result.current.port.submit(next!.session.id, 'Three')).rejects.toThrow('read unavailable');
    await act(async () => { expect(await h.result.current.port.submit(next!.session.id, 'Three')).toMatchObject({ status: 'confirmed' }); });
  });
  it.each(['foreign ACK', 'newer local name', 'scope ABA'])('does not confirm %s as this name save', async kind => {
    const h = setup(), ready = await h.result.current.port.prepare(identity);
    let resolve!: (value: StudioV2MappingSaveResult) => void;
    api.update.mockImplementationOnce(() => new Promise<StudioV2MappingSaveResult>(done => { resolve = done; }));
    const pending = h.result.current.port.submit(ready!.session.id, 'First');
    await waitFor(() => expect(api.update).toHaveBeenCalledTimes(1));
    if (kind === 'newer local name') h.stateRef.current = workbenchV2Reducer(h.stateRef.current, { type: 'updateMapping', pointId: 'P-native', patch: { display_name: 'Newer' } });
    if (kind === 'scope ABA') { h.rerender({ token: {}, generallyAllowed: false }); h.rerender({ token: {}, generallyAllowed: false }); }
    await act(async () => { resolve({ ...record, display_name: 'First', ...(kind === 'foreign ACK' ? { tag_id: 'FOREIGN' } : {}), runtime_apply_status: 'applied' }); });
    expect(await pending).toMatchObject({ status: 'unknown' });
    if (kind === 'newer local name') expect(h.stateRef.current.mappings['P-native'].display_name).toBe('Newer');
    expect(api.update).toHaveBeenCalledTimes(1);
  });
  it('never lends the name proof to another mapping or a changed non-name field in the same rule', async () => {
    const h = setup(), ready = await h.result.current.port.prepare(identity);
    let resolve!: (value: StudioV2MappingSaveResult) => void;
    api.update.mockImplementationOnce(() => new Promise<StudioV2MappingSaveResult>(done => { resolve = done; }));
    const pending = h.result.current.port.submit(ready!.session.id, 'One');
    await waitFor(() => expect(api.update).toHaveBeenCalledTimes(1));
    const point = h.stateRef.current.points[0], mapping = h.stateRef.current.mappings[point.id];
    expect(h.result.current.allowPreparedMapping(point, mapping)).toBe(true);
    expect(h.result.current.allowPreparedMapping({ ...point, id: 'other' }, { ...mapping, mapping_id: 'other' })).toBe(false);
    for (const patch of [{ target_type: 'float64' as const }, { enabled: false }, { scale: 99 }, { display_name: 'Other' }, { mapping_id: undefined }]) {
      expect(h.result.current.allowPreparedMapping(point, { ...mapping, ...patch })).toBe(false);
    }
    const otherPoint = { ...point, id: 'other' }, otherMapping = { ...mapping, point_id: 'other', mapping_id: 'other', persisted_point_id: 'other', tag_id: 'other' };
    h.stateRef.current = { ...h.stateRef.current, points: [...h.stateRef.current.points, otherPoint], mappings: { ...h.stateRef.current.mappings, other: otherMapping } };
    const before = h.stateRef.current, action: WorkbenchV2Action = { type: 'updateMapping', pointId: 'other', patch: { display_name: 'Other name' } };
    h.stateRef.current = workbenchV2Reducer(before, action); h.result.current.original.afterMappingAction(before, action, h.stateRef.current);
    expect(api.update).toHaveBeenCalledTimes(1);
    await act(async () => { record = { ...record, display_name: 'One' }; resolve({ ...record, runtime_apply_status: 'applied' }); await pending; });
    expect(api.update).toHaveBeenCalledTimes(1); expect(api.create).not.toHaveBeenCalled(); expect(api.remove).not.toHaveBeenCalled();
  });
  it('does not overwrite unrelated local state changed during source reads', async () => {
    const h = setup(); let resolve!: (value: ReturnType<typeof read>) => void;
    api.read.mockImplementationOnce(() => new Promise<ReturnType<typeof read>>(done => { resolve = done; }));
    const preparing = h.result.current.port.prepare(identity);
    await waitFor(() => expect(api.read).toHaveBeenCalledTimes(1));
    const settings = { ...h.stateRef.current.settings, theme: 'dark' as const };
    h.stateRef.current = { ...h.stateRef.current, settings };
    resolve(read()); expect(await preparing).toMatchObject({ status: 'ready' });
    expect(h.stateRef.current.settings).toBe(settings); expect(api.update).not.toHaveBeenCalled();
  });
  it('does not merge a source whose prepare completes after unmount', async () => {
    const h = setup(); let resolve!: (value: typeof record[]) => void;
    api.list.mockImplementationOnce(() => new Promise<typeof record[]>(done => { resolve = done; }));
    const preparing = h.result.current.port.prepare(identity); h.unmount(); resolve([record]);
    expect(await preparing).toBeUndefined(); expect(h.stateRef.current.devices).toEqual([]); expect(h.dispatch).not.toHaveBeenCalled();
  });
  it('keeps an uncertain name write paused when general MC authority later becomes ready', async () => {
    const h = setup(), ready = await h.result.current.port.prepare(identity);
    api.update.mockRejectedValueOnce(new Error('lost response'));
    await act(async () => { expect(await h.result.current.port.submit(ready!.session.id, 'Uncertain')).toMatchObject({ status: 'unknown' }); });
    const token = {};
    h.rerender({ token, generallyAllowed: true });
    await act(async () => { await Promise.resolve(); });
    expect(api.update).toHaveBeenCalledTimes(1); expect(h.stateRef.current.mappings['P-native'].save_state).toBe('save-error');
  });
  it('returns blocked rather than waiting behind an existing same-rule paused queue', async () => {
    const h = setup(), ready = await h.result.current.port.prepare(identity);
    const before = h.stateRef.current, action: WorkbenchV2Action = { type: 'updateMapping', pointId: 'P-native', patch: { unit: 'pending unit' } };
    h.stateRef.current = workbenchV2Reducer(before, action); h.result.current.original.afterMappingAction(before, action, h.stateRef.current);
    expect(h.result.current.original.canSubmitPreparedMapping('R')).toBe(false);
    expect(await h.result.current.port.submit(ready!.session.id, 'Blocked')).toMatchObject({ status: 'blocked' });
    expect(api.update).not.toHaveBeenCalled(); expect(h.stateRef.current.mappings['P-native'].unit).toBe('pending unit');
  });
  it('rejects duplicate or foreign source identity and preserves an existing dirty mapping', async () => {
    const h = setup(); api.list.mockResolvedValueOnce([record, record]);
    expect(await h.result.current.port.prepare(identity)).toBeUndefined();
    expect(await h.result.current.port.prepare({ ...identity, tag_id: 'FOREIGN' })).toBeUndefined();
    const ready = await h.result.current.port.prepare(identity); expect(ready?.status).toBe('ready');
    h.stateRef.current = workbenchV2Reducer(h.stateRef.current, { type: 'updateMapping', pointId: 'P-native', patch: { unit: 'dirty unit' } });
    expect(await h.result.current.port.prepare(identity)).toBeUndefined();
    expect(h.stateRef.current.mappings['P-native'].unit).toBe('dirty unit'); expect(api.update).not.toHaveBeenCalled();
  });
});
