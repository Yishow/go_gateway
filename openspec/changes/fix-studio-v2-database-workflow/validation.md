# Review checkpoint：三個 Studio v2 資料庫 change

## 2026-09-15 C2 實作啟動

- 使用者已授權封存 C1 並開始 C2。C1 8/8 完成，封存至 `../archive/2026-09-15-fix-studio-v2-database-result-truthfulness/`；正式規格同步 3 條新增及 2 條修改，無 cleanup warning。封存前 404 個工作樹檔案與上一輪通過完整檢查的版本相符。
- 本次只執行 Implementation Contract 階段 B：1.2、1.3、2.1、2.2、2.4。階段 D 的 5 項待辦尚未解鎖，仍需 C3 3.1–3.10 的實作及驗證。
- C2 preflight clean；analyze 0 Critical／0 Warning／36 Suggestion（補充具體 Example）。開始時 0/10，尚未把進行中的修改列為完成。
- 新 checkpoint：`/private/tmp/go-gateway-c1-archive-c2-20260915/`，保存修改前內容／hash、C1 archive 結果、C2 RED／GREEN logs 及進度。未 stage、commit、merge、push 或部署。

## 2026-09-16 階段 B 續作（1.3、2.2、2.4）

- 開始時工作樹無法建置（2.2 半成品：handler 缺 membership 欄位、Step4Database 仍傳舊 props）；本輪修復後完成 1.3、2.2、2.4，本案進度 5/10。2.3 需 `implement-studio-v2-verified-schema-setup` 3.1–3.10 完成（目前 0/10），3.3、3.4、4.1、4.3 依序等待。
- 2.2：記錄方案 Get／Update／Delete 改用工作區範圍查詢，缺少與跨工作區同為安全 404；Create／Update 只接受本工作區已存測量項目、設備屬於工作區、點位啟用且設備識別一致的成員（422 `RECORDING_PLAN_MEMBERS_INVALID`；讀取失敗或未注入 reader 時 503 fail closed）；production route 只注入非 nil reader。前端以已存測量項目組成員，讀取失敗顯示重試且不顯示建立表單，多方案需明確選擇。RED：`c2-2.2-go-red.log`、`c2-2.2-wiring-red.log`、`c2-2.2-fe-code-red.log`。
- 1.3：Step 4 預覽簽章納入工作區、設備、點位、對應與 row group；啟用中／唯讀時停用預覽與建表；舊回覆不寫入新畫面，remount 保留建表鎖。RED：`c2-1.3-red.log`（啟用中預覽未停用、切換點位後舊預覽仍顯示）。
- 2.4：工作區新增 `database_setup_revision`；連線、欄位對應與工作區參照在同一本地 SQLite 交易寫入，交易先核對 `expected_setup_revision`，舊或缺少版本回 409 `revision_mismatch` 且零 mutation；連線探測在交易外執行。前端以單一佇列依序儲存，每筆帶上一筆回覆的版本。故障注入以 SQLite trigger 使工作區寫入失敗，驗證連線、欄位對應與工作區皆回到原狀，首次儲存失敗不留孤兒連線。RED：`c2-2.4-go-red.log`、`c2-2.4-fe-red.log`。Mutation：rollback 改 commit 時三個故障注入測試失敗；佇列加一次重試時「不自動重試」測試失敗；回復後逐位元相同並通過。
- 證據檔位於本 session scratchpad（`/private/tmp/claude-501/…/scratchpad/`），屬暫存空間。

| 結果 | 檢查 |
| --- | --- |
| PASS | gofmt、`go build ./...`、`go vet ./...`、`go test ./...`（44 packages）、`golangci-lint run ./...`（0 issues） |
| PASS | 前端 `npm run lint`、`npx vitest run`（160 files／927 tests）、`npm run build` |
| PASS | `make check-lines`（無 >500；Step4Database.tsx 320、safeJson.ts 322、useStudioV2DatabaseAutosave.ts 324 為 >300 警告，後者由 371 降至 324）、`git diff --check` |
| NOT RUN | 瀏覽器／E2E、真實 PostgreSQL／MySQL、跨程序並行儲存、現場 PLC；未 commit、未部署 |

