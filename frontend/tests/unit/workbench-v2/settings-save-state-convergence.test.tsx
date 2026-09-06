import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { act, renderHook, waitFor } from '@testing-library/react';
import { useReducer } from 'react';
import type { ReactNode } from 'react';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import { useSettingsOperations } from '../../../src/features/datalink/workbench-v2/settings/useSettingsOperations';
import { INITIAL_STATE, workbenchV2Reducer } from '../../../src/features/datalink/workbench-v2/state/useWorkbenchV2State';
import { settingsAPI } from '../../../src/services/datalink';

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

type SettingResponse = { key: string; value: unknown; description: string; updated_at: string };

function renderSettings() {
  return renderHook(() => {
    const [state, dispatch] = useReducer(workbenchV2Reducer, INITIAL_STATE);
    return { state, operations: useSettingsOperations(state, dispatch) };
  }, { wrapper });
}

describe('Settings save state reflects only real save requests', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    vi.mocked(settingsAPI.listItems).mockResolvedValue([]);
  });

  it('does not fake an in-flight save when a Modbus Share field is edited', async () => {
    const { result } = renderSettings();

    act(() => { result.current.operations.onUpdateModbusShare({ port: 5030 }); });

    expect(result.current.state.settings.modbus_share.port).toBe(5030);
    expect(result.current.state.settings.modbus_share.save_state).not.toBe('saving');
    expect(settingsAPI.updateKey).not.toHaveBeenCalled();
  });

  it('clears a recorded save error when a Modbus Share field is edited', async () => {
    vi.mocked(settingsAPI.updateKey).mockRejectedValue({ error: { code: 'settings_update_failed', retryable: true, request_id: 'failed-save' } });
    const { result } = renderSettings();

    await act(async () => { await result.current.operations.onSave(); });
    await waitFor(() => expect(result.current.state.settings.modbus_share.save_state).toBe('save-error'));

    act(() => { result.current.operations.onUpdateModbusShare({ port: 5031 }); });

    expect(result.current.state.settings.modbus_share.save_error).toBeNull();
    expect(result.current.state.settings.modbus_share.save_state).not.toBe('saving');
  });

  it('does not fake an in-flight save when settings are reset to defaults', () => {
    const { result } = renderSettings();

    act(() => { result.current.operations.onReset(); });

    expect(result.current.state.settings.modbus_share.save_state).not.toBe('saving');
    expect(settingsAPI.updateKey).not.toHaveBeenCalled();
  });

  it('moves through saving and lands on saved for a dispatched save', async () => {
    let resolveSave!: (value: SettingResponse) => void;
    vi.mocked(settingsAPI.updateKey)
      .mockReturnValueOnce(new Promise<SettingResponse>((resolve) => { resolveSave = resolve; }))
      .mockResolvedValue({ key: 'setting', value: null, description: '', updated_at: '' });
    const { result } = renderSettings();

    act(() => { void result.current.operations.onSave(); });
    expect(result.current.state.settings.modbus_share.save_state).toBe('saving');

    await act(async () => { resolveSave({ key: 'write_precision', value: null, description: '', updated_at: '' }); });
    await waitFor(() => expect(result.current.state.settings.modbus_share.save_state).toBe('saved'));
    expect(result.current.state.settings.modbus_share.save_error).toBeNull();
  });
});

