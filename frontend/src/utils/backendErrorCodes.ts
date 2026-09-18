/**
 * Canonical list of backend error codes the UI is allowed to retain.
 *
 * Single source shared by `typedErrors.ts` (`getSafeErrorMessage`) and
 * `safeJson.ts` (`normalizeTypedEnvelope`); both must derive their allowlist
 * from this list so they cannot drift apart. Codes come from the backend
 * `internal/api/handlers/typed_errors.go` `ErrCode*` constants.
 */
export const BACKEND_ERROR_CODES = [
  'preview_invalid_request',
  'preview_unavailable',
  'preview_stream_closed',
  'runtime_device_not_found',
  'runtime_snapshot_unavailable',
  'runtime_stream_unavailable',
  'workspace_not_ready',
  'activation_failed',
  'settings_unavailable',
  'settings_update_failed',
  'settings_invalid',
  'activation_request_invalid',
  'readiness_blocked',
  'RECORDING_SCHEMA_NOT_IMPLEMENTED',
  'RECORDING_TEST_WRITE_NOT_IMPLEMENTED',
  'RECORDING_PLAN_MEMBERS_INVALID',
  'SCHEMA_CONFIRMATION_REQUIRED',
  'WORKSPACE_SCHEMA_PREPARATION_REQUIRED',
  'RECORDING_SCHEMA_PREVIEW_STALE',
  'RECORDING_SCHEMA_PREVIEW_EXPIRED',
  'RECORDING_SCHEMA_PREVIEW_NOT_FOUND',
  'RECORDING_SCHEMA_OPERATION_MISMATCH',
  'RECORDING_SCHEMA_OPERATION_BUSY',
  'RECORDING_SCHEMA_OPERATION_UNKNOWN',
  'RECORDING_SCHEMA_PREFIX_INVALID',
  'RECORDING_SCHEMA_INCOMPATIBLE',
  'RECORDING_SCHEMA_PERMISSION_DENIED',
  'RECORDING_SCHEMA_TARGET_UNAVAILABLE',
  'RECORDING_OPERATION_NOT_FOUND',
  'modbus_share_save_incomplete',
  'modbus_share_revision_conflict',
  'modbus_share_hydration_required',
  'modbus_share_disabled',
  'modbus_share_range_collision',
  'modbus_share_capacity_exceeded',
  'modbus_share_listener_bind_failed',
  'modbus_share_workspace_scope',
  'modbus_share_reconcile_failed',
  'modbus_share_dirty_unknown',
  'modbus_share_projection_required',
  'modbus_share_invalid_geometry',
  'validation',
  'revision_mismatch',
  'not_found',
  'internal',
] as const;

export type BackendErrorCode = (typeof BACKEND_ERROR_CODES)[number];
