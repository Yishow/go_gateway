import { pointAPI, pollingGroupAPI } from '../../../../services/datalink';
import type { SourceIdentity } from '../../../../pages/datalink/workbench-v2/rangeSourceSubmission';
export interface RecordingRangeScheduleRow { id: string; name: string; intervalMs: number; enabled: boolean; points: number }
/** Display-only snapshot. Never grants save, selection, activation or runtime authority. */
export async function readRecordingRangeSchedule(deviceId: string, sources: readonly SourceIdentity[]): Promise<RecordingRangeScheduleRow[] | undefined> {
  if (!sources.length || sources.some(source => source.device_id !== deviceId || !source.point_id) ||
    new Set(sources.map(source => source.point_id)).size !== sources.length) return;
  const [points, groups] = await Promise.all([Promise.all(sources.map(source => pointAPI.get(source.point_id))), pollingGroupAPI.list()]);
  const rows = new Map<string, RecordingRangeScheduleRow>();
  for (const [index, point] of points.entries()) {
    if (point.id !== sources[index].point_id || point.device_id !== deviceId || !point.polling_group_id) return;
    const matching = groups.filter(group => group.id === point.polling_group_id);
    if (matching.length !== 1 || !Number.isFinite(matching[0].interval_ms) || matching[0].interval_ms <= 0) return;
    const group = matching[0], previous = rows.get(group.id);
    rows.set(group.id, { id: group.id, name: group.name, intervalMs: group.interval_ms, enabled: group.enabled, points: (previous?.points ?? 0) + 1 });
  }
  return [...rows.values()];
}
