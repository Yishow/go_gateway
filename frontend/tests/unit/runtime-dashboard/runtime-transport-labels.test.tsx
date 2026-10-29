import { cleanup, fireEvent, render, screen } from '@testing-library/react';
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';
import i18n from '../../../src/i18n/config';
import { LiveStateBanner } from '../../../src/features/datalink/runtime-dashboard/LiveStateBanner';
import { RealtimeLogsPanel } from '../../../src/features/datalink/runtime-dashboard/RealtimeLogsPanel';
const receipt = Date.parse('2026-10-29T10:00:00Z');
const logs = [{ timestamp: '2026-10-29T10:00:05Z', level: 'warn' as const, message: 'Runtime stream connection lost. Reconnecting.' }];
beforeEach(async () => { await i18n.changeLanguage('en'); });
afterEach(cleanup);
describe('connection events and truthful snapshot presentation', () => {
  it.each(['loading', 'degraded', 'error', 'empty', 'empty-workspace', 'missing-device-context'] as const)('never claims LIVE for %s connection events', (routeState) => {
    const { container } = render(<RealtimeLogsPanel logs={logs} routeState={routeState} />);
    expect(container).not.toHaveTextContent(/\bLIVE\b/);
    expect(container.querySelector('.animate-pulse')).not.toBeInTheDocument();
    expect(screen.getByRole('heading', { name: 'Connection events' })).toBeInTheDocument();
    expect(container).toHaveTextContent('Browser transport and point-state notifications, not SQL delivery receipts or process logs.');
    expect(container).not.toHaveTextContent('Diagnostic logs');
  });
  it('claims LIVE only for the shared verified live decision, without SQL or all-points success claims', () => {
    const { container } = render(<><RealtimeLogsPanel logs={logs} routeState="live" /><LiveStateBanner routeState="live" snapshotReceivedAt={receipt} snapshotStale={false} /></>);
    expect(screen.getByTestId('runtime-dashboard-logs-panel')).toHaveTextContent(/\bLIVE\b/);
    expect(container).toHaveTextContent('Stream and snapshot are ready; this does not prove every point is fresh, collector success or SQL delivery.');
    expect(container).not.toHaveTextContent('SQL committed');
  });
  it('retains the 10:00 browser receipt stale after 10:00:05 disconnect and keeps both recovery actions', () => {
    const retry = vi.fn(async () => undefined); const reconnect = vi.fn();
    const view = render(<><LiveStateBanner routeState="live" snapshotReceivedAt={receipt} snapshotStale={false} onRetry={retry} onReconnect={reconnect} /><RealtimeLogsPanel routeState="live" /></>);
    view.rerender(<><LiveStateBanner routeState="degraded" snapshotReceivedAt={receipt} snapshotStale onRetry={retry} onReconnect={reconnect} /><RealtimeLogsPanel routeState="degraded" logs={logs} /></>);
    const time = screen.getByTestId('runtime-snapshot-received-at'); expect(time).toHaveAttribute('datetime', '2026-10-29T10:00:00.000Z');
    expect(time).toHaveTextContent('2026-10-29T10:00:00.000Z');
    expect(screen.getByTestId('runtime-dashboard-live-state-banner')).toHaveTextContent('Snapshot received at (browser)');
    expect(screen.getByTestId('runtime-dashboard-live-state-banner')).toHaveTextContent('Stale snapshot');
    expect(document.body).not.toHaveTextContent(/\bLIVE\b/);
    fireEvent.click(screen.getByRole('button', { name: 'Reconnect live stream' })); expect(reconnect).toHaveBeenCalledTimes(1); expect(retry).not.toHaveBeenCalled();
    fireEvent.click(screen.getByRole('button', { name: 'Retry snapshot' })); expect(retry).toHaveBeenCalledTimes(1);
    expect(time).not.toHaveTextContent('10:00:05');
  });
  it.each([undefined, null, 0, -1, Number.NaN, Number.POSITIVE_INFINITY, 9e18])('shows unavailable for missing/invalid receipt %s without inventing a timestamp', (value) => {
    const view = render(<LiveStateBanner routeState="degraded" snapshotReceivedAt={value} snapshotStale />);
    expect(screen.queryByTestId('runtime-snapshot-received-at')).not.toBeInTheDocument(); expect(view.container).toHaveTextContent('Snapshot receipt time unavailable');
    expect(view.container).not.toHaveTextContent(/Invalid Date|server as-of/);
  });
  it.each(['en', 'zh-TW'])('describes degraded monitoring without promising stream-only recovery in %s', async (locale) => {
    await i18n.changeLanguage(locale);
    const { container } = render(<LiveStateBanner routeState="degraded" snapshotReceivedAt={receipt} snapshotStale />);
    expect(container).toHaveTextContent(locale === 'en'
      ? 'Monitoring degraded. Retained snapshots remain stale until a successful snapshot read and a ready stream.'
      : '監看已降級；保留的快照會維持過期，直到成功讀取快照且串流就緒。');
    expect(container).not.toHaveTextContent('until the stream recovers');
  });
  it('defaults legacy logs-only callers to non-live rather than unconditional LIVE', () => {
    const { container } = render(<RealtimeLogsPanel logs={logs} />); expect(container).not.toHaveTextContent(/\bLIVE\b/);
  });
  it('uses localized Traditional Chinese connection and receipt labels', async () => {
    await i18n.changeLanguage('zh-TW');
    render(<><RealtimeLogsPanel logs={logs} routeState="degraded" /><LiveStateBanner routeState="degraded" snapshotReceivedAt={receipt} snapshotStale /></>);
    expect(screen.getByRole('heading', { name: '連線事件' })).toBeInTheDocument();
    expect(screen.getByTestId('runtime-dashboard-live-state-banner')).toHaveTextContent('快照接收時間（瀏覽器）');
    expect(screen.getByTestId('runtime-dashboard-live-state-banner')).toHaveTextContent('過期快照');
    expect(document.body).not.toHaveTextContent(/\bLIVE\b/);
  });
});
