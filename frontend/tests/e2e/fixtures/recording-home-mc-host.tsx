import { useRef } from 'react';
import { createRoot } from 'react-dom/client';
import * as Dialog from '@radix-ui/react-dialog';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import i18n from '../../../src/i18n/config';
import { RecordingDraftNavigation } from '../../../src/features/datalink/workbench-v2/recording-editor/RecordingDraftNavigation';
import { RecordingHomeSavedRangeEditor } from '../../../src/features/datalink/workbench-v2/recording-home/RecordingHomeSavedRangeEditor';
import type { HomeRuleObserver } from '../../../src/features/datalink/workbench-v2/recording-home/RecordingHomeSourceHost';
import { hydrateStudioV2Rule } from '../../../src/features/datalink/workbench-v2/state/studioV2RuleAutosave';
import { hydrateStudioV2Device } from '../../../src/features/datalink/workbench-v2/state/studioV2DeviceAutosave';
import { INITIAL_STATE } from '../../../src/features/datalink/workbench-v2/state/useWorkbenchV2State';
import { useStudioV2PersistedMCPoints } from '../../../src/pages/datalink/workbench-v2/useStudioV2PersistedMCPoints';
import type { Device, SourceRuleRecord } from '../../../src/types/datalink';
import '../../../src/index.css';
import '../../../src/styles/semantic-modal.css';
import '../../../src/features/datalink/workbench-v2/styles/workbench-v2.css';
import '../../../src/features/datalink/workbench-v2/recording-home/recording-home.css';

await i18n.changeLanguage('zh-TW');
document.documentElement.classList.add('dark');
const at = '2026-10-29T00:00:00Z';
const device: Device = { id: 'D', name: 'MC', description: '', protocol: 'mc_3e', status: 'draft',
  connection_config: '{"host":"fixture.invalid","port":5000}', last_test_at: null, last_test_success: null,
  last_test_error: '', created_at: at, updated_at: at };
const rule: SourceRuleRecord = { id: 'R', workspace_id: 'W', device_id: 'D', revision_id: 'R2', start_address: 'D0',
  count: 32, data_type: 'int16', naming_prefix: 'TAG_', enabled: true, skipped_addresses: [], locked: false,
  origin: 'manual', created_at: at, updated_at: at };
const rules = [rule], devices = [device], scope = { workspace_id: 'W', workspace_revision: 'W2' };
const local = hydrateStudioV2Rule(rule);
const state = { ...INITIAL_STATE, devices: [hydrateStudioV2Device(device)], rules: [local] };
const client = new QueryClient({ defaultOptions: { queries: { retry: false } } });

function Host() {
  const evidence = useStudioV2PersistedMCPoints(true, rules, devices, scope);
  const leave = useRef<((proceed: () => void) => void) | undefined>(undefined);
  const observer = useRef<HomeRuleObserver | undefined>(undefined);
  return <RecordingDraftNavigation>
    <Dialog.Root open><Dialog.Portal><Dialog.Overlay className="gw-modal-overlay" />
      <Dialog.Content data-workbench-v2 className="gw-modal gw-home-source-modal" data-size="form" data-source-kind="points">
        <header className="gw-modal-header"><div><Dialog.Title>MC · 設定讀取點位</Dialog.Title>
          <Dialog.Description>這些是共用來源設定；只有按下儲存才會送出變更。</Dialog.Description></div></header>
        <RecordingHomeSavedRangeEditor target={{ workspaceId: 'W', deviceId: 'D', ruleId: 'R', kind: 'points', epoch: 1 }}
          state={state} rule={local} mcEvidence={evidence.evidenceByRule.R} dispatch={() => { throw new Error('Unexpected write'); }}
          leave={leave} suspended={false} list={null} onLeave={() => undefined} onClose={() => undefined} headerSlot={null} onLocked={() => undefined}
          ruleObserverRef={observer} onRetryShare={() => undefined} onOpenShareSettings={() => undefined} />
      </Dialog.Content>
    </Dialog.Portal></Dialog.Root>
  </RecordingDraftNavigation>;
}
function Root() {
  return <QueryClientProvider client={client}><Host /></QueryClientProvider>;
}
createRoot(document.getElementById('root')!).render(<Root />);
