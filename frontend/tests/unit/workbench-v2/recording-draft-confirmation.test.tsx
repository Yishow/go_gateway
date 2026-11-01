import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { act, fireEvent, render, screen, waitFor } from '@testing-library/react';
import { beforeEach, expect, it, vi } from 'vitest';
import { GroupSchemaPanel } from '../../../src/features/datalink/workbench-v2/steps/step4/writeGroup/GroupSchemaPanel';
import { studioV2WorkspaceWriteGroupsAPI } from '../../../src/services/studioV2WorkspaceWriteGroups';
import { savedGroup } from '../../fixtures/writeGroupState';
import { managedSchemaPreviewToken } from '../../fixtures/managedSchemaPreview';
vi.mock('react-i18next', () => ({ useTranslation: () => ({ t: (key: string) => key }) }));
vi.mock('../../../src/services/studioV2WorkspaceWriteGroups', () => ({ studioV2WorkspaceWriteGroupsAPI: { schemaPreview: vi.fn(), schemaApply: vi.fn(), schemaOperation: vi.fn() } }));
function mount() {
  const group = savedGroup(), client = new QueryClient({ defaultOptions: { queries: { retry: false }, mutations: { retry: false } } });
  const tree = (dirty: boolean) => <QueryClientProvider client={client}><GroupSchemaPanel group={group} workspaceId={group.workspace_id} workspaceRevision="workspace-1" destinationMatches dirty={dirty} readonly={false} /></QueryClientProvider>;
  const view = render(tree(false)); return { change: (dirty: boolean) => view.rerender(tree(dirty)) };
}
beforeEach(() => { vi.clearAllMocks(); sessionStorage.clear(); vi.mocked(studioV2WorkspaceWriteGroupsAPI.schemaPreview).mockResolvedValue(managedSchemaPreviewToken); });
it('editing then reverting cannot revive a previous schema confirmation', async () => {
  const view = mount(); fireEvent.click(screen.getByTestId('group-schema-preview'));
  await screen.findByTestId('group-schema-preview-result');
  view.change(true); view.change(false);
  await waitFor(() => expect(screen.queryByTestId('group-schema-preview-result')).not.toBeInTheDocument());
  expect(screen.queryByTestId('group-schema-apply')).not.toBeInTheDocument();
  expect(studioV2WorkspaceWriteGroupsAPI.schemaApply).not.toHaveBeenCalled();
});
it('a late preview from before an edit cannot authorize the reverted configuration', async () => {
  let resolve!: (value: typeof managedSchemaPreviewToken) => void;
  vi.mocked(studioV2WorkspaceWriteGroupsAPI.schemaPreview).mockReturnValue(new Promise(r => { resolve = r; }));
  const view = mount(); fireEvent.click(screen.getByTestId('group-schema-preview'));
  view.change(true); view.change(false);
  await act(async () => resolve(managedSchemaPreviewToken));
  expect(screen.queryByTestId('group-schema-preview-result')).not.toBeInTheDocument();
  expect(studioV2WorkspaceWriteGroupsAPI.schemaApply).not.toHaveBeenCalled();
});
