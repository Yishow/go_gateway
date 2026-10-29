import { useState } from 'react';
import { afterEach, describe, expect, expectTypeOf, it, vi } from 'vitest';
import { cleanup, fireEvent, render, screen } from '@testing-library/react';
import { studioV2WorkspaceWriteGroupsAPI } from '@/services/studioV2WorkspaceWriteGroups';
import { studioV2RecordingStartAPI } from '@/services/studioV2RecordingStart';
import { resolveStudioSelection, type StudioSelectionEvidence, type StudioSelectionEvidenceProvider,
  type StudioSelectionRequest, type StudioSelectionResult } from '@/features/datalink/workbench-v2/state/studioSelection';

const request: StudioSelectionRequest = { workspaceId: 'W1', step: 4, deviceId: 'A', groupId: 'G1', groupRevisionHint: 'R1' };
type Mutable<T> = T extends object ? { -readonly [K in keyof T]: Mutable<T[K]> } : T;
type CurrentEvidence = Mutable<Extract<StudioSelectionEvidence, { state: 'current' }>>;
function evidence(): CurrentEvidence {
  return { state: 'current', workspace: { id: 'W1', ordered_device_ids: ['A', 'B'] }, workspaceRevision: 'W3',
    devices: [{ id: 'A' }, { id: 'B' }],
    rules: [{ currentRevision: 'S2', value: { id: 'RULE1', workspace_id: 'W1', device_id: 'A', revision_id: 'S2' } }],
    mappings: [{ id: 'MAP1', workspace_id: 'W1', device_id: 'A', rule_id: 'RULE1', point_id: 'P1', tag_id: 'T1' }],
    groups: [{ currentRevision: 'R2', workspaceRevision: 'W3', value: { id: 'G1', workspace_id: 'W1', revision: 'R2', status: 'ready',
      members: [{ device_id: 'A', point_id: 'P1', tag_id: 'T1' }] } }] };
}
const select = (value: StudioSelectionEvidence, selection: StudioSelectionRequest = request) =>
  resolveStudioSelection('W1', selection, { readCurrent: () => value });
afterEach(() => { cleanup(); vi.restoreAllMocks(); });

