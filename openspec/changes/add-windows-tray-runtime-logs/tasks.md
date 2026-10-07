## 1. 先建立可測的模式與 lifecycle 邊界

- [ ] 1.1 以 Red-Green-Refactor 建立 desktop／console／explicit headless 模式及data-source resolver：涵蓋 `internal/config/config.go`、`cmd/test_ui/harness_config.go` 的既有優先序、absolute aliases、exe/cwd歧義、新庫確認、取消與無權限；驗證不默默新建／搬移資料，Linux/macOS CLI語意不變
- [ ] 1.2 [after: 1.1] 以 Red-Green-Refactor 把 `cmd/test_ui/main.go`、`server_runtime.go`、`service_wiring.go` 的啟動錯誤及狀態接到可注入lifecycle介面：先bind再啟動採集，服務及assets確認後才ready；port衝突、missing assets、runtime degraded、early fatal皆有獨立結果且不kill他程序

## 2. 先安全投影，再加入有限保留

- [ ] 2.1 建立safe event schema與bounded broker，先用secret fixture、未知raw error、8KiB截斷、count/byte双邊界、burst/admission drop、sequence identity並行測試RED，再完成模板／欄位白名單、ring與drop counters的GREEN及race檢查；所有新sink只收safe events
- [ ] 2.2 [after: 2.1, 1.1] 實作非阻塞bounded diagnostic file sink與3檔15MiB輪替：測同目錄多DB各自namespace、LOG_FILE/backup碰撞與DB/WAL/SHM alias拒絕、symlink/reparse競態、初始無權限、disk full、writer stall、恢復gap及shutdown flush；只輪替owned files，native early error只顯示確實保存的檔案，不影響DB/outbox/receipt
- [ ] 2.3 [after: 2.2, 1.2] 盤點並接入production standard log／slog／Gin safe access與recovery／HTTP ErrorLog／startup/shutdown必要producer；替必要runtime錯誤提供safe code，禁止raw stdout/stderr tee；測秘密不進任何新sink、log endpoint不自我放大、console formatter及 `/debug/logs` 原契約另行維持，記錄未捕捉範圍

## 3. 有本機門禁的 API 與網頁

- [ ] 3.1 [after: 2.3] 在 `internal/api/router.go` 與新handler加入log-only guard及snapshot；以IPv4/IPv6、無Origin正常GET、同源Origin、Host/port/DNS rebinding、LAN、proxy headers、Sec-Fetch-Site same-site無Origin／重複／malformed與wildcard CORS繞過負測試驗證，再加入bounded filter/search/page/cursor及safe錯誤；其他route政策不變
- [ ] 3.2 [after: 3.1] 以 Red-Green-Refactor 完成snapshot→bounded replay→live原子邊界與SSE：測同時append、鎖不跨IO、replay大於live queue／replay-limit gap、相同filter語意與無matching時progress、Last-Event-ID衝突、future cursor、retention/restart gap、8-client上限、slow consumer隔離、5秒write deadline及並行cancel/shutdown無洩漏；重連ID採字串不丟JS精度
- [ ] 3.3 [after: 3.2] 依UI測試先行新增 `/studio/logs` 的lazy route與Studio/runtime入口，涵蓋level/source/search、pause/resume/tail/clear、browser雙容量界線、A→B遲到回覆、單一reconnect owner、403停止重試、plain-text XSS及en／zh-TW鍵盤可達狀態；執行focused Vitest及embedded Playwright

## 4. Windows 本機操作與安全退出

- [ ] 4.1 [after: 1.2] 先以並行launch／等價path／stale PID／owner死亡／headless owner測試RED，再加入Windows OS-backed database owner guard及same-user/session ACL IPC；二次開啟只能向已驗owner要求固定setup action，無任意URL/command、無port推定owner、無他程序終止
- [ ] 4.2 [after: 4.1, 2.3, 3.3] 加入Windows-native tray、shell open及error/version dialogs；測初始tray失敗退出、固定右鍵選單、狀態分離、AUTO_OPEN_BROWSER、IPv6 URL、LAN-only logs不可用、browser close不退出、Explorer重建不重啟runtime；同步native實機驗證無console閃窗
- [ ] 4.3 [after: 4.2] 以 Red-Green-Refactor 完成idempotent shutdown：先測tray取消／確認／signal race、HTTP mutation gate與SSE取消、runtime intake→group worker→Share/connections/DB順序、deadline pending phase、等待及二次確認force exit；以可丟棄DB核對accepted backlog/unknown在重啟後保留，不close活worker使用中的DB或宣稱全部送完

## 5. 打包、回歸與實機交接

- [ ] 5.1 [after: 4.3] 更新Windows desktop／console建置選項並內嵌icon與完整SPA：`scripts/build.ps1`、`Makefile`及必要 `start.ps1` 模式契約保持一致；以fresh clone/build-script tests驗證frontend失敗closed、WindowsGUI PE與console artifact分流，Linux/macOS不吃Windows flag
- [ ] 5.2 [after: 5.1] 對最終source執行 `go test ./...`、`go vet ./...`、`golangci-lint run ./...`、受影響並行package的race tests、frontend lint／完整Vitest／build／Playwright、`make check-lines`、`git diff --check` 與官方OpenSpec strict validate；保留命令、SHA、結果及未執行項目，focused不得冒充全套
- [ ] 5.3 [after: 5.2] 在Windows互動桌面驗收portable單exe、無Node/Go環境、無console/閃窗、tray keyboard/overflow、Explorer restart、duplicate launch、DB歧義／只讀路徑、port衝突、startup fatal、離線log與disk fault、正常/逾時退出；另測GUI exe explicit headless在cmd/PowerShell/supervisor的redirect、wait、exit status與AUTO_OPEN_BROWSER override，核對legacy console及Linux/macOS CLI；相同artifact的embedded logs＋durable recovery證據要可追溯
- [ ] 5.4 [after: 5.3] 更新使用與交接文件：build命令／選單／模式／資料路徑選擇／local-only限制／logging容量與不可捕捉範圍／timeout與rollback；若修改 `docs/technical/studio-surface-inventory/` 則依repo規則同步changelog DB。逐scenario對照證據後才討論archive；缺實機平台、field或正式DB驗證明列NOT RUN，不在本案自行部署

本檔全部為後續實作任務。規格文件完成、CLI artifacts顯示done或通過strict validate，皆不使上述checkbox完成。數值為design的契約上限，非此輪效能測量。
