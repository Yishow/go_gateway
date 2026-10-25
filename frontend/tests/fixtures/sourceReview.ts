import { fireEvent, screen } from '@testing-library/react';

/**
 * Source lists start collapsed behind their count. Opens one the way an operator does, with a real click on its
 * expand button, and leaves it alone when it is already open. `prefix` picks the Circuit list (default) or the
 * legacy Basic list (`basic-recording-sources`).
 */
export function expandSources(prefix = 'circuit-profile-sources'): void {
  const toggle = screen.getByTestId(`${prefix}-toggle`);
  if (toggle.getAttribute('aria-expanded') === 'false') fireEvent.click(toggle);
}

/** Opens the per source status list of one applied group, the way an operator does, with a real click. */
export function expandStatus(groupId: string): void {
  expandSources(`circuit-status-sources-${groupId}`);
}
