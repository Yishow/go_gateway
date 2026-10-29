import { beforeEach, describe, expect, it, vi } from 'vitest';
import {
  candidateKey, defaultCircuitProfileLocalDraft, deriveCircuitId,
  loadCircuitProfileLocalDraft, saveCircuitProfileLocalDraft,
} from '@/features/datalink/workbench-v2/state/circuitProfileBasic';
import { buildGroupCandidates } from '@/features/datalink/workbench-v2/state/writeGroup/candidates';
import { ensureOrUpdateCircuitProfileGroups } from '@/features/datalink/workbench-v2/steps/step4/circuitProfileActions';
import { savedGroup, savedGroupState } from '../../fixtures/writeGroupState';
import type { WriteGroup, WriteGroupProfileKind } from '@/types/studioV2WriteGroup';

const candidates = buildGroupCandidates(savedGroupState()).candidates;
function profileGroup(circuit: string, kind: WriteGroupProfileKind = 'instantaneous'): WriteGroup {
  return savedGroup({
    id: `group-${circuit}-${kind}`, layout_version: 'managed_long_v1', storage_mode: 'plain',
    recording_profile: { kind, interval: kind === 'cumulative' ? 10 : 15, unit: kind === 'cumulative' ? 'minute' : 'second' },
    members: [{ ...savedGroup().members[0], circuit_id: circuit, circuit_name: '原名稱' }],
    destination: { ...savedGroup().destination, storage_strategy: 'managed' },
  });
}
beforeEach(() => window.sessionStorage.clear());

describe('circuit identity and source selection survive editing', () => {
  it('gives two explicitly created Chinese circuits distinct opaque identities', () => {
    expect(deriveCircuitId('空壓機')).not.toBe(deriveCircuitId('照明'));
  });
  it('keeps persisted circuit identity when its label is renamed or cleared', () => {
    expect(deriveCircuitId('新名稱', 'opaque-circuit-A')).toBe('opaque-circuit-A');
    expect(deriveCircuitId('', 'opaque-circuit-A')).toBe('opaque-circuit-A');
  });
  it('does not invent an identity for unassigned raw tags', () => {
    expect(deriveCircuitId('')).toBe('');
  });
  it('does not pair instantaneous A with cumulative B on reload', () => {
    const draft = defaultCircuitProfileLocalDraft(candidates, [profileGroup('A'), profileGroup('B', 'cumulative')]);
    expect(draft.circuit_id).toBe('A');
    expect(draft.instantaneous_enabled).toBe(true);
    expect(draft.cumulative_enabled).toBe(false);
  });
  it('keeps a deliberately empty selection after reload', () => {
    const draft = { ...defaultCircuitProfileLocalDraft(candidates, []), selected_keys: [], instantaneous_keys: [], cumulative_keys: [] };
    saveCircuitProfileLocalDraft('ws-1', 'dev-1', draft);
    expect(loadCircuitProfileLocalDraft('ws-1', 'dev-1', candidates, []).selected_keys).toEqual([]);
    expect(loadCircuitProfileLocalDraft('ws-1', 'dev-1', candidates, []).instantaneous_keys).toEqual([]);
    expect(loadCircuitProfileLocalDraft('ws-1', 'dev-1', candidates, []).cumulative_keys).toEqual([]);
  });
  it('keeps deliberately empty legacy selection when dedicated arrays are absent', () => {
    window.sessionStorage.setItem('wbv2.circuit-profile-draft.ws-1%2Fdev-1', JSON.stringify({
      circuit_id: 'cir-legacy', circuit_name: 'Legacy', selected_keys: [],
    }));
    const loaded = loadCircuitProfileLocalDraft('ws-1', 'dev-1', candidates, []);
    expect(loaded.selected_keys).toEqual([]);
    expect(loaded.instantaneous_keys).toEqual([]);
    expect(loaded.cumulative_keys).toEqual([]);
  });
  it('treats dedicated profile arrays as single authority even if selected_keys is contradictory', () => {
    const key0 = candidateKey(candidates[0]);
    const key1 = candidateKey(candidates[1]);
    window.sessionStorage.setItem('wbv2.circuit-profile-draft.ws-1%2Fdev-1', JSON.stringify({
      circuit_id: 'cir-1',
      circuit_name: 'Feeder',
      selected_keys: [key0],
      instantaneous_keys: [key0],
      cumulative_keys: [key1],
    }));
    const loaded = loadCircuitProfileLocalDraft('ws-1', 'dev-1', candidates, []);
    expect(loaded.instantaneous_keys).toEqual([key0]);
    expect(loaded.cumulative_keys).toEqual([key1]);
    expect(loaded.selected_keys).toEqual([key0, key1]);
  });
  it('preserves explicit empty per profile independently', () => {
    const key0 = candidateKey(candidates[0]);
    window.sessionStorage.setItem('wbv2.circuit-profile-draft.ws-1%2Fdev-1', JSON.stringify({
      circuit_id: 'cir-1',
      circuit_name: 'Feeder',
      instantaneous_keys: [],
      cumulative_keys: [key0],
    }));
    const loaded = loadCircuitProfileLocalDraft('ws-1', 'dev-1', candidates, []);
    expect(loaded.instantaneous_keys).toEqual([]);
    expect(loaded.cumulative_keys).toEqual([key0]);
    expect(loaded.selected_keys).toEqual([key0]);
  });
  it('stores a new circuit identity once and reuses it on reload', () => {
    const draft = { ...defaultCircuitProfileLocalDraft(candidates, []), circuit_id: deriveCircuitId('空壓機'), circuit_name: '空壓機' };
    saveCircuitProfileLocalDraft('ws-1', 'dev-1', draft);
    expect(loadCircuitProfileLocalDraft('ws-1', 'dev-1', candidates, []).circuit_id).toBe(draft.circuit_id);
    expect(loadCircuitProfileLocalDraft('ws-1', 'dev-1', candidates, []).selected_keys).toContain(candidateKey(candidates[0]));
  });
});

