import { useEffect, useRef } from 'react';
import { createRoot } from 'react-dom/client';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import i18n from '../../../src/i18n/config';
import { RecordingDraftNavigation } from '../../../src/features/datalink/workbench-v2/recording-editor/RecordingDraftNavigation';
import { RecordingHomeSourceContent } from '../../../src/features/datalink/workbench-v2/recording-home/RecordingHomeSourceContent';
import { RecordingHomeSourceHost } from '../../../src/features/datalink/workbench-v2/recording-home/RecordingHomeSourceHost';
import { hydrateStudioV2Device } from '../../../src/features/datalink/workbench-v2/state/studioV2DeviceAutosave';
import { hydrateStudioV2Rule } from '../../../src/features/datalink/workbench-v2/state/studioV2RuleAutosave';
import { INITIAL_STATE } from '../../../src/features/datalink/workbench-v2/state/useWorkbenchV2State';
import DatalinkWorkbenchV2Page from '../../../src/pages/datalink/workbench-v2/DatalinkWorkbenchV2Page';
import type { DeviceSourcePort } from '../../../src/pages/datalink/workbench-v2/deviceSourceSubmission';
import type { Device, SourceRuleRecord } from '../../../src/types/datalink';
import '../../../src/index.css';
import '../../../src/features/datalink/workbench-v2/styles/workbench-v2.css';
import '../../../src/features/datalink/workbench-v2/recording-home/recording-home.css';

await i18n.changeLanguage('zh-TW');
document.documentElement.classList.add('dark');
const at = '2026-10-29T00:00:00Z';
const device: Device = { id: 'D', name: 'Fixture Modbus', description: '', protocol: 'modbus_tcp', status: 'draft',
  connection_config: '{"host":"fixture.invalid","port":502,"slave_id":1}', last_test_at: null, last_test_success: null,
  last_test_error: '', created_at: at, updated_at: at };
const rule: SourceRuleRecord = { id: 'R', workspace_id: 'W', device_id: 'D', revision_id: 'R2', start_address: '40001',
  count: 16, data_type: 'int16', naming_prefix: 'BLOCK_', enabled: true, skipped_addresses: [], locked: false,
  origin: 'manual', created_at: at, updated_at: at };
const state = { ...INITIAL_STATE, devices: [hydrateStudioV2Device(device)], rules: [hydrateStudioV2Rule(rule)] };
const source: DeviceSourcePort = { submit: async () => { throw new Error('Unexpected write'); },
  inspect: async () => undefined, observations: () => [] };
const client = new QueryClient({ defaultOptions: { queries: { retry: false } } });

function Entry({ open, workspaceChanged, autoOpen }: { open: () => void; workspaceChanged: (id: string) => void; autoOpen: boolean }) {
  const opened = useRef(false);
  useEffect(() => { workspaceChanged('W'); if (autoOpen && !opened.current) { opened.current = true; open(); } }, [autoOpen, open, workspaceChanged]);
  return <main data-testid="recording-home" className="gw-recording-home"><h1>記錄工作台</h1><button type="button" onClick={open}>讀取點位</button></main>;
}
function Root() {
  const params = new URLSearchParams(window.location.search);
  if (params.has('home')) return <QueryClientProvider client={client}><RecordingDraftNavigation>
    <DatalinkWorkbenchV2Page presentation="recordings" locatorIntent={{ kind: 'absent' }} presentationLocation={{ search: '?view=recordings', hash: '#home' }} />
  </RecordingDraftNavigation></QueryClientProvider>;
  return <QueryClientProvider client={client}><RecordingDraftNavigation><RecordingHomeSourceHost
    renderSource={controls => <RecordingHomeSourceContent key={controls.target.epoch} {...controls} state={state}
      dispatch={() => { throw new Error('Unexpected write'); }} source={source} suspended={false}
      onRetryShare={() => undefined} onOpenShareSettings={() => undefined} />}>
    {(open, workspaceChanged) => <Entry open={() => open('W', 'D', params.has('inspect-connection') ? 'connection' : 'points')}
      workspaceChanged={workspaceChanged} autoOpen={params.has('inspect-connection')} />}
  </RecordingHomeSourceHost></RecordingDraftNavigation></QueryClientProvider>;
}
createRoot(document.getElementById('root')!).render(<Root />);
