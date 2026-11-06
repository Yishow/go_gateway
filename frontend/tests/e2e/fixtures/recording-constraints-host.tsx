import type { AxiosAdapter } from 'axios';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { createRoot } from 'react-dom/client';
import i18n from '../../../src/i18n/config';
import { GroupEditor } from '../../../src/features/datalink/workbench-v2/steps/step4/writeGroup/GroupEditor';
import { buildGroupCandidates } from '../../../src/features/datalink/workbench-v2/state/writeGroup/candidates';
import { RecordingDraftNavigation } from '../../../src/features/datalink/workbench-v2/recording-editor/RecordingDraftNavigation';
import { SimpleRecordingShell } from '../../../src/features/datalink/workbench-v2/shell/SimpleRecordingShell';
import { api } from '../../../src/services/datalinkClient';
import { studioV2DatalinkApi } from '../../../src/services/studioV2Workspace';
import type { WriteGroup } from '../../../src/types/studioV2WriteGroup';
import type { WriteGroupCalculationCapabilities } from '../../../src/types/studioV2WriteGroupCalculation';
import { savedGroup, savedGroupState } from '../../fixtures/writeGroupState';
import { quietDelivery } from '../../fixtures/cutoverFixtures';
import { fixtureLifecycle, type LifecycleMockState } from './recording-lifecycle-state';
import '../../../src/index.css';
import '../../../src/features/datalink/workbench-v2/styles/workbench-v2.css';
import '@fontsource/inter/400.css';
import '@fontsource/inter/500.css';
import '@fontsource/inter/600.css';
import '@fontsource/jetbrains-mono/400.css';

// This fixture is presentation evidence only: identities, capabilities and observations are simulated.
const scenario = new URLSearchParams(location.search).get('scenario') ?? 'new';
if (!['new', 'long', 'legacy', 'instant'].includes(scenario)) throw new Error('Unknown fixture scenario');
await i18n.changeLanguage('zh-TW');
document.documentElement.classList.add('dark');
const state = savedGroupState();
state.devices[0].name = '模擬設備 A';
const connector = { ...state.db.connector, kind: 'postgres' as const, name: '模擬目的地', database: 'mock_only', schema: 'public' };
const base = savedGroup();
const destination = { ...base.destination, database: 'mock_only', table_schema: 'public' };
const long: WriteGroup = savedGroup({ id: 'G1', name: '模擬長表記錄', record_id: '11111111-1111-4111-8111-111111111111',
  layout_version: 'managed_long_v1', storage_mode: 'plain', recording_profile: { kind: 'mixed', result_shape: 'member-calculation-v1' },
  destination: { ...destination, storage_strategy: 'managed', table_name: 'fixture_long' },
  members: base.members.map((member, index) => ({ ...member, required: false, calculation: index ? 'mean' : 'latest',
    cadence: { interval: index ? 60 : 15, unit: 'second', alignment: 'utc_epoch' }, source_value_type: 'float64' })),
  row_policy: { allowed_lateness_seconds: 0, incomplete_policy: 'partial' }, write_policy: { mode: 'append', dedupe_capability: 'receipt' } });
const legacy: WriteGroup = savedGroup({ id: 'G1', name: '模擬舊式記錄', recording_profile: undefined,
  destination: { ...destination, table_name: 'fixture_legacy' }, row_policy: { ...base.row_policy, interval_seconds: 60 } });
const instant: WriteGroup = savedGroup({ id: 'G1', name: '模擬瞬時值記錄', layout_version: 'managed_wide_profile_v1',
  recording_profile: { kind: 'instantaneous', interval: 15, unit: 'second', sample_reduction: 'latest' },
  destination: { ...destination, storage_strategy: 'managed', table_name: 'fixture_instant' },
  members: base.members.map(member => ({ ...member, required: false, value_role: 'gauge', source_value_type: 'float64' })),
  row_policy: { interval_seconds: 15, allowed_lateness_seconds: 0, incomplete_policy: 'partial' }, write_policy: { mode: 'append', dedupe_capability: 'receipt' } });
