import { fireEvent, render, screen, waitFor } from '@testing-library/react';
import { MemoryRouter } from 'react-router-dom';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import LocalModbusWorkbenchPage from '../LocalModbusWorkbenchPage';

const { mockModbusShareAPI } = vi.hoisted(() => ({
  mockModbusShareAPI: {
    status: vi.fn(),
    start: vi.fn(),
    stop: vi.fn(),
    listMappings: vi.fn(),
    upsertMapping: vi.fn(),
    deleteMapping: vi.fn(),
    writeTagValue: vi.fn(),
    sync: vi.fn(),
  },
}));

vi.mock('../../../services/datalink', () => ({
  modbusShareAPI: mockModbusShareAPI,
}));

vi.mock('../../../hooks/datalink/useTags', () => ({
  useTagsQuery: () => ({
    data: [
      { id: 'tag-1', key: 'TEMP_1' },
      { id: 'tag-2', key: 'PRESS_1' },
    ],
  }),
}));

describe('LocalModbusWorkbenchPage', () => {
  const renderPage = (entry = '/datalink/local-modbus') =>
    render(
      <MemoryRouter initialEntries={[entry]}>
        <LocalModbusWorkbenchPage />
      </MemoryRouter>
    );

  beforeEach(() => {
    vi.clearAllMocks();
    mockModbusShareAPI.status.mockResolvedValue({
      enabled: false,
      port: 5020,
      address: '',
      bind_state: 'fail',
      mapping_count: 0,
    });
    mockModbusShareAPI.listMappings.mockResolvedValue([]);
    mockModbusShareAPI.start.mockResolvedValue({
      enabled: true,
      port: 5030,
      address: '127.0.0.1:5030',
      bind_state: 'pass',
      mapping_count: 0,
    });
    mockModbusShareAPI.stop.mockResolvedValue({
      enabled: false,
      port: 0,
      address: '',
      bind_state: 'fail',
      mapping_count: 0,
    });
    mockModbusShareAPI.sync.mockResolvedValue({ updated: 0, skipped: 0, errors: [] });
  });

  it('starts local server with specified port', async () => {
    renderPage();

    await screen.findByText('Stopped');
    const input = await screen.findByPlaceholderText('Server Port');
    fireEvent.change(input, { target: { value: '5030' } });
    fireEvent.click(screen.getByRole('button', { name: 'Start Server' }));

    await waitFor(() => {
      expect(mockModbusShareAPI.start).toHaveBeenCalledWith(5030);
    });
  });

  it('stops local server when currently enabled', async () => {
    mockModbusShareAPI.status.mockResolvedValue({
      enabled: true,
      port: 5020,
      address: '127.0.0.1:5020',
      bind_state: 'pass',
      mapping_count: 0,
    });

    renderPage();

    const stopButton = await screen.findByRole('button', { name: 'Stop Server' });
    fireEvent.click(stopButton);

    await waitFor(() => {
      expect(mockModbusShareAPI.stop).toHaveBeenCalledTimes(1);
    });
  });

  it('blocks sync when register conflicts exist', async () => {
    mockModbusShareAPI.status.mockResolvedValue({
      enabled: true,
      port: 5020,
      address: '127.0.0.1:5020',
      bind_state: 'pass',
      mapping_count: 2,
    });
    mockModbusShareAPI.listMappings.mockResolvedValue([
      { tag_id: 'tag-1', register: 12, data_type: 'float32', updated_at: new Date().toISOString() },
      { tag_id: 'tag-2', register: 12, data_type: 'float32', updated_at: new Date().toISOString() },
    ]);

    renderPage();

    const syncButton = await screen.findByRole('button', { name: 'Sync from Mappings' });
    expect(syncButton).toBeDisabled();
    expect(mockModbusShareAPI.sync).not.toHaveBeenCalled();
  });

  it('provides dashboard return link and preserves section context', async () => {
    renderPage('/datalink/local-modbus?section=settings');
    const returnLink = await screen.findByRole('link', { name: '返回 Dashboard' });
    expect(returnLink).toHaveAttribute('href', '/datalink?section=settings');
    expect(screen.getByText('回跳區段: settings')).toBeInTheDocument();
  });
});
