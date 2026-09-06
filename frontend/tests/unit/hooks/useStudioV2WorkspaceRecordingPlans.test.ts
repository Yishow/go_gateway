import { beforeEach, describe, expect, it, vi } from 'vitest';
import { studioV2WorkspaceKeys } from '@/hooks/datalink/keys';
import {
  useApplySchemaMutation,
  useCreateRecordingPlanMutation,
  useDeleteRecordingPlanMutation,
  usePreviewSchemaMutation,
  useStudioV2ConnectorCapabilitiesQuery,
  useStudioV2WorkspaceRecordingPlansQuery,
  useTestWritePlanMutation,
  useUpdateRecordingPlanMutation,
} from '@/hooks/datalink/useStudioV2WorkspaceRecordingPlans';
import { studioV2WorkspaceRecordingPlansAPI } from '@/services/studioV2WorkspaceRecordingPlans';
import type { RecordingPlan, SchemaPreviewToken, TestWriteResult } from '@/types/recordingPlan';

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

vi.mock('@/services/studioV2WorkspaceRecordingPlans', () => ({
  studioV2WorkspaceRecordingPlansAPI: {
    list: vi.fn(),
    get: vi.fn(),
    capabilities: vi.fn(),
    create: vi.fn(),
    update: vi.fn(),
    remove: vi.fn(),
    schemaPreview: vi.fn(),
    schemaApply: vi.fn(),
    testWrite: vi.fn(),
  },
}));

describe('useStudioV2WorkspaceRecordingPlans hooks', () => {
  beforeEach(() => {
    useMutationMock.mockReset();
    useQueryMock.mockReset();
    invalidateQueriesMock.mockReset();
  });

  it('queries recording plans with dedicated key', () => {
    useQueryMock.mockReturnValue({ data: [] });
    useStudioV2WorkspaceRecordingPlansQuery(true, 'dev-1');

    expect(useQueryMock).toHaveBeenCalledWith(
      expect.objectContaining({
        queryKey: studioV2WorkspaceKeys.recordingPlans('dev-1'),
        enabled: true,
        retry: false,
      })
    );
  });

  it('queries connector capabilities with kind', () => {
    useQueryMock.mockReturnValue({ data: [] });
    useStudioV2ConnectorCapabilitiesQuery(true, 'sqlite');

    expect(useQueryMock).toHaveBeenCalledWith(
      expect.objectContaining({
        queryKey: studioV2WorkspaceKeys.recordingPlanCapabilities('sqlite'),
        enabled: true,
      })
    );
  });

  it('creates and updates recording plan with invalidation', async () => {
    useMutationMock.mockReturnValue({ mutateAsync: vi.fn() });

    // Create
    useCreateRecordingPlanMutation();
    const createOptions = useMutationMock.mock.calls[0]?.[0] as {
      mutationFn: (plan: Partial<RecordingPlan>) => Promise<unknown>;
      onSuccess: () => Promise<void>;
    };
    const mockPlan = { id: 'plan-1', name: 'Plan 1' } as RecordingPlan;
    vi.mocked(studioV2WorkspaceRecordingPlansAPI.create).mockResolvedValueOnce(mockPlan);
    await expect(createOptions.mutationFn({ name: 'Plan 1' })).resolves.toEqual(mockPlan);
    await createOptions.onSuccess();
    expect(invalidateQueriesMock).toHaveBeenCalledWith({
      queryKey: studioV2WorkspaceKeys.recordingPlans(),
    });

    // Update
    useUpdateRecordingPlanMutation();
    const updateOptions = useMutationMock.mock.calls.at(-1)?.[0] as {
      mutationFn: (args: { id: string; plan: Partial<RecordingPlan> }) => Promise<unknown>;
      onSuccess: () => Promise<void>;
    };
    vi.mocked(studioV2WorkspaceRecordingPlansAPI.update).mockResolvedValueOnce(mockPlan);
    await expect(updateOptions.mutationFn({ id: 'plan-1', plan: { name: 'Plan 1 Updated' } })).resolves.toEqual(mockPlan);
    await updateOptions.onSuccess();
    expect(invalidateQueriesMock).toHaveBeenCalledWith({
      queryKey: studioV2WorkspaceKeys.recordingPlans(),
    });
  });

  it('deletes recording plan with invalidation', async () => {
    useMutationMock.mockReturnValue({ mutateAsync: vi.fn() });

    useDeleteRecordingPlanMutation();
    const deleteOptions = useMutationMock.mock.calls[0]?.[0] as {
      mutationFn: (id: string) => Promise<void>;
      onSuccess: () => Promise<void>;
    };

    vi.mocked(studioV2WorkspaceRecordingPlansAPI.remove).mockResolvedValueOnce();
    await expect(deleteOptions.mutationFn('plan-1')).resolves.toBeUndefined();
    await deleteOptions.onSuccess();
    expect(invalidateQueriesMock).toHaveBeenCalledWith({
      queryKey: studioV2WorkspaceKeys.recordingPlans(),
    });
  });

  it('previews, applies schema and performs test write', async () => {
    useMutationMock.mockReturnValue({ mutateAsync: vi.fn() });

    // Preview
    usePreviewSchemaMutation();
    const previewOptions = useMutationMock.mock.calls[0]?.[0] as {
      mutationFn: (req: unknown) => Promise<SchemaPreviewToken>;
    };
    const mockToken: SchemaPreviewToken = {
      token: 'tok-123',
      workspace_id: 'ws-1',
      plan_id: 'plan-1',
      plan_revision: '1',
      connector_id: 'conn-1',
      table_prefix: 'gw_record_',
      statements: ['CREATE TABLE ...'],
      expires_at: '2026-09-07T04:00:00Z',
      created_at: '2026-09-07T03:50:00Z',
    };
    vi.mocked(studioV2WorkspaceRecordingPlansAPI.schemaPreview).mockResolvedValueOnce(mockToken);
    await expect(previewOptions.mutationFn({ plan_id: 'plan-1', connector_id: 'conn-1', dialect: 'sqlite' })).resolves.toEqual(mockToken);

    // Apply
    useApplySchemaMutation();
    const applyOptions = useMutationMock.mock.calls.at(-1)?.[0] as {
      mutationFn: (token: string) => Promise<{ applied: boolean; message: string }>;
      onSuccess: () => Promise<void>;
    };
    vi.mocked(studioV2WorkspaceRecordingPlansAPI.schemaApply).mockResolvedValueOnce({ applied: true, message: 'Applied successfully' });
    await expect(applyOptions.mutationFn('tok-123')).resolves.toEqual({ applied: true, message: 'Applied successfully' });
    await applyOptions.onSuccess();
    expect(invalidateQueriesMock).toHaveBeenCalledWith({
      queryKey: studioV2WorkspaceKeys.recordingPlans(),
    });

    // Test Write
    useTestWritePlanMutation();
    const testWriteOptions = useMutationMock.mock.calls.at(-1)?.[0] as {
      mutationFn: (req: unknown) => Promise<TestWriteResult>;
    };
    const mockTestResult: TestWriteResult = {
      status: 'success',
      record_id: 'test-123',
      table: 'gw_record_samples',
      observed_at: '2026-09-07T03:50:00Z',
      delivered_at: '2026-09-07T03:50:00.012Z',
      message: 'Verified readback and cleaned up',
    };
    vi.mocked(studioV2WorkspaceRecordingPlansAPI.testWrite).mockResolvedValueOnce(mockTestResult);
    await expect(testWriteOptions.mutationFn({ plan_id: 'plan-1' })).resolves.toEqual(mockTestResult);
  });
});