describe('canonical setup selection', () => {
  it('selects current G1/R2 and retains R1 only as a historical hint', () => {
    const readCurrent = vi.fn(evidence);
    expect(resolveStudioSelection('W1', request, { readCurrent })).toEqual({ state: 'selected',
      target: { workspaceId: 'W1', step: 4, deviceId: 'A', groupId: 'G1', workspaceRevision: 'W3', groupRevision: 'R2' },
      history: { groupId: 'G1', revisionHint: 'R1', relation: 'different' } });
    expect(readCurrent).toHaveBeenCalledOnce();
    expect(readCurrent).toHaveBeenCalledWith(request);
  });
  it('refuses a foreign request before reading foreign details', () => {
    const readCurrent = vi.fn(evidence);
    expect(resolveStudioSelection('W1', { workspaceId: 'W2', step: 4, groupId: 'G2' }, { readCurrent }))
      .toEqual({ state: 'unavailable', reason: 'scope_unavailable' });
    expect(readCurrent).not.toHaveBeenCalled();
  });
  it('requires evidence rather than trusting a canonical-looking request', () => {
    expect(resolveStudioSelection('W1', request)).toEqual({ state: 'unavailable', reason: 'evidence_unavailable' });
  });
  it.each([undefined, { state: 'unknown' }])('degrades absent/unknown provider proof to unavailable: %j', (value) => {
    expect(select(value as StudioSelectionEvidence)).toEqual({ state: 'unavailable', reason: 'evidence_unavailable' });
  });
  it('does not reject a current canonical group solely because an unselected historical member was removed', () => {
    const current = evidence(); current.devices = [{ id: 'A' }]; current.workspace.ordered_device_ids = ['A'];
    current.groups![0] = { ...current.groups![0], value: { ...current.groups![0].value,
      members: [{ device_id: 'A', point_id: 'P1', tag_id: 'T1' }, { device_id: 'REMOVED', point_id: 'OLD', tag_id: 'OLD' }] } };
    expect(select(current)).toMatchObject({ state: 'selected', target: { deviceId: 'A', groupId: 'G1', groupRevision: 'R2' } });
    expect(select(current, { workspaceId: 'W1', step: 4, groupId: 'G1' })).toEqual({ state: 'selected',
      target: { workspaceId: 'W1', step: 4, groupId: 'G1', groupRevision: 'R2', workspaceRevision: 'W3' } });
  });
  it('rejects an unknown group lifecycle status rather than assuming it is a current saved target', () => {
    const current = evidence(); current.groups![0] = { ...current.groups![0], value: { ...current.groups![0].value, status: 'future' as 'ready' } };
    expect(select(current)).toEqual({ state: 'unavailable', reason: 'target_unavailable' });
  });
  it('needs no entity query evidence when the operator only selects a workspace step', () => {
    expect(select({ state: 'current', workspace: { id: 'W1', ordered_device_ids: [] }, workspaceRevision: '' },
      { workspaceId: 'W1', step: 1, deviceId: undefined })).toEqual({ state: 'selected',
      target: { workspaceId: 'W1', step: 1, workspaceRevision: '' } });
  });
  it('refuses a missing request before calling the provider', () => {
    const readCurrent = vi.fn(evidence);
    expect(resolveStudioSelection('W1', undefined as unknown as StudioSelectionRequest, { readCurrent }))
      .toEqual({ state: 'unavailable', reason: 'invalid_request' });
    expect(readCurrent).not.toHaveBeenCalled();
  });
  it('does not invent a missing workspace revision from a group revision', () => {
    const current = evidence(); current.workspaceRevision = undefined as unknown as string;
    expect(select(current)).toEqual({ state: 'unavailable', reason: 'revision_unavailable' });
  });
  it.each(['rule_id', 'device_id', 'point_id', 'tag_id'] as const)('requires canonical mapping relationship field %s', (key) => {
    const current = evidence(); current.mappings![0] = { ...current.mappings![0], [key]: '' };
    expect(select(current, { workspaceId: 'W1', step: 3, mappingId: 'MAP1' }))
      .toEqual({ state: 'unavailable', reason: 'target_unavailable' });
  });
  it.each(['missing_workspace', 'missing_device', 'parent_conflict'] as const)('requires rule ownership independently of a mapping claim: %s', (kind) => {
    const current = evidence(); current.rules![0] = { ...current.rules![0], value: { ...current.rules![0].value,
      ...(kind === 'missing_workspace' ? { workspace_id: undefined } : { device_id: kind === 'missing_device' ? '' : 'B' }) } };
    expect(select(current, { workspaceId: 'W1', step: 3, mappingId: 'MAP1' })).toEqual({ state: 'unavailable', reason: 'target_unavailable' });
  });
  it('uses the managed owner for an owned empty saved group but does not infer an owner for an empty custom group', () => {
    const current = evidence(); current.groups![0] = { ...current.groups![0], value: { ...current.groups![0].value,
      members: [], basic_managed_device_id: 'A' } };
    expect(select(current)).toMatchObject({ state: 'selected', target: { deviceId: 'A', groupId: 'G1' } });
    current.groups![0] = { ...current.groups![0], value: { ...current.groups![0].value, basic_managed_device_id: undefined } };
    expect(select(current)).toEqual({ state: 'unavailable', reason: 'target_unavailable' });
    expect(select(current, { workspaceId: 'W1', step: 4, groupId: 'G1' })).toMatchObject({ state: 'selected', target: { groupId: 'G1' } });
  });
  it('opens a current group-only repair target without selecting a deleted managed device', () => {
    const current = evidence(); current.devices = []; current.workspace.ordered_device_ids = [];
    current.groups![0] = { ...current.groups![0], value: { ...current.groups![0].value, basic_managed_device_id: 'REMOVED' } };
    expect(select(current, { workspaceId: 'W1', step: 4, groupId: 'G1' })).toEqual({ state: 'selected',
      target: { workspaceId: 'W1', step: 4, groupId: 'G1', groupRevision: 'R2', workspaceRevision: 'W3' } });
  });
  it('does not let a historical hint establish device membership in the current saved group', () => {
    const current = evidence(); current.groups![0] = { ...current.groups![0], value: { ...current.groups![0].value,
      members: [{ device_id: 'B', point_id: 'P2', tag_id: 'T2' }] } };
    expect(select(current)).toEqual({ state: 'unavailable', reason: 'target_unavailable' });
  });

  it('exposes typed navigation only, never URLs or mutation callbacks', () => {
    expectTypeOf<StudioSelectionRequest['step']>().toEqualTypeOf<1 | 2 | 3 | 4>();
    expectTypeOf<keyof StudioSelectionRequest>().toEqualTypeOf<'workspaceId' | 'step' | 'deviceId' | 'ruleId' | 'mappingId' | 'groupId' | 'groupRevisionHint'>();
    expectTypeOf<keyof StudioSelectionEvidenceProvider>().toEqualTypeOf<'readCurrent'>();
    expectTypeOf<{ workspaceId: string; step: 5 }>().not.toExtend<StudioSelectionRequest>();
    expectTypeOf<{ step: 4 }>().not.toExtend<StudioSelectionRequest>();
  });

  it.each([1, 2, 3, 4] as const)('allows workspace-only step %s without choosing the first device/group', (step) => {
    expect(select(evidence(), { workspaceId: 'W1', step })).toEqual({ state: 'selected',
      target: { workspaceId: 'W1', step, workspaceRevision: 'W3' } });
  });
  it.each([1, 2, 3, 4] as const)('retains independent owned rule/group viewing context at step %s without claiming group membership', (step) => {
    const current = evidence(); current.groups![0] = { ...current.groups![0], value: { ...current.groups![0].value,
      members: [{ device_id: 'A', point_id: 'P_OTHER_RULE', tag_id: 'T_OTHER_RULE' }] } };
    expect(select(current, { workspaceId: 'W1', step, ruleId: 'RULE1', groupId: 'G1' })).toEqual({ state: 'selected',
      target: { workspaceId: 'W1', step, workspaceRevision: 'W3', deviceId: 'A', ruleId: 'RULE1', ruleRevision: 'S2', groupId: 'G1', groupRevision: 'R2' } });
    expect(select(current, { workspaceId: 'W1', step, ruleId: 'RULE1', groupId: 'G1', mappingId: 'MAP1' }))
      .toEqual({ state: 'unavailable', reason: 'target_unavailable' });
  });

  it.each([
    [{ deviceId: 'A' }, { deviceId: 'A' }],
    [{ ruleId: 'RULE1' }, { deviceId: 'A', ruleId: 'RULE1', ruleRevision: 'S2' }],
    [{ mappingId: 'MAP1' }, { deviceId: 'A', ruleId: 'RULE1', ruleRevision: 'S2', mappingId: 'MAP1' }],
    [{ groupId: 'G1' }, { groupId: 'G1', groupRevision: 'R2' }],
  ])('resolves only canonical owned parent identities for %j', (ids, resolved) => {
    expect(select(evidence(), { workspaceId: 'W1', step: 3, ...ids })).toEqual({ state: 'selected',
      target: { workspaceId: 'W1', step: 3, workspaceRevision: 'W3', ...resolved } });
  });

  it.each(['R2', 'R99'])('treats hint %s only as equality/difference, not existence or chronology', (groupRevisionHint) => {
    expect(select(evidence(), { ...request, groupRevisionHint })).toMatchObject({ state: 'selected',
      target: { groupRevision: 'R2' }, history: { revisionHint: groupRevisionHint, relation: groupRevisionHint === 'R2' ? 'current' : 'different' } });
  });

  it('rereads current revisions on each request, never substitutes a hint for the current saved version', () => {
    const first = evidence(); const second = evidence();
    second.groups![0] = { ...second.groups![0], currentRevision: 'R3', value: { ...second.groups![0].value, revision: 'R3' } };
    const readCurrent = vi.fn().mockReturnValueOnce(first).mockReturnValueOnce(second);
    expect(resolveStudioSelection('W1', request, { readCurrent })).toMatchObject({ target: { groupRevision: 'R2' } });
    expect(resolveStudioSelection('W1', request, { readCurrent })).toMatchObject({ target: { groupRevision: 'R3' }, history: { revisionHint: 'R1' } });
    expect(readCurrent).toHaveBeenCalledTimes(2);
  });

  it.each([
    { workspaceId: '' }, { workspaceId: ' W1' }, { step: 0 }, { step: 5 }, { step: '4' },
    { deviceId: '' }, { ruleId: ' ' }, { mappingId: null }, { groupId: 'G1 ' }, { groupRevisionHint: '' },
    { groupId: undefined, groupRevisionHint: 'R1' },
  ])('fails closed on malformed typed-boundary input %j', (patch) => {
    const readCurrent = vi.fn(evidence);
    expect(resolveStudioSelection('W1', { ...request, ...patch } as StudioSelectionRequest, { readCurrent }))
      .toEqual({ state: 'unavailable', reason: 'invalid_request' });
    expect(readCurrent).not.toHaveBeenCalled();
  });
  it('rejects missing current workspace context', () => {
    const readCurrent = vi.fn(evidence);
    expect(resolveStudioSelection('', request, { readCurrent })).toEqual({ state: 'unavailable', reason: 'scope_unavailable' });
    expect(readCurrent).not.toHaveBeenCalled();
  });

  it.each(['loading', 'error', 'stale', 'unavailable'] as const)('preserves safe evidence degradation: %s', (state) => {
    expect(select({ state })).toEqual({ state: 'unavailable', reason: `evidence_${state}` });
  });
  it('contains provider errors without exposing raw backend details', () => {
    expect(resolveStudioSelection('W1', request, { readCurrent: () => { throw new Error('postgres://private/foreign'); } }))
      .toEqual({ state: 'unavailable', reason: 'evidence_error' });
  });
  it('refuses a provider snapshot for a foreign workspace without forwarding data', () => {
    const current = evidence(); current.workspace.id = 'W2';
    expect(select(current)).toEqual({ state: 'unavailable', reason: 'scope_unavailable' });
  });

  it.each(['devices', 'rules', 'mappings', 'groups'] as const)('requires explicit current %s evidence for a selected chain', (field) => {
    const current = evidence(); delete current[field];
    expect(select(current, { ...request, ruleId: 'RULE1', mappingId: 'MAP1' })).toEqual({ state: 'unavailable', reason: 'target_unavailable' });
  });
  it.each([
    { deviceId: 'MISSING' }, { ruleId: 'MISSING' }, { mappingId: 'MISSING' }, { groupId: 'MISSING' },
  ])('never falls back to another entity for a missing ID %j', (patch) => {
    expect(select(evidence(), { ...request, ...patch })).toEqual({ state: 'unavailable', reason: 'target_unavailable' });
  });
  it.each(['not_owned', 'not_present', 'duplicate_record', 'duplicate_ownership'] as const)('requires unambiguous device ownership and existence: %s', (kind) => {
    const current = evidence();
    if (kind === 'not_owned') current.workspace.ordered_device_ids = ['B'];
    if (kind === 'not_present') current.devices = [{ id: 'B' }];
    if (kind === 'duplicate_record') current.devices = [{ id: 'A' }, { id: 'A' }];
    if (kind === 'duplicate_ownership') current.workspace.ordered_device_ids = ['A', 'A'];
    expect(select(current)).toEqual({ state: 'unavailable', reason: 'target_unavailable' });
  });
  it.each(['rule', 'mapping', 'group'] as const)('rejects foreign %s ownership even when IDs match', (kind) => {
    const current = evidence();
    if (kind === 'rule') current.rules![0] = { ...current.rules![0], value: { ...current.rules![0].value, workspace_id: 'W2' } };
    if (kind === 'mapping') current.mappings![0] = { ...current.mappings![0], workspace_id: 'W2' };
    if (kind === 'group') current.groups![0] = { ...current.groups![0], value: { ...current.groups![0].value, workspace_id: 'W2' } };
    expect(select(current, { ...request, ruleId: 'RULE1', mappingId: 'MAP1' })).toEqual({ state: 'unavailable', reason: 'target_unavailable' });
  });
  it.each(['rules', 'mappings', 'groups'] as const)('rejects duplicate %s IDs instead of selecting the first', (kind) => {
    const current = evidence();
    if (kind === 'rules') current.rules = [...current.rules!, current.rules![0]];
    if (kind === 'mappings') current.mappings = [...current.mappings!, current.mappings![0]];
    if (kind === 'groups') current.groups = [...current.groups!, current.groups![0]];
    expect(select(current, { ...request, ruleId: 'RULE1', mappingId: 'MAP1' })).toEqual({ state: 'unavailable', reason: 'target_unavailable' });
  });
  it.each(['rule_device', 'mapping_device', 'mapping_rule', 'group_device', 'group_point', 'group_tag'] as const)('refuses a conflicting relationship: %s', (kind) => {
    const current = evidence();
    if (kind === 'rule_device') current.rules![0] = { ...current.rules![0], value: { ...current.rules![0].value, device_id: 'B' } };
    if (kind === 'mapping_device') current.mappings![0] = { ...current.mappings![0], device_id: 'B' };
    if (kind === 'mapping_rule') current.mappings![0] = { ...current.mappings![0], rule_id: 'RULE2' };
    if (kind.startsWith('group')) current.groups![0] = { ...current.groups![0], value: { ...current.groups![0].value,
      members: [{ device_id: kind === 'group_device' ? 'B' : 'A', point_id: kind === 'group_point' ? 'P2' : 'P1', tag_id: kind === 'group_tag' ? 'T2' : 'T1' }] } };
    expect(select(current, { ...request, ruleId: 'RULE1', mappingId: 'MAP1' })).toEqual({ state: 'unavailable', reason: 'target_unavailable' });
  });

  it.each(['rule_empty', 'rule_head_empty', 'group_empty', 'group_head_empty'] as const)('never invents required backend revision proof: %s', (kind) => {
    const current = evidence();
    if (kind === 'rule_empty') current.rules![0] = { ...current.rules![0], value: { ...current.rules![0].value, revision_id: undefined } };
    if (kind === 'rule_head_empty') current.rules![0] = { ...current.rules![0], currentRevision: '' };
    if (kind === 'group_empty') current.groups![0] = { ...current.groups![0], value: { ...current.groups![0].value, revision: '' } };
    if (kind === 'group_head_empty') current.groups![0] = { ...current.groups![0], currentRevision: '' };
    expect(select(current, { ...request, ruleId: 'RULE1' })).toEqual({ state: 'unavailable', reason: 'revision_unavailable' });
  });
  it.each(['rule', 'group', 'workspace'] as const)('rejects stale %s revisions instead of silently selecting an old cached record', (kind) => {
    const current = evidence();
    if (kind === 'rule') current.rules![0] = { ...current.rules![0], currentRevision: 'S3' };
    if (kind === 'group') current.groups![0] = { ...current.groups![0], currentRevision: 'R3' };
    if (kind === 'workspace') current.groups![0] = { ...current.groups![0], workspaceRevision: 'W2' };
    expect(select(current, { ...request, ruleId: 'RULE1' })).toEqual({ state: 'unavailable', reason: 'revision_conflict' });
  });
  it('allows the backend-proven empty initial workspace revision without synthesizing one', () => {
    const current = evidence(); current.workspaceRevision = '';
    expect(select(current, { workspaceId: 'W1', step: 1 })).toEqual({ state: 'selected', target: { workspaceId: 'W1', step: 1, workspaceRevision: '' } });
  });
  it('does not select a tombstone or same-name replacement group', () => {
    const current = evidence(); current.groups![0] = { ...current.groups![0], value: { ...current.groups![0].value, status: 'deleted' } };
    const replacement = { ...current.groups![0], value: { ...current.groups![0].value, id: 'G2', status: 'ready' as const, name: 'Same name' } };
    current.groups = [...current.groups!, replacement];
    expect(select(current)).toEqual({ state: 'unavailable', reason: 'target_unavailable' });
  });
  it.each(['draft', 'ready', 'running', 'disabled'] as const)('selection of saved status %s is only navigation, not action availability', (status) => {
    const current = evidence(); current.groups![0] = { ...current.groups![0], value: { ...current.groups![0].value, status } };
    expect(select(current)).toEqual({ state: 'selected', target: { workspaceId: 'W1', step: 4, deviceId: 'A', groupId: 'G1', workspaceRevision: 'W3', groupRevision: 'R2' },
      history: { groupId: 'G1', revisionHint: 'R1', relation: 'different' } });
  });
  it('retains selected A for a shared A+B group without claiming exclusive ownership or filtering the group', () => {
    const current = evidence(); current.groups![0] = { ...current.groups![0], value: { ...current.groups![0].value,
      members: [{ device_id: 'A', point_id: 'P1', tag_id: 'T1' }, { device_id: 'B', point_id: 'P2', tag_id: 'T2' }] } };
    const before = structuredClone(current);
    expect(select(current)).toMatchObject({ state: 'selected', target: { deviceId: 'A', groupId: 'G1' } });
    expect(current).toEqual(before);
  });
  it('enforces an exclusive managed-device owner instead of treating it as a shared group', () => {
    const current = evidence(); current.groups![0] = { ...current.groups![0], value: { ...current.groups![0].value, basic_managed_device_id: 'B' } };
    expect(select(current)).toEqual({ state: 'unavailable', reason: 'target_unavailable' });
  });
  it('copies and freezes accepted identity/history and sends only frozen typed request fields to the provider', () => {
    const current = evidence(); const caller = { ...request, injected: 'private detail' };
    const result = resolveStudioSelection('W1', caller, { readCurrent: (received) => {
      expect(received).toEqual(request); expect(Object.isFrozen(received)).toBe(true); return current;
    } });
    caller.groupId = 'G2'; current.groups![0] = { ...current.groups![0], value: { ...current.groups![0].value, revision: 'R9' } };
    expect(result).toMatchObject({ state: 'selected', target: { groupId: 'G1', groupRevision: 'R2' }, history: { revisionHint: 'R1' } });
    expect(Object.isFrozen(result)).toBe(true);
    if (result.state !== 'selected') throw new Error('Expected accepted selection');
    expect(Object.isFrozen(result.target)).toBe(true); expect(Object.isFrozen(result.history)).toBe(true);
    expect(result.target).not.toHaveProperty('groupRevisionHint'); expect(result.target).not.toHaveProperty('injected');
    expect(Object.isFrozen(select({ state: 'error' }))).toBe(true);
  });
  it('never forwards backend-private fields or accepts a same-name group substitution', () => {
    const current = evidence(); const raw = { ...current.groups![0].value, name: 'SECRET', destination: { database: '/private/path' }, ready: true };
    current.groups![0] = { ...current.groups![0], value: raw };
    expect(JSON.stringify(select(current))).not.toMatch(/SECRET|private|destination|ready/);
    current.groups![0] = { ...current.groups![0], value: { ...raw, id: 'G2' } };
    expect(select(current)).toEqual({ state: 'unavailable', reason: 'target_unavailable' });
  });
});

