# 非瀏覽器實作驗證

## 結論

三項 requirement、八項 scenario 均有實作與測試對應；沒有未處理的程式契約偏差。非瀏覽器任務完成後為 10/12；tasks.md 的 3.1 與 3.3 仍未完成，依使用者指示留待四項 change 完成後集中在 Mac browser 驗證。不能 archive 或宣稱 browser 通過。

## 完整性與待辦

- CRITICAL：tasks.md 的 3.1 尚無 Mac browser 實測。完成首頁／編輯／返回、Back/Forward、scope switch、refresh error、retry 與 cold/warm request/bytes/timing 後再勾選。
- CRITICAL：tasks.md 的 3.3 最終驗收與 archive 尚未完成。取得上述 browser 證據並重新 verify/review 後才交接 archive。
- WARNING：全 Go test 與全 lint 並非全綠，詳見下列既有問題證據。不能將 focused pass 擴稱全倉通過。
- 無缺失 requirement、未覆蓋 scenario、未處理設計偏差或需排除的功能測試。

## 契約、測試與執行

| Requirement | 實作位置 | 有測試的 scenarios／example | 執行結果 |
| --- | --- | --- | --- |
| Compact home source observations preserve canonical evidence | internal/datalink/workspace/recording_home_sources*.go、internal/api/handlers/studio_v2_recording_home_sources*.go、internal/api/router.go | Complete unchanged sources、Inconsistent or unavailable sources；共享 testdata/recording_home/source_projection.json 49 cases 包含 W1/D1/M1/P1/foreign D2 的 exact example；handler guards、production wiring、ReadSnapshot readonly regression | 本次 focused Go tests/vet/new-from-base lint 通過；TS parity 隨全 frontend 通過 |
| Return navigation reuses observed presentation without reusing authority | frontend/src/hooks/datalink/useRecordingHomeSources.ts、recordingHomeSourceCache.ts、RecordingHomeRoute.tsx、recordingHomePresentation.ts | Same workspace returns、Scope changes and old replies arrive、Refresh fails or sources mutate；source-summary、source-switching、source-summary-concurrent、source-mutation-workspace、source-modal tests 覆蓋 observed/stale、A→B→A、ignored abort、失敗、render abort 與 mutation owner | 本次 frontend 全套 655 files / 7,567 tests 通過；browser 未執行 |
| Home data demand stays bounded and complete | compact service/cache、route visible demand、recording-editor/recordingSourceDiscovery.ts | Large source fixture is projected、Source transport migration、Cache context isolation and eviction；route request count/25 detail ceiling、49 fixtures、QueryClient/backend/full base path 隔離、八 scope/第九 eviction、source mutation regressions | 本次 frontend 全套通過；20-device/2,000-point 實際 fixture 最終 guard recheck 為 2,336 bytes、canonical parity、零寫入 |

ReadSnapshot 既有實作只委派 Repository.Get；不新增 bootstrap accessor。Editor adapter 僅原樣重定位與 caller import/name 更新，完整來源 APIs 及既有 locator/await/readback authority 保留。Digest 僅表示 presentation observation，不是 atomic snapshot 或 command proof。

## 執行證據與限制

證據根目錄：`/workspace/shared/gateway-implementation-evidence/snapshots/`。命令與版本、fixture、RED/GREEN、量測原始資料詳見 non-browser-evidence.md。

- passed-current：`npm run lint`、`npm test -- --run`（655 files / 7,567 tests）、`npm run build`；frontend-{lint,test,build}.log 與 frontend-exit-statuses.txt。僅既有 chunk size warnings。
- passed-current：`go test ./internal/api ./internal/api/handlers ./internal/datalink/workspace ./cmd/test_ui -run 'TestRecordingHomeSources|TestReadSnapshot|TestProductionRecordingHomeSources' -count=1`。
- passed-current：`go vet ./internal/api ./internal/api/handlers ./internal/datalink/workspace ./cmd/test_ui`。
- passed-current：`golangci-lint run --timeout=10m --new-from-rev 0f404446 ./internal/api ./internal/api/handlers ./internal/datalink/workspace ./cmd/test_ui`，0 issues。上述三項為 backend-skips-final-{green,vet,lint}.log，覆蓋最終三個 skip guard 變更檔。
- passed-current：`git diff --check`、`make check-lines`。API docs 由官方 swag 模組 v1.16.6 執行 `make gen-docs` 產生，不手寫 generated files。
- prior full run：`go test -p4 -json ./...` 為 FAIL，14,762 pass / 523 skip / 1 fail；recordingplan concurrent restore failure 經 21 package variants／912 source+embed files 與原始 main byte identity、isolated 20-run 確認既有問題。`go vet -p4 ./...` PASS。其後三檔 guard 修補以 current affected checks 補足，不能冒稱完整 suite 已重跑。
- prior full lint：41 issues，39 位於與原始 main 相同的 29 repository files，2 位於 npm flatted 3.4.4 的 Go helper。沒有本 change 問題，不修改無關程式或 node_modules。完整證據在相鄰 final-backend 目錄。
- not-run：Mac browser / E2E browser acceptance、現場設備、external DB opt-in。跳過不等於通過；雲端 Chromium socket 限制不繞過、不擴權。
- 同一 SQLite fixture 的 paired timings 在最終 narrow guard 前量測：compact median 196.771 ms / 2,336 bytes；legacy source 417.249 ms / 2,394,466 bytes；legacy home 582.428 ms / 3,057,450 bytes。4,086 為 read-service calls，非 SQL statement count。最終 guard 單次重查仍 byte-identical、零 writes；378.953 ms 受並行 frontend load，不替代 paired comparison。

## 範圍與獨立檢閱

以 `0f404446976b1791b51327b2d5dd3b4b6ff12e8a` 為 explicit base，95 source entries 包含 rename/deletion；16 test paths 保留 preexisting_dirty 限制，path attribution 不證明逐 hunk ownership。完整 patches/provenance 與最終 snapshot check 記錄在 final-scope.json / final-scope-check.json；commit manifest 包含 rename 兩端。

獨立 report-only review 的 correctness、efficiency、simplification/reuse、conventions 四面向為 0 unresolved findings。Null/non-array skip Warning 已以十項 regression、四個預期失敗案例再轉綠修復。詳見 review-report.md。
