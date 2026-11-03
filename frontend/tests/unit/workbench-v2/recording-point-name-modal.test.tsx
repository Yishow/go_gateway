import { useRecordingPointNameBuffer } from '../../../src/features/datalink/workbench-v2/recording-editor/useRecordingPointNameBuffer';
import { RecordingDraftNavigation } from '../../../src/features/datalink/workbench-v2/recording-editor/RecordingDraftNavigation';
import type { PointNameOutcome, RecordingPointNamePort } from '../../../src/pages/datalink/workbench-v2/recordingPointNameSubmission';
import { act, fireEvent, render, renderHook, screen, waitFor, within } from '@testing-library/react';
import { expect, it, vi } from 'vitest';
import { useState } from 'react';
import { RecordingPointSelection } from '../../../src/features/datalink/workbench-v2/recording-editor/RecordingPointSelection';
import type { EditorMember } from '../../../src/features/datalink/workbench-v2/state/writeGroup/draft';
import type { GroupCandidate } from '../../../src/features/datalink/workbench-v2/state/writeGroup/candidates';
vi.mock('react-i18next', () => ({ useTranslation: () => ({ t: (key: string, values?: { count?: number }) => values?.count !== undefined ? `${key}:${values.count}` : key }) }));
const rows: GroupCandidate[] = [1, 2].map(n => ({ key: `P${n}|T${n}`, device_id: 'D', point_id: `P${n}`, tag_id: `T${n}`, tag_key: `tag${n}`, label: `Point ${n}`, device_name: 'PLC', address: `D${n}`, target_type: 'int16' }));
it('renames within the active picker without losing another pending point or the recording draft, and cancel writes nothing', async () => {
  const session = { id: 'S', identity: { device_id: 'D', point_id: 'P1', tag_id: 'T1' }, displayName: 'Point 1' };
  const name = { prepare: vi.fn(async () => ({ session, status: 'ready' as const })), release: vi.fn(), inspect: vi.fn(), submit: vi.fn(async (_id: string, value: string) => ({ session, status: 'confirmed' as const, currentName: value })) };
  const port = { submit: vi.fn(), inspect: vi.fn(), observations: () => [], name };
  function Harness() {
    const [members, setMembers] = useState<EditorMember[]>([]);
    return <><input aria-label="Recording draft" defaultValue="Keep me" /><output data-testid="members">{JSON.stringify(members)}</output>
      <RecordingPointSelection disabled={false} scope="W/G" workspaceId="W" candidates={rows} members={members} onChange={setMembers} deviceSource={port} devices={[{ id: 'D', name: 'PLC' }]} /></>;
  }
  render(<Harness />); fireEvent.click(screen.getByTestId('recording-points-open'));
  fireEvent.change(screen.getByLabelText('picker.search'), { target: { value: 'Point' } });
  fireEvent.click(screen.getByRole('checkbox', { name: 'Point 2' }));
  expect(screen.getByRole('textbox', { name: 'name.label Point 1' }).closest('.recording-point-picker__identity')).toContainElement(screen.getByRole('checkbox', { name: 'Point 1' }));
  const body = screen.getByRole('dialog').querySelector<HTMLElement>('.recording-point-picker__body')!;
  body.scrollTop = 230; fireEvent.scroll(body);
  act(() => screen.getByRole('textbox', { name: 'name.label Point 1' }).focus());
  expect(screen.getAllByRole('dialog')).toHaveLength(1);
  await waitFor(() => expect(screen.getByRole('textbox', { name: 'name.label Point 1' })).not.toHaveAttribute('readonly'));
  expect(screen.getByLabelText('picker.search')).toBeVisible();
  fireEvent.change(screen.getByRole('textbox', { name: 'name.label Point 1' }), { target: { value: 'Discarded' } });
  expect(screen.getByLabelText('picker.search')).toBeDisabled();
  expect(screen.getByLabelText('picker.device')).toBeDisabled();
  expect(screen.getByRole('button', { name: 'picker.add:1' })).toBeDisabled();
  expect(screen.getByRole('textbox', { name: 'name.label Point 2' })).toBeDisabled();
  fireEvent.click(within(screen.getByTestId('inline-point-name-P1')).getByRole('button', { name: 'picker.cancel' }));
  expect(screen.getByRole('textbox', { name: 'name.label Point 1' })).toHaveFocus();
  expect(screen.getByRole('dialog').querySelector<HTMLElement>('.recording-point-picker__body')!.scrollTop).toBe(230);
  expect(name.submit).not.toHaveBeenCalled(); expect(screen.getByRole('checkbox', { name: 'Point 2' })).toBeChecked();
  expect(screen.getByLabelText('picker.search')).toHaveValue('Point');
  act(() => { screen.getByLabelText('picker.search').focus(); screen.getByRole('textbox', { name: 'name.label Point 1' }).focus(); });
  await waitFor(() => expect(screen.getByRole('textbox', { name: 'name.label Point 1' })).not.toHaveAttribute('readonly'));
  expect(screen.getByLabelText('picker.search')).toBeVisible();
  fireEvent.change(screen.getByRole('textbox', { name: 'name.label Point 1' }), { target: { value: 'Point motor' } });
  fireEvent.click(screen.getByRole('button', { name: 'name.save' }));
  await screen.findByRole('textbox', { name: 'name.label Point motor' });
  expect(screen.getAllByRole('dialog')).toHaveLength(1);
  expect(screen.getByRole('checkbox', { name: 'Point 2' })).toBeChecked();
  expect(screen.getByRole('button', { name: 'picker.add:1' })).toBeEnabled();
  expect(screen.getByLabelText('picker.search')).toHaveValue('Point');
  expect(screen.getByTestId('members')).toHaveTextContent('[]'); expect(screen.getByLabelText('Recording draft')).toHaveValue('Keep me');
  expect(name.submit).toHaveBeenCalledExactlyOnceWith('S', 'Point motor'); expect(port.submit).not.toHaveBeenCalled();
});

