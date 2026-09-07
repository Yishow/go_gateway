## Context

main 的 writer 有 grouped map 與 flush loop；這不能單獨證明斷電可恢復。
新路徑對 enabled recording plan 使用耐久入口；既有未遷移 output 保持原契約且不能冒充可靠記錄模式。

## Goals / Non-Goals

目標：有界資源下可恢復、可追查；重送不重複資料或用量；任一目的地失敗不阻擋其他目的地。
不做：多資料庫分散式交易、無限磁碟、硬體損壞零遺失承諾、把所有其他網路輸出納入本次實作。

## Decisions

1. 先將帶穩定身份的樣本及 source sequence 落本機 durable journal，再對 recording intake 回成功。
   計算器按該 journal 推進，checkpoint/result/outbox 在單一 transaction 提交。
   磁碟與交易設定須符合 Windows/Linux/ARM 支援範圍，測試 fsync／正常關機／異常終止。
2. outbox 以 destination_id、destination_revision、plan_revision、record_id、calculation_revision 識別。
   record_id 重試不變；目的地路徑、表、型別與定義快照隨紀錄凍結。
   destination 更新不改舊佇列的意義；憑證輪替可對相同 identity 授權更新，不將資料改送別台。
3. managed 目的地使用同一 DB transaction 寫 receipt 與資料。
   斷線發生在 commit 後、回應前：恢復先查 receipt/unique key，不盲目再插入。
   原子性依 adapter 實際能力；trigger、非交易表、不可控制外部副作用不承諾單次效果。
   custom 表缺 durable unique key/ledger 能力則標 unsupported-reliable，不能默默退回無保護重送。
4. 同一 series 的輸出保持順序；不同目的地與設備使用有界併發，不因慢目標阻塞全部。
   暫時性失敗用有界指數退避與 jitter；密碼／權限／schema 錯誤進 blocked，修復後以相同 identity 重試。
   poison record 隔離、保留原值與原因；不跳過後假稱整個時段完整。
5. 分開 pending、sending、delivered、retrying、blocked、quarantined；
   last_observed_at、last_local_commit_at、last_delivered_at、pending_count/bytes/oldest_age 全部是真實數值。
   delivered 不代表 readback_verified；錯誤 code/action 安全顯示，原始帳密／DSN 不進日誌與前端。
6. 明確 disk quota、high-water、critical-water、排程公平性與 per-destination backlog 限制。
   建議預警80%、危急95%，皆可調整；不等到磁碟滿才提示。
   預設滿額時停止接受受影響 plan 的新可靠紀錄並標 data-loss-risk，不偷偷刪未送資料。
   現有設備採集／Modbus 分享是否仍運作分開顯示，不能以分享正常掩飾記錄暫停。
   管理者明確選擇丟棄才可釋放 backlog，留下時間範圍、筆數與計量不完整事件。
7. journal 在所有必要 consumer checkpoint 已前進且 evidence retention 滿足後才回收。
   慢目的地可占用自己的有界 backlog；receipt 保留期限至少涵蓋最大重播期限。
   超過 receipt 去重保存範圍的人工重播需重新驗證，不能宣稱一定不重複。
8. 啟動沿用既有 workspace revision/readiness/autosave barrier。
   plan 與 stream 只能使用已套用版次；草稿保存不立即改動運行中的記錄。
   套用時檢查 epoch、schema、queue 與 calculator boundary；舊排程停止接新樣本後切換，
   已接受舊樣本照舊版處理，來源 span/probe 與 Share token 保持原 gate。
   某設備不能啟用時，其他被使用者選定可獨立啟用的設備保留結果，不顯示全成功。
9. 正常關機有有界 drain；逾時留 durable backlog 不強制清空。
   重啟先 recovery/lease/checkpoint，再開始新 worker，防止兩個實例同時消費相同 outbox。
   若部署允許兩個 Gateway 共享寫出身份，須有租約/fencing；本機 journal 不宣稱跨主機 HA。

## Risks / Trade-offs

耐久寫入增加磁碟 IO，須用 batch transaction 與有界 queue，而不是降低可靠性卻不告知。
長時間外部失聯可能超過容量；容量不足是可見失敗，不是無限暫存承諾。
既有資料庫觸發器可能產生外部副作用，端到端唯一效果需限制能力與說明。

## Migration Plan

新增 durable tables、receipt ledger 與 migration，默認只给新 plan。
舊 grouped 記憶體資料切換前 flush 或明確標示不可恢復範圍，不能捏造已遷移成功。
回退先停止新 intake，保留所有 backlog 與 checkpoints；不得刪除待送資料。

## Open Questions

現場磁碟配額與可承受暫停時間由部署設定確認；預設保資料、不靜默丟棄。
