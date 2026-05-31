export type StudioV2ActivationResultStatus = 'success' | 'failed';

export interface StudioV2ActivationResult {
  device_id: string;
  status: StudioV2ActivationResultStatus;
  message: string;
}

export interface StudioV2ActivationResponse {
  workspace_id: string;
  results: StudioV2ActivationResult[];
  message?: string;
}
