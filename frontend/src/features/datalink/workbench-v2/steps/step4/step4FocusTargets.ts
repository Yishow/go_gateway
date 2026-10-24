/** Opens the destination settings card of Step 4 and puts focus on its control, the repair entry for a destination cause. */
export function focusDestinationSettings(): void {
  const toggle = document.querySelector<HTMLButtonElement>('[data-testid="database-connection-card-toggle"]');
  if (!toggle) return;
  if (toggle.getAttribute('aria-expanded') === 'false') toggle.click();
  toggle.focus();
}

/** Moves focus into the counter confirmation Basic shows for exactly this source; false when it is not on screen. */
export function focusCounterConfirmation(tagId: string): boolean {
  const entry = document.querySelector<HTMLElement>(`[data-testid="circuit-counter-confirmation-source-${tagId}"]`);
  const control = entry?.querySelector<HTMLElement>('select, input, button');
  if (!entry || !control) return false;
  entry.scrollIntoView?.({ block: 'center' });
  control.focus();
  return true;
}

/** Moves focus to the existing workspace readiness and activation area once Advanced has opened. */
export function focusWorkspaceReadiness(): void {
  window.setTimeout(() => document.getElementById('step4-workspace-readiness')?.focus(), 0);
}
