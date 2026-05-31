import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { fireEvent, render, screen, waitFor } from '@testing-library/react';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import DatalinkWorkbenchV2Page from '../../../src/pages/datalink/workbench-v2/DatalinkWorkbenchV2Page';
import { studioV2WorkspaceAPI } from '../../../src/services/studioV2Workspace';
import { studioV2WorkspaceDevicesAPI } from '../../../src/services/studioV2WorkspaceDevices';
import { studioV2RulesAPI } from '../../../src/services/studioV2Rules';
import { studioV2MappingsAPI } from '../../../src/services/studioV2Mappings';
import { studioV2WorkspaceDatabaseAPI } from '../../../src/services/studioV2WorkspaceDatabase';

const mappingPoints = [
  {
    id: 'rule-A-p-0',
    device_id: 'dev-A',
    rule_id: 'rule-A',
    rule_name: 'Line A Registers',
    name: 'A_0',
    address: '40001',
    data_type: 'int16',
    function: 'holding_register',
    width: 1,
    enabled: true,
    skipped: false,
    _rule_scale: 1,
    _rule_offset: 0,
  },
  {
    id: 'rule-B-p-0',
    device_id: 'dev-B',
    rule_id: 'rule-B',
    rule_name: 'Line B Registers',
    name: 'B_0',
    address: '40011',
    data_type: 'int16',
    function: 'holding_register',
    width: 1,
    enabled: true,
    skipped: false,
    _rule_scale: 1,
    _rule_offset: 0,
  },
];

const shiftedMappingPoints = [
  {
    ...mappingPoints[0],
    address: '40021',
    name: 'A_20',
  },
  mappingPoints[1],
];

vi.mock('../../../src/features/datalink/workbench-v2/shell/WorkbenchV2Shell', () => ({
  WorkbenchV2Shell: ({ state, actions }: any) => (
    <div data-testid="mapping-autosave-shell">
      <button
        type="button"
        data-testid="prime-mappings"
        onClick={() => actions.dispatch({ type: 'initMappingsForPoints', points: mappingPoints })}
      />
      <button
        type="button"
        data-testid="shift-rule-A-points"
        onClick={() => actions.dispatch({ type: 'initMappingsForPoints', points: shiftedMappingPoints })}
      />
      {Object.entries(state.mappings).map(([pointId, mapping]: any) => (
        <div key={pointId}>
          <div data-testid={`mapping-tag-key-${pointId}`}>{mapping.tag_key}</div>
          <div data-testid={`mapping-local-tag-key-${pointId}`}>{mapping.local_value?.tag_key ?? ''}</div>
          <div data-testid={`mapping-persisted-tag-key-${pointId}`}>{mapping.persisted_value?.tag_key ?? ''}</div>
          <div data-testid={`mapping-save-state-${pointId}`}>{mapping.save_state}</div>
          <div data-testid={`mapping-save-error-${pointId}`}>{mapping.save_error ?? ''}</div>
        </div>
      ))}
      <button
        type="button"
        data-testid="make-rule-A-valid"
        onClick={() => actions.dispatch({
          type: 'updateMapping',
          pointId: 'rule-A-p-0',
          patch: { tag_key: 'line.a.saved' },
        })}
      />
      <button
        type="button"
        data-testid="make-rule-A-invalid"
        onClick={() => actions.dispatch({
          type: 'updateMapping',
          pointId: 'rule-A-p-0',
          patch: { tag_key: '' },
        })}
      />
      <button
        type="button"
        data-testid="make-rule-A-fail"
        onClick={() => actions.dispatch({
          type: 'updateMapping',
          pointId: 'rule-A-p-0',
          patch: { tag_key: 'line.a.broken' },
        })}
      />
      <button
        type="button"
        data-testid="make-rule-B-valid"
        onClick={() => actions.dispatch({
          type: 'updateMapping',
          pointId: 'rule-B-p-0',
          patch: { tag_key: 'line.b.saved' },
        })}
      />
      <button
        type="button"
        data-testid="make-rule-A-save-fail"
        onClick={() => actions.dispatch({
          type: 'updateRule',
          ruleId: 'rule-A',
          patch: { start_address: '40021' },
        })}
      />
    </div>
  ),
}));

