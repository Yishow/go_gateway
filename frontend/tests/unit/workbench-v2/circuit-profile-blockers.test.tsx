import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { fireEvent, render, screen, waitFor, within } from '@testing-library/react';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import { CircuitProfileSchemaPreparation } from '@/features/datalink/workbench-v2/steps/step4/CircuitProfileSchemaPreparation';
import { studioV2WorkspaceWriteGroupsAPI } from '@/services/studioV2WorkspaceWriteGroups';
import type { WriteGroup } from '@/types/studioV2WriteGroup';
import { longProfileGroup } from '../../fixtures/circuitBatchPanel';

vi.mock('react-i18next', () => ({
  useTranslation: () => ({ t: (key: string, options?: unknown) => (options && typeof options === 'object' && 'name' in options ? `${key}:${(options as { name: string }).name}` : key) }),
}));
vi.mock('@/services/studioV2WorkspaceWriteGroups', () => ({
  studioV2WorkspaceWriteGroupsAPI: { schemaPreview: vi.fn(), schemaApply: vi.fn(), schemaOperation: vi.fn() },
}));

const api = vi.mocked(studioV2WorkspaceWriteGroupsAPI);
const typed = (code: string) => Object.assign(new Error('Request failed with status code 422 postgres://user:secret@db/prod'), { response: { status: 422, data: { error: { code } } } });
const previewOk = { operation_id: 'op-ok', token: 'token-ok', statements: ['CREATE TABLE secret_table (x int)'], tables: [], table_prefix: 'gw', group_layout: undefined } as never;

const inst = (over: Partial<WriteGroup> = {}): WriteGroup => ({ ...longProfileGroup({ id: 'grp-i', circuitId: 'cir-a', circuitName: 'Feeder A', interval: 15, points: [0, 1] }), name: 'Feeder A (Instantaneous)', ...over });
const cum = (over: Partial<WriteGroup> = {}): WriteGroup => ({ ...longProfileGroup({ id: 'grp-c', circuitId: 'cir-a', circuitName: 'Feeder A', interval: 10, points: [0, 1], kind: 'cumulative' }), name: 'Feeder A (Cumulative)', ...over });
const moved = (group: WriteGroup, revision: string): WriteGroup => ({ ...group, destination: { ...group.destination, connector_revision: revision } });

function mount(groups: WriteGroup[], extra: Record<string, unknown> = {}) {
  const client = new QueryClient({ defaultOptions: { queries: { retry: false }, mutations: { retry: false } } });
  const element = (list: WriteGroup[]) => (
    <QueryClientProvider client={client}>
      <CircuitProfileSchemaPreparation
        workspaceId="ws-1" workspaceRevision="wrev-1" preparedGroups={list} disabled={false} readonly={false} isPreparing={false}
        destinationMatches onPrepare={vi.fn()} onApplied={vi.fn()} {...extra}
      />
    </QueryClientProvider>
  );
  const view = render(element(groups));
  return { ...view, again: (list: WriteGroup[]) => view.rerender(element(list)) };
}
const blockers = () => screen.queryAllByTestId(/^circuit-blocker-(?!repair)[a-z_]+-\d+$/);

beforeEach(() => {
  vi.clearAllMocks();
  window.sessionStorage.clear();
  api.schemaOperation.mockRejectedValue(new Error('not used'));
});

