# Editor 3.3：主題、參數說明與可達性驗收

## 範圍與來源

- 本文交接 `streamline-recording-task-editor` 3.3，不代表七個 change 已完成。
- 沿用 DESIGN.md、固定 Taste app-screens 參考及既有 ThemeProvider；沒有新增主題系統、表單資料權威或採集引擎。
- Theme 最終組裝清單為 88 paths；Help 為 33 paths。組裝清單不是每檔所有修改的唯一歸屬。
- Step4Database 的樣式屬本項；其 runtime evidence wrapper 兩處接線另屬 4.1 的重讀修復，保留獨立四檔證據，不歸作 theme 邏輯。
- Browser 後最後 diff check 僅發現 TopBar 第 85 行字串外的一個行尾空格；v9 僅移除 1 byte U+0020，保存 v8 Browser 原始 snapshot 與 before/after 證明。沒有為無語意變動重跑測試／build／Browser，最終 assembly 不冒稱原 Browser byte-identical。
- Help v3 manifest SHA256：`93e326b5699185d5226aceb8b49120d193d802fc2b14f0b95a93f1c348e5e9b8`。
- Theme assembled manifest SHA256：`518a67f4bfefae0359e9a75087cf9e6efd3385c595676bc5d76f7e75510a9aa4`。
- 原始命令、凍結內容、圖片與 FAIL 均保留於本次 Linux 驗收資料；本文件不以檔名代替實際 viewport、語系或主題證據。

## 已實作

- Workbench 外殼與實際可見設定表單採共用 light/dark semantic tokens，保留 Follow System；移除原強制深色背景及外部字型請求。
- 用實際 sticky header 高度設定文件 scroll-padding-top；頁面切換沿既有 committed view，沒有 pointer/focus handler 偷移第一下點擊。
- Settings SaveBar 依自身高度與 bottom inset 設 scroll-padding-bottom；ResizeObserver／resize 更新，卸載還原原值及 priority，不改上方 padding。
- 手機表單不整頁橫溢；新增連接器按鈕、名稱／badge 可合理排列。四個既有設定開關重用原 Toggle；Connector Row 開關亦有名稱、switch role、checked 狀態及 44px 命中區。
- 原 Share pending、interaction lock、microtask 防重入及 lifecycle 保護保留，沒有為視覺測試啟用 Share。
- 參數說明使用非 modal popover；不新增第二個 dialog，不提交或改寫表單。Field 保留 label 關聯，按鈕不嵌在會誤啟動欄位的 label 行為內。
- 說明涵蓋連線、協議／serial、範圍、每點字序、計算、記錄目的地與進階參數；各語系 81 keys。
- coverage 的 93 是 consumer/key mappings，**不是 93 個實際 DOM placements**；可見數量依分支／列數而變。原去重統計漏掉 DB 欄重複按鈕，已以實際 render 回歸移除。
- 語系仍是 Settings 草稿；只有合法 locale 保存 ACK 或首次已完成讀取的合法儲存值才套用既有 i18n。未 Save、失敗或卸載後舊 ACK 不冒稱已套用；背景 refetch 不重新清除草稿。

## 直接測試及獨立 review

- Theme：15 個 token/theme 案例與 29 個 scroll/component/caller 案例；初版 71 個 TS/TSX 檔案 lint 通過。
- Help：原 90 個不同 focused cases 的有效結果沿用；重複 DB help／visual viewport 寬度兩項各真 RED→GREEN，僅跑兩項及改檔 lint。
- QA34 第一次整合型別檢查發現兩個新 test typing 錯誤，原 FAIL 保留；只修完整 connector fixture 與 EventTarget dispatch，兩直接案例／lint 通過，沒有 any 或 ts-ignore。
- Locale／Settings 窄修：9 個新案例＋21 個直接既有 caller 回歸；後續 hover、RWD、開關語意、四開關與原兩 Share caller、footer padding 分別只跑其直接新增／受影響案例及 lint。
- 各凍結增量由獨立 reviewer 核對：主題、Help、locale、hover 對比、RWD、switch 與 footer cleanup；最後兩檔 footer packet SHA256：`5d6916fad11a4e82d7e7ba21458aedfc238a3ce1c09ba9238d4e577ec87b125e`。
- 原資料權威、確認、安全與計算測試未因純樣式改動重跑。沒有為本項跑 Go race 或全套測試。

