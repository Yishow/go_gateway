import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { act, renderHook } from '@testing-library/react';
import type { ReactNode } from 'react';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import { useApplySchemaMutation } from '@/hooks/datalink/useStudioV2WorkspaceRecordingPlans';
import {
  useWriteGroupTestWriteMutation,
  useWriteGroupTestWritePreviewMutation,
} from '@/hooks/datalink/useStudioV2WriteGroupTestWrite';
import { studioV2WorkspaceRecordingPlansAPI } from '@/services/studioV2WorkspaceRecordingPlans';
import { studioV2WorkspaceWriteGroupsAPI } from '@/services/studioV2WorkspaceWriteGroups';
import { RecordingPlanResponseError } from '@/utils/recordingPlanJson';

vi.mock('@/services/studioV2WorkspaceRecordingPlans', () => ({
  studioV2WorkspaceRecordingPlansAPI: {
    schemaApplyConfirmed: vi.fn(),
  },
}));

vi.mock('@/services/studioV2WorkspaceWriteGroups', () => ({
  studioV2WorkspaceWriteGroupsAPI: {
    testWritePreview: vi.fn(),
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

const groupConfirmation = {
  token: 'tok-1',
  operation_id: 'op-1',
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
    vi.mocked(studioV2WorkspaceWriteGroupsAPI.testWritePreview).mockReset();
    vi.mocked(studioV2WorkspaceWriteGroupsAPI.testWrite).mockReset();
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

  it('does not retry group test write after a typed failure response', async () => {
    const failure = new RecordingPlanResponseError('write-group test write', {
      success: false,
      error: {
        code: 'RECORDING_TEST_WRITE_PLAN_UNRESOLVED',
        action: 'use the write group test write instead',
        retryable: false,
        request_id: 'req-501-write',
      },
    }, 'failed');
    vi.mocked(studioV2WorkspaceWriteGroupsAPI.testWrite).mockRejectedValueOnce(failure);
    const { result } = renderHook(() => useWriteGroupTestWriteMutation(), { wrapper: createWrapper() });

    await act(async () => {
      await expect(result.current.mutateAsync({ groupId: 'group-1', confirmation: groupConfirmation })).rejects.toBe(failure);
    });

    expect(studioV2WorkspaceWriteGroupsAPI.testWrite).toHaveBeenCalledTimes(1);
  });

  it('does not retry group test write after transport loss', async () => {
    const failure = new Error('transport lost');
    vi.mocked(studioV2WorkspaceWriteGroupsAPI.testWrite).mockRejectedValueOnce(failure);
    const { result } = renderHook(() => useWriteGroupTestWriteMutation(), { wrapper: createWrapper() });

    await act(async () => {
      await expect(result.current.mutateAsync({ groupId: 'group-1', confirmation: groupConfirmation })).rejects.toBe(failure);
    });

    expect(studioV2WorkspaceWriteGroupsAPI.testWrite).toHaveBeenCalledTimes(1);
  });

  it('does not retry group test write preview', async () => {
    const failure = new Error('transport lost');
    vi.mocked(studioV2WorkspaceWriteGroupsAPI.testWritePreview).mockRejectedValueOnce(failure);
    const { result } = renderHook(() => useWriteGroupTestWritePreviewMutation(), { wrapper: createWrapper() });

    await act(async () => {
      await expect(result.current.mutateAsync('group-1')).rejects.toBe(failure);
    });

    expect(studioV2WorkspaceWriteGroupsAPI.testWritePreview).toHaveBeenCalledTimes(1);
  });
});
