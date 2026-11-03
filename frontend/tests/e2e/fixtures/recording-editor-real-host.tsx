import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { createRoot } from 'react-dom/client';
import i18n from '../../../src/i18n/config';
import { GroupEditor } from '../../../src/features/datalink/workbench-v2/steps/step4/writeGroup/GroupEditor';
import { buildGroupCandidates } from '../../../src/features/datalink/workbench-v2/state/writeGroup/candidates';
import { RecordingDraftNavigation } from '../../../src/features/datalink/workbench-v2/recording-editor/RecordingDraftNavigation';
import { SimpleRecordingShell } from '../../../src/features/datalink/workbench-v2/shell/SimpleRecordingShell';
import { savedGroup, savedGroupState } from '../../fixtures/writeGroupState';
import '../../../src/index.css';
import '../../../src/features/datalink/workbench-v2/styles/workbench-v2.css';
import '@fontsource/inter/400.css';
import '@fontsource/inter/500.css';
import '@fontsource/inter/600.css';
import '@fontsource/jetbrains-mono/400.css';

await i18n.changeLanguage('zh-TW'); document.documentElement.classList.add(new URLSearchParams(location.search).get('theme') === 'light' ? 'light' : 'dark');
const state = savedGroupState(), original = savedGroup();
const group = savedGroup({ destination: { ...original.destination, storage_strategy: 'managed', table_schema: 'public', table_name: 'gw_group_fixture' } });
const client = new QueryClient({ defaultOptions: { queries: { retry: false }, mutations: { retry: false } } });

createRoot(document.getElementById('root')!).render(<QueryClientProvider client={client}>
  <main data-workbench-v2 className="bg-canvas min-h-screen font-sans text-[var(--gw-text)] antialiased"><RecordingDraftNavigation>
    <SimpleRecordingShell mode="located" workspaceId="ws-1" onBack={() => undefined} onAdvanced={() => undefined}>
      <GroupEditor presentation="simple" group={group} workspaceId="ws-1" workspaceRevision="workspace-1" connector={state.db.connector}
        candidates={buildGroupCandidates(state).candidates} excluded={[]} readonly={false} onSaved={() => undefined} onClose={() => undefined}
        onReload={() => undefined} fetchLatest={async () => undefined} />
    </SimpleRecordingShell>
  </RecordingDraftNavigation></main>
</QueryClientProvider>);
