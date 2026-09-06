export type SemanticKind =
  | 'gauge'
  | 'rate'
  | 'counter'
  | 'signed_counter'
  | 'delta'
  | 'state'
  | 'event'
  | 'text';

export type QualityFlag = 'good' | 'missing' | 'stale' | 'invalid' | 'uncertain' | 'bad';

export interface CounterPolicy {
  reset_threshold?: number;
  max_roll_over?: number;
  allow_negative: boolean;
}

export interface MeasurementDefinition {
  id: string;
  workspace_id: string;
  device_id: string;
  point_id: string;
  tag_id?: string;
  equipment_id: string;
  definition_revision: string;
  source_binding_revision: string;
  series_epoch: string;
  name: string;
  quantity: string;
  unit?: string;
  semantic_kind: SemanticKind;
  numeric_encoding?: string;
  counter_policy?: CounterPolicy;
  state_map?: Record<string, string>;
  bitmask_labels?: Record<number, string>;
  created_at?: string;
  updated_at?: string;
}

export interface SampleEnvelope {
  sample_id: string;
  workspace_id?: string;
  measurement_id: string;
  series_epoch: string;
  definition_revision: string;
  source_binding_revision: string;
  acquisition_id: string;
  source_sequence?: number;
  observed_at: string;
  received_at: string;
  time_origin: string;
  value_type: string;
  value: unknown;
  value_str?: string;
  raw_value?: unknown;
  quality: QualityFlag;
  quality_reason?: string;
}

export interface MixedItem {
  item_id: string;
  name: string;
  address?: string;
  register_offset: number;
  bit_offset?: number;
  bit_length?: number;
  data_type: string;
  data_format?: string;
  target_data_type?: string;
  scale_multiplier?: number;
  scale_offset?: number;
  quantity?: string;
  unit?: string;
  semantic_kind: SemanticKind;
  counter_policy?: CounterPolicy;
  state_map?: Record<string, string>;
  bitmask_labels?: Record<number, string>;
}

export interface MeasurementTemplate {
  id: string;
  version: string;
  name: string;
  description: string;
  layout_mode: 'homogeneous' | 'mixed';
  items: MixedItem[];
}

export interface TemplateApplyPreview {
  device_id: string;
  device_name: string;
  template_id: string;
  proposed_items: MixedItem[];
  definitions: MeasurementDefinition[];
  confirmed: boolean;
  needs_review: boolean;
  review_notes?: string[];
}
