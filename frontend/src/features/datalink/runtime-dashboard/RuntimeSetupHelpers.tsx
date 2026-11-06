import type { MouseEvent } from 'react';
import type { StudioV2WorkspaceReadinessIssue } from '../../../types/studioV2WorkspaceReadiness';

export function buildStudioV2FocusTarget(issue: StudioV2WorkspaceReadinessIssue | null): string {
  if (!issue) {
    return '/studio/v2';
  }

  const params = new URLSearchParams({
    step: String(readinessIssueOwnerStepNumber(issue)),
    focus: 'readiness',
    issue: issue.code,
  });

  return `/studio/v2?${params.toString()}`;
}

export function readinessIssueOwnerStepNumber(issue: StudioV2WorkspaceReadinessIssue): 1 | 2 | 3 | 4 {
  switch (issue.code) {
    case 'tag-missing':
    case 'mapping-missing':
      return 3;
    case 'database-target-missing':
      return 4;
    case 'device-connect-required':
    case 'device-probe-required':
      return 1;
    default:
      return readinessStepNumber(issue.step);
  }
}

export function readinessStepLabel(issue: StudioV2WorkspaceReadinessIssue): string {
  return `Step ${readinessIssueOwnerStepNumber(issue)}`;
}

export function fixInStudioLabel(
  t: (key: string, options?: Record<string, unknown>) => string,
  stepLabel: string,
): string {
  const fallback = `Fix in ${stepLabel}`;
  const translated = t('setup.fixInStudio', {
    step: stepLabel,
    defaultValue: fallback,
  });
  return translated === 'setup.fixInStudio' ? fallback : translated;
}

export function handleStudioTargetClick(
  navigateTo: ((target: string) => void) | undefined,
  target: string,
): ((event: MouseEvent<HTMLAnchorElement>) => void) | undefined {
  if (!navigateTo) {
    return undefined;
  }

  return (event) => {
    event.preventDefault();
    navigateTo(target);
  };
}

export function readinessStepNumber(step: StudioV2WorkspaceReadinessIssue['step']): 1 | 2 | 3 | 4 {
  switch (step) {
    case 'Step 2':
      return 2;
    case 'Step 3':
      return 3;
    case 'Step 4':
      return 4;
    case 'Step 1':
    default:
      return 1;
  }
}

export function formatReadinessStatus(
  ready: boolean | undefined,
  blockingCount: number,
  t: (key: string, fallback: string) => string,
): string {
  if (ready) {
    return t('setup.readiness.ready', 'Ready');
  }
  return `${blockingCount} ${t('setup.readiness.blockerUnit', 'blockers')}`;
}

export function RuntimeSetupFact({ label, value }: { label: string; value: string }) {
  return (
    <div>
      <dt className="text-[11px] uppercase tracking-wider text-slate-500">{label}</dt>
      <dd className="mt-0.5 break-words font-medium text-slate-200">{value}</dd>
    </div>
  );
}

export function readinessActionFor(
  issue: StudioV2WorkspaceReadinessIssue,
  t: (key: string, fallback: string) => string,
): string {
  switch (issue.code) {
    case 'point-missing':
      return t('setup.actions.pointMissing', 'Fix in Step 2: re-save this source rule to rebuild the missing derived point.');
    case 'database-target-missing':
      return t('setup.actions.databaseTargetMissing', 'Fix in Step 4: create a database target for this derived point.');
    case 'tag-missing':
      return t('setup.actions.tagMissing', 'Fix in Step 3: save the tag generated for this derived point.');
    case 'mapping-missing':
      return t('setup.actions.mappingMissing', 'Fix in Step 3: save the point-to-tag mapping.');
    case 'device-connect-required':
    case 'device-probe-required':
      return t('setup.actions.deviceProbeRequired', 'Fix in Step 1: pass connection/probe readiness before activation.');
    default:
      return t('setup.actions.default', 'Open Studio V2 and resolve the owning step shown above.');
  }
}

const KNOWN_READINESS_ISSUE_DEFAULTS: Readonly<Record<string, string>> = {
  'schema-confirmation-required': '尚未確認並套用資料表結構；開始記錄前請先檢查並確認。',
  'write-group-apply-required': '記錄設定已保存，尚未套用至資料庫與執行階段。',
  'database-target-missing': '尚未建立此點位的資料庫記錄目標。',
  'source-revision-stale': '來源點位或規則已更新，請重新保存此記錄設定。',
  'schema-check-blocked': '請先完成上方阻擋問題，才能檢查資料表結構。',
};

const UNKNOWN_READINESS_ISSUE_DEFAULT = '需要處理的設定問題';

export function describeReadinessIssue(
  issue: { code: string; message: string },
  tConstraints: (key: string, options: { defaultValue: string }) => string,
): { known: boolean; primary: string; diagnostic: string } {
  const trimmedMessage = issue.message.trim();
  const knownDefault = Object.hasOwn(KNOWN_READINESS_ISSUE_DEFAULTS, issue.code)
    ? KNOWN_READINESS_ISSUE_DEFAULTS[issue.code] : undefined;
  if (knownDefault) {
    const key = `readinessIssue.${issue.code}`;
    const translated = tConstraints(key, { defaultValue: knownDefault });
    const primary = translated && translated !== key ? translated : knownDefault;
    return {
      known: true,
      primary,
      diagnostic: trimmedMessage ? `${issue.code}：${primary}（${trimmedMessage}）` : `${issue.code}：${primary}`,
    };
  }
  const unknownKey = 'readinessIssue.unknown';
  const translatedUnknown = tConstraints(unknownKey, { defaultValue: UNKNOWN_READINESS_ISSUE_DEFAULT });
  const primary =
    translatedUnknown && translatedUnknown !== unknownKey
      ? translatedUnknown
      : UNKNOWN_READINESS_ISSUE_DEFAULT;
  return {
    known: false,
    primary,
    diagnostic: trimmedMessage ? `${issue.code}：${primary}（${trimmedMessage}）` : `${issue.code}：${primary}`,
  };
}
