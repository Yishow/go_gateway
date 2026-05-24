import type { Device, RuntimeValueEvent, SourceRuleRecord } from '../../../types/datalink';
import type {
  QuickFlowModelInput,
  QuickFlowPointCard,
  QuickFlowStepDefinition,
  QuickFlowViewModel,
} from './quickflowTypes';

export const QUICKFLOW_STEPS: QuickFlowStepDefinition[] = [
  { id: 'device', labelKey: 'quickflow.steps.device.label', descriptionKey: 'quickflow.steps.device.description' },
  { id: 'points', labelKey: 'quickflow.steps.points.label', descriptionKey: 'quickflow.steps.points.description' },
  { id: 'output', labelKey: 'quickflow.steps.output.label', descriptionKey: 'quickflow.steps.output.description' },
  { id: 'launch', labelKey: 'quickflow.steps.launch.label', descriptionKey: 'quickflow.steps.launch.description' },
];

const SAMPLE_VALUES = {
  'point-40021': 74.2,
  'point-40022': 71.9,
  'point-40030': 1480,
};

export const SAMPLE_DEVICE: Device = {
  id: 'sample-device-1',
  name: '混料站 PLC-07',
  description: 'Line 7 mixer',
  protocol: 'modbus_tcp',
  status: 'active',
  connection_config: JSON.stringify({ host: '192.168.18.7', port: 502, unit_id: 1 }),
  last_test_at: '2026-05-09T09:41:11Z',
  last_test_success: true,
  last_test_error: '',
  created_at: '',
  updated_at: '',
};

export const SAMPLE_RULE: SourceRuleRecord = {
  id: 'sample-rule-1',
  device_id: 'sample-device-1',
  start_address: '40021',
  count: 12,
  data_type: 'int16',
  naming_prefix: 'MIX',
  enabled: true,
  locked: false,
  origin: 'manual',
  skipped_addresses: [],
  created_at: '',
  updated_at: '',
};

function parseConnectionConfig(device: Device) {
  try {
    const parsed = JSON.parse(device.connection_config) as Record<string, unknown>;
    return {
      host: typeof parsed.host === 'string' ? parsed.host : '192.168.18.7',
      port: typeof parsed.port === 'number' ? parsed.port : 502,
      unitId:
        typeof parsed.unit_id === 'number'
          ? parsed.unit_id
          : typeof parsed.unitID === 'number'
            ? parsed.unitID
            : 1,
    };
  } catch {
    return { host: '192.168.18.7', port: 502, unitId: 1 };
  }
}

function getValueDisplay(event: RuntimeValueEvent | undefined, fallback: number, unit: string) {
  const raw = event?.transformed_value ?? event?.raw_value ?? fallback;
  const value = typeof raw === 'number' ? raw.toFixed(unit === 'rpm' ? 0 : 1) : String(raw);
  return `${value} ${unit}`;
}

function buildPoints(liveValues: Record<string, RuntimeValueEvent>): QuickFlowPointCard[] {
  return [
    {
      id: 'point-40021',
      label: 'Tank Temp',
      address: '40021',
      displayValue: getValueDisplay(liveValues['point-40021'], SAMPLE_VALUES['point-40021'], '°C'),
      status: 'ok',
    },
    {
      id: 'point-40022',
      label: 'Outlet Temp',
      address: '40022',
      displayValue: getValueDisplay(liveValues['point-40022'], SAMPLE_VALUES['point-40022'], '°C'),
      status: 'ok',
    },
    {
      id: 'point-40030',
      label: 'Motor RPM',
      address: '40030',
      displayValue: getValueDisplay(liveValues['point-40030'], SAMPLE_VALUES['point-40030'], 'rpm'),
      status: 'warn',
    },
  ];
}

function formatProtocolLabel(protocol: Device['protocol']) {
  if (protocol === 'modbus_tcp') return 'Modbus TCP';
  if (protocol === 'modbus_rtu') return 'Modbus RTU';
  if (protocol === 'modbus_udp') return 'Modbus UDP';
  if (protocol === 'fatek_fbs') return 'FATEK FBS';
  return protocol;
}

export function buildQuickFlowViewModel({
  device,
  sourceRule,
  runtimeState,
  liveValues,
}: QuickFlowModelInput): QuickFlowViewModel {
  const connection = parseConnectionConfig(device);
  const protocolLabel = formatProtocolLabel(device.protocol);
  const endpoint = `${protocolLabel} · ${connection.host}:${connection.port}`;
  const unitLabel = `Unit ID ${connection.unitId}`;
  const points = buildPoints(liveValues);
  const healthyFlow = runtimeState === 'connected';

  return {
    device,
    sourceRule,
    endpoint,
    unitLabel,
    runtimeState,
    metrics: [
      { labelKey: 'quickflow.metrics.points', value: '12', detail: 'Holding Register' },
      { labelKey: 'quickflow.metrics.interval', value: '1.0s', detail: 'mixer-fast-lane' },
      { labelKey: 'quickflow.metrics.throughput', value: '24', detail: 'pkt/min' },
      { labelKey: 'quickflow.metrics.lastWrite', value: '09:41:26', detail: 'SQLite committed' },
    ],
    points,
    outputs: [
      { id: 'sqlite', title: 'SQLite', detail: 'mixing_runtime.db · line7_samples', status: 'ok' },
      { id: 'postgres', title: 'PostgreSQL 欄位待補', detail: 'ops_edge.mixer_live · motor_rpm missing', status: 'error' },
    ],
    checks: [
      { id: 'polling', title: '輪詢群組', detail: 'mixer-fast-lane · 1 秒一次', status: 'ok' },
      { id: 'binding', title: '輸出綁定', detail: 'SQLite 已就緒，PostgreSQL 尚有 schema 缺口', status: 'warn' },
      { id: 'launch', title: '啟動策略', detail: '先補欄位再啟動寫入，避免 partial sync', status: 'error' },
    ],
    pipelineNodes: [
      { id: 'connector', title: 'Modbus TCP', detail: `${connection.host}:${connection.port}`, status: healthyFlow ? 'ok' : 'warn', metric: 'connect 18ms', animated: healthyFlow },
      { id: 'probe', title: 'Probe', detail: `Holding Register ${sourceRule.start_address}-${Number(sourceRule.start_address) + sourceRule.count - 1}`, status: healthyFlow ? 'ok' : 'warn', metric: `${sourceRule.count} points`, animated: healthyFlow },
      { id: 'sqlite', title: 'SQLite', detail: 'line7_samples', status: 'ok', metric: 'write ok', animated: healthyFlow },
      { id: 'postgres', title: 'PostgreSQL 欄位待補', detail: 'ops_edge.mixer_live', status: 'error', metric: 'schema mismatch', animated: false },
    ],
    blocker: 'motor_rpm 欄位尚未建立，啟動前需補齊 PostgreSQL schema。',
  };
}