describe('One cause of the same destination is explained once, with every affected profile', () => {
  it('shows a single Timescale explanation naming both profiles and no repeated per group error', async () => {
    api.schemaPreview.mockRejectedValue(typed('timescale_unavailable'));
    const onRepairDestination = vi.fn();
    mount([inst(), cum()], { onRepairDestination });
    const card = await screen.findByTestId('circuit-blocker-timescale_unavailable-0');
    expect(blockers()).toHaveLength(1);
    expect(card).toHaveTextContent('step4.basic.blockers.timescale_unavailable.title');
    expect(card).toHaveTextContent('Feeder A (Instantaneous)');
    expect(card).toHaveTextContent('Feeder A (Cumulative)');
    expect(screen.queryAllByTestId('group-schema-error')).toHaveLength(0);
    fireEvent.click(within(card).getByTestId('circuit-blocker-repair-destination'));
    expect(onRepairDestination).toHaveBeenCalledTimes(1);
    expect(document.body.textContent).not.toMatch(/secret|postgres:/);
  });

  it('keeps the same cause separate when the groups are frozen to different destination revisions', async () => {
    api.schemaPreview.mockRejectedValue(typed('timescale_unavailable'));
    mount([inst(), moved(cum(), 'crev-other')]);
    await waitFor(() => expect(blockers()).toHaveLength(2));
    for (const card of blockers()) expect(card.textContent?.match(/Feeder A \((Instantaneous|Cumulative)\)/g)).toHaveLength(1);
  });

  it('keeps different causes of one destination separate', async () => {
    api.schemaPreview.mockImplementation(async (id) => { throw typed(id === 'grp-i' ? 'timescale_unavailable' : 'storage_owner_mismatch'); });
    mount([inst(), cum()]);
    expect(await screen.findByTestId(/^circuit-blocker-timescale_unavailable-/)).toBeInTheDocument();
    expect(await screen.findByTestId(/^circuit-blocker-storage_owner_mismatch-/)).toBeInTheDocument();
    expect(blockers()).toHaveLength(2);
  });

  it('opens the exact group for a repair that belongs to a group', async () => {
    api.schemaPreview.mockRejectedValue(typed('schema_repair_blocked'));
    const onOpenAdvanced = vi.fn();
    mount([inst(), cum()], { onOpenAdvanced });
    const card = await screen.findByTestId('circuit-blocker-schema_repair_blocked-0');
    fireEvent.click(within(card).getByTestId('circuit-blocker-repair-group-grp-c'));
    expect(onOpenAdvanced).toHaveBeenCalledWith('grp-c');
  });
});

describe('The explanation never confirms or hides anything for another group', () => {
  it('lets one profile confirm its own schema while the other stays blocked, and sends only that confirmation', async () => {
    api.schemaPreview.mockImplementation(async (id) => { if (id === 'grp-i') return previewOk; throw typed('timescale_unavailable'); });
    api.schemaApply.mockResolvedValue({ operation_id: 'op-ok', status: 'succeeded' } as never);
    mount([inst(), cum()]);
    await screen.findByTestId('circuit-blocker-timescale_unavailable-0');
    expect(screen.getByTestId('circuit-blocker-timescale_unavailable-0')).not.toHaveTextContent('Feeder A (Instantaneous)');
    const applies = screen.getAllByTestId('group-schema-apply');
    expect(applies).toHaveLength(1);
    fireEvent.click(applies[0]);
    await waitFor(() => expect(api.schemaApply).toHaveBeenCalledTimes(1));
    expect(api.schemaApply.mock.calls[0][0]).toBe('grp-i');
  });

  it('leaves an error without a recognised destination cause visible on its own group', async () => {
    api.schemaPreview.mockImplementation(async (id) => { throw id === 'grp-i' ? typed('timescale_unavailable') : new Error('boom'); });
    mount([inst(), cum()]);
    await screen.findByTestId('circuit-blocker-timescale_unavailable-0');
    await waitFor(() => expect(screen.getAllByTestId('group-schema-error')).toHaveLength(1));
    expect(blockers()).toHaveLength(1);
  });

  it('drops the explanation once the group scope changes and its new preview succeeds', async () => {
    api.schemaPreview.mockRejectedValueOnce(typed('timescale_unavailable'));
    const view = mount([inst()]);
    await screen.findByTestId('circuit-blocker-timescale_unavailable-0');
    api.schemaPreview.mockResolvedValue(previewOk);
    view.again([{ ...inst(), revision: 'rev-grp-i-2' }]);
    await waitFor(() => expect(blockers()).toHaveLength(0));
    expect(await screen.findByTestId('group-schema-apply')).toBeInTheDocument();
  });
});
