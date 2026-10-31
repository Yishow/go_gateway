import type { Mapping, Point, Rule } from '../../../features/datalink/workbench-v2/state/types';
import type { SourceRuleRecord } from '../../../types/datalink';
import type { StudioV2RuntimeAppliedRecord } from '../../../types/studioV2RuntimeApply';
import type { StudioV2MappingSaveResult } from '../../../services/studioV2Mappings';
/** Optional request observers. Legacy queue/mutation ownership remains unchanged. */
export interface RuleQueueObservation {
  acknowledged(rule: Rule, saved: StudioV2RuntimeAppliedRecord<SourceRuleRecord>): boolean;
  saved(rule: Rule, saved: StudioV2RuntimeAppliedRecord<SourceRuleRecord>): void;
  failed(rule: Rule): void;
}
export interface MappingQueueObservation {
  acknowledged(point: Point, mapping: Mapping, saved: StudioV2MappingSaveResult): boolean;
  saved(point: Point, mapping: Mapping, saved: StudioV2MappingSaveResult): void;
  failed(point: Point): void;
}
