import { beforeEach, describe, expect, it, vi } from 'vitest';
import { studioV2WorkspaceKeys } from '@/hooks/datalink/keys';
import {
  useApplyMeasurementTemplateMutation,
  useCreateStudioV2MeasurementMutation,
  useDeleteStudioV2MeasurementMutation,
  usePreviewMeasurementTemplateMutation,
  useStudioV2MeasurementTemplatesQuery,
  useStudioV2WorkspaceMeasurementsQuery,
  useUpdateStudioV2MeasurementMutation,
} from '@/hooks/datalink/useStudioV2WorkspaceMeasurements';
import { studioV2WorkspaceMeasurementsAPI } from '@/services/studioV2WorkspaceMeasurements';
import type { MeasurementDefinition, TemplateApplyPreview } from '@/types/measurement';

const useMutationMock = vi.fn();
const useQueryMock = vi.fn();
const invalidateQueriesMock = vi.fn();

vi.mock('@tanstack/react-query', async () => {
  const actual = await vi.importActual<typeof import('@tanstack/react-query')>('@tanstack/react-query');
  return {
    ...actual,
    useQuery: (options: unknown) => useQueryMock(options),
    useMutation: (options: unknown) => useMutationMock(options),
    useQueryClient: () => ({
      invalidateQueries: invalidateQueriesMock,
    }),
  };
});

vi.mock('@/services/studioV2WorkspaceMeasurements', () => ({
  studioV2WorkspaceMeasurementsAPI: {
    list: vi.fn(),
    listTemplates: vi.fn(),
    previewTemplate: vi.fn(),
    applyTemplate: vi.fn(),
    create: vi.fn(),
    update: vi.fn(),
    remove: vi.fn(),
  },
}));

describe('useStudioV2WorkspaceMeasurements hooks', () => {
  beforeEach(() => {
    useMutationMock.mockReset();
    useQueryMock.mockReset();
    invalidateQueriesMock.mockReset();
  });

  it('queries measurements with dedicated key', async () => {
    useQueryMock.mockReturnValue({ data: [] });
    useStudioV2WorkspaceMeasurementsQuery(true, 'dev-1');

    expect(useQueryMock).toHaveBeenCalledWith(
      expect.objectContaining({
        queryKey: studioV2WorkspaceKeys.measurements('dev-1'),
        enabled: true,
        retry: false,
      })
    );
  });

  it('queries templates with dedicated key', async () => {
    useQueryMock.mockReturnValue({ data: [] });
    useStudioV2MeasurementTemplatesQuery(true);

    expect(useQueryMock).toHaveBeenCalledWith(
      expect.objectContaining({
        queryKey: studioV2WorkspaceKeys.measurementTemplates(),
        enabled: true,
      })
    );
  });

  it('previews and applies template with invalidation', async () => {
    useMutationMock.mockReturnValue({ mutateAsync: vi.fn() });

    usePreviewMeasurementTemplateMutation();
    const previewOptions = useMutationMock.mock.calls[0]?.[0] as {
      mutationFn: (req: unknown) => Promise<unknown>;
    };

    const mockPreview: TemplateApplyPreview = {
      device_id: 'dev-1',
      device_name: 'Power Meter',
      template_id: 'three_phase_power_meter_v1',
      proposed_items: [],
      definitions: [],
      confirmed: true,
      needs_review: false,
    };
    vi.mocked(studioV2WorkspaceMeasurementsAPI.previewTemplate).mockResolvedValueOnce(mockPreview);
    await expect(
      previewOptions.mutationFn({ template_id: 'three_phase_power_meter_v1', device_id: 'dev-1' })
    ).resolves.toEqual(mockPreview);

    useApplyMeasurementTemplateMutation();
    const applyOptions = useMutationMock.mock.calls.at(-1)?.[0] as {
      mutationFn: (preview: TemplateApplyPreview) => Promise<unknown>;
      onSuccess: () => Promise<void>;
    };

    vi.mocked(studioV2WorkspaceMeasurementsAPI.applyTemplate).mockResolvedValueOnce({
      applied: true,
      count: 8,
    });
    await expect(applyOptions.mutationFn(mockPreview)).resolves.toEqual({ applied: true, count: 8 });

    await applyOptions.onSuccess();
    expect(invalidateQueriesMock).toHaveBeenCalledWith({
      queryKey: studioV2WorkspaceKeys.measurements(),
    });
  });

  it('creates and updates measurement', async () => {
    useMutationMock.mockReturnValue({ mutateAsync: vi.fn() });

    useCreateStudioV2MeasurementMutation();
    const createOptions = useMutationMock.mock.calls[0]?.[0] as {
      mutationFn: (def: Partial<MeasurementDefinition>) => Promise<unknown>;
      onSuccess: () => Promise<void>;
    };

    const createdDef = { id: 'meas-1', name: 'Voltage' } as MeasurementDefinition;
    vi.mocked(studioV2WorkspaceMeasurementsAPI.create).mockResolvedValueOnce(createdDef);
    await expect(createOptions.mutationFn({ name: 'Voltage' })).resolves.toEqual(createdDef);

    await createOptions.onSuccess();
    expect(invalidateQueriesMock).toHaveBeenCalledWith({
      queryKey: studioV2WorkspaceKeys.measurements(),
    });

    useUpdateStudioV2MeasurementMutation();
    const updateOptions = useMutationMock.mock.calls.at(-1)?.[0] as {
      mutationFn: (payload: { id: string; def: Partial<MeasurementDefinition> }) => Promise<unknown>;
      onSuccess: () => Promise<void>;
    };

    vi.mocked(studioV2WorkspaceMeasurementsAPI.update).mockResolvedValueOnce({
      measurement: createdDef,
      epoch_transitioned: false,
    });
    await expect(updateOptions.mutationFn({ id: 'meas-1', def: { name: 'New Voltage' } })).resolves.toEqual({
      measurement: createdDef,
      epoch_transitioned: false,
    });

    await updateOptions.onSuccess();
    expect(invalidateQueriesMock).toHaveBeenCalledWith({
      queryKey: studioV2WorkspaceKeys.measurements(),
    });
  });

  it('deletes measurement with invalidation', async () => {
    useMutationMock.mockReturnValue({ mutateAsync: vi.fn() });

    useDeleteStudioV2MeasurementMutation();
    const deleteOptions = useMutationMock.mock.calls[0]?.[0] as {
      mutationFn: (id: string) => Promise<void>;
      onSuccess: () => Promise<void>;
    };

    vi.mocked(studioV2WorkspaceMeasurementsAPI.remove).mockResolvedValueOnce();
    await expect(deleteOptions.mutationFn('meas-1')).resolves.toBeUndefined();

    await deleteOptions.onSuccess();
    expect(invalidateQueriesMock).toHaveBeenCalledWith({
      queryKey: studioV2WorkspaceKeys.measurements(),
    });
  });
});
