# 本輪起草 Review 與檢查紀錄

日期：2026-10-03。基準：`d2c22ef625903264ec6d74f6fd4bf76c5a48f916`。

## 一次 review 的範圍

本輪為同一助理的一次完整起草 review，不是獨立 reviewer、產品 code review 或現場驗收。檢查六案 proposal/design/tasks/delta、上輪十個問題的唯一責任、依賴與既有 schema/ownership/quality/Share 合約。其後只做該輪修正及結構重驗，不再展開新的全面盤點。

## 核對結果與已納入修正

| 核對點 | 結果／修正 |
| --- | --- |
| 上輪問題是否遺漏 | R1-R6、U1-U4 均有唯一 owner，F 只作整合驗收；沒有以隱藏多 entity UI 代替 C 的修復 |
| 離線恢復是否誇大舊版能力 | A 明定舊版缺描述只有原 source/mapping/connector 可證明相同時才可補證；否則 blocked 並保留資料，不猜 schema 或宣稱已支援離線新 intake |
| 單一 device_id 是否誤套進階跨設備表 | D 的新 metadata 保證限新基本 per-device managed 表；advanced 仍按原逐 member/列身分，不填假的共同設備 |
| 為快出第一列而改資料語意 | E 保留既有 60 秒預設與完整 UTC 桶；F 的較快故障 fixture 需在 UI 明選 10 秒，不改 production default |
| 計時是否只記數字便算完成 | F 明定兩 DB 各三次、每次不超過 300 秒才通過受控首次使用條件；不得刪 correctness assertions 或預建目標表換取達標 |
| Schema 與 start 責任是否重疊 | D 管確認式 schema；E 只協調已確認結構的 start，沿用同 ledger 並區分 action，不自建 token/job/queue |
| Custom 及 Share 是否被新基本預設破壞 | 保留既有自訂表 metadata 合約、進階編輯與 Share-only gates，不 ALTER 使用者表、不自動拆/轉群組 |
| 是否把規格交付當成功實作 | 50 項 implementation/verification tasks 皆未勾選；不 archive，不寫產品通過或部署聲明 |

這次 review 後未保留已知的阻擋性草稿範圍／責任衝突；這不代表官方 CLI、完整 repository 或產品已驗證。

## 實際工具與檢查界線

- 起草時 GitHub 已讀取 main ref、原始碼／既有規格來源與 changes inventory，以基準 tree 準備只新增文件的提交；先前寫入第四案物件時被工具安全檢查阻擋，該次未建立 commit 或更新 main。使用者再次授權提交後，重新讀回 main 仍是本輪基準；保留此歷史紀錄，實際提交結果以 Git 歷史為準。未接觸開發機器 WIP。
- DevSpace 連線失敗，Desktop Commander 無線上設備；容器無法解析 GitHub/npm registry。沒有假稱已在使用者 checkout 執行。
- 參照 OpenSpec 官方 spec-driven schema 與 repo `.openspec.yaml` 的 schema/created 格式；未使用 Spectra 冒充使用者指定的 OpenSpec CLI。
- `npm exec --offline --yes --package=@fission-ai/openspec -- openspec --version` 回 ENOTCACHED；正常 npm exec 取套件因 EAI_AGAIN 失敗。因此官方 `openspec validate --strict` 為 **NOT RUN**，實作前必須逐案補跑。
- 本輪只對待提交文字建立隔離 mirror 做 diff、行數與結構檢查；不是完整 repository build/test。結構檢查不宣稱等價於官方 OpenSpec validator。

## 文件檢查結果

- `git diff --cached --check` 與 `git diff --check`：PASS，範圍為待提交的 35 個新增文字檔。
- `make check-lines`：使用與 main blob SHA 完全相同的 Makefile、script、ignore，於隔離 mirror 執行 PASS；repo 的 `openspec/*` 豁免使原 gate 只檢查兩份總覽/review 文件，未冒稱它覆蓋所有 specs。
- 額外完整草稿行數檢查：35 檔皆低於 300 行，最大 89 行；不修改原 gate 或 ignore。
- 明確標示為自訂的結構檢查：6 個 change、9 份 delta、15 個 requirements、61 個 scenarios、50 個未完成 tasks；schema/created、必要章節、WHEN/THEN、ADDED description 長度、連結與依賴 DAG 通過，無競爭 requirement owner。
- 5 個 MODIFIED requirement 標題對照已讀取的既有規格來源；未變更的情境和合約保留，沒有直接改寫 `openspec/specs/` 或 archive。
- 與本輪新檔相連的 15 個 Markdown links 已檢查：新檔在 mirror 可解析；兩個既有來源連結已透過 GitHub 讀取確認。不是全 repository 文件連結掃描。

## 未執行

產品程式修改、Go/前端 build/lint/tests、實際 project race、UI→SQL、資料 migration、PLC/LAN/SCADA、正式 DB、部署與長時間 soak 全未執行。本輪沒有把上輪隔離重現或既有 archive PASS 算成本次結果。

## 提交與後續邊界

只允許本批六個 change 目錄與 `docs/plans/studio-v2-flow-completion/` 的文件進 commit；原 main 產品、既有正式 specs、archive、工具設定不改。提交前再查 main，使用單一 parent 的新增 commit，僅 fast-forward 更新；完成後核對 parent、changed paths 及 main head。
