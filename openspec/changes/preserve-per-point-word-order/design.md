## Context

基準 6f4a5fc 的 schema.Point 已有 DataFormat；schema.SourceRuleLink 只有 rule/address/point identity，沒有 override 意圖。sourcerule.Service.Update 使用 next.DataFormat 更新所有保留的衍生點；collector 的 scheduled/manual read 傳 point.DataFormat。
Modbus executeRead 使用 request→connection；MC3E.executeRead 使用連線 dataConverter；FatekConnector.Read 使用 convertFatekValue，按 component width 拼接，忽略 request format。此差異須以 production Read fixture 驗證，不能以 mixed-layout helper 存在作為完成。
point-catalog、source-rule-runtime、source-rule-target-datatype 與 archived optimize-point-configuration-ux／rework-datalink-rule-persistence-and-point-tag-flow 已掃描；本案唯一擁有逐點格式語意，streamline-recording-task-editor 消費本契約，不重做 decoder。

## Goals / Non-Goals

**Goals:** 同設備不同點可持久化不同格式，範圍編輯不抹除；讀回可辨識來源、有效格式及是否適用；舊設定升級不改變既有讀值。

**Non-Goals:** 不新增 driver、不改歷史 SQL／量測、不改 scale／unit／排程、不重寫 mixed-layout、不新增部署／測試環境；不做新編輯器視覺布局；不改協議寫入 wire semantics。

## Decisions

### 點位意圖是唯一 authority

在 canonical Point 設定增加 format_mode（legacy／inherit／override）與明確 override 格式，並以版本辨識新契約；rule link 只作 stable ownership，不另存 competing override。舊 DataFormat 保留相容輸入／projection，不以值等於規則預設推測 inherit。新點預設 inherit，override 接受 ABCD/BADC/CDAB/DCBA；reset 是明確 inherit 操作，未傳欄位則不變。API 讀回 mode、configured_format、effective_format、effective_source、applicable、point/rule/workspace revision，來源未知時不虛構有效值。

### 舊版語意以 legacy 保留

未標記舊點維持 legacy：Modbus 保留 point→connection→原預設；MC3E 保留 connection→原 CDAB 預設；FATEK 保留原 component-width 拼接（不能說它原本採 connection）。保存的是解析政策，不是猜測人工意圖或把當下連線值固定成 override。無關修改、restart、規則更新不得偷偷轉成新模式；明確改字序／reset 才轉換，回傳前後 effective 差異供 UI 確認。舊 API 明確傳 data_format 時視為顯式修改、經相同 guard；省略保持，空值為 reset。legacy 與新 API 的同一修改不可由不同路徑繞過 guard。只有 operator/API boundary 的顯式輸入才能轉意圖；sourcerule materializer 與 rollbackUpdatedPoints 必須走保留／還原意圖的內部路徑，不可把內部 DataFormat projection 誤認成使用者覆寫。

### 規則衍生以 identity 保留且版本化

multiword inherit precedence 是 rule 非空格式→connection 非空格式→協議原生預設；standalone point 沒 rule 則從 connection 起。覆寫依 workspace/device/rule/point ID 及 canonical address identity 綁定，不依顯示名稱或表格列序。保留點 ID 的 rename/count/skip/enable/default 更新保留 override；新地址 inherit；被移除點不自動把 override 借給日後新建同址點。skip→unskip產生新identity時inherit；start-address滑動時只保留重疊且identity仍在的點，不按index搬移覆寫。datatype/width 變更使 override 不適用時整個設定修改拒絕，要求明確 reset，不能丟掉意圖。

### 共用正規化後逐協議解碼（依實際讀取路徑保真）

**Supersedes**: preserve-per-point-word-order / 共用正規化後逐協議解碼

僅澄清舊段單word敘述的適用路徑，multiword契約不變。2026-10-30查核：`lib/hsllogic/register_conversion.go` 的 `DataConverter.RegistersToValue` 對int16／uint16直接回傳首register（signed cast除外）；`internal/datalink/connector/adapters/mc3e_read.go` 與 `internal/datalink/connector/adapters/modbus_parse_helpers.go` 使用此路徑。`lib/hsllogic/bytetransform.go` 的 `ByteTransform.TransformUint16` 才在直接byte呼叫時交換BADC／DCBA。把raw 0x0102的adapter讀值預期為0x0201是已推翻的假設，不是應實作的RED目標；相容回歸應鎖定0x0102，direct helper仍為0x0201。

沿 collector→ReadRequest→正式 adapter.Read 傳不可變 resolved format；不得改共享 connection converter，避免不同點互相污染。Modbus TCP/RTU/UDP、MC3E、FATEK 都以 fixture 實測。32-bit／64-bit signed/unsigned/float 是設計範圍；bool、16-bit 與 string 不顯示多暫存器 selector，拒絕新的 override，legacy 仍按舊解碼。單word的「不適用」只指多word selector，必須按實際driver／呼叫路徑保留語意：legacy及新inherit的Modbus／MC3E register-value路徑，int16／uint16維持原register值（int16僅重解signed bits），不因BADC／DCBA新增byte swap；FATEK維持原生component讀值。直接byte helper ByteTransform.TransformUint16既有BADC／DCBA交換bytes的行為另行保持，不推論所有adapter Read都會交換。不新增16-bit override、不強制ABCD或其他format metadata、不擴張解碼行為。ABCD 是高字在前，BADC 交換每個16-bit字內bytes，CDAB 反轉整個 scalar 的16-bit words，DCBA 同時反轉words與字內bytes；64-bit 不可只交換首兩字。FATEK 先依原 component width（16/32）抽出16-bit words，再依上述 scalar semantics 重排，不改ReadCount單位或地址跨度。

