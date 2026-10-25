import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { createRoot } from 'react-dom/client';
import i18n from 'i18next';
import { initReactI18next } from 'react-i18next';
import '../../../src/index.css';
import '../../../src/features/datalink/workbench-v2/styles/workbench-v2.css';
import zh from '../../../src/i18n/locales/zh-TW/workbench-v2.json';
import review from '../../../src/i18n/locales/zh-TW/studio-review.json';
import { BasicRecordingPanel } from '../../../src/features/datalink/workbench-v2/steps/step4/BasicRecordingPanel';
import { studioV2WorkspaceWriteGroupsAPI as groupsAPI } from '../../../src/services/studioV2WorkspaceWriteGroups';
import { studioV2RecordingStartAPI as startAPI } from '../../../src/services/studioV2RecordingStart';
import { studioV2WorkspaceMeasurementsAPI as measurementsAPI } from '../../../src/services/studioV2WorkspaceMeasurements';
import { savedGroup, savedGroupState } from '../../fixtures/writeGroupState';

// Browser layout/interaction evidence only. Service calls are simulated here and all writes fail closed.
await i18n.use(initReactI18next).init({ lng: 'zh-TW', resources: { 'zh-TW': { 'workbench-v2': zh, 'studio-review': review } } });
const base = savedGroupState();
const devices = [base.devices[0], { ...base.devices[0], id: 'dev-b', name: '測試設備 B' }];
devices[0] = { ...devices[0], name: '測試設備 A' };
const points = devices.flatMap((device) => Array.from({ length: 32 }, (_, index) => ({
  ...base.points[0], id: `${device.id}-point-${index}`, device_id: device.id, name: `來源 ${index + 1}`, address: String(40001 + index),
})));
const mappings = Object.fromEntries(points.map((point) => [point.id, {
  ...base.mappings[base.points[0].id], point_id: point.id, persisted_point_id: point.id, device_id: point.device_id,
  tag_id: `tag-${point.id}`, display_name: point.name,
}]));
const state = savedGroupState({ devices, points, mappings });
const applied = new URLSearchParams(window.location.search).get('applied') === '1';
const legacy = savedGroup({
  id: 'entry-legacy-a', workspace_id: 'entry-fixture', name: '設備 A 原有快照', basic_managed_device_id: devices[0].id,
  layout_version: 'managed_wide_v1', status: applied ? 'running' : 'draft', applied_revision: applied ? 'legacy-applied' : '',
  members: points.filter((point) => point.device_id === devices[0].id).map((point, index) => ({
    device_id: point.device_id, point_id: point.id, tag_id: `tag-${point.id}`, source_revision: 'source-fixture',
    mapping_revision: 'mapping-fixture', target_column: `v_${index}`, required: true,
  })),
  row_policy: { interval_seconds: 60, allowed_lateness_seconds: 0, incomplete_policy: 'skip_row' },
});
const initialLegacy = structuredClone(legacy);
const writeCalls: string[] = [];
const blockWrite = async () => { writeCalls.push('blocked'); throw new Error('Read-only component fixture'); };
groupsAPI.list = async () => ({ workspace_id: 'entry-fixture', workspace_revision: 'entry-revision', groups: [structuredClone(legacy)] });
groupsAPI.delivery = async () => { throw new Error('No live SQL evidence in this UI fixture'); };
groupsAPI.create = blockWrite;
groupsAPI.update = blockWrite;
groupsAPI.apply = blockWrite;
groupsAPI.disable = blockWrite;
groupsAPI.remove = blockWrite;
groupsAPI.ensureCircuitProfile = blockWrite;
groupsAPI.schemaPreview = blockWrite;
groupsAPI.schemaApply = blockWrite;
startAPI.ensureBasic = blockWrite;
startAPI.start = blockWrite;
measurementsAPI.list = async () => [];
Object.assign(window, { recordingEntryEvidence: { initialLegacy, writeCalls, currentLegacy: legacy } });

createRoot(document.getElementById('root')!).render(
  <QueryClientProvider client={new QueryClient({ defaultOptions: { queries: { retry: false }, mutations: { retry: false } } })}>
    <main data-workbench-v2="true" className="min-h-screen bg-slate-950 p-4 text-slate-100 sm:p-6">
      <div className="mx-auto max-w-7xl space-y-4">
        <header><h1 className="text-xl font-semibold">儲存並開始記錄</h1><p className="text-xs text-slate-400">隔離元件驗證：模擬服務，沒有設備或資料庫連線。</p></header>
        <BasicRecordingPanel state={state} workspaceId="entry-fixture" readonly={false} onOpenAdvanced={() => undefined} />
      </div>
    </main>
  </QueryClientProvider>,
);
