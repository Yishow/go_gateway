# 五項既有需求修復驗證

基準 main `524555f14cad3da9d7c0cd064851f263e148a32b`，執行環境 macOS arm64。執行證據與內容SHA256：`docs/releases/evidence/studio-five-review-2026-10-23/{manifest,observations}.json`；詳述命令、失敗原因及限制：`docs/releases/2026-10-23-studio-five-review-fixes.md`。

## 完整性、正確性及一致性

五個requirements、十二個scenarios均已實作且有測試。無另列Example標題；R1/R4直接列出的D3/34464/1/100000與D0–D7/40001–40008值皆出現在測試。沒有排除的scenario，無以純樣式或成本理由省略分支。任務完成以實際assert與執行紀錄判斷，發佈SHA／CI另查，不以checkbox冒稱已push。

| Requirement／Scenario | 實作位置 | 覆蓋與執行 |
| --- | --- | --- |
| Scalar point reads／Multiword MC point across entrypoints | `connector/point_read.go:8`，collector polling/manual及device read-now | `point_count_contract_test.go` 真MC 18 scalar入口組合，含int64 max與uint64 9007199254740993；R1 RED→GREEN，focused race PASS；真UI三桶SQL PASS。 |
| Scalar point reads／Direct batch and other protocols | 同helper；MC adapter原乘法未改 | direct Count2→兩uint32、4words；Modbus TCP/RTU/UDP與FATEK捕獲word count；完整Go test PASS。 |
| Confirmed template／Explicit confirmation | `MeasurementTemplateSelector.tsx:55`，沿既有measurement backend | selector unit＋真API preview false／false apply422／UItrue200；PASS。 |
| Confirmed template／Changed scope or failed operation | 同selector key/device fence、request seq及安全error | 切設備/範本、過期回應、preview/apply斷線/409/500 unit；瀏覽器真API與錯誤注入7cases；PASS。 |
| Confirmed template／Offline intent and reload | measurements hook明確always＋retry:false；selector載入保留DOM | production defaults hook測3cases含恢復online1100ms不重送；reload DOM回歸；真browser offline錯誤後re-preview PASS。 |
| Focused header／No verified database commit | `FocusedDeviceHeader.tsx:139`與RuntimeDashboardPage真streamState | focused-header-evidence 5cases：waiting、Share-only、stopped、failed、unavailable；真UI不宣稱即時SQL且返回既有Step4證據 PASS。 |
| Complete occupied range／Multiword occupied end | `RangeSummary.tsx:32`既有parser＋最後span | 四float32為4項8words，D0–D7及40001–40008；RED→GREEN。 |
| Complete occupied range／Other widths and radices | 同上，不更改point生成 | occupied-range13cases含64-bit、word/bit、XF→X10；既有plain range一併PASS。 |
| Basic required／Required gap remains blocked | `BasicRecordingPanel.tsx:344`提示與現有editor窄入口 | basic-required-policy測saved有效required數與unsaved預設；Go snapshot required即使partial仍skip；真MC mapping-failed下required bucket skipped。PASS。 |
| Basic required／Optional gap and recovery | 沿既有required/completeness及snapshot，不更改品質保護 | 真UI明確取消D6 required→partial→save→preview→confirm→apply；真SQLite D6=NULL且mapping-failed provenance；恢復後D6=0，PASS。 |
| Basic required／Edit existing exact managed members | `columns.ts:75`限定managed SQLite generated uint16/32 TEXT；GroupEditor傳scope | columns unit驗managed可用，custom、foreign name/type仍拒絕；真UI保存保留原assigned columns。PASS。 |
| Basic required／Confirm unchanged schema | `schema_preview.go:116`在digest/儲存前把nil正規化[] | schema_noop_json_test先失敗null，修後PASS；真UI可預覽空DDL、明確確認並套用新group revision。 |

路徑縮寫：Go `internal/datalink/`；UI `frontend/src/features/datalink/`及`frontend/src/hooks/datalink/`，測試在各Go package旁與`frontend/tests/unit/`。

## 最後執行狀態

- passed-current：完整 `go test ./...`、`go vet ./...`、`golangci-lint run ./...`（0 issues）。
- passed-current：frontend lint、typecheck、215 files / 1378 tests、build；Playwright18。
- passed-current：最新binary之真MC socket＋臨時SQLite完整流程；三桶assert嚴格檢查非NULL、數字型別及unsigned精確TEXT；optional/recovery、頁面errors空、owned程序正常exit0。
- passed-current：最新binary template錯誤恢復7cases，無page errors。
- passed-current：diff whitespace、檔案行數及Spectra validate/analyze（無Critical/Warning；abstract scenario具體例子為Suggestion）。
- not-run：PostgreSQL真DB、Windows原生、現場PLC/正式DB、真人新手驗收。它們不在本次五修復的現場執行範圍，不當成已測通過。

## Review及修復

獨立review覆蓋R1 wire count、token/digest/持久化、R2–R5狀態與error分支，以及fixture隔離。找到並修好：header evidence URL缺focus、connected文案不準、ReactQuery預設retry可能重送、SQL NULL與0經JS減法混淆的測試assert。新增no-op預覽/managed TEXT修復均經獨立review；未發現剩餘產品缺陷。範圍含新檔，由task6 explicit_files補齊，不只依git diff追蹤舊檔。

一致性：保留現有四步與backend service，不新增報表平台、監控架構或新流程抽象。空preview仍需token、operation、scope/revision與確認。required未暗改；custom、歷史欄位、precision、unknown不重送等保護未放寬。

最後task用於歸檔準備，Git提交/推送/main遠端/CI及原工作樹同步由交付時實查。Windows active change不合併至本change亦不歸檔。
