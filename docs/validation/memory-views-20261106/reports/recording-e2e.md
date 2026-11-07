> 歷史報告閱覽副本：正規化本機絕對路徑與附件連結；保留當時結論與未驗項目。原始位元組及完整附件見 `local-only/original-evidence.tar.gz`，不要將歷史 PASS 當成新版重驗。

# 完整記錄 E2E skill 與本輪實證

## 結論

已建立並安裝 `$validate-go-gateway-recording-e2e`，包含完整 UI 旅程、hc Terra 修復流程、Docker／HSL 隔離方式、證據清單工具與踩坑筆記。skill 官方格式檢查與 11 項 CLI 回歸通過。

實際從空工作區 UI 新增 MC 3E Binary、Modbus TCP 設備與點位，經記錄設定、schema 預覽及明確確認，兩個來源均已寫入 SQLite 與 Docker PostgreSQL，並對上真實 SQL、來源身分、applied／frozen revision 及雙方 receipt。每個 provider 一個群組含兩設備六個 int16 點，維持新草稿 latest／15 秒預設，沒有縮短週期。

**正常 UI→SQL 主線通過；整份擴充驗收清單並非全數 PASS。** mixed 試寫為現有能力限制，PG stop 只證明觀察期間穩定，未證明 ACK 即刻停止 intake；首頁背景重查／3～5 秒計時本輪未重跑。`candidate/checks.json` 為 36 PASS／4 BLOCKED／3 NOT_RUN，嚴格清單檢查預期 exit 1；不能把 skill 自測通過當成所有產品項目通過。

## 結果與證據

| 項目 | 結果 | 附件 |
| --- | --- | --- |
| UI 新增、取消不殘留、保存來源、連線準備、編輯重載 | MC／Modbus PASS | `terra/fresh-fixed/evidence/run.json` |
| SQLite schema 預覽無 DDL、明確建立 | PASS | `terra/fresh-fixed/sqlite-schema-*-request.json` |
| PostgreSQL schema 預覽無 DDL、明確建立 | PASS | `terra/fresh-fixed/pg-schema-*-redacted.json` |
| SQLite 真實 SQL | 3 連續完整桶，每桶六值 PASS | `parent-checks/sqlite-three-buckets.json` |
| PostgreSQL 真實 SQL | 3 連續完整桶，每桶六值 PASS，父代理重跑 | `terra/fresh-fixed/pg-parent-rerun.json` |
| 重載後持續寫入 | SQLite、PG PASS | `terra/fresh-fixed/*-before-reload.json`、`*-after-reload.json` |
| 停止後穩定 | SQLite 10 行、PG 7 行維持不變；有界觀察 PASS | `terra/fresh-fixed/*-stop-settled*.json` |
| PG ACK 即刻停止 intake | NOT_PROVEN；不推定所有後續提交均是 pre-ACK backlog | `terra/fresh-fixed/pg-stop-semantics.json`、`terra/REPORT.md` |
| mixed 試寫 execution／readback／cleanup | BLOCKED；preview 422 安全拒絕，沒有試寫成功 | `testwrite-analysis/REPORT.md`、`terra/fresh-fixed/*test-preview*.json` |
| 五記憶體視圖功能／int16 真值 | 兩協議共 10 案例；原值及恢復新值均 PASS，父代理另跑 10／10 | `memory-views/run-S8JJ2T/`、`memory-views/run-G3eXnY/` |
| 暫停、錯誤、新值復原 | MC／Modbus PASS；GET、SSE、UI、獨立 TCP 封包一致 | `memory-views/recovery-wpVEgh/`、`parent-checks/independent-wire.json` |
| 單獨 stale、手機／暗色、多 word／完整型別矩陣 | 本輪 NOT_RUN | `memory-views/REPORT.md` |
| 原圖外觀逐點重新核對 | 本輪 NOT_RUN，沒有修掉先前已回報的結構差異 | 前次 `revalidation-20261106-7a01e7d5/REPORT.md` 與 `reference-checklist.json` |
| 自有服務與 worker 清理 | PASS，Gateway／HSL／PG exit 0，自有 Docker container 已移除 | `parent-checks/cleanup.json` |

SQLite 完整 due 為 06:30:30／06:30:45／06:31:00 UTC；PG 為 06:42:00／06:42:15／06:42:30。啟動截斷窗不計。兩者均核對 device／point／tag、target column、source／mapping revision、sample provenance、record ID → outbox effect key → 內外 receipt 的 payload digest、connector revision 與 frozen runtime descriptor。這裡是群組 `wg_delivery_*`，不是 legacy `gw_delivery_receipts`。

來源原值 MC `123,0,-42`、Modbus `456,0,-17`。HSL 暫停約 7 秒後兩者皆 GET last_error／SSE quality bad／UI `!`；恢復後種入 MC D0=321、Modbus 40001=654，TCP 原始 request/response、Gateway GET 與 UI 都收到新值。原生 TCP 由父代理再實發一次，排除只讀 simulator 內部 register API 及 UI 快取。實測為 error recovery，沒有把它改稱 stale PASS。

