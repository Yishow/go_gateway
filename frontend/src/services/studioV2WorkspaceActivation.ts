import type { APIResponse } from '../types/datalink';
import type { StudioV2ActivationResponse } from '../types/studioV2Activation';
import { studioV2DatalinkApi } from './studioV2Workspace';
import { normalizeTypedEnvelope } from '../utils/safeJson';

export interface StudioV2ActivationRequest {
  workspace_revision: string;
  settings_revision: string;
  readiness_token: string;
  pending_saves?: number;
  save_error?: string;
  stale?: boolean;
}

export class StudioV2ActivationBarrierError extends Error {
  readonly code: string;
  readonly retryable: boolean;
  readonly action?: string;
  readonly request_id?: string;

  constructor(code: string, retryable = true, action?: string, requestId?: string) {
    super(code);
    this.name = 'StudioV2ActivationBarrierError';
    this.code = code;
    this.retryable = retryable;
    this.action = action;
    this.request_id = requestId;
  }
}

export const studioV2WorkspaceActivationAPI = {
  async activate(request?: StudioV2ActivationRequest): Promise<StudioV2ActivationResponse> {
    if (!request) {
      throw new StudioV2ActivationBarrierError('modbus_share_save_incomplete');
    }
    if (request.stale) {
      throw new StudioV2ActivationBarrierError('modbus_share_revision_conflict');
    }
    if (
      !request.workspace_revision ||
      !request.settings_revision ||
      !request.readiness_token ||
      (request.pending_saves ?? 0) > 0 ||
      request.save_error
    ) {
      throw new StudioV2ActivationBarrierError('modbus_share_save_incomplete');
    }

    const { stale: _stale, ...payload } = request;
    try {
      const res = await studioV2DatalinkApi.post<APIResponse<StudioV2ActivationResponse>>(
        '/studio-v2/workspace/activate',
        payload,
      );
      return res.data.data!;
    } catch (error) {
      const response = typeof error === 'object' && error !== null && 'response' in error
        ? error.response : undefined;
      const data = typeof response === 'object' && response !== null && 'data' in response
        ? response.data : undefined;
      const status = typeof response === 'object' && response !== null && 'status' in response
        ? response.status : undefined;
      const envelope = normalizeTypedEnvelope(data);
      const code = envelope.code
        ?? (status === 409 ? 'modbus_share_revision_conflict'
          : status === 422 ? 'modbus_share_save_incomplete'
            : 'activation_failed');
      throw new StudioV2ActivationBarrierError(
        code,
        envelope.retryable ?? status !== 422,
        envelope.action,
        envelope.requestId,
      );
    }
  },
};
