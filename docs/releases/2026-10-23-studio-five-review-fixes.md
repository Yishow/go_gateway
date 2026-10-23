# Studio V2 五項 review 修復驗證

基準：`524555f14cad3da9d7c0cd064851f263e148a32b`。本次只修 R1–R5，沿用四步、受管／自訂群組、型別精度、DDL 明確確認、operation ledger、版本及 unknown 防重送保護。

## 實作與可觀察結果

| 項目 | 修正 | 實際證據 |
| --- | --- | --- |
| R1 | MC 單點 Count 用 typed value 數1；保留 adapter 乘寬度與 direct batch。Modbus/FATEK維持words契約。 | 真MC線路測 periodic/manual/read-now 的32/64-bit scalar，含int64 max、uint64 9007199254740993；direct batch兩個uint32。真UI→SQL三桶。 |
| R2 | 真preview=false，只有明確確認送true；scope切換忽略舊回應；安全錯誤、reload；離線不排隊且不自動retry。 | 真API false preview／未確認422／確認200；真瀏覽器503載入、preview/apply各409/500/offline共7情境恢復；production QueryClient defaults測失敗及重連仍只送一次。 |
| R3 | Header只顯示既有availability/running/stream；SQL沿既有Step4證據入口。 | 等待、Share-only、停止、交付失敗、unavailable測試；真UI可返回Step4的Basic記錄證據。 |
| R4 | 完整占用終點含最後項寬度，沿既有stride及addressParser，不改點位。 | MC/Modbus、32/64-bit、word/bit、十六進位進位13情境；4 float32為4項8words，D0–D7／40001–40008。 |
| R5 | 說明required優先、顯示必填數、窄入口開目前群組既有成員設定；不暗改optional。 | required失效略過；明確改optional＋partial、保存＋重新預覽確認＋套用後SQL為NULL、provenance=mapping-failed；恢復後SQL值0。 |

R5整合另修兩個既有阻礙：只對canonical managed SQLite uint16/uint32之generated nullable TEXT欄位沿用後端精確儲存契約，不放寬custom檢查；相容no-op preview的statements在digest與持久化前正規化為`[]`，保留token/ledger與明確確認，不放寬前端parser。

## 此次執行

全部為目前內容的 passed-current，非歷史報告推定。完整命令結果及來源SHA256見 [manifest](evidence/studio-five-review-2026-10-23/manifest.json)，實際SQL摘要見 [observations](evidence/studio-five-review-2026-10-23/observations.json)。

| 命令／範圍 | 結果 |
| --- | --- |
| `go test ./...` | PASS，全部packages |
| `go vet ./...` | PASS |
| `golangci-lint run ./...` | PASS，0 issues |
| R1 focused race與獨立review race重複5次 | PASS |
| `npm run lint`、`npx tsc --noEmit` | PASS |
| `npm test -- --run` | PASS，215 files / 1378 tests |
| `npm run build` | PASS；Vite既有chunk size提示保留 |
| `npm run test:e2e` | PASS，18 tests，含Studio/test/gateway及logs契約 |
| `five-review-flow.mjs`，目前嵌入UI＋真MC socket＋temp SQLite | PASS，三桶、required/optional/恢復與證據返回 |
| `five-template-errors.mjs`，目前嵌入UI＋真API/傳輸錯誤注入 | PASS，7 cases、0 page errors、owned gateway正常exit0 |

新測試先失敗再修：MC18項scalar錯、範本確認／scope／錯誤、Header5項、range8項、Basic2項；整合另重現離線暫停、繼承retry=1重送、managed TEXT誤判、no-op statements:null。失敗記錄保留本機，未刪除或當成PASS；成功結果由最後一次執行獨立確認。

重跑方式（只可使用隔離fixture）：先`npm run build`、同步`frontend/dist/.`至`cmd/test_ui/static/`並build `./cmd/test_ui`，再以`F_GATEWAY_BINARY`與絕對`F_EVIDENCE_DIR`執行兩個腳本。MC腳本自行建立loopback模擬器及temp DB並關閉自有程序；errors腳本須給該腳本產生的`F_EXISTING_WORK`，不得指向正式資料。

## 規格與review

五份delta共12個scenario均有測試對應，沒有以視覺樣式不測的例外；逐項映射見本change的verification.md。獨立review找到範本retry預設、證據連結focus參數、以及SQL null→0的JS測試誤判，均已修正並重驗。最終review未發現剩餘產品缺陷。

任務最後一項改為歸檔與發佈準備，避免以歸檔前打勾冒稱已commit/push；實際Git/remote/CI與原工作樹同步結果於交付時另行核實。

## 限制

- 本機macOS自動化驗收，並非真人新手驗收；未操作現場PLC、正式DB或部署。
- 真SQL只驗SQLite；PostgreSQL真連線與Windows原生tray/lifecycle/logs仍not-run，不因cross-compile或Studio E2E標PASS。
- R1線路矩陣以CDAB驗mixed32/64-bit；其他byte order沿既有測試，不宣稱本次逐組真設備驗收。
- MC三桶含累積uint32 100000與單位/縮放；64-bit精度由真MC線路測試確認，本次未宣稱64-bit完整SQL三桶。
- Windows active change保持未完成。CI狀態須按實際提交SHA查核；本機測試不代表GitHub CI通過。
