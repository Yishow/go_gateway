import { describe, expect, it, vi } from 'vitest';
import { recordingHomeGuidance } from '@/features/datalink/workbench-v2/recording-home/recordingHomeGuidance';
import { fixture, group } from './recording-home-guidance-fixtures';

describe('one home guidance resolver for the five entry fixtures', () => {
  it.each([
    [fixture(), 'missing_devices', 'connect_device'],
    [fixture(['D1']), 'missing_points', 'configure_points'],
    [fixture(['D1'], true), 'missing_recordings', 'choose_sources'],
    [fixture(['D1'], true, [group('C', 'draft')]), 'draft', 'recording'],
    [fixture(['D1'], true, [group('B')]), 'recordings', 'recording'],
  ] as const)('projects the canonical fixture into %s', (input, reason, target) => {
    const result = recordingHomeGuidance(input);
    expect(result.availability).toEqual({ state: 'available' });
    expect(result.reason).toBe(reason);
    expect(result.primaryAction?.kind).toBe('navigate');
    expect(result.primaryAction?.target).toBe(target);
  });
  it('configures the existing device instead of reconnecting it', () => {
    expect(recordingHomeGuidance(fixture(['D1'])).primaryAction).toEqual({ kind: 'navigate',
      scope: { workspaceId: 'W1', workspaceRevision: 'W2', epoch: 2 }, target: 'configure_points', deviceId: 'D1' });
  });
});

describe('projection purity and deterministic selection', () => {
  it('does not freeze or mutate query-owner input and detaches/freeze its output', () => {
    const input = fixture(['D1'], true, [group('B')]);
    const before = JSON.stringify(input);
    const result = recordingHomeGuidance(input);
    expect(JSON.stringify(input)).toBe(before);
    expect(Object.isFrozen(input)).toBe(false);
    expect(Object.isFrozen(result)).toBe(true);
    expect(Object.isFrozen(result.recordingItems)).toBe(true);
    expect(Object.isFrozen(result.primaryAction)).toBe(true);
    expect(result.primaryAction?.scope).not.toBe(input.evidence.scope);
  });
  it('orders draft continuations by identity, not response order or display name', () => {
    const z = group('Z', 'draft'); z.name = 'AAA';
    const a = group('A', 'draft'); a.name = 'ZZZ'; a.members = [];
    const result = recordingHomeGuidance(fixture(['D1'], true, [z, a]));
    expect(result.primaryAction).toMatchObject({ groupId: 'A', section: 'sources' });
    expect(result.draftItems.map(item => item.id)).toEqual(['A', 'Z']);
  });
  it('selects a deterministic existing device when several lack points', () => {
    expect(recordingHomeGuidance(fixture(['Z', 'A'])).primaryAction).toMatchObject({ target: 'configure_points', deviceId: 'A' });
  });
});


describe('read-only guidance boundary', () => {
  it('accepts deeply frozen owner input and never calls the network while selecting guidance', () => {
    const input = fixture(['D1'], true, [group('A')]);
    const freeze = (value: unknown): void => {
      if (value !== null && typeof value === 'object') { Object.values(value).forEach(freeze); Object.freeze(value); }
    };
    freeze(input);
    const network = vi.spyOn(globalThis, 'fetch').mockImplementation(() => { throw new Error('navigation must not request a mutation'); });
    try {
      expect(recordingHomeGuidance(input).primaryAction).toMatchObject({ kind: 'navigate', target: 'recording' });
      expect(network).not.toHaveBeenCalled();
    } finally { network.mockRestore(); }
  });
});