### 原子更新與安全失敗

新增 scoped point-format mutation，攜 workspace_id/device_id/rule_id/point_id、expected_workspace_revision、expected_rule_revision、expected_point_revision、request_id、format_mode、configured_format。standalone 不帶 rule_revision。沿既有 workspace ownership/CAS 邊界補 point revision；同 request_id+payload 重播回原結果，異 payload 衝突。403 foreign scope、409 stale/conflict、422 unsupported combination、503 uncertain persistence；無成功回應前不得宣稱 runtime 已套用。持久化與衍生 revision/candidate 失敗時 rollback，狀態不確定則停用受影響採集並安全回報，不假裝繼續正常。舊 observation/late responses 不覆蓋新revision顯示；格式變更沿既有 rule revision/readiness invalidation，不弱化 Share gate。measurement.SourceRevision 與 AcquisitionConfigFingerprint 必須加版本及format_mode、effective decode policy、rule/connection default來源identity/revision；不能只hash舊Point.DataFormat。collector installPointLocked／BindPointObservation／poll snapshot共用同份resolved metadata。例：MC legacy儲存ABCD但有效CDAB→overrideABCD，雖舊DataFormat沒變，新fingerprint必須不同；舊lease/proof/late sample不得綁新設定，歷史accepted sample保留其舊版本，不重新解碼。

## Implementation Contract

- In scope：內部設定 additive migration、rule/point/API/collector/現有 adapters，以及 frontend shared types/services；Out of scope：上述 Non-Goals 與其他 UI change 所有布局。
- 新點 inherit 與已存在 legacy 必須可區分；新rule預設與pointoverride均可讀回。UI 使用者可看到「沿用規則／沿用連線／此點自訂／舊版相容」及實際有效格式，legacy無法可靠判定時顯示相容政策，不給假的ABCD。
- 不以同格式值認定inherit。例：規則ABCD、A override ABCD、B inherit；規則改CDAB後A仍ABCD、B改CDAB。reset A後才變CDAB。
- 32-bit uint32 raw words [0x0102,0x0304] 在ABCD是0x01020304，在CDAB是0x03040102；另依既有 plcAdapterWireVectors 驗證 BADC/DCBA 與64-bit符號、float。FATEK D/DR各自覆蓋component count。
- 缺word、unknown format、unsupported type、foreign/stale mutation均不得產生 QualityGood 的虛假零值；回typed safe error，不洩漏endpoint或raw exception。
- 測試矩陣：各現有協議×六種multiword numeric type×四格式；同連線交錯A/B、scheduled/manual、restart、rule update/rematerialize、legacy fixture、reset、rollback、重播、跨scope/stale、缺word。數字是設計條件，非既有測試通過數。
- TDD：先重現覆寫流失及MC/FATEK忽略request失敗，再最小實作、重構；每個task只跑受影響測試，整個功能完成再一次完整檢查。使用既有fake transport與可丟棄設定DB，不碰真PLC／正式SQL。
- 驗收必須經 cmd/test_ui production wiring 的rule save→point reload→collector→adapter→ingestor路徑，不只helper單測；記錄實際通過／未執行，不把proposal通過當功能完成。

## Risks / Trade-offs

- 舊資料無意圖證據 → legacy policy保留原語意，不猜backfill；長期相容欄位增加複雜度但避免錯讀。
- FATEK 32-bit component與MC count單位不同 → fixture鎖定wire count及word正規化，禁止重複倍乘。
- 同連線併發解碼污染 → per-request immutable converter，race相關測試。
- rollback不完全 → fail closed並保留診斷，不以新值覆蓋歷史。
- 範圍增加 → 本案只交格式能力，簡易UI依賴獨立change，不擴scale或mixedlayout。

## Migration Plan

Expand：先新增nullable mode／revision與相容reader，既有點缺值視legacy，舊有效解碼fixture不變。Migrate：新API建立inherit或override，明確編輯才轉換，逐協議接線並保持legacy分支。Contract：所有使用端經同一resolver後移除重複推導，不刪舊資料或強迫全量轉換。部署前使用可丟棄設定DB演練重啟與失敗；回退舊binary前若有新override則由受管理回退程序及相容檢查阻止不相容downgrade，提供已驗證設定備份還原程序；不聲稱可控制任意舊binary啟動，也不新增部署engine，不靜默丟新意圖。

## Open Questions

無阻塞產品決策；實作仍需以逐協議fixture證明以上能力，未達矩陣的組合不得宣稱支援。
