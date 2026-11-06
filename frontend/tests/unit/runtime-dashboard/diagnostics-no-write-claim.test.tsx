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
  it('keeps an accepted database snapshot awaiting write evidence', () => {
    render(<RuntimeDiagnosticsPanel selectedDeviceId="a" databaseDelivery={[{ device_id: 'a', point_id: 'p1', tag_id: 'tag-1',
      status: 'accepted', stages: ['collected', 'mapped', 'buffer_accepted'], observed_at: '2026-11-04T10:00:01Z' }]} />);
    const panel = screen.getByTestId('runtime-dashboard-diagnostics-panel');
    expect(panel).toHaveTextContent('diagnostics.awaitingWriteEvidence');
    expect(panel).not.toHaveTextContent('db_write_succeeded');
    expect(panel).toHaveTextContent('Latest successful stage: Unavailable');
  });
});
