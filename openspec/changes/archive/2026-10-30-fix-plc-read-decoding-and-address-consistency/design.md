## Context

本提案承接 `docs/HslDocs_Categorized/`、`lib/` 與實際採集路徑的參數查核。已以本機模擬回應重現五類問題，未連線真 PLC。提案建立時基準為 `76a853d71fda23ee3a1d3cf9bf801e45f85f24f2`；實作前仍須重查 diff 與來源。

| 已重現問題 | 現行入口／符號 | 觀察結果 |
| --- | --- | --- |
| BADC/DCBA 遺漏 byte swap | `lib/hsllogic/bytetransform.go` 的 `RegistersToUint32`；`lib/hsllogic/register_conversion.go` 的 `registersToInt64` | float32 的 BADC words `803F 0000` 被解為約 `-5.785636e-39`，而非 1 |
| FATEK typed 解碼 | `internal/datalink/connector/adapters/fatek_helpers.go` 的 `convertFatekValue` | D0 的 `3F800000` 回 int 16256；DR0 回約 `1.0653532e9`，而非 float32 1；uint64 只剩首 word |
| MC radix 差異 | `offsetAddress`、前端 `AddressParser` 與 `lib/hsllogic/address_parser_mitsubishi.go` 的 `parseMitsubishiAddress` | W0 的舊展開包含 W10，runtime 實際 offset 為 16；X10 實際 offset 為 8 |
| Modbus Function 被覆蓋 | `internal/datalink/connector/adapters/modbus_parse_helpers.go` 的 `parseModbusAddress` | Address 100 + Function 03 的 wire request 是 FC01、offset 99 |
| FATEK malformed Cmd 44 | `internal/protocol/fatek/client.go` 的 `ReadStatus` | 空、短及含 X 的成功框架 payload 被當作有效結果 |

上述數值為查核實測。下列測試向量、數量上限與進位表為設計契約，不代表真機驗收。暫存重現檔不作為交接依賴；實作時須在 repo 建立可重跑測試。

MC 現有 root spec 明定 X/Y/B hex、D/W/M decimal；這與 runtime parser 相反，必須用 MODIFIED delta 明確變更。`SyncDerivedPointState` 不重跑 planner，但 `sourcerule.Service.Update` 會按文字 address 比對並清理舊 point/mapping，因此只改 planner 會引入資料關聯回歸。

## Goals / Non-Goals

**Goals:**

- In scope：五類已確認錯值／錯位址／錯誤成功問題，以及完成修復必要的 FATEK scalar count、MC hex 建點與舊規則保護。
- 維持既有 connector、ReadRequest／ReadResult、point→tag 架構；成功值型別與錯誤品質可由呼叫端直接驗證。
- 保留既有 connection defaults、point ID、tag/mapping 關係與歷史資料；揭露 Modbus raw precedence 的實際目標變化。

**Non-Goals:**

- Out of scope：FATEK station 0／255 政策、RT／RC 功能開通、MC PCNo／IONo 零值政策、新 PLC 家族／area／protocol／transport。
- 不新增 FATEK data_format、string／bitfield 解析、PLC 家族選擇器、地址版本欄位、自動資料 migration 或舊規則轉換 UI。
- 不改 MC runtime 的 area code、X/Y decimal fallback、port／timeout／station／byte order 預設值；不依 guide 的 Y/F code 表覆寫現有 decoder。
- 不修改 Python 範例或 HSL 套件、不引入 .NET／外部 dependency，不重算歷史值、不修復其他 runtime／recording WIP。
- 本次只起草；後續本機 simulator 驗收也不能代替真 PLC、Windows、正式 DB 或現場驗收。

## Decisions

### 共用 byte order 解碼保留既有 16 位元與預設值

將 `RegistersToInt32`／`RegistersToUint32` 與 `registersToInt64` 的完整 words 轉為每 word 高 byte 在前的 bytes，再重用已具四種 ordering 的 `TransformUint32`／`TransformUint64`，避免維護另一套不完整 switch。signed 與浮點由完整 bits reinterpret，不能做數值 cast。

