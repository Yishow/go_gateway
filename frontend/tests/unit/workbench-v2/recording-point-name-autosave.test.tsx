import { act, fireEvent, render, screen, waitFor } from '@testing-library/react';
import { expect, it, vi } from 'vitest';
import { RecordingPointSelection } from '../../../src/features/datalink/workbench-v2/recording-editor/RecordingPointSelection';
import type { GroupCandidate } from '../../../src/features/datalink/workbench-v2/state/writeGroup/candidates';

vi.mock('react-i18next', () => ({ useTranslation: () => ({ t: (key: string) => key }) }));

it('saves a changed point name on blur without a separate save button', async () => {
  const candidate: GroupCandidate = { key: 'P|T', device_id: 'D', point_id: 'P', tag_id: 'T', tag_key: 'point', label: 'Point', device_name: 'PLC', address: 'D29', target_type: 'int16' };
  const session = { id: 'S', identity: { device_id: 'D', point_id: 'P', tag_id: 'T' }, displayName: 'Point' };
  const name = { prepare: vi.fn(async () => ({ session, status: 'ready' as const })), release: vi.fn(), inspect: vi.fn(),
    submit: vi.fn(async (_id: string, value: string) => ({ session, status: 'confirmed' as const, currentName: value })) };
  render(<RecordingPointSelection disabled={false} scope="W/G" workspaceId="W" candidates={[candidate]} members={[]} onChange={vi.fn()}
    deviceSource={{ submit: vi.fn(), inspect: vi.fn(), observations: () => [], name }} devices={[{ id: 'D', name: 'PLC' }]} />);
  fireEvent.click(screen.getByTestId('recording-points-open'));
  const input = screen.getByRole('textbox', { name: 'name.label Point' });
  act(() => input.focus());
  await waitFor(() => expect(input).not.toHaveAttribute('readonly'));
  fireEvent.change(input, { target: { value: 'Motor' } });
  fireEvent.blur(input);
  await waitFor(() => expect(name.submit).toHaveBeenCalledExactlyOnceWith('S', 'Motor'));
  expect(screen.queryByRole('button', { name: 'name.save' })).not.toBeInTheDocument();
});

it('keeps typing during source preparation and saves after the field loses focus', async () => {
  const candidate: GroupCandidate = { key: 'P|T', device_id: 'D', point_id: 'P', tag_id: 'T', tag_key: 'point', label: 'Point', device_name: 'PLC', address: 'D29', target_type: 'int16' };
  const session = { id: 'S', identity: { device_id: 'D', point_id: 'P', tag_id: 'T' }, displayName: 'Point' };
  let prepare!: (value: { session: typeof session; status: 'ready' }) => void;
  const name = { prepare: vi.fn(() => new Promise<{ session: typeof session; status: 'ready' }>(resolve => { prepare = resolve; })), release: vi.fn(), inspect: vi.fn(),
    submit: vi.fn(async (_id: string, value: string) => ({ session, status: 'confirmed' as const, currentName: value })) };
  render(<RecordingPointSelection disabled={false} scope="W/G" workspaceId="W" candidates={[candidate]} members={[]} onChange={vi.fn()}
    deviceSource={{ submit: vi.fn(), inspect: vi.fn(), observations: () => [], name }} devices={[{ id: 'D', name: 'PLC' }]} />);
  fireEvent.click(screen.getByTestId('recording-points-open'));
  const input = screen.getByRole('textbox', { name: 'name.label Point' });
  act(() => input.focus());
  expect(input).not.toHaveAttribute('readonly');
  fireEvent.change(input, { target: { value: 'Motor' } });
  fireEvent.blur(input);
  expect(name.submit).not.toHaveBeenCalled();
  await act(async () => prepare({ session, status: 'ready' }));
  await waitFor(() => expect(name.submit).toHaveBeenCalledExactlyOnceWith('S', 'Motor'));
});

