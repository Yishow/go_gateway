# Studio V2 記錄生命週期交付驗收報告

更新日期：2026-10-28。範圍：`manage-recording-lifecycle-safely` task 4.3。

**現有工作樹的本地檢查已通過；42 個正式瀏覽器案例與所有截圖均未執行。** 使用者於 2026-10-28 15:27 UTC 明確批准本次略過實際 browser／截圖驗收後繼續交付。此批准只調整本次交付門檻，不把未執行項目改列 PASS，也不放寬產品保護。

**獨立發布候選已完成下列限定檢查。** 候選以 `606cb95ef5f3b8ef52da7ce195f874058faf9b53` 為基準，只納入本 change，兩個共用前端檔案排除後續 authoring 內容。下列原工作樹結果與候選結果分開呈現；本報告本身不代表 GitHub 發布或部署成功。

相關文件：[操作指南](studio-v2-recording-lifecycle.md)、[生命週期規格](../../openspec/changes/manage-recording-lifecycle-safely/specs/recording-lifecycle-management/spec.md)、[設計與本次驗收例外](../../openspec/changes/manage-recording-lifecycle-safely/design.md)、[任務狀態](../../openspec/changes/manage-recording-lifecycle-safely/tasks.md)。

## 獨立發布候選結果

- 前端完整 unit／integration：**303 檔／2,290 項 PASS**；不同於含後續修改的原工作樹 315／2,942。
- TypeScript、Vite build、正式 42 案例 discovery、明確涵蓋 E2E 的型別檢查 PASS；實際 browser 仍 NOT RUN。
- 完整 ESLint 指令兩次以 137／Killed 結束，原因未證實，原紀錄保留。改以相同設定、相同完整檔案集合分 14 批序列檢查，**829 檔全覆蓋、0 errors、2 個既有 console warnings**；沒有改規則或額外排除檔案。
- 後端 9 個相關 package 範圍：**8 個測試 package PASS，Swagger package 無測試**；一般測試共有 4,974 個通過事件、448 個 SKIP 事件（含子案例，不等於頂層測試數）。外部資料庫 fixtures 與真設備測試略過，不冒稱實機通過。
- 後端候選相依比較後，受後續排除內容影響的 API／主程式兩個 package 另在精確候選重跑，tests／vet PASS；其餘依相同來源與相依指紋重用本次結果，不重複累加測試數。
- 相關 package 的 go vet 與 282 個 Go 檔案格式／來源指紋核對通過；golangci-lint 不可用，維持 NOT RUN，vet 不作為它的替代。
- 前端 865 個 source／test 檔案指紋保持不變。候選差異、行數、分層來源歸屬與限定獨立審查通過。
- 未重跑全 repository race／真設備／部署驗收；既有逐任務 race 與 disposable DB 證據仍限於各自來源與範圍，原逾時、失敗、SKIP 保留。

## 結果分類與授權邊界

| 分類 | 本報告的意義 |
| --- | --- |
| PASS | 指定程式、範圍與環境中的檢查實際執行成功；不可擴大為另一驗證層的成功。 |
| NOT RUN | 該項沒有實際執行，沒有通過或失敗的行為結果。 |
| SKIP | 本次依明示授權略過；實際執行事實仍是 NOT RUN，不是測試成功。 |
| 待驗證 | 尚未建立對應最終候選版本的檢查結果。 |

本次 browser 例外獨立於先前 SQL Server 同機驗收例外。原本的能力矩陣、來源完整性、版本／CAS、擁有權、明確 schema 確認、就緒狀態與未來邊界保護全部保留；counter replay 與 cutover boundary 的後端前置條件也沒有因此略過。

先前凍結的準備矩陣、案例結果及來源雜湊是 **15:27 UTC 批准前的歷史證據**；其中「例外尚未回覆／任務未完成」描述的是當時狀態。本報告只更新本次授權與交付處置，不修改原始雜湊或把任何 NOT RUN 改成 PASS。隔離候選若採用不同來源組合，原 manifest 也不能當作該候選的驗證。

