import { fireEvent, render, screen, waitFor } from '@testing-library/react';
import { afterEach, expect, it, vi } from 'vitest';
import { AcquisitionPlanPreviewPanel } from '@/features/datalink/workbench-v2/acquisition-preview/AcquisitionPlanPreviewPanel';
import { studioV2DatalinkApi } from '@/services/studioV2Workspace';
import wire from '../../fixtures/acquisition-preview-run15.json';

// Exercise the real strict service, hook and panel; only HTTP transport is replaced.
// Fixture is the unmodified run15 GET233 response, not a hand-built DTO.
vi.mock('@/services/studioV2Workspace', () => ({ studioV2DatalinkApi: { get: vi.fn() } }));
vi.mock('react-i18next', () => ({ useTranslation: () => ({ t: (key: string, values?: Record<string, unknown>) => `${key}${values ? ` ${JSON.stringify(values)}` : ''}` }) }));
afterEach(() => vi.resetAllMocks());

it('renders the captured installed wire response through the strict service and refreshes its exact head', async () => {
  const value = wire.data;
  const group = value.groups[0];
  const scope = { workspace_id: value.workspace_id, device_id: value.device_id, expected_workspace_revision: value.workspace_revision };
  vi.mocked(studioV2DatalinkApi.get).mockResolvedValue({ data: wire });
  render(<AcquisitionPlanPreviewPanel scope={scope} />);
  expect(await screen.findByTestId(`acquisition-planned-${group.polling_group_id}`)).toHaveTextContent(/^1$/);
  expect(screen.getByTestId(`acquisition-points-${group.polling_group_id}`)).toHaveTextContent(/^2$/);
  expect(screen.getByText(/wireUnits.word.*"count":4/)).toBeInTheDocument();
  expect(screen.getByTestId(`acquisition-wire-attempts-${group.polling_group_id}`)).toHaveTextContent(/^1$/);
  expect(screen.getByTestId(`acquisition-logical-${group.polling_group_id}`)).toHaveTextContent('"succeeded":1,"failed":0,"canceled":0');
  expect(screen.queryByText('unavailable')).not.toBeInTheDocument();
  expect(studioV2DatalinkApi.get).toHaveBeenNthCalledWith(1, '/studio-v2/workspace/acquisition-plan-preview', expect.objectContaining({ params: scope }));
  fireEvent.click(screen.getByRole('button', { name: 'refresh' }));
  await waitFor(() => expect(studioV2DatalinkApi.get).toHaveBeenCalledTimes(2));
  expect(studioV2DatalinkApi.get).toHaveBeenNthCalledWith(2, '/studio-v2/workspace/acquisition-plan-preview', expect.objectContaining({ params: { ...scope, expected_preview_revision: value.preview_revision } }));
  expect(await screen.findByTestId(`acquisition-wire-attempts-${group.polling_group_id}`)).toHaveTextContent(/^1$/);
});
