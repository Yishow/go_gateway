# 採集驗收發現的狀態摘要成本

歸屬 integrate-protocol-aware-acquisition-batching 4.2；本文件只記已完成的局部診斷／修復，不代表完整負載矩陣通過。

## 量測條件及原問題

Linux、正式 runtime wiring、SQLite 可丟棄設定庫、兩台設備各 128 點、真 1ms RTT、1 秒採集週期、3 個 warm scan 加首讀；GOMAXPROCS=1、GOGC=50。使用正式 adapter、原 scheduler、ingestor 及 pipeline。

原驗收 SampleSink decorator 未轉發 SampleInterest，會讓未配置記錄的來源也進 pipeline。舊 capture-all 結果、FAIL／loss 保留為該條件歷史，不冒充未加觀測器的正式效能。修正僅限測試 wrapper：performance 模式轉發原 pipeline.WantsSample，原刻意 capture-all wiring fixtures 保留原意。

忠實 interest 版本仍觀察到 native 1024 點、collected 38、queue 986、output loss 0。CPU profile 的巢狀累計 70.74% 在逐點 publishDerivedStatus → RuntimeStatusSnapshot → RuntimeProjection 整 workspace 查詢；45.55%／37.40% 的 mapping 路徑包含於其中，不得相加。Profile 含 setup，不是隔離查詢 benchmark。

## 最窄產品修復

僅 internal/datalink/runtime/ingestor.go 移除每點同步重建狀態摘要的呼叫及 dead helper，沿用原有一秒 statusLoop。不新增 queue、cache、worker 或查詢權威。

fresh GET、單設備訂閱的初始 fresh read、逐點 PointState、ValueEvent 品質、diagnostics 與 SampleSink 路徑不變。摘要使用同一 derived model；同 refresh 間隔短暫 error→good 可以合併，因此 status stream 不是無損狀態轉換稽核。Ticker 加計算時間不構成硬性一秒內送達保證。

## 已驗證項目

- 真 RED：有／無 status subscriber 的四個輸入點各造成四次 workspace projection 查詢。修後為零，仍逐點保留 ValueEvent 與 quality。
- 精確新測試及相關 status／projection／sink／lifecycle 選測通過；精確 race、runtime vet、官方 golangci-lint 2.6.2 三檔增量檢查通過，lint 零問題。
- fresh API／初始訂閱、周期 stale／持續 fault／open breaker 與 snapshot 一致、原 Stop 等待 worker 結束均有回歸。獨立唯讀審查無確定 blocker。
- 三檔凍結：ingestor.go SHA256 1d4e4c7ae3119ee52518f106cfc1f930818452ad8f5aec767a42ddd0b3e46618；兩個新測試及完整命令／hash 位於 status-cadence-evidence。

## 修後相同短樣本

同 faithful pilot 修前／修後 process CPU 為 3.290／1.052 秒，collected 38／809，cutoff queue 986／215；兩次 native 1024、output loss 0。這是短樣本原始量測，不推估全年／所有設備效能，也不將尚在 queue 的資料稱已交付。

修後實收並核對 788 份 ValueEvent；該事件通道有界且可丟棄、沒有 AcquisitionID，不是每個 scheduled 點的完整值證據。native observer 證品質及 identity；paired adapter 的全點值等價證據另列。沒有配置記錄的 pipeline 正常拒收 interest，不因此宣稱 SQL committed。

另有真 status subscriber 短案例：初始 seed／bootstrap／fresh API 0.671 秒單列；其後 3.092 秒、CPU 0.890 秒，各設備收到兩次狀態摘要。native 768 = collected 529 + queue 239 + output loss 0。這不是 cold 分布或正式 100-scan 矩陣案例。

## 未完成與證據

完整 faithful scheduled 矩陣仍待執行；舊 capture-all 1/18 僅歷史條件，不混入新版完成數。Strategy adapter 矩陣原 12/18 與 scheduled 計時分開。不得由本次短測宣稱高負載無損、全部採集已入庫或整個 change 完成。

原始紀錄組：batching42-evidence 的 interest-pilot、interest-after-status-fix、status-subscriber-pilot、scheduled-evidence-boundaries；status-cadence-evidence 保留 RED、正常／race／vet／lint、精確 patch 及審查 hash。修後量測前後三個 runtime 檔案 hash 一致。
