import type { RecordingStartOperation } from '../types/studioV2RecordingStart';
import { boundedString } from './safeJson';

function record(value: unknown): value is Record<string, unknown> {
  return typeof value === 'object' && value !== null && !Array.isArray(value);
}

function text(value: unknown): string {
  const result = boundedString(value);
  if (!result) throw new Error('Invalid recording start response');
  return result;
}

function boolean(value: unknown): boolean {
  if (typeof value !== 'boolean') throw new Error('Invalid recording start progress');
  return value;
}

function array(value: unknown): unknown[] {
  if (!Array.isArray(value) || value.length > 64) throw new Error('Invalid recording start scope');
  return value;
}

/** Parses only safe progress facts; ledger owners and internal intent remain server-side. */
export function parseRecordingStartOperation(value: unknown): RecordingStartOperation | null {
  if (!record(value)) return null;
  try {
    const status = text(value.status);
    const stage = text(value.stage);
    const digest = text(value.intent_digest);
    if (value.action !== 'recording_start' ||
      !['pending', 'running', 'succeeded', 'partial', 'failed', 'unknown'].includes(status) ||
      !['save', 'readiness', 'apply', 'activation', 'complete'].includes(stage) ||
      !/^[0-9a-f]{64}$/.test(digest)) return null;
    const groups = array(value.groups).map((entry) => {
      if (!record(entry)) throw new Error('Invalid group progress');
      return {
        group_id: text(entry.group_id), group_revision: text(entry.group_revision),
        ...(entry.applied_revision ? { applied_revision: text(entry.applied_revision) } : {}),
        saved: boolean(entry.saved), ready: boolean(entry.ready), applied: boolean(entry.applied),
      };
    });
    const setupRevision = value.setup_revision === '' && groups.length === 0 ? '' : text(value.setup_revision);
    const devices = array(value.devices).map((entry) => {
      if (!record(entry)) throw new Error('Invalid device progress');
      return {
        device_id: text(entry.device_id), activated: boolean(entry.activated),
        ...(entry.reason ? { reason: text(entry.reason) } : {}),
      };
    });
    const deviceIds = array(value.device_ids).map(text);
    if (deviceIds.length === 0 || new Set(deviceIds).size !== deviceIds.length ||
      new Set(groups.map((group) => group.group_id)).size !== groups.length ||
      new Set(devices.map((device) => device.device_id)).size !== devices.length ||
      devices.length !== deviceIds.length || devices.some((device) => !deviceIds.includes(device.device_id))) return null;
    if (status === 'succeeded' && (stage !== 'complete' ||
      groups.some((group) => !group.saved || !group.ready || !group.applied || !group.applied_revision) ||
      devices.some((device) => !device.activated))) return null;
    return {
      operation_id: text(value.operation_id), action: 'recording_start',
      status: status as RecordingStartOperation['status'],
      stage: stage as RecordingStartOperation['stage'], intent_digest: digest,
      workspace_id: text(value.workspace_id), setup_revision: setupRevision, device_ids: deviceIds,
      groups, devices, created_at: text(value.created_at), updated_at: text(value.updated_at),
      ...(value.reason ? { reason: text(value.reason) } : {}),
      ...(value.next_action ? { next_action: text(value.next_action) } : {}),
    };
  } catch {
    return null;
  }
}
