import type { WorkbenchV2State } from '../state/types';

type RuntimeDashboardDeviceState = Pick<
  WorkbenchV2State,
  'selectedRuleId' | 'rules' | 'devices'
>;

const uuidPattern =
  /^[0-9a-f]{8}-[0-9a-f]{4}-[1-5][0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$/i;

function isPersistedBackendDeviceID(deviceID: string | null | undefined): deviceID is string {
  return typeof deviceID === 'string' && uuidPattern.test(deviceID);
}

export function resolveRuntimeDashboardDevice(
  state: RuntimeDashboardDeviceState,
  confirmedDeviceIds?: readonly string[],
): string | null {
  if (confirmedDeviceIds && confirmedDeviceIds.length > 0) {
    const uniqueConfirmedDeviceIds = Array.from(new Set(confirmedDeviceIds));
    const [confirmedDeviceId] = uniqueConfirmedDeviceIds;
    return uniqueConfirmedDeviceIds.length === 1 && isPersistedBackendDeviceID(confirmedDeviceId)
      ? confirmedDeviceId
      : null;
  }

  if (state.selectedRuleId) {
    const selectedRule = state.rules.find((rule) => rule.id === state.selectedRuleId);
    if (selectedRule && isPersistedBackendDeviceID(selectedRule.device_id)) {
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
  if (
    enabledRuleDeviceIds.length === 1 &&
    isPersistedBackendDeviceID(enabledRuleDeviceIds[0])
  ) {
    return enabledRuleDeviceIds[0];
  }

  if (
    state.devices.length === 1 &&
    isPersistedBackendDeviceID(state.devices[0]?.id)
  ) {
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
