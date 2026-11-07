import { cleanup, fireEvent, render, screen, within } from '@testing-library/react';
import { afterEach, beforeEach, expect, it, vi } from 'vitest';
import i18n from '@/i18n/config';
import { RecordingHomeMemoryView } from '@/features/datalink/workbench-v2/recording-home/RecordingHomeMemoryView';

beforeEach(async () => { await i18n.changeLanguage('en'); });
afterEach(cleanup);

it.each(['grid', 'navigator', 'map', 'chapters', 'detail'])('%s keeps partial ownership but never renders or observes an unproven string reading', layout => {
  const onOpen = vi.fn(), onAdd = vi.fn();
  const unproven = 'UNPROVEN_STRING_READING';
  const { container } = render(<RecordingHomeMemoryView protocol="modbus_tcp" status="ready" checking={false}
    inputs={[{ id: 'S', label: 'String start', start_address: '40001', count: 1, data_type: 'string', skipped_addresses: [] }]}
    values={{ status: 'live', byCell: new Map([['HR:1:string', { value: unproven }]]) }}
    onAdd={onAdd} onOpen={onOpen} onBack={vi.fn()} onClose={vi.fn()} />);
  fireEvent.click(screen.getByTestId(`home-memory-layout-${layout}`));
  const start = screen.getByTestId('home-memory-cell-40001');
  expect(start).toHaveAttribute('data-state', 'used');
  expect(start).toHaveAccessibleName(expect.stringContaining('span after this address cannot be confirmed'));
  expect(container).not.toHaveTextContent(unproven);
  expect(container.querySelector('[data-value-key="HR:1:string"]')).toBeNull();
  fireEvent.click(start);
  expect(onOpen).toHaveBeenCalledExactlyOnceWith('S');
  const remainder = screen.getByTestId('home-memory-cell-40002');
  expect(remainder).toHaveAttribute('data-state', 'unconfirmed');
  fireEvent.click(remainder);
  expect(onAdd).not.toHaveBeenCalled();
  if (layout === 'map' || layout === 'detail') {
    fireEvent.mouseMove(remainder);
    const inspector = screen.getByTestId(layout === 'map' ? 'home-memory-summary' : 'home-memory-detail');
    expect(within(inspector).queryByTestId('home-memory-detail-action')).not.toBeInTheDocument();
  }
});
