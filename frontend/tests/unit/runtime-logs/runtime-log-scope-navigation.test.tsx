import { act, cleanup, render, screen, waitFor } from '@testing-library/react';
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';
import RuntimeLogsPage from '@/features/runtime-logs/RuntimeLogsPage';
import { runtimeLogTransport } from '@/services/runtimeLogs';
import { parseStudioPresentation } from '@/services/studioRecordingView';
import type { RuntimeLogSnapshot } from '@/types/runtimeLogs';
import i18n from '@/i18n/config';
import { MemoryRouter } from 'react-router-dom';
import { AppRoutes } from '@/App';
import { record, snapshot } from './fixtures';
vi.mock('@/services/runtimeLogs', async original => ({ ...await original<typeof import('@/services/runtimeLogs')>(),
  runtimeLogTransport: { snapshot: vi.fn(), stream: vi.fn() },
}));
const search = (group = 'G') => `?entry=recording&workspace_id=W&group_id=${group}&group_revision=R`;
beforeEach(async () => {
  vi.resetAllMocks(); await i18n.changeLanguage('en');
  vi.mocked(runtimeLogTransport.snapshot).mockResolvedValue({ ...snapshot, records: [] });
  vi.mocked(runtimeLogTransport.stream).mockImplementation(() => new Promise(() => undefined));
});
afterEach(cleanup);
describe('recording-scoped log navigation', () => {
  it('consumes the actual application route search and hash rather than only test props', async () => {
    render(<MemoryRouter initialEntries={[`/studio/logs${search()}#actual-route`]}><AppRoutes /></MemoryRouter>);
    await waitFor(() => expect(runtimeLogTransport.stream).toHaveBeenCalledOnce());
    expect(vi.mocked(runtimeLogTransport.snapshot).mock.calls[0][0]).toMatchObject({ workspaceId: 'W', groupId: 'G' });
    expect(screen.getByRole('link', { name: 'Open setup' })).toHaveAttribute('href', expect.stringContaining('group_id=G'));
    expect(screen.getByRole('link', { name: 'Open setup' }).getAttribute('href')).toMatch(/#actual-route$/);
  });
  it('starts the original session with exact record filters and a safe canonical return', async () => {
    render(<RuntimeLogsPage search={search()} hash="#same-row" />);
    await waitFor(() => expect(runtimeLogTransport.stream).toHaveBeenCalledOnce());
    expect(runtimeLogTransport.snapshot).toHaveBeenCalledWith(expect.objectContaining({ workspaceId: 'W', groupId: 'G' }), expect.any(AbortSignal));
    expect(runtimeLogTransport.stream).toHaveBeenCalledWith(expect.objectContaining({ workspaceId: 'W', groupId: 'G' }), expect.any(String), expect.any(AbortSignal), expect.any(Function));
    const target = new URL(screen.getByRole('link', { name: 'Open setup' }).getAttribute('href')!, 'http://localhost');
    expect(target.hash).toBe('#same-row');
    expect(parseStudioPresentation(target.search)).toMatchObject({ presentation: 'recording', locatorIntent: {
      kind: 'locator', locator: { workspace_id: 'W', group_id: 'G', group_revision: 'R' },
    } });
  });
  it.each([`${search()}&group_id=other`, '?entry=recording&workspace_id=W&group_id=G',
    `${search()}&return=https://outside.invalid`, '?entry=recording&workspace_id=%ZZ&group_id=G&group_revision=R'])('rejects invalid scoped entry without falling back to a global session: %s', async value => {
    render(<RuntimeLogsPage search={value} />);
    expect(screen.getByRole('alert')).toBeVisible();
    await act(async () => {});
    expect(runtimeLogTransport.snapshot).not.toHaveBeenCalled(); expect(runtimeLogTransport.stream).not.toHaveBeenCalled();
  });
  it('cancels the prior record session on location changes and ignores late replies', async () => {
    let finish!: (value: RuntimeLogSnapshot) => void;
    vi.mocked(runtimeLogTransport.snapshot).mockImplementationOnce(() => new Promise(resolve => { finish = resolve; }));
    const view = render(<RuntimeLogsPage search={search('A')} />);
    await waitFor(() => expect(runtimeLogTransport.snapshot).toHaveBeenCalledOnce());
    view.rerender(<RuntimeLogsPage search={search('B')} />);
    await waitFor(() => expect(runtimeLogTransport.stream).toHaveBeenCalledOnce());
    expect(vi.mocked(runtimeLogTransport.snapshot).mock.calls[0][1].aborted).toBe(true);
    await act(async () => finish({ ...snapshot, records: [record('99', 'Older success')] }));
    expect(document.querySelector('[data-log-cursor="instance-a:99"]')).toBeNull();
    expect(runtimeLogTransport.stream).toHaveBeenCalledOnce();
    expect(vi.mocked(runtimeLogTransport.stream).mock.calls[0][0].groupId).toBe('B');
    view.rerender(<RuntimeLogsPage search={search('A')} />);
    await waitFor(() => expect(runtimeLogTransport.stream).toHaveBeenCalledTimes(2));
    expect(vi.mocked(runtimeLogTransport.stream).mock.calls[0][2].aborted).toBe(true);
    expect(vi.mocked(runtimeLogTransport.stream).mock.calls[1][0].groupId).toBe('A');
  });
});
