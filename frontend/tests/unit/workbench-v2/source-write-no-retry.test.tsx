import { studioV2RulesAPI, studioV2WorkspaceDevicesAPI, resetDeviceAutosaveMocks } from './device-autosave-page.testHarness';
import { act, renderHook } from '@testing-library/react';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import type { PropsWithChildren } from 'react';
import { beforeEach, expect, it } from 'vitest';
import { queryClient } from '@/lib/queryClient';
import { useCreateStudioV2RuleMutation, useUpdateStudioV2RuleMutation, useDeleteStudioV2RuleMutation } from '@/hooks/datalink/useStudioV2Rules';
import { useCreateStudioV2WorkspaceDeviceMutation, useUpdateStudioV2WorkspaceDeviceMutation, useDeleteStudioV2WorkspaceDeviceMutation, useUpdateStudioV2WorkspaceDeviceAvailabilityMutation, useUpdateStudioV2WorkspaceDeviceOrderMutation } from '@/hooks/datalink/useStudioV2WorkspaceDevices';
function useWrites(): Record<string, () => Promise<unknown>> {
  const create = useCreateStudioV2WorkspaceDeviceMutation(), update = useUpdateStudioV2WorkspaceDeviceMutation();
  const remove = useDeleteStudioV2WorkspaceDeviceMutation(), availability = useUpdateStudioV2WorkspaceDeviceAvailabilityMutation(), order = useUpdateStudioV2WorkspaceDeviceOrderMutation();
  const createRule = useCreateStudioV2RuleMutation(), updateRule = useUpdateStudioV2RuleMutation(), removeRule = useDeleteStudioV2RuleMutation();
  return {
    create: () => create.mutateAsync({ name: 'PLC', protocol: 'modbus_tcp', connection_config: { host: '192.0.2.1', port: 502 } }),
    update: () => update.mutateAsync({ deviceId: 'D', request: { name: 'New PLC' } }),
    remove: () => remove.mutateAsync('D'),
    updateAvailability: () => availability.mutateAsync({ deviceId: 'D', request: { availability_status: 'unavailable', availability_reason: 'invalid' } }),
    updateOrder: () => order.mutateAsync(['D']),
    createRule: () => createRule.mutateAsync({ id: 'R', device_id: 'D', start_address: '40001', count: 1, data_type: 'int16', naming_prefix: 'R', enabled: true }),
    updateRule: () => updateRule.mutateAsync({ ruleId: 'R', request: { enabled: false } }),
    removeRule: () => removeRule.mutateAsync('R'),
  };
}
const operations = [
  ['create', studioV2WorkspaceDevicesAPI.create], ['update', studioV2WorkspaceDevicesAPI.update], ['remove', studioV2WorkspaceDevicesAPI.remove],
  ['updateAvailability', studioV2WorkspaceDevicesAPI.updateAvailability], ['updateOrder', studioV2WorkspaceDevicesAPI.updateOrder],
  ['createRule', studioV2RulesAPI.create], ['updateRule', studioV2RulesAPI.update], ['removeRule', studioV2RulesAPI.remove],
] as const;
beforeEach(() => { resetDeviceAutosaveMocks(); for (const [, service] of operations) service.mockReset(); });
it.each(operations)('%s never automatically retries an uncertain write and still accepts a later explicit operation', async (operation, service) => {
  const client = new QueryClient({ defaultOptions: queryClient.getDefaultOptions() });
  const wrapper = ({ children }: PropsWithChildren) => <QueryClientProvider client={client}>{children}</QueryClientProvider>;
  const hook = renderHook(useWrites, { wrapper });
  const failure = { response: { status: 500, data: { error: { code: 'unknown_failure' } } } };
  service.mockRejectedValue(failure);
  await act(async () => { await expect(hook.result.current[operation]()).rejects.toEqual(failure); });
  expect(service).toHaveBeenCalledTimes(1);
  service.mockResolvedValue({});
  await act(async () => { await hook.result.current[operation](); });
  expect(service).toHaveBeenCalledTimes(2);
});
