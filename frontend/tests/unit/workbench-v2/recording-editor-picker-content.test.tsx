import { useState } from 'react';
import { describe, expect, it, vi } from 'vitest';
import { act, fireEvent, render, screen, within } from '@testing-library/react';
import { RecordingPointPicker } from '../../../src/features/datalink/workbench-v2/recording-editor/RecordingPointPicker';
import type { PickerLabels, RecordingPointPickerProps } from '../../../src/features/datalink/workbench-v2/recording-editor/RecordingPointPicker';
import { createPicker, sourceKey } from '../../../src/features/datalink/workbench-v2/recording-editor/recordingPickerState';
import type { GroupCandidate } from '../../../src/features/datalink/workbench-v2/state/writeGroup/candidates';
const labels: PickerLabels = {
  rename: 'Change name',
  title: 'Select points', search: 'Search points', device: 'Device', pointType: 'Point type', allDevices: 'All devices', allTypes: 'All types',
  types: { gauge: 'Gauge', counter: 'Counter', state: 'State' }, unknownType: 'Unknown meaning', dataType: 'Data type', meaningUnknown: 'Some points have no measurement purpose yet', point: 'Point', address: 'Address', alreadyAdded: 'Already added', horizontalHint: 'Scroll sideways for more columns', unavailable: 'Check source', cancel: 'Cancel',
  loading: 'Loading sources', failed: 'Cannot load sources', retry: 'Retry', empty: 'No source points', noMatches: 'No matches', clearFilters: 'Clear filters',
  newSource: 'Add source', invalidPending: 'Revalidate pending points', previous: 'Previous', next: 'Next', pageSize: 'Rows per page',
  matching: (n, a) => `${n} matches, ${a} already added`, pending: n => `${n} pending`, hidden: n => `${n} outside this page`, add: n => `Add ${n} points`, page: (n, t) => `Page ${n} of ${t}`, removePending: identity => `Remove unavailable point ${identity.point_id}`,
};
const points: GroupCandidate[] = Array.from({ length: 31 }, (_, i) => ({ key: `p${i}`, point_id: `p${i}`, tag_id: `t${i}`, tag_key: `tag-${i}`, device_id: i < 2 ? 'd1' : 'd2', device_name: i < 2 ? 'PLC 1' : 'PLC 2', label: `Point ${i}`, address: '40001', target_type: 'float32' }));
function props(): RecordingPointPickerProps {
  return { state: createPicker('scope'), projection: { scope: 'scope', candidates: points, members: [points[0]], evidence: Object.fromEntries(points.map(p => [sourceKey(p), { availability: 'eligible', semantic: p.point_id === 'p1' ? 'counter' : 'gauge' }])) }, devices: [{ id: 'd1', name: 'PLC 1' }, { id: 'd2', name: 'PLC 2' }, { id: 'empty', name: 'Empty PLC' }], status: 'loaded', labels, onChange: vi.fn(), onCancel: vi.fn(), onAdd: vi.fn(), onSource: vi.fn(), onRetry: vi.fn() };
}
function Controlled({ value }: { value: RecordingPointPickerProps }) {
  const [state, onChange] = useState(value.state);
  return <RecordingPointPicker {...value} state={state} onChange={onChange} />;
}
describe('unmounted picker content', () => {
  it('locks existing members, stages additions and keeps count/actions together', () => {
    const value = props(); render(<Controlled value={value} />);
    expect(screen.getByRole('button', { name: 'Add 0 points' })).toBeDisabled();
    expect(screen.getByRole('checkbox', { name: 'Point 0' })).toBeDisabled();
    fireEvent.click(screen.getByRole('checkbox', { name: 'Point 1' }));
    const footer = screen.getByRole('contentinfo');
    expect(within(footer).getByText('1 pending')).toBeVisible();
    fireEvent.change(screen.getByLabelText('Device'), { target: { value: 'd2' } });
    expect(screen.queryByRole('checkbox', { name: 'Point 1' })).not.toBeInTheDocument();
    expect(within(footer).getByText('1 outside this page')).toBeVisible();
    fireEvent.click(within(footer).getByRole('button', { name: 'Add 1 points' }));
    expect(value.onAdd).toHaveBeenCalledWith([{ device_id: 'd1', point_id: 'p1', tag_id: 't1' }]);
  });
  it('uses one table for search/device/meaning and keeps zero-point devices filterable', () => {
    const value = props(); render(<Controlled value={value} />);
    fireEvent.change(screen.getByLabelText('Search points'), { target: { value: 'Point 1' } });
    fireEvent.change(screen.getByLabelText('Point type'), { target: { value: 'counter' } });
    expect(screen.getAllByRole('checkbox')).toHaveLength(1);
    expect(screen.getByText('1 matches, 0 already added')).toBeVisible();
    fireEvent.change(screen.getByLabelText('Device'), { target: { value: 'empty' } });
    expect(screen.getByText('No matches')).toBeVisible();
    fireEvent.click(screen.getByRole('button', { name: 'Clear filters' }));
    fireEvent.change(screen.getByLabelText('Device'), { target: { value: 'empty' } });
    expect(screen.getByText('No source points')).toBeVisible();
    fireEvent.click(screen.getByRole('button', { name: 'Add source' }));
    expect(value.onSource).toHaveBeenCalledWith('empty');
  });
  it('cancels/Escapes pending only, without backdrop dismissal or any add callback', () => {
    const value = props(); render(<Controlled value={value} />);
    fireEvent.click(screen.getByRole('checkbox', { name: 'Point 1' }));
    fireEvent.keyDown(screen.getByLabelText('Search points'), { key: 'Escape' });
    expect(value.onCancel).toHaveBeenCalledWith('escape');
    expect(screen.getByRole('checkbox', { name: 'Point 1' })).not.toBeChecked();
    expect(screen.getByRole('checkbox', { name: 'Point 0' })).toBeChecked();
    fireEvent.click(screen.getByRole('button', { name: 'Cancel' }));
    expect(value.onCancel).toHaveBeenCalledWith('cancel');
    fireEvent.click(screen.getByRole('region', { name: 'Select points' }));
    expect(value.onCancel).toHaveBeenCalledTimes(2);
    expect(value.onAdd).not.toHaveBeenCalled();
  });
  it.each(['loading', 'failed'] as const)('does not show %s as empty or permit stale Add', status => {
    const value = props(); value.status = status; value.state.pending = [points[1]];
    render(<RecordingPointPicker {...value} />);
    expect(screen.getByText(status === 'loading' ? labels.loading : labels.failed)).toBeVisible();
    expect(screen.queryByText(labels.empty)).not.toBeInTheDocument();
    expect(screen.getByRole('button', { name: 'Add 1 points' })).toBeDisabled();
    expect(screen.queryAllByRole('checkbox')).toHaveLength(0);
    if (status === 'failed') { fireEvent.click(screen.getByRole('button', { name: 'Retry' })); expect(value.onRetry).toHaveBeenCalledOnce(); }
  });
  it('bounds rows and omits pagination for a short result', () => {
    render(<Controlled value={props()} />);
    expect(screen.getAllByRole('checkbox')).toHaveLength(25);
    fireEvent.click(screen.getByRole('button', { name: 'Next' }));
    expect(screen.getAllByRole('checkbox')).toHaveLength(6);
    fireEvent.change(screen.getByLabelText('Device'), { target: { value: 'd1' } });
    expect(screen.getAllByRole('checkbox')).toHaveLength(2);
    expect(screen.queryByRole('button', { name: 'Next' })).not.toBeInTheDocument();
  });
  it('remembers and restores same-scope picker scroll/focus without creating a modal', () => {
    const value = props(); value.state = { ...value.state, scrollTop: 120, focusKey: 'search' };
    const view = render(<RecordingPointPicker {...value} />);
    expect(screen.getByLabelText('Search points')).toHaveFocus();
    const body = view.container.querySelector('.recording-point-picker__body') as HTMLDivElement;
    expect(body.scrollTop).toBe(120);
    fireEvent.scroll(body, { target: { scrollTop: 220 } });
    expect(value.onChange).toHaveBeenLastCalledWith({ ...value.state, scrollTop: 220 });
    expect(screen.queryByRole('dialog')).not.toBeInTheDocument();
  });

});

