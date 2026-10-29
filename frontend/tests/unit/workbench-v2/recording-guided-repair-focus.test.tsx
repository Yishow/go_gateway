import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { cleanup, fireEvent, render, screen } from '@testing-library/react';
import { afterEach, expect, it, vi } from 'vitest';
import { CircuitProfileBasicPanel } from '@/features/datalink/workbench-v2/steps/step4/CircuitProfileBasicPanel';
import { studioV2WorkspaceWriteGroupsAPI as groups } from '@/services/studioV2WorkspaceWriteGroups';
import { studioV2WorkspaceMeasurementsAPI as measurements } from '@/services/studioV2WorkspaceMeasurements';
import { studioV2RecordingStartAPI as start } from '@/services/studioV2RecordingStart';
import { savedGroupState } from '../../fixtures/writeGroupState';
import { expandSources } from '../../fixtures/sourceReview';
vi.mock('react-i18next', () => ({ useTranslation: () => ({ t: (key: string) => key }) }));
vi.mock('@/services/studioV2WorkspaceWriteGroups', () => ({ studioV2WorkspaceWriteGroupsAPI: {
  list: vi.fn(async () => ({ workspace_id: 'ws-1', workspace_revision: 'head-1', groups: [] })),
  ensureCircuitProfile: vi.fn(), update: vi.fn(), schemaPreview: vi.fn(), schemaApply: vi.fn(),
} }));
vi.mock('@/services/studioV2WorkspaceMeasurements', () => ({ studioV2WorkspaceMeasurementsAPI: {
  list: vi.fn(async () => { throw new Error('measurement lookup unavailable'); }), create: vi.fn(), update: vi.fn(),
} }));
vi.mock('@/services/studioV2RecordingStart', () => ({ studioV2RecordingStartAPI: { start: vi.fn(), operation: vi.fn() } }));
afterEach(cleanup);
it('focuses the actual failed measurement lookup repair without issuing a GET or any owner command', async () => {
  const client = new QueryClient({ defaultOptions: { queries: { retry: false } } });
  render(<QueryClientProvider client={client}><CircuitProfileBasicPanel state={savedGroupState()} workspaceId="ws-1" readonly={false} /></QueryClientProvider>);
  expandSources(); fireEvent.change(screen.getByTestId('circuit-profile-role-tag-0'), { target: { value: 'counter_raw' } });
  fireEvent.change(screen.getByTestId('circuit-profile-circuit-name'), { target: { value: 'Counter' } });
  const failure = await screen.findByTestId('circuit-profile-measurements-error');
  const originalQueries = vi.mocked(measurements.list).mock.calls.length;
  const primary = screen.getByTestId('recording-next-action'); expect(primary).toBeEnabled();
  fireEvent.click(primary);
  expect(failure.querySelector('button')).toHaveFocus();
  expect(measurements.list).toHaveBeenCalledTimes(originalQueries);
  expect(groups.ensureCircuitProfile).not.toHaveBeenCalled(); expect(groups.update).not.toHaveBeenCalled();
  expect(groups.schemaPreview).not.toHaveBeenCalled(); expect(groups.schemaApply).not.toHaveBeenCalled(); expect(start.start).not.toHaveBeenCalled();
});
