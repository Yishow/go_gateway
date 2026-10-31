import type { Mapping, MappingValue, Point, WorkbenchV2State } from '../../../features/datalink/workbench-v2/state/types';
type CanMutateRule = (ruleId: string) => boolean;

export function canAutosaveMappingRule(state: WorkbenchV2State, point: Point, canMutateRule: CanMutateRule): boolean {
  const ownerRule = state.rules.find((rule) => rule.id === point.rule_id);
  return Boolean(ownerRule?.persisted && ownerRule.save_state === 'saved' && canMutateRule(point.rule_id));
}

export function shouldAutosaveMappingPatch(patch: Partial<Mapping>): boolean {
  return (
    patch.tag_key !== undefined ||
    patch.display_name !== undefined ||
    patch.unit !== undefined ||
    patch.target_type !== undefined ||
    patch.scale !== undefined ||
    patch.offset !== undefined ||
    patch.enabled !== undefined
  );
}

export function sameMappingValue(left?: MappingValue, right?: MappingValue): boolean {
  return (
    left?.tag_key === right?.tag_key &&
    left?.display_name === right?.display_name &&
    left?.unit === right?.unit &&
    left?.target_type === right?.target_type &&
    left?.scale === right?.scale &&
    left?.offset === right?.offset &&
    left?.enabled === right?.enabled
  );
}

export function sameHydratedMapping(current: Mapping, next: Mapping): boolean {
  return (
    current.tag_key === next.tag_key &&
    current.display_name === next.display_name &&
    current.unit === next.unit &&
    current.target_type === next.target_type &&
    current.scale === next.scale &&
    current.offset === next.offset &&
    current.enabled === next.enabled &&
    current.mapping_id === next.mapping_id &&
    current.workspace_id === next.workspace_id &&
    current.rule_id === next.rule_id &&
    current.device_id === next.device_id &&
    current.address === next.address &&
    current.tag_id === next.tag_id &&
    current.persisted_point_id === next.persisted_point_id &&
    current.persisted === next.persisted &&
    current.save_state === next.save_state &&
    current.save_error === next.save_error &&
    sameMappingValue(current.local_value, next.local_value) &&
    sameMappingValue(current.persisted_value, next.persisted_value)
  );
}

