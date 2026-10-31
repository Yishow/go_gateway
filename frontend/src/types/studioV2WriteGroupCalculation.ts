import type { WriteGroupMemberCalculation, WriteGroupLayoutVersion } from './studioV2WriteGroup';

/** Backend preparation support is not proof that any saved record is ready. */
export interface WriteGroupCalculationCapabilities {
  version: 1;
  result_shape: 'member-calculation-v1';
  draft_supported: boolean;
  activation_supported: boolean;
  editor_supported: boolean;
  reason: string;
  methods: readonly WriteGroupMemberCalculation[];
  cadence: {
    alignment: 'utc_epoch';
    units: readonly ('second' | 'minute')[];
    minimum: number;
    max_seconds: number;
    max_minutes: number;
    window_bounds_checked: true;
  };
  schema_dialects: readonly ('sqlite' | 'postgres')[];
  schema_layouts: readonly WriteGroupLayoutVersion[];
  unsupported_combinations: readonly {
    method: WriteGroupMemberCalculation;
    source_transform: 'exact_counter_scale';
    phase: 'draft' | 'runtime';
    reason: string;
  }[];
}
