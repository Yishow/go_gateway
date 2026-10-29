import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { cleanup, fireEvent, render, screen, waitFor, within } from '@testing-library/react';
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';
import i18n from '../../../src/i18n/config';
import { RuntimeSettingsLink } from '../../../src/features/datalink/runtime-dashboard/RuntimeSettingsLink';
import { RuntimeRecordingRulesPanel } from '../../../src/features/datalink/runtime-dashboard/RuntimeRecordingRulesPanel';
import { RuntimeSourceRuleList } from '../../../src/features/datalink/runtime-dashboard/RuntimeSourceRuleList';
import { RuntimeSetupContextPanel } from '../../../src/features/datalink/runtime-dashboard/RuntimeSetupContextPanel';
import { studioV2DatalinkApi } from '../../../src/services/studioV2Workspace';
import type { EntityLocator } from '../../../src/types/runtimeRuleEvidence';
import { evidenceContext, evidenceScope, ruleEvidenceFixture } from '../services/runtimeRuleEvidenceFixture';
vi.mock('../../../src/services/studioV2Workspace', () => ({ studioV2DatalinkApi: { get: vi.fn() } }));
const base: EntityLocator = { step: 1, workspace_id: 'W', device_id: 'D', rule_id: null, mapping_id: null, group_id: null, group_revision: null };
function requestParams(value: unknown): Record<string, unknown> {
  if (value === null || typeof value !== 'object' || Array.isArray(value)) throw new Error('Expected scoped request parameters');
  return Object.fromEntries(Object.entries(value));
}
const clients: QueryClient[] = [];
function mount(view: React.ReactNode) { const client = new QueryClient({ defaultOptions: { queries: { retry: false, gcTime: 0 } } }); clients.push(client); return render(<QueryClientProvider client={client}>{view}</QueryClientProvider>); }
beforeEach(async () => { await i18n.changeLanguage('en'); vi.mocked(studioV2DatalinkApi.get).mockReset(); });
afterEach(() => { cleanup(); clients.splice(0).forEach((client) => client.clear()); });
describe('runtime exact settings links', () => {
  it.each([
    [base, '/studio/v2?entry=locator&workspace_id=W&step=1&device_id=D'],
    [{ ...base, step: 2, rule_id: 'R' }, '/studio/v2?entry=locator&workspace_id=W&step=2&device_id=D&rule_id=R'],
    [{ ...base, step: 3, mapping_id: 'M' }, '/studio/v2?entry=locator&workspace_id=W&step=3&device_id=D&mapping_id=M'],
    [{ ...base, step: 4, group_id: 'G', group_revision: 'A' }, '/studio/v2?entry=locator&workspace_id=W&step=4&device_id=D&group_id=G&group_revision=A'],
  ] as [EntityLocator, string][])('uses only exact owning-step IDs %j without fetching', (locator, target) => {
    const navigate = vi.fn(); render(<RuntimeSettingsLink locator={locator} navigateTo={navigate}>Settings</RuntimeSettingsLink>);
    const link = screen.getByRole('link', { name: 'Settings' }); expect(link).toHaveAttribute('href', target);
    link.focus(); expect(link).toHaveFocus(); fireEvent.click(link, { detail: 0 }); expect(navigate).toHaveBeenCalledWith(target);
    expect(studioV2DatalinkApi.get).not.toHaveBeenCalled();
  });
  it.each([{ ctrlKey: true }, { metaKey: true }, { shiftKey: true }, { altKey: true }, { button: 1 }])('preserves browser modifiers %j', (event) => {
    const navigate = vi.fn(); render(<RuntimeSettingsLink locator={base} navigateTo={navigate}>Settings</RuntimeSettingsLink>);
    fireEvent.click(screen.getByRole('link'), event); expect(navigate).not.toHaveBeenCalled();
  });
  it.each([null, { ...base, device_id: 'secret.example' }, { ...base, step: 3, rule_id: 'R', mapping_id: 'M' }] as (EntityLocator | null)[])('renders no URL for invalid locator %j', (locator) => {
    const { container } = render(<RuntimeSettingsLink locator={locator}>Settings</RuntimeSettingsLink>);
    expect(screen.queryByRole('link')).not.toBeInTheDocument(); expect(container.innerHTML).not.toContain('secret.example');
  });
  it('offers device settings only from matching explicit safe scope', () => {
    const props = { setupContext: { source_rules: [], mappings: [], database_targets: [] }, selectedDevice: { device_id: 'D', name: 'Meter', protocol: 'modbus_tcp', running: false, availability_status: 'available' as const }, evidenceScope };
    vi.mocked(studioV2DatalinkApi.get).mockReturnValue(new Promise(() => undefined));
    const rendered = mount(<RuntimeSetupContextPanel {...props} />);
    expect(screen.getByRole('link', { name: 'Open device settings' })).toHaveAttribute('href', '/studio/v2?entry=locator&workspace_id=W&step=1&device_id=D');
    rendered.rerender(<QueryClientProvider client={clients[0]}><RuntimeSetupContextPanel {...props} evidenceScope={{ ...evidenceScope, device_id: 'E' }} /></QueryClientProvider>);
    expect(screen.queryByRole('link', { name: 'Open device settings' })).not.toBeInTheDocument();
  });
  it('keeps saved C and effective A hints separate, without per-card reads or intake claims', async () => {
    const view = ruleEvidenceFixture(); view.groups[0].saved!.members.push({ ...view.groups[0].saved!.members[0], device_id: 'Other' });
    vi.mocked(studioV2DatalinkApi.get).mockResolvedValue({ data: { success: true, data: evidenceContext(view) } });
    mount(<RuntimeRecordingRulesPanel scope={evidenceScope} selectedDeviceId="D" />);
    const saved = await screen.findByRole('link', { name: 'Open saved group settings' });
    expect(saved).toHaveAttribute('href', '/studio/v2?entry=locator&workspace_id=W&step=4&device_id=D&group_id=G&group_revision=C');
    expect(screen.getByRole('link', { name: 'Inspect effective group settings' })).toHaveAttribute('href', '/studio/v2?entry=locator&workspace_id=W&step=4&device_id=D&group_id=G&group_revision=A');
    expect(screen.getByText(/Shared group settings/)).toBeInTheDocument();
    expect(screen.getByText(/Revision hints are informational/)).toBeInTheDocument();
    expect(studioV2DatalinkApi.get).toHaveBeenCalledTimes(1);
    expect(screen.getByRole('region', { name: 'Recording rules' })).not.toHaveTextContent('SQL committed');
  });
  it.each(['stale', 'unavailable'] as const)('withholds settings actions for %s recording observation', async (freshness) => {
    const view = ruleEvidenceFixture(); view.freshness = freshness;
    vi.mocked(studioV2DatalinkApi.get).mockResolvedValue({ data: { success: true, data: evidenceContext(view) } });
    mount(<RuntimeRecordingRulesPanel scope={evidenceScope} selectedDeviceId="D" />);
    await screen.findByText('meter_draft'); expect(screen.queryByRole('link')).not.toBeInTheDocument();
    expect(studioV2DatalinkApi.get).toHaveBeenCalledTimes(1);
  });
  it('keeps saved settings available without inventing an effective revision link', async () => {
    const view = ruleEvidenceFixture(); view.groups[0].proof = 'unavailable';
    vi.mocked(studioV2DatalinkApi.get).mockResolvedValue({ data: { success: true, data: evidenceContext(view) } });
    mount(<RuntimeRecordingRulesPanel scope={evidenceScope} selectedDeviceId="D" />);
    await screen.findByRole('link', { name: 'Open saved group settings' });
    expect(screen.queryByRole('link', { name: 'Inspect effective group settings' })).not.toBeInTheDocument();
    expect(screen.queryByText('meter_old')).not.toBeInTheDocument();
  });
  it('localizes exact settings links without exposing entity labels in URLs', async () => {
    await i18n.changeLanguage('zh-TW'); const view = ruleEvidenceFixture();
    vi.mocked(studioV2DatalinkApi.get).mockResolvedValue({ data: { success: true, data: evidenceContext(view) } });
    mount(<RuntimeRecordingRulesPanel scope={evidenceScope} selectedDeviceId="D" />);
    const link = await screen.findByRole('link', { name: '開啟已儲存群組設定' });
    expect(link).toHaveAttribute('href', '/studio/v2?entry=locator&workspace_id=W&step=4&device_id=D&group_id=G&group_revision=C');
    expect(link.getAttribute('href')).not.toContain('meter_draft');
  });
  it('uses source page identities for rule and mapping links without detail fan-out', async () => {
    vi.mocked(studioV2DatalinkApi.get).mockImplementation(async (_url, config) => {
      const kind = requestParams(config?.params).source_list;
      return { data: { success: true, data: { workspace_id: 'W', source_page: { contract_version: 1, ...evidenceScope, kind, query: '', limit: 50, list_revision: 'L1', next_cursor: null,
        rules: kind === 'rules' ? [{ id: 'R', revision_id: 'R1', device_id: 'D', start_address: '40001', count: 1, data_type: 'int16', naming_prefix: 'Meter', enabled: true }] : [],
        mappings: kind === 'mappings' ? [{ id: 'M', rule_id: 'R', rule_revision: 'R1', device_id: 'D', point_id: 'P', address: '40001', enabled: true }] : [] } } } };
    });
    mount(<RuntimeSourceRuleList scope={evidenceScope} />);
    await waitFor(() => expect(screen.getByRole('link', { name: 'Open rule settings' })).toBeInTheDocument());
    expect(screen.getByRole('link', { name: 'Open rule settings' })).toHaveAttribute('href', '/studio/v2?entry=locator&workspace_id=W&step=2&device_id=D&rule_id=R');
    expect(screen.getByRole('link', { name: 'Open mapping settings' })).toHaveAttribute('href', '/studio/v2?entry=locator&workspace_id=W&step=3&device_id=D&mapping_id=M');
    expect(studioV2DatalinkApi.get).toHaveBeenCalledTimes(2);
    expect(within(screen.getByRole('region', { name: 'Source rules' })).getByText('Saved-only list; effective configuration requires separate evidence.')).toBeInTheDocument();
  });
});
