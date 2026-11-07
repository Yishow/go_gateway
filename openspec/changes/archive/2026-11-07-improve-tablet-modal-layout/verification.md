# 驗證與審查結案

## 結論與版本

本案 7/7 任務完成，3 個 requirement 均有實作及測試對照。以 db91bbdb 為基準，驗證候選 66079373aacf4426bd81f34bfa464f531dd85a66。後續封存提交只改文件與同步規格，不改已驗證產品或測試來源。

本次範圍未發現未修復的 Critical／Warning／Suggestion。這不是整庫 E2E 全通過聲明，也不代表使用者已完成真機驗收。

## 執行證據

- passed-prior：雲端 `npm test -- --run`，683 files／7741 tests，`unit-release.log`、`unit-release.exit=0`。最後後續改動只有 E2E 尺寸斷言，單元來源相同。
- passed-prior：雲端 `npm run lint` 及 `npm run build`，`lint-translate.log`、`build-translate.log`、`lint-build-translate.exit=0`。其後只更新合法 E2E 尺寸斷言；候選正式 Mac 流程也成功編譯與執行。
- passed-prior：`bash scripts/check_file_lines.sh --base db91bbdb` 通過；ImportExportDialog.tsx 311 行、index.css 302 行為警告，未達 500 行阻擋門檻。
- passed-prior：Mac 原專案目錄、clean-source evidence gate，候選 66079373，225/225 無 skip。包括 memory 174、modal 38、dirty／cancel／history 13。執行結果由主實作者當輪確認，測試來源及 SHA 對應本次候選。
- passed-prior：Mac agent-browser，五視圖 × 390／768／1024／1440 × 明暗，共 40 組；另 1366 coarse touch 10 組使用 Playwright 真 hasTouch tap，不是 agent-browser。無橫向 overflow、footer 在界內，選取格 76px、footer 穩定、legacy 位移 0。主實作者實際抽查截圖無文字重疊。
- 五張原始參考圖已由主實作者查看，採構圖與資訊層級對照，不聲稱像素一致。768px 參考與長選取截圖亦已查看。

以上證據位於執行工作區 `/workspace/shared/tablet-ui-evidence`，正式瀏覽器來源由主實作者管理。報告保留來源層級，不冒稱本報告作者親自執行 Mac 操作。單元、建置與靜態審查不能替代真瀏覽器證據。

## 規格對照

| Requirement／scenario | 實作位置 | 檢查與結果 |
| --- | --- | --- |
| New configuration popups are centered viewport-sized modals；Opening a configuration popup | responsive-modal.css:3、ui/dialog.tsx、legacy primitives | responsive-modal.spec.ts:6 的六類 modal、non-modal 排除及 :34 旋轉／桌面回復；38 組正式 pass |
| Classified modals and contextual return state preserve safety；A mobile or tablet viewport opens a modal | responsive-modal.css、ImportExportDialog、MappingPayloadDialog | responsive-modal.spec.ts:6、:25、:34、:46、:64 覆蓋 390/768/1024、1366 coarse 邊界、低高度、鍵盤內捲、16px 文字；真實 viewport 過關 |
| A short confirmation is displayed | ConfirmDialog、Radix Dialog | responsive-modal.spec.ts confirm/radix-confirm、鍵盤讀完整內文、Escape；recording-home-memory-sizing.spec.ts:42 保留其他來源與桌面確認尺寸 |
| Returning after an item was removed | 既有 return-focus／identity 邏輯未改 | recording-modal-return.test.tsx、recording-home-return-scope.test.tsx、locator-return-after-source-save.test.tsx 的 removed／foreign／conflict 回歸在 7741 全單元內通過；另13個 dirty/cancel/history browser 通過 |
| Help is collapsed | 既有 ParameterHelp，非 modal 資訊排除 | parameter-help.test.tsx 指標／鍵盤／Escape 行為與不提交／不編輯保護，parameter-help-consumers.test.tsx 內容契約；全單元通過，未改 help 資料或風險內容 |
| Selected memory facts remain readable without expanding address rows；A long selected description is displayed | RecordingHomeMemoryCells、RecordingHomeMemoryView、兩份 memory CSS | recording-memory-tablet-selection.spec.ts:6 選取前後高度與 footer/toolbar 不移動、完整資訊可捲、五 layout；memory 單元及174組瀏覽器過關 |
| A continuation is visible after a window change | canonical memoryValueKey、每個 used physical cell 與 inspection 可見訂閱 | recording-memory-visuals.spec.ts:146 跨列／頁完整 physical focus 與 original key；:170 unknown/conflict；continuity、values 單元覆蓋 fail-closed；同源唯一值 owner 不變 |

