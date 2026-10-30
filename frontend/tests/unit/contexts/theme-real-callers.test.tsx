import { useEffect, useState } from 'react';
import { cleanup, fireEvent, render, screen, waitFor } from '@testing-library/react';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { MemoryRouter, Route, Routes } from 'react-router-dom';
import { afterEach, expect, it, vi } from 'vitest';
import { ThemeProvider } from '@/contexts/ThemeContext';
import { ThemeToggle } from '@/components/ThemeToggle';
import DatalinkLayout from '@/layouts/DatalinkLayout';
import { settingsAPI } from '@/services/datalink';
vi.mock('@/services/datalink', () => ({ settingsAPI: { listItems: vi.fn(), updateKey: vi.fn() } }));
vi.mock('react-i18next', () => ({ useTranslation: () => ({ t: (key: string) => key, i18n: { language: 'en', changeLanguage: vi.fn() } }) }));
afterEach(() => { cleanup(); localStorage.clear(); vi.restoreAllMocks(); });
it('keeps the actual legacy layout and header caller on one theme without remounting the outlet or saving business settings', async () => {
  localStorage.clear();
  vi.mocked(settingsAPI.listItems).mockResolvedValue([{ key: 'theme', value: 'light', description: '', updated_at: '' }]);
  vi.mocked(settingsAPI.updateKey).mockImplementation(async (key, value) => ({ key, value, description: '', updated_at: '' }));
  const mounts = vi.fn();
  function EditorDraft() {
    const [value, setValue] = useState('unsaved');
    useEffect(() => { mounts(); }, []);
    return <input aria-label="Existing draft" value={value} onChange={event => setValue(event.target.value)} />;
  }
  const client = new QueryClient({ defaultOptions: { queries: { retry: false } } });
  const { container, unmount } = render(<QueryClientProvider client={client}><ThemeProvider><ThemeToggle />
    <MemoryRouter initialEntries={['/datalink']}><Routes><Route element={<DatalinkLayout />}><Route path="/datalink" element={<EditorDraft />} /></Route></Routes></MemoryRouter>
  </ThemeProvider></QueryClientProvider>);
  await waitFor(() => expect(document.documentElement).toHaveClass('light'));
  expect(settingsAPI.updateKey).not.toHaveBeenCalled();
  const draft = screen.getByLabelText('Existing draft');
  fireEvent.change(draft, { target: { value: 'keep my draft' } });
  fireEvent.click(screen.getByRole('button', { name: 'layout.toggleTheme' }));
  await waitFor(() => expect(document.documentElement).toHaveClass('dark'));
  expect(container.querySelector('.min-h-screen')).toHaveClass('bg-[#0f172a]');
  expect(screen.getByLabelText('Existing draft')).toBe(draft);
  expect(draft).toHaveValue('keep my draft'); expect(mounts).toHaveBeenCalledTimes(1);
  await waitFor(() => expect(settingsAPI.updateKey).toHaveBeenCalledWith('theme', 'dark'));
  expect(vi.mocked(settingsAPI.updateKey).mock.calls).toEqual([['theme', 'dark']]);
  unmount(); client.clear();
});
