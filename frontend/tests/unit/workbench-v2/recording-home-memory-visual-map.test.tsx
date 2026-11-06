import { cleanup, fireEvent, render, screen } from '@testing-library/react';
import { afterEach, beforeEach, expect, it, vi } from 'vitest';
import i18n from '@/i18n/config';
import { RecordingHomeMemoryView } from '@/features/datalink/workbench-v2/recording-home/RecordingHomeMemoryView';
beforeEach(async () => { await i18n.changeLanguage('en'); });
afterEach(cleanup);
function mount(unknown = false) {
  const onAdd = vi.fn(), onOpen = vi.fn();
  render(<RecordingHomeMemoryView protocol="mc_3e" status="ready" checking={false}
    inputs={[{ id: 'W', label: 'Link registers', start_address: 'WFE', count: 1, data_type: unknown ? 'string' : 'int64', skipped_addresses: [], proven: unknown ? undefined : ['WFE'] }]}
    values={{ status: 'live', byCell: new Map([['W:254:int64', { value: '9223372036854775807', stale: true }]]) }}
    onOpen={onOpen} onAdd={onAdd} onBack={vi.fn()} onClose={vi.fn()} />);
  fireEvent.click(screen.getByTestId('home-memory-layout-map'));
  fireEvent.change(screen.getByRole('textbox'), { target: { value: 'W0' } });
  fireEvent.click(screen.getByRole('button', { name: 'Go' }));
  return { onAdd, onOpen };
}
it('renders a distinct compact summary with the complete continuation reading and safe owner action', () => {
  const { onOpen } = mount();
  fireEvent.change(screen.getByTestId('home-memory-map-select'), { target: { value: '255' } });
  const summary = screen.getByTestId('home-memory-summary');
  expect(screen.queryByTestId('home-memory-detail')).not.toBeInTheDocument();
  expect(summary.querySelector('dl')).toBeNull();
  expect(summary).toHaveAttribute('data-memory-inspection');
  expect(summary).toHaveTextContent('WFF'); expect(summary).toHaveTextContent('WFE');
  expect(summary).toHaveTextContent('int64'); expect(summary).toHaveTextContent('4 registers');
  expect(summary).toHaveTextContent('9223372036854775807'); expect(summary).toHaveTextContent('Stale reading');
  expect(summary.querySelector('[data-value-key]')).toHaveAttribute('data-value-key', 'W:254:int64');
  fireEvent.click(screen.getByTestId('home-memory-detail-action'));
  expect(onOpen).toHaveBeenCalledExactlyOnceWith('W');
});
it('keeps all 256 canonical slots and exact first, interior and last slider positions', () => {
  const { onAdd, onOpen } = mount();
  expect(document.querySelectorAll('.gw-memory-map-strip > span')).toHaveLength(256);
  const slider = screen.getByTestId('home-memory-map-select');
  for (const [index, address] of [[0, 'W0'], [127, 'W7F'], [255, 'WFF']] as const) {
    fireEvent.change(slider, { target: { value: String(index) } });
    expect(slider).toHaveAttribute('aria-valuetext', address);
    expect(screen.getByTestId(`home-memory-cell-${address}`)).toHaveAttribute('data-selected', 'true');
  }
  expect(onAdd).not.toHaveBeenCalled(); expect(onOpen).not.toHaveBeenCalled();
});
it('does not offer a summary action for unconfirmed cells', () => {
  mount(true);
  expect(screen.getByTestId('home-memory-summary')).toHaveTextContent('Unconfirmed');
  expect(screen.queryByTestId('home-memory-detail-action')).not.toBeInTheDocument();
});