## 真實 Browser 結果

工具為官方 chrome-devtools-mcp；僅本機隔離 fixture、HSL 與 Gateway。沒有把 jsdom 或 tool success 當實際像素驗收。

- QA34：Light／Dark 的 1440、768、390 取樣可讀；System 隨官方 colorScheme 切換。部分 emulate 會重設未帶 viewport，報告已按實際幾何更正圖片檔名的侷限。
- QA34：說明可用 Enter／Space 開啟、Escape 只關 note 並回焦點；原來源 dialog 及未送出內容仍在。點真正未被 note 蓋住的 input 可外關且不偷焦點；底部說明可向上呈現。
- QA33 原正常 Shift+Tab 遮擋，在 QA34 以相同正常路徑閉合：Choose 可見，第一下 pointer click 及獨立 Enter 均進入結果配置，不需 PageUp。
- 手動捲動把焦點藏到 header 後，對隱形 UID 的點擊不視為正常可見點擊契約；原觀察保留，沒有加入危險的穿透或 pointer 捲動補丁。
- QA35：繁中／英文只有 Save ACK 後切換，重新載入及獨立新 browser context 的已存繁中 hydration 通過；雙語說明實測。Save All 依原契約寫多個設定，不冒稱只寫 locale；測試明確保持 auto_start=false、Share disabled。
- QA35：390 短確認依內容自高約 329px，並非硬撐 75% 高；Keep／Discard 沿原 owner。
- QA35→36：原 Settings 390 橫溢、逐字按鈕、Row 開關變形均保 FAIL→修復→真回歸；Light／Dark390 的 HTML scrollWidth/clientWidth 均 375，內容區均 351。
- QA36→37：其餘四個 Settings 開關由高橢圓／匿名按鈕改為正常橫軌、有名稱與 checked 的 switch；真圖片／AX 通過。
- QA37→38：自然 Breaker→Tab AutoStart→Tab Share 原被 SaveBar 遮擋；修後 Light／Dark390 與 1440 的 Tab／Shift+Tab、最後控制均完整可見，無 PageDown workaround。離開 Settings 後 bottom padding 還原 auto。
- QA34–38 最後導航與 console 持續檢查，沒有 React #185 或未處理 JS；QA35–38 最後 console 無 error/warn。
- QA34 有一次 Step4 mount 舊 runtime head 409，並非 Advanced click 觸發；後有新 head 200，但當時 body 未留。其 explicit retry authority 修復及 QA35 原始 groups→matching head 讀取證據另歸 4.1，不把 409 刪除或改稱刻意負向。

## 未執行／限制

- 真實觸控 tap／pinch：工具沒有相應事件操作；僅官方 touch emulation 的 layout／44px 命中區有證據，不能冒稱真 touch 事件通過。
- 真瀏覽器 zoom：一次官方 Control++ 嘗試前後幾何不變，列 NOT RUN；visualViewport offset／resize 寬度直接回歸已測，但不代替實機 zoom。
- 768 是相關頁面取樣，非所有欄位、每種組合的全笛卡兒矩陣。
- 沒有順帶驗所有其他次級頁、實體 PLC、serial 現場、Windows 或 PostgreSQL live。
- 本項沒有跑全部產品測試；最終全套仍留七個 change 實作完成後。

## 安全與後續

- 每輪新 fixture 保留原資料庫；最後相關輪次 target 仍不存在，來源為 draft，沒有 Group Apply／Start／DDL 或修改正式設備。
- 每輪停止自有服務並驗 source hashes、舊 DB hashes 與 backup integrity；原失敗、初始狀態及中間修復證據不覆蓋。
- QA38 最終 494 source hashes／40 舊 DB 不變，03:48:17 UTC 全停；RUN38-QA.txt SHA256：`b10380a2886458aac824e5716b7636bdfde784a62b6bae20aa67612d2acde13c`。
- 3.3 的 Browser 項逐項列 PASS／NOT RUN；仍需 4.1 批次驗證整合與 4.2 獨立權威／草稿／確認交接審查，之後才開 lifecycle。
