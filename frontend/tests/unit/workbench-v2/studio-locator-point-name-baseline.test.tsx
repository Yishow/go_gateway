import { act, renderHook } from '@testing-library/react';
import { expect, it, vi } from 'vitest';
import { studioV2LocatorHydration, useStudioV2LocatorHydration } from '@/pages/datalink/workbench-v2/studioV2LocatorHydration';
import { INITIAL_STATE, type useWorkbenchV2State } from '@/features/datalink/workbench-v2/state/useWorkbenchV2State';
import { hydrateStudioV2Mapping, toStudioV2MappingUpdateRequest } from '@/features/datalink/workbench-v2/state/studioV2MappingAutosave';
import type { WorkbenchV2State } from '@/features/datalink/workbench-v2/state/types';
import { locatorRead, nativeMapping } from './studio-locator-fixtures';
function harness() {
  const entry = locatorRead(4), ref = { current: structuredClone(INITIAL_STATE) as WorkbenchV2State };
  const actions: ReturnType<typeof useWorkbenchV2State> = { state: INITIAL_STATE, dispatch: vi.fn(), setView: vi.fn(), setCurrent: vi.fn(), completeStep: vi.fn(),
    toggleSidebar: vi.fn(), toggleSummaryRail: vi.fn(), setSidebarCollapsed: vi.fn(), setShowSummaryRail: vi.fn(), resetFlow: vi.fn(), selectRule: vi.fn() };
  const hook = renderHook(({ selected }) => useStudioV2LocatorHydration(actions, ref, false, selected), { initialProps: { selected: entry } });
  const projected = studioV2LocatorHydration(locatorRead(3))!;
  const point = projected.points[0], mapping = projected.mappings[point.id];
  const proof = { workspaceId: 'W', workspaceRevision: 'W2', identity: { device_id: 'D', point_id: 'P-native', tag_id: 'T' }, mappingId: 'M', ruleId: 'R', ruleRevision: 'R2', point, mapping };
  const merge = { devices: projected.devices, rules: projected.rules, points: projected.points, mappings: projected.mappings };
  const read = () => {
    const accept = hook.result.current.capturePointNameRead(ref.current, proof, merge);
    ref.current = { ...ref.current, ...merge }; act(() => accept?.()); hook.rerender({ selected: entry }); return accept;
  };
  return { ...hook, entry, ref, proof, merge, read };
}
it('accepts only the fresh missing source read and preserves a later unrelated draft', () => {
  const h = harness(); expect(h.read()).toBeTypeOf('function'); expect(h.result.current.canLeave).toBe(true);
  h.ref.current = { ...h.ref.current, db: { ...h.ref.current.db, connector: { ...h.ref.current.db.connector, name: 'Unsent' } } };
  h.rerender({ selected: h.entry }); expect(h.result.current.canLeave).toBe(false);
});
it.each(['foreign-revision', 'preexisting-dirty', 'replaced-row', 'late-entry'] as const)('does not clean a read for %s', reason => {
  const h = harness();
  if (reason === 'preexisting-dirty') h.ref.current = { ...h.ref.current, db: { ...h.ref.current.db, connector: { ...h.ref.current.db.connector, name: 'Unsent' } } };
  if (reason === 'replaced-row') h.ref.current = { ...h.ref.current, devices: [{ ...h.merge.devices[0], name: 'Unsent D' }] };
  const accept = h.result.current.capturePointNameRead(h.ref.current, { ...h.proof, workspaceRevision: reason === 'foreign-revision' ? 'FOREIGN' : 'W2' }, h.merge);
  if (reason === 'late-entry') h.rerender({ selected: locatorRead(4, 'B') });
  h.ref.current = { ...h.ref.current, ...h.merge }; act(() => accept?.()); h.rerender({ selected: reason === 'late-entry' ? locatorRead(4, 'B') : h.entry });
  expect(h.result.current.canLeave).toBe(false);
});
it.each(['own', 'foreign', 'new-name', 'other-dirty', 'late-entry'] as const)('keeps own-name ACK separate from read acceptance: %s', reason => {
  const h = harness(); h.read();
  const submitted = { ...toStudioV2MappingUpdateRequest(h.proof.point, h.proof.mapping), display_name: 'Confirmed name' };
  const accept = h.result.current.capturePointNameSave(h.proof, submitted);
  expect(accept).toBeTypeOf('function');
  const saved = { ...nativeMapping(), ...submitted, runtime_apply_status: 'not_running' as const };
  const latest = hydrateStudioV2Mapping(h.proof.point, saved);
  h.ref.current = { ...h.ref.current, mappings: { ...h.ref.current.mappings, [h.proof.point.id]: reason === 'new-name' ? { ...latest, display_name: 'Newer unsent' } : latest } };
  h.rerender({ selected: h.entry }); expect(h.result.current.canLeave).toBe(false);
  if (reason === 'other-dirty') h.ref.current = { ...h.ref.current, db: { ...h.ref.current.db, connector: { ...h.ref.current.db.connector, name: 'Unsent DB' } } };
  if (reason === 'late-entry') h.rerender({ selected: locatorRead(4, 'B') });
  act(() => accept!({ ...saved, ...(reason === 'foreign' ? { id: 'FOREIGN' } : {}) }));
  h.rerender({ selected: reason === 'late-entry' ? locatorRead(4, 'B') : h.entry });
  expect(h.result.current.canLeave).toBe(reason === 'own');
});
