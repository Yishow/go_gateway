# 最終 Mac 驗收腳本：雲端交付

本次基準 c296c2cdb7ab5db5e673c9fc6538cdd9c6326ffb。只新增15個測試／fixture／helper檔案；產品、設定與A owned server helper不變。4.1/4.2/4.3是authoring，3.1真Mac與3.3最終驗收/archive仍未完成。

## 本次已執行

- 7 files / 67 unit tests PASS：18 routed fixture/ledger、39 seed/SQLite/sentinel、10既有Aownership。
- Focused ESLint、明確包含E2E的TypeScript、Playwright discovery11cases/2files、diff/line、Spectra analyze/validate通過；五項既有example Suggestions非阻擋。
- 真production build/start/seed/HTTP/SQLite smoke通過：20draft devices、20disabled rules、2000disabled/unread points，0mapping/tag/connector，owned TCPsentinel0connections。兩次compact讀取各2306JSONbytes，canonical20devices與空candidates/excluded相等、digest穩定、全部49tablehash一致。
- 初始19unsafe-source/origin cases先失敗再修復。合法aligned ACK與hydrated-ready但disabled listener的harness假設有獨立regression修正；沒有產品activation問題。
- 獨立review修復late-response斷言過早問題：先等terminaldelivery再assertA2/fresh，最後再次檢查；regression先失敗再通過。SQLite失敗時保留before/afterevidence亦有regression。最終四面向review0未解findings。

## 已撰寫但未執行的 browser assertions

- 真/studio/v2 routed frontend配受控API：60/1000groups、filters/page2/scroll/focus、Back/Forward、同scope返回1compactGET、source modal Close/Cancel/Escape0新增GET、成功ruleACK獨立1revalidation。
- Locator409/503拒絕操作，freshcompact不授權write；503→unknown→fresh的唯讀retry；home禁止raw來源transport，editor保留既有discovery；A/B兩種late-release順序保留finished/cancelled/not-delivered區分。
- 真embeddedfrontend+productionAPI：cold與warmreturns、modalClose、passiveCDP/ResourceTiming、screenshots、configurationtablehash。沒有requestinterception；legacyeditor/modalGET可能改device timestamps，獨立記錄不混入homeinterval。

Discovery不是browser通過。真history/focus/scroll/counts、browser timing與screenshots全部NOT RUN。EqualSQLitehash只證明無persisted差異，不證明無transient/rolledback/write-restore。2306bytes不可當作同機before/after latency。Fixture為disabled/unmapped，不替代原49個mapped/identityparitycases。

## Mac 指令

在已commit最終SHA的Mac frontend，以專用PlaywrightChromium執行。需Go1.25.5與Node內建node:sqlite（22.13+，建議已驗證24.22.1）：

```sh
npm ci
npx playwright install chromium
npx playwright test tests/e2e/recording-home-snapshot-return.spec.ts tests/e2e/recording-home-snapshot-embedded.spec.ts --project=chromium --workers=1 --retries=0 --trace=on
```

遇權限或capability限制停止，不擴權。Seed只用verifiedownedorigin的正式device/ruleAPIs，不呼叫test/activate/poll或外部DB/PLC，不接日常瀏覽器。先保留trace/ledger與真失敗，再做scope內修復/review。

## 證據定位

雲端 evidence root：`/workspace/shared/gateway-implementation-evidence/snapshots/browser/`。
- aggregate-{unit,lint,types,discovery}.log；focused-types.json；product-fingerprint-comparison.json（4084files不變）。
- routed/handoff.md、touched-manifest.json與RED/GREENlogs。
- embedded/implementation-report.md、final-files.sha256、http-sqlite-smoke.json、build-identity.json與attempt-specificlogs。
- commit-plan/含完整候選內容/diff/provenance/verification/review/Maccommands；正式snapshot以其plan.json與scope-check.json為準。

完整frontendD667files7623tests是先前結果，此次未重跑。完整Go既有test/lint failures與全E2Etypes六項mixed-review-host既有錯誤保持明列，不改無關程式、不宣稱allgreen。沒有commit/push/archive由本worker執行。
