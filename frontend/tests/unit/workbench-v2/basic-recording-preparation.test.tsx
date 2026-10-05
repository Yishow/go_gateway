import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { act, fireEvent, render, screen, waitFor } from '@testing-library/react';
import { beforeEach, expect, it, vi } from 'vitest';
import { BasicSchemaPreparation } from '@/features/datalink/workbench-v2/steps/step4/BasicSchemaPreparation';
import { BasicRecordingPanel } from '@/features/datalink/workbench-v2/steps/step4/BasicRecordingPanel';
import { ConnectorSection } from '@/features/datalink/workbench-v2/steps/step4/ConnectorSection';
import { BasicRecordingEvidence } from '@/features/datalink/workbench-v2/steps/step4/BasicRecordingEvidence';
import { savedGroup, savedGroupState } from '../../fixtures/writeGroupState';
const mocks = vi.hoisted(() => ({ ensure: vi.fn(), start: vi.fn(), refetch: vi.fn(), schema: vi.fn() }));
vi.mock('react-i18next', () => ({ useTranslation: () => ({ t: (key: string) => key }) }));
vi.mock('@/hooks/datalink/useStudioV2WriteGroups', () => ({ useWriteGroupsQuery: () => ({ data: { workspace_id: 'ws-1', workspace_revision: 'wrev-1', groups: [] }, refetch: mocks.refetch }), useWriteGroupDeliveryQuery: () => ({ data: undefined }) }));
vi.mock('@/hooks/datalink/useStudioV2RecordingStart', () => ({ useEnsureBasicManagedMutation: () => ({ mutateAsync: mocks.ensure }), useRecordingStartMutation: () => ({ mutateAsync: mocks.start }), useRecordingStartOperationQuery: () => ({ data: undefined }) }));
vi.mock('@/features/datalink/workbench-v2/steps/step4/writeGroup/GroupSchemaPanel', () => ({ GroupSchemaPanel: (props: unknown) => { mocks.schema(props); return <div data-testid="reused-group-schema">Existing token/confirm panel</div>; } }));
const managed = savedGroup({ basic_managed_device_id: 'dev-1', destination: { ...savedGroup().destination, storage_strategy: 'managed' } });
beforeEach(() => { vi.clearAllMocks(); window.sessionStorage.clear(); mocks.ensure.mockResolvedValue({ group: managed, workspace_revision: 'wrev-2' }); });
function panel(readonly = false) { const client = new QueryClient({ defaultOptions: { queries: { retry: false } } }); return render(<QueryClientProvider client={client}><BasicRecordingPanel state={savedGroupState()} workspaceId="ws-1" readonly={readonly} /></QueryClientProvider>); }
it('prepares a local managed group in Basic and reuses the schema-token panel without starting or DDL on render', async () => {
 panel(); expect(mocks.ensure).not.toHaveBeenCalled(); expect(mocks.start).not.toHaveBeenCalled(); expect(mocks.schema).not.toHaveBeenCalled();
 fireEvent.click(screen.getByTestId('basic-recording-prepare'));
 await screen.findByTestId('reused-group-schema');
 expect(mocks.ensure).toHaveBeenCalledTimes(1); expect(mocks.start).not.toHaveBeenCalled();
 expect(mocks.schema.mock.calls.at(-1)?.[0]).toMatchObject({ group: managed, workspaceId: 'ws-1', workspaceRevision: 'wrev-2', destinationMatches: true, readonly: false });
 expect(screen.queryByTestId('step4-advanced-recording')).not.toBeInTheDocument();
});
it('preparation respects readonly and never ensures a group', () => { panel(true); expect(screen.getByTestId('basic-recording-prepare')).toBeDisabled(); expect(mocks.ensure).not.toHaveBeenCalled(); });
it('shows one effective group configuration and disables unsupported saved and new kinds', () => {
 const onKindChange = vi.fn(), onSelectConnector = vi.fn();
 render(<ConnectorSection connector={savedGroupState().db.connector} connectors={[{ id: 'mysql-pool', name: 'Unsupported saved DB', kind: 'mysql', enabled: true } as never]} onKindChange={onKindChange} onSelectConnector={onSelectConnector} onUpdateConnector={vi.fn()} />);
 for (const kind of ['mysql', 'sqlserver']) expect(screen.getByText(`step4.kind_${kind}`).closest('button')).toBeDisabled();
 expect(screen.getByRole('option', { name: /Unsupported saved DB/ })).toBeDisabled();
 expect(screen.queryByLabelText('step4.field_table')).not.toBeInTheDocument();
 expect(screen.queryByLabelText(/step4.write_mode_upsert/)).not.toBeInTheDocument();
 expect(screen.queryByRole('spinbutton')).not.toBeInTheDocument();
});
it('retains receipt identity inside diagnostics while primary recording status stays visible', async () => {
 render(<BasicRecordingEvidence intervalSeconds={10} deliveryError={false} group={managed} committedEffect={{ group_revision: 'rev-1', connector_revision: 'crev-1', record_id: 'opaque-record', effect_key: 'opaque-effect', payload_digest: 'digest', committed_at: '2026-10-05T00:00:00Z' }} />);
 const identity = screen.getByText('opaque-effect'); expect(identity.closest('details')).not.toBeNull();
 expect(screen.getByText('step4.basic.committed_verified')).toBeVisible();
 await waitFor(() => expect(screen.getByTestId('basic-recording-committed-effect')).toBeInTheDocument());
});

it('ignores a deferred Basic preparation after the device scope changes', async () => {
 let resolve!: (value: { group: typeof managed; workspace_revision: string }) => void;
 const pending = new Promise<{ group: typeof managed; workspace_revision: string }>((accept) => { resolve = accept; });
 mocks.ensure.mockReturnValueOnce(pending);
 const props = { workspaceId: 'ws-1', workspaceRevision: 'wrev-1', disabled: false, readonly: false, destinationMatches: true, ensureGroup: mocks.ensure, onApplied: mocks.refetch };
 const view = render(<BasicSchemaPreparation key="device-a" {...props} />);
 fireEvent.click(screen.getByTestId('basic-recording-prepare'));
 view.rerender(<BasicSchemaPreparation key="device-b" {...props} />);
 await act(async () => { resolve({ group: managed, workspace_revision: 'wrev-2' }); await pending; });
 expect(screen.queryByTestId('reused-group-schema')).not.toBeInTheDocument();
 expect(mocks.schema).not.toHaveBeenCalled();
 expect(screen.getByTestId('basic-recording-prepare')).not.toBeDisabled();
});
