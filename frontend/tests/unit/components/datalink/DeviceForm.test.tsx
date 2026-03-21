import { fireEvent, render, screen, waitFor } from '@testing-library/react';
import { beforeEach, describe, expect, it, vi } from 'vitest';

import DeviceForm from '../../../../src/components/datalink/DeviceForm';
import type { Device, ProtocolInfo } from '../../../../src/types/datalink';

const protocolList: ProtocolInfo[] = [
  {
    type: 'mc_3e',
    name: 'Mitsubishi MC 3E',
    description: '三菱 MC 3E 協議 (Binary)',
    config_schema: '{}',
  },
];

const mocks = vi.hoisted(() => ({
  listProtocolsMock: vi.fn(),
  getSettingsMock: vi.fn(),
}));

vi.mock('../../../../src/services/datalink', () => ({
  protocolAPI: {
    list: mocks.listProtocolsMock,
  },
  settingsAPI: {
    get: mocks.getSettingsMock,
  },
}));

vi.mock('react-i18next', () => ({
  useTranslation: () => ({
    t: (key: string) => key,
  }),
}));

vi.mock('../../../../src/utils/logger', () => ({
  logger: {
    error: vi.fn(),
  },
}));

describe('DeviceForm', () => {
  beforeEach(() => {
    mocks.listProtocolsMock.mockResolvedValue(protocolList);
    mocks.getSettingsMock.mockResolvedValue({
      default_retry_count: 3,
      default_retry_delay: 1000,
    });
  });

  it('submits legacy mc_3e data_format values as supported byte order strings', async () => {
    const onSubmit = vi.fn().mockResolvedValue(undefined);
    const device: Device = {
      id: 'device-1',
      name: 'MC Device',
      description: '',
      protocol: 'mc_3e',
      status: 'draft',
      connection_config: JSON.stringify({
        host: '10.0.0.20',
        port: 5000,
        data_format: 'binary',
      }),
      last_test_at: null,
      last_test_success: null,
      last_test_error: '',
      created_at: '',
      updated_at: '',
    };

    render(<DeviceForm device={device} onSubmit={onSubmit} onCancel={vi.fn()} />);

    fireEvent.click(screen.getByRole('button', { name: 'device.saveDevice' }));

    await waitFor(() => {
      expect(onSubmit).toHaveBeenCalledWith(
        expect.objectContaining({
          connection_config: expect.objectContaining({
            data_format: 'CDAB',
          }),
        }),
      );
    });
  });
});
