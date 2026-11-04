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

await i18n.changeLanguage('zh-TW'); document.documentElement.classList.add('dark');
const state = savedGroupState(), original = savedGroup();
const connector = { ...state.db.connector, kind: 'postgres' as const, schema: 'public', database: 'test_database' };
const group = savedGroup({ destination: { ...original.destination, storage_strategy: 'managed', database: 'test_database', table_schema: 'public', table_name: 'gw_group_fixture' } });
const other = savedGroup({ id: 'other-group', destination: { ...group.destination, table_name: 'gw_other_owner' } });
const client = new QueryClient({ defaultOptions: { queries: { retry: false }, mutations: { retry: false } } });

createRoot(document.getElementById('root')!).render(<QueryClientProvider client={client}>
  <main data-workbench-v2 className="bg-canvas min-h-screen font-sans text-[var(--gw-text)] antialiased"><RecordingDraftNavigation>
    <SimpleRecordingShell mode="located" workspaceId="ws-1" onBack={() => undefined} onAdvanced={() => undefined}>
      <GroupEditor presentation="simple" group={group} workspaceId="ws-1" workspaceRevision="workspace-1" connector={connector}
        configuredDestinations={{ data: { workspace_id: 'ws-1', workspace_revision: 'workspace-1', groups: [group, other] }, loading: false, failed: false, retry: () => undefined }}
        candidates={buildGroupCandidates(state).candidates} excluded={[]} readonly={false} onSaved={() => undefined} onClose={() => undefined}
        onReload={() => undefined} fetchLatest={async () => undefined} />
    </SimpleRecordingShell>
  </RecordingDraftNavigation></main>
</QueryClientProvider>);