it('cancels a stale blur save when the field is focused again during preparation', async () => {
  const candidate: GroupCandidate = { key: 'P|T', device_id: 'D', point_id: 'P', tag_id: 'T', tag_key: 'point', label: 'Point', device_name: 'PLC', address: 'D29', target_type: 'int16' };
  const session = { id: 'S', identity: { device_id: 'D', point_id: 'P', tag_id: 'T' }, displayName: 'Point' };
  let prepare!: (value: { session: typeof session; status: 'ready' }) => void;
  const name = { prepare: vi.fn(() => new Promise<{ session: typeof session; status: 'ready' }>(resolve => { prepare = resolve; })), release: vi.fn(), inspect: vi.fn(),
    submit: vi.fn(async (_id: string, value: string) => ({ session, status: 'confirmed' as const, currentName: value })) };
  render(<RecordingPointSelection disabled={false} scope="W/G" workspaceId="W" candidates={[candidate]} members={[]} onChange={vi.fn()}
    deviceSource={{ submit: vi.fn(), inspect: vi.fn(), observations: () => [], name }} devices={[{ id: 'D', name: 'PLC' }]} />);
  fireEvent.click(screen.getByTestId('recording-points-open'));
  const input = screen.getByRole('textbox', { name: 'name.label Point' });
  act(() => input.focus());
  fireEvent.change(input, { target: { value: 'First draft' } });
  fireEvent.blur(input);
  fireEvent.focus(input);
  fireEvent.change(input, { target: { value: 'Still editing' } });
  await act(async () => prepare({ session, status: 'ready' }));
  expect(name.submit).not.toHaveBeenCalled();
  expect(input).toHaveFocus();
  fireEvent.blur(input);
  await waitFor(() => expect(name.submit).toHaveBeenCalledExactlyOnceWith('S', 'Still editing'));
});

it('keeps an explicit Enter save intent while preparation is pending', async () => {
  const candidate: GroupCandidate = { key: 'P|T', device_id: 'D', point_id: 'P', tag_id: 'T', tag_key: 'point', label: 'Point', device_name: 'PLC', address: 'D29', target_type: 'int16' };
  const session = { id: 'S', identity: { device_id: 'D', point_id: 'P', tag_id: 'T' }, displayName: 'Point' };
  let prepare!: (value: { session: typeof session; status: 'ready' }) => void;
  const name = { prepare: vi.fn(() => new Promise<{ session: typeof session; status: 'ready' }>(resolve => { prepare = resolve; })), release: vi.fn(), inspect: vi.fn(),
    submit: vi.fn(async (_id: string, value: string) => ({ session, status: 'confirmed' as const, currentName: value })) };
  render(<RecordingPointSelection disabled={false} scope="W/G" workspaceId="W" candidates={[candidate]} members={[]} onChange={vi.fn()}
    deviceSource={{ submit: vi.fn(), inspect: vi.fn(), observations: () => [], name }} devices={[{ id: 'D', name: 'PLC' }]} />);
  fireEvent.click(screen.getByTestId('recording-points-open'));
  const input = screen.getByRole('textbox', { name: 'name.label Point' });
  act(() => input.focus());
  fireEvent.change(input, { target: { value: 'Motor' } });
  fireEvent.keyDown(input, { key: 'Enter' });
  expect(name.submit).not.toHaveBeenCalled();
  await act(async () => prepare({ session, status: 'ready' }));
  await waitFor(() => expect(name.submit).toHaveBeenCalledExactlyOnceWith('S', 'Motor'));
});

