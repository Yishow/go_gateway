import * as React from 'react';
import { useReducer, useCallback } from 'react';
import type { ConnectionTestResult } from '../../../../types/datalink';
import type { WorkbenchV2State, Device, Rule, Mapping, Point } from './types';
import { getDefaultConfig } from './protocols';
import type { ProtocolId } from './types';
import { ruleReducer } from './ruleReducer';
import { mappingReducer } from './mappingReducer';
import { dbReducer } from './dbReducer';
import { settingsReducer } from './settingsReducer';
import type { DbConnector, DbTarget, CommitLog, Settings, SettingsConnector } from './types';

export type WorkbenchV2Action =
  | { type: 'SET_VIEW'; payload: 'flow' | 'settings' }
  | { type: 'SET_CURRENT'; payload: 1 | 2 | 3 | 4 }
  | { type: 'COMPLETE_STEP'; payload: number }
  | { type: 'TOGGLE_SIDEBAR' }
  | { type: 'TOGGLE_SUMMARY_RAIL' }
  | { type: 'SET_SIDEBAR_COLLAPSED'; payload: boolean }
  | { type: 'SET_SHOW_SUMMARY_RAIL'; payload: boolean }
  | { type: 'RESET_FLOW' }
  | { type: 'SET_STATE'; payload: Partial<WorkbenchV2State> }
  | { type: 'addDevice'; device: Device }
  | { type: 'removeDevice'; deviceId: string }
  | { type: 'updateDevice'; deviceId: string; patch: Partial<Device> }
  | { type: 'updateDeviceConfig'; deviceId: string; patch: Partial<Device['config']> }
  | { type: 'renameDevice'; deviceId: string; name: string }
  | { type: 'changeDeviceProtocol'; deviceId: string; protocol: ProtocolId }
  | { type: 'startDeviceTest'; deviceId: string }
  | { type: 'resolveDeviceTest'; deviceId: string; result: ConnectionTestResult }
  | { type: 'completeDeviceTest'; deviceId: string; totalLatency: number }
  | { type: 'failDeviceTest'; deviceId: string; stageId: string; message: string }
  | { type: 'addRule'; rule: Rule }
  | { type: 'removeRule'; ruleId: string }
  | { type: 'updateRule'; ruleId: string; patch: Partial<Rule> }
  | { type: 'renameRule'; ruleId: string; name: string }
  | { type: 'toggleRuleEnabled'; ruleId: string }
  | { type: 'updateRuleSkipped'; ruleId: string; skippedAddresses: string[] }
  | { type: 'toggleRuleSkippedAddress'; ruleId: string; address: string }
  | { type: 'toggleRuleShareEnabled'; ruleId: string }
  | { type: 'updateRuleShareStart'; ruleId: string; shareStart: number | null }
  | { type: 'updateRuleShareStride'; ruleId: string; shareStride: number | null }
  | { type: 'selectRule'; ruleId: string | null }
  | { type: 'initMappingsForPoints'; points: Point[] }
  | { type: 'updateMapping'; pointId: string; patch: Partial<Mapping> }
  | { type: 'toggleMappingEnabled'; pointId: string }
  | { type: 'bulkApplyTransform'; fromPointId: string; fields: ('scale' | 'offset' | 'target_type')[] }
  | { type: 'updateDbConnector'; patch: Partial<DbConnector> }
  | { type: 'upsertDbTarget'; pointId: string; target: DbTarget }
  | { type: 'updateDbTarget'; pointId: string; patch: Partial<DbTarget> }
  | { type: 'autoAssignDbTargets'; targets: Record<string, DbTarget> }
  | { type: 'startCommit' }
  | { type: 'appendCommitLog'; log: CommitLog }
  | { type: 'completeCommit' }
  | { type: 'resetCommit' }
  | { type: 'updateSettings'; patch: Partial<Settings> }
  | { type: 'updateSettingsSection'; section: keyof Settings; patch: Record<string, any> }
  | { type: 'addConnector' }
  | { type: 'updateConnector'; id: string; patch: Partial<SettingsConnector> }
  | { type: 'removeConnector'; id: string }
  | { type: 'startConnectorTest'; id: string }
  | {
      type: 'completeConnectorTest';
      id: string;
      result: {
        status: 'unknown' | 'ready' | 'unreachable' | 'auth_failed' | 'error';
        last_check_at: string;
        last_check_error?: string;
      };
    }
  | { type: 'resetSettingsToDefaults' };

