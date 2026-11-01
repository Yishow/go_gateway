import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { act, cleanup, fireEvent, render, screen, waitFor } from '@testing-library/react';
import { createMemoryRouter, RouterProvider } from 'react-router-dom';
import { afterEach, expect, it, vi } from 'vitest';
import { AppRoutes } from '@/App';
import { RecordingRouterNavigation } from '@/features/datalink/workbench-v2/recording-editor/RecordingRouterNavigation';
import { useWorkbenchV2State } from '@/features/datalink/workbench-v2/state/useWorkbenchV2State';
import type { Step4Database } from '@/features/datalink/workbench-v2/steps/step4';
import type { ComponentProps } from 'react';
import { savedGroupState } from '../../fixtures/writeGroupState';
import { locatorRead } from './studio-locator-fixtures';
vi.mock('react-i18next', () => ({ useTranslation: () => ({ t: (key: string) => key }) }));
vi.mock('@/hooks/datalink/useStudioV2Workspace', () => ({ useStudioV2WorkspaceQuery: () => ({ data: locatorRead(4).workspace, isSuccess: true, isLoading: false, isError: false, refetch: vi.fn() }) }));
vi.mock('@/hooks/datalink/useModbusShareStatus', () => ({ useModbusShareStatusQuery: () => ({ data: undefined, refetch: vi.fn() }) }));
vi.mock('@/hooks/datalink/useStudioV2WorkspaceAuditHistory', () => ({ useStudioV2WorkspaceAuditHistoryQuery: () => ({ data: { entries: [] }, isError: false }) }));
vi.mock('@/hooks/datalink/useStudioV2WorkspaceActivation', () => ({ useActivateStudioV2WorkspaceMutation: () => ({ mutateAsync: vi.fn() }) }));
vi.mock('@/hooks/datalink/useStudioLocatorEntry', () => ({ useStudioLocatorEntry: () => ({ data: { ok: true, value: locatorRead(4) }, isFetching: false, isError: false, refetch: vi.fn() }) }));
vi.mock('@/pages/datalink/workbench-v2/useStudioV2AutosaveState', () => ({ useStudioV2AutosaveState: () => {
  const actions = useWorkbenchV2State(savedGroupState()); const query = { isLoading: false, isError: false, refetch: vi.fn() };
  return { state: actions.state, actions, workspaceHydrated: true, locatorHydration: { ready: true, blocked: false },
    devicesQuery: query, rulesQuery: query, mappingsQuery: query, databaseConfigQuery: query, databaseTargetsQuery: query };
} }));
vi.mock('@/features/datalink/workbench-v2/steps/step4', () => ({ Step4Database: (props: ComponentProps<typeof Step4Database>) =>
  <button onClick={() => props.onSavedLocation?.(locatorRead(4).entry.write_groups!.groups[0])}>Confirmed saved location</button> }));
afterEach(cleanup);
it('actual App → Page → Shell passes replace to the router, preserving hash without a new-draft history entry', async () => {
  const client = new QueryClient({ defaultOptions: { queries: { retry: false } } });
  const router = createMemoryRouter([{ path: '/prior', element: <div>Prior page</div> }, { path: '*', element:
    <QueryClientProvider client={client}><RecordingRouterNavigation><AppRoutes /></RecordingRouterNavigation></QueryClientProvider> }],
  { initialEntries: ['/prior', '/studio/v2?view=recording&intent=new#saved'], initialIndex: 1 });
  render(<RouterProvider router={router} />);
  fireEvent.click(await screen.findByText('Confirmed saved location'));
  await waitFor(() => expect(router.state.location.search).toContain('entry=locator'));
  expect(router.state.location.search).toContain('group_id=G'); expect(router.state.location.search).toContain('group_revision=G3');
  expect(router.state.location.hash).toBe('#saved'); expect(router.state.historyAction).toBe('REPLACE');
  await act(async () => { await router.navigate(-1); });
  expect(await screen.findByText('Prior page')).toBeInTheDocument();
  client.clear();
});
