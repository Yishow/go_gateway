import { describe, expect, it } from 'vitest';
import i18n from '@/i18n/config';
import { describeReadinessIssue } from '@/features/datalink/runtime-dashboard/RuntimeSetupHelpers';

describe('unknown readiness codes', () => {
  it.each([
    ['constructor', 'zh-TW', '需要處理的設定問題'],
    ['toString', 'zh-TW', '需要處理的設定問題'],
    ['__proto__', 'zh-TW', '需要處理的設定問題'],
    ['constructor', 'en', 'A setting needs attention'],
    ['toString', 'en', 'A setting needs attention'],
    ['__proto__', 'en', 'A setting needs attention'],
  ])('treats %s as an unknown code in %s', (code, language, generic) => {
    const translate = i18n.getFixedT(language, 'recording-editor-constraints');
    expect(describeReadinessIssue({ code, message: 'Server note' }, translate)).toEqual({
      known: false,
      primary: generic,
      diagnostic: `${code}：${generic}（Server note）`,
    });
    expect(describeReadinessIssue({ code, message: '  ' }, translate)).toEqual({
      known: false,
      primary: generic,
      diagnostic: `${code}：${generic}`,
    });
  });
});
