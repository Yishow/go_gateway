import { useState } from 'react';
import { afterEach, describe, expect, it, vi } from 'vitest';
import { cleanup, fireEvent, render, screen, within } from '@testing-library/react';
import { RecordingHomeCollection } from '@/features/datalink/workbench-v2/recording-home/RecordingHomeCollection';
import { initialCollectionState, type HomeDeviceSummary, type RecordingCollectionState } from '@/features/datalink/workbench-v2/recording-home/recordingHomeCollections';
import { readHomePresentation, rememberHomePresentation } from '@/features/datalink/workbench-v2/recording-home/recordingHomePresentation';
import { recordingHomeGuidance } from '@/features/datalink/workbench-v2/recording-home/recordingHomeGuidance';
import { fixture, group } from './recording-home-guidance-fixtures';
vi.mock('react-i18next', () => ({ useTranslation: () => ({ t: (key: string, value?: Record<string, unknown>) => key === 'pageOf' ? `${value?.page}/${value?.pages} (${value?.count})` : key === 'recordPointCount' ? `${value?.count} recording points` : key }) }));
afterEach(() => { cleanup(); vi.restoreAllMocks(); vi.unstubAllGlobals(); });
const devices: HomeDeviceSummary[] = Array.from({ length: 100 }, (_, index) => ({ id: `D${index}`, name: `Device ${index}`, protocol: 'mc_3e', availability: 'available' }));
const groups = Array.from({ length: 300 }, (_, index) => ({ ...group(`G${String(index).padStart(3, '0')}`, index % 3 === 0 ? 'draft' : 'running'), name: `Recording ${String(index).padStart(3, '0')}`,
  members: [{ ...group('G').members[0], device_id: `D${index % 100}` }], destination: { ...group('G').destination, table_name: `table_${index}` } }));
