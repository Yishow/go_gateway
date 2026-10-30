import { act, fireEvent, screen, waitFor } from '@testing-library/react';
import { describe, expect, it, beforeEach, vi } from 'vitest';
import { studioV2WorkspaceKeys } from '../../../src/hooks/datalink/keys';
import {
  candidateSnapshot,
  compatibleAddresses,
  mappingFixture,
  modbusShareAPI,
  renderPage,
  setupPersistedMCMocks,
  sourceRuleCandidateAPI,
  studioV2ServiceMocks,
  workspaceFixture,
} from './mc-persisted-hydration.fixtures';

vi.mock('../../../src/features/datalink/workbench-v2/shell/WorkbenchV2Shell', () => ({
  WorkbenchV2Shell: ({
    state,
    actions,
  }: {
    state: {
      rules: Array<{ id: string; mc_persisted_points?: { status: string } }>;
      mappings: Record<string, { tag_key: string; save_state?: string }>;
    };
    actions: { dispatch: (action: unknown) => void };
  }) => {
    const rule = state.rules.find((item) => item.id === 'rule-mc');
    const mapping = state.mappings['legacy-point-0'];
    return (
      <div data-testid="mc-mapping-queue-shell">
        <div data-testid="mc-mapping-status">{rule?.mc_persisted_points?.status ?? 'missing'}</div>
        <div data-testid="mc-mapping-value">{mapping?.tag_key ?? 'missing'}:{mapping?.save_state ?? 'missing'}</div>
        <button type="button" data-testid="edit-mc-mapping" onClick={() => actions.dispatch({ type: 'updateMapping', pointId: 'legacy-point-0', patch: { tag_key: 'edited.tag' } })} />
      </div>
    );
  },
}));

vi.mock('../../../src/services/studioV2Workspace', async () => {
  const { studioV2ServiceMocks: mocks } = await import('./helpers/studioV2ServiceMocks');
  return { studioV2WorkspaceAPI: mocks.workspace };
});
vi.mock('../../../src/services/studioV2WorkspaceDevices', async () => {
  const { studioV2ServiceMocks: mocks } = await import('./helpers/studioV2ServiceMocks');
  return { studioV2WorkspaceDevicesAPI: mocks.devices };
});
vi.mock('../../../src/services/studioV2Rules', async () => {
  const { studioV2ServiceMocks: mocks } = await import('./helpers/studioV2ServiceMocks');
  return { studioV2RulesAPI: mocks.rules };
});
vi.mock('../../../src/services/studioV2Mappings', async () => {
  const { studioV2ServiceMocks: mocks } = await import('./helpers/studioV2ServiceMocks');
  return { studioV2MappingsAPI: mocks.mappings };
});
vi.mock('../../../src/services/studioV2WorkspaceDatabase', async () => {
  const { studioV2ServiceMocks: mocks } = await import('./helpers/studioV2ServiceMocks');
  return { studioV2WorkspaceDatabaseAPI: mocks.database };
});
vi.mock('../../../src/services/sourceRuleCandidates', () => ({
  sourceRuleCandidateAPI: { get: vi.fn() },
}));
vi.mock('../../../src/services/datalink', () => ({
  modbusShareAPI: { status: vi.fn(), reconcile: vi.fn() },
}));

describe('persisted MC mapping queue recovery', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    setupPersistedMCMocks();
  });

  it('resends the latest mapping draft after an in-flight response crosses scope', async () => {
    let workspaceRevision = 'workspace-rev-1';
    let resolveOldMapping: ((value: ReturnType<typeof mappingFixture>) => void) | undefined;
    let resolveNewCandidate: ((value: ReturnType<typeof candidateSnapshot>) => void) | undefined;
    vi.mocked(studioV2ServiceMocks.workspace.get).mockImplementation(async () => workspaceFixture({
      ordered_device_ids: ['dev-mc', 'dev-modbus'],
      modbus_share: { hydration_state: 'ready', workspace_id: 'workspace-1', workspace_revision: workspaceRevision, readiness: true,
        status: { workspace_id: 'workspace-1', enabled: false, port: 5020, address: '127.0.0.1', bind_state: 'disabled', mapping_count: 0, workspace_revision: workspaceRevision } },
    }));
    vi.mocked(modbusShareAPI.status).mockImplementation(async () => ({
      workspace_id: 'workspace-1', enabled: false, port: 5020, address: '127.0.0.1', bind_state: 'disabled', mapping_count: 0, workspace_revision: workspaceRevision,
    }));
    vi.mocked(sourceRuleCandidateAPI.get).mockImplementation(async (_ruleID, scope) => {
      if (scope?.expected_workspace_revision === 'workspace-rev-2') {
        return new Promise((resolve) => { resolveNewCandidate = resolve; });
      }
      return candidateSnapshot(compatibleAddresses, { workspace_revision: scope?.expected_workspace_revision, revision_id: scope?.revision_id });
    });
    vi.mocked(studioV2ServiceMocks.mappings.update).mockImplementation(async (_mappingID, request) => {
      if (!resolveOldMapping) {
        return new Promise<ReturnType<typeof mappingFixture>>((resolve) => { resolveOldMapping = resolve; });
      }
      return mappingFixture({ id: 'legacy-mapping-0', point_id: 'legacy-point-0', rule_id: 'rule-mc', device_id: 'dev-mc', address: 'W0', tag_key: request.tag_key });
    });

    const queryClient = renderPage();
    await waitFor(() => expect(screen.getByTestId('mc-mapping-status')).toHaveTextContent('compatible'));
    fireEvent.click(screen.getByTestId('edit-mc-mapping'));
    await waitFor(() => expect(studioV2ServiceMocks.mappings.update).toHaveBeenCalledTimes(1));

    workspaceRevision = 'workspace-rev-2';
    await act(async () => {
      await queryClient.invalidateQueries({ queryKey: studioV2WorkspaceKeys.bootstrap() });
      await queryClient.invalidateQueries({ queryKey: studioV2WorkspaceKeys.sourceRules() });
    });
    await waitFor(() => expect(screen.getByTestId('mc-mapping-status')).toHaveTextContent('checking'));
    resolveOldMapping?.(mappingFixture({ id: 'legacy-mapping-0', point_id: 'legacy-point-0', rule_id: 'rule-mc', device_id: 'dev-mc', address: 'W0', tag_key: 'server.old' }));
    await waitFor(() => expect(screen.getByTestId('mc-mapping-value')).toHaveTextContent('edited.tag:draft-invalid'));

    resolveNewCandidate?.(candidateSnapshot(compatibleAddresses, { workspace_revision: 'workspace-rev-2' }));
    await waitFor(() => expect(screen.getByTestId('mc-mapping-status')).toHaveTextContent('compatible'));
    await waitFor(() => expect(studioV2ServiceMocks.mappings.update).toHaveBeenCalledTimes(2));
    expect(studioV2ServiceMocks.mappings.update).toHaveBeenLastCalledWith('legacy-mapping-0', expect.objectContaining({ tag_key: 'edited.tag' }));
  });
});
