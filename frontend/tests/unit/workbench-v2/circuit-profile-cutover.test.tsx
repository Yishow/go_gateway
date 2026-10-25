import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { act, fireEvent, render, screen, waitFor, within } from '@testing-library/react';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import { Step4Database } from '@/features/datalink/workbench-v2/steps/step4/Step4Database';
import { studioV2RecordingStartAPI } from '@/services/studioV2RecordingStart';
import { studioV2WorkspaceWriteGroupsAPI } from '@/services/studioV2WorkspaceWriteGroups';
import type { WriteGroup } from '@/types/studioV2WriteGroup';
import type { WriteGroupCutover } from '@/types/studioV2WriteGroupCutover';
import { httpFailure } from '../../fixtures/circuitBatchPanel';
import { DIGEST, member, quietDelivery, recordOf, REVIEWED_BOUNDARY, reviewOf, target, view, wide, wideMeanTarget } from '../../fixtures/cutoverFixtures';
import { savedGroupState, TABLE_COLUMNS } from '../../fixtures/writeGroupState';

// Real hooks, the real Step 4 host and the real Advanced group section; only the service boundary is simulated. This is
// frontend evidence about what is sent and shown, not SQL: the cutover semantics are proven by the backend cutover tests.
vi.mock('react-i18next', () => ({ useTranslation: () => ({ t: (key: string, options?: Record<string, unknown>) => (options && typeof options === 'object' && 'name' in options ? `${key}:${String(options.name)}` : key) }) }));
vi.mock('@/features/datalink/workbench-v2/steps/step4/useStep4TargetColumns', () => ({
  useStep4TargetColumns: () => ({ columns: TABLE_COLUMNS, status: 'exists', refetch: vi.fn() }),
}));
vi.mock('@/services/studioV2WorkspaceWriteGroups', () => ({
  studioV2WorkspaceWriteGroupsAPI: {
    list: vi.fn(), readiness: vi.fn(), delivery: vi.fn(), create: vi.fn(), update: vi.fn(), apply: vi.fn(), disable: vi.fn(), remove: vi.fn(),
    get: vi.fn(), schemaPreview: vi.fn(), schemaApply: vi.fn(), schemaOperation: vi.fn(), ensureCircuitProfile: vi.fn(),
    reviewCutover: vi.fn(), applyCutover: vi.fn(), getCutover: vi.fn(), listCutovers: vi.fn(),
  },
}));
vi.mock('@/services/studioV2WorkspaceMeasurements', () => ({ studioV2WorkspaceMeasurementsAPI: { list: vi.fn(async () => []) } }));
vi.mock('@/services/studioV2RecordingStart', () => ({ studioV2RecordingStartAPI: { start: vi.fn(), operation: vi.fn(), ensureBasic: vi.fn() } }));

const api = vi.mocked(studioV2WorkspaceWriteGroupsAPI);
const startAPI = vi.mocked(studioV2RecordingStartAPI);
const BOUNDARY_LOCAL = '2099-01-01T00:00:00';
/** Everything that changes a group, a writer or the schema. A cutover must reach the server only through its own API. */
const OTHER_WRITES = [api.create, api.update, api.apply, api.disable, api.remove, api.ensureCircuitProfile, api.schemaApply, startAPI.start, startAPI.ensureBasic];
const expectNoOtherWrites = () => OTHER_WRITES.forEach((write) => expect(write).not.toHaveBeenCalled());

let groups: WriteGroup[];
let recorded: WriteGroupCutover[];
beforeEach(() => {
  vi.clearAllMocks();
  window.sessionStorage.clear();
  groups = [wide(), target()];
  recorded = [];
  api.list.mockImplementation(async () => ({ workspace_id: 'ws-1', workspace_revision: 'wrev-1', groups: structuredClone(groups) }));
  api.delivery.mockResolvedValue(quietDelivery as never);
  api.readiness.mockResolvedValue({ workspace_id: 'ws-1', group_id: 'x', group_revision: 'r', applied_revision: '', config_ready: true, schema_ready: true, ready: true, issues: [] } as never);
  api.listCutovers.mockImplementation(async () => structuredClone(recorded));
  api.reviewCutover.mockResolvedValue(reviewOf());
  api.applyCutover.mockResolvedValue(recordOf());
});

