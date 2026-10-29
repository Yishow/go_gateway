import * as React from 'react';
import { QueryClient, QueryClientProvider, useQueryClient } from '@tanstack/react-query';
import { studioV2WorkspaceKeys } from '@/hooks/datalink/keys';
import { clickEnabled, guidedLifecycle, guidedReadiness } from '../../fixtures/guidedRecordingPreparation';
import { parseRecordingStartOperation } from '@/utils/studioV2RecordingStartJson';
import { render, screen, waitFor } from '@testing-library/react';
import { beforeEach, expect, it, vi } from 'vitest';
import { CircuitProfileBasicPanel } from '@/features/datalink/workbench-v2/steps/step4/CircuitProfileBasicPanel';
import { savedGroup, savedGroupState } from '../../fixtures/writeGroupState';
import type { WriteGroup } from '@/types/studioV2WriteGroup';
import type { RecordingStartOperation } from '@/types/studioV2RecordingStart';

const { query, measurements, ensure, update, start, operationQuery, schemaRender } = vi.hoisted(() => ({
  query: { data: { workspace_id: 'ws-1', workspace_revision: 'ws-old', groups: [] as WriteGroup[] }, refetch: vi.fn() },
  measurements: { data: [] as unknown[] | undefined, isLoading: false, isError: false, refetch: vi.fn() },
  ensure: vi.fn(), update: vi.fn(), start: vi.fn(), schemaRender: vi.fn(),
  operationQuery: { data: undefined as RecordingStartOperation | undefined, refetch: vi.fn() },
}));
vi.mock('react-i18next', () => ({ useTranslation: () => ({ t: (key: string) => key }) }));
function useCanonicalGroups() { const client = useQueryClient(); React.useLayoutEffect(() => { client.setQueryData(studioV2WorkspaceKeys.writeGroups(), query.data); }, [client]); return query; }
function canonical(id: string) { const group = query.data.groups.find(item => item.id === id); if (!group) throw new Error('Missing canonical group'); return group; }
vi.mock('@/services/studioV2WriteGroupLifecycle', () => ({ studioV2WriteGroupLifecycleAPI: { lifecycle: vi.fn(async (id: string) => guidedLifecycle(canonical(id), query.data.workspace_revision)) } }));
vi.mock('@/services/studioV2WorkspaceWriteGroups', () => ({ studioV2WorkspaceWriteGroupsAPI: { readiness: vi.fn(async (id: string) => guidedReadiness(canonical(id), query.data.workspace_revision)) } }));
vi.mock('@/hooks/datalink/useStudioV2WriteGroups', async importOriginal => ({
  ...await importOriginal<typeof import('@/hooks/datalink/useStudioV2WriteGroups')>(),
  useWriteGroupsQuery: useCanonicalGroups,
  useWriteGroupDeliveryQuery: () => ({ data: undefined, isError: false, isLoading: false }),
  useEnsureCircuitProfileMutation: () => ({ mutateAsync: ensure, isPending: false }),
  useUpdateWriteGroupMutation: () => ({ mutateAsync: update, isPending: false }),
}));
vi.mock('@/hooks/datalink/useStudioV2RecordingStart', () => ({
  useRecordingStartMutation: () => ({ mutateAsync: start, isPending: false }),
  useRecordingStartOperationQuery: () => operationQuery,
}));
vi.mock('@/hooks/datalink/useStudioV2WorkspaceMeasurements', () => ({
  useStudioV2WorkspaceMeasurementsQuery: () => measurements,
  useCreateStudioV2MeasurementMutation: () => ({ mutateAsync: vi.fn() }),
  useUpdateStudioV2MeasurementMutation: () => ({ mutateAsync: vi.fn() }),
}));
vi.mock('@/features/datalink/workbench-v2/steps/step4/writeGroup/GroupSchemaPanel', async importOriginal => {
  const actual = await importOriginal<typeof import('@/features/datalink/workbench-v2/steps/step4/writeGroup/GroupSchemaPanel')>();
  return { GroupSchemaPanel: (props: React.ComponentProps<typeof actual.GroupSchemaPanel>) => { schemaRender(props); return <actual.GroupSchemaPanel {...props} />; } };
});

