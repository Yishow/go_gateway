## 1. 先建立失敗回歸
- [x] 1.1 驗證 Source preserving exact mapping conversion：新增 default UI target、executeCast 邊界與 production uint64-to-SQL 回歸，串行 focused Go/Vitest 記錄改前失敗（RED）；預期值 9007199254740993，不只 codec。
- [x] 1.2 驗證 Device scoped live point identity：新增 A.D0=215／B.D0=187、切換 stale event、名稱型別單位與 waiting 回歸，單 worker Vitest 記錄改前失敗（RED）。

## 2. 修正與回查
- [x] 2.1 [after: 1.1] 落實保留型別並拒絕失真；buildDefaultMapping 保留 source type，executeCast 拒絕非法／溢位／分數／非有限／失真，既有映射不自動改寫；上述 focused tests GREEN。
- [x] 2.2 [after: 1.2] 落實依選取設備與 point 身份配對；selected device caller、event、metadata 都隔離，表格可讀名稱／值／型別單位；上述切換 tests GREEN。
- [x] 2.5 [after: 2.1] 補償 review 校正（不改已完成歷史），涵蓋「保留合法縮放與舊步驟順序」及「保留合法縮放與舊步驟順序的確認邊界」：neutral source／明確 scale float64 default、new scale→checked integer cast、saved step order；先 RED 再 GREEN，合法 int16×0.5+10 production default-to-SQL 與 tag 型別一致。
- [ ] 2.3 [after: 2.5, 2.2] 進行 scope review→fix→rereview，git diff --check、make check-lines 與受影響 suite 通過；錯值不當 good 寫 SQL 的證據存 verification.md。
- [ ] 2.4 [after: 2.3] 先通知主對話協調低負載時段，再執行 repo required gates／simulator 同址切換實測，最後 aggregate exact HEAD review；未執行環境驗收明列 blocker，不勾虛假完成。
