import type { WorkbenchV2State } from '../state/types';

type RuntimeDashboardDeviceState = Pick<
  WorkbenchV2State,
  'selectedRuleId' | 'rules' | 'devices'
>;

export function resolveRuntimeDashboardDevice(
  state: RuntimeDashboardDeviceState,
): string | null {
  if (state.selectedRuleId) {
    const selectedRule = state.rules.find((rule) => rule.id === state.selectedRuleId);
    if (selectedRule) {
      return selectedRule.device_id;
    }
  }

  const enabledRuleDeviceIds = Array.from(
    new Set(
      state.rules
        .filter((rule) => rule.enabled)
        .map((rule) => rule.device_id),
    ),
  );
  if (enabledRuleDeviceIds.length === 1) {
    return enabledRuleDeviceIds[0];
  }

  if (state.devices.length === 1) {
    return state.devices[0].id;
  }

  return null;
}

export function buildRuntimeDashboardTarget(deviceId: string | null): string {
  if (!deviceId) {
    return '/studio/runtime';
  }

  const params = new URLSearchParams({ device_id: deviceId });
  return `/studio/runtime?${params.toString()}`;
}
