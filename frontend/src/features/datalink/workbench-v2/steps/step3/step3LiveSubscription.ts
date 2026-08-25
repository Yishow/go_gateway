import type { Point } from '../../state/types';

export interface DevicePointGroup {
  deviceId: string;
  streamPointIds: string[];
}

export interface Step3LiveSubscription {
  localPointId: string;
  deviceId: string;
  address: string;
  streamPointId: string;
}

export function buildSubscriptions(
  points: Point[],
  persistedPointIdByLocalPointId: Record<string, string>,
): Step3LiveSubscription[] {
  return points.map((point) => ({
    localPointId: point.id,
    deviceId: point.device_id,
    address: point.address,
    streamPointId: persistedPointIdByLocalPointId[point.id] ?? point.id,
  }));
}

export function buildGroups(subscriptions: Step3LiveSubscription[]): DevicePointGroup[] {
  const groups = new Map<string, string[]>();

  subscriptions.forEach((subscription) => {
    const current = groups.get(subscription.deviceId) ?? [];
    current.push(subscription.streamPointId);
    groups.set(subscription.deviceId, current);
  });

  return Array.from(groups.entries())
    .sort(([left], [right]) => left.localeCompare(right))
    .map(([deviceId, pointIds]) => ({
      deviceId,
      streamPointIds: pointIds,
    }));
}
