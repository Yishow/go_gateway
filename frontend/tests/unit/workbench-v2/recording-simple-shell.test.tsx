import { useState, type ComponentProps } from 'react';
import { cleanup, fireEvent, render, screen } from '@testing-library/react';
import { createMemoryRouter, RouterProvider, useLocation, useNavigate } from 'react-router-dom';
import { afterEach, describe, expect, it, vi } from 'vitest';
import { WorkbenchV2Shell } from '@/features/datalink/workbench-v2/shell/WorkbenchV2Shell';
import { useWorkbenchV2State } from '@/features/datalink/workbench-v2/state/useWorkbenchV2State';
import { RecordingRouterNavigation } from '@/features/datalink/workbench-v2/recording-editor/RecordingRouterNavigation';
import { useRecordingDraftGuard } from '@/features/datalink/workbench-v2/recording-editor/RecordingDraftNavigation';
import type { Step4Database } from '@/features/datalink/workbench-v2/steps/step4';
const captured = vi.hoisted(() => ({ props: undefined as ComponentProps<typeof Step4Database> | undefined }));
vi.mock('react-i18next', () => ({ useTranslation: () => ({ t: (key: string) => key }) }));
vi.mock('@/features/datalink/workbench-v2/steps/step4', () => ({ Step4Database: function CanonicalStep4(props: ComponentProps<typeof Step4Database>) {
  captured.props = props;
  const [dirty, setDirty] = useState(false);
  useRecordingDraftGuard({ version: String(dirty), dirty, discard: () => setDirty(false) });
  return <div data-testid="single-step4-owner"><button onClick={() => setDirty(true)}>Edit draft</button></div>;
} }));
afterEach(cleanup);
function RoutedShell() {
  const location = useLocation(), navigate = useNavigate(), actions = useWorkbenchV2State();
  if (new URLSearchParams(location.search).get('view') === 'recordings') return <div data-testid="returned-home">Home</div>;
  return <WorkbenchV2Shell state={actions.state} actions={actions} workspaceId="W" recordingIntent="new"
    presentationLocation={{ search: location.search, hash: location.hash }} navigateTo={navigate} />;
}
function mount() {
  const router = createMemoryRouter([{ path: '/studio/v2', element: <RecordingRouterNavigation><RoutedShell /></RecordingRouterNavigation> }],
    { initialEntries: ['/studio/v2?view=recording&intent=new#same-position'] });
  render(<RouterProvider router={router} />); return router;
}
describe('simple route is one screen, not a second wizard', () => {
  it('mounts Step4 directly once, without Basic, sidebar or step navigation', () => {
    mount(); expect(screen.getAllByTestId('single-step4-owner')).toHaveLength(1);
    expect(captured.props).toMatchObject({ presentation: 'simple', initialSelection: 'new', workspaceId: 'W' });
    expect(screen.queryByTestId('sidebar-rail-aside')).not.toBeInTheDocument();
    expect(screen.queryByText('STEP 01 / 04')).not.toBeInTheDocument();
    expect(captured.props?.activateWorkspace).toBeUndefined();
  });
  it('returns through the existing route dirty guard and preserves hash', async () => {
    const router = mount(); fireEvent.click(screen.getByText('Edit draft')); fireEvent.click(screen.getByText('simple.back'));
    expect(await screen.findByRole('alertdialog')).toBeInTheDocument();
    fireEvent.click(screen.getByText('navigation.keep'));
    expect(screen.getByTestId('single-step4-owner')).toBeInTheDocument();
    fireEvent.click(screen.getByText('simple.back')); fireEvent.click(await screen.findByText('navigation.discard'));
    expect(await screen.findByTestId('returned-home')).toBeInTheDocument();
    expect(router.state.location.search).toBe('?view=recordings'); expect(router.state.location.hash).toBe('#same-position');
  });
});
