import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { cleanup, fireEvent, render, screen, waitFor, within } from '@testing-library/react';
import { useState } from 'react';
import { afterEach, beforeAll, beforeEach, describe, expect, it, vi } from 'vitest';
import i18n from '@/i18n/config';
import { RecordingHomeCollection } from '@/features/datalink/workbench-v2/recording-home/RecordingHomeCollection';
import { anomalySummaryStatus, type AnomalySummaryStatus } from '@/features/datalink/workbench-v2/recording-home/AcquisitionAnomalyButton';
import { initialCollectionState, type HomeDeviceSummary, type RecordingCollectionState } from '@/features/datalink/workbench-v2/recording-home/recordingHomeCollections';
import { recordingHomeGuidance } from '@/features/datalink/workbench-v2/recording-home/recordingHomeGuidance';
import { AcquisitionAnomalyError, listDeviceAcquisitionAnomalies } from '@/services/acquisitionAnomalies';
import type { AcquisitionAnomalyItem, AcquisitionAnomalySummary } from '@/types/acquisitionAnomaly';
import { fixture, group } from './recording-home-guidance-fixtures';

vi.mock('@/services/acquisitionAnomalies', async () => {
  const actual = await vi.importActual<typeof import('@/services/acquisitionAnomalies')>('@/services/acquisitionAnomalies');
  return { ...actual, listDeviceAcquisitionAnomalies: vi.fn(), getAcquisitionAnomalySummary: vi.fn() };
});

const devices: HomeDeviceSummary[] = [
  { id: 'MC', name: 'MC 主機', protocol: 'mc_3e', availability: 'available' },
  { id: 'MB', name: 'Modbus 表頭', protocol: 'modbus_tcp', availability: 'available' },
  { id: 'NEW', name: '新設備', protocol: 'modbus_tcp', availability: 'unknown' },
];
const groups = [{ ...group('G1', 'running'), members: [{ ...group('G1').members[0], device_id: 'MC' }] }];
const guidance = recordingHomeGuidance(fixture(devices.map(device => device.id), false, groups));
const summary: AcquisitionAnomalySummary = { items: [{ device_id: 'MB', active_count: 0 }, { device_id: 'MC', active_count: 2 }], total_active: 2 };
const stale: AcquisitionAnomalyItem = {
  id: 11, device_id: 'MC', point_id: 'p5', point_name: 'Temp', point_address: 'D5', kind: 'stale', state: 'ended',
  started_at: '2026-11-05T01:00:00Z', ended_at: '2026-11-05T01:01:02Z', duration_ms: 62_000, last_value: '42', detail_codes: ['no-successful-read'],
};
const disconnected: AcquisitionAnomalyItem = {
  id: 12, device_id: 'MC', point_id: null, point_name: null, point_address: null, kind: 'disconnected', state: 'active',
  started_at: '2026-11-05T02:00:00Z', ended_at: null, duration_ms: null, last_value: null, detail_codes: ['connection-failed'],
};

function Host({ anomalies, mode = 'devices' }: { anomalies?: AnomalySummaryStatus; mode?: RecordingCollectionState['mode'] }) {
  const [state, setState] = useState<RecordingCollectionState>({ ...initialCollectionState(), mode });
  return <RecordingHomeCollection groups={groups} devices={devices} state={state} guidance={guidance} observedIds={[]} unverifiedObserved={false}
    sourcesAvailable restored={{ state, workspaceRevision: 'W', scrollTop: 0, focusId: null }} anomalySummary={anomalies}
    onChange={setState} onOpen={vi.fn()} onOpenDevice={vi.fn()} onRemember={vi.fn()} onRecheck={vi.fn()} />;
}
function renderHost(props: Parameters<typeof Host>[0]) {
  const client = new QueryClient({ defaultOptions: { queries: { retry: false } } });
  return render(<QueryClientProvider client={client}><Host {...props} /></QueryClientProvider>);
}
const ready: AnomalySummaryStatus = { state: 'ready', summary };
const anomalyButton = (scope: HTMLElement, deviceId: string) => within(scope).getByTestId(`home-device-anomalies-${deviceId}`);

beforeAll(async () => { await i18n.changeLanguage('zh-TW'); });
beforeEach(() => vi.mocked(listDeviceAcquisitionAnomalies).mockReset());
afterEach(cleanup);

describe('anomaly summary status projection', () => {
  it('fails closed: an error, even with older data, is unknown; no data yet is checking', () => {
    expect(anomalySummaryStatus({ isError: true, data: summary })).toEqual({ state: 'unavailable' });
    expect(anomalySummaryStatus({ isError: false, data: undefined })).toEqual({ state: 'checking' });
    expect(anomalySummaryStatus({ isError: false, data: summary })).toEqual({ state: 'ready', summary });
  });
});

