import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { fireEvent, render, screen } from '@testing-library/react';
import { describe, expect, it, vi } from 'vitest';
import { acquisitionPreviewWorkspaceScope } from '@/features/datalink/workbench-v2/acquisition-preview/acquisitionPreviewScope';
import { Step1Device } from '@/features/datalink/workbench-v2/steps/step1/Step1Device';
import { INITIAL_STATE } from '@/features/datalink/workbench-v2/state/useWorkbenchV2State';
import type { Device } from '@/features/datalink/workbench-v2/state/types';
import type { AcquisitionPlanPreviewPanelProps } from '@/features/datalink/workbench-v2/acquisition-preview/AcquisitionPlanPreviewPanel';
import { locatorRead } from './studio-locator-fixtures';
vi.mock('react-i18next', () => ({ useTranslation: () => ({ t: (key: string) => key }) }));
vi.mock('@/features/datalink/workbench-v2/acquisition-preview/AcquisitionPlanPreviewPanel', () => ({
  AcquisitionPlanPreviewPanel: (props: AcquisitionPlanPreviewPanelProps) => <output data-testid="preview-mounted">{JSON.stringify(props)}</output>,
}));
const device: Device = { id: 'D', name: 'PLC', description: '', protocol: 'modbus_tcp', config: {}, status: 'draft', test: null,
  persisted: true, save_state: 'saved', availability_status: 'available' };
const context = { workspace_id: 'W', expected_workspace_revision: 'R' };
function mount(value = device) {
  const client = new QueryClient({ defaultOptions: { queries: { retry: false } } });
  return render(<QueryClientProvider client={client}><Step1Device state={{ ...INITIAL_STATE, devices: [value] }}
    dispatch={vi.fn()} onContinue={vi.fn()} previewWorkspace={context} /></QueryClientProvider>);
}
describe('existing source surface acquisition scope', () => {
  it('uses the database setup revision, never the distinct source/Share revision', () => {
    expect(acquisitionPreviewWorkspaceScope({ workspaceId: 'W',
      shareStatus: { workspace_id: 'W', workspace_revision: 'source-revision-8c577' },
      databaseConnector: { workspace_id: 'W', setup_revision: 'database-revision-54b559', persisted: true, save_state: 'saved' } }))
      .toEqual({ workspace_id: 'W', expected_workspace_revision: 'database-revision-54b559' });
  });
  it('accepts an explicitly hydrated empty database head but never infers it from a missing head', () => {
    const databaseConnector = { workspace_id: 'W', persisted: true, save_state: 'saved' as const };
    expect(acquisitionPreviewWorkspaceScope({ workspaceId: 'W', databaseConnector: { ...databaseConnector, setup_revision: '' } }))
      .toEqual({ workspace_id: 'W', expected_workspace_revision: '' });
    expect(acquisitionPreviewWorkspaceScope({ workspaceId: 'W', databaseConnector })).toBeUndefined();
  });
  it('fails closed when only a source/Share head exists', () => {
    expect(acquisitionPreviewWorkspaceScope({ workspaceId: 'W', shareStatus: { workspace_id: 'W', workspace_revision: 'source-only' } })).toBeUndefined();
  });
  it('does not borrow another workspace database head', () => {
    expect(acquisitionPreviewWorkspaceScope({ workspaceId: 'W', databaseConnector:
      { workspace_id: 'OTHER', setup_revision: 'db', persisted: true, save_state: 'saved' } })).toBeUndefined();
  });
  it('keeps a proven database setup head when unrelated connection fields are being edited', () => {
    expect(acquisitionPreviewWorkspaceScope({ workspaceId: 'W', databaseConnector:
      { workspace_id: 'W', setup_revision: 'db-last-ack', persisted: true, save_state: 'draft-invalid' } }))
      .toEqual({ workspace_id: 'W', expected_workspace_revision: 'db-last-ack' });
  });
  it('takes the exact locator read revision rather than global share status', () => {
    const locator = locatorRead(1);
    expect(acquisitionPreviewWorkspaceScope({ workspaceId: locator.workspace.id, locatorEntry: locator,
      shareStatus: { workspace_id: locator.workspace.id, workspace_revision: 'GLOBAL' } }))
      .toEqual({ workspace_id: locator.workspace.id, expected_workspace_revision: locator.entry.workspace_revision });
  });
  it('does not fall back to global scope while the locator is suspended or inconsistent', () => {
    const locator = locatorRead(1);
    expect(acquisitionPreviewWorkspaceScope({ workspaceId: locator.workspace.id, locatorEntry: locator, locatorSuspended: true,
      shareStatus: { workspace_id: locator.workspace.id, workspace_revision: 'GLOBAL' } })).toBeUndefined();
    expect(acquisitionPreviewWorkspaceScope({ workspaceId: 'OTHER', locatorEntry: locator })).toBeUndefined();
  });
  it('mounts the read-only preview for exactly the saved selected source', () => {
    mount(); expect(JSON.parse(screen.getByTestId('preview-mounted').textContent!)).toMatchObject({ scope: { ...context, device_id: 'D' } });
  });
  it.each(['saving', 'draft-invalid', 'save-error', 'idle'] as const)('does not preview typed or unconfirmed settings (%s)', save_state => {
    mount({ ...device, save_state });
    expect(JSON.parse(screen.getByTestId('preview-mounted').textContent!)).toEqual({ unavailableReason: 'unsaved' });
  });
  it('gives the full-height diagnostic and preview separate intrinsic grid rows above the footer', () => {
    mount();
    const preview = screen.getByTestId('preview-mounted');
    const diagnostic = screen.getByRole('region', { name: 'step1.panels.diagnostic' });
    expect(preview.parentElement).toBe(diagnostic.parentElement);
    expect(preview.parentElement).toHaveClass('grid', 'content-start', 'gap-6');
    expect(preview.parentElement).not.toContainElement(screen.getByTestId('btn-continue-step1'));
  });
  it('does not reuse a device-specific locator proof after selecting a different local device', () => {
    const locator = locatorRead(1); if (locator.selection.state !== 'selected') throw new Error('fixture');
    const client = new QueryClient({ defaultOptions: { queries: { retry: false } } });
    render(<QueryClientProvider client={client}><Step1Device state={{ ...INITIAL_STATE, devices: [device, { ...device, id: 'B' }] }}
      dispatch={vi.fn()} onContinue={vi.fn()} locatorSelection={locator.selection}
      previewWorkspace={{ workspace_id: 'W', expected_workspace_revision: 'W2' }} /></QueryClientProvider>);
    fireEvent.click(screen.getByTestId('device-tab-B'));
    expect(JSON.parse(screen.getByTestId('preview-mounted').textContent!)).toEqual({ unavailableReason: 'scope' });
  });
  it('does not preview a nonpersisted source even if a local saved flag is present', () => {
    mount({ ...device, persisted: false });
    expect(JSON.parse(screen.getByTestId('preview-mounted').textContent!)).toEqual({ unavailableReason: 'unsaved' });
  });
});
