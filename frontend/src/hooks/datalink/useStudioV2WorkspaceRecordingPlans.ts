import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { studioV2WorkspaceKeys } from './keys';
import {
  studioV2WorkspaceRecordingPlansAPI,
  type SchemaPreviewRequest,
  type TestWriteRequest,
} from '../../services/studioV2WorkspaceRecordingPlans';
import type { RecordingPlan } from '../../types/recordingPlan';

export function useStudioV2WorkspaceRecordingPlansQuery(enabled: boolean, deviceId?: string) {
  return useQuery({
    queryKey: studioV2WorkspaceKeys.recordingPlans(deviceId),
    queryFn: () => studioV2WorkspaceRecordingPlansAPI.list(deviceId),
    enabled,
    retry: false,
    refetchOnWindowFocus: false,
  });
}

export function useStudioV2ConnectorCapabilitiesQuery(enabled: boolean, kind?: string) {
  return useQuery({
    queryKey: studioV2WorkspaceKeys.recordingPlanCapabilities(kind ?? 'all'),
    queryFn: () => studioV2WorkspaceRecordingPlansAPI.capabilities(kind),
    enabled,
    retry: false,
    refetchOnWindowFocus: false,
  });
}

export function useCreateRecordingPlanMutation() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (plan: Partial<RecordingPlan>) => studioV2WorkspaceRecordingPlansAPI.create(plan),
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: studioV2WorkspaceKeys.recordingPlans() });
    },
  });
}

export function useUpdateRecordingPlanMutation() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ id, plan }: { id: string; plan: Partial<RecordingPlan> }) =>
      studioV2WorkspaceRecordingPlansAPI.update(id, plan),
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: studioV2WorkspaceKeys.recordingPlans() });
    },
  });
}

export function useDeleteRecordingPlanMutation() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (id: string) => studioV2WorkspaceRecordingPlansAPI.remove(id),
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: studioV2WorkspaceKeys.recordingPlans() });
    },
  });
}

export function usePreviewSchemaMutation() {
  return useMutation({
    mutationFn: (req: SchemaPreviewRequest) => studioV2WorkspaceRecordingPlansAPI.schemaPreview(req),
  });
}

export function useApplySchemaMutation() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (token: string) => studioV2WorkspaceRecordingPlansAPI.schemaApply(token),
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: studioV2WorkspaceKeys.recordingPlans() });
    },
  });
}

export function useTestWritePlanMutation() {
  return useMutation({
    mutationFn: (req: TestWriteRequest) => studioV2WorkspaceRecordingPlansAPI.testWrite(req),
  });
}
