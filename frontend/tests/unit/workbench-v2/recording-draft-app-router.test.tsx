import { StrictMode } from 'react';
import { act, fireEvent, render, screen, waitFor } from '@testing-library/react';
import { createMemoryRouter, Link, RouterProvider, useLocation, useResolvedPath } from 'react-router-dom';
import { expect, it, vi } from 'vitest';
import { AppRouteBoundary } from '../../../src/App';
vi.mock('react-i18next', () => ({ useTranslation: () => ({ t: (key: string) => key }) }));
vi.mock('../../../src/pages/datalink/workbench-v2/DatalinkWorkbenchV2Page', () => ({ default: function Page(props: { presentation?: string; locatorIntent?: { kind: string }; runtimeReturnFocus?: { step: number } }) {
  const location = useLocation(), parent = useResolvedPath('../runtime'), self = useResolvedPath('.');
  return <><p data-testid="location">{location.pathname}{location.search}{location.hash}</p><p data-testid="presentation">{props.presentation}</p>
    <p data-testid="locator">{props.locatorIntent?.kind}</p><p data-testid="step">{props.runtimeReturnFocus?.step}</p>
    <output data-testid="relative">{parent.pathname}|{self.pathname}</output><Link to="?query=next#focus">Relative query</Link></>;
} }));
vi.mock('../../../src/features/datalink/runtime-dashboard/RuntimeDashboardRoute', () => ({ RuntimeDashboardRoute: () => <p>Runtime</p> }));
vi.mock('../../../src/router/gateway', () => ({ GatewayEntryRoute: () => <p>Gateway</p>, GatewayCreateEntryRedirect: () => <p>Create</p>, GatewayQuickSetupRoute: () => <p>Setup</p>, GatewayExpertWorkbenchRoute: () => <p>Expert</p> }));
function mount(entry: string) {
  const router = createMemoryRouter([{ path: '*', element: <AppRouteBoundary /> }], { initialEntries: [entry] });
  render(<StrictMode><RouterProvider router={router} /></StrictMode>); return router;
}
it.each(['/', '/does-not-exist', '/datalink/workbench'])('keeps %s entry redirect with query/hash and no dirty prompt', async entry => {
  const router = mount(`${entry}?view=recordings#saved`);
  await waitFor(() => expect(router.state.location.pathname).toBe('/studio/v2'));
  expect(await screen.findByTestId('location')).toHaveTextContent('/studio/v2?view=recordings#saved');
  expect(screen.getByTestId('presentation')).toHaveTextContent('recordings'); expect(screen.queryByRole('alertdialog')).not.toBeInTheDocument();
});
it('preserves deep-link locator and ordinary relative query navigation under the catchall', async () => {
  const router = mount('/studio/v2?entry=locator&workspace_id=W&step=4&group_id=G&group_revision=R#member');
  expect(await screen.findByTestId('locator')).toHaveTextContent('locator');
  expect(screen.getByTestId('relative')).toHaveTextContent('/runtime|/studio/v2');
  fireEvent.click(screen.getByText('Relative query'));
  await waitFor(() => expect(router.state.location.search).toBe('?query=next'));
  expect(router.state.location.pathname).toBe('/studio/v2'); expect(router.state.location.hash).toBe('#focus');
  await act(async () => { await router.navigate('/studio/runtime'); }); expect(await screen.findByText('Runtime')).toBeInTheDocument();
});
it('preserves readiness return and legacy feature route navigation without recording guards', async () => {
  const router = mount('/studio/v2?focus=readiness&step=2&issue=source-required');
  expect(await screen.findByTestId('step')).toHaveTextContent('2');
  await act(async () => { await router.navigate('/gateway/entry'); }); expect(await screen.findByText('Gateway')).toBeInTheDocument();
});
