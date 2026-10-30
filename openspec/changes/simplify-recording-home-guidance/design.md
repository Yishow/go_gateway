## Context

基準為 6f4a5fc。已讀 AppRoutes、DatalinkWorkbenchV2Page、RecordingNextAction、recordingWorkflowView。既有 workspace/query、autosave 與 scoped owner 為 authority；首頁是投影，不是新的控制平面。inventory 的部分缺口敘述早於現行程式，依上述 source 與既有契約判斷，不將舊記錄當此次測試結果。

## Design Source

來源為使用者已看過的 Gateway 簡易記錄 HTML demo v2（2026-10-30），不是正式後端或完整驗收。原樣快取於 .spectra/design-cache/simplify-recording-home-guidance/gateway-simple-ui-v2.dc.html；SHA256 d6987f1c8749db5d0ff9898544e3b2d505e98a99774b63d936a75daab2ffa306。cache 不提交。

| 項目 | 原樣擷取值 |
| --- | --- |
| 畫布／卡片／文字 | #f6f5f0 / #fffefa / #263b3b |
| 次要字／邊線／主色 | #708080 / #dce3dc / #27766f |
| 按鈕／卡片圓角 | 8px / 12px |
| 主內容寬度／內距 | max-width 1170px / 48px 32px 100px |
| 卡片內距／主要間距 | 24px / 24px |
| 圖示 | 原稿 Unicode 字符，沒有命名圖示組；正式版沿用專案 lucide-react 的語意圖示並有可讀標籤 |

上述數字是快取實測樣式，不是效能承諾。新情境、鍵盤與窄畫面可用性優先於逐像素照搬；demo 的暫存資料、位址簡化與假回報不得搬入正式服務。

## Goals / Non-Goals

**Goals:** 同一 Studio V2 內提供首頁與記錄編輯入口，首次與日常操作都有下一步；以真實證據與 scope 顯示目前結果。

**Non-Goals:** 不新增 V3、workspace、SQL 瀏覽器、引擎或路由 /studio；不改預設首頁、不取消原導引、runtime/logs/test 工具、不啟動設備或建表。兩個主畫面是本次設計條件，置中 modal與診斷頁不算新增主流程。

## Decisions

### 首頁優先序與同時存在的狀態

以純 resolver 接收 bootstrap 狀態、可證明 scope、來源清單、具名草稿、group/lifecycle/delivery 摘要。輸出 primaryAction、attentionItems、draftItems、recordingItems 與 availability。載入中／錯誤／scope 不可證明先獨立顯示 unavailable；可讀取時，未決 unknown operation 或異常優先，其次草稿，接著無設備→連接設備、有設備無可用點→設定點位、有點無記錄→建立記錄，最後正常摘要。優先序只決定提示，不刪除其他卡片或更改 persisted lifecycle。

### 同路由可選入口與既有識別保持

沿用 /studio/v2，使用明確 view=recordings 選擇新首頁；未指定維持舊導引。本批新增入口不改 root/wildcard 導向。先對原始 query 做完整 encoding/duplicate 檢查，再只接受單一 view=recordings 或 view=advanced presentation hint，剝離這個已驗證 hint 後交給原 parseStudioEntityLocator；不擴展 locator payload 權限。含 locator 意圖但有未知、重複或 malformed view/key 時維持 invalid，不降級一般首頁。無 locator 意圖的未知 view 才依原入口處理。合法 entity locator 與 runtime focus 的精確操作意圖優先於 view；只有完全無 locator 意圖的未知 view 才忽略並沿用既有入口；任何 malformed locator 仍維持 fail-closed。返回保存 query/hash 與焦點；不可由 URL 宣稱新的 mutation 權限。

### 結果證據而非狀態綠燈

最近 SQL committed 時間與目的地源自匹配 group/revision 的 delivery 證據；歷史成功和目前 freshness 分別展示。未收到結果是等待或未知；saved、start acknowledged、collecting、SQL committed 不等價。不承諾一般 delivery DTO 沒提供的 raw values；值展示僅在現有明確 readback 證據可用時標明來源。