環境準備曾遇官方瀏覽器套件回傳 HTML 而非可用封裝、既有 Chromium 啟動的 OS socket 權限失敗、受支援執行途徑在測試開始前失敗，以及受管瀏覽器無法連到測試頁的 loopback 位址。這些是 setup／環境阻擋，不是產品案例 FAIL。沒有為此調整系統安全限制、外露服務或換用另一台主機。略過權限來自明示批准，不由失敗自行推定。

## 已實際執行的本地檢查

下表是 2026-10-28 **既有工作樹**的證據摘要，與上節獨立發布候選的結果分開呈現。

| 檢查層 | 實際結果 | 不代表什麼 |
| --- | --- | --- |
| 完整前端 unit／integration 序列執行 | PASS：315 個測試檔、2,942 項測試 | Vitest 不含 `tests/e2e`；不是瀏覽器、SQL 或現場驗收。 |
| 完整前端 ESLint | PASS：0 errors；保留既有 E2E 檔的 2 個 console warnings | 不證明互動或畫面正確。 |
| 最終 TypeScript 與 Vite build | PASS；保留既有大於 500 kB 的 chunk warnings | 不代表實際頁面已在瀏覽器開啟。 |
| 正式 Playwright discovery | PASS：3 個 spec 檔、42 個案例；使用已安裝的 Playwright 1.64.1 | `--list` 不啟動瀏覽器，不執行案例。 |
| E2E／fixture host 的明確範圍型別檢查與變更檔 lint | PASS；包含正式 spec、helper 與 test-only host TS／TSX | 一般前端 tsconfig 不含 E2E，不能只用一般 typecheck 代替。 |
| 嚴格 mock API 合約及實際 host 的 jsdom 測試 | PASS：2 個檔案、41 項（29 個合約、12 個 host） | 是較低層測試，不是 41 項 browser PASS。 |
| CAS 修復的 focused 回歸 | PASS：33 項；保留變更欄位拒絕的 guard 檢查 | 不取代實際 browser 中的最後確認。 |
| Test-only fixture build | PASS | 只證明可建置，不是畫面或 HTTP production wiring 驗收。 |
| 既有工作樹空白與行數檢查 | PASS；既有超過 300 行的警告保留 | 不證明隔離發布候選已完成相同檢查。 |
| 隔離發布候選檢查 | 見上節獨立候選結果 | 不把原工作樹結果無條件轉移到另一版本。 |

Focused 測試與完整套件可能重疊，以上測試數不可相加成「唯一通過總數」。完整前端的 2,942 項也不能計入 42 個未執行的正式 browser 案例。

### 保留失敗與修復經過

- 首次完整前端執行以 **exit 137／Killed** 結束，終止前可見測試失敗。原因未獲證明，不標成已確認 OOM；後來的序列成功不抹除原失敗。
- SQL fixture focused run 曾因不可用的既有暫存目錄設定，讓 19 項在 assertions 前失敗。這是 setup failure，不是產品行為 TDD RED；改為可攜的每測試暫存目錄後 19／19 通過，原 SQL assertions 保留。
- 中斷的完整執行另曾出現 writeGroupSection 時序失敗；其後 focused 17 項及完整套件通過。未據此宣稱找到原始根因或完成產品修復。
- 實際 host 的 jsdom 測試揭露真正的 CAS 比較缺陷：相同欄位經 production parser 重建後順序不同，原本以 JSON 字串比較誤判 stale。修復改為語意比較，保留所有版本欄位檢查；沒有重排 mock 回應或放寬預期來掩飾。
- 修復 CAS 後，terminal remount 又揭露 test fixture 缺少 G2 的 ordinary schema-preview 回應。fixture 僅補上已完成操作、屬於本 scope 且 current CAS 吻合的 G2；reserved G2、stale scope、未知 ID 與普通 schema apply 仍拒絕。這是測試支援修正，不是新增產品權限。