it('does not replace an Enter intent with blur when the field is refocused before preparation finishes', async () => {
  const candidate: GroupCandidate = { key: 'P|T', device_id: 'D', point_id: 'P', tag_id: 'T', tag_key: 'point', label: 'Point', device_name: 'PLC', address: 'D29', target_type: 'int16' };
  const session = { id: 'S', identity: { device_id: 'D', point_id: 'P', tag_id: 'T' }, displayName: 'Point' };
  let prepare!: (value: { session: typeof session; status: 'ready' }) => void;
  const name = { prepare: vi.fn(() => new Promise<{ session: typeof session; status: 'ready' }>(resolve => { prepare = resolve; })), release: vi.fn(), inspect: vi.fn(),
    submit: vi.fn(async (_id: string, value: string) => ({ session, status: 'confirmed' as const, currentName: value })) };
  render(<RecordingPointSelection disabled={false} scope="W/G" workspaceId="W" candidates={[candidate]} members={[]} onChange={vi.fn()}
    deviceSource={{ submit: vi.fn(), inspect: vi.fn(), observations: () => [], name }} devices={[{ id: 'D', name: 'PLC' }]} />);
  fireEvent.click(screen.getByTestId('recording-points-open'));
  const input = screen.getByRole('textbox', { name: 'name.label Point' });
  act(() => input.focus());
  fireEvent.change(input, { target: { value: 'Motor' } });
  fireEvent.keyDown(input, { key: 'Enter' });
  fireEvent.blur(input);
  fireEvent.focus(input);
  await act(async () => prepare({ session, status: 'ready' }));
  await waitFor(() => expect(name.submit).toHaveBeenCalledExactlyOnceWith('S', 'Motor'));
});

it('preserves the draft and refuses to overwrite a name that changed during preparation', async () => {
  const candidate: GroupCandidate = { key: 'P|T', device_id: 'D', point_id: 'P', tag_id: 'T', tag_key: 'point', label: 'Point', device_name: 'PLC', address: 'D29', target_type: 'int16' };
  const session = { id: 'S', identity: { device_id: 'D', point_id: 'P', tag_id: 'T' }, displayName: 'Changed elsewhere' };
  let prepare!: (value: { session: typeof session; status: 'ready' }) => void;
  const name = { prepare: vi.fn(() => new Promise<{ session: typeof session; status: 'ready' }>(resolve => { prepare = resolve; })),
    release: vi.fn(), inspect: vi.fn(), submit: vi.fn() };
  render(<RecordingPointSelection disabled={false} scope="W/G" workspaceId="W" candidates={[candidate]} members={[]} onChange={vi.fn()}
    deviceSource={{ submit: vi.fn(), inspect: vi.fn(), observations: () => [], name }} devices={[{ id: 'D', name: 'PLC' }]} />);
  fireEvent.click(screen.getByTestId('recording-points-open'));
  const input = screen.getByRole('textbox', { name: 'name.label Point' });
  act(() => input.focus()); fireEvent.change(input, { target: { value: 'My draft' } }); fireEvent.blur(input);
  await act(async () => prepare({ session, status: 'ready' }));
  expect(screen.getByText('name.blocked')).toBeVisible();
  expect(input).toHaveValue('My draft'); expect(input).toHaveAttribute('readonly');
  expect(name.submit).not.toHaveBeenCalled();
});

it.each(['unavailable', 'rejected'] as const)('keeps typed text when source preparation is %s', async failure => {
  const candidate: GroupCandidate = { key: 'P|T', device_id: 'D', point_id: 'P', tag_id: 'T', tag_key: 'point', label: 'Point', device_name: 'PLC', address: 'D29', target_type: 'int16' };
  let finish!: () => void;
  const name = { prepare: vi.fn(() => new Promise<undefined>((resolve, reject) => { finish = () => failure === 'rejected' ? reject(new Error('unavailable')) : resolve(undefined); })),
    release: vi.fn(), inspect: vi.fn(), submit: vi.fn() };
  render(<RecordingPointSelection disabled={false} scope="W/G" workspaceId="W" candidates={[candidate]} members={[]} onChange={vi.fn()}
    deviceSource={{ submit: vi.fn(), inspect: vi.fn(), observations: () => [], name }} devices={[{ id: 'D', name: 'PLC' }]} />);
  fireEvent.click(screen.getByTestId('recording-points-open'));
  const input = screen.getByRole('textbox', { name: 'name.label Point' });
  act(() => input.focus()); fireEvent.change(input, { target: { value: 'My draft' } }); fireEvent.blur(input);
  await act(async () => finish());
  expect(screen.getByText('name.blocked')).toBeVisible();
  expect(input).toHaveValue('My draft'); expect(name.submit).not.toHaveBeenCalled();
});

