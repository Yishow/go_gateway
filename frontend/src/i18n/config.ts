import i18n from 'i18next';
import { initReactI18next } from 'react-i18next';
import LanguageDetector from 'i18next-browser-languagedetector';
import enTranslations from './locales/en/common.json';
import enTagReviewTranslations from './locales/en/tag-review.json';
import zhTWTranslations from './locales/zh-TW/common.json';
import zhTWTagReviewTranslations from './locales/zh-TW/tag-review.json';
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
        translation: mergeTranslations(enTranslations, enTagReviewTranslations),
      },
      'zh-TW': {
        translation: mergeTranslations(zhTWTranslations, zhTWTagReviewTranslations),
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
