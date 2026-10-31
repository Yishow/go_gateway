import { fireEvent, render, screen, waitFor } from '@testing-library/react';
import { MemoryRouter } from 'react-router-dom';
import { expect, it, vi } from 'vitest';
import PointForm from '../../../../src/components/datalink/PointForm';
import type { Point } from '../../../../src/types/datalink';

vi.mock('../../../../src/hooks/datalink/useDevices', () => ({
  useDevicesQuery: () => ({ data: [{ id: 'device', name: 'Device' }] }),
}));
vi.mock('../../../../src/hooks/datalink/usePollingGroups', () => ({
  usePollingGroupsQuery: () => ({ data: [{ id: 'group', name: 'Group' }] }),
}));

it('ordinary edit preserves existing override by omitting all format intent and guard fields', async () => {
  const point: Point = {
    id: 'point', device_id: 'device', name: 'Before', description: '', address: 'HR0',
    data_type: 'uint32', data_format: 'CDAB', format_mode: 'override', format_override: 'CDAB',
    format_version: 1, revision_id: 'captured-point', enabled: true, polling_group_id: 'group',
    last_value: null, last_read_at: '', last_error: '', error_count: 0, created_at: '', updated_at: '',
  };
  const submit = vi.fn().mockResolvedValue(undefined);
  const { container } = render(<MemoryRouter><PointForm point={point} onSubmit={submit} onCancel={vi.fn()} /></MemoryRouter>);
  fireEvent.change(screen.getByDisplayValue('Before'), { target: { value: 'After' } });
  fireEvent.submit(container.querySelector('form')!);
  await waitFor(() => expect(submit).toHaveBeenCalledTimes(1));
  expect(submit.mock.calls[0][0]).toMatchObject({ name: 'After' });
  for (const key of ['data_format', 'format_mode', 'configured_format', 'format_override', 'expected_workspace_revision', 'expected_point_revision', 'request_id']) {
    expect(submit.mock.calls[0][0]).not.toHaveProperty(key);
  }
});
