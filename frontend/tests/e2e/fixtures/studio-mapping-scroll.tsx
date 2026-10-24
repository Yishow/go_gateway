import { createRoot } from 'react-dom/client';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { MemoryRouter } from 'react-router-dom';
import i18n from 'i18next';
import { initReactI18next } from 'react-i18next';
import '../../../src/index.css';
import '../../../src/features/datalink/workbench-v2/styles/workbench-v2.css';
import zh from '../../../src/i18n/locales/zh-TW/workbench-v2.json';
import { WorkbenchV2Shell } from '../../../src/features/datalink/workbench-v2/shell/WorkbenchV2Shell';
import { INITIAL_STATE, useWorkbenchV2State } from '../../../src/features/datalink/workbench-v2/state/useWorkbenchV2State';
import { DEFAULT_RULE } from '../../../src/features/datalink/workbench-v2/state/defaults';
import { DeviceListContext } from '../../../src/features/datalink/workbench-v2/state/deviceColors';
import type { Device } from '../../../src/features/datalink/workbench-v2/state/types';
import type { StudioV2WorkspaceReadinessSummary } from '../../../src/types/studioV2WorkspaceReadiness';

await i18n.use(initReactI18next).init({
  lng: 'zh-TW', resources: { 'zh-TW': { 'workbench-v2': zh } },
});
const devices: Device[] = [{ id: 'plc', name: '測試設備', description: '', protocol: 'modbus_tcp', config: {}, status: 'draft', test: null }];
const queryClient = new QueryClient({ defaultOptions: { queries: { retry: false } } });
const readiness: StudioV2WorkspaceReadinessSummary = {
  ready: false, blocking_count: 1, warning_count: 0,
  issues: [{ code: 'schema-confirmation-required', severity: 'blocking', step: 'Step 4', scope: 'workspace', message: '請到 Step 4 確認資料表。' }],
};
const history = Array.from({ length: 3 }, (_, index) => ({
  id: `event-${index}`, workspace_id: 'fixture', event_type: 'database_config_saved',
  result: 'success', scope: 'workspace', occurred_at: '2026-10-24T00:00:00Z', created_at: '2026-10-24T00:00:00Z',
}));

function Preview() {
  const actions = useWorkbenchV2State({
    ...INITIAL_STATE, current: 3, completed: new Set([1, 2]), devices,
    rules: [{ ...DEFAULT_RULE, id: 'rule', device_id: 'plc', count: 48 }],
    points: [], mappings: {},
  });
  return <div data-workbench-v2="true"
    className="bg-canvas min-h-screen font-sans text-slate-100">
    <DeviceListContext.Provider value={devices}>
      <WorkbenchV2Shell state={actions.state} actions={actions} workspaceReadiness={readiness} workspaceAuditHistory={history} />
    </DeviceListContext.Provider>
  </div>;
}

createRoot(document.getElementById('root')!).render(
  <QueryClientProvider client={queryClient}>
    <MemoryRouter><Preview /></MemoryRouter>
  </QueryClientProvider>,
);
