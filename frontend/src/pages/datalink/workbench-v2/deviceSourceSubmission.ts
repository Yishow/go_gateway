import type { RecordingPointNamePort } from './recordingPointNameSubmission';
import type { RecordingPointOrderPort } from '../../../features/datalink/workbench-v2/recording-editor/useRecordingPointOrder';
import type { RecentConnectionPort } from '../../../features/datalink/workbench-v2/recording-editor/recentConnectionDefaults';
import type { RangeSourcePort } from './rangeSourceSubmission';
import type { Device } from '../../../features/datalink/workbench-v2/state/types';
import type { DeviceAcquisitionInput } from '../../../types/acquisitionSettings';
import type { StudioV2WorkspaceDeviceRecord } from '../../../services/studioV2WorkspaceDevices';
import type { StudioLocatorEntryRead } from '../../../types/studioV2LocatorEntry';
import type { WriteGroupMutationResponse } from '../../../types/studioV2WriteGroup';

export interface LocatorGroupSaveAcknowledgement {
  previous: StudioLocatorEntryRead;
  acknowledgement: WriteGroupMutationResponse;
  current: StudioLocatorEntryRead;
}

export type DeviceSourceBuffer = Pick<Device, 'name' | 'description' | 'protocol' | 'config' | 'acquisition'> & {
  /** Raw form text for saved-device acquisition settings; replaced by `acquisition` when parsed. */
  acquisitionInput?: DeviceAcquisitionInput;
  /** Raw random limit text; parsed into `config.acquisition.limits.random_points`. */
  randomPointsInput?: string;
};
export interface DeviceSourceRequest {
  workspaceId: string;
  sessionId: string;
  requestId: string;
  generation: number;
  deviceId?: string;
  buffer: DeviceSourceBuffer;
}
export interface DeviceSourceOutcome {
  request: DeviceSourceRequest;
  status: 'pending' | 'confirmed' | 'failed' | 'unknown' | 'blocked' | 'persisted-unavailable';
  device?: StudioV2WorkspaceDeviceRecord;
}
/** Local observations, not a backend operation ledger or CAS permission. */
export interface DeviceSourcePort {
  name?: RecordingPointNamePort;
  pointOrder?: RecordingPointOrderPort;
  range?: RangeSourcePort;
  /** Page-session prefill only; no persistence or source-write authority. */
  recentConnections?: RecentConnectionPort;
  /** Advance an own-save read basis only; never grants a source write. */
  acknowledgeLocatorGroupSave?(proof: LocatorGroupSaveAcknowledgement): boolean;
  prepare?(workspaceId: string): Promise<boolean>;
  submit(request: DeviceSourceRequest): Promise<DeviceSourceOutcome>;
  inspect(requestId: string): Promise<DeviceSourceOutcome | undefined>;
  observations(): readonly DeviceSourceOutcome[];
}
