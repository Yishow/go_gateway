import type { ReactNode } from 'react';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { act, renderHook } from '@testing-library/react';
import { afterEach, expect, it, vi } from 'vitest';
import { studioV2WorkspaceKeys } from '@/hooks/datalink/keys';
import { useDeleteStudioV2WorkspaceDeviceMutation } from '@/hooks/datalink/useStudioV2WorkspaceDevices';
import { modbusShareStatusKey } from '@/hooks/datalink/useModbusShareStatus';
import { studioV2DatalinkApi } from '@/services/studioV2Workspace';

afterEach(() => vi.restoreAllMocks());

it.each([
  ['success', undefined],
  ['share refresh failed after commit', 'DEVICE_DELETED_SHARE_SCOPE_REFRESH_FAILED'],
  ['connection close failed after commit', 'DEVICE_DELETED_CONNECTION_CLOSE_FAILED'],
])('refreshes device and recording evidence when %s', async (_description, code) => {
  const client = new QueryClient({ defaultOptions: { mutations: { retry: false } } });
  const keys = [studioV2WorkspaceKeys.bootstrap(), studioV2WorkspaceKeys.devices(), studioV2WorkspaceKeys.writeGroups(), modbusShareStatusKey];
  keys.forEach(key => client.setQueryData(key, { before: true }));
  const del = vi.spyOn(studioV2DatalinkApi, 'delete');
  if (code) del.mockRejectedValueOnce({ response: { status: 500, data: { success: false, error: { code } } } });
  else del.mockResolvedValueOnce({ data: { success: true } });
  const wrapper = ({ children }: { children: ReactNode }) => <QueryClientProvider client={client}>{children}</QueryClientProvider>;
  const { result, unmount } = renderHook(() => useDeleteStudioV2WorkspaceDeviceMutation('W1'), { wrapper });
  let error: unknown;
  await act(async () => { try { await result.current.mutateAsync('D1'); } catch (cause) { error = cause; } });
  if (code) expect(error).toMatchObject({ outcome: 'deleted', code, action: 'reload' });
  else expect(error).toBeUndefined();
  for (const key of keys) expect(client.getQueryState(key)?.isInvalidated).toBe(true);
  unmount(); client.clear();
});
