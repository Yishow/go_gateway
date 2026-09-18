import { useMutation } from '@tanstack/react-query';
import {
  studioV2WorkspaceActivationAPI,
  type StudioV2ActivationRequest,
} from '../../services/studioV2WorkspaceActivation';
import type { StudioV2ActivationResponse } from '../../types/studioV2Activation';

export function useActivateStudioV2WorkspaceMutation() {
  return useMutation<StudioV2ActivationResponse, Error, StudioV2ActivationRequest | undefined>({
    mutationFn: (request?: StudioV2ActivationRequest) => studioV2WorkspaceActivationAPI.activate(request),
    retry: false,
  });
}
