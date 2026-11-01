import { useState } from 'react';
import { act, cleanup, fireEvent, render, screen, waitFor } from '@testing-library/react';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { afterEach, expect, it, vi } from 'vitest';
import { ThemeProvider, useTheme } from '@/contexts/ThemeContext';
import { settingsAPI } from '@/services/datalink';
import { Input } from '@/features/datalink/workbench-v2/components/inputs';
vi.mock('@/services/datalink', () => ({ settingsAPI: { listItems: vi.fn(), updateKey: vi.fn() } }));
function DirtyForm() {
  const { setPreference, theme } = useTheme();
  const [text, setText] = useState('draft');
  return <div data-workbench-v2><Input aria-label="Interval text" value={text} onChange={e => setText(e.target.value)} />
    <output>{theme}</output>{(['light', 'dark', 'system'] as const).map(value =>
      <button type="button" key={value} onClick={() => setPreference(value)}>{value}</button>)}</div>;
}
afterEach(() => { cleanup(); vi.unstubAllGlobals(); localStorage.clear(); });
it('preserves incomplete input across explicit and live system themes without saving form data', async () => {
  let dark = false;
  const listeners = new Set<() => void>();
  localStorage.clear(); vi.clearAllMocks();
  vi.stubGlobal('matchMedia', () => ({ get matches() { return dark; },
    addEventListener: (_: string, fn: () => void) => listeners.add(fn),
    removeEventListener: (_: string, fn: () => void) => listeners.delete(fn) }));
  vi.mocked(settingsAPI.listItems).mockResolvedValue([]);
  vi.mocked(settingsAPI.updateKey).mockImplementation(async (key, value) => ({ key, value, description: '', updated_at: '' }));
  render(<QueryClientProvider client={new QueryClient()}><ThemeProvider><DirtyForm /></ThemeProvider></QueryClientProvider>);
  fireEvent.change(screen.getByRole('textbox'), { target: { value: '1.' } });
  for (const preference of ['dark', 'light', 'system']) {
    fireEvent.click(screen.getByRole('button', { name: preference }));
    await waitFor(() => expect(document.documentElement.style.colorScheme).toBe(preference === 'dark' ? 'dark' : 'light'));
    expect(screen.getByRole('textbox')).toHaveValue('1.');
  }
  await waitFor(() => expect(settingsAPI.updateKey).toHaveBeenLastCalledWith('theme', 'auto'));
  const writes = vi.mocked(settingsAPI.updateKey).mock.calls.length;
  act(() => { dark = true; listeners.forEach(fn => fn()); });
  expect(document.documentElement.style.colorScheme).toBe('dark');
  expect(screen.getByRole('textbox')).toHaveValue('1.');
  expect(settingsAPI.updateKey).toHaveBeenCalledTimes(writes);
  expect(vi.mocked(settingsAPI.updateKey).mock.calls.every(([key]) => key === 'theme')).toBe(true);
});
