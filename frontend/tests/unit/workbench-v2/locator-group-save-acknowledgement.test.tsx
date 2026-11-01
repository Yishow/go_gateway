import { act, renderHook } from '@testing-library/react';
import { beforeEach, expect, it, vi } from 'vitest';
import { useLocatorDeviceSourceLease } from '@/pages/datalink/workbench-v2/useLocatorDeviceSourceLease';
import type { LocatorGroupSaveAcknowledgement } from '@/pages/datalink/workbench-v2/deviceSourceSubmission';
import { studioV2WorkspaceAPI } from '@/services/studioV2Workspace';
import { studioV2WorkspaceDevicesAPI } from '@/services/studioV2WorkspaceDevices';
import { studioV2WorkspaceWriteGroupsAPI } from '@/services/studioV2WorkspaceWriteGroups';
import { locatorRead, nativeDevice } from './studio-locator-fixtures';
vi.mock('@/services/studioV2Workspace', () => ({ studioV2WorkspaceAPI: { get: vi.fn() } }));
vi.mock('@/services/studioV2WorkspaceDevices', () => ({ studioV2WorkspaceDevicesAPI: { list: vi.fn() } }));
vi.mock('@/services/studioV2WorkspaceWriteGroups', () => ({ studioV2WorkspaceWriteGroupsAPI: { list: vi.fn() } }));
function entry(revision = 'W2', groupRevision = 'G3') {
  const read = locatorRead(4); read.workspace.ordered_device_ids = ['D'];
  read.entry.workspace_revision = revision; read.entry.write_groups!.workspace_revision = revision;
  read.entry.write_groups!.groups[0].revision = groupRevision; read.entry.group_heads.G = groupRevision;
  read.entry.editor_hydrated = true;
  if (read.selection.state === 'selected') Object.assign(read.selection.target, { workspaceRevision: revision, groupRevision });
  return read;
}
function proof(previous = entry(), current = entry('W3', 'G4')): LocatorGroupSaveAcknowledgement {
  return { previous, current, acknowledgement: { workspace_revision: current.entry.workspace_revision, group: structuredClone(current.entry.write_groups!.groups[0]) } };
}
function setup() {
  const original = entry();
  return { original, ...renderHook(({ read, enabled }) => useLocatorDeviceSourceLease(read, 'W', enabled), { initialProps: { read: original, enabled: true } }) };
}
function server(read = entry()) {
  vi.mocked(studioV2WorkspaceAPI.get).mockResolvedValue(read.workspace);
  vi.mocked(studioV2WorkspaceWriteGroupsAPI.list).mockResolvedValue(read.entry.write_groups!);
  vi.mocked(studioV2WorkspaceDevicesAPI.list).mockResolvedValue([nativeDevice()]);
}
beforeEach(() => { vi.resetAllMocks(); server(); });
it('acknowledges exact own readback without granting or rewriting the original locator, then reacquires through both brackets', async () => {
  const hook = setup(), original = JSON.stringify(hook.original), next = proof(); server(next.current);
  await act(async () => { expect(hook.result.current.acknowledgeLocatorGroupSave(next)).toBe(true); });
  expect(hook.result.current.allowed).toBe(false); expect(studioV2WorkspaceAPI.get).not.toHaveBeenCalled();
  expect(JSON.stringify(hook.original)).toBe(original);
  await act(async () => { expect(await hook.result.current.acquire('W')).toBe(true); });
  expect(studioV2WorkspaceAPI.get).toHaveBeenCalledTimes(2); expect(studioV2WorkspaceWriteGroupsAPI.list).toHaveBeenCalledTimes(2);
  expect(hook.result.current.allowed).toBe(true);
  const newer = proof(next.current, entry('W4', 'G5')); server(newer.current);
  await act(async () => { expect(hook.result.current.acknowledgeLocatorGroupSave(newer)).toBe(true); });
  expect(hook.result.current.allowed).toBe(false);
  expect(hook.result.current.acknowledgeLocatorGroupSave(next)).toBe(false);
});
it('retains an existing source receipt recovery eligibility without treating the group acknowledgement as a grant', async () => {
  const hook = setup();
  await act(async () => { expect(await hook.result.current.acquire('W')).toBe(true); });
  const next = proof(); server(next.current);
  await act(async () => { expect(hook.result.current.acknowledgeLocatorGroupSave(next)).toBe(true); });
  expect(hook.result.current.allowed).toBe(false);
  await act(async () => { expect(await hook.result.current.readback()).toEqual([nativeDevice()]); });
  expect(hook.result.current.allowed).toBe(true);
});
it.each(['prior-workspace', 'prior-revision', 'prior-group', 'prior-group-revision', 'ack-workspace', 'ack-group', 'ack-revision', 'ack-empty',
  'current-workspace', 'current-group', 'current-revision', 'current-head', 'unhydrated', 'foreign-list', 'deleted'])('rejects %s proof and cannot grant a newer server scope', async kind => {
  const hook = setup(), next = proof();
  if (next.previous.selection.state !== 'selected' || next.current.selection.state !== 'selected') throw new Error('fixture');
  if (kind === 'prior-workspace') Object.assign(next.previous.selection.target, { workspaceId: 'OTHER' });
  if (kind === 'prior-revision') next.previous.entry.workspace_revision = 'OLD';
  if (kind === 'prior-group') Object.assign(next.previous.selection.target, { groupId: 'OTHER' });
  if (kind === 'prior-group-revision') Object.assign(next.previous.selection.target, { groupRevision: 'OLD' });
  if (kind === 'ack-workspace') next.acknowledgement.group.workspace_id = 'OTHER';
  if (kind === 'ack-group') next.acknowledgement.group.id = 'OTHER';
  if (kind === 'ack-revision') next.acknowledgement.group.revision = 'UNOBSERVED';
  if (kind === 'ack-empty') next.acknowledgement.workspace_revision = '';
  if (kind === 'current-workspace') next.current.workspace.id = 'OTHER';
  if (kind === 'current-group') Object.assign(next.current.selection.target, { groupId: 'OTHER' });
  if (kind === 'current-revision') next.current.entry.workspace_revision = 'UNOBSERVED';
  if (kind === 'current-head') next.current.entry.group_heads.G = 'UNOBSERVED';
  if (kind === 'unhydrated') Reflect.deleteProperty(next.current.entry, 'editor_hydrated');
  if (kind === 'foreign-list') next.current.entry.write_groups!.workspace_id = 'OTHER';
  if (kind === 'deleted') next.current.entry.write_groups!.groups[0].status = 'deleted';
  expect(hook.result.current.acknowledgeLocatorGroupSave(next)).toBe(false);
  server(entry('W3', 'G4'));
  await act(async () => { expect(await hook.result.current.acquire('W')).toBe(false); });
  expect(hook.result.current.allowed).toBe(false);
});
it.each(['entry', 'suspension', 'roundtrip'])('a captured acknowledgement cannot borrow a newer owner after %s', kind => {
  const hook = setup(), acknowledge = hook.result.current.acknowledgeLocatorGroupSave;
  hook.rerender({ read: kind === 'suspension' ? hook.original : entry(), enabled: kind !== 'suspension' });
  if (kind === 'roundtrip') hook.rerender({ read: hook.original, enabled: true });
  expect(acknowledge(proof())).toBe(false); expect(hook.result.current.allowed).toBe(false);
});
it('invalidates an unfinished old acquisition without promoting its late result', async () => {
  const hook = setup(); let finish!: (value: ReturnType<typeof nativeDevice>[]) => void;
  vi.mocked(studioV2WorkspaceDevicesAPI.list).mockImplementationOnce(() => new Promise(resolve => { finish = resolve; }));
  let pending!: Promise<boolean>;
  await act(async () => { pending = hook.result.current.acquire('W'); await Promise.resolve(); await Promise.resolve(); });
  const next = proof(); server(next.current);
  await act(async () => { expect(hook.result.current.acknowledgeLocatorGroupSave(next)).toBe(true); finish([nativeDevice()]); expect(await pending).toBe(false); });
  expect(hook.result.current.allowed).toBe(false);
  await act(async () => { expect(await hook.result.current.acquire('W')).toBe(true); });
});
