# 規格與測試對照

32 個 Scenario、5 個 Example，均分類為 **covered by a test**；沒有 test-scope exclusion 或 uncovered gap。這是 source-level 對照，實際執行範圍／結果見 [validation.md](validation.md)。

## 測試縮寫

| 縮寫 | Repo-relative 測試檔案 |
| --- | --- |
| Bits | `lib/hsllogic/plc_byte_order_regression_test.go` |
| Wire | `internal/datalink/connector/adapters/plc_read_decoding_regression_test.go` |
| MCCount | `internal/datalink/collector/point_count_contract_test.go` |
| FCount | `internal/datalink/collector/fatek_point_count_contract_test.go`、`internal/datalink/connector/point_read_test.go` |
| FValue | `internal/datalink/connector/adapters/plc_fatek_decoding_regression_test.go` |
| FStatus | `internal/protocol/fatek/status_payload_test.go`、`internal/datalink/connector/adapters/plc_fatek_status_regression_test.go` |
| FAcquire | `internal/datalink/collector/fatek_status_acquisition_test.go` |
| Raw | `internal/datalink/connector/adapters/plc_modbus_address_regression_test.go` |
| MCPlan | `internal/datalink/sourcerule/mc_address_planner_regression_test.go`、`internal/datalink/sourcerule/protocol_address_test.go` |
| MCGuard | `internal/datalink/sourcerule/mc_radix_compatibility_test.go`、`mc_radix_runtime_reconcile_test.go`（同目錄） |
| MCPoint | `internal/datalink/point/mc_hex_address_test.go` |
| Parser | `frontend/tests/unit/utils/addressParser.test.ts` |
| Derive | `frontend/tests/unit/workbench-v2/sourceRule.test.ts` |
| Hydrate | `frontend/tests/unit/workbench-v2/mc-persisted-hydration.test.tsx` |
| Queue | `frontend/tests/unit/workbench-v2/mc-persisted-mapping-queue.test.tsx` |
| View | `frontend/tests/unit/workbench-v2/mc-persisted-view.test.tsx` |
| Draft | `frontend/tests/unit/workbench-v2/mc-persisted-draft-preservation.test.ts` |
| Locator | `frontend/tests/unit/workbench-v2/studio-locator-mc-hydration.test.tsx` |
| Next | `frontend/tests/unit/workbench-v2/mc-next-free-address.test.ts` |
| Bounds | `internal/datalink/sourcerule/mc_address_bounds_test.go`、`frontend/tests/unit/workbench-v2/mc-address-bounds.test.tsx` |
| Ready | `frontend/tests/unit/workbench-v2/mc-persisted-readiness-regression.test.ts` |
| NoWrite | `internal/datalink/sourcerule/mc_validation_boundary_test.go` |

## protocol-connectors：18 個 Scenario

| Scenario | covered by a test：具體函式／案例 |
| --- | --- |
| Multiword MC point across entrypoints | MCCount `TestMCScalarPointCountAcrossEntrypoints`；Wire `TestPLCReadDecoding_MC3EAdapter` |
| Direct batch and other protocols remain unchanged | MCCount `TestMCDirectTypedBatchCountIsPreserved`、`TestOtherProtocolPointWordCountsRemainUnchanged`；Wire `TestPLCReadDecoding_MCDirectTypedBatchCount` |
| FATEK scalar count follows native width across entrypoints | FCount `TestFatekScalarPointCountAcrossEntrypoints`、`TestPointReadCountForAddress` |
| Invalid FATEK scalar address does not produce a request | FCount `TestFatekInvalidScalarAddressFailsBeforeReadAcrossEntrypoints` |
| Four byte orders decode the same float32 | Bits `TestPLCByteOrderRegisterDecoding/float32`；Wire MC／Modbus adapter tests |
| Four byte orders preserve a complete uint64 | Bits `TestPLCByteOrderRegisterDecoding/uint64`；Wire MC／Modbus adapter tests |
| Signed and float64 values retain bits | Bits `TestPLCByteOrderRegisterDecoding` 的 int32／int64／float64；`TestPLCByteOrderSingleWordRegisterDecodingUnchanged` |
| Float32 has the same bits in word and doubleword areas | FValue `TestPLCFatekDecoding_RealReadTypedValues`：D_float32／DR_float32 |
| Complete 64-bit values in both native widths | FValue 同函式 D／DR 的 int64／uint64／float64 |
| A typed scalar cannot be completed | FValue `TestPLCFatekDecoding_RejectsIncompleteOrUnsupportedValues`：D／DR int64／uint64／float64 short、D_float32_short、unsupported |
| Explicit raw holding-register request stays a holding-register request | Raw `TestPLCModbusAddress_ReadExplicitFunctionThroughAdapters`；parser explicit raw register case |
| Raw boundary offsets are preserved | Raw parser 0／65535、`TestPLCModbusAddress_RejectsInvalidRawAddresses` 與 read wire cases |
| Traditional and alias addresses retain prefix semantics | Raw `TestPLCModbusAddress_ParseAliasesAndTraditional` |
| No-function fallback remains compatible | Raw 同函式 no-function 1／100／0／10000／65535；overflow/no-wire case |
| Persisted short numeric points expose the corrected precedence | Raw `TestPLCModbusAddress_PersistedShortNumericPointPreservesIdentity` |
| Malformed discrete data is rejected | FStatus `TestFatekStatusPayloadRejectsMalformed`、`TestPLCFatekStatusPayloadRejectsMalformed` |
| Valid scalar and batch responses remain supported | FStatus `TestFatekStatusPayloadAcceptsValid`、`TestFatekStatusPayloadRejectsInvalidCountBeforeWire`、adapter valid cases |
| Malformed payload cannot become a valid acquisition sample | FAcquire `TestFatekMalformedStatusPayloadAcrossEntrypoints`：3 malformed payloads × periodic／manual／read-now |

