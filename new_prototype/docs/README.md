# Datalink Workbench - Page Specs

每個 md 檔都是一份**完整的 prompt**，描述對應頁面的版面、互動與資料流。可以直接餵給 AI 重生對應頁面，或拿來做設計 review。

## 目錄

| 文件 | 對應檔案 | 主題 |
|---|---|---|
| [_overview.md](./_overview.md) | (全) | 架構、state shape、版面、檔案結構 |
| [step1-device.md](./step1-device.md) | `step1-device.jsx` | Step 1 · 新增裝置 (多裝置 tab) |
| [step2-rule.md](./step2-rule.md) | `step2-rule.jsx` | Step 2 · 接入規則 (多規則 + 記憶體網格 + Share 位址) |
| [step3-mapping.md](./step3-mapping.md) | `step3-mapping.jsx` | Step 3 · 點位映射 (Tag + 個別 scale/offset + 轉換管線) |
| [step4-database.md](./step4-database.md) | `step4-database.jsx` | Step 4 · 儲存資料庫 (Connector + 欄位映射 + Commit) |
| [settings.md](./settings.md) | `settings.jsx` | 設定頁 (連接器池 / 時序 / 排程 / Share / 介面) |

## 文件結構慣例

每份頁面文件包含：

1. **目的** — 解決什麼問題、對應哪些後端 API
2. **完整 Prompt** — 一段足以重建該頁的描述 (可直接餵 AI)
3. **State 介面** — props、local state、衍生計算
4. **互動規範** — 表格式列出每個行為的觸發與結果
5. **視覺重點** — 配色、間距、特殊處理
6. **驗收標準** — checkbox 列表的可驗證項目

## 建議閱讀順序

1. 先看 `_overview.md` 了解全局
2. 依步驟順序看 step1 → step4
3. 最後看 `settings.md`
