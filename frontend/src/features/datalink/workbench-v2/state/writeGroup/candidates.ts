import type { Point, TargetType, WorkbenchV2State } from '../types';

/** One saved Tag the operator may put into a write group, with the persisted IDs the backend stores. */
export interface GroupCandidate {
  /** Stable key inside the editor: the persisted point and tag IDs. */
  key: string;
  device_id: string;
  point_id: string;
  tag_id: string;
  tag_key: string;
  label: string;
  device_name: string;
  address: string;
  target_type: TargetType;
}

export type CandidateExclusionReason = 'device-not-saved' | 'mapping-not-saved' | 'tag-missing';

export interface ExcludedCandidate {
  label: string;
  reason: CandidateExclusionReason;
}

export interface GroupCandidates {
  candidates: GroupCandidate[];
  excluded: ExcludedCandidate[];
}

function pointLabel(point: Point, tagKey: string): string {
  return tagKey || point.name || point.address;
}

/**
 * Only fully saved sources can join a group: the backend stores persisted
 * device, point and tag IDs, never UI draft IDs. Anything else is reported with
 * the reason so the operator can save it first instead of finding it missing.
 */
export function buildGroupCandidates(state: WorkbenchV2State): GroupCandidates {
  const candidates: GroupCandidate[] = [];
  const excluded: ExcludedCandidate[] = [];
  for (const point of state.points) {
    if (!point.enabled || point.skipped) continue;
    const rule = state.rules.find((candidate) => candidate.id === point.rule_id);
    const mapping = state.mappings[point.id];
    if (rule?.enabled === false || mapping?.enabled === false) continue;
    const label = pointLabel(point, mapping?.tag_key ?? '');
    const device = state.devices.find((candidate) => candidate.id === point.device_id);
    if (!device?.persisted || device.save_state !== 'saved') {
      excluded.push({ label, reason: 'device-not-saved' });
      continue;
    }
    if (!mapping || !mapping.persisted || mapping.save_state !== 'saved' || !mapping.persisted_point_id) {
      excluded.push({ label, reason: 'mapping-not-saved' });
      continue;
    }
    if (!mapping.tag_id) {
      excluded.push({ label, reason: 'tag-missing' });
      continue;
    }
    candidates.push({
      key: `${mapping.persisted_point_id}|${mapping.tag_id}`,
      device_id: point.device_id, point_id: mapping.persisted_point_id, tag_id: mapping.tag_id,
      tag_key: mapping.tag_key, label, device_name: device.name, address: point.address, target_type: mapping.target_type,
    });
  }
  return { candidates, excluded };
}
