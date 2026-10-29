import { createRoot } from 'react-dom/client';
import { BrowserRouter } from 'react-router-dom';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import i18n from '../../../src/i18n/config';
import { studioV2DatalinkApi } from '../../../src/services/studioV2Workspace';
import { guidedScenarioFrom } from './recording-guided-state';
import { RecordingGuidedHostView } from './recording-guided-host-view';
import '../../../src/index.css';
import '../../../src/features/datalink/workbench-v2/styles/workbench-v2.css';

studioV2DatalinkApi.defaults.baseURL = '/api/v1/datalink';
const params = new URLSearchParams(location.search);
await i18n.changeLanguage(params.get('lang') === 'zh-TW' ? 'zh-TW' : 'en');
const client = new QueryClient({ defaultOptions: { queries: { retry: false }, mutations: { retry: false } } });
const root = document.getElementById('root');
if (!root) throw new Error('Missing guide fixture root');
createRoot(root).render(<QueryClientProvider client={client}><BrowserRouter>
  <RecordingGuidedHostView scenario={guidedScenarioFrom(params.get('scenario'))} />
</BrowserRouter></QueryClientProvider>);
