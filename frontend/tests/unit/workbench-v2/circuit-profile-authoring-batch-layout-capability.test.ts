import { describe, expect, it, vi } from 'vitest';
import { executeBatchPlan, type FrozenBatchPlan, type FrozenBatchTarget } from '@/features/datalink/workbench-v2/steps/step4/circuitProfileBatchExecutor';
import type { CircuitProfileCapabilityContext } from '@/features/datalink/workbench-v2/state/circuitProfileCapabilities';
import type { WriteGroupUpdateRequest } from '@/types/studioV2WriteGroup';
import { longProfileGroup } from '../../fixtures/circuitBatchPanel';
import { savedGroup } from '../../fixtures/writeGroupState';

function context(): CircuitProfileCapabilityContext {
  return {
    connectorKind: 'mysql', destination: { connector_id: 'conn-1', connector_revision: 'crev-1' },
    capabilities: { mysql: { isLoading: false, isFetching: false, isError: false, data: {
      connector_id: 'conn-1', connector_revision: 'crev-1', database: 'recording_test',
      server_family: 'mysql', server_version: '8.0', capability_revision: 'cap-1',
      connected: true, metadata_status: 'verified', transactional_engine: true, supports_managed_long: true,
      supports_native_partition: false, supports_group_test_write: false, supports_missing_repair: false, reasons: [],
    } } },
  };
}
function target(id: string): FrozenBatchTarget {
  const group = longProfileGroup({ id: `group-${id}`, circuitId: `cir-${id}`, circuitName: id, interval: 15, points: [0, 1] });
  return { circuit_id: `cir-${id}`, circuit_name: id, instantaneousGroup: group, anyGroup: group };
}
function input() {
  const plan: FrozenBatchPlan = {
    workspaceId: 'ws-1', workspaceRevision: 'wrev-1', connectorRevision: 'crev-1',
    applyInst: true, instInterval: 30, applyCum: false, cumInterval: 10, targets: [target('a'), target('b')],
  };
  return {
    plan, capabilityContext: context(), ensureCircuitProfile: vi.fn(),
    updateWriteGroup: vi.fn(async ({ id, request }: { id: string; request: WriteGroupUpdateRequest }) => ({
      workspace_revision: `next-${id}`, group: savedGroup({ ...request.group, id, revision: `revision-${id}`,
        destination: { ...request.group.destination, database: 'recording_test' },
        members: request.group.members.map((member) => ({ ...member, source_revision: 's1', mapping_revision: 'm1' })),
      }),
    })),
  };
}

