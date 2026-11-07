import { createInstance } from 'i18next';
import enHome from '@/i18n/locales/en/recording-home.json';
import zhHome from '@/i18n/locales/zh-TW/recording-home.json';
import enHomeStatus from '@/i18n/locales/en/recording-home-status.json';
import zhHomeStatus from '@/i18n/locales/zh-TW/recording-home-status.json';
import enHomeDelete from '@/i18n/locales/en/recording-home-delete.json';
import zhHomeDelete from '@/i18n/locales/zh-TW/recording-home-delete.json';
import enLifecycle from '@/i18n/locales/en/recording-lifecycle.json';
import zhLifecycle from '@/i18n/locales/zh-TW/recording-lifecycle.json';
import enDialog from '@/i18n/locales/en/recording-lifecycle-dialog.json';
import zhDialog from '@/i18n/locales/zh-TW/recording-lifecycle-dialog.json';
import { mergeTranslations } from '@/i18n/mergeTranslations';

// A disposable preview instance: no language detector, backend, cache or global plugin.
// Resource JSON is shared with actual components; no production config side effects.
export const previewI18n = createInstance();
export const previewTranslationsReady = previewI18n.init({
  lng: 'en', fallbackLng: 'en',
  resources: {
    en: { 'recording-home': mergeTranslations(enHome, enHomeStatus, enHomeDelete), 'recording-lifecycle': mergeTranslations(enLifecycle, enDialog) },
    'zh-TW': { 'recording-home': mergeTranslations(zhHome, zhHomeStatus, zhHomeDelete), 'recording-lifecycle': mergeTranslations(zhLifecycle, zhDialog) },
  },
  interpolation: { escapeValue: false },
  react: { useSuspense: false },
});
