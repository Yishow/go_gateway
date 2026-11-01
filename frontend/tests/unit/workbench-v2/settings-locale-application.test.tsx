import { DiagnosticsSection } from '@/features/datalink/workbench-v2/settings/DiagnosticsSection';
import { SchedulerSection } from '@/features/datalink/workbench-v2/settings/SchedulerSection';
import { ModbusShareSection } from '@/features/datalink/workbench-v2/settings/ModbusShareSection';
import { ConnectorRow } from '@/features/datalink/workbench-v2/settings/ConnectorRow';
import type { SettingsConnector } from '@/features/datalink/workbench-v2/state/types';
import { useReducer } from 'react';
import { act, cleanup, fireEvent, render, screen, waitFor } from '@testing-library/react';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { afterEach, beforeEach, expect, it, vi } from 'vitest';
import { settingsKeys } from '@/hooks/datalink/keys';
import { settingsAPI, dbTargetAPI } from '@/services/datalink';
import { UiSection } from '@/features/datalink/workbench-v2/settings/UiSection';
import { SettingsPage } from '@/features/datalink/workbench-v2/settings/SettingsPage';
import { useSettingsOperations } from '@/features/datalink/workbench-v2/settings/useSettingsOperations';
import { INITIAL_STATE, workbenchV2Reducer } from '@/features/datalink/workbench-v2/state/useWorkbenchV2State';
const locale = vi.hoisted(() => {
  const listeners = new Set<() => void>();
  const owner = { language: 'zh-TW', on: (_: string, fn: () => void) => { listeners.add(fn); },
    off: (_: string, fn: () => void) => { listeners.delete(fn); },
    changeLanguage: vi.fn(async (next: string) => { owner.language = next; listeners.forEach(fn => fn()); }), listeners };
  return owner;
});
vi.mock('react-i18next', () => ({ useTranslation: () => ({ i18n: locale, t: (key: string, fallback?: unknown) => typeof fallback === 'string' ? fallback : key }) }));
vi.mock('@/services/datalink', () => ({ settingsAPI: { listItems: vi.fn(), updateKey: vi.fn() }, dbTargetAPI: { listConnectors: vi.fn() } }));
beforeEach(() => {
  vi.clearAllMocks(); locale.language = 'zh-TW'; locale.listeners.clear();
  vi.mocked(settingsAPI.listItems).mockResolvedValue([]); vi.mocked(dbTargetAPI.listConnectors).mockResolvedValue([]);
  vi.mocked(settingsAPI.updateKey).mockImplementation(async (key, value) => ({ key, value, description: '', updated_at: '' }));
});
afterEach(cleanup);
function Form({ page = false }: { page?: boolean }) {
  const [state, dispatch] = useReducer(workbenchV2Reducer, INITIAL_STATE);
  const operations = useSettingsOperations(state, dispatch);
  return page ? <SettingsPage state={state} dispatch={dispatch} /> : <>
    <UiSection settings={state.settings.general} onChange={operations.onUpdateGeneral} />
    <button onClick={() => void operations.onSave()}>Save settings</button>
    <button onClick={operations.onReset}>Reset draft</button>
    <output data-testid="locale-draft">{state.settings.general.locale}</output>
    <output data-testid="save-state">{operations.isSaving ? 'pending' : 'idle'}</output>
  </>;
}
function mount(page = false, client = new QueryClient({ defaultOptions: { queries: { retry: false }, mutations: { retry: false } } })) { return render(<QueryClientProvider client={client}><Form page={page} /></QueryClientProvider>); }
function choose(value: string) { fireEvent.change(screen.getAllByRole('combobox')[1], { target: { value } }); }
async function save() { fireEvent.click(screen.getByText('Save settings')); await waitFor(() => expect(screen.getByTestId('save-state')).toHaveTextContent('idle')); }
it('keeps locale editing and reset draft-only, applying only the saved locale ACK', async () => {
  mount(); choose('en'); expect(screen.getByTestId('locale-draft')).toHaveTextContent('en');
  expect(locale.changeLanguage).not.toHaveBeenCalled(); expect(settingsAPI.updateKey).not.toHaveBeenCalled();
  fireEvent.click(screen.getByText('Reset draft')); expect(locale.changeLanguage).not.toHaveBeenCalled();
  choose('en'); await save(); expect(locale.changeLanguage).toHaveBeenCalledExactlyOnceWith('en');
});
it('does not change the active language when locale persistence fails', async () => {
  vi.mocked(settingsAPI.updateKey).mockImplementation(async (key, value) => { if (key === 'locale') throw new Error('failed'); return { key, value, description: '', updated_at: '' }; });
  mount(); choose('en'); await save(); expect(locale.changeLanguage).not.toHaveBeenCalled(); expect(locale.language).toBe('zh-TW');
});
it('keeps acknowledged locale when a subsequent setting fails rather than pretending the partial save rolled back', async () => {
  vi.mocked(settingsAPI.updateKey).mockImplementation(async (key, value) => { if (key === 'addr_format') throw new Error('later failure'); return { key, value, description: '', updated_at: '' }; });
  mount(); choose('en'); await save(); expect(locale.changeLanguage).toHaveBeenCalledExactlyOnceWith('en');
});
it('ignores a locale ACK from an unmounted settings owner after a newer page saved its choice', async () => {
  let release!: () => void; let held = false;
  vi.mocked(settingsAPI.updateKey).mockImplementation(async (key, value) => { if (key === 'locale' && !held) { held = true; await new Promise<void>(resolve => { release = resolve; }); } return { key, value, description: '', updated_at: '' }; });
  const first = mount(); choose('en'); fireEvent.click(screen.getByText('Save settings')); await waitFor(() => expect(release).toBeTypeOf('function'));
  first.unmount(); mount(); await save(); expect(locale.changeLanguage).toHaveBeenCalledExactlyOnceWith('zh-TW');
  await act(async () => release()); expect(locale.changeLanguage).toHaveBeenCalledTimes(1);
});
it('applies an actual stored locale at hydration without writing settings', async () => {
  vi.mocked(settingsAPI.listItems).mockResolvedValue([{ key: 'locale', value: 'en', description: '', updated_at: '' }]);
  mount(true); await screen.findByTestId('settings-page'); await waitFor(() => expect(locale.changeLanguage).toHaveBeenCalledWith('en'));
  expect(settingsAPI.updateKey).not.toHaveBeenCalled();
});
it.each([{ items: [] }, { items: [{ key: 'locale', value: 'invalid', description: '', updated_at: '' }] }])('does not treat absent or invalid stored locale as a saved default', async ({ items }) => {
  vi.mocked(settingsAPI.listItems).mockResolvedValue(items); mount(true); await screen.findByTestId('settings-page'); expect(locale.changeLanguage).not.toHaveBeenCalled();
});
it('waits for initial cached settings to refresh, then never rehydrates over an editable draft', async () => {
  const client = new QueryClient({ defaultOptions: { queries: { retry: false } } });
  client.setQueryData(settingsKeys.items(), [{ key: 'locale', value: 'en', description: '', updated_at: '' }]);
  let release!: () => void;
  vi.mocked(settingsAPI.listItems).mockImplementationOnce(() => new Promise(resolve => { release = () => resolve([{ key: 'locale', value: 'zh-TW', description: '', updated_at: '' }]); }));
  mount(true, client); await waitFor(() => expect(release).toBeTypeOf('function'));
  await act(async () => { await Promise.resolve(); });
  expect(screen.queryByTestId('settings-page')).not.toBeInTheDocument(); expect(locale.changeLanguage).not.toHaveBeenCalled();
  await act(async () => release()); await screen.findByTestId('settings-page');
  expect(locale.changeLanguage).toHaveBeenCalledExactlyOnceWith('zh-TW');
  fireEvent.change(screen.getByLabelText('顯示語系'), { target: { value: 'en' } });
  vi.mocked(settingsAPI.listItems).mockResolvedValue([{ key: 'locale', value: 'zh-TW', description: '', updated_at: '' }]);
  await act(async () => { await client.invalidateQueries({ queryKey: settingsKeys.items() }); });
  expect(screen.getByLabelText('顯示語系')).toHaveValue('en'); expect(locale.changeLanguage).toHaveBeenCalledTimes(1);
});
it('uses shared semantic colors on the Settings surfaces needed to change language and theme (not a pixel assertion)', async () => {
  const { readFileSync } = await import('node:fs');
  for (const name of ['SettingsHeader', 'ConnectorPoolSection', 'ConnectorRow', 'TimeseriesSection', 'SchedulerSection', 'ModbusShareSection', 'UiSection', 'ApiSection', 'DiagnosticsSection', 'SaveBar', 'SettingsStatus']) {
    const source = readFileSync(`${process.cwd()}/src/features/datalink/workbench-v2/settings/${name}.tsx`, 'utf8');
    expect(source, name).not.toMatch(/(?:text|bg|border)-(?:gray|slate|blue|cyan|emerald|red)-\d|text-white|bg-gradient|backdrop-blur/);
    if (name === 'ConnectorPoolSection') {
      expect(source).toContain('bg-[var(--gw-primary)] hover:bg-[var(--gw-primary-hover)] text-[var(--gw-on-primary)]');
    }
  }
});
it('keeps narrow Settings connector actions and names bounded with a separate touch target and switch track (not a pixel assertion)', async () => {
  const { readFileSync } = await import('node:fs');
  const source = (name: string) => readFileSync(`${process.cwd()}/src/features/datalink/workbench-v2/settings/${name}.tsx`, 'utf8');
  const pool = source('ConnectorPoolSection');
  const row = source('ConnectorRow');
  expect(pool).toContain('flex flex-col sm:flex-row');
  expect(pool).toContain('shrink-0 whitespace-nowrap');
  expect(row).toContain('flex min-w-0 flex-1 flex-wrap items-center gap-2');
  expect(row).toContain('min-w-0 w-full sm:w-auto');
  expect(row).toContain('inline-flex h-11 w-11 shrink-0 items-center justify-center');
  expect(row).toContain('relative inline-flex h-4 w-7 shrink-0 items-center rounded-full');
});
it('exposes connector switch name and checked state without changing its update contract', () => {
  const connector: SettingsConnector = { id: 'test', name: 'Local archive', kind: 'sqlite', host: '', port: 0, database: ':memory:', username: '', schema: '', table: '', enabled: true, status: 'unknown', default_write_interval_seconds: 1 };
  const onUpdate = vi.fn();
  const view = render(<ConnectorRow connector={connector} onUpdate={onUpdate} onRemove={vi.fn()} onTest={vi.fn()} />);
  const toggle = screen.getByRole('switch', { name: '啟用 Local archive', checked: true });
  fireEvent.click(toggle); expect(onUpdate).toHaveBeenCalledExactlyOnceWith({ enabled: false });
  view.rerender(<ConnectorRow connector={{ ...connector, enabled: false, status: 'testing' }} onUpdate={onUpdate} onRemove={vi.fn()} onTest={vi.fn()} />);
  expect(screen.getByRole('switch', { name: '啟用 Local archive', checked: false })).toBeDisabled();
});
it.each([
  { key: 'enable_audit_log' as const, label: '記錄操作審計日誌' },
  { key: 'enable_debug_panel' as const, label: '顯示開發者調試面板' },
  { key: 'auto_start' as const, label: '開機自動啟動' },
])('uses bounded shared Settings switch for $key with the existing boolean patch', ({ key, label }) => {
  const onChange = vi.fn();
  const component = (checked: boolean) => key === 'auto_start'
    ? <SchedulerSection settings={{ ...INITIAL_STATE.settings.scheduler, auto_start: checked }} onChange={onChange} />
    : <DiagnosticsSection settings={{ ...INITIAL_STATE.settings.general, [key]: checked }} onChange={onChange} />;
  const view = render(component(false));
  const toggle = screen.getByRole('switch', { name: label, checked: false });
  expect(toggle.firstElementChild).toHaveStyle({ width: '40px', height: '20px' });
  fireEvent.click(toggle); expect(onChange).toHaveBeenLastCalledWith({ [key]: true });
  view.rerender(component(true));
  fireEvent.click(screen.getByRole('switch', { name: label, checked: true }));
  expect(onChange).toHaveBeenLastCalledWith({ [key]: false });
});
it('uses bounded shared Settings switch for Share while preserving pending and repeat guards', async () => {
  const onChange = vi.fn();
  const component = (enabled: boolean, pending: boolean) => <ModbusShareSection settings={{ ...INITIAL_STATE.settings.modbus_share, enabled }} onChange={onChange} pending={pending} />;
  const view = render(component(false, true));
  const toggle = screen.getByRole('switch', { name: 'Local Modbus Share', checked: false });
  expect(toggle).toBeDisabled(); expect(toggle.parentElement).toHaveAttribute('aria-busy', 'true');
  fireEvent.click(toggle); expect(onChange).not.toHaveBeenCalled();
  view.rerender(component(false, false));
  expect(toggle.firstElementChild).toHaveStyle({ width: '40px', height: '20px' });
  fireEvent.click(toggle); fireEvent.click(toggle);
  expect(onChange).toHaveBeenCalledExactlyOnceWith({ enabled: true });
  await act(async () => { await Promise.resolve(); });
  view.rerender(component(true, false)); fireEvent.click(screen.getByRole('switch', { name: 'Local Modbus Share', checked: true }));
  expect(onChange).toHaveBeenLastCalledWith({ enabled: false });
});
