export type StudioV2AvailabilityStatus = 'available' | 'unavailable';

export interface StudioV2Availability {
  availability_status: StudioV2AvailabilityStatus;
  availability_reason?: string;
  running: boolean;
}

export interface StudioV2AvailabilityRequest {
  availability_status: StudioV2AvailabilityStatus;
  availability_reason?: string;
}
