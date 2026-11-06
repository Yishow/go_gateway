import React from 'react';
import { act, fireEvent, render, renderHook, screen, waitFor } from '@testing-library/react';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { afterEach, describe, expect, it, vi } from 'vitest';
import * as devices from '@/hooks/datalink/useDevices';
import * as points from '@/hooks/datalink/usePoints';
import * as mappings from '@/hooks/datalink/useMappings';
import * as tags from '@/hooks/datalink/useTags';
import * as rules from '@/hooks/datalink/useSourceRules';
import { deviceAPI, pointAPI, mappingAPI, sourceRuleAPI, tagAPI } from '@/services/datalink';
import { studioV2WorkspaceKeys } from '@/hooks/datalink/keys';
import { sourceEntry, sourcePartition } from '@/hooks/datalink/recordingHomeSourceCache';
import { studioV2DatalinkApi } from '@/services/studioV2Workspace';
import { CompleteStep } from '@/components/datalink/wizard/steps/CompleteStep';
const owners = [
  [tags.useCreateTagMutation, tagAPI, 'create', {}],
  [tags.useUpdateTagMutation, tagAPI, 'update', { id: 't', data: {} }],
  [tags.useDeleteTagMutation, tagAPI, 'delete', 't'],
  [tags.useActivateTagMutation, tagAPI, 'activate', 't'],
  [tags.useRetireTagMutation, tagAPI, 'retire', 't'],
  [tags.useBatchCreateTagsMutation, tagAPI, 'batchCreate', []],
  [devices.useCreateDeviceMutation, deviceAPI, 'create', {}],
  [devices.useUpdateDeviceMutation, deviceAPI, 'update', { id: 'd', data: {} }],
  [devices.useDeleteDeviceMutation, deviceAPI, 'delete', 'd'],
  [devices.useToggleDeviceStatusMutation, deviceAPI, 'disable', { id: 'd', currentStatus: 'active' }],
  [devices.useToggleDeviceStatusMutation, deviceAPI, 'activate', { id: 'd', currentStatus: 'disabled' }],
  [points.useCreatePointMutation, pointAPI, 'create', {}],
  [points.useUpdatePointMutation, pointAPI, 'update', { id: 'p', data: {} }],
  [points.useDeletePointMutation, pointAPI, 'delete', 'p'],
  [mappings.useCreateMappingMutation, mappingAPI, 'create', {}],
  [mappings.useUpdateMappingMutation, mappingAPI, 'update', { id: 'm', data: {} }],
  [mappings.useDeleteMappingMutation, mappingAPI, 'delete', 'm'],
  [mappings.useToggleMappingStatusMutation, mappingAPI, 'update', { id: 'm', enabled: true }],
  [rules.useCreateSourceRuleMutation, sourceRuleAPI, 'create', {}],
  [rules.useUpdateSourceRuleMutation, sourceRuleAPI, 'update', { id: 'r', data: {} }],
  [rules.useDeleteSourceRuleMutation, sourceRuleAPI, 'delete', 'r'],
  [rules.useEnableSourceRuleMutation, sourceRuleAPI, 'enable', 'r'],
  [rules.useDisableSourceRuleMutation, sourceRuleAPI, 'disable', 'r'],
] as const;
function setup() {
  const client = new QueryClient({ defaultOptions: { mutations: { retry: false }, queries: { retry: false } } });
  client.setQueryData(studioV2WorkspaceKeys.bootstrap(), { id: 'a' });
  const a = sourceEntry(sourcePartition(client), 'a'), b = sourceEntry(sourcePartition(client), 'b');
  const wrapper = ({ children }: { children: React.ReactNode }) => <QueryClientProvider client={client}>{children}</QueryClientProvider>;
  return { client, a, b, wrapper };
}
afterEach(() => vi.restoreAllMocks());
describe('generic source mutation observation owners', () => {
  for (const [hook, api, method, variables] of owners) {
    it(`${hook.name}/${method} invalidates dispatch scope on ACK without a read barrier`, async () => {
      const { client, a, b, wrapper } = setup();
      let resolve!: (value: unknown) => void;
      const request = vi.spyOn(api as unknown as Record<string, (...args: unknown[]) => Promise<unknown>>, method).mockImplementation(() => new Promise(r => { resolve = r; }) as never);
      const read = vi.spyOn(client, 'fetchQuery');
      const { result } = renderHook(hook as unknown as () => { mutateAsync: (input: unknown) => Promise<unknown> }, { wrapper });
      let pending!: Promise<unknown>;
      act(() => { pending = result.current.mutateAsync(variables as never); });
      await waitFor(() => expect(request).toHaveBeenCalledTimes(1));
      expect(a.generation).toBe(0);
      client.setQueryData(studioV2WorkspaceKeys.bootstrap(), { id: 'b' });
      const ack = { id: 'saved', device_id: 'd' };
      await act(async () => { resolve(ack); expect(await pending).toBe(ack); });
      expect(a.generation).toBe(1);
      expect(b.generation).toBe(0);
      expect(read).not.toHaveBeenCalled();
      client.clear();
    });
    it(`${hook.name}/${method} failure preserves source observation`, async () => {
      const { client, a, b, wrapper } = setup();
      vi.spyOn(api as unknown as Record<string, (...args: unknown[]) => Promise<unknown>>, method).mockRejectedValue(new Error('save failed') as never);
      const { result } = renderHook(hook as unknown as () => { mutateAsync: (input: unknown) => Promise<unknown> }, { wrapper });
      await act(async () => { await expect(result.current.mutateAsync(variables as never)).rejects.toThrow('save failed'); });
      expect([a.generation, b.generation]).toEqual([0, 0]);
      client.clear();
    });
  }
  it('wizard mapping create invalidates after successful save', async () => {
    const { client, a, wrapper } = setup();
    const onComplete = vi.fn();
    vi.spyOn(mappingAPI, 'create').mockResolvedValue({ id: 'm' } as never);
    render(<CompleteStep formData={{ deviceId: 'd', pointId: 'p', tagId: 't', transformPipeline: [], confirmed: true }} onComplete={onComplete} />, { wrapper });
    fireEvent.click(screen.getByRole('button', { name: /建立並啟用映射/ }));
    await waitFor(() => expect(onComplete).toHaveBeenCalledOnce());
    expect(a.generation).toBe(1);
    client.clear();
  });
  for (const [hook, api, method, variables] of [
    [devices.useCheckReadinessMutation, deviceAPI, 'checkReadiness', 'd'],
    [devices.useTestConnectionMutation, deviceAPI, 'testConnection', 'd'],
    [devices.useTestDraftConnectionMutation, deviceAPI, 'testDraftConnection', {}],
    [points.usePollPointMutation, pointAPI, 'pollNow', 'p'],
    [points.useBatchPollPointsMutation, pointAPI, 'batchPoll', ['p']],
    [mappings.usePreviewMappingMutation, mappingAPI, 'preview', {}],
    [mappings.useValidatePipelineMutation, mappingAPI, 'validatePipeline', []],
  ] as const) {
    it(`${hook.name} does not invalidate configuration observations`, async () => {
      const { client, a, b, wrapper } = setup();
      vi.spyOn(api as unknown as Record<string, (...args: unknown[]) => Promise<unknown>>, method).mockResolvedValue({} as never);
      const { result } = renderHook(hook as unknown as () => { mutateAsync: (input: unknown) => Promise<unknown> }, { wrapper });
      await act(async () => { await result.current.mutateAsync(variables as never); });
      expect([a.generation, b.generation]).toEqual([0, 0]);
      client.clear();
    });
  }
  it('wizard mapping failure preserves observation and never completes', async () => {
    const { client, a, wrapper } = setup();
    const onComplete = vi.fn();
    vi.spyOn(mappingAPI, 'create').mockRejectedValue(new Error('save failed'));
    render(<CompleteStep formData={{ deviceId: 'd', pointId: 'p', tagId: 't', transformPipeline: [], confirmed: true }} onComplete={onComplete} />, { wrapper });
    fireEvent.click(screen.getByRole('button', { name: /建立並啟用映射/ }));
    await screen.findByText('save failed');
    expect(a.generation).toBe(0);
    expect(onComplete).not.toHaveBeenCalled();
    client.clear();
  });

  for (const [hook, api, method, variables] of owners.filter(([hook]) =>
    [devices.useToggleDeviceStatusMutation, mappings.useToggleMappingStatusMutation,
      tags.useActivateTagMutation, tags.useRetireTagMutation].includes(hook as typeof devices.useToggleDeviceStatusMutation))) {
    it(`${hook.name}/${method} captures scope and backend before optimistic cancellation`, async () => {
      const { client, a, b, wrapper } = setup();
      const previousBaseURL = studioV2DatalinkApi.defaults.baseURL;
      let release!: () => void;
      const cancel = vi.spyOn(client, 'cancelQueries').mockImplementationOnce(() => new Promise<void>(resolve => { release = resolve; }));
      const request = vi.spyOn(api as unknown as Record<string, (...args: unknown[]) => Promise<unknown>>, method).mockResolvedValue({ id: 'saved' });
      const { result } = renderHook(hook as unknown as () => { mutateAsync: (input: unknown) => Promise<unknown> }, { wrapper });
      try {
        let pending!: Promise<unknown>;
        act(() => { pending = result.current.mutateAsync(variables); });
        await waitFor(() => expect(cancel).toHaveBeenCalledTimes(1));
        expect(request).not.toHaveBeenCalled();
        client.setQueryData(studioV2WorkspaceKeys.bootstrap(), { id: 'b' });
        studioV2DatalinkApi.defaults.baseURL = '/other-backend/api';
        const other = sourceEntry(sourcePartition(client), 'b');
        await act(async () => { release(); await pending; });
        expect(a.generation).toBe(1);
        expect([b.generation, other.generation]).toEqual([0, 0]);
      } finally {
        studioV2DatalinkApi.defaults.baseURL = previousBaseURL;
        client.clear();
      }
    });
  }

});
