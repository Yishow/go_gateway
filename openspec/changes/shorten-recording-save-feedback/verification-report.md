# Verify：shorten-recording-save-feedback

## 摘要

| 面向 | 結果 |
| --- | --- |
| Completeness | 9/11 tasks完成；browser/archive兩項仍未完成 |
| Correctness | 3/3 requirements 有對應實作，6/6 scenarios 有測試，1/1 example 使用原值 |
| Coherence | config-only ACK/refresh 邊界、既有 reducer 與 consequential waits 保留 |

## 範圍與證據層級

Base與HEAD為 `406490c8ff1f801c27ddc491d9f29debbd12d523`；source scope為 `touched_tracking`／`resolved`，51個文字檔，未包含既有 `DESIGN.md`。已檢查 source snapshot `998030ac80ea03b95498bc3955f168aab1b3ba57`；之後只新增本 change 的驗證報告，重新 capture/check 的最後 identity保存在雲端交接。測試範圍有1610檔 source/test/config SHA256 manifest，完成後1610檔全部比對一致。

`database-save-feedback.test.tsx` 有 preexisting-dirty 警告：後續 task baseline不能單獨證明每個 hunk。它在B HEAD不存在，完整新增內容已由C各項 explicit-file provenance與獨立 reviewer 檢查；提交計畫仍顯示完整內容。

測試存在與實際執行分開記錄。Browser案例已寫好但未執行，不得當成通過。

## Requirement／Scenario／Example 對應

### Configuration save acknowledgement is separate from refresh

- 實作：`useStudioV2WorkspaceDatabase.ts:101-119`、`recordingConfigRefresh.ts:46-97`、`useStudioV2DatabaseAutosave.ts:234-303`。
- Valid acknowledgement precedes refresh：`database-save-feedback.test.tsx` 的 held-three-read及synthetic案例；mutation保存成功仍保留 runtime apply_failed，不把它當啟用或SQL證據。
- Acknowledgement lacks confirmation evidence：同檔缺workspace／identity/setup revisions及unknown expected owner案例；`database-save-feedback-owner.test.tsx`與browser malformed ACK案例。
- Example deferred refreshes：精確ACK1ms、三個各200ms，下一事件迴圈同時開始，ACK不等待600ms鏈，完成前refreshing。將刷新await回mutation的mutation check會失敗。
- 執行：focused passed-current；完整662檔／7600項 passed-current。Browser not-run。

### Background refresh remains scoped and observable

- 實作：`recordingConfigRefresh.ts` scoped generation/cancellation/coalescing與error projection；`useStudioV2Workspace.ts`、config/target hooks及services傳遞captured baseURL/AbortSignal；`recordingHomeSourceCache.ts:68-83`保持current owner authority。
- Refresh fails after confirmed save：`database-save-feedback-presentation.test.tsx`折疊摘要可見安全error/retry；`database-save-feedback-page.test.tsx`已確認同owner retained shell，initial malformed bootstrap仍error；`database-save-feedback.test.tsx`保存成功不變失敗、唯讀retry、外部GET恢復需三類證據。
- New edit or workspace wins：`database-save-refresh-scope.test.tsx` late GET、PUT期間GET、origin/workspace分離、generation/coalescing；`database-save-feedback-owner.test.tsx`one-way／A→B→A／origin／unmount；既有database-connector-input-race保留更新文字與unknown不重送。
- 執行：focused passed-current；獨立review recovery重現2/2 passed-current；最終完整662檔／7600項 passed-current。Browser not-run。

### Save acceleration preserves uncertain and consequential operation barriers

- 實作：僅named config mutation更換completion；target及`useInvalidateWriteGroups`等helper原等待不變；Step4使用既有readonly／action guards；autosave保留unknown intent、CAS與reducer acknowledgement。
- Timeout conflict or repeated action：既有database-autosave-page、database-save-queue、database-connector-input-race、activation-barrier、reload-recovery測試，加新scoped owner與ACK資料驗證。無ACK不自動補送。
- Another operation uses shared helpers：`database-save-feedback-barriers.test.tsx`明確held兩個讀取，target與write-group promise仍pending；移除await的mutation check失敗，原helper無產品diff。既有locator及draft-guard回歸保持。
- Close/Cancel/Back/Forward：production guard unit回歸及20個真實router Playwright authoring覆蓋；database autosave原本沒有獨立navigation guard，未在本change新增或假稱已有。
- 執行：affected regression178/178、provider44/44、MC owner16/16；最終完整662檔／7600項 passed-current。Browser not-run。

全部6個scenario與1個example均有可檢查測試；沒有test-scope exclusion或未覆蓋scenario。這不等於所有browser斷言已執行。

## Coherence 與 review

- 沿用B home API及唯一query authority；generic owner fallback更新有active B／inactive A與failed-current回歸。
- scoped read adapter不改API schema；target、write-group、schema、activation、test-write無縮短成功界線。
- 保存文案與en/zh-TW更新一致；raw exception／DSN不加入新的operator summary。
- 獨立review的sticky recovery與owner round-trip問題已修復並重測，見`review-report.md`。目前無未解決程式缺陷。

## 執行清單

| 命令／驗收 | 狀態 |
| --- | --- |
| `npm run build` | passed-current，包含TypeScript；保留既有chunk size提示 |
| `npm run lint` | passed-current |
| `npm test -- --run --maxWorkers=2` | passed-current，662檔／7600項，527.93秒，exit0；未更改default timeout或斷言 |
| `git diff --check`、`make check-lines` | passed-current；包含最終報告 |
| Playwright `--list` | passed-current，20個case；僅discovery，未啟動browser |
| fixture Vitest、E2E-inclusive typecheck | passed-current，fixture6/6包含於完整unit suite；最終E2E-inclusive typecheck通過 |
| Mac browser／screenshots／SQL／field | not-run |
| Archive／push／部署 | not-run，不在雲端階段執行 |

無後端修改，沒有重跑完整Go；B的既有Go／lint問題不作本change已修復宣稱。

## CRITICAL：保留的未完成項目

1. `tasks.md`3.1：最終Mac browser合併驗收未執行。依使用者順序，全部change實作後執行20個authored案例、記錄ACK/refresh/readback實際時間、截圖及必要修復。
2. `tasks.md`3.3：最終browser證據、必要修復與verify/review後才archive；目前不得封存。

這些Completeness事項不等於雲端實作可以當作browser驗收完成。最終Mac驗收後重跑verify/review，再決定archive。
