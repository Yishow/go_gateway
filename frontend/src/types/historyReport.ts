export type Resolution = 'raw' | '1m' | '1h' | '1d';

export interface HistoryQuery {
  workspace_id?: string;
  plan_id: string;
  measurement_ids?: string[];
  start_time?: string;
  end_time?: string;
  resolution?: Resolution;
  limit?: number;
  cursor?: string;
  revision_cutoff?: number;
}

export interface HistoryPoint {
  measurement_id: string;
  observed_at: string;
  value_numeric?: number;
  value_string?: string;
  time_weighted_mean?: number;
  sampled_min?: number;
  sampled_max?: number;
  usage_delta?: number;
  quality: string;
  coverage_ratio: number;
  is_estimated: boolean;
  is_provisional: boolean;
}

export interface HistoryReport {
  workspace_id: string;
  plan_id: string;
  points: HistoryPoint[];
  total_count: number;
  next_cursor?: string;
  generated_at: string;
}