## 正式瀏覽器案例與執行狀態

正式來源分為[操作與歷史](../../frontend/tests/e2e/recording-lifecycle-actions.spec.ts)、[重建與切換](../../frontend/tests/e2e/recording-lifecycle-transitions.spec.ts)、[中斷與復原](../../frontend/tests/e2e/recording-lifecycle-recovery.spec.ts)，共用[瀏覽器 helper](../../frontend/tests/e2e/recording-lifecycle.helpers.ts)。

以下 **21 個參數案例各有兩個 viewport，合計 42 案例**。兩個結果欄的「NOT RUN／SKIP」均指沒有實際執行、依本次授權略過，不是測試 runner 已跑出的成功或 skip 統計。

| 案例與參數 | 正式驗收斷言 | 1024×720 | 1600×1000 |
| --- | --- | --- | --- |
| LC-A01 draft-wide | Basic／Advanced 操作可用性一致；草稿與已套用版本分開 | NOT RUN／SKIP | NOT RUN／SKIP |
| LC-A01 applied-wide | 同一 server action authority 與保存版本 | NOT RUN／SKIP | NOT RUN／SKIP |
| LC-A01 disabled-wide | 停用狀態及可用操作一致 | NOT RUN／SKIP | NOT RUN／SKIP |
| LC-A01 deleted-wide | 已移除管理入口、唯讀歷史與明確重建 | NOT RUN／SKIP | NOT RUN／SKIP |
| LC-A02 applied-wide | host 就緒後 Cancel／Escape／Back 無新副作用；鍵盤焦點限制與還原 | NOT RUN／SKIP | NOT RUN／SKIP |
| LC-A03 applied→disabled→deleted | 停止保留目的地；移除需確認；tombstone 與不同待送單位可查 | NOT RUN／SKIP | NOT RUN／SKIP |
| LC-A04 blocked-wide | stale connector 阻擋移除，提供原因與原保存版本的停止路徑 | NOT RUN／SKIP | NOT RUN／SKIP |
| LC-A05 applied-wide | 受管群組回到相同設備／profile 的 Basic 審核，不直接 apply／start | NOT RUN／SKIP | NOT RUN／SKIP |
| LC-T01 wide→long | 完整影響、新 schema 計畫、明確 future cutover；各確認雙擊一次 | NOT RUN／SKIP | NOT RUN／SKIP |
| LC-T01 long→wide | 反向布局使用相同分階段確認契約 | NOT RUN／SKIP | NOT RUN／SKIP |
| LC-T01 long→long／conn-2 | 同布局／plain mode 換目的地，不重用原 schema proof | NOT RUN／SKIP | NOT RUN／SKIP |
| LC-T02 deleted-unapplied | 明確首次開始與 future B，不虛構舊 cutover 或 SQL 交付 | NOT RUN／SKIP | NOT RUN／SKIP |
| LC-T03 deleted-basic | Basic tombstone 固定完整形狀，只保留不同 ID 的新世代 | NOT RUN／SKIP | NOT RUN／SKIP |
| LC-T04 review-blocked-wide | 初次 blocker／stale 顯示處理指引，無 reserve／DDL 確認 | NOT RUN／SKIP | NOT RUN／SKIP |
| LC-T05 applied-wide | 過去 B 拒絕；更改已審核 B 必須重新審核 | NOT RUN／SKIP | NOT RUN／SKIP |
| LC-T06 applied-wide | 最後階段 CAS 失效保留原 root，不送 activation confirmation | NOT RUN／SKIP | NOT RUN／SKIP |
| LC-T07 applied-wide | readiness 不足時保留 schema 準備，不提供最後開始 | NOT RUN／SKIP | NOT RUN／SKIP |
| LC-R01 schema ACK loss／unknown | 切模式、reload、remount 保留同 root；無新的普通或 successor preview | NOT RUN／SKIP | NOT RUN／SKIP |
| LC-R02 preview ACK loss | 僅 server 證明全部未提交時，明確同 root repreview，舊 preview 歷史保留 | NOT RUN／SKIP | NOT RUN／SKIP |
| LC-R03 status failure | 原未知操作保留，不悄悄重送任何效果 | NOT RUN／SKIP | NOT RUN／SKIP |
| LC-R04 activation ACK loss | 查回原 terminal；remount 後遲到失敗及 reload 不使結果回退 | NOT RUN／SKIP | NOT RUN／SKIP |

