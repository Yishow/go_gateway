import { useMutation } from '@tanstack/react-query';
import { studioV2WorkspaceActivationAPI } from '../../services/studioV2WorkspaceActivation';

export function useActivateStudioV2WorkspaceMutation() {
  return useMutation({
    mutationFn: () => studioV2WorkspaceActivationAPI.activate(),
  });
}
