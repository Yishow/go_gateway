import { act, fireEvent, render, screen, waitFor, within } from '@testing-library/react';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import { AcquisitionPlanPreviewPanel } from '@/features/datalink/workbench-v2/acquisition-preview/AcquisitionPlanPreviewPanel';
import { readAcquisitionPlanPreview } from '@/services/acquisitionPlanPreview';
import type { AcquisitionActual, AcquisitionPlanPreview, AcquisitionPlanPreviewResult, AcquisitionPreviewGroup, AcquisitionPreviewScope } from '@/types/acquisitionPlanPreview';
vi.mock('@/services/acquisitionPlanPreview', () => ({ readAcquisitionPlanPreview: vi.fn() }));
vi.mock('react-i18next', () => ({ useTranslation: () => ({ t: (key: string, values?: Record<string, unknown>) => `${key}${values ? ` ${JSON.stringify(values)}` : ''}` }) }));
const read = vi.mocked(readAcquisitionPlanPreview);
const scope: AcquisitionPreviewScope = { workspace_id: 'W', device_id: 'D', expected_workspace_revision: 'R' };
function group(id = 'G'): AcquisitionPreviewGroup {
  return { polling_group_id: id, poll_interval_ms: 1000, plan_revision: 'P', declared_plan_revision: 'P', capability_revision: 'C',
    basis: 'declared_no_connection', acquisition_point_count: 4, planned_requests: 3,
    unique_wire_units: [{ area: 'DR', function: 'read', address_unit_bits: 16, unique_address_units: 3 },
      { area: 'WM', function: 'read', address_unit_bits: 1, unique_address_units: 16 }],
    point_errors: [{ point_id: 'rejected-point', code: 'invalid_address' }], fallback_reasons: ['capability_unavailable'],
    actual: { status: 'unavailable', reason: 'no_matching_scheduled_scan' } };
}
function preview(groups = [group()]): AcquisitionPlanPreview {
  return { workspace_id: 'W', device_id: 'D', workspace_revision: 'R', status: 'available',
    preview_revision: 'HEAD', estimate: true, groups, output_loss: { status: 'unavailable', reason: 'no_retained_detail' } };
}
function observed(): Extract<AcquisitionActual, { status: 'available' }> {
  return { status: 'available', acquisition_id: 'A', observed_at: '2026-10-31T12:00:00Z', declared_plan_revision: 'P', initial_plan_revision: 'P',
    phases: [{ phase: 'initial', plan_revision: 'P' }, { phase: 'fallback', plan_revision: 'SAFE' }], wire_attempts: 5,
    retry_wire_attempts: 2, fallback_wire_attempts: 3, logical_requests: { status: 'observed', succeeded: 1, failed: 1, canceled: 1 },
    cycle_latency_ms: 1200, overrun: true, last_good_age_ms: null };
}
function reply(value: AcquisitionPlanPreview = preview()) { read.mockResolvedValue({ ok: true, value }); }
function deferred() { let resolve!: (value: AcquisitionPlanPreviewResult) => void;
  const promise = new Promise<AcquisitionPlanPreviewResult>(yes => { resolve = yes; }); return { promise, resolve }; }
beforeEach(() => { read.mockReset(); });

