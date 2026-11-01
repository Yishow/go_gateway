# 採集預覽的最新已保存版本重新讀取

本修復承接原採集預覽 consumer 及整體驗收，不新增 change、planner、telemetry registry 或 editable head 權威。QA29 schema-only 成功不等於這個獨立 Refresh 行為原已完成。

## 原問題及原因

QA29 明確 schema Apply 後正常回 Step1，acquisition-preview GET78 仍帶 Apply 前的 DB setup head，收到 409 並顯示 stale／Refresh。這不是刻意負向，也不是 schema-preview POST 失敗；當次沒有按 Refresh，因此不能宣稱已恢復。

原 ordinary preview scope 來自 connector.setup_revision。schema Apply 的 writeGroups／bootstrap／databaseMetadata invalidation 不會更新這個 editable DB owner；普通 bootstrap 沒有 DB setup head，而 databaseConfig hydration 另以 id／updated_at 去重。舊 Refresh 只遞增 refreshId，重送原 scope head，沒有重新取得 canonical DB head 的路徑。

## 最窄修復

- 只有 ordinary、persisted／saved selected-device 的明確 Refresh，先用既有 write-groups 唯讀服務及 strict DTO parser，核對 workspace 並取得有效 canonical DB head（保留原合法 legacy 空 revision）。
- 該 head 只是這次 preview 的 detached read receipt；不寫 connector CAS、setup revision、autosave、draft、dispatch 或 query cache，不從 409 或 Share／runtime-context head 盲目 rebase。
- 初次讀取仍使用 supplied exact scope；locator／default panel 保留原 exact authority。scope 改變、suspend、unmount、較新 Refresh 會取消或拒絕較舊回覆。
- Native cached preview head 必須同 owner 及 resolved DB revision 才重用；R→R2 不能借舊 native head。後端 DB-head、source lease、installed-only 與零 PLC I/O 契約未改。
- 對未 installed 的來源，真 200 installation_unavailable 是正確結果，不補 available／0 或假實測。

## 已執行驗證

原碼三個真 RED 後，五個 acquisition 相關測試檔 68／68 PASS，六檔 ESLint／精確差異檢查通過。涵蓋 R→R2、foreign／malformed workspace、late scope／suspend／unmount／較新 Refresh、locator exact、read failure、legacy empty head；真 Step1 dirty connector 測試驗零 dispatch／cache write。

主工作樹第一次 typecheck 被同時開發的 3.2 未完成測試阻擋，未宣稱通過；QA30 用已驗 baseline362 加六檔修復的精確 union364，乾淨 typecheck、Vite／static sync／Go embed build 通過，沒有混入 3.2 WIP。

獨立六檔 frozen-packet review 無具體 findings、6／6 hashes 相符。Reviewer 的 CLI 截斷管線失敗保留；root 另以完整 JSON 執行 scope/check，snapshot 79185000e4dbd499b8f07138abdec8c275f85677、exit0。這不冒稱整個 change 已審完或已歸屬全部新路徑。

兩次最初 focused Vitest 曾短暫與 scheduled d10/p128/r1 重疊，原共享負載證據保留，該矩陣案例另有唯一 quiet 重跑；沒有刪除原 loss／timing 結果或重跑其他已有效案例。

## 真 Browser closure

QA30 原自有 fixture27 pair 在同 session 明確 rename Save53／200，canonical DB head 由 c9cc2574… 變成 46fa5319…。原 schema 仍有效，因此沒有硬造 noop DDL。這輪 head advance 是 rename Save，原 QA29 才是 schema Apply，兩者不混稱同一操作。

正常 Step1 自動 preview67／409 仍使用舊 head，真 DOM 顯 stale。明確點 Refresh 後依序 GET69 write-groups／200（完整 canonical body 已保留）→ preview70 使用同新 DB head、同 workspace／device，返回 200 installation_unavailable；DOM 改為真 unavailable，沒有 fake available 或數值。

沒有 reload、Start、connector Save 或另一個無關 Save 來掩蓋 Refresh 恢復。Refresh 前後七個設定表完全一致；connector 與前置 Save 前的 identity／CAS row 也相同。沒有 DDL，target 備份 hash 不變、兩表仍各 0 rows。內部 cache／dispatch 不變是單元／靜態證據，不冒稱 Browser 曾直接檢視 React state。

最後導航成功，console 無 React #185／uncaught，原67／409已明確恢復；外部 font 拒連另列。自有服務 21:53:58.816 UTC 停止，其他28份舊DB及 source364 hashes 符合，備份 integrity 通過。

## 證據與後續

- preview-refresh-evidence 保留四個 baseline 原檔（核對 QA29 manifest）、兩個原缺 new files、六檔 freeze、task-only patch、全部 logs 及界線。
- QA30 報告 SHA256：deeac48ab7ce95dc9b715c0abf95d73e9378c49c22c73b040b2233a688fe254c。
- QA30 union364 SHA256：90a210a634be219a8f81f96f68cd9a2e694bd5197f99600e2ec25634d9970b66。
- 此缺口已閉合，原4.1整體終驗仍需在全部整合後保留對應回歸；不因此勾完 batching4.2 或全部七案，不宣稱物理PLC或 SQL資料寫入。
