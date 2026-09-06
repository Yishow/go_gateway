import { describe, it, expect, vi } from 'vitest';
import { render, screen, fireEvent } from '@testing-library/react';
import { KindSelector } from '../../../src/features/datalink/workbench-v2/steps/step4/KindSelector';
import { WriteStrategy } from '../../../src/features/datalink/workbench-v2/steps/step4/WriteStrategy';
import {
  databaseConnectorNeedsPassword,
  isDatabaseConnectorIdentityChange,
} from '../../../src/features/datalink/workbench-v2/steps/step4/step4DatabaseHelpers';
import {
  hydrateStudioV2DatabaseConnector,
  isStudioV2DatabaseConnectorValid,
} from '../../../src/features/datalink/workbench-v2/state/studioV2DatabaseAutosave';
import type { DbConnector } from '../../../src/features/datalink/workbench-v2/state/types';

// Mock react-i18next
vi.mock('react-i18next', () => ({
  useTranslation: () => ({
    t: (key: string, options?: { interval?: number }) => {
      if (options && options.interval !== undefined) {
        return `${key}_interval_${options.interval}`;
      }
      return key;
    },
  }),
}));

describe('Step 4 Database Components', () => {
  describe('KindSelector', () => {
    it('應渲染 4 個 db kind 卡片，且點擊時呼叫 onChange', () => {
      const onChange = vi.fn();
      render(<KindSelector value="postgres" onChange={onChange} />);

      expect(screen.getByText('step4.kind_postgres')).toBeInTheDocument();
      expect(screen.getByText('step4.kind_sqlite')).toBeInTheDocument();
      expect(screen.getByText('step4.kind_mysql')).toBeInTheDocument();
      expect(screen.getByText('step4.kind_sqlserver')).toBeInTheDocument();

      const sqliteBtn = screen.getByText('step4.kind_sqlite').closest('button');
      expect(sqliteBtn).not.toBeNull();
      fireEvent.click(sqliteBtn!);
      expect(onChange).toHaveBeenCalledWith('sqlite');
    });

    it('在 disabled 為 true 時應禁用所有按鈕', () => {
      render(<KindSelector value="postgres" onChange={() => { }} disabled={true} />);
      const buttons = screen.getAllByRole('button');
      buttons.forEach(btn => {
        expect(btn).toBeDisabled();
      });
    });
  });

  describe('WriteStrategy', () => {
    it('應能渲染兩種策略 radio 與 interval 輸入，點擊時呼叫對應 callback', () => {
      const onWriteModeChange = vi.fn();
      const onWriteIntervalChange = vi.fn();

      render(
        <WriteStrategy
          writeMode="insert"
          writeIntervalSeconds={5}
          onWriteModeChange={onWriteModeChange}
          onWriteIntervalChange={onWriteIntervalChange}
        />
      );

      const upsertRadio = screen.getByLabelText(/step4.write_mode_upsert/);
      fireEvent.click(upsertRadio);
      expect(onWriteModeChange).toHaveBeenCalledWith('upsert');

      const intervalInput = screen.getByRole('spinbutton');
      fireEvent.change(intervalInput, { target: { value: '10' } });
      expect(onWriteIntervalChange).toHaveBeenCalledWith(10);
    });
  });
});

describe('Connector pool selection does not carry a stale credential', () => {
  const baseConnector: DbConnector = {
    kind: 'postgres',
    name: 'Connector A',
    host: 'db-a.internal',
    port: 5432,
    database: 'metrics_a',
    username: 'writer_a',
    password: 'secret-a',
    schema: 'public',
    table: 'sensor_readings',
    write_mode: 'insert',
    write_interval_seconds: 5,
    timestamp_column: 'ts',
    status: 'ready',
  };

  it('flags a password as required when the identity changes and the pool password is redacted', () => {
    const poolConnector = { kind: 'mysql' as const, host: 'db-b.internal', port: 3306, database: 'metrics_b', username: 'writer_b', password: '' };

    expect(isDatabaseConnectorIdentityChange(baseConnector, poolConnector)).toBe(true);
    expect(databaseConnectorNeedsPassword(poolConnector.kind)).toBe(true);
    expect(isStudioV2DatabaseConnectorValid({
      ...baseConnector,
      ...poolConnector,
      password_required: true,
    })).toBe(false);
  });

  it('accepts the connector once the operator re-enters a password', () => {
    expect(isStudioV2DatabaseConnectorValid({
      ...baseConnector,
      host: 'db-b.internal',
      password: 'secret-b',
      password_required: false,
    })).toBe(true);
  });

  it('treats an unchanged identity as no re-entry required', () => {
    expect(isDatabaseConnectorIdentityChange(baseConnector, baseConnector)).toBe(false);
  });

  it('does not require a password for sqlite', () => {
    expect(databaseConnectorNeedsPassword('sqlite')).toBe(false);
  });
});

describe('Hydration preserves the pending credential re-entry flag', () => {
  it('keeps password_required across a backend hydration', () => {
    const current: DbConnector = {
      kind: 'mysql',
      name: 'Connector B',
      host: 'db-b.internal',
      port: 3306,
      database: 'metrics_b',
      username: 'writer_b',
      password: '',
      password_required: true,
      schema: '',
      table: 'sensor_readings',
      write_mode: 'insert',
      write_interval_seconds: 5,
      timestamp_column: 'ts',
      status: 'ready',
    };

    const hydrated = hydrateStudioV2DatabaseConnector({
      id: 'database-config-1',
      workspace_id: 'workspace-1',
      kind: 'mysql',
      name: 'Connector B',
      host: 'db-b.internal',
      port: 3306,
      database: 'metrics_b',
      username: 'writer_b',
      schema: '',
      table: 'sensor_readings',
      write_mode: 'insert',
      write_interval_seconds: 5,
      timestamp_column: 'ts',
      status: 'ready',
      save_state: 'saved',
      created_at: '2026-09-06T00:00:00Z',
      updated_at: '2026-09-06T00:00:00Z',
    }, current);

    expect(hydrated.password_required).toBe(true);
    expect(isStudioV2DatabaseConnectorValid(hydrated)).toBe(false);
  });
});
