import { useCallback, useRef } from 'react';
import type { EntityLocator } from '../../../types/runtimeRuleEvidence';
import type { StudioLocatorEntryRead } from '../../../types/studioV2LocatorEntry';
import { buildStudioEntityLocatorTarget } from '../../../services/studioEntityLocator';

type Recovery = { accepted: StudioLocatorEntryRead; target: string; rejected: string; returned: boolean; ruleRevision?: string };
/** One rejected-navigation recovery, never a source cache or a new persistence proof. */
export function useStudioLocatorRecovery(locator: EntityLocator | null, fresh: StudioLocatorEntryRead | undefined,
  navigateTo?: (target: string, options?: { replace?: boolean }) => void) {
  const token = useRef<Recovery | undefined>(undefined);
  const target = locator ? buildStudioEntityLocatorTarget(locator) : null;
  let current = token.current;
  // A different request (including the same B after returning to D) is a new viewing attempt.
  if (current && (current.returned ? target !== current.target : target !== current.target && target !== current.rejected)) {
    token.current = undefined; current = undefined;
  }
  const original = current?.accepted;
  const selected = fresh?.selection.state === 'selected' ? fresh.selection.target : undefined;
  const owned = original?.selection.state === 'selected' ? original.selection.target : undefined;
  const restores = Boolean(current && target === current.target && fresh && original && selected && owned &&
    buildStudioEntityLocatorTarget(fresh.locator) === current.target && fresh.workspace.id === original.workspace.id &&
    fresh.workspace.id === fresh.locator.workspace_id && selected.workspaceId === owned.workspaceId &&
    selected.workspaceRevision === owned.workspaceRevision && selected.step === owned.step &&
    selected.deviceId === owned.deviceId && selected.ruleId === owned.ruleId && selected.ruleRevision === current.ruleRevision && selected.mappingId === owned.mappingId &&
    selected.groupId === owned.groupId && selected.groupRevision === owned.groupRevision);
  if (restores && current) current.returned = true;
  const restoreRejected = useCallback((blocked: boolean, accepted?: StudioLocatorEntryRead, retainedRuleRevision?: string) => {
    if (!blocked || !navigateTo || !target || token.current || !fresh || !accepted || accepted.selection.state !== 'selected' || accepted.selection.target.step > 2 || fresh.selection.state !== 'selected') return;
    const originalTarget = buildStudioEntityLocatorTarget(accepted.locator);
    if (!originalTarget || target === originalTarget || accepted.workspace.id !== accepted.locator.workspace_id ||
        fresh.workspace.id !== accepted.workspace.id || fresh.selection.target.workspaceId !== accepted.workspace.id ||
        fresh.selection.target.workspaceRevision !== accepted.selection.target.workspaceRevision) return;
    token.current = { accepted, target: originalTarget, rejected: target, returned: false, ruleRevision: retainedRuleRevision ?? accepted.selection.target.ruleRevision };
    navigateTo(originalTarget, { replace: true });
  }, [fresh, navigateTo, target]);
  // Failed/missing/foreign/revision-changed reads never get the retained entry.
  // The old object is retained only to preserve its existing editor and ACK token.
  const returning = Boolean(current && target === current.target);
  return { unavailable: returning && !restores, entry: returning ? restores ? original : undefined : fresh, refused: returning && Boolean(fresh) && !restores, restoreRejected };
}
