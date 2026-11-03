import { beforeEach, expect, it, vi } from 'vitest';
import { getDefaultConfig } from '../../../src/features/datalink/workbench-v2/state/protocols';
import { hydrateStudioV2Device, toStudioV2DeviceUpdateRequest } from '../../../src/features/datalink/workbench-v2/state/studioV2DeviceAutosave';
import { prepareRecordingPointName } from '../../../src/pages/datalink/workbench-v2/prepareRecordingPointName';
import { studioV2DatalinkApi, studioV2WorkspaceAPI } from '../../../src/services/studioV2Workspace';
import { studioV2MappingsAPI } from '../../../src/services/studioV2Mappings';
import { sourceRuleCandidateAPI } from '../../../src/services/sourceRuleCandidates';
import { INITIAL_STATE } from '../../../src/features/datalink/workbench-v2/state/useWorkbenchV2State';
import type { WorkbenchV2State } from '../../../src/features/datalink/workbench-v2/state/types';
import { locatorRead, nativeDevice, nativeMapping, nativeRule } from './studio-locator-fixtures';
vi.mock('../../../src/services/studioV2Workspace', () => ({ studioV2DatalinkApi: { get: vi.fn() }, studioV2WorkspaceAPI: { get: vi.fn() } }));
vi.mock('../../../src/services/studioV2Mappings', () => ({ studioV2MappingsAPI: { list: vi.fn() } }));
vi.mock('../../../src/services/sourceRuleCandidates', () => ({ sourceRuleCandidateAPI: { getTags: vi.fn() } }));
const identity = { device_id: 'D', point_id: 'P-native', tag_id: 'T' };
let state: WorkbenchV2State;
function fixture(protocol: 'modbus_tcp' | 'mc_3e' = 'modbus_tcp') {
  const read = locatorRead(3), address = protocol === 'mc_3e' ? 'D0' : '40001';
  read.entry.devices = [{ ...nativeDevice(), protocol }];
  read.entry.source_rules = [{ ...nativeRule(), start_address: address, count: 1, data_type: 'float32' }];
  read.entry.mappings = [{ ...nativeMapping(), address, target_type: 'float32', scale: 1, offset: 0, display_name: 'SRC0 (40001)' }];
  const wire = { ...read.workspace, locator_entry: read.entry };
  vi.mocked(studioV2DatalinkApi.get).mockResolvedValue({ data: { success: true, data: wire } });
  vi.mocked(studioV2MappingsAPI.list).mockResolvedValue(read.entry.mappings);
  return read;
}
function mcFixture() {
  const read = fixture('mc_3e');
  vi.mocked(studioV2WorkspaceAPI.get).mockResolvedValue({ ...read.workspace, modbus_share: { hydration_state: 'ready', readiness: true, workspace_revision: 'SOURCE-REV' } });
  vi.mocked(sourceRuleCandidateAPI.getTags).mockResolvedValue({ source_rule_id: 'R', workspace_id: 'W', workspace_revision: 'SOURCE-REV', revision_id: 'R2',
    tags: { status: 'ready', candidates: [{ id: 'C', identity: { source_rule_id: 'R', candidate_type: 'tags', candidate_kind: 'tag', derived_from_rule_address: 'D0' }, proposed_signature: 'sig', address: 'D0', point_id: 'P-native', tag_id: 'T', mapping_id: 'M', tag_key: 'temperature', display_name: 'MC', data_type: 'float32', status: 'active' }] } });
  return read;
}
beforeEach(() => { vi.resetAllMocks(); state = { ...structuredClone(INITIAL_STATE), devices: [], rules: [], points: [], mappings: {} }; });
it('uses the real Step3 D+M locator contract for a newly saved float32 range, then verifies its returned rule', async () => {
  fixture();
  const prepared = await prepareRecordingPointName('W', identity, () => state, () => true);
  expect(prepared?.preparation).toMatchObject({ mappingId: 'M', ruleId: 'R', ruleRevision: 'R2', identity });
  expect(studioV2DatalinkApi.get).toHaveBeenCalledExactlyOnceWith('/studio-v2/workspace', expect.objectContaining({ params: { entry: 'locator', workspace_id: 'W', step: '3', device_id: 'D', mapping_id: 'M' } }));
  expect(prepared?.merge?.mappings['P-native']).toMatchObject({ mapping_id: 'M', tag_id: 'T', target_type: 'float32', display_name: 'SRC0 (40001)' });
  expect(studioV2WorkspaceAPI.get).not.toHaveBeenCalled(); expect(sourceRuleCandidateAPI.getTags).not.toHaveBeenCalled();
});
it('still rejects a foreign mapping in the actual scoped locator response', async () => {
  const read = fixture(); read.entry.mappings![0] = { ...read.entry.mappings![0], tag_id: 'FOREIGN' };
  expect(await prepareRecordingPointName('W', identity, () => state, () => true)).toBeUndefined();
});
it('reads the MC source revision from the original workspace authority, never the detached database setup revision', async () => {
  mcFixture();
  expect((await prepareRecordingPointName('W', identity, () => state, () => true))?.preparation.mcEvidence).toMatchObject({ status: 'compatible', scope: { workspace_revision: 'SOURCE-REV' } });
  expect(sourceRuleCandidateAPI.getTags).toHaveBeenCalledWith('R', { workspace_id: 'W', expected_workspace_revision: 'SOURCE-REV', revision_id: 'R2' });
});

