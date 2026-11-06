import { useReducer, useRef } from 'react';
import { createRoot } from 'react-dom/client';
import * as Dialog from '@radix-ui/react-dialog';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import type { AxiosAdapter } from 'axios';
import i18n from '../../../src/i18n/config';
import { RecordingDraftNavigation } from '../../../src/features/datalink/workbench-v2/recording-editor/RecordingDraftNavigation';
import { RecordingHomeSavedRangeEditor } from '../../../src/features/datalink/workbench-v2/recording-home/RecordingHomeSavedRangeEditor';
import type { HomeRuleObserver } from '../../../src/features/datalink/workbench-v2/recording-home/RecordingHomeSourceHost';
import { hydrateStudioV2Rule } from '../../../src/features/datalink/workbench-v2/state/studioV2RuleAutosave';
import { hydrateStudioV2Device } from '../../../src/features/datalink/workbench-v2/state/studioV2DeviceAutosave';
import { INITIAL_STATE, workbenchV2Reducer, type WorkbenchV2Action } from '../../../src/features/datalink/workbench-v2/state/useWorkbenchV2State';
import { useStudioV2RuleAutosave } from '../../../src/pages/datalink/workbench-v2/useStudioV2RuleAutosave';
import { studioV2DatalinkApi } from '../../../src/services/studioV2Workspace';
import { api } from '../../../src/services/datalinkClient';
import type { StudioV2UpdateSourceRuleRequest } from '../../../src/services/studioV2Rules';
import type { SourceRuleRecord } from '../../../src/types/datalink';
import { locatorRead, nativeDevice, nativeRule } from '../../unit/workbench-v2/studio-locator-fixtures';
import '../../../src/index.css';
import '../../../src/styles/semantic-modal.css';
import '../../../src/features/datalink/workbench-v2/styles/workbench-v2.css';
import '../../../src/features/datalink/workbench-v2/recording-home/recording-home.css';

// This standalone fixture intercepts every API call; no device or database is contacted.
await i18n.changeLanguage('zh-TW');
const scenario = new URLSearchParams(location.search).get('scenario') ?? 'normalize';
if (!['normalize', 'large', 'overlap', 'reject', 'ambiguous'].includes(scenario)) throw new Error('Unknown fixture scenario');
const protocol = scenario === 'overlap' ? 'mc_3e' : 'fatek_fbs';
const device = { ...nativeDevice(), protocol, name: '本機模擬設備', connection_config: JSON.stringify({
  host: 'fixture.invalid', port: 5000, acquisition: { version: 1,
    profile: protocol === 'mc_3e' ? 'mitsubishi-q-iqr-3e-binary-v1' : 'fatek-fbs-ascii-v1' },
}) };
const items = scenario === 'overlap' ? [
  { item_id: 'wide', name: 'Wide', address: 'D3', register_offset: 0, data_type: 'uint32', semantic_kind: 'gauge' },
  { item_id: 'other', name: 'Other', address: 'D5', register_offset: 0, data_type: 'int16', semantic_kind: 'gauge' },
] : Array.from({ length: scenario === 'large' ? 70 : 1 }, (_, index) => ({
  item_id: `item-${index}`, name: `Item ${index}`, address: `R${index * 2}`, register_offset: 0,
  data_type: 'int16', semantic_kind: 'gauge',
}));
let saved: SourceRuleRecord = { ...nativeRule(), start_address: items[0].address, data_type: items[0].data_type,
  count: items.length, layout_mode: 'mixed', mixed_items: items };
