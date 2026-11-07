import { useState } from 'react';
import { afterEach, describe, expect, it, vi } from 'vitest';
import { cleanup, fireEvent, render, screen, within } from '@testing-library/react';
import { RecordingHome } from '@/features/datalink/workbench-v2/recording-home/RecordingHome';
import { recordingHomeGuidance, type RecordingHomeGuidanceInput } from '@/features/datalink/workbench-v2/recording-home/recordingHomeGuidance';
import { recordingHomeResultGuidance } from '@/features/datalink/workbench-v2/recording-home/recordingHomeResultGuidance';
import { recordingHomeResultSummary } from '@/features/datalink/workbench-v2/recording-home/recordingHomeResultSummary';
import { homeOwners } from '@/features/datalink/workbench-v2/recording-home/recordingHomeGuidanceEvidence';
import { recordingHomeProjection } from '@/features/datalink/workbench-v2/recording-home/recordingHomeProjection';
import { RecordingHomeCollection } from '@/features/datalink/workbench-v2/recording-home/RecordingHomeCollection';
import { initialCollectionState } from '@/features/datalink/workbench-v2/recording-home/recordingHomeCollections';
import { parseStudioPresentation } from '@/services/studioRecordingView';
import type { RecordingOwnerObservations } from '@/features/datalink/workbench-v2/state/recordingOwnerObservation';
import type { WriteGroup } from '@/types/studioV2WriteGroup';
import { fixture, group, groupScope, scope } from './recording-home-guidance-fixtures';
vi.mock('react-i18next', () => ({ useTranslation: () => ({ t: (key: string) => key }) }));
afterEach(cleanup);

/** Mirrors readRecordingHomeOwners for a removed group whose schema/start history is not retained
 * in this browser while the server test-write recovery list is complete and empty. */
function unreadableHistory(value: WriteGroup, start: 'unavailable' | 'unresolved' = 'unavailable'): RecordingOwnerObservations {
  const bound = { ...groupScope(value), deviceId: 'D1' };
  return [
    { owner: 'schema', scope: bound, state: 'unavailable', reason: 'unreadable' },
    start === 'unresolved'
      ? { owner: 'start', scope: bound, state: 'unresolved', reason: 'unknown', operationId: 'OP-START' }
      : { owner: 'start', scope: bound, state: 'unavailable', reason: 'unreadable' },
    { owner: 'test_write', scope: bound, state: 'observed', result: 'none' },
  ];
}
function removedWithDrafts(start: 'unavailable' | 'unresolved' = 'unavailable'): RecordingHomeGuidanceInput {
  const removed = group('TEST', 'deleted');
  const drafts = ['A', 'B', 'C'].map(id => group(id, 'draft'));
  const input = fixture(['D1'], true, [removed, ...drafts]);
  return { ...input, owners: { ...input.owners,
    TEST: { state: 'available', scope: groupScope(removed), value: unreadableHistory(removed, start) } } };
}
/** Same composition as the actual route: base guidance, result summaries, then result guidance. */
function routeGuidance(input: RecordingHomeGuidanceInput) {
  const base = recordingHomeGuidance(input);
  const summaries = Object.fromEntries(base.recordingItems.map(item => [item.id, recordingHomeResultSummary(item.evidence, {
    owner: homeOwners(scope, item.evidence, input.owners?.[item.id]).inspection })]));
  return recordingHomeResultGuidance(base, summaries);
}

describe('a removed recording whose operation history is unconfirmed', () => {
  it('does not take over the primary message, yet stays a named unconfirmed item with its specific gaps', () => {
    for (const guidance of [recordingHomeGuidance(removedWithDrafts()), routeGuidance(removedWithDrafts())]) {
      expect(guidance.availability).toEqual({ state: 'available' });
      expect(guidance.reason).toBe('draft');
      expect(guidance.primaryAction).toMatchObject({ target: 'recording', groupId: 'A' });
      const removed = guidance.attentionItems.find(item => item.id === 'TEST');
      expect(removed).toMatchObject({ reason: 'evidence_unavailable', name: 'Recording TEST' });
      // Lifecycle and delivery were read; only schema/start history is unconfirmed. An empty
      // test-write list is never stretched into proof that schema/start work has ended.
      expect(removed?.unconfirmed).toEqual({ lifecycle: false, delivery: false, owners: ['schema', 'start'] });
    }
  });

  it('a known unresolved operation on the removed recording remains the primary, traceable guidance', () => {
    const guidance = routeGuidance(removedWithDrafts('unresolved'));
    expect(guidance.reason).toBe('operation_recovery');
    expect(guidance.primaryAction).toMatchObject({ target: 'recording', groupId: 'TEST', section: 'operation',
      intent: { kind: 'query_operation', operationOwner: 'start', operationId: 'OP-START' } });
  });

  it('does not demote an unconfirmed recording that is still active', () => {
    const live = group('LIVE'); const draft = group('A', 'draft');
    const input = fixture(['D1'], true, [live, draft]);
    const guidance = routeGuidance({ ...input, owners: { ...input.owners,
      LIVE: { state: 'available', scope: groupScope(live), value: unreadableHistory(live) } } });
    expect(guidance.reason).toBe('evidence_unavailable');
    expect(guidance.primaryAction).toMatchObject({ groupId: 'LIVE' });
  });
});

