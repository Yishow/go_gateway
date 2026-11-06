import { cleanup, render, screen, within } from '@testing-library/react';
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';
import { I18nextProvider } from 'react-i18next';
import i18n from '@/i18n/config';
import { GroupReadinessPanel } from '@/features/datalink/workbench-v2/steps/step4/writeGroup/GroupReadinessPanel';
import { RuntimeSetupContextPanel } from '@/features/datalink/runtime-dashboard/RuntimeSetupContextPanel';
import { describeReadinessIssue } from '@/features/datalink/runtime-dashboard/RuntimeSetupHelpers';
import { GroupMemberTable } from '@/features/datalink/workbench-v2/steps/step4/writeGroup/GroupMemberTable';
import { RecordingPointPicker } from '@/features/datalink/workbench-v2/recording-editor/RecordingPointPicker';
import { createPicker } from '@/features/datalink/workbench-v2/recording-editor/recordingPickerState';
import { recordingPickerLabels } from '@/features/datalink/workbench-v2/recording-editor/recordingPickerLabels';
import { buildGroupCandidates } from '@/features/datalink/workbench-v2/state/writeGroup/candidates';
import { savedGroup, savedGroupState, TABLE_COLUMNS } from '../../fixtures/writeGroupState';
import type { WriteGroupReadiness } from '@/types/studioV2WriteGroup';
import type { StudioV2RuntimeSetupContext } from '@/types/studioV2RuntimeContext';

afterEach(cleanup);
beforeEach(async () => {
  await i18n.changeLanguage('zh-TW');
});

const KNOWN_CASES = [
  {
    code: 'schema-confirmation-required',
    zh: '尚未確認並套用資料表結構；開始記錄前請先檢查並確認。',
    en: 'preview and explicitly confirm the saved managed table before applying',
  },
  {
    code: 'write-group-apply-required',
    zh: '記錄設定已保存，尚未套用至資料庫與執行階段。',
    en: 'saved write group revision has not been applied to runtime',
  },
  {
    code: 'database-target-missing',
    zh: '尚未建立此點位的資料庫記錄目標。',
    en: 'derived point is missing its persisted database target',
  },
  {
    code: 'source-revision-stale',
    zh: '來源點位或規則已更新，請重新保存此記錄設定。',
    en: 'saved source or mapping revision is stale',
  },
  {
    code: 'schema-check-blocked',
    zh: '請先完成上方阻擋問題，才能檢查資料表結構。',
    en: 'destination schema check is blocked by configuration issues',
  },
] as const;

