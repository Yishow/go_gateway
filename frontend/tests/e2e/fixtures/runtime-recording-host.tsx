import { createRoot } from 'react-dom/client';
import { BrowserRouter, useLocation, useNavigate } from 'react-router-dom';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { RuntimeDashboardRoute } from '../../../src/features/datalink/runtime-dashboard/RuntimeDashboardRoute';
import DatalinkWorkbenchV2Page from '../../../src/pages/datalink/workbench-v2/DatalinkWorkbenchV2Page';
import { parseStudioEntityLocator } from '../../../src/services/studioEntityLocator';
import { ToastProvider } from '../../../src/contexts/ToastContext';
import { api } from '../../../src/services/datalinkClient';
import { studioV2DatalinkApi } from '../../../src/services/studioV2Workspace';
import i18n from '../../../src/i18n/config';
import '../../../src/index.css';
import '../../../src/features/datalink/workbench-v2/styles/workbench-v2.css';

/** Actual production route and native locator receiver; only HTTP and EventSource are fixtures. */
export function RuntimeRecordingHostView() {
  const location = useLocation(); const navigate = useNavigate();
  return <ToastProvider><main className="min-h-screen bg-slate-950 text-slate-100">
    <p data-testid="runtime-recording-fixture-label" className="p-2 text-xs text-amber-200">Real runtime/setup components; mock HTTP and stream data. No device, database or SQL verification.</p>
    {location.pathname === '/studio/v2'
      ? <DatalinkWorkbenchV2Page navigateTo={navigate} locatorIntent={parseStudioEntityLocator(location.search)} />
      : location.pathname === '/studio/runtime' || location.pathname.endsWith('/runtime-recording-host.html')
        ? <RuntimeDashboardRoute /> : <p role="alert">Unsupported fixture destination.</p>}
  </main></ToastProvider>;
}

const roots = typeof document === 'undefined' ? [] : document.querySelectorAll('#runtime-recording-root');
if (roots.length === 1) {
  const root = roots[0]; const params = new URLSearchParams(window.location.search);
  const lang = params.get('lang') === 'zh-TW' ? 'zh-TW' : 'en'; const path = params.get('route') === 'settings' ? '/studio/v2' : '/studio/runtime';
  for (const key of ['lang', 'route', 'scenario']) params.delete(key);
  if (path === '/studio/runtime' && !params.has('device_id')) params.set('device_id', 'A');
  window.history.replaceState(null, '', `${path}?${params.toString()}`);
  studioV2DatalinkApi.defaults.baseURL = '/api/v1/datalink';
  api.defaults.baseURL = '/api/v1/datalink';
  const client = new QueryClient({ defaultOptions: { queries: { retry: false }, mutations: { retry: false } } });
  void i18n.changeLanguage(lang).then(() => createRoot(root).render(<QueryClientProvider client={client}><BrowserRouter><RuntimeRecordingHostView /></BrowserRouter></QueryClientProvider>));
}