- Review 修正（同日）：同列在前一筆回覆前再次提交時保留較新值並重送，列上待確認值依序清除（1.2／2.4）；改選另一個已存連線視為分組範圍變更（2.4）；排隊中的連線儲存改用最新 identity 與設定版本，回覆不覆蓋較新的表單值（2.4）；`SchemaSetupSection` 在 StrictMode 重設掛載旗標（1.3）；停用連線改回 422「請啟用連線」而非 404（2.1）；欄位入列即標為儲存中（2.4）。9 個新測試修正前失敗（`review-go-red.log`、`review-fe-red.log`），修正後 `go test ./...` 44 packages、golangci-lint 0 issues、vitest 161 files／934 tests、lint、build 皆通過。
- 第二輪 review 修正（同日）：儲存期間改選另一個已存連線時，回覆只採用設定版本，不再把選擇蓋回舊連線（2.4）；列上待確認值在儲存或被拒絕時都清掉較舊的值，並允許改回較早的值立即送出，重複提交只擋最近一次的值（1.2）。2 個新測試修正前失敗（`review2-red.log`），修正後 vitest 161 files／936 tests、tsc、lint、build、`make check-lines` 皆通過；後端未變更。
- 相容性：已有設定版本的工作區，舊客戶端未送 `expected_setup_revision` 會被 409 拒絕；首次儲存（版本為空）不受影響。Swagger 未收錄 `database-config`／`database-targets`，本輪無需重生；Studio inventory 的 Step 4 描述仍正確，欄位契約同步留給 4.3。

以下各節是先前提案審查的歷史紀錄；其「僅文件／未實作」限制不代表後續已授權的實作仍未開始。

## 本次範圍與基準

2026-09-15 依使用者確認，只審查並修正三份 change 的 proposal、design、delta specs 與 tasks，另更新其來源及發布交接文件。產品程式、主規格、資料庫、部署、Git index 與 commit 都不在本次修改範圍。

起始本機 `main`／HEAD：`569da9f98e52b2635638d85250a4af16601d4364`，起始工作樹乾淨。原提案依據 `1a0311c8e8db9c62fe4388f0701ba38afe552ff7`；兩基準之間只有提案文件差異。未查遠端 main，後續實作前仍須核對當時工作樹。

## 修正內容

| 問題 | 本次文件修正 |
| --- | --- |
| 假成功、建表與進度需求在多案重複；archive 可能覆蓋另一案較完整的條文 | 每條 requirement 只由一個 change 負責；移交四項重複待辦，保留原編號及 Supersedes 紀錄 |
| C2/C4/C5 沒有發布，卻被當成可執行前置 | 改為已存在的 truthfulness→workflow 前置→schema→workflow 後段，明列跨案待辦 gate |
| 一般建表「保持相容」未區分階段，generic handler 缺少明確任務 | 安全封鎖階段保留已實作功能；schema 階段同時遷移 workspace、generic 與 activation 入口，封住無確認 mutation |
| 同操作進行中同時寫成 409／202，operation_id 來源未定義 | 後端 preview 發給固定 operation_id；同操作 running 回 202、終態回 200，別操作佔同範圍才回 409 |
| 試寫要求持久操作編號，卻缺少預覽、確認與查詢的交接責任 | workflow 3.3 沿用 schema 的持久機制擴充 `test_write`，明定 token 種類隔離、試寫預覽、版本檢查及狀態查詢／重複要求契約 |
| 任務缺少可解析需求／設計對照與完整交接契約 | 補上 Implementation Contract 及具名需求／設計對照；舊 token 相容讀取與拒絕條件列入任務 |
| 前端共用解析器不認識新增安全 code/action | 在 truthfulness 的 Impact／design／tasks 加入 safeJson／typedErrors、精確允許值、翻譯及完整 envelope 測試 |
| 連線版本、scope resolver、operation 查詢與故障測試責任不完整 | 指定 workflow 前置及 schema 的型別／repository／migration／route／hook 責任，列出失敗與重啟矩陣 |
| 舊動畫規格被誤當現行程式現況，空結果成功文案仍缺證據 | 保留目前 activation 資料流；規格／待辦聚焦後端證據、空結果、部分成功、重載及 runtime 導航 |
| 待辦將 router 測試放到 handler 目錄，並把既有 route 檔寫成待新增 | truthfulness 1.1 改沿用 `internal/api/` 的 router fixture；schema 3.4 明確更新既有 route 檔並新增 status handler |

