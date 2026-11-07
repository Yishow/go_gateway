# 最終視覺驗收與封存依據（2026-11-07）

## 身分與核對範圍

最終Mac來源 `0ffe74b57939bed5ec9fbee677585abe3639b85b`，clean before/after。證據包 `go_gateway-final-e2e-evidence-0ffe74b5.zip`，SHA256 `d2824fd3459665c7639dffe8ec499ede80f35fa0658c5c591c7d83e4672b561f`。包內 `final-0ffe74b5/d/results.json` SHA256 `72c4b937235af9f8495f5849a0fd66c760cd63278a537eb05d7302c80261acfd`；全部43sourcehash、90finalartifacts、5原圖fullSHA/尺寸、30基礎圖metadata/尺寸及39visual結果身分皆吻合。

Mac macOS27.2 arm64、Node24.22.1、npm12.2.1、Playwright1.64.1；agent-browser0.41.2以專用Chromium156 headless shell執行，非日常Chrome/profile。後續A只有兩處註解／EOF空白清理，轉譯一致並重驗路由/lint/build；D來源未後改。

官方verify/review對照4requirements、8scenarios。7個有測試；「Browser execution is blocked」是驗收流程契約，依test-scope排除產品unit斷言，過去blocked與此次真Mac執行分列，不是未覆蓋產品分支。Spec無正式Example，不捏造數量。正確性、效率、重用簡化、慣例review無未解產品finding；validate零errors/warnings。

原tracking未移入整合checkout，scope為approximated_base；以4個D提交選40路徑並檢視最終共用hunks，不將路徑清單說成獨占所有權。原reference-manifest與歷史non-browser/review/verify保留；本檔和tasks取代其當前pending結論。

## 實際 D172 與手動檢查

原始命令保留於包內 `final-0ffe74b5/d/run-record.txt`，選取 `recording-memory-visuals.spec.ts` 及 `recording-home-memory-.*\.spec\.ts`，`--project=chromium --workers=1 --retries=0 --trace=on`。172/172 PASS，零failed/skipped/retry/flaky，284.852s為suite時間。

| 套件 | 真執行案例 |
| --- | ---: |
| browser regressions / density / grid | 4 / 20 / 5 |
| hover / 320px layout / navigation | 12 / 1 / 30 |
| presentations / recovery / sizing | 15 / 30 / 5 |
| native slider / routed memory view / new visuals | 6 / 5 / 39 |

- agent-browser真擷取30圖：五版×light/dark×1440×1000、768×1000、390×844。記錄者逐張直接對照五張原圖；review以相同hash重用完整比較並抽查像素，不冒稱第二次全38張人工比較。
- 8補圖含phone Map/Detail底部、unknown hatch加keyboard focus、light/dark中點native slider。底部摘要、inspector、完整focus outline與thumb可見，沒有clipping。
- Grid full-width；Navigator left／phone above；Map完整256slots上條與下摘要；Chapters全寬章節帶在grid上方；Detail desktop右／phone下。五版不是同面板換標題。
- Quiet mineral取代原圖漸層／glow是DESIGN明示的有意差異；原圖示範數值、計數、健康狀態不作產品authority。全部30圖有MOCK SAMPLE，不代表真HSL／設備／SQL讀值。
- 768px Chapters light/dark仍有完整兩列與樣本警告；未要求phone Chapters也有兩列。Phone Map每個occupancy segment仍可分辨。
- 30圖3188個已採樣文字對比最低4.7776:1；30個control boundary最低5.5786:1；focus最低4.9335:1。額外native slider六個viewport/theme真paint boundary light5.578596772:1、dark4.730172457:1；沒有宣稱全站WCAG或未採樣native所有狀態通過。
- genuine touch slider、0/127/255→W0/W7F/WFF、33項long chapter、long-name/large-value、unknown/conflict、4-register跨列／250slot跨頁、low-height768×420與phone44px/14px門檻都有已執行斷言。

## 安全與共用context

| 要求 | 程式及驗收 |
| --- | --- |
| distinct structures | panels/layouts/CSS單一DESIGN權威；30matrix與presentation/long-list/256slot測試 |
| physical meaning/safety | canonicalmap不改；每physicalcell可focus；join/continuation依真columns；partial不洩漏不可信string reading；unknown無假成功動作；2/4register與跨page測試 |
| shared interaction context | 一份grid與parent live hook；area/window/focus/owner共用；hover不移tabstop，切layout保持context；toolbarblur不吞click；30navigation/12hover/30recovery |
| current original evidence | 原圖fullSHA、每圖fixture/theme/viewport與finalSHAclean guard；全矩陣與8補圖真pixel檢查；blocked歷史未當passed |

真production embedded host補驗位於 `manual-final-0ffe74b5/`：36個native Tab stops正反循環、wrap不離dialog、dirty Keep/Discard、重複Escape、memory Escape／connection Cancel回原trigger、unknown護欄及真Back/Forward。唯一syntheticPUT在client先abort，未送backend；read-current不能證實原PUT成功。49SQLite表rowhash前後完全相同、sentinel0，owned服務／專用browser已關閉。

## 品質及界線

- passed-prior：同final來源full frontend681files/7738tests、full lint、productionbuild與affectedE2Etypes；set-e sessionexit0、frozen source與pre/postblob/43hash核對一致。
- passed-current：本次artifact/source/pixel evidence身份核對、官方validate/scope、文件diff與linegate。歷史完整E2Etypecheck債不冒稱全綠。
- 全批Go1.25.5 Linux fulltest60packages/14778tests PASS，524tests/5packages skip另列；vetPASS。全庫lint39既存問題涉及29個原baseline相同檔，new-from-base0issues。
- 本change沒有修改API/runtime或全域theme，原add-switchable-memory-views完成紀錄保留。沒有真PLC/SCADA/HSL、外部DB/SQLdelivery、Windows實機或部署驗收聲明。

3.1–3.7與4.2均有有效完成證據；DESIGN僅更新驗收狀態，不改色彩／佈局契約。官方archive將套用本案delta一次；沒有push。