function renderHost() {
  const client = new QueryClient({ defaultOptions: { queries: { retry: false }, mutations: { retry: false } } });
  const view_ = render(<QueryClientProvider client={client}><Step4Database state={savedGroupState()} dispatch={vi.fn()} workspaceId="ws-1" /></QueryClientProvider>);
  return { ...view_, client };
}
/** Opens the panel the way an operator reaches it: the Advanced group section of Step 4. */
async function chooseForward(options: { review?: boolean } = {}) {
  const host = renderHost();
  await screen.findByTestId('cutover-panel');
  fireEvent.click(await screen.findByTestId('cutover-target-tgt-1'));
  fireEvent.change(screen.getByTestId('cutover-boundary'), { target: { value: BOUNDARY_LOCAL } });
  if (options.review !== false) {
    fireEvent.click(screen.getByTestId('cutover-review-button'));
    await screen.findByTestId('cutover-review');
  }
  return host;
}

describe('A reviewed cutover is reachable from the Advanced group section and sends nothing before confirmation', () => {
  it('reviews the exact selection, shows the difference, and applies only the reviewed boundary text and digest', async () => {
    await chooseForward();
    expect(api.applyCutover).not.toHaveBeenCalled();
    expect(api.reviewCutover).toHaveBeenCalledTimes(1);
    const [sourceId, reviewRequest] = api.reviewCutover.mock.calls[0];
    expect(sourceId).toBe('wide-1');
    expect(reviewRequest).toMatchObject({
      kind: 'forward', workspace_id: 'ws-1', expected_workspace_revision: 'wrev-1', expected_source_revision: 'rev-wide', expected_connector_revision: 'crev-1',
      targets: [{ group_id: 'tgt-1', expected_revision: 'rev-tgt-1-1' }], boundary: new Date(BOUNDARY_LOCAL).toISOString(),
    });
    expect(reviewRequest.review_digest).toBeUndefined();

    const review = within(screen.getByTestId('cutover-review'));
    expect(screen.getByTestId('cutover-review-boundary')).toHaveTextContent(REVIEWED_BOUNDARY); // microseconds, never rounded to milliseconds
    expect(review.getByTestId('cutover-review-source-wide-1')).toHaveTextContent('PLC A raw');
    expect(review.getByTestId('cutover-review-target-tgt-1')).toHaveTextContent('step4.cutover.layout_long');
    expect(review.getByTestId('cutover-moved')).toHaveTextContent('step4.cutover.moved');
    expect(review.getByText('step4.cutover.explain_truncation')).toBeInTheDocument();
    expect(review.getByText('step4.cutover.explain_drain')).toBeInTheDocument();
    expect(review.getByText('step4.cutover.explain_baseline')).toBeInTheDocument();
    expect(api.applyCutover).not.toHaveBeenCalled();

    const listCallsBeforeApply = api.list.mock.calls.length;
    fireEvent.click(screen.getByTestId('cutover-confirm'));
    await screen.findByTestId('cutover-record');
    await waitFor(() => expect(api.list.mock.calls.length).toBeGreaterThan(listCallsBeforeApply)); // canonical groups are re-read, not patched locally
    expect(api.applyCutover).toHaveBeenCalledTimes(1);
    const [appliedSource, applyRequest] = api.applyCutover.mock.calls[0];
    expect(appliedSource).toBe('wide-1');
    expect(applyRequest).toMatchObject({ ...reviewRequest, boundary: REVIEWED_BOUNDARY, review_digest: DIGEST });
    expect(screen.getByTestId('cutover-record-id')).toHaveTextContent('cut-1');
    expect(screen.getByTestId('cutover-record-boundary')).toHaveTextContent(REVIEWED_BOUNDARY);
    expect(screen.getByTestId('cutover-record-started')).toHaveTextContent('Feeder (instantaneous)');
    expect(screen.getByTestId('cutover-record-ended')).toHaveTextContent('PLC A raw');
    expect(screen.getByTestId('cutover-scheduled')).toHaveTextContent('step4.cutover.scheduled'); // a future boundary is scheduled, not already collecting
    expectNoOtherWrites();
  });

  it.each([
    ['the target selection', async () => { fireEvent.click(screen.getByTestId('cutover-target-tgt-1')); }],
    ['the boundary', async () => { fireEvent.change(screen.getByTestId('cutover-boundary'), { target: { value: '2099-01-02T00:00:00' } }); }],
  ])('marks the review stale and blocks confirmation when %s changes', async (_name, change) => {
    await chooseForward();
    expect(screen.getByTestId('cutover-confirm')).toBeEnabled();
    await act(async () => { await change(); });
    expect(await screen.findByTestId('cutover-review-stale')).toBeInTheDocument();
    expect(screen.getByTestId('cutover-confirm')).toBeDisabled();
    expect(api.applyCutover).not.toHaveBeenCalled();
  });

  it('marks the review stale when a saved group changes after it, and needs a fresh review to confirm again', async () => {
    const { client } = await chooseForward();
    groups = [wide({ revision: 'rev-wide-2' }), target()];
    await act(async () => { await client.invalidateQueries(); });
    expect(await screen.findByTestId('cutover-review-stale')).toBeInTheDocument();
    expect(screen.getByTestId('cutover-confirm')).toBeDisabled();
    fireEvent.click(screen.getByTestId('cutover-review-button'));
    await waitFor(() => expect(api.reviewCutover).toHaveBeenCalledTimes(2));
    expect(api.reviewCutover.mock.calls[1][1].expected_source_revision).toBe('rev-wide-2');
    await waitFor(() => expect(screen.getByTestId('cutover-confirm')).toBeEnabled());
    expect(api.applyCutover).not.toHaveBeenCalled();
  });

  it('never takes a fresh revision with an old digest: a changed review is applied only with its own digest', async () => {
    await chooseForward();
    fireEvent.change(screen.getByTestId('cutover-boundary'), { target: { value: '2099-01-02T00:00:00' } });
    api.reviewCutover.mockResolvedValue(reviewOf({ boundary: '2099-01-02T00:00:00.000001Z', digest: 'b'.repeat(64) }));
    fireEvent.click(screen.getByTestId('cutover-review-button'));
    await waitFor(() => expect(screen.getByTestId('cutover-review-boundary')).toHaveTextContent('2099-01-02T00:00:00.000001Z'));
    fireEvent.click(await screen.findByTestId('cutover-confirm'));
    await waitFor(() => expect(api.applyCutover).toHaveBeenCalledTimes(1));
    expect(api.applyCutover.mock.calls[0][1]).toMatchObject({ boundary: '2099-01-02T00:00:00.000001Z', review_digest: 'b'.repeat(64) });
  });
});

