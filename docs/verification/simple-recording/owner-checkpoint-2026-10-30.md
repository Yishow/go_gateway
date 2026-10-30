# 固定原型與正式實作確認

此清單記錄 adopt-simple-recording-entry 1.1 的實作前決策，不是最終產品驗收。

## 固定版本與決定

- Site v9：<https://gateway-simple-recording-preview.yishow-chen.chatgpt.site>。
- 原型 commit：`c7698bf137dc0d966e2be19c0ba6b06110da02d4`。
- HTML SHA256：`48adf85eb9820fd092523de0f6e9e510cee132e41c60a424f3b40e4d29e8af9b`。
- 版本與需求對照見 [整合需求覆核](requirements-audit-2026-10-30.md)。該文件保留發布當時尚未取得實作確認的歷史狀態。
- 2026-10-30 18:44 UTC 的 owner 訊息及主責判斷完整記於 [正式實作決定](implementation-decision-2026-10-30.md)。Owner 表示「若不用議你就完整驗到完成」；主責判斷不需繼續討論布局，開始既有七案的正式實作。
- 這是明確指示及已記錄的判斷，不以沉默、舊版 v4 或過往 browser waiver 代替確認。沒有宣稱 owner 逐頁操作過正式產品。

## 已提供的情境與保留的驗收

| 項目 | 固定版本／規格對照 | 正式產品必要驗收 |
|---|---|---|
| 首次建立與已有來源 | 設備、點位、記錄、目的地、檢查與明確開始 | 真服務接線；loading/error/unknown 不冒充空白；零隱性啟動 |
| 選點與來源設定 | 唯一加法選點入口、同表搜尋／設備／類型篩選、獨立加入計數、設備／範圍分步保存 | 取消保留既有與非法輸入；已存來源不回滾；沒有第一台設備默選 |
| 協議及參數 | 原型與 editor 的既有協議保真矩陣、進階入口、MQTT disabled | 逐欄 round-trip 與原 validator；模擬器證據不冒充真 PLC |
| 視覺與操作 | DESIGN.md 的 system/light/dark、一般 75vw × 75dvh modal、手機近滿版、小確認自動高度 | 1440/768/390、短視窗、鍵盤、焦點返回、對比與 reduced motion 的本批 browser 證據 |
| 記錄計算 | 每 member latest/mean/delta、各自 UTC epoch 對齊週期、寬表 not_due=NULL／長表只到期列 | 真可丟棄 SQL、窗口與 coverage、baseline、closure、重啟與重播；不補 0 或沿用增加量 |
| 已保存記錄與異常 | 草稿／正常／異常、修改／停恢／複製／移除／前向表格變更 | server availability、scope/revision、安全確認、unknown 先查原 operation |
| 深連結與退回 | 同 Studio V2、既有 locator/runtime/進階交接 | 明確意圖先於首頁；invalid 不降級；預設入口切換與資料設定分離 |

手機與小確認框尺寸例外已在 DESIGN.md 明確定義；不新增未選設計候選。實作期發現的資料正確性、作用域或使用流程缺陷必須修復並重測，不以「上線後再修」免除。

## 不得由本確認推論的結果

1. 原型的 jsdom／fixture 通過不是原生 browser、視覺、SQL 或協議伺服器通過。
2. Owner 18:52 UTC 再次要求 browser 發現問題須回修並重測；此要求沒有豁免。
3. 本清單解除開始正式實作的待決狀態，不開啟正式預設首頁。必要產品驗收尚未通過時維持 opt-in／原預設入口。
4. 入口回退只切換呈現，不撤銷已保存設定、不刪歷史、不重送未決命令。
5. 分段 commit／同步 main 已獲授權，但不代表部署、真 PLC 或正式 DB 操作授權。做的期間只跑相關測試，七案完成才統一總測。

沒有新的布局決策需在實作前再次詢問；仍需把上述必要驗收逐項完成並留證據，不能以此清單關閉未執行的測試。
