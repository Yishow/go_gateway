import { afterEach, describe, expect, it, vi } from 'vitest';
import { createRecordingHomeSnapshotAPI, SNAPSHOT_BASE, SNAPSHOT_SOURCE } from '../../e2e/fixtures/recording-home-snapshot-api';
import { recordingHomeSourcesAPI } from '../../../src/services/recordingHomeSources';
import { studioV2DatalinkApi } from '../../../src/services/studioV2Workspace';
import { api as datalinkAPI } from '../../../src/services/datalinkClient';
import { readStudioLocatorEntry, readStudioLocatorEditorEntry } from '../../../src/services/studioV2LocatorEntry';
import { readRecordingSourceDiscovery } from '../../../src/features/datalink/workbench-v2/recording-editor/recordingSourceDiscovery';
import { parseWriteGroupLifecycle } from '../../../src/utils/studioV2WriteGroupLifecycleJson';
import { parseWriteGroupListData } from '../../../src/utils/studioV2WriteGroupJson';
import { recordingCollection, initialCollectionState } from '../../../src/features/datalink/workbench-v2/recording-home/recordingHomeCollections';

function connectMock() {
  const fixture = createRecordingHomeSnapshotAPI();
  const get = async (path: string, options?: { params?: unknown }) => {
    const url = new URL(`/api/v1/datalink${path}`, 'http://fixture.invalid');
    if (options?.params && typeof options.params === 'object') {
      Object.entries(options.params).forEach(([key, value]) => { if (typeof value !== 'string') throw new Error('Unexpected query value'); url.searchParams.set(key, value); });
    }
    const reply = fixture.handle({ method: 'GET', url: url.href });
    if (reply.status >= 400) throw { response: { status: reply.status } };
    return { data: reply.body };
  };
  vi.spyOn(studioV2DatalinkApi, 'get').mockImplementation(get);
  if (datalinkAPI !== studioV2DatalinkApi) vi.spyOn(datalinkAPI, 'get').mockImplementation(get);
  return fixture;
}
afterEach(() => vi.restoreAllMocks());
describe('production parsers accept the routed snapshot fixture', () => {
  it('compact summary and full editor discovery project the same canonical sources', async () => {
    const fixture = connectMock();
    const compact = await recordingHomeSourcesAPI.read('A'); expect(compact.state).toBe('available');
    fixture.setPhase('editor');
    const discovery = await readRecordingSourceDiscovery({ workspaceId: 'A', workspaceRevision: 'A1', epoch: 1 });
    expect(discovery.state).toBe('available');
    if (discovery.state !== 'available' || compact.state !== 'available') throw new Error('Fixture did not satisfy production readers');
    expect({ devices: compact.devices, candidates: compact.candidates, excluded: compact.excluded }).toEqual(discovery.value);
    expect(fixture.ledger.filter(call => call.action.startsWith('raw:'))).toHaveLength(8);
    expect(fixture.ledger.filter(call => call.unexpected)).toEqual([]);
  });
  it('complete exact locator envelopes hydrate native device/rule and editor owners', async () => {
    connectMock();
    for (const step of [1, 2] as const) {
      const result = await readStudioLocatorEntry({ workspace_id: 'A', step, device_id: 'D-A', rule_id: step === 2 ? 'R-A' : null,
        mapping_id: null, group_id: null, group_revision: null });
      expect(result).toMatchObject({ ok: true, value: { entry: { workspace_revision: 'A1' } } });
    }
    const locator = { workspace_id: 'A', step: 4 as const, device_id: null, rule_id: null, mapping_id: null,
      group_id: 'A-0001', group_revision: 'A-G1-r1' };
    expect(await readStudioLocatorEntry(locator)).toMatchObject({ ok: true });
    expect(await readStudioLocatorEditorEntry(locator, { workspaceRevision: 'A1', groupRevision: 'A-G1-r1' }))
      .toMatchObject({ ok: true, value: { entry: { editor_hydrated: true } } });
    expect(await readStudioLocatorEditorEntry(locator, { workspaceRevision: 'A1', groupRevision: 'OLD' })).toEqual({ ok: false, reason: 'stale-reference' });
  });
  it('list and every lifecycle response retain exact group and workspace identities', () => {
    const fixture = createRecordingHomeSnapshotAPI();
    const list = fixture.handle({ method: 'GET', url: `${SNAPSHOT_BASE}/write-groups` }).body as { data: unknown };
    const parsed = parseWriteGroupListData(list.data); expect(parsed?.groups).toHaveLength(64);
    for (const group of parsed!.groups) {
      const reply = fixture.handle({ method: 'GET', url: `${SNAPSHOT_BASE}/write-groups/${group.id}/lifecycle` }).body as { data: unknown };
      expect(parseWriteGroupLifecycle(reply.data)).toMatchObject({ workspace_id: 'A', workspace_revision: 'A1', group_id: group.id,
        group_revision: group.revision, actions: expect.arrayContaining([expect.objectContaining({ action: 'edit', allowed: true })]) });
    }
  });
  it.each([60, 1000])('has %i matching records plus four genuinely excluded filter controls', count => {
    const fixture = createRecordingHomeSnapshotAPI({ groupCount: count }), state = fixture.snapshot();
    const data = (fixture.handle({ method: 'GET', url: `${SNAPSHOT_SOURCE}?workspace_id=A` }).body as { data: { devices: [] } }).data;
    const filters = { ...initialCollectionState(), keyword: 'Line A', deviceId: 'D-A', status: 'ready' as const, destination: 'conn-1', page: 2 };
    const collection = recordingCollection(state.groups, data.devices, filters);
    expect(collection.filteredCount).toBe(count); expect(collection.rows).toHaveLength(25);
    expect(collection.rows[0].id).toBe('A-0026');
    for (const key of ['keyword', 'deviceId', 'status', 'destination'] as const) {
      expect(recordingCollection(state.groups, data.devices, { ...filters, [key]: '' }).filteredCount).toBe(count + 1);
    }
  });
});
