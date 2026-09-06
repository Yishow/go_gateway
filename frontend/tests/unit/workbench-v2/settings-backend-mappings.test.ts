import { describe, expect, it } from 'vitest';
import { buildUpdateConnectorRequest } from '../../../src/features/datalink/workbench-v2/settings/backendMappings';
import { getDbKindPatch } from '../../../src/features/datalink/workbench-v2/state/dbSchemas';

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

describe('Connector credential is not carried across a kind change', () => {
  const baseConnector = {
    id: 'conn-1',
    name: 'Pool connector',
    kind: 'postgres' as const,
    host: 'db.internal',
    port: 5432,
    database: 'metrics',
    username: 'writer',
    password: 'postgres-secret',
    schema: 'public',
    table: 'sensor_readings',
    enabled: true,
    status: 'ready' as const,
    default_write_interval_seconds: 5,
  };

  it('clears the password locally when the database kind changes', () => {
    const patch = getDbKindPatch('mysql', baseConnector);

    expect(patch.password).toBe('');
    expect(patch.kind).toBe('mysql');
  });

  it('asks the backend to clear the stored password while none has been re-entered', () => {
    const request = buildUpdateConnectorRequest({
      ...baseConnector,
      kind: 'mysql',
      password: '',
      password_required: true,
    });

    expect(request.clear_password).toBe(true);
    expect(request.connection_config).not.toHaveProperty('password');
  });

  it('sends the new password and stops asking for a clear once re-entered', () => {
    const request = buildUpdateConnectorRequest({
      ...baseConnector,
      kind: 'mysql',
      password: 'mysql-secret',
      password_required: false,
    });

    expect(request.clear_password).toBeUndefined();
    expect(request.connection_config).toMatchObject({ password: 'mysql-secret' });
  });
});
