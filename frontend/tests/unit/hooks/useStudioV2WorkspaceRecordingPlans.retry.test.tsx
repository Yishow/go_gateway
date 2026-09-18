import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { act, renderHook } from '@testing-library/react';
import type { ReactNode } from 'react';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import {
  useApplySchemaMutation,
  useTestWritePlanMutation,
} from '@/hooks/datalink/useStudioV2WorkspaceRecordingPlans';
import { studioV2WorkspaceRecordingPlansAPI } from '@/services/studioV2WorkspaceRecordingPlans';
import { RecordingPlanResponseError } from '@/utils/recordingPlanJson';

vi.mock('@/services/studioV2WorkspaceRecordingPlans', () => ({
  studioV2WorkspaceRecordingPlansAPI: {
    schemaApplyConfirmed: vi.fn(),
    testWrite: vi.fn(),
  },
}));

const confirmation = {
  token: 'tok-1',
  operation_id: 'op-1',
  expected_workspace_revision: 'setup-1',
  expected_plan_revision: 'rev-1',
  expected_connector_revision: 'identity-1',
};

function createWrapper() {
  const queryClient = new QueryClient({
    defaultOptions: {
      mutations: { retry: 1, retryDelay: 0 },
      queries: { retry: false },
    },
  });
  return function QueryWrapper({ children }: { children: ReactNode }) {
    return <QueryClientProvider client={queryClient}>{children}</QueryClientProvider>;
  };
}

describe('recording mutation retry policy', () => {
  beforeEach(() => {
    vi.mocked(studioV2WorkspaceRecordingPlansAPI.schemaApplyConfirmed).mockReset();
    vi.mocked(studioV2WorkspaceRecordingPlansAPI.testWrite).mockReset();
  });

  it('does not retry schema apply after a 501 response', async () => {
    const failure = new RecordingPlanResponseError('recording schema apply', {
      success: false,
      error: {
        code: 'RECORDING_SCHEMA_NOT_IMPLEMENTED',
        action: 'wait_for_supported_operation',
        retryable: false,
        request_id: 'req-501',
      },
    }, 'failed');
    vi.mocked(studioV2WorkspaceRecordingPlansAPI.schemaApplyConfirmed).mockRejectedValueOnce(failure);
    const { result } = renderHook(() => useApplySchemaMutation(), { wrapper: createWrapper() });

    await act(async () => {
      await expect(result.current.mutateAsync(confirmation)).rejects.toBe(failure);
    });

    expect(studioV2WorkspaceRecordingPlansAPI.schemaApplyConfirmed).toHaveBeenCalledTimes(1);
  });

  it('does not retry schema apply after transport loss', async () => {
    const failure = new Error('transport lost');
    vi.mocked(studioV2WorkspaceRecordingPlansAPI.schemaApplyConfirmed).mockRejectedValueOnce(failure);
    const { result } = renderHook(() => useApplySchemaMutation(), { wrapper: createWrapper() });

    await act(async () => {
      await expect(result.current.mutateAsync({ ...confirmation, token: 'tok-transport' })).rejects.toBe(failure);
    });

    expect(studioV2WorkspaceRecordingPlansAPI.schemaApplyConfirmed).toHaveBeenCalledTimes(1);
  });

  it('does not retry test write after a 501 response', async () => {
    const failure = new RecordingPlanResponseError('recording test write', {
      success: false,
      error: {
        code: 'RECORDING_TEST_WRITE_NOT_IMPLEMENTED',
        action: 'wait_for_supported_operation',
        retryable: false,
        request_id: 'req-501-write',
      },
    }, 'failed');
    vi.mocked(studioV2WorkspaceRecordingPlansAPI.testWrite).mockRejectedValueOnce(failure);
    const { result } = renderHook(() => useTestWritePlanMutation(), { wrapper: createWrapper() });

    await act(async () => {
      await expect(result.current.mutateAsync({ plan_id: 'plan-1' })).rejects.toBe(failure);
    });

    expect(studioV2WorkspaceRecordingPlansAPI.testWrite).toHaveBeenCalledTimes(1);
  });

  it('does not retry test write after transport loss', async () => {
    const failure = new Error('transport lost');
    vi.mocked(studioV2WorkspaceRecordingPlansAPI.testWrite).mockRejectedValueOnce(failure);
    const { result } = renderHook(() => useTestWritePlanMutation(), { wrapper: createWrapper() });

    await act(async () => {
      await expect(result.current.mutateAsync({ plan_id: 'plan-transport' })).rejects.toBe(failure);
    });

    expect(studioV2WorkspaceRecordingPlansAPI.testWrite).toHaveBeenCalledTimes(1);
  });
});
