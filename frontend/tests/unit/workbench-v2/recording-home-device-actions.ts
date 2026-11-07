import { fireEvent, screen, waitFor, within } from '@testing-library/react';

/** Open the devices tab and return the requested action for the same device identity. */
export async function homeDeviceAction(id = 'D', kind: 'connection' | 'points' = 'points'): Promise<HTMLElement> {
  const tab = await screen.findByRole('button', { name: /^(Devices|設備|devicesTab)$/ });
  if (tab.getAttribute('aria-pressed') !== 'true') fireEvent.click(tab);
  return waitFor(() => {
    const table = screen.getByTestId('home-management-table');
    const connection = within(table).getAllByRole('button', { name: /^(Connection|連線設定|deviceConnection)$/ })
      .find(button => button.getAttribute('data-home-focus') === `device-${id}`);
    const row = connection?.closest('tr');
    if (!connection || !row) throw new Error(`Device ${id} is not in the current device list`);
    return kind === 'connection' ? connection : within(row).getByRole('button', { name: /^(Read points|讀取點位|devicePoints)$/ });
  });
}