/** Isolated test receiver, not production route/query wiring or session persistence. */
function Receiver({ selection, provider }: { selection: StudioSelectionRequest; provider: StudioSelectionEvidenceProvider }) {
  const [result, setResult] = useState<StudioSelectionResult>();
  return <><button type="button" onClick={() => setResult(resolveStudioSelection('W1', selection, provider))}>Open selection</button>
    <p role="status">{result?.state === 'selected' ? `${result.target.groupId} / ${result.target.groupRevision}` : result?.state === 'unavailable' ? 'Selection unavailable' : 'No selection'}</p></>;
}
describe('isolated selection receiver boundary', () => {
  it('selects only on intent; render, repeated selection, Back-like rerender and unmount never save/apply/start or write storage', () => {
    const forbidden = () => { throw new Error('Mutation during selection'); };
    const saves = vi.spyOn(studioV2WorkspaceWriteGroupsAPI, 'update').mockImplementation(forbidden);
    const creates = vi.spyOn(studioV2WorkspaceWriteGroupsAPI, 'create').mockImplementation(forbidden);
    const applies = vi.spyOn(studioV2WorkspaceWriteGroupsAPI, 'apply').mockImplementation(forbidden);
    const starts = vi.spyOn(studioV2RecordingStartAPI, 'start').mockImplementation(forbidden);
    const storage = vi.spyOn(Storage.prototype, 'setItem').mockImplementation(() => undefined);
    const readCurrent = vi.fn(evidence); const provider = { readCurrent };
    const rendered = render(<Receiver selection={request} provider={provider} />);
    expect(screen.getByRole('status')).toHaveTextContent('No selection'); expect(readCurrent).not.toHaveBeenCalled();
    fireEvent.click(screen.getByRole('button', { name: 'Open selection' }));
    expect(screen.getByRole('status')).toHaveTextContent('G1 / R2');
    rendered.rerender(<Receiver selection={{ ...request, groupRevisionHint: 'R2' }} provider={provider} />);
    expect(readCurrent).toHaveBeenCalledOnce();
    fireEvent.click(screen.getByRole('button', { name: 'Open selection' })); rendered.unmount();
    expect(readCurrent).toHaveBeenCalledTimes(2);
    for (const spy of [saves, creates, applies, starts, storage]) expect(spy).not.toHaveBeenCalled();
  });
  it('shows unavailable for W2/G2 without foreign data, selecting another group or reading its provider', () => {
    const readCurrent = vi.fn(() => { throw new Error('Foreign private detail'); });
    render(<Receiver selection={{ workspaceId: 'W2', step: 4, groupId: 'G2' }} provider={{ readCurrent }} />);
    fireEvent.click(screen.getByRole('button', { name: 'Open selection' }));
    expect(screen.getByRole('status')).toHaveTextContent('Selection unavailable');
    expect(document.body).not.toHaveTextContent(/G1|G2|Foreign private detail/); expect(readCurrent).not.toHaveBeenCalled();
  });
});
