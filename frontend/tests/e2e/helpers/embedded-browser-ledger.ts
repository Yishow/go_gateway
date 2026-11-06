import type { Page } from '@playwright/test';

/** Passive CDP observation only: routing would disable the browser HTTP cache. */
export async function browserLedger(page: Page) {
  const session = await page.context().newCDPSession(page);
  await session.send('Network.enable');
  const events: { kind: string; value: unknown }[] = [];
  for (const name of ['requestWillBeSent', 'responseReceived', 'requestServedFromCache', 'dataReceived', 'loadingFinished', 'loadingFailed'] as const) {
    session.on(`Network.${name}`, value => events.push({ kind: name, value }));
  }
  return { events, dispose: () => session.detach() };
}

export async function performanceLedger(page: Page) {
  return page.evaluate(() => ({
    timeOrigin: performance.timeOrigin,
    navigation: performance.getEntriesByType('navigation').map(entry => entry.toJSON()),
    resources: performance.getEntriesByType('resource').map(entry => entry.toJSON()),
    totals: ([...performance.getEntriesByType('navigation'), ...performance.getEntriesByType('resource')] as PerformanceResourceTiming[]).reduce((sum, entry) => ({
      requestCount: sum.requestCount + 1, encodedBodyBytes: sum.encodedBodyBytes + entry.encodedBodySize,
      decodedBodyBytes: sum.decodedBodyBytes + entry.decodedBodySize, transferBytes: sum.transferBytes + entry.transferSize,
    }), { requestCount: 0, encodedBodyBytes: 0, decodedBodyBytes: 0, transferBytes: 0 }),
  }));
}

export function summarize(values: number[]) {
  const sorted = [...values].sort((a, b) => a - b);
  return { raw: values, median: sorted[Math.floor(sorted.length / 2)], p95: sorted[Math.ceil(sorted.length * 0.95) - 1] };
}
