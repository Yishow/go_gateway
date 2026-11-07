import { cleanup, fireEvent, render, screen, within } from '@testing-library/react';
import { afterEach, expect, it, vi } from 'vitest';
import { RecordingHomeDeviceStrip } from '@/features/datalink/workbench-v2/recording-home/RecordingHomeDeviceStrip';
import { RecordingHomeDeviceActions } from '@/features/datalink/workbench-v2/recording-home/RecordingHomeDeviceActions';

vi.mock('react-i18next', () => ({ useTranslation: () => ({ t: (key: string, options?: { name?: string }) => options?.name ? `${key} ${options.name}` : key }) }));
afterEach(cleanup);

const devices = [{ id: 'D1', name: 'mc_new', protocol: 'mc', availability: 'available' as const }];
function mount(available = true) {
  const handlers = { select: vi.fn(), open: vi.fn(), create: vi.fn(), remove: vi.fn() };
  render(<RecordingHomeDeviceStrip devices={devices} observations={new Map()} selected="" now={Date.now()}
    onSelect={handlers.select} onOpenDevice={handlers.open} onCreateDevice={handlers.create}
    renderDeviceActions={device => <RecordingHomeDeviceActions device={device} available={available} onDelete={handlers.remove} />} />);
  return handlers;
}

it('offers edit, view and one-click delete on each card, plus an add-device card', () => {
  const handlers = mount();
  const strip = screen.getByRole('region', { name: 'deviceStatus.strip' });
  fireEvent.click(within(strip).getByRole('button', { name: 'deviceStatus.editNamed mc_new' }));
  expect(handlers.open).toHaveBeenLastCalledWith('D1', 'connection');
  fireEvent.click(within(strip).getByRole('button', { name: 'deviceStatus.viewNamed mc_new' }));
  expect(handlers.open).toHaveBeenLastCalledWith('D1', 'points');
  const remove = within(strip).getByRole('button', { name: 'deleteDialogs.device.deleteNamed mc_new' });
  fireEvent.click(remove);
  expect(handlers.remove).toHaveBeenCalledExactlyOnceWith(devices[0]);
  expect(remove).toHaveFocus();
  expect(screen.queryByRole('menu')).not.toBeInTheDocument();
  fireEvent.click(within(strip).getByRole('button', { name: 'createDevice' }));
  expect(handlers.create).toHaveBeenCalledOnce();
  fireEvent.click(screen.getByTestId('device-status-D1'));
  expect(handlers.select).toHaveBeenCalledExactlyOnceWith('D1');
  expect(handlers.open).toHaveBeenCalledTimes(2);
});

it('keeps delete disabled with its reason while the device list is stale', () => {
  const handlers = mount(false);
  const remove = screen.getByRole('button', { name: 'deleteDialogs.device.deleteNamed mc_new' });
  expect(remove).toBeDisabled();
  expect(remove).toHaveAttribute('title', 'deleteDialogs.device.stale');
  expect(screen.getByText('deleteDialogs.device.stale')).toBeVisible();
  expect(remove).toHaveAccessibleDescription('deleteDialogs.device.stale');
  fireEvent.click(remove);
  expect(handlers.remove).not.toHaveBeenCalled();
});

it('keeps the summary bounded to four cards, keeps a selected device visible, and expands or searches the full list', () => {
  const many = Array.from({ length: 9 }, (_, index) => ({ id: `D${index}`, name: `Device ${index}`, protocol: 'mc', availability: index === 1 ? 'unknown' : 'available' }));
  render(<RecordingHomeDeviceStrip devices={many} observations={new Map()} selected="D7" now={Date.now()} onSelect={vi.fn()} onOpenDevice={vi.fn()} />);
  const strip = screen.getByRole('region', { name: 'deviceStatus.strip' });
  expect(within(strip).getAllByRole('article')).toHaveLength(5);
  expect(within(strip).getByTestId('device-status-D7')).toHaveAttribute('aria-pressed', 'true');
  expect(within(strip).getByTestId('device-status-D1').querySelector('[data-availability-badge]')).toBeTruthy();
  expect(within(strip).getByTestId('device-status-D0').querySelector('[data-availability-badge]')).toBeNull();
  const toggle = within(strip).getByRole('button', { name: 'deviceStatus.showAll' });
  expect(toggle).toHaveAttribute('aria-expanded', 'false');
  fireEvent.click(toggle);
  expect(within(strip).getAllByRole('article')).toHaveLength(9);
  fireEvent.click(within(strip).getByRole('button', { name: 'deviceStatus.showLess' }));
  fireEvent.change(within(strip).getByRole('searchbox', { name: 'deviceKeyword' }), { target: { value: 'device 8' } });
  expect(within(strip).getAllByRole('article').map(card => card.querySelector('strong')?.textContent)).toEqual(['Device 8']);
  fireEvent.change(within(strip).getByRole('searchbox', { name: 'deviceKeyword' }), { target: { value: 'nothing' } });
  expect(within(strip).getByRole('status')).toHaveTextContent('noResults');
});

it('clears device search and expansion when the list shrinks to the summary size', () => {
  const many = Array.from({ length: 6 }, (_, index) => ({ id: `D${index}`, name: `Device ${index}`, protocol: 'mc', availability: 'available' }));
  const view = (list: typeof many) => <RecordingHomeDeviceStrip devices={list} observations={new Map()} selected="" now={0} onSelect={vi.fn()} />;
  const { rerender } = render(view(many));
  fireEvent.change(screen.getByRole('searchbox', { name: 'deviceKeyword' }), { target: { value: 'Device 5' } });
  expect(screen.getAllByRole('article')).toHaveLength(1);
  rerender(view(many.slice(0, 4)));
  expect(screen.queryByRole('searchbox', { name: 'deviceKeyword' })).not.toBeInTheDocument();
  expect(screen.getAllByRole('article')).toHaveLength(4);
  rerender(view(many));
  expect(screen.getByRole('searchbox', { name: 'deviceKeyword' })).toHaveValue('');
  expect(screen.getAllByRole('article')).toHaveLength(4);
});
