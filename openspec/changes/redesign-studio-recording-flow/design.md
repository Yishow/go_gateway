## Context

沿用 Studio surface inventory 的主線：/studio/v2 設定，/studio/runtime 完成後觀察，/test 工程工具。
四步導航保留；不是增成十幾個頁面，也不是重建已退休 /studio。
範例、示意值與實測值必須明確區分。Step4 目前 sample-column fallback 是此次明確取代的行為。

## Goals / Non-Goals

目標：不懂資料庫的人能完成正確記錄，熟手仍可展開自訂表、SQL 預覽與進階策略。
不做：移除安全驗證、掩飾部分失敗、讓瀏覽器代算用量、未完成的報表占位頁。

## Decisions

1. 共用固定動作列：上一步、當前儲存狀態、下一個具名動作。每一階段一个主要按鈕。
   required/optional、欄位例子、單位與當前待修問題就近顯示；不只用顏色傳意。
   密度預設一般，進階收合；不讓側欄、摘要和多份重複表格一起擠壓主要內容。
   鍵盤可完成核心操作，focus/label/error association 正確，桌面窄視窗也可用。
2. Step1「連接設備」：設備名称、通訊方式與必要參數；預設參數標明來源，IP 例子不當已填值。
   可存未完成草稿；實際啟用仍需既有協議 probe。未測成功的設備顯示原因，不假稱已連線。
   多台設備可明確選本次要啟用哪些，未選者保留草稿，而非一台離線導致全部無法設定。
3. Step2「選擇資料」：選設備與範圍，顯示每個項目占哪些位址，支援 uniform/mixed。
   以清單勾選與批次工具為主，網格是進階輔助，不要求先學 Shift/Ctrl。
   逐項 format、倍率例外由明確 editor 編輯；超出 span／重疊即指出哪兩項衝突。
4. Step3「確認數值與用途」：顯示原值→換算→單位，再選「畫曲線」「計算用量」「記狀態」等。
   可同時開曲線與用量等多個輸出；rate 積分與 counter 差分不得混為同一項。
   範本自动帶入建議，首次與不同型號需逐項確認。保留名稱/內部 key 的不同角色；key 自动生成。
   批次套用先比對來源型號、角色、span、倍率、用途；顯示將新增／修改／跳過／需確認的項目。
   使用者自行改過的欄位不被範本靜默覆蓋；重套範本不生成重複 measurement/plan。
5. Step4「設定記錄」：先用自然語句問用途與頻率／保存期限，接著選目的地。
   顯示已存連線／新增連線；不用跳去 settings 建一遍又回來填一遍。
   區分讀取、明細、摘要、送出 batch、保存期限，畫面提供一句可理解摘要與估算容量。
   單純 Local Modbus 轉發模式跳過資料庫；資料庫與轉發啟用狀態獨立且不繞過各自 readiness。
6. Connection→真正選表／managed 建表→項目與欄位預覽→試寫→開始持續記錄。
   managed 使用穩定角色與身份生成欄位，不需要每項人工配對。
   自訂表僅以真實 metadata 提建議；名稱相似但單位／型別不明標「請確認」。
   保留原已確認 binding 需重新驗證合法性；沒有合適欄位就留空，不依欄位順序亂填。
7. 操作狀態：儲存中／已存草稿／有變更待套用／記錄中／部分失敗／待補送。
   autosave 失敗、重載、切頁、關閉、兩分頁修改與 revision conflict 有可行動處理；
   密碼不放 localStorage；未持久化內容離頁需警告。關閉瀏覽器不停止已啟用後端工作。
   執行中修改先預覽影响，確認後才 apply；畫面永久顯示目前 applied version 與待套用 version。
8. Final summary 清楚列設備、項目數、用途、頻率、資料表、估算／缺項與權限。
   啟動與 schema/test-write 使用有界等待、取消等待／重查結果與後端 identity；
   逾時不代表沒有副作用，不直接從頭重送。局部成功分項顯示，可只重試失敗項。
   成功需帶真正 record/time/device/value 與 written/verified 差異，不能只看 collector running。
9. 中英語系完整，主要 zh-TW 文案：接入規則→讀取資料設定、Tag→資料名稱／內部識別名稱、
   Row group→一起記錄的資料、DDL→建表指令、Runtime→執行狀態、upsert→只更新目前資料。
   技術錯誤藏在安全詳細資訊，正常畫面說明「哪個資料、出了什麼事、如何修復」。

## Risks / Trade-offs

文字簡化不能隱藏資料意義與精度；第一次確認稍多，後續可安全套範本。
接口未完成時可以獨立測試 UI，但禁止正式入口以模擬成功代替真 API。
此 change 明確修改原 sample defaults 與 fallback 規格，保留其他 schema／憑證安全修正。

## Migration Plan

既有使用者載入原設定不自動升級用途；可預覽後採用新方案。
保留 routing identity，替換現有工作台元件，不新增競爭的 studio v3。
回退 UI 時後端已啟用的計畫仍可查狀態並安全停止；不能遺失待送與 checkpoint。
實作若更新 inventory 文件，按 CLAUDE.md 另外記 changelog.sqlite；本次不修改該目錄。

## Open Questions

無阻擋問題。實際預設頻率與保存空間由 plan review 决定；模板只提供可辨識的建議。