LC-T01–T07 使用 `Asia/Taipei`（UTC+08），其餘使用 UTC。固定測試時鐘為 `2030-01-01T00:00:00Z`；future B 為 `01:00Z`，在臺北輸入／顯示為 **09:00**。helper 依瀏覽器當地 Date 欄位填入 datetime-local，並從 DOM 值轉回 UTC 核對，不把 UTC 字串裁切後冒充當地時間。這些時區 assertions 同樣未在實際 browser 執行。

### 真實 host 與模擬資料的界線

[Test-only entry](../../frontend/tests/e2e/fixtures/recording-lifecycle-host.html) 掛載真正的 Step4Database、Basic／Advanced、production controller／dialog、QueryClient 與 i18n；沒有複製假 UI。只有 [API transport](../../frontend/tests/e2e/fixtures/recording-lifecycle-api.ts) 與測試資料使用受限 mock，不連資料庫或設備。mock ledger／snapshot 是測試判斷依據，不是 production HTTP、SQL、runtime 擁有權或寫入證據。

既有 Basic host 會對部分非 deleted 的寬表設定自動取得 **ordinary schema preview**。測試先等待其 DOM 結果或明確 dirty blocker，保留初始 token issuance 的完整 ledger，再比較 lifecycle dialog 是否新增副作用；不是宣稱整個 host 從不預覽。ordinary preview、successor preview 與 schema-confirm 分開計數；preview token 不等於 DDL，未列入的路由或寫入嘗試也不能因 mock 拒絕而被隱藏。

## 規格情境對照與未覆蓋邊界

以下對照[完整規格](../../openspec/changes/manage-recording-lifecycle-safely/specs/recording-lifecycle-management/spec.md)；所有列出的 browser 案例仍為 NOT RUN／SKIP。僅呈現 UI 的案例，不能取代右欄的後端契約驗證。

| 規格情境 | 對應 browser 案例 | 必須分開判斷的範圍 |
| --- | --- | --- |
| Edit while the old revision is recording | A01 | 僅區分保存／套用版本及入口；真正保存普通編輯的持久化仍屬其他測試。 |
| An unavailable action has a repair route | A04 | typed reason、保存版本 stop 的後端驗證。 |
| Stop with queued deliveries | A03 | 真實 delivery worker 仍送往凍結目的地，不由 UI mock 證明。 |
| Ownership proof fails during removal | A04 的不可用入口 | 缺失 immutable owner 的原子拒絕屬後端。 |
| Recreate a removed recording | A01、T02、T03 | canonical lineage／新世代持久化。 |
| Recreate a Basic recording with the same complete shape | T03 | 完整來源／不可變語意與擁有權 proof。 |
| Concurrent recreation and retry | T01、R01、R04 | 單畫面重複點擊不等於跨 client transaction concurrency。 |
| Bidirectional managed layout transition | T01 兩方向 | browser fixture 為 SQLite mock，不代表真實 SQLite／PostgreSQL SQL。 |
| Unsupported layout or shape | T04 blocker UI | 完整跨引擎／mode／shape 拒絕矩陣屬後端。 |
| Same-mode long destination replacement | T01 conn-2 | 各引擎獨立 capability／schema／receipt 驗證，不以 mock 通用化。 |
| Cancel an initial successor review | A02 | 從已完成 host baseline 起，無新 reservation／token／DDL。 |
| Prepare exactly one reserved successor | T01、T03 | backend reservation 原子性及 immutable replay。 |
| Reject foreign schema authority | 無注入 foreign token 的 browser 案例 | production parser／service 及後端 scope guards。 |
| Preparation fails | T04、R01、R03 | unknown 與確定拒絕、外部 DDL 結果分開。 |
| Confirmed transition crosses the boundary | T01、T05、T06 的 B 審核 | 真實 fence、observed_at handoff、lateness 與 runtime 排程。 |
| Successor counter establishes its own baseline | 無 browser 數值斷言 | 首筆 baseline／NULL delta 與不跨 gap 補量屬後端。 |
| A reviewed revision becomes stale | T04、T05、T06 | source／target／canonical head 的後端 CAS。 |
| Remove or transition with mixed pending evidence | A02、A03、T01、T03 | 樣本 3、outbox 4、blocked 2、unknown 1 各自標示，不相加。 |
| Lost confirmation response | R01、R04 | server ledger 原 root readback 與重放身分。 |
| Operator leaves a pending transition | R01、R04 | session 中復原，不等於跨新瀏覽器或清除儲存的保證。 |
| Browser execution is skipped under the explicit delivery exception | 全部 42 案例 | 本次批准略過；不建立任何 browser PASS。 |

