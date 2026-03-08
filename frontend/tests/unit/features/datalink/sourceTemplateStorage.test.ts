import { describe, expect, it, beforeEach } from 'vitest';
import {
  SOURCE_TEMPLATE_STORAGE_KEY,
  SOURCE_TEMPLATE_SCHEMA_VERSION,
  isTemplateStale,
  loadSourceTemplates,
  saveSourceTemplates,
  upgradeTemplates,
  type SourceTemplateRecord,
} from '@/features/datalink/sourceTemplateStorage';

describe('sourceTemplateStorage', () => {
  beforeEach(() => {
    localStorage.clear();
  });

  it('returns empty array when storage is empty', () => {
    expect(loadSourceTemplates()).toEqual([]);
  });

  it('saves and loads templates for restore-after-reload flow', () => {
    const templates: SourceTemplateRecord[] = [
      {
        id: 'tpl-1',
        name: 'INT batch',
        dataType: 'int16',
        count: 5,
        startAddress: '40001',
        updatedAt: '2026-02-13T00:00:00.000Z',
        lastUsedAt: '2026-02-13T00:00:00.000Z',
        version: SOURCE_TEMPLATE_SCHEMA_VERSION,
      },
      {
        id: 'tpl-2',
        name: 'FLOAT batch',
        dataType: 'float32',
        count: 10,
        startAddress: '40101',
        updatedAt: '2026-02-13T00:00:01.000Z',
        lastUsedAt: '2026-02-13T00:00:01.000Z',
        version: SOURCE_TEMPLATE_SCHEMA_VERSION,
      },
    ];

    saveSourceTemplates(templates);
    const restored = loadSourceTemplates();
    expect(restored).toEqual(templates);
  });

  it('falls back to empty array on invalid json payload', () => {
    localStorage.setItem(SOURCE_TEMPLATE_STORAGE_KEY, '{broken-json');
    expect(loadSourceTemplates()).toEqual([]);
  });

  it('normalizes legacy template without version/lastUsedAt as stale', () => {
    localStorage.setItem(
      SOURCE_TEMPLATE_STORAGE_KEY,
      JSON.stringify([
        {
          id: 'legacy-1',
          name: 'Legacy',
          dataType: 'int16',
          count: 3,
          startAddress: '40001',
          updatedAt: '2026-02-13T00:00:00.000Z',
        },
      ])
    );

    const [legacy] = loadSourceTemplates();
    expect(legacy.lastUsedAt).toBe('2026-02-13T00:00:00.000Z');
    expect(legacy.version).toBe(1);
    expect(isTemplateStale(legacy)).toBe(true);
  });

  it('upgrades stale templates to latest schema version', () => {
    const stale: SourceTemplateRecord[] = [
      {
        id: 'stale-1',
        name: 'Stale',
        dataType: 'int16',
        count: 2,
        startAddress: '40001',
        updatedAt: '2026-02-13T00:00:00.000Z',
        lastUsedAt: '2026-02-13T00:00:00.000Z',
        version: 1,
      },
    ];

    const upgraded = upgradeTemplates(stale);
    expect(upgraded[0].version).toBe(SOURCE_TEMPLATE_SCHEMA_VERSION);
    expect(isTemplateStale(upgraded[0])).toBe(false);
  });
});