const guidance = recordingHomeGuidance(fixture(devices.map(device => device.id), false, groups));
function Host({ initial = initialCollectionState(), observedIds = [], openDevice = vi.fn() }: { initial?: RecordingCollectionState; observedIds?: string[]; openDevice?: (id: string, kind: 'connection' | 'points') => void }) {
  const [state, setState] = useState(initial);
  return <RecordingHomeCollection groups={groups} devices={devices} state={state} guidance={guidance} observedIds={observedIds} unverifiedObserved={false}
    sourcesAvailable restored={{ state: initial, workspaceRevision: 'W2', scrollTop: 0, focusId: null }}
    onChange={setState} onOpen={vi.fn()} onOpenDevice={openDevice} onRemember={vi.fn()} onRecheck={vi.fn()} />;
}
describe('compact home management', () => {
  it('300 records render only 25/50 rows and the panel never switches to a device table', () => {
    render(<Host />); expect(within(screen.getByTestId('home-management-table')).getAllByRole('row')).toHaveLength(26);
    expect(screen.queryByTestId('home-device-summary')).not.toBeInTheDocument();
    fireEvent.change(screen.getByLabelText('pageSize'), { target: { value: '50' } });
    expect(within(screen.getByTestId('home-management-table')).getAllByRole('row')).toHaveLength(51);
    fireEvent.click(screen.getByText('next')); expect(screen.getByTestId('home-page-status')).toHaveTextContent('2/6 (300)');
    expect(screen.queryByRole('button', { name: 'devicesTab' })).not.toBeInTheDocument();
    expect(screen.getByTestId('home-management-table')).toHaveAttribute('data-mode', 'recordings');
  });
  it('keyword and filters find a named recording without scanning pages, with clear-filter recovery', () => {
    render(<Host />); fireEvent.change(screen.getByLabelText('keyword'), { target: { value: 'table_299' } });
    expect(screen.getByText('Recording 299')).toBeInTheDocument();
    expect(screen.getByTestId('home-page-status')).toHaveTextContent('1/1 (1)');
    fireEvent.change(screen.getByLabelText('savedStatus'), { target: { value: 'draft' } });
    expect(screen.getByText('noResults')).toBeInTheDocument(); expect(screen.queryByText('emptyCollection')).not.toBeInTheDocument();
    fireEvent.click(screen.getByText('clearFilters')); expect(screen.getByTestId('home-page-status')).toHaveTextContent('1/12 (300)');
  });
  it('observed-only filter cannot turn unobserved records into a health claim', () => {
    render(<Host observedIds={['G299']} />); fireEvent.click(screen.getByLabelText('observedFilter'));
    expect(screen.getByText('Recording 299')).toBeInTheDocument();
    expect(screen.getByText('observedFilterBoundary')).toBeInTheDocument();
    expect(screen.getByText('observationBoundary')).toBeInTheDocument();
  });
  it('keeps observed issues available outside an unrelated keyword filter and clears that filter to manage them', () => {
    const g = group('G0'); const input = fixture(['D1'], true, [g]); const item = input.evidence.delivery?.G0;
    if (item?.state !== 'available') throw new Error('fixture');
    const fault = recordingHomeGuidance({ ...input, evidence: { ...input.evidence, delivery: { G0: { ...item, value: { ...item.value, intake: { state: 'blocked' } } } } } });
    const changed = vi.fn(); render(<RecordingHomeCollection groups={[g]} devices={devices} state={{ ...initialCollectionState(), keyword: 'absent' }} guidance={fault}
      observedIds={['G0']} unverifiedObserved={false} sourcesAvailable restored={{ state: initialCollectionState(), workspaceRevision: 'W2', scrollTop: 0, focusId: null }}
      onChange={changed} onOpen={vi.fn()} onOpenDevice={vi.fn()} onRemember={vi.fn()} onRecheck={vi.fn()} />);
    expect(screen.getByText('Recording G0')).toBeInTheDocument();
    fireEvent.click(screen.getByText('viewObserved')); expect(changed).toHaveBeenCalledWith(expect.objectContaining({ keyword: '', observedOnly: true, page: 1 }));
  });
  it('restores filters only within the same query context/workspace and keeps IDs as presentation only', () => {
    const context = {}; const saved = { state: { ...initialCollectionState(), keyword: 'needle', page: 9, observedOnly: true }, workspaceRevision: 'old', scrollTop: 200, focusId: 'record-removed', observedIds: ['removed'] };
    rememberHomePresentation(context, 'W1', saved);
    expect(readHomePresentation(context, 'W1', 'new').state.keyword).toBe('needle');
    expect(readHomePresentation(context, 'W2', 'new').state.keyword).toBe('');
    expect(readHomePresentation({}, 'W1', 'new').state.keyword).toBe('');
  });
  it('clamps a removed final page and restores scroll/focus only to an existing rendered identity', () => {
    const state = { ...initialCollectionState(), page: 12 }; const restored = { state, workspaceRevision: 'old', scrollTop: 80, focusId: 'record-G026' };
    const props = { groups: groups.slice(0, 28), devices, state, guidance, observedIds: [], unverifiedObserved: false, restored,
      onChange: vi.fn(), onOpen: vi.fn(), onOpenDevice: vi.fn(), onRemember: vi.fn(), onRecheck: vi.fn() };
    const view = render(<RecordingHomeCollection {...props} sourcesAvailable={false} />);
    const panel = screen.getByTestId('home-collection-scroll'); Object.defineProperty(panel, 'scrollHeight', { value: 500 }); Object.defineProperty(panel, 'clientHeight', { value: 100 });
    view.rerender(<RecordingHomeCollection {...props} sourcesAvailable />);
    expect(screen.getByTestId('home-page-status')).toHaveTextContent('2/2 (28)');
    expect(within(screen.getByText('Recording 026').closest('tr')!).getByRole('button', { name: 'openRecord' })).toHaveFocus(); expect(panel.scrollTop).toBe(80);
  });
  it('retains restored focus when the browser emits its restoration scroll event', () => {
    const state = { ...initialCollectionState(), page: 2 }, onRemember = vi.fn();
    const restored = { state, workspaceRevision: 'W2', scrollTop: 80, focusId: 'record-G026' };
    const props = { groups, devices, state, guidance, observedIds: [], unverifiedObserved: false, restored,
      onChange: vi.fn(), onOpen: vi.fn(), onOpenDevice: vi.fn(), onRemember, onRecheck: vi.fn() };
    const view = render(<RecordingHomeCollection {...props} sourcesAvailable={false} />);
    const panel = screen.getByTestId('home-collection-scroll');
    Object.defineProperty(panel, 'scrollHeight', { value: 500 }); Object.defineProperty(panel, 'clientHeight', { value: 100 });
    view.rerender(<RecordingHomeCollection {...props} sourcesAvailable />);
    const trigger = document.querySelector('[data-home-focus="record-G026"]');
    expect(trigger).toHaveFocus(); expect(panel.scrollTop).toBe(80);
    // Browsers dispatch scroll after the layout-effect scrollTop assignment and focus restoration.
    fireEvent.scroll(panel);
    expect(onRemember).toHaveBeenLastCalledWith(80, 'record-G026', 2);
    expect(trigger).toHaveFocus();
    // A user scroll after moving focus away must not resurrect the old recording focus.
    screen.getByLabelText('keyword').focus(); panel.scrollTop = 120; fireEvent.scroll(panel);
    expect(onRemember).toHaveBeenLastCalledWith(120, null, 2);
    expect(screen.getByLabelText('keyword')).toHaveFocus();
  });
  it.each([[783, 115], [-52, -60], [100, 0]])('reveals restored focus at viewport y=%i by moving only the document', (top, movement) => {
    vi.stubGlobal('innerHeight', 720);
    const scroll = vi.spyOn(window, 'scrollBy').mockImplementation((...args: unknown[]) => { top -= (args[0] as { top?: number }).top ?? 0; });
    const state = { ...initialCollectionState(), page: 2 };
    const props = { groups, devices, state, guidance, observedIds: [], unverifiedObserved: false,
      restored: { state, workspaceRevision: 'W2', scrollTop: 80, focusId: 'record-G026' },
      onChange: vi.fn(), onOpen: vi.fn(), onOpenDevice: vi.fn(), onRemember: vi.fn(), onRecheck: vi.fn() };
    const view = render(<RecordingHomeCollection {...props} sourcesAvailable={false} />);
    const panel = screen.getByTestId('home-collection-scroll');
    Object.defineProperty(panel, 'scrollHeight', { value: 500 }); Object.defineProperty(panel, 'clientHeight', { value: 100 });
    const trigger = document.querySelector<HTMLElement>('[data-home-focus="record-G026"]')!;
    vi.spyOn(trigger, 'getBoundingClientRect').mockImplementation(() => ({ top, bottom: top + 44, left: 0, right: 80,
      width: 80, height: 44, x: 0, y: top, toJSON: () => ({}) }));
    view.rerender(<RecordingHomeCollection {...props} sourcesAvailable />);
    expect(trigger).toHaveFocus(); expect(panel.scrollTop).toBe(80);
    if (movement) expect(scroll).toHaveBeenCalledExactlyOnceWith({ top: movement, behavior: 'instant' });
    else expect(scroll).not.toHaveBeenCalled();
  });
  it('does not resurrect a removed focus target or describe stale observations as zero faults', () => {
    const state = initialCollectionState(); render(<RecordingHomeCollection groups={groups.slice(0, 1)} devices={devices} state={{ ...state, observedOnly: true }} guidance={guidance}
      observedIds={[]} unverifiedObserved sourcesAvailable restored={{ state, workspaceRevision: 'old', scrollTop: 80, focusId: 'record-removed' }}
      onChange={vi.fn()} onOpen={vi.fn()} onOpenDevice={vi.fn()} onRemember={vi.fn()} onRecheck={vi.fn()} />);
    expect(screen.getByText('observationNeedsRefresh')).toBeInTheDocument(); expect(screen.getByRole('heading', { name: 'saved' })).toHaveFocus();
    expect(screen.getByText('recheckVisible')).toBeInTheDocument(); expect(screen.queryByText('removed')).not.toBeInTheDocument();
  });
});


