import type { DataType } from '../../types/datalink';

export const SOURCE_TEMPLATE_STORAGE_KEY = 'pipeline-studio-source-templates-v1';

export interface SourceTemplateRecord {
  id: string;
  name: string;
  dataType: DataType;
  count: number;
  startAddress: string;
  updatedAt: string;
}

export function loadSourceTemplates(storage: Storage = window.localStorage): SourceTemplateRecord[] {
  try {
    const raw = storage.getItem(SOURCE_TEMPLATE_STORAGE_KEY);
    if (!raw) return [];
    const parsed = JSON.parse(raw);
    return Array.isArray(parsed) ? (parsed as SourceTemplateRecord[]) : [];
  } catch {
    return [];
  }
}

export function saveSourceTemplates(
  templates: SourceTemplateRecord[],
  storage: Storage = window.localStorage
): void {
  storage.setItem(SOURCE_TEMPLATE_STORAGE_KEY, JSON.stringify(templates));
}
