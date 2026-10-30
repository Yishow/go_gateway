import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { render } from '@testing-library/react';
import { createElement } from 'react';
import { vi } from 'vitest';
import DatalinkWorkbenchV2Page from '../../../src/pages/datalink/workbench-v2/DatalinkWorkbenchV2Page';
import { sourceRuleCandidateAPI } from '../../../src/services/sourceRuleCandidates';
import { modbusShareAPI } from '../../../src/services/datalink';
import { studioV2ServiceMocks } from './helpers/studioV2ServiceMocks';
import {
  deviceFixture,
  dbTargetFixture,
  mappingFixture,
  ruleFixture,
  sqliteConfigFixture,
  workspaceFixture,
} from './helpers/workbenchV2PageHarness';
import type { SourceRuleCandidateSnapshotView } from '../../../src/types/sourceRuleCandidates';

export { deviceFixture, dbTargetFixture, mappingFixture, ruleFixture, sqliteConfigFixture, workspaceFixture, studioV2ServiceMocks, sourceRuleCandidateAPI, modbusShareAPI };
export const legacyAddresses = [
  'W0', 'W1', 'W2', 'W3', 'W4', 'W5', 'W6', 'W7', 'W8', 'W9',
  'W10', 'W11', 'W12', 'W13', 'W14', 'W15', 'W16',
];
export const compatibleAddresses = ['W0', 'W1', 'W2', 'W3', 'W4', 'W5', 'W6', 'W7', 'W8', 'W9', 'WA', 'WB', 'WC', 'WD', 'WE', 'WF', 'W10'];

export function candidateSnapshot(
  addresses = legacyAddresses,
  overrides: Partial<Pick<SourceRuleCandidateSnapshotView, 'workspace_id' | 'workspace_revision' | 'revision_id'>> = {},
): SourceRuleCandidateSnapshotView {
  return {
    source_rule_id: 'rule-mc',
    workspace_id: overrides.workspace_id ?? 'workspace-1',
    workspace_revision: overrides.workspace_revision ?? 'workspace-rev-1',
    revision_id: overrides.revision_id ?? 'rule-rev-1',
    tags: {
      status: 'ready',
      candidates: addresses.map((address, index) => ({
        id: `candidate-${index}`,
        identity: {
          source_rule_id: 'rule-mc',
          candidate_type: 'tags',
          candidate_kind: 'source_point',
          derived_from_rule_address: 'W0',
        },
        proposed_signature: `signature-${index}`,
        address,
        point_id: `legacy-point-${index}`,
        tag_id: `legacy-tag-${index}`,
        mapping_id: `legacy-mapping-${index}`,
        tag_key: `legacy.w${index}`,
        display_name: `Legacy W${index}`,
        data_type: 'int16',
        status: 'active',
      })),
    },
    database_outputs: { status: 'ready', candidates: [] },
    local_modbus_outputs: { status: 'ready', candidates: [] },
  };
}

export function renderPage() {
  const queryClient = new QueryClient({ defaultOptions: { queries: { retry: false } } });
  render(createElement(QueryClientProvider, { client: queryClient }, createElement(DatalinkWorkbenchV2Page)));
  return queryClient;
}

export function setupPersistedMCMocks(): void {
  vi.mocked(studioV2ServiceMocks.workspace.get).mockResolvedValue(
    workspaceFixture({
      ordered_device_ids: ['dev-mc', 'dev-modbus'],
      modbus_share: {
        hydration_state: 'ready',
        workspace_id: 'workspace-1',
        workspace_revision: 'workspace-rev-1',
        readiness: true,
        status: {
          workspace_id: 'workspace-1',
          enabled: false,
          port: 5020,
          address: '127.0.0.1',
          bind_state: 'disabled',
          mapping_count: 0,
          workspace_revision: 'workspace-rev-1',
        },
      },
    }),
  );
  vi.mocked(studioV2ServiceMocks.devices.list).mockResolvedValue([
    deviceFixture({ id: 'dev-mc', protocol: 'mc_3e' }),
    deviceFixture({ id: 'dev-modbus', protocol: 'modbus_tcp' }),
  ]);
  vi.mocked(studioV2ServiceMocks.rules.list).mockResolvedValue([
    ruleFixture({
      id: 'rule-mc',
      device_id: 'dev-mc',
      workspace_id: 'workspace-1',
      start_address: 'W0',
      count: 17,
      revision_id: 'rule-rev-1',
      naming_prefix: 'LEGACY_',
    }),
    ruleFixture({
      id: 'rule-modbus',
      device_id: 'dev-modbus',
      start_address: '40001',
      count: 1,
      revision_id: 'rule-modbus-rev-1',
      naming_prefix: 'MODBUS_',
    }),
  ]);
  vi.mocked(studioV2ServiceMocks.mappings.list).mockResolvedValue([
    ...legacyAddresses.map((address, index) => mappingFixture({
      id: `legacy-mapping-${index}`,
      point_id: `legacy-point-${index}`,
      rule_id: 'rule-mc',
      device_id: 'dev-mc',
      address,
      tag_id: `legacy-tag-${index}`,
      tag_key: `legacy.w${index}`,
    })),
    mappingFixture({
      id: 'modbus-mapping-0',
      point_id: 'modbus-point-0',
      rule_id: 'rule-modbus',
      device_id: 'dev-modbus',
      address: '40001',
      tag_id: 'modbus-tag-0',
      tag_key: 'modbus.40001',
    }),
  ]);
  vi.mocked(studioV2ServiceMocks.database.getConfig).mockResolvedValue(sqliteConfigFixture());
  vi.mocked(studioV2ServiceMocks.database.listTargets).mockResolvedValue([dbTargetFixture({ point_id: 'legacy-point-0' })]);
  vi.mocked(sourceRuleCandidateAPI.get).mockResolvedValue(candidateSnapshot());
  vi.mocked(modbusShareAPI.status).mockResolvedValue({
    workspace_id: 'workspace-1',
    enabled: false,
    port: 5020,
    address: '127.0.0.1',
    bind_state: 'disabled',
    mapping_count: 0,
    workspace_revision: 'workspace-rev-1',
  });
}
