import type {
  Device,
  RuntimeStreamConnectionState,
  RuntimeValueEvent,
  SourceRuleRecord,
} from '../../../types/datalink';

export type QuickFlowStepId = 'device' | 'points' | 'output' | 'launch';

export interface QuickFlowStepDefinition {
  id: QuickFlowStepId;
  labelKey: string;
  descriptionKey: string;
}

export interface QuickFlowMetric {
  labelKey: string;
  value: string;
  detail: string;
}

export interface QuickFlowPointCard {
  id: string;
  label: string;
  address: string;
  displayValue: string;
  status: 'ok' | 'warn';
}

export interface QuickFlowOutputTarget {
  id: string;
  title: string;
  detail: string;
  status: 'ok' | 'warn' | 'error';
}

export interface QuickFlowLaunchCheck {
  id: string;
  title: string;
  detail: string;
  status: 'ok' | 'warn' | 'error';
}

export interface QuickFlowPipelineNode {
  id: string;
  title: string;
  detail: string;
  status: 'ok' | 'warn' | 'error' | 'idle';
  metric: string;
  animated: boolean;
}

export interface QuickFlowViewModel {
  device: Device;
  sourceRule: SourceRuleRecord;
  endpoint: string;
  unitLabel: string;
  runtimeState: RuntimeStreamConnectionState;
  metrics: QuickFlowMetric[];
  points: QuickFlowPointCard[];
  outputs: QuickFlowOutputTarget[];
  checks: QuickFlowLaunchCheck[];
  pipelineNodes: QuickFlowPipelineNode[];
  blocker: string;
}

export interface QuickFlowModelInput {
  device: Device;
  sourceRule: SourceRuleRecord;
  runtimeState: RuntimeStreamConnectionState;
  liveValues: Record<string, RuntimeValueEvent>;
}
