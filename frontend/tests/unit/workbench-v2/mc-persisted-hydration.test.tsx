import { act, fireEvent, screen, waitFor } from '@testing-library/react';
import { describe, expect, it, beforeEach, vi } from 'vitest';
import { studioV2WorkspaceKeys } from '../../../src/hooks/datalink/keys';
import {
  compatibleAddresses,
  candidateSnapshot,
  legacyAddresses,
  mappingFixture,
  ruleFixture,
  renderPage,
  setupPersistedMCMocks,
  sqliteConfigFixture,
  sourceRuleCandidateAPI,
  studioV2ServiceMocks,
  modbusShareAPI,
  workspaceFixture,
} from './mc-persisted-hydration.fixtures';
import type { SourceRuleCandidateSnapshotView } from '../../../src/types/sourceRuleCandidates';

vi.mock('../../../src/features/datalink/workbench-v2/shell/WorkbenchV2Shell', () => ({
  WorkbenchV2Shell: ({
    state,
    actions,
  }: {
    actions: { dispatch: (action: unknown) => void };
    state: {
      rules: Array<{ id: string; save_state?: string; mc_persisted_points?: { status: string; scope?: { workspace_id: string; workspace_revision: string; revision_id: string }; points: Array<{ point_id: string; address: string }> } }>;
      points: Array<{ id: string; rule_id: string; address: string }>;
      mappings: Record<string, { mapping_id?: string; persisted_point_id?: string; address?: string; tag_key?: string; save_state?: string }>;
      db: { connector: { table: string; save_state?: string }; targets: Record<string, { column_name: string; save_state?: string }> };
    };
  }) => {
    const rule = state.rules.find((item) => item.id === 'rule-mc');
    const points = state.points.filter((point) => point.rule_id === 'rule-mc');
    const mappings = Object.values(state.mappings).filter((mapping) => mapping.address?.startsWith('W'));
    return (
      <div data-testid="mc-persisted-shell">
        <div data-testid="mc-persisted-status">{rule?.mc_persisted_points?.status ?? 'missing'}</div>
        <div data-testid="mc-persisted-scope">{rule?.mc_persisted_points?.scope ? `${rule.mc_persisted_points.scope.workspace_id}@${rule.mc_persisted_points.scope.workspace_revision}@${rule.mc_persisted_points.scope.revision_id}` : 'missing'}</div>
        <div data-testid="mc-rule-save-state">{rule?.save_state ?? 'missing'}</div>
        <div data-testid="mc-persisted-points">{points.map((point) => `${point.id}:${point.address}`).join('|')}</div>
        <div data-testid="mc-persisted-mappings">{mappings.map((mapping) => `${mapping.mapping_id}:${mapping.persisted_point_id}:${mapping.address}:${mapping.tag_key}:${mapping.save_state}`).join('|')}</div>
        <div data-testid="mc-db-connector">{state.db.connector.table}:{state.db.connector.save_state}</div>
        <div data-testid="mc-db-target">{state.db.targets['legacy-point-0']?.column_name ?? 'missing'}:{state.db.targets['legacy-point-0']?.save_state ?? 'missing'}</div>
        <div data-testid="unaffected-points">{state.points.filter((point) => point.rule_id === 'rule-modbus').map((point) => `${point.id}:${point.address}`).join('|')}</div>
        <button type="button" data-testid="edit-mc-rule" onClick={() => actions.dispatch({ type: 'updateRule', ruleId: 'rule-mc', patch: { naming_prefix: 'EDITED_' } })} />
        <button type="button" data-testid="edit-mc-scale" onClick={() => actions.dispatch({ type: 'updateRule', ruleId: 'rule-mc', patch: { scale_multiplier: 2 } })} />
        <button type="button" data-testid="edit-mc-range" onClick={() => actions.dispatch({ type: 'updateRule', ruleId: 'rule-mc', patch: { count: 16 } })} />
        <button type="button" data-testid="edit-mc-mapping" onClick={() => actions.dispatch({ type: 'updateMapping', pointId: 'legacy-point-0', patch: { tag_key: 'edited.tag' } })} />
        <button type="button" data-testid="edit-mc-db-connector" onClick={() => actions.dispatch({ type: 'updateDbConnector', patch: { table: 'edited_table' } })} />
        <button type="button" data-testid="edit-mc-db-target" onClick={() => actions.dispatch({ type: 'updateDbTarget', pointId: 'legacy-point-0', patch: { column_name: 'edited_column' } })} />
        <button type="button" data-testid="remove-mc-points" onClick={() => actions.dispatch({ type: 'initMappingsForPoints', points: [] })} />
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

describe('persisted MC point hydration', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    setupPersistedMCMocks();
  });
  it('hydrates persisted MC point and mapping identities from the scoped candidate GET', async () => {
    renderPage();
    await waitFor(() => expect(screen.getByTestId('mc-persisted-shell')).toBeInTheDocument());
    await waitFor(() => expect(sourceRuleCandidateAPI.get).toHaveBeenCalledWith('rule-mc', {
      workspace_id: 'workspace-1',
      expected_workspace_revision: 'workspace-rev-1',
      revision_id: 'rule-rev-1',
    }));
    expect(screen.getByTestId('mc-persisted-points')).toHaveTextContent(
      legacyAddresses.map((address, index) => `legacy-point-${index}:${address}`).join('|'),
    );
    expect(screen.getByTestId('mc-persisted-mappings')).toHaveTextContent(
      legacyAddresses.map((address, index) => `legacy-mapping-${index}:legacy-point-${index}:${address}:legacy.w${index}:saved`).join('|'),
    );
    expect(screen.getByTestId('unaffected-points')).toHaveTextContent('rule-modbus-p-0:40001');
  });
  it('blocks rule, mapping, and dependent database mutations while persisted MC evidence conflicts', async () => {
    renderPage();
    await waitFor(() => expect(screen.getByTestId('mc-persisted-status')).toHaveTextContent('conflict'));
    fireEvent.click(screen.getByTestId('edit-mc-rule'));
    fireEvent.click(screen.getByTestId('edit-mc-mapping'));
    fireEvent.click(screen.getByTestId('edit-mc-db-connector'));
    fireEvent.click(screen.getByTestId('edit-mc-db-target'));
    fireEvent.click(screen.getByTestId('remove-mc-points'));
    await waitFor(() => {
      expect(studioV2ServiceMocks.rules.update).not.toHaveBeenCalled();
      expect(studioV2ServiceMocks.mappings.update).not.toHaveBeenCalled();
      expect(studioV2ServiceMocks.mappings.create).not.toHaveBeenCalled();
      expect(studioV2ServiceMocks.database.updateConfig).not.toHaveBeenCalled();
      expect(studioV2ServiceMocks.database.upsertTarget).not.toHaveBeenCalled();
      expect(studioV2ServiceMocks.mappings.remove).not.toHaveBeenCalled();
    });
  });
  it.each([
    ['forbidden', () => vi.mocked(sourceRuleCandidateAPI.get).mockRejectedValueOnce({ response: { status: 403 } })],
    ['revision conflict', () => vi.mocked(sourceRuleCandidateAPI.get).mockRejectedValueOnce({ response: { status: 409 } })],
    ['service unavailable', () => vi.mocked(sourceRuleCandidateAPI.get).mockRejectedValueOnce({ response: { status: 503 } })],
    ['wrong revision', () => vi.mocked(sourceRuleCandidateAPI.get).mockResolvedValueOnce(candidateSnapshot(legacyAddresses, { revision_id: 'rule-rev-old' }))],
    ['incomplete', () => vi.mocked(sourceRuleCandidateAPI.get).mockResolvedValueOnce(candidateSnapshot(legacyAddresses.slice(0, -1)))],
    ['missing points', () => vi.mocked(sourceRuleCandidateAPI.get).mockResolvedValueOnce(candidateSnapshot([]))],
  ])('keeps persisted MC blocked for %s candidate evidence', async (_caseName, arrange) => {
    arrange();
    renderPage();
    await waitFor(() => expect(screen.getByTestId('mc-persisted-status')).toHaveTextContent('unavailable'));
    expect(screen.getByTestId('mc-persisted-points')).toHaveTextContent('');
    fireEvent.click(screen.getByTestId('edit-mc-rule'));
    fireEvent.click(screen.getByTestId('edit-mc-mapping'));
    fireEvent.click(screen.getByTestId('edit-mc-db-target'));
    await waitFor(() => {
      expect(studioV2ServiceMocks.rules.update).not.toHaveBeenCalled();
      expect(studioV2ServiceMocks.mappings.update).not.toHaveBeenCalled();
      expect(studioV2ServiceMocks.database.upsertTarget).not.toHaveBeenCalled();
    });
  });
  it('blocks mutations while scoped evidence is pending and retains unaffected points', async () => {
    vi.mocked(sourceRuleCandidateAPI.get).mockImplementationOnce(() => new Promise(() => undefined));
    renderPage();
    await waitFor(() => expect(screen.getByTestId('mc-persisted-status')).toHaveTextContent('checking'));
    expect(screen.getByTestId('mc-persisted-points')).toHaveTextContent('');
    expect(screen.getByTestId('unaffected-points')).toHaveTextContent('rule-modbus-p-0:40001');
    fireEvent.click(screen.getByTestId('edit-mc-rule'));
    fireEvent.click(screen.getByTestId('edit-mc-db-connector'));
    await waitFor(() => {
      expect(studioV2ServiceMocks.rules.update).not.toHaveBeenCalled();
      expect(studioV2ServiceMocks.database.updateConfig).not.toHaveBeenCalled();
    });
  });
  it('treats a persisted MC server row without revision as unavailable', async () => {
    vi.mocked(studioV2ServiceMocks.rules.list).mockResolvedValueOnce([
      ruleFixture({ id: 'rule-mc', device_id: 'dev-mc', revision_id: undefined, start_address: 'W0', count: 17 }),
    ]);
    renderPage();
    await waitFor(() => expect(screen.getByTestId('mc-persisted-status')).toHaveTextContent('unavailable'));
    fireEvent.click(screen.getByTestId('edit-mc-rule'));
    await waitFor(() => expect(studioV2ServiceMocks.rules.update).not.toHaveBeenCalled());
  });
  it('keeps a locally hydrated MC rule blocked if the server list drops its row', async () => {
    vi.mocked(sourceRuleCandidateAPI.get).mockResolvedValueOnce(candidateSnapshot(compatibleAddresses));
    let ruleListCall = 0;
    vi.mocked(studioV2ServiceMocks.rules.list).mockImplementation(async () => {
      ruleListCall += 1;
      return ruleListCall === 1
        ? [
          ruleFixture({ id: 'rule-mc', device_id: 'dev-mc', revision_id: 'rule-rev-1', start_address: 'W0', count: 17 }),
          ruleFixture({ id: 'rule-modbus', device_id: 'dev-modbus', revision_id: 'rule-modbus-rev-1', start_address: '40001', count: 1 }),
        ]
        : [ruleFixture({ id: 'rule-modbus', device_id: 'dev-modbus', revision_id: 'rule-modbus-rev-1', start_address: '40001', count: 1 })];
    });
    const queryClient = renderPage();
    await waitFor(() => expect(screen.getByTestId('mc-persisted-status')).toHaveTextContent('compatible'));
    await act(async () => {
      await queryClient.invalidateQueries({ queryKey: studioV2WorkspaceKeys.sourceRules() });
    });
    await waitFor(() => expect(studioV2ServiceMocks.rules.list).toHaveBeenCalledTimes(2));
    await waitFor(() => expect(screen.getByTestId('mc-persisted-status')).toHaveTextContent('unavailable'));
    fireEvent.click(screen.getByTestId('edit-mc-rule'));
    await waitFor(() => expect(studioV2ServiceMocks.rules.update).not.toHaveBeenCalled());
  });
  it('allows an ordinary update when the stored MC address set is compatible', async () => {
    let revision = 1;
    let workspaceRevision = 'workspace-rev-1';
    let currentRule = ruleFixture({ id: 'rule-mc', device_id: 'dev-mc', start_address: 'W0', count: 17, revision_id: 'rule-rev-1' });
    let resolveFirstUpdate: ((value: ReturnType<typeof ruleFixture>) => void) | undefined;
    let holdFirstUpdate = true;
    vi.mocked(studioV2ServiceMocks.workspace.get).mockImplementation(async () => workspaceFixture({
      ordered_device_ids: ['dev-mc', 'dev-modbus'],
      modbus_share: { hydration_state: 'ready', workspace_id: 'workspace-1', workspace_revision: workspaceRevision, readiness: true,
        status: { workspace_id: 'workspace-1', enabled: false, port: 5020, address: '127.0.0.1', bind_state: 'disabled', mapping_count: 0, workspace_revision: workspaceRevision } },
    }));
    vi.mocked(modbusShareAPI.status).mockImplementation(async () => ({
      workspace_id: 'workspace-1', enabled: false, port: 5020, address: '127.0.0.1', bind_state: 'disabled', mapping_count: 0, workspace_revision: workspaceRevision,
    }));
    vi.mocked(studioV2ServiceMocks.rules.list).mockImplementation(async () => [currentRule, ruleFixture({ id: 'rule-modbus', device_id: 'dev-modbus', start_address: '40001', count: 1, revision_id: 'rule-modbus-rev-1' })]);
    vi.mocked(sourceRuleCandidateAPI.get).mockImplementation(async (_ruleID, scope) => candidateSnapshot(currentRule.count === 16 ? compatibleAddresses.slice(0, 16) : compatibleAddresses, { workspace_revision: scope?.expected_workspace_revision, revision_id: scope?.revision_id ?? currentRule.revision_id }));
    vi.mocked(studioV2ServiceMocks.rules.update).mockImplementation(async (_ruleID, request) => {
      revision += 1;
      currentRule = ruleFixture({ ...currentRule, ...request, revision_id: `rule-rev-${revision}` });
      workspaceRevision = 'workspace-rev-2';
      if (holdFirstUpdate) {
        holdFirstUpdate = false;
        return new Promise<ReturnType<typeof ruleFixture>>((resolve) => { resolveFirstUpdate = resolve; });
      }
      return currentRule;
    });
    const queryClient = renderPage();
    await waitFor(() => expect(screen.getByTestId('mc-persisted-status')).toHaveTextContent('compatible'));
    fireEvent.click(screen.getByTestId('edit-mc-rule'));
    await waitFor(() => expect(studioV2ServiceMocks.rules.update).toHaveBeenCalledTimes(1));
    await act(async () => {
      await queryClient.invalidateQueries({ queryKey: studioV2WorkspaceKeys.bootstrap() });
      await queryClient.invalidateQueries({ queryKey: studioV2WorkspaceKeys.sourceRules() });
    });
    await waitFor(() => expect(sourceRuleCandidateAPI.get).toHaveBeenCalledWith('rule-mc', expect.objectContaining({ expected_workspace_revision: 'workspace-rev-2', revision_id: 'rule-rev-2' })));
    await waitFor(() => expect(screen.getByTestId('mc-persisted-status')).toHaveTextContent('compatible'));
    resolveFirstUpdate?.(currentRule);
    await waitFor(() => expect(screen.getByTestId('mc-rule-save-state')).toHaveTextContent('saved'));
    await waitFor(() => expect(screen.getByTestId('mc-persisted-status')).toHaveTextContent('compatible'));
    fireEvent.click(screen.getByTestId('edit-mc-scale'));
    await waitFor(() => expect(studioV2ServiceMocks.rules.update).toHaveBeenCalledTimes(2));
    await waitFor(() => expect(screen.getByTestId('mc-persisted-status')).toHaveTextContent('compatible'));
    fireEvent.click(screen.getByTestId('edit-mc-range'));
    await waitFor(() => expect(studioV2ServiceMocks.rules.update).toHaveBeenCalledTimes(3));
    await waitFor(() => expect(screen.getByTestId('mc-persisted-status')).toHaveTextContent('compatible'));
  });
  it('does not let a late old-scope response restore writable persisted MC state', async () => {
    let resolveOld: ((snapshot: SourceRuleCandidateSnapshotView) => void) | undefined;
    let resolveNew: ((snapshot: SourceRuleCandidateSnapshotView) => void) | undefined;
    vi.mocked(sourceRuleCandidateAPI.get)
      .mockImplementationOnce(() => new Promise((resolve) => { resolveOld = resolve; }))
      .mockImplementationOnce(() => new Promise((resolve) => { resolveNew = resolve; }));
    let workspaceRead = 0;
    vi.mocked(studioV2ServiceMocks.workspace.get).mockImplementation(async () => {
      workspaceRead += 1;
      return workspaceFixture({
        ordered_device_ids: ['dev-mc', 'dev-modbus'],
        modbus_share: {
          hydration_state: 'ready',
          workspace_id: 'workspace-1',
          workspace_revision: workspaceRead > 1 ? 'workspace-rev-2' : 'workspace-rev-1',
          readiness: true,
          status: {
            workspace_id: 'workspace-1',
            enabled: false,
            port: 5020,
            address: '127.0.0.1',
            bind_state: 'disabled',
            mapping_count: 0,
            workspace_revision: workspaceRead > 1 ? 'workspace-rev-2' : 'workspace-rev-1',
          },
        },
      });
    });
    vi.mocked(modbusShareAPI.status).mockImplementation(async () => ({
      workspace_id: 'workspace-1', enabled: false, port: 5020, address: '127.0.0.1', bind_state: 'disabled', mapping_count: 0,
      workspace_revision: workspaceRead > 1 ? 'workspace-rev-2' : 'workspace-rev-1',
    }));
    const queryClient = renderPage();
    await waitFor(() => expect(sourceRuleCandidateAPI.get).toHaveBeenCalledWith('rule-mc', expect.objectContaining({ workspace_id: 'workspace-1', expected_workspace_revision: 'workspace-rev-1' })));
    await act(async () => {
      await queryClient.invalidateQueries({ queryKey: studioV2WorkspaceKeys.bootstrap() });
    });
    await waitFor(() => expect(sourceRuleCandidateAPI.get).toHaveBeenCalledWith('rule-mc', expect.objectContaining({ expected_workspace_revision: 'workspace-rev-2' })));
    expect(screen.getByTestId('mc-persisted-status')).toHaveTextContent('checking');
    resolveOld?.(candidateSnapshot());
    await waitFor(() => expect(screen.getByTestId('mc-persisted-status')).toHaveTextContent('checking'));
    expect(screen.getByTestId('mc-persisted-points')).toHaveTextContent('');
    fireEvent.click(screen.getByTestId('edit-mc-rule'));
    await waitFor(() => expect(studioV2ServiceMocks.rules.update).not.toHaveBeenCalled());
    resolveNew?.(candidateSnapshot(legacyAddresses, { workspace_revision: 'workspace-rev-2' }));
  });
  it('re-hydrates when the cached candidate scope changes even if revision and identities match', async () => {
    let workspaceID = 'workspace-1';
    vi.mocked(studioV2ServiceMocks.workspace.get).mockImplementation(async () => workspaceFixture({
      id: workspaceID,
      ordered_device_ids: ['dev-mc', 'dev-modbus'],
      modbus_share: {
        hydration_state: 'ready', workspace_id: workspaceID, workspace_revision: 'workspace-rev-1', readiness: true,
        status: { workspace_id: workspaceID, enabled: false, port: 5020, address: '127.0.0.1', bind_state: 'disabled', mapping_count: 0, workspace_revision: 'workspace-rev-1' },
      },
    }));
    vi.mocked(modbusShareAPI.status).mockImplementation(async () => ({
      workspace_id: workspaceID, enabled: false, port: 5020, address: '127.0.0.1', bind_state: 'disabled', mapping_count: 0, workspace_revision: 'workspace-rev-1',
    }));
    vi.mocked(sourceRuleCandidateAPI.get).mockImplementation(async (_ruleID, scope) => candidateSnapshot(compatibleAddresses, {
      workspace_id: scope?.workspace_id,
      workspace_revision: scope?.expected_workspace_revision,
      revision_id: scope?.revision_id,
    }));
    const queryClient = renderPage();
    await waitFor(() => expect(screen.getByTestId('mc-persisted-scope')).toHaveTextContent('workspace-1@workspace-rev-1@rule-rev-1'));
    workspaceID = 'workspace-2';
    await act(async () => {
      await queryClient.invalidateQueries({ queryKey: studioV2WorkspaceKeys.bootstrap() });
    });
    await waitFor(() => expect(sourceRuleCandidateAPI.get).toHaveBeenCalledWith('rule-mc', expect.objectContaining({ workspace_id: 'workspace-2' })));
    await waitFor(() => expect(screen.getByTestId('mc-persisted-scope')).toHaveTextContent('workspace-2@workspace-rev-1@rule-rev-1'));
  });
  it('retains a blocked mapping draft and resumes it only after matching evidence returns', async () => {
    let workspaceID = 'workspace-1';
    let workspaceRevision = 'workspace-rev-1';
    vi.mocked(studioV2ServiceMocks.workspace.get).mockImplementation(async () => workspaceFixture({
      id: workspaceID, ordered_device_ids: ['dev-mc', 'dev-modbus'],
      modbus_share: { hydration_state: 'ready', workspace_id: workspaceID, workspace_revision: workspaceRevision, readiness: true,
        status: { workspace_id: workspaceID, enabled: false, port: 5020, address: '127.0.0.1', bind_state: 'disabled', mapping_count: 0, workspace_revision: workspaceRevision } },
    }));
    vi.mocked(modbusShareAPI.status).mockImplementation(async () => ({
      workspace_id: workspaceID, enabled: false, port: 5020, address: '127.0.0.1', bind_state: 'disabled', mapping_count: 0, workspace_revision: workspaceRevision,
    }));
    vi.mocked(studioV2ServiceMocks.mappings.update).mockResolvedValue(mappingFixture({ id: 'legacy-mapping-0', point_id: 'legacy-point-0', rule_id: 'rule-mc', device_id: 'dev-mc', address: 'W0', tag_key: 'edited.tag' }));
    vi.mocked(sourceRuleCandidateAPI.get).mockImplementation(async (_ruleID, scope) => scope?.expected_workspace_revision === 'workspace-rev-2'
      ? candidateSnapshot(compatibleAddresses, { workspace_id: workspaceID, workspace_revision: workspaceRevision, revision_id: 'rule-rev-1' })
      : candidateSnapshot());
    const queryClient = renderPage();
    await waitFor(() => expect(screen.getByTestId('mc-persisted-status')).toHaveTextContent('conflict'));
    fireEvent.click(screen.getByTestId('edit-mc-mapping'));
    await waitFor(() => expect(studioV2ServiceMocks.mappings.update).not.toHaveBeenCalled());
    expect(screen.getByTestId('mc-persisted-mappings')).toHaveTextContent('edited.tag');
    workspaceID = 'workspace-2';
    workspaceRevision = 'workspace-rev-2';
    await act(async () => {
      await queryClient.invalidateQueries({ queryKey: studioV2WorkspaceKeys.bootstrap() });
    });
    await waitFor(() => expect(screen.getByTestId('mc-persisted-status')).toHaveTextContent('compatible'));
    await waitFor(() => expect(studioV2ServiceMocks.mappings.update).toHaveBeenCalledTimes(1));
  });
  it('does not apply late rule, mapping, or database acknowledgements after scope changes', async () => {
    let workspaceID = 'workspace-1';
    let workspaceRevision = 'workspace-rev-1';
    const workspaceForScope = () => workspaceFixture({
      id: workspaceID,
      ordered_device_ids: ['dev-mc', 'dev-modbus'],
      modbus_share: {
        hydration_state: 'ready', workspace_id: workspaceID, workspace_revision: workspaceRevision, readiness: true,
        status: { workspace_id: workspaceID, enabled: false, port: 5020, address: '127.0.0.1', bind_state: 'disabled', mapping_count: 0, workspace_revision: workspaceRevision },
      },
    });
    vi.mocked(studioV2ServiceMocks.workspace.get).mockImplementation(async () => workspaceForScope());
    vi.mocked(modbusShareAPI.status).mockImplementation(async () => ({
      workspace_id: workspaceID, enabled: false, port: 5020, address: '127.0.0.1', bind_state: 'disabled', mapping_count: 0, workspace_revision: workspaceRevision,
    }));
    let resolveRule: ((value: ReturnType<typeof ruleFixture>) => void) | undefined;
    let resolveMapping: ((value: ReturnType<typeof mappingFixture>) => void) | undefined;
    let resolveDatabase: ((value: ReturnType<typeof sqliteConfigFixture>) => void) | undefined;
    let resolveNewCandidate: ((value: SourceRuleCandidateSnapshotView) => void) | undefined;
    vi.mocked(studioV2ServiceMocks.rules.update).mockImplementation(() => new Promise((resolve) => { resolveRule = resolve; }));
    vi.mocked(studioV2ServiceMocks.mappings.update).mockImplementation(() => new Promise((resolve) => { resolveMapping = resolve; }));
    vi.mocked(studioV2ServiceMocks.database.updateConfig).mockImplementation(() => new Promise((resolve) => { resolveDatabase = resolve; }));
    vi.mocked(sourceRuleCandidateAPI.get).mockImplementation(async (_ruleID, scope) => {
      if (scope?.expected_workspace_revision === 'workspace-rev-2') {
        return new Promise((resolve) => { resolveNewCandidate = resolve; });
      }
      return candidateSnapshot(compatibleAddresses);
    });
    const queryClient = renderPage();
    await waitFor(() => expect(screen.getByTestId('mc-persisted-status')).toHaveTextContent('compatible'));
    fireEvent.click(screen.getByTestId('edit-mc-mapping'));
    fireEvent.click(screen.getByTestId('edit-mc-db-connector'));
    fireEvent.click(screen.getByTestId('edit-mc-rule'));
    await waitFor(() => {
      expect(studioV2ServiceMocks.mappings.update).toHaveBeenCalledTimes(1);
      expect(studioV2ServiceMocks.database.updateConfig).toHaveBeenCalledTimes(1);
      expect(studioV2ServiceMocks.rules.update).toHaveBeenCalledTimes(1);
    });
    workspaceID = 'workspace-2';
    workspaceRevision = 'workspace-rev-2';
    await act(async () => {
      await queryClient.invalidateQueries({ queryKey: studioV2WorkspaceKeys.bootstrap() });
    });
    await waitFor(() => expect(screen.getByTestId('mc-persisted-status')).toHaveTextContent('checking'));
    resolveRule?.(ruleFixture({ id: 'rule-mc', device_id: 'dev-mc', start_address: 'W0', count: 17, revision_id: 'rule-rev-2', naming_prefix: 'SERVER_' }));
    resolveMapping?.(mappingFixture({ id: 'legacy-mapping-0', point_id: 'legacy-point-0', rule_id: 'rule-mc', device_id: 'dev-mc', address: 'W0', tag_key: 'server.tag' }));
    resolveDatabase?.(sqliteConfigFixture({ table: 'server_table', updated_at: '2026-10-30T00:01:00Z' }));
    await waitFor(() => {
      expect(studioV2ServiceMocks.rules.update).toHaveBeenCalledTimes(1);
      expect(studioV2ServiceMocks.mappings.update).toHaveBeenCalledTimes(1);
      expect(studioV2ServiceMocks.database.updateConfig).toHaveBeenCalledTimes(1);
    });
    expect(screen.getByTestId('mc-db-connector')).toHaveTextContent('edited_table');
    expect(screen.getByTestId('mc-persisted-mappings')).toHaveTextContent('edited.tag');
    resolveNewCandidate?.(candidateSnapshot(compatibleAddresses, { workspace_id: workspaceID, workspace_revision: 'workspace-rev-2' }));
    await waitFor(() => expect(screen.getByTestId('mc-persisted-status')).toHaveTextContent('compatible'));
    expect(studioV2ServiceMocks.rules.update).toHaveBeenCalledTimes(1);
    expect(studioV2ServiceMocks.mappings.update).toHaveBeenCalledTimes(1);
    expect(studioV2ServiceMocks.database.updateConfig).toHaveBeenCalledTimes(1);
  });
});
