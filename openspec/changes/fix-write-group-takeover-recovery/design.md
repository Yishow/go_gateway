## Context
已確認 canonical WriteGroup／journal／outbox／receipt production 接線。遷移未 disable legacy，Owns 依 active boundary；ResolveQuarantine 僅 Store 可呼叫。基準 d41638b。

## Goals / Non-Goals
**Goals:** 一次有效接管後永不意外復活 legacy，操作員可修復 blocked/quarantine head。
**Non-Goals:** 不重做 migration、不新 daemon／平台，不更改 frozen destination、不放寬 unknown／DDL／dedupe，不新 protocol／報表／counter 衍生能耗前置／多 tag event 模型。

## Decisions
### 持久接管與活動 worker 分離
在首次有效 canonical intake 接管時建立 durable connector/tag ownership（含 group／applied revision 來源），先持久化再接管。讀取錯誤 fail closed；worker disable/delete/retire/change destination 不釋放歷史接管。首次有效接管前與 failed apply 不奪走 legacy。重用既有 config DB／migration repository，避免平行配置模型。
### 以原 effect 身份修復 head
group scoped delivery attention rows 只回安全 code 與 opaque effect key；retry／skip 以 workspace/group/effect scope、expected state 做原子檢查，保留 actor/decision/time 的 audit。retry 沿用 frozen revision/destination/payload/hash/effect；skip 保留 payload 並明示資料未送達。unknown 拒絕 retry/skip，不推測外部 commit。

## Implementation Contract
1. production chain tests 先證明 retired boundary 後 legacy 回寫；GREEN 需 migration 後成功接管、disable、delete、目的地変更、重啟皆維持舊 writer 抑制；未接管／失敗接管仍允許原 legacy。
2. attention read、resolution mutation 綁當前 workspace 與 group，foreign effect／state stale／unknown 拒絕；mutation 原子记录 resolution provenance。重複请求不重送，state 轉移與 audit 同 transaction。
3. UI 顯示 blocked/quarantined 安全原因，修好後重試或明確確認 skip；unknown 只顯示需查證。按鈕 pending/failed 可見，不洩 credentials/DSN/raw exception；使用現有 GroupDeliveryStrip、service/hooks/i18n。
4. missingtable、permission blocked、badrow poison head 的 focused tests 驗證 next row 被釋放及 frozen identity 不改；managed receipt dedupe 仍成立。
5. 串行 RED→GREEN→review→fix→rereview；Go -p 1、GOMAXPROCS=2、Vitest 一 worker。repo required gates／browser／重驗證前先主對話協調；只 owned disposable DB 與 simulator。
6. scope 限 takeover/recovery modules、直接 API/UI calls、相鄰 tests／schema migration；無正式 DB/PLC/push/merge/archive/deploy。git diff --check、line gate、scope gates、最後 exact HEAD review，未做項目明列。

## Risks / Trade-offs
- 首次接管持久化時序 crash → transaction／restart 回歸，持久化失敗不開 intake。
- 只依 enabled group 查 ownership 會復活 legacy → 永久接管紀錄與 active worker 分離，disable/delete/restart tests。
- 外部 commit 未知 → unknown 不可 operator retry/skip，沿用 receipt reconcile。
- head 修復 race → expected state/CAS 與 audit 同 transaction，foreign scope／雙請求 tests。
