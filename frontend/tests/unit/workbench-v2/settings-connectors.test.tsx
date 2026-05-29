import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest';
import { render, screen, fireEvent, act } from '@testing-library/react';
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
  });

  describe('ConnectorPoolSection', () => {
    beforeEach(() => {
      vi.useFakeTimers();
    });

    afterEach(() => {
      vi.useRealTimers();
      vi.restoreAllMocks();
    });

    it('triggers add and remove actions', () => {
      const dispatch = vi.fn();
      render(<ConnectorPoolSection connectors={INITIAL_STATE.settings.connectors} dispatch={dispatch} />);

      fireEvent.click(screen.getByText('新增連接器'));
      expect(dispatch).toHaveBeenCalledWith({ type: 'addConnector' });

      fireEvent.click(screen.getByText('刪除'));
      expect(dispatch).toHaveBeenCalledWith({ type: 'removeConnector', id: 'conn-prod' });
    });

    it('performs successful mock connection test', () => {
      const dispatch = vi.fn();
      
      // Spy Math.random to return 0.5 (which is >= 0.15, meaning success)
      vi.spyOn(Math, 'random').mockReturnValue(0.5);

      render(<ConnectorPoolSection connectors={INITIAL_STATE.settings.connectors} dispatch={dispatch} />);

      fireEvent.click(screen.getByText('測試連線'));
      expect(dispatch).toHaveBeenCalledWith({ type: 'startConnectorTest', id: 'conn-prod' });

      // Fast forward time by 900ms
      act(() => {
        vi.advanceTimersByTime(900);
      });

      expect(dispatch).toHaveBeenLastCalledWith({
        type: 'completeConnectorTest',
        id: 'conn-prod',
        result: expect.objectContaining({
          status: 'ready',
        }),
      });
    });

    it('performs failed mock connection test', () => {
      const dispatch = vi.fn();
      
      // Spy Math.random to return 0.05 (which is < 0.15, meaning fail)
      vi.spyOn(Math, 'random').mockReturnValue(0.05);

      render(<ConnectorPoolSection connectors={INITIAL_STATE.settings.connectors} dispatch={dispatch} />);

      fireEvent.click(screen.getByText('測試連線'));
      expect(dispatch).toHaveBeenCalledWith({ type: 'startConnectorTest', id: 'conn-prod' });

      // Fast forward time by 900ms
      act(() => {
        vi.advanceTimersByTime(900);
      });

      expect(dispatch).toHaveBeenLastCalledWith({
        type: 'completeConnectorTest',
        id: 'conn-prod',
        result: expect.objectContaining({
          status: 'unreachable',
        }),
      });
    });
  });
});
