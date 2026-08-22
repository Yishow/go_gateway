import { describe, expect, it } from 'vitest';
import {
  INSPECTOR_SELECTION_NONE,
  deviceStatusToReadiness,
  tagStatusToReadiness,
  type InspectorSelection,
} from '@/pages/datalink/workbench/workbenchTypes';

describe('deviceStatusToReadiness', () => {
  it.each([
    ['active', 'ready'],
    ['draft', 'draft'],
    ['disabled', 'blocked'],
  ] as const)('maps DeviceStatus "%s" → WorkbenchReadiness "%s"', (status, expected) => {
    expect(deviceStatusToReadiness(status)).toBe(expected);
  });
});

describe('tagStatusToReadiness', () => {
  it.each([
    ['active', true, 'ready'],
    ['active', false, 'partial'],
    ['draft', true, 'draft'],
    ['draft', false, 'draft'],
    ['retired', true, 'blocked'],
    ['retired', false, 'blocked'],
  ] as const)(
    'maps TagStatus "%s" (hasOutputMapping=%s) → WorkbenchReadiness "%s"',
    (status, hasMapping, expected) => {
      expect(tagStatusToReadiness(status, hasMapping)).toBe(expected);
    },
  );
});

describe('InspectorSelection', () => {
  it('INSPECTOR_SELECTION_NONE is kind "none"', () => {
    expect(INSPECTOR_SELECTION_NONE.kind).toBe('none');
  });

  it('discriminated union covers all kinds', () => {
    const selections: InspectorSelection[] = [
      { kind: 'none' },
      { kind: 'device', deviceId: 'd1' },
      { kind: 'rule', ruleId: 'r1' },
      { kind: 'span', spanAddress: '40001', ruleId: 'r1' },
      { kind: 'tag', tagId: 't1', pointId: 'p1' },
      { kind: 'outputCandidate', tagId: 't1', target: 'modbus' },
    ];

    const kinds = selections.map((s) => s.kind);
    expect(kinds).toEqual(['none', 'device', 'rule', 'span', 'tag', 'outputCandidate']);
  });

  it('span selection can have optional ruleId', () => {
    const withRule: InspectorSelection = { kind: 'span', spanAddress: '40001', ruleId: 'r1' };
    const withoutRule: InspectorSelection = { kind: 'span', spanAddress: '40001' };

    expect(withRule.kind).toBe('span');
    expect(withoutRule.kind).toBe('span');
    if (withRule.kind === 'span') {
      expect(withRule.ruleId).toBe('r1');
    }
    if (withoutRule.kind === 'span') {
      expect(withoutRule.ruleId).toBeUndefined();
    }
  });

  it('outputCandidate selection can have optional target', () => {
    const withTarget: InspectorSelection = {
      kind: 'outputCandidate',
      tagId: 't1',
      target: 'database',
    };
    const withoutTarget: InspectorSelection = {
      kind: 'outputCandidate',
      tagId: 't1',
    };

    if (withTarget.kind === 'outputCandidate') {
      expect(withTarget.target).toBe('database');
    }
    if (withoutTarget.kind === 'outputCandidate') {
      expect(withoutTarget.target).toBeUndefined();
    }
  });
});