const DEFAULT_DEVICE: Device = {
  id: 'dev-01',
  name: 'PLC-生產線-01',
  description: 'Modbus TCP PLC (Line A 主控)',
  protocol: 'modbus_tcp',
  config: { host: '192.168.1.100', port: 502, slave_id: 1, timeout: 5 },
  status: 'draft',
  test: null,
  persisted: false,
  save_state: 'idle',
  save_error: null,
  runtime_apply_status: null,
  runtime_apply_message: null,
  availability_status: 'available',
  availability_reason: null,
  running: false,
};

const DEFAULT_RULE: Rule = {
  id: 'rule-01',
  device_id: 'dev-01',
  name: 'Holding Registers',
  start_address: '40001',
  count: 8,
  data_type: 'int16',
  naming_prefix: 'SENSOR_',
  enabled: true,
  scale_multiplier: 0.1,
  scale_offset: 0,
  data_format: '',
  skipped_addresses: [],
  share_enabled: true,
  share_start_register: 40001,
  share_stride: null,
  persisted: false,
  save_state: 'idle',
  save_error: null,
};

export const INITIAL_STATE: WorkbenchV2State = {
  view: 'flow',
  current: 1,
  completed: new Set<number>(),
  sidebarCollapsed: false,
  showSummaryRail: true,
  devices: [DEFAULT_DEVICE],
  rules: [DEFAULT_RULE],
  selectedRuleId: 'rule-01',
  points: [],
  mappings: {},
  db: {
    connector: {
      kind: 'postgres',
      name: 'TimeSeries Prod',
      host: 'tsdb.internal',
      port: 5432,
      database: 'gateway_metrics',
      username: 'gw_writer',
      schema: 'public',
      table: 'sensor_readings',
      write_mode: 'insert',
      write_interval_seconds: 5,
      timestamp_column: 'ts',
      status: 'ready',
      persisted: false,
      save_state: 'idle',
      save_error: null,
    },
    targets: {},
  },
  settings: {
    connectors: [
      {
        id: 'conn-prod',
        name: 'TimeSeries Prod',
        kind: 'postgres',
        host: 'tsdb.internal',
        port: 5432,
        database: 'gateway_metrics',
        username: 'gw_writer',
        schema: 'public',
        table: 'sensor_readings',
        enabled: true,
        status: 'ready',
        last_check_at: '2026-05-29T03:00:00.000Z',
        default_write_interval_seconds: 5,
      },
    ],
    timeseries: {
      write_precision: 'millisecond',
      partition_interval: 'daily',
      batch_size: 500,
      retention_days: 90,
    },
    scheduler: {
      default_interval_ms: 1000,
      default_retry_count: 3,
      default_retry_delay_ms: 500,
      breaker_threshold: 10,
      auto_start: true,
    },
    modbus_share: {
      enabled: true,
      bind_address: '0.0.0.0',
      port: 5020,
      slave_id: 1,
      base_register: 40001,
    },
    general: {
      theme: 'dark',
      locale: 'zh-TW',
      addr_format: 'modbus',
      api_base: 'http://localhost:8080',
      api_version: 'v1',
      timeout_seconds: 30,
      log_level: 'info',
      sse_heartbeat_seconds: 15,
      enable_debug_panel: false,
      enable_audit_log: true,
    },
  },
  committed: false,
};

/**
 * 刪除設備時連動清理相關的 Rules、Points、Mappings 與 Database Targets
 * 
 * @param state 當前狀態
 * @param deviceId 被刪除的設備 ID
 * @returns 變更後的狀態
 */
export function cascadeRemoveDevice(state: WorkbenchV2State, deviceId: string): WorkbenchV2State {
  const devices = state.devices.filter((d) => d.id !== deviceId);
  const removedRuleIds = new Set(
    state.rules.filter((r) => r.device_id === deviceId).map((r) => r.id)
  );
  const rules = state.rules.filter((r) => r.device_id !== deviceId);
  const removedPointIds = new Set(
    state.points.filter((p) => p.device_id === deviceId || removedRuleIds.has(p.rule_id)).map((p) => p.id)
  );
  const points = state.points.filter((p) => p.device_id !== deviceId && !removedRuleIds.has(p.rule_id));

  const mappings = { ...state.mappings };
  Object.keys(mappings).forEach((pointId) => {
    if (removedPointIds.has(pointId)) {
      delete mappings[pointId];
    }
  });

  const dbTargets = { ...state.db.targets };
  Object.keys(dbTargets).forEach((pointId) => {
    if (removedPointIds.has(pointId)) {
      delete dbTargets[pointId];
    }
  });

  let selectedRuleId = state.selectedRuleId;
  if (selectedRuleId && removedRuleIds.has(selectedRuleId)) {
    selectedRuleId = rules.length > 0 ? rules[0].id : null;
  }

  return {
    ...state,
    devices,
    rules,
    selectedRuleId,
    points,
    mappings,
    db: {
      ...state.db,
      targets: dbTargets,
    },
  };
}

