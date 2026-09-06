## 1. 演算法 Red
- [ ] 1.1 將 acceptance.md 的 counter、gap、wrap、signed、delta、rate、state 範例轉成固定時鐘 fixtures。
- [ ] 1.2 Red：測試時間加權分子/有效時長、min/max、0/null、品質中斷與視窗邊界。
- [ ] 1.3 Red：測試 counter 起始、相等、重置、多次回捲不明、極大值精度、異常跳值與 late input。
- [ ] 1.4 Red：測試 signed/delta/event 去重、缺邊界、跨日、DAG 單位與循環、state unknown 時長。

## 2. Green 與原子保存
- [ ] 2.1 實作確定性 gauge/rate 聚合與有效區段 coverage，禁止取 writer 最後一筆代替輸入。
- [ ] 2.2 實作 counter/signed/delta baseline、epoch、異常隔離與 evidence ledger。
- [ ] 2.3 實作不完整區間、已知小計與明確啟用的估算；實測與估算來源不可混合。
- [ ] 2.4 實作速率積分、state/bitmask/event/text 處理及白名單衍生計算。
- [ ] 2.5 實作 watermark/寬限、有限重算與 calculation_revision，checkpoint/result/outbox 同 transaction。
- [ ] 2.6 加入崩潰於各提交邊界的故障注入測試，重啟不得重複計量或跳過一段。

## 3. Refactor 與整合
- [ ] 3.1 將純數學、狀態機、時間邊界與持久化分離，避免單一巨大服務檔。
- [ ] 3.2 API 結果包含 coverage、estimated、epoch、evidence、calculation_revision 與安全原因。
- [ ] 3.3 與 2、4 聯測多目的地、晚到更正與過期證據；同一來源不重複採集或雙算。
- [ ] 3.4 執行 Go test/vet/lint、精度 property tests、固定時區測試及 strict spec 驗證，留下全部反例證據。
