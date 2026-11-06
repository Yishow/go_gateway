## 1. 共用設計與語意

- [x] 1.1 在 optimize-embedded-route-loading 的 owner 邊界穩定後核對五張原圖 cache，落實單一設計權威與有意取捨：DESIGN.md 的 recipes/85% exception 與實際 CSS 對照不漂移；Visual acceptance uses current original-reference evidence 的 reference ID／SHA／viewport／fixture 矩陣可追溯。
- [x] 1.2 [after: 1.1] 先補跨列／頁、unknown/conflict、完整 datatype 與狀態優先的失敗測試，再依共享卡片與分明構圖校正共用 cells；驗證 Typed cells preserve physical address meaning and safety，執行 memory-continuity/map/scope/recovery/view tests。
- [x] 1.3 [after: 1.2] 校正 grid/navigator 的層級、範圍 identity 與真選取，滿足 Memory views expose distinct inspection structures；memory-layouts/panels tests 及對照清單驗證不是同一白面板加標題。

## 2. 分批校正其餘構圖

- [x] 2.1 [after: 1.3] 校正 map 的完整 window strip／緊湊摘要，256-slot 與 slider 精確定位不退化；memory-panels/layouts/keyboard tests 證明無裝飾抽樣與選址偏移。
- [x] 2.2 [after: 2.1] 將 chapters 校正為 grid 上方全寬章節帶，保留 bounded 長清單及所有 range facts 可達；memory-panels/layouts 測試驗證展開／定位／long-list，不捏造健康 badge。
- [x] 2.3 [after: 2.2] 校正 detail inspector、手機順序與長值可達性，保持 Layout changes retain one shared interaction context；memory-hover/keyboard/observer-states/values tests 證明 hover 不移 tab stop、切換保持 context 且無第二 hook。

## 3. 真實視覺與互動驗收

- [ ] 3.1 [after: 2.3] 依逐圖與互動雙重驗收，於四案雲端實作與測試撰寫完成後，在最終合併 Mac browser 擷取與比對 grid 的 light/dark × 1440/768/390 六個基礎狀態，記 reference/SHA/viewport/fixture 與 intentional difference/defect；驗證 DESIGN.md 構圖，若環境阻擋保持未完成、不拿歷史證據替代。
- [ ] 3.2 [after: 3.1] 依逐圖與互動雙重驗收，於四案雲端實作與測試撰寫完成後，在最終合併 Mac browser 擷取與比對 navigator 的 light/dark × 1440/768/390 六個基礎狀態，記 reference/SHA/viewport/fixture 與 intentional difference/defect；驗證 DESIGN.md 構圖，若環境阻擋保持未完成、不拿歷史證據替代。
- [ ] 3.3 [after: 3.2] 依逐圖與互動雙重驗收，於四案雲端實作與測試撰寫完成後，在最終合併 Mac browser 擷取與比對 map 的 light/dark × 1440/768/390 六個基礎狀態，記 reference/SHA/viewport/fixture 與 intentional difference/defect；驗證 DESIGN.md 構圖，若環境阻擋保持未完成、不拿歷史證據替代。
- [ ] 3.4 [after: 3.3] 依逐圖與互動雙重驗收，於四案雲端實作與測試撰寫完成後，在最終合併 Mac browser 擷取與比對 chapters 的 light/dark × 1440/768/390 六個基礎狀態，記 reference/SHA/viewport/fixture 與 intentional difference/defect；驗證 DESIGN.md 構圖，若環境阻擋保持未完成、不拿歷史證據替代。
- [ ] 3.5 [after: 3.4] 依逐圖與互動雙重驗收，於四案雲端實作與測試撰寫完成後，在最終合併 Mac browser 擷取與比對 detail 的 light/dark × 1440/768/390 六個基礎狀態，記 reference/SHA/viewport/fixture 與 intentional difference/defect；驗證 DESIGN.md 構圖，若環境阻擋保持未完成、不拿歷史證據替代。
- [ ] 3.6 [after: 3.5] 補 low-height、長名稱／大值、unknown/conflict、四格跨列／頁視覺情境，實測對比、85% sizing 與 safe-area；recording-home-memory-sizing/presentations specs 與差異清單證明控制及完整資訊可達。
- [ ] 3.7 [after: 3.6] 執行 recording-memory-visuals.spec.ts 與既有 navigation/hover/recovery matrix，覆蓋切換、觸控、keyboard、Close/Cancel、Back/Forward、重複點擊與 dirty guard；修正後重跑受影響畫面／互動。
- [x] 3.8 [after: 2.3] 完成 frontend lint/test/build、git diff --check、make check-lines 與 Spectra verify/review，核對 DESIGN.md 單一權威、無 global theme/API/runtime 改動，舊 add-switchable-memory-views checks 原封保留，新未完成驗收不得勾選或 archive。

任務切片以每項不超過一小時為規劃估算，不是工期承諾；若實作探索證明超過，先把同一行為拆成可獨立測試的垂直批次並保留 dependencies，不擴張 scope。browser 3.1–3.7 與 archive 必須等最終合併 Mac 驗證，不阻塞已授權的雲端非 browser 實作／QA。

## 4. 最終合併交接

- [x] 4.1 [after: 2.3] 在雲端撰寫 recording-memory-visuals.spec.ts 與必要 fixtures，保留 30 基礎畫面、幾何／對比／長值／邊界與互動斷言；只執行 discovery、lint、含 E2E 型別檢查，不啟動 browser，真 modal/history 交由 A/B routed harness 驗證。
- [ ] 4.2 [after: 3.7, 3.8, 4.1] 最終 Mac browser 與逐圖差異驗收完成後核對新證據，執行 Spectra verify/review，修復並重驗後才 archive；目前保持未完成。
