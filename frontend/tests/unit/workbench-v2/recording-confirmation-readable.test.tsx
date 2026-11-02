import { cleanup, fireEvent, render, screen, within } from '@testing-library/react';
import { afterEach, beforeEach, expect, it, vi } from 'vitest';
import { I18nextProvider } from 'react-i18next';
import i18n from '@/i18n/config';
import { GroupReadinessPanel } from '@/features/datalink/workbench-v2/steps/step4/writeGroup/GroupReadinessPanel';
import { SavedRecordingActions } from '@/features/datalink/workbench-v2/steps/step4/writeGroup/lifecycle/SavedRecordingActions';
import type { WriteGroupReadiness } from '@/types/studioV2WriteGroup';
import type { WriteGroupLifecycle } from '@/types/studioV2WriteGroupLifecycle';
afterEach(cleanup);
beforeEach(async () => { await i18n.changeLanguage('zh-TW'); });
const scope = { workspaceId: 'W', workspaceRevision: 'WV', groupId: 'G', groupRevision: 'GV' };
const known = 'preview and explicitly confirm the saved managed table before applying';
function readiness(): WriteGroupReadiness {
  return { workspace_id: 'W', workspace_revision: 'WV', group_id: 'G', group_revision: 'GV', applied_revision: '', config_ready: true, schema_ready: false, ready: false,
    issues: [{ code: 'schema-confirmation-required', severity: 'blocking', step: 'Step 4', scope: 'G', message: known }] };
}
it('explains the exact compact schema blocker in Chinese without hiding the original issue or changing readiness', () => {
  const value = readiness();
  render(<I18nextProvider i18n={i18n}><GroupReadinessPanel readiness={value} loading={false} failed={false} dirty={false} onRetry={vi.fn()} {...{ compact: true }} /></I18nextProvider>);
  expect(screen.getByText('請先檢查已儲存的目的地資料表，並明確確認準備內容，再開始記錄。', { exact: false })).toBeVisible();
  const original = screen.getByText(`schema-confirmation-required: ${known}`);
  expect(original).not.toBeVisible(); fireEvent.click(original.closest('details')!.querySelector('summary')!); expect(original).toBeVisible();
  expect(value.ready).toBe(false); expect(value.issues[0].message).toBe(known);
});
it('keeps unknown warnings visible as unverified issues rather than inventing a blocker translation', () => {
  const value = readiness(); value.ready = true; value.schema_ready = true; value.issues = [{ code: 'future-warning', severity: 'warning', step: 'Step 4', scope: 'G', message: 'Original warning from server' }];
  render(<I18nextProvider i18n={i18n}><GroupReadinessPanel readiness={value} loading={false} failed={false} dirty={false} onRetry={vi.fn()} {...{ compact: true }} /></I18nextProvider>);
  expect(screen.getByText('有一項記錄狀態需要確認。請展開下方詳細資訊，查看原始原因。', { exact: false })).toBeVisible();
  const original = screen.getByText('future-warning: Original warning from server'); expect(original).not.toBeVisible();
  fireEvent.click(original.closest('details')!.querySelector('summary')!); expect(original).toBeVisible(); expect(value.ready).toBe(true);
});
it('groups only identical compact action reasons while every original permission and accessible reason remains intact', () => {
  const onAction = vi.fn();
  const lifecycle: WriteGroupLifecycle = { workspace_id: 'W', workspace_revision: 'WV', group_id: 'G', group_revision: 'GV', applied_revision: '', status: 'ready', connector_id: 'C', connector_revision: 'CV',
    history: { canonical: true, current_group_id: 'G', generations: [] }, actions: ['edit', 'stop', 'remove', 'recreate', 'transition'].map(action => ({ action: action as WriteGroupLifecycle['actions'][number]['action'], allowed: action === 'edit', reason_code: 'successor_reserved', next_action: 'review_lifecycle', requires_review: false })) };
  const view = render(<I18nextProvider i18n={i18n}><SavedRecordingActions lifecycle={lifecycle} scope={scope} disabled={false} onAction={onAction} {...{ compact: true }} /></I18nextProvider>);
  expect(screen.getAllByTestId(/^saved-recording-reason-/)).toHaveLength(1);
  const reason = screen.getByTestId('saved-recording-reason-stop'); expect(reason).toBeVisible();
  for (const id of ['group-disable', 'group-delete', 'saved-recording-recreate', 'saved-recording-transition']) {
    const button = screen.getByTestId(id); expect(button).toBeDisabled(); expect(button).toHaveAttribute('aria-describedby', reason.id);
    expect(button.className).not.toContain('opacity-40'); fireEvent.click(button);
  }
  expect(onAction).not.toHaveBeenCalled(); fireEvent.click(screen.getByTestId('saved-recording-edit'));
  expect(onAction).toHaveBeenCalledExactlyOnceWith({ ...scope, action: 'edit' });
  view.rerender(<I18nextProvider i18n={i18n}><SavedRecordingActions lifecycle={{ ...lifecycle, workspace_revision: 'FOREIGN' }} scope={scope} disabled={false} onAction={onAction} {...{ compact: true }} /></I18nextProvider>);
  within(screen.getByTestId('saved-recording-actions')).getAllByRole('button').forEach(button => expect(button).toBeDisabled());
});
it('keeps the original guided readiness text when compact is absent', () => {
  render(<I18nextProvider i18n={i18n}><GroupReadinessPanel readiness={readiness()} loading={false} failed={false} dirty={false} onRetry={vi.fn()} /></I18nextProvider>);
  expect(screen.getByText(known, { exact: false })).toBeVisible();
});
