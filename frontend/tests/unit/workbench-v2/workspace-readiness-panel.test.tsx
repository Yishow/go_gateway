import { describe, expect, it } from 'vitest';
import type { TFunction } from 'i18next';
import {
  compactReadinessGuidance,
  groupWorkspaceReadinessIssues,
  readinessIssueLabel,
  readinessIssueMessage,
  summarizeWorkspaceReadiness,
} from '../../../src/features/datalink/workbench-v2/components/WorkspaceReadinessPanel';
import type { StudioV2WorkspaceReadinessSummary } from '../../../src/types/studioV2WorkspaceReadiness';

describe('workspace readiness grouping', () => {
  it('groups repeated point-level blockers into one operator-facing bucket but keeps the displayed blocker count truthful', () => {
    const summary: StudioV2WorkspaceReadinessSummary = {
      ready: false,
      blocking_count: 8,
      warning_count: 0,
      issues: Array.from({ length: 8 }, (_, index) => ({
        code: 'mapping-missing',
        severity: 'blocking' as const,
        step: 'Step 3' as const,
        scope: `point-${index + 1}`,
        message: 'derived point is missing its persisted mapping',
      })),
    };

    const groupedIssues = groupWorkspaceReadinessIssues(summary.issues);
    const groupedSummary = summarizeWorkspaceReadiness(summary);

    expect(groupedIssues).toHaveLength(1);
    expect(groupedIssues[0]).toMatchObject({
      count: 8,
      issue: {
        code: 'mapping-missing',
        severity: 'blocking',
        step: 'Step 3',
        message: 'derived point is missing its persisted mapping',
      },
    });
    expect(groupedSummary.blockingCount).toBe(8);
    expect(groupedSummary.warningCount).toBe(0);
  });

  it('keeps distinct issue codes as separate operator-facing blockers', () => {
    const summary: StudioV2WorkspaceReadinessSummary = {
      ready: false,
      blocking_count: 3,
      warning_count: 1,
      issues: [
        {
          code: 'mapping-missing',
          severity: 'blocking',
          step: 'Step 3',
          scope: 'point-1',
          message: 'derived point is missing its persisted mapping',
        },
        {
          code: 'mapping-missing',
          severity: 'blocking',
          step: 'Step 3',
          scope: 'point-2',
          message: 'derived point is missing its persisted mapping',
        },
        {
          code: 'database-target-missing',
          severity: 'blocking',
          step: 'Step 4',
          scope: 'point-2',
          message: 'derived point is missing its persisted database target',
        },
        {
          code: 'database-connector-unreachable',
          severity: 'warning',
          step: 'Step 4',
          scope: 'db-main',
          message: 'database connector requires attention',
        },
      ],
    };

    const groupedSummary = summarizeWorkspaceReadiness(summary);

    expect(groupedSummary.blockingCount).toBe(3);
    expect(groupedSummary.warningCount).toBe(1);
    expect(groupedSummary.issues).toHaveLength(3);
  });

  it('translates the grouped blocker label, message, and compact action guidance', () => {
    const issue = {
      code: 'mapping-missing',
      severity: 'blocking' as const,
      step: 'Step 3' as const,
      scope: 'point-1',
      message: 'derived point is missing its persisted mapping',
    };
    const t = ((key: string, fallbackOrOptions?: unknown) => {
      const translations: Record<string, string> = {
        'workspace_readiness.issue_labels.mapping-missing': '映射缺失',
        'workspace_readiness.issue_messages.mapping-missing': '衍生點位尚未綁定已保存的 Tag 映射。',
        'workspace_readiness.solution_step3': '到 Step 3 補齊點位與 Tag 映射，讓資料能投影到 runtime。',
      };
      if (key === 'workspace_readiness.compact_fix') {
        const options = fallbackOrOptions as { label: string; action: string };
        return `先修正「${options.label}」：${options.action}`;
      }
      return translations[key] ?? String(fallbackOrOptions ?? key);
    }) as unknown as TFunction<'workbench-v2'>;

    expect(readinessIssueLabel(issue, t)).toBe('映射缺失');
    expect(readinessIssueMessage(issue, t)).toBe('衍生點位尚未綁定已保存的 Tag 映射。');
    expect(
      compactReadinessGuidance({ issue, count: 1 }, t),
    ).toBe('先修正「映射缺失」：到 Step 3 補齊點位與 Tag 映射，讓資料能投影到 runtime。');
  });

  it('keeps row-group duplicate column blockers distinct in readiness guidance', () => {
    const issue = {
      code: 'database-row-group-column-conflict',
      severity: 'blocking' as const,
      step: 'Step 4' as const,
      scope: 'temperature_c',
      message: 'database column is reused across row groups',
    };
    const t = ((key: string, fallbackOrOptions?: unknown) => {
      const translations: Record<string, string> = {
        'workspace_readiness.issue_labels.database-row-group-column-conflict': 'Row group 欄位衝突',
        'workspace_readiness.issue_messages.database-row-group-column-conflict': '同一欄位被不同 row group 共用，請調整群組或欄位。',
        'workspace_readiness.solution_step4': '回到 Step 4 修正資料庫 row group 與欄位配置。',
      };
      if (key === 'workspace_readiness.compact_fix') {
        const options = fallbackOrOptions as { label: string; action: string };
        return `${options.label}: ${options.action}`;
      }
      return translations[key] ?? String(fallbackOrOptions ?? key);
    }) as unknown as TFunction<'workbench-v2'>;

    expect(readinessIssueLabel(issue, t)).toBe('Row group 欄位衝突');
    expect(readinessIssueMessage(issue, t)).toBe('同一欄位被不同 row group 共用，請調整群組或欄位。');
    expect(compactReadinessGuidance({ issue, count: 1 }, t)).toBe('Row group 欄位衝突: 回到 Step 4 修正資料庫 row group 與欄位配置。');
  });
});
