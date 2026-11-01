# Scoped Start 與原請求查詢：後端局部交接

歸屬 editor 3.2。七檔後端接縫及相關測試已驗；前端 saved-record owner、真 Browser、完整 mean／delta 窗口及 SQL/readback 仍待整合，不因此勾完 3.2。

## 已定位的 mixed delta 啟動問題

真 production wiring fixture 的兩設備 mixed 記錄可通過 readiness，但 scoped Start 在 Apply checkpoint 回滾。相同 revision 的原 WriteGroups.Apply 診斷可成功；根因為 recordingStartCheckpoint 只把 tag types 傳給 managedSourceDigest，漏掉 mixed delta 必要的 counter proof。

只在原 checkpoint 的同一 SQL transaction 透過既有 counterBindingEvidenceInTx 取得 proof，遇 issue fail closed，再計相同 managed source digest。沒有在 transaction 外讀另個快照、巢狀呼叫 ManagedSchemaScope 或改 counter 算法；legacy digest 不變。暫時 direct-Apply 診斷分支已移除，最後測試仍走真 scoped Start。

## 原請求唯讀查詢

新增 POST /studio-v2/workspace/recording-start/lookup，body 為完整原 RecordingStartRequest，沿原 normalize／digest、既存 workspace repo.Get、ledger.FindRecordingStart、decode／原 Operation DTO。POST 用來避免 readiness 資料與完整請求出現在 URL，不是再送一次 Start。

不 GetOrCreate、不 Claim／Resume／Apply／Activate／probe，不以目前 source revision 重新驗歷史請求。原 workspace／request-ID／payload digest／action／detail 校驗保留；not found 或查詢失敗仍 unknown，不能當安全重送許可。沒有新 store、ledger、migration 或前端私有 operation-ID 算法。

## 實際證據

- 真 HTTP、production services、loopback Modbus peer、SQLite 設定及目標庫；兩設備一 mixed（兩個 mean、一個 delta），第三設備未選。
- 原 scoped Start 成功，精確三個 member PointIDs installed／native Good；同設備額外 source 未裝，第三設備仍 draft。只有原一筆 group／record，沒有隱藏 Basic record；replay／GET 同 operation，changed intent 409。
- 五個真負向：wrong workspace、stale group、stale connector、foreign member、沒有 schema proof。沒有來源開始、沒有錯誤 Apply 效果。
- Lookup 使用 SQLite BEFORE INSERT／UPDATE／DELETE abort triggers 覆蓋 workspace、ledger、groups、devices、points、tags、mappings、measurements，證明成功查詢零寫入；歷史 head drift 可查、缺 workspace 不初始化、foreign／digest／action／corrupt 拒絕、running／unknown／partial 不 resume。
- 精確新增測試及七個相關 legacy regressions 通過；race：cmd/test_ui 32.359 秒、internal/api 47.098 秒；四個受影響 package vet 通過。官方 golangci-lint 2.6.2 最初只有四個 import 格式及一個 fixture if/switch 問題，修後 exact seven-path patch 零問題；最後五個新增測試再通過（2.434／1.695 秒）。
- 獨立七檔唯讀審查零 findings、最終 hashes 相符；format-only 修正未改來源啟動或查詢語意。沒有宣稱重跑全套。

## 界線與識別

Native Good／operation succeeded 不是完整 15／60 秒結果窗口、SQL row／receipt 或 Browser 證據。這些由原 3.2 前端及 QA 後續閉合；Start ACK 只能進 waiting-first-write。API Swagger annotation 已寫，generated Swagger 由 root 最終統一重生，不能宣稱目前 generated docs 已更新此路徑。

七檔清單、原始 baseline、所有 logs 與完整交接在 mixed-start-evidence；manifest.json SHA256 afb65677669486a6c0e1535fdb05eaacec10ab8a0ecd5fdb427aaea81fbffbe5，exact.patch SHA256 c8b29f0f657e38b79ac8ce5a228c68c8d01a8ae81bb0ca63a6299f5477a4c2e7。沒有 commit／push／部署由 worker 執行。
