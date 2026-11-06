## 1. 保留來源觀察契約

- [x] 1.1 在 optimize-embedded-route-loading shared page split 完成後，建立 Compact home source observations preserve canonical evidence 的 parity fixtures：相同 device/rule/mapping/point 投影、duplicate/foreign/missing/source-change/read-error；測試先行，確認舊 projectSources 與新 server projection 的 candidates/excluded 完全相同。
- [x] 1.2 [after: 1.1] 依唯讀來源觀察與完整性重用已存在的 Service.ReadSnapshot existing-workspace read accessor，只委派 Repository.Get；read_snapshot_test.go 回歸與 mutation check 證明不存在時不建立、read error 不修復、不觸發 reconcile，並確認 nil／取消／detached copy 保護。
- [x] 1.3 [after: 1.2] 完成兩輪 canonical compact projection/digest 與 identity checks；recording_home_sources_test.go 的 Go/TS parity fixtures、變動／duplicate／foreign負案例證明 available/unknown 界線一致。
- [x] 1.4 [after: 1.3] 新增 compact handler 的 scope、安全 DTO 與 typed error，studio_v2_recording_home_sources_handler_test.go 驗證 GET 零 mutation、404/400/503、敏感欄位排除、nil dependency unavailable。
- [x] 1.5 [after: 1.4] 正式 router 接線並產生 API 文件；router integration test 驗證 cmd/test_ui service graph 可達，make gen-docs 與生成檔 diff 驗證 annotations 契約。

## 2. 首頁漸進遷移

- [x] 2.1 [after: 1.5] 先增新 service/query adapter 並以失敗案例鎖定 Return navigation reuses observed presentation without reusing authority，再依快取呈現與請求隔離遷移 source read；recording-home-source-summary.test.tsx 覆蓋 A→B→A、取消、late response、失敗、QueryClient/backend 同 ID 隔離、第九 scope eviction 及 cache fresh/stale 分離。
- [x] 2.2 [after: 2.1] 完成 Home data demand stays bounded and complete 與漸進替換與有界首頁：一次 revalidation 一個摘要 GET，保留 visible-page detail ceiling、filter/page/scroll/focus，成功 device/rule/point/mapping mutation 精確 invalidation 後不可假 fresh；來源／route／evidence／facts／collections 測試及 request-count assertions 通過。
- [x] 2.3 [after: 2.2] 同 fixture 比對摘要 bytes 與 canonical candidates/excluded，記 server reads/time 不冒稱都改善；parity、unknown guards 與回歸通過後才移除首頁舊雙輪 HTTP adapter 路徑；因 GroupEditor 仍用於明確 locator-scoped discovery，原樣移至 recording-editor/recordingSourceDiscovery.ts 並精確更新 editor caller／tests，保留49-fixture parity 與 editor caller regression，首頁 request-count 斷言證明不呼叫它；原 full source APIs 不刪除，npm run build 與 focused Go/frontend tests 證明可交付。

## 3. 返回導覽驗收

- [ ] 3.1 [after: 2.3] 依雲端實作與最終 Mac 瀏覽器驗收，在全部四項 change 非瀏覽器實作完成後集中於 Mac browser 驗證首頁→編輯→返回、Back/Forward、scope switch、refresh error、重複 retry 無 mutation 且保留 presentation，記 cold/warm request/bytes/timing；受限制時保持未完成並回報，不擴權或重試受限雲端 Chromium。
- [x] 3.2 [after: 2.3, 3.4] 完成 scoped source review、affected Go/API tests、frontend lint/test/build、git diff --check、make check-lines、Spectra verify/review，非瀏覽器品質門檻與 browser／field 狀態分列；browser 未跑不得當成通過，並交接 exact invalidation keys 給 shorten-recording-save-feedback。

任務切片以每項不超過一小時為規劃估算，不是工期承諾；若實作探索證明超過，先把同一行為拆成可獨立測試的垂直批次並保留 dependencies，不擴張 scope。所有 checkbox 目前未執行。

- [ ] 3.3 [after: 3.1, 3.2] 最終確認 Mac browser 證據與全部非瀏覽器檢查有效，Spectra verify/review 無未處理阻擋後，才交接 archive；非瀏覽器品質通過後允許先做 scoped commit 並繼續其他實作，不能宣稱 browser 通過。

變更註記：原 3.1 雲端 browser／不使用 Mac 限制依最新使用者授權由設計的雲端實作與最終 Mac 瀏覽器驗收取代，不視為已完成。

契約澄清：source inspection 發現 GroupEditor 的明確來源 discovery 是舊 adapter 的既有 caller。2.3 僅重定位 editor ownership，不刪除其既有 authority／await／readback，也不改用 compact evidence；新舊路徑須明確歸屬與提交。

- [x] 3.4 [after: 2.3] 修復review與parity驗證發現：ECMAScript BOM／NEL label trim精確一致；persisted null／non-array skips與mapped empty skips不得提升available，合法unconfigured empty保留。以49-fixture parity、十項persisted-source skip regression、四package focused Go tests／vet與new-from-base lint確認，逐檔歸屬review修復，不更動已完成任務歷史。
