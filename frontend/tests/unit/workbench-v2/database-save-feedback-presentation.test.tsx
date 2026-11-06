import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { act, fireEvent, render, screen, waitFor } from '@testing-library/react';
import { expect, it, vi } from 'vitest';
import { Step4SupportPanels } from '@/features/datalink/workbench-v2/steps/step4/Step4SupportPanels';
import { beginConfigSave, confirmConfigSave, configReadContext, configReadKey } from '@/hooks/datalink/recordingConfigRefresh';
import { studioV2WorkspaceKeys } from '@/hooks/datalink/keys';
import { savedGroupState } from '../../fixtures/writeGroupState';
vi.mock('react-i18next', () => ({ useTranslation: () => ({ t: (key: string) => key }) }));
it('keeps saved confirmation and safe read-only retry visible in the collapsed connection card', async () => {
  const client = new QueryClient(), context = configReadContext('ws-1');
  // An observed failing read is projected separately from the connector save state.
  client.setQueryData(configReadKey(studioV2WorkspaceKeys.databaseConfig(), context), { id: 'db-1' });
  const captured = beginConfigSave(client, 'ws-1');
  const failure = vi.spyOn(client, 'refetchQueries').mockRejectedValue(new Error('password=SECRET'));
  confirmConfigSave(client, captured);
  await act(async () => {});
  render(<QueryClientProvider client={client}><Step4SupportPanels connector={savedGroupState().db.connector}
    onUpdateConnector={vi.fn()} onKindChange={vi.fn()} destinationReady /></QueryClientProvider>);
  expect(screen.getByTestId('database-connector-save-status')).toHaveTextContent('step4.configuration_saved');
  await waitFor(() => expect(screen.getByTestId('database-connector-refresh-status')).toHaveTextContent('step4.config_refresh_failed'));
  expect(screen.getByLabelText('step4.field_name')).not.toBeVisible();
  expect(document.body.textContent).not.toContain('SECRET');
  failure.mockResolvedValue();
  fireEvent.click(screen.getByRole('button', { name: 'step4.config_refresh_retry' }));
  await waitFor(() => expect(screen.queryByTestId('database-connector-refresh-status')).not.toBeInTheDocument());
});