describe('picker recovery after projection or result changes', () => {
  it.each(['stale', 'invalid', 'unknown'] as const)('removes a %s pending point without losing other pending or allowing invalid re-add', availability => {
    const value = props(); value.state.pending = [points[1], points[2]];
    value.projection.evidence = { ...value.projection.evidence, [sourceKey(points[1])]: { availability, reason: 'Recheck source' } };
    render(<Controlled value={value} />);
    expect(screen.getByRole('button', { name: 'Add 2 points' })).toBeDisabled();
    expect(within(screen.getByRole('contentinfo')).getByRole('button', { name: 'Remove unavailable point p1' })).toBeEnabled();
    const invalid = screen.getByRole('checkbox', { name: 'Point 1' });
    expect(invalid).toBeEnabled();
    fireEvent.click(invalid);
    expect(invalid).not.toBeChecked();
    expect(invalid).toBeDisabled();
    expect(screen.getByRole('checkbox', { name: 'Point 2' })).toBeChecked();
    expect(screen.getByRole('checkbox', { name: 'Point 0' })).toBeChecked();
    expect(screen.getByRole('checkbox', { name: 'Point 0' })).toBeDisabled();
    expect(screen.getByText('1 pending')).toBeVisible();
    fireEvent.click(screen.getByRole('button', { name: 'Add 1 points' }));
    expect(value.onAdd).toHaveBeenCalledWith([points[2]]);
    expect(value.onCancel).not.toHaveBeenCalled();
  });
  it.each([
    ['Device', 'd2'], ['Search points', 'Point 1'], ['Point type', 'counter'], ['Rows per page', '50'],
  ])('physically resets scroll after changing %s, without stealing current focus', (label, target) => {
    const value = props(); value.state = { ...value.state, scrollTop: 120, focusKey: 'search' };
    const view = render(<Controlled value={value} />);
    const body = view.container.querySelector('.recording-point-picker__body') as HTMLDivElement;
    expect(body.scrollTop).toBe(120);
    fireEvent.scroll(body, { target: { scrollTop: 220 } });
    const control = screen.getByLabelText(label);
    control.focus();
    fireEvent.change(control, { target: { value: target } });
    expect(body.scrollTop).toBe(0);
    if (label !== 'Rows per page') expect(control).toHaveFocus();
  });
  it('physically resets page scroll and retains pending selection', () => {
    const value = props(); value.state.pending = [points[1]];
    const view = render(<Controlled value={value} />);
    const body = view.container.querySelector('.recording-point-picker__body') as HTMLDivElement;
    fireEvent.scroll(body, { target: { scrollTop: 220 } });
    fireEvent.click(screen.getByRole('button', { name: 'Next' }));
    expect(body.scrollTop).toBe(0);
    expect(screen.getByText('1 pending')).toBeVisible();
    expect(screen.getByText('1 outside this page')).toBeVisible();
  });
});


