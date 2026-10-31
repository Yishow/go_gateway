import type { Mapping, Point, Rule } from '../../../features/datalink/workbench-v2/state/types';
import type { SourceRuleRecord } from '../../../types/datalink';
import type { StudioV2RuntimeAppliedRecord } from '../../../types/studioV2RuntimeApply';
import type { StudioV2MappingSaveResult } from '../../../services/studioV2Mappings';
/** Optional request observers. Legacy queue/mutation ownership remains unchanged. */
export interface RuleQueueObservation {
  queued?(rule: Rule, generation: number): void;
  submitted?(rule: Rule, generation: number): void;
  acknowledged(rule: Rule, saved: StudioV2RuntimeAppliedRecord<SourceRuleRecord>): boolean;
  saved(rule: Rule, saved: StudioV2RuntimeAppliedRecord<SourceRuleRecord>): void;
  failed(rule: Rule, error?: unknown): void;
}
export interface MappingQueueObservation {
  acknowledged(point: Point, mapping: Mapping, saved: StudioV2MappingSaveResult): boolean;
  saved(point: Point, mapping: Mapping, saved: StudioV2MappingSaveResult): void;
  failed(point: Point): void;
}

/** Only the server's pre-write validation rejection proves this attempt did not write. */
export function isKnownRuleValidationRejection(error: unknown): boolean {
  if (!error || typeof error !== 'object' || !('response' in error)) return false;
  const response = error.response;
  if (!response || typeof response !== 'object' || Array.isArray(response) || !('status' in response) || response.status !== 400 || !('data' in response)) return false;
  const data = response.data;
  if (!data || typeof data !== 'object' || Array.isArray(data) || !('error' in data)) return false;
  const serverError = data.error;
  // An exact bounded server discriminator, never a client/root fallback or message.
  return Boolean(serverError && typeof serverError === 'object' && !Array.isArray(serverError) && 'code' in serverError && serverError.code === 'validation');
}
