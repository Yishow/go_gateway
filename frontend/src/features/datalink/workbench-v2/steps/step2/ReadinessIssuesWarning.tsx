import * as React from 'react';
import { useTranslation } from 'react-i18next';
import type { RuleReadinessIssue, RuleReadinessReason } from '../../state/sourceRule';

interface ReadinessIssuesWarningProps {
  issues: RuleReadinessIssue[];
}

const missingReasons: RuleReadinessReason[] = ['unknown_device', 'deleted_device'];
const fallbackCopy: Record<RuleReadinessReason, { title: string; action: string }> = {
  unknown_device: {
    title: '找不到此規則的目前設備歸屬。',
    action: '請選擇目前設備後再繼續。',
  },
  deleted_device: {
    title: '此規則所屬設備已被刪除。',
    action: '請選擇替代設備或先還原設備後再繼續。',
  },
  invalid_address: {
    title: '位址不符合所選協議格式。',
    action: '請修正該規則的協議位址後再繼續。',
  },
};

function identityLabel(value: string, fallback: string): string {
  return value.trim() || fallback;
}

export const ReadinessIssuesWarning: React.FC<ReadinessIssuesWarningProps> = ({ issues }) => {
  const { t } = useTranslation('workbench-v2');
  if (issues.length === 0) {
    return null;
  }

  const groups: RuleReadinessReason[] = [...missingReasons, 'invalid_address'];
  return (
    <div
      className="rounded-lg border border-red-500/30 bg-red-500/10 px-3 py-2 text-xs text-red-200"
      role="alert"
      data-testid="readiness-issues-warning"
    >
      {groups.map((reason, index) => {
        const reasonIssues = issues.filter((issue) => issue.reason === reason);
        if (reasonIssues.length === 0) {
          return null;
        }

        const groupClass = index > 0 ? 'mt-3' : undefined;
        const warningTestId = reason === 'unknown_device'
          ? 'readiness-missing-device-warning'
          : reason === 'invalid_address'
          ? 'readiness-invalid-address-warning'
          : `readiness-${reason}-warning`;
        return (
          <div key={reason} className={groupClass} data-testid={warningTestId} data-readiness-reason={reason}>
            <div className="font-semibold">
              {t(`step2.summary.${reason}`, fallbackCopy[reason].title)}
            </div>
            <div className="mt-1">
              {t(`step2.summary.${reason}_action`, fallbackCopy[reason].action)}
            </div>
            <ul className="mt-1 list-disc pl-4">
              {reasonIssues.map((issue) => {
                const device = identityLabel(issue.deviceId, t('step2.summary.unknown_device_identity', 'unknown device'));
                const address = identityLabel(issue.startAddress, t('step2.summary.unknown_address', 'unknown address'));
                return (
                  <li
                    key={issue.ruleId}
                    data-testid={reason === 'invalid_address'
                      ? `readiness-issue-invalid-address-${issue.ruleId}`
                      : `readiness-issue-missing-device-${issue.ruleId}`}
                    data-readiness-reason={issue.reason}
                  >
                    <span data-testid={`readiness-issue-${issue.reason}-${issue.ruleId}`}>
                      {issue.rule.name} · {device} · {address}
                    </span>
                  </li>
                );
              })}
            </ul>
          </div>
        );
      })}
    </div>
  );
};

export default ReadinessIssuesWarning;
