import * as React from 'react';
import type { WriteGroup } from '../../../../../types/studioV2WriteGroup';
import type { CircuitProfileLocalDraft } from '../../state/circuitProfileBasic';
import { applyUpdatedGroupsToActiveDraft, syncBatchDraftsAfterSave } from '../../state/circuitProfileBatch';
import type { GroupCandidate } from '../../state/writeGroup/candidates';

interface BatchHostInput {
  /** Identity of the Basic scope rendering this host (workspace, device, circuit and destination). */
  scopeKey: string;
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
  const onSaved = async (updated: WriteGroup[], nextRevision: string) => {
    if (input.workspaceId && input.deviceId && updated.length > 0) {
      // The launch scope owns the storage it saved into, even if the screen has moved on.
      syncBatchDraftsAfterSave(input.workspaceId, input.deviceId, input.candidates, input.savedLongGroups, updated);
    }
    if (mounted.current && latestScope.current === launchScope) {
      input.setRevisionOverride({ base: input.queriedRevision, revision: nextRevision });
      input.setLocalDraft((prev) => applyUpdatedGroupsToActiveDraft(prev, input.activeCircuitId, updated, input.savedLongGroups));
    }
    await reload();
  };

  return { busy, setHeld, reload, onSaved };
}
