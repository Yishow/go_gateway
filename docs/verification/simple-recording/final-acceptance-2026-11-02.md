# 七案最終驗收記錄

日期：2026-11-02 UTC。此文件由實作、測試與雲端主旅程的實際收據整理；歷史失敗不覆寫，版本沿用與本輪執行分開。

> 原七案功能、最後完整 FE、官方建置與選定雲端驗收已閉合；歷史／環境未驗限制及尚未同步的發布狀態如下。

## 範圍與版本

原七案為 `simplify-recording-home-guidance`、`streamline-recording-task-editor`、`simplify-recording-lifecycle-recovery`、`adopt-simple-recording-entry`、`preserve-per-point-word-order`、`support-per-member-recording-calculation`、`integrate-protocol-aware-acquisition-batching`。不新增記錄引擎、連線管理平台或產品功能分支。

- 最新呈現修復／官方建置版本：`5025bdf609dc030193d64b505d5281bc02cde401`；同版最後574 files／6743 tests全PASS。前版 `e1184b19` 的完整PASS與其後像素FAIL均保留，不以舊結果取代這次重跑。
- 完整 Go／vet 實測版本：`0175b8e8a83ddebd4a0a30a2ebd9a90ebe59bff6`。其後只由官方 formatter 重排一個測試檔的 import；import 集合及非 import bytes 完全相同，另跑該檔原 7 top／8 leaf 與完整 lint。
- MAIN 是同一自有測試設定／歷史，跨 `c88afdfe`、`b21be972`、`cffaf969`、`ca0fd23a`、`a8a0b2ce` 修復後續跑的完整主旅程；不是聲稱每一步都在最後提交重跑。
- 最新已確認遠端 `main` 仍為 `17972a409f5118d7acb3b0c34742ebc78ee13303`。後續 native commits 尚未發布：既有同步授權下本批遭自動核准中止，依返回要求待使用者再次明確確認；本機驗收與 GitHub 同步分開。

## 完整程式驗證

| 驗證 | 實際結果／版本 |
| --- | --- |
| `go test -p 1 -count=1 -timeout=15m -json ./...` | 0175 完整單次：64 packages＝59 PASS／5 no-test；5,367 top PASS、12,538 unique leaf PASS、520 leaf SKIP、0 FAIL；未用 test-result cache |
| `go vet -p 1 ./...` | 0175 PASS |
| 官方 golangci-lint 2.6.2 `run --concurrency=1 --timeout=10m ./...` | 原0175僅新發現測試的 goimports 1 FAIL；官方純 import 格式後完整重跑為 0 issues，另原7 top／8 leaf PASS；不稱格式後再跑完整 Go |
| `npm run lint` | 5025 PASS（01:17:30–01:17:47 UTC） |
| `npm test -- --run --maxWorkers=1` | 5025：574 files／6,743 tests 全 PASS，01:17:55–01:30:24 UTC；6,370 tracked inputs 前後完全相同 |
| `node scripts/build-single.mjs --os linux --arch amd64` | 5025：01:13:15–01:13:36 UTC PASS；strict／Vite／static／Go embed 全過，110 dist/static files 相同 |

完整 Go 的 520 skips 是未提供的 PostgreSQL fixture/DSN 293、verified TLS 124、MySQL/MariaDB 68、SQL Server 26、child helpers 3、缺 mock 2、scheduled opt-in 1、SQLite 不適用 1、既有未實作 duplicate-name case 1、硬體不可達 1。此輪 BuiltSPA 已直接 PASS。這些 skip 不等同現場 DB／TLS／PLC 通過；後述 MAIN 的自有 SQLite／PostgreSQL E2E 亦不等於全 engine fixture 全跑。

原 `assembler_mixed_windows_test.go` 因 Go OS 檔名規則在 Linux 被排除，並非 skip。以相同內容改為中性 `assembler_mixed_window_test.go` 後，原7 top／8 leaf 已由本輪完整 Go 真正執行。沒有新增斷言、更改語意或把舊漏測結果算成已跑。

全 lint 通過包含19行具體說明的 Go 規則例外與原 FE 語意 memo 例外；不是零 waiver。Go 的 bare sentinel 比較保留 wrapped/joined 真 cause，獨立 cleanup／fresh projection context 邊界不被請求取消或舊快取污染。精確位置、理由與測試索引另列例外表，沒有關閉全域規則或排除整檔。

