## Context
目前 Basic managed group 與 recording-start 已 production 接通；F fresh-ui 從 Advanced 呼叫 GroupSchemaPanel 才能準備空白 DB。舊 connector write controls 與非 group schema setup 在同畫面造成雙設定誤解。基準 d41638b。

## Goals / Non-Goals
**Goals:** 四步主線完整 Basic schema、安全 capability、單一有效 group 設定與簡潔 operator 顯示。
**Non-Goals:** 不重建前端、不新 model/API/driver/protocol/daemon/report/event 製程模型；不放寬 DDL、unknown、receipt safety；不拿正式 DB 或 PLC 驗證。

## Decisions
### Basic 重用安全 schema panel
將 GroupSchemaPanel 接入 selected Basic group；沿用 read-only preview、revision/token/明確 confirm、operation ledger/recovery。ensure group 是 local configuration，不隱式外部 DDL；preview 後取消無建表，changed scope 失效。成功後只 refresh 原 group readiness，不創第二份 schema 狀態。
### 連線與生效記錄各有唯一設定
ConnectorSection 僅 connection identity/schema/credential；group 持有 members/storage/interval。移除 Step4 主線舊 nonpreview SchemaSetupSection 與舊 table/write strategy/interval 控制；若其他 engineering callers 必須保留，用明確 caller prop 隔離。KindSelector 根據既有 capability，只支援 group SQLite/Postgres；pool 不支援選項同樣 disabled。
### 主畫面與診斷分層
只直接必要調整：主畫面名稱、值、型別單位、記錄 stages/repair actions；raw IDs、revisions、API payload 等放既有 diagnostics/disclosure。保留真實 saved/applied/accepted/committed 區分與可行動 error，不把診斷 raw backend exceptions 放 operator DOM。

## Implementation Contract
1. Basic 空白 owned disposable SQLite 由選測量點、point→tag、managed group 到 schema preview/confirm/test/start 全程完成；Advanced 非必要。每點依型別佔 register span，不把八 words 當八 points。
2. Preview 明列實際 DDL/目的地作用域；confirm 必须顯式，cancel/missing/stale/foreign token 不建表；沿用既有 Backend tests 不放寬。
3. canonical group 是寫入 authority；主線不再提供會409的 nonpreview schema action 與不生效的 connector table/UPSERT/interval。Unsupported MySQL/SQLServer 不可選，安全說明 en/zh-TW 並覆蓋 pool selection。
4. Main display 可見名稱／值／type+unit／recording state；technical identity/revision/payload 在 diagnostic disclosure，測試可驗展開後仍可追蹤。同址設備隔離由 fix-basic-value-device-integrity 負責；recovery 動作由 fix-write-group-takeover-recovery 負責，不重複更改其邏輯。
5. UI 先 focused Vitest RED，再 GREEN、scope review→fix→rereview。重驗證前主對話協調，串行低 workers，不開 Docker/browser/build/full suites 直到時段確認。只本地 commit。
6. git diff --check、line gate、affected tests；協調後 repo required checks 與 simulator F acceptance：八量測點定時一列、uint64 預設精確、同址雙設備、停用重啟、legacy migration、DB fault recovery。最後 aggregate exact HEAD review。真 PLC/Windows/LAN/長跑未做明列。

## Risks / Trade-offs
- panel 與 Basic async state scope race → 切 group/device/connector token invalidation focused tests。
- 重複設定造成錯誤預期 → 主線只留實際 group authority，舊 caller 明確隔離。
- disclosure 不小心隱藏錯誤 → operator stages/repair/error 保留主畫面，focused DOM assertions。
- Mac concurrent load → 明列需协调的 gates；不自行啟重壓驗證。
