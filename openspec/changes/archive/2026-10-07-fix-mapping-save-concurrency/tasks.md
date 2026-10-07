## 1. 修正與回歸

- [x] 1.1 滿足 Concurrent owned mapping saves preserve rule integrity，依「規則協調涵蓋 production mutation 與 save-confirm」與「連結替換是 repository 原子操作」交付同規則 save-confirm/source mutation 協調與 repository 原子連結替換，先保留真 SQL 八列並行/DELETE窗口/rollback 失敗回歸，再證明 pool 1/4、不同 rule、edit/sync/reapply/delete 及 stale source 保護通過。
- [x] 1.2 滿足 Rule queued autosave retains latest drafts and actionable failures，依「前端 rule queue 與 typed failure」交付前端按 rule 排隊、快改保留最後草稿、失敗明確 retry 與安全 status/code，先用 React deferred request 測試重現批次/快改/失敗，再驗證不同 rule 與 uint64/scale 不退步。

## 2. 整合與交接

- [x] 2.1 [after: 1.1, 1.2] 隔離正常 UI 多列 batch 保存/preview→真 SQLite 持久 readback，保留 exact source/binary/UI 證據並清理 owned 資源，不操作現場3222/5173。
- [x] 2.2 [after: 2.1] 完成 review→fix→rereview、scenario coverage、Go/frontend 全 gates、OpenSpec/line checks、本地 commit 與 exact HEAD aggregate review，回報兩樹狀態及現場版本未更新。
