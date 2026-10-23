import { render, screen, cleanup } from '@testing-library/react';
import { afterEach, describe, expect, it, vi } from 'vitest';
import { FocusedDeviceHeader } from '@/features/datalink/runtime-dashboard/FocusedDeviceHeader';
import type { RuntimeStreamConnectionState } from '@/types/datalink';
vi.mock('react-i18next', () => ({ useTranslation: () => ({ t: (key: string, fallback?: string) => typeof fallback === 'string' ? fallback : key }) }));
afterEach(cleanup);
describe('focused header evidence', () => {
  it.each([
    ['waiting-first-row', true, 'available', 'connecting'],
    ['share-only', true, 'available', 'live'],
    ['stopped', false, 'available', 'disconnected'],
    ['delivery-failed', true, 'available', 'error'],
    ['unavailable', false, 'unavailable', 'stale'],
  ] as const)('%s does not infer database delivery', (_case, running, availability, streamState) => {
    const device = { device_id: 'A', name: 'PLC', protocol: 'mc_3e', running, availability_status: availability };
    render(<FocusedDeviceHeader selectedDeviceId="A" selectedDevice={device} devices={[device]} onSelectDevice={() => {}} {...{ streamState: streamState as RuntimeStreamConnectionState }} />);
    expect(screen.getByTestId('runtime-dashboard-header')).not.toHaveTextContent('資料採集與儲存進行中');
    expect(screen.getByTestId('runtime-dashboard-header')).not.toHaveTextContent('即時寫入資料庫');
    expect(screen.getByTestId('runtime-dashboard-header')).not.toHaveTextContent('即時遙測串流已啟動');
    expect(screen.getByTestId('runtime-header-stream')).toHaveTextContent(`header.streamStates.${streamState}`);
    expect(screen.getByTestId('runtime-header-recording-evidence')).toHaveAttribute('href', '/studio/v2?step=4&focus=readiness');
  });
});
