# 實作與驗證證據

## 當前狀態（最新指示）

本案非browser實作、完整前端檢查、Go web檢查及49-path review完成，現在為9/11。3.3已完成；僅3.2真browser與3.4最終驗收／archive待辦。依使用者最新指示，非browser相關變更可先獨立commit，繼續後續三案，全部實作完成後一起在Mac驗收。雲端不得擴權，browser尚未通過，archive仍保留。以下舊8/10與「未繼續下一案」是前次checkpoint，不是目前的停止要求。

## 1.1 實作前基準（2026-11-06）

- GitHub main 與本地 HEAD 均為 `0f54740e2df083ff0e6695897bfa4963039d5fda`；原有 `DESIGN.md` 視覺提案增補保留，不屬本案。
- 環境：雲端 Linux；Vite 8.3.5、React/React DOM 19.3.1、TypeScript 7.0.2（lint alias 6.0.2）、Playwright 1.64.1、Go 1.25.5。
- 固定 fixture：production index.html 的 `/studio/v2` recording-home 靜態 import graph，無 API fixture／瀏覽器狀態。
- `npm run build -- --manifest`：PASS；完整 manifest 中所有引用可解析。
- eager JS 共 6 檔，原始 **2,467,102 bytes**；相同檔案 Python gzip level=6、mtime=0 為 **695,754 bytes**。後者是可重現壓縮估算，不是網路 transfer 或互動時間。
- entry `assets/index-CfsFF1WQ.js` 原始 2,186,023 bytes；其餘初始依賴：createLucideIcon、jsx-runtime、react、runtimeOperatorSafety、uiVersion。
- 以實作前 handler 原碼及相同 dist 的 `httptest` 記錄 GET/HEAD：JS identity 與 gzip request 均送未壓縮 2,186,023 bytes；Cache-Control 為 `public, max-age=31536000`，無 Content-Encoding、Vary、ETag。HTML 2,550 bytes，沒有 Cache-Control。缺資產 404。
- 這是 Go HTTP handler 契約測試，沒有 socket、實際網路傳輸或 browser timings。完整 build/manifest/headers 原始紀錄保留在本次雲端交接證據中，未提交生成 bundles。
- 真瀏覽器 BLOCKED：雲端 `/usr/bin/chromium` 正常 headless 啟動並指定可寫 profile 後仍 exit 134，`socket() failed: Operation not permitted (1)`。未擴權、停用安全機制或改用 Mac；cold/warm、互動時間、五次 median/p95、deep-link/chunk failure/version browser 驗收均 NOT RUN。

## 1.2–2.5 失敗回歸與實作

- 路由先新增延遲載入／意圖／草稿導航與錯誤回復測試。初始 3 個失敗含一個測試路徑設定問題；後續修正測試路徑，再以原版 HEAD 內容重測確認 eager import 真正失敗，還原最終實作後通過。沒有把設定錯誤當行為 RED。
- 2.1 原樣抽離 advanced owner 後：5 檔 19 項 routing／source／draft navigation tests 與 production build PASS。
- 2.3 先測 negotiation：缺 gzip/Vary、拒絕 identity 卻回 200/206 的失敗已記錄。2.4 實作固定資產 gzip cache，無 request-created keys；`go test ./internal/web -count=1`、`go vet ./internal/web` PASS。
- 2.5 先測 validators/cache：缺 ETag、HTML no-cache/304、未 hashed JS 錯誤長快取的失敗已記錄；實作後 package tests、vet、diff check PASS。
- HEAD、MIME、Range、404 與 API/SSE 保留行為另以 deliberate MIME／全域 gzip middleware 變異證明測試會失敗，已還原並通過。
- Fixture gzip retained payload 158 bytes／eligible original 7,169 bytes；這不是完整 process RAM，production graph 成本另列。

## 2.2 最終載入分離

- 24 檔 108 項 focused tests PASS（預設 6 workers）；含真正 rejected import、安全錯誤、query/hash/locator、dirty Back/Forward、home 保留、延遲 modal Close 與 reopen epoch。
- cold module transform 曾超過既有 owner tests 的 1 秒等待；owner-only 測試改對抽離後 owner，3 份 home owner interaction suite 預載真實模組。沒有弱化 assertion／提高 timeout；獨立 delayed-import tests 保留冷載入驗收。
- `AdvancedWorkbenchV2Page` 的 bootstrap／owner body 與原 HEAD 除 default export 外完全一致。

## 3.1 最終交付與比較

