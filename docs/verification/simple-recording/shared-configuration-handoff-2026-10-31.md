# 新舊記錄呈現共用配置與草稿交接

歸屬 streamline-recording-task-editor 2.3；不含後續選表、schema／Start 整合或 lifecycle。沿用 2.2 同一導航 guard、原 session draft owner、role 建構、routine mapping、API 與 operation recovery，沒有新 storage／queue／ledger。

## 可恢復草稿與未保存輸入

原 Basic session draft 保留 scoped 整數週期及 completeness policy；Circuit 原 owner 保留原始週期文字、名稱、來源選取、roles 與配置選項。能由原 loader 實際讀回恢復的意圖，仍可正常 Next／Back／Advanced／reload，不因尚未 POST group 就強迫丟棄。

只有未被原 storage 保留的 live intent 才提示 Keep／Discard。Discard 還原最後可恢復的值，不刪已保存 session draft 或已提交 operation。session 恢復不代表後端已保存、設備已驗證或可開始；pending／unknown 仍明示，離開不等於取消或重送。

同 workspace／device／circuit 的表單 owner 保留 storage-failed 原始輸入，跨 cached-read failure、成功 reread、destination 權威撤回／新 revision 不覆寫。已有表單保 mounted、顯示錯誤並拒絕依賴新證據的操作；初始無資料仍不可用。保留輸入不保留舊準備證據，preparedGroups、consent、command scope 仍隨原 revision／destination 失效。

三個原 connector dispatch handlers 未改，不在每字元輸入彈離開確認；真正換 device／circuit owner 才走同 guard。四個 source slices memo 保留 React #185 修復。Basic／Circuit 面板分別 497／446 行，採 cohesive extraction，沒有超長檔增行。

## 測試與獨立審查

- Storage 寫入失敗、正常 offline raw 恢復、cached-read 保表單、revision 不覆寫、Basic seed／Share 切換／destination 權威 loss 均有精確 RED → GREEN。
- 13 個依存檔最初 197 案例中 196 PASS，原 pending navigation 案例改明確確認離開，保留 late A/B identity／原值／不重送斷言。後續受影響 9 檔 134 PASS；其餘未變結果沿用，不稱一次全套執行。
- 最後準備證據失效修正為 2 檔 17 PASS；7 歸屬檔 lint、typecheck、差異檢查通過。獨立審查依原 task baseline／增量 patch／hash，無剩餘 blocker。

## Browser 首次保存發現與修復

run23 的未存／離線／raw／同 session reload 與真 Gateway 暫停造成讀取失敗後保留同表單及 Retry 通過。但首次 Circuit Save 後輸入 raw `1e` 再開 Advanced，誤觸離開提示。原 FAIL 保留；reload 後成功不是修復證據。

原 loader 會補 canonical member_roles；UI 原 defaultMemberRole 已將同 candidate 的缺省角色呈現成相同有效角色。舊 hook 整份 JSON 比較錯把缺省／明示等價當 dirty。只修改原 hook 的有效角色比較，沿用同 candidate identity 的原角色規則；所有 explicit keys、其他 raw 欄位與 pending 均保留比較，沒有改 storage 或角色權威。

真 RED 重現此邊界；修後 15 案例、兩檔 lint／typecheck 通過，真正不同 canonical State、storage failure／invalid raw／unknown 操作保護仍有效。兩檔增量獨立 review 通過。

run24 使用 fresh profile：VA 明示 State、VB 缺省 Measurement，第一次成功 Save-only 後輸入 `1e`，不先 reload 即可開同 profile Advanced、沒有錯誤 Leave 提示。返回 groups／Step1／Step4 保留 raw、來源、兩種角色；SQL group ID／revision／members 無額外改寫。最後導航及 console 無 React #185／uncaught。

前置 connector 修改後第一次 ensure 遭舊 workspace revision 409，原表單及明確錯誤保留；取得 canonical 新 head 後由明確 Save 成功，未自動重送。全程沒有 Apply／Start，目標資料庫未建立。Step3 原導航會產生 points/poll 與 mapping preview，因此不宣稱零 PLC I/O；read telemetry 不等於開始來源或寫入記錄。

## 證據與界線

- v2 七檔 manifest：8b5881cf7d7a13ec91d12da93fc751e1b6ee84ae935cd63e1d6650d2a800f8d1。
- 證據組：editor23-shared-configuration-evidence、gateway-native-23／24。run24 354 檔來源 hash 符合，26 份舊 DB 未變，backup integrity 通過；自有服務 20:05:55 UTC 已停止。
- Legacy managed_wide_v1 的 edit-existing 需要真既有 legacy group；安全 QA 副本沒有，現 UI 建的是 managed_wide_profile_v1，不能冒充。該 legacy 分支有正式形狀的單元／整合覆蓋，但本輪 Browser 未驗，不用 SQL／API 塞造假 identity。
- sessionStorage failure、accepted-write unknown／late ACK 為單元／整合證據，沒有用 Browser script 注入假故障。真 Browser 的 Gateway 暫停是讀取失敗，不冒稱已接受寫入未知結局。
- 後續 2.4 選表與 3.2 明確來源 Start／成功窗口／SQL readback 仍各自未完成。全套測試及 main 同步留至七案結束。
