import { screen, waitFor } from '@testing-library/react';

/** Return the requested action on the device's card in the home device strip (same device identity). */
export async function homeDeviceAction(id = 'D', kind: 'connection' | 'points' = 'points'): Promise<HTMLElement> {
  return waitFor(() => {
    const card = screen.getByTestId(`device-status-${id}`).closest('article');
    const action = card?.querySelector<HTMLElement>(`[data-device-action="${kind}"]`);
    if (!action) throw new Error(`Device ${id} has no ${kind} action in the device strip`);
    return action;
  });
}
