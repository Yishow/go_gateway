import * as React from 'react';
import type { WriteGroup } from '../../../../../types/studioV2WriteGroup';
import { groupCircuitId, type CircuitProfileLocalDraft } from '../../state/circuitProfileBasic';
import { applyUpdatedGroupsToActiveDraft, syncBatchDraftsAfterSave } from '../../state/circuitProfileBatch';
import type { GroupCandidate } from '../../state/writeGroup/candidates';

interface BatchHostInput {
  /** Identity of the Basic scope rendering this host (workspace, device, circuit and destination). */
  scopeKey: string;
  /** Current native controls only; used to avoid replacing later operator edits at settlement. */
  nativeConsent?: string;
  /** Component-local versions incremented only by the native edit handler, never by hydration/navigation. */
  nativeEdits?: ReadonlyMap<string, number>;
  workspaceId?: string;
  deviceId: string;
  candidates: GroupCandidate[];
  savedLongGroups: WriteGroup[];
  activeCircuitId?: string;
  queriedRevision: string;
  refetchGroups: () => Promise<{ isError?: boolean; error?: unknown } | undefined>;
  setRevisionOverride: (value: { base: string; revision: string }) => void;
  setLocalDraft: (updater: (prev: CircuitProfileLocalDraft) => CircuitProfileLocalDraft) => void;
}

/**
 * Host side of the batch editor. `onSaved` is rebuilt on every render, so the copy an editor run keeps
 * knows the scope that launched it; its reply is compared with the newest scope before touching state.
 */
export function useCircuitProfileBatchHost(input: BatchHostInput) {
  const [busy, setBusy] = React.useState(false);
  const mounted = React.useRef(true);
  const latestScope = React.useRef(input.scopeKey);
  latestScope.current = input.scopeKey;
  const consentKey = JSON.stringify([input.scopeKey, input.nativeConsent]);
  const consent = React.useRef({ key: consentKey, epoch: 0 });
  if (consent.current.key !== consentKey) consent.current = { key: consentKey, epoch: consent.current.epoch + 1 };
  const latestNativeEdits = React.useRef(input.nativeEdits); latestNativeEdits.current = input.nativeEdits;
  React.useEffect(() => { mounted.current = true; return () => { mounted.current = false; }; }, []);
  // Every run holds the lock once and releases it once, so an old run never unlocks a newer one.
  const holds = React.useRef(0);
  const setHeld = React.useCallback((held: boolean) => {
    holds.current = Math.max(0, holds.current + (held ? 1 : -1));
    if (mounted.current) setBusy(holds.current > 0);
  }, []);

  const reload = async () => {
    const result = await input.refetchGroups();
    if (result?.isError || result?.error) throw new Error('Write groups could not be reloaded.');
  };

  const launchScope = input.scopeKey;
  const launchConsentEpoch = consent.current.epoch;
  const launchNativeEdits = new Map(input.nativeEdits);
  const unchangedConsent = () => input.nativeConsent === undefined || consent.current.epoch === launchConsentEpoch;
  const unchangedDraft = (circuitId = '') => {
    if (!input.nativeEdits) return true; // Existing standalone callers retain their native behavior.
    const key = JSON.stringify([input.workspaceId ?? '', input.deviceId, circuitId]);
    return latestNativeEdits.current !== undefined && (latestNativeEdits.current.get(key) ?? 0) === (launchNativeEdits.get(key) ?? 0);
  };
  const onSaved = async (updated: WriteGroup[], nextRevision: string) => {
    if (input.workspaceId && input.deviceId && updated.length > 0 && (!input.nativeEdits || mounted.current)) {
      // An unmounted guided owner cannot verify a newer instance's edits; canonical acceptance/reload still proceed.
      // A mounted launch scope owns its storage even if the screen has moved on.
      // Protect every actually edited target, including another viewed circuit; untouched accepted targets still sync.
      const toSync = updated.filter(group => unchangedDraft(groupCircuitId(group)));
      syncBatchDraftsAfterSave(input.workspaceId, input.deviceId, input.candidates, input.savedLongGroups, toSync);
    }
    if (mounted.current && latestScope.current === launchScope) {
      input.setRevisionOverride({ base: input.queriedRevision, revision: nextRevision });
      input.setLocalDraft((prev) => unchangedConsent() && unchangedDraft(input.activeCircuitId) ? applyUpdatedGroupsToActiveDraft(prev, input.activeCircuitId, updated, input.savedLongGroups) : prev);
    }
    await reload();
  };

  return { busy, setHeld, reload, onSaved };
}
