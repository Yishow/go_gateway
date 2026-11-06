import type { RecordingHomeSourceAvailable } from '@/types/recordingHomeSources';
/** Independent response fixture for Home; editor full-source reads stay separate. */
export function homeSources(devices: readonly { id: string; name: string; protocol?: string | null }[], workspace_id = 'W'): RecordingHomeSourceAvailable {
  return { workspace_id, state: 'available', source_digest: 'a'.repeat(64), candidates: [], excluded: [],
    devices: devices.map(({ id, name, protocol }) => ({ id, name, protocol: protocol ?? 'unknown', availability: 'unknown' })) };
}