32／64 位元的 BADC 都是每 word 內交換 byte；CDAB 是整個 scalar 的 word 順序反轉；DCBA 同時反轉 word 與 byte。ABCD/CDAB 的現有結果不變，單 word int16／uint16 的既有解讀不變。不更動寫入編碼，也不把此次修復擴成所有低階短 buffer API 的重設計。

替代方案：只在 Modbus adapter 補 swap 會留下 MC 錯值；複製全部 switch 則保留兩套語意，因此採共用既有 byte transform。

### FATEK typed scalar 與原生元件 count 分開處理

`convertFatekValue` 改為回傳 `(any, error)`，按 `ComponentType.Width` 取得足夠 bits：16 位元區需 1／2／4 個元件才能形成 16／32／64 位元 scalar；32 位元區需 1／1／2 個元件。沿用目前 Int32 的高 word 在前慣例，不新增可配置 endian。float32／float64 以 IEEE-754 bits 解碼，int64／uint64 保留完整整數，不經 float64 中轉；不足或未支援型別明確失敗，不退回首值。

在 `internal/datalink/connector/point_read.go` 新增 address-aware `PointReadCountForAddress(protocol, address, dataType) (int, error)`，供 periodic、manual、read-now 使用。先解析既有 FATEK 元件符號與寬度再算 count，DR/DD 等原生 32 位元 float32 讀 1 個、64 位元讀 2 個；D/R 的既有 word count、Modbus word count、MC typed count 維持不變。無效位址在送出讀取前失敗；沿用既有錯誤傳遞路徑。

舊 `PointReadCount` 在遷移期間保留；直接 `ReadRequest.Count` 對 FATEK 仍代表元件數，不 reinterpret 為 typed value 數。`FatekClient.ReadRegisters` 的元件數、長度與上限檢查保持有效。

替代方案：只修 float cast 無法處理 64 位元及原生 32 位元的多讀；全域把 FATEK Count 改成 typed count 會破壞直接 batch 呼叫，因此限 scalar 呼叫遷移。

### MC 新規則採現行 runtime 可確定解析的 radix

| Area（既有 planner 集合） | 新規則 suffix | 設計向量：offset 16 |
| --- | --- | --- |
| X、Y | octal，僅 0–7 | X0 → X20、Y0 → Y20 |
| W、B | hexadecimal，0–9/A–F | W0 → W10、B0 → B10 |
| D、M | decimal | D0 → D16、M0 → M16 |

前端 `AddressParser` 與 backend `offsetAddress` 同步 parse／format／validate／expand；W 的 type 維持 word，B/X/Y/M 維持 bit。MC 不能沿用「遇到第一個數字才切 prefix」的 `splitAlphaNumeric`，否則 WA/BFF 無法解析；只辨識表中既有單字母 area，再依 radix 解析 suffix。FATEK 的既有 parser 規則保持不變。

`point.validateAndNormalizeAddress` 保留目前接受的格式，僅增加 W/B 的 hex suffix 詞形（包括 WA、W1A、BFF），避免把任意字串視為合法。這是通用 point storage 的詞形檢查；MC area／radix 語意由 protocol planner／connector 負責，不新增區域支援。

runtime parser 的 X8／X9 decimal fallback 暫不改。新規則嚴格拒絕 X8、YA、WG、Z0；既有直接讀取仍有相容例外，不能聲稱所有入口接受集合完全相同。

替代方案：依舊規格重寫 runtime 會直接改變舊 X10/W10 的實際讀取位置；引入家族設定會擴大產品契約。此次保留已執行的物理位址語意，明列 planner 的 breaking change。

### MC 舊規則更新在任何 mutation 前阻擋隱含重建

在 `Service.Update` 現有 rule mutation lock 內，先以「儲存中的 start/count/data type/skipped addresses」按新 grammar 計算目前應有集合，再與 persisted links 的正規化 address 集合比較；不能用這次 patch 的新範圍來判定舊資料是否相容。

若新 grammar 拒絕舊範圍，或 address set 不一致，回傳包覆 `ErrValidation`、帶 rule/address context 的可辨識原因 `mc_address_plan_conflict`。guard 必須早於 `repo.Update`、orphan cleanup、point/tag/mapping create/update/delete、`ReplaceLinks` 與 revision 變動。即使只改名稱、scale，或 patch 明確提供新範圍，皆不授權隱含搬移；本 change 不提供 bypass。讓使用者看見既有規則需另外規劃轉換，禁止回報已儲存。

