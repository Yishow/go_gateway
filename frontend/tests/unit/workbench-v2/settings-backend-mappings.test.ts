import { describe, expect, it } from 'vitest';
import { buildUpdateConnectorRequest } from '../../../src/features/datalink/workbench-v2/settings/backendMappings';

describe('settings backend mappings', () => {
  it('includes mysql password when building connector update requests', () => {
    expect(
      buildUpdateConnectorRequest({
        id: 'conn-mysql',
        name: 'MySQL Primary',
        kind: 'mysql',
        host: '127.0.0.1',
        port: 3306,
        database: 'gateway_metrics',
        username: 'root',
        password: 'secret',
        schema: '',
        table: 'sensor_values',
        enabled: true,
        status: 'unknown',
        default_write_interval_seconds: 5,
      }),
    ).toEqual({
      name: 'MySQL Primary',
      kind: 'mysql',
      enabled: true,
      default_write_interval_seconds: 5,
      connection_config: {
        host: '127.0.0.1',
        port: 3306,
        database: 'gateway_metrics',
        user: 'root',
        password: 'secret',
        table: 'sensor_values',
      },
    });
  });
});
