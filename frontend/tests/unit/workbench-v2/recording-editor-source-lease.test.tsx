import { act, renderHook } from '@testing-library/react';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import { useLocatorDeviceSourceLease } from '../../../src/pages/datalink/workbench-v2/useLocatorDeviceSourceLease';
import { studioV2WorkspaceAPI } from '../../../src/services/studioV2Workspace';
import { studioV2WorkspaceDevicesAPI } from '../../../src/services/studioV2WorkspaceDevices';
import { studioV2WorkspaceWriteGroupsAPI } from '../../../src/services/studioV2WorkspaceWriteGroups';
import type { StudioLocatorEntryRead } from '../../../src/types/studioV2LocatorEntry';
import type { StudioV2Workspace } from '../../../src/types/studioV2Workspace';
import type { StudioV2WorkspaceDeviceRecord } from '../../../src/services/studioV2WorkspaceDevices';

vi.mock('../../../src/services/studioV2Workspace', () => ({ studioV2WorkspaceAPI: { get: vi.fn() } }));
vi.mock('../../../src/services/studioV2WorkspaceDevices', () => ({ studioV2WorkspaceDevicesAPI: { list: vi.fn() } }));
function workspaceFixture(patch: Partial<StudioV2Workspace> = {}): StudioV2Workspace { return { id: 'workspace-1', kind: 'single', status: 'ready', ordered_device_ids: ['dev-A'], created_at: '', updated_at: '', ...patch }; }
function deviceFixture(): StudioV2WorkspaceDeviceRecord { return { id: 'dev-A', name: 'PLC', description: '', protocol: 'modbus_tcp', connection_config: '{}', status: 'draft', last_test_at: null, last_test_success: null, last_test_error: '', created_at: '', updated_at: '' }; }

