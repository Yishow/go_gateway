import * as Dialog from '@radix-ui/react-dialog';
import { useState } from 'react';
import { act, fireEvent, render, screen, within } from '@testing-library/react';
import { createMemoryRouter, Outlet, RouterProvider } from 'react-router-dom';
import { expect, it, vi } from 'vitest';
import { useRecordingDraftGuard } from '../../../src/features/datalink/workbench-v2/recording-editor/RecordingDraftNavigation';
import { RecordingRouterNavigation } from '../../../src/features/datalink/workbench-v2/recording-editor/RecordingRouterNavigation';
vi.mock('react-i18next', () => ({ useTranslation: () => ({ t: (key: string) => key }) }));
function Source() {
  const [text, setText] = useState('unsent source');
  const navigation = useRecordingDraftGuard({ version: text, dirty: true, host: true });
  return <Dialog.Root open><Dialog.Portal><Dialog.Content role={navigation.confirming ? 'alertdialog' : 'dialog'}>
    {navigation.confirmation ?? <><Dialog.Title>Source</Dialog.Title><Dialog.Description>Source fields</Dialog.Description>
      <input aria-label="source input" value={text} onChange={event => setText(event.target.value)} /></>}
  </Dialog.Content></Dialog.Portal></Dialog.Root>;
}
it('browser Back uses the already open source modal rather than stacking another active dialog', async () => {
  const router = createMemoryRouter([{ element: <RecordingRouterNavigation><Outlet /></RecordingRouterNavigation>, children: [
    { path: '/source', element: <Source /> }, { path: '/home', element: <p>Home</p> },
  ] }], { initialEntries: ['/home', '/source'], initialIndex: 1 });
  render(<RouterProvider router={router} />);
  fireEvent.change(screen.getByLabelText('source input'), { target: { value: 'Invalid source text kept' } });
  await act(async () => { await router.navigate(-1); });
  await screen.findByText('navigation.keep');
  const confirmations = screen.getAllByRole('alertdialog', { hidden: true }); expect(confirmations).toHaveLength(1);
  expect(screen.queryByRole('dialog')).not.toBeInTheDocument();
  expect(document.getElementById(confirmations[0].getAttribute('aria-describedby') ?? '')).toHaveTextContent('navigation.description');
  expect(within(confirmations[0]).getByText('navigation.keep')).toHaveFocus();
  fireEvent.click(within(confirmations[0]).getByText('navigation.keep'));
  expect(screen.getByLabelText('source input')).toHaveValue('Invalid source text kept');
  expect(router.state.location.pathname).toBe('/source');
});
