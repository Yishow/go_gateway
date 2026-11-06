import { cleanup, fireEvent, render, screen } from '@testing-library/react';
import { afterEach, beforeEach, expect, it, vi } from 'vitest';
import i18n from '@/i18n/config';
import { RecordingHomeMemoryView } from '@/features/datalink/workbench-v2/recording-home/RecordingHomeMemoryView';
beforeEach(async () => { await i18n.changeLanguage('en'); });
afterEach(cleanup);
function mount(type: 'int64' | 'string' = 'int64') {
  render(<RecordingHomeMemoryView protocol="modbus_tcp" status="ready" checking={false}
    inputs={[{ id: 'A', label: 'Accumulator', start_address: '40001', count: 1, data_type: type, skipped_addresses: [] }]}
    values={{ status: 'live', byCell: new Map([[`HR:1:${type}`, { value: '9223372036854775807', stale: true }]]) }}
    onOpen={vi.fn()} onAdd={vi.fn()} onBack={vi.fn()} onClose={vi.fn()} />);
  fireEvent.click(screen.getByTestId('home-memory-layout-detail'));
  return screen.getByTestId('home-memory-detail');
}
it('groups type and footprint apart from reading and freshness without nested fact cards', () => {
  const detail = mount();
  const type = detail.querySelector('[data-detail-group="type"]');
  const reading = detail.querySelector('[data-detail-group="reading"]');
  expect(type).toHaveTextContent('int64'); expect(type).toHaveTextContent('4 registers');
  expect(type).toHaveTextContent('40001');
  expect(reading).toHaveTextContent('9223372036854775807'); expect(reading).toHaveTextContent('Stale reading');
  expect(detail.querySelector('[data-detail-group="actions"]')).toContainElement(screen.getByTestId('home-memory-detail-action'));
});
it('never subscribes to or exposes a reading for an unconfirmed partial string width', () => {
  const detail = mount('string');
  expect(detail).toHaveTextContent('Unconfirmed');
  expect(detail).not.toHaveTextContent('9223372036854775807');
  expect(detail.querySelector('[data-value-key]')).toBeNull();
});