function mountInlineGuard(name: RecordingPointNamePort) {
  render(<RecordingDraftNavigation><RecordingPointSelection disabled={false} scope="W/G" workspaceId="W" candidates={rows} members={[]} onChange={vi.fn()}
    deviceSource={{ submit: vi.fn(), inspect: vi.fn(), observations: () => [], name }} devices={[{ id: 'D', name: 'PLC' }]} /></RecordingDraftNavigation>);
  fireEvent.click(screen.getByTestId('recording-points-open'));
  fireEvent.click(screen.getByRole('checkbox', { name: 'Point 2' }));
  const body = screen.getByRole('dialog').querySelector<HTMLElement>('.recording-point-picker__body')!;
  body.scrollTop = 230; fireEvent.scroll(body);
  act(() => screen.getByRole('textbox', { name: 'name.label Point 1' }).focus());
}
const inlineSession = { id: 'S', identity: { device_id: 'D', point_id: 'P1', tag_id: 'T1' }, displayName: 'Point 1' };
it('keeps dirty inline name and pending selection on Keep, and explicit Discard writes nothing', async () => {
  const name: RecordingPointNamePort = { prepare: vi.fn(async () => ({ session: inlineSession, status: 'ready' as const })), submit: vi.fn(), inspect: vi.fn(), release: vi.fn() };
  mountInlineGuard(name);
  const input = screen.getByRole('textbox', { name: 'name.label Point 1' });
  await waitFor(() => expect(input).not.toHaveAttribute('readonly'));
  fireEvent.change(input, { target: { value: 'Unsent' } });
  fireEvent.keyDown(input, { key: 'Escape' });
  expect(screen.getByRole('alertdialog')).toBeVisible(); fireEvent.click(screen.getByRole('button', { name: 'navigation.keep' }));
  expect(input).toHaveValue('Unsent'); expect(screen.getByRole('checkbox', { name: 'Point 2' })).toBeChecked();
  expect(screen.getByRole('dialog').querySelector<HTMLElement>('.recording-point-picker__body')!.scrollTop).toBe(230);
  fireEvent.click(within(screen.getByRole('contentinfo')).getByRole('button', { name: 'picker.cancel' }));
  fireEvent.click(screen.getByRole('button', { name: 'navigation.discard' }));
  expect(screen.queryByRole('dialog')).not.toBeInTheDocument(); expect(name.submit).not.toHaveBeenCalled();
});
it('keeps pending and unknown inline name visible on picker close and offers original-result inspection', async () => {
  let resolve: ((value: PointNameOutcome) => void) | undefined;
  const name: RecordingPointNamePort = { prepare: vi.fn(async () => ({ session: inlineSession, status: 'ready' as const })),
    submit: vi.fn(() => new Promise<PointNameOutcome>(done => { resolve = done; })), inspect: vi.fn(async () => ({ session: inlineSession, status: 'unknown' as const, currentName: 'Observed only' })), release: vi.fn() };
  mountInlineGuard(name);
  const input = screen.getByRole('textbox', { name: 'name.label Point 1' });
  await waitFor(() => expect(input).not.toHaveAttribute('readonly'));
  fireEvent.change(input, { target: { value: 'Submitted' } }); fireEvent.click(screen.getByRole('button', { name: 'name.save' }));
  fireEvent.click(within(screen.getByRole('contentinfo')).getByRole('button', { name: 'picker.cancel' }));
  expect(screen.getByRole('dialog')).toBeVisible(); expect(input).toHaveValue('Submitted');
  fireEvent.pointerDown(input); expect(name.prepare).toHaveBeenCalledTimes(1); expect(name.submit).toHaveBeenCalledTimes(1);
  await act(async () => resolve?.({ session: inlineSession, status: 'unknown' }));
  fireEvent.keyDown(input, { key: 'Escape' });
  expect(screen.getByRole('dialog')).toBeVisible(); expect(screen.getByText('name.unknown')).toBeVisible();
  expect(screen.getByRole('button', { name: 'name.finishInspect' })).toBeVisible();
  fireEvent.pointerDown(input); expect(name.prepare).toHaveBeenCalledTimes(1); expect(name.submit).toHaveBeenCalledTimes(1);
  expect(screen.getByRole('checkbox', { name: 'Point 2' })).toBeChecked();
  fireEvent.click(screen.getByRole('button', { name: 'name.inspect' })); await waitFor(() => expect(name.inspect).toHaveBeenCalledExactlyOnceWith('S'));
  expect(name.submit).toHaveBeenCalledOnce(); expect(name.release).not.toHaveBeenCalled();
});

