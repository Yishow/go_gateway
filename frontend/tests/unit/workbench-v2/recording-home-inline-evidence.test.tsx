import { cleanup, fireEvent, render, screen, waitFor, within } from '@testing-library/react';
import { afterEach, describe, expect, it, vi } from 'vitest';
import { RecordingHomeCollection } from '@/features/datalink/workbench-v2/recording-home/RecordingHomeCollection';
import { initialCollectionState } from '@/features/datalink/workbench-v2/recording-home/recordingHomeCollections';
import { recordingHomeGuidance } from '@/features/datalink/workbench-v2/recording-home/recordingHomeGuidance';
import { RecordingHomeResultFacts } from '@/features/datalink/workbench-v2/recording-home/RecordingHomeResultFacts';
import type { recordingHomeResultSummary } from '@/features/datalink/workbench-v2/recording-home/recordingHomeResultSummary';
import { fixture, group } from './recording-home-guidance-fixtures';

vi.mock('react-i18next', () => ({ useTranslation: () => ({ t: (key: string) => key }) }));
afterEach(cleanup);

type Summary = ReturnType<typeof recordingHomeResultSummary>;
const base: Summary = {
  saved: true, configuration: 'not_requested', effective: 'unknown', effectiveRevision: null,
  observedAt: null, requestedBoundary: null, sql: 'unknown', lastCommittedAt: null,
  historicalDestination: null, collection: 'unknown', current: 'unknown', owners: 'unknown',
  readback: 'unknown', cleanup: 'unknown',
};

describe('recording row write evidence', () => {
  it.each([
    [{}, ['unconfirmed', 'unconfirmed', 'unconfirmed', 'unconfirmed']],
    [{ configuration: 'requested', effective: 'requested_effective', sql: 'revision_receipt', readback: 'test_write_verified' }, ['complete', 'complete', 'complete', 'complete']],
    [{ configuration: 'requested', effective: 'requested_not_effective', sql: 'historical_commit' }, ['complete', 'attention', 'unconfirmed', 'unconfirmed']],
    [{ configuration: 'requested', effective: 'requested_effective', current: 'delivery_attention' }, ['complete', 'complete', 'attention', 'unconfirmed']],
  ] satisfies Array<[Partial<Summary>, string[]]>)('maps the specification example %j to four labelled stages', (fields, states) => {
    const rowFields: Partial<Summary> = fields;
    render(<RecordingHomeResultFacts id="R" name="Recording R" summary={{ ...base, ...fields }} />);
    const stages = [...screen.getByTestId('home-result-R').querySelectorAll<HTMLElement>('[data-write-stage]')];
    expect(stages.map(stage => stage.dataset.writeStage)).toEqual(['configuration', 'activation', 'sql', 'readback']);
    expect(stages.map(stage => stage.dataset.stageState)).toEqual(states);
    stages.forEach(stage => expect(stage).toHaveTextContent(/rowEvidence\./));
    if (rowFields.sql === 'historical_commit') expect(stages[2]).toHaveTextContent('rowEvidence.historical');
    if (rowFields.readback === 'test_write_verified') expect(stages[3]).toHaveTextContent('rowEvidence.testWriteReadback');
  });

  it.each([
    [{ sql: 'effective_revision_receipt' }, 'complete'],
    [{ sql: 'previous_revision_receipt' }, 'unconfirmed'],
    [{ sql: 'awaiting_first_write' }, 'unconfirmed'],
    [{ sql: 'revision_receipt', current: 'delivery_attention' }, 'attention'],
  ] satisfies Array<[Partial<Summary>, string]>)('keeps SQL evidence distinct for %j', (fields, expected) => {
    render(<RecordingHomeResultFacts id="R" name="Recording R" summary={{ ...base, ...fields }} />);
    const stage = screen.getByTestId('home-result-R').querySelector('[data-write-stage="sql"]');
    expect(stage).toHaveAttribute('data-stage-state', expected);
    if ((fields as Partial<Summary>).sql === 'previous_revision_receipt') expect(stage).toHaveTextContent('rowEvidence.historical');
  });

  it('makes the original observation, revision, owner and cleanup facts reachable from the focused stages', () => {
    render(<RecordingHomeResultFacts id="R" name="Recording R" summary={{ ...base, effectiveRevision: 'R3', observedAt: '2026-10-30T10:34:00Z',
      requestedBoundary: '2026-10-30T10:30:00Z', owners: 'inspected', readback: 'test_write_verified', cleanup: 'test_write_completed' }} />);
    const result = screen.getByTestId('home-result-R');
    const expected = [
      ['configuration', '2026-10-30T10:30:00Z'], ['activation', 'R3'], ['sql', 'result.owners.inspected'], ['readback', 'result.testWriteCleanup'],
    ];
    for (const [name, detail] of expected) {
      const stage = result.querySelector<HTMLElement>(`[data-write-stage="${name}"]`)!;
      stage.focus();
      expect(stage).toHaveFocus();
      const tooltip = within(stage).getByRole('tooltip', { hidden: true });
      expect(stage).toHaveAttribute('aria-describedby', tooltip.id);
      expect(tooltip).toHaveTextContent(detail);
    }
  });

  it('shows the last SQL commit as a relative time beside the SQL fact', () => {
    vi.spyOn(Date, 'now').mockReturnValue(Date.parse('2026-10-30T10:32:30Z'));
    render(<RecordingHomeResultFacts id="R" name="Recording R" summary={{ ...base, sql: 'historical_commit', lastCommittedAt: '2026-10-30T10:32:00Z' }} />);
    expect(screen.getByTestId('home-result-R').querySelector('time[datetime="2026-10-30T10:32:00Z"]')).toHaveTextContent('rowEvidence.secondsAgo');
    vi.restoreAllMocks();
  });
});

