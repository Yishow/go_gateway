export type StudioV2ActivationResultStatus = 'pending' | 'success' | 'failed' | 'skipped';

export interface StudioV2ActivationResult {
  device_id: string;
  status: StudioV2ActivationResultStatus;
  message?: string;
  code?: string;
  action?: string;
  retryable?: boolean;
  request_id?: string;
}

export interface StudioV2ActivationResponse {
  workspace_id: string;
  results: StudioV2ActivationResult[];
  message?: string;
  code?: string;
  action?: string;
  retryable?: boolean;
  request_id?: string;
}
