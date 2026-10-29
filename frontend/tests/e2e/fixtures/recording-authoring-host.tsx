import { useReducer, useState } from 'react';
import { createRoot } from 'react-dom/client';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import i18n from '../../../src/i18n/config';
import { Step4Database } from '../../../src/features/datalink/workbench-v2/steps/step4/Step4Database';
import { workbenchV2Reducer } from '../../../src/features/datalink/workbench-v2/state/useWorkbenchV2State';
import { studioV2DatalinkApi } from '../../../src/services/studioV2Workspace';
import { recordingAuthoringHostState, authoringScenarioFrom } from './recording-authoring-state';
import '../../../src/index.css';
import '../../../src/features/datalink/workbench-v2/styles/workbench-v2.css';

// The actual production host/services/hooks render here. Node-side page.route owns all mock evidence.
studioV2DatalinkApi.defaults.baseURL = '/api/v1/datalink';
const params = new URLSearchParams(location.search);
const scenario = authoringScenarioFrom(params.get('scenario'));
await i18n.changeLanguage(params.get('lang') === 'zh-TW' ? 'zh-TW' : 'en');
const client = new QueryClient({ defaultOptions: { queries: { retry: false }, mutations: { retry: false } } });
function FixtureHost() {
  const [state, dispatch] = useReducer(workbenchV2Reducer, scenario, recordingAuthoringHostState);
  const [mount, setMount] = useState(0);
  if (params.get('away') === '1') return <main data-testid="fixture-away"><h1>Authoring fixture navigation away</h1>
    <a href={`?scenario=${scenario}`}>Return to authoring fixture</a></main>;
  return <QueryClientProvider client={client}>
    <main data-workbench-v2="true" className="min-h-screen bg-slate-950 p-4 text-slate-100 sm:p-6">
      <div className="mx-auto max-w-7xl space-y-4">
        <header className="space-y-2">
          <h1 className="text-xl font-semibold">Recording authoring — mock API / test data</h1>
          <p>Actual Step4Database; no database, device or SQL verification. Scenario: {scenario}</p>
          <div className="flex gap-3">
            <button type="button" data-testid="fixture-remount-host" onClick={() => setMount(value => value + 1)}>Remount host</button>
            <button type="button" data-testid="fixture-refresh-connector-revision" onClick={() => dispatch({ type: 'updateDbConnector',
              patch: { identity_revision: 'crev-2', save_state: 'saved', persisted: true } })}>Receive externally saved connector revision 2</button>
            <a data-testid="fixture-leave" href={`?scenario=${scenario}&lang=${params.get('lang') ?? 'en'}&away=1`}>Leave fixture</a>
          </div>
        </header>
        <Step4Database key={mount} state={state} dispatch={dispatch} workspaceId="ws-1" />
      </div>
    </main>
  </QueryClientProvider>;
}
createRoot(document.getElementById('root')!).render(<FixtureHost />);
