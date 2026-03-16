import { beforeEach, describe, expect, it } from 'vitest';
import {
  loadSourceTemplates,
  saveSourceTemplates,
} from '@/features/datalink/sourceTemplateStorage';
import {
  applyTemplateToPlanner,
  createTemplateFromPlanner,
  isTemplateRecordContractValid,
  normalizeNamingPrefix,
  upsertTemplateRecord,
} from '@/features/datalink/sourcePlannerContract';

describe('sourcePlannerContract', () => {
  beforeEach(() => {
    localStorage.clear();
  });

  it('normalizes batch naming prefix by rule', () => {
    expect(normalizeNamingPrefix(' line a temp ')).toBe('LINE-A-TEMP');
    expect(normalizeNamingPrefix('$$$')).toBe('SRC');
  });

  it('creates template from planner and applies template back to planner defaults', () => {
    const template = createTemplateFromPlanner({
      templateName: 'Line A Float',
      draft: {
        dataType: 'float32',
        count: 10,
        startAddress: '40001',
      },
      preferredViewMode: 'live',
      now: '2026-02-13T00:00:00.000Z',
    });

    expect(template.id).toBe('line-a-float');
    expect(isTemplateRecordContractValid(template)).toBe(true);
    expect(template.preferredViewMode).toBe('live');
    expect(applyTemplateToPlanner(template)).toEqual({
      dataType: 'float32',
      count: 10,
      startAddress: '40001',
    });
  });

  it('supports planner-template persistence round trip', () => {
    const created = createTemplateFromPlanner({
      templateName: 'Line A Int',
      draft: {
        dataType: 'int16',
        count: 5,
        startAddress: '40010',
      },
      preferredViewMode: 'plan',
      now: '2026-02-13T00:00:00.000Z',
    });

    const records = upsertTemplateRecord([], created);
    saveSourceTemplates(records);
    const restored = loadSourceTemplates();

    expect(restored).toHaveLength(1);
    expect(restored[0]).toEqual(created);
    expect(applyTemplateToPlanner(restored[0])).toEqual({
      dataType: 'int16',
      count: 5,
      startAddress: '40010',
    });
  });
});
