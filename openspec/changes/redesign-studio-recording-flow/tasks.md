## 1. UI 測試驅動與 Red
- [x] 1.1 先加入四步核心使用者流程、鍵盤操作、欄位錯誤關联與窄視窗測試。
- [x] 1.2 先加入多台電表與 MC 混合範圍的角色確認、逐項例外、範本重套不覆蓋不重複測試。
- [x] 1.3 先鎖定真實 metadata、無欄位不猜、切換連線清預覽、取消/逾時不重複寫的測試。

## 2. 綠燈四步與白話文案
- [x] 2.1 整理 shell 固定動作列、密度與進階折疊，保留 routing 與既有 readiness。
- [x] 2.2 完成設備清單與本次啟用選擇、混合位址 editor、原值換算预覽與用途選擇。
- [x] 2.3 串接範本 batch preview/apply、修正冲突定位與局部重試，不用 browser-only成功狀態。
- [x] 2.4 完成記錄方案、頻率／保存／容量摘要、儲存模式與目的地共用設定。
- [x] 2.5 重排連線→選/建表→結果預覽→試寫→啟動；移除 production sample schema 與 index fallback。
- [x] 2.6 統一 zh-TW/en 用詞、按鈕、錯誤與空白狀態，技術欄移入進階而不刪安全資訊。

## 3. 狀態與回歸
- [x] 3.1 整合 autosave 接續、雙分頁 CAS、草稿／applied版本、密碼保護與離頁提示。
- [x] 3.2 整合真實採集／本機保存／送達／查回狀態、局部啟動與過期預覽保護。
- [x] 3.3 保留 Local Modbus only 流程、/studio/runtime handoff、/test 主要互動與 unknown-route 行為。
- [x] 3.4 Refactor、frontend lint/test/build、Playwright 與截图，既有安全 regression 全部通過。
- [x] 3.5 以未用過產品者驗收不用 SQL 完成案例，記錄錯誤修復路徑；strict spec 驗證後更新證據。