// Reducer 狀態轉移函數
export function workbenchV2Reducer(state: WorkbenchV2State, action: WorkbenchV2Action): WorkbenchV2State {
  switch (action.type) {
    case 'SET_VIEW':
      return { ...state, view: action.payload };
    case 'SET_CURRENT':
      return { ...state, current: action.payload };
    case 'COMPLETE_STEP': {
      const completed = new Set(state.completed);
      completed.add(action.payload);
      return { ...state, completed };
    }
    case 'TOGGLE_SIDEBAR':
      return { ...state, sidebarCollapsed: !state.sidebarCollapsed };
    case 'TOGGLE_SUMMARY_RAIL':
      return { ...state, showSummaryRail: !state.showSummaryRail };
    case 'SET_SIDEBAR_COLLAPSED':
      return { ...state, sidebarCollapsed: action.payload };
    case 'SET_SHOW_SUMMARY_RAIL':
      return { ...state, showSummaryRail: action.payload };
    case 'RESET_FLOW':
      return {
        ...state,
        view: 'flow',
        current: 1,
        completed: new Set<number>(),
        committed: false,
      };
    case 'SET_STATE':
      return { ...state, ...action.payload };
    case 'addDevice':
      return {
        ...state,
        devices: [...state.devices, action.device],
      };
    case 'removeDevice':
      return cascadeRemoveDevice(state, action.deviceId);
    case 'updateDevice':
      return {
        ...state,
        devices: state.devices.map((d) =>
          d.id === action.deviceId ? { ...d, ...action.patch } : d
        ),
      };
    case 'updateDeviceConfig':
      return {
        ...state,
        devices: state.devices.map((d) =>
          d.id === action.deviceId
            ? { ...d, config: { ...d.config, ...action.patch } }
            : d
        ),
      };
    case 'renameDevice':
      return {
        ...state,
        devices: state.devices.map((d) =>
          d.id === action.deviceId ? { ...d, name: action.name } : d
        ),
      };
    case 'changeDeviceProtocol':
      return {
        ...state,
        devices: state.devices.map((d) =>
          d.id === action.deviceId
            ? {
                ...d,
                protocol: action.protocol,
                config: getDefaultConfig(action.protocol),
                status: 'draft',
                test: null,
              }
            : d
        ),
      };
    case 'startDeviceTest': {
      return {
        ...state,
        devices: state.devices.map((d) =>
          d.id === action.deviceId
            ? {
                ...d,
                status: 'draft',
                test: {
                  status: 'running',
                  stages: {
                    connect: { status: 'running' },
                    probe: { status: 'pending' },
                  },
                },
              }
            : d
        ),
      };
    }
    case 'resolveDeviceTest':
      return {
        ...state,
        devices: state.devices.map((d) => {
          if (d.id !== action.deviceId || !d.test) return d;
          const connectMessage = action.result.connect?.message || action.result.connect?.error;
          const probeMessage = action.result.probe?.message || action.result.probe?.error;
          return {
            ...d,
            status: action.result.success ? 'tested' : 'draft',
            test: {
              status: action.result.success ? 'success' : 'failed',
              latency_ms: action.result.latency_ms,
              stages: {
                connect: {
                  status: action.result.connect?.status ?? 'failed',
                  latency_ms: action.result.connect?.latency_ms,
                  message: connectMessage,
                },
                probe: {
                  status: action.result.probe?.status ?? 'skipped',
                  latency_ms: action.result.probe?.latency_ms,
                  message: probeMessage,
                },
              },
              tested_at: new Date().toISOString(),
            },
          };
        }),
      };
    case 'completeDeviceTest':
      return {
        ...state,
        devices: state.devices.map((d) => {
          if (d.id !== action.deviceId || !d.test) return d;
          return {
            ...d,
            status: 'tested',
            test: {
              ...d.test,
              status: 'success',
              latency_ms: action.totalLatency,
              tested_at: new Date().toISOString(),
            },
          };
        }),
      };
    case 'failDeviceTest':
      return {
        ...state,
        devices: state.devices.map((d) => {
          if (d.id !== action.deviceId || !d.test) return d;
          const updatedStages = {
            ...d.test.stages,
            [action.stageId]: { status: 'failed' as const, message: action.message },
          };
          return {
            ...d,
            status: 'draft',
            test: {
              ...d.test,
              status: 'failed',
              stages: updatedStages,
              tested_at: new Date().toISOString(),
            },
          };
        }),
      };
    case 'selectRule':
    case 'addRule':
    case 'removeRule':
    case 'updateRule':
    case 'renameRule':
    case 'toggleRuleEnabled':
    case 'updateRuleSkipped':
    case 'toggleRuleSkippedAddress':
    case 'toggleRuleShareEnabled':
    case 'updateRuleShareStart':
    case 'updateRuleShareStride':
      return ruleReducer(state, action);
    case 'initMappingsForPoints':
    case 'updateMapping':
    case 'toggleMappingEnabled':
    case 'bulkApplyTransform':
      return mappingReducer(state, action);
    case 'updateDbConnector':
    case 'upsertDbTarget':
    case 'updateDbTarget':
    case 'autoAssignDbTargets':
    case 'startCommit':
    case 'appendCommitLog':
    case 'completeCommit':
    case 'resetCommit':
      return dbReducer(state, action);
    case 'updateSettings':
    case 'updateSettingsSection':
    case 'addConnector':
    case 'updateConnector':
    case 'removeConnector':
    case 'startConnectorTest':
    case 'completeConnectorTest':
    case 'resetSettingsToDefaults':
      return settingsReducer(state, action);
    default:
      return state;
  }
}