it('reopens the same focused inline name after Cancel and confirmed Save without duplicate preparation', async () => {
  let savedName = 'Point 1';
  const name: RecordingPointNamePort = { prepare: vi.fn(async () => ({ session: { ...inlineSession, displayName: savedName }, status: 'ready' as const })),
    submit: vi.fn(async (_id, value) => { savedName = value; return { session: inlineSession, status: 'confirmed' as const, currentName: value }; }),
    inspect: vi.fn(), release: vi.fn() };
  render(<RecordingDraftNavigation><RecordingPointSelection disabled={false} scope="W/G" workspaceId="W" candidates={rows} members={[]} onChange={vi.fn()}
    deviceSource={{ submit: vi.fn(), inspect: vi.fn(), observations: () => [], name }} devices={[{ id: 'D', name: 'PLC' }]} /></RecordingDraftNavigation>);
  fireEvent.click(screen.getByTestId('recording-points-open'));
  fireEvent.click(screen.getByRole('checkbox', { name: 'Point 2' }));
  fireEvent.change(screen.getByLabelText('picker.search'), { target: { value: 'Point' } });
  const body = screen.getByRole('dialog').querySelector<HTMLElement>('.recording-point-picker__body')!;
  body.scrollTop = 230; fireEvent.scroll(body);
  const input = screen.getByRole('textbox', { name: 'name.label Point 1' });
  fireEvent.pointerDown(input); act(() => input.focus());
  await waitFor(() => expect(input).not.toHaveAttribute('readonly'));
  expect(name.prepare).toHaveBeenCalledTimes(1);
  fireEvent.change(input, { target: { value: 'Discarded' } });
  fireEvent.click(within(screen.getByTestId('inline-point-name-P1')).getByRole('button', { name: 'picker.cancel' }));
  expect(input).toHaveFocus(); expect(input).toHaveAttribute('readonly'); expect(name.submit).not.toHaveBeenCalled();
  // Keep the real DOM activeElement; do not blur or synthesize another focus event.
  fireEvent.pointerDown(input);
  await waitFor(() => expect(input).not.toHaveAttribute('readonly'));
  expect(name.prepare).toHaveBeenCalledTimes(2);
  fireEvent.change(input, { target: { value: 'Point renamed' } });
  fireEvent.click(within(screen.getByTestId('inline-point-name-P1')).getByRole('button', { name: 'name.save' }));
  const confirmed = await screen.findByRole('textbox', { name: 'name.label Point renamed' });
  await waitFor(() => expect(confirmed).toHaveAttribute('readonly'));
  expect(confirmed).toHaveFocus();
  fireEvent.pointerDown(confirmed);
  await waitFor(() => expect(confirmed).not.toHaveAttribute('readonly'));
  expect(name.prepare).toHaveBeenCalledTimes(3); expect(name.submit).toHaveBeenCalledExactlyOnceWith('S', 'Point renamed');
  expect(screen.getByRole('checkbox', { name: 'Point 2' })).toBeChecked();
  expect(screen.getByLabelText('picker.search')).toHaveValue('Point'); expect(body.scrollTop).toBe(230);
});

