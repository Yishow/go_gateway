import { useState } from 'react';
import { act, fireEvent, render, screen, waitFor } from '@testing-library/react';
import { createMemoryRouter, Outlet, RouterProvider, useNavigate } from 'react-router-dom';
import { expect, it, vi } from 'vitest';
import { useRecordingDraftGuard } from '../../../src/features/datalink/workbench-v2/recording-editor/RecordingDraftNavigation';
import { RecordingRouterNavigation } from '../../../src/features/datalink/workbench-v2/recording-editor/RecordingRouterNavigation';
vi.mock('react-i18next', () => ({ useTranslation: () => ({ t: (key: string) => key }) }));
function Editor() {
  const [value, setValue] = useState('draft'); const [scope, setScope] = useState('device-a'); const navigate = useNavigate();
  const guard = useRecordingDraftGuard({ version: scope + value, dirty: true });
  return <><input aria-label="intent" value={value} onChange={event => setValue(event.target.value)} /><p>{scope}</p>
    <button onClick={() => guard.request(() => setScope('device-b'))}>Switch scope</button>
    <button onClick={() => navigate('/other')}>Home</button></>;
}
function mount(entries: string[], index: number) {
  const router = createMemoryRouter([{ element: <RecordingRouterNavigation><Outlet /></RecordingRouterNavigation>, children: [
    { path: '/editor', element: <Editor /> }, { path: '/other', element: <p>Other page</p> },
  ] }], { initialEntries: entries, initialIndex: index });
  render(<RouterProvider router={router} />); return router;
}
it.each([['Back', ['/other', '/editor'], 1, -1], ['Forward', ['/editor', '/other'], 0, 1]] as const)('%s retains all intent until discard and leaves history traversal intact', async (_label, entries, index, delta) => {
  const router = mount([...entries], index);
  fireEvent.change(screen.getByLabelText('intent'), { target: { value: '1e invalid but intentional' } });
  await act(async () => { await router.navigate(delta); });
  expect(await screen.findByRole('alertdialog')).toBeInTheDocument();
  fireEvent.click(screen.getByRole('button', { name: 'navigation.keep' }));
  expect(screen.getByLabelText('intent')).toHaveValue('1e invalid but intentional');
  expect(router.state.location.pathname).toBe('/editor');
  await act(async () => { await router.navigate(delta); });
  fireEvent.click(await screen.findByRole('button', { name: 'navigation.discard' }));
  await waitFor(() => expect(router.state.location.pathname).toBe('/other'));
});
it('return home and internal scope changes share retain/discard semantics', async () => {
  const router = mount(['/editor'], 0);
  fireEvent.click(screen.getByText('Switch scope')); fireEvent.click(screen.getByText('navigation.keep'));
  expect(screen.getByText('device-a')).toBeInTheDocument();
  fireEvent.click(screen.getByText('Switch scope')); fireEvent.click(screen.getByText('navigation.discard'));
  expect(screen.getByText('device-b')).toBeInTheDocument();
  fireEvent.click(screen.getByText('Home')); expect(await screen.findByRole('alertdialog')).toBeInTheDocument();
  fireEvent.click(screen.getByText('navigation.keep')); expect(router.state.location.pathname).toBe('/editor');
});
it('a changed draft invalidates an outstanding leave confirmation instead of authorizing the newer intent', async () => {
  const router = mount(['/editor'], 0); fireEvent.click(screen.getByText('Home'));
  await screen.findByRole('alertdialog');
  fireEvent.change(screen.getByLabelText('intent'), { target: { value: 'Newer intent' } });
  await waitFor(() => expect(screen.queryByRole('alertdialog')).not.toBeInTheDocument());
  expect(router.state.location.pathname).toBe('/editor'); expect(screen.getByLabelText('intent')).toHaveValue('Newer intent');
});
