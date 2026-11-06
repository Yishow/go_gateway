import type { CreateSourceRuleRequest, SourceRuleRecord } from '../../../../types/datalink';
import type { StudioV2UpdateSourceRuleRequest } from '../../../../services/studioV2Rules';
import type { Rule } from './types';
import { asScatteredDataType, toPayloadMixedItems } from './scatteredRange';

const VALID_DATA_FORMATS = new Set(['', 'ABCD', 'BADC', 'CDAB', 'DCBA']);

function isNonEmptyString(value: unknown): value is string {
  return typeof value === 'string' && value.trim() !== '';
}

function isNonNegativeNumber(value: unknown): value is number {
  return typeof value === 'number' && Number.isFinite(value) && value >= 0;
}

export function canHydrateStudioV2Rule(record: SourceRuleRecord | null | undefined): record is SourceRuleRecord {
  if (!record) {
    return false;
  }

  return (
    isNonEmptyString(record.id) &&
    typeof record.device_id === 'string' &&
    typeof record.start_address === 'string' &&
    isNonNegativeNumber(record.count) &&
    isNonEmptyString(record.naming_prefix) &&
    Array.isArray(record.skipped_addresses)
  );
}

export function isStudioV2RuleValid(rule: Rule): boolean {
  if (rule.layout_mode === 'mixed') {
    return Boolean(
      rule.device_id.trim() !== '' &&
      Array.isArray(rule.mixed_items) &&
      rule.mixed_items.length > 0 &&
      rule.mixed_items.every(
        (item) =>
          typeof item.address === 'string' &&
          item.address.trim() !== '' &&
          item.data_type !== 'string' &&
          VALID_DATA_FORMATS.has(item.data_format ?? ''),
      ),
    );
  }
  return Boolean(
    rule.device_id.trim() !== '' &&
    rule.start_address.trim() !== '' &&
    rule.count > 0 &&
    VALID_DATA_FORMATS.has(rule.data_format),
  );
}

export function toStudioV2RuleCreateRequest(rule: Rule): CreateSourceRuleRequest {
  if (rule.layout_mode === 'mixed') {
    const items = toPayloadMixedItems(rule.mixed_items ?? []);
    const first = items[0];
    return {
      id: rule.id,
      device_id: rule.device_id,
      start_address: first?.address ?? rule.start_address,
      count: Math.max(1, items.length),
      data_type: first ? asScatteredDataType(first.data_type) : rule.data_type,
      naming_prefix: rule.naming_prefix,
      enabled: rule.enabled,
      skipped_addresses: [],
      layout_mode: 'mixed',
      mixed_items: items,
    };
  }
  return {
    id: rule.id,
    device_id: rule.device_id,
    start_address: rule.start_address,
    count: rule.count,
    data_type: rule.data_type,
    naming_prefix: rule.naming_prefix,
    enabled: rule.enabled,
    skipped_addresses: [...rule.skipped_addresses],
    scale_multiplier: rule.scale_multiplier,
    scale_offset: rule.scale_offset,
    data_format: rule.data_format,
    share_enabled: rule.share_enabled,
    share_start_register: rule.share_start_register,
    share_stride: rule.share_stride,
  };
}

export function toStudioV2RuleUpdateRequest(rule: Rule): StudioV2UpdateSourceRuleRequest {
  if (rule.layout_mode === 'mixed') {
    return {
      device_id: rule.device_id,
      naming_prefix: rule.naming_prefix,
      enabled: rule.enabled,
      layout_mode: 'mixed',
      mixed_items: toPayloadMixedItems(rule.mixed_items ?? []),
      skipped_addresses: [],
      share_enabled: false,
    };
  }
  return {
    device_id: rule.device_id,
    layout_mode: 'homogeneous',
    start_address: rule.start_address,
    count: rule.count,
    data_type: rule.data_type,
    naming_prefix: rule.naming_prefix,
    enabled: rule.enabled,
    skipped_addresses: [...rule.skipped_addresses],
    scale_multiplier: rule.scale_multiplier,
    scale_offset: rule.scale_offset,
    data_format: rule.data_format,
    share_enabled: rule.share_enabled,
    share_start_register: rule.share_start_register,
    share_stride: rule.share_stride,
  };
}

export function hydrateStudioV2Rule(record: SourceRuleRecord): Rule {
  return {
    id: record.id,
    device_id: record.device_id,
    workspace_id: record.workspace_id,
    revision_id: record.revision_id,
    name: record.naming_prefix || record.id,
    start_address: record.start_address,
    count: record.count,
    data_type: record.data_type,
    naming_prefix: record.naming_prefix,
    enabled: record.enabled,
    scale_multiplier: record.scale_multiplier ?? 1,
    scale_offset: record.scale_offset ?? 0,
    data_format: (record.data_format ?? '') as Rule['data_format'],
    skipped_addresses: [...record.skipped_addresses],
    share_enabled: record.share_enabled ?? false,
    share_start_register: record.share_start_register ?? null,
    share_stride: record.share_stride ?? null,
    ...(record.layout_mode ? { layout_mode: record.layout_mode } : {}),
    ...(record.mixed_items ? { mixed_items: toPayloadMixedItems(record.mixed_items) } : {}),
    ...(record.mixed_items_error ? { mixed_items_error: record.mixed_items_error } : {}),
    persisted: true,
    save_state: 'saved',
    save_error: null,
  };
}

export function rebindDraftRulesToWorkspaceDevice(rules: Rule[], deviceID: string): Rule[] {
  return rules.map((rule) => ({
    ...rule,
    device_id: rule.persisted ? rule.device_id : deviceID,
    save_state: rule.save_state ?? 'idle',
    save_error: rule.save_error ?? null,
  }));
}