## 真實雲端主旅程

使用自有 Linux、HslMacSimulator、SQLite、官方 PostgreSQL 17.11 與瀏覽器，loopback／新自有資料庫，非正式資料或現場 PLC。每次啟動核官方 binary／HTML／index／lazy page SHA；避免把 restored tab 的舊 JavaScript 當最新版本。

### 三協議與 SQLite

- Modbus TCP float32 12.5、MC D0 uint16 125、FATEK D0 uint16 3725；正式 UI 建立／保存來源，真 point poll 與 UI 同值，非後台補畫面。
- 首次明確 Start（後續另有一次下列 Resume），三值於完整5秒 joint 窗寫入 SQLite，SQL 與 receipt 對應。
- dirty Keep/Discard、Stop Cancel/Escape 零停止 mutation、確認停止、fresh Check/Confirm/Resume；保存獨立副本，只移除未開始的副本，原記錄／歷史不變。
- 寬表→長表於2026-11-01 22:56 UTC邊界；S1長表→S2於同日23:00 UTC邊界。舊末窗與新首窗對齊，無新的 Start 取代 forward 操作。
- S1 最終113寬表＋144長表／257 receipts 完全保留；S2 自然續寫，原411列保留至864列／864 receipts。數量是記錄的備份時點，不冒停止瞬間所有 pending flush 已完成。

### PostgreSQL 與精確差值

- 正常 UI 選 Settings 已保存 P1，不寫 legacy workspace connection config，不改 A；新 B Save 成功並轉 canonical locator。
- 已保存目的地卡、摘要與 scoped DTO 一致顯示 PostgreSQL／P1／public；卡片唯讀，前向變更仍走原獨立確認流程。
- mean5秒為12.5；counter15秒首窗建立1000基準，不寫原累積值或0；1011產生精確差值11，payload/effect/digest 與 PG receipt 對應。
- P1→P2於2026-11-01 23:53 UTC切換；P1末窗結束於邊界，新P2使用1200 fresh baseline，不借前代1011。其後1217、1240真 SQL差值17、23，mean仍12.5。
- 2026-11-01 23:59:29–50Z 有界全量快照 P1為59列／59receipts，P2為107／107（未達2000上限，未截斷）；最後 PG dumps另存。
- 這是同 engine 的 SQLite與PostgreSQL forward，未宣稱跨 engine 轉移、所有DB支援或 mean全覆蓋；無完整採集證據時仍顯示 coverage-unverified。

### 導航、畫面與安全收束

1440／768／390、scoped logs真正W/G query與canonical返回／BackForward、可達操作與關閉流程保留。最後console無JS exception，仍有既有form id/name提醒。曾有工具自動捲動讓按鈕被footer擋，正常捲動後可達，明列工具操作觀察而非虛構CSS修復。

每輪停止原owned controller／PG／HSL／Gateway／Chrome／MCP，核PID退出、SQLite一致備份/integrity、PG dump與原7份資料庫hash不變。原 MAIN停止時間00:00:12.886Z；最後resume06於01:16:39.672Z全部停止、01:17:17核六PID全退出；27個GET均200、零UI寫入，console僅既有form id/name提醒。其停止dump P1維持59、P2 179→213，S2 1014→1092，所有prior rows/receipts精確保留、S1及old7資料庫不變。

## 真失敗及修復

1. FATEK draft diagnostics漏原canonical mapper：station別名未轉station_no。只修caller，原後端schema／站0／已存相容保留，正常raw-default UI復驗成功。
2. Step3把persisted ID不等於local ID當必要條件，合法MC同ID被漏poll。改用既有mapping/snapshot證據，未存UUID仍拒，late scope不投影，SSE retry原回歸保留。
3. 新記錄目的地漏canonical connectors query且誤寫legacy設定。new/copy改選既有canonical目標，只改原draft；stale/disabled/revision/ABA仍拒，A不變。
4. 保存後目的地卡仍呈現外層工作區S1。改以matching saved scope安全唯讀卡，不改全域設定或forward權威。
5. 受管長表誤顯寬表缺欄提示。只排除managed+managed_long_v1的wide-specific proposal，wide真缺欄／custom仍提示，mapping/readiness/Start guard未動。e118單頁實載像素證總proposal已閉，但member Status仍錯顯缺欄，保留真FAIL；5025已一併傳遞layout，僅managed_long略過本地逐member physical-column missing/incompatible呈現，真正issues/conflict仍紅，10相關tests及4檔lint通過，5025於resume06實載三asset SHA一致，1440/390每列皆「已指定」、無false missing/incompatible／proposal，P2卡正確；同版完整FE 574 files／6743 tests已於01:30:24 UTC全PASS。
6. 最後47df完整FE有3 FAIL：query hook侵入原mixed early readonly、以及fixture漏分類既有採集POST。只將hook移入已接納child、fixture精確限制已定位PA1；原唯讀斷言／設定寫入拒絕不放寬。21相關case通過；e118 該輪完整 FE 574 files／6742 tests 全通過，最終 5025 則為 574 files／6743 tests 全通過。