it('shows safe source-review rejection and cancels only that rejected name while preserving pending points and search', async () => {
  const name: RecordingPointNamePort = { prepare: vi.fn(async () => ({ session: inlineSession, status: 'ready' as const })),
    submit: vi.fn(async () => ({ session: inlineSession, status: 'unknown' as const, reason: 'source-review' as const, canDiscardRejected: true })),
    discardRejected: vi.fn().mockResolvedValueOnce(false).mockResolvedValueOnce(true), inspect: vi.fn(), release: vi.fn() };
  mountInlineGuard(name);
  const input = screen.getByRole('textbox', { name: 'name.label Point 1' });
  await waitFor(() => expect(input).not.toHaveAttribute('readonly'));
  fireEvent.change(input, { target: { value: 'Rejected name' } }); fireEvent.click(screen.getByRole('button', { name: 'name.save' }));
  await screen.findByText('name_source_review'); expect(screen.queryByText('name.unknown')).not.toBeInTheDocument();
  fireEvent.click(screen.getByRole('button', { name: 'cancel_rejected_name' }));
  await waitFor(() => expect(name.discardRejected).toHaveBeenCalledTimes(1));
  expect(input).toHaveValue('Rejected name'); expect(screen.getByRole('button', { name: 'name.inspect' })).toBeVisible();
  await waitFor(() => expect(screen.getByRole('button', { name: 'cancel_rejected_name' })).toBeEnabled());
  fireEvent.click(screen.getByRole('button', { name: 'cancel_rejected_name' }));
  await waitFor(() => expect(input).toHaveValue('Point 1'));
  expect(screen.getByRole('checkbox', { name: 'Point 2' })).toBeChecked(); expect(screen.getByRole('button', { name: 'picker.add:1' })).toBeEnabled();
  expect(name.submit).toHaveBeenCalledTimes(1); expect(name.release).not.toHaveBeenCalled();
});

it('one-flights rejected cancellation in the buffer and recovers from a thrown cancel read', async () => {
  let resolve!: (value: boolean) => void;
  const name: RecordingPointNamePort = { prepare: vi.fn(async () => ({ session: inlineSession, status: 'ready' as const })),
    submit: vi.fn(async () => ({ session: inlineSession, status: 'unknown' as const, reason: 'source-review' as const, canDiscardRejected: true })),
    discardRejected: vi.fn().mockRejectedValueOnce(new Error('read unavailable')).mockImplementationOnce(() => new Promise<boolean>(done => { resolve = done; })), inspect: vi.fn(), release: vi.fn() };
  const h = renderHook(() => useRecordingPointNameBuffer('W/G', name));
  await act(async () => { await h.result.current.begin(rows[0]); });
  act(() => h.result.current.change('Rejected')); await act(async () => { await h.result.current.save(); });
  await act(async () => { await h.result.current.discardRejected(); });
  expect(h.result.current.pending).toBe(false); expect(h.result.current.value?.name).toBe('Rejected');
  let first: Promise<void>;
  act(() => { first = h.result.current.discardRejected(); void h.result.current.discardRejected(); void h.result.current.inspect(); });
  expect(name.discardRejected).toHaveBeenCalledTimes(2); expect(name.inspect).not.toHaveBeenCalled();
  await act(async () => { resolve(true); await first!; });
  expect(h.result.current.value).toBeUndefined(); expect(name.submit).toHaveBeenCalledTimes(1);
});