let writes = 0;
const ledger: { method: string; path: string; request?: StudioV2UpdateSourceRuleRequest }[] = [];
const adapter: AxiosAdapter = async config => {
  const method = (config.method ?? 'GET').toUpperCase(), path = config.url ?? '';
  if (method === 'GET' && path === '/studio-v2/workspace') {
    ledger.push({ method, path });
    const read = locatorRead(2);
    read.entry.devices = [device]; read.entry.source_rules = [saved]; read.entry.rule_heads.R = saved.revision_id;
    return { data: { success: true, data: { ...read.workspace, locator_entry: read.entry } }, status: 200, statusText: 'OK', headers: {}, config };
  }
  if (method !== 'PUT' || path !== '/studio-v2/workspace/source-rules/R') throw new Error('Fixture rejected an unexpected API call');
  const request = JSON.parse(String(config.data)) as StudioV2UpdateSourceRuleRequest;
  ledger.push({ method, path, request }); writes += 1;
  if (scenario === 'ambiguous') throw new Error('Fixture simulated an ambiguous transport failure');
  if (scenario === 'reject' && writes === 1) throw Object.assign(new Error('Fixture pre-write rejection'), {
    response: { status: 400, data: { error: { code: 'validation', message: 'fixture private diagnostic' } } },
  });
  const first = request.mixed_items?.[0];
  if (!first || !request.mixed_items) throw new Error('Fixture requires a mixed request');
  saved = { ...saved, ...request, start_address: first.address ?? '', count: request.mixed_items.length,
    data_type: first.data_type, revision_id: `R${writes + 2}` };
  return { data: { success: true, data: { ...saved, runtime_apply_status: 'not_running' } }, status: 200, statusText: 'OK', headers: {}, config };
};
studioV2DatalinkApi.defaults.adapter = adapter;
api.defaults.adapter = adapter;
Object.assign(window, { mixedReviewFixture: { ledger, saved: () => saved } });
const initial = { ...INITIAL_STATE, devices: [hydrateStudioV2Device(device)], rules: [hydrateStudioV2Rule(saved)] };
const client = new QueryClient({ defaultOptions: { queries: { retry: false }, mutations: { retry: false } } });

function Host() {
  const [state, dispatch] = useReducer(workbenchV2Reducer, initial);
  const stateRef = useRef(state), observer = useRef<HomeRuleObserver | undefined>(undefined);
  const leave = useRef<((proceed: () => void) => void) | undefined>(undefined);
  const autosave = useStudioV2RuleAutosave({ dispatch }, stateRef, undefined, {
    submitted: (rule, generation) => observer.current?.submitted?.(rule, generation), acknowledged: () => true,
    saved: (rule, record) => observer.current?.saved(rule, record), failed: (rule, error) => observer.current?.failed(rule, error),
  }, 'W');
  const send = (action: WorkbenchV2Action) => {
    stateRef.current = workbenchV2Reducer(stateRef.current, action); dispatch(action); autosave.afterRuleAction(action);
  };
  const evidence = protocol === 'mc_3e' ? { status: 'compatible' as const,
    scope: { workspace_id: 'W', workspace_revision: 'W2', revision_id: 'R2' },
    points: items.map((item, index) => ({ point_id: `P${index}`, address: item.address })) } : undefined;
  return <RecordingDraftNavigation><Dialog.Root open><Dialog.Portal><Dialog.Overlay className="gw-modal-overlay" />
    <Dialog.Content data-workbench-v2 className="gw-modal gw-home-source-modal" data-size="form" data-source-kind="points">
      <header className="gw-modal-header"><div><Dialog.Title>分散位址修復驗證 · {scenario}</Dialog.Title>
        <Dialog.Description>本機模擬 API 與資料；不連接設備或資料庫。</Dialog.Description></div></header>
      <RecordingHomeSavedRangeEditor target={{ workspaceId: 'W', deviceId: 'D', ruleId: 'R', kind: 'points', epoch: 1 }}
        state={state} rule={state.rules[0]} mcEvidence={evidence} dispatch={send} leave={leave} suspended={false} list={null}
        onLeave={() => undefined} onClose={() => undefined} headerSlot={null} onLocked={() => undefined} ruleObserverRef={observer}
        onRetryShare={() => undefined} onOpenShareSettings={() => undefined} />
    </Dialog.Content></Dialog.Portal></Dialog.Root></RecordingDraftNavigation>;
}
createRoot(document.getElementById('root')!).render(<QueryClientProvider client={client}><Host /></QueryClientProvider>);
