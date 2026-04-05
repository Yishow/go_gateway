import { DATALINK_DATA_TYPES, type DataType, type ProtocolType } from '../../types/datalink';
import {
  SOURCE_TEMPLATE_SCHEMA_VERSION,
  type SourceTemplateCapabilitySnapshot,
  type SourceTemplateRecord,
  type SourceTemplateViewMode,
} from './sourceTemplateStorage';

/** 來源規劃／範本允許之資料型態（與 Point CRUD 一致）。 */
export const SOURCE_PLANNER_ALLOWED_DATA_TYPES = DATALINK_DATA_TYPES;

export type DataTypeGroupEntry = {
  value: DataType;
  supported: boolean;
  disabledReasonKey?: string;
};

export type DataTypeGroup = {
  labelKey: string;
  types: DataTypeGroupEntry[];
};

export const SOURCE_PLANNER_DATA_TYPE_GROUPS: DataTypeGroup[] = [
  {
    labelKey: 'workbench.source.planner.dataTypeGroup.singleWord',
    types: [
      { value: 'bool', supported: true },
      { value: 'int16', supported: true },
      { value: 'uint16', supported: true },
      { value: 'string', supported: true },
    ],
  },
  {
    labelKey: 'workbench.source.planner.dataTypeGroup.32bit',
    types: [
      { value: 'int32', supported: true },
      { value: 'uint32', supported: true },
      { value: 'float32', supported: true },
    ],
  },
  {
    labelKey: 'workbench.source.planner.dataTypeGroup.64bit',
    types: [
      { value: 'int64', supported: true },
      { value: 'uint64', supported: true },
      { value: 'float64', supported: true },
    ],
  },
];

export interface SourcePlannerDraft {
  dataType: DataType;
  count: number;
  startAddress: string;
  namingPrefix: string;
}

/**
 * 解析規劃器「數量」輸入字串或已同步為數字的 state。
 *
 * - 空白（含僅空白字元）視為尚未輸入完成，回傳 `null`。
 * - 僅接受有限且為整數、且 ≥ 1 的數值；否則回傳 `null`。
 *
 * @param raw 使用者於輸入框內的文字，或規劃器內以 `number` 保存的數量（與受控 `<input type="number">` 對齊）
 * @returns 有效整數數量，或表示無效／留空時的 `null`
 */
export function parsePlannerCountInput(raw: string | number): number | null {
  const trimmed = String(raw).trim();
  if (trimmed === '') {
    return null;
  }
  const n = Number(trimmed);
  if (!Number.isFinite(n) || !Number.isInteger(n) || n < 1) {
    return null;
  }
  return n;
}

export function normalizeNamingPrefix(prefix: string): string {
  const compact = prefix.trim().toUpperCase();
  const sanitized = compact.replace(/[^A-Z0-9_-]+/g, '-').replace(/-+/g, '-').replace(/^-|-$/g, '');
  return sanitized || 'SRC';
}

/**
 * 依設備協定回傳來源規劃器「命名前綴」的建議預設值（大寫縮寫）。
 * 用於切換設備或開啟規則建立器時，讓預設 Point 名稱與協定語意一致（例如 Modbus TCP／RTU／UDP → MBT，三菱 MC Protocol → MC）。
 *
 * @param protocol 設備的 {@link ProtocolType}
 * @returns 建議前綴字串（尚未經 {@link normalizeNamingPrefix} 正規化，通常已符合慣例字元集）
 */
export function getDefaultNamingPrefixForProtocol(protocol: ProtocolType): string {
  switch (protocol) {
    case 'modbus_tcp':
    case 'modbus_udp':
    case 'modbus_rtu':
      return 'MBT';
    case 'mc_3e':
      return 'MC';
    case 'fatek_fbs':
      return 'FBS';
    case 'mqtt':
      return 'MQTT';
  }
}

export function normalizeTemplateId(name: string): string {
  return name.trim().toLowerCase().replace(/\s+/g, '-').replace(/[^a-z0-9_-]/g, '');
}

export function createTemplateFromPlanner(input: {
  templateName: string;
  draft: Pick<SourcePlannerDraft, 'dataType' | 'count' | 'startAddress'>;
  preferredViewMode?: SourceTemplateViewMode;
  capabilitySnapshot?: SourceTemplateCapabilitySnapshot;
  now?: string;
}): SourceTemplateRecord {
  const now = input.now ?? new Date().toISOString();
  const normalizedName = input.templateName.trim();

  return {
    id: normalizeTemplateId(normalizedName),
    name: normalizedName,
    dataType: input.draft.dataType,
    count: input.draft.count,
    startAddress: input.draft.startAddress.trim().toUpperCase(),
    preferredViewMode: input.preferredViewMode ?? 'plan',
    capabilitySnapshot: input.capabilitySnapshot,
    updatedAt: now,
    lastUsedAt: now,
    version: SOURCE_TEMPLATE_SCHEMA_VERSION,
  };
}

export function upsertTemplateRecord(
  previous: SourceTemplateRecord[],
  next: SourceTemplateRecord,
  maxSize = 20
): SourceTemplateRecord[] {
  const existingIndex = previous.findIndex((item) => item.id === next.id);
  if (existingIndex === -1) return [next, ...previous].slice(0, maxSize);

  const copied = [...previous];
  copied[existingIndex] = next;
  return copied;
}

export function applyTemplateToPlanner(template: SourceTemplateRecord): Pick<SourcePlannerDraft, 'dataType' | 'count' | 'startAddress'> {
  return {
    dataType: template.dataType,
    count: template.count,
    startAddress: template.startAddress,
  };
}

export function isTemplateRecordContractValid(template: SourceTemplateRecord): boolean {
  if (!template.id || !template.name) return false;
  if (!Number.isInteger(template.count) || template.count <= 0 || template.count > 200) return false;
  if (!template.startAddress.trim()) return false;
  if (!Number.isFinite(template.version) || template.version < 1) return false;
  if (
    template.preferredViewMode !== undefined &&
    !['plan', 'live', 'link'].includes(template.preferredViewMode)
  ) {
    return false;
  }
  if (!(SOURCE_PLANNER_ALLOWED_DATA_TYPES as readonly DataType[]).includes(template.dataType)) {
    return false;
  }
  return true;
}
