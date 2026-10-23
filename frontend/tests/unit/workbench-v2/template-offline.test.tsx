import type { PropsWithChildren } from 'react';
import { QueryClient, QueryClientProvider, onlineManager } from '@tanstack/react-query';
import { act, renderHook, waitFor } from '@testing-library/react';
import { afterEach, expect, it, vi } from 'vitest';
import { useApplyMeasurementTemplateMutation, usePreviewMeasurementTemplateMutation, useStudioV2MeasurementTemplatesQuery } from '@/hooks/datalink/useStudioV2WorkspaceMeasurements';
import { queryClient } from '@/lib/queryClient';
import { studioV2WorkspaceMeasurementsAPI } from '@/services/studioV2WorkspaceMeasurements';
vi.mock('@/services/studioV2WorkspaceMeasurements', () => ({ studioV2WorkspaceMeasurementsAPI: { previewTemplate: vi.fn(), applyTemplate: vi.fn(), listTemplates: vi.fn() } }));
afterEach(() => { onlineManager.setOnline(true); vi.clearAllMocks(); });
const preview = { device_id: 'A', device_name: 'A', template_id: 'T', confirmed: true, needs_review: false, proposed_items: [], definitions: [] };
function wrapper({ children }: PropsWithChildren) {
  const client = new QueryClient({ defaultOptions: queryClient.getDefaultOptions() });
  return <QueryClientProvider client={client}>{children}</QueryClientProvider>;
}
it.each(['preview', 'apply'] as const)('%s fails offline instead of silently queuing consent until reconnect', async kind => {
  const api = kind === 'preview' ? studioV2WorkspaceMeasurementsAPI.previewTemplate : studioV2WorkspaceMeasurementsAPI.applyTemplate;
  vi.mocked(api).mockRejectedValue(new Error('offline'));
  onlineManager.setOnline(false);
  const useOperation = kind === 'preview' ? usePreviewMeasurementTemplateMutation : useApplyMeasurementTemplateMutation;
  const { result } = renderHook(() => useOperation(), { wrapper });
  act(() => result.current.mutate(preview));
  await waitFor(() => expect(result.current.isError).toBe(true), { timeout: 3000 });
  expect(api).toHaveBeenCalledOnce();
  expect(result.current.isPaused).toBe(false);
  act(() => onlineManager.setOnline(true));
  await new Promise(resolve => setTimeout(resolve, 1100));
  expect(api).toHaveBeenCalledOnce();
});
it('template loading fails offline with an error state instead of staying paused', async () => {
  vi.mocked(studioV2WorkspaceMeasurementsAPI.listTemplates).mockRejectedValue(new Error('offline'));
  onlineManager.setOnline(false);
  const { result } = renderHook(() => useStudioV2MeasurementTemplatesQuery(true), { wrapper });
  await waitFor(() => expect(result.current.isError).toBe(true));
  expect(result.current.fetchStatus).not.toBe('paused');
});
