import { useTranslation } from 'react-i18next';

/**
 * 語言切換組件
 * 
 * 提供語言選擇下拉選單，支援：
 * - 英文 (en)
 * - 繁體中文 (zh-TW)
 * 
 * 語言設定會自動儲存到 localStorage
 */
export default function LanguageSwitcher() {
  const { i18n, t } = useTranslation();

  const languages = [
    { code: 'en', label: 'English' },
    { code: 'zh-TW', label: '繁體中文' },
  ];

  const handleLanguageChange = (langCode: string) => {
    i18n.changeLanguage(langCode);
  };

  const currentLanguage = i18n.language === 'zh' ? 'zh-TW' : i18n.language.startsWith('en') ? 'en' : i18n.language;

  return (
    <div className="relative">
      <select
        value={currentLanguage}
        onChange={(e) => handleLanguageChange(e.target.value)}
        className="bg-slate-800 border border-slate-700 rounded-lg px-3 py-2 text-sm text-slate-200 focus:outline-none focus:border-blue-500 focus-visible:ring-2 focus-visible:ring-blue-500 cursor-pointer appearance-none pr-8"
        aria-label={t('layout.languageSelector')}
      >
        {languages.map((lang) => (
          <option key={lang.code} value={lang.code}>
            {lang.label}
          </option>
        ))}
      </select>
      <div className="absolute right-2 top-1/2 -translate-y-1/2 pointer-events-none">
        <svg className="w-4 h-4 text-slate-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 9l-7 7-7-7" />
        </svg>
      </div>
    </div>
  );
}
