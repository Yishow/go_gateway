import { useEffect, useRef } from 'react';
import { useRecordingDraftGuard, useRecordingNavigation } from '../../recording-editor/RecordingDraftNavigation';
import { clearBasicLegacyDraft, loadBasicLegacyDraft } from '../../state/basicLegacyDraft';
import { candidateKey, defaultMemberRole, loadCircuitProfileLocalDraft, type CircuitProfileLocalDraft } from '../../state/circuitProfileBasic';
import type { GroupCandidate } from '../../state/writeGroup/candidates';
import type { WriteGroup } from '../../../../../types/studioV2WriteGroup';

/** Observe the original Basic control/storage owner. No second draft or operation owner. */
export function useBasicLegacyDraftNavigation(input: {
  scopeKey: string; group?: WriteGroup; visible: boolean; pending: boolean;
  interval: number; policy: string; setInterval: (value: number) => void; setPolicy: (value: string) => void;
}) {
  const { scopeKey, group, visible, pending, interval, policy, setInterval, setPolicy } = input;
  const hasGroup = Boolean(group);
  const savedInterval = group?.row_policy.interval_seconds ?? 60;
  const savedPolicy = group?.row_policy.incomplete_policy === 'partial' ? 'partial' : 'skip_row';
  const controls = useRef({ interval, policy }); controls.current = { interval, policy };
  const hydrated = useRef<{ scope: string; groupId?: string; interval: number; policy: string } | undefined>(undefined);
  useEffect(() => {
    const draft = scopeKey ? loadBasicLegacyDraft(scopeKey) : undefined;
    if (draft && hasGroup && (draft.intervalSeconds === undefined || draft.intervalSeconds === savedInterval) &&
      (draft.incompletePolicy === undefined || draft.incompletePolicy === savedPolicy)) clearBasicLegacyDraft(scopeKey);
    const next = { scope: scopeKey, groupId: group?.id, interval: draft?.intervalSeconds ?? savedInterval, policy: draft?.incompletePolicy ?? savedPolicy };
    const previous = hydrated.current; hydrated.current = next;
    const current = controls.current;
    if (previous?.scope === scopeKey && previous.groupId === group?.id &&
      (current.interval !== previous.interval || current.policy !== previous.policy) &&
      (current.interval !== next.interval || current.policy !== next.policy)) return;
    setInterval(next.interval); setPolicy(next.policy);
  }, [scopeKey, group?.id, hasGroup, savedInterval, savedPolicy, setInterval, setPolicy]);
  const stored = loadBasicLegacyDraft(scopeKey);
  const restoredInterval = stored?.intervalSeconds ?? savedInterval, restoredPolicy = stored?.incompletePolicy ?? savedPolicy;
  return useRecordingDraftGuard({ version: JSON.stringify([scopeKey, interval, policy, restoredInterval, restoredPolicy]),
    dirty: visible && (interval !== restoredInterval || policy !== restoredPolicy), pending: visible && pending,
    discard: () => { setInterval(restoredInterval); setPolicy(restoredPolicy); },
  });
}

/** Original scoped storage readback permits offline navigation without claiming backend persistence/readiness. */
export function useCircuitDraftNavigation(input: {
  workspaceId: string; deviceId: string; circuitId?: string; candidates: GroupCandidate[]; groups: WriteGroup[];
  explicit: boolean; draft: CircuitProfileLocalDraft; pending: boolean; setDraft: (value: CircuitProfileLocalDraft) => void;
}) {
  const { workspaceId, deviceId, circuitId, candidates, groups, explicit, draft, pending, setDraft } = input;
  const restored = loadCircuitProfileLocalDraft(workspaceId, deviceId, candidates, groups, circuitId, explicit);
  // The original form treats an omitted role as defaultMemberRole; canonical Save readback
  // makes that same role explicit. Compare that UI meaning, not map insertion/omission.
  // Retain all explicit keys and every other raw field so genuine loss still blocks.
  const recoverableIntent = (value: CircuitProfileLocalDraft) => {
    const roles = { ...value.member_roles };
    for (const candidate of candidates) roles[candidateKey(candidate)] ??= defaultMemberRole(candidate);
    return JSON.stringify({ ...value, member_roles: Object.fromEntries(Object.entries(roles).sort(([a], [b]) => a.localeCompare(b))) });
  };
  return useRecordingDraftGuard({ version: JSON.stringify([workspaceId, deviceId, circuitId, draft, restored]),
    dirty: recoverableIntent(draft) !== recoverableIntent(restored), pending,
    discard: () => { setDraft(restored); },
  });
}

/** A temporary destination denial cannot replace a previously admitted database form with Share-only UI. */
export function useShareOnlyPresentation(workspaceId: string, savedDatabase: boolean, canonicalLoaded: boolean, shareConfigured: boolean) {
  const { blocked } = useRecordingNavigation();
  const admitted = useRef({ workspaceId, database: false });
  if (admitted.current.workspaceId !== workspaceId) admitted.current = { workspaceId, database: false };
  if (savedDatabase && canonicalLoaded) admitted.current.database = true;
  return shareConfigured && !savedDatabase && !(admitted.current.database && blocked());
}