## Implementation Contract

共同 apply 前置：四案規格完成審查、補齊 demo，並已取得 owner 對本次固定版本關鍵頁面／路徑的確認。此 gate 在第一案產品實作前執行，不等到第四案；demo確認不代表正式瀏覽器或實機驗收通過。

- 新 frontend/src/features/datalink/workbench-v2/recording-home/ 放投影與視圖；重用 useStudioV2WorkspaceQuery、既有 scoped groups 與 lifecycle/delivery hook，不另存第二份工作區狀態。輸入過期或 scope 不符不接受，較舊非同步回覆不覆蓋目前 scope。
- 點 primaryAction 只導航至對應 editor/修復/查狀態；真正 start、schema、停止等仍由後續 change 的明確 review/confirmation 命令執行。
- five-state fixtures 與混合情境均透過同一 resolver/render；不得各自 hard-code 五頁。fixtures 只進 demo/tests，不進 production bootstrap。
- 單一主要下一步之外，草稿、正常、異常記錄仍可到達；各項有名稱、原因與對應動作。設備區可折疊，異常不藏在進階內。
- RecordingHome 測試先行，包含空白／已有設備無點／有點無記錄／草稿／正常／異常／混合／loading／error／unknown，並鎖定導航零 mutation、舊 query/hash、locator、返回焦點及 stale response。
- 本批相關 unit/router tests、lint/typecheck、line gate；整功能完成才跑完整前端與嵌入式 build。瀏覽器能力目前受限，jsdom 不算視覺或實機驗收；本次 UI gate 不繼承歷史 waiver。
- 依賴：使用已完成 guide-step4-recording-workflow、expose-runtime-recording-rules、make-runtime-logs-actionable。交接下批 editor 與 lifecycle；本批不得把只導航的空 editor 宣稱完整功能。

## Risks / Trade-offs

- [混合狀態被單一標籤遮住] → priority 只控制主提示，全部具名項目仍可開啟。
- [舊入口、locator 或 runtime return 被搶走] → 明確意圖優先、route contract 回歸；尚未切預設。
- [簡化造成成功誤判] → receipt/revision/freshness 必須分開，unknown 提供查狀態而非自動重送。
- [原始 demo 太簡略] → 後續 demo 補生命週期及 owner 頁面 checkpoint；不照搬模擬引擎。

## 2026-10-30 需求更新（尚待版面選項確認）

根目錄 `DESIGN.md` 為唯一人讀設計規範；依 Taste app-screens 與 designtoken.md 的 primitive → semantic → component 分層，沿用既有 ThemeProvider 與字型，不另建主題系統。舊 demo 的淺色色值僅是歷史快取證據，不再是實作目標。

新簡易介面的彈出設定改為置中 modal，modal 本身寬度為目前 browser viewport 的 75%、高度為 75%（75vw × 75dvh），內容內捲；主頁仍正常 RWD，不是頁面縮為75%。已選手機近滿版與小型確認較小例外，尺寸分類依根 DESIGN.md；一般設定仍75% viewport。固定header/footer、contexthelp與同scope返回保留搜尋／篩選／捲動納入；選項8已選為dev/test專用樣式驗收頁，不放主選單，3/5/6/7未選，不新增tasks。文案簡潔、標籤保留意義，繁中／英文與數值／單位不能因版面省略。

共同 owner gate 必須用這次修訂後固定 demo 版本；過去 demo 檢視不視為確認。現處需求收集階段，未 apply、未做產品驗收；確定的協議保真與主題要求先納入，未選建議不擴 scope。

### 系統主題與單一語意 tokens

**Supersedes**: datalink-workbench-v2-shell / Design Tokens 落地：CSS variables + Tailwind theme.extend

取代本輪新 surfaces 的 dark-only／light-only限制，保留原 scope 隔離、字型與既有 brand 映射；不全站重寫。根 DESIGN.md 定義 raw primitive、light/dark semantic 與 component state，正式 CSS 單一值來源、TS僅adapter。舊淺色demo快取不刪，以歷史來源標示。

