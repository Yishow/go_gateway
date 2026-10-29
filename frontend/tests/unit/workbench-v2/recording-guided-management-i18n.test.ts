import { expect, it } from 'vitest';
import en from '@/i18n/locales/en/workbench-v2-recording-guided.json';
import zh from '@/i18n/locales/zh-TW/workbench-v2-recording-guided.json';
it('labels saved management and Advanced as status rather than workspace activation in both languages', () => {
  expect(en).toHaveProperty('management.advancedTitle', 'Advanced settings and workspace status');
  expect(zh).toHaveProperty('management.advancedTitle', '進階設定與工作區狀態');
  expect(en).toHaveProperty('management.workspaceStatus'); expect(zh).toHaveProperty('management.workspaceStatus');
  for (const key of ['savedNameFallback', 'tableUnavailable', 'openDetails', 'actionsUnavailable']) {
    expect(en).toHaveProperty(`management.${key}`); expect(zh).toHaveProperty(`management.${key}`);
  }
});
