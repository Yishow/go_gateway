# 記錄草稿導航驗收交接

範圍：streamline-recording-task-editor 2.2。既有 Basic／Circuit 本地草稿接同一 guard 屬後續 2.3；來源開始及成功窗口寫入仍由 3.2 驗收。本次不新增持久化草稿引擎或來源啟動入口。

## 實作及安全邊界

- 共用暫存導航 guard 處理 Close、Escape、首頁／步驟切換、SPA 歷史及原生 reload 提醒；保留維持當前欄位，明確捨棄只清未送出意圖。
- 已保存設備、範圍、逐點格式及已提交操作仍由原 owner 管理，不因捨棄記錄而回滾。pending／unknown 不冒稱已取消，不盲目重送。
- 原 Save 自 preflight 至 settlement 保持 busy／scope 保護；舊回覆或 finally 不解鎖新 scope，捨棄後尚未送出的 POST 不得發出。
- 快取重讀失敗保留已掛載表單，顯示錯誤並拒絕需要新證據的操作；重試不覆寫未保存欄位。schema consent 隨輸入 generation 失效。
- 路由採受支援的 blocker，同一 catchall data router 保留原路由／providers；沒有 history monkeypatch、第二份 queue 或 ledger。

## 實際測試

- 原 24 個相關檔案、220 案例以初次執行及精確失敗案例重跑閉合，不宣稱一次完整全綠。未變 21 檔結果沿用；定位修正保留原 request、DOM 及 unknown receipt 斷言。
- 27 檔 ESLint、前端 typecheck、差異檢查通過；v1 獨立唯讀審查通過。沒有跑全專案測試。
- Linux 真實 Browser run20／20b：dirty Close／Escape／Step1／SPA Back 保留及捨棄、reload 原生確認、巢狀來源返回、已存 range 不回滾、離線 preflight 後捨棄不送 group POST、Gateway 暫停造成真快取讀取失敗後保草稿及明確重試恢復。短確認檢查包含桌面及手機尺寸。
- IME、dirty SPA Forward、foreign workspace、已接受寫入後 late ACK／unknown 為相關單元／整合測試證據，不宣稱已做真 Browser 故障注入。

## 新增來源首次點擊缺陷

run21 在正式頁面重現：三列來源／已選一列時按鈕部分裁切，第一次點擊只取得焦點，第二次才進入設備選擇；Tab 後仍被裁切，Enter 可進入。保留原 FAIL，不歸咎於等待不足。

精確回歸測試證明原 focus-only rerender 將原生 scrollTop 29 回寫為舊值 0。只分離既有 scroll restoration 與初始 focus restoration effects；保留 scope、page、filter 及返回焦點契約，未改業務 owner。

- 修前新增案例失敗；修後 3 個相關檔案 32 案例、兩檔 lint、typecheck、差異檢查通過，兩檔增量獨立審查通過。
- run22 同三列／一列已選及原裁切位置，首次點擊即進入設備選擇；Tab 自動滾動至 29、按鈕完整可見，Enter 一次成功。返回、Cancel／重新開啟保留 scroll 29 與來源焦點。
- run22 再驗 Escape／Keep／Discard、最後回 Step1，console 無 React #185。347 檔來源雜湊一致，23 份舊 DB 不變；本次沒有儲存來源／記錄或啟動採集，測試目標 DB 未建立。服務已停止。

## 證據識別與後續

- v1 manifest：5fcef3e5a627edb9b410c52122e7d2c130c0710f2ea60e681b70a3f92a2ccf5b。
- v2 manifest：a0a9d07e71a2a005975808fcc85b170a8128395f92939c540127ffed3d0f9ebc，28 個歸屬檔案。
- 證據組：editor22-draft-navigation-evidence；gateway-native-20、20b、21、22。原 FAIL、增量 patch、命令及原始 Browser 紀錄均保留。
- React #185 的原 run13 因果修復及 run14 輸入一致性驗證沿用；本次 Circuit 四個 source slices memo 未改動。
- 下一項 2.3 只將原 Basic／Circuit 草稿 owner 接同一 guard，保留 role、offline draft、routine mapping 及原配置。完整七案結束後才跑全套及同步 main。