describe('A cutover that cannot be applied says why and opens the exact group to repair', () => {
  it('shows a target that is not ready, opens it for repair and offers no confirmation', async () => {
    api.reviewCutover.mockResolvedValue(reviewOf({ applicable: false, issues: ['target_not_ready'], targets: [view(target(), { schema_confirmed: false })], unselected: [], remaining: [] }));
    await chooseForward();
    expect(screen.getByTestId('cutover-issue-target_not_ready')).toHaveTextContent('step4.cutover.issues.target_not_ready');
    expect(screen.getByTestId('cutover-review-target-tgt-1')).toHaveTextContent('step4.cutover.schema_pending');
    expect(screen.getByTestId('cutover-confirm')).toBeDisabled();
    fireEvent.click(screen.getByTestId('cutover-repair-target-tgt-1'));
    expect(await screen.findByTestId('group-name')).toHaveValue('Feeder (instantaneous)');
    expect(api.applyCutover).not.toHaveBeenCalled();
    expectNoOtherWrites();
  });

  it('explains sources that stay, opens the current group and does not remove or change any member', async () => {
    api.reviewCutover.mockResolvedValue(reviewOf({ applicable: false, issues: ['unselected_sources', 'remaining_draft_mismatch'] }));
    await chooseForward();
    expect(screen.getByTestId('cutover-unselected')).toHaveTextContent('step4.cutover.unselected');
    expect(screen.getByTestId('cutover-issue-remaining_draft_mismatch')).toBeInTheDocument();
    fireEvent.click(screen.getAllByTestId('cutover-repair-source')[0]);
    expect(await screen.findByTestId('group-name')).toHaveValue('PLC A raw');
    expect(screen.getByTestId('group-interval')).toHaveValue(60);
    expectNoOtherWrites();
    expect(api.applyCutover).not.toHaveBeenCalled();
  });

  it('keeps a failed review harmless and free of server text', async () => {
    api.reviewCutover.mockRejectedValue(httpFailure(409, 'revision_mismatch', 'conflict postgres://user:secret@db/prod'));
    await chooseForward({ review: false });
    fireEvent.click(screen.getByTestId('cutover-review-button'));
    expect(await screen.findByTestId('cutover-review-error')).toHaveTextContent('step4.cutover.errors.revision_mismatch');
    expect(document.body.textContent).not.toMatch(/secret|postgres:/);
    expect(screen.queryByTestId('cutover-confirm')).not.toBeInTheDocument();
    expect(api.applyCutover).not.toHaveBeenCalled();
  });
});

