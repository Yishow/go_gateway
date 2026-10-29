import type { StudioV2LocatorReadScope } from '../../types/studioV2LocatorEntry';
import { studioLocatorReadKey } from '../../services/studioEntityLocator';
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { studioV2WorkspaceKeys } from './keys';
import {
  studioV2WorkspaceMeasurementsAPI,
  type PreviewTemplateRequest,
} from '../../services/studioV2WorkspaceMeasurements';
import type { MeasurementDefinition, TemplateApplyPreview } from '../../types/measurement';

const measurementFamilyKey = [...studioV2WorkspaceKeys.all, 'measurements'] as const;

interface MeasurementQueryScope {
  workspaceId: string;
  workspaceRevision: string;
}

export function useStudioV2WorkspaceMeasurementsQuery(
  enabled: boolean, deviceId?: string, scope?: MeasurementQueryScope, locatorScope?: StudioV2LocatorReadScope,
) {
  const selectedLocatorScope = locatorScope === undefined ? undefined : { workspaceId: locatorScope.workspaceId, workspaceRevision: locatorScope.workspaceRevision };
  return useQuery({
    // This fences late reads in the cache; only server mutation CAS can prove a current revision.
    queryKey: [...studioV2WorkspaceKeys.measurements(deviceId), ...(scope ? [scope.workspaceId, scope.workspaceRevision] : []), ...studioLocatorReadKey(selectedLocatorScope)],
    queryFn: async () => {
      const definitions = await (selectedLocatorScope ? studioV2WorkspaceMeasurementsAPI.list(deviceId, selectedLocatorScope) : studioV2WorkspaceMeasurementsAPI.list(deviceId));
      const hasForeignScope = scope && definitions.some((definition) => definition.workspace_id !== scope.workspaceId ||
        (deviceId !== undefined && definition.device_id !== deviceId));
      if (hasForeignScope) {
        throw new Error('Measurement definitions did not match the requested scope.');
      }
      return definitions;
    },
    enabled,
    retry: false,
    refetchOnWindowFocus: false,
  });
}

export function useStudioV2MeasurementTemplatesQuery(enabled: boolean) {
  return useQuery({
    queryKey: studioV2WorkspaceKeys.measurementTemplates(),
    networkMode: 'always',
    queryFn: () => studioV2WorkspaceMeasurementsAPI.listTemplates(),
    enabled,
    retry: false,
    refetchOnWindowFocus: false,
  });
}

export function usePreviewMeasurementTemplateMutation() {
  return useMutation({
    networkMode: 'always',
    retry: false,
    mutationFn: (request: PreviewTemplateRequest) =>
      studioV2WorkspaceMeasurementsAPI.previewTemplate(request),
  });
}

export function useApplyMeasurementTemplateMutation() {
  const queryClient = useQueryClient();

  return useMutation({
    networkMode: 'always',
    retry: false,
    mutationFn: (preview: TemplateApplyPreview) =>
      studioV2WorkspaceMeasurementsAPI.applyTemplate(preview),
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: measurementFamilyKey });
    },
  });
}

export function useCreateStudioV2MeasurementMutation() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (def: Partial<MeasurementDefinition>) =>
      studioV2WorkspaceMeasurementsAPI.create(def),
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: measurementFamilyKey });
    },
  });
}

export function useUpdateStudioV2MeasurementMutation() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ id, def }: { id: string; def: Partial<MeasurementDefinition> }) =>
      studioV2WorkspaceMeasurementsAPI.update(id, def),
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: measurementFamilyKey });
    },
  });
}

export function useDeleteStudioV2MeasurementMutation() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (id: string) => studioV2WorkspaceMeasurementsAPI.remove(id),
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: measurementFamilyKey });
    },
  });
}
