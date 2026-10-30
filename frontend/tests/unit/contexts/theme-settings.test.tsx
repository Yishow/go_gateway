import { useReducer } from 'react';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { act, cleanup, fireEvent, render, screen, waitFor } from '@testing-library/react';
import { afterEach, beforeEach, expect, it, vi } from 'vitest';
import { ThemeProvider, useTheme } from '@/contexts/ThemeContext';
import { UiSection } from '@/features/datalink/workbench-v2/settings/UiSection';
import { useSettingsOperations } from '@/features/datalink/workbench-v2/settings/useSettingsOperations';
import { INITIAL_STATE, workbenchV2Reducer } from '@/features/datalink/workbench-v2/state/useWorkbenchV2State';
import { SettingsPage } from '@/features/datalink/workbench-v2/settings/SettingsPage';
import { ThemeToggle } from '@/components/ThemeToggle';
import { settingsAPI, dbTargetAPI } from '@/services/datalink';
vi.mock('react-i18next', () => ({ useTranslation: () => ({ t: (key: string, fallback?: unknown) => typeof fallback === 'string' ? fallback : key }) }));
vi.mock('@/services/datalink', () => ({ settingsAPI: { listItems: vi.fn(), updateKey: vi.fn() }, dbTargetAPI: { listConnectors: vi.fn() } }));
afterEach(cleanup);
beforeEach(() => {
  vi.clearAllMocks(); localStorage.clear();
  vi.mocked(settingsAPI.listItems).mockResolvedValue([]);
  vi.mocked(dbTargetAPI.listConnectors).mockResolvedValue([]);
  vi.mocked(settingsAPI.updateKey).mockImplementation(async (key, value) => ({ key, value, description: '', updated_at: '' }));
});
function Harness() {
  const [state, dispatch] = useReducer(workbenchV2Reducer, { ...INITIAL_STATE, settings: { ...INITIAL_STATE.settings,
    general: { ...INITIAL_STATE.settings.general, theme: 'dark' as const }, scheduler: { ...INITIAL_STATE.settings.scheduler, default_interval_ms: 9876 },
    modbus_share: { ...INITIAL_STATE.settings.modbus_share, port: 6550 } } });
  const operations = useSettingsOperations(state, dispatch); const theme = useTheme();
  return <><UiSection settings={state.settings.general} onChange={operations.onUpdateGeneral} />
    <button onClick={() => theme.setPreference('light')}>Header light</button><button onClick={() => { void operations.onSave(); }}>Whole save</button>
    <output data-testid="draft">{state.settings.general.theme}/{state.settings.scheduler.default_interval_ms}/{state.settings.modbus_share.port}</output></>;
}
function mount() { render(<QueryClientProvider client={new QueryClient({ defaultOptions: { queries: { retry: false }, mutations: { retry: false } } })}><ThemeProvider><Harness /></ThemeProvider></QueryClientProvider>); }
it('theme-only UI selection updates general draft and writes only theme, leaving other drafts untouched', async () => {
  mount(); fireEvent.change(screen.getByLabelText('視覺主題'), { target: { value: 'light' } });
  await waitFor(() => expect(settingsAPI.updateKey).toHaveBeenCalledWith('theme', 'light'));
  expect(vi.mocked(settingsAPI.updateKey).mock.calls).toEqual([['theme', 'light']]);
  expect(screen.getByTestId('draft')).toHaveTextContent('light/9876/6550');
});
it('global explicit choice cannot be overwritten by a later full Save carrying an old dark draft', async () => {
  mount(); fireEvent.click(screen.getByText('Header light')); await waitFor(() => expect(settingsAPI.updateKey).toHaveBeenCalledWith('theme', 'light'));
  expect(screen.getByTestId('draft')).toHaveTextContent('dark/9876/6550');
  fireEvent.click(screen.getByText('Whole save'));
  await waitFor(() => expect(vi.mocked(settingsAPI.updateKey).mock.calls.filter(([key]) => key === 'theme')).toHaveLength(2));
  expect(vi.mocked(settingsAPI.updateKey).mock.calls.filter(([key]) => key === 'theme').map(([, value]) => value)).toEqual(['light', 'light']);
  expect(vi.mocked(settingsAPI.updateKey).mock.calls.some(([key]) => key === 'modbus_share')).toBe(true);
  expect(vi.mocked(settingsAPI.updateKey).mock.calls.some(([key]) => key === 'default_interval_ms')).toBe(true);
});
it('pending/error stays visible and retry only writes theme without flushing drafts', async () => {
  vi.mocked(settingsAPI.updateKey).mockRejectedValueOnce(new Error('secret'));
  mount(); fireEvent.change(screen.getByLabelText('視覺主題'), { target: { value: 'light' } });
  expect(screen.getByText('正在儲存主題偏好…')).toHaveAttribute('role', 'status');
  await waitFor(() => expect(screen.getByRole('alert')).toHaveTextContent('未儲存'));
  fireEvent.click(screen.getByRole('button', { name: '重試儲存主題' }));
  await waitFor(() => expect(screen.queryByRole('alert')).not.toBeInTheDocument());
  expect(vi.mocked(settingsAPI.updateKey).mock.calls).toEqual([['theme', 'light'], ['theme', 'light']]);
});

