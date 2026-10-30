# 2026-10-30 正式實作決定與驗證界線

## Owner 最新決定

2026-10-30 18:44 UTC，owner 原話（訊息 `Sentinel_02c70478b3a081919f4f2fc5a4840c8f`）：

> 我覺得應該差不多，上線後再修正應該也行，還是你確為還需要再議？實話說。若不用議你就完整驗到完成，以下原則記在記憶：做的期間只跑相關測試，直到完成才跑全部。

主責回覆（`Sentinel_73821135cf9c81919fe255407b61e94b`）判斷：目前設計足以開始正式實作，不需繼續討論布局；依既有七案完成實作、測試、修復及已授權的 GitHub main 交付。已知 correctness 問題須在部署前修復，不能將「上線後再修正」當成略過已知缺陷的授權。

## 固定採用參考

- 私人預覽：<https://gateway-simple-recording-preview.yishow-chen.chatgpt.site>，Site v9。
- HTML SHA256：`48adf85eb9820fd092523de0f6e9e510cee132e41c60a424f3b40e4d29e8af9b`。
- 原型來源 commit：`c7698bf137dc0d966e2be19c0ba6b06110da02d4`（短碼 `c7698bf`）。
- 本決定記錄基於產品 checkout `011e8037`；原型來源 commit 不是產品實作證據。

本次確認解除的是「是否繼續布局討論／是否開始正式實作」的待決狀態。它不宣稱完成原生瀏覽器視覺、RWD、焦點、操作旅程、真實可丟棄 SQL write/readback/cleanup 或正式產品驗收，也不免除 adoption 的必要 gate；未滿足時仍維持 opt-in。

## 七案共用測試節奏

適用首頁、編輯器、生命週期復原、預設入口採用、逐點字序、逐 member 計算及協議感知合批七案。

1. 做的期間只跑直接受影響及相關使用端測試；必要 focused integration、production wiring、race、SQL 與 acceptance 仍按各 task 完成，並記錄 PASS／FAIL／NOT RUN。
2. 七案整批實作完成後、最終交付前，才統一執行涉及前後端的完整 test／lint／vet／build 與嵌入式交付檢查；不是每完成一案便重跑全套。
3. 全套後若再修改相關程式、依賴或設定，依影響重驗必要檢查，不拿舊結果代表最終版本。
4. 不將 jsdom、原型 fixture 或 CLI artifact 驗證當成瀏覽器、SQL、真 PLC、正式 DB、Windows 或現場部署驗收。

## 歷史與本次紀錄

[先前需求覆核](requirements-audit-2026-10-30.md)保留原樣，該時點 owner 尚未確認與 apply 尚未開始的描述是歷史事實，不是本決定後的現況。本文僅更新決策與測試節奏，不追認未執行的驗收，不完成任何產品 task。

本次 artifact ingest 只修正首頁／編輯器／生命週期／採用四份 proposal 的 Modified 路徑分組語法（逐一確認實檔後分別以反引號引用），並澄清字序案 4.3 的最終全套時點；不改功能語意、工作項數、依賴或 checkbox。
