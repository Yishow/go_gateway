# 本次文件驗證

日期：2026-10-01。範圍僅proposal／design／delta specs／tasks／docs；沒有修改產品程式、生成前端或部署。

## 已執行

- 以fresh clone核對remote main `4c00af81fb51b3239b4e0d0d9c6b83680c808b36`，建立獨立工作分支，沒有覆蓋使用者工作目錄
- 暫讀repo最後原版openspec-propose（來源見evidence），隔離安裝官方 `@fission-ai/openspec@1.2.0`，實際version輸出1.2.0；未跑init/update、未安裝Spectra/spxa、未修改tracked skills或tool config
- 六案均執行真實CLI `new change`、`status --json`，依 `instructions proposal/design/specs/tasks --json` 的模板和依賴產生artifacts；proposal先寫，再design/specs，最後tasks
- 六案最終status的proposal/design/specs/tasks全為done，applyRequires=tasks；這只表示artifacts齊全。新增產品tasks共54項，全未勾選；新案9份delta specs、26條requirements、97個scenarios（文件計數，不是測試通過數）
- 六案與移交後舊active change分別通過官方 `openspec validate <change> --strict --json`，沒有issues
- 獨立review完成，已補offline draft navigation與完整group lifecycle/tombstone/backlog保護驗收；沒有擴成新change
- 自訂文件檢查核對48個變更檔的相對連結、scope、capability/spec一致性、MODIFIED原需求存在、task依賴無環、active delta唯一owner、原完成task與managed Spectra block未變；無問題。這是補充檢查，不取代官方CLI
- `git diff --check` 與 `make check-lines` 通過；staged diff一併核對；commit後按base SHA重跑完整批次
- 人工cross-change review核對authority、時間與row identity、retry/unknown、scope/revisions、完整/partial、legacy移交；完整good custom rows允許local provenance，不強制額外外部metadata；silent buckets由time-driven closure顯示no_data

## 檢查的邊界

- OpenSpec1.2.0沒有Spectra analyze命令，未執行也未宣稱其通過；跨文件一致性由source核對及獨立review補充，不能把CLI結構檢查稱為產品正確性證明
- 原active change保留5個原完成task文字；5個未完任務移交D/E/F，另留未完closeout gate。沒有新增任何產品完成checkbox或archive
- Go test/vet/golangci-lint、Vitest、frontend build/lint、Playwright、真實SQL、simulator／field tests本次未執行：只有文件變更，這些驗收是後續A–F的實作任務
- Windows／ARM、external/embedded browser、LAN、真PLC、SCADA、production DB、長時間性能與現場部署均未驗證

## 重跑命令

需已安裝官方OpenSpec1.2.0；此命令不會實作或archive：

```sh
for change in \
  unify-studio-v2-write-group-contract \
  enforce-write-group-sample-semantics \
  wire-durable-write-group-delivery \
  implement-confirmed-write-group-test-write \
  simplify-studio-v2-four-step-setup \
  validate-studio-v2-device-to-sql \
  fix-studio-v2-database-workflow
do
  OPENSPEC_TELEMETRY=0 openspec status --change "$change" --json
  OPENSPEC_TELEMETRY=0 openspec validate "$change" --strict --json
done
git diff --check
make check-lines
```

最終git publication SHA與review結果於交付時核對；文件內不預先寫入尚不存在的commit或CI通過結果。


## 2026-10-03 實作與最後交付

原10-01僅文件的紀錄保留；此段記錄後續使用者授權依A→F實作驗收的結果。六案54/54tasks完成，舊案五項移交需求與5.1closeout核對完成；未部署或新增commit。

- Go一般／tagged完整test、vet、lint均PASS，兩lint0issues；前端lint／181files1090tests／build／Playwright13/13PASS；harness24/24PASS。
- SQLite／PostgreSQL真UI→Modbus→SQL七型別、品質／quota／disk／poison／CAS／worker／重啟／lostcommit／ownedtestcleanup皆各自通過；Linuxamd64容器有獨立source/build/UI/SQLPASS，清理錯誤修後SQLite與PG另補驗PASS。
- decimal依使用者確認只驗codec/SQL，不宣稱設備採集支援；field/正式DB/部署/Windows/nativeARM/nativeLinux硬體/embedded/LAN/PLC/SCADA/longsoak仍NOTRUN。
- 最後F scope snapshot `04662a850f0dbdcc92b3ffa234c8e0989c5501e5`（243paths，44binary不作程式內容驗證）；舊案 `1e4cc71cda9ecc0595caac591ab82cc8076b0099`（85paths）均current。四lens review／verify完成，無confirmed未修缺陷；dirtybaseline僅能證明檔案歸屬。
- F與舊案preview均無未完成task、無warnings；corearchive各執行一次，specs各套用一次，cleanupwarnings皆空。實際archive位置如下；links於移動後調整，沒有改寫需求語意。
- [2026-10-03-validate-studio-v2-device-to-sql](../../../openspec/changes/archive/2026-10-03-validate-studio-v2-device-to-sql/validation.md)；specs：studio-v2-device-to-sql-acceptance。
- [2026-10-03-fix-studio-v2-database-workflow](../../../openspec/changes/archive/2026-10-03-fix-studio-v2-database-workflow/validation.md)；specs：datalink-workbench-v2-step4-database, recording-database-setup。

完整requirements/scenarios、測試／source/witness身分、清理與未驗界線見 [final-verification.md](final-verification.md) 與 [handoff-closeout.md](handoff-closeout.md)。
