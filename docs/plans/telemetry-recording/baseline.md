# 基準、衝突與來源

## 已核對的主分支

main = 4762f566ea28a1dbb4bbe313b156462d4af52517；與本對話先前檢查一致。
該提交歸檔 Studio V2 設定與 MySQL 目標修正。當次 `openspec/changes` 只有 archive，沒有同名進行中 change。
本組只加提案與規畫文件，不改 canonical specs、不開啟 PLC 或資料庫操作。

## 讀取證據（全部以該 SHA 為基準）

- AGENTS.md、CLAUDE.md、.agents/skills/openspec-propose/SKILL.md、openspec/config.yaml。
- point-catalog、source-rule-runtime、edge-analytics、database-target-workbench、database-output-delivery。
- datalink-workbench-v2-step1-device 的 Continue gate；Step4 的 connector、auto-assignment 契約。
- frontend/src/features/datalink/workbench-v2/state/types.ts 與 types-settings.ts：現有規則、映射與DB模型。
- Step4Database、TargetMappingTable、dbSchemas、autoAssignTargets：前一輪同SHA確認的欄位呈現與fallback。
- internal/datalink/dbtarget/writer.go：WriteTagValue、grouped map、flush與單點外部寫入路徑。
- internal/datalink/dbtarget/service_mysql_inspection.go 與 2026-09-06 MySQL修正歸檔的檢索證據。
- docs/technical/studio-surface-inventory/START_HERE.md、context.json、CURRENT_STATE.md：路由和產品定位。

以上是原始碼／規格審查，不是執行中系統或實機測試；不能以檔案存在推定功能全數驗收。
部分inventory仍寫多設備handoff待接，與更近期同SHA程式實作不完全同步；不把那些舊敘述直接當缺陷結論。
原README／AGENTS偏重SQLite/PostgreSQL，而main有MySQL修正；採能力矩陣，不抹掉既有MySQL進展。

## 明確處理的規格變更

| 舊要求 | 本組調整 | Delta 擁有者 |
| --- | --- | --- |
| 範例連線首次即ready、sample schema自動配對 | 初始未設定，實際metadata與確認配對 | 5 |
| 欄位不足仍wrap分配造成conflict | 留unmatched，讓人修正，不製造錯誤配對 | 5 |
| 所有設備成功才可繼續 | 明確所選scope全部成功；未選仍是草稿 | 5 |
| 所有資料庫輸出都由grouped mapping規畫 | 新managed plan；legacy/custom grouped仍保留 | 2 |
| 資料庫僅列兩種kind | 既有能力保留，新模式逐kind驗證後顯示 | 2 |
| delivery只看目前enabled | 新intake仍如此；已耐久接受的record不因後續disable消失 | 4 |

source-rule candidate、ownership、probe、Share readiness、原API安全要求保持適用。
新plan草稿與明確套用不改舊API原行為；不以新API繞過workspace activation barrier。
每個 modified requirement 只有一份 change 修改，避免六份提案互相覆寫。

## 外部一手資料（查閱2026-09-07）

- [OpenTelemetry Metrics Data Model](https://opentelemetry.io/docs/specs/otel/metrics/data-model/)：
  借鏡量測、累積／區間與時間身份的區分；本組不是照搬OTel傳輸模型或導入該套件。
- [Inductive Automation Mitsubishi PLC文件](https://support.inductiveautomation.com/hc/en-us/articles/16517576753165-Understanding-Mitsubishi-PLCs)：
  其driver文件說明word與跨word資料型別；現場仍以PLC程式／型號手冊確認，不硬編他牌driver限制。
- [PostgreSQL Constraints](https://www.postgresql.org/docs/current/ddl-constraints.html)：
  使用穩定unique key作資料身份保護；其他adapter需以各自交易能力實測，不能類推全數支援。
- [OpenSpec](https://github.com/Fission-AI/OpenSpec)：spec-driven artifacts與strict驗證的工具來源。

頻率、保留期限、容量門檻、效能目標與錯誤保守政策是本提案的設計決定，不是假裝外部標準的固定規定。
