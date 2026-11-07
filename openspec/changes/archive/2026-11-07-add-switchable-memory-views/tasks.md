## 1. 共用控制與型別呈現

- [x] 1.1 以測試先行實作共用狀態與安全權威：Layout changes preserve shared memory context、Every layout preserves safe cell actions；五種切換保持位址與原 click 契約，執行參數化 layout 單元測試與現有 scope/values 測試。
- [x] 1.2 [after: 1.1] 完成共用格子與型別連接：Typed values retain physical address geometry；先加入跨列／頁與 string-width 失敗案例，再實作連接與型別標記，執行 memory map/view/keyboard 測試。

## 2. 五種版面與驗收

- [x] 2.1 [after: 1.2] 完成五種版面與共同控制：Layout-specific navigation shares the current data；實作導航、地圖、章節、詳情與 locale，參數化測試每個真實動作與相同資料投影。
- [x] 2.2 [after: 2.1] 驗證 Themes and responsive access preserve functionality；明暗與 1440/768/390 逐功能瀏覽器檢查，保存五圖對照截圖，修正實際差異並重測焦點、wrap、observer；取得本輪執行確認且環境允許後，以保留 HSL 的 MC／Modbus typed seed/readback 驗證即時值更新與失聯，未執行明列 pending。
- [x] 2.3 [after: 2.1] 執行 spectra-review 及修復，完整 frontend lint/test/build、diff 與 line gate，保留明確通過／未執行／blocked 證據及最終 commit 交接。

驗證紀錄（2026-11-05）：Mac 隔離環境完成五圖對照、明暗 1440/768/390 的 118 項瀏覽器驗證及焦點／wrap 重測；observer 重綁與四格跨列／頁由完整單元測試驗證。保留 HSL 的 MC／Modbus typed seed/readback、即時值變動與失聯恢復均已實測；私人截圖與協議紀錄保存在驗收目錄。