describe('home acquisition anomaly button', () => {
  it.each(['devices', 'recordings'] as const)('shows the three states next to read points when entering from %s mode', mode => {
    renderHost({ anomalies: ready, mode });
    if (mode === 'recordings') fireEvent.click(screen.getByRole('button', { name: '設備' }));
    const scope = screen.getByTestId('home-management-table');
    const mc = anomalyButton(scope, 'MC'), mb = anomalyButton(scope, 'MB'), unknown = anomalyButton(scope, 'NEW');
    expect(mc).toHaveTextContent('採集異常 2');
    expect(mc).toHaveAttribute('data-anomaly-state', 'active');
    expect(mb).toHaveTextContent(/^採集異常$/);
    expect(mb).toHaveAttribute('data-anomaly-state', 'none');
    expect(unknown).toHaveTextContent('狀態未知');
    expect(unknown).toHaveAttribute('data-anomaly-state', 'unknown');
    expect(mc.previousElementSibling).toHaveTextContent('讀取點位');
  });

  it('a failed summary marks every device unknown instead of normal', () => {
    renderHost({ anomalies: { state: 'unavailable' } });
    for (const id of ['MC', 'MB', 'NEW']) {
      const button = anomalyButton(screen.getByTestId('home-management-table'), id);
      expect(button).toHaveAttribute('data-anomaly-state', 'unknown');
      expect(button).toHaveTextContent('狀態未知');
    }
  });

  it('without a summary source the status is unknown, and while reading it says checking', () => {
    renderHost({});
    expect(anomalyButton(screen.getByTestId('home-management-table'), 'MC')).toHaveAttribute('data-anomaly-state', 'unknown');
    cleanup();
    renderHost({ anomalies: { state: 'checking' } });
    expect(anomalyButton(screen.getByTestId('home-management-table'), 'MC')).toHaveAttribute('data-anomaly-state', 'checking');
  });
});

