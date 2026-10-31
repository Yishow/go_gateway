/** Server-owned decoding evidence. Never resolve an effective format in the client. */
export type PointWordFormat = 'ABCD' | 'BADC' | 'CDAB' | 'DCBA';
export type PointFormatMode = 'legacy' | 'inherit' | 'override';
export interface PointFormatResolution {
  resolved: boolean;
  format_mode: PointFormatMode;
  configured_format: '' | PointWordFormat;
  effective_format: '' | PointWordFormat;
  effective_source: 'point' | 'connection' | 'native' | 'rule' | 'override';
  effective_policy: 'legacy_modbus' | 'legacy_mc3e' | 'legacy_fatek_components' | 'inherit_scalar' | 'override_scalar';
  native_compatibility: boolean;
  applicable: boolean;
  format_version: 0 | 1;
  point_revision: string;
  workspace_revision: string;
  rule_revision: string;
  connection_revision: string;
}
export interface PointFormatScope {
  workspace_id: string;
  device_id: string;
  point_id: string;
  rule_id?: string;
}
export interface PointFormatSnapshot extends PointFormatScope {
  rule_id: string;
  /** Source/Share authority, NOT DatabaseSetupRevision. Empty is never inferred. */
  workspace_revision: string;
  workspace_revision_origin: 'initial' | 'persisted';
  rule_revision: string;
  point_revision: string;
  connection_revision: string;
  editable: boolean;
  resolution: PointFormatResolution;
}
export type PointFormatIntent =
  | { format_mode: 'inherit'; configured_format: '' }
  | { format_mode: 'override'; configured_format: PointWordFormat };
export interface PointFormatOrdinaryPatch {
  name?: string;
  description?: string;
  address?: string;
  data_type?: 'bool' | 'int16' | 'int32' | 'int64' | 'uint16' | 'uint32' | 'uint64' | 'float32' | 'float64' | 'string';
  enabled?: boolean;
  polling_group_id?: string;
}
export type PointFormatOperation = Readonly<PointFormatScope & PointFormatIntent & {
  expected_workspace_revision: string;
  expected_rule_revision?: string;
  expected_point_revision: string;
  request_id: string;
  patch?: Readonly<PointFormatOrdinaryPatch>;
}>;
export interface PointFormatMutationResult {
  point: {
    id: string; device_id: string; revision_id: string;
    format_mode: 'inherit' | 'override'; format_override: '' | PointWordFormat;
    format_version: 1; data_format: string;
  };
  resolution: PointFormatResolution;
  workspace_revision: string;
  /** Configuration is saved; runtime is not claimed applied or healthy. */
  runtime_apply_status: 'pending';
  invalidation_status: 'pending' | 'invalidated';
}
