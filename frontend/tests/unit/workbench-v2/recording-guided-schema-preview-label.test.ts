import { expect, it } from 'vitest';
import en from '@/i18n/locales/en/workbench-v2.json';
import zh from '@/i18n/locales/zh-TW/workbench-v2.json';

it.each([
  ['English', en, 'Preview creates confirmation data for this operation but does not change tables. Managed structures are created only after explicit confirmation.'],
  ['Traditional Chinese', zh, '預覽會建立本次操作的確認資料，但不會變更資料表；明確確認後才會建立受管結構。'],
] as const)('distinguishes preview operation data from explicitly confirmed schema changes in %s', (_language, resource, expected) => {
  expect(resource.step4.group.schema.confirm_note).toBe(expected);
});