describe('acquisition anomaly modal', () => {
  it('shows the elapsed duration of an ongoing event', async () => {
    const now = vi.spyOn(Date, 'now').mockReturnValue(Date.parse('2026-11-05T02:10:00Z'));
    try {
      vi.mocked(listDeviceAcquisitionAnomalies).mockResolvedValue({ items: [disconnected], total: 1, limit: 25, offset: 0 });
      renderHost({ anomalies: ready });
      fireEvent.click(anomalyButton(screen.getByTestId('home-management-table'), 'MC'));
      const table = await within(await screen.findByRole('dialog')).findByRole('table');
      const [, activeRow] = within(table).getAllByRole('row');
      expect(within(activeRow).getAllByRole('cell')[4]).toHaveTextContent('10 分');
    } finally {
      now.mockRestore();
    }
  });

  it('opens a classified form modal listing events with start, end, duration, address and last value', async () => {
    vi.mocked(listDeviceAcquisitionAnomalies).mockResolvedValue({ items: [disconnected, stale], total: 2, limit: 25, offset: 0 });
    renderHost({ anomalies: ready });
    fireEvent.click(anomalyButton(screen.getByTestId('home-management-table'), 'MC'));
    const dialog = await screen.findByRole('dialog');
    expect(dialog).toHaveAttribute('data-size', 'form');
    expect(within(dialog).getByRole('heading', { level: 2 })).toHaveTextContent('MC 主機');
    expect(listDeviceAcquisitionAnomalies).toHaveBeenCalledWith('MC', { limit: 25, offset: 0 });
    const table = await within(dialog).findByRole('table');
    for (const header of ['類型', '狀態', '開始', '結束', '持續時間', '點位位址', '最後數值']) expect(within(table).getByRole('columnheader', { name: header })).toBeInTheDocument();
    const [, activeRow, endedRow] = within(table).getAllByRole('row');
    expect(activeRow).toHaveTextContent('斷線');
    expect(activeRow).toHaveTextContent('進行中');
    expect(activeRow).toHaveTextContent('整台設備');
    expect(activeRow).toHaveTextContent('無紀錄');
    expect(endedRow).toHaveTextContent('未更新');
    expect(endedRow).toHaveTextContent('已結束');
    expect(endedRow).toHaveTextContent('D5');
    expect(endedRow).toHaveTextContent('42');
    expect(endedRow).toHaveTextContent('1 分 2 秒');
    expect(within(endedRow).getAllByRole('cell')[3]?.querySelector('time')).toHaveAttribute('datetime', '2026-11-05T01:01:02Z');
  });

  it('filters by kind and state and returns to the first page', async () => {
    vi.mocked(listDeviceAcquisitionAnomalies).mockResolvedValue({ items: [stale], total: 30, limit: 25, offset: 0 });
    renderHost({ anomalies: ready });
    fireEvent.click(anomalyButton(screen.getByTestId('home-management-table'), 'MC'));
    const dialog = await screen.findByRole('dialog');
    await within(dialog).findByRole('table');
    fireEvent.click(within(dialog).getByRole('button', { name: '下一頁' }));
    await waitFor(() => expect(listDeviceAcquisitionAnomalies).toHaveBeenLastCalledWith('MC', { limit: 25, offset: 25 }));
    fireEvent.change(within(dialog).getByLabelText('類型'), { target: { value: 'stale' } });
    await waitFor(() => expect(listDeviceAcquisitionAnomalies).toHaveBeenLastCalledWith('MC', { kind: 'stale', limit: 25, offset: 0 }));
    fireEvent.change(within(dialog).getByLabelText('狀態'), { target: { value: 'active' } });
    await waitFor(() => expect(listDeviceAcquisitionAnomalies).toHaveBeenLastCalledWith('MC', { kind: 'stale', state: 'active', limit: 25, offset: 0 }));
  });

  it('shows an actionable error with retry instead of an empty list, then recovers', async () => {
    vi.mocked(listDeviceAcquisitionAnomalies).mockRejectedValueOnce(new AcquisitionAnomalyError('unavailable', true));
    renderHost({ anomalies: ready });
    fireEvent.click(anomalyButton(screen.getByTestId('home-management-table'), 'MC'));
    const dialog = await screen.findByRole('dialog');
    const alert = await within(dialog).findByRole('alert');
    expect(alert).toHaveTextContent('目前無法讀取採集異常');
    expect(within(dialog).queryByText(/沒有採集異常紀錄|沒有符合/)).toBeNull();
    expect(within(dialog).queryByRole('table')).toBeNull();
    vi.mocked(listDeviceAcquisitionAnomalies).mockResolvedValueOnce({ items: [stale], total: 1, limit: 25, offset: 0 });
    fireEvent.click(within(alert).getByRole('button', { name: '重試' }));
    expect(await within(dialog).findByRole('table')).toHaveTextContent('D5');
    expect(within(dialog).queryByRole('alert')).toBeNull();
  });

  it('closes on Escape and returns focus to the anomaly button', async () => {
    vi.mocked(listDeviceAcquisitionAnomalies).mockResolvedValue({ items: [stale], total: 1, limit: 25, offset: 0 });
    renderHost({ anomalies: ready, mode: 'recordings' });
    fireEvent.click(screen.getByRole('button', { name: '設備' }));
    const trigger = anomalyButton(screen.getByTestId('home-management-table'), 'MC');
    trigger.focus();
    fireEvent.click(trigger);
    const dialog = await screen.findByRole('dialog');
    await within(dialog).findByRole('table');
    fireEvent.keyDown(dialog, { key: 'Escape' });
    await waitFor(() => expect(screen.queryByRole('dialog')).toBeNull());
    await waitFor(() => expect(document.activeElement).toBe(anomalyButton(screen.getByTestId('home-management-table'), 'MC')));
  });

  it('distinguishes no events from no matches only after a successful read', async () => {
    vi.mocked(listDeviceAcquisitionAnomalies).mockResolvedValue({ items: [], total: 0, limit: 25, offset: 0 });
    renderHost({ anomalies: ready });
    fireEvent.click(anomalyButton(screen.getByTestId('home-management-table'), 'MB'));
    const dialog = await screen.findByRole('dialog');
    expect(await within(dialog).findByText('這台設備沒有採集異常紀錄。')).toBeInTheDocument();
    fireEvent.change(within(dialog).getByLabelText('類型'), { target: { value: 'frozen' } });
    expect(await within(dialog).findByText('沒有符合篩選條件的異常。')).toBeInTheDocument();
  });

  it('joins multiple detail reasons with localized separator in zh-TW and en', async () => {
    const multiReasonItem: AcquisitionAnomalyItem = {
      ...disconnected,
      detail_codes: ['connection-failed', 'no-successful-read'],
    };
    vi.mocked(listDeviceAcquisitionAnomalies).mockResolvedValue({ items: [multiReasonItem], total: 1, limit: 25, offset: 0 });
    await i18n.changeLanguage('zh-TW');
    renderHost({ anomalies: ready });
    fireEvent.click(anomalyButton(screen.getByTestId('home-management-table'), 'MC'));
    const dialogZh = await screen.findByRole('dialog');
    expect(await within(dialogZh).findByText('連線失敗；沒有成功讀取')).toBeInTheDocument();
    cleanup();

    await i18n.changeLanguage('en');
    renderHost({ anomalies: ready });
    fireEvent.click(anomalyButton(screen.getByTestId('home-management-table'), 'MC'));
    const dialogEn = await screen.findByRole('dialog');
    expect(await within(dialogEn).findByText('Connection failed; No successful read')).toBeInTheDocument();
    cleanup();
    await i18n.changeLanguage('zh-TW');
  });
});
