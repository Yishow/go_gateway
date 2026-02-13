import type { DataType } from '../../types/datalink';

export const SOURCE_TEMPLATE_STORAGE_KEY = 'pipeline-studio-source-templates-v1';
export const SOURCE_TEMPLATE_SCHEMA_VERSION = 2;

export interface SourceTemplateRecord {
  id: string;
  name: string;
  dataType: DataType;
  count: number;
  startAddress: string;
  updatedAt: string;
  lastUsedAt: string;
  version: number;
}

export function loadSourceTemplates(storage: Storage = window.localStorage): SourceTemplateRecord[] {
  try {
    const raw = storage.getItem(SOURCE_TEMPLATE_STORAGE_KEY);
    if (!raw) return [];
    const parsed = JSON.parse(raw);
    if (!Array.isArray(parsed)) return [];
    return parsed
      .map((item) => normalizeTemplateRecord(item))
      .filter((item): item is SourceTemplateRecord => item !== null);
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

export function isTemplateStale(template: SourceTemplateRecord): boolean {
  return template.version < SOURCE_TEMPLATE_SCHEMA_VERSION;
}

export function upgradeTemplates(templates: SourceTemplateRecord[]): SourceTemplateRecord[] {
  const now = new Date().toISOString();
  return templates.map((template) => ({
    ...template,
    version: SOURCE_TEMPLATE_SCHEMA_VERSION,
    updatedAt: template.updatedAt || now,
    lastUsedAt: template.lastUsedAt || template.updatedAt || now,
  }));
}

function normalizeTemplateRecord(raw: unknown): SourceTemplateRecord | null {
  if (!raw || typeof raw !== 'object') return null;
  const rec = raw as Partial<SourceTemplateRecord>;

  if (
    typeof rec.id !== 'string' ||
    typeof rec.name !== 'string' ||
    typeof rec.dataType !== 'string' ||
    typeof rec.count !== 'number' ||
    typeof rec.startAddress !== 'string' ||
    typeof rec.updatedAt !== 'string'
  ) {
    return null;
  }

  const inferredLastUsedAt =
    typeof rec.lastUsedAt === 'string' && rec.lastUsedAt ? rec.lastUsedAt : rec.updatedAt;
  const inferredVersion =
    typeof rec.version === 'number' && Number.isFinite(rec.version) ? rec.version : 1;

  return {
    id: rec.id,
    name: rec.name,
    dataType: rec.dataType as DataType,
    count: rec.count,
    startAddress: rec.startAddress,
    updatedAt: rec.updatedAt,
    lastUsedAt: inferredLastUsedAt,
    version: inferredVersion,
  };
}
