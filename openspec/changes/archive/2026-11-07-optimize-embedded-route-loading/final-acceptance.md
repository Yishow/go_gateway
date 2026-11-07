# 最終驗收與封存依據（2026-11-07）

## 版本與證據身分

- 最終 Mac 真瀏覽器來源：`0ffe74b57939bed5ec9fbee677585abe3639b85b`，執行前後工作樹與 index 乾淨。
- 證據包：`go_gateway-final-e2e-evidence-0ffe74b5.zip`，SHA256 `d2824fd3459665c7639dffe8ec499ede80f35fa0658c5c591c7d83e4672b561f`。
- 包內 `final-0ffe74b5/a/results.json` SHA256：`8950802b4d83de983a6712b9a49f08002604158a4a6e3a99761f5922b44e0c57`；run-record、source manifest、43 個修復檔案 hashes/modes 已逐項核對。
- Mac：macOS 27.2 arm64、Node 24.22.1、npm 12.2.1、Go 1.27.3、Playwright 1.64.1、專用 Chromium 156.0.8078.4。不是日常 Chrome，也不是雲端 replay。
- 本次收尾只清除 `AdvancedWorkbenchV2Page.tsx` 一個註解行尾空白與多餘 EOF 空白行。TypeScript 6 lint-compatible transpile（移除註解、無 source map）前後輸出與 diagnostics 完全相同；另以真正 production TypeScript 7/Vite build 通過核對。未聲稱 Mac 在這個後續文件／空白提交重跑。
- 歷史 evidence.md 的 browser pending、舊 task 數字與環境限制均保留為當時紀錄；當前完成依據以本檔及 tasks.md 為準。

## Spectra verify / review

官方 spxa 0.1.1 / core 3.0.0 已執行 verify 工作流、規格覆蓋、review 四個面向及 scope snapshot 核對。初始 scope 因原 touched tracking 未移入整合 checkout 而為 `approximated_base`；依四個 A 提交選出 63 個意圖相關路徑，檢視最終內容及後續重疊修改，不把路徑歸屬冒充每個 hunk 的獨占所有權。

四項要求、六個情境、兩個例子均有實作與測試，沒有以 test-scope 排除的情境。產品 review 無 Critical、Warning 或 Suggestion；唯一完整基準 diff-check 空白問題已修復。最後新增的任務追蹤只歸屬該空白檔，不追認其他歷史來源。

| 要求 | 實作與斷言 | 真正執行依據 |
| --- | --- | --- |
| 完整 embedded graph | Vite manifest/vendor groups、App routes、完整 static sync；vendor-groups／routing／delivery shell | Mac route smoke、production build；本輪 shell 21/21 |
| 延後 advanced/editor | 輕量入口、DeferredWorkbenchBoundary、外層 navigation/source owner；delay/reject/query/hash/dirty guards | Mac route delay/failure/history；單元完整批次與清理後 routing 5/5 |
| HTTP 與版本語意 | embed encoding/cache、GET/HEAD/Range/q=0/406、representation ETag、HTML revalidate、hash immutable、404、API/SSE 隔離 | Go web 測試與本輪全庫 Go；Mac 雙版本同 origin cache 驗收 |
| 分開傳輸與互動證據 | passive CDP/ResourceTiming、完整 graph、獨立故障注入 | Mac 5 cold + 5 warm raw ledger；原歷史同 fixture/toolchain build/HTTP 比較單獨保留 |

具體 vendor group 名稱與 `Accept-Encoding: identity;q=0, gzip;q=0` → CSS 406 的 Example 值均與測試一致。升級使用兩份明示的 cache-protocol variant，非 A-only 之前／之後效能比較。

## 實際 Mac A 結果

執行 `npx playwright test`，選取 `embedded-frontend-delivery.spec.ts`、`embedded-route-loading.spec.ts`、`embedded-cache-upgrade.spec.ts`、`embedded-loading-metrics.spec.ts`，參數 `--project=chromium --workers=1 --retries=0 --trace=on`。

22/22 PASS，0 failed/skipped/flaky/retry；66.619s 是套件執行時間，不是效能基準。量測為真 embedded fresh-empty SQLite（0 device/point），不代表大型來源、live reading 或現場吞吐量。

