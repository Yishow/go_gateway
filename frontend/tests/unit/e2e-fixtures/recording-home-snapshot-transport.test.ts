import { describe, expect, it, vi } from 'vitest';
import type { Page, Route } from '@playwright/test';
import { installSnapshotTransport } from '../../e2e/fixtures/recording-home-snapshot-transport';
import { createRecordingHomeSnapshotAPI, SNAPSHOT_BASE } from '../../e2e/fixtures/recording-home-snapshot-api';
const origin = 'http://127.0.0.1:4173';
type Handler = (route: Route) => Promise<void> | void;
async function install() {
  const api = createRecordingHomeSnapshotAPI();
  const init = vi.fn<(script: () => void) => Promise<void>>().mockResolvedValue(undefined);
  const register = vi.fn<(pattern: string, handler: Handler) => Promise<void>>().mockResolvedValue(undefined);
  const unexpected = await installSnapshotTransport({ addInitScript: init, on: vi.fn(), route: register } as unknown as Page, api, origin);
  return { api, unexpected, init: init.mock.calls[0][0], pattern: register.mock.calls[0][0], handle: register.mock.calls[0][1] };
}
function transport(url: string, method = 'GET', resourceType = 'fetch') {
  const fulfill = vi.fn().mockResolvedValue(undefined), abort = vi.fn().mockResolvedValue(undefined), proceed = vi.fn().mockResolvedValue(undefined);
  const request = { url: () => url, method: () => method, resourceType: () => resourceType, postData: () => null };
  return { route: { request: () => request, fulfill, abort, continue: proceed } as unknown as Route, fulfill, abort, proceed };
}
describe('snapshot routed transport isolation', () => {
  it('opts out of react-grab before startup and intercepts all transport', async () => {
    const fixture = await install();
    try { fixture.init(); expect(Reflect.get(window, '__REACT_GRAB_DISABLED__')).toBe(true);
      expect(localStorage.getItem('i18nextLng')).toBe('en'); expect(fixture.pattern).toBe('**/*');
    } finally { Reflect.deleteProperty(window, '__REACT_GRAB_DISABLED__'); }
  });
  it.each([
    ['https://www.react-grab.com/api/version?source=browser', 'GET', 'fetch'],
    ['https://fonts.googleapis.com/css2?family=Geist', 'GET', 'stylesheet'],
    [`https://outside.invalid${SNAPSHOT_BASE}`, 'GET', 'fetch'],
    [`https://outside.invalid${SNAPSHOT_BASE}`, 'PUT', 'fetch'],
  ])('aborts external %s before it reaches the fixture API', async (url, method, type) => {
    const fixture = await install(), request = transport(url, method, type); await fixture.handle(request.route);
    expect(request.abort).toHaveBeenCalledWith('blockedbyclient'); expect(request.fulfill).not.toHaveBeenCalled();
    expect(request.proceed).not.toHaveBeenCalled(); expect(fixture.api.ledger).toEqual([]);
    expect(fixture.unexpected).toHaveLength(1);
  });
  it('continues app assets but still rejects unknown same-origin API calls', async () => {
    const fixture = await install(), asset = transport(`${origin}/src/main.tsx`, 'GET', 'script'), unknown = transport(`${origin}/api/version`);
    await fixture.handle(asset.route); await fixture.handle(unknown.route);
    expect(asset.proceed).toHaveBeenCalledOnce(); expect(asset.fulfill).not.toHaveBeenCalled();
    expect(unknown.fulfill).toHaveBeenCalledWith(expect.objectContaining({ status: 404 }));
    expect(fixture.api.ledger).toHaveLength(1); expect(fixture.api.ledger[0].unexpected).toBe(true);
  });
});
