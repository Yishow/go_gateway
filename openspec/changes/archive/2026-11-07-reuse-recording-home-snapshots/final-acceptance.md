# 最終驗收與封存依據（2026-11-07）

## 版本與證據身分

Mac 最終來源為 `0ffe74b57939bed5ec9fbee677585abe3639b85b`，執行前後 clean。證據包 `go_gateway-final-e2e-evidence-0ffe74b5.zip`，SHA256 `d2824fd3459665c7639dffe8ec499ede80f35fa0658c5c591c7d83e4672b561f`；包內 `final-0ffe74b5/b/` 保留原始 results、run-record、exit與ledger，source manifest43檔hash/mode均與整合程式一致。之後 A 只清兩處註解／EOF空白，轉譯一致且路由、lint/build重驗；B來源沒有後改。

Mac macOS27.2 arm64、Node24.22.1、npm12.2.1、Go1.27.3、Playwright1.64.1、專用Chromium156 headless shell。此文件記錄實際執行結果的身分核對，並非在雲端重跑browser。舊 non-browser-evidence、browser-harness-evidence、verification/review中的 pending 為歷史checkpoint；當前結論以本檔與tasks為準。

## 規格與程式審查

官方 verify/review 工作流已完成。原 touched tracking 未帶入整合checkout，scope `approximated_base` 只代表全批次近似範圍；按 B 提交及最終共用修改選取意圖相關內容，不假造獨占hunk歸屬。三個 requirements、八個 scenarios、一個 mapping owner example 已逐一對照實作、斷言與真正執行，沒有未覆蓋規格行為。正確性、效率、重用／簡化與專案慣例 review 無未解產品發現。

| 要求 | 實作／斷言 | 執行依據 |
| --- | --- | --- |
| Canonical compact evidence | workspace ReadSnapshot、完整兩輪source projection、typed400/404/503、無敏感設定／讀值DTO、49組Go/TS parity、owner mismatch負例 | Go測試與全庫test；真embedded canonical digest、20devices/2000points、49表hash |
| 返回只重用呈現，不重用權限 | QueryClient/backend/full-baseURL/workspace/epoch隔離；成功mutation捕捉與失效；當前locator proof仍驗證；返回focus/scroll與visible target | 60/1000groups routed editor/history、409/503、兩次GET-only retry、A→B→A、頁面／篩選／焦點斷言 |
| 需求有界且完整 | home一次compactGET，visible-page lifecycle/delivery上限；editor保留舊明確discovery；八個scope容量與第九eviction | 每次25 lifecycle+25delivery；零off-page／duplicate；modalClose零額外compactGET；unit cache隔離／eviction |

B readiness修復將觀察與明確檢查寫入分開，保留原啟用／scope保護。大LIST parser支援完整1000項。返回焦點只對當前visit目標做可見性觀察、只調document，不改集合scroll；使用者操作或讀取完成即停止，不新增全頁位置快取。

## 實際 Mac B11

命令：`npx playwright test tests/e2e/recording-home-snapshot-return.spec.ts tests/e2e/recording-home-snapshot-embedded.spec.ts --project=chromium --workers=1 --retries=0 --trace=on`。11/11 PASS，0 failed/skipped/retry/flaky；48.007s只是suite時間。

- 10個routed案例用明示mock API、真 `/studio/v2` router；1個embedded用真正Go/frontend與owned可丟棄SQLite。
- 60/1000group cases保留四種篩選、第2頁、集合scroll差≤1px、原焦點與viewport完整可見，含真Back/Forward及延後layout位移。
- 新home visit恰一個compactGET；普通source modal Close/Cancel/Escape仍同mount，零新增compactGET；成功source save獨立觸發一次刷新，Close不重複。
- 兩種A→B→A晚回覆release順序，實際browser均為cancelled；忽略abort但仍delivered的fencing證明另由unit提供，不把取消誤稱late-delivered。
- 失敗retry只GET、不送mutation；409/503當前locator proof拒絕不提升compact為操作權限。

真embedded以正式API seed20 disabled draft devices／20 disabled rules／2000 disabled points，不啟動test/connect/poll/activate，也沒有externalDB或PLC。

| 區段 | API requests | compact GET | compact body／transfer bytes | duration ms |
| --- | ---: | ---: | ---: | ---: |
| cold home | 7 | 1 | 2306 / 2606 | 102.0 |
| warm editor return | 5 | 1 | 2306 / 2606 | 134.5 |
| warm history forward | 5 | 1 | 2306 / 2606 | 135.9 |
| ordinary source modal Close | 0 | 0 | 0 / 0 | 不適用 |

每區段GET-only、全部49表row hashes相同、protocol sentinel0。完整request/ResourceTiming ledger在原results內，不把compact payload等同全頁流量，也不把一次duration稱穩定效能比較。

歷史20/2000 fixture的compact2336 bytes、舊14GET3057450 bytes及canonical parity保留於non-browser-evidence.md；它是當時同fixture/toolchain的HTTP測量，不是本次final2306bytes同一次對照，不跨平台拼接宣稱固定速度提升。表hash相同證明無持久化差異，不代表偵測所有瞬時或回滾write。

## 真 production guard補證

同最終SHA另一個真embedded host的36個native tab stops正反循環、dirty Close/Keep/Discard、重複Escape、memory Escape與connection Cancel、unknown及真Back/Forward均通過。唯一synthetic PUT在client端先abort，未進backend；read-current不能證明原PUT成功。原始guard events/HAR與49表before/after保留於包內 `manual-final-0ffe74b5/`，sentinel0、owned清理完成。

## 最終品質與界線

- passed-prior：同凍結final來源full frontend681files/7738tests、full lint、production build、受影響E2Etypes；執行者set-e sessionexit0與pre/postblob、43檔manifest匹配。
- passed-current：Go1.25.5 Linux `go test -p 4 -json ./...`，60packages、14778test/subtest PASS；524tests與5packages skip分列。`go vet -p 4 ./...` PASS。
- 全庫lint39項既存問題涉及29檔，均與原0f54740e byte-identical；new-from-rev0f54740e同全庫0issues。沒有壓掉規則或擴張修復其他功能。
- 完整基準diff-check、line gate及本案validate通過。沒有把mock、disabledfixture或單元測試等同外部DB delivery、PLC/SCADA/HSL現場驗收；Windows實機及全庫歷史E2Etypes不是本次通過聲明。

3.1與3.3已有完成依據；官方archive將只套用本案delta一次。未push／部署，其他active changes保留。