相容的新規則可沿既有更新流程運作；此 guard 不重設其他協議 lifecycle。startup/reload 的 `SyncDerivedPointState` 保持依 persisted links 同步狀態，不重算地址、不建立替代 points。檢查 healthy fixture 的 rule/link/point/tag/mapping IDs、address、關聯與 revisions；拒絕更新後必須與前一 snapshot 相同。

替代方案：靜默重新展開會失去 mapping；自動猜測舊使用者原意無法區分設計位址與實際位址。因此本次拒絕不相容更新，轉換另案授權。

### MC 已儲存點位以 scoped snapshot 作為前端事實來源

一般 Step2 hydration 在 `useStudioV2AutosaveState` 用 `deriveAllPoints` 重新生成點位；`useAllPoints`／`useRulePoints` 與 `RangeSummary` 亦重算位址。`buildBootstrapMappings` 以 rule/address 配對，所以只改 parser 會把既有 mapping 誤配到 synthetic points。

沿用 `sourceRuleCandidateAPI.get(ruleId, scope)` 的查詢 `GET /datalink/source-rules/:id/candidates`，傳入現有 workspace_id、expected_workspace_revision、rule revision_id；新 hook 使用自有 query key，至少包含 rule ID 及上述三個 scope 欄位，workspace/revision 切換不沿用舊結果；不得沿用目前只按 rule ID 快取且不傳 scope 的 `useSourceRuleCandidatesQuery`。前端不得為載入而呼叫 POST candidates/recompute 或任何 apply。既有 GET 的 `GetCandidateViewAtRevision` 在特定條件會刷新 database-output snapshot；此次保留該行為，不宣稱 endpoint 完全無副作用，也不把它擴成 MC tags 地址或 identity 的重建。回應的 `tags.candidates[]` 已含 persisted address、point_id、tag_id、mapping_id；以 rule 自身 data_type 建立 source point metadata，不能將 candidate 的 target data type 當 source width。

新增小型 `mcPersistedPoints` adapter 與 scoped query hook；local Rule state 持有對應 revision 的 persisted points，以及 checking／compatible／conflict／unavailable 狀態，不新增 server schema。對 persisted MC rule，只用此集合建立點位與對應 mapping；新草稿才依 planner 生成。取得 snapshot 後，比對新 planner 的集合，差異為 conflict；缺失、403/409/503、revision 不符或資料不完整皆為 unavailable，尚未取得則 checking。

checking／conflict／unavailable 的 rule 禁止自動 rule save、mapping 建立／重綁／刪除與後續以該假設為基礎的設定送出；其他設備及未變更草稿保留。conflict 顯示已儲存的位址與安全的需另行轉換訊息；無足夠資料時顯示無法確認，不能 fallback 成新 planner 結果或 saved/ready。使用 en／zh-TW 既有 i18n，raw backend exception 不進 operator DOM。

Step2 selectors／point grid／RangeSummary／occupied-conflict 檢查共用同一來源；對不連續或無法按現行語法解析的舊集合，不顯示推算的連續區間或「無衝突」，改用 persisted 列表及 blocked/unknown 狀態。scope/revision 改變先清除可寫資格，再重新讀取；不得讓先前 compatible snapshot 授權新 revision 的 autosave。

替代方案：只靠 backend guard 擋 Update，仍擋不住錯誤預覽及獨立 mapping mutation；新增 migration API 則無必要。既有 scoped GET 的 tags snapshot 足以提供事實，資料不齊時明確阻擋。

### Modbus raw numeric 依明確 Function 使用 zero-based offset

`parseModbusAddress` 的決策順序如下；read 與 write 共用此 parser，必須回歸驗證兩種 caller：