it('releases a late preparation after an edit is cancelled before it finishes', async () => {
  const candidate: GroupCandidate = { key: 'P|T', device_id: 'D', point_id: 'P', tag_id: 'T', tag_key: 'point', label: 'Point', device_name: 'PLC', address: 'D29', target_type: 'int16' };
  const session = { id: 'S', identity: { device_id: 'D', point_id: 'P', tag_id: 'T' }, displayName: 'Point' };
  let prepare!: (value: { session: typeof session; status: 'ready' }) => void;
  const name = { prepare: vi.fn(() => new Promise<{ session: typeof session; status: 'ready' }>(resolve => { prepare = resolve; })),
    release: vi.fn(), inspect: vi.fn(), submit: vi.fn() };
  render(<RecordingPointSelection disabled={false} scope="W/G" workspaceId="W" candidates={[candidate]} members={[]} onChange={vi.fn()}
    deviceSource={{ submit: vi.fn(), inspect: vi.fn(), observations: () => [], name }} devices={[{ id: 'D', name: 'PLC' }]} />);
  fireEvent.click(screen.getByTestId('recording-points-open'));
  const input = screen.getByRole('textbox', { name: 'name.label Point' });
  act(() => input.focus()); fireEvent.change(input, { target: { value: 'Discarded' } }); fireEvent.keyDown(input, { key: 'Escape' });
  await act(async () => prepare({ session, status: 'ready' }));
  expect(name.release).toHaveBeenCalledExactlyOnceWith('S');
  expect(name.submit).not.toHaveBeenCalled(); expect(input).toHaveValue('Point');
});

it('finishes the first save before editing the newly focused point', async () => {
  const candidates: GroupCandidate[] = [1, 2].map(n => ({ key: `P${n}|T${n}`, device_id: 'D', point_id: `P${n}`, tag_id: `T${n}`,
    tag_key: `point${n}`, label: `Point ${n}`, device_name: 'PLC', address: `D${n}`, target_type: 'int16' }));
  const sessions = candidates.map((candidate, index) => ({ id: `S${index + 1}`, identity: candidate, displayName: candidate.label }));
  let complete!: (value: { session: typeof sessions[number]; status: 'confirmed'; currentName: string }) => void;
  const name = { prepare: vi.fn(async (identity: { point_id: string }) => ({ session: sessions.find(session => session.identity.point_id === identity.point_id)!, status: 'ready' as const })),
    release: vi.fn(), inspect: vi.fn(), submit: vi.fn((_id: string, value: string) => new Promise<{ session: typeof sessions[number]; status: 'confirmed'; currentName: string }>(resolve => { complete = result => resolve({ ...result, currentName: value }); })) };
  render(<RecordingPointSelection disabled={false} scope="W/G" workspaceId="W" candidates={candidates} members={[]} onChange={vi.fn()}
    deviceSource={{ submit: vi.fn(), inspect: vi.fn(), observations: () => [], name }} devices={[{ id: 'D', name: 'PLC' }]} />);
  fireEvent.click(screen.getByTestId('recording-points-open'));
  const first = screen.getByRole('textbox', { name: 'name.label Point 1' }), second = screen.getByRole('textbox', { name: 'name.label Point 2' });
  act(() => first.focus()); await waitFor(() => expect(name.prepare).toHaveBeenCalledTimes(1));
  await waitFor(() => expect(first).not.toHaveAttribute('readonly'));
  fireEvent.change(first, { target: { value: 'Motor' } });
  act(() => second.focus());
  await waitFor(() => expect(name.submit).toHaveBeenCalledExactlyOnceWith('S1', 'Motor'));
  await act(async () => complete({ session: sessions[0], status: 'confirmed', currentName: 'Motor' }));
  await waitFor(() => expect(name.prepare).toHaveBeenCalledTimes(2));
  await waitFor(() => expect(second).not.toHaveAttribute('readonly'));
  expect(second).toHaveFocus();
  fireEvent.change(second, { target: { value: 'Pump' } });
  fireEvent.keyDown(second, { key: 'Enter' });
  await waitFor(() => expect(name.submit).toHaveBeenCalledTimes(2));
  expect(name.submit).toHaveBeenLastCalledWith('S2', 'Pump');
  await act(async () => complete({ session: sessions[1], status: 'confirmed', currentName: 'Pump' }));
});
