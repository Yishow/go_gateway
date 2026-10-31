import type { SourceRuleRecord } from '../../../../types/datalink';
import type { WorkbenchV2Action } from './useWorkbenchV2State';
import type { Rule } from './types';
import { hydrateStudioV2Rule } from './studioV2RuleAutosave';

const persistedFields = [
  'device_id', 'start_address', 'count', 'data_type', 'naming_prefix', 'enabled',
  'skipped_addresses', 'scale_multiplier', 'scale_offset', 'data_format',
  'share_enabled', 'share_start_register', 'share_stride',
] as const satisfies readonly (keyof Rule)[];
export type RuleFieldGenerations = Partial<Record<keyof Rule, number>>;
const same = (left: unknown, right: unknown) => JSON.stringify(left) === JSON.stringify(right);

/** Track user intent, including linked resets and A→B→A edits during a request. */
export function ruleSaveFields(action: WorkbenchV2Action): readonly (keyof Rule)[] {
  switch (action.type) {
    case 'addRule': return persistedFields;
    case 'updateRule': {
      const fields: (keyof Rule)[] = persistedFields.filter(field => action.patch[field] !== undefined);
      if (fields.some(field => ['start_address', 'count', 'data_type'].includes(field))) fields.push('skipped_addresses');
      return fields;
    }
    case 'toggleRuleEnabled': return ['enabled'];
    case 'updateRuleSkipped': case 'toggleRuleSkippedAddress': return ['skipped_addresses'];
    case 'toggleRuleShareEnabled': return ['share_enabled'];
    case 'updateRuleShareStart': return ['share_start_register'];
    case 'updateRuleShareStride': return ['share_stride'];
    default: return [];
  }
}

/** Only untouched submitted fields may take this request's canonical response. */
export function canonicalRulePatch(submitted: Rule, latest: Rule, saved: SourceRuleRecord, generation: number, edits: RuleFieldGenerations): Partial<Rule> {
  const canonical = hydrateStudioV2Rule(saved);
  return Object.fromEntries(persistedFields
    .filter(field => (edits[field] ?? 0) <= generation && same(latest[field], submitted[field]))
    .map(field => [field, canonical[field]]));
}