const state = savedGroupState();
function group(): WriteGroup {
  return savedGroup({
    id: 'group-inst', layout_version: 'managed_long_v1', storage_mode: 'plain',
    members: [{ ...savedGroup().members[0], circuit_id: 'opaque-A', circuit_name: 'Water', target_column: '', required: false, measurement_id: '' }],
    recording_profile: { kind: 'instantaneous', interval: 15, unit: 'second' },
    destination: { ...savedGroup().destination, storage_strategy: 'managed', schema_revision: 'rev-1', schema_digest: 'schema-proof' },
    row_policy: { interval_seconds: 15, allowed_lateness_seconds: 0, incomplete_policy: 'partial' },
    write_policy: { mode: 'append', dedupe_capability: 'receipt', recreate_missing: true },
  });
}
function operation(): RecordingStartOperation {
  const parsed = parseRecordingStartOperation({
    operation_id: 'operation-A', workspace_id: 'ws-1', setup_revision: 'ws-old', action: 'recording_start', status: 'succeeded', stage: 'complete',
    intent_digest: 'a'.repeat(64), device_ids: ['dev-1'], devices: [{ device_id: 'dev-1', activated: true }],
    groups: [{ group_id: 'group-inst', group_revision: 'rev-1', applied_revision: 'rev-1', saved: true, applied: true, ready: true }],
    created_at: '2026-10-24T00:00:00Z', updated_at: '2026-10-24T00:00:01Z',
  });
  if (!parsed) throw new Error('Invalid operation fixture'); return parsed;
}
function panel(props: Partial<React.ComponentProps<typeof CircuitProfileBasicPanel>> = {}) {
  return <QueryClientProvider client={new QueryClient({ defaultOptions: { queries: { retry: false } } })}><CircuitProfileBasicPanel state={state} workspaceId="ws-1" readonly={false} {...props} /></QueryClientProvider>;
}
beforeEach(() => {
  vi.clearAllMocks(); window.sessionStorage.clear();
  query.data = { workspace_id: 'ws-1', workspace_revision: 'ws-old', groups: [group()] };
  query.refetch.mockImplementation(async () => ({ data: query.data }));
  ensure.mockImplementation(async () => ({ group: group(), workspace_revision: 'ws-write' }));
  update.mockImplementation(async () => ({ group: group(), workspace_revision: 'ws-write' }));
  start.mockResolvedValue(operation());
  measurements.data = []; measurements.isLoading = false; measurements.isError = false;
  operationQuery.data = undefined;
});

it('keeps a saved SQLite main schema clean when the connector schema is empty', async () => {
  query.data.groups[0].destination.table_schema = 'main';
  const sqliteState = { ...state, db: { ...state.db, connector: { ...state.db.connector, kind: 'sqlite' as const, schema: '' } } };
  render(panel({ state: sqliteState }));
  await waitFor(() => expect(schemaRender).toHaveBeenCalled());
  expect(schemaRender.mock.lastCall?.[0].dirty).toBe(false);
  // Guide contract: a clean, schema-verified saved scope offers review, not a duplicate save command.
  await waitFor(() => expect(screen.getByTestId('circuit-profile-start')).toBeEnabled());
  expect(screen.queryByTestId('circuit-profile-prepare')).not.toBeInTheDocument();
  expect(update).not.toHaveBeenCalled();
  expect(ensure).not.toHaveBeenCalled();
});

it('does not send the database setup revision as a Share workspace proof', async () => {
  render(panel());
  await clickEnabled('circuit-profile-start');
  expect(start).not.toHaveBeenCalled();
  await clickEnabled('recording-start-confirm');
  await waitFor(() => expect(start).toHaveBeenCalledOnce());
  expect(start.mock.calls[0][0]).toMatchObject({ expected_workspace_revision: 'ws-old' });
  expect(start.mock.calls[0][0]).not.toHaveProperty('workspace_revision');
  expect(start.mock.calls[0][0]).not.toHaveProperty('settings_revision');
  expect(start.mock.calls[0][0]).not.toHaveProperty('readiness_token');
});

it('uses the separate Share projection proof when all its revisions are available', async () => {
  render(panel({ readinessToken: 'share-proof', settingsRevision: 'settings-7', shareWorkspaceRevision: 'share-workspace-9' }));
  await clickEnabled('circuit-profile-start');
  expect(start).not.toHaveBeenCalled();
  await clickEnabled('recording-start-confirm');
  await waitFor(() => expect(start).toHaveBeenCalledOnce());
  expect(start.mock.calls[0][0]).toMatchObject({ expected_workspace_revision: 'ws-old', readiness_token: 'share-proof', settings_revision: 'settings-7', workspace_revision: 'share-workspace-9' });
});
