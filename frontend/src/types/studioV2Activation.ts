export type StudioV2ActivationResultStatus = 'pending' | 'success' | 'failed' | 'skipped';

export interface StudioV2ActivationResult {
  device_id: string;
  status: StudioV2ActivationResultStatus;
  message?: string;
  code?: string;
  action?: string;
  retryable?: boolean;
  request_id?: string;
  operation_id?: string;
}

export interface StudioV2ActivationResponse {
  workspace_id: string;
  results: StudioV2ActivationResult[];
  /** Client certainty for an error response; never inferred from presentation timing. */
  outcome?: 'failed' | 'unconfirmed';
  message?: string;
  code?: string;
  action?: string;
  retryable?: boolean;
  request_id?: string;
  operation_id?: string;
}

/** Existing server state used for navigation, without proving a new activation. */
export interface StudioV2ActivationRecovery {
  workspace_id: string;
  devices: Array<{ device_id: string; running: boolean }>;
  operation_id?: string;
}
