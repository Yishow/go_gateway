import { useState } from 'react';
import { afterEach, describe, expect, it, vi } from 'vitest';
import { cleanup, fireEvent, render, screen, within } from '@testing-library/react';
import { RecordingHomeCollection } from '@/features/datalink/workbench-v2/recording-home/RecordingHomeCollection';
import { initialCollectionState, type HomeDeviceSummary, type RecordingCollectionState } from '@/features/datalink/workbench-v2/recording-home/recordingHomeCollections';
import { readHomePresentation, rememberHomePresentation } from '@/features/datalink/workbench-v2/recording-home/recordingHomePresentation';
import { recordingHomeGuidance } from '@/features/datalink/workbench-v2/recording-home/recordingHomeGuidance';
import { fixture, group } from './recording-home-guidance-fixtures';
vi.mock('react-i18next', () => ({ useTranslation: () => ({ t: (key: string, value?: Record<string, unknown>) => key === 'pageOf' ? `${value?.page}/${value?.pages} (${value?.count})` : key === 'recordPointCount' ? `${value?.count} recording points` : key }) }));
afterEach(cleanup);
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
  it('300 records render only25/50 rows, devices summary only4, and same-page device management stays bounded', () => {
    render(<Host />); expect(within(screen.getByTestId('home-management-table')).getAllByRole('row')).toHaveLength(26);
    expect(within(screen.getByTestId('home-device-summary')).getAllByRole('listitem')).toHaveLength(4);
    fireEvent.change(screen.getByLabelText('pageSize'), { target: { value: '50' } });
    expect(within(screen.getByTestId('home-management-table')).getAllByRole('row')).toHaveLength(51);
    fireEvent.click(screen.getByText('next')); expect(screen.getByTestId('home-page-status')).toHaveTextContent('2/6 (300)');
    fireEvent.click(screen.getByRole('button', { name: 'devicesTab' }));
    expect(within(screen.getByTestId('home-management-table')).getAllByRole('row')).toHaveLength(51);
    expect(screen.getByTestId('home-page-status')).toHaveTextContent('1/2 (100)');
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
  it('does not resurrect a removed focus target or describe stale observations as zero faults', () => {
    const state = initialCollectionState(); render(<RecordingHomeCollection groups={groups.slice(0, 1)} devices={devices} state={{ ...state, observedOnly: true }} guidance={guidance}
      observedIds={[]} unverifiedObserved sourcesAvailable restored={{ state, workspaceRevision: 'old', scrollTop: 80, focusId: 'record-removed' }}
      onChange={vi.fn()} onOpen={vi.fn()} onOpenDevice={vi.fn()} onRemember={vi.fn()} onRecheck={vi.fn()} />);
    expect(screen.getByText('observationNeedsRefresh')).toBeInTheDocument(); expect(screen.getByRole('heading', { name: 'saved' })).toHaveFocus();
    expect(screen.getByText('recheckVisible')).toBeInTheDocument(); expect(screen.queryByText('removed')).not.toBeInTheDocument();
  });
});


describe('demo device summary stays inside the simple home', () => {
  it('places saved recordings before the device summary and uses an explicit points action', () => {
    const openDevice = vi.fn(); render(<Host openDevice={openDevice} />);
    const table = screen.getByTestId('home-management-table');
    const summary = screen.getByTestId('home-device-summary');
    expect(table.compareDocumentPosition(summary) & Node.DOCUMENT_POSITION_FOLLOWING).toBeTruthy();
    const row = within(summary).getAllByRole('listitem')[0];
    expect(within(row).queryByRole('button', { name: 'Device 0' })).not.toBeInTheDocument();
    expect(within(row).getByText('Device 0').tagName).toBe('STRONG');
    fireEvent.click(within(row).getByRole('button', { name: 'devicePoints' }));
    expect(openDevice).toHaveBeenCalledExactlyOnceWith('D0', 'points');
  });
  it('offers distinct connection and points actions on the same-page device list', () => {
    const openDevice = vi.fn(); render(<Host openDevice={openDevice} />);
    fireEvent.click(screen.getByRole('button', { name: 'devicesTab' }));
    const row = within(screen.getByTestId('home-management-table')).getAllByRole('row')[1];
    expect(within(row).queryByRole('button', { name: 'Device 0' })).not.toBeInTheDocument();
    fireEvent.click(within(row).getByRole('button', { name: 'deviceConnection' }));
    expect(openDevice).toHaveBeenLastCalledWith('D0', 'connection');
    fireEvent.click(within(row).getByRole('button', { name: 'devicePoints' }));
    expect(openDevice).toHaveBeenLastCalledWith('D0', 'points');
    expect(openDevice).toHaveBeenCalledTimes(2);
    expect(screen.getByTestId('home-page-status')).toHaveTextContent('1/4 (100)');
  });
});

