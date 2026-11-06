import { describe, expect, it, vi } from 'vitest';
import { INITIAL_STATE } from '../../../src/features/datalink/workbench-v2/state/useWorkbenchV2State';
import { hydrateStudioV2Device } from '../../../src/features/datalink/workbench-v2/state/studioV2DeviceAutosave';
import { hydrateStudioV2Rule } from '../../../src/features/datalink/workbench-v2/state/studioV2RuleAutosave';
import { deriveAllPoints } from '../../../src/features/datalink/workbench-v2/state/sourceRule';
import { mergePersistedMCState } from '../../../src/features/datalink/workbench-v2/state/mcPersistedPoints';
import { buildGroupCandidates } from '../../../src/features/datalink/workbench-v2/state/writeGroup/candidates';
import type { Device, SourceRuleRecord, StudioV2WorkspaceMappingRecord } from '../../../src/types/datalink';
import type { Route } from '@playwright/test';
import { DatabaseSaveFeedbackApi, databaseSavePaths, readKinds } from '../../e2e/fixtures/database-save-feedback-api';

// Node-only contract checks of the browser fixture. No page, browser, HTTP server or external transport.
function route(path: string, method = 'GET', body: Record<string, unknown> = {}) {
  const fulfill = vi.fn<(response: { status?: number; json: unknown }) => Promise<void>>().mockResolvedValue(undefined);
  const abort = vi.fn<(reason?: string) => Promise<void>>().mockResolvedValue(undefined);
  const request = { url: () => `http://127.0.0.1:4173${path}`, method: () => method, postDataJSON: () => body, failure: () => null };
  return { value: { request: () => request, fulfill, abort } as unknown as Route, fulfill, abort };
}
async function save(api: DatabaseSaveFeedbackApi, name: string) {
  const ack = api.prepareSave(), put = route(databaseSavePaths.config, 'PUT', { name, connector_id: 'db-W' });
  const pending = api.handle(put.value); await ack.received;
  expect(put.fulfill).not.toHaveBeenCalled(); ack.release(); await pending;
  return put;
}

describe('database-save browser API fixture', () => {
  it('bootstraps an actually eligible saved group candidate through production hydration', async () => {
    const api = new DatabaseSaveFeedbackApi();
    async function get<T>(suffix: string) {
      const request = route(`${databaseSavePaths.bootstrap}/${suffix}`); await api.handle(request.value);
      return (request.fulfill.mock.calls[0][0].json as { data: T }).data;
    }
    const devices = (await get<Device[]>('devices')).map(hydrateStudioV2Device);
    const rules = (await get<SourceRuleRecord[]>('source-rules')).map(hydrateStudioV2Rule);
    const points = deriveAllPoints(rules, { D: 'modbus_tcp' });
    const state = { ...INITIAL_STATE, devices, rules, points, mappings: {} };
    const hydrated = mergePersistedMCState(state, {}, await get<StudioV2WorkspaceMappingRecord[]>('mappings'), []);
    expect(buildGroupCandidates({ ...state, ...hydrated }).candidates).toEqual([expect.objectContaining({ device_id: 'D', point_id: 'P-native', tag_id: 'T' })]);
  });
  it('emits complete scoped ACK evidence with separate failed runtime apply', async () => {
    const api = new DatabaseSaveFeedbackApi(), put = await save(api, 'One');
    expect(put.fulfill).toHaveBeenCalledWith({ json: { success: true, data: expect.objectContaining({
      id: 'db-W', workspace_id: 'W', identity_revision: 'identity-1', setup_revision: 'setup-1', name: 'One', runtime_apply_status: 'apply_failed',
    }) } });
    expect(api.writes).toHaveLength(1); expect(api.unexpected).toEqual([]);
  });
  it('holds three independent reads and captures their old snapshots before a newer ACK', async () => {
    const api = new DatabaseSaveFeedbackApi(); await save(api, 'First');
    const held = api.holdRefreshes(), routes = readKinds.map(kind => route(databaseSavePaths[kind]));
    const requests = routes.map(item => api.handle(item.value));
    await Promise.all(readKinds.map(kind => held[kind].started));
    await save(api, 'Second');
    held.targets.release(); await held.targets.finished;
    expect(routes[1].fulfill).toHaveBeenCalledOnce();
    expect(routes[0].fulfill).not.toHaveBeenCalled(); expect(routes[2].fulfill).not.toHaveBeenCalled();
    held.bootstrap.fail(); await held.bootstrap.finished;
    expect(routes[2].fulfill).toHaveBeenCalledWith(expect.objectContaining({ status: 503 }));
    held.config.release(); await Promise.all(requests);
    expect(routes[0].fulfill).toHaveBeenCalledWith({ status: 200, json: { success: true, data: expect.objectContaining({ name: 'First', setup_revision: 'setup-1' }) } });
    expect(api.config).toMatchObject({ name: 'Second', setup_revision: 'setup-2' });
  });
  it('simulates lost ACK as persisted state plus aborted transport, not fake success', async () => {
    const api = new DatabaseSaveFeedbackApi(), ack = api.prepareSave('lost');
    const put = route(databaseSavePaths.config, 'PUT', { name: 'Committed' }), pending = api.handle(put.value);
    await ack.received; ack.release(); await pending;
    expect(put.abort).toHaveBeenCalledWith('connectionclosed'); expect(put.fulfill).not.toHaveBeenCalled();
    const get = route(databaseSavePaths.config); await api.handle(get.value);
    expect(get.fulfill).toHaveBeenCalledWith({ status: 200, json: { success: true, data: expect.objectContaining({ name: 'Committed', setup_revision: 'setup-1' }) } });
    expect(api.writes).toHaveLength(1);
  });
  it('omits revision evidence deliberately instead of repairing malformed ACKs', async () => {
    const api = new DatabaseSaveFeedbackApi(), ack = api.prepareSave('missing-setup');
    const put = route(databaseSavePaths.config, 'PUT', { name: 'Unknown' }), pending = api.handle(put.value);
    await ack.received; ack.release(); await pending;
    const response = put.fulfill.mock.calls[0][0].json as { data: Record<string, unknown> };
    expect(response.data).not.toHaveProperty('setup_revision'); expect(api.config.setup_revision).toBe('setup-1');
  });
  it('aborts and records every unmatched API and unplanned mutation', async () => {
    const api = new DatabaseSaveFeedbackApi(), unknown = route('/api/v1/datalink/unexpected'), mutation = route(databaseSavePaths.config, 'PUT');
    await api.handle(unknown.value); await api.handle(mutation.value);
    expect(unknown.abort).toHaveBeenCalledWith('blockedbyclient'); expect(mutation.abort).toHaveBeenCalledWith('blockedbyclient');
    expect(api.unexpected).toEqual(['GET /api/v1/datalink/unexpected', 'unplanned config PUT']);
  });
});