vi.mock('../../../src/services/studioV2Workspace', () => ({
  studioV2WorkspaceAPI: {
    get: vi.fn(),
  },
}));

vi.mock('../../../src/services/studioV2WorkspaceDevices', () => ({
  studioV2WorkspaceDevicesAPI: {
    list: vi.fn(),
    create: vi.fn(),
    update: vi.fn(),
    remove: vi.fn(),
    updateOrder: vi.fn(),
  },
}));

vi.mock('../../../src/services/studioV2Rules', () => ({
  studioV2RulesAPI: {
    list: vi.fn(),
    create: vi.fn(),
    update: vi.fn(),
    remove: vi.fn(),
  },
}));

vi.mock('../../../src/services/studioV2Mappings', () => ({
  studioV2MappingsAPI: {
    list: vi.fn(),
    create: vi.fn(),
    update: vi.fn(),
    remove: vi.fn(),
  },
}));

vi.mock('../../../src/services/studioV2WorkspaceDatabase', () => ({
  studioV2WorkspaceDatabaseAPI: {
    getConfig: vi.fn(),
    updateConfig: vi.fn(),
    listTargets: vi.fn(),
    upsertTarget: vi.fn(),
  },
}));

function renderPage() {
  const queryClient = new QueryClient({
    defaultOptions: {
      queries: {
        retry: false,
      },
    },
  });

  render(
    <QueryClientProvider client={queryClient}>
      <DatalinkWorkbenchV2Page />
    </QueryClientProvider>,
  );
}