describe('every selected batch profile is capability checked before the first mutation', () => {
  it('rejects a later native-monthly target without first updating a compatible plain target', async () => {
    const run = input();
    run.plan.targets[1].instantaneousGroup!.storage_mode = 'native_monthly';
    const original = structuredClone(run.plan);
    await expect(executeBatchPlan(run)).rejects.toThrow('Recording layout is unavailable for this batch.');
    expect(run.updateWriteGroup).not.toHaveBeenCalled();
    expect(run.ensureCircuitProfile).not.toHaveBeenCalled();
    expect(run.plan).toEqual(original);
  });

  it('checks storage inherited by a missing cumulative profile before any create', async () => {
    const run = input();
    run.plan.applyInst = false;
    run.plan.applyCum = true;
    run.plan.targets[1].anyGroup.storage_mode = 'native_monthly';
    await expect(executeBatchPlan(run)).rejects.toThrow('Recording layout is unavailable for this batch.');
    expect(run.updateWriteGroup).not.toHaveBeenCalled();
    expect(run.ensureCircuitProfile).not.toHaveBeenCalled();
  });

  it.each(['missing-context', 'missing-destination', 'missing-id', 'missing-revision', 'missing-kind', 'plan-revision', 'target-id', 'target-revision', 'capability-revision', 'fetching', 'error'] as const)(
    'rejects %s without a mutation or implicit SQLite authority', async (condition) => {
      const run = input();
      if (condition === 'missing-context') {
        const { capabilityContext: _context, ...unscoped } = run;
        await expect(executeBatchPlan(unscoped)).rejects.toThrow('Recording layout is unavailable for this batch.');
      } else {
        if (condition === 'missing-destination') run.capabilityContext.destination = undefined;
        if (condition === 'missing-id') run.capabilityContext.destination!.connector_id = '';
        if (condition === 'missing-revision') run.capabilityContext.destination!.connector_revision = '';
        if (condition === 'missing-kind') run.capabilityContext.connectorKind = undefined;
        if (condition === 'plan-revision') run.plan.connectorRevision = 'stale-revision';
        if (condition === 'target-id') run.plan.targets[1].anyGroup.destination.connector_id = 'foreign-connector';
        if (condition === 'target-revision') run.plan.targets[1].anyGroup.destination.connector_revision = 'stale-revision';
        if (condition === 'capability-revision') run.capabilityContext.capabilities!.mysql!.data!.connector_revision = 'stale-revision';
        if (condition === 'fetching') run.capabilityContext.capabilities!.mysql!.isFetching = true;
        if (condition === 'error') run.capabilityContext.capabilities!.mysql!.isError = true;
        await expect(executeBatchPlan(run)).rejects.toThrow('Recording layout is unavailable for this batch.');
      }
      expect(run.updateWriteGroup).not.toHaveBeenCalled();
      expect(run.ensureCircuitProfile).not.toHaveBeenCalled();
    },
  );

  it('does not convert a saved wide profile into a batch long profile', async () => {
    const run = input();
    run.capabilityContext = { connectorKind: 'sqlite', destination: { connector_id: 'conn-1', connector_revision: 'crev-1' } };
    run.plan.targets[1].anyGroup.layout_version = 'managed_wide_profile_v1';
    await expect(executeBatchPlan(run)).rejects.toMatchObject({ reason: 'batch_layout_unsupported' });
    expect(run.updateWriteGroup).not.toHaveBeenCalled();
    expect(run.plan.targets[1].anyGroup.layout_version).toBe('managed_wide_profile_v1');
  });

  it('checks only requested profiles, preserving IDs and revision chaining', async () => {
    const run = input();
    const unrequested = longProfileGroup({ id: 'cum-a', circuitId: 'cir-a', circuitName: 'a', interval: 10, points: [0, 1], kind: 'cumulative' });
    unrequested.storage_mode = 'native_monthly';
    run.plan.targets[0].cumulativeGroup = unrequested;
    await executeBatchPlan(run);
    expect(run.updateWriteGroup.mock.calls.map(([request]) => request.id)).toEqual(['group-a', 'group-b']);
    expect(run.updateWriteGroup.mock.calls[0][0].request).toMatchObject({
      expected_workspace_revision: 'wrev-1', expected_group_revision: 'rev-group-a-1', expected_connector_revision: 'crev-1',
    });
    expect(run.updateWriteGroup.mock.calls[1][0].request.expected_workspace_revision).toBe('next-group-a');
    expect(unrequested.storage_mode).toBe('native_monthly');
    expect(run.ensureCircuitProfile).not.toHaveBeenCalled();
  });

  it('rejects an unsupported requested cumulative profile even when the instantaneous profile is valid', async () => {
    const run = input();
    run.plan.applyCum = true;
    run.plan.targets[1].cumulativeGroup = longProfileGroup({ id: 'cum-b', circuitId: 'cir-b', circuitName: 'b', interval: 10, points: [0, 1], kind: 'cumulative' });
    run.plan.targets[1].cumulativeGroup.storage_mode = 'native_monthly';
    await expect(executeBatchPlan(run)).rejects.toMatchObject({ reason: 'storage_unsupported' });
    expect(run.updateWriteGroup).not.toHaveBeenCalled();
    expect(run.ensureCircuitProfile).not.toHaveBeenCalled();
  });

  it('allows verified native storage without changing it to plain', async () => {
    const run = input();
    run.plan.targets[1].instantaneousGroup!.storage_mode = 'native_monthly';
    run.capabilityContext.capabilities!.mysql!.data!.supports_native_partition = true;
    await executeBatchPlan(run);
    expect(run.updateWriteGroup.mock.calls[1][0].request.group.storage_mode).toBe('native_monthly');
    expect(run.updateWriteGroup.mock.calls[1][0].id).toBe('group-b');
  });

  it.each(['current', 'stale'] as const)('binds the legacy explicit connector-revision argument to the %s context', async (revision) => {
    const run = input();
    run.plan.connectorRevision = undefined;
    const execution = executeBatchPlan({ ...run, connectorRevision: revision === 'current' ? 'crev-1' : 'old-revision' });
    if (revision === 'stale') {
      await expect(execution).rejects.toMatchObject({ reason: 'capability_unavailable' });
      expect(run.updateWriteGroup).not.toHaveBeenCalled();
    } else {
      await execution;
      expect(run.updateWriteGroup.mock.calls[0][0].request.expected_connector_revision).toBe('crev-1');
    }
  });


  it.each([false, true])('binds SQL Server batch native storage to supports_native_partition=%s', async (supported) => {
    const run = input();
    run.plan.targets[1].instantaneousGroup!.storage_mode = 'native_monthly';
    run.capabilityContext = {
      connectorKind: 'sqlserver', destination: { connector_id: 'conn-1', connector_revision: 'crev-1' },
      capabilities: { sqlserver: { isLoading: false, isFetching: false, isError: false, data: {
        connector_id: 'conn-1', connector_revision: 'crev-1', database: 'recording_test', schema: 'dbo',
        server_version: '16.0.4165.4', server_edition: 'Developer Edition', engine_edition: 3, compatibility_level: 160,
        connected: true, metadata_status: 'verified', capability_revision: 'cap-1', supports_managed_long: true,
        supports_native_partition: supported, supports_group_test_write: false, supports_missing_repair: false, reasons: [],
        permissions: { verified: true, database_view_definition: true, schema_view_definition: true, create_table: true,
          alter_schema: true, select_schema: true, insert_schema: true, delete_schema: true, alter_any_dataspace: true },
      } } },
    };
    if (supported) {
      await executeBatchPlan(run);
      expect(run.updateWriteGroup).toHaveBeenCalledTimes(2);
      expect(run.updateWriteGroup.mock.calls[1][0].request.group.storage_mode).toBe('native_monthly');
    } else {
      await expect(executeBatchPlan(run)).rejects.toMatchObject({ reason: 'storage_unsupported' });
      expect(run.updateWriteGroup).not.toHaveBeenCalled();
    }
    expect(run.ensureCircuitProfile).not.toHaveBeenCalled();
  });

});
