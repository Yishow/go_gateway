# 七案規格與審查對帳

日期：2026-11-02。產品依據原七案規格，不修改驗收條文以配合結果；沒有新增 change 或 archive。

## 規格映射

| Change | Requirements／Scenarios／Examples | 實作與可追溯映射 |
| --- | --- | --- |
| simplify-recording-home-guidance | 12／35／0 | `home-editor-conformance.md/json`：同scope事實、未知不造值、排序／有界總覽、明確新增／返回；原route/owner測試與Run06實載bundle驗證 |
| streamline-recording-task-editor | 17／72／0 | 同上；來源owner、逐點字序／算法、保存／檢查／確認／Start、dirty與canonical；MAIN三協議／兩DB旅程及每個真缺陷的最窄相關測試 |
| simplify-recording-lifecycle-recovery | 5／16／0 | `lifecycle-task-map.md`逐11task：availability、停止／移除、fresh恢復、副本權威清除、forward歷史、同operationunknown、results/repair/logscope；獨立owner review packets與分版本Browser |
| adopt-simple-recording-entry | 8／16／0 | `adoption/conformance.md`、`independent-review.md`：明確意圖優先／invalid failclosed、原owner重用、預設採用／advanced回退、設計與證據分層。交付gate等政策scenario以版本化owner／驗收紀錄人工核，不造產品permission branch |
| preserve-per-point-word-order | 9／20／1 | `calculation-order-conformance.md/json`：inherit/override、codec/protocol、migration/down-grade、預覽與權威；唯一正式Example逐值對應既有golden／direct assertions |
| support-per-member-recording-calculation | 11／25／0 | 同上：UTC對齊、due聯集、NULL/not_due、counter首baseline與新generation、failure/partial/gap/reset。新發現Linux漏7case已中性命名並真執行，不將MAIN5s/15s替代spec15s/60s精確fixture |
| integrate-protocol-aware-acquisition-batching | 10／41／0 | `TASK-4.2-EVIDENCE-INDEX.md`、`task42-map.json`、逐file semantic-map、current15scheduled／6bench，原planner/adapter/fault/Stop/取消/身份直接測試；不以zero-sink壓力當SQL無損 |

表內數量按原 delta spec 的正式 heading 計算，共72 requirements／225 scenarios／1 Example。映射描述實作與被檢視斷言，不等同每個scenario都有新的Browser／硬體實測。政策型scenario的test-scope排除與未提供環境skip分開；執行版本、命令、結果詳見最終驗收記錄。

## 審查範圍與限制

- 使用原repo Spectra verify/review/test-scope 指令與四個review面向（正確性、效率、重用、慣例）。沒有根據目前乾淨worktree假設完整歷史diff已重新審過。
- 5025 最終文書收尾的CLI scope：home/lifecycle/word-order/member缺可信pre-implementation base，狀態 `insufficient`；不猜base，不擴掃無關檔。功能完整性以原spec、current code/tests、已封packet及真驗收證據逐條判斷；**完整歷史hunk歸屬仍未由CLI證明**。
- editor的tracked scope僅4paths，不能代替全部歷史實作；adoption scope101paths有4個preexisting-dirty歸屬限制。獨立review已覆必要主旅程production與權威delta，52測試檔未宣稱全部重新逐行審查；它們有本輪完整FE/Go執行證據。
- batching早先CLI fallback scope673paths屬近似baseline範圍，不能把它們全部歸4.2。完成task14時已用官方CLI明確列19個原授path，最終scope為19paths（不是全dirty tree）。只以原19paths凍結包、逐phase依賴hash、39策略/11race逐diff語意review與正式量測作精確證據；原大量其他lint變更不借此歸屬。
- 性能paired/race採用歷史實測＋逐檔語意審查，原版本結果與current量測分列；微小正常路徑指令成本未重新paired量化。未額外創造zero-drop硬門檻，也不把有drop隱藏成全值交付。
- 最後新增query hook早拒回歸由獨立review定位，保原無Provider拒絕斷言；runtime poll fixture只接收精確已定位點位，原設定寫入/unknown拒絕不放寬。Go隱藏檔名與import格式問題均保留真FAIL／修復收據。

## 完成與發布分開

最後5025完整FE、官方build、同版長表真像素及本機收束皆通過後，已用官方CLI逐項更新原14個pending checkbox；原七案合計88/88 tasks完成。各項皆依實證，而非依工具scope是否完整決定功能已驗；各change仍保留可檢視，不archive。官方artifact對應產品版本，後續純文書提交不重建或偽稱另一次完整程式測試。既有同步授權下，本批遭自動核准中止，依返回要求待使用者再次明確確認；不由本機驗收狀態推定已發布。

| Change | Tasks |
| --- | --- |
| `simplify-recording-lifecycle-recovery` | 11/11 |
| `integrate-protocol-aware-acquisition-batching` | 14/14 |
| `adopt-simple-recording-entry` | 8/8 |
| `streamline-recording-task-editor` | 15/15 |
| `simplify-recording-home-guidance` | 13/13 |
| `preserve-per-point-word-order` | 13/13 |
| `support-per-member-recording-calculation` | 14/14 |

最終 verify：Completeness 88/88，72 requirements／225 scenarios／1 Example 已按既有映射對帳；未發現未實作的必要功能或新的阻擋性缺陷。Correctness 與 Coherence 以具名測試、真旅程、獨立review及本輪窄修驗證支持；上述歷史歸屬／部分diff review限制仍列 WARNING，不改稱完整歷史audit。七案最後scope snapshot check均PASS，狀態insufficient者仍原樣保留。
