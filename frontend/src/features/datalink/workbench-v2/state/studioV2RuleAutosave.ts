import type { CreateSourceRuleRequest, SourceRuleRecord } from '../../../../types/datalink';
import type { StudioV2UpdateSourceRuleRequest } from '../../../../services/studioV2Rules';
import type { Rule } from './types';

const VALID_DATA_FORMATS = new Set(['', 'ABCD', 'BADC', 'CDAB', 'DCBA']);

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
    share_enabled: false,
    share_start_register: null,
    share_stride: null,
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