describe('Settings save always converges to a terminal state', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    vi.mocked(settingsAPI.listItems).mockResolvedValue([]);
  });

  /** 讓 modbus_share 這筆 entry 的回應可由測試控制，其餘 entry 立即成功。 */
  function deferModbusShareEntry() {
    let settle!: { resolve: (value: SettingResponse) => void; reject: (error: unknown) => void };
    const deferred = new Promise<SettingResponse>((resolve, reject) => { settle = { resolve, reject }; });
    vi.mocked(settingsAPI.updateKey).mockImplementation((key: string) => (
      key === 'modbus_share'
        ? deferred
        : Promise.resolve({ key, value: null, description: '', updated_at: '' })
    ));
    return settle;
  }

  it('surfaces a failure even when the operator edits another field mid-flight', async () => {
    const settle = deferModbusShareEntry();
    const { result } = renderSettings();

    act(() => { void result.current.operations.onSave(); });
    await waitFor(() => expect(settingsAPI.updateKey).toHaveBeenCalledWith('modbus_share', expect.anything()));
    act(() => { result.current.operations.onUpdateGeneral({ theme: 'light' }); });

    await act(async () => {
      settle.reject({ error: { code: 'settings_update_failed', retryable: true, request_id: 'superseded-save' } });
      await Promise.resolve();
    });

    await waitFor(() => expect(result.current.state.settings.modbus_share.save_state).toBe('save-error'));
    expect(result.current.state.settings.modbus_share.save_error).toBeTruthy();
    expect(result.current.operations.operationError?.requestId).toBe('superseded-save');

    const callsBeforeRetry = vi.mocked(settingsAPI.updateKey).mock.calls.length;
    act(() => { result.current.operations.onRetry(); });
    await waitFor(() => expect(vi.mocked(settingsAPI.updateKey).mock.calls.length).toBeGreaterThan(callsBeforeRetry));
  });

  it('reaches the saved state even when the operator edits another field mid-flight', async () => {
    const settle = deferModbusShareEntry();
    const { result } = renderSettings();

    act(() => { void result.current.operations.onSave(); });
    await waitFor(() => expect(settingsAPI.updateKey).toHaveBeenCalledWith('modbus_share', expect.anything()));
    act(() => { result.current.operations.onUpdateGeneral({ theme: 'light' }); });

    await act(async () => {
      settle.resolve({ key: 'modbus_share', value: null, description: '', updated_at: '' });
      await Promise.resolve();
    });

    await waitFor(() => expect(result.current.state.settings.modbus_share.save_state).toBe('saved'));
    expect(result.current.state.settings.modbus_share.save_error).toBeNull();
  });

  it('adopts the server revision even when the operator edits mid-flight', async () => {
    const settle = deferModbusShareEntry();
    const { result } = renderSettings();

    act(() => { void result.current.operations.onSave(); });
    await waitFor(() => expect(settingsAPI.updateKey).toHaveBeenCalledWith('modbus_share', expect.anything()));
    act(() => { result.current.operations.onUpdateGeneral({ theme: 'light' }); });

    await act(async () => {
      settle.resolve({ key: 'modbus_share', value: { new_settings_revision: 'rev-2' }, description: '', updated_at: '' });
      await Promise.resolve();
    });

    await waitFor(() => expect(result.current.state.settings.modbus_share.save_state).toBe('saved'));
    // revision 是伺服器指派的 CAS 憑證：漏採用會讓後續存檔永久撞上 revision 衝突。
    expect(result.current.state.settings.modbus_share.settings_revision).toBe('rev-2');
    expect(result.current.state.settings.modbus_share.expected_settings_revision).toBe('rev-2');
  });

  it('stays in flight while a coalesced second save is still queued', async () => {
    // 每一次 modbus_share 請求都可獨立控制，才能觀察兩批之間的狀態窗口。
    const resolvers: Array<(value: SettingResponse) => void> = [];
    vi.mocked(settingsAPI.updateKey).mockImplementation((key: string) => {
      if (key !== 'modbus_share') {
        return Promise.resolve({ key, value: null, description: '', updated_at: '' });
      }
      return new Promise<SettingResponse>((resolve) => { resolvers.push(resolve); });
    });
    const { result } = renderSettings();

    act(() => { void result.current.operations.onSave(); });
    await waitFor(() => expect(resolvers).toHaveLength(1));
    act(() => { void result.current.operations.onSave(); });

    // 第一批完成，第二批仍在佇列中：此時不得宣告 saved。
    await act(async () => {
      resolvers[0]({ key: 'modbus_share', value: null, description: '', updated_at: '' });
      await Promise.resolve();
    });
    await waitFor(() => expect(resolvers).toHaveLength(2));
    expect(result.current.state.settings.modbus_share.save_state).toBe('saving');
    expect(result.current.operations.isSaving).toBe(true);

    await act(async () => {
      resolvers[1]({ key: 'modbus_share', value: null, description: '', updated_at: '' });
      await Promise.resolve();
    });
    await waitFor(() => expect(result.current.state.settings.modbus_share.save_state).toBe('saved'));
    expect(result.current.operations.isSaving).toBe(false);
  });

  it('sends the adopted revision on the next save instead of a stale one', async () => {
    vi.mocked(settingsAPI.updateKey).mockImplementation((key: string) => Promise.resolve(
      key === 'modbus_share'
        ? { key, value: { new_settings_revision: 'rev-2' }, description: '', updated_at: '' }
        : { key, value: null, description: '', updated_at: '' },
    ));
    const { result } = renderSettings();

    await act(async () => { await result.current.operations.onSave(); });
    act(() => { result.current.operations.onUpdateGeneral({ theme: 'light' }); });
    await act(async () => { await result.current.operations.onSave(); });

    const shareCalls = vi.mocked(settingsAPI.updateKey).mock.calls.filter(([key]) => key === 'modbus_share');
    expect(shareCalls).toHaveLength(2);
    expect(shareCalls.at(-1)?.[1]).toEqual(expect.objectContaining({
      settings_revision: 'rev-2',
      expected_settings_revision: 'rev-2',
    }));
  });
});