function ActualPage() {
  const [state, dispatch] = useReducer(workbenchV2Reducer, INITIAL_STATE);
  return <><ThemeToggle /><SettingsPage state={state} dispatch={dispatch} /><output data-testid="page-theme">{state.settings.general.theme}</output></>;
}
it('actual SettingsPage hydration ignores injected dark and keeps header choice in the general draft', async () => {
  localStorage.setItem('theme', 'light');
  render(<QueryClientProvider client={new QueryClient({ defaultOptions: { queries: { retry: false } } })}><ThemeProvider><ActualPage /></ThemeProvider></QueryClientProvider>);
  await waitFor(() => expect(screen.getByLabelText('視覺主題')).toHaveValue('light'));
  expect(screen.getByTestId('page-theme')).toHaveTextContent('light'); expect(settingsAPI.updateKey).not.toHaveBeenCalled();
  fireEvent.click(screen.getByRole('button', { name: 'layout.toggleTheme' }));
  await waitFor(() => expect(screen.getByLabelText('視覺主題')).toHaveValue('dark'));
  expect(screen.getByTestId('page-theme')).toHaveTextContent('dark');
  expect(vi.mocked(settingsAPI.updateKey).mock.calls).toEqual([['theme', 'dark']]);
});

it('a theme choice during an in-flight full Save is not reverted when that old snapshot reaches its theme entry', async () => {
  let release!: () => void; let held = false;
  vi.mocked(settingsAPI.updateKey).mockImplementation((key, value) => {
    if (key !== 'theme' && !held) { held = true; return new Promise(resolve => { release = () => resolve({ key, value, description: '', updated_at: '' }); }); }
    return Promise.resolve({ key, value, description: '', updated_at: '' });
  });
  mount(); fireEvent.click(screen.getByText('Whole save')); await waitFor(() => expect(held).toBe(true));
  fireEvent.click(screen.getByText('Header light')); await waitFor(() => expect(settingsAPI.updateKey).toHaveBeenCalledWith('theme', 'light'));
  await act(async () => release());
  await waitFor(() => expect(vi.mocked(settingsAPI.updateKey).mock.calls.filter(([key]) => key === 'theme')).toHaveLength(2));
  expect(vi.mocked(settingsAPI.updateKey).mock.calls.filter(([key]) => key === 'theme').map(([, value]) => value)).toEqual(['light', 'light']);
});
function mountPage() { render(<QueryClientProvider client={new QueryClient({ defaultOptions: { queries: { retry: false }, mutations: { retry: false } } })}><ThemeProvider><ActualPage /></ThemeProvider></QueryClientProvider>); }
it('confirmed reset is explicit draft-only theme intent; a later full Save persists the default, not obsolete light', async () => {
  localStorage.setItem('theme', 'light'); mountPage(); await waitFor(() => expect(screen.getByLabelText('視覺主題')).toHaveValue('light'));
  const confirm = vi.spyOn(window, 'confirm').mockReturnValue(true);
  fireEvent.click(screen.getByRole('button', { name: '重設為預設' }));
  expect(screen.getByTestId('page-theme')).toHaveTextContent('dark'); expect(screen.getByLabelText('視覺主題')).toHaveValue('dark');
  expect(settingsAPI.updateKey).not.toHaveBeenCalled(); expect(localStorage.getItem('theme')).toBe('light');
  fireEvent.click(screen.getByRole('button', { name: '儲存所有設定' }));
  await waitFor(() => expect(settingsAPI.updateKey).toHaveBeenCalledWith('theme', 'dark')); confirm.mockRestore();
});
it('canceling reset changes neither global theme nor draft nor persistence', async () => {
  localStorage.setItem('theme', 'light'); mountPage(); await waitFor(() => expect(screen.getByLabelText('視覺主題')).toHaveValue('light'));
  const confirm = vi.spyOn(window, 'confirm').mockReturnValue(false); fireEvent.click(screen.getByRole('button', { name: '重設為預設' }));
  expect(screen.getByTestId('page-theme')).toHaveTextContent('light'); expect(screen.getByLabelText('視覺主題')).toHaveValue('light');
  expect(settingsAPI.updateKey).not.toHaveBeenCalled(); confirm.mockRestore();
});
it('a full Save started before reset cannot silently persist the newer reset draft', async () => {
  localStorage.setItem('theme', 'light'); let release!: () => void; let held = false;
  vi.mocked(settingsAPI.updateKey).mockImplementation((key, value) => {
    if (key !== 'theme' && !held) { held = true; return new Promise(resolve => { release = () => resolve({ key, value, description: '', updated_at: '' }); }); }
    return Promise.resolve({ key, value, description: '', updated_at: '' });
  });
  mountPage(); await waitFor(() => expect(screen.getByLabelText('視覺主題')).toHaveValue('light'));
  fireEvent.click(screen.getByRole('button', { name: '儲存所有設定' })); await waitFor(() => expect(held).toBe(true));
  const confirm = vi.spyOn(window, 'confirm').mockReturnValue(true); fireEvent.click(screen.getByRole('button', { name: '重設為預設' }));
  await act(async () => release());
  expect(screen.getByLabelText('視覺主題')).toHaveValue('dark');
  expect(vi.mocked(settingsAPI.updateKey).mock.calls.filter(([key]) => key === 'theme')).toEqual([]);
  confirm.mockRestore();
});
it('reset fences an already-issued theme reply, remains unsaved, and a later user choice still wins the next full Save', async () => {
  localStorage.setItem('theme', 'dark'); let release!: () => void;
  vi.mocked(settingsAPI.updateKey).mockImplementationOnce((key, value) => new Promise(resolve => { release = () => resolve({ key, value, description: '', updated_at: '' }); }));
  mountPage(); await waitFor(() => expect(screen.getByLabelText('視覺主題')).toHaveValue('dark'));
  fireEvent.change(screen.getByLabelText('視覺主題'), { target: { value: 'light' } });
  await waitFor(() => expect(settingsAPI.updateKey).toHaveBeenCalledWith('theme', 'light'));
  const store = vi.spyOn(Storage.prototype, 'setItem'); const confirm = vi.spyOn(window, 'confirm').mockReturnValue(true);
  fireEvent.click(screen.getByRole('button', { name: '重設為預設' }));
  expect(store).not.toHaveBeenCalled(); expect(screen.getByLabelText('視覺主題')).toHaveValue('dark');
  await act(async () => release());
  expect(screen.getByLabelText('視覺主題')).toHaveValue('dark'); expect(screen.getByText('主題偏好尚未儲存。')).toBeInTheDocument();
  expect(settingsAPI.updateKey).toHaveBeenCalledTimes(1);
  fireEvent.change(screen.getByLabelText('視覺主題'), { target: { value: 'light' } });
  await waitFor(() => expect(settingsAPI.updateKey).toHaveBeenCalledTimes(2));
  fireEvent.click(screen.getByRole('button', { name: '儲存所有設定' }));
  await waitFor(() => expect(vi.mocked(settingsAPI.updateKey).mock.calls.filter(([key]) => key === 'theme')).toHaveLength(3));
  expect(vi.mocked(settingsAPI.updateKey).mock.calls.filter(([key]) => key === 'theme').map(([, value]) => value)).toEqual(['light', 'light', 'light']);
  store.mockRestore(); confirm.mockRestore();
});
