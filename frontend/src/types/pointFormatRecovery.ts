/** One explicit recovery of the original persisted format operation; no activation authority. */
export interface PointFormatRecoveryRequest {
  readonly workspace_id: string;
  readonly device_id: string;
  readonly point_id: string;
  readonly rule_id?: string;
  readonly request_id: string;
  readonly expected_workspace_revision: string;
  readonly expected_point_revision: string;
  readonly expected_rule_revision?: string;
  readonly expected_connection_revision: string;
  readonly expected_settings_revision: string;
}
export interface PointFormatRecoveryResult {
  outcome: 'recovered' | 'replayed';
  runtime_state: 'installed' | 'not_started';
  share_applied: boolean;
  hydration: { state: 'ready'; workspace_id: string; workspace_revision: string; settings_revision: string; readiness: true };
}