本次 delta 沒有 `Example` 區塊，無逐列數值範例可遺漏。上述 scenarios 全有測試對照；未以單一表面樣式偏好增加無關斷言。精確 768×420、390×844、1024×768、1366 touch 及桌面 85% 尺寸屬跨 modal 契約，已檢查。原列表行動版橫向捲動產品修正未更動。

## 已修復而非未決問題

前輪 modal 34/38 的四個 legacy 失敗：最佳化後仍有 Tailwind 負百分比位移，1dea5a58 明確重設 --tw-translate-x/y:0；正式38/38通過。

前輪 memory 173/174：低高度舊斷言要求752×404，與新全螢幕契約矛盾；66079373 改為原點0、768×420，保留全部可達性檢查；正式174/174通過。

## 獨立 Spectra review

引用 `spectra-review-66079373.md`，scope snapshot `59ba38dc678f6fdd94c087b1c4829106cee6d647`，resolved／touched_tracking；25份產品／測試／DESIGN文字全部檢查，另6份 artifacts 作意圖脈絡。四個角度 correctness、efficiency、simplification/reuse、convention 均完成，0 Critical／0 Warning／0 Suggestion。審查者已用相同base確認snapshot current；此後未改產品來源。

## 限制與分開保留的失敗

- 更廣 guided-navigation 試跑為25pass、8fail、1interrupt、6notrun，未混入本案225。失敗屬 GD-N03/04/06/09 四類，主實作者在同一 Mac 原目錄切回 main db91bbdb、390 viewport 逐類重現，確認非本次引入；本案不擴範圍修復。
- 原生 iOS 軟鍵盤、動態安全區／實體 iPad與手機尚未驗證；Chrome viewport/touch context 不能取代真機。
- 未連真設備、真 DB 或驗證硬體／資料庫端到端；相關瀏覽器 fixture 為 UI-only mock sample，沒有宣稱真設備讀寫成功。
- 本案不修改 backend、協議、API、資料權威、ownership／ready token。未重跑後端，不宣稱後端已驗。

綜合目前有效證據，可封存本案與一次性套用兩份 delta specs；上述限制隨封存保留。

## Mac 原始證據定位

測試位置：`/Users/yishow/prj/go_gateway/frontend`。證據根目錄：`/Users/yishow/prj/go_gateway-tablet-validation-evidence`。

- clean fix/tablet-modal-layout@66079373：`npx playwright test tests/e2e/responsive-modal.spec.ts tests/e2e/recording-memory-tablet-selection.spec.ts --workers=2 --reporter=list --output=<root>/final38-output`；final38.log、final38.exit=0，38 pass。
- 同SHA：`npx playwright test "recording-home-memory-|recording-memory-visuals|recording-home-memory-slider|recording-picker-cancel-name|recording-home-snapshot-return" --workers=2 --reporter=list --output=<root>/final187-output`；final187.log、final187.exit=0，187 pass。只使用list reporter；frontend/playwright-report是舊報告，不能當本次HTML。
- clean main db91bbdb：`npx playwright test tests/e2e/recording-guided-navigation.spec.ts --grep "390.*GD-N0(3|4|6|9)" --workers=2 --reporter=list --output=<root>/baseline-guard-output`；baseline-guards.log／exit=1與error-context.md，四類均重現。候選擴展結果在main-guards.log。
- agent-main/run.sh、run.log、40張PNG為415bce56的五視圖390/768/1024/1440明暗click實測。415→660 memory產品來源未變，因此保留passed-prior；變更僅legacy translate與低高度測試。
- coarse1366-*十張PNG為Playwright hasTouch tap。最終660由agent-browser追加legacy768與長值focus768，圖為final660-legacy-768.png、final660-long-selected-768.png。不得把40+10統稱agent-browser50組。
