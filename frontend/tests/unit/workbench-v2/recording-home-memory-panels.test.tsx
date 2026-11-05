import { cleanup, fireEvent, render, screen } from '@testing-library/react';
import { afterEach, beforeEach, expect, it, vi } from 'vitest';
import i18n from '@/i18n/config';
import { RecordingHomeMemoryView } from '@/features/datalink/workbench-v2/recording-home/RecordingHomeMemoryView';
import type { MemoryRuleInput } from '@/features/datalink/workbench-v2/recording-home/recordingHomeMemoryMap';
beforeEach(async () => { await i18n.changeLanguage('en'); });
afterEach(cleanup);
const inputs: MemoryRuleInput[] = [
  { id: 'I', label: 'Pressure', start_address: '40001', count: 2, data_type: 'int16', skipped_addresses: [] },
  { id: 'F', label: 'Flow', start_address: '40301', count: 2, data_type: 'float32', skipped_addresses: [] },
];
function view(rules = inputs) {
  const onAdd = vi.fn(), onOpen = vi.fn();
  render(<RecordingHomeMemoryView protocol="modbus_tcp" inputs={rules} status="ready" checking={false}
    values={{ status: 'live', byCell: new Map([['HR:301:float32', { value: '12.5', readAt: '2026-11-04T12:00:00Z' }]]) }}
    onAdd={onAdd} onOpen={onOpen} onBack={vi.fn()} onClose={vi.fn()} />);
  return { onAdd, onOpen };
}
it('range navigation and chapters jump the shared window without opening a rule', () => {
  const actions = view();
  fireEvent.click(screen.getByRole('button', { name: 'Range navigation' }));
  expect(screen.getByTestId('home-memory-range-nav')).toHaveTextContent('Pressure');
  fireEvent.click(screen.getByTestId('home-memory-range-jump-F'));
  expect(screen.getByTestId('home-memory-cell-40301')).toBeInTheDocument();
  fireEvent.click(screen.getByRole('button', { name: 'Range chapters' }));
  expect(screen.getByTestId('home-memory-chapter-I')).toHaveTextContent('Pressure');
  fireEvent.click(screen.getByTestId('home-memory-range-jump-I'));
  expect(screen.getByTestId('home-memory-cell-40001')).toBeInTheDocument();
  expect(actions.onOpen).not.toHaveBeenCalled(); expect(actions.onAdd).not.toHaveBeenCalled();
});
it('horizontal map selects without invoking an action and shares selection with detail', () => {
  const { onAdd, onOpen } = view();
  fireEvent.click(screen.getByRole('button', { name: 'Horizontal map' }));
  fireEvent.change(screen.getByTestId('home-memory-map-select'), { target: { value: '3' } });
  expect(screen.getByTestId('home-memory-map')).toHaveTextContent('40004');
  expect(screen.getByTestId('home-memory-cell-40004')).toHaveAttribute('data-selected', 'true');
  fireEvent.click(screen.getByRole('button', { name: 'Grid with detail' }));
  expect(screen.getByTestId('home-memory-detail')).toHaveTextContent('40004');
  expect(onAdd).not.toHaveBeenCalled(); expect(onOpen).not.toHaveBeenCalled();
  fireEvent.click(screen.getByTestId('home-memory-detail-action'));
  expect(onAdd).toHaveBeenCalledExactlyOnceWith('40004');
});
it('detail resolves a continuation reading from the correct starting register', () => {
  const { onOpen } = view();
  fireEvent.change(screen.getByLabelText('Go to address'), { target: { value: '40302' } });
  fireEvent.click(screen.getByRole('button', { name: 'Go' }));
  fireEvent.click(screen.getByRole('button', { name: 'Grid with detail' }));
  const detail = screen.getByTestId('home-memory-detail');
  expect(detail).toHaveTextContent('40302'); expect(detail).toHaveTextContent('40301'); expect(detail).toHaveTextContent('12.5');
  fireEvent.click(screen.getByTestId('home-memory-detail-action'));
  expect(onOpen).toHaveBeenCalledExactlyOnceWith('F');
});
it('detail cannot create unknown cells', () => {
  const { onAdd, onOpen } = view([{ ...inputs[0], data_type: 'string' }]);
  fireEvent.click(screen.getByRole('button', { name: 'Horizontal map' }));
  fireEvent.change(screen.getByTestId('home-memory-map-select'), { target: { value: '3' } });
  fireEvent.click(screen.getByRole('button', { name: 'Grid with detail' }));
  expect(screen.queryByTestId('home-memory-detail-action')).not.toBeInTheDocument();
  expect(onAdd).not.toHaveBeenCalled(); expect(onOpen).not.toHaveBeenCalled();
});
it('detail reports read failure instead of presenting an old value as current', () => {
  render(<RecordingHomeMemoryView protocol="modbus_tcp" inputs={inputs} status="ready" checking={false}
    values={{ status: 'live', byCell: new Map([['HR:1:int16', { value: '99', error: 'untrusted backend detail' }]]) }}
    onAdd={vi.fn()} onOpen={vi.fn()} onBack={vi.fn()} onClose={vi.fn()} />);
  fireEvent.click(screen.getByRole('button', { name: 'Grid with detail' }));
  const detail = screen.getByTestId('home-memory-detail');
  expect(detail.querySelector('.gw-memory-detail-reading')).toHaveTextContent('!');
  expect(detail).toHaveTextContent('Read failed'); expect(detail).not.toHaveTextContent('untrusted backend detail');
});
it('detail selection buttons remain in bounds and conflict opens the owner chooser', () => {
  const { onOpen, onAdd } = view([inputs[0], { ...inputs[0], id: 'Overlap', label: 'Overlap' }]);
  fireEvent.click(screen.getByRole('button', { name: 'Grid with detail' }));
  expect(screen.getByTestId('home-memory-detail-previous')).toBeDisabled();
  fireEvent.click(screen.getByTestId('home-memory-detail-next'));
  expect(screen.getByTestId('home-memory-detail')).toHaveTextContent('40002');
  fireEvent.click(screen.getByTestId('home-memory-detail-action'));
  expect(screen.getByTestId('home-memory-owners')).toHaveTextContent('Overlap');
  expect(onOpen).not.toHaveBeenCalled(); expect(onAdd).not.toHaveBeenCalled();
});
