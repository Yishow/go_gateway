# 最終驗收與封存依據（2026-11-07）

## 版本與完整性

最終Mac來源 `0ffe74b57939bed5ec9fbee677585abe3639b85b`，before/after clean。`go_gateway-final-e2e-evidence-0ffe74b5.zip` SHA256 `d2824fd3459665c7639dffe8ec499ede80f35fa0658c5c591c7d83e4672b561f`；包內final-0ffe74b5/c的results/run-record/exit與20份HTTP ledger、43個來源hash及138個包內容hash均已核對。A之後只清除兩處註解／EOF空白，有轉譯一致及新routing/lint/build證據；C程式未再修改。

Mac macOS27.2 arm64、Node24.22.1、npm12.2.1、Playwright1.64.1、專用Chromium156 headless shell。使用真router加明示mock API，不是外部真DB保存或SQLreadback測量。歷史verification/review與non-browser-evidence的pending保留，當前結果以本檔及tasks為準。

官方Spectra verify/review完成，3個requirements、6個scenarios、1個精確example都有實作和測試，無排除或未覆蓋缺口。原tracking未移入整合checkout，scope為approximated_base；依C兩個提交與必要重疊caller檢視60路徑，沒有假稱所有hunk獨占歸C。正確性、效率、簡化重用、慣例四面review無未解產品finding，validate有效且無errors/warnings。

## 要求與執行對照

| 要求 | 實作／斷言 | 實際驗收 |
| --- | --- | --- |
| ACK與refresh分離 | named config mutation捕捉owner/baseURL；ACK須workspace/connector/identity/setup revisions完整；原reducer保留較新輸入 | 有效ACK已saved且三GET仍held；六種缺證據ACK保守失敗；1msACK/200ms三讀的精確unit example |
| refresh有scope且可觀察 | QueryClient metadata、同步取消與失效、generation/owner fence、三類並行讀、same-key合併 | config/targets/bootstrap逐項失敗，saved不變失敗；GET-only retry、較新編輯、跨workspace及A→B→A／origin單元 |
| unknown與後果性操作保護 | retry:false、typed409 vsunknown、原queue/reducer；target/writegroup/schema/lifecycle等待保留 | timeout/lostACK不重送、409明確再編輯才送；原recording dirty Close/Escape/Back/Forward保護 |

精確example使用fake clock ACK1ms、三GET各200ms，acknowledgedAt=1，reads started=[2,2,2]；不是以近似時間替代。Configuration saved不等於activation、SQLcommitted、readbackverified或cleanup。

## 實際 Mac C20

`npx playwright test tests/e2e/recording-database-save-feedback.spec.ts --project=chromium --workers=1 --retries=0 --trace=on`：20/20 PASS，0failed/skipped/retry/flaky，29.610s是suite時間。

| 從submit起算 | config→targets→bootstrap | bootstrap→targets→config |
| --- | ---: | ---: |
| fixture ACK fulfill | 10ms | 9ms |
| browser PUT response | 17ms | 15ms |
| 三GET開始 | 17ms | 16ms |
| saved且三GET仍held的UI斷言 | 34ms | 32ms |
| configuration GET釋放完成 | 174ms | 193ms |
| targets GET釋放完成 | 228ms | 163ms |
| bootstrap GET釋放完成 | 255ms | 161ms |

人工held reads只用於建立事件順序；以上不是效能基準或外部SQLreadback時間。Refresh完成事件是fixture fulfill後記錄，不是精確paint timestamp。

- 三種refresh fail後retry均只有1PUT；409只在明確新編輯後第2PUT；timeout/lostACK各1PUT，不自動replay。
- 較新intent以新revision送出，舊GET不覆寫；workspace B不被A晚read取代。
- unknown防replay是既有mounted owner界線；重開GET觀察server值不證明失落ACK成功，亦不代表未送新文字已保存。
- 保留的是原recording-draft navigation guard，沒有偷增database-only guard或durable未知意圖資料層。Close/Escape Keep/Discard與真browserBack/Forward斷言均執行。
- 真embedded人工補驗另有client-before-backend abortedPUT、36tabfocus與unknown護欄，49表hash不變／sentinel0；此補證不冒充C保存mock為真SQL操作。

## 品質與交接

- passed-prior：同finalfrozen來源full frontend681files/7738tests、full lint/build、affectedE2Etypes，執行set-e sessionexit0與pre/postblob及43hash一致。
- passed-current：收尾完整基準diff-check、linegate、validate；A純空白routing5/5、lint/build再驗。
- 本批後端全庫Go1.25.5 test60packages/14778tests PASS，524tests/5packages skip分列；vetPASS。全庫lint39既存項/29檔與baseline完全同bytes，new-from-base0。C不修改後端，該結果不是外部DB或設備成功聲明。
- Full歷史E2Etypecheck、Windows實機、現場PLC/SCADA/HSL、真外部DB、SQLreadback/cleanup與部署未宣稱通過。

3.1與3.3的browser、必要修復及最終verify/review都已完成；接續官方archive套用delta一次與限定commit，不push。歷史失敗與中途證據不覆寫。
