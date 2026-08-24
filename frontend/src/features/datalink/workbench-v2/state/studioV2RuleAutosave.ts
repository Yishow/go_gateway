import type { CreateSourceRuleRequest, SourceRuleRecord } from '../../../../types/datalink';
import type { StudioV2UpdateSourceRuleRequest } from '../../../../services/studioV2Rules';
import type { Rule } from './types';

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
  return Boolean(
    rule.device_id.trim() !== '' &&
    rule.start_address.trim() !== '' &&
    rule.count > 0 &&
    VALID_DATA_FORMATS.has(rule.data_format),
  );
}

export function toStudioV2RuleCreateRequest(rule: Rule): CreateSourceRuleRequest {
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
  return {
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