describe('The outcome of the one confirmed request is never guessed', () => {
  it('shows a refused apply as not applied, keeps the old group and requires a fresh review', async () => {
    api.applyCutover.mockRejectedValue(httpFailure(409, 'revision_mismatch', 'conflict postgres://user:secret@db/prod'));
    await chooseForward();
    fireEvent.click(screen.getByTestId('cutover-confirm'));
    expect(await screen.findByTestId('cutover-failed')).toHaveTextContent('step4.cutover.failed');
    expect(screen.getByTestId('cutover-failed')).toHaveTextContent('step4.cutover.errors.revision_mismatch');
    expect(document.body.textContent).not.toMatch(/secret|postgres:/);
    expect(screen.queryByTestId('cutover-record')).not.toBeInTheDocument();
    expect(screen.queryByTestId('cutover-confirm')).not.toBeInTheDocument();
    expect(api.applyCutover).toHaveBeenCalledTimes(1);
    expectNoOtherWrites();
  });

  it('treats a lost answer as unknown, finds the recorded cutover by what was reviewed and ignores unrelated ones', async () => {
    api.applyCutover.mockRejectedValue(Object.assign(new TypeError('Failed to fetch'), { outcome: 'unconfirmed' }));
    await chooseForward();
    fireEvent.click(screen.getByTestId('cutover-confirm'));
    expect(await screen.findByTestId('cutover-unknown')).toBeInTheDocument();
    expect(screen.queryByTestId('cutover-failed')).not.toBeInTheDocument();
    expect(screen.queryByTestId('cutover-record')).not.toBeInTheDocument();
    expect(screen.queryByTestId('cutover-confirm')).not.toBeInTheDocument();

    recorded = [recordOf({ id: 'cut-other', review_digest: 'c'.repeat(64) }), recordOf({ id: 'cut-late', boundary: '2099-01-03T00:00:00Z' })];
    fireEvent.click(screen.getByTestId('cutover-check'));
    expect(await screen.findByTestId('cutover-check-none')).toBeInTheDocument();
    expect(screen.queryByTestId('cutover-record')).not.toBeInTheDocument();

    recorded = [...recorded, recordOf({ id: 'cut-mine' })];
    fireEvent.click(screen.getByTestId('cutover-check'));
    expect(await screen.findByTestId('cutover-record-id')).toHaveTextContent('cut-mine');
    expect(api.applyCutover).toHaveBeenCalledTimes(1); // checking never applies anything
    expect(api.reviewCutover).toHaveBeenCalledTimes(1);
    expectNoOtherWrites();
  });

  it('resends only the identical reviewed request, and only after a check found nothing', async () => {
    api.applyCutover.mockRejectedValueOnce(httpFailure(502));
    await chooseForward();
    fireEvent.click(screen.getByTestId('cutover-confirm'));
    await screen.findByTestId('cutover-unknown');
    expect(screen.queryByTestId('cutover-resend')).not.toBeInTheDocument();
    fireEvent.click(screen.getByTestId('cutover-check'));
    await screen.findByTestId('cutover-check-none');
    fireEvent.click(screen.getByTestId('cutover-resend'));
    await screen.findByTestId('cutover-record');
    expect(api.applyCutover).toHaveBeenCalledTimes(2);
    expect(api.applyCutover.mock.calls[1]).toEqual(api.applyCutover.mock.calls[0]);
    expect(api.reviewCutover).toHaveBeenCalledTimes(1);
  });
});

