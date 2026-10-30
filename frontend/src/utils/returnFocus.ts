const regionOwnerKey: unique symbol = Symbol('recording-return-region-owner');
type FocusRegion = HTMLElement & { [regionOwnerKey]?: unknown };
/** Attach the existing host owner identity; no presentation or command state is stored. */
export function markReturnRegion(node: HTMLElement | null, owner: unknown): void { if (node) (node as FocusRegion)[regionOwnerKey] = owner; }
/** Ephemeral DOM-only ticket. Caller retains all scope/revision/operation authority. */
export interface ReturnFocusTicket {
  owner: unknown;
  trigger: HTMLElement | null;
  scopeNode: HTMLElement | null;
  scope: string | null;
  region: HTMLElement | null;
  regionScope: string | null;
  regionOwner: unknown;
  heading: HTMLElement | null;
}
export function eligibleFocusTarget(target: HTMLElement | null | undefined): target is HTMLElement {
  if (!target?.isConnected || target.matches(':disabled, [aria-disabled="true"]')) return false;
  for (let node: HTMLElement | null = target; node; node = node.parentElement) {
    if (node.hidden || node.hasAttribute('inert') || node.getAttribute('aria-hidden') === 'true') return false;
    // Closed disclosure content is natively hidden even when computed CSS says otherwise.
    if (node !== target && node.tagName === 'DETAILS' && !(node as HTMLDetailsElement).open) {
      const summary: Element | undefined = [...node.children].find(child => child.tagName === 'SUMMARY');
      if (!summary?.contains(target)) return false;
    }
    const style = node.ownerDocument.defaultView?.getComputedStyle(node);
    if (style?.display === 'none' || style?.visibility === 'hidden' || style?.visibility === 'collapse') return false;
  }
  return true;
}
export function captureReturnFocus(owner: unknown, trigger: HTMLElement | null): ReturnFocusTicket {
  const scopeNode = trigger?.closest<HTMLElement>('[data-focus-scope]') ?? null;
  const region = trigger?.closest<HTMLElement>('[data-focus-return-region]') ?? null;
  return { owner, trigger, scopeNode, region, regionOwner: region ? (region as FocusRegion)[regionOwnerKey] : undefined, regionScope: region?.dataset.focusReturnRegion ?? null, heading: region?.querySelector<HTMLElement>('[data-focus-return-heading]') ?? null, scope: scopeNode?.dataset.focusScope ?? null };
}
export function restoreReturnFocus(ticket: ReturnFocusTicket | null | undefined, owner: unknown, fallback?: HTMLElement | null): 'trigger' | 'fallback' | 'none' {
  if (!ticket) return 'none';
  const sameOwner = ticket.owner === owner;
  const currentScope = !ticket.scopeNode || ticket.scopeNode.isConnected && ticket.scopeNode.dataset.focusScope === ticket.scope;
  if (sameOwner && currentScope && eligibleFocusTarget(ticket.trigger)) { ticket.trigger.focus({ preventScroll: true }); return 'trigger'; }
  if (sameOwner && eligibleFocusTarget(fallback)) { fallback.focus({ preventScroll: true }); return 'fallback'; }
  if (ticket.region?.isConnected && (sameOwner || ticket.regionOwner !== undefined && (ticket.region as FocusRegion)[regionOwnerKey] === ticket.regionOwner) && ticket.region.dataset.focusReturnRegion === ticket.regionScope && eligibleFocusTarget(ticket.heading)) { ticket.heading.focus({ preventScroll: true }); return 'fallback'; }
  return 'none';
}
