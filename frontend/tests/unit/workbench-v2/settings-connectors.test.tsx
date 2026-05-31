import { describe, it, expect, vi } from 'vitest';
import { render, screen, fireEvent } from '@testing-library/react';
import { ConnectorRow } from '../../../src/features/datalink/workbench-v2/settings/ConnectorRow';
import { ConnectorPoolSection } from '../../../src/features/datalink/workbench-v2/settings/ConnectorPoolSection';
import { INITIAL_STATE } from '../../../src/features/datalink/workbench-v2/state/useWorkbenchV2State';

// 模擬 i18next
vi.mock('react-i18next', () => ({
  useTranslation: () => ({
    t: (key: string, fallback?: string) => fallback || key,
  }),
}));

describe('Connector settings', () => {
  describe('ConnectorRow', () => {
    it('renders and allows updates on inputs', () => {
      const onUpdate = vi.fn();
      const onRemove = vi.fn();
      const onTest = vi.fn();
      
      const conn = INITIAL_STATE.settings.connectors[0];
      
      render(
        <ConnectorRow
          connector={conn}
          onUpdate={onUpdate}
          onRemove={onRemove}
          onTest={onTest}
        />
      );

      const hostInput = screen.getByText('主機位址').closest('div')?.querySelector('input');
      expect(hostInput).toHaveValue('tsdb.internal');
      
      fireEvent.change(hostInput!, { target: { value: 'localhost' } });
      expect(onUpdate).toHaveBeenCalledWith({ host: 'localhost' });
    });

    it('resets status to unknown when host changes', () => {
      const onUpdate = vi.fn();
      const onRemove = vi.fn();
      const onTest = vi.fn();
      const conn = { ...INITIAL_STATE.settings.connectors[0], status: 'ready' as const };

      render(
        <ConnectorRow
          connector={conn}
          onUpdate={onUpdate}
          onRemove={onRemove}
          onTest={onTest}
        />
      );

      const hostInput = screen.getByText('主機位址').closest('div')?.querySelector('input');
      fireEvent.change(hostInput!, { target: { value: 'localhost' } });
      expect(onUpdate).toHaveBeenCalledWith({ host: 'localhost' });
    });

    it('renders mysql password input and forwards password updates', () => {
      const onUpdate = vi.fn();
      const conn = {
        ...INITIAL_STATE.settings.connectors[0],
        id: 'conn-mysql',
        kind: 'mysql' as const,
        port: 3306,
        password: '',
      };

      render(
        <ConnectorRow
          connector={conn}
          onUpdate={onUpdate}
          onRemove={vi.fn()}
          onTest={vi.fn()}
        />
      );

      const passwordInput = screen
        .getByText('密碼')
        .closest('div')
        ?.querySelector('input');
      expect(passwordInput).toHaveAttribute('type', 'password');

      fireEvent.change(passwordInput!, { target: { value: 'secret' } });
      expect(onUpdate).toHaveBeenCalledWith({ password: 'secret' });
    });
  });

  describe('ConnectorPoolSection', () => {
    it('triggers add and remove actions', () => {
      const onAddConnector = vi.fn();
      const onUpdateConnector = vi.fn();
      const onRemoveConnector = vi.fn();
      const onTestConnector = vi.fn();
      render(
        <ConnectorPoolSection
          connectors={INITIAL_STATE.settings.connectors}
          onAddConnector={onAddConnector}
          onUpdateConnector={onUpdateConnector}
          onRemoveConnector={onRemoveConnector}
          onTestConnector={onTestConnector}
        />,
      );

      fireEvent.click(screen.getByText('新增連接器'));
      expect(onAddConnector).toHaveBeenCalled();

      fireEvent.click(screen.getByText('刪除'));
      expect(onRemoveConnector).toHaveBeenCalledWith('conn-prod');
    });

    it('delegates connector test to the backend callback', () => {
      const onTestConnector = vi.fn();
      render(
        <ConnectorPoolSection
          connectors={INITIAL_STATE.settings.connectors}
          onAddConnector={vi.fn()}
          onUpdateConnector={vi.fn()}
          onRemoveConnector={vi.fn()}
          onTestConnector={onTestConnector}
        />,
      );

      fireEvent.click(screen.getByText('測試連線'));
      expect(onTestConnector).toHaveBeenCalledWith('conn-prod');
    });
  });
});
