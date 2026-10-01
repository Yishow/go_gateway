# 跨案資料契約（擬議，尚未實作）

各案design/spec為可驗收要求；本檔統一名稱，避免交接時各自另造模型。

## A：WriteGroup是唯一設定authority

- `id`、`workspace_id`、`revision`、`applied_revision`：不可猜測的穩定ID與CAS版本
- `members[]`：每項保存`device_id`、`point_id`、`tag_id`、source/mapping revision、target column與required flag；point/tag須真實persisted且屬於workspace。`measurement_id`可選，不能用pointID代填
- `destination`：已存connector ID和identity revision、database/schema/table、managed或custom storage strategy、schema verification revision/digest
- `row_policy`：`interval_seconds`、`allowed_lateness_seconds`、每member的`max_age_seconds`、`incomplete_policy`、entity-key mapping、value/quality/provenance column映射
- `write_policy`：新history群組以穩定record identity追蹤append；target有verified dedupe時可安全重試，不具此能力的custom append明示limited並在uncertain後停止核對；latest/upsert只對有verified unique key的顯式選項開放，不推測為timestamp-only
- `migration`：來源legacy IDs、review result、adapter version；不把legacy records改名後就宣稱等價
- `operation`：保存/套用由workspace、group、connector expected revisions共同gate，與既有settings revision/readiness token保持一致

基本raw group只要求來源型別與映射正確，不要求semantic_kind、unit、report或retention表單。選擇usage/aggregation等進階能力才引用已確認的measurement definition，這批不新增那些能力。

群組rename保持stable ID及既有delivery identity；member/column/policy/destination編輯先存draft，explicit valid Apply於下一bucket邊界切換，不改寫舊bucket。disable停止新intake；delete採tombstone，保留accepted payload/revisions/receipts/ownership與查詢直到安全解決，不重用ID、不cascade丟backlog。legacy格式無法保護歸屬就拒delete並提供disable，實體purge不在本案。

## B：Sample與snapshot

`sample_id`、`acquisition_id`、workspace/device/point/tag、source revision、`observed_at`、`received_at`、`time_origin`、value type/encoding、value、quality及reason一起傳遞；可重用SampleEnvelope，但不能省略runtime需要的身分。協議沒有來源時間時，用backend真正完成該次採集的時間並標記origin，不把稍後flush時間偽裝成觀察時間。

新basic群組為**週期snapshot，不是區間平均、電量或every-sample history**。interval為已存正整數秒；預設從現有設定沿用，無設定時建議15秒並在保存前展示。15秒是設計預設，不是性能量測或採集頻率。

UTC bucket為`[start,end)`，`start=floor(observed_at/interval)*interval`。在`end+allowed_lateness`關閉；預設lateness=0，非零必須保存明確值。每member選該bucket內最大的observed_at；同時刻用穩定sample_id作確定排序（字典序較大者），不能按到達次序覆蓋。到達關閉後的樣本記late，不重開已封存row、不悄悄修改歷史。每個applied group以time-driven tick關閉到期bucket；整個bucket完全沒有sample也記一次scoped no_data/skipped狀態，不能因未建立memory bucket而消失。

預設`incomplete_policy=skip_row`且所有enabled member為required：任一required缺值、bad、無法解碼或過期即不寫該row，留下group/bucket/reason。freshness以`end-observed_at <= max_age`，預設max_age=interval且絕不跨bucket補值。若顯式選partial，只有target nullable且可保存逐member quality/provenance時開放；否則拒絕設定。partial以SQL NULL及quality reason表示缺值，不用0或上一筆good值。

值保留bool、text、signed/unsigned integer及decimal；大整數與decimal用明確encoding穿過JSON/SQL，不經JS Number截斷。NaN/Inf或SQL overflow為bad／blocked，絕不宣告good。所有群組在本地durable envelope／journal保存member observation/quality與record_id/target identity的關聯。全good且完整的custom row不強制增加外部逐欄metadata或companion表；外部schema只需符合明確選擇的partial模式或dedupe策略。UI必須揭露品質證據只在本地、SQL表本身未攜帶逐欄資訊的限制。

`record_id`由workspace、group ID、group revision、entity key（無業務分區時用固定group scope）及bucket_start導出；不可只有timestamp。命名空間同時納入目標scope作destination effect key。同一group跨device的members是同一row的不同columns；不同entity/群組即使同bucket也必須不同row。共享同一column的legacy布局只有保證row identity不同才合法，不自動合成會覆蓋的wide row。

## C：接受、交付與重啟

- sample ACK在local journal transaction成功後；open bucket狀態能由journal/checkpoint恢復
- final row、凍結payload/destination/group revisions、outbox intent及consumed checkpoint同交易提交；不在target Exec前刪除唯一副本
- destination sender以verified unique effect key或transactional receipt把target row與receipt同target transaction寫入。receipt僅在local DB不夠防範外部commit回覆遺失
- 無法提供dedupe/receipt的custom表：未曾送出的明確失敗可retry；可能已commit則unknown／blocked，需核對同record，不無限重插
- 設定變更只影響新的intake；accepted backlog保留原revision及destination。credential失效則blocked待修復同identity，不能改送新endpoint
- quota warning/hard limit需設定且測量使用量；hard limit拒絕新ACK並揭露loss-risk／affected scope。已接收資料不可因retry次數到達而刪除
- poison永久失敗進quarantine並保留payload與安全reason；依同group/entity順序阻擋其後繼，其他partition繼續；operator明確處置後才推進
- retry/backoff、queue bytes、batch與shutdown timeout有設定上限；不在本提案宣稱吞吐/SLA。worker用transaction claim＋lease/fencing或等價exclusive ownership
- 狀態至少分collecting、local_durable、queued/retrying、blocked/unknown、sql_committed；verified僅由實際readback產生。缺少觀察結果顯示unknown，禁止nil error直接轉success

## D：test_write operation

同一operation ledger保存kind、operation ID、workspace/group/connector revisions、target scope、payload digest、expiry及claim；建表token不能跨kind使用。舊plan endpoint需解析至同group，無法一對一解析回422，不猜第一筆。

preview零target mutation。確認後才用production row codec/sender寫入帶operation專屬identity的test row，並在相同scope核對實際values/types，再只清理operation自有row／metadata。無法安全標識或清理的custom表須在preview揭露並阻擋自動試寫，不廣泛刪除。

保留舊契約的`written_verified / written_unverified / failed / unknown`，獨立`cleanup_status=not_attempted / cleaned / failed / unknown`；清理後保存durable receipt，重送不能重寫。running=202；保存結果=200；並行或stale/expired首次claim=409；缺欄=400；foreign/unknown=404；wrong kind/unsupported=422。已保存operation的重送先查結果，不因token過期再次執行。
