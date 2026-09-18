import { describe, expect, it, vi } from 'vitest';
import { fireEvent, render, screen } from '@testing-library/react';
import { ConnectorSection } from '@/features/datalink/workbench-v2/steps/step4/ConnectorSection';
import { INITIAL_STATE } from '@/features/datalink/workbench-v2/state/useWorkbenchV2State';

vi.mock('react-i18next', () => ({
  useTranslation: () => ({ t: (key: string, fallback?: string | { defaultValue?: string }) => typeof fallback === 'string' ? fallback : fallback?.defaultValue ?? key }),
}));

describe('SavedConnectorIdentity credential controls', () => {
  it('provides explicit password removal and resets removal when typing a replacement', () => {
    const update = vi.fn();
    render(<ConnectorSection connector={{ ...INITIAL_STATE.db.connector, kind: 'postgres' }} onUpdateConnector={update} onKindChange={vi.fn()} />);
    fireEvent.click(screen.getByRole('checkbox', { name: 'step4.clear_password' }));
    expect(update).toHaveBeenLastCalledWith({ clear_password: true, password: '', password_required: false });
    fireEvent.change(screen.getByLabelText('step4.field_password'), { target: { value: ' replacement ' } });
    expect(update).toHaveBeenLastCalledWith({ password: ' replacement ', clear_password: false });
  });

  it('blocks credential removal and disabled saved connections', () => {
    const connector = { ...INITIAL_STATE.db.connector, kind: 'postgres' as const };
    render(<ConnectorSection connector={connector} disabled onUpdateConnector={vi.fn()} onKindChange={vi.fn()} connectors={[{
      ...connector, id: 'disabled-target', enabled: false, status: 'ready', default_write_interval_seconds: 5,
    }]} />);
    expect(screen.getByRole('checkbox', { name: 'step4.clear_password' })).toBeDisabled();
    expect(screen.getByRole('option', { name: /POSTGRES/ })).toBeDisabled();
  });
});
