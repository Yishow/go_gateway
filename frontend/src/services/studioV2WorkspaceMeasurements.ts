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
  async list(deviceId?: string): Promise<MeasurementDefinition[]> {
    const params = deviceId ? { device_id: deviceId } : undefined;
    const res = await studioV2DatalinkApi.get<APIResponse<MeasurementDefinition[]>>(
      '/studio-v2/workspace/measurements',
      { params }
    );
    return res.data.data ?? [];
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
