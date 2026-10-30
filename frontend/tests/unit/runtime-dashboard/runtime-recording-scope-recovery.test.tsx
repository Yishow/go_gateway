import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { act, cleanup, fireEvent, render, screen, waitFor } from '@testing-library/react';
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';
import i18n from '../../../src/i18n/config';
import { useStudioV2RuntimeContextQuery } from '../../../src/hooks/datalink/useStudioV2RuntimeContext';
import { RuntimeRecordingRulesPanel } from '../../../src/features/datalink/runtime-dashboard/RuntimeRecordingRulesPanel';
import { studioV2DatalinkApi } from '../../../src/services/studioV2Workspace';
import { ruleEvidenceFixture } from '../services/runtimeRuleEvidenceFixture';

vi.mock('../../../src/services/studioV2Workspace', () => ({ studioV2DatalinkApi: { get: vi.fn() } }));
const clients: QueryClient[] = [];
let revision = 'W1';
let contextUnavailable = false;
let evidenceFailure = 'stale-reference';

function ScopedMonitor() {
  const context = useStudioV2RuntimeContextQuery();
  const scope = context.data?.workspace_revision ? {
    workspace_id: context.data.workspace_id, workspace_revision: context.data.workspace_revision, device_id: 'D',
  } : null;
  return <RuntimeRecordingRulesPanel scope={scope} selectedDeviceId="D" />;
}

function mount() {
  const client = new QueryClient({ defaultOptions: { queries: { retry: false, gcTime: 0, staleTime: 30_000 } } });
  clients.push(client);
  render(<QueryClientProvider client={client}><ScopedMonitor /></QueryClientProvider>);
}

const calls = () => vi.mocked(studioV2DatalinkApi.get).mock.calls;
const requestedRevisions = () => calls().map(([, config]) => {
  const params = config?.params as Record<string, unknown> | undefined;
  return params?.workspace_revision;
});

beforeEach(async () => {
  revision = 'W1'; contextUnavailable = false; evidenceFailure = 'stale-reference';
  await i18n.changeLanguage('en');
  vi.mocked(studioV2DatalinkApi.get).mockReset();
  vi.mocked(studioV2DatalinkApi.get).mockImplementation(async (_url, config) => {
    const params = config?.params as Record<string, unknown> | undefined;
    if (!params) {
      if (contextUnavailable) throw new Error('private backend diagnostic');
      return { data: { success: true, data: { workspace_id: 'W', workspace_revision: revision, devices: [], default_device_id: null, setup: null } } };
    }
    if (params.workspace_revision !== revision) {
      throw { response: { data: { error: { code: evidenceFailure, message: 'private backend diagnostic' } } } };
    }
    const view = ruleEvidenceFixture({ workspace_id: 'W', workspace_revision: revision, device_id: 'D' });
    view.groups[0].effective!.destination.table = revision === 'W1' ? 'meter_old' : 'meter_new';
    return { data: { success: true, data: { workspace_id: 'W', rule_evidence: view } } };
  });
});
afterEach(() => { cleanup(); clients.splice(0).forEach((client) => client.clear()); });

describe('recording evidence workspace revision recovery', () => {
  it('refreshes the current workspace before reading evidence at its new revision', async () => {
    mount(); await screen.findByText('meter_old');
    revision = 'W2';
    fireEvent.click(screen.getByRole('button', { name: 'Refresh recording evidence' }));
    await screen.findByText('meter_new');
    expect(screen.queryByText('meter_old')).not.toBeInTheDocument();
    expect(requestedRevisions()).toEqual([undefined, 'W1', 'W1', undefined, 'W2']);
    expect(calls().every(([url]) => url === '/studio-v2/workspace/runtime-context')).toBe(true);
  });

  it('keeps a failed context refresh unavailable without retry loops or old evidence', async () => {
    mount(); await screen.findByText('meter_old');
    revision = 'W2'; contextUnavailable = true;
    fireEvent.click(screen.getByRole('button', { name: 'Refresh recording evidence' }));
    await waitFor(() => expect(calls()).toHaveLength(4));
    expect(await screen.findByRole('alert')).toHaveTextContent('Recording configuration changed');
    expect(screen.queryByText('meter_old')).not.toBeInTheDocument();
    expect(document.body.textContent).not.toContain('private backend diagnostic');
    expect(requestedRevisions()).toEqual([undefined, 'W1', 'W1', undefined]);
  });

  it('hides stale recording facts while the workspace refresh is still pending', async () => {
    mount(); await screen.findByText('meter_old');
    revision = 'W2';
    const get = vi.mocked(studioV2DatalinkApi.get);
    const original = get.getMockImplementation();
    if (!original) throw new Error('expected configured runtime API fixture');
    let releaseContext!: () => void;
    const pendingContext = new Promise<void>((resolve) => { releaseContext = resolve; });
    get.mockImplementation(async (url, config) => {
      if (!config?.params) await pendingContext;
      return original(url, config);
    });

    fireEvent.click(screen.getByRole('button', { name: 'Refresh recording evidence' }));
    try {
      await waitFor(() => expect(calls()).toHaveLength(4));
      expect(await screen.findByRole('alert')).toHaveTextContent('Recording configuration changed');
      expect(screen.queryByText('meter_old')).not.toBeInTheDocument();
      expect(requestedRevisions()).toEqual([undefined, 'W1', 'W1', undefined]);
    } finally {
      await act(async () => { releaseContext(); await pendingContext; });
    }
    await screen.findByText('meter_new');
    expect(requestedRevisions()).toEqual([undefined, 'W1', 'W1', undefined, 'W2']);
  });

  it.each(['foreign-scope', 'unavailable'])('does not refresh scope for %s failures', async (reason) => {
    mount(); await screen.findByText('meter_old');
    revision = 'W2'; evidenceFailure = reason;
    fireEvent.click(screen.getByRole('button', { name: 'Refresh recording evidence' }));
    await screen.findByRole('alert');
    expect(requestedRevisions()).toEqual([undefined, 'W1', 'W1']);
    expect(screen.queryByText('meter_new')).not.toBeInTheDocument();
  });
});
