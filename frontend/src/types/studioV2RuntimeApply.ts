export type StudioV2RuntimeApplyStatus = 'not_running' | 'applied' | 'apply_failed';

export interface StudioV2RuntimeApply {
  runtime_apply_status: StudioV2RuntimeApplyStatus;
  runtime_apply_message?: string;
}

export type StudioV2RuntimeAppliedRecord<T> = T & StudioV2RuntimeApply;
