## Purpose

Preserve explicit per-point multi-register decoding choices on one device across rule edits and runtime execution. Distinguish inherited defaults, explicit overrides, and legacy decoding policies without silently changing existing measurements.

## ADDED Requirements

### Requirement: Point format intent and effective order

The point configuration service SHALL persist distinct legacy, inherit, and override modes. New points SHALL default to inherit; explicit overrides SHALL remain explicit even when equal to the current default. Reads SHALL expose configured format, effective format or unresolved status, effective source, applicability, and revisions. Reset SHALL explicitly select inherit; omitted fields SHALL leave intent unchanged.

#### Scenario: Override equal to default survives default change
- **WHEN** a rule changes from ABCD to CDAB with point A explicitly overriding ABCD and point B inheriting
- **THEN** A SHALL remain ABCD and B SHALL become CDAB
- **AND** resetting A SHALL make it inherit CDAB without changing B

#### Scenario: Resolve inherited format
- **WHEN** an inherited multiword rule-owned point is resolved
- **THEN** the resolver SHALL select nonempty rule format, otherwise connection format, otherwise the protocol's native default
- **AND** a standalone point SHALL omit the rule step and report its actual source

### Requirement: Preserve legacy decoding semantics

Missing legacy intent metadata SHALL NOT be interpreted as proof of inheritance or explicit override. Upgrade, restart, unrelated edits, and rule synchronization SHALL preserve existing effective decoder policy until an explicit format edit or reset. Legacy Modbus SHALL retain point-to-connection fallback, MC3E SHALL retain connection decoding, and FATEK SHALL retain its component-width assembly. Historical recorded data SHALL NOT be rewritten.

#### Scenario: Existing MC3E point differs from connection
- **WHEN** a legacy MC3E point stores ABCD while its connection uses CDAB and the service upgrades or restarts
- **THEN** its decoding SHALL remain CDAB under the legacy policy
- **AND** only an explicit guarded override to ABCD SHALL change its effective decoding

#### Scenario: Existing FATEK ignored format remains compatible
- **WHEN** a legacy FATEK point contains a format previously ignored by its decoder
- **THEN** upgrade SHALL preserve original component assembly rather than activate that stored format
- **AND** the response SHALL identify legacy compatibility rather than falsely label it a user override

### Requirement: Preserve overrides across rule revisions

Rule synchronization SHALL bind intent to canonical workspace, device, rule, and point identity, not display names or row position. Updates retaining a point SHALL preserve its explicit override. New identities SHALL inherit; removed identities SHALL NOT donate overrides to newly created points at the same address. An incompatible type or width change SHALL reject the mutation until the operator explicitly resets the affected override.

#### Scenario: Resize and rename a rule
- **WHEN** a rule is renamed, resized, disabled, or regenerated while preserving point identities
- **THEN** surviving point overrides SHALL remain unchanged through persistence, restart, and runtime projection
- **AND** added points SHALL inherit and removed-point overrides SHALL NOT leak to future identities

#### Scenario: Skip and range slide preserve identity boundaries
- **WHEN** an overridden point is skipped then recreated by unskip, or the range start address moves
- **THEN** a newly created identity SHALL inherit and SHALL NOT receive the removed point's override
- **AND** surviving overlapping identities SHALL retain their own overrides regardless of row index

#### Scenario: Type change would discard override
- **WHEN** a multiword point with an override is changed to an inapplicable type
- **THEN** the whole mutation SHALL fail with an actionable validation error and no partial persistence

### Requirement: Protocol decoder parity

Production scheduled and manual reads for existing Modbus TCP, RTU, UDP, MC3E, and FATEK SHALL honor independently resolved point formats without mutating shared connection converters. Supported multiword numeric types SHALL include signed and unsigned 32-bit and 64-bit integers and 32-bit and 64-bit floating point. ABCD SHALL preserve words and bytes, BADC SHALL swap bytes within each 16-bit word, CDAB SHALL reverse all 16-bit words of the scalar, and DCBA SHALL perform both operations. FATEK component normalization SHALL preserve its existing address/count units.

#### Scenario: Different orders on one connection
- **WHEN** points A and B on one connection read raw words [0x0102, 0x0304] as uint32 using ABCD and CDAB respectively
- **THEN** A SHALL decode as 0x01020304 and B as 0x03040102
- **AND** interleaved reads SHALL NOT change either point's effective order

#### Scenario: Four-word value and component width
- **WHEN** a uint64 CDAB point reads normalized words [0x0708, 0x0506, 0x0304, 0x0102]
- **THEN** its value SHALL be 0x0102030405060708
- **AND** FATEK 16-bit and 32-bit component paths SHALL produce the same normalized scalar without altering their required component count

### Requirement: Applicability and safe decoding failure

The service SHALL expose applicability by protocol, address component width, and datatype. New multi-register overrides on bool, single-word integers, strings, unknown formats, or unsupported combinations SHALL be rejected rather than ignored. Existing legacy configurations SHALL retain their original behavior. Insufficient raw data SHALL yield bad quality and a safe actionable error, not a successful zero or truncated value.

#### Scenario: Inapplicable new override
- **WHEN** an operator submits a new CDAB override for uint16
- **THEN** the mutation SHALL be rejected before persistence and existing intent SHALL remain unchanged

