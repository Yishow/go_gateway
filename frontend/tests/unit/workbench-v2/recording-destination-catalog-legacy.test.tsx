import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { act, render, waitFor } from '@testing-library/react';
import { expect, it, vi } from 'vitest';
import { RecordingDestinationFields } from '@/features/datalink/workbench-v2/steps/step4/writeGroup/RecordingDestinationFields';
import { newDraft } from '@/features/datalink/workbench-v2/state/writeGroup/draft';
import { readRecordingDestinationCatalog } from '@/services/studioV2DatabaseCatalog';
import { savedGroupState } from '../../fixtures/writeGroupState';

vi.mock('react-i18next', () => ({ useTranslation: () => ({ t: (key: string) => key }) }));
vi.mock('@/services/studioV2DatabaseCatalog', () => ({ readRecordingDestinationCatalog: vi.fn(async () => ({ default_schema: 'main', schemas: [] })) }));

it('reads a saved legacy workspace with an explicitly empty revision', async () => {
  const connector = savedGroupState().db.connector;
  const draft = newDraft({ connector_id: 'conn-1', connector_revision: 'crev-1', table_schema: '', table_name: '' });
  const client = new QueryClient({ defaultOptions: { queries: { retry: false } } });
  render(<QueryClientProvider client={client}><RecordingDestinationFields workspaceId="ws-1" workspaceRevision="" connector={connector}
    draft={draft} readonly={false} locked={false} onChange={() => undefined} onStorageStrategy={() => undefined} /></QueryClientProvider>);
  await waitFor(() => expect(readRecordingDestinationCatalog).toHaveBeenCalledWith(expect.objectContaining({ workspaceRevision: '' }), undefined));
});

it('does not request a catalog for an immutable managed destination', async () => {
  vi.mocked(readRecordingDestinationCatalog).mockClear();
  const connector = savedGroupState().db.connector;
  const draft = { ...newDraft({ connector_id: 'conn-1', connector_revision: 'crev-1', table_schema: '', table_name: 'gw_table' }), storage_strategy: 'managed' as const };
  const client = new QueryClient({ defaultOptions: { queries: { retry: false } } });
  render(<QueryClientProvider client={client}><RecordingDestinationFields workspaceId="ws-1" workspaceRevision="W" connector={connector}
    scopedPresentation draft={draft} readonly={false} locked onChange={() => undefined} onStorageStrategy={() => undefined} /></QueryClientProvider>);
  await act(async () => {});
  expect(readRecordingDestinationCatalog).not.toHaveBeenCalled();
});