vi.mock('../../../src/services/studioV2WorkspaceWriteGroups', () => ({ studioV2WorkspaceWriteGroupsAPI: { list: vi.fn() } }));
const workspaces = vi.mocked(studioV2WorkspaceAPI), devices = vi.mocked(studioV2WorkspaceDevicesAPI), groups = vi.mocked(studioV2WorkspaceWriteGroupsAPI);
function entry(revision = 'db-revision-1'): StudioLocatorEntryRead {
  return { workspace: workspaceFixture(), locator: { step: 4, workspace_id: 'workspace-1', device_id: null, rule_id: null, mapping_id: null, group_id: 'group-1', group_revision: 'g1' },
    entry: { workspace_revision: revision, hydration_scope: { complete_workspace: false, device_ids: [], rule_ids: [], mapping_ids: [], group_ids: ['group-1'] }, rule_heads: {}, group_heads: { 'group-1': 'g1' }, revalidated_mapping_ids: [] },
    selection: { state: 'selected', target: { workspaceId: 'workspace-1', workspaceRevision: revision, step: 4, groupId: 'group-1', groupRevision: 'g1' } } };
}
function setup(value = entry()) { return renderHook(({ read, enabled, workspaceId }) => useLocatorDeviceSourceLease(read, workspaceId, enabled), { initialProps: { read: value, enabled: true, workspaceId: 'workspace-1' } }); }
function deferred<T>() { let resolve!: (value: T) => void; const promise = new Promise<T>(yes => { resolve = yes; }); return { resolve, promise }; }
beforeEach(() => {
  vi.clearAllMocks();
  workspaces.get.mockReset().mockResolvedValue(workspaceFixture());
  devices.list.mockReset().mockResolvedValue([deviceFixture()]);
  groups.list.mockReset().mockResolvedValue({ workspace_id: 'workspace-1', workspace_revision: 'db-revision-1', groups: [] });
});
describe('explicit locator device source read lease, not CAS', () => {
  it('does no reads until explicit acquisition and brackets complete devices with matching DB revision reads', async () => {
    const current = setup(); expect(current.result.current.allowed).toBe(false); expect(workspaces.get).not.toHaveBeenCalled();
    await act(async () => { expect(await current.result.current.acquire('workspace-1')).toBe(true); });
    expect(current.result.current.allowed).toBe(true);
    expect(workspaces.get).toHaveBeenCalledTimes(2); expect(groups.list).toHaveBeenCalledTimes(2); expect(devices.list).toHaveBeenCalledOnce();
    expect(workspaces.get.mock.invocationCallOrder[0]).toBeLessThan(devices.list.mock.invocationCallOrder[0]);
    expect(workspaces.get.mock.invocationCallOrder[1]).toBeGreaterThan(devices.list.mock.invocationCallOrder[0]);
  });
  it.each(['foreign', 'stale', 'empty', 'changed', 'missing-device', 'duplicate-device'] as const)('fails closed for %s evidence without returning candidate authority', async mode => {
    const current = setup();
    if (mode === 'foreign') workspaces.get.mockResolvedValue(workspaceFixture({ id: 'other' }));
    if (mode === 'stale') groups.list.mockResolvedValue({ workspace_id: 'workspace-1', workspace_revision: 'stale', groups: [] });
    if (mode === 'empty') groups.list.mockResolvedValue({ workspace_id: 'workspace-1', workspace_revision: '', groups: [] });
    if (mode === 'changed') groups.list.mockResolvedValueOnce({ workspace_id: 'workspace-1', workspace_revision: 'db-revision-1', groups: [] }).mockResolvedValue({ workspace_id: 'workspace-1', workspace_revision: 'db-revision-2', groups: [] });
    if (mode === 'missing-device') devices.list.mockResolvedValue([]);
    if (mode === 'duplicate-device') devices.list.mockResolvedValue([deviceFixture(), deviceFixture()]);
    await act(async () => { expect(await current.result.current.acquire('workspace-1')).toBe(false); });
    expect(current.result.current.allowed).toBe(false);
    expect(await current.result.current.readback()).toBeUndefined();
  });
  it('allows a stable postmutation revision advance, but not a mid-read change', async () => {
    const current = setup(); await act(async () => { expect(await current.result.current.acquire('workspace-1')).toBe(true); });
    groups.list.mockResolvedValue({ workspace_id: 'workspace-1', workspace_revision: 'after-own-create', groups: [] });
    await act(async () => { expect(await current.result.current.readback()).toEqual([deviceFixture()]); });
    groups.list.mockResolvedValueOnce({ workspace_id: 'workspace-1', workspace_revision: 'after-own-create', groups: [] }).mockResolvedValue({ workspace_id: 'workspace-1', workspace_revision: 'concurrent-change', groups: [] });
    await act(async () => { expect(await current.result.current.readback()).toBeUndefined(); });
    expect(current.result.current.allowed).toBe(false);
  });
  it('rejects late acquisition after entry replacement, suspension, and A to B to A', async () => {
    const read = entry(); const current = setup(read); const held = deferred<ReturnType<typeof deviceFixture>[]>(); devices.list.mockReturnValueOnce(held.promise);
    let pending!: Promise<boolean>; act(() => { pending = current.result.current.acquire('workspace-1'); });
    await act(async () => { await Promise.resolve(); await Promise.resolve(); });
    current.rerender({ read: entry(), enabled: true, workspaceId: 'workspace-2' });
    current.rerender({ read, enabled: false, workspaceId: 'workspace-1' });
    current.rerender({ read, enabled: true, workspaceId: 'workspace-1' });
    await act(async () => { held.resolve([deviceFixture()]); expect(await pending).toBe(false); });
    expect(current.result.current.allowed).toBe(false);
  });
  it('does not acquire for foreign requested workspace or suspended locator', async () => {
    const current = setup(); await act(async () => { expect(await current.result.current.acquire('other')).toBe(false); });
    current.rerender({ read: entry(), enabled: false, workspaceId: 'workspace-1' });
    await act(async () => { expect(await current.result.current.acquire('workspace-1')).toBe(false); });
    expect(workspaces.get).not.toHaveBeenCalled();
  });
  it('invalidates a previously granted lease when canonical readback fails', async () => {
    const current = setup(); await act(async () => { expect(await current.result.current.acquire('workspace-1')).toBe(true); });
    devices.list.mockRejectedValue(new Error('unavailable'));
    await act(async () => { expect(await current.result.current.readback()).toBeUndefined(); });
    expect(current.result.current.allowed).toBe(false);
  });
});

describe('same-owner read-only recovery', () => {
  it.each(['mismatch', 'failure'])('recovers after %s without reacquiring the obsolete locator revision', async mode => {
    const read = entry(), current = setup(read);
    await act(async () => { expect(await current.result.current.acquire('workspace-1')).toBe(true); });
    groups.list.mockResolvedValue({ workspace_id: 'workspace-1', workspace_revision: 'after-create', groups: [] });
    if (mode === 'mismatch') workspaces.get.mockResolvedValueOnce(workspaceFixture({ id: 'foreign' }));
    else devices.list.mockRejectedValueOnce(new Error('offline'));
    await act(async () => { expect(await current.result.current.readback()).toBeUndefined(); });
    expect(current.result.current.allowed).toBe(false);
    await act(async () => { expect(await current.result.current.readback()).toEqual([deviceFixture()]); });
    expect(current.result.current.allowed).toBe(true);
    current.rerender({ read, enabled: true, workspaceId: 'foreign' });
    expect(await current.result.current.readback()).toBeUndefined();
    current.rerender({ read, enabled: true, workspaceId: 'workspace-1' });
    expect(await current.result.current.readback()).toBeUndefined();
    expect(current.result.current.allowed).toBe(false);
  });
});
