# i18n 國際化說明

本專案已整合 i18n（國際化）功能，支援英文和繁體中文。

## 支援的語言

- **英文 (en)**: 預設語言
- **繁體中文 (zh-TW)**: 繁體中文

## 檔案結構

```
frontend/src/
├── i18n/
│   ├── config.ts                    # i18n 配置檔案
│   └── locales/
│       ├── en/
│       │   └── common.json          # 英文翻譯
│       └── zh-TW/
│           └── common.json          # 繁體中文翻譯
└── components/
    └── LanguageSwitcher.tsx         # 語言切換組件
```

## 使用方式

### 在組件中使用翻譯

```tsx
import { useTranslation } from 'react-i18next';

function MyComponent() {
  const { t } = useTranslation();
  
  return (
    <div>
      <h1>{t('nav.dashboard')}</h1>
      <button>{t('common.save')}</button>
    </div>
  );
}
```

### 語言切換

語言切換器已整合到 `DatalinkLayout` 的頂部導航欄中。使用者可以透過下拉選單切換語言。

語言設定會自動儲存到 `localStorage`，下次訪問時會自動載入上次選擇的語言。

### 新增翻譯

1. 在 `frontend/src/i18n/locales/en/common.json` 新增英文翻譯
2. 在 `frontend/src/i18n/locales/zh-TW/common.json` 新增繁體中文翻譯
3. 在組件中使用 `t('your.key.path')` 來引用翻譯

### 翻譯鍵值命名規範

- 使用命名空間組織翻譯鍵值（如 `nav.*`, `device.*`, `tag.*`）
- 使用小寫字母和點號分隔（如 `device.saveDevice`）
- 保持鍵值名稱簡潔且具描述性

## 自動語言偵測

系統會依以下順序自動偵測語言：

1. `localStorage` 中儲存的語言設定
2. 瀏覽器的語言設定
3. HTML 標籤的 `lang` 屬性
4. 預設語言（英文）

## 注意事項

- 所有使用者可見的文字都應該使用翻譯功能
- 避免在組件中硬編碼文字
- 新增新功能時，記得同時新增英文和繁體中文翻譯