1. HR/IR/C/CO/DI alias：沿用 alias 的 area 與 zero-based offset，保留歷史上由 prefix 決定功能碼的行為。
2. 恰為 5 位、首碼 0/1/3/4 且後四位 0001–9999 的傳統編號：area 由 prefix 決定、offset 減 1；即使 Function 有預設值，也維持 prefix 優先。
3. 其他純數字且有明確 Function：以完整數值作 raw offset，合法範圍 0–65535，Function 不被推測 area 覆蓋。
4. 未提供 Function 的非傳統純數字：保留現行 fallback；1–9999 仍為 coil 的 one-based 舊格式，其餘沿用既有 FC03／範圍檢查。此次只修明確 Function 的歧義。

保持現有 Function 正規化與支援清單；負數、非數字、overflow 不送 wire。傳統格式範圍內的 40001 + FC04 仍按原行為讀 holding offset 0；要表達 raw offset 40001 不新增語法，該歧義保留為限制。6 位擴充傳統格式不在此次新增支援範圍。

既有 point API 會將空 Function 正規化為 03，因此儲存為短數字 + 03 的舊點位無法辨識當時意圖，部署新版本後實際目標會改變。此為明確揭露的 BREAKING，不承諾其物理目標不變，也不自動改 Address／Function／point ID。替代方案：繼續讓數值區間勝出會使明確 raw Function 無效；新增持久化 address-mode/migration UI 超出本次修復。

### FATEK Cmd 44 嚴格驗證後才形成成功值

在 `FatekClient.ReadStatus` 的 `execute`／frame status／LRC 驗證之後，要求 payload byte 數恰為 requested count，每個 byte 僅能為 ASCII 0 或 1；空、短、長、非 binary 全部回 error，不回部分 bool slice。保留 count 1–255 的既有設計上限。

adapter 將錯誤傳回 `ReadResult` 的 bad quality、無有效 Value，periodic／manual／read-now 沿既有 read-failed 路徑處理。可以產生 failure evidence，但不得產生 quality-good 的新樣本或以 false／nil 假裝成功；不更動 downstream recording policy。替代方案：補 false 或截短會把壞資料變成可採集數值，因此直接拒絕。

## Implementation Contract

**可觀察行為與介面：** ReadRequest／ReadResult JSON、point/rule schema 與 connection defaults 不變；只新增內部 address-aware count helper、FATEK 轉換 error、MC rule validation reason 與前端暫存的 persisted-point 狀態；candidates GET 沿用既有 API，沒有新增 endpoint。成功數值須為請求的 Go 型別，不能只驗證 JSON 數字看似相同。MC 新規則的 preview、derived point address、runtime 解出的 offset 一致；舊規則依 persisted 位址顯示及執行，不以新 planner 假裝已轉換。

**設計驗收向量：**

| 目標 | Input | Expected |
| --- | --- | --- |
| shared float32 | ABCD `3F80 0000`；BADC `803F 0000`；CDAB `0000 3F80`；DCBA `0000 803F` | float32 1；MC/Modbus 都經真實 adapter 解碼 |
| shared uint64 | ABCD `0102 0304 0506 0708`；BADC `0201 0403 0605 0807`；CDAB `0708 0506 0304 0102`；DCBA `0807 0605 0403 0201` | uint64 `0x0102030405060708` |
| signed/float64 | 含 sign bit 的 32/64 位元整數；float64 bits `3FF0000000000000`，四種 order | signed bits 保留；float64 1；ABCD/CDAB 不回歸 |
| FATEK | D0 兩個 16-bit components `3F80 0000`；DR0 一個 32-bit component `3F800000` | typed float32 1，wire count 分別 2／1 |
| FATEK 64-bit | D0 四個 components `0102 0304 0506 0708`；DR0 兩個 components `01020304 05060708` | 同一 uint64；所有 scalar caller 的 wire count 分別 4／2 |
| MC create | W0 count 17、int16；WA 作 start；BFF offset 1 | W0…W9,WA…WF,W10 全部可建點；WA；B100 |
| MC octal | X0 offset 16；X7 offset 1；Y0 offset 16 | X20；X10；Y20，runtime offset 分別 16／8／16 |
| MC conflict | 舊 W0 rule links 為 W0…W16，已有 tag/mapping；更新 name 或 scale | validation conflict；所有持久化 snapshot 不變 |
| MC reload UI | 舊 W0 rule + scoped persisted snapshot；缺失／stale snapshot | 顯示既有 W0…W16；不自動儲存、不產生或重綁 mapping；資料不齊時 blocked，不推算成功 |
| Modbus raw | Address 100、Function 03；0/65535、Function 04 | FC03 offset100；FC04 offset0/65535；65536 拒絕 |
| Modbus legacy | 00001、10001、30001、40001；HR100/IR100/C100/CO100/DI100 | 傳統 area offset0；alias offset100；既有 prefix 優先 |
| FATEK status | request count2；payload empty/1/101/1X；合法 01 | 前四者 error + bad quality；01 為 false/true；count1 回 scalar，count255 接受合法完整資料 |