describe('home message names what is unconfirmed', () => {
  const location = { search: '?view=recordings', hash: '#saved' };
  it('keeps the draft message and lists the removed item, its status and gaps with a management check entry', () => {
    const navigate = vi.fn();
    render(<RecordingHome guidance={routeGuidance(removedWithDrafts())} location={location} onNavigate={navigate} onRetry={vi.fn()} />);
    const guidanceRegion = screen.getByTestId('recording-home-guidance');
    expect(within(guidanceRegion).getByText('reason.draft')).toBeInTheDocument();
    expect(screen.queryByText('reason.evidence_unavailable')).not.toBeInTheDocument();
    const notice = screen.getByRole('region', { name: 'unconfirmed.title' });
    const row = within(notice).getByText('Recording TEST').closest('li')!;
    expect(within(row).getByText('savedStatuses.deleted')).toBeInTheDocument();
    expect(within(row).getByText('unconfirmed.owner.schema')).toBeInTheDocument();
    expect(within(row).getByText('unconfirmed.owner.start')).toBeInTheDocument();
    expect(within(row).queryByText('unconfirmed.owner.test_write')).not.toBeInTheDocument();
    expect(within(row).getByText('unconfirmed.removed')).toBeInTheDocument();
    fireEvent.click(within(row).getByRole('button', { name: 'unconfirmed.inspect' }));
    // A removed item opens its original management entry, never a deleted editor.
    expect(navigate).toHaveBeenLastCalledWith('/studio/v2?step=output&target=database&view=advanced#saved');
    fireEvent.click(within(row).getByRole('button', { name: 'unconfirmed.logs' }));
    expect(navigate).toHaveBeenLastCalledWith('/studio/logs?entry=recording&workspace_id=W1&group_id=TEST&group_revision=TEST2#saved');
  });

  it('shows a known unresolved operation concretely and opens its operation query', () => {
    const navigate = vi.fn();
    render(<RecordingHome guidance={routeGuidance(removedWithDrafts('unresolved'))} location={location} onNavigate={navigate} onRetry={vi.fn()} />);
    expect(within(screen.getByTestId('recording-home-guidance')).getByText('reason.operation_recovery')).toBeInTheDocument();
    const row = within(screen.getByRole('region', { name: 'unconfirmed.title' })).getByText('Recording TEST').closest('li')!;
    expect(within(row).getByText('unconfirmed.unresolved.start')).toBeInTheDocument();
    fireEvent.click(within(row).getByRole('button', { name: 'unconfirmed.inspect' }));
    expect(parseStudioPresentation(new URL(navigate.mock.calls.at(-1)![0], 'http://localhost').search)).toMatchObject({
      repairIntent: { groupId: 'TEST', section: 'operation', operation: { kind: 'query_operation', owner: 'start', operationId: 'OP-START' } },
    });
  });

  it('an item-level gap does not read as the whole home failing; a whole-read gap still does', () => {
    const live = group('LIVE'); const input = fixture(['D1'], true, [live]);
    render(<RecordingHome guidance={routeGuidance({ ...input, owners: { LIVE: { state: 'available', scope: groupScope(live), value: unreadableHistory(live) } } })}
      location={location} onRetry={vi.fn()} />);
    expect(screen.getByText('reason.evidence_unconfirmed')).toBeInTheDocument();
    expect(screen.queryByText('reason.evidence_unavailable')).not.toBeInTheDocument();
    expect(within(screen.getByRole('region', { name: 'unconfirmed.title' })).getByText('Recording LIVE')).toBeInTheDocument();
    cleanup();
    const failed = fixture(); render(<RecordingHome guidance={recordingHomeGuidance({ ...failed, evidence: { ...failed.evidence, workspace: { state: 'error', scope } } })}
      location={location} onRetry={vi.fn()} />);
    expect(screen.getByText('reason.evidence_unavailable')).toBeInTheDocument();
  });

  it('a missing lifecycle or delivery read is named as such', () => {
    const live = group('LIVE'); const input = fixture(['D1'], true, [live]);
    const evidence = { ...input.evidence, lifecycle: {}, delivery: { LIVE: { state: 'error' as const, scope: groupScope(live) } } };
    const guidance = routeGuidance({ ...input, evidence });
    expect(recordingHomeProjection(evidence).groups.state).toBe('available');
    expect(guidance.attentionItems[0].unconfirmed).toEqual({ lifecycle: true, delivery: true, owners: [] });
    render(<RecordingHome guidance={guidance} location={location} onRetry={vi.fn()} />);
    const row = within(screen.getByRole('region', { name: 'unconfirmed.title' })).getByText('Recording LIVE').closest('li')!;
    expect(within(row).getByText('unconfirmed.read.lifecycle')).toBeInTheDocument();
    expect(within(row).getByText('unconfirmed.read.delivery')).toBeInTheDocument();
  });
});