## 未知回覆與性能的證據界線

Unknown02為歷史版本一次Start成功後截斷ACK的真Browser證據；保同operation、Escape/Keep/Close/BackForward後先lookup，無replacement mutation，原規格恢復場景通過。初始ACK斷後同request_id/body在新socket出現wire duplicate仍保留，proxy阻擋第二forward；不稱Browser驗到server第二次去重，也不稱wire0。沒有第三次故障注入。後續只讀完成投影Run06以實載asset證據閉合；cold cache無完成證據仍unknown，不將舊成功當目前來源健康。

Batching current c88測原15缺格、564 device-cases／56,964 required samples，每device1cold+100warm，CPU7,8/GOMAXPROCS2/cadence1秒；15身份契約PASS，10格達零outputdrop、5格未達。原規格允許可觀測有界drop，此為容量目標未達而非新造硬門檻；sink0／無recording interest，不推論DB丟同數量或所有值交付等價。6 planner cells各事前固定1秒×5，共30rep，完整保ns/op/B/op/alloc，非p95。

本次驗收的paired18與窄race採用歷史實測＋39/11逐檔語意審查；不冒當前重新paired測速／全race。正常路徑小guard/tuple成本未重新paired量化。最後1306 production Go/module/lintconfig與c88相同，FE修復與測試中性命名不重跑未變性能矩陣。

## 可追溯收據與發布

收據索引相對本任務保存根，僅供查核，不是公開下載承諾：

- `final-complete-gates-0175b8e8/FINAL-SUMMARY.json` SHA256 `66a6b3333ffd589c42b7fe5ee21a28833ef6ec7becff22d2a027fdcde8f1d5a1`（保留該輪 goimports FAIL）
- `final-complete-gates-0175b8e8-format-closure/FINAL-SUMMARY.json` SHA256 `e97d276fbceb6b8111a4230b4ae079b1df5ac379964297ea8d7aba88fff8ab23`（純格式增量及完整 lint 閉合）
- `final-fe-aggregate-5025bdf6/`：完整 FE log SHA256 `cef6b355b72ac7c52d8daf0b8692c7ea9fd04768eba19bb7ce689193e8a61793`；`final-build-artifacts.json` SHA256 `4c04b4133a18377c9dfad3d68ea18e78f1f6e3434bcda8b9ec2f236c4e7461c7`；binary SHA256 `699aa07cd578ac596120682dab9d79380b730dd5edbe5973d022f288957775ee`
- `batching42-final-c88afdfe/CURRENT-INDEX.json`、`semantic-map/`
- `linux-qa-tools/evidence/final-c88-main/`、`final-main-resume01`至`04/`，每輪有raw/asset/SQL/receipt/backup/cleanup紀錄
- `linux-qa-tools/evidence/final-main-resume06/RUN-MAIN-RESUME06-QA.txt` SHA256 `e9c8d18b90fb14a9a19d1f5a8553c5bfe3ad2c2a5973095571c56c4166c28707`，保留resume05像素FAIL與06閉合
- `linux-qa-tools/evidence/lifecycle-unknown-02/ACCEPTANCE-SCOPE-ADDENDUM.txt`
- 七案逐requirement/scenario/example、review歸屬與scope工具限制另列conformance摘要；沒有拿乾淨tree冒完整歷史差異審查。

不納入raw資料庫、私有profile、credentials、完整敏感payload或大型截圖。Change不archive，保留可檢視；原七案依實證完成88/88 tasks，並保留歷史歸屬與未驗環境限制。既有同步授權下，本批遭自動核准中止，依返回要求待使用者再次明確確認；不代推、不改路繞拒絕。
