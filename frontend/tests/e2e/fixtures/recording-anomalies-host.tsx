import type { AxiosAdapter } from 'axios';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { useState } from 'react';
import { createRoot } from 'react-dom/client';
import i18n from '../../../src/i18n/config';
import { RecordingHomeCollection } from '../../../src/features/datalink/workbench-v2/recording-home/RecordingHomeCollection';
import { anomalySummaryStatus } from '../../../src/features/datalink/workbench-v2/recording-home/AcquisitionAnomalyButton';
import { initialCollectionState, type HomeDeviceSummary } from '../../../src/features/datalink/workbench-v2/recording-home/recordingHomeCollections';
import { recordingHomeGuidance } from '../../../src/features/datalink/workbench-v2/recording-home/recordingHomeGuidance';
import { useAcquisitionAnomalySummaryQuery } from '../../../src/hooks/datalink/useAcquisitionAnomalies';
import { api } from '../../../src/services/datalinkClient';
import { fixture, group } from '../../unit/workbench-v2/recording-home-guidance-fixtures';
import '../../../src/index.css';
import '../../../src/features/datalink/workbench-v2/styles/workbench-v2.css';
import '../../../src/features/datalink/workbench-v2/recording-home/recording-home.css';

// Presentation evidence only. The Axios adapter rejects all non-fixture requests and all writes.
await i18n.changeLanguage('zh-TW');
const scenario = new URLSearchParams(location.search).get('scenario') ?? 'normal';
if (!['normal', 'summary-error', 'list-error'].includes(scenario)) throw new Error('Unknown scenario');
const devices: HomeDeviceSummary[] = [
  { id: 'MC', name: '模擬 MC 主機', protocol: 'mc_3e', availability: 'available' },
  { id: 'MB', name: '模擬 Modbus 表頭', protocol: 'modbus_tcp', availability: 'available' },
  { id: 'NEW', name: '模擬新設備', protocol: 'modbus_tcp', availability: 'unknown' },
];
const groups = [{ ...group('G1'), name: '模擬記錄', members: [{ ...group('G1').members[0], device_id: 'MC' }] }];
const guidance = recordingHomeGuidance(fixture(devices.map(device => device.id), false, groups));
const now = Date.now();
const events = [
  { id: 12, device_id: 'MC', point_id: 'p5', point_name: '溫度', point_address: 'D5', kind: 'stale', state: 'active',
    started_at: new Date(now - 125_000).toISOString(), ended_at: null, duration_ms: null, last_value: '42.7', detail: 'no-successful-read' },
  { id: 11, device_id: 'MC', point_id: 'p10', point_name: '壓力', point_address: 'D10', kind: 'frozen', state: 'ended',
    started_at: '2026-11-05T01:00:00Z', ended_at: '2026-11-05T01:01:02Z', duration_ms: 62_000, last_value: '17.4', detail: 'value-unchanged' },
];
const ledger: { method: string; path: string; params: unknown }[] = [];
let listCalls = 0;
const unavailable = () => Object.assign(new Error('fixture unavailable'), {
  response: { status: 500, data: { error: { code: 'acquisition_anomaly_unavailable' } } },
});
const adapter: AxiosAdapter = async config => {
  const method = (config.method ?? 'GET').toUpperCase(), path = config.url ?? '';
  ledger.push({ method, path, params: config.params });
  if (method !== 'GET') throw new Error('Fixture forbids writes');
  let data: unknown;
  if (path === '/acquisition-anomalies/summary') {
    if (scenario === 'summary-error') throw unavailable();
    data = { items: [{ device_id: 'MC', active_count: 1 }, { device_id: 'MB', active_count: 0 }], total_active: 1 };
  } else if (path === '/devices/MC/acquisition-anomalies') {
    listCalls += 1;
    if (scenario === 'list-error' && listCalls === 1) throw unavailable();
    const params = config.params as { kind?: string; state?: string; limit: number; offset: number };
    const matching = events.filter(event => (!params.kind || event.kind === params.kind) && (!params.state || event.state === params.state));
    data = { items: matching.slice(params.offset, params.offset + params.limit), total: matching.length, limit: params.limit, offset: params.offset };
  } else throw new Error(`Fixture rejects unrecognised read: ${path}`);
  return { data: { success: true, data }, status: 200, statusText: 'OK', headers: {}, config };
};
api.defaults.adapter = adapter;
Object.assign(window, { acquisitionAnomalyFixtureLedger: ledger });
const client = new QueryClient({ defaultOptions: { queries: { retry: false } } });

function Host() {
  const [state, setState] = useState(initialCollectionState);
  const summary = useAcquisitionAnomalySummaryQuery();
  return <main data-workbench-v2 className="gw-recording-home">
    <p data-testid="fixture-notice">驗證畫面 · 模擬 API 與資料 · 不連接設備或資料庫</p>
    <RecordingHomeCollection groups={groups} devices={devices} state={state} guidance={guidance} observedIds={[]} unverifiedObserved={false}
      sourcesAvailable restored={{ state, workspaceRevision: 'fixture', scrollTop: 0, focusId: null }} anomalySummary={anomalySummaryStatus(summary)}
      onChange={setState} onOpen={() => undefined} onOpenDevice={() => undefined} onRemember={() => undefined} onRecheck={() => undefined} />
  </main>;
}

createRoot(document.getElementById('root')!).render(<QueryClientProvider client={client}><Host /></QueryClientProvider>);
