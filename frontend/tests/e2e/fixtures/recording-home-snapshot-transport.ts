import type { Page, Request } from '@playwright/test';
import type { RecordingHomeSnapshotAPI, SnapshotCall } from './recording-home-snapshot-api';

/** Browser transport only. Unknown app APIs still fail the fixture ledger; external requests never receive fake success. */
export async function installSnapshotTransport(page: Page, api: RecordingHomeSnapshotAPI, baseURL: string) {
  const origin = new URL(baseURL).origin, unexpected: string[] = [];
  const calls = new WeakMap<Request, SnapshotCall>();
  page.on('requestfinished', request => { const call = calls.get(request); if (call) call.delivery = 'finished'; });
  page.on('requestfailed', request => {
    const call = calls.get(request); if (!call) return;
    call.failure = request.failure()?.errorText ?? 'unknown request failure';
    call.delivery = /abort|cancel/i.test(call.failure) ? 'cancelled' : 'not-delivered';
  });
  await page.route('**/*', async route => {
    const request = route.request(), url = new URL(request.url());
    if (url.origin !== origin || !url.pathname.startsWith('/api/')) {
      const asset = ['script', 'stylesheet', 'font', 'image', 'manifest'].includes(request.resourceType());
      const document = request.resourceType() === 'document' && url.pathname === '/studio/v2';
      if (url.origin === origin && request.method() === 'GET' && (asset || document)) return route.continue();
      unexpected.push(`non-fixture transport ${request.method()} ${url.origin}${url.pathname}`);
      return route.abort('blockedbyclient');
    }
    const reply = api.handle({ method: request.method(), url: request.url(),
      body: request.postData() ? request.postDataJSON() : undefined });
    calls.set(request, reply.call);
    await reply.wait;
    try {
      await route.fulfill({ status: reply.status, json: reply.body });
      if (reply.call.delivery === 'pending') reply.call.delivery = 'fulfilled';
    } catch (error) {
      if (reply.call.delivery !== 'cancelled') reply.call.delivery = 'not-delivered';
      reply.call.failure = String(error);
    }
  });
  await page.addInitScript(() => {
    // Native development-overlay opt-out, not fabricated replies to third-party services.
    Reflect.set(window, '__REACT_GRAB_DISABLED__', true);
    localStorage.setItem('i18nextLng', 'en');
  });
  return unexpected;
}
