import { DATALINK_DATA_TYPES, type DataType, type ProtocolType } from '../../../../types/datalink';
import type { MixedItem, SemanticKind } from '../../../../types/measurement';
import { addressParser, getDefaultPlannerStartAddress } from '../../../../utils/addressParser';
import { randomProfileLimit, readAcquisitionDeclaration } from './acquisitionSettings';

export const SCATTERED_ITEM_DATA_TYPES: readonly Exclude<DataType, 'string'>[] = [
  'bool', 'int16', 'uint16', 'int32', 'uint32', 'float32', 'int64', 'uint64', 'float64',
];

export function asScatteredDataType(value: string | undefined, fallback: Exclude<DataType, 'string'> = 'int16'): Exclude<DataType, 'string'> {
  if (value && value !== 'string' && (DATALINK_DATA_TYPES as readonly string[]).includes(value)) {
    return value as Exclude<DataType, 'string'>;
  }
  return fallback;
}

export type ScatteredDisabledReason = 'no_device' | 'protocol_not_supported' | 'profile_not_enabled';

export interface ScatteredRangeSupport {
  supported: boolean;
  effectiveLimit: number;
  profile: string;
  profileLimit: number;
  disabledReason?: ScatteredDisabledReason;
}

export interface ScatteredDeviceLike {
  protocol: ProtocolType;
  config?: Record<string, unknown>;
  connection_config?: string | Record<string, unknown>;
}

function extractDeviceConfig(device: ScatteredDeviceLike): Record<string, unknown> {
  if (device.config && typeof device.config === 'object') return device.config;
  if (typeof device.connection_config === 'string' && device.connection_config.trim() !== '') {
    try {
      const parsed: unknown = JSON.parse(device.connection_config);
      if (parsed && typeof parsed === 'object' && !Array.isArray(parsed)) return parsed as Record<string, unknown>;
    } catch {
      return {};
    }
  }
  if (device.connection_config && typeof device.connection_config === 'object') return device.connection_config;
  return {};
}

export function resolveScatteredRangeSupport(device?: ScatteredDeviceLike | null): ScatteredRangeSupport {
  if (!device) {
    return { supported: false, effectiveLimit: 0, profile: 'undeclared', profileLimit: 0, disabledReason: 'no_device' };
  }
  if (device.protocol !== 'mc_3e' && device.protocol !== 'fatek_fbs') {
    return { supported: false, effectiveLimit: 0, profile: 'undeclared', profileLimit: 0, disabledReason: 'protocol_not_supported' };
  }
  const declared = readAcquisitionDeclaration(extractDeviceConfig(device));
  const profileLimit = randomProfileLimit(declared.profile);
  if (profileLimit <= 0) {
    return { supported: false, effectiveLimit: 0, profile: declared.profile, profileLimit: 0, disabledReason: 'profile_not_enabled' };
  }
  const effectiveLimit =
    declared.randomPoints !== undefined &&
    Number.isInteger(declared.randomPoints) &&
    declared.randomPoints >= 1 &&
    declared.randomPoints <= profileLimit
      ? declared.randomPoints
      : profileLimit;

  return { supported: true, effectiveLimit, profile: declared.profile, profileLimit };
}

export interface ScatteredRandomUsage {
  usedPoints: number;
  effectiveLimit: number;
  ratio: number;
  isWarning: boolean;
  isOverLimit: boolean;
  status: 'normal' | 'warning' | 'over_limit';
}

function itemWordWidth(dataType: string): number {
  switch (dataType) {
    case 'int32':
    case 'uint32':
    case 'float32':
      return 2;
    case 'int64':
    case 'uint64':
    case 'float64':
      return 4;
    default:
      return 1;
  }
}

export function calculateScatteredRandomUsage(
  protocol: ProtocolType | undefined,
  items: ReadonlyArray<Pick<MixedItem, 'address' | 'data_type'>>,
  effectiveLimit: number,
): ScatteredRandomUsage {
  let usedPoints = 0;

  if (protocol === 'mc_3e') {
    let unparsedPoints = 0;
    const units = new Set<string>();
    for (const item of items) {
      const is64Bit =
        item.data_type === 'int64' || item.data_type === 'uint64' || item.data_type === 'float64';
      const is32Bit =
        item.data_type === 'int32' || item.data_type === 'uint32' || item.data_type === 'float32';
      const addr = (item.address ?? '').trim().toUpperCase();
      if (addr !== '') {
        try {
          const parsed = addressParser.parse(addr, 'mc_3e');
          if (parsed.type === 'bit') {
            const wordAddr = Math.floor(parsed.startNumber / 16) * 16;
            units.add(`${parsed.area}:${wordAddr}:16`);
            continue;
          }
          if (is64Bit) {
            for (let n = 0; n < 4; n += 1) {
              units.add(`${parsed.area}:${parsed.startNumber + n}:16`);
            }
          } else if (is32Bit) {
            units.add(`${parsed.area}:${parsed.startNumber}:32`);
          } else {
            units.add(`${parsed.area}:${parsed.startNumber}:16`);
          }
          continue;
        } catch {
          // Fall through to unparsed count when address is being typed
        }
      }
      unparsedPoints += is64Bit ? 4 : 1;
    }
    usedPoints = units.size + unparsedPoints;
  } else {
    for (const item of items) {
      usedPoints += itemWordWidth(item.data_type);
    }
  }

  const ratio = effectiveLimit > 0 ? usedPoints / effectiveLimit : 0;
  const isOverLimit = effectiveLimit > 0 && usedPoints > effectiveLimit;
  const isWarning = effectiveLimit > 0 && ratio >= 0.9;
  const status = isOverLimit ? 'over_limit' : isWarning ? 'warning' : 'normal';

  return {
    usedPoints,
    effectiveLimit,
    ratio,
    isWarning,
    isOverLimit,
    status,
  };
}

