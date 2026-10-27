import { describe, expect, it } from 'vitest';
import { getSafeErrorMessage } from '../../../src/utils/typedErrors';

const mockT = (key: string, options?: Record<string, unknown>) => {
  const translations: Record<string, string> = {
    'errors.partition_unavailable': 'Native monthly partition objects are unavailable; choose plain storage.',
    'errors.select_plain_storage': 'Choose plain storage',
    'errors.adapter_unverified': 'The SQL Server recording adapter is unverified for this operation.',
    'errors.exact_value_unrepresentable': 'The source value cannot be represented in the configured destination format.',
    'errors.schema_operation_partial': 'The schema operation partially succeeded; verify the database catalog.',
    'errors.storage_owner_mismatch': 'The target table is owned by a different configuration.',
    'errors.partition_maintenance_blocked': 'Partition maintenance is blocked. Check partition function and scheme configuration.',
    'errors.generic_failure': `Generic failure. Request ID: ${options?.requestId || 'unknown'}`,
  };
  return translations[key] ?? (options?.defaultValue as string) ?? key;
};

describe('SQL Server typed error presentation', () => {
  it('maps partition_unavailable and action select_plain_storage safely without leaking raw DSN text', () => {
    const errorPayload = {
      response: {
        data: {
          success: false,
          error: {
            code: 'partition_unavailable',
            message: 'sqlserver://sa:secret_password@10.0.0.12:1433?database=gw_corp failed: DDL denied',
            action: 'select_plain_storage',
            request_id: 'req-sql-partition-404',
            retryable: false,
          },
        },
      },
    };

    const result = getSafeErrorMessage(errorPayload, mockT);
    expect(result.code).toBe('partition_unavailable');
    expect(result.action).toBe('Choose plain storage');
    expect(result.message).toBe('Native monthly partition objects are unavailable; choose plain storage.');
    expect(result.message).not.toContain('secret_password');
    expect(result.message).not.toContain('sqlserver://');
    expect(result.requestId).toBe('req-sql-partition-404');
    expect(result.retryable).toBe(false);
  });

  it('maps actual SQL Server contract error codes without leaking internal diagnostics', () => {
    for (const code of [
      'adapter_unverified',
      'exact_value_unrepresentable',
      'schema_operation_partial',
      'storage_owner_mismatch',
      'partition_maintenance_blocked',
    ]) {
      const errorPayload = {
        code,
        message: 'private connection diagnostic with raw credentials server=sql1;user=admin;pwd=secret',
        request_id: `req-${code}`,
        retryable: false,
      };

      const result = getSafeErrorMessage(errorPayload, mockT);
      expect(result.code).toBe(code);
      expect(result.title).toBe(code);
      expect(result.message).not.toContain('private connection diagnostic');
      expect(result.message).not.toContain('secret');
      expect(result.requestId).toBe(`req-${code}`);
    }
  });
});