const group = scenario === 'long' ? long : scenario === 'legacy' ? legacy : scenario === 'instant' ? instant : undefined;
const other = savedGroup({ ...instant, id: 'fixture-other', name: '模擬其他記錄', destination: { ...instant.destination, table_name: 'fixture_owned' } });
const capabilities: WriteGroupCalculationCapabilities = {
  version: 1, result_shape: 'member-calculation-v1', editor_supported: true, draft_supported: true, activation_supported: false,
  reason: 'record_readiness_required', methods: ['latest', 'mean', 'delta'],
  cadence: { alignment: 'utc_epoch', units: ['second', 'minute'], minimum: 1, max_seconds: 600, max_minutes: 10, window_bounds_checked: true },
  schema_dialects: ['postgres'], schema_layouts: ['managed_wide_profile_v1', 'managed_long_v1'], unsupported_combinations: [],
};
const ledger: { method: string; path: string; handled: boolean }[] = [];
// Both Axios clients are intercepted in-process. Every write and every unrecognised read is rejected.
const adapter: AxiosAdapter = async config => {
  const method = (config.method ?? 'GET').toUpperCase(), path = config.url ?? '';
  const entry = { method, path, handled: false }; ledger.push(entry);
  if (method !== 'GET') throw new Error('Fixture forbids all writes');
  let data: unknown;
  if (path === '/db-targets/connectors') data = [{ id: 'conn-1', identity_revision: 'crev-1', name: connector.name, kind: 'postgres', enabled: true,
    status: 'unknown', connection_config: { database: 'mock_only', schema: 'public' }, default_write_interval_seconds: 15 }];
  else if (path.endsWith('/database-catalog')) data = { workspace_id: 'ws-1', workspace_revision: 'workspace-1', connector_id: 'conn-1', connector_revision: 'crev-1',
    kind: 'postgres', default_schema: 'public', schemas: [{ name: 'public', tables: ['fixture_available', 'fixture_owned', 'fixture_long', 'fixture_legacy', 'fixture_instant'] }] };
  else if (path.endsWith('/database-metadata') && group) data = { workspace_id: 'ws-1', connector_id: 'conn-1', connector_revision: 'crev-1', database: 'mock_only',
    schema: 'public', table: group.destination.table_name, inspection_status: 'exists', columns: group.members.map(member => ({ name: member.target_column, data_type: 'double precision', nullable: true, primary_key: false })) };
  else if (path.endsWith('/readiness') && group) data = { workspace_id: 'ws-1', workspace_revision: 'workspace-1', group_id: group.id, group_revision: group.revision,
    applied_revision: group.applied_revision, config_ready: true, schema_ready: false, ready: false,
    issues: [{ code: 'schema-confirmation-required', severity: 'blocking', step: 'Step 4', scope: group.id, message: 'Mock: confirm the table structure before starting.' }] };
  else if (path.endsWith('/lifecycle') && group) {
    const scope: LifecycleMockState = { scenario: 'applied-long', source: group, workspaceRevision: 'workspace-1', lineageRevision: 1 };
    const lifecycle = fixtureLifecycle(scope, group);
    // These never-applied fixtures offer neither successor action; no speculative seed is needed.
    delete lifecycle.successor_seed;
    data = lifecycle;
  } else if (path.endsWith('/delivery') && group) data = { ...quietDelivery, group_id: group.id, intake: { state: 'not_running' } };
  else if (path.endsWith('/database-operations/test-writes')) data = { workspace_id: 'ws-1', complete: true, operations: [] };
  else if (path.endsWith('/measurements')) data = [];
  else throw new Error(`Fixture rejects unrecognised read: ${path}`);
  entry.handled = true;
  return { data: { success: true, data }, status: 200, statusText: 'OK', headers: {}, config };
};
api.defaults.baseURL = '/api/v1/datalink';
studioV2DatalinkApi.defaults.baseURL = '/api/v1/datalink';
api.defaults.adapter = adapter;
studioV2DatalinkApi.defaults.adapter = adapter;
Object.assign(window, { recordingConstraintsFixtureLedger: ledger });
const client = new QueryClient({ defaultOptions: { queries: { retry: false }, mutations: { retry: false } } });

createRoot(document.getElementById('root')!).render(<QueryClientProvider client={client}>
  <main data-workbench-v2 className="bg-canvas min-h-screen font-sans text-[var(--gw-text)] antialiased">
    <p data-testid="fixture-notice" className="px-4 py-2 text-xs">驗收測試畫面 · 模擬資料 · 不連線、不儲存 · {scenario}</p>
    <RecordingDraftNavigation><SimpleRecordingShell mode="located" workspaceId="ws-1" onBack={() => undefined} onAdvanced={() => undefined}>
      <GroupEditor presentation="simple" group={group} workspaceId="ws-1" workspaceRevision="workspace-1" connector={connector}
        deviceSource={{ submit: async () => { throw new Error('Fixture forbids source writes'); }, inspect: async () => undefined, observations: () => [] }}
        sourceDevices={[{ id: 'dev-1', name: '模擬設備 A' }]}
        calculationCapabilities={capabilities} configuredDestinations={{ data: { workspace_id: 'ws-1', workspace_revision: 'workspace-1', groups: [other, ...(group ? [group] : [])] }, loading: false, failed: false, retry: () => undefined }}
        candidates={buildGroupCandidates(state).candidates} excluded={[]} readonly={false} onCopy={() => undefined} onSaved={() => undefined}
        onClose={() => undefined} onReload={() => undefined} fetchLatest={async () => { throw new Error('Fixture forbids saves'); }} />
    </SimpleRecordingShell></RecordingDraftNavigation>
  </main>
</QueryClientProvider>);
