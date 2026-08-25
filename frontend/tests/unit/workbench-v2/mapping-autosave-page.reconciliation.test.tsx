import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { fireEvent, render, screen, waitFor } from '@testing-library/react';
import { describe, expect, it, vi } from 'vitest';
import DatalinkWorkbenchV2Page from '../../../src/pages/datalink/workbench-v2/DatalinkWorkbenchV2Page';
import { studioV2RulesAPI } from '../../../src/services/studioV2Rules';
import { studioV2MappingsAPI } from '../../../src/services/studioV2Mappings';
import {
  mappingFixture,
  mappingPoints,
  registerMappingAutosaveFixtures,
  shiftedMappingPoints,
  runtimeApplied,
  type WorkbenchV2ShellMockProps,
} from './mapping-autosave-page.testHarness';

function renderPage() {
  const queryClient = new QueryClient({ defaultOptions: { queries: { retry: false } } });
  render(
    <QueryClientProvider client={queryClient}>
      <DatalinkWorkbenchV2Page />
    </QueryClientProvider>,
  );
}

vi.mock('../../../src/features/datalink/workbench-v2/shell/WorkbenchV2Shell', () => ({
  WorkbenchV2Shell: ({
    state,
    actions,
  }: WorkbenchV2ShellMockProps<{
    mappings: Record<
      string,
      {
        tag_key: string;
        local_value?: { tag_key?: string } | null;
        persisted_value?: { tag_key?: string } | null;
        save_state: string;
        save_error?: string | null;
      }
    >;
  }>) => (
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
      {Object.entries(state.mappings).map(([pointId, mapping]) => (
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
      <button
        type="button"
        data-testid="disable-all-mappings"
        onClick={() => actions.dispatch({
          type: 'setAllMappingsEnabled',
          enabled: false,
        })}
      />
      <button
        type="button"
        data-testid="enable-all-mappings"
        onClick={() => actions.dispatch({
          type: 'setAllMappingsEnabled',
          enabled: true,
        })}
      />
    </div>
  ),
}));

vi.mock('../../../src/services/studioV2Workspace', async () => {
  const { studioV2ServiceMocks } = await import('./mapping-autosave-page.testHarness');
  return { studioV2WorkspaceAPI: studioV2ServiceMocks.workspace };
});

vi.mock('../../../src/services/studioV2WorkspaceDevices', async () => {
  const { studioV2ServiceMocks } = await import('./mapping-autosave-page.testHarness');
  return { studioV2WorkspaceDevicesAPI: studioV2ServiceMocks.devices };
});

vi.mock('../../../src/services/studioV2Rules', async () => {
  const { studioV2ServiceMocks } = await import('./mapping-autosave-page.testHarness');
  return { studioV2RulesAPI: studioV2ServiceMocks.rules };
});

vi.mock('../../../src/services/studioV2Mappings', async () => {
  const { studioV2ServiceMocks } = await import('./mapping-autosave-page.testHarness');
  return { studioV2MappingsAPI: studioV2ServiceMocks.mappings };
});

vi.mock('../../../src/services/studioV2WorkspaceDatabase', async () => {
  const { studioV2ServiceMocks } = await import('./mapping-autosave-page.testHarness');
  return { studioV2WorkspaceDatabaseAPI: studioV2ServiceMocks.database };
});

registerMappingAutosaveFixtures();

describe('DatalinkWorkbenchV2Page mapping autosave orchestration', () => {
  it('isolates one mapping row save failure without blocking another valid row', async () => {
    vi.mocked(studioV2MappingsAPI.list).mockResolvedValue([
      mappingFixture(),
      mappingFixture({
        id: 'mapping-B',
        point_id: 'point-B',
        rule_id: 'rule-B',
        device_id: 'dev-B',
        address: '40011',
        tag_id: 'tag-B',
        tag_key: 'line.b.persisted',
        display_name: 'Line B Temp',
      }),
    ]);
    vi.mocked(studioV2MappingsAPI.update).mockImplementation(async (mappingId, payload) => {
      if (mappingId === 'mapping-A') {
        throw new Error('save failed');
      }

      return runtimeApplied(mappingFixture({
        id: mappingId,
        point_id: 'point-B',
        rule_id: 'rule-B',
        device_id: 'dev-B',
        address: '40011',
        tag_id: 'tag-B',
        tag_key: String(payload.tag_key),
        display_name: 'Line B Temp',
      }));
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
      expect(screen.getByTestId('mapping-save-error-rule-A-p-0')).toHaveTextContent('errors.autosave_failed');
      expect(screen.getByTestId('mapping-local-tag-key-rule-A-p-0')).toHaveTextContent('line.a.broken');
      expect(screen.getByTestId('mapping-persisted-tag-key-rule-A-p-0')).toHaveTextContent('line.a.persisted');
    });
  });

  it('recreates the row mapping instead of reusing a stale mapping_id after the rule address changes', async () => {
    vi.mocked(studioV2MappingsAPI.list).mockResolvedValue([mappingFixture()]);
    vi.mocked(studioV2MappingsAPI.create).mockResolvedValueOnce(
      runtimeApplied(mappingFixture({
        id: 'mapping-A-new',
        point_id: 'point-A-new',
        address: '40021',
        tag_id: 'tag-A-new',
        tag_key: 'line.a.saved',
      })),
    );

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
