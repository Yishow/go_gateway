import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { act, renderHook, waitFor } from '@testing-library/react';
import { useReducer } from 'react';
import type { ReactNode } from 'react';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import { useSettingsOperations } from '../../../src/features/datalink/workbench-v2/settings/useSettingsOperations';
import { INITIAL_STATE, workbenchV2Reducer } from '../../../src/features/datalink/workbench-v2/state/useWorkbenchV2State';
import { dbTargetAPI, settingsAPI } from '../../../src/services/datalink';

vi.mock('react-i18next', () => ({
  useTranslation: () => ({ t: (key: string, options?: { defaultValue?: string }) => options?.defaultValue ?? key }),
}));
vi.mock('../../../src/services/datalink', () => ({
  settingsAPI: { updateKey: vi.fn(), listItems: vi.fn() },
  dbTargetAPI: { createConnector: vi.fn(), updateConnector: vi.fn(), deleteConnector: vi.fn(), testConnector: vi.fn(), listConnectors: vi.fn() },
}));

function wrapper({ children }: { children: ReactNode }) {
  return <QueryClientProvider client={new QueryClient({ defaultOptions: { mutations: { retry: false } } })}>{children}</QueryClientProvider>;
}

describe('Settings operations race ownership', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('does not let an older connector error replace the newer operation owner', async () => {
    vi.mocked(settingsAPI.listItems).mockResolvedValue([]);
    let rejectFirst!: (error: unknown) => void;
    let rejectSecond!: (error: unknown) => void;
    vi.mocked(dbTargetAPI.createConnector)
      .mockReturnValueOnce(new Promise((_, reject) => { rejectFirst = reject; }))
      .mockReturnValueOnce(new Promise((_, reject) => { rejectSecond = reject; }));
    const { result } = renderHook(() => useSettingsOperations(INITIAL_STATE, vi.fn()), { wrapper });

    act(() => { void result.current.onAddConnector(); });
    act(() => { void result.current.onAddConnector(); });
    act(() => { rejectFirst({ error: { code: 'settings_update_failed', retryable: true, request_id: 'old-request' } }); });
    await waitFor(() => expect(dbTargetAPI.createConnector).toHaveBeenCalledTimes(2));
    await waitFor(() => expect(result.current.operationError).toBeNull());

    act(() => { rejectSecond({ error: { code: 'settings_update_failed', retryable: true, request_id: 'new-request' } }); });
    await waitFor(() => expect(result.current.operationError?.requestId).toBe('new-request'));
  });

  it('keeps an independent connector failure visible while a settings save is newer', async () => {
    vi.mocked(settingsAPI.listItems).mockResolvedValue([]);
    vi.mocked(settingsAPI.updateKey).mockResolvedValue({ key: 'setting', value: null, description: '', updated_at: '' });
    let rejectConnector!: (error: unknown) => void;
    vi.mocked(dbTargetAPI.updateConnector).mockReturnValueOnce(new Promise((_, reject) => {
      rejectConnector = reject;
    }));
    const { result } = renderHook(() => useSettingsOperations(INITIAL_STATE, vi.fn()), { wrapper });

    act(() => { result.current.onUpdateConnector('conn-prod', { name: 'new name' }); });
    act(() => { void result.current.onSave(); });
    act(() => { rejectConnector({ error: { code: 'connector_update_failed', retryable: true, request_id: 'connector-request' } }); });

    await waitFor(() => expect(result.current.operationError?.requestId).toBe('connector-request'));
  });

  it('releases loading after a stale save settles and allows saving the latest state', async () => {
    vi.mocked(settingsAPI.listItems).mockResolvedValue([]);
    let resolveFirst!: (value: { key: string; value: unknown; description: string; updated_at: string }) => void;
    const firstRequest = new Promise<{ key: string; value: unknown; description: string; updated_at: string }>((resolve) => {
      resolveFirst = resolve;
    });
    vi.mocked(settingsAPI.updateKey)
      .mockReturnValueOnce(firstRequest)
      .mockResolvedValue({ key: 'setting', value: null, description: '', updated_at: '' });

    const { result } = renderHook(() => {
      const [state, dispatch] = useReducer(workbenchV2Reducer, INITIAL_STATE);
      return { state, operations: useSettingsOperations(state, dispatch) };
    }, { wrapper });

    act(() => { void result.current.operations.onSave(); });
    expect(result.current.operations.isSaving).toBe(true);
    act(() => { result.current.operations.onUpdateGeneral({ theme: 'light' }); });
    resolveFirst({ key: 'write_precision', value: null, description: '', updated_at: '' });

    await waitFor(() => expect(result.current.operations.isSaving).toBe(false));
    const callsAfterStaleSave = vi.mocked(settingsAPI.updateKey).mock.calls.length;
    act(() => { void result.current.operations.onSave(); });
    await waitFor(() => expect(vi.mocked(settingsAPI.updateKey).mock.calls.length).toBeGreaterThan(callsAfterStaleSave));
    expect(vi.mocked(settingsAPI.updateKey).mock.calls.at(-1)).toEqual(['enable_audit_log', true]);
    expect(vi.mocked(settingsAPI.updateKey).mock.calls.some(([key, value]) => key === 'theme' && value === 'light')).toBe(true);
  });

  it('serializes Save A and Save B so A cannot complete after B', async () => {
    vi.mocked(settingsAPI.listItems).mockResolvedValue([]);
    let resolveA!: (value: { key: string; value: unknown; description: string; updated_at: string }) => void;
    let resolveB!: (value: { key: string; value: unknown; description: string; updated_at: string }) => void;
    vi.mocked(settingsAPI.updateKey)
      .mockReturnValueOnce(new Promise((resolve) => { resolveA = resolve; }))
      .mockReturnValueOnce(new Promise((resolve) => { resolveB = resolve; }))
      .mockResolvedValue({ key: 'setting', value: null, description: '', updated_at: '' });
    const { result } = renderHook(() => {
      const [state, dispatch] = useReducer(workbenchV2Reducer, INITIAL_STATE);
      return { state, operations: useSettingsOperations(state, dispatch) };
    }, { wrapper });

    act(() => { void result.current.operations.onSave(); });
    act(() => { result.current.operations.onUpdateGeneral({ theme: 'light' }); });
    act(() => { void result.current.operations.onSave(); });
    await waitFor(() => expect(vi.mocked(settingsAPI.updateKey).mock.calls).toHaveLength(1));
    resolveA({ key: 'write_precision', value: null, description: '', updated_at: '' });
    await waitFor(() => expect(vi.mocked(settingsAPI.updateKey).mock.calls).toHaveLength(2));
    expect(result.current.operations.isSaving).toBe(true);
    resolveB({ key: 'write_precision', value: null, description: '', updated_at: '' });
    await waitFor(() => expect(result.current.operations.isSaving).toBe(false));
  });

  it('serializes rapid connector adds and keeps both successful rows', async () => {
    vi.mocked(settingsAPI.listItems).mockResolvedValue([]);
    let resolveFirst!: (value: ReturnType<typeof makeConnector>) => void;
    function makeConnector(id: string) {
      return { id, name: `Connector ${id}`, kind: 'postgres' as const, connection_config: {}, enabled: true, status: 'ready' as const, last_check_error: '', created_at: '', updated_at: '', default_write_interval_seconds: 5 };
    }
    vi.mocked(dbTargetAPI.createConnector)
      .mockReturnValueOnce(new Promise((resolve) => { resolveFirst = resolve; }))
      .mockResolvedValueOnce(makeConnector('conn-2'));
    const { result } = renderHook(() => {
      const [state, dispatch] = useReducer(workbenchV2Reducer, INITIAL_STATE);
      return { state, operations: useSettingsOperations(state, dispatch) };
    }, { wrapper });

    const first = result.current.operations.onAddConnector();
    const second = result.current.operations.onAddConnector();
    await waitFor(() => expect(vi.mocked(dbTargetAPI.createConnector).mock.calls).toHaveLength(1));
    await act(async () => {
      resolveFirst(makeConnector('conn-1'));
      await first;
      await second;
    });
    expect(vi.mocked(dbTargetAPI.createConnector).mock.calls).toHaveLength(2);
    expect(result.current.state.settings.connectors.map((connector) => connector.id)).toEqual(['conn-prod', 'conn-1', 'conn-2']);
  });

  it('does not resurrect a deleted connector on a later add', async () => {
    vi.mocked(settingsAPI.listItems).mockResolvedValue([]);
    function makeConnector(id: string) {
      return { id, name: `Connector ${id}`, kind: 'postgres' as const, connection_config: {}, enabled: true, status: 'ready' as const, last_check_error: '', created_at: '', updated_at: '', default_write_interval_seconds: 5 };
    }
    vi.mocked(dbTargetAPI.createConnector)
      .mockResolvedValueOnce(makeConnector('conn-a'))
      .mockResolvedValueOnce(makeConnector('conn-b'));
    vi.mocked(dbTargetAPI.deleteConnector).mockResolvedValue(undefined as never);
    const { result } = renderHook(() => {
      const [state, dispatch] = useReducer(workbenchV2Reducer, INITIAL_STATE);
      return { state, operations: useSettingsOperations(state, dispatch) };
    }, { wrapper });

    await act(async () => { await result.current.operations.onAddConnector(); });
    await waitFor(() => expect(result.current.state.settings.connectors.map((connector) => connector.id)).toContain('conn-a'));

    act(() => { result.current.operations.onRemoveConnector('conn-a'); });
    await waitFor(() => expect(result.current.state.settings.connectors.map((connector) => connector.id)).not.toContain('conn-a'));

    await act(async () => { await result.current.operations.onAddConnector(); });

    await waitFor(() => expect(result.current.state.settings.connectors.map((connector) => connector.id)).toContain('conn-b'));
    expect(result.current.state.settings.connectors.map((connector) => connector.id)).not.toContain('conn-a');
  });

  it('merges same-render connector patches into the second request payload', async () => {
    vi.mocked(settingsAPI.listItems).mockResolvedValue([]);
    let resolveFirst!: (value: ReturnType<typeof makeConnector>) => void;
    function makeConnector(name: string) {
      return { id: 'conn-prod', name, kind: 'postgres' as const, connection_config: { host: 'new-host' }, enabled: true, status: 'ready' as const, last_check_error: '', created_at: '', updated_at: '', default_write_interval_seconds: 5 };
    }
    vi.mocked(dbTargetAPI.updateConnector)
      .mockReturnValueOnce(new Promise((resolve) => { resolveFirst = resolve; }))
      .mockResolvedValueOnce(makeConnector('second'));
    const { result } = renderHook(() => useSettingsOperations(INITIAL_STATE, vi.fn()), { wrapper });

    act(() => {
      result.current.onUpdateConnector('conn-prod', { name: 'first' });
      result.current.onUpdateConnector('conn-prod', { host: 'new-host' } as never);
    });
    await waitFor(() => expect(dbTargetAPI.updateConnector).toHaveBeenCalledTimes(2));
    expect(dbTargetAPI.updateConnector).toHaveBeenNthCalledWith(2, 'conn-prod', expect.objectContaining({ name: 'first' }));
    resolveFirst(makeConnector('first'));
  });
});