describe('read-only canonical acquisition preview consumer', () => {
  it('does not fetch without proven scope and explains unsaved settings', () => {
    render(<AcquisitionPlanPreviewPanel unavailableReason="unsaved" />);
    expect(screen.getByText('unsaved')).toBeInTheDocument();
    expect(read).not.toHaveBeenCalled();
    expect(screen.queryByTestId('acquisition-planned-G')).not.toBeInTheDocument();
  });
  it('uses the existing source card and styled secondary refresh control', async () => {
    reply(); render(<AcquisitionPlanPreviewPanel scope={scope} />);
    await screen.findByTestId('acquisition-group-G');
    const card = screen.getByRole('region', { name: 'title' });
    expect(card).toHaveClass('wbv2-surface');
    expect(screen.getByRole('button', { name: 'refresh' })).toHaveClass('inline-flex', 'text-xs', 'border');
  });
  it('shows every original group cadence and estimates separately from unavailable actual evidence', async () => {
    reply(preview([group(), { ...group('SLOW'), poll_interval_ms: 5000, planned_requests: 7 }]));
    render(<AcquisitionPlanPreviewPanel scope={scope} />);
    expect(await screen.findByTestId('acquisition-planned-G')).toHaveTextContent('3');
    expect(screen.getByTestId('acquisition-planned-SLOW')).toHaveTextContent('7');
    expect(screen.getByTestId('acquisition-cadence-G')).toHaveTextContent('1000');
    expect(screen.getByTestId('acquisition-cadence-SLOW')).toHaveTextContent('5000');
    expect(screen.getAllByText('actualUnavailable')).toHaveLength(2);
    expect(screen.queryByTestId('acquisition-wire-attempts-G')).not.toBeInTheDocument();
    expect(read.mock.calls[0][0]).toEqual(scope);
  });
  it('keeps bit cells and 16-bit address cells distinct, retaining planner errors and fallback reasons', async () => {
    reply(); render(<AcquisitionPlanPreviewPanel scope={scope} />);
    const row = await screen.findByTestId('acquisition-group-G');
    expect(within(row).getByTestId('acquisition-points-G')).toHaveTextContent('4');
    expect(within(row).getByText(/wireUnits.word.*"count":3/)).toBeInTheDocument();
    expect(within(row).getByText(/wireUnits.bit.*"count":16/)).toBeInTheDocument();
    expect(within(row).getByText('invalid_address').closest('details')).not.toHaveAttribute('open');
    expect(within(row).getByText('capability_unavailable').closest('details')).not.toHaveAttribute('open');
  });
  it('separates overlapping wire retry/fallback subsets from logical request outcomes', async () => {
    reply(preview([{ ...group(), actual: observed() }])); render(<AcquisitionPlanPreviewPanel scope={scope} />);
    expect(await screen.findByTestId('acquisition-wire-attempts-G')).toHaveTextContent('5');
    expect(screen.getByTestId('acquisition-retry-G')).toHaveTextContent('2');
    expect(screen.getByTestId('acquisition-fallback-G')).toHaveTextContent('3');
    expect(screen.getByTestId('acquisition-logical-G')).toHaveTextContent('logicalCounts');
    expect(screen.getByText('wireSubsetNote')).toBeInTheDocument();
    expect(screen.getByText('lastGoodUnavailable')).toBeInTheDocument();
  });
  it('renders an observed true zero without inventing logical outcomes for not_started', async () => {
    reply(preview([{ ...group(), actual: { ...observed(), wire_attempts: 0, retry_wire_attempts: 0,
      fallback_wire_attempts: 0, logical_requests: { status: 'not_started' } } }]));
    render(<AcquisitionPlanPreviewPanel scope={scope} />);
    expect(await screen.findByTestId('acquisition-wire-attempts-G')).toHaveTextContent('0');
    expect(screen.getByText('logicalNotStarted')).toBeInTheDocument();
    expect(screen.queryByTestId('acquisition-logical-G')).not.toBeInTheDocument();
  });
  it('marks retained loss detail as partial rather than a scoped total or zero history', async () => {
    const value = preview(); if (value.status === 'available') value.output_loss = { status: 'partial', history_complete: false,
      retained_detail: [{ point_id: 'point-1', losses: 2 }] };
    reply(value); render(<AcquisitionPlanPreviewPanel scope={scope} />);
    expect(await screen.findByText('lossPartial')).toBeInTheDocument();
    expect(screen.getByText(/lossRow.*"count":2/)).toBeInTheDocument();
  });
  it('clears an estimate on stale 409, then explicitly reacquires without silently rebasing the old request', async () => {
    read.mockResolvedValueOnce({ ok: true, value: preview() }).mockResolvedValueOnce({ ok: false, reason: 'stale-reference' })
      .mockResolvedValueOnce({ ok: true, value: preview() });
    render(<AcquisitionPlanPreviewPanel scope={scope} />);
    await screen.findByTestId('acquisition-planned-G');
    fireEvent.click(screen.getByRole('button', { name: 'refresh' }));
    expect(screen.queryByTestId('acquisition-planned-G')).not.toBeInTheDocument();
    await screen.findByText('stale');
    expect(read.mock.calls[1][0]).toEqual({ ...scope, expected_preview_revision: 'HEAD' });
    expect(read).toHaveBeenCalledTimes(2);
    fireEvent.click(screen.getByRole('button', { name: 'refresh' }));
    await screen.findByTestId('acquisition-planned-G');
    expect(read.mock.calls[2][0]).toEqual(scope);
  });
  it.each(['foreign-scope', 'unavailable'] as const)('clears earlier data on %s without leaking raw errors', async reason => {
    read.mockResolvedValueOnce({ ok: true, value: preview() }).mockResolvedValueOnce({ ok: false, reason });
    render(<AcquisitionPlanPreviewPanel scope={scope} />); await screen.findByTestId('acquisition-planned-G');
    fireEvent.click(screen.getByRole('button', { name: 'refresh' }));
    await screen.findByText(reason === 'foreign-scope' ? 'foreign' : 'unavailable');
    expect(screen.queryByTestId('acquisition-planned-G')).not.toBeInTheDocument();
  });
  it('ignores late old-scope responses and aborts their read', async () => {
    const old = deferred(); read.mockReturnValueOnce(old.promise).mockResolvedValueOnce({ ok: true,
      value: { workspace_id: 'W', device_id: 'OTHER', workspace_revision: 'R', status: 'unavailable', reason: 'not_installed' } });
    const view = render(<AcquisitionPlanPreviewPanel scope={scope} />);
    await waitFor(() => expect(read).toHaveBeenCalledTimes(1));
    const signal = read.mock.calls[0][1];
    view.rerender(<AcquisitionPlanPreviewPanel scope={{ ...scope, device_id: 'OTHER' }} />);
    await screen.findByText('not_installed');
    expect(signal?.aborted).toBe(true);
    await act(async () => old.resolve({ ok: true, value: preview() }));
    expect(screen.queryByTestId('acquisition-planned-G')).not.toBeInTheDocument();
  });
  it('hides old results immediately when saved settings become dirty', async () => {
    reply(); const view = render(<AcquisitionPlanPreviewPanel scope={scope} />); await screen.findByTestId('acquisition-planned-G');
    view.rerender(<AcquisitionPlanPreviewPanel unavailableReason="unsaved" />);
    expect(screen.queryByTestId('acquisition-planned-G')).not.toBeInTheDocument();
    expect(screen.getByText('unsaved')).toBeInTheDocument(); expect(read).toHaveBeenCalledTimes(1);
  });
  it('rejects a response bound to another workspace revision', async () => {
    reply({ ...preview(), workspace_revision: 'OLD' }); render(<AcquisitionPlanPreviewPanel scope={scope} />);
    await screen.findByText('unavailable'); expect(screen.queryByTestId('acquisition-planned-G')).not.toBeInTheDocument();
  });
  it('does not revive an old receipt after switching away and back to the same scope', async () => {
    reply(); const view = render(<AcquisitionPlanPreviewPanel scope={scope} />); await screen.findByTestId('acquisition-planned-G');
    const pending = deferred(); read.mockReturnValue(pending.promise);
    view.rerender(<AcquisitionPlanPreviewPanel scope={{ ...scope, device_id: 'OTHER' }} />);
    view.rerender(<AcquisitionPlanPreviewPanel scope={scope} />);
    expect(screen.queryByTestId('acquisition-planned-G')).not.toBeInTheDocument();
    expect(screen.getByText('loading')).toBeInTheDocument();
    expect(read.mock.calls.at(-1)?.[0]).toEqual(scope);
  });
  it('does not display actual values whose plan fingerprint no longer matches the group estimate', async () => {
    reply(preview([{ ...group(), actual: { ...observed(), initial_plan_revision: 'OLD' } }]));
    render(<AcquisitionPlanPreviewPanel scope={scope} />);
    await screen.findByTestId('acquisition-planned-G');
    expect(screen.getByText('actualUnavailable')).toBeInTheDocument();
    expect(screen.queryByTestId('acquisition-wire-attempts-G')).not.toBeInTheDocument();
  });
  it('keeps proven zero planned requests instead of treating zero as missing', async () => {
    reply(preview([{ ...group(), acquisition_point_count: 0, planned_requests: 0, unique_wire_units: [], point_errors: [], fallback_reasons: [] }]));
    render(<AcquisitionPlanPreviewPanel scope={scope} />);
    expect(await screen.findByTestId('acquisition-planned-G')).toHaveTextContent('0');
    expect(screen.getByTestId('acquisition-points-G')).toHaveTextContent('0');
  });
  it('handles rejected reads as unavailable rather than exposing an exception', async () => {
    read.mockRejectedValue(new Error('secret endpoint must not render'));
    render(<AcquisitionPlanPreviewPanel scope={scope} />);
    await screen.findByText('unavailable'); expect(screen.queryByText(/secret endpoint/)).not.toBeInTheDocument();
  });
  it('shows canonical request limits and keeps unknown frame cost unavailable rather than zero', async () => {
    reply(preview([{ ...group(), request_limits: [{ kind: 'contiguous', area: 'D', function: 'read', native_unit_limit: 64,
      request_frame_bytes: 12, response_frame_bytes: 133, frame_cost_known: true }, { kind: 'random', area: 'DR', function: 'read',
      native_unit_limit: 16, request_frame_bytes: 0, response_frame_bytes: 0, frame_cost_known: false }] }]));
    render(<AcquisitionPlanPreviewPanel scope={scope} />);
    const details = await screen.findByTestId('acquisition-limits-G');
    expect(within(details).getByText(/limitRow.*"limit":64/)).toBeInTheDocument();
    expect(within(details).getByText(/frameCost.*"request":12.*"response":133/)).toBeInTheDocument();
    expect(within(details).getByText('frameCostUnavailable')).toBeInTheDocument();
  });
  it('shows unavailable installed context without making empty estimates', async () => {
    reply({ workspace_id: 'W', device_id: 'D', workspace_revision: 'R', status: 'unavailable', reason: 'context_unproven' });
    render(<AcquisitionPlanPreviewPanel scope={scope} />);
    expect((await screen.findByText('context_unproven')).closest('details')).not.toHaveAttribute('open'); expect(screen.queryByTestId('acquisition-planned-G')).not.toBeInTheDocument();
  });
});