詳細分工以 [design.md 的 Implementation Contract](./design.md#implementation-contract) 為準。去重後共 28 項待辦（workflow 10、truthfulness 8、schema 10），全部未完成；減少數量不是完成實作。

## 本機文件驗證

主代理已回讀實際 source／主規格及全部變更 diff，並對最後版本重跑以下文件檢查。未沿用子代理或初稿的檢查結果作為完成依據。

| Change | `spectra analyze <change> --json` | `spectra validate <change>` | 產品待辦 |
| --- | --- | --- | --- |
| `fix-studio-v2-database-result-truthfulness` | PASS：0 Critical／0 Warning；13 Suggestion | PASS | 0/8 |
| `fix-studio-v2-database-workflow` | PASS：0 Critical／0 Warning；36 Suggestion | PASS | 0/10 |
| `implement-studio-v2-verified-schema-setup` | PASS：0 Critical／0 Warning；17 Suggestion | PASS | 0/10 |

- **PASS**：`git diff --check`、主代理跨案完整性檢查。18 條需求均有待辦對照且只有一份 active delta；28 個本案／跨案 task nodes 無缺失依賴與循環；相對文件連結皆存在，三案皆有 Implementation Contract。
- **PASS**：`spectra list --json` 的 completedTasks 全為 0；`git diff --cached --name-only` 為空。
- **保留的非阻擋建議**：65 項是建議為 scenario 增加獨立 `Example`，既有 WHEN／THEN 與具名測試仍是驗收依據；1 項是 schema scenario 中 `may have executed`，描述回覆遺失後是否已執行的不確定事實，並非可選的安全要求。這些建議沒有宣稱已清除。


## 限制與交接

- **NOT RUN（最終文件交付）**：未對本次文件變更執行 Go／前端產品測試、lint、build、瀏覽器、真資料庫、硬體與現場驗收；這些產品驗證仍由各案實作待辦要求執行。已撤回版本的測試紀錄見下節，不列為目前產品驗收證據。
- `spectra-review` 初始三案 scope 均為 `approximated_worktree / insufficient`，沒有已驗證的產品實作比較基準；snapshot 檢查為 current。此次交付是 artifact review／修正，不宣稱完成程式碼 review。
- 最先可執行的是 `fix-studio-v2-database-result-truthfulness`；另外兩案必須先核對 design 指定的階段前置與測試證據，不能因 validate 通過就整案直接 apply。
- 未 archive、未 stage、未 commit。修改前文件備份：`/private/tmp/go-gateway-change-review-20260915/before/`；本次 diff 完成後保存在同目錄的 `review.patch`。這些是本機暫存復原依據，可能被系統清理。

## 續作範圍校正

續作曾誤延伸到 truthfulness 產品實作，已停止並收回 31 份既有產品檔案的修改、移出 8 份新檔，且恢復所有產品待辦為未完成。收回後的工作樹 diff 已與續作前 `before.patch` 逐位元比對相同，再完成待辦定位與試寫交接契約修正；最後交付仍限 16 份文件。

額外實作、測試紀錄與收回依據保留於 `/private/tmp/go-gateway-c1-20260915/` 及 `/private/tmp/go-gateway-scope-recovery-20260915/`，不當作現行程式通過驗收。三案已重新 analyze／validate，結果如上；最後文件 diff 為後者的 `final-document-review.patch`。未部署。後續若另行授權實作，須重新建立當時的比較基準，不沿用此次已撤回的完成紀錄。

## 初稿檢查紀錄（歷史）

初稿在無 Spectra CLI 的環境手工起草，原先記錄 115 個簡單結構檢查、13 條 requirements、53 個 scenarios、14 個 workflow 待辦；這些是當時文件計數，不是目前文件或產品測試通過數。當時未執行官方 analyze／validate／park，也未執行產品測試。本次本機結果以上節為準，沒有回填虛構的歷史 CLI 證據。
