import { render, screen } from '@testing-library/react';
import { describe, expect, it, vi } from 'vitest';
import { RecordingSelectedMembers } from '../../../src/features/datalink/workbench-v2/recording-editor/RecordingSelectedMembers';
import { GroupProfileMembers } from '../../../src/features/datalink/workbench-v2/steps/step4/writeGroup/GroupProfileSummary';
import { newDraft, type EditorMember } from '../../../src/features/datalink/workbench-v2/state/writeGroup/draft';
import type { GroupCandidate } from '../../../src/features/datalink/workbench-v2/state/writeGroup/candidates';
vi.mock('react-i18next', () => ({ useTranslation: () => ({ t: (key: string) => key }) }));
const member: EditorMember = { key: 'P|T', device_id: 'D', point_id: 'P', tag_id: 'T', target_column: 'temperature', required: true, entity_key: '' };
const candidate: GroupCandidate = { ...member, tag_key: 'temperature', label: 'Temperature', device_name: 'PLC A', address: 'D0', target_type: 'float32' };
const renderSource = (source: Pick<EditorMember, 'device_id' | 'point_id'>) => <span data-testid={`source-order-${source.point_id}`}>Canonical {source.device_id}/{source.point_id}: override CDAB</span>;
describe('canonical point details remain attached to existing selected rows', () => {
  it('renders the canonical source control in the normal selected identity cell without changing membership', () => {
    const onChange = vi.fn();
    render(<RecordingSelectedMembers candidates={[candidate]} members={[member]} columns={[]} metadataStatus="exists"
      suggestions={{ assignments: {}, unmatched: [] }} issues={[]} showEntityKey={false} disabled={false} onChange={onChange}
      onRemove={vi.fn()} labels={{ remove: () => 'Remove', unavailable: 'Unavailable', undo: 'Undo' }} renderSource={renderSource} />);
    expect(screen.getByTestId('source-order-P')).toHaveTextContent('Canonical D/P: override CDAB');
    expect(screen.getByTestId('group-member-row-P')).toContainElement(screen.getByTestId('source-order-P'));
    expect(onChange).not.toHaveBeenCalled();
  });
  it.each(['managed_long_v1', 'managed_wide_profile_v1'] as const)('renders the same canonical source control in %s without a calculation/layout conversion', layout_version => {
    const draft = { ...newDraft({ connector_id: 'C', connector_revision: 'C1', table_schema: '', table_name: 'readings' }), members: [member], layout_version, recording_profile: { kind: 'instantaneous' as const, interval: 1, unit: 'second' as const, sample_reduction: 'latest' as const } };
    const before = structuredClone(draft);
    render(<GroupProfileMembers draft={draft} candidates={[candidate]} workspaceId="W" readonly={false} renderSource={renderSource} />);
    expect(screen.getByTestId('source-order-P')).toHaveTextContent('Canonical D/P: override CDAB');
    expect(draft).toEqual(before);
  });
});
