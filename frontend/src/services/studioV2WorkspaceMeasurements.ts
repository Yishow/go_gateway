import type { StudioV2LocatorReadScope } from '../types/studioV2LocatorEntry';
import { captureStudioLocatorReadScope, studioLocatorReadParams, assertStudioLocatorWorkspaceBinding, throwStudioLocatorReadError } from './studioEntityLocator';
import type { APIResponse } from '../types/datalink';
import type {
  MeasurementDefinition,
  MeasurementTemplate,
  TemplateApplyPreview,
} from '../types/measurement';
import { studioV2DatalinkApi } from './studioV2Workspace';

export interface PreviewTemplateRequest {
  template_id: string;
  device_id: string;
  device_name?: string;
  base_address?: string;
}

export interface UpdateMeasurementResponse {
  measurement: MeasurementDefinition;
  epoch_transitioned: boolean;
}

export const studioV2WorkspaceMeasurementsAPI = {
  async list(deviceId?: string, locatorScope?: StudioV2LocatorReadScope): Promise<MeasurementDefinition[]> {
    const selected = captureStudioLocatorReadScope(locatorScope);
    const params = selected ? { ...(deviceId ? { device_id: deviceId } : {}), ...studioLocatorReadParams(selected) } : deviceId ? { device_id: deviceId } : undefined;
    const res = await studioV2DatalinkApi.get<APIResponse<MeasurementDefinition[]>>(
      '/studio-v2/workspace/measurements',
      { params }
    ).catch(error => throwStudioLocatorReadError(error, selected));
    if (selected && (res.data.success !== true || !Array.isArray(res.data.data))) throw new Error('Scoped settings read unavailable.');
    const definitions = res.data.data ?? [];
    if (selected) {
      if (!Array.isArray(definitions)) throw new Error('Scoped settings read unavailable.');
      for (const definition of definitions) { assertStudioLocatorWorkspaceBinding(definition, selected);
        if (!definition || definition.workspace_id !== selected.workspaceId || (deviceId !== undefined && definition.device_id !== deviceId)) throw new Error('Scoped settings read unavailable.'); }
    }
    return definitions;
  },

  async listTemplates(): Promise<MeasurementTemplate[]> {
    const res = await studioV2DatalinkApi.get<APIResponse<MeasurementTemplate[]>>(
      '/studio-v2/workspace/measurements/templates'
    );
    return res.data.data ?? [];
  },

  async previewTemplate(request: PreviewTemplateRequest): Promise<TemplateApplyPreview> {
    const res = await studioV2DatalinkApi.post<APIResponse<TemplateApplyPreview>>(
      '/studio-v2/workspace/measurements/templates/preview',
      request
    );
    return res.data.data!;
  },

  async applyTemplate(preview: TemplateApplyPreview): Promise<{ applied: boolean; count: number }> {
    const res = await studioV2DatalinkApi.post<APIResponse<{ applied: boolean; count: number }>>(
      '/studio-v2/workspace/measurements/templates/apply',
      preview
    );
    return res.data.data!;
  },

  async create(def: Partial<MeasurementDefinition>): Promise<MeasurementDefinition> {
    const res = await studioV2DatalinkApi.post<APIResponse<MeasurementDefinition>>(
      '/studio-v2/workspace/measurements',
      def
    );
    return res.data.data!;
  },

  async update(id: string, def: Partial<MeasurementDefinition>): Promise<UpdateMeasurementResponse> {
    const res = await studioV2DatalinkApi.put<APIResponse<UpdateMeasurementResponse>>(
      `/studio-v2/workspace/measurements/${id}`,
      def
    );
    return res.data.data!;
  },

  async remove(id: string): Promise<void> {
    await studioV2DatalinkApi.delete(`/studio-v2/workspace/measurements/${id}`);
  },
};
