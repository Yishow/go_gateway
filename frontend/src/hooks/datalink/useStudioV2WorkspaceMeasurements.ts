import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { studioV2WorkspaceKeys } from './keys';
import {
  studioV2WorkspaceMeasurementsAPI,
  type PreviewTemplateRequest,
} from '../../services/studioV2WorkspaceMeasurements';
import type { MeasurementDefinition, TemplateApplyPreview } from '../../types/measurement';

export function useStudioV2WorkspaceMeasurementsQuery(enabled: boolean, deviceId?: string) {
  return useQuery({
    queryKey: studioV2WorkspaceKeys.measurements(deviceId),
    queryFn: () => studioV2WorkspaceMeasurementsAPI.list(deviceId),
    enabled,
    retry: false,
    refetchOnWindowFocus: false,
  });
}

export function useStudioV2MeasurementTemplatesQuery(enabled: boolean) {
  return useQuery({
    queryKey: studioV2WorkspaceKeys.measurementTemplates(),
    queryFn: () => studioV2WorkspaceMeasurementsAPI.listTemplates(),
    enabled,
    retry: false,
    refetchOnWindowFocus: false,
  });
}

export function usePreviewMeasurementTemplateMutation() {
  return useMutation({
    mutationFn: (request: PreviewTemplateRequest) =>
      studioV2WorkspaceMeasurementsAPI.previewTemplate(request),
  });
}

export function useApplyMeasurementTemplateMutation() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (preview: TemplateApplyPreview) =>
      studioV2WorkspaceMeasurementsAPI.applyTemplate(preview),
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: studioV2WorkspaceKeys.measurements() });
    },
  });
}

export function useCreateStudioV2MeasurementMutation() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (def: Partial<MeasurementDefinition>) =>
      studioV2WorkspaceMeasurementsAPI.create(def),
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: studioV2WorkspaceKeys.measurements() });
    },
  });
}

export function useUpdateStudioV2MeasurementMutation() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ id, def }: { id: string; def: Partial<MeasurementDefinition> }) =>
      studioV2WorkspaceMeasurementsAPI.update(id, def),
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: studioV2WorkspaceKeys.measurements() });
    },
  });
}

export function useDeleteStudioV2MeasurementMutation() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (id: string) => studioV2WorkspaceMeasurementsAPI.remove(id),
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: studioV2WorkspaceKeys.measurements() });
    },
  });
}
