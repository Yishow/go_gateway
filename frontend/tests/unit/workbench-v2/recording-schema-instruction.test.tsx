import { cleanup, render, screen, within } from '@testing-library/react';
import { afterEach, describe, expect, it, vi } from 'vitest';
import { I18nextProvider } from 'react-i18next';
import i18n from '@/i18n/config';
import { RuntimeSetupContextPanel } from '@/features/datalink/runtime-dashboard/RuntimeSetupContextPanel';
import { GroupReadinessPanel } from '@/features/datalink/workbench-v2/steps/step4/writeGroup/GroupReadinessPanel';
import type { WriteGroupReadiness } from '@/types/studioV2WriteGroup';
import type { StudioV2RuntimeSetupContext } from '@/types/studioV2RuntimeContext';

const message = 'preview and explicitly confirm the saved managed table before applying';
const issue = { code: 'schema-confirmation-required', severity: 'blocking' as const, step: 'Step 4' as const, scope: 'group-1', message };
const readiness: WriteGroupReadiness = { workspace_id: 'ws-1', workspace_revision: 'wr-1', group_id: 'group-1', group_revision: 'gr-1',
  applied_revision: '', config_ready: true, schema_ready: false, ready: false, issues: [issue] };
const setup: StudioV2RuntimeSetupContext = { evidence_state: 'saved-only', readiness_summary: { ready: false, blocking_count: 1, warning_count: 0, issues: [issue] },
  source_rules: [], mappings: [], database_config: null, database_targets: [] };

afterEach(async () => { cleanup(); await i18n.changeLanguage('zh-TW'); });

describe('schema readiness instruction', () => {
  it.each([
    ['zh-TW', '尚未確認並套用資料表結構；開始記錄前請先檢查並確認。', '診斷'],
    ['en', 'The table structure has not been confirmed and applied. Check and confirm it before starting recording.', 'Diagnostics'],
  ])('instructs the operator to check and confirm before starting in both surfaces (%s)', async (language, instruction, diagnosticLabel) => {
    await i18n.changeLanguage(language);
    render(<I18nextProvider i18n={i18n}>
      <GroupReadinessPanel compact readiness={readiness} loading={false} failed={false} dirty={false} onRetry={vi.fn()} />
      <RuntimeSetupContextPanel setupContext={setup} selectedDevice={null} />
    </I18nextProvider>);
    for (const container of [screen.getByTestId('group-readiness-issues'), screen.getByTestId('runtime-dashboard-setup-context')]) {
      const primary = within(container).getByText(instruction);
      expect(primary).toBeVisible();
      expect(primary).not.toHaveTextContent(issue.code);
      expect(primary).not.toHaveTextContent(message);
      const diagnostic = within(container).getByText(`${diagnosticLabel}：schema-confirmation-required：${instruction}（${message}）`);
      expect(diagnostic).toBeVisible();
      expect(within(container).getAllByText(/schema-confirmation-required|preview and explicitly confirm/)).toEqual([diagnostic]);
    }
  });
});
