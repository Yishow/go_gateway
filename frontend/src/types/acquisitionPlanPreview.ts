/** Read-only installed acquisition estimates; never recording membership or SQL writes. */
export interface AcquisitionPreviewScope {
  workspace_id: string;
  device_id: string;
  expected_workspace_revision: string;
  expected_preview_revision?: string;
}
export interface AcquisitionWireUnits {
  area: string;
  function: string;
  address_unit_bits: 1 | 16;
  unique_address_units: number;
}
export type AcquisitionActual = { status: 'unavailable'; reason: string } | {
  status: 'available';
  acquisition_id: string;
  observed_at: string;
  declared_plan_revision: string;
  initial_plan_revision: string;
  phases: { phase: 'initial' | 'safe' | 'fallback'; plan_revision: string }[];
  wire_attempts: number;
  retry_wire_attempts: number;
  fallback_wire_attempts: number;
  logical_requests: { status: 'not_started' } | { status: 'observed'; succeeded: number; failed: number; canceled: number };
  cycle_latency_ms: number;
  overrun: boolean;
  last_good_age_ms: number | null;
};
export interface AcquisitionPreviewGroup {
  polling_group_id: string;
  poll_interval_ms: number;
  plan_revision: string;
  declared_plan_revision: string;
  capability_revision: string;
  basis: 'declared_no_connection' | 'declared_current_connection' | 'safe_current_connection';
  acquisition_point_count: number;
  unique_wire_units: AcquisitionWireUnits[];
  planned_requests: number;
  request_limits?: { kind: 'contiguous' | 'random'; area: string; function: string; native_unit_limit: number; request_frame_bytes: number; response_frame_bytes: number; frame_cost_known: boolean }[];
  point_errors: { point_id: string; code: string }[];
  fallback_reasons: string[];
  actual: AcquisitionActual;
}
export type AcquisitionPlanPreview = {
  workspace_id: string;
  device_id: string;
  workspace_revision: string;
} & ({ status: 'unavailable'; reason: string } | {
  status: 'available';
  preview_revision: string;
  estimate: true;
  groups: AcquisitionPreviewGroup[];
  output_loss: { status: 'unavailable'; reason: string } | {
    status: 'partial'; history_complete: false;
    retained_detail: { point_id: string; losses: number }[];
  };
});
export type AcquisitionPlanPreviewResult = { ok: true; value: AcquisitionPlanPreview } |
  { ok: false; reason: 'foreign-scope' | 'stale-reference' | 'unavailable' };
