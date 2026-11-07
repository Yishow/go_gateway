import { describe, expect, it, vi } from 'vitest';
import type { Page, Route } from '@playwright/test';
import { normalizeTypedEnvelope } from '../../../src/utils/safeJson';
import { DatabaseSaveFeedbackApi, databaseSavePaths } from '../../e2e/fixtures/database-save-feedback-api';

// Exercise the installed routing boundary without a browser or external network.
const origin = 'http://127.0.0.1:4173';
type Handler = (route: Route) => Promise<void> | void;
async function install(api: DatabaseSaveFeedbackApi) {
  const addInitScript = vi.fn<(script: () => void) => Promise<void>>().mockResolvedValue(undefined);
  const register = vi.fn<(pattern: string, handler: Handler) => Promise<void>>().mockResolvedValue(undefined);
  await api.install({ addInitScript, on: vi.fn(), route: register } as unknown as Page);
  expect(register).toHaveBeenCalledWith('**/*', expect.any(Function));
  return { init: addInitScript.mock.calls[0][0], handle: register.mock.calls[0][1] };
}
function transport(url: string, method = 'GET', resourceType = 'fetch') {
  const fulfill = vi.fn().mockResolvedValue(undefined), abort = vi.fn().mockResolvedValue(undefined);
  const proceed = vi.fn().mockResolvedValue(undefined);
  const request = { url: () => url, method: () => method, resourceType: () => resourceType,
    failure: () => null, postDataJSON: () => ({ connector_id: 'db-W', name: 'Rejected edit' }) };
  return { route: { request: () => request, fulfill, abort, continue: proceed } as unknown as Route,
    fulfill, abort, proceed };
}

describe('database-save fixture transport isolation', () => {
  it('disables the development overlay before App startup instead of mocking its external calls', async () => {
    const { init } = await install(new DatabaseSaveFeedbackApi());
    try {
      init();
      expect(Reflect.get(window, '__REACT_GRAB_DISABLED__')).toBe(true);
      expect(localStorage.getItem('i18nextLng')).toBe('zh-TW');
    } finally { Reflect.deleteProperty(window, '__REACT_GRAB_DISABLED__'); }
  });

  it.each([
    ['https://www.react-grab.com/api/version?source=browser', 'GET', 'fetch'],
    ['https://fonts.googleapis.com/css2?family=Geist:wght@500&display=swap', 'GET', 'stylesheet'],
    [`https://outside.invalid${databaseSavePaths.config}`, 'GET', 'fetch'],
    [`https://outside.invalid${databaseSavePaths.config}`, 'PUT', 'fetch'],
  ])('aborts and records external %s without any successful fixture response', async (url, method, type) => {
    const api = new DatabaseSaveFeedbackApi(), { handle } = await install(api);
    const request = transport(url, method, type);
    await handle(request.route);
    expect(request.abort).toHaveBeenCalledWith('blockedbyclient');
    expect(request.fulfill).not.toHaveBeenCalled(); expect(request.proceed).not.toHaveBeenCalled();
    expect(api.writes).toEqual([]);
    expect(api.unexpected).toEqual([`non-fixture transport ${method} ${new URL(url).origin}${new URL(url).pathname}`]);
  });

  it('continues only same-origin assets and still rejects unknown local APIs', async () => {
    const api = new DatabaseSaveFeedbackApi(), { handle } = await install(api);
    const asset = transport(`${origin}/src/main.tsx`, 'GET', 'script');
    const unknown = transport(`${origin}/api/version`);
    await handle(asset.route); await handle(unknown.route);
    expect(asset.proceed).toHaveBeenCalledOnce(); expect(asset.fulfill).not.toHaveBeenCalled();
    expect(unknown.abort).toHaveBeenCalledWith('blockedbyclient');
    expect(unknown.fulfill).not.toHaveBeenCalled(); expect(unknown.proceed).not.toHaveBeenCalled();
    expect(api.unexpected).toEqual(['GET /api/version']);
  });

  it('returns the backend revision rejection contract without changing persisted state', async () => {
    const api = new DatabaseSaveFeedbackApi(), { handle } = await install(api), before = structuredClone(api.config);
    const plan = api.prepareSave('conflict'), request = transport(`${origin}${databaseSavePaths.config}`, 'PUT');
    const pending = handle(request.route); await plan.received; plan.release(); await pending;
    // internal/api/handlers/studio_v2_workspace_database_errors.go emits this typed 409.
    expect(request.fulfill).toHaveBeenCalledWith({ status: 409, json: { success: false,
      error: { code: 'revision_mismatch', message: 'database setup changed; reload the saved configuration',
        request_id: 'fixture-config-conflict', retryable: false, action: 'reload' } } });
    expect(normalizeTypedEnvelope({ response: { data: request.fulfill.mock.calls[0][0].json } }))
      .toMatchObject({ code: 'revision_mismatch', action: 'reload', retryable: false });
    expect(api.config).toEqual(before); expect(api.unexpected).toEqual([]);
    expect(request.proceed).not.toHaveBeenCalled(); expect(request.abort).not.toHaveBeenCalled();
  });
});