**Supersedes**: datalink-workbench-v2-settings / 不處理「主題切換」實作 light mode（v2 仍 dark-only；主題 select 只是存設定但 UI 不換）。

ThemeProvider preference/resolved分離；auto→system，raw server theme key存在才是已保存意圖，mapper預設dark不算。首次優先server有效偏好→新版local明確偏好→legacy light/dark→system；legacy值無法辨識手動或舊OS快取，保守保留且提供跟隨系統。使用者新選擇不被舊hydration覆蓋；system監聽OS change並cleanup，override不被OS改。server theme是global system_settings，非workspace；workspace切換不重設theme。theme-only操作只用updateKey('theme',value) adapter，不用會寫modbus_share/scheduler等整份設定的persistSettingsSnapshot；OS change不持久化。保存pending/error如實顯示，無storage/matchMedia安全fallback。所有新modal/portal/nativecontrol與新頁相同主題，不重設草稿／operation。

### 已選 modal 與返回引導契約

依 DESIGN.md 尺碼分類：一般popup75vw×75dvh；mobile<640px或short viewport<480px近滿版；小型短確認較小autoheight。第二主畫面editor維持RWD，不是popup。modal固定header/footer、body內捲；criticalrisk/validation/consent不能藏contexthelp。返回搜尋/篩選/捲動僅同scope、entity與view，revision變更重驗identity，不能復活舊selection/availability/receipt。dirtyguard、focus與pending操作語意不變。選項8已選為dev/test專用樣式驗收頁，不放主選單，未選候選不入實作；此需求更新不是owner對新版demo的視覺核准。

### 開發驗收樣式頁不成為第三主流程

已選8：正式僅dev/test build guarded入口，使用actualcomponents與同一tokens，展示light/dark、button/input/status/error/disabled/focus/modal等狀態。不放主選單，不新增operator主頁、不寫資料或啟用、不收credentials，不用fake production bootstrap。production不bundle開發頁；用現有Vite buildmode guard，無新服務。demo可專用非主選單link/hash，明示mock。preview只是輔助驗收，不替代實際page journeys或browser gate。

## v4 owner 未驗收

v4配色、密度與部分行為被使用者否決；既有mock測試不是UI確認。本輪先修demo及規格，修訂固定版本重新確認後才可apply；已以DESIGN.md Quiet mineral候選替換舊blue/slate；新版仍待owner視覺驗收，不當核准品牌。

### 大量設備與規則可定位而非無限卡片

**Supersedes**: simplify-recording-home-guidance / 首頁優先序與同時存在的狀態

只取代無界卡片呈現，保留原resolver狀態優先序與全部項目可達；設備最多4列摘要及同頁完整管理view mode，不增第三主頁。完整記錄／設備列表支援名稱、設備、資料表keyword及設備／狀態／目的地filter、穩定canonical-ID tie-break排序、25/50分頁。短摘要不塞分頁；attention仍明顯，不因篩選而假稱全部健康。空集合與查無結果不同，有清除篩選。

首頁使用既有workspace/device與useWriteGroupsQuery資料的client projection；目前write-groups是全量讀，不宣稱network bounded或server全域search。runtime既有useRuntimeSourcePage/query scope/list_revision/cursor只單設備來源rules，可沿用定位模式但不能冒充跨device索引。切scope重設query結果authority，revision變動重驗selection、clamp頁次；返回先filter/page再scroll/focus，較舊response不能復活刪除項目。

### Quiet mineral 取代被否決的 v4 配色

**Supersedes**: simplify-recording-home-guidance / 系統主題與單一語意 tokens

只替換該決策內未核准blue/slate視覺mapping，保留單一provider、theme-only persistence及token authority。DESIGN.md新版Quiet mineral完整primitive/semantic/state映射是v5候選，不宣稱owner已同意；light暖中性／dark中性charcoal、teal單主色，沒有navy。字型載入與CJKfallback沿用，字級24/18/14，panel8/control6/modal12px為設計條件。先靜態對比，再實際視覺斷點；過去v4離線PASS不是本版視覺證據。
