export interface RuntimeTruthState {
  state: 'ready' | 'empty' | 'degraded' | 'unavailable' | 'stale' | string;
  empty: boolean;
  degraded: boolean;
  unavailable: boolean;
  stale: boolean;
  reason?: string;
}
