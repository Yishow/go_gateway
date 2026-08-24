import type { Device, Mapping, Rule, WorkbenchV2State } from '../../../features/datalink/workbench-v2/state/types';
import { isStep2Ready } from '../../../features/datalink/workbench-v2/state/sourceRule';
import type { ProtocolType } from '../../../types/datalink';

export function inferHydratedProgress(
  devices: Device[],
  rules: Rule[],
  mappings: Record<string, Mapping>,
): Pick<WorkbenchV2State, 'current' | 'completed'> {
  const deviceProtocolMap: Record<string, ProtocolType> = Object.fromEntries(
    devices.map((device) => [device.id, device.protocol]),
  );
  const step2Ready = isStep2Ready(rules, deviceProtocolMap);
  const completed = new Set<number>();

  if (devices.length > 0) {
    completed.add(1);
  }
  if (step2Ready) {
    completed.add(2);
  }
  if (
    step2Ready &&
    Object.values(mappings).some((mapping) => mapping.persisted && mapping.enabled && mapping.tag_id)
  ) {
    completed.add(3);
  }

  const current: WorkbenchV2State['current'] = completed.has(3)
    ? 4
    : completed.has(2)
    ? 3
    : completed.has(1)
    ? 2
    : 1;
  return { current, completed };
}
