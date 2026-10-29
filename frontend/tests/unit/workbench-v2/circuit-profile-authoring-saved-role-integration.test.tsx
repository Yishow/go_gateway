import { studioV2WorkspaceKeys } from '@/hooks/datalink/keys';
import { guidedReadiness, guidedLifecycle } from '../../fixtures/guidedRecordingPreparation';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { fireEvent, render, screen } from '@testing-library/react';
import { createInstance } from 'i18next';
import { I18nextProvider } from 'react-i18next';
import { beforeEach, expect, it, vi } from 'vitest';
import { BasicRecordingPanel } from '@/features/datalink/workbench-v2/steps/step4/BasicRecordingPanel';
import type { WriteGroup, WriteGroupValueRole } from '@/types/studioV2WriteGroup';
import productionI18n from '@/i18n/config';
import { savedGroupState } from '../../fixtures/writeGroupState';
import { longProfileGroup } from '../../fixtures/circuitBatchPanel';
import { expandSources } from '../../fixtures/sourceReview';
const groupsQuery = { data: { workspace_id: 'ws-1', workspace_revision: 'wrev-1', groups: [] as WriteGroup[] }, isLoading: false, isError: false, refetch: vi.fn() };
const mutations = { ensure: vi.fn(), update: vi.fn(), start: vi.fn(), ensureLegacy: vi.fn(), schema: vi.fn(), measurement: vi.fn() };
let canonicalClient: QueryClient;
const evidenceRefetch = vi.fn();
const evidence = new Map<string, { readiness: ReturnType<typeof guidedReadiness>; lifecycle: ReturnType<typeof guidedLifecycle> }>();
function currentEvidence(id: string) {
  const group = groupsQuery.data.groups.find(value => value.id === id);
  if (!group) return undefined;
  const key = JSON.stringify([groupsQuery.data.workspace_revision, group]);
  if (!evidence.has(key)) evidence.set(key, { readiness: guidedReadiness(group, groupsQuery.data.workspace_revision), lifecycle: guidedLifecycle(group, groupsQuery.data.workspace_revision) });
  return evidence.get(key)!;
}
function publishGroups() {
  canonicalClient?.setQueryData(studioV2WorkspaceKeys.writeGroups(), groupsQuery.data);
  return { data: groupsQuery.data };
}
vi.mock('@/hooks/datalink/useStudioV2WriteGroupLifecycle', () => ({
  useWriteGroupLifecycleQuery: (_workspace: string, id: string) => ({ data: currentEvidence(id)?.lifecycle, isLoading: false, isFetching: false, isError: false, refetch: evidenceRefetch }),
}));
vi.mock('@/hooks/datalink/useStudioV2WriteGroups', () => ({
  useWriteGroupReadinessQuery: (id: string) => ({ data: currentEvidence(id)?.readiness, isLoading: false, isFetching: false, isError: false, refetch: evidenceRefetch }),
  useWriteGroupsQuery: () => groupsQuery, useWriteGroupDeliveryQuery: () => ({ data: undefined }),
  useUpdateWriteGroupMutation: () => ({ mutateAsync: mutations.update, isPending: false }),
  useEnsureCircuitProfileMutation: () => ({ mutateAsync: mutations.ensure, isPending: false }),
  useWriteGroupSchemaPreviewMutation: () => ({ mutateAsync: vi.fn(), isPending: false }),
  useWriteGroupSchemaApplyMutation: () => ({ mutateAsync: mutations.schema, isPending: false }),
  useWriteGroupSchemaOperationQuery: () => ({ data: undefined, refetch: vi.fn() }),
}));
vi.mock('@/hooks/datalink/useStudioV2RecordingStart', () => ({
  useEnsureBasicManagedMutation: () => ({ mutateAsync: mutations.ensureLegacy, isPending: false }),
  useRecordingStartMutation: () => ({ mutateAsync: mutations.start, isPending: false }),
  useRecordingStartOperationQuery: () => ({ data: undefined, refetch: vi.fn() }),
}));
vi.mock('@/hooks/datalink/useStudioV2WorkspaceMeasurements', () => ({
  useStudioV2WorkspaceMeasurementsQuery: () => ({ data: [], isLoading: false, isError: false, refetch: vi.fn() }),
  useCreateStudioV2MeasurementMutation: () => ({ mutateAsync: mutations.measurement, isPending: false }),
  useUpdateStudioV2MeasurementMutation: () => ({ mutateAsync: mutations.measurement, isPending: false }),
}));
const english = 'Saved source roles cannot currently be changed. Layout or destination review preserves this role and its frozen binding.';
const chinese = '目前不支援變更已儲存的來源角色。布局或目的地審查會保留此角色與凍結綁定。';
function saved(role: WriteGroupValueRole) {
  const group = longProfileGroup({ id: 'G1', circuitId: 'cir-saved', circuitName: 'Saved', interval: 15, points: [0, 1] });
  return { ...group, recording_profile: { kind: 'instantaneous' as const, interval: 15, unit: 'second' as const, sample_reduction: 'mean' as const }, members: group.members.map((member, index) => ({ ...member, value_role: index === 0 ? role : 'gauge' as const })) };
}
async function mount(lng = 'en', readonly = false) {
  const state = savedGroupState();
  state.devices.push({ ...state.devices[0], id: 'dev-2', name: 'PLC B' });
  const i18n = createInstance();
  await i18n.init({ lng, resources: productionI18n.options.resources });
  const client = new QueryClient({ defaultOptions: { queries: { retry: false } } });
  canonicalClient = client; client.setQueryDefaults(studioV2WorkspaceKeys.writeGroups(), { structuralSharing: false }); publishGroups();
  const onOpenAdvanced = vi.fn();
  const content = () => <I18nextProvider i18n={i18n}><QueryClientProvider client={client}><BasicRecordingPanel state={state} workspaceId="ws-1" readonly={readonly} onOpenAdvanced={onOpenAdvanced} /></QueryClientProvider></I18nextProvider>;
  return { ...render(content()), content, onOpenAdvanced };
}
function noMutation() { Object.values(mutations).forEach((mutation) => expect(mutation).not.toHaveBeenCalled()); }
beforeEach(() => { vi.clearAllMocks(); sessionStorage.clear(); groupsQuery.data.groups = []; groupsQuery.refetch.mockImplementation(async () => publishGroups()); });
it.each([['en', english], ['zh-TW', chinese]])('explains saved-role unsupported status accessibly in %s', async (lng, description) => {
  groupsQuery.data.groups = [saved('gauge')]; await mount(lng); expandSources();
  const role = screen.getByTestId('circuit-profile-role-tag-0');
  expect(role).toBeDisabled(); expect(role).toHaveValue('gauge');
  expect(role).toHaveAccessibleDescription(description);
  expect(screen.getByTestId('circuit-profile-role-tag-0-locked')).toHaveTextContent(description);
  noMutation();
});
it.each(['state', 'gauge', 'counter_raw'] as const)('defends saved %s against a programmatically dispatched change (not ordinary disabled HTML interaction)', async (role) => {
  const group = saved(role); const original = structuredClone(group); groupsQuery.data.groups = [group];
  await mount(); expandSources();
  const select = screen.getByTestId('circuit-profile-role-tag-0');
  expect(select).toBeDisabled();
  fireEvent.change(select, { target: { value: role === 'state' ? 'gauge' : 'state' } });
  expect(select).toHaveValue(role); expect(groupsQuery.data.groups[0]).toEqual(original); noMutation();
});
it('keeps saved role/source binding through supported Advanced navigation, device return and remount', async () => {
  const group = saved('gauge'); const frozen = structuredClone(group); groupsQuery.data.groups = [group];
  const view = await mount(); expandSources();
  fireEvent.click(screen.getByTestId('basic-recording-open-advanced'));
  expect(view.onOpenAdvanced).toHaveBeenCalledTimes(1);
  noMutation();
  fireEvent.click(screen.getByRole('button', { name: 'PLC B' })); fireEvent.click(screen.getByRole('button', { name: 'PLC A' })); expandSources();
  expect(screen.getByTestId('circuit-profile-role-tag-0')).toHaveValue('gauge');
  expect(screen.getByTestId('circuit-profile-role-tag-0')).toBeDisabled();
  view.unmount(); await mount(); expandSources();
  expect(screen.getByTestId('circuit-profile-role-tag-0')).toHaveAccessibleDescription(english);
  expect(groupsQuery.data.groups[0]).toEqual(frozen); noMutation();
});
it('leaves a new editable role without a misleading saved-role explanation', async () => {
  await mount(); expandSources();
  const select = screen.getByTestId('circuit-profile-role-tag-0');
  expect(select).toBeEnabled(); expect(screen.queryByTestId('circuit-profile-role-tag-0-locked')).not.toBeInTheDocument();
  fireEvent.change(select, { target: { value: 'state' } }); expect(select).toHaveValue('state'); noMutation();
});
it('also rejects a synthetic role change when the unsaved form is readonly', async () => {
  await mount('en', true); expandSources();
  const select = screen.getByTestId('circuit-profile-role-tag-0'); expect(select).toBeDisabled();
  fireEvent.change(select, { target: { value: 'state' } }); expect(select).toHaveValue('gauge'); noMutation();
});