describe('A rollback is a new reviewed boundary, never the old one', () => {
  it('lists the switch, reviews a rollback at a new boundary and applies it with its own digest', async () => {
    recorded = [recordOf()];
    renderHost();
    await screen.findByTestId('cutover-panel');
    fireEvent.click(screen.getByTestId('cutover-mode-rollback'));
    fireEvent.click(await screen.findByTestId('cutover-rollback-cut-1'));
    fireEvent.change(screen.getByTestId('cutover-boundary'), { target: { value: '2099-02-01T00:00:00' } });
    const rollbackReview = reviewOf({ kind: 'rollback', boundary: '2099-02-01T00:00:00.654321Z', digest: 'd'.repeat(64), targets: [], moved: [member(0), member(1)], unselected: [], remaining: [] });
    api.reviewCutover.mockResolvedValue(rollbackReview);
    fireEvent.click(screen.getByTestId('cutover-review-button'));
    await screen.findByTestId('cutover-review');
    const [sourceId, request] = api.reviewCutover.mock.calls[0];
    expect(sourceId).toBe('wide-1');
    expect(request).toMatchObject({ kind: 'rollback', cutover_id: 'cut-1', boundary: new Date('2099-02-01T00:00:00').toISOString() });
    expect(request.targets).toBeUndefined();
    expect(request.boundary).not.toBe(REVIEWED_BOUNDARY);
    expect(api.applyCutover).not.toHaveBeenCalled();

    api.applyCutover.mockResolvedValue(recordOf({ id: 'cut-2', kind: 'rollback', prior_cutover_id: 'cut-1', boundary: rollbackReview.boundary, review_digest: rollbackReview.digest, started: [{ group_id: 'wide-1', revision: 'rev-wide-3' }], ended: [{ group_id: 'tgt-1', revision: 'rev-tgt-1-1' }] }));
    fireEvent.click(screen.getByTestId('cutover-confirm'));
    await screen.findByTestId('cutover-record');
    expect(api.applyCutover.mock.calls[0][1]).toMatchObject({ kind: 'rollback', cutover_id: 'cut-1', boundary: '2099-02-01T00:00:00.654321Z', review_digest: 'd'.repeat(64) });
    expect(screen.getByTestId('cutover-record-boundary')).toHaveTextContent('2099-02-01T00:00:00.654321Z');
    expectNoOtherWrites();
  });

  it('offers nothing to roll back when no switch was recorded, and a switch already rolled back is not offered again', async () => {
    renderHost();
    await screen.findByTestId('cutover-panel');
    fireEvent.click(screen.getByTestId('cutover-mode-rollback'));
    expect(await screen.findByTestId('cutover-no-rollback')).toBeInTheDocument();
    expect(screen.getByTestId('cutover-review-button')).toBeDisabled();
  });

  it('drops the review of one mode when the other is chosen', async () => {
    await chooseForward();
    fireEvent.click(screen.getByTestId('cutover-mode-rollback'));
    expect(screen.queryByTestId('cutover-review')).not.toBeInTheDocument();
    expect(api.applyCutover).not.toHaveBeenCalled();
  });
});

describe('A prepared new wide mean target is chosen and reviewed like any other target', () => {
  it('offers the saved unapplied wide mean group, names it in the reviewed request and labels it as a wide layout', async () => {
    groups = [wide(), wideMeanTarget()];
    api.reviewCutover.mockResolvedValue(reviewOf({ targets: [view(wideMeanTarget(), { profile_kind: 'instantaneous' })] }));
    renderHost();
    await screen.findByTestId('cutover-panel');
    fireEvent.click(await screen.findByTestId('cutover-target-tgt-w'));
    fireEvent.change(screen.getByTestId('cutover-boundary'), { target: { value: BOUNDARY_LOCAL } });
    fireEvent.click(screen.getByTestId('cutover-review-button'));
    await screen.findByTestId('cutover-review');

    expect(api.reviewCutover).toHaveBeenCalledTimes(1);
    expect(api.reviewCutover.mock.calls[0][1]).toMatchObject({ kind: 'forward', targets: [{ group_id: 'tgt-w', expected_revision: 'rev-tgt-w-1' }], boundary: new Date(BOUNDARY_LOCAL).toISOString() });
    expect(within(screen.getByTestId('cutover-review')).getByTestId('cutover-review-target-tgt-w')).toHaveTextContent('step4.cutover.layout_wide');
    expect(api.applyCutover).not.toHaveBeenCalled();
    expectNoOtherWrites();
  });
});