it('separates saved sources, destination and unknown writing facts without inventing runtime health', () => {
  render(<Host />); const table = screen.getByTestId('home-management-table');
  expect(within(table).getAllByRole('columnheader').map(node => node.textContent)).toEqual(['name', 'savedStatus', 'recordSources', 'destinationFilter', 'recordWrites', 'actions']);
  const cells = within(within(table).getAllByRole('row')[1]).getAllByRole('cell');
  expect(cells[2]).toHaveTextContent('Device 0'); expect(cells[2]).toHaveTextContent('1 recording points');
  expect(cells[3]).toHaveTextContent('table_0'); expect(cells[4]).toHaveTextContent('statusUnknown');
  expect(within(cells[5]).getByRole('button', { name: 'openRecord' })).toBeInTheDocument();
  expect(cells[4]).not.toHaveTextContent(/healthy|rows\/s|0 faults/i);
});

it('keeps recording identity plain and exposes exactly one open action per row', () => {
 render(<Host />); const row = within(screen.getByTestId('home-management-table')).getAllByRole('row')[1];
 expect(within(row).getByText('Recording 000').tagName).not.toBe('BUTTON');
 expect(within(row).queryByRole('button', { name: 'Recording 000' })).not.toBeInTheDocument();
 expect(within(row).getAllByRole('button')).toHaveLength(1);
 expect(within(row).getByRole('button', { name: 'openRecord' })).toHaveAttribute('data-home-focus', 'record-G000');
 expect(within(row).getByText('destinationIdentity').closest('details')).toBeInTheDocument();
});

it('keeps clear and observed controls together without dropping unknown source state', () => {
 render(<Host />); const clear = screen.getByRole('button', { name: 'clearFilters' });
 expect(clear.parentElement).toBe(screen.getByLabelText('observedFilter').parentElement?.parentElement);
 const row = within(screen.getByTestId('home-management-table')).getAllByRole('row')[1];
 expect(within(row).getAllByRole('cell')[2]).toHaveTextContent('result.current.unknown');
 expect(within(row).getAllByRole('cell')[1]).not.toHaveTextContent('statusUnknown');
});

describe('Home separates recording and equipment intent', () => {
  it('keeps a single device management surface with both exact actions beside each device', () => {
    const openDevice = vi.fn(); render(<Host openDevice={openDevice} />);
    fireEvent.click(screen.getByRole('button', { name: 'devicesTab' }));
    expect(screen.queryByTestId('home-device-summary')).not.toBeInTheDocument();
    expect(screen.getByText('devicesPurpose')).toBeVisible();
    const row = screen.getByText('Device 0').closest('tr'); if (!row) throw new Error('missing row');
    fireEvent.click(within(row).getByRole('button', { name: 'deviceConnection' }));
    fireEvent.click(within(row).getByRole('button', { name: 'devicePoints' }));
    expect(openDevice.mock.calls).toEqual([['D0', 'connection'], ['D0', 'points']]);
    fireEvent.click(screen.getByRole('button', { name: 'recordingsTab' }));
    expect(screen.getByTestId('home-management-table')).toHaveAttribute('data-mode', 'recordings');
  });
});