PG 停用 ACK 後最後一列的 due／commit 晚於 ACK，但來源 observed／received 早於 ACK。沒有逐筆 accepted-at 的直接證據，因此不宣稱同步停止或契約違反；詳見 stop semantics 交接。試寫缺 entity key 是現有 mixed 能力邊界，規格要求無法證明 operation ownership 時拒絕；沒有移除 guard、修改表或新增未授權 capability。

## 本輪產品修正

透過 hc 的 `recording-e2e-terra-1106`（GPT-6.1 Terra、max）重現真實 MC 新保存範圍卡在準備中：connect／probe 成功但沒有 enable PUT。取得兩檔 ownership 後，先取得原版 genuine RED，再保留重新驗證的 proof 到 queue ACK，沒有移除 identity／revision／MC gate。

- `frontend/src/pages/datalink/workbench-v2/useRecordingRangeSubmission.ts`
- `frontend/tests/unit/workbench-v2/recording-editor-range-mc-prepare.test.tsx`（新增，含診斷間 revision 改變拒寫）

失敗候選與修復 RED 保留於 `failed-candidate/`、`candidate/repairs/`。修後隔離建置的新 binary SHA-256：`432307aba7a9347047bc0b4bc63e88c06a00fabf55cefe9c69454fd0ee9e7a1b`，base HEAD `7a01e7d5e0299e8709063cf01e0dc8e21bda097c` 加 11 檔 dirty overlay。父代理核對 binary、所有 static 檔，以及實際 HTTP 提供的 index／主 JS／CSS；與 manifest 一致。

其他 hc Terra 分別處理 skill 工具反例、試寫 422／父代理 SQL verifier 假設，以及五視圖／故障恢復；責任與報告分開保留。產品碼只有上述兩個本輪路徑新增修改。原有 9 檔 WIP 9／9 byte 雜湊未變、候選 overlay 11 檔最後未漂移、index 空，未提交或推送。

## 已執行驗證

- 前端完整測試：644 檔、7,391 tests PASS；lint PASS；本輪新候選 frontend build／Go embedded build PASS。
- 父代理另跑 MC／strict ACK／lease 22 tests、range prepare lifecycle／submission，皆 exit 0。
- `git diff --check` PASS；行數檢查 PASS，3 個既有 WIP 檔有 300 行警告，沒有 500 行阻擋。
- Skill 官方 `quick_validate.py` PASS；11 項 CLI 回歸 PASS。系統 Python 缺 PyYAML，使用 owned 暫存 venv，未改全域或產品依賴。
- 父代理直接唯讀 SQL 對證兩 provider、實發兩協議 TCP、重跑恢復後 10 視圖；報告及原始結果保留。正常產品 setup 都是 UI，API 僅觀察，未以 fixture API／SQL INSERT 代替。
- 秘密遮蔽 helper 自測 PASS；複製證據前檢查本輪 generated password 不存在於待交付文字，排除 private runtime env／設定 DB／build source。

## Skill 與使用

安裝：`$SKILLS/validate-go-gateway-recording-e2e/SKILL.md`。

新任務可呼叫 `$validate-go-gateway-recording-e2e`，指定協議與資料庫；依 `references/journey.md` 執行，遇錯用 hc Terra，依 `references/lessons.md` 避免已知陷阱。`scripts/evidence.py` 只檢查來源／清單／附件，不代替真實操作。

Skill 經驗已記錄空 DB revision、既有範圍啟用、背景重掛、MC proof 被清除、表單操作順序與 fill、DTO 形狀、mixed 試寫限制、真正 receipt authority、停用後歷史版本、simulator API 與 TCP 的區別、error 與 stale、串流 Live 與讀值成功的區別。

`terra/recording-e2e.cjs` 與 `recording-ui.cjs` 是本輪收斂的輔助 runner；**整合版尚未從新工作區一條命令重跑**。已實跑的來源 runner、PG draft／stop helper 及其餘逐段真實 UI 證據分別列明，不把語法檢查當完整 runner PASS。重跑前需重新建立隔離服務、空設定 DB／目的地、種值與核對執行檔；舊 `/tmp` 路徑及 ports 不是固定環境。

## 規格與剩餘界線

相關規格：repo `openspec/specs/studio-v2-device-to-sql-acceptance/spec.md`、`recording-database-setup/spec.md`、`recording-lifecycle-management/spec.md`、`per-member-recording-calculation/spec.md`；本輪沒有改寫規格、勾完成或 archive。

本輪是一個 fresh workspace、兩設備六個 int16 點、latest／15 秒、SQLite／PostgreSQL 的代表性旅程。不是正式 spec 的三次 fresh run、多 entity 列、bool／string／精確 uint64／全部型別驗收；也不是 SQL Server／MySQL、Windows、真 PLC／LAN／現場驗收。原圖外觀差異仍應依前次逐點報告處理，不能由這次功能與數值 PASS 覆蓋。

HSL／Gateway／Docker 的執行環境及憑證僅屬本輪 owned namespace。Docker 原啟動未加 run label，清理改以記錄的 hc argv、唯一名稱、建立時間、固定 container ID 與 loopback binding 核對；已自動移除 owned container。所有自有服務及 worker 已關；證據與輔助腳本保留，未動使用者開發服務或其他工作者。
