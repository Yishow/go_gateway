import { cleanup, fireEvent, render, screen } from '@testing-library/react';
import { afterEach, beforeEach, expect, it, vi } from 'vitest';
import i18n from '@/i18n/config';
import { RecordingHomeMemoryView } from '@/features/datalink/workbench-v2/recording-home/RecordingHomeMemoryView';
import type { MemoryLiveValue, MemoryValuesStatus } from '@/features/datalink/workbench-v2/recording-home/recordingHomeMemoryValues';
beforeEach(async () => { await i18n.changeLanguage('en'); });
afterEach(cleanup);
function mount(live?: MemoryLiveValue, status: MemoryValuesStatus = 'live') {
  render(<RecordingHomeMemoryView protocol="modbus_tcp" status="ready" checking={false}
    inputs={[{ id: 'wide', label: 'Wide value', start_address: '40009', count: 1, data_type: 'int64', skipped_addresses: [] }]}
    values={{ status, byCell: new Map(live ? [['HR:9:int64', live]] : []) }}
    onAdd={vi.fn()} onOpen={vi.fn()} onBack={vi.fn()} onClose={vi.fn()} />);
  return screen.getByTestId('home-memory-cell-40009');
}
it('provides complete focused values including continuations without another tab stop', () => {
  const full = '9223372036854775807';
  mount({ value: full });
  for (const address of ['40009', '40010', '40011', '40012']) {
    const cell = screen.getByTestId(`home-memory-cell-${address}`);
    fireEvent.focus(cell);
    const disclosure = cell.querySelector('.gw-memory-focus-information');
    expect(disclosure).not.toBeNull();
    expect(disclosure).toHaveTextContent(full);
    expect(disclosure).toHaveTextContent('int64');
    expect(disclosure).toHaveTextContent(address);
    expect(disclosure).not.toHaveAttribute('tabindex');
    expect(disclosure).toHaveAttribute('data-value-key', 'HR:9:int64');
    expect(cell).toHaveAttribute('tabindex', '0');
    expect(document.querySelectorAll('.gw-memory-grid [tabindex="0"]')).toHaveLength(1);
  }
});
it.each([
  [{ value: '0' }, 'live', 'Stored reading'],
  [{ value: '17', stale: true }, 'live', 'Stale reading'],
  [{ value: '17', error: 'private transport details' }, 'live', 'Read failed'],
  [undefined, 'loading', 'Waiting for reading'],
  [undefined, 'live', 'No reading'],
] as const)('announces truthful reading freshness for %j', (live, status, label) => {
  const cell = mount(live, status);
  expect(cell.getAttribute('aria-label')).toContain(label);
  expect(cell.getAttribute('aria-label')).not.toContain('private transport');
  if (live && 'error' in live) expect(cell.getAttribute('aria-label')).not.toContain('17');
});

it('discloses a continuation reading when its original start is outside the page', () => {
  render(<RecordingHomeMemoryView protocol="modbus_tcp" status="ready" checking={false}
    inputs={[{ id: 'edge', label: 'Page crossing', start_address: '40249', count: 1, data_type: 'int64', skipped_addresses: [] }]}
    values={{ status: 'live', byCell: new Map([['HR:249:int64', { value: '9223372036854775807', stale: true }]]) }}
    onAdd={vi.fn()} onOpen={vi.fn()} onBack={vi.fn()} onClose={vi.fn()} />);
  fireEvent.change(screen.getByRole('textbox'), { target: { value: '40001' } });
  fireEvent.click(screen.getByRole('button', { name: 'Go' }));
  fireEvent.click(screen.getByRole('button', { name: 'Next' }));
  expect(screen.queryByTestId('home-memory-cell-40249')).not.toBeInTheDocument();
  const continuation = screen.getByTestId('home-memory-cell-40251');
  fireEvent.focus(continuation);
  const information = continuation.querySelector('.gw-memory-focus-information');
  expect(information).toHaveTextContent('40249');
  expect(information).toHaveTextContent('9223372036854775807');
  expect(information).toHaveTextContent('Stale reading');
  expect(continuation).toHaveAttribute('tabindex', '0');
});