/**
 * 封裝了 React useReducer 的 Workbench V2 狀態 Hook
 * 
 * 落地設計決策：「State 管理：本地 React state + 可序列化形狀」
 */
export function useWorkbenchV2State(initialState: WorkbenchV2State = INITIAL_STATE) {
  const [state, dispatch] = useReducer(workbenchV2Reducer, initialState);

  // 在 mount 時從 localStorage 同步
  React.useEffect(() => {
    try {
      if (typeof window !== 'undefined' && window.localStorage) {
        const collapsed = window.localStorage.getItem('wbv2_sidebar_collapsed');
        if (collapsed !== null) {
          dispatch({ type: 'SET_SIDEBAR_COLLAPSED', payload: collapsed === 'true' });
        }
        const showSummary = window.localStorage.getItem('wbv2_show_summary_rail');
        if (showSummary !== null) {
          dispatch({ type: 'SET_SHOW_SUMMARY_RAIL', payload: showSummary === 'true' });
        }
      }
    } catch (e) {
      // eslint-disable-next-line no-console
      console.warn('localStorage is not available:', e);
    }
  }, []);

  const setView = useCallback((view: 'flow' | 'settings') => {
    dispatch({ type: 'SET_VIEW', payload: view });
  }, []);

  const setCurrent = useCallback((current: 1 | 2 | 3 | 4) => {
    dispatch({ type: 'SET_CURRENT', payload: current });
  }, []);

  const completeStep = useCallback((stepId: number) => {
    dispatch({ type: 'COMPLETE_STEP', payload: stepId });
  }, []);

  const toggleSidebar = useCallback(() => {
    dispatch({ type: 'TOGGLE_SIDEBAR' });
  }, []);

  const toggleSummaryRail = useCallback(() => {
    dispatch({ type: 'TOGGLE_SUMMARY_RAIL' });
  }, []);

  const setSidebarCollapsed = useCallback((collapsed: boolean) => {
    dispatch({ type: 'SET_SIDEBAR_COLLAPSED', payload: collapsed });
  }, []);

  const setShowSummaryRail = useCallback((show: boolean) => {
    dispatch({ type: 'SET_SHOW_SUMMARY_RAIL', payload: show });
  }, []);

  const resetFlow = useCallback(() => {
    dispatch({ type: 'RESET_FLOW' });
  }, []);

  const selectRule = useCallback((ruleId: string | null) => {
    dispatch({ type: 'selectRule', ruleId });
  }, []);

  return {
    state,
    setView,
    setCurrent,
    completeStep,
    toggleSidebar,
    toggleSummaryRail,
    setSidebarCollapsed,
    setShowSummaryRail,
    resetFlow,
    selectRule,
    dispatch,
  };
}
