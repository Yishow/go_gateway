import { describe, expect, it, beforeEach } from 'vitest';
import {
  SOURCE_TEMPLATE_STORAGE_KEY,
  loadSourceTemplates,
  saveSourceTemplates,
  type SourceTemplateRecord,
} from '../sourceTemplateStorage';

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
      },
      {
        id: 'tpl-2',
        name: 'FLOAT batch',
        dataType: 'float32',
        count: 10,
        startAddress: '40101',
        updatedAt: '2026-02-13T00:00:01.000Z',
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
});