#### Scenario: Single-word byte semantics remain compatible
- **WHEN** a legacy or newly inherited single-word point is read
- **THEN** its actual driver and call-path semantics SHALL remain unchanged: Modbus and MC3E register-value reads SHALL preserve the int16 or uint16 register value, with signed interpretation for int16, without introducing BADC or DCBA byte swaps
- **AND** FATEK SHALL retain native component decoding; direct ByteTransform.TransformUint16 calls SHALL retain their existing BADC or DCBA byte swaps without transferring that behavior to register-value adapter reads
- **AND** hiding the multiword selector SHALL NOT force ABCD or other format metadata, discard existing data-format metadata, or enable new single-word overrides

##### Example: Register-value reads differ from direct byte helpers
- **GIVEN** a single register 0x0102 and BADC or DCBA metadata
- **WHEN** legacy or newly inherited Modbus or MC3E uint16 reads use RegistersToValue
- **THEN** the adapter value SHALL remain 0x0102, while a direct TransformUint16 call on bytes [0x01, 0x02] SHALL retain its existing 0x0201 result

#### Scenario: Short response
- **WHEN** a numeric decoder receives fewer words or components than its declared scalar requires
- **THEN** it SHALL return bad quality without publishing a fabricated good value

### Requirement: Scope and revision guarded updates

Format updates SHALL verify ownership and expected workspace, point, and applicable rule revisions. Identical request identifiers and payloads SHALL replay the original result; reused identifiers with different payloads SHALL conflict. Persistence, rule candidates, and runtime invalidation SHALL remain consistent or fail closed. Old API inputs SHALL use the same guard and resolver, with omitted values preserving intent and explicit empty values resetting inheritance.

#### Scenario: Stale or foreign edit
- **WHEN** a format mutation carries a stale revision or foreign point identity
- **THEN** the service SHALL reject it without changing stored or runtime formats
- **AND** it SHALL return a safe conflict or ownership error without endpoint credentials

#### Scenario: Internal projection is not operator intent
- **WHEN** a source-rule materializer or rollback path writes a point format projection
- **THEN** it SHALL preserve or restore the original intent and SHALL NOT interpret that internal write as an explicit operator override

#### Scenario: Retry and partial failure
- **WHEN** an identical mutation is retried after its first success
- **THEN** it SHALL return the same outcome without duplicate revision advancement
- **AND** a failed persistence or projection update SHALL rollback or mark affected acquisition unavailable until consistent state is recovered

### Requirement: Version decoding identities and observation proofs

Source revisions and acquisition fingerprints SHALL version and include format intent, effective decoding policy, and applicable default-source identity and revision. Collector installation, observation binding, and polling snapshots SHALL use the same resolved metadata. Previously accepted samples SHALL retain their original version without reinterpretation.

#### Scenario: Identical stored format changes effective decoding
- **WHEN** an MC3E point storing ABCD changes from legacy connection-based CDAB decoding to explicit ABCD
- **THEN** its source revision and acquisition fingerprint SHALL change even though the old DataFormat string is unchanged
- **AND** old observation leases, proofs, and late samples SHALL NOT bind to the new configuration

### Requirement: Production evidence and consumer handoff

Completion SHALL require fixtures through production rule persistence, point reload, collector, adapter Read, and ingestion paths across the supported protocol/type/format matrix, plus restart, legacy compatibility, reset, scope, stale revision, retry, and rollback cases. The simple recording editor SHALL consume the same contract rather than introduce UI-only format state. Tests SHALL use fake transports and disposable configuration databases, not real PLCs or production databases.

#### Scenario: Helper passes but runtime omits point format
- **WHEN** a helper unit test passes but a production adapter or collector path still ignores the requested point format
- **THEN** acceptance SHALL remain incomplete and the consumer SHALL NOT claim that combination is supported

### Requirement: 已知格式保存須有顯式且受版本保護的恢復入口

系統 SHALL 在既有 reconcile 操作提供明確的 point-format recovery intent，使用原操作 receipt 與保存後的 scope/revision 證據，從 server canonical source 重建指定點的有效讀取設定。此操作 SHALL NOT 隱含啟動 recording 或 Share listener；普通 activation SHALL 繼續要求恢復後的 fresh readiness token。系統 SHALL NOT 要求每次合法格式修改都重啟整個服務。

#### Scenario: 保存後恢復與開始採集分離
- **WHEN** 使用者對 invalidated receipt 明確要求恢復，且同一 point/source/default identity、clean witness 與 lifecycle generation 仍成立
- **THEN** 系統恢復指定點的 profile/installation，回傳可供原 activation 再檢查的 fresh readiness 證據，而不隱含開始 recording

#### Scenario: 同 revision 的 unknown 撤銷恢復
- **WHEN** restore 計畫後發生同 revision MarkDirty，或任何 settings/membership/pending/failed generation 改變
- **THEN** 包含 aligned/replay 在內的所有成功出口拒絕舊 attempt，不清除 foreign dirty、不發布 Ready、不覆蓋新 installation

#### Scenario: 成功操作重播與停用 Share
- **WHEN** 同一恢復操作已成功且當前 ready 證據仍一致
- **THEN** 重播唯讀返回原 token，不重安裝或發新 token
- **WHEN** Share 已停用且合法 active point 要求恢復
- **THEN** 只恢復該 point runtime，明確表示 Share 未投影且 listener 未啟動，仍拒絕 settings epoch 變更
