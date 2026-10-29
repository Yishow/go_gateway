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

  it('fences optional workspace/revision scopes without changing the device-filtered API request', async () => {
    useStudioV2WorkspaceMeasurementsQuery(true, 'dev-1', { workspaceId: 'ws-1', workspaceRevision: 'revision-2' });
    const options = useQueryMock.mock.calls[0][0] as { queryKey: string[]; queryFn: () => Promise<MeasurementDefinition[]> };
    expect(options.queryKey).toEqual(['studio-v2-workspace', 'measurements', 'dev-1', 'ws-1', 'revision-2']);
    vi.mocked(studioV2WorkspaceMeasurementsAPI.list).mockResolvedValueOnce([]);
    await expect(options.queryFn()).resolves.toEqual([]);
    expect(studioV2WorkspaceMeasurementsAPI.list).toHaveBeenCalledWith('dev-1');
  });

  it.each([
    { workspace_id: 'foreign-workspace', device_id: 'dev-1' },
    { workspace_id: 'ws-1', device_id: 'foreign-device' },
  ])('rejects a mixed response containing a foreign scope %j', async (foreign) => {
    useStudioV2WorkspaceMeasurementsQuery(true, 'dev-1', { workspaceId: 'ws-1', workspaceRevision: 'revision-2' });
    const options = useQueryMock.mock.calls[0][0] as { queryFn: () => Promise<MeasurementDefinition[]> };
    const definition: MeasurementDefinition = {
      id: 'measurement-1', workspace_id: 'ws-1', device_id: 'dev-1', point_id: 'pt-0', tag_id: 'tag-0',
      name: 'Counter', equipment_id: 'cir-1', definition_revision: 'definition-1', source_binding_revision: 'binding-1',
      series_epoch: 'epoch-1', quantity: '', unit: 'count', semantic_kind: 'counter',
    };
    vi.mocked(studioV2WorkspaceMeasurementsAPI.list).mockResolvedValueOnce([definition, { ...definition, ...foreign }]);
    await expect(options.queryFn()).rejects.toThrow('Measurement definitions did not match the requested scope.');
  });

  it('keeps legacy unscoped queries compatible with workspace-wide results', async () => {
    useStudioV2WorkspaceMeasurementsQuery(true);
    const options = useQueryMock.mock.calls[0][0] as { queryKey: string[]; queryFn: () => Promise<MeasurementDefinition[]> };
    expect(options.queryKey).toEqual(['studio-v2-workspace', 'measurements', 'all']);
    vi.mocked(studioV2WorkspaceMeasurementsAPI.list).mockResolvedValueOnce([]);
    await expect(options.queryFn()).resolves.toEqual([]);
    expect(studioV2WorkspaceMeasurementsAPI.list).toHaveBeenCalledWith(undefined);
  });

  it('allows a workspace-scoped read without a device filter and preserves disabled queries', async () => {
    useStudioV2WorkspaceMeasurementsQuery(false, undefined, { workspaceId: 'ws-1', workspaceRevision: 'revision-2' });
    const options = useQueryMock.mock.calls[0][0] as { enabled: boolean; queryFn: () => Promise<MeasurementDefinition[]> };
    expect(options.enabled).toBe(false);
    const definition: MeasurementDefinition = {
      id: 'measurement-1', workspace_id: 'ws-1', device_id: 'dev-2', point_id: 'pt-0', tag_id: 'tag-0',
      name: 'Counter', equipment_id: 'cir-1', definition_revision: 'definition-1', source_binding_revision: 'binding-1',
      series_epoch: 'epoch-1', quantity: '', unit: 'count', semantic_kind: 'counter',
    };
    vi.mocked(studioV2WorkspaceMeasurementsAPI.list).mockResolvedValueOnce([definition]);
    const result = await options.queryFn();
    expect(result[0].workspace_id).toBe('ws-1');
    expect(result[0].device_id).toBe('dev-2');
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
      queryKey: ['studio-v2-workspace', 'measurements'],
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
      queryKey: ['studio-v2-workspace', 'measurements'],
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
      queryKey: ['studio-v2-workspace', 'measurements'],
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
      queryKey: ['studio-v2-workspace', 'measurements'],
    });
  });
});
