import type { WriteGroup } from '../../../../types/studioV2WriteGroup';
import type { WorkbenchV2State } from '../state/types';
import { buildGroupCandidates } from '../state/writeGroup/candidates';

/** Only explicit saved managed membership; never adds other workspace devices or converts a group. */
export function savedRecordingStartDevices(group: WriteGroup): string[] {
  const ids = [...new Set(group.members.map(member => member.device_id))].sort();
  return group.destination.storage_strategy === 'managed' && group.status !== 'deleted' && ids.length > 0 && ids.length <= 64 &&
    ids.every(id => Boolean(id) && !id.includes('/')) ? ids : [];
}
export function savedRecordingStartSourceIdentity(group: WriteGroup, state: WorkbenchV2State): string {
  const ids = new Set(group.members.map(member => member.device_id));
  return JSON.stringify([state.db.connector, state.devices.filter(device => ids.has(device.id)),
    state.rules.filter(rule => ids.has(rule.device_id)),
    buildGroupCandidates(state).candidates.filter(candidate => group.members.some(member => member.device_id === candidate.device_id && member.point_id === candidate.point_id && member.tag_id === candidate.tag_id))]);
}
export function savedRecordingStartSourcesSaved(group: WriteGroup, state: WorkbenchV2State): boolean {
  const ids = savedRecordingStartDevices(group), connector = state.db.connector;
  const candidates = buildGroupCandidates(state).candidates;
  return Boolean(ids.length && ids.every(id => state.devices.some(device => device.id === id && device.persisted && device.save_state === 'saved')) &&
    connector.persisted && connector.save_state === 'saved' &&
    connector.connector_id === group.destination.connector_id && connector.identity_revision === group.destination.connector_revision &&
    state.rules.filter(rule => ids.includes(rule.device_id)).every(rule => rule.persisted && rule.save_state === 'saved') &&
    group.members.every(member => candidates.some(candidate => candidate.device_id === member.device_id && candidate.point_id === member.point_id && candidate.tag_id === member.tag_id)));
}