function renderRow() {
  const item = group('R');
  item.destination.table_name = 'gw_long_3fe66f8c5c98acb3';
  const state = initialCollectionState();
  render(<RecordingHomeCollection groups={[item]} devices={[]} state={state} guidance={recordingHomeGuidance(fixture([], false, [item]))}
    summaries={{ R: base }} observedIds={[]} unverifiedObserved={false} sourcesAvailable restored={{ state, workspaceRevision: 'W2', scrollTop: 0, focusId: null }}
    onChange={vi.fn()} onOpen={vi.fn()} onOpenDevice={vi.fn()} onRemember={vi.fn()} onRecheck={vi.fn()} />);
  return screen.getByText('Recording R').closest('tr')!;
}

describe('destination identity', () => {
  it('copies the exact schema and table without a disclosure', async () => {
    const writeText = vi.fn().mockResolvedValue(undefined);
    Object.defineProperty(navigator, 'clipboard', { configurable: true, value: { writeText } });
    const row = renderRow();
    expect(row.querySelectorAll('details')).toHaveLength(0);
    expect(within(row).getByTitle('public.gw_long_3fe66f8c5c98acb3')).toHaveTextContent('public.gw_long_3fe66f8c5c98acb3');
    fireEvent.click(within(row).getByRole('button', { name: 'rowEvidence.copyTable' }));
    await waitFor(() => expect(writeText).toHaveBeenCalledWith('public.gw_long_3fe66f8c5c98acb3'));
    expect(within(row).getByRole('status')).toHaveTextContent('rowEvidence.copied');
    expect(row).not.toHaveTextContent('DB1');
    expect(within(row).getByTitle('product.connector: DB1')).toBeInTheDocument();
  });
  it('selects the exact name and explains the failure when clipboard access is refused', async () => {
    const writeText = vi.fn().mockRejectedValue(new Error('private clipboard error'));
    Object.defineProperty(navigator, 'clipboard', { configurable: true, value: { writeText } });
    const row = renderRow();
    fireEvent.click(within(row).getByRole('button', { name: 'rowEvidence.copyTable' }));
    expect(await within(row).findByRole('status')).toHaveTextContent('rowEvidence.copyFailed');
    expect(window.getSelection()?.toString()).toBe('public.gw_long_3fe66f8c5c98acb3');
    expect(row).not.toHaveTextContent('private clipboard error');
  });
});