## 截圖與畫面驗收

**沒有實際 browser 截圖。** 1024×720／1600×1000 的水平溢出、dialog 標題／頁尾可達、長內容捲動、錯誤可讀、鍵盤焦點限制／還原與真實畫面均為 NOT RUN／SKIP。jsdom、build 成功或示意圖片不可替代。

保留的測試只會在未來真實執行時，透過 Playwright `page.screenshot` 與 `test.info().outputPath` 產生、附加實際畫面，並標明 mock API／test data。預定畫面包含 Basic／Advanced、選變更、影響與 schema plan、首次開始／cutover future B、unknown recovery 與 terminal result。

## 已知未解決風險

**一般開始記錄的 A→B→A 遲到成功競態仍未修復。** 在晚到 POST 的 invalidation 完成前由 A 切 B 再切回 A，可能使 A 在 operation ID 已保存後仍顯示 unknown。受控排序在既有版本及本次來源皆可重現；該重現沒有第二次 POST 或錯誤設備歸屬。

這不是本次 successor dialog 的 CAS 修復，亦不因完整測試 PASS 而變成已修復。unknown 不能被解讀為確定未開始；需依原操作查核。本次沒有修改一般開始記錄的產品程式，後續修復範圍仍待處理。

[SQL Server 說明](studio-v2-sqlserver-circuit-recording.md) 中既有長表功能的歷史 PASS 不代表本次 successor 實機驗收；本次實際 SQL 寫入／讀回／cleanup 依其獨立同機例外保留 NOT RUN。真設備、正式 DB、LAN、部署及現場驗收不在本報告的 PASS 範圍。

## 日後重跑與發布候選核對

在具備相容瀏覽器與允許執行的環境中，從 `frontend/` 重跑保留的正式案例：

```sh
npx playwright test \
  tests/e2e/recording-lifecycle-actions.spec.ts \
  tests/e2e/recording-lifecycle-transitions.spec.ts \
  tests/e2e/recording-lifecycle-recovery.spec.ts --project=chromium
```

只加 `--list` 時仍只是 discovery。真正執行後須逐案例記錄 PASS／FAIL／NOT RUN，查看真實截圖，再依結果修復及重驗；不能把本報告的授權 SKIP 自動改成 PASS。

發布時須再核對候選檔案內容與遠端 commit 的 parent／tree／blob，並查詢該 commit 的 CI 狀態。本報告的候選測試結果不等於已發布或已部署；後續程式修改需依影響重驗。原工作樹的 315 檔／2,942 項與候選 303 檔／2,290 項結果、本次 browser 例外及所有未執行風險均隨交付保留。
