import { beforeEach, describe, expect, it, vi } from 'vitest';
import { readRecordingHomeSources } from '@/features/datalink/workbench-v2/recording-home/recordingHomeSources';
import { studioV2WorkspaceAPI } from '@/services/studioV2Workspace';
import { studioV2WorkspaceDevicesAPI } from '@/services/studioV2WorkspaceDevices';
import { studioV2RulesAPI } from '@/services/studioV2Rules';
import { studioV2MappingsAPI } from '@/services/studioV2Mappings';
import { pointAPI } from '@/services/datalink';
vi.mock('@/services/studioV2Workspace', () => ({ studioV2WorkspaceAPI: { get: vi.fn() } }));
vi.mock('@/services/studioV2WorkspaceDevices', () => ({ studioV2WorkspaceDevicesAPI: { list: vi.fn() } }));
vi.mock('@/services/studioV2Rules', () => ({ studioV2RulesAPI: { list: vi.fn() } }));
vi.mock('@/services/studioV2Mappings', () => ({ studioV2MappingsAPI: { list: vi.fn() } }));
vi.mock('@/services/datalink', () => ({ pointAPI: { list: vi.fn() } }));
const scope = { workspaceId: 'W', workspaceRevision: '', epoch: 1 };
const workspace = { id: 'W', kind: 'single' as const, status: 'ready' as const, ordered_device_ids: ['D'], created_at: '', updated_at: '' };
const device = { id: 'D', name: 'Device' } as Awaited<ReturnType<typeof studioV2WorkspaceDevicesAPI.list>>[number];
const rule = { id: 'R', workspace_id: 'W', device_id: 'D', revision_id: 'source-revision', enabled: true, skipped_addresses: [], start_address: 'D0', count: 1, data_type: 'float32', naming_prefix: 'Temperature', locked: false, origin: 'manual', created_at: '', updated_at: '' } satisfies Awaited<ReturnType<typeof studioV2RulesAPI.list>>[number];
const mapping = { id: 'M', workspace_id: 'W', point_id: 'P', rule_id: 'R', device_id: 'D', address: 'D0', tag_id: 'T', tag_key: 'temperature', display_name: 'Temperature', target_type: 'float32' as const, enabled: true } as Awaited<ReturnType<typeof studioV2MappingsAPI.list>>[number];
const point = { id: 'P', device_id: 'D', address: 'D0', enabled: true, last_value: 'never show', last_error: 'never show' } as Awaited<ReturnType<typeof pointAPI.list>>[number];
beforeEach(() => {
  vi.resetAllMocks(); vi.mocked(studioV2WorkspaceAPI.get).mockResolvedValue(workspace);
  vi.mocked(studioV2WorkspaceDevicesAPI.list).mockResolvedValue([device]); vi.mocked(studioV2RulesAPI.list).mockResolvedValue([rule]);
  vi.mocked(studioV2MappingsAPI.list).mockResolvedValue([mapping]); vi.mocked(pointAPI.list).mockResolvedValue([point]);
});
describe('home readonly canonical source observation', () => {
  it('uses persisted point identities and the existing candidate resolver without deriving protocol addresses', async () => {
    const result = await readRecordingHomeSources(scope);
    expect(result).toEqual({ state: 'available', scope, value: { candidates: [{ key: 'P|T', device_id: 'D', point_id: 'P', tag_id: 'T', tag_key: 'temperature', label: 'Temperature', device_name: 'Device', address: 'D0', target_type: 'float32' }], excluded: [] } });
    expect(JSON.stringify(result)).not.toContain('never show'); expect(pointAPI.list).toHaveBeenCalledTimes(2);
  });
  it('actual disabled points do not become source candidates', async () => {
    vi.mocked(pointAPI.list).mockResolvedValue([{ ...point, enabled: false }]);
    expect(await readRecordingHomeSources(scope)).toMatchObject({ state: 'available', value: { candidates: [] } });
  });
  it('fresh workspace with no sources is an observed empty result', async () => {
    vi.mocked(studioV2WorkspaceAPI.get).mockResolvedValue({ ...workspace, status: 'empty', ordered_device_ids: [] });
    vi.mocked(studioV2WorkspaceDevicesAPI.list).mockResolvedValue([]); vi.mocked(studioV2RulesAPI.list).mockResolvedValue([]);
    vi.mocked(studioV2MappingsAPI.list).mockResolvedValue([]); vi.mocked(pointAPI.list).mockResolvedValue([]);
    expect(await readRecordingHomeSources(scope)).toEqual({ state: 'available', scope, value: { candidates: [], excluded: [] } });
  });
  it('source revision changes during the read remain unknown even if Share would be unchanged', async () => {
    vi.mocked(studioV2RulesAPI.list).mockResolvedValueOnce([rule]).mockResolvedValue([{ ...rule, revision_id: 'changed' }]);
    expect(await readRecordingHomeSources(scope)).toEqual({ state: 'unknown', scope });
  });
  it('membership changes during reads are unknown, not a coherent fabricated workspace', async () => {
    vi.mocked(studioV2WorkspaceAPI.get).mockResolvedValueOnce(workspace).mockResolvedValue({ ...workspace, ordered_device_ids: [] });
    expect(await readRecordingHomeSources(scope)).toEqual({ state: 'unknown', scope });
  });
  it.each([
    { workspace_id: 'foreign' }, { device_id: 'foreign' }, { point_id: 'missing' }, { address: 'different' }, { rule_id: 'missing' },
  ])('rejects unsupported mapping ownership %j', async patch => {
    vi.mocked(studioV2MappingsAPI.list).mockResolvedValue([{ ...mapping, ...patch }]);
    expect(await readRecordingHomeSources(scope)).toEqual({ state: 'unknown', scope });
  });
  it('source failure returns a safe error with captured scope, never details or empty candidates', async () => {
    vi.mocked(studioV2MappingsAPI.list).mockRejectedValue(new Error('password=secret private.host'));
    expect(await readRecordingHomeSources(scope)).toEqual({ state: 'error', scope });
  });
  it('late completion keeps the originally captured epoch even if its caller object changes', async () => {
    const caller = { ...scope }; let resolve!: (value: typeof workspace) => void;
    vi.mocked(studioV2WorkspaceAPI.get).mockImplementationOnce(() => new Promise(done => { resolve = done; }));
    const pending = readRecordingHomeSources(caller); caller.epoch = 2; caller.workspaceId = 'other';
    resolve(workspace); expect(await pending).toMatchObject({ state: 'available', scope });
  });
  it('other devices values are not projected into this workspace', async () => {
    vi.mocked(pointAPI.list).mockResolvedValue([point, { ...point, id: 'foreign', device_id: 'other' }]);
    const result = await readRecordingHomeSources(scope);
    expect(result.state === 'available' && result.value.candidates).toHaveLength(1);
  });
  it('100 devices and 300 source rules use ten bulk reads, not per-device requests or polling', async () => {
    const devices = Array.from({ length: 100 }, (_, index) => ({ ...device, id: `D${index}`, name: `Device ${index}` }));
    const rules = Array.from({ length: 300 }, (_, index) => ({ ...rule, id: `R${index}`, device_id: `D${index % 100}` }));
    const mappings = rules.map((item, index) => ({ ...mapping, id: `M${index}`, rule_id: item.id, device_id: item.device_id, point_id: `P${index}`, tag_id: `T${index}` }));
    const points = mappings.map(item => ({ ...point, id: item.point_id, device_id: item.device_id }));
    vi.mocked(studioV2WorkspaceAPI.get).mockResolvedValue({ ...workspace, ordered_device_ids: devices.map(item => item.id) });
    vi.mocked(studioV2WorkspaceDevicesAPI.list).mockResolvedValue(devices); vi.mocked(studioV2RulesAPI.list).mockResolvedValue(rules);
    vi.mocked(studioV2MappingsAPI.list).mockResolvedValue(mappings); vi.mocked(pointAPI.list).mockResolvedValue(points);
    const result = await readRecordingHomeSources(scope);
    expect(result.state === 'available' && result.value.candidates).toHaveLength(300);
    for (const read of [studioV2WorkspaceAPI.get, studioV2WorkspaceDevicesAPI.list, studioV2RulesAPI.list, studioV2MappingsAPI.list, pointAPI.list]) {
      expect(read).toHaveBeenCalledTimes(2);
    }
  });
  it('failed bulk reads do not retry or multiply requests across devices', async () => {
    vi.mocked(studioV2RulesAPI.list).mockRejectedValue(new Error('unavailable'));
    expect(await readRecordingHomeSources(scope)).toEqual({ state: 'error', scope });
    for (const read of [studioV2WorkspaceAPI.get, studioV2WorkspaceDevicesAPI.list, studioV2RulesAPI.list, studioV2MappingsAPI.list, pointAPI.list]) {
      expect(read).toHaveBeenCalledTimes(1);
    }
  });

});