async function saveCircuit(circuitId: string, existing: WriteGroup[], circuitName = '新名稱') {
  const ensure = vi.fn(async () => ({ group: profileGroup(circuitId), workspace_revision: 'ws-new' }));
  const update = vi.fn(async () => ({ group: profileGroup(circuitId), workspace_revision: 'ws-new' }));
  await ensureOrUpdateCircuitProfileGroups({
    workspaceId: 'ws-1', initialWorkspaceRevision: 'ws-old', deviceName: 'Device', connectorKind: 'sqlite',
    enabledProfiles: [{ profile: 'instantaneous', intervalValue: 15, intervalSeconds: 15 }],
    storageMode: 'plain', recreateMissing: true, circuitId, circuitName,
    destination: savedGroup().destination, selectedCandidates: [candidates[0]], measurementsByPointKey: new Map(),
    savedLongGroups: existing, preparedGroups: [], ensureCircuitProfile: ensure, updateWriteGroup: update,
  });
  return { ensure, update };
}

describe('canonical group selection', () => {
  it('creates circuit B instead of updating circuit A on the same device', async () => {
    const { ensure, update } = await saveCircuit('B', [profileGroup('A')]);
    expect(update).not.toHaveBeenCalled();
    expect(ensure).toHaveBeenCalledOnce();
  });
  it('renames the exact circuit through its existing group without replacing identity', async () => {
    const { ensure, update } = await saveCircuit('A', [profileGroup('A')]);
    expect(ensure).not.toHaveBeenCalled();
    expect(update).toHaveBeenCalledWith(expect.objectContaining({
      id: 'group-A-instantaneous', request: expect.objectContaining({
        group: expect.objectContaining({ members: [expect.objectContaining({ circuit_id: 'A', circuit_name: '新名稱' })] }),
      }),
    }));
  });
});
