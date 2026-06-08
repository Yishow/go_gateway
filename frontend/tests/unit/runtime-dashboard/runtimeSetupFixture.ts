import type { StudioV2RuntimeSetupContext } from '../../../src/types/studioV2RuntimeContext';

export function createRuntimeSetupFixture(): StudioV2RuntimeSetupContext {
  return {
    readiness_summary: {
      ready: false,
      blocking_count: 2,
      warning_count: 0,
      issues: [{
        code: 'database-target-missing',
        severity: 'blocking',
        step: 'Step 4',
        scope: 'point-1',
        message: 'derived point is missing its persisted database target',
      }, {
        code: 'tag-missing',
        severity: 'blocking',
        step: 'Step 3',
        scope: 'point-2',
        message: 'derived point is missing its persisted tag',
      }],
    },
    source_rules: [{
      id: 'rule-1',
      device_id: 'device-A',
      start_address: '40001',
      count: 8,
      data_type: 'int16',
      naming_prefix: 'LINE_',
      enabled: true,
      revision_id: 'rule-v1',
    }],
    mappings: [{
      id: 'mapping-1',
      rule_id: 'rule-1',
      device_id: 'device-A',
      point_id: 'point-1',
      address: '40001',
      tag_id: 'tag-1',
      tag_key: 'line01.temp.inlet',
      display_name: '入口溫度',
      unit: '°C',
      target_type: 'int16',
      enabled: true,
      status: 'active',
    }],
    database_config: {
      id: 'db-1',
      name: 'PostgreSQL Connector',
      kind: 'postgres',
      database: 'gateway_metrics',
      schema: 'public',
      table: 'sensor_readings',
      write_mode: 'insert',
      write_interval_seconds: 5,
      status: 'ready',
    },
    database_targets: [],
  };
}
