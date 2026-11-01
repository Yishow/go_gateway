# Editor 4.1：批次驗證交接

## 本次驗證方式

本項承接已完成 1.4、3.3，不另開產品功能。按使用者要求與 docs/testing.md，開發期間只跑直接受影響案例；來源、契約與證據仍有效者沿用，不把全套、race 或整個 frontend 測試當每次修改的固定流程。七個 change 全部完成後才執行最終完整檢查。

目前這是 editor 批次交接，不代表 lifecycle、adoption 或採集性能矩陣完成。

## 有效證據索引

- 來源／範圍與 canonical 保存：`editor-range-handoff-2026-10-31.md`、`editor-protocol-parity-handoff-2026-10-31.md`；保留唯一設備、rule、mapping、Point/tag 保存 owner、明確 ACK／讀回與 matching modal 的邊界。
- Planner 預覽／每點字序：`acquisition-preview-handoff-2026-10-31.md` 及既有逐點字序交接。正式 actual 與 unavailable 分開，不能以估算補 0；字序不由 UI 解碼。
- 同頁草稿、失敗恢復、原 guided 路徑：`draft-navigation-handoff-2026-10-31.md`、`shared-configuration-handoff-2026-10-31.md`，包含原 mapping/autosave/source 相關直接回歸。
- 目的地選取：`configured-destination-handoff-2026-10-31.md`；同名不同 scope 不混、已有寬表歸屬保護、選擇本身不寫資料。
- 每成員計算：`member-calculation-handoff-2026-10-31.md`；每個 member 的 UTC 對齊窗口、not_due、first baseline 與真失敗原因不同，不帶值或補 0。
- 準備與明確開始：`preparation-flow-handoff-2026-10-31.md`、`scoped-start-backend-handoff-2026-10-31.md`、`scoped-start-ui-handoff-2026-10-31.md`。
- QA31 真兩個 configured logical devices／同 simulator 的 15秒平均12.5、30秒平均25、60秒用量19完整窗口與 SQL／effect／digest／record ID 對齊；不是兩台實體 PLC。
- QA32 僅重驗後續 UI caller 修復：唯一 Start、同 session Close 列表與 fresh reopen，未重跑已通過完整 60秒窗口。
- 外殼、表單、參數說明、語系與 keyboard/RWD：`theme-parameter-help-handoff-2026-11-01.md`。QA34–38 逐項修復、直接案例及獨立 review；未以重跑整張矩陣替代影響判斷。

上列手冊內較早的 pending／FAIL 為歷史，需閱讀後續 closure；不刪除原失敗，也不把歷史 pending 再當成新的未驗需求。

## 兩個唯讀重讀接線修復

### 採集預覽

原 QA29 最後 Step1 GET78 使用舊 DB head 返回409，當時未點 Refresh，並非刻意負向。`acquisition-preview-refresh-2026-10-31.md` 記錄後續窄修及 QA30：明確 Refresh 先讀 canonical writeGroups head，再用該合法 head 讀預覽。回覆 installation_unavailable 是真實結果，沒有偽造 available，也不代表 0 request。保留 legacy 合法空 head、晚回覆／foreign guards；不寫 draft/CAS 或自動啟動。

### Step4 runtime evidence

QA34 GET140 發生在 Step4 mount（Advanced click 在後），cached head 安全拒絕409後有新 head 200；當時後者 body 未保存，不能由 status 推定 UI 已接納。當時0 saved groups，沒有驗到 saved summary 的 explicit Retry。

真正接線問題是 scope 由 writeGroups 提供，原重讀卻只刷新 runtimeContext。精確三個 source＋一個 test 修復如下：

- shared hook 接受穩定的 exact authority key，預設 runtimeContext 行為保留。
- Step4 explicit 重讀先取 canonical groups；只有 same head 才 refetch 原 exact evidence，新 head 由既有 observer 派生新 key，不手送舊 scope。
- workspace/device owner、current revision、generation 與 mount fences 保留，失敗／foreign fail-closed；不 dispatch draft／connector CAS，不 Start、不 schema mutation。
- 真 RED：6個預期行為失敗／1 default pass。GREEN：18/18（10直接＋8 shared）、四檔 lint PASS；獨立 review 四檔 SHA 與權威邊界無 blocker。
- QA35 由未啟動、已存一筆記錄的隔離 fixture 真點原「重讀證據」：request42 canonical groups→request43 matching head evidence200，兩者 raw body 保存；只有 saved intent，effective=null，零 mutation，最後 console/navigation 通過。
- 該 Browser 案例是 **same-head**，不冒充 R→R2；delayed/newer/failure/foreign 由直接測試證明。wrapper-specific device/ABA 與已有 fresh evidence 後的 pending/failure hiding 為 source-review 證據，沒有另外命名的直接案例。
- 原 mount cached head 仍可安全409；本修復保證正確權威重讀，不宣稱取消所有防禦性409。

## React #185 與最後導航

原 run10 FAIL 保留；run11–12 單純未重現不能當成已定位。run13 於 DB 路徑逐字輸入取得 stack，定位 CircuitProfileBasicPanel 將整 state 當 candidates memo 依賴，造成每字元重建 hydration；僅改四個真 source slices。原真 RED→GREEN、run14 原87字路徑及247字名稱／274字路徑的 DOM／PUT／ACK／SQL／reload 對照已閉合因果回歸。

後續 source form／Start／Close／返回及 QA34–38 最後 console 持續驗證；QA35–38 最後 console 無 error/warn。沒有把 SQL 成功代替頁面導航驗收，亦沒有在本項重做已有效的長字串測試。

## 最終檢查與限制

- QA38 同一合併來源的 strict TypeScript→Vite→static sync→Go embed PASS，494 凍結 source hashes 一致；最後實際頁面導航及 console 通過。
- Browser 後 TopBar 只刪字串外一個行尾空格，保 v8原始 bytes 與精確1-byte等價證明；不冒稱 v9 與 Browser 原 bytes 相同，也不為此重跑型別或 Browser。
- root `git diff --check`：PASS。初次單一行尾空格 FAIL 與修正證據保留。
- root `make check-lines`：PASS；該次檢查526個變更文字檔，既有超長檔未增行。這是文字行數守衛，不是全套產品測試。
- 每批直接 lint、型別、測試與 Browser 證據按真實執行／沿用區分；沒有全套 Go／frontend 測試或本項 race。
- 真 tap/pinch 與真正 browser zoom 未驗成功，工具限制及 emulation／單元證據的界線保留於 3.3 交接。
- 沒有宣稱所有次頁、現場 serial、實體 PLC、Windows、PostgreSQL live 已驗。
- 下一步為 4.2 獨立檢視權威、草稿、confirmation 失效與進階交接；通過後才接 lifecycle，不重做已完成的採集／計算引擎。