| 樣本 | cold frame / settled / actionable ms | warm frame / settled / actionable ms |
| --- | --- | --- |
| 1 | 76.2 / 80.0 / 110.0 | 45.0 / 48.3 / 76.7 |
| 2 | 76.0 / 79.6 / 112.4 | 44.6 / 47.8 / 79.5 |
| 3 | 76.5 / 80.4 / 107.8 | 45.2 / 48.3 / 79.4 |
| 4 | 76.2 / 80.4 / 107.4 | 44.3 / 47.4 / 76.0 |
| 5 | 77.1 / 81.4 / 112.3 | 44.7 / 47.8 / 77.4 |

- cold frame median/p95 76.2/77.1ms、settled 80.4/81.4、actionable 110.0/112.4；warm 為 44.7/45.2、47.8/48.3、77.4/79.5。這些是 automation 觀察上界，不是 paint 或設備讀值 latency。
- 每次 cold 全部 transfer 591,713 bytes，warm 9,344；各有 22 個 Performance entries，entry 不是網路重抓次數。
- 14 個 hashed assets cold transfer 579,481 bytes，warm 每個都是 0；9 個 JS decoded 1,401,447、encoded 451,625 bytes。Advanced chunk 均未在首頁抓取。
- 歷史 A 同工具／fixture：eager JS 2,467,102 → 1,393,373；Python gzip6 估算 695,754 → 432,389；Go HTTP identity 2,467,102 → gzip 436,624。它們與本次 Mac 數據分列，不拼接宣稱跨平台速度提升。
- 最終雲端 graph eager JS 1,401,447、同 gzip6 估算 435,538，仍低於歷史 graph 基準；Mac/Linux 的 9 個資產 hash 不同，沒有宣稱 bit-for-bit 相同。
- 原固定 cache payload 936,305 bytes、HeapAlloc 約 994KB 是當時 A build 的歷史測量，不是 final process RSS。實作依有限 embedded graph 建 cache，不接受 request-created keys。
- `cold-1.png` 至 `cold-5.png`、`warm-1.png` 至 `warm-5.png` 已在原 Mac `final-0ffe74b5/a/test-results/embedded-loading-metrics-E-9d192--retain-raw-browser-ledgers-chromium/` 唯讀核實保留；各 82,165 bytes、1280×798、SHA256 `7b495a714af28e82a25f38cd2c1b65064fe128c1cb11c2fefeed2074b59d5ec6`。它們是相同 rendered state 的十次截圖，不是十張不同畫面；便攜包未包含 test-results，完整 traces/videos/PNG 仍保留原 Mac。

## 最終品質檢查與界線

- passed-prior：同最終來源前端全量 681 files / 7,738 tests、全量 lint、production build、受影響 E2E types。執行者確認 set-e session exit0；frozen source/pre-post commit blob及43檔 manifest 一致，非從 log 時間猜測 SHA。
- passed-current：清理後 routing 5/5、檔案 ESLint、production build；完整基準 `git diff --check 0f54740e` 與 line gate。保留既有 >500kB chunk warning及行數提醒。
- passed-current：Go 1.25.5 Linux `go test -p 4 -json ./...`，60 packages、14,778 test/subtest PASS；524 test/subtest與5 package skip另列，沒有當成pass。`go vet -p 4 ./...` PASS。
- 全庫 `golangci-lint run --concurrency 4 ./...` 仍 FAIL：39項，涉及29檔均與原基準 `0f54740e` byte-identical；`--new-from-rev 0f54740e` 同全庫0 issues。首次 goimports cache 無法寫入已改用工作區可寫 XDG cache重跑，沒有改系統權限或抑制規則。
- 不宣稱全庫歷史 E2E types、Windows 實機、PLC/SCADA/HSL、外部真實 DB、部署或全球 WCAG 已驗收。沒有 push。

3.2 真瀏覽器條件與 3.4 最後核對已有完整證據；archive 只由官方核心預覽及交易套用規格一次，不手工同步 specs 或預先清除 tracking。
