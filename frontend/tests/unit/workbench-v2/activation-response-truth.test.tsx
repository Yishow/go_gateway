import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { act, renderHook } from '@testing-library/react';
import type { ReactNode } from 'react';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import { useActivateStudioV2WorkspaceMutation } from '@/hooks/datalink/useStudioV2WorkspaceActivation';
import { studioV2WorkspaceActivationAPI } from '@/services/studioV2WorkspaceActivation';
import { studioV2DatalinkApi } from '@/services/studioV2Workspace';

vi.mock('@/services/studioV2Workspace', () => ({ studioV2DatalinkApi: { post: vi.fn() } }));

const request = { workspace_revision: 'workspace-rev-3', settings_revision: 'settings-rev-7', readiness_token: 'ready-token-1' };

describe('BackendBackedActivationResults response boundary', () => {
  beforeEach(() => vi.resetAllMocks());

  it.each([
    null,
    { data: { workspace_id: 'workspace-1', results: [] } },
    { success: false, data: { workspace_id: 'workspace-1', results: [{ device_id: 'dev-1', status: 'success' }] } },
    { success: true, data: { workspace_id: 'workspace-1', results: [{ device_id: 'dev-1', status: 'invented' }] } },
    { success: true, data: { results: [] } },
    { success: true, data: { workspace_id: 'workspace-1' } },
  ])('classifies malformed nominal 200 as unconfirmed: %j', async (body) => {
    vi.mocked(studioV2DatalinkApi.post).mockResolvedValueOnce({ data: body } as never);
    await expect(studioV2WorkspaceActivationAPI.activate(request)).rejects.toMatchObject({
      outcome: 'unconfirmed', retryable: false,
    });
  });

  it('retains an existing operation id after a malformed response without raw diagnostics', async () => {
    vi.mocked(studioV2DatalinkApi.post).mockResolvedValueOnce({ data: {
      success: true, data: { operation_id: 'op-activation-1', results: 'invalid' },
      message: 'raw internal diagnostic',
    } } as never);
    const error = await studioV2WorkspaceActivationAPI.activate(request).catch((failure: unknown) => failure);
    expect(error).toMatchObject({ outcome: 'unconfirmed', retryable: false, operation_id: 'op-activation-1' });
    expect(String(error)).not.toContain('raw internal diagnostic');
  });

  it.each([false, true])('preserves safe typed metadata from a malformed 200 with success=%s', async (success) => {
    vi.mocked(studioV2DatalinkApi.post).mockResolvedValueOnce({ data: {
      success, data: { results: 'invalid' }, error: {
        code: 'modbus_share_revision_conflict', action: 'retry', retryable: true,
        request_id: 'req-malformed-full-0123456789', operation_id: 'op-malformed-typed',
        message: 'raw database diagnostic',
      },
    } } as never);
    const error = await studioV2WorkspaceActivationAPI.activate(request).catch((failure: unknown) => failure);
    expect(error).toMatchObject({
      code: 'modbus_share_revision_conflict', action: 'retry', outcome: 'unconfirmed', retryable: false,
      request_id: 'req-malformed-full-0123456789', operation_id: 'op-malformed-typed',
    });
    expect(String(error)).not.toContain('raw database diagnostic');
  });

  it('does not treat response loss as proof that activation failed without side effects', async () => {
    vi.mocked(studioV2DatalinkApi.post).mockRejectedValueOnce({
      message: 'raw transport diagnostic', operation_id: 'op-lost-response',
    });
    await expect(studioV2WorkspaceActivationAPI.activate(request)).rejects.toMatchObject({
      outcome: 'unconfirmed', retryable: false, operation_id: 'op-lost-response',
    });
  });

  it('preserves a definite revision rejection and its safe request metadata', async () => {
    vi.mocked(studioV2DatalinkApi.post).mockRejectedValueOnce({ response: {
      status: 409, data: { success: false, error: {
        code: 'modbus_share_revision_conflict', retryable: true,
        request_id: 'req-revision-full', message: 'raw revision diagnostic',
      } },
    } });
    await expect(studioV2WorkspaceActivationAPI.activate(request)).rejects.toMatchObject({
      code: 'modbus_share_revision_conflict', outcome: 'failed', retryable: true, request_id: 'req-revision-full',
    });
  });

  it.each([502, 504])('keeps an HTTP %i gateway failure unconfirmed without automatic retries', async (status) => {
    vi.mocked(studioV2DatalinkApi.post).mockRejectedValueOnce({ response: {
      status, data: { error: { request_id: 'req-gateway-full', retryable: true } },
    } });
    await expect(studioV2WorkspaceActivationAPI.activate(request)).rejects.toMatchObject({
      outcome: 'unconfirmed', retryable: false, request_id: 'req-gateway-full',
    });
  });

  it('does not automatically repeat activation when global mutation retries are enabled', async () => {
    vi.mocked(studioV2DatalinkApi.post).mockRejectedValue(new Error('response lost'));
    const client = new QueryClient({ defaultOptions: { mutations: { retry: 1, retryDelay: 0 } } });
    const wrapper = ({ children }: { children: ReactNode }) => (
      <QueryClientProvider client={client}>{children}</QueryClientProvider>
    );
    const { result } = renderHook(() => useActivateStudioV2WorkspaceMutation(), { wrapper });
    await act(async () => {
      await expect(result.current.mutateAsync(request)).rejects.toMatchObject({ outcome: 'unconfirmed' });
    });
    expect(studioV2DatalinkApi.post).toHaveBeenCalledTimes(1);
  });
});
