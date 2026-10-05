## 1. Basic preparation
- [ ] 1.1 先新增 Complete basic preparation with one effective configuration 與 Basic preparation retains schema and capability safety UI 回歸：空白 DB、cancel/stale scope/明確 confirm、unsupported kind/pool、主画面/diagnostics，單 worker Vitest RED。
- [ ] 1.2 [after: 1.1] 落實 Basic 重用安全 schema panel；selected Basic group 可 preview/confirm/recover，不須 Advanced；上述 focused tests GREEN，backend token guards 不放寬。
- [ ] 1.3 [after: 1.2] 落實連線與生效記錄各有唯一設定；移除主線舊 table/write strategy/interval/nonpreview schema，capability disabled kind/pool 與 en/zh-TW 說明；focused Step4／TestPage caller regressions GREEN。
- [ ] 1.4 [after: 1.3] 落實主畫面與診斷分層；主画面保留名稱/值/type+unit/真實 stages/errors，technical IDs/revisions/payload 放 disclosure；DOM focused tests GREEN。

## 2. 回查與整體驗收
- [ ] 2.1 [after: 1.4] scope review→fix→rereview、git diff --check、make check-lines 與 affected tests；verification.md 明列差異對應與測試證據。
- [ ] 2.2 [after: 2.1] 主對話協調後更新 F fresh-ui 走 Basic，串行 simulator 八量測點定時一列／default uint64／同址 A215/B187／disable-restart／legacy migration／DB recovery，UI 截圖與 repo required gates；最后 aggregate exact HEAD review，未做真 PLC/Windows/LAN/長跑明列。
