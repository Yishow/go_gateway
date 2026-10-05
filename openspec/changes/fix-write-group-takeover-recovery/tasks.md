## 1. 接管回歸與修正
- [x] 1.1 新增 Durable legacy writer takeover production-chain regression：遷移 enabled legacy target，接管後 disable/delete/destination change/retire/restart 不回寫；未有效接管仍可 legacy 寫；focused Go 先 RED。
- [x] 1.2 [after: 1.1] 落實持久接管與活動 worker 分離；首次有效接管先保存 connector/tag/group/revision ownership，失敗 fail closed，所有 lifecycle focused Go GREEN。

## 2. 故障 head 恢復
- [x] 2.1 [after: 1.2] 新增 Traceable scoped delivery head resolution RED：missing table/permission、poison row、foreign scope、stale state、duplicate decision、unknown 拒絕與 frozen identity tests。
- [x] 2.2 [after: 2.1] 落實以原 effect 身份修復 head；group attention/read 與原子 audit retry/skip，unknown 拒絕，production chain frozen destination/receipt dedupe tests GREEN。
- [x] 2.3 [after: 2.2] operator UI 先 Vitest RED，再實作安全原因、retry、明確 skip、pending/error/unknown states，單 worker GREEN；現有 group authority 不另建狀態模型。

## 3. 檢查與驗收
- [x] 3.1 [after: 2.3] scope review→fix→rereview、git diff --check、make check-lines、affected suite 通過，verification.md 記錄 RED/GREEN 与 provenance assertions。
- [x] 3.2 [after: 3.1] 向主對話協調時段後串行 repo required gates、simulator 八量測點定時一列、停用重啟、legacy migration 與 DB 故障恢復；aggregate exact HEAD review，未做現場／平台驗收明列限制。