describe('missing pending source recovery', () => {
  it('removes a disappeared identity individually, preserving other pending and draft members', () => {
    const value = props(); value.state.pending = [points[1], points[2]];
    const view = render(<Controlled value={value} />);
    value.projection = { ...value.projection, candidates: points.filter(point => point !== points[1]) };
    view.rerender(<Controlled value={value} />);
    expect(screen.queryByRole('checkbox', { name: 'Point 1' })).not.toBeInTheDocument();
    expect(screen.getByRole('button', { name: 'Add 2 points' })).toBeDisabled();
    fireEvent.click(screen.getByRole('button', { name: 'Remove unavailable point p1' }));
    expect(screen.getByRole('checkbox', { name: 'Point 2' })).toBeChecked();
    expect(screen.getByRole('checkbox', { name: 'Point 0' })).toBeChecked();
    expect(screen.getByText('1 pending')).toBeVisible();
    fireEvent.click(screen.getByRole('button', { name: 'Add 1 points' }));
    expect(value.onAdd).toHaveBeenCalledWith([points[2]]);
    expect(value.onCancel).not.toHaveBeenCalled();
  });
  it('does not mislabel merely filtered valid pending as disappeared', () => {
    const value = props(); value.state.pending = [points[1], points[2]];
    render(<Controlled value={value} />);
    fireEvent.change(screen.getByLabelText('Device'), { target: { value: 'd2' } });
    expect(screen.getByText('1 outside this page')).toBeVisible();
    expect(screen.queryByRole('button', { name: 'Remove unavailable point p1' })).not.toBeInTheDocument();
    expect(screen.getByRole('button', { name: 'Add 2 points' })).toBeEnabled();
  });
});


