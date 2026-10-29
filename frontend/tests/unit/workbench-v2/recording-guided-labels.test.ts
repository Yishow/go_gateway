import { createElement, type ReactNode } from 'react';
import { createInstance } from 'i18next';
import { I18nextProvider, initReactI18next } from 'react-i18next';
import { cleanup, renderHook } from '@testing-library/react';
import { afterEach, describe, expect, it } from 'vitest';
import { recordingGuidedLabels } from '@/features/datalink/workbench-v2/steps/step4/recordingGuidedLabels';
import { useRecordingNextActionLabels } from '@/features/datalink/workbench-v2/steps/step4/recordingNextActionLabels';
import type { RecordingNextActionLabels } from '@/features/datalink/workbench-v2/steps/step4/RecordingNextAction';
import en from '@/i18n/locales/en/workbench-v2-recording-guided.json';
import zh from '@/i18n/locales/zh-TW/workbench-v2-recording-guided.json';

function lookup(resource: unknown, key: string): unknown {
  return key.split('.').reduce<unknown>((node, part) => typeof node === 'object' && node !== null && Object.hasOwn(node, part)
    ? (node as Record<string, unknown>)[part] : undefined, resource);
}
function keys(resource: unknown, prefix = ''): string[] {
  return resource && typeof resource === 'object' ? Object.entries(resource).flatMap(([key, value]) =>
    keys(value, prefix ? `${prefix}.${key}` : key)) : [prefix];
}
function translatedValue(resource: unknown, key: string): string {
  const value = lookup(resource, key);
  expect(typeof value).toBe('string'); expect(value).not.toBe(key); expect(String(value).trim()).not.toBe('');
  return String(value);
}
// These sections have different production owners; they are explicitly covered below, not discarded.
const separateSections = new Set(['next', 'management', 'completion']);
const preparationKeys = ['fixControls', 'startTitle', 'startDescription', 'deliverySeparate', 'confirmStart', 'cancel', 'repair', 'preparedOnly'];
const sectionKeys = {
  // Actual saved-card wording also has recording-guided-management-i18n.test.ts coverage.
  management: ['advancedTitle', 'workspaceStatus', 'savedNameFallback', 'tableUnavailable', 'openDetails', 'actionsUnavailable'],
  // Actual translated completion rendering belongs to recording-guided-completion-summary.test.tsx.
  completion: ['heading', 'activated', 'incompleteDevices', 'incompleteGroups', 'noneActivated', 'deviceFallback', 'groupFallback',
    'deliverySeparate', 'continue'],
};
afterEach(cleanup);
describe('registered guided label contract', () => {
  it.each([['English', en], ['Traditional Chinese', zh]] as const)('provides exactly the summary/detail helper labels in %s', (_name, resource) => {
    const requested: string[] = [];
    const labels = recordingGuidedLabels(key => {
      requested.push(key); return translatedValue(resource, key);
    });
    expect(labels.unsaved).not.toBe(labels.effective); expect(labels.saved).not.toBe(labels.effective);
    expect(new Set(requested)).toEqual(new Set(keys(resource).filter(key => !separateSections.has(key.split('.')[0]))));
  });
  it.each([['en', en], ['zh-TW', zh]] as const)('provides the actual typed next-action hook and preparation labels in %s', async (language, resource) => {
    const i18n = createInstance();
    await i18n.use(initReactI18next).init({ lng: language, fallbackLng: false,
      resources: { [language]: { 'workbench-v2-recording-guided': resource } } });
    const expected: RecordingNextActionLabels = { heading: resource.next.heading, scopes: resource.next.scopes,
      actions: resource.next.actions, reasons: resource.next.reasons };
    const { result } = renderHook(useRecordingNextActionLabels, {
      wrapper: ({ children }: { children: ReactNode }) => createElement(I18nextProvider, { i18n, children }),
    });
    expect(result.current).toEqual(expected);
    expect(keys(resource.next).sort()).toEqual([...keys(expected), ...preparationKeys].sort());
    keys(resource.next).forEach(key => translatedValue(resource, `next.${key}`));
  });
  it.each([['English', en], ['Traditional Chinese', zh]] as const)('accounts for every separately owned management/completion label in %s', (_name, resource) => {
    for (const [section, expected] of Object.entries(sectionKeys)) {
      expect(keys(lookup(resource, section)).sort()).toEqual([...expected].sort());
      expected.forEach(key => translatedValue(resource, `${section}.${key}`));
    }
  });
  it('keeps both locale leaf sets identical', () => expect(keys(en).sort()).toEqual(keys(zh).sort()));
});
