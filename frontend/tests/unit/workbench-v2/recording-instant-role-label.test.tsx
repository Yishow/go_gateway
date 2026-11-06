import { act, cleanup, render, screen, within } from '@testing-library/react';
import { afterEach, beforeEach, describe, expect, it } from 'vitest';
import { I18nextProvider } from 'react-i18next';
import i18n from '@/i18n/config';
import { GroupProfileMembers } from '@/features/datalink/workbench-v2/steps/step4/writeGroup/GroupProfileSummary';
import { buildGroupCandidates } from '@/features/datalink/workbench-v2/state/writeGroup/candidates';
import { groupToDraft } from '@/features/datalink/workbench-v2/state/writeGroup/draft';
import type { WriteGroupSampleReduction } from '@/types/studioV2WriteGroup';
import { savedGroup, savedGroupState } from '../../fixtures/writeGroupState';

const state = savedGroupState();
const candidates = buildGroupCandidates(state).candidates;
const base = savedGroup();
function instant(reduction?: WriteGroupSampleReduction) {
  return groupToDraft(savedGroup({
    layout_version: 'managed_wide_profile_v1',
    recording_profile: { kind: 'instantaneous', interval: 15, unit: 'second', ...(reduction ? { sample_reduction: reduction } : {}) },
    members: [
      { ...base.members[0], value_role: 'gauge' },
      { ...base.members[1], value_role: 'state' },
      { ...base.members[0], point_id: 'pt-2', tag_id: 'tag-2', value_role: 'counter_raw' },
    ],
  }));
}
function view(reduction?: WriteGroupSampleReduction) {
  return <I18nextProvider i18n={i18n}>
    <GroupProfileMembers draft={instant(reduction)} candidates={candidates} workspaceId="ws-1" readonly />
  </I18nextProvider>;
}
beforeEach(async () => { await i18n.changeLanguage('zh-TW'); });
afterEach(async () => { cleanup(); await i18n.changeLanguage('zh-TW'); });

describe('instantaneous source role labels', () => {
  it.each(['latest', undefined] as const)('shows a gauge as latest for %s, including the legacy default', reduction => {
    render(view(reduction));
    const sources = within(screen.getByTestId('group-profile-members'));
    expect(sources.getByText('量測值（最新值）')).toBeVisible();
    expect(sources.queryByText('量測值（平均）')).not.toBeInTheDocument();
    expect(sources.getByText('狀態（最新值）')).toBeVisible();
    expect(sources.getByText('計數器原值（最新值）')).toBeVisible();
  });

  it('updates the gauge label with the draft reduction, and localizes both states', async () => {
    const { rerender } = render(view('latest'));
    const sources = () => within(screen.getByTestId('group-profile-members'));
    expect(sources().getByText('量測值（最新值）')).toBeVisible();
    rerender(view('mean'));
    expect(sources().getByText('量測值（平均）')).toBeVisible();
    expect(sources().getByText('狀態（最新值）')).toBeVisible();
    expect(sources().getByText('計數器原值（最新值）')).toBeVisible();
    await act(async () => { await i18n.changeLanguage('en'); });
    rerender(view('mean'));
    expect(sources().getByText('Measurement (mean)')).toBeVisible();
    rerender(view('latest'));
    expect(sources().getByText('Measurement (latest)')).toBeVisible();
  });
});
