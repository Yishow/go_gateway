import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { fireEvent, render, screen } from '@testing-library/react';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import { CutoverPanel } from '@/features/datalink/workbench-v2/steps/step4/writeGroup/CutoverPanel';
import { studioV2WorkspaceWriteGroupsAPI } from '@/services/studioV2WorkspaceWriteGroups';
import { reviewOf, target, wide } from '../../fixtures/cutoverFixtures';

vi.mock('react-i18next', () => ({ useTranslation: () => ({ t: (key: string) => key }) }));
vi.mock('@/services/studioV2WorkspaceWriteGroups', () => ({
  studioV2WorkspaceWriteGroupsAPI: { reviewCutover: vi.fn(), applyCutover: vi.fn(), listCutovers: vi.fn(async () => []) },
}));

const api = vi.mocked(studioV2WorkspaceWriteGroupsAPI);
const BOUNDARY_LOCAL = '2099-01-01T00:00:00';

beforeEach(() => {
  vi.clearAllMocks();
  api.listCutovers.mockResolvedValue([]);
  api.reviewCutover.mockResolvedValue(reviewOf());
});

describe('The panel respects the page locks', () => {
  it('locks every control, a new review and the confirmation when the page becomes read-only after a review', async () => {
    const client = new QueryClient({ defaultOptions: { queries: { retry: false }, mutations: { retry: false } } });
    const element = (readonly: boolean) => (
      <QueryClientProvider client={client}>
        <CutoverPanel groups={[wide(), target()]} workspaceId="ws-1" workspaceRevision="wrev-1" connectorRevision="crev-1" readonly={readonly} onOpenGroup={vi.fn()} />
      </QueryClientProvider>
    );
    const view_ = render(element(false));
    fireEvent.click(await screen.findByTestId('cutover-target-tgt-1'));
    fireEvent.change(screen.getByTestId('cutover-boundary'), { target: { value: BOUNDARY_LOCAL } });
    expect(screen.getByTestId('cutover-review-button')).toBeEnabled();
    fireEvent.click(screen.getByTestId('cutover-review-button'));
    await screen.findByTestId('cutover-review');
    expect(screen.getByTestId('cutover-confirm')).toBeEnabled();

    view_.rerender(element(true));
    for (const id of ['cutover-mode-forward', 'cutover-mode-rollback', 'cutover-source', 'cutover-target-tgt-1', 'cutover-boundary', 'cutover-review-button', 'cutover-confirm']) expect(screen.getByTestId(id), id).toBeDisabled();
    fireEvent.click(screen.getByTestId('cutover-review-button'));
    fireEvent.click(screen.getByTestId('cutover-confirm'));
    expect(api.reviewCutover).toHaveBeenCalledTimes(1);
    expect(api.applyCutover).not.toHaveBeenCalled();
  });

  it.each([
    ['workspace', { workspaceRevision: 'wrev-2' }],
    ['destination connector', { connectorRevision: 'crev-2' }],
  ])('marks a review stale when the %s revision changes, and never confirms it', async (_name, change) => {
    const client = new QueryClient({ defaultOptions: { queries: { retry: false }, mutations: { retry: false } } });
    const element = (extra: { workspaceRevision?: string; connectorRevision?: string }) => (
      <QueryClientProvider client={client}>
        <CutoverPanel groups={[wide(), target()]} workspaceId="ws-1" workspaceRevision="wrev-1" connectorRevision="crev-1" readonly={false} onOpenGroup={vi.fn()} {...extra} />
      </QueryClientProvider>
    );
    const view_ = render(element({}));
    fireEvent.click(await screen.findByTestId('cutover-target-tgt-1'));
    fireEvent.change(screen.getByTestId('cutover-boundary'), { target: { value: BOUNDARY_LOCAL } });
    fireEvent.click(screen.getByTestId('cutover-review-button'));
    await screen.findByTestId('cutover-review');
    expect(screen.getByTestId('cutover-confirm')).toBeEnabled();
    view_.rerender(element(change));
    expect(await screen.findByTestId('cutover-review-stale')).toBeInTheDocument();
    expect(screen.getByTestId('cutover-confirm')).toBeDisabled();
    expect(api.applyCutover).not.toHaveBeenCalled();
  });

  it('explains that a running group is needed when none can be switched', async () => {
    const client = new QueryClient({ defaultOptions: { queries: { retry: false }, mutations: { retry: false } } });
    render(
      <QueryClientProvider client={client}>
        <CutoverPanel groups={[target()]} workspaceId="ws-1" workspaceRevision="wrev-1" connectorRevision="crev-1" readonly={false} onOpenGroup={vi.fn()} />
      </QueryClientProvider>,
    );
    expect(await screen.findByTestId('cutover-no-source')).toBeInTheDocument();
    expect(screen.queryByTestId('cutover-review-button')).not.toBeInTheDocument();
  });
});
