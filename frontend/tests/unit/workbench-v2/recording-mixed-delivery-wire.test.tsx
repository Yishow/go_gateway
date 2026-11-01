import { render, screen } from '@testing-library/react';
import { expect, it, vi } from 'vitest';
import wire from '../../fixtures/recording-mixed-delivery-run18.json';
import { savedGroup } from '../../fixtures/writeGroupState';
import { studioV2DatalinkApi } from '../../../src/services/studioV2Workspace';
import { studioV2WorkspaceWriteGroupsAPI } from '../../../src/services/studioV2WorkspaceWriteGroups';
import { parseWriteGroupDeliveryData } from '../../../src/utils/studioV2WriteGroupDeliveryJson';
import { RecordingMemberResult } from '../../../src/features/datalink/workbench-v2/recording-editor/RecordingMemberResult';
vi.mock('../../../src/services/studioV2Workspace', () => ({ studioV2DatalinkApi: { get: vi.fn() } }));
vi.mock('react-i18next', () => ({ useTranslation: () => ({ t: (key: string) => key }) }));

it('accepts the complete run18 wire through outer parser and service into honest skipped/not-due result UI', async () => {
  expect(parseWriteGroupDeliveryData(wire.data)).not.toBeNull();
  vi.mocked(studioV2DatalinkApi.get).mockResolvedValueOnce({ data: wire });
  const delivery = await studioV2WorkspaceWriteGroupsAPI.delivery(wire.data.group_id);
  render(<RecordingMemberResult group={savedGroup({ id: wire.data.group_id, applied_revision: wire.data.source_evidence.group_revision,
    record_id: '9644d4a5-0b78-4b8f-9f34-f334c2ef78da', recording_profile: { kind: 'mixed', result_shape: 'member-calculation-v1' } })} delivery={delivery} />);
  const counterPoint = wire.data.source_evidence.sources[0].point_id;
  expect(screen.getByTestId(`member-result-${counterPoint}-latest`)).toHaveTextContent('result.not_due');
  expect(screen.getByTestId(`member-result-${counterPoint}-latest`)).toHaveTextContent('result.null');
  expect(screen.queryByTestId(`member-result-window-${counterPoint}-latest`)).not.toBeInTheDocument();
  expect(screen.getAllByText('result.delivery.window_skipped')).toHaveLength(3);
  expect(screen.getAllByText(/result.missing/)).toHaveLength(2);
  expect(screen.queryByText('result.committed')).not.toBeInTheDocument();
  expect(screen.queryByText('result.unavailable')).not.toBeInTheDocument();
});
it('rejects emitted mixed rows with absent effect and payload proof', () => {
  const value = structuredClone(wire.data);
  for (const source of value.source_evidence.sources) {
    source.latest.calculation.row_emitted = true;
    source.latest.calculation.outcome_kind = 'row'; source.latest.delivery = 'queued';
  }
  expect(parseWriteGroupDeliveryData(value)).toBeNull();
});
it('rejects forged not-due window sentinels', () => {
  const value = structuredClone(wire.data);
  value.source_evidence.sources[0].latest.window_start = '2026-10-31T17:34:00Z';
  expect(parseWriteGroupDeliveryData(value)).toBeNull();
});
it('does not permit blank proof identities for another profile', () => {
  const value = structuredClone(wire.data);
  value.source_evidence.sources[0].latest.profile_kind = 'cumulative';
  expect(parseWriteGroupDeliveryData(value)).toBeNull();
});
