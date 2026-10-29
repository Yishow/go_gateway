import { useReducer, useState } from 'react';
import { createRoot } from 'react-dom/client';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import i18n from '../../../src/i18n/config';
import { Step4Database } from '../../../src/features/datalink/workbench-v2/steps/step4/Step4Database';
import { workbenchV2Reducer } from '../../../src/features/datalink/workbench-v2/state/useWorkbenchV2State';
import { studioV2DatalinkApi } from '../../../src/services/studioV2Workspace';
import { recordingLifecycleHostState, scenarioFrom } from './recording-lifecycle-state';
import '../../../src/index.css';
import '../../../src/features/datalink/workbench-v2/styles/workbench-v2.css';

// TEST-ONLY host. Playwright intercepts the exact API paths; no service or controller is replaced.
// Force this isolated entry to its own origin even when a developer has a remote VITE_API_BASE_URL.
studioV2DatalinkApi.defaults.baseURL = '/api/v1/datalink';
await i18n.changeLanguage(new URLSearchParams(location.search).get('lang') === 'zh-TW' ? 'zh-TW' : 'en');
const scenario = scenarioFrom(new URLSearchParams(location.search).get('scenario'));
const client = new QueryClient({ defaultOptions: { queries: { retry: false }, mutations: { retry: false } } });
function FixtureHost() {
  const [state, dispatch] = useReducer(workbenchV2Reducer, undefined, recordingLifecycleHostState);
  const [mount, setMount] = useState(0);
  return <QueryClientProvider client={client}>
    <main data-workbench-v2="true" className="min-h-screen bg-slate-950 p-4 text-slate-100 sm:p-6">
      <div className="mx-auto max-w-7xl space-y-4">
        <header className="space-y-2">
          <h1 className="text-xl font-semibold">Recording lifecycle — mock API / test data</h1>
          <p className="text-sm text-slate-300">Real Step4Database host; no database, device, or SQL delivery evidence. Scenario: {scenario}</p>
          <button type="button" className="rounded border border-slate-500 px-3 py-1" data-testid="fixture-remount-host" onClick={() => setMount(value => value + 1)}>Remount test host</button>
        </header>
        <Step4Database key={mount} state={state} dispatch={dispatch} workspaceId="ws-1" />
      </div>
    </main>
  </QueryClientProvider>;
}
createRoot(document.getElementById('root')!).render(<FixtureHost />);
