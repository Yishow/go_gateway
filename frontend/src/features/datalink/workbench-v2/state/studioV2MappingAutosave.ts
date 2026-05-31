import type { StudioV2WorkspaceMappingRecord } from '../../../../types/datalink';
import type { Mapping, MappingValue, Point, TargetType } from './types';
import type { StudioV2WorkspaceMappingRequest } from '../../../../services/studioV2Mappings';

const VALID_TARGET_TYPES = new Set<TargetType>([
  'bool',
  'int16',
  'int32',
  'int64',
  'uint16',
  'uint32',
  'uint64',
  'float32',
  'float64',
  'string',
]);

function sameMappingValue(left?: MappingValue, right?: MappingValue): boolean {
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

export function mappingValueFromRecord(record: StudioV2WorkspaceMappingRecord): MappingValue {
  return {
    tag_key: record.tag_key,
    display_name: record.display_name,
    unit: record.unit,
    target_type: record.target_type,
    scale: record.scale,
    offset: record.offset,
    enabled: record.enabled,
  };
}

export function mappingValueFromState(mapping: Mapping): MappingValue {
  return {
    tag_key: mapping.tag_key,
    display_name: mapping.display_name,
    unit: mapping.unit,
    target_type: mapping.target_type,
    scale: mapping.scale,
    offset: mapping.offset,
    enabled: mapping.enabled,
  };
}

export function isStudioV2MappingValid(mapping: Mapping): boolean {
  return Boolean(
    mapping.tag_key.trim() !== '' &&
    VALID_TARGET_TYPES.has(mapping.target_type),
  );
}

export function hydrateStudioV2Mapping(point: Point, record: StudioV2WorkspaceMappingRecord, current?: Mapping): Mapping {
  const persistedValue = mappingValueFromRecord(record);
  const keepLocal = Boolean(
    current && (
      current.save_state === 'draft-invalid' ||
      current.save_state === 'save-error' ||
      current.save_state === 'saving' ||
      (
        current.persisted_value &&
        !sameMappingValue(current.local_value, current.persisted_value)
      )
    ),
  );
  const localValue = keepLocal && current?.local_value
    ? current.local_value
    : persistedValue;

  return {
    point_id: point.id,
    tag_key: localValue.tag_key,
    display_name: localValue.display_name,
    unit: localValue.unit,
    target_type: localValue.target_type,
    scale: localValue.scale,
    offset: localValue.offset,
    enabled: localValue.enabled,
    mapping_id: record.id,
    workspace_id: record.workspace_id,
    rule_id: record.rule_id,
    device_id: record.device_id,
    address: record.address,
    tag_id: record.tag_id,
    persisted_point_id: record.point_id,
    persisted: true,
    local_value: localValue,
    persisted_value: persistedValue,
    save_state: keepLocal ? current?.save_state ?? 'idle' : 'saved',
    save_error: keepLocal ? current?.save_error ?? null : null,
  };
}

export function toStudioV2MappingCreateRequest(point: Point, mapping: Mapping): StudioV2WorkspaceMappingRequest {
  return {
    rule_id: point.rule_id,
    address: point.address,
    tag_key: mapping.tag_key,
    display_name: mapping.display_name,
    unit: mapping.unit,
    target_type: mapping.target_type,
    scale: mapping.scale,
    offset: mapping.offset,
    enabled: mapping.enabled,
  };
}

export function toStudioV2MappingUpdateRequest(point: Point, mapping: Mapping): StudioV2WorkspaceMappingRequest {
  return toStudioV2MappingCreateRequest(point, mapping);
}