## datalink-workbench-v2-step2-rule：14 個 Scenario

| Scenario | covered by a test：具體函式／案例 |
| --- | --- |
| Derive points for MC 3E device | Derive 的 MC D0、D100 stride、X0、FATEK prefix cases；MCPlan 的 W0/count17 create／persist cases |
| Backend address offset matches protocol radix | MCPlan `TestOffsetAddress_MCUsesCorrectedRadixAndRuntimeOffsets` |
| Backend rejects invalid protocol address | MCPlan invalid create／update fail-before-persistence，含 `rule_id`、`start_address` assertion |
| MC 3E hexadecimal offset crosses a digit boundary | MCPlan corrected radix vectors；Parser `should advance MC3E addresses using area-specific radix` |
| Invalid protocol address fails closed | MCPlan invalid create／update；Parser invalid radix cases；Derive invalid protocol returns empty；View／Step2 rule-scoped error |
| Hexadecimal suffixes survive point creation | MCPlan W0/count17 persist；MCPoint `TestServiceCreate_AllowsMCHexWordAndBitAddressShapes`；Parser expansion |
| Metadata update cannot recreate legacy points | MCGuard `TestServiceUpdate_MCRadixLegacyRangeConflictPreservesCompleteFixture`：name／scale；Hydrate conflict mutation gate |
| A range patch does not authorize legacy migration | MCGuard 同函式 range patch；`TestServiceUpdateWithRuntimeReconcile_MCRadixConflictDoesNotRestore` |
| Startup uses persisted points | MCGuard `TestServiceSyncDerivedPointState_MCRadixLegacyLinksRemainAuthoritative`；Hydrate scoped identity case |
| Compatible rules retain normal updates | MCGuard `TestServiceUpdate_MCAddressPlanCompatibleRuleStillUpdatesNormally`；Hydrate ordinary update case |
| Reload does not reinterpret legacy addresses | Hydrate scoped point/mapping identity；View actual identity／legacy range cases；真實 browser reload |
| Unavailable or stale evidence blocks synthetic success | Hydrate 403／409／503／missing／incomplete／pending／missing revision／server row drop；View retention／unavailable cases |
| Summaries do not assert a fabricated continuous range | View legacy range／unavailable／disabled summary cases；Step2／occupied-range regression；真實 legacy browser |
| Scope changes revoke writable evidence | Hydrate old-scope response／workspace revision／late ACK cases；Queue in-flight response crosses scope case |

`D100/count4/int16` 沒有完全相同的單一 fixture；D0/count4/int16 與 D100/count2/int32 分別驗證 prefix 與 stride 行為，分類為同一行為的覆蓋變體，不宣稱已有該 exact fixture。

## 5 個 Example：33 列

| Example | 列數 | 分類與逐列對應 |
| --- | ---: | --- |
| Float32 word sequences | 4 | covered：Bits／Wire 的 ABCD `3F80 0000`、BADC `803F 0000`、CDAB `0000 3F80`、DCBA `0000 803F` |
| Uint64 word sequences | 4 | covered：Bits／Wire 的 ABCD／BADC／CDAB／DCBA 獨立 uint64 vectors，完整 `0x0102030405060708` |
| Preserved prefixed addresses | 8 | covered：Raw parser 的 00001／10001／30001／40001／HR100／IR100／C100 與 CO100／DI100 |
| Multi-protocol point derivation | 7 | covered：Parser Modbus int16 expansion；Derive Modbus float32、MC D0/int16、D100/int32、D0/float32、X0/bool、FATEK R0/int16 |
| Protocol radix boundary | 10 | covered：Parser／MCPlan 的 X0+16、Y7+1、W0+10、W0+16、WA+1、BFF+1、D0+16、M0+16、R0+16、40001+2 |

## Review 補齊的分支與邊界

不新增 requirement 或 scenario；以下補足上述原契約的 caller 與錯誤路徑。

| 原契約 | 補充 assertion |
| --- | --- |
| Unavailable or stale evidence blocks synthetic success | Draft：沒有 server record 的新 mapping／target 於 checking→unavailable→compatible 仍保留值、原 point ID；既有 persisted／Modbus drafts 同時不變。 |
| Reload does not reinterpret legacy addresses | Locator：cold MC mapping 與 editor hydration 保留 native point ID／X8／WA，不需要暫存 Step2 evidence，且來源型別不取 target type；非 MC 控制仍通過。 |
| Protocol radix boundary／point derivation | Next：X0/count8→X10、W0/count17→W11、WA→WB、BFF→B100；legacy persisted W0…W16→W17，缺 evidence 不臆造空位；Modbus／FATEK fallback 控制。 |
| Invalid protocol address fails closed | Bounds：0xFFFFFF 最大值、超界字串、offset／expansion overflow、Create 無 repository writes。Ready：new／persisted local invalid start 與 range 均 blocked，legacy conflict 保留 native identity。 |
| Hexadecimal suffixes survive point creation | Ready：WA／WF→W (Word)、BFF→B (Bit)，全字母 suffix 不會變成 area；X8 native legacy identity 不變。 |
| Compatible rules retain normal updates | NoWrite：W0 compatible rule 改 X8 時 ErrValidation、零 restore writes、完整 fixture 不變；實際寫入後注入 ErrValidation 仍有 snapshot compensation。 |

全部仍為 covered by a test；最終測試執行與真 PLC／Windows／正式 DB 的 NOT RUN 分別記於 validation，不以 source inspection 取代實測。
