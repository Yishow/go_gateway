import type { WorkbenchV2State } from './types';

/** Removes a device and all draft/runtime state owned by its rules and points. */
export function cascadeRemoveDevice(state: WorkbenchV2State, deviceId: string): WorkbenchV2State {
  const devices = state.devices.filter((d) => d.id !== deviceId);
  const removedRuleIds = new Set(
    state.rules.filter((r) => r.device_id === deviceId).map((r) => r.id)
  );
  const rules = state.rules.filter((r) => r.device_id !== deviceId);
  const removedPointIds = new Set(
    state.points.filter((p) => p.device_id === deviceId || removedRuleIds.has(p.rule_id)).map((p) => p.id)
  );
  const points = state.points.filter((p) => p.device_id !== deviceId && !removedRuleIds.has(p.rule_id));

  const mappings = { ...state.mappings };
  Object.keys(mappings).forEach((pointId) => {
    if (removedPointIds.has(pointId)) {
      delete mappings[pointId];
    }
  });

  const dbTargets = { ...state.db.targets };
  Object.keys(dbTargets).forEach((pointId) => {
    if (removedPointIds.has(pointId)) {
      delete dbTargets[pointId];
    }
  });

  let selectedRuleId = state.selectedRuleId;
  if (selectedRuleId && removedRuleIds.has(selectedRuleId)) {
    selectedRuleId = rules.length > 0 ? rules[0].id : null;
  }

  return {
    ...state,
    devices,
    rules,
    selectedRuleId,
    points,
    mappings,
    db: {
      ...state.db,
      targets: dbTargets,
    },
  };
}
