## 1. 共同入口修正

- [x] 1.1 完成「共同入口與獨立編輯」及 Purpose-led layout with one advanced authority：舊寬表 A 與新 B 預設共用用途及瞬時／累積表單，舊編輯必須明確開啟；先以共同入口回歸案例觀察行為失敗，再修正 BasicRecordingPanel／CircuitProfileBasicPanel／ExistingDeviceRecording 與兩語系文案，驗證 A=47、B=33 及舊編輯草稿 73 各自保留且不 mutation。

## 2. 既有保護與有界交付

- [x] 2.1 完成「沿用伺服器權威」：開始只送審核的新群組，舊未套用草稿不加入，server 拒絕來源接手時不改／停／刪舊群組；保留 Share-only、Advanced 往返、原值、來源收合、操作復原與只準備的零開始保護。以 recording-workflow-unified、compatibility、migration-prepare、migration-timescale、basic 與 summary 既有／新增測試驗證，不把原本通過的保護冒稱新增 RED。
- [x] 2.2 完成「有界驗收與交付」：同步既有操作指南，使用隔離真瀏覽器 fixture 在 1024×720／1600×1000 驗共同入口、A/B 草稿、舊編輯往返、重載、鍵盤與零寫入，明列模擬服務限制；核對前端 tests、tsc、lint、build、git diff --check 與 make check-lines，不操作真設備／資料庫，不納入外案 WIP。
