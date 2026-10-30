# 2026-10-30 簡易記錄整合需求覆核

## 結論與版本

今日需求已納入七個既有提案，不只五種狀態或外觀調整。原型覆核找到三個會遺失使用者意圖的問題，已修復並重測；另補逐點下次區間截止。**這是離線原型與規格的進展，不是正式產品完成。**

- 產品查核基準：`41ed7bda7c851d85e95f853fccd147fba7446d44`；本文件之外未修改產品實作。
- 原 v8 HTML SHA256：`767269f10ab3fd375310a8f3846d573579016e47086fa619530d8e1a09beb428`。
- 修復候選 HTML SHA256：`48adf85eb9820fd092523de0f6e9e510cee132e41c60a424f3b40e4d29e8af9b`。
- [同一私人預覽站](https://gateway-simple-recording-preview.yishow-chen.chatgpt.site)已發布 **Site v9**（2026-10-30 18:30:38 UTC，部署成功）；來源commit `c7698bf137dc0d966e2be19c0ba6b06110da02d4`。此commit是原型來源，不是產品實作commit。
- 最終原型runner：**26/26 suites exit 0**，執行前後HTML hash不變；獨立operator **11 journeys PASS**；boundary/intent **3 PASS**。這些集合有重疊，不相加為總測試數。
- 七案共 **88 個產品工作項未完成**，並非88個純後端項目。新版 owner 頁面確認仍 pending，未開始產品 apply、未切換預設入口。

## 需求對照

「已準備」表示規格或demo已具備，不表示runtime完成；demo「通過」僅指所述離線事件測試。所有原生瀏覽器視覺／RWD欄仍未驗證。

| 使用者需求 | 原型／規格證據 | 正式程式或缺口 | 使用者影響／目前判定 |
|---|---|---|---|
| 設備→點位→記錄→DB，只有總覽與編輯兩主畫面 | `DESIGN.md`；首頁／編輯器案；首次建立與已有來源journey | 新簡易入口尚未接正式服務 | 已準備；設備管理留在總覽，不另建第三套產品 |
| 完整生命週期，不只五個狀態 | demo涵蓋空白、草稿、loading/error、waiting、stale、unknown、修改、停恢、複製、移除、接續新表；lifecycle案 | 既有server availability／scope／revision仍為authority | 已準備；不可將local status判斷直接搬成正式操作權限 |
| 保存／開始／SQL receipt／readback／cleanup分開 | demo等待首筆與歷史成功分列；lifecycle recovery | `cmd/test_ui/service_wiring.go`已有group pipeline，不能據此證明新mixed完成 | 新整合待實作；手動成功fixture不是SQL證據 |
| 保留協議、進階參數與型別 | 編輯器design協議矩陣；demo連線分型、能力說明，MQTT禁用 | 原進階10型別、位址、縮放、skip、Mapping與排程不可有損 | 已準備；demo只模擬部分型別，不宣稱實機parity |
| 保留既有路由／深連結 | `frontend/src/App.tsx`保留V2/runtime/test/gateway | `services/studioEntityLocator.ts:32–36`仍拒絕`step=output&target=modbus`相容意圖；adoption已有修復任務 | routes存在但全相容未完成；不恢復已退役`/studio` |
| system/light/dark、強字型／tokens | `DESIGN.md`單一ThemeProvider契約；demo OS listener、偏好、tokens | server raw-key hydration、late response、theme-only保存待正式整合 | 已準備；不得複製demo成第二主題系統 |
| RWD／75%置中modal／手機近滿／小confirm | `DESIGN.md`尺寸分類；demo data-size form/confirm CSS | 真1440/768/390、短視窗、焦點trap／返回、縮放未驗 | 未驗證視覺；CSS或jsdom不等於好用 |
| 繁中簡潔、每參數ⓘ、密集常顯設定 | demo help registry、connection-grid、point-grid；編輯器help契約 | 正式parameter coverage與原生版面仍待驗 | 已準備；錯誤／風險／必填不可藏help |
| 同協議session預填、可清除、不帶secrets | demo whitelist、epoch/generation、submode、清除；新設備保存事件 | 正式成功回覆／workspace隔離待接 | demo已準備；清除後未保存名稱dirty缺陷已修 |
| 可搜尋／篩選／分頁100設備300記錄 | demo 25/50有界表、設備摘要4台、規模fixture | 未作production效能量測 | demo旅程通過；不是100設備runtime性能承諾 |
| 重用已設定DB表，長表共用／寬表owner | demo configured-table選擇與owner守護；`workspace/write_group_long_storage.go`scope派生表名 | 新mixed仍需fresh schema／type／nullable／ownership證明 | 已準備；不掃DB、不複製舊proof，不把同名當相容 |
| 每點ABCD/CDAB保留且真的解碼 | 字序案legacy/inherit/override、range identity與proof；demo逐點設定 | `sourcerule/service.go`仍重寫衍生format；`adapters/mc3e_read.go`用connection converter | demo／規格已準備；新runtime能力未完成 |
| 主頁只列已選，獨立加法picker | demo單入口、同表搜尋／設備／類型、locked既有點、pending/footer/Add0 | 正式加入需scope/revision重驗 | demo已準備；不改來源，不在取消時改整份record |
| 唯一來源建立流程，空設備可續建 | demo contextual setup、明確設備／範圍分步保存、新點不自動選 | 正式unknown／late-save settlement仍待接 | demo旅程通過；取消後一步不回滾已保存來源 |
| draft／非法文字／Back／Cancel／remove/undo/readd | demo會員意圖與picker分層；三項本批RED→GREEN | demo reload前保存限制有明示，正式durable draft待整合 | 本批缺陷已修；不可將session draft當reload可恢復 |
| 每record-member latest／mean／delta互相獨立 | member計算案；demo不同record使用同Point的journey | `workspace/write_group.go:59–88`無member calculation/interval | demo／規格已準備；新22/12/8 SQL同列未實作 |
| 5／15／60／600秒獨立UTC對齊 | member案；demo固定clock與逐member截止 | `write_group_profile.go:77–84`仍為group interval | demo已準備；真正獨立閉窗／restart/replay待實作 |
| wide未到期NULL/not_due；long只到期；首基準/gap/error誠實 | demo正常／首窗／必要與可選基準預覽；member缺值契約 | baseline與closure同tx、新shape/schema與SQL readback待實作 | 已準備；gap/reset目前文字說明，非操作或runtime驗收 |
| 窗口／coverage／next due與實際SQL時間不同 | demo新增每點下次截止，明示固定08:00:07台北；窗與coverage可見 | 真row與receipt commit時間配對待接 | demo補足；不是目前時鐘、PLC讀取或SQL提交保證 |
| 真合批／去重／公平／限額／retry／freshness／量測 | 合批案planner、request preview、fault/load矩陣 | `collector/scheduler_dispatch.go:140–145`仍逐點；`scheduler_polling.go:13–53`逐點Read | 僅提案；不能稱MC/FATEK random已接，不能宣稱3→1或加速 |
| 進階保真與owner採用gate | adoption明確intent優先與opt-in回退 | 新UI與完整runtime未接，新版owner未確認 | 保留gate；overnight工作授權不等於視覺確認 |

## 本輪修復與驗證邊界

| 原版問題 | 最小修復 | 修復候選證據 |
|---|---|---|
| 僅輸入非法member間隔後取消，沒有dirty確認便離開 | 將raw非法文字／錯誤納入draft dirty，不污染最後有效值 | native click/change；Cancel、Back、scenario、beforeunload、picker取消 |
| 保存草稿後重開，非法間隔文字消失 | session draft保存／還原非法input buffer | 保存→resume保原文字與error；有效值未被0覆蓋 |
| 清除連線預填後，未保存設備名稱被當baseline，關閉無提示 | 清除只重設連線欄位基準，保留名稱dirty | name/config→clear→取消；一個active source dialog |
| 沒有逐member下一到期呈現 | 使用相同固定demo as-of推算截止，非法間隔不估算 | 5/15/60/600秒同時對齊與非法抑制；不假稱SQL時間 |

原版失敗資料與修復後結果均保留於本批外部QA證據包：原v8事件套件10項為7 PASS／3 FAIL，修復候選獨立operator為11 PASS（含後補四種間隔同時呈現）。沒有將舊失敗刪去或改寫成通過。首次建立觀察為26次動作（14 click＋12 edits），採預設值輔助，非正式產品操作下限。

## 跨功能風險與下一步

1. **選表×混合算法**：長表可共用不代表新結果語意相容；required但未到期的wide欄仍須nullable。method／interval變動必須使舊consent失效，必要時新表接續，舊raw歷史保留。
2. **字序×合批×窗口**：合批可共享raw wire資料，不可共享一種decode或合併Point identity；每Point保原observed time／quality／revision，再由member窗口閉合，不用SQL提交時間重寫窗口。
3. **來源×記錄草稿**：record draft、pendingAdd、source buffer、submitted operation各自保存／取消；已送出的unknown只查原operation，不自動補償刪來源或重送。
4. **簡單×進階**：latest不強制counter政策；需delta才驗confirmed policy。未支持組合有原因與原進階入口，不另造V3／引擎／新資料權威。
5. **通過×完成**：下一正式階段依既有七案與owner gate進行；需fake transport到真可丟棄SQL readback、durable replay及獨立批次效能量測，不連真PLC／正式DB。

## 檢查來源與未執行項

- 七案Spectra analyze使用既有`spxa`：Coverage／Consistency／Gaps／Localization均無發現；Ambiguity僅Suggestion（首頁31、編輯器50、lifecycle16、adoption9、字序12、member17、合批13），主要為Example格式與may措辭，不是runtime判定。
- 本批原型事件測試為jsdom，包含native click與named-form properties模擬；不等於原生browser、visual、focus或layout證據。
- 原生瀏覽器、1440/768/390截圖、Windows、PLC、真實SQL與production性能：**NOT RUN**。使用者18:28 UTC明確要求真瀏覽器→修復→重測；此要求仍未完成，不以本文件或jsdom結果關閉。
- 本批產品變更僅本文件，未重跑無關產品test/lint/build；`git diff --check`與`make check-lines`通過（1個文件）。既有完整產品測試數不移作本批證據。
