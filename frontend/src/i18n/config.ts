import i18n from 'i18next';
import { initReactI18next } from 'react-i18next';
import LanguageDetector from 'i18next-browser-languagedetector';
import enTranslations from './locales/en/common.json';
import enDatabaseOutputTranslations from './locales/en/database-output.json';
import enTagReviewTranslations from './locales/en/tag-review.json';
import zhTWTranslations from './locales/zh-TW/common.json';
import zhTWDatabaseOutputTranslations from './locales/zh-TW/database-output.json';
import zhTWTagReviewTranslations from './locales/zh-TW/tag-review.json';
import enWorkbenchV2Translations from './locales/en/workbench-v2.json';
import enWorkbenchV2CircuitSQLServer from './locales/en/workbench-v2-circuit-sqlserver.json';
import enWorkbenchV2RecordingAuthoring from './locales/en/workbench-v2-recording-authoring.json';
import enWorkbenchV2RecordingGuided from './locales/en/workbench-v2-recording-guided.json';
import enRuntimeDashboardTranslations from './locales/en/runtime-dashboard.json';
import zhTWWorkbenchV2Translations from './locales/zh-TW/workbench-v2.json';
import zhTWWorkbenchV2CircuitSQLServer from './locales/zh-TW/workbench-v2-circuit-sqlserver.json';
import zhTWWorkbenchV2RecordingAuthoring from './locales/zh-TW/workbench-v2-recording-authoring.json';
import zhTWWorkbenchV2RecordingGuided from './locales/zh-TW/workbench-v2-recording-guided.json';
import zhTWRuntimeDashboardTranslations from './locales/zh-TW/runtime-dashboard.json';
import enMappingErrors from './locales/en/mapping-errors.json';
import zhTWMappingErrors from './locales/zh-TW/mapping-errors.json';
import enRuntimeLogs from './locales/en/runtime-logs.json';
import enStudioReview from './locales/en/studio-review.json';
import zhTWStudioReview from './locales/zh-TW/studio-review.json';
import zhTWRuntimeLogs from './locales/zh-TW/runtime-logs.json';
import enRecordingLifecycleDialog from './locales/en/recording-lifecycle-dialog.json';
import zhTWRecordingLifecycleDialog from './locales/zh-TW/recording-lifecycle-dialog.json';
import enRecordingLifecycle from './locales/en/recording-lifecycle.json';
import zhTWRecordingLifecycle from './locales/zh-TW/recording-lifecycle.json';
import { logger } from '../utils/logger';
import { mergeTranslations } from './mergeTranslations';

/**
 * i18n 國際化配置
 * 
 * 支援語言：
 * - en: 英文
 * - zh-TW: 繁體中文
 * 
 * 預設語言：英文
 * 自動偵測瀏覽器語言設定
 * 
 * 使用 initReactI18next 以確保與 React 19 的兼容性
 * 必須在 .init() 之前調用 .use(initReactI18next)
 */
i18n
  .use(LanguageDetector) // 自動偵測瀏覽器語言
  .use(initReactI18next) // 初始化 React i18next 插件，確保正確的 React hooks 集成
  .init({
    resources: {
      en: {
        translation: mergeTranslations(
          enTranslations,
          enTagReviewTranslations,
          enDatabaseOutputTranslations,
        ),
        'mapping-errors': enMappingErrors,
        'workbench-v2': mergeTranslations(
          enWorkbenchV2Translations,
          enWorkbenchV2CircuitSQLServer,
          enWorkbenchV2RecordingAuthoring,
        ),
        'workbench-v2-recording-guided': enWorkbenchV2RecordingGuided,
        'runtime-dashboard': enRuntimeDashboardTranslations,
        'runtime-logs': enRuntimeLogs,
        'studio-review': enStudioReview,
        'recording-lifecycle': mergeTranslations(enRecordingLifecycle, enRecordingLifecycleDialog),
      },
      'zh-TW': {
        translation: mergeTranslations(
          zhTWTranslations,
          zhTWTagReviewTranslations,
          zhTWDatabaseOutputTranslations,
        ),
        'mapping-errors': zhTWMappingErrors,
        'workbench-v2': mergeTranslations(
          zhTWWorkbenchV2Translations,
          zhTWWorkbenchV2CircuitSQLServer,
          zhTWWorkbenchV2RecordingAuthoring,
        ),
        'workbench-v2-recording-guided': zhTWWorkbenchV2RecordingGuided,
        'runtime-dashboard': zhTWRuntimeDashboardTranslations,
        'runtime-logs': zhTWRuntimeLogs,
        'studio-review': zhTWStudioReview,
        'recording-lifecycle': mergeTranslations(zhTWRecordingLifecycle, zhTWRecordingLifecycleDialog),
      },
    },
    fallbackLng: 'en', // 預設語言為英文
    debug: false, // 開發時可設為 true 以查看除錯訊息
    
    interpolation: {
      escapeValue: false, // React 已經會處理 XSS，不需要額外轉義
    },
    
    detection: {
      // 偵測順序：localStorage > navigator > htmlTag > fallback
      order: ['localStorage', 'navigator', 'htmlTag'],
      caches: ['localStorage'], // 將語言設定儲存到 localStorage
    },
    
    // React 19 兼容性設置
    react: {
      useSuspense: false, // 禁用 Suspense，避免 React 19 兼容性問題
    },
  })
  .catch((err) => {
    logger.error('i18n initialization failed:', err);
  });

export default i18n;