it('separates saved sources, destination and unknown writing facts without inventing runtime health', () => {
  render(<Host />); const table = screen.getByTestId('home-management-table');
  expect(within(table).getAllByRole('columnheader').map(node => node.textContent)).toEqual(['name', 'recordSources', 'destinationFilter', 'recordWrites', 'actions']);
  const cells = within(within(table).getAllByRole('row')[1]).getAllByRole('cell');
  expect(cells[0]).toHaveTextContent('Recording 000'); expect(cells[0].querySelector('[data-status-badge]')).toBeTruthy();
  expect(cells[1]).toHaveTextContent('Device 0'); expect(cells[1]).toHaveTextContent('1 recording points');
  expect(cells[2]).toHaveTextContent('table_0'); expect(cells[3]).toHaveTextContent('statusUnknown');
  expect(within(cells[4]).getByRole('button', { name: 'openRecord' })).toBeInTheDocument();
  expect(cells[3]).not.toHaveTextContent(/healthy|rows\/s|0 faults/i);
});

it('keeps recording identity plain, with open and copy; a draft without delete authority shows no extra actions', () => {
 render(<Host />); const row = within(screen.getByTestId('home-management-table')).getAllByRole('row')[1];
 expect(within(row).getByText('Recording 000').tagName).not.toBe('BUTTON');
 expect(within(row).queryByRole('button', { name: 'Recording 000' })).not.toBeInTheDocument();
 expect(within(row).getAllByRole('button')).toHaveLength(2);
 expect(within(row).getByRole('button', { name: 'openRecord' })).toHaveAttribute('data-home-focus', 'record-G000');
 expect(row.querySelectorAll('details')).toHaveLength(0);
 expect(within(row).getByRole('button', { name: 'rowEvidence.copyTable' })).toBeInTheDocument();
});

it('keeps clear and observed controls together without dropping unknown source state', () => {
 render(<Host />); const clear = screen.getByRole('button', { name: 'clearFilters' });
 expect(clear.parentElement).toBe(screen.getByLabelText('observedFilter').parentElement?.parentElement);
 const row = within(screen.getByTestId('home-management-table')).getAllByRole('row')[1];
 expect(within(row).getAllByRole('cell')[1]).toHaveTextContent('result.current.unknown');
 expect(within(row).getAllByRole('cell')[0]).not.toHaveTextContent('statusUnknown');
});