describe('Readiness issues are shown in plain language (Task 3.2)', () => {
  it.each(KNOWN_CASES)(
    'GroupReadinessPanel renders localized primary text and diagnostic code/zh/en for known code $code',
    ({ code, zh, en }) => {
      const readiness: WriteGroupReadiness = {
        workspace_id: 'ws-1',
        workspace_revision: 'w1',
        group_id: 'g-1',
        group_revision: 'gr-1',
        applied_revision: '',
        config_ready: false,
        schema_ready: false,
        ready: false,
        issues: [{ code, severity: 'blocking', step: 'Step 4', scope: 'g-1', message: en }],
      };

      render(
        <I18nextProvider i18n={i18n}>
          <GroupReadinessPanel
            compact
            readiness={readiness}
            loading={false}
            failed={false}
            dirty={false}
            onRetry={vi.fn()}
          />
        </I18nextProvider>,
      );

      const issuesList = screen.getByTestId('group-readiness-issues');
      expect(within(issuesList).getByText(zh)).toBeVisible();
      expect(
        within(issuesList).getByText(`${code}：${zh}（${en}）`, { exact: false }),
      ).toBeVisible();
      expect(issuesList.querySelectorAll('.recording-readiness-issue__code')).toHaveLength(1);
      expect(issuesList.querySelectorAll('.sr-only')).toHaveLength(0);
      expect(issuesList.textContent).not.toContain(`${code} · ${en}`);
    },
  );

  it.each(['', '   '])('omits empty parentheses in diagnostics when issue message is %j (F4)', (emptyMessage) => {
    const setupContext: StudioV2RuntimeSetupContext = {
      evidence_state: 'saved-only',
      readiness_summary: {
        ready: false,
        blocking_count: 1,
        warning_count: 0,
        issues: [
          {
            code: 'schema-confirmation-required',
            severity: 'blocking',
            step: 'Step 4',
            scope: 'scope-1',
            message: emptyMessage,
          },
        ],
      },
      source_rules: [],
      mappings: [],
      database_config: null,
      database_targets: [],
    };

    render(
      <I18nextProvider i18n={i18n}>
        <RuntimeSetupContextPanel setupContext={setupContext} selectedDevice={null} />
      </I18nextProvider>,
    );

    const panel = screen.getByTestId('runtime-dashboard-setup-context');
    expect(
      within(panel).getByText('診斷：schema-confirmation-required：尚未確認並套用資料表結構；開始記錄前請先檢查並確認。'),
    ).toBeVisible();
    expect(panel.textContent).not.toContain('（');
    expect(panel.textContent).not.toContain('（）');
    expect(panel.textContent).not.toContain('()');
  });

  it('GroupReadinessPanel renders generic fallback for unknown code and keeps code + raw message in diagnostics', () => {
    const readiness: WriteGroupReadiness = {
      workspace_id: 'ws-1',
      workspace_revision: 'w1',
      group_id: 'g-1',
      group_revision: 'gr-1',
      applied_revision: '',
      config_ready: false,
      schema_ready: false,
      ready: false,
      issues: [
        {
          code: 'custom-future-code',
          severity: 'blocking',
          step: 'Step 4',
          scope: 'g-1',
          message: 'raw English server explanation',
        },
      ],
    };

    render(
      <I18nextProvider i18n={i18n}>
        <GroupReadinessPanel
          compact
          readiness={readiness}
          loading={false}
          failed={false}
          dirty={false}
          onRetry={vi.fn()}
        />
      </I18nextProvider>,
    );

    const issuesList = screen.getByTestId('group-readiness-issues');
    expect(within(issuesList).getByText('需要處理的設定問題')).toBeVisible();
    expect(
      within(issuesList).getByText('custom-future-code：需要處理的設定問題（raw English server explanation）', {
        exact: false,
      }),
    ).toBeVisible();
  });

  it.each(KNOWN_CASES)(
    'RuntimeSetupContextPanel renders localized primary text and diagnostic code/zh/en for known code $code',
    ({ code, zh, en }) => {
      const setupContext: StudioV2RuntimeSetupContext = {
        evidence_state: 'saved-only',
        readiness_summary: {
          ready: false,
          blocking_count: 1,
          warning_count: 0,
          issues: [{ code, severity: 'blocking', step: 'Step 4', scope: 'scope-1', message: en }],
        },
        source_rules: [],
        mappings: [],
        database_config: null,
        database_targets: [],
      };

      render(
        <I18nextProvider i18n={i18n}>
          <RuntimeSetupContextPanel setupContext={setupContext} selectedDevice={null} />
        </I18nextProvider>,
      );

      const panel = screen.getByTestId('runtime-dashboard-setup-context');
      expect(within(panel).getByText(zh)).toBeVisible();
      expect(within(panel).queryByText(code, { exact: true })).not.toBeInTheDocument();
      expect(within(panel).getByText(`診斷：${code}：${zh}（${en}）`)).toBeVisible();
    },
  );

  it('RuntimeSetupContextPanel renders generic fallback for unknown code and keeps code + raw message in diagnostics', () => {
    const setupContext: StudioV2RuntimeSetupContext = {
      evidence_state: 'saved-only',
      readiness_summary: {
        ready: false,
        blocking_count: 1,
        warning_count: 0,
        issues: [
          {
            code: 'unknown-runtime-issue',
            severity: 'blocking',
            step: 'Step 4',
            scope: 'scope-1',
            message: 'unmapped backend message',
          },
        ],
      },
      source_rules: [],
      mappings: [],
      database_config: null,
      database_targets: [],
    };

    render(
      <I18nextProvider i18n={i18n}>
        <RuntimeSetupContextPanel setupContext={setupContext} selectedDevice={null} />
      </I18nextProvider>,
    );

    const panel = screen.getByTestId('runtime-dashboard-setup-context');
    expect(within(panel).getByText('需要處理的設定問題')).toBeVisible();
    expect(
      within(panel).getByText('診斷：unknown-runtime-issue：需要處理的設定問題（unmapped backend message）', {
        exact: false,
      }),
    ).toBeVisible();
  });

  it.each([
    ['zh-TW', '需要處理的設定問題'],
    ['en', 'A setting needs attention'],
  ])('keeps the localized fallback in diagnostics for an unknown code without a message (%s)', (language, fallback) => {
    const translate = i18n.getFixedT(language, 'recording-editor-constraints');
    expect(describeReadinessIssue({ code: 'future-code', message: '  ' }, translate)).toMatchObject({
      primary: fallback,
      diagnostic: `future-code：${fallback}`,
    });
  });

  it('sets autoComplete="off" on recording point search inputs', () => {
    const state = savedGroupState();
    const { candidates } = buildGroupCandidates(state);
    const group = savedGroup();

    const { unmount } = render(
      <I18nextProvider i18n={i18n}>
        <GroupMemberTable
          presentation="simple"
          managed
          candidates={candidates}
          excluded={[]}
          members={group.members.map((m) => ({
            key: `${m.point_id}:${m.tag_id}`,
            device_id: m.device_id,
            point_id: m.point_id,
            tag_id: m.tag_id,
            entity_key: m.entity_key ?? '',
            target_column: m.target_column,
            required: m.required,
          }))}
          columns={TABLE_COLUMNS}
          metadataStatus="exists"
          suggestions={{ assignments: {}, unmatched: [] }}
          issues={[]}
          showEntityKey={false}
          disabled={false}
          selectionControls={{ removeLabel: () => '移除', onRemove: vi.fn() }}
          onChange={vi.fn()}
        />
      </I18nextProvider>,
    );

    expect(screen.getByTestId('group-member-search')).toHaveAttribute('autoComplete', 'off');
    unmount();

    const pickerState = createPicker('scope-1');
    const tPicker = i18n.getFixedT('zh-TW', 'recording-editor');
    render(
      <I18nextProvider i18n={i18n}>
        <RecordingPointPicker
          state={pickerState}
          projection={{ scope: 'scope-1', candidates, members: [], evidence: {} }}
          devices={[{ id: 'dev-1', name: 'PLC 1' }]}
          status="loaded"
          labels={recordingPickerLabels(tPicker)}
          onChange={vi.fn()}
          onCancel={vi.fn()}
          onAdd={vi.fn()}
          onSource={vi.fn()}
          onRetry={vi.fn()}
        />
      </I18nextProvider>,
    );

    const pickerSearch = document.querySelector('input[data-picker-focus="search"]');
    expect(pickerSearch).not.toBeNull();
    expect(pickerSearch).toHaveAttribute('autoComplete', 'off');
  });
});
