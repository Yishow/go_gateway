import { StrictMode } from 'react';
import { act, cleanup, fireEvent, render, screen, waitFor } from '@testing-library/react';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';
import { ThemeProvider, useTheme } from '@/contexts/ThemeContext';
import { settingsAPI } from '@/services/datalink';
vi.mock('@/services/datalink', () => ({ settingsAPI: { listItems: vi.fn(), updateKey: vi.fn() } }));
let dark = true;
const listeners = new Set<() => void>();
function Probe() {
  const value = useTheme();
  return <><output data-testid="theme">{value.theme}/{value.preference}/{value.saveState}</output>
    {(['light', 'dark', 'system'] as const).map(preference => <button key={preference} onClick={() => value.setPreference(preference)}>{preference}</button>)}</>;
}
function mount(client = new QueryClient({ defaultOptions: { queries: { retry: false } } })) {
  return render(<QueryClientProvider client={client}><ThemeProvider><Probe /></ThemeProvider></QueryClientProvider>);
}
function os(value: boolean) { act(() => { dark = value; listeners.forEach(listener => listener()); }); }
beforeEach(() => {
  localStorage.clear(); dark = true; listeners.clear(); vi.clearAllMocks();
  vi.stubGlobal('matchMedia', vi.fn(() => ({ get matches() { return dark; }, addEventListener: (_: string, callback: () => void) => listeners.add(callback), removeEventListener: (_: string, callback: () => void) => listeners.delete(callback) })));
  vi.mocked(settingsAPI.listItems).mockResolvedValue([]);
  vi.mocked(settingsAPI.updateKey).mockImplementation(async (key, value) => ({ key, value, description: '', updated_at: '' }));
});
afterEach(() => { cleanup(); vi.unstubAllGlobals(); });
describe('one global preference owner', () => {
  it('absent raw theme follows live OS without saving resolved colors; explicit override and system are intentional writes only', async () => {
    mount(); await waitFor(() => expect(settingsAPI.listItems).toHaveBeenCalled());
    expect(screen.getByTestId('theme')).toHaveTextContent('dark/system/idle'); os(false);
    expect(screen.getByTestId('theme')).toHaveTextContent('light/system/idle');
    expect(settingsAPI.updateKey).not.toHaveBeenCalled(); expect(localStorage.getItem('theme')).toBeNull();
    fireEvent.click(screen.getByText('dark')); await waitFor(() => expect(screen.getByTestId('theme')).toHaveTextContent('dark/dark/saved'));
    os(true); os(false); expect(screen.getByTestId('theme')).toHaveTextContent('dark/dark/saved');
    fireEvent.click(screen.getByText('system')); await waitFor(() => expect(settingsAPI.updateKey).toHaveBeenLastCalledWith('theme', 'auto'));
    expect(screen.getByTestId('theme')).toHaveTextContent('light/system');
  });
  it('valid raw server preference wins legacy; absent raw key preserves legacy instead of mapper default dark', async () => {
    localStorage.setItem('theme', 'light'); mount(); await waitFor(() => expect(settingsAPI.listItems).toHaveBeenCalled());
    expect(screen.getByTestId('theme')).toHaveTextContent('light/light'); cleanup();
    vi.mocked(settingsAPI.listItems).mockResolvedValue([{ key: 'theme', value: 'dark', description: '', updated_at: '' }]); mount();
    await waitFor(() => expect(screen.getByTestId('theme')).toHaveTextContent('dark/dark'));
  });
  it('late hydration cannot overwrite a newer explicit choice', async () => {
    let resolve!: (items: Awaited<ReturnType<typeof settingsAPI.listItems>>) => void;
    vi.mocked(settingsAPI.listItems).mockReturnValue(new Promise(done => { resolve = done; })); mount();
    fireEvent.click(screen.getByText('light')); await act(async () => resolve([{ key: 'theme', value: 'dark', description: '', updated_at: '' }]));
    expect(screen.getByTestId('theme')).toHaveTextContent('light/light');
  });
  it('failed persistence is not saved success; API/storage/media absence remains usable', async () => {
    vi.stubGlobal('matchMedia', undefined); vi.mocked(settingsAPI.updateKey).mockRejectedValue(new Error('private endpoint data'));
    mount(); expect(screen.getByTestId('theme')).toHaveTextContent('light/system');
    fireEvent.click(screen.getByText('dark')); await waitFor(() => expect(screen.getByTestId('theme')).toHaveTextContent('dark/dark/error'));
    expect(screen.queryByText('private endpoint data')).not.toBeInTheDocument();
  });
});
it('versioned explicit local preference wins legacy, invalid raw values do not invent server intent', async () => {
  const { THEME_PREFERENCE_KEY } = await import('@/contexts/themePreference');
  localStorage.setItem('theme', 'light'); localStorage.setItem(THEME_PREFERENCE_KEY, JSON.stringify({ version: 1, preference: 'system' }));
  vi.mocked(settingsAPI.listItems).mockResolvedValue([{ key: 'theme', value: 'invented', description: '', updated_at: '' }]);
  mount(); await waitFor(() => expect(settingsAPI.listItems).toHaveBeenCalled());
  expect(screen.getByTestId('theme')).toHaveTextContent('dark/system/idle'); expect(settingsAPI.updateKey).not.toHaveBeenCalled();
});
it('a new actual app QueryClient fences late hydration from the previous context, while ordinary child navigation keeps the owner', async () => {
  let oldReply!: (items: Awaited<ReturnType<typeof settingsAPI.listItems>>) => void;
  vi.mocked(settingsAPI.listItems).mockReturnValueOnce(new Promise(resolve => { oldReply = resolve; })).mockResolvedValue([{ key: 'theme', value: 'light', description: '', updated_at: '' }]);
  const first = new QueryClient({ defaultOptions: { queries: { retry: false } } }); const second = new QueryClient({ defaultOptions: { queries: { retry: false } } });
  const view = mount(first);
  view.rerender(<QueryClientProvider client={second}><ThemeProvider><Probe /><div>Workspace B</div></ThemeProvider></QueryClientProvider>);
  await waitFor(() => expect(screen.getByTestId('theme')).toHaveTextContent('light/light'));
  await act(async () => oldReply([{ key: 'theme', value: 'dark', description: '', updated_at: '' }]));
  expect(screen.getByTestId('theme')).toHaveTextContent('light/light');
  const count = vi.mocked(settingsAPI.listItems).mock.calls.length;
  view.rerender(<QueryClientProvider client={second}><ThemeProvider><Probe /><div>Workspace C</div></ThemeProvider></QueryClientProvider>);
  expect(settingsAPI.listItems).toHaveBeenCalledTimes(count); expect(screen.getByTestId('theme')).toHaveTextContent('light/light');
});
it('rapid choices serialize theme writes and retain latest pending status until its own completion', async () => {
  const completions: Array<() => void> = [];
  vi.mocked(settingsAPI.updateKey).mockImplementation((key, value) => new Promise(resolve => completions.push(() => resolve({ key, value, description: '', updated_at: '' }))));
  mount(); fireEvent.click(screen.getByText('light')); await waitFor(() => expect(completions).toHaveLength(1));
  fireEvent.click(screen.getByText('dark')); expect(settingsAPI.updateKey).toHaveBeenCalledTimes(1);
  await act(async () => completions[0]()); await waitFor(() => expect(completions).toHaveLength(2));
  expect(screen.getByTestId('theme')).toHaveTextContent('dark/dark/pending');
  await act(async () => completions[1]()); expect(screen.getByTestId('theme')).toHaveTextContent('dark/dark/saved');
  expect(vi.mocked(settingsAPI.updateKey).mock.calls).toEqual([['theme', 'light'], ['theme', 'dark']]);
});
it('removes OS listeners and tolerates inaccessible storage and throwing matchMedia', async () => {
  const storage = vi.spyOn(Storage.prototype, 'getItem').mockImplementation(() => { throw new Error('denied'); });
  vi.stubGlobal('matchMedia', () => { throw new Error('unavailable'); }); mount();
  expect(screen.getByTestId('theme')).toHaveTextContent('light/system'); cleanup(); storage.mockRestore();
  vi.stubGlobal('matchMedia', () => ({ matches: false, addEventListener: (_: string, callback: () => void) => listeners.add(callback), removeEventListener: (_: string, callback: () => void) => listeners.delete(callback) }));
  mount(); expect(listeners.size).toBe(1); cleanup(); expect(listeners.size).toBe(0);
});
it('an old-context save completion cannot mark the new context saved, and writes to the same API remain ordered', async () => {
  const completions: Array<() => void> = [];
  vi.mocked(settingsAPI.updateKey).mockImplementation((key, value) => new Promise(resolve => completions.push(() => resolve({ key, value, description: '', updated_at: '' }))));
  const view = mount(); fireEvent.click(screen.getByText('dark')); await waitFor(() => expect(completions).toHaveLength(1));
  const other = new QueryClient({ defaultOptions: { queries: { retry: false } } });
  view.rerender(<QueryClientProvider client={other}><ThemeProvider><Probe /></ThemeProvider></QueryClientProvider>);
  fireEvent.click(screen.getByText('light')); expect(settingsAPI.updateKey).toHaveBeenCalledTimes(1);
  await act(async () => completions[0]()); await waitFor(() => expect(completions).toHaveLength(2));
  expect(screen.getByTestId('theme')).toHaveTextContent('light/light/pending');
  await act(async () => completions[1]()); expect(screen.getByTestId('theme')).toHaveTextContent('light/light/saved');
});

it('StrictMode replay leaves one OS subscription and never turns the resolved color into saved intent', async () => {
  render(<StrictMode><QueryClientProvider client={new QueryClient({ defaultOptions: { queries: { retry: false } } })}><ThemeProvider><Probe /></ThemeProvider></QueryClientProvider></StrictMode>);
  await waitFor(() => expect(settingsAPI.listItems).toHaveBeenCalled());
  expect(listeners.size).toBe(1); os(false); expect(screen.getByTestId('theme')).toHaveTextContent('light/system');
  expect(settingsAPI.updateKey).not.toHaveBeenCalled(); expect(localStorage.length).toBe(0);
  cleanup(); expect(listeners.size).toBe(0);
});