export function defaultSemanticKindForType(dataType: string): SemanticKind {
  return dataType === 'bool' ? 'state' : 'gauge';
}

export function toPayloadMixedItems(items: ReadonlyArray<MixedItem>): MixedItem[] {
  const seenIds = new Set<string>();
  return items.map((item, index) => {
    let itemId = (item.item_id ?? '').trim() || `item-${index + 1}`;
    if (seenIds.has(itemId)) {
      itemId = `${itemId}-${index + 1}`;
    }
    seenIds.add(itemId);
    const dataType = asScatteredDataType(item.data_type);
    const semanticKind: SemanticKind = item.semantic_kind || defaultSemanticKindForType(dataType);
    return {
      item_id: itemId,
      name: (item.name ?? '').trim(),
      address: (item.address ?? '').trim().toUpperCase(),
      register_offset: 0,
      data_type: dataType,
      ...(item.data_format ? { data_format: item.data_format } : {}),
      ...(item.quantity !== undefined ? { quantity: item.quantity } : {}),
      ...(item.unit !== undefined ? { unit: item.unit } : {}),
      semantic_kind: semanticKind,
    };
  });
}

export function createDefaultScatteredItem(
  index: number,
  existingItems: ReadonlyArray<MixedItem> = [],
  protocol: ProtocolType = 'mc_3e',
  fallbackAddress?: string,
  fallbackDataType: DataType = 'int16',
): MixedItem {
  const safeType = asScatteredDataType(fallbackDataType);
  let nextAddress = (fallbackAddress ?? '').trim().toUpperCase() || getDefaultPlannerStartAddress(protocol);
  const lastItem = existingItems[existingItems.length - 1];
  const lastAddress = (lastItem?.address ?? '').trim().toUpperCase();
  if (lastItem && lastAddress !== '') {
    try {
      const width = itemWordWidth(lastItem.data_type);
      nextAddress = addressParser.offset(lastAddress, width, protocol);
    } catch {
      nextAddress = lastAddress;
    }
  }
  const existingIds = new Set(existingItems.map((item) => item.item_id));
  let seq = index + 1;
  let itemId = `item-${seq}`;
  while (existingIds.has(itemId)) {
    seq += 1;
    itemId = `item-${seq}`;
  }
  return {
    item_id: itemId,
    name: '',
    address: nextAddress,
    register_offset: 0,
    data_type: safeType,
    semantic_kind: defaultSemanticKindForType(safeType),
  };
}

export type ScatteredItemIssueCode =
  | 'empty_address'
  | 'invalid_address'
  | 'unsupported_string'
  | 'bit_requires_bool'
  | 'overlap';

export interface ScatteredValidationResult {
  valid: boolean;
  itemIssues: Record<string, { code: ScatteredItemIssueCode; conflictWithAddress?: string }>;
  overlapPairs: Array<[string, string]>;
}

export function validateScatteredItems(
  protocol: ProtocolType | undefined,
  items: ReadonlyArray<MixedItem>,
): ScatteredValidationResult {
  if (!Array.isArray(items) || items.length === 0) {
    return { valid: false, itemIssues: {}, overlapPairs: [] };
  }

  const itemIssues: Record<string, { code: ScatteredItemIssueCode; conflictWithAddress?: string }> = {};
  const overlapPairs: Array<[string, string]> = [];
  const spans: Array<{ itemId: string; address: string; area: string; start: number; end: number }> = [];

  items.forEach((item, index) => {
    const itemId = (item.item_id ?? '').trim() || `item-${index + 1}`;
    const addr = (item.address ?? '').trim().toUpperCase();
    if (!addr) {
      itemIssues[itemId] = { code: 'empty_address' };
      return;
    }
    if (item.data_type === 'string') {
      itemIssues[itemId] = { code: 'unsupported_string' };
      return;
    }
    if (!protocol) {
      return;
    }
    if (!addressParser.validate(addr, protocol).valid) {
      itemIssues[itemId] = { code: 'invalid_address' };
      return;
    }
    try {
      const parsed = addressParser.parse(addr, protocol);
      if (parsed.type === 'bit' && item.data_type !== 'bool') {
        itemIssues[itemId] = { code: 'bit_requires_bool' };
        return;
      }
      const width = parsed.type === 'bit' ? 1 : itemWordWidth(item.data_type);
      const start = parsed.startNumber;
      const end = start + width;
      if (protocol === 'mc_3e' && width > 1) {
        addressParser.offset(addr, width - 1, protocol);
      }
      for (const prev of spans) {
        if (prev.area === parsed.area && prev.start < end && start < prev.end) {
          overlapPairs.push([prev.itemId, itemId]);
          if (!itemIssues[prev.itemId]) {
            itemIssues[prev.itemId] = { code: 'overlap', conflictWithAddress: addr };
          }
          itemIssues[itemId] = { code: 'overlap', conflictWithAddress: prev.address };
        }
      }
      spans.push({ itemId, address: addr, area: parsed.area, start, end });
    } catch {
      itemIssues[itemId] = { code: 'invalid_address' };
    }
  });

  return {
    valid: Object.keys(itemIssues).length === 0,
    itemIssues,
    overlapPairs,
  };
}