describe('same-scope membership refresh recovery', () => {
  it('removes pending that becomes an existing member without removing the member or other pending', () => {
    const value = props(); value.state.pending = [points[1], points[2]];
    const view = render(<Controlled value={value} />);
    value.projection = { ...value.projection, members: [points[0], points[1]] };
    view.rerender(<Controlled value={value} />);
    const refreshedMember = screen.getByRole('checkbox', { name: 'Point 1' });
    expect(refreshedMember).toBeDisabled();
    expect(screen.getByRole('button', { name: 'Add 2 points' })).toBeDisabled();
    expect(screen.getByText('2 pending')).toBeVisible();
    fireEvent.click(within(screen.getByRole('contentinfo')).getByRole('button', { name: 'Remove unavailable point p1' }));
    expect(refreshedMember).toBeChecked();
    expect(refreshedMember).toBeDisabled();
    expect(value.projection.members).toEqual([points[0], points[1]]);
    expect(screen.getByRole('checkbox', { name: 'Point 2' })).toBeChecked();
    fireEvent.click(screen.getByRole('button', { name: 'Add 1 points' }));
    expect(value.onAdd).toHaveBeenCalledWith([points[2]]);
    expect(value.onCancel).not.toHaveBeenCalled();
  });
});


it('focus-only rerender does not undo native scroll before the scroll event is delivered', () => {
  const value = props(); const view = render(<Controlled value={value} />);
  const body = view.container.querySelector('.recording-point-picker__body') as HTMLDivElement;
  const source = screen.getByRole('button', { name: 'Add source' });
  // Native scrollIntoView/focus can move the viewport before React receives onScroll.
  // Keep the pointer target stationary through focus capture, not at the old stored offset.
  body.scrollTop = 29;
  act(() => source.focus());
  expect(source).toHaveFocus();
  expect(body.scrollTop).toBe(29);
  fireEvent.click(source);
  expect(value.onSource).toHaveBeenCalledOnce();
  expect(value.onSource).toHaveBeenCalledWith(undefined);
  expect(value.onAdd).not.toHaveBeenCalled();
});

it('keeps 32 existing members distinct from pending additions and shows actual types without repeating unknown meaning', () => {
  const value = props();
  const rows = Array.from({ length: 48 }, (_, i) => ({ ...points[0], key: `P${i}`, point_id: `P${i}`, tag_id: `T${i}`, label: `TAG_${i}`, address: `D${i}` }));
  value.state.pageSize = 50;
  value.projection = { scope: value.state.scope, candidates: rows, members: rows.slice(0, 32), evidence: Object.fromEntries(rows.map(row => [sourceKey(row), { availability: 'eligible' }])) };
  render(<Controlled value={value} />);
  expect(screen.getAllByRole('columnheader')).toHaveLength(4);
  expect(screen.getByRole('columnheader', { name: 'Data type' })).toBeVisible();
  expect(screen.getAllByText(labels.meaningUnknown)).toHaveLength(1);
  expect(screen.queryByText(labels.unknownType)).not.toBeInTheDocument();
  const existing = screen.getByRole('checkbox', { name: 'TAG_29' });
  expect(existing).toBeChecked(); expect(existing).toBeDisabled();
  expect(within(existing.closest('tr')!).queryByText(labels.alreadyAdded)).not.toBeInTheDocument();
  expect(existing).toHaveAccessibleDescription(labels.alreadyAdded);
  expect(existing.closest('label')).toHaveClass('recording-point-picker__point');
  expect(within(existing.closest('tr')!).getByText('float32')).toBeVisible();
  expect(screen.getByRole('button', { name: 'Add 0 points' })).toBeDisabled();
  fireEvent.click(screen.getByRole('checkbox', { name: 'TAG_32' }));
  expect(screen.getByRole('button', { name: 'Add 1 points' })).toBeEnabled();
  fireEvent.click(screen.getByRole('button', { name: 'Cancel' }));
  expect(value.onAdd).not.toHaveBeenCalled();
  expect(screen.getByRole('checkbox', { name: 'TAG_29' })).toBeChecked();
});


