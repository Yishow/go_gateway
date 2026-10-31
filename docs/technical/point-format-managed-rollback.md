# 點位字序：受管理回退與設定備份演練

此程序延伸既有 [設定備份／回復邊界](studio-v2-write-groups.md#備份與回復邊界)，由部署負責人執行；不是部署 engine，也不能阻止任意人直接啟動舊 binary。本案自動驗證只用可丟棄資料庫；正式部署、真 PLC／serial 硬體及正式歷史 DB 沒有因此取得驗收。

## Expand → migrate → contract

- **Expand**：048 只新增 nullable mode/override、version/revision。舊資料不從 `data_format` 猜成 override；缺 mode、version 0 仍是 legacy。migration 重跑由既有 Migrator ledger 管理，不能直接重跑 ALTER SQL。
- **Migrate**：新點是 inherit，明確設定才 override。各 protocol adapter 消費同一 immutable ReadProfile；來源 rule identity、版本、lease 與 fingerprint 綁定當次安裝。
- **Contract**：正式來源／manual／scheduled 解析走共用 resolver；保留 legacy reader，不刪除舊欄或強迫全量轉換。這不是可將新版意圖靜默投影回舊欄的授權。

## 先檢查相容性，不能先換 binary

1. 確認實際設定 DB、目前與目標 binary、所有 migration/feature 能力。先正常停止 intake／writer 並確認 ownership 已結束；timeout／unknown 不算已停止。
2. 保存目前設定、journal／outbox／receipt 與原產物，備份須完整且沿原 secrets 權限管理。不可將現場備份提交 Git。
3. 在**一致備份副本**執行點位格式檢查：

   ```sh
   python3 scripts/check_point_format_downgrade.py \
     --sqlite-file /absolute/path/to/verified-config-copy.db \
     --target-format-version 0
   ```

   此指令只讀、使用同一 SQLite snapshot，不建立不存在的 DB、不修改意圖、不部署。退出碼 **0** 僅代表 point-format 這一個維度可供 legacy reader；**2** 為意圖不相容；**1** 為證據不足／schema 或讀取錯誤。輸出只有安全計數，不列點位名稱、設定內容或 credentials。

4. 任何 override 均阻擋舊 reader。inherit 也保守阻擋：舊 reader 不理解 rule/connection 來源，即使當下 `data_format` 恰好相同，也不是持續等價的證據。未知 mode/version、損壞或部分 migration 均不能放行。
5. 非 0 **立即停止受管理回退程序**。保留新版服務／產物及設定，交部署負責人處理；不得清空 mode/override、把新點改成 legacy、刪 receipt 或強寫舊 projection 來過檢查。
6. 即使退出 0，仍須核對其他模組、schema、owned writer 與 backlog 相容性，並依原部署程序完成確認。**0 不等於整體版本可回退**；此檢查不控制啟動器、任意舊 binary 或其他 process。

## 備份與還原：不改歷史 SQL

使用 SQLite backup API／`.backup` 建立一致副本，或在所有連線確實關閉後保存所需 DB/WAL/SHM；不能只複製正在寫入的主 `.db`。本次演練用 Python `sqlite3.Connection.backup`，不新增自動 restore 入口。

還原驗證一律先寫到**新的可丟棄路徑**，執行 `integrity_check`，比對設定、mode/format/revision 與原備份；原目前 DB 保留不覆蓋。啟動新版副本時仍由既有 Migrator／runtime hydration 驗證，不假造 Ready 或重解歷史樣本。

若備份後已有 accepted sample、receipt、outbox 或外部提交效果，不能用舊 DB 覆蓋目前 DB。binary rollback 與資料 rollback 是兩件事；保留歷史 SQL、原 frozen source identity 與原目的地，絕不重新 INSERT／重新解碼歷史資料來「補齊」。本案檢查器不提供任何 destructive down／restore 功能。

## 此案驗證與交接矩陣

| 範圍 | 證據層級 |
|---|---|
| Modbus TCP、MC3E、FATEK D：6 numeric × 4 orders | cmd 真 source save → SQL reload → scheduler／manual → 原 adapter → ingestor，72 cases |
| Modbus UDP：6 × 4 | 相同 production chain＋真 localhost UDP peer，24 cases |
| Modbus RTU：6 × 4 | 相同 production chain＋僅 raw byte transport 模擬，client RTU framing／CRC 與 decoder 保留，24 cases；非實體 serial |
| FATEK DR | adapter 層支援，但 source authoring 明確拒絕；未擴 parser，不宣稱 DR source 成功矩陣 |
| zero／short response、legacy numeric／inherit16、reset／restart／rollback | 各既有 focused／integration tests；壞 response 不當成 Good 0 |
| scope／lease／SourceRevision／fingerprint、成功回執與未知結果 | 保存、recovery、typed sink、SQL member source identity、原 HTTP 反例均有相關驗證 |
| 048 expand、legacy 相容／modern downgrade 拒絕 | 原 Migrator tests＋此 script fixtures；readonly bytes 不變 |
| 一致設定備份、還原至新副本、不改歷史 SQL | 可丟棄 SQLite backup API 與獨立 history DB hash 比對 |
| 實機、正式 DB、任意舊 binary 啟動攔截 | 未執行／非本工具能力 |

測試入口：`python3 -m unittest discover -s scripts/tests -p test_point_format_downgrade.py`。Go 受影響測試與 test matrices 依各 checkpoint 記錄；七案完成前只執行相關測試，完整前後端總測留整批最後，不把本案交接等同產品全部驗收。
