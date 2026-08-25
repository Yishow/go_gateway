import { beforeEach, describe, expect, it, vi } from 'vitest';
import { modbusShareAPI } from '../../../src/services/modbusShare';
import { api } from '../../../src/services/datalinkClient';

vi.mock('../../../src/services/datalinkClient', () => ({
  api: {
    post: vi.fn(),
    get: vi.fn(),
  },
}));

const status = {
  enabled: true,
  running: true,
  port: 5020,
  address: '127.0.0.1',
  bind_state: 'pass' as const,
  mapping_count: 0,
};

describe('modbusShareAPI lifecycle requests', () => {
  beforeEach(() => vi.clearAllMocks());

  it('sends the expected settings revision when starting', async () => {
    vi.mocked(api.post).mockResolvedValue({ data: { success: true, data: status } });

    await modbusShareAPI.start({ port: 5020, expected_settings_revision: 'settings-7' });

    expect(api.post).toHaveBeenCalledWith('/modbus-share/start', {
      port: 5020,
      expected_settings_revision: 'settings-7',
    });
  });

  it('sends the expected settings revision when stopping', async () => {
    vi.mocked(api.post).mockResolvedValue({ data: { success: true, data: { ...status, running: false } } });

    await modbusShareAPI.stop('settings-8');

    expect(api.post).toHaveBeenCalledWith('/modbus-share/stop', {
      expected_settings_revision: 'settings-8',
    });
  });
});
