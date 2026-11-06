import type { GroupCandidates } from '../features/datalink/workbench-v2/state/writeGroup/candidates';
import type { HomeDeviceSummary } from '../features/datalink/workbench-v2/recording-home/recordingHomeCollections';
export interface RecordingHomeSourceSummary extends GroupCandidates { devices: HomeDeviceSummary[] }
export interface RecordingHomeSourceAvailable extends RecordingHomeSourceSummary {
  workspace_id: string; state: 'available'; source_digest: string;
}
export type RecordingHomeSourceObservation = RecordingHomeSourceAvailable | {
  workspace_id: string; state: 'unknown'; reason_code: 'source-inconsistent' | 'source-changed';
};