每項修復先新增獨立預期值的回歸測試，在未改產品程式時確認失敗原因，再修到通過；wire fixture 不得用待測 production encoder 產生預期值，以免對稱錯誤互相抵銷。保留既有 direct MC batch、Modbus/FATEK 16-bit word count、正常 Cmd44/46、source-rule mapping lifecycle 回歸。

**驗證範圍：** 使用本機 simulator、memory repository／可丟棄 SQLite。後端至少涵蓋 `lib/hsllogic`、`internal/protocol/fatek`、connector/adapters、collector、device、point、sourcerule；前端涵蓋 addressParser 與 Step2 的新建、儲存重載、occupied-range／已儲存位址。依 `docs/testing.md` 執行選測、受影響 packages 的 vet/lint、前端型別與檔案 lint；使用真實瀏覽器及本機 simulator 驗證 Step2，不接 LAN PLC。沒有執行的檢查標為 NOT RUN，不由 unit test 推論 UI／真機通過。

**範圍與停止條件：** 不加入新 schema、轉換 UI、外部寫入或持續性資料更動；遇到必須修改範圍外 runtime/recording policy 才能符合上述契約，先提出具體失敗證據與最小依賴，僅暫停依賴項目。五類契約與必要回歸滿足後停止，不全面重寫 parser。

## Risks / Trade-offs

- [MC 舊規則不能直接編輯] → 以可辨識 conflict 保存原資料；既有採集繼續使用原 points，另案處理轉換，不默認刪規則重建。
- [X/Y family radix 差異] → 本提案只對目前 runtime 可確定解析子集合承諾一致；新規則 strict octal，舊 direct fallback 是已知相容邊界，並非所有三菱家族的通用宣告。
- [Modbus 舊短數字點位改變目標] → 明列 breaking，實作交接列出受影響形狀與 before/after wire；部署前另取得目標設備授權與盤點，不以此次提案授權現場切換。
- [byte order 修正改變過去錯誤值] → 驗證獨立 known-bit vectors 與 adapter 路徑，保留預設值；不回填或重新解讀歷史資料。
- [共用長檔限制] → 原有超長檔不得增行；新增 count／guard 行為抽成有責任邊界的小 helper，測試用獨立檔。僅搬移直接被修復的分支，不為行數作周邊重構；以 line gate 確認。

## Migration Plan

1. Expand：先新增回歸測試及內部 helper；每個行為先確認未修產品上的預期失敗，再在同一批完成修復與 focused green。
2. Migrate：先讓 MC guard 生效，再切換 backend planner／合法建點；前端 parser 與 persisted hydration、selectors／顯示保護全部完成後才算可交付；FATEK typed converter 就緒後遷移 scalar count 的三個 caller，直接 batch 不改。
3. Contract：所有 caller 遷移且測試通過後，移除本次取代的重複 word-order／首值 fallback 分支；舊公開 helper 若仍有使用者則保留，不做通用 API 清理。
4. 本機交付：測試、靜態檢查、UI simulator 證據與未驗項目分開記錄，不 commit／push／部署。後續現場 rollout 與短數字 Modbus 盤點需另外授權。
5. Rollback：此次無 persistent schema/data migration，可回復本次程式變更；新 grammar 建立的 WA/X20 等規則不能假設舊版本可編輯，回退前須盤點並停止相關設定編輯。不得透過反向 migration 刪除點位、映射或歷史資料。

## Open Questions

本次實作所需決策已定。FATEK station 0／255、RT／RC 及 MC 家族選擇均明確排除，不能在 apply 時自行補成產品功能；真機接受新讀取結果的現場判定不屬於本機完成證據。