describe('DatalinkWorkbenchV2Page mapping autosave orchestration', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    vi.mocked(studioV2WorkspaceAPI.get).mockResolvedValue({
      id: 'workspace-1',
      kind: 'single',
      status: 'ready',
      ordered_device_ids: ['dev-B', 'dev-A'],
      created_at: '2026-05-30T00:00:00Z',
      updated_at: '2026-05-30T00:00:00Z',
    });
    vi.mocked(studioV2WorkspaceDevicesAPI.list).mockResolvedValue([
      {
        id: 'dev-B',
        name: 'Line B PLC',
        description: '',
        protocol: 'modbus_tcp',
        status: 'draft',
        connection_config: '{"host":"192.168.10.11","port":502,"slave_id":2,"timeout":5}',
        last_test_at: null,
        last_test_success: null,
        last_test_error: '',
        created_at: '2026-05-30T00:00:00Z',
        updated_at: '2026-05-30T00:00:00Z',
      },
      {
        id: 'dev-A',
        name: 'Line A PLC',
        description: '',
        protocol: 'modbus_tcp',
        status: 'draft',
        connection_config: '{"host":"192.168.10.10","port":502,"slave_id":1,"timeout":5}',
        last_test_at: null,
        last_test_success: null,
        last_test_error: '',
        created_at: '2026-05-30T00:00:00Z',
        updated_at: '2026-05-30T00:00:00Z',
      },
    ]);
    vi.mocked(studioV2RulesAPI.list).mockResolvedValue([
      {
        id: 'rule-B',
        device_id: 'dev-B',
        workspace_id: 'workspace-1',
        start_address: '40011',
        count: 1,
        data_type: 'int16',
        naming_prefix: 'B_',
        enabled: true,
        locked: false,
        origin: 'manual',
        skipped_addresses: [],
        revision_id: 'rev-B',
        scale_multiplier: 1,
        scale_offset: 0,
        data_format: '',
        created_at: '2026-05-30T00:00:00Z',
        updated_at: '2026-05-30T00:00:00Z',
      },
      {
        id: 'rule-A',
        device_id: 'dev-A',
        workspace_id: 'workspace-1',
        start_address: '40001',
        count: 1,
        data_type: 'int16',
        naming_prefix: 'A_',
        enabled: true,
        locked: false,
        origin: 'manual',
        skipped_addresses: [],
        revision_id: 'rev-A',
        scale_multiplier: 1,
        scale_offset: 0,
        data_format: '',
        created_at: '2026-05-30T00:00:00Z',
        updated_at: '2026-05-30T00:00:00Z',
      },
    ]);
    vi.mocked(studioV2MappingsAPI.list).mockResolvedValue([]);
    vi.mocked(studioV2WorkspaceDatabaseAPI.getConfig).mockResolvedValue(null);
    vi.mocked(studioV2WorkspaceDatabaseAPI.listTargets).mockResolvedValue([]);
  });

  it('hydrates persisted mapping rows onto the same rule/address rows on step init', async () => {
    vi.mocked(studioV2MappingsAPI.list).mockResolvedValue([
      {
        id: 'mapping-B',
        workspace_id: 'workspace-1',
        point_id: 'point-B',
        rule_id: 'rule-B',
        device_id: 'dev-B',
        address: '40011',
        tag_id: 'tag-B',
        tag_key: 'line.b.persisted',
        display_name: 'Line B Temp',
        unit: 'C',
        target_type: 'float64',
        scale: 1,
        offset: 0,
        enabled: true,
        save_state: 'saved',
        created_at: '2026-05-30T00:00:00Z',
        updated_at: '2026-05-30T00:00:00Z',
      },
      {
        id: 'mapping-A',
        workspace_id: 'workspace-1',
        point_id: 'point-A',
        rule_id: 'rule-A',
        device_id: 'dev-A',
        address: '40001',
        tag_id: 'tag-A',
        tag_key: 'line.a.persisted',
        display_name: 'Line A Temp',
        unit: 'C',
        target_type: 'float64',
        scale: 1,
        offset: 0,
        enabled: true,
        save_state: 'saved',
        created_at: '2026-05-30T00:00:00Z',
        updated_at: '2026-05-30T00:00:00Z',
      },
    ]);

    renderPage();

    await waitFor(() => {
      expect(screen.getByTestId('workbench-v2-root')).toBeInTheDocument();
    });

    fireEvent.click(screen.getByTestId('prime-mappings'));

    await waitFor(() => {
      expect(screen.getByTestId('mapping-tag-key-rule-B-p-0')).toHaveTextContent('line.b.persisted');
      expect(screen.getByTestId('mapping-tag-key-rule-A-p-0')).toHaveTextContent('line.a.persisted');
    });
  });

  it('saves one valid draft mapping row and marks it saved', async () => {
    vi.mocked(studioV2MappingsAPI.create).mockResolvedValueOnce({
      id: 'mapping-A',
      workspace_id: 'workspace-1',
      point_id: 'point-A',
      rule_id: 'rule-A',
      device_id: 'dev-A',
      address: '40001',
      tag_id: 'tag-A',
      tag_key: 'line.a.saved',
      display_name: '入口溫度',
      unit: '°C',
      target_type: 'float64',
      scale: 1,
      offset: 0,
      enabled: true,
      save_state: 'saved',
      runtime_apply_status: 'not_running',
      created_at: '2026-05-30T00:00:00Z',
      updated_at: '2026-05-30T00:00:00Z',
    });

    renderPage();

    await waitFor(() => {
      expect(screen.getByTestId('workbench-v2-root')).toBeInTheDocument();
    });

    fireEvent.click(screen.getByTestId('prime-mappings'));
    fireEvent.click(screen.getByTestId('make-rule-A-valid'));

    await waitFor(() => {
      expect(studioV2MappingsAPI.create).toHaveBeenCalledWith(
        expect.objectContaining({
          rule_id: 'rule-A',
          address: '40001',
          tag_key: 'line.a.saved',
        }),
      );
    });

    await waitFor(() => {
      expect(screen.getByTestId('mapping-save-state-rule-A-p-0')).toHaveTextContent('saved');
    });
  });

  it('keeps invalid local mapping edits without overwriting the backend', async () => {
    renderPage();

    await waitFor(() => {
      expect(screen.getByTestId('workbench-v2-root')).toBeInTheDocument();
    });

    fireEvent.click(screen.getByTestId('prime-mappings'));
    fireEvent.click(screen.getByTestId('make-rule-A-invalid'));

    await waitFor(() => {
      expect(screen.getByTestId('mapping-tag-key-rule-A-p-0')).toHaveTextContent('');
      expect(screen.getByTestId('mapping-save-state-rule-A-p-0')).toHaveTextContent('draft-invalid');
    });

    expect(studioV2MappingsAPI.create).not.toHaveBeenCalled();
    expect(studioV2MappingsAPI.update).not.toHaveBeenCalled();
  });

  it('isolates one mapping row save failure without blocking another valid row', async () => {
    vi.mocked(studioV2MappingsAPI.list).mockResolvedValue([
      {
        id: 'mapping-A',
        workspace_id: 'workspace-1',
        point_id: 'point-A',
        rule_id: 'rule-A',
        device_id: 'dev-A',
        address: '40001',
        tag_id: 'tag-A',
        tag_key: 'line.a.persisted',
        display_name: 'Line A Temp',
        unit: 'C',
        target_type: 'float64',
        scale: 1,
        offset: 0,
        enabled: true,
        save_state: 'saved',
        created_at: '2026-05-30T00:00:00Z',
        updated_at: '2026-05-30T00:00:00Z',
      },
      {
        id: 'mapping-B',
        workspace_id: 'workspace-1',
        point_id: 'point-B',
        rule_id: 'rule-B',
        device_id: 'dev-B',
        address: '40011',
        tag_id: 'tag-B',
        tag_key: 'line.b.persisted',
        display_name: 'Line B Temp',
        unit: 'C',
        target_type: 'float64',
        scale: 1,
        offset: 0,
        enabled: true,
        save_state: 'saved',
        created_at: '2026-05-30T00:00:00Z',
        updated_at: '2026-05-30T00:00:00Z',
      },
    ]);
    vi.mocked(studioV2MappingsAPI.update).mockImplementation(async (mappingId, payload) => {
      if (mappingId === 'mapping-A') {
        throw new Error('save failed');
      }

      return {
        id: mappingId,
        workspace_id: 'workspace-1',
        point_id: 'point-B',
        rule_id: 'rule-B',
        device_id: 'dev-B',
        address: '40011',
        tag_id: 'tag-B',
        tag_key: String(payload.tag_key),
        display_name: 'Line B Temp',
        unit: 'C',
        target_type: 'float64',
        scale: 1,
        offset: 0,
        enabled: true,
        save_state: 'saved',
        runtime_apply_status: 'not_running',
        created_at: '2026-05-30T00:00:00Z',
        updated_at: '2026-05-30T00:00:00Z',
      } as any;
    });

    renderPage();

    await waitFor(() => {
      expect(screen.getByTestId('workbench-v2-root')).toBeInTheDocument();
    });

    fireEvent.click(screen.getByTestId('prime-mappings'));

    await waitFor(() => {
      expect(screen.getByTestId('mapping-tag-key-rule-A-p-0')).toHaveTextContent('line.a.persisted');
      expect(screen.getByTestId('mapping-tag-key-rule-B-p-0')).toHaveTextContent('line.b.persisted');
    });

    fireEvent.click(screen.getByTestId('make-rule-A-fail'));
    fireEvent.click(screen.getByTestId('make-rule-B-valid'));

    await waitFor(() => {
      expect(studioV2MappingsAPI.update).toHaveBeenCalledWith(
        'mapping-A',
        expect.objectContaining({ tag_key: 'line.a.broken' }),
      );
      expect(studioV2MappingsAPI.update).toHaveBeenCalledWith(
        'mapping-B',
        expect.objectContaining({ tag_key: 'line.b.saved' }),
      );
    });

    await waitFor(() => {
      expect(screen.getByTestId('mapping-save-state-rule-A-p-0')).toHaveTextContent('save-error');
      expect(screen.getByTestId('mapping-save-state-rule-B-p-0')).toHaveTextContent('saved');
      expect(screen.getByTestId('mapping-save-error-rule-A-p-0')).toHaveTextContent('save failed');
      expect(screen.getByTestId('mapping-local-tag-key-rule-A-p-0')).toHaveTextContent('line.a.broken');
      expect(screen.getByTestId('mapping-persisted-tag-key-rule-A-p-0')).toHaveTextContent('line.a.persisted');
    });
  });

  it('recreates the row mapping instead of reusing a stale mapping_id after the rule address changes', async () => {
    vi.mocked(studioV2MappingsAPI.list).mockResolvedValue([
      {
        id: 'mapping-A',
        workspace_id: 'workspace-1',
        point_id: 'point-A',
        rule_id: 'rule-A',
        device_id: 'dev-A',
        address: '40001',
        tag_id: 'tag-A',
        tag_key: 'line.a.persisted',
        display_name: 'Line A Temp',
        unit: 'C',
        target_type: 'float64',
        scale: 1,
        offset: 0,
        enabled: true,
        save_state: 'saved',
        created_at: '2026-05-30T00:00:00Z',
        updated_at: '2026-05-30T00:00:00Z',
      },
    ]);
    vi.mocked(studioV2MappingsAPI.create).mockResolvedValueOnce({
      id: 'mapping-A-new',
      workspace_id: 'workspace-1',
      point_id: 'point-A-new',
      rule_id: 'rule-A',
      device_id: 'dev-A',
      address: '40021',
      tag_id: 'tag-A-new',
      tag_key: 'line.a.saved',
      display_name: 'Line A Temp',
      unit: 'C',
      target_type: 'float64',
      scale: 1,
      offset: 0,
      enabled: true,
      save_state: 'saved',
      runtime_apply_status: 'not_running',
      created_at: '2026-05-30T00:00:00Z',
      updated_at: '2026-05-30T00:00:00Z',
    });

    renderPage();

    await waitFor(() => {
      expect(screen.getByTestId('workbench-v2-root')).toBeInTheDocument();
    });

    fireEvent.click(screen.getByTestId('prime-mappings'));

    await waitFor(() => {
      expect(screen.getByTestId('mapping-tag-key-rule-A-p-0')).toHaveTextContent('line.a.persisted');
    });

    fireEvent.click(screen.getByTestId('shift-rule-A-points'));
    fireEvent.click(screen.getByTestId('make-rule-A-valid'));

    await waitFor(() => {
      expect(studioV2MappingsAPI.create).toHaveBeenCalledWith(
        expect.objectContaining({
          rule_id: 'rule-A',
          address: '40021',
          tag_key: 'line.a.saved',
        }),
      );
    });

    expect(studioV2MappingsAPI.update).not.toHaveBeenCalledWith(
      'mapping-A',
      expect.anything(),
    );
  });

  it('does not autosave mappings while the owning rule is in save-error', async () => {
    vi.mocked(studioV2RulesAPI.update).mockRejectedValueOnce(new Error('rule save failed'));

    renderPage();

    await waitFor(() => {
      expect(screen.getByTestId('workbench-v2-root')).toBeInTheDocument();
    });

    fireEvent.click(screen.getByTestId('make-rule-A-save-fail'));

    await waitFor(() => {
      expect(studioV2RulesAPI.update).toHaveBeenCalledWith(
        'rule-A',
        expect.objectContaining({ start_address: '40021' }),
      );
    });

    fireEvent.click(screen.getByTestId('shift-rule-A-points'));
    fireEvent.click(screen.getByTestId('make-rule-A-valid'));

    await waitFor(() => {
      expect(studioV2MappingsAPI.create).not.toHaveBeenCalled();
      expect(studioV2MappingsAPI.update).not.toHaveBeenCalled();
    });
  });
});
