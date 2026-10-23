import { render, screen } from '@testing-library/react';
import { describe, it, expect, vi } from 'vitest';
import { RuntimeDiagnosticsPanel } from '../../../src/features/datalink/runtime-dashboard/RuntimeDiagnosticsPanel';
vi.mock('react-i18next', () => ({ useTranslation: () => ({ t: (key: string, fallback?: unknown) => typeof fallback === 'string' ? fallback : key }) }));
describe('diagnostic evidence boundaries', () => {
  it.each(['collector', 'modbus_share_delivery', undefined])('does not infer SQL commit from %s without an error', stage => {
    render(<RuntimeDiagnosticsPanel selectedDeviceId="a" diagnostics={[{ scope: 'device:a', device_id: 'a', latest_successful_stage: stage, stages: [] }]} />);
    const panel = screen.getByTestId('runtime-dashboard-diagnostics-panel');
    expect(panel).toHaveTextContent('diagnostics.noErrorDetected');
    expect(panel).toHaveTextContent('diagnostics.awaitingWriteEvidence');
    expect(panel).not.toHaveTextContent('所有點位皆已');
    expect(panel).not.toHaveTextContent('database_delivery');
  });
});