it('orders all matching points by address first, then sorts by a chosen column without losing selections', () => {
  const value = props();
  const rows = points.map((point, index) => ({ ...point, address: `D${30 - index}` }));
  value.projection = { ...value.projection, candidates: rows, members: [rows[30]],
    evidence: Object.fromEntries(rows.map(row => [sourceKey(row), { availability: 'eligible' }])) };
  render(<Controlled value={value} />);
  const addresses = () => within(screen.getByRole('table')).getAllByRole('row').slice(1).map(row => within(row).getAllByRole('cell')[2].textContent);
  expect(addresses().slice(0, 4)).toEqual(['D0', 'D1', 'D2', 'D3']);
  expect(screen.getByRole('columnheader', { name: 'Address' })).toHaveAttribute('aria-sort', 'ascending');
  expect(screen.getByRole('checkbox', { name: 'Point 30' })).toBeDisabled();
  fireEvent.click(screen.getByRole('checkbox', { name: 'Point 29' }));
  fireEvent.click(screen.getByRole('button', { name: 'Next' }));
  fireEvent.click(screen.getByRole('button', { name: 'Address' }));
  expect(screen.getByText('Page 1 of 2')).toBeVisible();
  expect(screen.getByRole('columnheader', { name: 'Address' })).toHaveAttribute('aria-sort', 'descending');
  expect(addresses().slice(0, 3)).toEqual(['D30', 'D29', 'D28']);
  expect(screen.getByText('1 pending')).toBeVisible();
  expect(screen.getByText('1 outside this page')).toBeVisible();
  fireEvent.click(screen.getByRole('button', { name: 'Add 1 points' }));
  expect(value.onAdd).toHaveBeenCalledWith([{ device_id: rows[29].device_id, point_id: 'p29', tag_id: 't29' }]);
});


it('sorts MC hexadecimal addresses and each visible column in a human-readable order', () => {
  const value = props();
  const rows: GroupCandidate[] = [
    { ...points[0], label: 'Point 30', address: 'W10', device_name: 'PLC 2', target_type: 'int16' },
    { ...points[1], label: 'Point 3', address: 'WA', device_name: 'PLC 10', target_type: 'float32' },
    { ...points[2], label: 'Point 10', address: 'W9', device_name: 'PLC 1', target_type: 'bool' },
  ];
  value.projection = { ...value.projection, candidates: rows, members: [],
    evidence: Object.fromEntries(rows.map(row => [sourceKey(row), { availability: 'eligible' }])) };
  render(<Controlled value={value} />);
  const cells = (index: number) => within(screen.getByRole('table')).getAllByRole('row').slice(1)
    .map(row => within(row).getAllByRole('cell')[index].textContent?.trim());
  expect(cells(2)).toEqual(['W9', 'WA', 'W10']);
  fireEvent.click(screen.getByRole('button', { name: 'Point' }));
  expect(cells(0)).toEqual(['Point 3', 'Point 10', 'Point 30']);
  fireEvent.click(screen.getByRole('button', { name: 'Device' }));
  expect(cells(1)).toEqual(['PLC 1', 'PLC 2', 'PLC 10']);
  fireEvent.click(screen.getByRole('button', { name: 'Data type' }));
  expect(cells(3)).toEqual(['bool', 'float32', 'int16']);
  expect(value.onAdd).not.toHaveBeenCalled();
});