it('allows an existing MC mapping when its older candidate snapshot has no tag or mapping ID', async () => {
  mcFixture();
  const snapshot = await sourceRuleCandidateAPI.getTags('R', { workspace_id: 'W', expected_workspace_revision: 'SOURCE-REV', revision_id: 'R2' });
  snapshot.tags.candidates = snapshot.tags.candidates.map(candidate => ({ ...candidate, tag_id: undefined, mapping_id: undefined }));
  vi.mocked(sourceRuleCandidateAPI.getTags).mockResolvedValue(snapshot);
  expect((await prepareRecordingPointName('W', identity, () => state, () => true))?.preparation).toMatchObject({ mappingId: 'M', identity });
});

it.each(['tag_id', 'mapping_id'] as const)('rejects an MC candidate with a conflicting %s', async field => {
  mcFixture();
  const snapshot = await sourceRuleCandidateAPI.getTags('R', { workspace_id: 'W', expected_workspace_revision: 'SOURCE-REV', revision_id: 'R2' });
  snapshot.tags.candidates = snapshot.tags.candidates.map(candidate => ({ ...candidate, [field]: 'FOREIGN' }));
  vi.mocked(sourceRuleCandidateAPI.getTags).mockResolvedValue(snapshot);
  expect(await prepareRecordingPointName('W', identity, () => state, () => true)).toBeUndefined();
});

it.each(['foreign-before', 'foreign-after', 'revision-after', 'candidate-rule-revision'])('rejects MC %s without borrowing the database setup revision', async kind => {
  const read = mcFixture(), valid = { ...read.workspace, modbus_share: { hydration_state: 'ready', readiness: true, workspace_revision: 'SOURCE-REV' } };
  if (kind === 'foreign-before') vi.mocked(studioV2WorkspaceAPI.get).mockResolvedValueOnce({ ...valid, id: 'FOREIGN' });
  if (kind === 'foreign-after') vi.mocked(studioV2WorkspaceAPI.get).mockResolvedValueOnce(valid).mockResolvedValueOnce({ ...valid, id: 'FOREIGN' });
  if (kind === 'revision-after') vi.mocked(studioV2WorkspaceAPI.get).mockResolvedValueOnce(valid).mockResolvedValueOnce({ ...valid, modbus_share: { ...valid.modbus_share, workspace_revision: 'NEW-SOURCE-REV' } });
  if (kind === 'candidate-rule-revision') {
    const snapshot = await sourceRuleCandidateAPI.getTags('R', { workspace_id: 'W', expected_workspace_revision: 'SOURCE-REV', revision_id: 'R2' });
    vi.mocked(sourceRuleCandidateAPI.getTags).mockResolvedValue({ ...snapshot, revision_id: 'FOREIGN-RULE' });
  }
  expect(await prepareRecordingPointName('W', identity, () => state, () => true)).toBeUndefined();
});

it('accepts the newly created MC device original local config after an equivalent server canonical key ordering', async () => {
  const read = mcFixture();
  const device = { ...read.entry.devices![0], connection_config: JSON.stringify({ host: '127.0.0.1', io_no: 1023, network_no: 0, pc_no: 255, port: 15000, station_no: 0, timeout: 5 }) };
  read.entry.devices = [device];
  state.devices = [{ ...hydrateStudioV2Device(device), config: { ...getDefaultConfig('mc_3e'), host: '127.0.0.1', port: 15000 } }];
  expect(toStudioV2DeviceUpdateRequest(state.devices[0])).toEqual(toStudioV2DeviceUpdateRequest(hydrateStudioV2Device(device)));
  expect(JSON.stringify(toStudioV2DeviceUpdateRequest(state.devices[0]))).not.toBe(JSON.stringify(toStudioV2DeviceUpdateRequest(hydrateStudioV2Device(device))));
  expect((await prepareRecordingPointName('W', identity, () => state, () => true))?.preparation.mcEvidence?.status).toBe('compatible');
});

it('still refuses a real MC local config change rather than treating it as server key ordering', async () => {
  const read = mcFixture();
  const device = { ...read.entry.devices![0], connection_config: JSON.stringify({ host: '127.0.0.1', io_no: 1023, network_no: 0, pc_no: 255, port: 15000, station_no: 0, timeout: 5 }) };
  read.entry.devices = [device];
  state.devices = [{ ...hydrateStudioV2Device(device), config: { ...getDefaultConfig('mc_3e'), host: '127.0.0.1', port: 15001 } }];
  expect(await prepareRecordingPointName('W', identity, () => state, () => true)).toBeUndefined();
});
