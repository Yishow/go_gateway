## Context

以總覽固定的 main 為基準。Update 保存 draft 並保留 AppliedRevision，但 Reconcile 以 status=ready 篩選；boundary 的恢復目前先查遠端表，ledger Restore 又隸屬特定 revision。until 的讀取早於同一 mutex 的取得。動機見 proposal.md。

## Goals / Non-Goals

修復現有採集與排空責任，不新增通用工作流、另一個 journal、跨程序採集協調或新交付模式。不得把「未 ACK」說成已持久保存；也不得因退役 UI 而抹掉資料。

## Decisions

1. 以有效 immutable applied snapshot 決定 intake，draft 只是待套用設定。保存名稱或 member 草稿不停止舊版本；再 Apply 的切換界線按既有 UTC 與共同區間規則決定。拒絕以 draft status 直接控制 runtime。
2. 在原 applied-version 資料中保存可重建 runtime 的已驗證描述：凍結的 source/tag 型別、欄位與編碼布局、目的身分、capability 與 digest。Apply 先 read-only 檢查，交易內再檢查原有 revisions；不保存憑證。新版本取得可靠描述之後，冷啟動不依賴遠端 InspectTable。
3. 啟動恢復分兩類：現在有效的版本可收新樣本；仍有未封桶已接受 journal 的舊版本只恢復 closure。掃描依既有 group/revision/checkpoint，不以目前群組清單的 status 排除歷史資料。已成 outbox 的資料仍交由既有 sender。
4. 切換、Disable、Delete 的責任截止必須從持久狀態重建；不可重啟後只看到記憶體 until 消失。同一版本再啟用不可沿用已退休 boundary；保留既有樣本／effect identity，不回灌已封桶資料。
5. until 的檢查與接收樣本在同一鎖域；切換也使用相同同步機制。鎖內只做必要本地一致性工作，不加入遠端 I/O。blocked 群組不授權 legacy writer 自動接手。

## Risks / Trade-offs

- 舊版本缺少完整恢復描述：只有可連線且能證明原 source/mapping/connector revisions 與歷史布局一致時，才能一次性補上可驗證描述；不能證明或當下離線時保留 journal 並明示 blocked/loss-risk，不從現在的 Tag、草稿或別的目的地猜測，也不聲稱舊版已具離線 intake 能力。
- 同時存在多個歷史版本：沿用有界分批恢復與現有 worker，不一次載入全部 accepted data。
- 停用與封桶交錯：用固定時鐘及 ACK/closure/CAS 故障點驗證，不靠 sleep 猜測。

## Migration Plan

只允許既有版本／生命週期儲存的 additive 延伸。保存原始版本、checkpoint、journal、outbox、receipt；重跑 migration 不重置資料。回退先停止新 intake，不用舊備份覆蓋 accepted data。更動 connector 後仍不得把舊 backlog 送到新 endpoint。

## Verification

先在 production service wiring 補 RED：草稿保存後跨兩桶、同語意再啟用、DB 離線冷啟動、ACK 後停用／刪除／切版崩潰與 until race。GREEN 後檢查 SQL、journal、effect key 和可查狀態；跨案真 UI 矩陣由最後驗收案整合，不能代替本案回歸。
