import { cleanup, render, screen, within } from '@testing-library/react';
import { afterEach, expect, it, vi } from 'vitest';
import { I18nextProvider } from 'react-i18next';
import i18n from '@/i18n/config';
import { GroupIdentityFields } from '@/features/datalink/workbench-v2/steps/step4/writeGroup/GroupIdentityFields';
import { GroupMemberTable } from '@/features/datalink/workbench-v2/steps/step4/writeGroup/GroupMemberTable';
import { groupToDraft } from '@/features/datalink/workbench-v2/state/writeGroup/draft';
import type { GroupCandidate } from '@/features/datalink/workbench-v2/state/writeGroup/candidates';
import { savedGroup } from '../../fixtures/writeGroupState';

afterEach(cleanup);
it.each([['en', 'System managed'], ['zh-TW', '系統管理']] as const)('renders actual %s managed table and column translations without exposing keys', async (language, label) => {
  await i18n.changeLanguage(language);
  const change = vi.fn(), group = savedGroup(); group.destination.storage_strategy = 'managed';
  const candidate: GroupCandidate = { key: 'point', device_id: 'dev-1', point_id: 'point', tag_id: 'tag', tag_key: 'temperature', label: 'Temperature', device_name: 'PLC', address: '40001', target_type: 'float64' };
  render(<I18nextProvider i18n={i18n}>
    <GroupIdentityFields scopedPresentation draft={groupToDraft(group)} readonly={false} locked onChange={change} onStorageStrategy={change} />
    <GroupMemberTable scopedPresentation managed candidates={[candidate]} excluded={[]} members={[{ ...candidate, entity_key: '', target_column: 'v_exact', required: true }]}
      columns={[{ name: 'v_exact', type: 'REAL', nullable: true, primary_key: false }]} metadataStatus="exists" suggestions={{ assignments: {}, unmatched: [] }}
      issues={[]} showEntityKey={false} disabled={false} onChange={change} />
  </I18nextProvider>);
  expect(within(screen.getByTestId('group-managed-table')).getByText(label)).toBeInTheDocument();
  expect(within(screen.getByTestId('group-managed-column-point')).getByText(label)).toBeInTheDocument();
  expect(screen.getByTestId('group-managed-table')).toHaveTextContent(group.destination.table_name);
  expect(screen.getByTestId('group-managed-column-point')).toHaveTextContent('v_exact');
  expect(document.body.textContent).not.toContain('simple.systemManaged'); expect(document.body.textContent).not.toContain('compact.systemManaged'); expect(change).not.toHaveBeenCalled();
});
