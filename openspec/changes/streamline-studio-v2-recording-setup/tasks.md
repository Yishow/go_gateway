## 1. 基本資料確認與預設群組

- [ ] 1.1 先補前端 RED（Four-step intent-led setup）：既有點位可產生真 Point-to-Tag、只確認必要資料；新基本設定每設備一群組，持久 workspace/device/role key 在 reload/重送不重複，同名設備與目的地變更不誤用，既有 advanced 群組不被轉換。
- [ ] 1.2 沿用既有 source/tag/group service 完成基本保存與 create-once/reuse，讓 1.1 GREEN；保留保存失敗及離線草稿，更新步驟文案。

## 2. 開始記錄的後端協調

- [ ] 2.1 補 RED（Recoverable scoped recording start）：保存→readiness→Apply→activation 各階段失敗、stale revisions、雙擊、失去回覆與重啟恢復；同 scope/request 必須同 operation，不影響鄰近群組；A ready/B incomplete 不阻擋，A 啟動/B healthy 不重投影或停止，恢復仍用原 scope。
- [ ] 2.2 用既有 services 與 operation ledger 實作狹窄 start 協調，讓 2.1 GREEN；不同 action 不混用 token，未確認 schema 絕不 DDL；同步安全 DTO、SDK 與 generated Swagger。
- [ ] 2.3 接 UI 的單一開始動作與可查進度，補 late response／切換群組／reload 回歸；部分成功不冒充 rollback 或整體成功。

## 3. Step 4 與交付事實

- [ ] 3.1 驗證 Basic recording controls hide incidental complexity：先補 RED 再移除基本路徑的 legacy targets/row_groups 必填依賴，接 D 的預覽/確認；自訂欄位、revision 與試寫移至可達的進階區。
- [ ] 3.2 驗證 First recording result is evidence-backed：主畫面呈現 interval、snapshot/缺值語意、等待首桶及採集/本地接受/SQL committed 三階段；測試 queued/running/test-write 不被當成已落庫。
- [ ] 3.3 驗證 Share-only 不需 DB、不啟動 writer，保留原 readiness gate；同步 en/zh-TW、鍵盤及 390/768/1440 互動回歸與必要 runtime handoff。

## 4. 整合驗證

- [ ] 4.1 執行受影響前後端套件及 repository 最低檢查，文件改到 surface inventory 時同步既有 changelog 工具，不自行新增管理制度。
- [ ] 4.2 review 確認無新協議、模式、框架、dashboard 重寫或第二套 ledger；修復本案發現後重跑，交給 F 實際驗收。
