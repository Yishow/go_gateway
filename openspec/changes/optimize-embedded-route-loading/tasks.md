## 1. 建立可比較的交付基準

- [x] 1.1 完成分層量測與交付的基準：固定 SHA/toolchain/fixture，記錄未壓縮 bytes、gzip bytes、初始 import graph 與 HTTP headers；以 frontend production build 與資源清單檢查證明 Performance evidence separates transfer and interaction，真 browser 尚未跑明列 blocked/NOT RUN。
- [x] 1.2 [after: 1.1] 先補失敗回歸，驗證 Initial recording load excludes unused editors：延遲 advanced chunk 時穩定 reading frame 仍在，deep-link/dirty Back/Forward/import failure 不丟 intent；執行 frontend/tests/unit/app-routing-lazy-load.test.tsx 並保留修正前失敗證據。

## 2. 分批減少首次成本

- [x] 2.1 [after: 1.2] 依輕量首頁與按需編輯先新增 AdvancedWorkbenchV2Page 相容模組，遷移原 advanced wiring 而不改 owner lifetime；上述 routing tests 與 npm run build 證明每批可建置。
- [x] 2.2 [after: 2.1] 移除首頁對 advanced-only 模組的 eager imports，核對 Frontend build emits a complete embedded asset graph 及 stable logical vendor groups；完整 asset manifest、routing tests、npm run build 證明所有保留路由可達，退休 /studio 不復活。
- [x] 2.3 [after: 1.1] 先建立 Static representations preserve HTTP and version semantics 的失敗測試，internal/web/embed_delivery_test.go 以 GET/HEAD、gzip;q=0、identity;q=0＋Range、404/406、API/SSE 隔離鎖定 negotiation，記錄 RED。
- [x] 2.4 [after: 2.3] 實作靜態傳輸與版本快取中的固定資產壓縮與 negotiation，涵蓋 GET/HEAD/Range/MIME/406，不變更 API/SSE；go test ./internal/web 通過對應案例並量測固定 cache bytes。
- [x] 2.5 [after: 2.4] 先補 validators/304/HTML revalidation/hash immutable 的失敗案例再完成版本快取，go test ./internal/web 證明不同 representation 不共用錯誤 ETag，缺失資產仍404。

## 3. 交付驗收與收斂

- [x] 3.1 [after: 2.2, 2.5] 執行 tests/shell/embedded-frontend-delivery.sh 與 embedded build，確認完整同步、placeholder、stale removal、missing input fail-closed，且壓縮 cache 固定有界並記 RAM 成本；比對基準證明 eager JS 與 transfer bytes 均減少。
- [ ] 3.2 [after: 3.1] 依最終 Mac 瀏覽器驗收，在四案非browser實作完成後，統一於Mac實測 cold/warm/deep-link/chunk failure/版本更新及五次同條件數據，執行 frontend/tests/e2e/embedded-frontend-delivery.spec.ts；保留真實截圖與原始值，未執行保持未完成，不拿replay或歷史Mac證據代替。
- [x] 3.3 [after: 3.1] 完成本案非browser實作收斂，移除僅本案migration adapter並review scope；完成affected Go tests、frontend lint/test/build、git diff --check、make check-lines及Spectra verify/review，修復可成立的發現，明列browser延後及未跑的跨平台風險。
- [ ] 3.4 [after: 3.2, 3.3, 3.5] 最終Mac驗收後對本案scope重新verify/review，修復並重跑受影響測試，確認browser證據與最終程式一致，全部任務實證完成後才進行archive；非browser品質通過的相關變更可先獨立commit並保留未完成驗收，不push或部署。

- [x] 3.5 [after: 3.3] 補齊 Static representations preserve HTTP and version semantics 的非瀏覽器錯誤路徑：任何預壓縮或 validator 建表讀取／走訪失敗，捨棄不完整衍生快取，保留可讀 identity 與 identity 被禁止時的406；先以 TestEmbeddedCacheFailureDiscardsPartialState 重現，再修復並通過 go test／go vet／golangci-lint ./internal/web、交付 shell tests、Spectra verify/review，不擴張 API／SSE 或 browser 範圍。

任務切片以每項不超過一小時為規劃估算，不是工期承諾；若實作探索證明超過，先把同一行為拆成可獨立測試的垂直批次並保留 dependencies，不擴張 scope。Checkbox 依本次實測更新；browser 阻擋與依賴任務維持未完成。

變更註記：依2026-11-06使用者最新指示，原3.2的「雲端browser／不使用Mac」限制被最終Mac驗收取代；原3.3的非browser品質工作提前獨立執行，browser後最終收斂保留於3.4。既有8個已完成任務描述與來源歸屬未變。

後續指示：「一樣照此實作，archive時先保留，最後一起上Mac,更新你的記憶，之後都依此方式，避免停止實作」。因此各案非browser實作及verify/review/fix完成後可先做範圍限定commit，archive與真browser驗收延至四案實作完成後；不因雲端browser限制停止後續實作。

非瀏覽器補查註記：官方 golangci-lint 補齊後發現建表錯誤處理與靜態寫法問題，新增3.5修復及回歸，保留所有已完成任務原描述；browser及archive仍待最終Mac驗收。
