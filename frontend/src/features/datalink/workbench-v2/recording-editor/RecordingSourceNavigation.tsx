import { createContext, useContext, useSyncExternalStore, type ReactNode } from 'react';
import type { DeviceSourcePort } from '../../../../pages/datalink/workbench-v2/deviceSourceSubmission';
import type { PointOrderObservation, RecordingPointOrderPort } from './useRecordingPointOrder';
import { useRecordingDraftGuard } from './RecordingDraftNavigation';
const empty: readonly PointOrderObservation[] = [];
const noSubscribe = () => () => undefined;
interface Owner { workspaceId: string; pointOrder?: RecordingPointOrderPort; observations?: DeviceSourcePort['observations'] }
const SourceOwner = createContext<Owner | undefined>(undefined);
interface Props { workspaceId: string; source?: DeviceSourcePort; pointOrder?: RecordingPointOrderPort; children: ReactNode }
function SourceGuard({ workspaceId, source, pointOrder }: Omit<Props, 'children'>) {
  const points = useSyncExternalStore(pointOrder?.subscribe ?? noSubscribe, () => pointOrder?.observations?.() ?? empty);
  const currentPoints = points.filter(value => value.scope.workspace_id === workspaceId);
  const device = source?.observations().filter(value => value.request.workspaceId === workspaceId && ['pending', 'unknown'].includes(value.status)) ?? [];
  const ranges = source?.range?.observations().filter(value => value.request.workspaceId === workspaceId && ['pending', 'unknown', 'persisted-unavailable'].includes(value.status)) ?? [];
  useRecordingDraftGuard({
    version: JSON.stringify([workspaceId, currentPoints.map(({ scope, entry }) => [scope, entry.status, entry.intent, entry.operation?.request_id]),
      device.map(value => [value.request.requestId, value.status]), ranges.map(value => [value.request.requestId, value.status])]),
    dirty: currentPoints.some(({ entry }) => entry.status === 'ready' && Boolean(entry.intent)),
    // A confirmed format save is terminal configuration evidence even while runtime apply is pending.
    pending: device.length > 0 || ranges.length > 0 || currentPoints.some(({ entry }) => ['saving', 'unknown', 'inspecting', 'recovering', 'recovery-unknown'].includes(entry.status)),
    discard: () => { for (const { scope } of currentPoints) pointOrder?.discardUnsent?.(scope); },
  });
  return null;
}
/** Page lifetime, not row lifetime. Only observes the original owners and discards unsent point intent.
 * Nested standalone editors share this registration instead of creating another source owner. */
export function RecordingSourceNavigation({ workspaceId, source, pointOrder = source?.pointOrder, children }: Props) {
  const inherited = useContext(SourceOwner);
  if ((!source && !pointOrder) || inherited?.workspaceId === workspaceId && inherited.pointOrder === pointOrder && inherited.observations === source?.observations) return children;
  return <SourceOwner.Provider value={{ workspaceId, pointOrder, observations: source?.observations }}>
    <SourceGuard workspaceId={workspaceId} source={source} pointOrder={pointOrder} />{children}
  </SourceOwner.Provider>;
}