- `bash tests/shell/embedded-frontend-delivery.sh`：21 tests PASS，stale removal、placeholder、缺 input／失敗 build fail-closed。
- `node scripts/build-single.mjs --os linux --arch amd64`：PASS，包含 TypeScript／Vite production build、完整 static sync、Go embedded binary。static 與 dist 108 檔逐一 bytes 相同，另保留 placeholder，無 stale child。這是 Linux build，不是假稱另外五個平台實測。
- 完整 manifest 103 entries、189 references 全部可解析；保留 route chunks，advanced／charts 不在首頁 eager closure，react/app/misc logical groups 存在，無依賴的 MUI group 未產生。
- 相同 toolchain／fixture eager JS：2,467,102 → **1,393,373 bytes（減少 43.52%）**；同方法 gzip estimate：695,754 → **432,389 bytes（減少 37.85%）**。graph-derived JS requests 6 → 9，是多個較小 shared chunks 的取捨，不是 browser 觀察。
- 實際 handler／httptest eager JS response body：原版 identity 2,467,102 → 新版 gzip **436,624 bytes**；Go gzip 與 Python estimate 分開記錄，沒有宣稱網路／首次可互動時間。
- 完整22份 eligible JS/CSS：原始3,452,154 bytes，固定 gzip payload936,305 bytes；130個validator字串8,580 bytes。3次固定 fixture 的 retained Go HeapAlloc delta 為993,896／993,928／993,864 bytes，含存活 map/string/slice allocations，均低於原始 eligible bytes。不是 RSS、總 process RAM 或 startup peak。
- 最終 SHA-256：`fda3580ad26bdbfa258234446be60201c7029b30471fd0028e5c1db2f2a11e80`（Linux embedded binary，生成檔不提交）。
- Vite 仍警告 entry／advanced chunks 大於500kB；本案已達相同比較門檻，不藉此擴大拆分所有編輯器。

## 前次 checkpoint（已由最新指示更新）

3.2 真 browser 仍 BLOCKED，3.3 依賴它，不能把整案標完成；完整frontend tests648檔7397項通過；verify仍因browser列未完成，review49paths無發現。未 archive、commit、push 或部署。

## 測試遷移與限制

- 抽離 owner 後，owner-only 測試改import新模組；Home／locator互動測試保留真實入口並預載已編譯Advanced，backend延遲／CAS／ACK assertions 不變。真chunk延遲、失敗、Close與epoch由獨立測試保留。
- 最後一輪相關25檔171項通過，所有調整未放寬assertion或timeout。曾發現locator fixture的doUnmock順序，已將該預載放回doUnmock之後並重驗。
- 前次Spectra verify：8/10，Critical僅3.2／3.3未完成與browser驗收未結案；49path獨立review無程式缺陷。
- 既有DESIGN.md增補未改；另外三案仍parked、0項完成。

## 最終可執行檢查

- `npm test -- --run`：648檔／7,397項全數通過，206.76秒，沒有skip或縮小全套範圍。
- 最終 `npm run lint`、`npm run build`（含TypeScript）、`go test ./internal/web -count=1`、`go vet ./internal/web`、`git diff --check`、`make check-lines`、Spectra validate均PASS。再次build的108份dist與已嵌入static逐byte相同，先前binary與量測仍適用。
- 前次Spectra verify為8/10：原task3.2真browser與原task3.3依賴結案未完成；spectra-review已重審49份source/test內容，0 Critical／0 Warning／0 Suggestion，scope snapshot current。
- 行數警告：Advanced376、locator context407、mapping reconciliation351、memory view343、source modal411、workspace boot469；歷史shell621未增加。後續分拆activation/status helper與test fixture groups，這次不為消除警告擴大scope。
- browser E2E、五次cold/warm timing、真正瀏覽器截圖、Windows/Mac、真PLC／DB與部署未執行；golangci-lint未安裝。不得以單元、HTTP或Linux build代替。
- 交付前再次核對遠端main仍0f54740e；沒有archive／commit／push，沒有擴權、關閉安全機制或Mac操作。

## 驗收順序更新

2026-11-06使用者明確要求全部非browser實作後再一起於Mac驗收，並保留archive以避免停止實作。前述「未archive／commit、後三案parked」是本次前段checkpoint，不是新的停止要求；後三案可依序繼續。雲端仍不得擴權。非browser品質已通過，browser與最後archive任務維持未完成，可先提交已授權的本案程式與證據。
