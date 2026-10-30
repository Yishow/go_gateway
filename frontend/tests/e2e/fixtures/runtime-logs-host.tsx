import { createRoot } from 'react-dom/client';
import { BrowserRouter, useLocation, useNavigate } from 'react-router-dom';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import RuntimeLogsPage from '../../../src/features/runtime-logs/RuntimeLogsPage';
import DatalinkWorkbenchV2Page from '../../../src/pages/datalink/workbench-v2/DatalinkWorkbenchV2Page';
import { parseStudioEntityLocator } from '../../../src/services/studioEntityLocator';
import { ToastProvider } from '../../../src/contexts/ToastContext';
import { studioV2DatalinkApi } from '../../../src/services/studioV2Workspace';
import { installRuntimeLogsStream } from './runtime-logs-stream';
import i18n from '../../../src/i18n/config';
import '../../../src/index.css';
import '../../../src/features/datalink/workbench-v2/styles/workbench-v2.css';

/** Real production routes and native receiver. Only HTTP and stream bytes are test fixtures. */
export function RuntimeLogsHostView() {
  const location = useLocation(), navigate = useNavigate();
  return <ToastProvider><div className="min-h-screen bg-slate-50 text-slate-900">
    <p data-testid="runtime-logs-fixture-label" className="border-b border-amber-300 bg-amber-50 p-2 text-xs text-amber-900">
      Real runtime/setup components; mock HTTP and stream data. No backend, device or SQL verification.
    </p>
    {location.pathname === '/studio/v2'
      ? <DatalinkWorkbenchV2Page navigateTo={navigate} locatorIntent={parseStudioEntityLocator(location.search)} />
      : location.pathname === '/studio/logs' ? <RuntimeLogsPage /> : <p role="alert">Unsupported fixture destination.</p>}
  </div></ToastProvider>;
}

const root = typeof document === 'undefined' ? null : document.getElementById('runtime-logs-root');
if (root) {
  window.__runtimeLogsStream = installRuntimeLogsStream(); // Before any component can request its stream.
  const params = new URLSearchParams(window.location.search);
  const lang = params.get('lang') === 'zh-TW' ? 'zh-TW' : 'en';
  for (const key of ['lang', 'scenario']) params.delete(key);
  window.history.replaceState(null, '', `/studio/logs${params.size ? `?${params}` : ''}`);
  studioV2DatalinkApi.defaults.baseURL = '/api/v1/datalink';
  const client = new QueryClient({ defaultOptions: { queries: { retry: false }, mutations: { retry: false } } });
  void i18n.changeLanguage(lang).then(() => createRoot(root).render(<QueryClientProvider client={client}><BrowserRouter><RuntimeLogsHostView /></BrowserRouter></QueryClientProvider>));
}