describe('unconfirmed list boundaries', () => {
  const location = { search: '?view=recordings', hash: '' };
  function manyUnconfirmed(count: number, unrequested = 0): RecordingHomeGuidanceInput {
    const listed = Array.from({ length: count }, (_, index) => group(`U${index}`));
    const absent = Array.from({ length: unrequested }, (_, index) => group(`Z${index}`));
    const input = fixture(['D1'], true, [...listed, ...absent]);
    const lifecycle = Object.fromEntries(Object.entries(input.evidence.lifecycle ?? {}).filter(([id]) => !id.startsWith('Z')));
    const delivery = Object.fromEntries(Object.entries(input.evidence.delivery ?? {}).filter(([id]) => !id.startsWith('Z')));
    return { evidence: { ...input.evidence, lifecycle, delivery }, owners: Object.fromEntries(listed.map(value =>
      [value.id, { state: 'available' as const, scope: groupScope(value), value: unreadableHistory(value) }])) };
  }
  it('lists only items whose details were read for this page, never items that were not requested', () => {
    render(<RecordingHome guidance={routeGuidance(manyUnconfirmed(1, 3))} location={location} onRetry={vi.fn()} />);
    const notice = screen.getByRole('region', { name: 'unconfirmed.title' });
    expect(within(notice).getByText('Recording U0')).toBeInTheDocument();
    expect(within(notice).queryByText('Recording Z0')).not.toBeInTheDocument();
  });
  it('shows four, then a keyboard disclosure with the rest from the same evidence', () => {
    render(<RecordingHome guidance={routeGuidance(manyUnconfirmed(6))} location={location} onRetry={vi.fn()} />);
    const notice = screen.getByRole('region', { name: 'unconfirmed.title' });
    const more = within(notice).getByText('unconfirmed.more');
    expect(more.tagName).toBe('SUMMARY');
    expect(more.closest('details')).not.toHaveAttribute('open');
    fireEvent.click(more);
    expect(more.closest('details')).toHaveAttribute('open');
    for (const id of ['U0', 'U1', 'U2', 'U3', 'U4', 'U5']) expect(within(notice).getByText(`Recording ${id}`)).toBeInTheDocument();
  });
});

it('collection search and filters never hide the unconfirmed list', () => {
  const guidance = routeGuidance(removedWithDrafts());
  const groups = guidance.recordingItems.map(item => item.evidence.group);
  function Host() {
    const [state, setState] = useState(initialCollectionState());
    return <RecordingHome guidance={guidance} location={{ search: '', hash: '' }} onRetry={vi.fn()}
      collection={<RecordingHomeCollection groups={groups} devices={[]} state={state} guidance={guidance} observedIds={[]} unverifiedObserved={false}
        sourcesAvailable restored={{ state, workspaceRevision: 'W2', scrollTop: 0, focusId: null }}
        onChange={setState} onOpen={vi.fn()} onOpenDevice={vi.fn()} onRemember={vi.fn()} onRecheck={vi.fn()} />} />;
  }
  render(<Host />);
  fireEvent.change(screen.getByLabelText('keyword'), { target: { value: 'Recording A' } });
  fireEvent.change(screen.getByLabelText('savedStatus'), { target: { value: 'draft' } });
  expect(within(screen.getByTestId('home-management-table')).queryByText('Recording TEST')).not.toBeInTheDocument();
  expect(within(screen.getByRole('region', { name: 'unconfirmed.title' })).getByText('Recording TEST')).toBeInTheDocument();
});

it('the unconfirmed list stays visible and traceable on the recording home', () => {
  const navigate = vi.fn();
  render(<RecordingHome guidance={routeGuidance(removedWithDrafts('unresolved'))} location={{ search: '', hash: '' }} onNavigate={navigate} onRetry={vi.fn()} />);
  const row = within(screen.getByRole('region', { name: 'unconfirmed.title' })).getByText('Recording TEST').closest('li')!;
  expect(within(row).getByText('unconfirmed.unresolved.start')).toBeInTheDocument();
  fireEvent.click(within(row).getByRole('button', { name: 'unconfirmed.inspect' }));
  expect(navigate.mock.calls.at(-1)?.[0]).toContain('group_id=TEST');
});
